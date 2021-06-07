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

var settings BackendSettings

var upgrader = websocket.Upgrader{}

var db *sql.DB

func main() {
	var ok bool
	settings, ok = getBackendSettings()
	if !ok {
		fmt.Println("ERROR READING SETTINGS FILE")
		return
	}

	fmt.Println("Server ready! :D")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", reverse)
	http.HandleFunc("/document", handleDocument)
	go http.ListenAndServe(":"+strconv.Itoa(int(settings.Server.Port)), nil)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", settings.Db.Host, settings.Db.Port, settings.Db.User, settings.Db.Password, settings.Db.Dbname)
	db, _ = sql.Open("postgres", psqlInfo)
	db.Ping()

	initialData()
	go cleanDocumentTokens()

	c := cron.New()
	c.AddFunc("30 0 * * *", updateCurrencyExchange)
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

	// AUTHENTICATION
	if !authentication(ws, r.RemoteAddr) {
		return
	}
	// END AUTHENTICATION

	for {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
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

		commandProcessor(command[0:commandSeparatorIndex], command[commandSeparatorIndex+1:], message[separatorIndex+1:], mt, ws)
	}
}

func authentication(ws *websocket.Conn, remoteAddr string) bool {
	// AUTHENTICATION
	for i := 0; i < 3; i++ {
		// Receive message
		mt, message, err := ws.ReadMessage()
		if err != nil {
			return false
		}

		// Remote the port from the address
		remoteAddr = remoteAddr[:strings.Index(remoteAddr, ":")]

		// Attempt login in DB
		var userLogin UserLogin
		json.Unmarshal(message, &userLogin)
		result := UserLoginResult{}
		if len(userLogin.Token) > 0 {
			t := LoginToken{Name: userLogin.Token, IpAddress: remoteAddr}
			result.Ok = t.checkLoginToken()
		} else {
			result = userLogin.login(remoteAddr)
		}

		// Return result to client (Ok + Token)
		data, _ := json.Marshal(result)
		ws.WriteMessage(mt, data)
		if result.Ok {
			return true
		}
	}
	// END AUTHENTICATION
	return false
}

func commandProcessor(instruction string, command string, message []byte, mt int, ws *websocket.Conn) {
	switch instruction {
	case "GET":
		instructionGet(command, string(message), mt, ws)
	case "INSERT":
		instructionInsert(command, message, mt, ws)
	case "UPDATE":
		instructionUpdate(command, message, mt, ws)
	case "DELETE":
		instructionDelete(command, string(message), mt, ws)
	case "NAME":
		instructionName(command, string(message), mt, ws)
	case "GETNAME":
		instructionGetName(command, string(message), mt, ws)
	case "DEFAULTS":
		instructionDefaults(command, string(message), mt, ws)
	case "LOCATE":
		instructionLocate(command, string(message), mt, ws)
	case "ACTION":
		instructionAction(command, string(message), mt, ws)
	case "SEARCH":
		instructionSearch(command, string(message), mt, ws)
	}
}

