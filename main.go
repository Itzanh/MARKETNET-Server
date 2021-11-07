package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

// Basic, static, server settings such as the DB password or the port.
var settings BackendSettings

// Http object for the websocket clients to conenect to.
var upgrader = websocket.Upgrader{}

// Database connection to PostgreSQL.
var db *sql.DB

// List of all the concurrent websocket connections to the server.
var connections []Connection

// Global cron instance
var c *cron.Cron

// List of all cron IDs. Key= Enterprise Id, Value= Array of objects with cron IDs.
var runningCrons map[int32]EnterpriseCronInfo = make(map[int32]EnterpriseCronInfo)

func main() {
	// read settings
	var ok bool
	settings, ok = getBackendSettings()
	if !ok {
		fmt.Println("ERROR READING SETTINGS FILE")
		return
	}

	// connect to PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", settings.Db.Host, settings.Db.Port, settings.Db.User, settings.Db.Password, settings.Db.Dbname)
	db, _ = sql.Open("postgres", psqlInfo) // control error
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// installation
	if !isParameterPresent("--dev-no-upgrade-database") {
		if !installDB() {
			os.Exit(1)
		}
	}

	// initial data
	settingsRecords := getSettingsRecords()
	for i := 0; i < len(settingsRecords); i++ {
		initialData(settingsRecords[i].Id)
	}
	if isParameterPresent("--install-only") {
		return
	}
	if isParameterPresent("--generate-demo-data") {
		for i := 0; i < len(settingsRecords); i++ {
			generateDemoData(settingsRecords[i].Id)
		}
	}

	// add a new enterprise by command line
	if isParameterPresent("--add-enterprise") {
		ok := addEnterpriseFromParameters()
		if ok {
			os.Exit(0)
		} else {
			os.Exit(3)
		}
	}

	// listen to requests
	fmt.Println("Server ready! :D")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", reverse)
	http.HandleFunc("/document", handleDocument)
	http.HandleFunc("/report", generateReport)
	http.HandleFunc("/export", handleExport)
	if isParameterPresent("--saas") {
		http.HandleFunc("/saas", handleEnterprise)
	}

	addHttpHandlerFuncions()
	if settings.Server.TLS.UseTLS {
		go http.ListenAndServeTLS(":"+strconv.Itoa(int(settings.Server.Port)), settings.Server.TLS.CrtPath, settings.Server.TLS.KeyPath, nil)
	} else {
		go http.ListenAndServe(":"+strconv.Itoa(int(settings.Server.Port)), nil)
	}

	// crons
	go cleanDocumentTokens()
	c = cron.New()
	for i := 0; i < len(settingsRecords); i++ {
		var enterpriseId int32 = settingsRecords[i].Id
		enterpriseCronInfo := EnterpriseCronInfo{}
		if settingsRecords[i].Currency != "_" {
			cronId, err := c.AddFunc(settingsRecords[i].CronCurrency, func() {
				updateCurrencyExchange(enterpriseId)
			})
			if err != nil {
				enterpriseCronInfo.CronCurrency = &cronId
			}
		}
		if settingsRecords[i].Ecommerce != "_" {
			e := ECommerce{Enterprise: settingsRecords[i].Id}
			cronId, err := c.AddFunc(settingsRecords[i].CronPrestaShop, e.ecommerceControllerImportFromEcommerce)
			if err != nil {
				enterpriseCronInfo.CronPrestaShop = &cronId
			}
		}
		cronId, err := c.AddFunc(settingsRecords[i].CronClearLabels, func() {
			deleteAllShippingTags(enterpriseId)
		})
		if err != nil {
			enterpriseCronInfo.CronClearLabels = cronId
		}
		runningCrons[enterpriseId] = enterpriseCronInfo
	}
	c.AddFunc(settings.Server.CronClearLogs, clearLogs)
	c.AddFunc("@every 1m", resetMaxRequestsPerEnterprise)
	c.Start()
	c.Run()

	// activation
	go activate()

	// idle wait to prevent the main thread from exiting
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func reverse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Client connected! " + r.RemoteAddr)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// AUTHENTICATION
	ok, userId, permissions, enterpriseId := authentication(ws, r.RemoteAddr)
	if !ok || permissions == nil {
		return
	}
	okFilter := userConnection(userId, r.RemoteAddr, enterpriseId)
	if !okFilter {
		return
	}
	s := getSettingsRecordById(enterpriseId)
	if len(getConnections(enterpriseId)) >= int(s.MaxConnections) {
		ws.Close()
		return
	}
	// END AUTHENTICATION
	c := Connection{Address: r.RemoteAddr, User: userId, ws: ws, enterprise: enterpriseId}
	c.addConnection()

	for {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
			c.deleteConnection()
			userDisconnected(userId)
			return
		}

		msg := string(message)
		separatorIndex := strings.Index(msg, "$")
		if separatorIndex < 0 {
			break
		}

		command := msg[0:separatorIndex]
		commandSeparatorIndex := strings.Index(command, ":")
		if commandSeparatorIndex < 0 {
			break
		}

		commandProcessor(command[0:commandSeparatorIndex], command[commandSeparatorIndex+1:], message[separatorIndex+1:], mt, ws, *permissions, userId, enterpriseId)
	}
}

// Ok, user id, user permissions, enterprise id
func authentication(ws *websocket.Conn, remoteAddr string) (bool, int32, *Permissions, int32) {
	var userId int32
	var enterpriseId int32
	// AUTHENTICATION
	var i int16 = 0
	for ; i < settings.Server.MaxLoginAttemps; i++ {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
			return false, 0, nil, 0
		}

		// Remote the port from the address
		if strings.Contains(remoteAddr, ":") {
			remoteAddr = remoteAddr[:strings.Index(remoteAddr, ":")]
		}

		// Attempt login in DB
		var userLogin UserLogin
		json.Unmarshal(message, &userLogin)
		result := UserLoginResult{}
		if len(userLogin.Token) > 0 {
			t := LoginToken{Name: userLogin.Token, IpAddress: remoteAddr}
			result.Ok, result.Permissions, userId, enterpriseId = t.checkLoginToken()
			if result.Ok && userId >= 0 {
				result.Language = getUserRow(userId).Language
			}
		} else {
			result, userId, enterpriseId = userLogin.login(remoteAddr)
		}

		// Google Authenticator
		user := getUserRow(userId)
		if len(userLogin.Token) == 0 && user.UsesGoogleAuthenticator {
			data, _ := json.Marshal(UserLoginResult{Ok: true, GoogleAuthenticator: true})
			ws.WriteMessage(mt, data)

			// Receive message
			_, message, err := ws.ReadMessage()
			if err != nil {
				return false, 0, nil, 0
			}

			ok := authenticateUserInGoogleAuthenticator(userId, enterpriseId, string(message))
			if !ok {
				data, _ := json.Marshal(UserLoginResult{GoogleAuthenticator: true})
				ws.WriteMessage(mt, data)
				return false, 0, nil, 0
			} else {
				// Return result to client (Ok + Token)
				data, _ := json.Marshal(result)
				ws.WriteMessage(mt, data)
				if result.Ok {
					return true, userId, result.Permissions, enterpriseId
				}
			}
		} else {

			// Return result to client (Ok + Token)
			data, _ := json.Marshal(result)
			ws.WriteMessage(mt, data)
			if result.Ok {
				return true, userId, result.Permissions, enterpriseId
			}

		}
	}
	// END AUTHENTICATION
	return false, 0, nil, 0
}

func commandProcessor(instruction string, command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions, userId int32, enterpriseId int32) {
	switch instruction {
	case "GET":
		instructionGet(command, string(message), mt, ws, permissions, enterpriseId)
	case "INSERT":
		instructionInsert(command, message, mt, ws, permissions, userId, enterpriseId)
	case "UPDATE":
		instructionUpdate(command, message, mt, ws, permissions, enterpriseId)
	case "DELETE":
		instructionDelete(command, string(message), mt, ws, permissions, enterpriseId)
	case "NAME":
		instructionName(command, string(message), mt, ws, enterpriseId)
	case "GETNAME":
		instructionGetName(command, string(message), mt, ws, enterpriseId)
	case "DEFAULTS":
		instructionDefaults(command, string(message), mt, ws, permissions, enterpriseId)
	case "LOCATE":
		instructionLocate(command, string(message), mt, ws, permissions, enterpriseId)
	case "ACTION":
		instructionAction(command, string(message), mt, ws, permissions, userId, enterpriseId)
	case "SEARCH":
		instructionSearch(command, string(message), mt, ws, permissions, enterpriseId)
	}
}

