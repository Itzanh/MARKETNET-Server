package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func main() {
	var ok bool
	settings, ok = getBackendSettings()
	if !ok {
		fmt.Println("ERROR READING SETTINGS FILE")
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", settings.Db.Host, settings.Db.Port, settings.Db.User, settings.Db.Password, settings.Db.Dbname)
	db, _ = sql.Open("postgres", psqlInfo) // control error
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server ready! :D")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", reverse)
	http.HandleFunc("/document", handleDocument)
	http.HandleFunc("/report", generateReport)
	http.HandleFunc("/export", handleExport)
	if settings.Server.TLS.UseTLS {
		go http.ListenAndServeTLS(":"+strconv.Itoa(int(settings.Server.Port)), settings.Server.TLS.CrtPath, settings.Server.TLS.KeyPath, nil)
	} else {
		go http.ListenAndServe(":"+strconv.Itoa(int(settings.Server.Port)), nil)
	}

	initialData()
	go cleanDocumentTokens()

	s := getSettingsRecord()
	c := cron.New()
	if s.Currency != "_" {
		c.AddFunc(s.CronCurrency, updateCurrencyExchange)
	}
	if s.Ecommerce == "P" {
		c.AddFunc(s.CronPrestaShop, importFromPrestaShop)
	}
	c.Start()
	c.Run()

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
	s := getSettingsRecord()
	if s.MaxConnections > 0 && len(connections) >= int(s.MaxConnections) {
		ws.Close()
		return
	}

	// AUTHENTICATION
	ok, userId, permissions := authentication(ws, r.RemoteAddr)
	if !ok {
		return
	}
	// END AUTHENTICATION
	c := Connection{Address: r.RemoteAddr, User: userId, ws: ws}
	c.addConnection()

	for {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
			c.deleteConnection()
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

		commandProcessor(command[0:commandSeparatorIndex], command[commandSeparatorIndex+1:], message[separatorIndex+1:], mt, ws, permissions)
	}
}

func authentication(ws *websocket.Conn, remoteAddr string) (bool, int16, Permissions) {
	var userId int16
	// AUTHENTICATION
	for i := 0; i < 3; i++ {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
			return false, 0, Permissions{}
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
			result.Ok, result.Permissions, userId = t.checkLoginToken()
			if result.Ok && userId >= 0 {
				result.Language = getUserRow(userId).Language
			}
		} else {
			result, userId = userLogin.login(remoteAddr)
		}

		// Return result to client (Ok + Token)
		data, _ := json.Marshal(result)
		ws.WriteMessage(mt, data)
		if result.Ok {
			return true, userId, result.Permissions
		}
	}
	// END AUTHENTICATION
	return false, 0, Permissions{}
}

func commandProcessor(instruction string, command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions) {
	switch instruction {
	case "GET":
		instructionGet(command, string(message), mt, ws, permissions)
	case "INSERT":
		instructionInsert(command, message, mt, ws, permissions)
	case "UPDATE":
		instructionUpdate(command, message, mt, ws, permissions)
	case "DELETE":
		instructionDelete(command, string(message), mt, ws, permissions)
	case "NAME":
		instructionName(command, string(message), mt, ws)
	case "GETNAME":
		instructionGetName(command, string(message), mt, ws)
	case "DEFAULTS":
		instructionDefaults(command, string(message), mt, ws, permissions)
	case "LOCATE":
		instructionLocate(command, string(message), mt, ws, permissions)
	case "ACTION":
		instructionAction(command, string(message), mt, ws, permissions)
	case "SEARCH":
		instructionSearch(command, string(message), mt, ws, permissions)
	}
}