func instructionGet(command string, message string, mt int, ws *websocket.Conn) {
	var found bool = true
	var data []byte
	switch command {
	case "SALES_ORDER":
		data, _ = json.Marshal(getSalesOrder())
	case "SALES_ORDER_PREPARATION":
		data, _ = json.Marshal(getSalesOrderPreparation())
	case "SALES_ORDER_AWAITING_SHIPPING":
		data, _ = json.Marshal(getSalesOrderAwaitingShipping())
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
	case "WAREHOUSE":
		data, _ = json.Marshal(getWarehouses())
	case "LANGUAGE":
		data, _ = json.Marshal(getLanguages())
	case "COUNTRY":
		data, _ = json.Marshal(getCountries())
	case "CITY":
		data, _ = json.Marshal(getCities())
	case "CUSTOMER":
		data, _ = json.Marshal(getCustomers())
	case "COLOR":
		data, _ = json.Marshal(getColor())
	case "SALES_INVOICE":
		data, _ = json.Marshal(getSalesInvoices())
	case "MANUFACTURING_ORDER_TYPE":
		data, _ = json.Marshal(getManufacturingOrderType())
	case "PACKAGES":
		data, _ = json.Marshal(getPackages())
	case "WAREHOUSE_MOVEMENTS":
		data, _ = json.Marshal(getWarehouseMovement())
	case "WAREHOUSE_WAREHOUSE_MOVEMENTS":
		data, _ = json.Marshal(getWarehouseMovementByWarehouse(message))
	case "SALES_DELIVERY_NOTES":
		data, _ = json.Marshal(getSalesDeliveryNotes())
	case "INCOTERMS":
		data, _ = json.Marshal(getIncoterm())
	case "CARRIERS":
		data, _ = json.Marshal(getCariers())
	case "SHIPPINGS":
		data, _ = json.Marshal(getShippings())
	case "USERS":
		data, _ = json.Marshal(getUser())
	case "GROUPS":
		data, _ = json.Marshal(getGroup())
	case "SUPPLIERS":
		data, _ = json.Marshal(getSuppliers())
	case "PURCHASE_ORDER":
		data, _ = json.Marshal(getPurchaseOrder())
	case "NEEDS":
		data, _ = json.Marshal(getNeeds())
	case "PURCHASE_DELIVERY_NOTES":
		data, _ = json.Marshal(getPurchaseDeliveryNotes())
	case "PURCHASE_INVOICES":
		data, _ = json.Marshal(getPurchaseInvoices())
	case "SETTINGS":
		data, _ = json.Marshal(getSettingsRecord())
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
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	switch command {
	case "MANUFACTURING_ORDER": // accepts the "0" value
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
		data, _ = json.Marshal(getSalesOrderDetail(int32(id)))
	case "STOCK":
		data, _ = json.Marshal(getStock(int32(id)))
	case "SALES_ORDER_DISCOUNT":
		data, _ = json.Marshal(getSalesOrderDiscounts(int32(id)))
	case "SALES_INVOICE_DETAIL":
		data, _ = json.Marshal(getSalesInvoiceDetail(int32(id)))
	case "SALES_ORDER_PACKAGING":
		data, _ = json.Marshal(getPackaging(int32(id)))
	case "SALES_DELIVERY_NOTES_DETAILS":
		data, _ = json.Marshal(getWarehouseMovementBySalesDeliveryNote(int32(id)))
	case "SHIPPING_PACKAGING":
		data, _ = json.Marshal(getPackagingByShipping(int32(id)))
	case "GET_USER_GROUPS":
		data, _ = json.Marshal(getUserGroups(int16(id)))
	case "PURCHASE_ORDER_DETAIL":
		data, _ = json.Marshal(getPurchaseOrderDetail(int32(id)))
	case "PURCHASE_DELIVERY_NOTES_DETAILS":
		data, _ = json.Marshal(getWarehouseMovementByPurchaseDeliveryNote(int32(id)))
	case "PURCHASE_INVOICE_DETAIL":
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
	}
	ws.WriteMessage(mt, data)
}

func instructionInsert(command string, message []byte, mt int, ws *websocket.Conn) {
	var ok bool
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
	case "WAREHOUSE":
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		ok = warehouse.insertWarehouse()
	case "LANGUAGE":
		var language Language
		json.Unmarshal(message, &language)
		ok = language.insertLanguage()
	case "COUNTRY":
		var country Country
		json.Unmarshal(message, &country)
		ok = country.insertCountry()
	case "CITY":
		var city City
		json.Unmarshal(message, &city)
		ok = city.insertCity()
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
	case "SALES_ORDER":
		var saleOrder SaleOrder
		json.Unmarshal(message, &saleOrder)
		ok = saleOrder.insertSalesOrder()
	case "SALES_ORDER_DETAIL":
		var saleOrderDetail SalesOrderDetail
		json.Unmarshal(message, &saleOrderDetail)
		ok = saleOrderDetail.insertSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		var saleOrderDiscount SalesOrderDiscount
		json.Unmarshal(message, &saleOrderDiscount)
		ok = saleOrderDiscount.insertSalesOrderDiscount()
	case "SALES_INVOICE":
		var saleInvoice SalesInvoice
		json.Unmarshal(message, &saleInvoice)
		ok, _ = saleInvoice.insertSalesInvoice()
	case "SALES_INVOICE_DETAIL":
		var salesInvoiceDetail SalesInvoiceDetail
		json.Unmarshal(message, &salesInvoiceDetail)
		ok = salesInvoiceDetail.insertSalesInvoiceDetail(true)
	case "MANUFACTURING_ORDER_TYPE":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(message, &manufacturingOrder)
		ok = manufacturingOrder.insertManufacturingOrder()
	case "PACKAGES":
		var packages Packages
		json.Unmarshal(message, &packages)
		ok = packages.insertPackage()
	case "SALES_ORDER_PACKAGING":
		var packaging Packaging
		json.Unmarshal(message, &packaging)
		ok = packaging.insertPackaging()
	case "SALES_ORDER_DETAIL_PACKAGED":
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackaged()
	case "SALES_ORDER_DETAIL_PACKAGED_EAN13":
		var salesOrderDetailPackaged SalesOrderDetailPackagedEAN13
		json.Unmarshal(message, &salesOrderDetailPackaged)
		ok = salesOrderDetailPackaged.insertSalesOrderDetailPackagedEAN13()
	case "WAREHOUSE_MOVEMENTS":
		var warehouseMovement WarehouseMovement
		json.Unmarshal(message, &warehouseMovement)
		ok = warehouseMovement.insertWarehouseMovement()
	case "SALES_DELIVERY_NOTES":
		var salesDeliveryNote SalesDeliveryNote
		json.Unmarshal(message, &salesDeliveryNote)
		ok, _ = salesDeliveryNote.insertSalesDeliveryNotes()
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
	case "USER":
		var userInsert UserInsert
		json.Unmarshal(message, &userInsert)
		ok = userInsert.insertUser()
	case "GROUP":
		var group Group
		json.Unmarshal(message, &group)
		ok = group.insertGroup()
	case "USER_GROUP":
		var userGroup UserGroup
		json.Unmarshal(message, &userGroup)
		ok = userGroup.insertUserGroup()
	case "SUPPLIER":
		var supplier Supplier
		json.Unmarshal(message, &supplier)
		ok = supplier.insertSupplier()
	case "PURCHASE_ORDER":
		var purchaseOrder PurchaseOrder
		json.Unmarshal(message, &purchaseOrder)
		ok, _ = purchaseOrder.insertPurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		ok, _ = purchaseOrderDetail.insertPurchaseOrderDetail()
	case "PURCHASE_DELIVERY_NOTE":
		var purchaseDeliveryNote PurchaseDeliveryNote
		json.Unmarshal(message, &purchaseDeliveryNote)
		ok, _ = purchaseDeliveryNote.insertPurchaseDeliveryNotes()
	case "PURCHASE_INVOICE":
		var purchaseInvoice PurchaseInvoice
		json.Unmarshal(message, &purchaseInvoice)
		ok, _ = purchaseInvoice.insertPurchaseInvoice()
	case "PURCHASE_INVOICE_DETAIL":
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		json.Unmarshal(message, &purchaseInvoiceDetail)
		ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(true)
	case "DOCUMENT_CONTAINER":
		var documentContainer DocumentContainer
		json.Unmarshal(message, &documentContainer)
		ok = documentContainer.insertDocumentContainer()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionUpdate(command string, message []byte, mt int, ws *websocket.Conn) {
	var ok bool
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
	case "WAREHOUSE":
		var warehouse Warehouse
		json.Unmarshal(message, &warehouse)
		ok = warehouse.updateWarehouse()
	case "LANGUAGE":
		var language Language
		json.Unmarshal(message, &language)
		ok = language.updateLanguage()
	case "COUNTRY":
		var country Country
		json.Unmarshal(message, &country)
		ok = country.updateCountry()
	case "CITY":
		var city City
		json.Unmarshal(message, &city)
		ok = city.updateCity()
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
	case "SALES_ORDER":
		var saleOrder SaleOrder
		json.Unmarshal(message, &saleOrder)
		ok = saleOrder.updateSalesOrder()
	case "SALES_ORDER_DETAIL":
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(message, &salesOrderDetail)
		ok = salesOrderDetail.updateSalesOrderDetail()
	case "MANUFACTURING_ORDER_TYPE":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(message, &manufacturingOrderType)
		ok = manufacturingOrderType.updateManufacturingOrderType()
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
	case "SHIPPING":
		var shipping Shipping
		json.Unmarshal(message, &shipping)
		ok = shipping.updateShipping()
	case "USER":
		var user User
		json.Unmarshal(message, &user)
		ok = user.updateUser()
	case "GROUP":
		var group Group
		json.Unmarshal(message, &group)
		ok = group.updateGroup()
	case "SUPPLIER":
		var supplier Supplier
		json.Unmarshal(message, &supplier)
		ok = supplier.updateSupplier()
	case "PURCHASE_ORDER":
		var PurchaseOrdep PurchaseOrder
		json.Unmarshal(message, &PurchaseOrdep)
		ok = PurchaseOrdep.updatePurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(message, &purchaseOrderDetail)
		ok = purchaseOrderDetail.updatePurchaseOrderDetail()
	case "SETTINGS":
		var settings Settings
		json.Unmarshal(message, &settings)
		ok = settings.updateSettingsRecord()
	case "DOCUMENT_CONTAINER":
		var documentContainer DocumentContainer
		json.Unmarshal(message, &documentContainer)
		ok = documentContainer.updateDocumentContainer()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
}

func instructionDelete(command string, message string, mt int, ws *websocket.Conn) {
	// attempt to delete from resources with alpha key, if the resource if not found,
	// parse the input as number and attemp to delete resource with numeric key
	var found bool = true
	var ok bool
	// ALPHA
	switch command {
	case "BILLING_SERIE":
		var serie BillingSerie
		serie.Id = message
		ok = serie.deleteBillingSerie()
	case "WAREHOUSE":
		var warehouse Warehouse
		warehouse.Id = message
		ok = warehouse.deleteWarehouse()
	case "USER_GROUP":
		var userGroup UserGroup
		json.Unmarshal([]byte(message), &userGroup)
		ok = userGroup.deleteUserGroup()
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
	case "CITY":
		var city City
		city.Id = int32(id)
		ok = city.deleteCity()
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
	case "SALES_ORDER":
		var saleOrder SaleOrder
		saleOrder.Id = int32(id)
		ok = saleOrder.deleteSalesOrder()
	case "SALES_ORDER_DETAIL":
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int32(id)
		ok = saleOrderDetail.deleteSalesOrderDetail()
	case "SALES_ORDER_DISCOUNT":
		var saleOrderDiscount SalesOrderDiscount
		saleOrderDiscount.Id = int32(id)
		ok = saleOrderDiscount.deleteSalesOrderDiscount()
	case "SALES_INVOICE":
		var salesInvoice SalesInvoice
		salesInvoice.Id = int32(id)
		ok = salesInvoice.deleteSalesInvoice()
	case "SALES_INVOICE_DETAIL":
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int32(id)
		ok = salesInvoiceDetail.deleteSalesInvoiceDetail()
	case "MANUFACTURING_ORDER_TYPE":
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int16(id)
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	case "MANUFACTURING_ORDER":
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		ok = manufacturingOrder.deleteManufacturingOrder()
	case "PACKAGES":
		var packages Packages
		packages.Id = int16(id)
		ok = packages.deletePackage()
	case "PACKAGING":
		var packaging Packaging
		packaging.Id = int32(id)
		ok = packaging.deletePackaging()
	case "WAREHOUSE_MOVEMENTS":
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		ok = warehouseMovement.deleteWarehouseMovement()
	case "SALES_DELIVERY_NOTES":
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int32(id)
		ok = salesDeliveryNote.deleteSalesDeliveryNotes()
	case "INCOTERM":
		var incoterm Incoterm
		incoterm.Id = int16(id)
		ok = incoterm.deleteIncoterm()
	case "CARRIER":
		var carrier Carrier
		carrier.Id = int16(id)
		ok = carrier.deleteCarrier()
	case "SHIPPING":
		var shipping Shipping
		shipping.Id = int32(id)
		ok = shipping.deleteShipping()
	case "USER":
		var user User
		user.Id = int16(id)
		ok = user.deleteUser()
	case "GROUP":
		var group Group
		group.Id = int16(id)
		ok = group.deleteGroup()
	case "SUPPLIER":
		var supplier Supplier
		supplier.Id = int32(id)
		ok = supplier.deleteSupplier()
	case "PURCHASE_ORDER":
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int32(id)
		ok = purchaseOrder.deletePurchaseOrder()
	case "PURCHASE_ORDER_DETAIL":
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int32(id)
		ok = purchaseOrderDetail.deletePurchaseOrderDetail()
	case "PURCHASE_DELIVERY_NOTE":
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int32(id)
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes()
	case "PURCHASE_INVOICE":
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int32(id)
		ok = purchaseInvoice.deletePurchaseInvoice()
	case "PURCHASE_INVOICE_DETAIL":
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int32(id)
		ok = purchaseInvoiceDetail.deletePurchaseInvoiceDetail()
	case "DOCUMENT_CONTAINER":
		var documentContainer DocumentContainer
		documentContainer.Id = int16(id)
		ok = documentContainer.deleteDocumentContainer()
	case "DOCUMENT":
		var document Document
		document.Id = int32(id)
		ok = document.deleteDocument()
	}
	data, _ := json.Marshal(ok)
	ws.WriteMessage(mt, data)
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
	case "CITY":
		var cityName CityNameQuery
		json.Unmarshal([]byte(message), &cityName)
		data, _ = json.Marshal(findCityByName(cityName))
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
	case "CITY":
		name = getNameCity(int32(id))
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

func instructionDefaults(command string, message string, mt int, ws *websocket.Conn) {
	// there are defaults that require an ID of a row, and there are defaults without parametres
	// attemps first respond to the parameterless, and if not found, parse the parameters and return

	var found bool = true
	var data []byte
	// ALPHA
	switch command {
	case "SALES_ORDER":
		data, _ = json.Marshal(getSaleOrderDefaults())
	case "PURCHASE_ORDER":
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
		data, _ = json.Marshal(getCustomerDefaults(int32(id)))
	case "SALES_ORDER_DETAIL":
		data, _ = json.Marshal(getOrderDetailDefaults(int32(id)))
	case "SUPPLIER":
		data, _ = json.Marshal(getSupplierDefaults(int32(id)))
	}
	ws.WriteMessage(mt, data)
}

func instructionLocate(command string, message string, mt int, ws *websocket.Conn) {
	var data []byte
	var found bool = true

	// PARAMETERLESS
	switch command {
	case "SALE_ORDER":
		data, _ = json.Marshal(locateSaleOrder())
	case "DOCUMENT_CONTAINER":
		data, _ = json.Marshal(locateDocumentContainer())
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
		data, _ = json.Marshal(locateAddressBySupplier(int32(id)))
	case "SALE_DELIVERY_NOTE":
		data, _ = json.Marshal(locateSalesDeliveryNotesBySalesOrder(int32(id)))
	}
	ws.WriteMessage(mt, data)
}

func instructionAction(command string, message string, mt int, ws *websocket.Conn) {
	var data []byte

	switch command {
	case "INVOICE_ALL_SALE_ORDER":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllSaleOrder(int32(id)))
	case "INVOICE_PARTIAL_SALE_ORDER":
		var invoiceInfo SalesOrderDetailInvoice
		json.Unmarshal([]byte(message), &invoiceInfo)
		data, _ = json.Marshal(invoiceInfo.invoicePartiallySaleOrder())
	case "GET_SALES_ORDER_RELATIONS":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesOrderRelations(int32(id)))
	case "GET_SALES_INVOICE_RELATIONS":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getSalesInvoiceRelations(int32(id)))
	case "TOGGLE_MANUFACTURING_ORDER":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(toggleManufactuedManufacturingOrder(int64(id)))
	case "MANUFACTURING_ORDER_ALL_SALE_ORDER":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(manufacturingOrderAllSaleOrder(int32(id)))
	case "MANUFACTURING_ORDER_PARTIAL_SALE_ORDER":
		var orderInfo SalesOrderDetailManufacturingOrder
		json.Unmarshal([]byte(message), &orderInfo)
		data, _ = json.Marshal(orderInfo.manufacturingOrderPartiallySaleOrder())
	case "DELETE_SALES_ORDER_DETAIL_PACKAGED":
		var salesOrderDetailPackaged SalesOrderDetailPackaged
		json.Unmarshal([]byte(message), &salesOrderDetailPackaged)
		data, _ = json.Marshal(salesOrderDetailPackaged.deleteSalesOrderDetailPackaged(true))
	case "DELIVERY_NOTE_ALL_SALE_ORDER":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllSaleOrder(int32(id))
		data, _ = json.Marshal(ok)
	case "DELIVERY_NOTE_PARTIALLY_SALE_ORDER":
		var noteInfo SalesOrderDetailDeliveryNote
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallySaleOrder())
	case "SHIPPING_SALE_ORDER":
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
	case "GET_SALES_DELIVERY_NOTE_RELATIONS":
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
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		ok, _ := deliveryNoteAllPurchaseOrder(int32(id))
		data, _ = json.Marshal(ok)
	case "GET_PURCHASE_ORDER_RELATIONS":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseOrderRelations(int32(id)))
	case "GET_INVOICE_ORDER_RELATIONS":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseInvoiceRelations(int32(id)))
	case "GET_PURCHASE_DELIVERY_NOTE_RELATIONS":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(getPurchaseDeliveryNoteRelations(int32(id)))
	case "DELIVERY_NOTE_PARTIALLY_PURCHASE_ORDER":
		var noteInfo PurchaseOrderDetailDeliveryNote
		json.Unmarshal([]byte(message), &noteInfo)
		data, _ = json.Marshal(noteInfo.deliveryNotePartiallyPurchaseOrder())
	case "INVOICE_ALL_PURCHASE_ORDER":
		id, err := strconv.Atoi(message)
		if err != nil {
			return
		}
		data, _ = json.Marshal(invoiceAllPurchaseOrder(int32(id)))
	case "INVOICE_PARTIAL_PURCHASE_ORDER":
		var invoiceInfo PurchaseOrderDetailInvoice
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
	}
	ws.WriteMessage(mt, data)
}

func instructionSearch(command string, message string, mt int, ws *websocket.Conn) {
	var data []byte
	switch command {
	case "CUSTOMER":
		data, _ = json.Marshal(searchCustomers(message))
	case "SUPPLER":
		data, _ = json.Marshal(searchSuppliers(message))
	case "PRODUCT":
		data, _ = json.Marshal(searchProduct(message))
	case "SHIPPING":
		data, _ = json.Marshal(searchShippings(message))
	case "SALES_ORDER":
		var salesOrderSearch SalesOrderSearch
		json.Unmarshal([]byte(message), &salesOrderSearch)
		data, _ = json.Marshal(salesOrderSearch.searchSalesOrder())
	case "SALES_INVOICE":
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchSalesInvoices())
	case "SALES_DELIVERY_NOTE":
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchSalesDelvieryNotes())
	case "PURCHASE_ORDER":
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseOrder())
	case "PURCHASE_INVOICE":
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseInvoice())
	case "PURCHASE_DELIVERY_NOTE":
		var orderSearch OrderSearch
		json.Unmarshal([]byte(message), &orderSearch)
		data, _ = json.Marshal(orderSearch.searchPurchaseDeliveryNote())
	}
	ws.WriteMessage(mt, data)
}