type PaginationQuery struct {
	Enterprise int32
	Offset     int64 `json:"offset"`
	Limit      int64 `json:"limit"`
}

func (q *PaginationQuery) isValid() bool {
	return !(q.Offset < 0 || q.Limit <= 0)
}

type OperationResult struct {
	Id        int64  `json:"id"`
	Code      uint16 `json:"code"`
	ExtraData string `json:"extraData"`
}

func instructionGet(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	var found bool = true
	var data []byte

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			var paginationQuery PaginationQuery
			json.Unmarshal([]byte(message), &paginationQuery)
			paginationQuery.Enterprise = enterpriseId
			data, _ = json.Marshal(paginationQuery.getAddresses())
		case "PRODUCT":
			data, _ = json.Marshal(getProduct(enterpriseId))
		case "PRODUCT_FAMILY":
			data, _ = json.Marshal(getProductFamilies(enterpriseId))
		case "BILLING_SERIE":
			data, _ = json.Marshal(getBillingSeries(enterpriseId))
		case "CURRENCY":
			data, _ = json.Marshal(getCurrencies(enterpriseId))
		case "PAYMENT_METHOD":
			data, _ = json.Marshal(getPaymentMethods(enterpriseId))
		case "LANGUAGE":
			data, _ = json.Marshal(getLanguages(enterpriseId))
		case "COUNTRY":
			data, _ = json.Marshal(getCountries(enterpriseId))
		case "STATE":
			data, _ = json.Marshal(getStates(enterpriseId))
		case "CUSTOMER":
			var paginationQuery PaginationQuery
			json.Unmarshal([]byte(message), &paginationQuery)
			paginationQuery.Enterprise = enterpriseId
			data, _ = json.Marshal(paginationQuery.getCustomers())
		case "COLOR":
			data, _ = json.Marshal(getColor(enterpriseId))
		case "PACKAGES":
			data, _ = json.Marshal(getPackages(enterpriseId))
		case "INCOTERMS":
			data, _ = json.Marshal(getIncoterm(enterpriseId))
		case "CARRIERS":
			data, _ = json.Marshal(getCariers(enterpriseId))
		case "SUPPLIERS":
			data, _ = json.Marshal(getSuppliers(enterpriseId))
		case "DOCUMENT_CONTAINER":
			data, _ = json.Marshal(getDocumentContainer(enterpriseId))
		case "DOCUMENTS":
			if message == "" {
				data, _ = json.Marshal(getDocuments(enterpriseId))
			} else {
				var document Document
				json.Unmarshal([]byte(message), &document)
				data, _ = json.Marshal(document.getDocumentsRelations(enterpriseId))
			}
		default:
			found = false
		}

		if found {
			ws.WriteMessage(mt, data)
			return
		} else {
			found = true
		}
	} // Masters

	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		data, _ = json.Marshal(paginationQuery.getSalesOrder(enterpriseId))
	case "SALES_ORDER_PREPARATION":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getSalesOrderPreparation(enterpriseId))
	case "SALES_ORDER_AWAITING_SHIPPING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getSalesOrderAwaitingShipping(enterpriseId))
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getWarehouses(enterpriseId))
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.Enterprise = enterpriseId
		data, _ = json.Marshal(paginationQuery.getSalesInvoices())
	case "MANUFACTURING_ORDER_TYPE":
		if (!permissions.Manufacturing) && (!permissions.Masters) {
			return
		}
		data, _ = json.Marshal(getManufacturingOrderType(enterpriseId))
	case "WAREHOUSE_MOVEMENTS":
		if !permissions.Warehouse {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.Enterprise = enterpriseId
		data, _ = json.Marshal(paginationQuery.getWarehouseMovement())
	case "WAREHOUSE_WAREHOUSE_MOVEMENTS":
		if !permissions.Warehouse {
			return
		}
		var warehouseMovementByWarehouse WarehouseMovementByWarehouse
		json.Unmarshal([]byte(message), &warehouseMovementByWarehouse)
		warehouseMovementByWarehouse.Enterprise = enterpriseId
		data, _ = json.Marshal(warehouseMovementByWarehouse.getWarehouseMovementByWarehouse())
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.Enterprise = enterpriseId
		data, _ = json.Marshal(paginationQuery.getSalesDeliveryNotes())
	case "SHIPPINGS":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getShippings(enterpriseId))
	case "SHIPPING_NOT_COLLECTED":
		data, _ = json.Marshal(getShippingsPendingCollected(enterpriseId))
	case "USERS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getUser(enterpriseId))
	case "GROUPS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getGroup(enterpriseId))
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrder(enterpriseId))
	case "NEEDS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getNeeds(enterpriseId))
	case "PURCHASE_DELIVERY_NOTES":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNotes(enterpriseId))
	case "PURCHASE_INVOICES":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoices(enterpriseId))
	case "SETTINGS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getSettingsRecordById(enterpriseId))
	case "CLIENT_SETTINGS":
		data, _ = json.Marshal(getSettingsRecordById(enterpriseId).censorSettings())
	case "PS_ZONES":
		if !permissions.PrestaShop {
			return
		}
		data, _ = json.Marshal(getPSZones(enterpriseId))
	case "TABLES":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getTableAndFieldInfo())
	case "CONNECTIONS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getConnections(enterpriseId))
	case "JOURNALS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getJournals(enterpriseId))
	case "ACCOUNTS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getAccounts(enterpriseId))
	case "ACCOUNTING_MOVEMENTS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getAccountingMovement(enterpriseId))
	case "CONFIG_ACCOUNTS_VAT":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getConfigAccountsVat(enterpriseId))
	case "PENDING_COLLECTION_OPERATIONS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getPendingColletionOperations(enterpriseId))
	case "SEARCH_COLLECTION_OPERATIONS":
		if !permissions.Accounting {
			return
		}
		var collectionOperationSearch CollectionOperationPaymentTransactionSearch
		json.Unmarshal([]byte(message), &collectionOperationSearch)
		data, _ = json.Marshal(searchCollectionOperations(collectionOperationSearch, enterpriseId))
	case "SEARCH_PAYMENT_TRANSACTIONS":
		if !permissions.Accounting {
			return
		}
		var collectionOperationSearch CollectionOperationPaymentTransactionSearch
		json.Unmarshal([]byte(message), &collectionOperationSearch)
		data, _ = json.Marshal(searchPaymentTransactions(collectionOperationSearch, enterpriseId))
	case "PENDING_PAYMENT_TRANSACTIONS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getPendingPaymentTransaction(enterpriseId))
	case "COUNTRIES_SALES_ORDERS_AMOUNT":
		var countriesSaleOrdersQuery CountriesSaleOrdersQuery
		json.Unmarshal([]byte(message), &countriesSaleOrdersQuery)
		data, _ = json.Marshal(countriesSaleOrdersQuery.countriesSaleOrdersAmount(enterpriseId))
	case "MANUFACTURING_ORDER_CREATED_MANUFACTURES_DAILY":
		data, _ = json.Marshal(manufacturingOrderCreatedManufacturedDaily(enterpriseId))
	case "DAILY_SHIPPING_QUANTITY":
		data, _ = json.Marshal(dailyShippingQuantity(enterpriseId))
	case "SHIPPING_BY_CARRIERS":
		data, _ = json.Marshal(shippingByCarriers(enterpriseId))
	case "API_KEYS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getApiKeys(enterpriseId))
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingPaginationQuery ManufacturingPaginationQuery
		json.Unmarshal([]byte(message), &manufacturingPaginationQuery)
		data, _ = json.Marshal(manufacturingPaginationQuery.getManufacturingOrder(enterpriseId))
	case "CONNECTION_LOG":
		if !permissions.Admin {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.Enterprise = enterpriseId
		data, _ = json.Marshal(paginationQuery.getConnectionLogs())
	case "CONNECTION_FILTERS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getConnectionFilters(enterpriseId))
	case "REPORT_TEMPLATE":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getReportTemplates(enterpriseId))
	default:
		found = false
	}

	if found {
		ws.WriteMessage(mt, data)
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	found = true
	switch command {
	case "MONTHLY_SALES_AMOUNT":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(monthlySalesAmount(year, enterpriseId))
	case "MONTHLY_SALES_QUANTITY":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(monthlySalesQuantity(year, enterpriseId))
	case "SALES_OF_A_PRODUCT_QUANTITY":
		data, _ = json.Marshal(salesOfAProductQuantity(int32(id), enterpriseId))
	case "SALES_OF_A_PRODUCT_AMOUNT":
		data, _ = json.Marshal(salesOfAProductAmount(int32(id), enterpriseId))
	case "DAYS_OF_SERVICE_SALE_ORDERS":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(daysOfServiceSaleOrders(year, enterpriseId))
	case "DAYS_OF_SERVICE_PURCHASE_ORDERS":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(daysOfServicePurchaseOrders(year, enterpriseId))
	case "PURCHASE_ORDERS_BY_MONTH_AMOUNT":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(purchaseOrdersByMonthAmount(year, enterpriseId))
	case "PAYMENT_METHODS_SALE_ORDERS_AMOUNT":
		var year *int16
		if id > 0 {
			aux := int16(id)
			year = &aux
		}
		data, _ = json.Marshal(paymentMethodsSaleOrdersAmount(year, enterpriseId))
	default:
		found = false
	}
	if found {
		ws.WriteMessage(mt, data)
		return
	}

	if id <= 0 {
		return
	}
	switch command {
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderDetail(int64(id), enterpriseId))
	case "STOCK":
		data, _ = json.Marshal(getStock(int32(id), enterpriseId))
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderDiscounts(int32(id), enterpriseId))
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceDetail(int64(id), enterpriseId))
	case "SALES_ORDER_PACKAGING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getPackaging(int64(id), enterpriseId))
	case "SALES_DELIVERY_NOTES_DETAILS":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementBySalesDeliveryNote(int64(id), enterpriseId))
	case "SHIPPING_PACKAGING":
		data, _ = json.Marshal(getPackagingByShipping(int64(id), enterpriseId))
	case "GET_USER_GROUPS":
		data, _ = json.Marshal(getUserGroups(int32(id), enterpriseId))
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderDetail(int64(id), enterpriseId))
	case "PURCHASE_DELIVERY_NOTES_DETAILS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementByPurchaseDeliveryNote(int64(id), enterpriseId))
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceDetail(int64(id), enterpriseId))
	case "PRODUCT_SALES_ORDER_PENDING":
		data, _ = json.Marshal(getProductSalesOrderDetailsPending(int32(id), enterpriseId))
	case "PRODUCT_PURCHASE_ORDER_PENDING":
		data, _ = json.Marshal(getProductPurchaseOrderDetailsPending(int32(id), enterpriseId))
	case "PRODUCT_SALES_ORDER":
		data, _ = json.Marshal(getProductSalesOrderDetails(int32(id), enterpriseId))
	case "PRODUCT_PURCHASE_ORDER":
		data, _ = json.Marshal(getProductPurchaseOrderDetails(int32(id), enterpriseId))
	case "PRODUCT_WAREHOUSE_MOVEMENT":
		data, _ = json.Marshal(getProductWarehouseMovement(int32(id), enterpriseId))
	case "SALES_ORDER_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderRow(int64(id)))
	case "SALES_INVOICE_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceRow(int64(id)))
	case "PURCHASE_ORDER_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderRow(int64(id)))
	case "PURCHASE_INVOICE_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceRow(int64(id)))
	case "PRODUCT_IMAGE":
		data, _ = json.Marshal(getProductImages(int32(id), enterpriseId))
	case "CUSTOMER_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getCustomerRow(int32(id)))
	case "SUPPLIER_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getSupplierRow(int32(id)))
	case "PALLETS":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getSalesOrderPallets(int64(id), enterpriseId))
	case "CUSTOMER_ADDRESSES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getCustomerAddresses(int32(id), enterpriseId))
	case "CUSTOMER_SALE_ORDERS":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getCustomerSaleOrders(int32(id), enterpriseId))
	case "SALES_DELIVERY_NOTE_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesDeliveryNoteRow(int64(id)))
	case "PURCHASE_DELIVERY_NOTE_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNoteRow(int64(id)))
	case "ACCOUNTING_MOVEMENT_DETAILS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getAccountingMovementDetail(int64(id), enterpriseId))
	case "ACCOUNTING_MOVEMENT_SALE_INVOICES":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getAccountingMovementSaleInvoices(int64(id)))
	case "ACCOUNTING_MOVEMENT_COLLECTION_OPERATION":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getColletionOperations(int64(id), enterpriseId))
	case "COLLECTION_OPERATION_CHARGES":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getCharges(int32(id), enterpriseId))
	case "ACCOUNTING_MOVEMENT_PAYMENT_TRANSACTIONS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getPaymentTransactions(int64(id), enterpriseId))
	case "PAYMENT_TRANSACTION_PAYMENTS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getPayments(int32(id), enterpriseId))
	case "ACCOUNTING_MOVEMENT_PURCHASE_INVOICES":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(getAccountingMovementPurchaseInvoices(int64(id)))
	case "SUPPLIER_ADDRESSES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getSupplierAddresses(int32(id), enterpriseId))
	case "SUPPLIER_PURCHASE_ORDERS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getSupplierPurchaseOrders(int32(id), enterpriseId))
	case "SHIPPING_TAGS":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getShippingTags(int64(id), enterpriseId))
	case "SALES_ORDER_DETAILS_FROM_PURCHASE_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderDetailsFromPurchaseOrderDetail(int64(id), enterpriseId))
	case "PURCHASES_ORDER_DETAILS_FROM_SALE_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getPurchasesOrderDetailsFromSaleOrderDetail(int32(id), enterpriseId))
	case "CONNECTION_FILTER_USERS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getConnectionFilterUser(int32(id), enterpriseId))
	case "ADDRESS_ROW":
		if !permissions.Masters {
			return
		}
		address := getAddressRow(int32(id))
		if address.enterprise != enterpriseId {
			return
		}
		data, _ = json.Marshal(address)
	case "ACCOUNTING_MOVEMENT_ROW":
		if !permissions.Accounting {
			return
		}
		address := getAccountingMovementRow(int64(id))
		if address.enterprise != enterpriseId {
			return
		}
		data, _ = json.Marshal(address)
	}
	ws.WriteMessage(mt, data)
}