func instructionGet(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
	var found bool = true
	var data []byte

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			data, _ = json.Marshal(getAddresses())
		case "PRODUCT":
			data, _ = json.Marshal(getProduct())
		case "PRODUCT_FAMILY":
			data, _ = json.Marshal(getProductFamilies())
		case "BILLING_SERIE":
			data, _ = json.Marshal(getBillingSeries())
		case "CURRENCY":
			data, _ = json.Marshal(getCurrencies())
		case "PAYMENT_METHOD":
			data, _ = json.Marshal(getPaymentMethods())
		case "LANGUAGE":
			data, _ = json.Marshal(getLanguages())
		case "COUNTRY":
			data, _ = json.Marshal(getCountries())
		case "STATE":
			data, _ = json.Marshal(getStates())
		case "CUSTOMER":
			data, _ = json.Marshal(getCustomers())
		case "COLOR":
			data, _ = json.Marshal(getColor())
		case "PACKAGES":
			data, _ = json.Marshal(getPackages())
		case "INCOTERMS":
			data, _ = json.Marshal(getIncoterm())
		case "CARRIERS":
			data, _ = json.Marshal(getCariers())
		case "SUPPLIERS":
			data, _ = json.Marshal(getSuppliers())
		case "DOCUMENT_CONTAINER":
			data, _ = json.Marshal(getDocumentContainer())
		case "DOCUMENTS":
			if message == "" {
				data, _ = json.Marshal(getDocuments())
			} else {
				var document Document
				json.Unmarshal([]byte(message), &document)
				data, _ = json.Marshal(document.getDocumentsRelations())
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
	}

	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrder())
	case "SALES_ORDER_PREPARATION":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getSalesOrderPreparation())
	case "SALES_ORDER_AWAITING_SHIPPING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getSalesOrderAwaitingShipping())
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getWarehouses())
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesInvoices())
	case "MANUFACTURING_ORDER_TYPE":
		if (!permissions.Manufacturing) && (!permissions.Masters) {
			return
		}
		data, _ = json.Marshal(getManufacturingOrderType())
	case "WAREHOUSE_MOVEMENTS":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getWarehouseMovement())
	case "WAREHOUSE_WAREHOUSE_MOVEMENTS":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementByWarehouse(message))
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesDeliveryNotes())
	case "SHIPPINGS":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getShippings())
	case "SHIPPING_NOT_COLLECTED":
		data, _ = json.Marshal(getShippingsPendingCollected())
	case "USERS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getUser())
	case "GROUPS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getGroup())
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrder())
	case "NEEDS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getNeeds())
	case "PURCHASE_DELIVERY_NOTES":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNotes())
	case "PURCHASE_INVOICES":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoices())
	case "SETTINGS":
		data, _ = json.Marshal(getSettingsRecord())
	case "PS_ZONES":
		if !permissions.PrestaShop {
			return
		}
		data, _ = json.Marshal(getPSZones())
	case "TABLES":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getTableAndFieldInfo())
	case "CONNECTIONS":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(getConnections())
	case "JOURNALS":
		data, _ = json.Marshal(getJournals())
	case "ACCOUNTS":
		data, _ = json.Marshal(getAccounts())
	case "ACCOUNTING_MOVEMENTS":
		data, _ = json.Marshal(getAccountingMovement())
	case "CONFIG_ACCOUNTS_VAT":
		data, _ = json.Marshal(getConfigAccountsVat())
	case "PENDING_COLLECTION_OPERATIONS":
		data, _ = json.Marshal(getPendingColletionOperations())
	case "PENDING_PAYMENT_TRANSACTIONS":
		data, _ = json.Marshal(getPendingPaymentTransaction())
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
	case "MANUFACTURING_ORDER": // accepts the "0" value
		if !permissions.Manufacturing {
			return
		}
		data, _ = json.Marshal(getManufacturingOrder(int16(id)))
		found = true
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
		data, _ = json.Marshal(getSalesOrderDetail(int32(id)))
	case "STOCK":
		data, _ = json.Marshal(getStock(int32(id)))
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderDiscounts(int32(id)))
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceDetail(int32(id)))
	case "SALES_ORDER_PACKAGING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(getPackaging(int32(id)))
	case "SALES_DELIVERY_NOTES_DETAILS":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementBySalesDeliveryNote(int32(id)))
	case "SHIPPING_PACKAGING":
		data, _ = json.Marshal(getPackagingByShipping(int32(id)))
	case "GET_USER_GROUPS":
		data, _ = json.Marshal(getUserGroups(int16(id)))
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderDetail(int32(id)))
	case "PURCHASE_DELIVERY_NOTES_DETAILS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getWarehouseMovementByPurchaseDeliveryNote(int32(id)))
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceDetail(int32(id)))
	case "PRODUCT_SALES_ORDER_PENDING":
		data, _ = json.Marshal(getProductSalesOrderDetailsPending(int32(id)))
	case "PRODUCT_PURCHASE_ORDER_PENDING":
		data, _ = json.Marshal(getProductPurchaseOrderDetailsPending(int32(id)))
	case "PRODUCT_SALES_ORDER":
		data, _ = json.Marshal(getProductSalesOrderDetails(int32(id)))
	case "PRODUCT_PURCHASE_ORDER":
		data, _ = json.Marshal(getProductPurchaseOrderDetails(int32(id)))
	case "PRODUCT_WAREHOUSE_MOVEMENT":
		data, _ = json.Marshal(getProductWarehouseMovement(int32(id)))
	case "SALES_ORDER_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesOrderRow(int32(id)))
	case "SALES_INVOICE_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceRow(int32(id)))
	case "PURCHASE_ORDER_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderRow(int32(id)))
	case "PURCHASE_INVOICE_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceRow(int32(id)))
	case "PRODUCT_IMAGE":
		data, _ = json.Marshal(getProductImages(int32(id)))
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
		data, _ = json.Marshal(getSalesOrderPallets(int32(id)))
	case "CUSTOMER_ADDRESSES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getCustomerAddresses(int32(id)))
	case "CUSTOMER_SALE_ORDERS":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getCustomerSaleOrders(int32(id)))
	case "SALES_DELIVERY_NOTE_ROW":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getSalesDeliveryNoteRow(int32(id)))
	case "PURCHASE_DELIVERY_NOTE_ROW":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNoteRow(int32(id)))
	case "ACCOUNTING_MOVEMENT_DETAILS":
		data, _ = json.Marshal(getAccountingMovementDetail(int64(id)))
	case "ACCOUNTING_MOVEMENT_SALE_INVOICES":
		data, _ = json.Marshal(getAccountingMovementSaleInvoices(int64(id)))
	case "ACCOUNTING_MOVEMENT_COLLECTION_OPERATION":
		data, _ = json.Marshal(getColletionOperations(int64(id)))
	case "COLLECTION_OPERATION_CHARGES":
		data, _ = json.Marshal(getCharges(int32(id)))
	case "ACCOUNTING_MOVEMENT_PAYMENT_TRANSACTIONS":
		data, _ = json.Marshal(getPaymentTransactions(int64(id)))
	case "PAYMENT_TRANSACTION_PAYMENTS":
		data, _ = json.Marshal(getPayments(int32(id)))
	case "ACCOUNTING_MOVEMENT_PURCHASE_INVOICES":
		data, _ = json.Marshal(getAccountingMovementPurchaseInvoices(int64(id)))
	case "SUPPLIER_ADDRESSES":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(getSupplierAddresses(int32(id)))
	case "SUPPLIER_PURCHASE_ORDERS":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getSupplierPurchaseOrders(int32(id)))
	}
	ws.WriteMessage(mt, data)
}

