package main

import (
	"database/sql"
	"encoding/base64"
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

const MAX_INT32 = 2147483647

// Basic, static, server settings such as the DB password or the port.
var settings BackendSettings

// Http object for the websocket clients to conenect to.
var upgrader = websocket.Upgrader{}

// Database connection to PostgreSQL.
var db *sql.DB

// List of all the concurrent websocket connections to the server.
var connections []Connection

// MUTEX FOR var connections []Connection: List of all the concurrent websocket connections to the server.
var connectionsMutex sync.Mutex

// Global cron instance
var c *cron.Cron

// List of all cron IDs. Key= Enterprise Id, Value= Array of objects with cron IDs.
var runningCrons map[int32]EnterpriseCronInfo = make(map[int32]EnterpriseCronInfo)

// MUTEX FOR: var runningCrons map[int32]EnterpriseCronInfo: List of all cron IDs. Key= Enterprise Id, Value= Array of objects with cron IDs.
var runningCronsMutex sync.Mutex

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

	// add passwords to blacklist
	if isParameterPresent("--add-pwd-blacklist") {
		addPasswordsToBlacklist()
		os.Exit(0)
	}

	// listen to requests
	fmt.Println("Server ready! :D")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", reverse)
	http.HandleFunc("/document", handleDocument)
	http.HandleFunc("/report", generateReport)
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
		if settingsRecords[i].CronSendcloudTracking != "" {
			cronId, err := c.AddFunc(settingsRecords[i].CronSendcloudTracking, func() {
				getShippingTrackingSendCloud(enterpriseId)
			})
			if err != nil {
				enterpriseCronInfo.CronSendcloudTracking = &cronId
			}
		}
		runningCrons[enterpriseId] = enterpriseCronInfo
	}
	c.AddFunc(settings.Server.CronClearLogs, clearLogs)
	c.AddFunc("@every 1m", resetMaxRequestsPerEnterprise)
	c.AddFunc("@every 1h", crashreporter)
	c.AddFunc("@every 5m", attemptToSendQueuedWebHooks)
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
		ws.Close()
		return
	}
	setUserDateLastLogin(userId)
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
				if result.Ok {
					if result.checkUserConnection(userId, remoteAddr, enterpriseId) {
						ws.WriteMessage(mt, data)
						return true, userId, result.Permissions, enterpriseId
					} else {
						data, _ := json.Marshal(result)
						ws.WriteMessage(mt, data)
						// END AUTHENTICATION
						return false, 0, nil, 0
					}
				} else {
					ws.WriteMessage(mt, data)
				}
			}
		} else {

			// Return result to client (Ok + Token)
			data, _ := json.Marshal(result)
			if result.Ok {
				if result.checkUserConnection(userId, remoteAddr, enterpriseId) {
					ws.WriteMessage(mt, data)
					return true, userId, result.Permissions, enterpriseId
				} else {
					data, _ := json.Marshal(result)
					ws.WriteMessage(mt, data)
					// END AUTHENTICATION
					return false, 0, nil, 0
				}
			} else {
				ws.WriteMessage(mt, data)
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
		instructionUpdate(command, message, mt, ws, permissions, enterpriseId, userId)
	case "DELETE":
		instructionDelete(command, string(message), mt, ws, permissions, enterpriseId, userId)
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
	enterprise int32
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
			paginationQuery.enterprise = enterpriseId
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
			paginationQuery.enterprise = enterpriseId
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
		paginationQuery.enterprise = enterpriseId
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
		paginationQuery.enterprise = enterpriseId
		data, _ = json.Marshal(paginationQuery.getWarehouseMovement())
	case "WAREHOUSE_WAREHOUSE_MOVEMENTS":
		if !permissions.Warehouse {
			return
		}
		var warehouseMovementByWarehouse WarehouseMovementByWarehouse
		json.Unmarshal([]byte(message), &warehouseMovementByWarehouse)
		warehouseMovementByWarehouse.enterprise = enterpriseId
		data, _ = json.Marshal(warehouseMovementByWarehouse.getWarehouseMovementByWarehouse())
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.enterprise = enterpriseId
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
	case "BENEFITS_STATISTICS":
		var benefitsStatisticsQuery BenefitsStatisticsQuery
		json.Unmarshal([]byte(message), &benefitsStatisticsQuery)
		data, _ = json.Marshal(benefitsStatisticsQuery.benefitsStatistics(enterpriseId))
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
	case "COMPLEX_MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var complexManufacturingPaginationQuery ComplexManufacturingPaginationQuery
		json.Unmarshal([]byte(message), &complexManufacturingPaginationQuery)
		data, _ = json.Marshal(complexManufacturingPaginationQuery.getComplexManufacturingOrder(enterpriseId))
	case "CONNECTION_LOG":
		if !permissions.Admin {
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal([]byte(message), &paginationQuery)
		paginationQuery.enterprise = enterpriseId
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
	case "REGISTER_TRANSACTIONAL_LOG":
		var query TransactionalLogQuery
		json.Unmarshal([]byte(message), &query)
		query.enterpriseId = enterpriseId
		data, _ = json.Marshal(query.getRegisterTransactionalLogs())
	case "EMAIL_LOGS":
		var search EmailLogSearch
		json.Unmarshal([]byte(message), &search)
		data, _ = json.Marshal(search.getEmailLogs(enterpriseId))
	case "POS_TERMINALS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getPOSTerminals(enterpriseId))
	case "PERMISSION_DICTIONARY":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getPermissionDictionary(enterpriseId))
	case "PERMISSION_DICTIONARY_GRUPS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getGroupsPermissionDictionary(enterpriseId, message))
	case "TRIAL_BALANCE":
		if !permissions.Accounting {
			return
		}
		var query TrialBalanceQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.getTrialBalance(enterpriseId))
	case "PRODUCT_SALES_ORDER":
		if !permissions.Masters {
			return
		}
		var query ProductSalesOrderDetailsQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(getProductSalesOrderDetails(query, enterpriseId))
	case "PRODUCT_PURCHASE_ORDER":
		if !permissions.Masters {
			return
		}
		var query ProductPurchaseOrderDetailsQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(getProductPurchaseOrderDetails(query, enterpriseId))
	case "PRODUCT_WAREHOUSE_MOVEMENT":
		if !permissions.Masters {
			return
		}
		var query ProductPurchaseOrderDetailsQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(getProductWarehouseMovement(query, enterpriseId))
	case "PRODUCT_MANUFACTURING_ORDERS":
		if !permissions.Masters {
			return
		}
		var query ProductManufacturingOrdersQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(getProductManufacturingOrders(query, enterpriseId))
	case "PRODUCT_COMPLEX_MANUFACTURING_ORDERS":
		if !permissions.Masters {
			return
		}
		var query ProductManufacturingOrdersQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(getProductComplexManufacturingOrders(query, enterpriseId))
	case "REPORT_TEMPLATE_TRANSLATION":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getReportTemplateTranslations(enterpriseId))
	case "HS_CODES":
		var query HSCodeQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.getHSCodes())
	case "REPORT_111":
		if !permissions.Accounting {
			return
		}
		var query Form111Query
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.execReportForm111(enterpriseId))
	case "REPORT_115":
		if !permissions.Accounting {
			return
		}
		var query Form115Query
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.execReportForm115(enterpriseId))
	case "INVENTORY":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getInventories(enterpriseId))
	case "INVENTORY_VALUATION":
		if !permissions.Accounting {
			return
		}
		var query InventoyValuationQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.getInventoyValuation(enterpriseId))
	case "WEBHOOK_SETTINGS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getWebHookSettings(enterpriseId))
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
		data, _ = json.Marshal(getSalesOrderDiscounts(int64(id), enterpriseId))
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
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getProductSalesOrderDetailsPending(int32(id), enterpriseId))
	case "PRODUCT_PURCHASE_ORDER_PENDING":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getProductPurchaseOrderDetailsPending(int32(id), enterpriseId))
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
	case "COMPLEX_MANUFACTURING_ORDERS_FROM_PURCHASE_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getComplexManufacturingOrdersFromPurchaseOrderDetail(int64(id), enterpriseId))
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
	case "SHIPPING_STATUS_HISTORY":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getShippingStatusHistory(enterpriseId, int64(id)))
	case "SALES_ORDER_DETAIL_DIGITAL_PRODUCT_DATA":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderDetailDigitalProductData(int64(id), enterpriseId))
	case "MANUFACTURING_ORDER_TYPE_COMPONENTS":
		if !permissions.Manufacturing {
			return
		}
		data, _ = json.Marshal(getManufacturingOrderTypeComponents(int32(id), enterpriseId))
	case "COMPLEX_MANUFACTURING_ORDER_MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		data, _ = json.Marshal(getComplexManufacturingOrderManufacturingOrder(int64(id), enterpriseId))
	case "MANUFACTURING_ORDER_TYPE_PRODUCTS":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getProductsByManufacturingOrderType(int32(id), enterpriseId))
	case "GROUP_PERMISSION_DICTIONARY":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getGroupPermissionDictionary(enterpriseId, int32(id)))
	case "PRODUCT_ACCOUNTS":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getProductAccounts(int32(id), enterpriseId))
	case "WAREHOUSE_MOVEMENT_RELATIONS":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementRelations(int64(id), enterpriseId))
	case "INVENTORY_PODUCTS":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getInventoryProducts(int32(id), enterpriseId))
	case "WEBHOOK_QUEUE":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getWebHookRequestQueue(enterpriseId, int32(id)))
	case "WEBHOOK_LOGS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getWebHookLogs(enterpriseId, int32(id)))
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
			operationResult = address.insertAddress(userId)
			found = true
		case "CUSTOMER":
			var customer Customer
			json.Unmarshal(message, &customer)
			customer.enterprise = enterpriseId
			operationResult = customer.insertCustomer(userId)
			found = true
		case "SUPPLIER":
			var supplier Supplier
			json.Unmarshal(message, &supplier)
			supplier.enterprise = enterpriseId
			operationResult = supplier.insertSupplier(userId)
			found = true
		}

		if found {
			data, _ := json.Marshal(operationResult)
			ws.WriteMessage(mt, data)
			return
		}
	} // Masters

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
			ok, _ = shipping.insertShipping(userId, nil)
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
		if !permissions.Sales || getUserPermission("CANT_MANUALLY_CREATE_SALE_ORDER", enterpriseId, userId) {
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal([]byte(message), &saleOrder)
		saleOrder.enterprise = enterpriseId
		ok, orderId := saleOrder.insertSalesOrder(userId)
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			order := getSalesOrderRow(orderId)
			returnData, _ = json.Marshal(order)
		}
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		json.Unmarshal(message, &saleOrderDetail)
		saleOrderDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(saleOrderDetail.insertSalesOrderDetail(userId))
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		json.Unmarshal(message, &salesInvoiceDetail)
		salesInvoiceDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(salesInvoiceDetail.insertSalesInvoiceDetail(nil, userId))
	case "SALES_INVOICE":
		if !permissions.Sales || getUserPermission("CANT_MANUALLY_CREATE_SALE_INVOICE", enterpriseId, userId) {
			return
		}
		var saleInvoice SalesInvoice
		json.Unmarshal(message, &saleInvoice)
		saleInvoice.enterprise = enterpriseId
		ok, invoiceId := saleInvoice.insertSalesInvoice(userId, nil)
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			invoice := getSalesInvoiceRow(invoiceId)
			returnData, _ = json.Marshal(invoice)
		}
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales || getUserPermission("CANT_MANUALLY_CREATE_SALE_DELIVERY_NOTE", enterpriseId, userId) {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		json.Unmarshal(message, &salesDeliveryNote)
		salesDeliveryNote.enterprise = enterpriseId
		ok, nodeId := salesDeliveryNote.insertSalesDeliveryNotes(userId, nil)
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
		ok, orderId := purchaseOrder.insertPurchaseOrder(userId, nil)
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
		ok, invoiceId := purchaseInvoice.insertPurchaseInvoice(userId, nil)
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			invoice := getPurchaseInvoiceRow(invoiceId)
			returnData, _ = json.Marshal(invoice)
		}
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases || getUserPermission("CANT_MANUALLY_CREATE_PURCHASE_DELIVERY_NOTE", enterpriseId, userId) {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		json.Unmarshal(message, &purchaseDeliveryNote)
		purchaseDeliveryNote.enterprise = enterpriseId
		ok, noteId := purchaseDeliveryNote.insertPurchaseDeliveryNotes(userId, nil)
		if !ok {
			returnData, _ = json.Marshal(nil)
		} else {
			note := getPurchaseDeliveryNoteRow(noteId)
			returnData, _ = json.Marshal(note)
		}
	case "MANUFACTURING_ORDER_TYPE_COMPONENTS":
		if !permissions.Manufacturing {
			return
		}
		var c ManufacturingOrderTypeComponents
		json.Unmarshal(message, &c)
		c.enterprise = enterpriseId
		ok, errorCode := c.insertManufacturingOrderTypeComponents()
		isValid := OkAndErrorCodeReturn{
			Ok:       ok,
			ErorCode: errorCode,
		}
		returnData, _ = json.Marshal(isValid)
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		purchaseOrderDetail.enterprise = enterpriseId
		ok, _ := purchaseOrderDetail.insertPurchaseOrderDetail(userId, nil)
		returnData, _ = json.Marshal(ok)
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		json.Unmarshal(message, &purchaseInvoiceDetail)
		purchaseInvoiceDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseInvoiceDetail.insertPurchaseInvoiceDetail(userId, nil))
	case "PRODUCT":
		if !permissions.Masters {
			return
		}
		if getUserPermission("CANT_CREATE_PRODUCT", enterpriseId, userId) {
			return
		}
		var product Product
		json.Unmarshal(message, &product)
		product.enterprise = enterpriseId
		returnData, _ = json.Marshal(product.insertProduct(userId))
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_MANUALLY_CREATE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(message, &manufacturingOrder)
		manufacturingOrder.UserCreated = userId
		manufacturingOrder.enterprise = enterpriseId
		returnData, _ = json.Marshal(manufacturingOrder.insertManufacturingOrder(userId, nil))
	case "MULTIPLE_MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_MANUALLY_CREATE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var manufacturingOrder MultipleManufacturingOrders
		json.Unmarshal(message, &manufacturingOrder)
		manufacturingOrder.Order.UserCreated = userId
		manufacturingOrder.Order.enterprise = enterpriseId
		returnData, _ = json.Marshal(manufacturingOrder.insertMultipleManufacturingOrders(userId))
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
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		json.Unmarshal(message, &saleOrderDiscount)
		saleOrderDiscount.enterprise = enterpriseId
		ok = saleOrderDiscount.insertSalesOrderDiscount(userId)
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "COMPLEX_MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_MANUALLY_CREATE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var complexManufacturingOrder ComplexManufacturingOrder
		json.Unmarshal(message, &complexManufacturingOrder)
		complexManufacturingOrder.UserCreated = userId
		complexManufacturingOrder.enterprise = enterpriseId
		ok, _ = complexManufacturingOrder.insertComplexManufacturingOrder(userId, nil)
	case "MULTIPLE_COMPLEX_MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_MANUALLY_CREATE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var complexManufacturingOrder MultipleComplexManufacturingOrders
		json.Unmarshal(message, &complexManufacturingOrder)
		complexManufacturingOrder.Order.UserCreated = userId
		complexManufacturingOrder.Order.enterprise = enterpriseId
		ok = complexManufacturingOrder.insertMultipleComplexManufacturingOrders(userId)
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
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackaged(userId)
	case "SALES_ORDER_DETAIL_PACKAGED_EAN13":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackagedEAN13
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackagedEAN13(enterpriseId, userId)
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		json.Unmarshal(message, &warehouseMovement)
		warehouseMovement.enterprise = enterpriseId
		ok = warehouseMovement.insertWarehouseMovement(userId, nil)
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
		ok = accountingMovement.insertAccountingMovement(userId, nil)
	case "ACCOUNTING_MOVEMENT_DETAIL":
		if !permissions.Accounting {
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		json.Unmarshal(message, &accountingMovementDetail)
		accountingMovementDetail.enterprise = enterpriseId
		ok = accountingMovementDetail.insertAccountingMovementDetail(userId, nil)
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
		ok = charges.insertCharges(userId)
	case "PAYMENT":
		if !permissions.Accounting {
			return
		}
		var payment Payment
		json.Unmarshal(message, &payment)
		payment.enterprise = enterpriseId
		ok = payment.insertPayment(userId)
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
	case "SALES_ORDER_DETAIL_DIGITAL_PRODUCT_DATA":
		if !permissions.Sales {
			return
		}
		var d SalesOrderDetailDigitalProductData
		json.Unmarshal(message, &d)
		ok = d.insertSalesOrderDetailDigitalProductData(enterpriseId)
	case "PERMISSION_DICTIONARY_GROUP":
		if !permissions.Admin {
			return
		}
		var d PermissionDictionaryGroup
		json.Unmarshal(message, &d)
		d.enterprise = enterpriseId
		ok = d.insertPermissionDictionaryGroup()
	case "PRODUCT_ACCOUNTS":
		if !permissions.Accounting {
			return
		}
		var d ProductAccount
		json.Unmarshal(message, &d)
		d.enterprise = enterpriseId
		ok = d.insertProductAccount()
	case "REPORT_TEMPLATE_TRANSLATION":
		if !permissions.Admin {
			return
		}
		var t ReportTemplateTranslation
		json.Unmarshal(message, &t)
		t.enterprise = enterpriseId
		ok = t.insertReportTemplateTranslation()
	case "INVENTORY":
		if !permissions.Warehouse {
			return
		}
		var i Inventory
		json.Unmarshal(message, &i)
		ok = i.insertInventory(enterpriseId)
	case "WEBHOOK_SETTINGS":
		if !permissions.Admin {
			return
		}
		var s WebHookSettings
		json.Unmarshal(message, &s)
		ok = s.insertWebHookSettings(enterpriseId)
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionUpdate(command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32, userId int32) {
	var ok bool
	var returnData []byte
	switch command {
	case "MANUFACTURING_ORDER_TYPE_COMPONENTS":
		if !permissions.Manufacturing {
			return
		}
		var c ManufacturingOrderTypeComponents
		json.Unmarshal(message, &c)
		c.enterprise = enterpriseId
		ok, errorCode := c.updateManufacturingOrderTypeComponents()
		isValid := OkAndErrorCodeReturn{
			Ok:       ok,
			ErorCode: errorCode,
		}
		returnData, _ = json.Marshal(isValid)
		ok = true
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(message, &salesOrderDetail)
		salesOrderDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(salesOrderDetail.updateSalesOrderDetail(userId))
		ok = true
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		purchaseOrderDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseOrderDetail.updatePurchaseOrderDetail(userId))
		ok = true
	case "PRODUCT":
		if getUserPermission("CANT_UPDATE_DELETE_PRODUCT", enterpriseId, userId) {
			return
		}
		var product Product
		json.Unmarshal(message, &product)
		product.enterprise = enterpriseId
		returnData, _ = json.Marshal(product.updateProduct(userId))
		ok = true
	}
	if ok {
		ws.WriteMessage(mt, returnData)
		return
	}

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
			ok = customer.updateCustomer(userId)
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
			ok = supplier.updateSupplier(userId)
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
	} // Masters

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
		ok = saleOrder.updateSalesOrder(userId)
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
		ok = shipping.updateShipping(userId)
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
		ok = purchaseOrder.updatePurchaseOrder(userId)
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
	case "SALES_ORDER_DETAIL_DIGITAL_PRODUCT_DATA":
		if !permissions.Sales {
			return
		}
		var d SalesOrderDetailDigitalProductData
		json.Unmarshal(message, &d)
		ok = d.updateSalesOrderDetailDigitalProductData(enterpriseId)
	case "POS_TERMINAL":
		if !permissions.Admin {
			return
		}
		var t POSTerminal
		json.Unmarshal(message, &t)
		t.enterprise = enterpriseId
		ok = t.updatePOSTerminal()
	case "API_KEY":
		if !permissions.Admin {
			return
		}
		var a ApiKey
		json.Unmarshal(message, &a)
		a.enterprise = enterpriseId
		ok = a.updateApiKey()
	case "PRODUCT_ACCOUNTS":
		if !permissions.Accounting {
			return
		}
		var d ProductAccount
		json.Unmarshal(message, &d)
		d.enterprise = enterpriseId
		ok = d.updateProductAccount()
	case "REPORT_TEMPLATE_TRANSLATION":
		if !permissions.Admin {
			return
		}
		var t ReportTemplateTranslation
		json.Unmarshal(message, &t)
		t.enterprise = enterpriseId
		ok = t.updateReportTemplateTranslation()
	case "WEBHOOK_SETTINGS":
		if !permissions.Admin {
			return
		}
		var s WebHookSettings
		json.Unmarshal(message, &s)
		ok = s.updateWebHookSettings(enterpriseId)
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionDelete(command string, message string, mt int, ws *websocket.Conn, permissions Permissions, enterpriseId int32, userId int32) {
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
	case "PERMISSION_DICTIONARY_GROUP":
		if !permissions.Admin {
			return
		}
		var p PermissionDictionaryGroup
		json.Unmarshal([]byte(message), &p)
		p.enterprise = enterpriseId
		ok = p.deletePermissionDictionaryGroup()
	case "REPORT_TEMPLATE_TRANSLATION":
		if !permissions.Admin {
			return
		}
		var t ReportTemplateTranslation
		json.Unmarshal([]byte(message), &t)
		t.enterprise = enterpriseId
		ok = t.deleteReportTemplateTranslation()
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
			ok = customer.deleteCustomer(userId)
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
			ok = supplier.deleteSupplier(userId)
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

	var returnData []byte
	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int64(id)
		saleOrder.enterprise = enterpriseId
		returnData, _ = json.Marshal(saleOrder.deleteSalesOrder(userId))
		found = true
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int64(id)
		saleOrderDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(saleOrderDetail.deleteSalesOrderDetail(userId, nil))
		found = true
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int64(id)
		salesInvoice.enterprise = enterpriseId
		returnData, _ = json.Marshal(salesInvoice.deleteSalesInvoice(userId))
		found = true
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int64(id)
		salesInvoiceDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(salesInvoiceDetail.deleteSalesInvoiceDetail(userId, nil))
		found = true
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int64(id)
		purchaseOrderDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseOrderDetail.deletePurchaseOrderDetail(userId, nil))
		found = true
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int64(id)
		purchaseOrder.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseOrder.deletePurchaseOrder(userId))
		found = true
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int64(id)
		purchaseInvoice.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseInvoice.deletePurchaseInvoice(userId, nil))
		found = true
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int64(id)
		purchaseInvoiceDetail.enterprise = enterpriseId
		returnData, _ = json.Marshal(purchaseInvoiceDetail.deletePurchaseInvoiceDetail(userId, nil))
		found = true
	case "PRODUCT":
		if getUserPermission("CANT_UPDATE_DELETE_PRODUCT", enterpriseId, userId) {
			return
		}
		var product Product
		product.Id = int32(id)
		product.enterprise = enterpriseId
		returnData, _ = json.Marshal(product.deleteProduct(userId))
		found = true
	}
	if found {
		ws.WriteMessage(mt, returnData)
		return
	}

	switch command {
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		saleOrderDiscount.Id = int32(id)
		saleOrderDiscount.enterprise = enterpriseId
		ok = saleOrderDiscount.deleteSalesOrderDiscount(userId)
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int32(id)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_DELETE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		manufacturingOrder.enterprise = enterpriseId
		ok = manufacturingOrder.deleteManufacturingOrder(userId, nil)
	case "COMPLEX_MANUFACTURING_ORDER":
		if !permissions.Manufacturing || getUserPermission("CANT_DELETE_MANUFACTURING_ORDERS", enterpriseId, userId) {
			return
		}
		var complexManufacturingOrder ComplexManufacturingOrder
		complexManufacturingOrder.Id = int64(id)
		complexManufacturingOrder.enterprise = enterpriseId
		ok = complexManufacturingOrder.deleteComplexManufacturingOrder(userId, nil)
	case "PACKAGING":
		if !permissions.Preparation {
			return
		}
		var packaging Packaging
		packaging.Id = int64(id)
		packaging.enterprise = enterpriseId
		ok = packaging.deletePackaging(enterpriseId, userId)
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		warehouseMovement.enterprise = enterpriseId
		ok = warehouseMovement.deleteWarehouseMovement(userId, nil)
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int64(id)
		salesDeliveryNote.enterprise = enterpriseId
		ok = salesDeliveryNote.deleteSalesDeliveryNotes(userId, nil).Ok
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		var shipping Shipping
		shipping.Id = int64(id)
		shipping.enterprise = enterpriseId
		ok = shipping.deleteShipping(userId)
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

	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int64(id)
		purchaseDeliveryNote.enterprise = enterpriseId
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes(userId, nil)
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
		ok = accountingMovement.deleteAccountingMovement(userId, nil)
	case "ACCOUNTING_MOVEMENT_DETAIL":
		if !permissions.Accounting {
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		accountingMovementDetail.Id = int64(id)
		accountingMovementDetail.enterprise = enterpriseId
		ok = accountingMovementDetail.deleteAccountingMovementDetail(userId, nil)
	case "CHARGES":
		if !permissions.Accounting {
			return
		}
		var charges Charges
		charges.Id = int32(id)
		charges.enterprise = enterpriseId
		ok = charges.deleteCharges(userId)
	case "PAYMENT":
		if !permissions.Accounting {
			return
		}
		var payment Payment
		payment.Id = int32(id)
		payment.enterprise = enterpriseId
		ok = payment.deletePayment(userId)
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
	case "SALES_ORDER_DETAIL_DIGITAL_PRODUCT_DATA":
		if !permissions.Sales {
			return
		}
		d := SalesOrderDetailDigitalProductData{}
		d.Id = int32(id)
		ok = d.deleteSalesOrderDetailDigitalProductData(enterpriseId)
	case "MANUFACTURING_ORDER_TYPE_COMPONENTS":
		if !permissions.Manufacturing {
			return
		}
		var c ManufacturingOrderTypeComponents
		c.Id = int32(id)
		c.enterprise = enterpriseId
		ok = c.deleteManufacturingOrderTypeComponents()
	case "PRODUCT_ACCOUNTS":
		if !permissions.Accounting {
			return
		}
		var d ProductAccount
		d.Id = int32(id)
		d.enterprise = enterpriseId
		ok = d.deleteProductAccount()
	case "INVENTORY":
		if !permissions.Warehouse {
			return
		}
		var i Inventory = Inventory{}
		i.Id = int32(id)
		i.enterprise = enterpriseId
		ok = i.deleteInventory(enterpriseId)
	case "WEBHOOK_SETTINGS":
		if !permissions.Admin {
			return
		}
		var s WebHookSettings = WebHookSettings{}
		s.Id = int32(id)
		ok = s.deleteWebHookSettings(enterpriseId)
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
		if !permissions.Sales {
			return
		}
		var query SaleOrderLocateQuery
		json.Unmarshal([]byte(message), &query)
		data, _ = json.Marshal(query.locateSaleOrder(enterpriseId))
	case "DOCUMENT_CONTAINER":
		if !permissions.Masters {
			return
		}
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
	case "LOCATE_ACCOUNT_SALES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateAccountForSales(enterpriseId))
	case "LOCATE_ACCOUNT_PURCHASES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateAccountForPurchases(enterpriseId))
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
	case "COLOR":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateColor(enterpriseId))
	case "PRODUCT_FAMILIES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(locateProductFamilies(enterpriseId))
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
		data, _ = json.Marshal(invoiceAllSaleOrder(int64(id), enterpriseId, userId))
	case "INVOICE_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallySaleOrder(enterpriseId, userId))
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
	case "TOGGLE_COMPLEX_MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleManufactuedComplexManufacturingOrder(int64(id), userId, enterpriseId))
	case "MANUFACTURING_ORDER_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderAllSaleOrder(int64(id), userId, enterpriseId) || complexManufacturingOrderAllSaleOrder(int64(id), userId, enterpriseId))
	case "MANUFACTURING_ORDER_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var orderInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &orderInfo)
		data, _ = json.Marshal(orderInfo.manufacturingOrderPartiallySaleOrder(userId, enterpriseId) || orderInfo.complexManufacturingOrderPartiallySaleOrder(userId, enterpriseId))
	case "DELETE_SALES_ORDER_DETAIL_PACKAGED":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal([]byte(message), &salesOrderDetailPackaged)
		salesOrderDetailPackaged.enterprise = enterpriseId
		data, _ = json.Marshal(salesOrderDetailPackaged.deleteSalesOrderDetailPackaged(userId, nil))
	case "DELIVERY_NOTE_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllSaleOrder(int64(id), enterpriseId, userId, nil)
		data, _ = json.Marshal(ok)
	case "DELIVERY_NOTE_PARTIALLY_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var noteInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallySaleOrder(enterpriseId, userId))
	case "SHIPPING_SALE_ORDER":
		if !permissions.Preparation {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(generateShippingFromSaleOrder(int64(id), enterpriseId, userId))
	case "TOGGLE_SHIPPING_SENT":
		if !permissions.Preparation {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleShippingSent(int64(id), enterpriseId, userId))
	case "SET_SHIPPING_COLLECTED":
		var shippings []int64
		json.Unmarshal([]byte(message), &shippings)
		data, _ = json.Marshal(setShippingCollected(shippings, enterpriseId, userId))
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
		ok, errorCode := generatePurchaseOrdersFromNeeds(needs, enterpriseId, userId)
		ret := OkAndErrorCodeReturn{
			Ok:       ok,
			ErorCode: errorCode,
		}
		data, _ = json.Marshal(ret)
	case "DELIVERY_NOTE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllPurchaseOrder(int64(id), enterpriseId, userId)
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
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallyPurchaseOrder(enterpriseId, userId))
	case "INVOICE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllPurchaseOrder(int64(id), enterpriseId, userId))
	case "INVOICE_PARTIAL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallyPurchaseOrder(enterpriseId, userId))
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
		data, _ = json.Marshal(p.updateProduct(userId))
	case "EMAIL":
		if !(permissions.Sales || permissions.Purchases) {
			return
		}
		var emailInfo EmailInfo
		json.Unmarshal([]byte(message), &emailInfo)
		data, _ = json.Marshal(emailInfo.sendEmail(enterpriseId))
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
		data, _ = json.Marshal(calculateMinimumStock(enterpriseId, userId))
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
		data, _ = json.Marshal(salesPostInvoices(invoiceIds, enterpriseId, userId))
	case "PURCHASE_POST_INVOICES":
		if !permissions.Accounting {
			return
		}
		var invoiceIds []int64
		json.Unmarshal([]byte(message), &invoiceIds)
		data, _ = json.Marshal(purchasePostInvoices(invoiceIds, enterpriseId, userId))
	case "MANUFACTURING_ORDER_TAG_PRINTED":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderTagPrinted(int64(id), userId, enterpriseId))
	case "COMPLEX_MANUFACTURING_ORDER_TAG_PRINTED":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(complexManufacturingOrderTagPrinted(int64(id), userId, enterpriseId))
	case "CANCEL_SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(cancelSalesOrderDetail(int64(id), enterpriseId, userId))
	case "CANCEL_PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(cancelPurchaseOrderDetail(int64(id), enterpriseId, userId))
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
		data, _ = json.Marshal(productGenerator.productGenerator(enterpriseId, userId))
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
		data, _ = json.Marshal(toggleSimplifiedInvoiceSalesInvoice(int64(id), enterpriseId, userId))
	case "MAKE_AMENDING_SALE_INVOICE":
		if !permissions.Sales {
			return
		}
		var makeAmendingInvoice MakeAmendingInvoice
		json.Unmarshal([]byte(message), &makeAmendingInvoice)
		data, _ = json.Marshal(makeAmendingSaleInvoice(makeAmendingInvoice.InvoiceId, enterpriseId, makeAmendingInvoice.Quantity, makeAmendingInvoice.Description, userId))
	case "MAKE_AMENDING_PURCHASE_INVOICE":
		if !permissions.Sales {
			return
		}
		var makeAmendingInvoice MakeAmendingInvoice
		json.Unmarshal([]byte(message), &makeAmendingInvoice)
		data, _ = json.Marshal(makeAmendingPurchaseInvoice(makeAmendingInvoice.InvoiceId, enterpriseId, makeAmendingInvoice.Quantity, makeAmendingInvoice.Description, userId))
	case "UPDATE_CURRENCY_DATA":
		if !permissions.Masters {
			return
		}
		updateCurrencyExchange(enterpriseId)
	case "SET_DIGITAL_SALES_ORDER_DETAIL_AS_SENT":
		if !permissions.Sales {
			return
		}
		var dat SetDigitalSalesOrderDetailAsSent
		json.Unmarshal([]byte(message), &dat)
		data, _ = json.Marshal(dat.setDigitalSalesOrderDetailAsSent(enterpriseId))
	case "GET_ENTERPRISE_LOGO":
		var dat map[string]string = make(map[string]string)
		logo, mimeType := getEnterpriseLogo(enterpriseId)
		dat["base64"] = base64.StdEncoding.EncodeToString(logo)
		dat["mimeType"] = mimeType
		data, _ = json.Marshal(dat)
	case "SET_ENTERPRISE_LOGO":
		var dat map[string]string = make(map[string]string)
		json.Unmarshal([]byte(message), &dat)
		logobase64, ok := dat["base64"]
		if !ok {
			return
		}
		logo, err := base64.StdEncoding.DecodeString(logobase64)
		if err != nil {
			return
		}
		data, _ = json.Marshal(setEnterpriseLogo(enterpriseId, logo))
	case "DELETE_ENTERPRISE_LOGO":
		data, _ = json.Marshal(deleteEnterpriseLogo(enterpriseId))
	case "POS_TERMINAL_REQUEST":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(posTerminalRequest(message, enterpriseId))
	case "POS_INSERT_NEW_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(posInsertNewSaleOrder(message, enterpriseId, userId))
	case "POS_SERVE_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(posServeSaleOrder(int64(id), enterpriseId, userId))
	case "POS_INSERT_NEW_SALE_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var info InsertNewSaleOrderDetail
		json.Unmarshal([]byte(message), &info)
		data, _ = json.Marshal(info.posInsertNewSaleOrderDetail(enterpriseId, userId))
	case "GET_EMPTY_API_KEY_PERMISSIONS_OBJECT":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(ApiKeyPermissions{})
	case "VAT_NUMBER_CHECK":
		if !permissions.Masters {
			return
		}
		var check CheckVatNumber
		json.Unmarshal([]byte(message), &check)
		if !check.isValid() {
			data, _ = json.Marshal(OkAndErrorCodeReturn{Ok: false})
		} else {
			data, _ = json.Marshal(checkVatNumber(check.CountryIsoCode2, check.VATNumber))
		}
	case "FINISH_INVENTORY":
		if !permissions.Warehouse {
			return
		}
		var i Inventory
		json.Unmarshal([]byte(message), &i)
		ok := i.finishInventory(userId, enterpriseId)
		data, _ = json.Marshal(ok)
	case "INSERT_UPDATE_DELETE_INVENTORY_PRODUCTS":
		if !permissions.Warehouse {
			return
		}
		var i InputInventoryProducts
		json.Unmarshal([]byte(message), &i)
		ok := i.insertUpdateDeleteInventoryProducts(enterpriseId)
		data, _ = json.Marshal(ok)
	case "INSERT_PRODUCT_FAMILY_INVENTORY_PRODUCTS":
		if !permissions.Warehouse {
			return
		}
		var i InputInventoryProducts
		json.Unmarshal([]byte(message), &i)
		ok := i.insertProductFamilyInventoryProducts(enterpriseId)
		data, _ = json.Marshal(ok)
	case "INSERT_ALL_PRODUCTS_INVENTORY_PRODUCTS":
		if !permissions.Warehouse {
			return
		}
		var i InputInventoryProducts
		json.Unmarshal([]byte(message), &i)
		ok := i.insertAllProductsInventoryProducts(enterpriseId)
		data, _ = json.Marshal(ok)
	case "DELETE_ALL_PRODUCTS_INVENTORY_PRODUCTS":
		if !permissions.Warehouse {
			return
		}
		var i InputInventoryProducts
		json.Unmarshal([]byte(message), &i)
		ok := i.deleteAllProductsInventoryProducts(enterpriseId)
		data, _ = json.Marshal(ok)
	case "WEBHOOK_SETTINGS_RENEW_AUTH_TOKEN":
		if !permissions.Admin {
			return
		}
		var s WebHookSettings
		json.Unmarshal([]byte(message), &s)
		ok := s.renewAuthToken(enterpriseId)
		data, _ = json.Marshal(ok)
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
		paginatedSearch.enterprise = enterpriseId
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
		salesOrderSearch.enterprise = enterpriseId
		data, _ = json.Marshal(salesOrderSearch.searchSalesOrder())
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchSalesInvoices())
	case "SALES_DELIVERY_NOTE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchSalesDelvieryNotes())
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchPurchaseOrder())
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.enterprise = enterpriseId
		data, _ = json.Marshal(orderSearch.searchPurchaseInvoice())
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		orderSearch.enterprise = enterpriseId
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
		paginatedSearch.enterprise = enterpriseId
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
		warehouseMovement.enterprise = enterpriseId
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