func instructionInsert(command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions, userId int32, enterpriseId int32) {
	var ok bool

	if permissions.Masters {
		var found bool
		var operationResult OperationResult
		switch command {
		case "ADDRESS":
			var address Address
			json.Unmarshal(message, &address)
			address.enterprise = enterpriseId
			operationResult = address.insertAddress()
			found = true
		case "CUSTOMER":
			var customer Customer
			json.Unmarshal(message, &customer)
			customer.enterprise = enterpriseId
			operationResult = customer.insertCustomer()
			found = true
		case "SUPPLIER":
			var supplier Supplier
			json.Unmarshal(message, &supplier)
			supplier.enterprise = enterpriseId
			operationResult = supplier.insertSupplier()
			found = true
		}

		if found {
			data, _ := json.Marshal(operationResult)
			ws.WriteMessage(mt, data)
			return
		}
	}

	if permissions.Masters {
		switch command {
		case "BILLING_SERIE":
			var serie BillingSerie
			json.Unmarshal(message, &serie)
			serie.enterprise = enterpriseId
			ok = serie.insertBillingSerie()
		case "CURRENCY":
			var currency Currency
			json.Unmarshal(message, &currency)
			currency.enterprise = enterpriseId
			ok = currency.insertCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			json.Unmarshal(message, &paymentMethod)
			paymentMethod.enterprise = enterpriseId
			ok = paymentMethod.insertPaymentMethod()
		case "LANGUAGE":
			var language Language
			json.Unmarshal(message, &language)
			language.enterprise = enterpriseId
			ok = language.insertLanguage()
		case "COUNTRY":
			var country Country
			json.Unmarshal(message, &country)
			country.enterprise = enterpriseId
			ok = country.insertCountry()
		case "STATE":
			var state State
			json.Unmarshal(message, &state)
			state.enterprise = enterpriseId
			ok = state.insertState()
		case "PRODUCT":
			var product Product
			json.Unmarshal(message, &product)
			product.enterprise = enterpriseId
			ok = product.insertProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			json.Unmarshal(message, &productFamily)
			productFamily.enterprise = enterpriseId
			ok = productFamily.insertProductFamily()
		case "COLOR":
			var color Color
			json.Unmarshal(message, &color)
			color.enterprise = enterpriseId
			ok = color.insertColor()
		case "PACKAGES":
			var packages Packages
			json.Unmarshal(message, &packages)
			packages.enterprise = enterpriseId
			ok = packages.insertPackage()
		case "INCOTERM":
			var incoterm Incoterm
			json.Unmarshal(message, &incoterm)
			incoterm.enterprise = enterpriseId
			ok = incoterm.insertIncoterm()
		case "CARRIER":
			var carrier Carrier
			json.Unmarshal(message, &carrier)
			carrier.enterprise = enterpriseId
			ok = carrier.insertCarrier()
		case "SHIPPING":
			var shipping Shipping
			json.Unmarshal(message, &shipping)
			shipping.enterprise = enterpriseId
			ok, _ = shipping.insertShipping()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			json.Unmarshal(message, &documentContainer)
			documentContainer.enterprise = enterpriseId
			ok = documentContainer.insertDocumentContainer()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			json.Unmarshal(message, &productImage)
			ok = productImage.insertProductImage(enterpriseId)
		}
	} // Masters

	var returnData []byte
	var found bool = true
	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal([]byte(message), &saleOrder)
		saleOrder.enterprise = enterpriseId
		ok, orderId := saleOrder.insertSalesOrder()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			order := getSalesOrderRow(orderId)
			returnData, _ = json.Marshal(order)
		}
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var saleInvoice SalesInvoice
		json.Unmarshal(message, &saleInvoice)
		saleInvoice.enterprise = enterpriseId
		ok, invoiceId := saleInvoice.insertSalesInvoice()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			invoice := getSalesInvoiceRow(invoiceId)
			returnData, _ = json.Marshal(invoice)
		}
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		json.Unmarshal(message, &salesDeliveryNote)
		salesDeliveryNote.enterprise = enterpriseId
		ok, nodeId := salesDeliveryNote.insertSalesDeliveryNotes()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			note := getSalesDeliveryNoteRow(nodeId)
			returnData, _ = json.Marshal(note)
		}
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		json.Unmarshal(message, &purchaseOrder)
		purchaseOrder.enterprise = enterpriseId
		ok, orderId := purchaseOrder.insertPurchaseOrder()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			order := getPurchaseOrderRow(orderId)
			returnData, _ = json.Marshal(order)
		}
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoice PurchaseInvoice
		json.Unmarshal(message, &purchaseInvoice)
		purchaseInvoice.enterprise = enterpriseId
		ok, invoiceId := purchaseInvoice.insertPurchaseInvoice()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			invoice := getPurchaseInvoiceRow(invoiceId)
			returnData, _ = json.Marshal(invoice)
		}
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		json.Unmarshal(message, &purchaseDeliveryNote)
		purchaseDeliveryNote.enterprise = enterpriseId
		ok, noteId := purchaseDeliveryNote.insertPurchaseDeliveryNotes()
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			note := getPurchaseDeliveryNoteRow(noteId)
			returnData, _ = json.Marshal(note)
		}
	default:
		found = false
	}
	if found {
		ws.WriteMessage(mt, returnData)
		return
	}

	switch command {
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		warehouse.enterprise = enterpriseId
		ok = warehouse.insertWarehouse()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		json.Unmarshal(message, &saleOrderDetail)
		saleOrderDetail.enterprise = enterpriseId
		ok = saleOrderDetail.insertSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		json.Unmarshal(message, &saleOrderDiscount)
		saleOrderDiscount.enterprise = enterpriseId
		ok = saleOrderDiscount.insertSalesOrderDiscount()
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		json.Unmarshal(message, &salesInvoiceDetail)
		salesInvoiceDetail.enterprise = enterpriseId
		ok = salesInvoiceDetail.insertSalesInvoiceDetail(true)
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(message, &manufacturingOrder)
		manufacturingOrder.UserCreated = userId
		manufacturingOrder.enterprise = enterpriseId
		ok = manufacturingOrder.insertManufacturingOrder()
	case "SALES_ORDER_PACKAGING":
		if !permissions.Preparation {
			return
		}
		var packaging Packaging
		json.Unmarshal(message, &packaging)
		packaging.enterprise = enterpriseId
		ok = packaging.insertPackaging()
	case "SALES_ORDER_DETAIL_PACKAGED":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal(message, &salesOrderDetailPackaged)
		salesOrderDetailPackaged.enterprise = enterpriseId
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackaged()
	case "SALES_ORDER_DETAIL_PACKAGED_EAN13":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackagedEAN13
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackagedEAN13(enterpriseId)
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		json.Unmarshal(message, &warehouseMovement)
		warehouseMovement.enterprise = enterpriseId
		ok = warehouseMovement.insertWarehouseMovement()
	case "USER":
		if !permissions.Admin {
			return
		}
		var userInsert UserInsert
		json.Unmarshal(message, &userInsert)
		ok = userInsert.insertUser(enterpriseId)
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		json.Unmarshal(message, &group)
		group.enterprise = enterpriseId
		ok = group.insertGroup()
	case "USER_GROUP":
		if !permissions.Admin {
			return
		}
		var userGroup UserGroup
		json.Unmarshal(message, &userGroup)
		ok = userGroup.insertUserGroup()
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		purchaseOrderDetail.enterprise = enterpriseId
		ok, _ = purchaseOrderDetail.insertPurchaseOrderDetail(true)
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		json.Unmarshal(message, &purchaseInvoiceDetail)
		purchaseInvoiceDetail.enterprise = enterpriseId
		ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(true)
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		json.Unmarshal(message, &pallet)
		pallet.enterprise = enterpriseId
		ok = pallet.insertPallet()
	case "JOURNAL":
		if !permissions.Accounting {
			return
		}
		var journal Journal
		json.Unmarshal(message, &journal)
		journal.enterprise = enterpriseId
		ok = journal.insertJournal()
	case "ACCOUNT":
		if !permissions.Accounting {
			return
		}
		var account Account
		json.Unmarshal(message, &account)
		account.enterprise = enterpriseId
		ok = account.insertAccount()
	case "ACCOUNTING_MOVEMENT":
		if !permissions.Accounting {
			return
		}
		var accountingMovement AccountingMovement
		json.Unmarshal(message, &accountingMovement)
		accountingMovement.enterprise = enterpriseId
		ok = accountingMovement.insertAccountingMovement()
	case "ACCOUNTING_MOVEMENT_DETAIL":
		if !permissions.Accounting {
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		json.Unmarshal(message, &accountingMovementDetail)
		accountingMovementDetail.enterprise = enterpriseId
		ok = accountingMovementDetail.insertAccountingMovementDetail()
	case "CONFIG_ACCOUNTS_VAT":
		if !permissions.Admin {
			return
		}
		var configAccountsVat ConfigAccountsVat
		json.Unmarshal(message, &configAccountsVat)
		configAccountsVat.enterprise = enterpriseId
		ok = configAccountsVat.insertConfigAccountsVat()
	case "CHARGES":
		if !permissions.Accounting {
			return
		}
		var charges Charges
		json.Unmarshal(message, &charges)
		charges.enterprise = enterpriseId
		ok = charges.insertCharges()
	case "PAYMENT":
		if !permissions.Accounting {
			return
		}
		var payment Payment
		json.Unmarshal(message, &payment)
		payment.enterprise = enterpriseId
		ok = payment.insertPayment()
	case "API_KEYS":
		if !permissions.Admin {
			return
		}
		var apiKey ApiKey
		json.Unmarshal(message, &apiKey)
		apiKey.UserCreated = userId
		apiKey.enterprise = enterpriseId
		ok = apiKey.insertApiKey()
	case "CONNECTION_FILTER":
		if !permissions.Admin {
			return
		}
		var filter ConnectionFilter
		json.Unmarshal(message, &filter)
		filter.enterprise = enterpriseId
		ok = filter.insertConnectionFilter()
	case "CONNECTION_FILTER_USER":
		if !permissions.Admin {
			return
		}
		var filter ConnectionFilterUser
		json.Unmarshal(message, &filter)
		ok = filter.insertConnectionFilterUser(enterpriseId)
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionUpdate(command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	var ok bool

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			var address Address
			json.Unmarshal(message, &address)
			address.enterprise = enterpriseId
			ok = address.updateAddress()
		case "BILLING_SERIE":
			var serie BillingSerie
			json.Unmarshal(message, &serie)
			serie.enterprise = enterpriseId
			ok = serie.updateBillingSerie()
		case "CURRENCY":
			var currency Currency
			json.Unmarshal(message, &currency)
			currency.enterprise = enterpriseId
			ok = currency.updateCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			json.Unmarshal(message, &paymentMethod)
			paymentMethod.enterprise = enterpriseId
			ok = paymentMethod.updatePaymentMethod()
		case "LANGUAGE":
			var language Language
			json.Unmarshal(message, &language)
			language.enterprise = enterpriseId
			ok = language.updateLanguage()
		case "COUNTRY":
			var country Country
			json.Unmarshal(message, &country)
			country.enterprise = enterpriseId
			ok = country.updateCountry()
		case "STATE":
			var state State
			json.Unmarshal(message, &state)
			state.enterprise = enterpriseId
			ok = state.updateState()
		case "CUSTOMER":
			var customer Customer
			json.Unmarshal(message, &customer)
			customer.enterprise = enterpriseId
			ok = customer.updateCustomer()
		case "PRODUCT":
			var product Product
			json.Unmarshal(message, &product)
			product.enterprise = enterpriseId
			ok = product.updateProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			json.Unmarshal(message, &productFamily)
			productFamily.enterprise = enterpriseId
			ok = productFamily.updateProductFamily()
		case "COLOR":
			var color Color
			json.Unmarshal(message, &color)
			color.enterprise = enterpriseId
			ok = color.updateColor()
		case "PACKAGES":
			var packages Packages
			json.Unmarshal(message, &packages)
			packages.enterprise = enterpriseId
			ok = packages.updatePackage()
		case "INCOTERM":
			var incoterm Incoterm
			json.Unmarshal(message, &incoterm)
			incoterm.enterprise = enterpriseId
			ok = incoterm.updateIncoterm()
		case "CARRIER":
			var carrier Carrier
			json.Unmarshal(message, &carrier)
			carrier.enterprise = enterpriseId
			ok = carrier.updateCarrier()
		case "SUPPLIER":
			var supplier Supplier
			json.Unmarshal(message, &supplier)
			supplier.enterprise = enterpriseId
			ok = supplier.updateSupplier()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			json.Unmarshal(message, &documentContainer)
			documentContainer.enterprise = enterpriseId
			ok = documentContainer.updateDocumentContainer()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			json.Unmarshal(message, &productImage)
			ok = productImage.updateProductImage(enterpriseId)
		}
	}

	switch command {
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		warehouse.enterprise = enterpriseId
		ok = warehouse.updateWarehouse()
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal(message, &saleOrder)
		saleOrder.enterprise = enterpriseId
		ok = saleOrder.updateSalesOrder()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(message, &salesOrderDetail)
		salesOrderDetail.enterprise = enterpriseId
		ok = salesOrderDetail.updateSalesOrderDetail()
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.updateManufacturingOrderType()
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		var shipping Shipping
		json.Unmarshal(message, &shipping)
		shipping.enterprise = enterpriseId
		ok = shipping.updateShipping()
	case "USER":
		if !permissions.Admin {
			return
		}
		var user User
		json.Unmarshal(message, &user)
		user.enterprise = enterpriseId
		ok = user.updateUser()
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		json.Unmarshal(message, &group)
		group.enterprise = enterpriseId
		ok = group.updateGroup()
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		json.Unmarshal(message, &purchaseOrder)
		purchaseOrder.enterprise = enterpriseId
		ok = purchaseOrder.updatePurchaseOrder()
	case "SETTINGS":
		var settings Settings
		json.Unmarshal(message, &settings)
		settings.Id = enterpriseId
		ok = settings.updateSettingsRecord()
	case "PS_ZONES":
		if !permissions.PrestaShop {
			return
		}
		var zone PSZoneWeb
		json.Unmarshal(message, &zone)
		zone.enterprise = enterpriseId
		ok = zone.updatePSZoneWeb()
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		json.Unmarshal(message, &pallet)
		pallet.enterprise = enterpriseId
		ok = pallet.updatePallet()
	case "JOURNAL":
		if !permissions.Accounting {
			return
		}
		var journal Journal
		json.Unmarshal(message, &journal)
		journal.enterprise = enterpriseId
		ok = journal.updateJournal()
	case "ACCOUNT":
		if !permissions.Accounting {
			return
		}
		var account Account
		json.Unmarshal(message, &account)
		account.enterprise = enterpriseId
		ok = account.updateAccount()
	case "CONNECTION_FILTER":
		if !permissions.Accounting {
			return
		}
		var filter ConnectionFilter
		json.Unmarshal(message, &filter)
		filter.enterprise = enterpriseId
		ok = filter.updateConnectionFilter()
	case "REPORT_TEMPLATE":
		if !permissions.Admin {
			return
		}
		var template ReportTemplate
		json.Unmarshal(message, &template)
		template.enterprise = enterpriseId
		ok = template.updateReportTemplate()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionDelete(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	// attempt to delete from resources with alpha key, if the resource if not found,
	// parse the input as number and attemp to delete resource with numeric key
	var found bool = true
	var ok bool
	// ALPHA
	switch command {
	case "BILLING_SERIE":
		if !permissions.Masters {
			return
		}
		var serie BillingSerie
		serie.Id = message
		serie.enterprise = enterpriseId
		ok = serie.deleteBillingSerie()
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		warehouse.Id = message
		warehouse.enterprise = enterpriseId
		ok = warehouse.deleteWarehouse()
	case "USER_GROUP":
		var userGroup UserGroup
		json.Unmarshal([]byte(message), &userGroup)
		ok = userGroup.deleteUserGroup()
	case "CONFIG_ACCOUNTS_VAT":
		if !permissions.Admin {
			return
		}
		id, err := strconv.ParseFloat(message, 32)
		if err != nil || id < 0 {
			return
		}
		var configAccountsVat ConfigAccountsVat
		configAccountsVat.VatPercent = float64(id)
		configAccountsVat.enterprise = enterpriseId
		ok = configAccountsVat.deleteConfigAccountsVat()
	case "CONNECTION_FILTER_USER":
		if !permissions.Admin {
			return
		}
		var filter ConnectionFilterUser
		json.Unmarshal([]byte(message), &filter)
		ok = filter.deleteConnectionFilterUser(enterpriseId)
	default:
		found = false
	}

	if found {
		data, _ := json.Marshal(ok)
		ws.WriteMessage(mt, data)
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil || id <= 0 {
		return
	}

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			var address Address
			address.Id = int32(id)
			address.enterprise = enterpriseId
			ok = address.deleteAddress()
		case "CURRENCY":
			var currency Currency
			currency.Id = int32(id)
			currency.enterprise = enterpriseId
			ok = currency.deleteCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			paymentMethod.Id = int32(id)
			paymentMethod.enterprise = enterpriseId
			ok = paymentMethod.deletePaymentMethod()
		case "LANGUAGE":
			var language Language
			language.Id = int32(id)
			language.enterprise = enterpriseId
			ok = language.deleteLanguage()
		case "COUNTRY":
			var country Country
			country.Id = int32(id)
			country.enterprise = enterpriseId
			ok = country.deleteCountry()
		case "STATE":
			var state State
			state.Id = int32(id)
			state.enterprise = enterpriseId
			ok = state.deleteState()
		case "CUSTOMER":
			var customer Customer
			customer.Id = int32(id)
			customer.enterprise = enterpriseId
			ok = customer.deleteCustomer()
		case "PRODUCT":
			var product Product
			product.Id = int32(id)
			product.enterprise = enterpriseId
			ok = product.deleteProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			productFamily.Id = int32(id)
			productFamily.enterprise = enterpriseId
			ok = productFamily.deleteProductFamily()
		case "COLOR":
			var color Color
			color.Id = int32(id)
			color.enterprise = enterpriseId
			ok = color.deleteColor()
		case "PACKAGES":
			var packages Packages
			packages.Id = int32(id)
			ok = packages.deletePackage()
		case "INCOTERM":
			var incoterm Incoterm
			incoterm.Id = int32(id)
			incoterm.enterprise = enterpriseId
			ok = incoterm.deleteIncoterm()
		case "CARRIER":
			var carrier Carrier
			carrier.Id = int32(id)
			carrier.enterprise = enterpriseId
			ok = carrier.deleteCarrier()
		case "SUPPLIER":
			var supplier Supplier
			supplier.Id = int32(id)
			supplier.enterprise = enterpriseId
			ok = supplier.deleteSupplier()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			documentContainer.Id = int32(id)
			documentContainer.enterprise = enterpriseId
			ok = documentContainer.deleteDocumentContainer()
		case "DOCUMENT":
			var document Document
			document.Id = int32(id)
			document.enterprise = enterpriseId
			ok = document.deleteDocument()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			productImage.Id = int32(id)
			ok = productImage.deleteProductImage(enterpriseId)
		}
	}

	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int64(id)
		saleOrder.enterprise = enterpriseId
		ok = saleOrder.deleteSalesOrder()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int64(id)
		saleOrderDetail.enterprise = enterpriseId
		ok = saleOrderDetail.deleteSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		saleOrderDiscount.Id = int32(id)
		saleOrderDiscount.enterprise = enterpriseId
		ok = saleOrderDiscount.deleteSalesOrderDiscount()
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int64(id)
		salesInvoice.enterprise = enterpriseId
		ok = salesInvoice.deleteSalesInvoice()
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int64(id)
		salesInvoiceDetail.enterprise = enterpriseId
		ok = salesInvoiceDetail.deleteSalesInvoiceDetail()
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int32(id)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		manufacturingOrder.enterprise = enterpriseId
		ok = manufacturingOrder.deleteManufacturingOrder()
	case "PACKAGING":
		if !permissions.Preparation {
			return
		}
		var packaging Packaging
		packaging.Id = int64(id)
		packaging.enterprise = enterpriseId
		ok = packaging.deletePackaging(enterpriseId)
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		warehouseMovement.enterprise = enterpriseId
		ok = warehouseMovement.deleteWarehouseMovement()
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int64(id)
		salesDeliveryNote.enterprise = enterpriseId
		ok = salesDeliveryNote.deleteSalesDeliveryNotes()
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		var shipping Shipping
		shipping.Id = int64(id)
		shipping.enterprise = enterpriseId
		ok = shipping.deleteShipping()
	case "USER":
		if !permissions.Admin {
			return
		}
		var user User
		user.Id = int32(id)
		user.enterprise = enterpriseId
		ok = user.deleteUser()
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		group.Id = int32(id)
		group.enterprise = enterpriseId
		ok = group.deleteGroup()
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int64(id)
		purchaseOrder.enterprise = enterpriseId
		ok = purchaseOrder.deletePurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int64(id)
		purchaseOrderDetail.enterprise = enterpriseId
		ok = purchaseOrderDetail.deletePurchaseOrderDetail()
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int64(id)
		purchaseDeliveryNote.enterprise = enterpriseId
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes()
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int64(id)
		purchaseInvoice.enterprise = enterpriseId
		ok = purchaseInvoice.deletePurchaseInvoice()
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int64(id)
		purchaseInvoiceDetail.enterprise = enterpriseId
		ok = purchaseInvoiceDetail.deletePurchaseInvoiceDetail()
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		pallet.Id = int32(id)
		pallet.enterprise = enterpriseId
		ok = pallet.deletePallet()
	case "JOURNAL":
		if !permissions.Accounting {
			return
		}
		var journal Journal
		journal.Id = int32(id)
		journal.enterprise = enterpriseId
		ok = journal.deleteJournal()
	case "ACCOUNT":
		if !permissions.Accounting {
			return
		}
		var account Account
		account.Id = int32(id)
		account.enterprise = enterpriseId
		ok = account.deleteAccount()
	case "ACCOUNTING_MOVEMENT":
		if !permissions.Accounting {
			return
		}
		var accountingMovement AccountingMovement
		accountingMovement.Id = int64(id)
		accountingMovement.enterprise = enterpriseId
		ok = accountingMovement.deleteAccountingMovement()
	case "ACCOUNTING_MOVEMENT_DETAIL":
		if !permissions.Accounting {
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		accountingMovementDetail.Id = int64(id)
		accountingMovementDetail.enterprise = enterpriseId
		ok = accountingMovementDetail.deleteAccountingMovementDetail()
	case "CHARGES":
		if !permissions.Accounting {
			return
		}
		var charges Charges
		charges.Id = int32(id)
		charges.enterprise = enterpriseId
		ok = charges.deleteCharges()
	case "PAYMENT":
		if !permissions.Accounting {
			return
		}
		var payment Payment
		payment.Id = int32(id)
		payment.enterprise = enterpriseId
		ok = payment.deletePayment()
	case "API_KEYS":
		if !permissions.Admin {
			return
		}
		var apiKey ApiKey
		apiKey.Id = int32(id)
		apiKey.enterprise = enterpriseId
		ok = apiKey.deleteApiKey()
	case "CONNECTION_FILTER":
		if !permissions.Admin {
			return
		}
		var filter ConnectionFilter
		filter.Id = int32(id)
		filter.enterprise = enterpriseId
		ok = filter.deleteConnectionFilter()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

type NameInt16 struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

type NameInt32 struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type NameString struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func instructionName(command string, message string, mt int, ws *websocket.Conn, enterpriseId int32) {
	var data []byte
	switch command {
	case "LANGUAGE":
		data, _ = json.Marshal(findLanguageByName(message, enterpriseId))
	case "CURRENCY":
		data, _ = json.Marshal(findCurrencyByName(message, enterpriseId))
	case "CUSTOMER":
		data, _ = json.Marshal(findCustomerByName(message, enterpriseId))
	case "COUNTRY":
		data, _ = json.Marshal(findCountryByName(message, enterpriseId))
	case "STATE":
		var cityName StateNameQuery
		json.Unmarshal([]byte(message), &cityName)
		data, _ = json.Marshal(findStateByName(cityName, enterpriseId))
	case "PAYMENT_METHOD":
		data, _ = json.Marshal(findPaymentMethodByName(message, enterpriseId))
	case "BILLING_SERIE":
		data, _ = json.Marshal(findBillingSerieByName(message, enterpriseId))
	case "PRODUCT_FAMILY":
		data, _ = json.Marshal(findProductFamilyByName(message, enterpriseId))
	case "COLOR":
		data, _ = json.Marshal(findColorByName(message, enterpriseId))
	case "PRODUCT":
		data, _ = json.Marshal(findProductByName(message, enterpriseId))
	case "WAREHOUSE":
		data, _ = json.Marshal(findWarehouseByName(message, enterpriseId))
	case "CARRIER":
		data, _ = json.Marshal(findCarrierByName(message, enterpriseId))
	case "SUPPLIER":
		data, _ = json.Marshal(findSupplierByName(message, enterpriseId))
	}
	ws.WriteMessage(mt, data)
}

func instructionGetName(command string, message string, mt int, ws *websocket.Conn, enterpriseId int32) {
	// attempt to get the name from resources with alpha key, if the resource if not found,
	// parse the input as number and attemp to get the name resource with numeric key
	var found bool = true
	var name string
	// ALPHA
	switch command {
	case "BILLING_SERIE":
		name = getNameBillingSerie(message, enterpriseId)
	case "WAREHOUSE":
		name = getNameWarehouse(message, enterpriseId)
	default:
		found = false
	}

	if found {
		ws.WriteMessage(mt, []byte(name))
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	switch command {
	case "LANGUAGE":
		name = getNameLanguage(int32(id), enterpriseId)
	case "CURRENCY":
		name = getNameCurrency(int32(id), enterpriseId)
	case "CUSTOMER":
		name = getNameCustomer(int32(id), enterpriseId)
	case "COUNTRY":
		name = getNameCountry(int32(id), enterpriseId)
	case "STATE":
		name = getNameState(int32(id), enterpriseId)
	case "PAYMENT_METHOD":
		name = getNamePaymentMethod(int32(id), enterpriseId)
	case "PRODUCT_FAMILY":
		name = getNameProductFamily(int32(id), enterpriseId)
	case "COLOR":
		name = getNameColor(int32(id), enterpriseId)
	case "ADDRESS":
		name = getAddressName(int32(id), enterpriseId)
	case "PRODUCT":
		name = getNameProduct(int32(id), enterpriseId)
	case "CARRIER":
		name = getNameCarrier(int32(id), enterpriseId)
	case "SALE_DELIERY_NOTE":
		name = getNameSalesDeliveryNote(int64(id), enterpriseId)
	case "SUPPLIER":
		name = getNameSupplier(int32(id), enterpriseId)
	}
	ws.WriteMessage(mt, []byte(name))
}

func instructionDefaults(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	// there are defaults that require an ID of a row, and there are defaults without parametres
	// attemps first respond to the parameterless, and if not found, parse the parameters and return

	var found bool = true
	var data []byte
	// ALPHA
	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSaleOrderDefaults(enterpriseId))
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderDefaults(enterpriseId))
	default:
		found = false
	}

	if found {
		ws.WriteMessage(mt, data)
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	switch command {
	case "CUSTOMER":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getCustomerDefaults(int32(id), enterpriseId))
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getOrderDetailDefaults(int32(id), enterpriseId))
	case "SUPPLIER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getSupplierDefaults(int32(id), enterpriseId))
	}
	ws.WriteMessage(mt, data)
}