func instructionInsert(command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions) {
	var ok bool

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			var address Address
			json.Unmarshal(message, &address)
			ok = address.insertAddress()
		case "BILLING_SERIE":
			var serie BillingSerie
			json.Unmarshal(message, &serie)
			ok = serie.insertBillingSerie()
		case "CURRENCY":
			var currency Currency
			json.Unmarshal(message, &currency)
			ok = currency.insertCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			json.Unmarshal(message, &paymentMethod)
			ok = paymentMethod.insertPaymentMethod()
		case "LANGUAGE":
			var language Language
			json.Unmarshal(message, &language)
			ok = language.insertLanguage()
		case "COUNTRY":
			var country Country
			json.Unmarshal(message, &country)
			ok = country.insertCountry()
		case "STATE":
			var state State
			json.Unmarshal(message, &state)
			ok = state.insertState()
		case "CUSTOMER":
			var customer Customer
			json.Unmarshal(message, &customer)
			ok = customer.insertCustomer()
		case "PRODUCT":
			var product Product
			json.Unmarshal(message, &product)
			ok = product.insertProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			json.Unmarshal(message, &productFamily)
			ok = productFamily.insertProductFamily()
		case "COLOR":
			var color Color
			json.Unmarshal(message, &color)
			ok = color.insertColor()
		case "PACKAGES":
			var packages Packages
			json.Unmarshal(message, &packages)
			ok = packages.insertPackage()
		case "INCOTERM":
			var incoterm Incoterm
			json.Unmarshal(message, &incoterm)
			ok = incoterm.insertIncoterm()
		case "CARRIER":
			var carrier Carrier
			json.Unmarshal(message, &carrier)
			ok = carrier.insertCarrier()
		case "SHIPPING":
			var shipping Shipping
			json.Unmarshal(message, &shipping)
			ok, _ = shipping.insertShipping()
		case "SUPPLIER":
			var supplier Supplier
			json.Unmarshal(message, &supplier)
			ok = supplier.insertSupplier()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			json.Unmarshal(message, &documentContainer)
			ok = documentContainer.insertDocumentContainer()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			json.Unmarshal(message, &productImage)
			ok = productImage.insertProductImage()
		}
	}

	switch command {
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		ok = warehouse.insertWarehouse()
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal(message, &saleOrder)
		ok, _ = saleOrder.insertSalesOrder()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		json.Unmarshal(message, &saleOrderDetail)
		ok = saleOrderDetail.insertSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		json.Unmarshal(message, &saleOrderDiscount)
		ok = saleOrderDiscount.insertSalesOrderDiscount()
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var saleInvoice SalesInvoice
		json.Unmarshal(message, &saleInvoice)
		ok, _ = saleInvoice.insertSalesInvoice()
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		json.Unmarshal(message, &salesInvoiceDetail)
		ok = salesInvoiceDetail.insertSalesInvoiceDetail(true)
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(message, &manufacturingOrder)
		ok = manufacturingOrder.insertManufacturingOrder()
	case "SALES_ORDER_PACKAGING":
		if !permissions.Preparation {
			return
		}
		var packaging Packaging
		json.Unmarshal(message, &packaging)
		ok = packaging.insertPackaging()
	case "SALES_ORDER_DETAIL_PACKAGED":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackaged()
	case "SALES_ORDER_DETAIL_PACKAGED_EAN13":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackagedEAN13
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackagedEAN13()
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		json.Unmarshal(message, &warehouseMovement)
		ok = warehouseMovement.insertWarehouseMovement()
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		json.Unmarshal(message, &salesDeliveryNote)
		ok, _ = salesDeliveryNote.insertSalesDeliveryNotes()
	case "USER":
		if !permissions.Admin {
			return
		}
		var userInsert UserInsert
		json.Unmarshal(message, &userInsert)
		ok = userInsert.insertUser()
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		json.Unmarshal(message, &group)
		ok = group.insertGroup()
	case "USER_GROUP":
		if !permissions.Admin {
			return
		}
		var userGroup UserGroup
		json.Unmarshal(message, &userGroup)
		ok = userGroup.insertUserGroup()
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		json.Unmarshal(message, &purchaseOrder)
		ok, _ = purchaseOrder.insertPurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		ok, _ = purchaseOrderDetail.insertPurchaseOrderDetail(true)
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		json.Unmarshal(message, &purchaseDeliveryNote)
		ok, _ = purchaseDeliveryNote.insertPurchaseDeliveryNotes()
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoice PurchaseInvoice
		json.Unmarshal(message, &purchaseInvoice)
		ok, _ = purchaseInvoice.insertPurchaseInvoice()
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		json.Unmarshal(message, &purchaseInvoiceDetail)
		ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(true)
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		json.Unmarshal(message, &pallet)
		ok = pallet.insertPallet()
	case "JOURNAL":
		var journal Journal
		json.Unmarshal(message, &journal)
		ok = journal.insertJournal()
	case "ACCOUNT":
		var account Account
		json.Unmarshal(message, &account)
		ok = account.insertAccount()
	case "ACCOUNTING_MOVEMENT":
		var accountingMovement AccountingMovement
		json.Unmarshal(message, &accountingMovement)
		ok = accountingMovement.insertAccountingMovement()
	case "ACCOUNTING_MOVEMENT_DETAIL":
		var accountingMovementDetail AccountingMovementDetail
		json.Unmarshal(message, &accountingMovementDetail)
		ok = accountingMovementDetail.insertAccountingMovementDetail()
	case "CONFIG_ACCOUNTS_VAT":
		var configAccountsVat ConfigAccountsVat
		json.Unmarshal(message, &configAccountsVat)
		ok = configAccountsVat.insertConfigAccountsVat()
	case "CHARGES":
		var charges Charges
		json.Unmarshal(message, &charges)
		ok = charges.insertCharges()
	case "PAYMENT":
		var payment Payment
		json.Unmarshal(message, &payment)
		ok = payment.insertPayment()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionUpdate(command string, message []byte, mt int, ws *websocket.Conn, permissions Permissions) {
	var ok bool

	if permissions.Masters {
		switch command {
		case "ADDRESS":
			var address Address
			json.Unmarshal(message, &address)
			ok = address.updateAddress()
		case "BILLING_SERIE":
			var serie BillingSerie
			json.Unmarshal(message, &serie)
			ok = serie.updateBillingSerie()
		case "CURRENCY":
			var currency Currency
			json.Unmarshal(message, &currency)
			ok = currency.updateCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			json.Unmarshal(message, &paymentMethod)
			ok = paymentMethod.updatePaymentMethod()
		case "LANGUAGE":
			var language Language
			json.Unmarshal(message, &language)
			ok = language.updateLanguage()
		case "COUNTRY":
			var country Country
			json.Unmarshal(message, &country)
			ok = country.updateCountry()
		case "STATE":
			var city State
			json.Unmarshal(message, &city)
			ok = city.updateState()
		case "CUSTOMER":
			var customer Customer
			json.Unmarshal(message, &customer)
			ok = customer.updateCustomer()
		case "PRODUCT":
			var product Product
			json.Unmarshal(message, &product)
			ok = product.updateProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			json.Unmarshal(message, &productFamily)
			ok = productFamily.updateProductFamily()
		case "COLOR":
			var color Color
			json.Unmarshal(message, &color)
			ok = color.updateColor()
		case "PACKAGES":
			var packages Packages
			json.Unmarshal(message, &packages)
			ok = packages.updatePackage()
		case "INCOTERM":
			var incoterm Incoterm
			json.Unmarshal(message, &incoterm)
			ok = incoterm.updateIncoterm()
		case "CARRIER":
			var incoterm Carrier
			json.Unmarshal(message, &incoterm)
			ok = incoterm.updateCarrier()
		case "SUPPLIER":
			var supplier Supplier
			json.Unmarshal(message, &supplier)
			ok = supplier.updateSupplier()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			json.Unmarshal(message, &documentContainer)
			ok = documentContainer.updateDocumentContainer()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			json.Unmarshal(message, &productImage)
			ok = productImage.updateProductImage()
		}
	}

	switch command {
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		ok = warehouse.updateWarehouse()
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal(message, &saleOrder)
		ok = saleOrder.updateSalesOrder()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(message, &salesOrderDetail)
		ok = salesOrderDetail.updateSalesOrderDetail()
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		ok = manufacturingOrderType.updateManufacturingOrderType()
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		var shipping Shipping
		json.Unmarshal(message, &shipping)
		ok = shipping.updateShipping()
	case "USER":
		if !permissions.Admin {
			return
		}
		var user User
		json.Unmarshal(message, &user)
		ok = user.updateUser()
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		json.Unmarshal(message, &group)
		ok = group.updateGroup()
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var PurchaseOrdep PurchaseOrder
		json.Unmarshal(message, &PurchaseOrdep)
		ok = PurchaseOrdep.updatePurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		ok = purchaseOrderDetail.updatePurchaseOrderDetail()
	case "SETTINGS":
		var settings Settings
		json.Unmarshal(message, &settings)
		ok = settings.updateSettingsRecord()
	case "PS_ZONES":
		if !permissions.PrestaShop {
			return
		}
		var zone PSZoneWeb
		json.Unmarshal(message, &zone)
		ok = zone.updatePSZoneWeb()
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		json.Unmarshal(message, &pallet)
		ok = pallet.updatePallet()
	case "JOURNAL":
		var journal Journal
		json.Unmarshal(message, &journal)
		ok = journal.updateJournal()
	case "ACCOUNT":
		var account Account
		json.Unmarshal(message, &account)
		ok = account.updateAccount()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionDelete(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
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
		ok = serie.deleteBillingSerie()
	case "WAREHOUSE":
		if !permissions.Warehouse {
			return
		}
		var warehouse Warehouse
		warehouse.Id = message
		ok = warehouse.deleteWarehouse()
	case "USER_GROUP":
		var userGroup UserGroup
		json.Unmarshal([]byte(message), &userGroup)
		ok = userGroup.deleteUserGroup()
	case "CONFIG_ACCOUNTS_VAT":
		id, err := strconv.ParseFloat(message, 32)
		if err != nil || id < 0 {
			return
		}
		var configAccountsVat ConfigAccountsVat
		configAccountsVat.VatPercent = float32(id)
		ok = configAccountsVat.deleteConfigAccountsVat()
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
			ok = address.deleteAddress()
		case "CURRENCY":
			var currency Currency
			currency.Id = int16(id)
			ok = currency.deleteCurrency()
		case "PAYMENT_METHOD":
			var paymentMethod PaymentMethod
			paymentMethod.Id = int16(id)
			ok = paymentMethod.deletePaymentMethod()
		case "LANGUAGE":
			var language Language
			language.Id = int16(id)
			ok = language.deleteLanguage()
		case "COUNTRY":
			var country Country
			country.Id = int16(id)
			ok = country.deleteCountry()
		case "STATE":
			var city State
			city.Id = int32(id)
			ok = city.deleteState()
		case "CUSTOMER":
			var customer Customer
			customer.Id = int32(id)
			ok = customer.deleteCustomer()
		case "PRODUCT":
			var product Product
			product.Id = int32(id)
			ok = product.deleteProduct()
		case "PRODUCT_FAMILY":
			var productFamily ProductFamily
			productFamily.Id = int16(id)
			ok = productFamily.deleteProductFamily()
		case "COLOR":
			var color Color
			color.Id = int16(id)
			ok = color.deleteColor()
		case "PACKAGES":
			var packages Packages
			packages.Id = int16(id)
			ok = packages.deletePackage()
		case "INCOTERM":
			var incoterm Incoterm
			incoterm.Id = int16(id)
			ok = incoterm.deleteIncoterm()
		case "CARRIER":
			var carrier Carrier
			carrier.Id = int16(id)
			ok = carrier.deleteCarrier()
		case "SUPPLIER":
			var supplier Supplier
			supplier.Id = int32(id)
			ok = supplier.deleteSupplier()
		case "DOCUMENT_CONTAINER":
			var documentContainer DocumentContainer
			documentContainer.Id = int16(id)
			ok = documentContainer.deleteDocumentContainer()
		case "DOCUMENT":
			var document Document
			document.Id = int32(id)
			ok = document.deleteDocument()
		case "PRODUCT_IMAGE":
			var productImage ProductImage
			productImage.Id = int32(id)
			ok = productImage.deleteProductImage()
		}
	}

	switch command {
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int32(id)
		ok = saleOrder.deleteSalesOrder()
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int32(id)
		ok = saleOrderDetail.deleteSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		if !permissions.Sales {
			return
		}
		var saleOrderDiscount SalesOrderDiscount
		saleOrderDiscount.Id = int32(id)
		ok = saleOrderDiscount.deleteSalesOrderDiscount()
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int32(id)
		ok = salesInvoice.deleteSalesInvoice()
	case "SALES_INVOICE_DETAIL":
		if !permissions.Sales {
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int32(id)
		ok = salesInvoiceDetail.deleteSalesInvoiceDetail()
	case "MANUFACTURING_ORDER_TYPE":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int16(id)
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		ok = manufacturingOrder.deleteManufacturingOrder()
	case "PACKAGING":
		if !permissions.Preparation {
			return
		}
		var packaging Packaging
		packaging.Id = int32(id)
		ok = packaging.deletePackaging()
	case "WAREHOUSE_MOVEMENTS":
		if !(permissions.Sales || permissions.Purchases || permissions.Warehouse) {
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		ok = warehouseMovement.deleteWarehouseMovement()
	case "SALES_DELIVERY_NOTES":
		if !permissions.Sales {
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int32(id)
		ok = salesDeliveryNote.deleteSalesDeliveryNotes()
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		var shipping Shipping
		shipping.Id = int32(id)
		ok = shipping.deleteShipping()
	case "USER":
		if !permissions.Admin {
			return
		}
		var user User
		user.Id = int16(id)
		ok = user.deleteUser()
	case "GROUP":
		if !permissions.Admin {
			return
		}
		var group Group
		group.Id = int16(id)
		ok = group.deleteGroup()
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int32(id)
		ok = purchaseOrder.deletePurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int32(id)
		ok = purchaseOrderDetail.deletePurchaseOrderDetail()
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int32(id)
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes()
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int32(id)
		ok = purchaseInvoice.deletePurchaseInvoice()
	case "PURCHASE_INVOICE_DETAIL":
		if !permissions.Purchases {
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int32(id)
		ok = purchaseInvoiceDetail.deletePurchaseInvoiceDetail()
	case "PALLET":
		if !permissions.Preparation {
			return
		}
		var pallet Pallet
		pallet.Id = int32(id)
		ok = pallet.deletePallet()
	case "JOURNAL":
		var journal Journal
		journal.Id = int16(id)
		ok = journal.deleteJournal()
	case "ACCOUNT":
		var account Account
		account.Id = int32(id)
		ok = account.deleteAccount()
	case "ACCOUNTING_MOVEMENT":
		var accountingMovement AccountingMovement
		accountingMovement.Id = int64(id)
		ok = accountingMovement.deleteAccountingMovement()
	case "ACCOUNTING_MOVEMENT_DETAIL":
		var accountingMovementDetail AccountingMovementDetail
		accountingMovementDetail.Id = int64(id)
		ok = accountingMovementDetail.deleteAccountingMovementDetail()
	case "CHARGES":
		var charges Charges
		charges.Id = int32(id)
		ok = charges.deleteCharges()
	case "PAYMENT":
		var payment Payment
		payment.Id = int32(id)
		ok = payment.deletePayment()
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

func instructionName(command string, message string, mt int, ws *websocket.Conn) {
	var data []byte
	switch command {
	case "LANGUAGE":
		data, _ = json.Marshal(findLanguageByName(message))
	case "CURRENCY":
		data, _ = json.Marshal(findCurrencyByName(message))
	case "CUSTOMER":
		data, _ = json.Marshal(findCustomerByName(message))
	case "COUNTRY":
		data, _ = json.Marshal(findCountryByName(message))
	case "STATE":
		var cityName StateNameQuery
		json.Unmarshal([]byte(message), &cityName)
		data, _ = json.Marshal(findStateByName(cityName))
	case "PAYMENT_METHOD":
		data, _ = json.Marshal(findPaymentMethodByName(message))
	case "BILLING_SERIE":
		data, _ = json.Marshal(findBillingSerieByName(message))
	case "PRODUCT_FAMILY":
		data, _ = json.Marshal(findProductFamilyByName(message))
	case "COLOR":
		data, _ = json.Marshal(findColorByName(message))
	case "PRODUCT":
		data, _ = json.Marshal(findProductByName(message))
	case "WAREHOUSE":
		data, _ = json.Marshal(findWarehouseByName(message))
	case "CARRIER":
		data, _ = json.Marshal(findCarrierByName(message))
	case "SUPPLIER":
		data, _ = json.Marshal(findSupplierByName(message))
	}
	ws.WriteMessage(mt, data)
}

func instructionGetName(command string, message string, mt int, ws *websocket.Conn) {
	// attempt to get the name from resources with alpha key, if the resource if not found,
	// parse the input as number and attemp to get the name resource with numeric key
	var found bool = true
	var name string
	// ALPHA
	switch command {
	case "BILLING_SERIE":
		name = getNameBillingSerie(message)
	case "WAREHOUSE":
		name = getNameWarehouse(message)
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
		name = getNameLanguage(int16(id))
	case "CURRENCY":
		name = getNameCurrency(int16(id))
	case "CUSTOMER":
		name = getNameCustomer(int32(id))
	case "COUNTRY":
		name = getNameCountry(int16(id))
	case "STATE":
		name = getNameState(int32(id))
	case "PAYMENT_METHOD":
		name = getNamePaymentMethod(int16(id))
	case "PRODUCT_FAMILY":
		name = getNameProductFamily(int16(id))
	case "COLOR":
		name = getNameColor(int16(id))
	case "ADDRESS":
		name = getAddressName(int32(id))
	case "PRODUCT":
		name = getNameProduct(int32(id))
	case "CARRIER":
		name = getNameCarrier(int16(id))
	case "SALE_DELIERY_NOTE":
		name = getNameSalesDeliveryNote(int32(id))
	case "SUPPLIER":
		name = getNameSupplier(int32(id))
	}
	ws.WriteMessage(mt, []byte(name))
}

func instructionDefaults(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
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
		data, _ = json.Marshal(getSaleOrderDefaults())
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderDefaults())
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
		data, _ = json.Marshal(getCustomerDefaults(int32(id)))
	case "SALES_ORDER_DETAIL":
		if !permissions.Sales {
			return
		}
		data, _ = json.Marshal(getOrderDetailDefaults(int32(id)))
	case "SUPPLIER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(getSupplierDefaults(int32(id)))
	}
	ws.WriteMessage(mt, data)
}

func instructionLocate(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
	var data []byte
	var found bool = true

	// PARAMETERLESS
	switch command {
	case "SALE_ORDER":
		data, _ = json.Marshal(locateSaleOrder())
	case "DOCUMENT_CONTAINER":
		data, _ = json.Marshal(locateDocumentContainer())
	case "LOCATE_ACCOUNT_CUSTOMER":
		data, _ = json.Marshal(locateAccountForCustomer())
	case "LOCATE_ACCOUNT_SUPPLIER":
		data, _ = json.Marshal(locateAccountForSupplier())
	case "LOCATE_ACCOUNT_BANKS":
		data, _ = json.Marshal(locateAccountForBanks())
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
		data, _ = json.Marshal(locateAddressByCustomer(int32(id)))
	case "ADDRESS_SUPPLIER":
		if !permissions.Purchases {
			return
		}
		data, _ = json.Marshal(locateAddressBySupplier(int32(id)))
	case "SALE_DELIVERY_NOTE":
		data, _ = json.Marshal(locateSalesDeliveryNotesBySalesOrder(int32(id)))
	}
	ws.WriteMessage(mt, data)
}

func instructionAction(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
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
		data, _ = json.Marshal(invoiceAllSaleOrder(int32(id)))
	case "INVOICE_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallySaleOrder())
	case "GET_SALES_ORDER_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesOrderRelations(int32(id)))
	case "GET_SALES_INVOICE_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceRelations(int32(id)))
	case "TOGGLE_MANUFACTURING_ORDER":
		if !permissions.Manufacturing {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleManufactuedManufacturingOrder(int64(id)))
	case "MANUFACTURING_ORDER_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderAllSaleOrder(int32(id)))
	case "MANUFACTURING_ORDER_PARTIAL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var orderInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &orderInfo)
		data, _ = json.Marshal(orderInfo.manufacturingOrderPartiallySaleOrder())
	case "DELETE_SALES_ORDER_DETAIL_PACKAGED":
		if !permissions.Preparation {
			return
		}
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal([]byte(message), &salesOrderDetailPackaged)
		data, _ = json.Marshal(salesOrderDetailPackaged.deleteSalesOrderDetailPackaged(true))
	case "DELIVERY_NOTE_ALL_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllSaleOrder(int32(id))
		data, _ = json.Marshal(ok)
	case "DELIVERY_NOTE_PARTIALLY_SALE_ORDER":
		if !permissions.Sales {
			return
		}
		var noteInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallySaleOrder())
	case "SHIPPING_SALE_ORDER":
		if !permissions.Preparation {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(generateShippingFromSaleOrder(int32(id)))
	case "TOGGLE_SHIPPING_SENT":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleShippingSent(int32(id)))
	case "SET_SHIPPING_COLLECTED":
		var shippings []int32
		json.Unmarshal([]byte(message), &shippings)
		data, _ = json.Marshal(setShippingCollected(shippings))
	case "GET_SALES_DELIVERY_NOTE_RELATIONS":
		if !permissions.Sales {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesDeliveryNoteRelations(int32(id)))
	case "USER_PWD":
		var userPassword UserPassword
		json.Unmarshal([]byte(message), &userPassword)
		data, _ = json.Marshal(userPassword.userPassword())
	case "USER_OFF":
		var user User
		json.Unmarshal([]byte(message), &user)
		data, _ = json.Marshal(user.offUser())
	case "PURCHASE_NEEDS":
		var needs []PurchaseNeed
		json.Unmarshal([]byte(message), &needs)
		data, _ = json.Marshal(generatePurchaseOrdersFromNeeds(needs))
	case "DELIVERY_NOTE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllPurchaseOrder(int32(id))
		data, _ = json.Marshal(ok)
	case "GET_PURCHASE_ORDER_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderRelations(int32(id)))
	case "GET_INVOICE_ORDER_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceRelations(int32(id)))
	case "GET_PURCHASE_DELIVERY_NOTE_RELATIONS":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNoteRelations(int32(id)))
	case "DELIVERY_NOTE_PARTIALLY_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var noteInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallyPurchaseOrder())
	case "INVOICE_ALL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllPurchaseOrder(int32(id)))
	case "INVOICE_PARTIAL_PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var invoiceInfo OrderDetailGenerate
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallyPurchaseOrder())
	case "INSERT_DOCUMENT":
		var document Document
		json.Unmarshal([]byte(message), &document)
		ok := document.insertDocument()
		if ok {
			data, _ = json.Marshal(document)
		} else {
			data, _ = json.Marshal(ok)
		}
	case "GRANT_DOCUMENT_ACCESS_TOKEN":
		data, _ = json.Marshal(grantDocumentAccessToken())
	case "GET_PRODUCT_ROW":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getProductRow(int32(id)))
	case "PRODUCT_EAN13":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		p := getProductRow(int32(id))
		p.generateBarcode()
		data, _ = json.Marshal(p.updateProduct())
	case "EMAIL":
		if !(permissions.Sales || permissions.Purchases) {
			return
		}
		var emailInfo EmailInfo
		json.Unmarshal([]byte(message), &emailInfo)
		data, _ = json.Marshal(emailInfo.sendEmail())
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
		data, _ = json.Marshal(exportToJSON(message))
	case "IMPORT_JSON":
		if !permissions.Admin {
			return
		}
		var importInfo ImportInfo
		json.Unmarshal([]byte(message), &importInfo)
		data, _ = json.Marshal(importInfo.importJson())
	case "REGENERATE_DRAGGED_STOCK":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(regenerateDraggedStock(message))
	case "REGENERATE_PRODUCT_STOCK":
		if !permissions.Warehouse {
			return
		}
		data, _ = json.Marshal(regenerateProductStock())
	case "DISCONNECT":
		if !permissions.Admin {
			return
		}
		data, _ = json.Marshal(disconnectConnection(message))
	case "PRESTASHOP":
		if !permissions.PrestaShop {
			return
		}
		importFromPrestaShop()
	case "CALCULATE_MINIMUM_STOCK":
		data, _ = json.Marshal(calculateMinimumStock())
	case "GENERATE_MANUFACTURIG_OR_PURCHASE_ORDERS_MINIMUM_STOCK":
		data, _ = json.Marshal(generateManufacturingOrPurchaseOrdersMinimumStock())
	case "SALES_POST_INVOICES":
		var invoiceIds []int32
		json.Unmarshal([]byte(message), &invoiceIds)
		data, _ = json.Marshal(salesPostInvoices(invoiceIds))
	case "PURCHASE_POST_INVOICES":
		var invoiceIds []int32
		json.Unmarshal([]byte(message), &invoiceIds)
		data, _ = json.Marshal(purchasePostInvoices(invoiceIds))
	}
	ws.WriteMessage(mt, data)
}

