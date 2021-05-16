package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "marketnet"
	password = "-.qu@.5vaqBE6GH"
	dbname   = "marketnet"
)

var upgrader = websocket.Upgrader{}

var db *sql.DB

func main() {
	fmt.Println("Hola Mundo! :D")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", reverse)
	go http.ListenAndServe(":12279", nil)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ = sql.Open("postgres", psqlInfo)
	db.Ping()

	for {

	}
}

func reverse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HERE!")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

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
		//fmt.Println(command[0:commandSeparatorIndex] + " " + command[commandSeparatorIndex+1:] + " " + string(message[separatorIndex+1:]))
	}
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
	}
}

func instructionGet(command string, message string, mt int, ws *websocket.Conn) {
	var found bool = true
	var data []byte
	switch command {
	case "SALES_ORDER":
		data, _ = json.Marshal(getSalesOrder())
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
	default:
		found = false
	}

	if found {
		ws.WriteMessage(mt, data)
		return
	}

	// NUMERIC
	id, err := strconv.Atoi(message)
	if err != nil || id <= 0 {
		return
	}
	switch command {
	case "SALES_ORDER_DETAIL":
		data, _ = json.Marshal(getSalesOrderDetail(int32(id)))
	case "STOCK":
		data, _ = json.Marshal(getStock(int32(id)))
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
	}
	ws.WriteMessage(mt, data)
}

func instructionLocate(command string, message string, mt int, ws *websocket.Conn) {
	var data []byte

	id, err := strconv.Atoi(message)
	if err != nil {
		return
	}
	switch command {
	case "ADDRESS":
		data, _ = json.Marshal(locateAddress(int32(id)))
	}
	ws.WriteMessage(mt, data)
}