func instructionLocate(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	var data []byte
	var found bool = true

	// PARAMETERLESS
	switch command {
	case "SALE_ORDER":
		data, _ = json.Marshal(locateSaleOrder(enterpriseId))
	case "DOCUMENT_CONTAINER":
		data, _ = json.Marshal(locateDocumentContainer(enterpriseId))
	case "LOCATE_ACCOUNT_CUSTOMER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateAccountForCustomer(enterpriseId))
	case "LOCATE_ACCOUNT_SUPPLIER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateAccountForSupplier(enterpriseId))
	case "LOCATE_ACCOUNT_BANKS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(locateAccountForBanks(enterpriseId))
	case "LOCATE_PRODUCT":
		if !permissions.Masters {
			return
		}
		var productLocateQuery ProductLocateQuery
		json.Unmarshal([]byte(message), &productLocateQuery)
		data, _ = json.Marshal(productLocateQuery.locateProduct(enterpriseId))
	case "LOCATE_CUSTOMER":
		if !permissions.Masters {
			return
		}
		var customerLocateQuery CustomerLocateQuery
		json.Unmarshal([]byte(message), &customerLocateQuery)
		data, _ = json.Marshal(customerLocateQuery.locateCustomers(enterpriseId))
	case "LOCATE_SUPPLIER":
		if !permissions.Masters {
			return
		}
		var supplierLocateQuery SupplierLocateQuery
		json.Unmarshal([]byte(message), &supplierLocateQuery)
		data, _ = json.Marshal(supplierLocateQuery.locateSuppliers(enterpriseId))
	case "CURRENCIES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateCurrency(enterpriseId))
	case "CARRIER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateCarriers(enterpriseId))
	case "PAYMENT_METHOD":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locatePaymentMethods(enterpriseId))
	case "BILLING_SERIE":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateBillingSeries(enterpriseId))
	default:
		found = false
	}

	if found {
		ws.WriteMessage(mt, data)
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	switch command {
	case "ADDRESS_CUSTOMER":
		data, _ = json.Marshal(locateAddressByCustomer(int32(id), enterpriseId))
	case "ADDRESS_SUPPLIER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(locateAddressBySupplier(int32(id), enterpriseId))
	case "SALE_DELIVERY_NOTE":
		data, _ = json.Marshal(locateSalesDeliveryNotesBySalesOrder(int64(id), enterpriseId))
	}
	ws.WriteMessage(mt, data)
}