func instructionSearch(command string, message string, mt int, ws *websocket.Conn, permissions Permissions) {
	var data []byte
	switch command {
	case "CUSTOMER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchCustomers(message))
	case "SUPPLER":
		if !permissions.Masters {
			return
		}
		data, _ = json.Marshal(searchSuppliers(message))
	case "PRODUCT":
		if !permissions.Masters {
			return
		}
		var productSearch ProductSearch
		json.Unmarshal([]byte(message), &productSearch)
		data, _ = json.Marshal(productSearch.searchProduct())
	case "SHIPPING":
		if !permissions.Preparation {
			return
		}
		data, _ = json.Marshal(searchShippings(message))
	case "SALES_ORDER":
		if !permissions.Sales {
			return
		}
		var salesOrderSearch SalesOrderSearch
		json.Unmarshal([]byte(message), &salesOrderSearch)
		data, _ = json.Marshal(salesOrderSearch.searchSalesOrder())
	case "SALES_INVOICE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchSalesInvoices())
	case "SALES_DELIVERY_NOTE":
		if !permissions.Sales {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchSalesDelvieryNotes())
	case "PURCHASE_ORDER":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseOrder())
	case "PURCHASE_INVOICE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseInvoice())
	case "PURCHASE_DELIVERY_NOTE":
		if !permissions.Purchases {
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseDeliveryNote())
	case "COUNTRY":
		data, _ = json.Marshal(searchCountries(message))
	case "STATE":
		data, _ = json.Marshal(searchStates(message))
	case "ADDRESS":
		data, _ = json.Marshal(searchAddresses(message))
	case "LANGUAGE":
		data, _ = json.Marshal(searchLanguages(message))
	case "WAREHOUSE_MOVEMENT":
		if !permissions.Warehouse {
			return
		}
		var warehouseMovement WarehouseMovementSearch
		json.Unmarshal([]byte(message), &warehouseMovement)
		data, _ = json.Marshal(warehouseMovement.searchWarehouseMovement())
	}
	ws.WriteMessage(mt, data)
}
