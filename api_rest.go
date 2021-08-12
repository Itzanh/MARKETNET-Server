package main

// THIS FILE CONTAILS ALL THE FUNCIONALITY FOR THE ERP'S REST API, WHICH CAN BE USED TO INTEGRATE MARKETNET INTO OTHER SOFTWARE WITHOUT FORKING MARKETNET.

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// If the REST API on the ERP is enabled, this funcion is called and adds the events listeners for the ERP's REST API to function.
func addHttpHandlerFuncions() {
	// sales
	http.HandleFunc("/api/sale_orders", apiSaleOrders)
	http.HandleFunc("/api/sale_order_details", apiSaleOrderDetails)
	http.HandleFunc("/api/sale_invoices", apiSaleInvoices)
	http.HandleFunc("/api/sale_invoice_details", apiSaleInvoiceDetals)
	http.HandleFunc("/api/sale_delivery_notes", apiSaleDeliveryNotes)
	// purchases
	http.HandleFunc("/api/purchase_orders", apiPurchaseOrders)
	http.HandleFunc("/api/purchase_order_details", apiPurchaseOrderDetails)
	http.HandleFunc("/api/purchase_invoices", apiPurchaseInvoices)
	http.HandleFunc("/api/purchase_invoice_details", apiPurchaseInvoiceDetails)
	http.HandleFunc("/api/purchase_delivery_notes", apiPurchaseDeliveryNotes)
	// masters
	http.HandleFunc("/api/customers", apiCustomers)
	http.HandleFunc("/api/suppliers", apiSuppliers)
	http.HandleFunc("/api/products", apiProducts)
	http.HandleFunc("/api/countries", apiCountries)
	http.HandleFunc("/api/states", apiStates)
	http.HandleFunc("/api/colors", apiColors)
	http.HandleFunc("/api/product_families", apiProductFamilies)
	http.HandleFunc("/api/addresses", apiAddresses)
	http.HandleFunc("/api/carriers", apiCarriers)
	http.HandleFunc("/api/billing_series", apiBillingSeries)
	http.HandleFunc("/api/currencies", apiCurrencies)
	http.HandleFunc("/api/payment_methods", apiPaymentMethods)
	http.HandleFunc("/api/languages", apiLanguages)
	http.HandleFunc("/api/packages", apiPackages)
	http.HandleFunc("/api/incoterms", apiIncoterms)
	// warehouse
	http.HandleFunc("/api/warehouses", apiWarehouses)
	http.HandleFunc("/api/warehouse_movements", apiWarehouseMovements)
	// manufacturing
	http.HandleFunc("/api/manufacturing_orders", apiManufacturingOrders)
	http.HandleFunc("/api/manufacturing_order_types", apiManufacturingOrderTypes)
	// preparation
	http.HandleFunc("/api/shippings", apiShipping)
	// stock
	http.HandleFunc("/api/stock", apiStock)
}