func instructionAction(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, userId int32, enterpriseId int32) {
	var data []byte

	switch command {
	case "INVOICE_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllSaleOrder(int64(id), enterpriseId))
	case "INVOICE_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallySaleOrder(enterpriseId))
	case "GET_SALES_ORDER_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesOrderRelations(int64(id), enterpriseId))
	case "GET_SALES_INVOICE_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceRelations(int64(id), enterpriseId))
	case "TOGGLE_MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleManufactuedManufacturingOrder(int64(id), userId, enterpriseId))
	case "MANUFACTURING_ORDER_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderAllSaleOrder(int64(id), userId, enterpriseId))
	case "MANUFACTURING_ORDER_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var orderInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &orderInfo)
		data, _ = json.Marshal(orderInfo.manufacturingOrderPartiallySaleOrder(userId, enterpriseId))
	case "DELETE_SALES_ORDER_DETAIL_PACKAGED":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal([]byte(message), &salesOrderDetailPackaged)
		salesOrderDetailPackaged.enterprise = enterpriseId
		data, _ = json.Marshal(salesOrderDetailPackaged.deleteSalesOrderDetailPackaged(true))
	case "DELIVERY_NOTE_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllSaleOrder(int64(id), enterpriseId)
		data, _ = json.Marshal(ok)
	case "DELIVERY_NOTE_PARTIALLY_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var noteInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallySaleOrder(enterpriseId))
	case "SHIPPING_SALE_ORDER":
		if !permissions.Preparation {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(generateShippingFromSaleOrder(int64(id), enterpriseId))
	case "TOGGLE_SHIPPING_SENT":
		if !permissions.Preparation {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleShippingSent(int64(id), 1))
	case "SET_SHIPPING_COLLECTED":
		var shippings []int64
		json.Unmarshal([]byte(message), &shippings)
		data, _ = json.Marshal(setShippingCollected(shippings, enterpriseId))
	case "GET_SALES_DELIVERY_NOTE_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesDeliveryNoteRelations(int64(id), enterpriseId))
	case "USER_PWD":
		if !permissions.Admin {
			return
		}
		var userPassword UserPassword
		json.Unmarshal([]byte(message), &userPassword)
		data, _ = json.Marshal(userPassword.userPassword(enterpriseId))
	case "USER_AUTO_PWD":
		// every user can use this
		var userPassword UserAutoPassword
		json.Unmarshal([]byte(message), &userPassword)
		data, _ = json.Marshal(userPassword.userAutoPassword(enterpriseId, userId))
	case "GET_CURRENT_USER":
		// every user can use this
		data, _ = json.Marshal(getUserRow(userId))
	case "USER_OFF":
		if !permissions.Admin {
			return
		}
		var user User
		json.Unmarshal([]byte(message), &user)
		user.enterprise = enterpriseId
		data, _ = json.Marshal(user.offUser())
	case "PURCHASE_NEEDS":
		if !permissions.Purchases {
			return
		}
		var needs []PurchaseNeed
		json.Unmarshal([]byte(message), &needs)
		data, _ = json.Marshal(generatePurchaseOrdersFromNeeds(needs, enterpriseId))
	case "DELIVERY_NOTE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllPurchaseOrder(int64(id), enterpriseId)
		data, _ = json.Marshal(ok)
	case "GET_PURCHASE_ORDER_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderRelations(int64(id), enterpriseId))
	case "GET_INVOICE_ORDER_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceRelations(int64(id), enterpriseId))
	case "GET_PURCHASE_DELIVERY_NOTE_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNoteRelations(int32(id), enterpriseId))
	case "DELIVERY_NOTE_PARTIALLY_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var noteInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallyPurchaseOrder(enterpriseId))
	case "INVOICE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllPurchaseOrder(int64(id), enterpriseId))
	case "INVOICE_PARTIAL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallyPurchaseOrder(enterpriseId))
	case "INSERT_DOCUMENT":
		var document Document
		json.Unmarshal([]byte(message), &document)
		document.enterprise = enterpriseId
		ok := document.insertDocument()
		if ok {
			data, _ = json.Marshal(document)
		} else {
			data, _ = json.Marshal(ok)
		}
	case "GRANT_DOCUMENT_ACCESS_TOKEN":
		data, _ = json.Marshal(grantDocumentAccessToken(enterpriseId))
	case "GET_PRODUCT_ROW":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		p := getProductRow(int32(id))
		if p.enterprise != enterpriseId {
			data = []byte("false")
			break
		}
		data, _ = json.Marshal(p)
	case "PRODUCT_EAN13":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		p := getProductRow(int32(id))
		if p.enterprise != enterpriseId {
			data = []byte("false")
			break
		}
		p.generateBarcode(enterpriseId)
		data, _ = json.Marshal(p.updateProduct())
	case "EMAIL":
		if !(permissions.Sales || permissions.Purchases) {
			return
		}
		var emailInfo EmailInfo
		json.Unmarshal([]byte(message), &emailInfo)
		data, _ = json.Marshal(emailInfo.sendEmail(enterpriseId))
	case "EXPORT":
		if !permissions.Admin {
			return
		}
		var exportInfo ExportInfo
		json.Unmarshal([]byte(message), &exportInfo)
		data, _ = json.Marshal(exportInfo.export())
	case "EXPORT_JSON":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(exportToJSON(message, enterpriseId))
	case "IMPORT_JSON":
		if !permissions.Admin {
			return
		}
		var importInfo ImportInfo
		json.Unmarshal([]byte(message), &importInfo)
		data, _ = json.Marshal(importInfo.importJson(enterpriseId))
	case "REGENERATE_DRAGGED_STOCK":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(regenerateDraggedStock(message, enterpriseId))
	case "REGENERATE_PRODUCT_STOCK":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(regenerateProductStock(enterpriseId))
	case "DISCONNECT":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(disconnectConnection(message, enterpriseId))
	case "PRESTASHOP":
		if !permissions.PrestaShop {
			return
		}
		importFromPrestaShop(enterpriseId)
	case "WOOCOMMERCE":
		if !permissions.Admin {
			return
		}
		importFromWooCommerce(enterpriseId)
	case "CALCULATE_MINIMUM_STOCK":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(calculateMinimumStock(enterpriseId))
	case "GENERATE_MANUFACTURIG_OR_PURCHASE_ORDERS_MINIMUM_STOCK":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(generateManufacturingOrPurchaseOrdersMinimumStock(userId, enterpriseId))
	case "SALES_POST_INVOICES":
		if !permissions.Accounting {
			return
		}
		var invoiceIds []int64
		json.Unmarshal([]byte(message), &invoiceIds)
		data, _ = json.Marshal(salesPostInvoices(invoiceIds, enterpriseId))
	case "PURCHASE_POST_INVOICES":
		if !permissions.Accounting {
			return
		}
		var invoiceIds []int64
		json.Unmarshal([]byte(message), &invoiceIds)
		data, _ = json.Marshal(purchasePostInvoices(invoiceIds, enterpriseId))
	case "MANUFACTURING_ORDER_TAG_PRINTED":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderTagPrinted(int64(id), userId, enterpriseId))
	case "CANCEL_SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(cancelSalesOrderDetail(int64(id), enterpriseId))
	case "API_KEYS":
		if !permissions.Admin {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		a := ApiKey{}
		a.Id = int32(id)
		a.enterprise = enterpriseId
		data, _ = json.Marshal(a.offApiKey())
	case "SHOPIFY":
		if !permissions.Admin {
			return
		}
		importFromShopify(enterpriseId)
	case "EVALUATE_PASSWORD_SECURE_CLOUD":
		data, _ = json.Marshal(evaluatePasswordSecureCloud(enterpriseId, message))
	case "PRODUCT_GENERATOR":
		if !permissions.Manufacturing {
			return
		}
		var productGenerator ProductGenerator
		json.Unmarshal([]byte(message), &productGenerator)
		data, _ = json.Marshal(productGenerator.productGenerator(enterpriseId))
	case "REGISTER_USER_IN_GOOGLE_AUTHENTICATOR":
		if !permissions.Admin {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(registerUserInGoogleAuthenticator(int32(id), enterpriseId))
	case "REMOVE_USER_IN_GOOGLE_AUTHENTICATOR":
		if !permissions.Admin {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(removeUserFromGoogleAuthenticator(int32(id), enterpriseId))
	case "TOGGLE_SIMPLIFIED_INVOICE_SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleSimplifiedInvoiceSalesInvoice(int64(id), enterpriseId))
	case "MAKE_AMENDING_SALE_INVOICE":
		if !permissions.Sales {
			return
		}
		var makeAmendingInvoice MakeAmendingInvoice
		json.Unmarshal([]byte(message), &makeAmendingInvoice)
		data, _ = json.Marshal(makeAmendingSaleInvoice(makeAmendingInvoice.InvoiceId, enterpriseId, makeAmendingInvoice.Quantity, makeAmendingInvoice.Description))
	case "MAKE_AMENDING_PURCHASE_INVOICE":
		if !permissions.Sales {
			return
		}
		var makeAmendingInvoice MakeAmendingInvoice
		json.Unmarshal([]byte(message), &makeAmendingInvoice)
		data, _ = json.Marshal(makeAmendingPurchaseInvoice(makeAmendingInvoice.InvoiceId, enterpriseId, makeAmendingInvoice.Quantity, makeAmendingInvoice.Description))
	}
	ws.WriteMessage(mt, data)
}

type PaginatedSearch struct {
	PaginationQuery
	Search string `json:"search"`
}

func instructionSearch(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32) {
	var data []byte
	switch command {
	case "CUSTOMER":
		if !permissions.Masters {
			return
		}
		var paginatedSearch PaginatedSearch
		json.Unmarshal([]byte(message), &paginatedSearch)
		paginatedSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(paginatedSearch.searchCustomers())
	case "SUPPLER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchSuppliers(message, enterpriseId))
	case "PRODUCT":
		if !permissions.Masters {
			return
		}
		var productSearch ProductSearch
		json.Unmarshal([]byte(message), &productSearch)
		data, _ = json.Marshal(productSearch.searchProduct(enterpriseId))
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(searchShippings(message, enterpriseId))
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var salesOrderSearch SalesOrderSearch
		json.Unmarshal([]byte(message), &salesOrderSearch)
		salesOrderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(salesOrderSearch.searchSalesOrder())
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchSalesInvoices())
	case "SALES_DELIVERY_NOTE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchSalesDelvieryNotes())
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchPurchaseOrder())
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchPurchaseInvoice())
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchPurchaseDeliveryNote())
	case "COUNTRY":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchCountries(message, enterpriseId))
	case "STATE":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchStates(message, enterpriseId))
	case "ADDRESS":
		if !permissions.Masters {
			return
		}
		var paginatedSearch PaginatedSearch
		json.Unmarshal([]byte(message), &paginatedSearch)
		paginatedSearch.Enterprise = enterpriseId
		data, _ = json.Marshal(paginatedSearch.searchAddresses())
	case "LANGUAGE":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchLanguages(message, enterpriseId))
	case "WAREHOUSE_MOVEMENT":
		if !permissions.Warehouse {
			return
		}
		var warehouseMovement WarehouseMovementSearch
		json.Unmarshal([]byte(message), &warehouseMovement)
		warehouseMovement.Enterprise = enterpriseId
		data, _ = json.Marshal(warehouseMovement.searchWarehouseMovement())
	case "ACCOUNT":
		if !permissions.Accounting {
			return
		}
		var accountSearch AccountSearch
		json.Unmarshal([]byte(message), &accountSearch)
		data, _ = json.Marshal(accountSearch.searchAccounts(enterpriseId))
	case "ACCOUNTING_MOVEMENTS":
		if !permissions.Accounting {
			return
		}
		data, _ = json.Marshal(searchAccountingMovements(message, enterpriseId))
	}
	ws.WriteMessage(mt, data)
}

func isParameterPresent(parameter string) bool {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == parameter {
			return true
		}
	}
	return false
}

func getParameterValue(parameter string) (string, bool) {
	for i := 1; i < len(os.Args); i++ {
		parameterValue := strings.Split(os.Args[i], "=")
		if len(parameterValue) == 2 && parameterValue[0] == parameter {
			return parameterValue[1], true
		}
	}
	return "", false
}