func apiSaleOrders(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getSalesOrder())
		w.Write(data)
		return
	case "POST":
		var saleOrder SaleOrder
		json.Unmarshal(body, &saleOrder)
		ok, _ = saleOrder.insertSalesOrder()
	case "PUT":
		var saleOrder SaleOrder
		json.Unmarshal(body, &saleOrder)
		ok = saleOrder.updateSalesOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int32(id)
		ok = saleOrder.deleteSalesOrder()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSaleOrderDetails(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getSalesOrderDetail(int32(id)))
		w.Write(data)
		return
	case "POST":
		var saleOrderDetail SalesOrderDetail
		json.Unmarshal(body, &saleOrderDetail)
		ok = saleOrderDetail.insertSalesOrderDetail()
	case "PUT":
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(body, &salesOrderDetail)
		ok = salesOrderDetail.updateSalesOrderDetail()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int32(id)
		ok = saleOrderDetail.deleteSalesOrderDetail()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSaleInvoices(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getSalesInvoices())
		w.Write(data)
		return
	case "POST":
		var saleInvoice SalesInvoice
		json.Unmarshal(body, &saleInvoice)
		ok, _ = saleInvoice.insertSalesInvoice()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int32(id)
		ok = salesInvoice.deleteSalesInvoice()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSaleInvoiceDetals(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getSalesInvoiceDetail(int32(id)))
		w.Write(data)
		return
	case "POST":
		var salesInvoiceDetail SalesInvoiceDetail
		json.Unmarshal(body, &salesInvoiceDetail)
		ok = salesInvoiceDetail.insertSalesInvoiceDetail(true)
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int32(id)
		ok = salesInvoiceDetail.deleteSalesInvoiceDetail()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSaleDeliveryNotes(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getSalesDeliveryNotes())
		w.Write(data)
		return
	case "POST":
		var salesDeliveryNote SalesDeliveryNote
		json.Unmarshal(body, &salesDeliveryNote)
		ok, _ = salesDeliveryNote.insertSalesDeliveryNotes()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int32(id)
		ok = salesDeliveryNote.deleteSalesDeliveryNotes()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getPurchaseOrder())
		w.Write(data)
		return
	case "POST":
		var purchaseOrder PurchaseOrder
		json.Unmarshal(body, &purchaseOrder)
		ok, _ = purchaseOrder.insertPurchaseOrder()
	case "PUT":
		var PurchaseOrdep PurchaseOrder
		json.Unmarshal(body, &PurchaseOrdep)
		ok = PurchaseOrdep.updatePurchaseOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int32(id)
		ok = purchaseOrder.deletePurchaseOrder()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPurchaseOrderDetails(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPurchaseOrderDetail(int32(id)))
		w.Write(data)
		return
	case "POST":
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(body, &purchaseOrderDetail)
		ok, _ = purchaseOrderDetail.insertPurchaseOrderDetail(true)
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int32(id)
		ok = purchaseOrderDetail.deletePurchaseOrderDetail()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPurchaseInvoices(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getPurchaseInvoices())
		w.Write(data)
		return
	case "POST":
		var purchaseInvoice PurchaseInvoice
		json.Unmarshal(body, &purchaseInvoice)
		ok, _ = purchaseInvoice.insertPurchaseInvoice()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int32(id)
		ok = purchaseInvoice.deletePurchaseInvoice()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPurchaseInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPurchaseInvoiceDetail(int32(id)))
		w.Write(data)
		return
	case "POST":
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		json.Unmarshal(body, &purchaseInvoiceDetail)
		ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(true)
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int32(id)
		ok = purchaseInvoiceDetail.deletePurchaseInvoiceDetail()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPurchaseDeliveryNotes(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getPurchaseDeliveryNotes())
		w.Write(data)
		return
	case "POST":
		var purchaseDeliveryNote PurchaseDeliveryNote
		json.Unmarshal(body, &purchaseDeliveryNote)
		ok, _ = purchaseDeliveryNote.insertPurchaseDeliveryNotes()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int32(id)
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCustomers(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getCustomers())
		w.Write(data)
		return
	case "POST":
		var customer Customer
		json.Unmarshal(body, &customer)
		ok = customer.insertCustomer()
	case "PUT":
		var customer Customer
		json.Unmarshal(body, &customer)
		ok = customer.updateCustomer()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customer Customer
		customer.Id = int32(id)
		ok = customer.deleteCustomer()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSuppliers(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getSuppliers())
		w.Write(data)
		return
	case "POST":
		var supplier Supplier
		json.Unmarshal(body, &supplier)
		ok = supplier.insertSupplier()
	case "PUT":
		var supplier Supplier
		json.Unmarshal(body, &supplier)
		ok = supplier.updateSupplier()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var supplier Supplier
		supplier.Id = int32(id)
		ok = supplier.deleteSupplier()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiProducts(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getProduct())
		w.Write(data)
		return
	case "POST":
		var product Product
		json.Unmarshal(body, &product)
		ok = product.insertProduct()
	case "PUT":
		var product Product
		json.Unmarshal(body, &product)
		ok = product.updateProduct()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var product Product
		product.Id = int32(id)
		ok = product.deleteProduct()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCountries(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getCountries())
		w.Write(data)
		return
	case "POST":
		var country Country
		json.Unmarshal(body, &country)
		ok = country.insertCountry()
	case "PUT":
		var country Country
		json.Unmarshal(body, &country)
		ok = country.updateCountry()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var country Country
		country.Id = int16(id)
		ok = country.deleteCountry()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiStates(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getStates())
		w.Write(data)
		return
	case "POST":
		var state State
		json.Unmarshal(body, &state)
		ok = state.insertState()
	case "PUT":
		var city State
		json.Unmarshal(body, &city)
		ok = city.updateState()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var city State
		city.Id = int32(id)
		ok = city.deleteState()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiColors(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getColor())
		w.Write(data)
		return
	case "POST":
		var color Color
		json.Unmarshal(body, &color)
		ok = color.insertColor()
	case "PUT":
		var color Color
		json.Unmarshal(body, &color)
		ok = color.updateColor()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var color Color
		color.Id = int16(id)
		ok = color.deleteColor()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiProductFamilies(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getProductFamilies())
		w.Write(data)
		return
	case "POST":
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		ok = productFamily.insertProductFamily()
	case "PUT":
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		ok = productFamily.updateProductFamily()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var productFamily ProductFamily
		productFamily.Id = int16(id)
		ok = productFamily.deleteProductFamily()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiAddresses(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getAddresses())
		w.Write(data)
		return
	case "POST":
		var address Address
		json.Unmarshal(body, &address)
		ok = address.insertAddress()
	case "PUT":
		var address Address
		json.Unmarshal(body, &address)
		ok = address.updateAddress()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var address Address
		address.Id = int32(id)
		ok = address.deleteAddress()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCarriers(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getCariers())
		w.Write(data)
		return
	case "POST":
		var carrier Carrier
		json.Unmarshal(body, &carrier)
		ok = carrier.insertCarrier()
	case "PUT":
		var incoterm Carrier
		json.Unmarshal(body, &incoterm)
		ok = incoterm.updateCarrier()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var carrier Carrier
		carrier.Id = int16(id)
		ok = carrier.deleteCarrier()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiBillingSeries(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getBillingSeries())
		w.Write(data)
		return
	case "POST":
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		ok = serie.insertBillingSerie()
	case "PUT":
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		ok = serie.updateBillingSerie()
	case "DELETE":
		var serie BillingSerie
		serie.Id = string(body)
		ok = serie.deleteBillingSerie()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCurrencies(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getCurrencies())
		w.Write(data)
		return
	case "POST":
		var currency Currency
		json.Unmarshal(body, &currency)
		ok = currency.insertCurrency()
	case "PUT":
		var currency Currency
		json.Unmarshal(body, &currency)
		ok = currency.updateCurrency()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var currency Currency
		currency.Id = int16(id)
		ok = currency.deleteCurrency()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPaymentMethods(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getPaymentMethods())
		w.Write(data)
		return
	case "POST":
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		ok = paymentMethod.insertPaymentMethod()
	case "PUT":
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		ok = paymentMethod.updatePaymentMethod()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var paymentMethod PaymentMethod
		paymentMethod.Id = int16(id)
		ok = paymentMethod.deletePaymentMethod()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiLanguages(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getLanguages())
		w.Write(data)
		return
	case "POST":
		var language Language
		json.Unmarshal(body, &language)
		ok = language.insertLanguage()
	case "PUT":
		var language Language
		json.Unmarshal(body, &language)
		ok = language.updateLanguage()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var language Language
		language.Id = int16(id)
		ok = language.deleteLanguage()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPackages(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getPackages())
		w.Write(data)
		return
	case "POST":
		var packages Packages
		json.Unmarshal(body, &packages)
		ok = packages.insertPackage()
	case "PUT":
		var packages Packages
		json.Unmarshal(body, &packages)
		ok = packages.updatePackage()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var packages Packages
		packages.Id = int16(id)
		ok = packages.deletePackage()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiIncoterms(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getIncoterm())
		w.Write(data)
		return
	case "POST":
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		ok = incoterm.insertIncoterm()
	case "PUT":
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		ok = incoterm.updateIncoterm()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var incoterm Incoterm
		incoterm.Id = int16(id)
		ok = incoterm.deleteIncoterm()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiWarehouses(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getWarehouses())
		w.Write(data)
		return
	case "POST":
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		ok = warehouse.insertWarehouse()
	case "PUT":
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		ok = warehouse.updateWarehouse()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var warehouse Warehouse
		warehouse.Id = string(body)
		ok = warehouse.deleteWarehouse()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiWarehouseMovements(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		data, _ := json.Marshal(paginationQuery.getWarehouseMovement())
		w.Write(data)
		return
	case "POST":
		var warehouseMovement WarehouseMovement
		json.Unmarshal(body, &warehouseMovement)
		ok = warehouseMovement.insertWarehouseMovement()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		ok = warehouseMovement.deleteWarehouseMovement()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiManufacturingOrders(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		var manufacturingPaginationQuery ManufacturingPaginationQuery
		json.Unmarshal(body, &manufacturingPaginationQuery)
		data, _ := json.Marshal(manufacturingPaginationQuery.getAllManufacturingOrders())
		w.Write(data)
		return
	case "POST":
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(body, &manufacturingOrder)
		manufacturingOrder.UserCreated = userId
		ok = manufacturingOrder.insertManufacturingOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		ok = manufacturingOrder.deleteManufacturingOrder()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiManufacturingOrderTypes(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getManufacturingOrderType())
		w.Write(data)
		return
	case "POST":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "PUT":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		ok = manufacturingOrderType.updateManufacturingOrderType()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int16(id)
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiShipping(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(getShippings())
		w.Write(data)
		return
	case "POST":
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		ok, _ = shipping.insertShipping()
	case "PUT":
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		ok = shipping.updateShipping()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var shipping Shipping
		shipping.Id = int32(id)
		ok = shipping.deleteShipping()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiStock(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// token
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok, userId := checkApiKey(token[0])
	if !ok || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getStock(int32(id)))
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusNotAcceptable)
}
