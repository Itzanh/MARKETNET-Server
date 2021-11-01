package main

// THIS FILE CONTAILS ALL THE FUNCIONALITY FOR THE ERP'S REST API, WHICH CAN BE USED TO INTEGRATE MARKETNET INTO OTHER SOFTWARE WITHOUT FORKING MARKETNET.

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Requests made from an enterprise in the last minute
// Key: enterprise ID
// Value: requests made
//
// Reset every 60 seconds by a cron
var requestsPerMinuteEnterprise map[int32]int32 = make(map[int32]int32)

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
	// accounting
	http.HandleFunc("/api/journal", apiJournal)
	http.HandleFunc("/api/account", apiAccount)
	http.HandleFunc("/api/accounting_movement", apiAccountingMovement)
	http.HandleFunc("/api/accounting_movement_detail", apiAccountingMovementDetail)
	http.HandleFunc("/api/collection_operation", apiCollectionOperation)
	http.HandleFunc("/api/charges", apiCharges)
	http.HandleFunc("/api/payment_transaction", apiPaymentTransaction)
	http.HandleFunc("/api/payment", apiPayments)
	http.HandleFunc("/api/post_sale_invoice", apiPostSaleInvoices)
	http.HandleFunc("/api/post_purchase_invoice", apiPostPurchaseInvoices)
}

func apiSaleOrders(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(paginationQuery.getSalesOrder(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var saleOrder SaleOrder
			json.Unmarshal(body, &saleOrder)
			saleOrder.enterprise = enterpriseId
			ok, _ = saleOrder.insertSalesOrder()
		} else if string(body[0]) == "[" {
			var saleOrders []SaleOrder
			json.Unmarshal(body, &saleOrders)
			for i := 0; i < len(saleOrders); i++ {
				saleOrders[i].enterprise = enterpriseId
				ok, _ = saleOrders[i].insertSalesOrder()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var saleOrder SaleOrder
		json.Unmarshal(body, &saleOrder)
		saleOrder.enterprise = enterpriseId
		ok = saleOrder.updateSalesOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int64(id)
		saleOrder.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getSalesOrderDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var saleOrderDetail SalesOrderDetail
			json.Unmarshal(body, &saleOrderDetail)
			saleOrderDetail.enterprise = enterpriseId
			ok = saleOrderDetail.insertSalesOrderDetail()
		} else if string(body[0]) == "[" {
			var saleOrderDetails []SalesOrderDetail
			json.Unmarshal(body, &saleOrderDetails)
			for i := 0; i < len(saleOrderDetails); i++ {
				saleOrderDetails[i].enterprise = enterpriseId
				ok = saleOrderDetails[i].insertSalesOrderDetail()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(body, &salesOrderDetail)
		salesOrderDetail.enterprise = enterpriseId
		ok = salesOrderDetail.updateSalesOrderDetail()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int64(id)
		saleOrderDetail.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		paginationQuery.Enterprise = enterpriseId
		data, _ := json.Marshal(paginationQuery.getSalesInvoices())
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var saleInvoice SalesInvoice
			json.Unmarshal(body, &saleInvoice)
			saleInvoice.enterprise = enterpriseId
			ok, _ = saleInvoice.insertSalesInvoice()
		} else if string(body[0]) == "[" {
			var saleInvoices []SalesInvoice
			json.Unmarshal(body, &saleInvoices)
			for i := 0; i < len(saleInvoices); i++ {
				saleInvoices[i].enterprise = enterpriseId
				ok, _ = saleInvoices[i].insertSalesInvoice()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int64(id)
		salesInvoice.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getSalesInvoiceDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var salesInvoiceDetail SalesInvoiceDetail
			json.Unmarshal(body, &salesInvoiceDetail)
			salesInvoiceDetail.enterprise = enterpriseId
			ok = salesInvoiceDetail.insertSalesInvoiceDetail(true)
		} else if string(body[0]) == "[" {
			var salesInvoiceDetails []SalesInvoiceDetail
			json.Unmarshal(body, &salesInvoiceDetails)
			for i := 0; i < len(salesInvoiceDetails); i++ {
				salesInvoiceDetails[i].enterprise = enterpriseId
				ok = salesInvoiceDetails[i].insertSalesInvoiceDetail(true)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int64(id)
		salesInvoiceDetail.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		paginationQuery.Enterprise = enterpriseId
		data, _ := json.Marshal(paginationQuery.getSalesDeliveryNotes())
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var salesDeliveryNote SalesDeliveryNote
			json.Unmarshal(body, &salesDeliveryNote)
			salesDeliveryNote.enterprise = enterpriseId
			ok, _ = salesDeliveryNote.insertSalesDeliveryNotes()
		} else if string(body[0]) == "[" {
			var salesDeliveryNotes []SalesDeliveryNote
			json.Unmarshal(body, &salesDeliveryNotes)
			for i := 0; i < len(salesDeliveryNotes); i++ {
				salesDeliveryNotes[i].enterprise = enterpriseId
				ok, _ = salesDeliveryNotes[i].insertSalesDeliveryNotes()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int64(id)
		salesDeliveryNote.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPurchaseOrder(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var purchaseOrder PurchaseOrder
			json.Unmarshal(body, &purchaseOrder)
			purchaseOrder.enterprise = enterpriseId
			ok, _ = purchaseOrder.insertPurchaseOrder()
		} else if string(body[0]) == "[" {
			var purchaseOrders []PurchaseOrder
			json.Unmarshal(body, &purchaseOrders)
			for i := 0; i < len(purchaseOrders); i++ {
				purchaseOrders[i].enterprise = enterpriseId
				ok, _ = purchaseOrders[i].insertPurchaseOrder()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var purchaseOrder PurchaseOrder
		json.Unmarshal(body, &purchaseOrder)
		purchaseOrder.enterprise = enterpriseId
		ok = purchaseOrder.updatePurchaseOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int64(id)
		purchaseOrder.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPurchaseOrderDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var purchaseOrderDetail PurchaseOrderDetail
			json.Unmarshal(body, &purchaseOrderDetail)
			purchaseOrderDetail.enterprise = enterpriseId
			ok, _ = purchaseOrderDetail.insertPurchaseOrderDetail(true)
		} else if string(body[0]) == "[" {
			var purchaseOrderDetails []PurchaseOrderDetail
			json.Unmarshal(body, &purchaseOrderDetails)
			for i := 0; i < len(purchaseOrderDetails); i++ {
				purchaseOrderDetails[i].enterprise = enterpriseId
				ok, _ = purchaseOrderDetails[i].insertPurchaseOrderDetail(true)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int64(id)
		purchaseOrderDetail.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPurchaseInvoices(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var purchaseInvoice PurchaseInvoice
			json.Unmarshal(body, &purchaseInvoice)
			purchaseInvoice.enterprise = enterpriseId
			ok, _ = purchaseInvoice.insertPurchaseInvoice()
		} else if string(body[0]) == "[" {
			var purchaseInvoices []PurchaseInvoice
			json.Unmarshal(body, &purchaseInvoices)
			for i := 0; i < len(purchaseInvoices); i++ {
				purchaseInvoices[i].enterprise = enterpriseId
				ok, _ = purchaseInvoices[i].insertPurchaseInvoice()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int64(id)
		purchaseInvoice.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPurchaseInvoiceDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var purchaseInvoiceDetail PurchaseInvoiceDetail
			json.Unmarshal(body, &purchaseInvoiceDetail)
			purchaseInvoiceDetail.enterprise = enterpriseId
			ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(true)
		} else if string(body[0]) == "[" {
			var purchaseInvoiceDetails []PurchaseInvoiceDetail
			json.Unmarshal(body, &purchaseInvoiceDetails)
			for i := 0; i < len(purchaseInvoiceDetails); i++ {
				purchaseInvoiceDetails[i].enterprise = enterpriseId
				ok = purchaseInvoiceDetails[i].insertPurchaseInvoiceDetail(true)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int64(id)
		purchaseInvoiceDetail.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPurchaseDeliveryNotes(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var purchaseDeliveryNote PurchaseDeliveryNote
			json.Unmarshal(body, &purchaseDeliveryNote)
			purchaseDeliveryNote.enterprise = enterpriseId
			ok, _ = purchaseDeliveryNote.insertPurchaseDeliveryNotes()
		} else if string(body[0]) == "[" {
			var purchaseDeliveryNotes []PurchaseDeliveryNote
			json.Unmarshal(body, &purchaseDeliveryNotes)
			for i := 0; i < len(purchaseDeliveryNotes); i++ {
				purchaseDeliveryNotes[i].enterprise = enterpriseId
				ok, _ = purchaseDeliveryNotes[i].insertPurchaseDeliveryNotes()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int64(id)
		purchaseDeliveryNote.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		paginationQuery.Enterprise = enterpriseId
		data, _ := json.Marshal(paginationQuery.getCustomers())
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var customer Customer
			json.Unmarshal(body, &customer)
			customer.enterprise = enterpriseId
			ok = customer.insertCustomer().Id > 0
		} else if string(body[0]) == "[" {
			var customers []Customer
			json.Unmarshal(body, &customers)
			for i := 0; i < len(customers); i++ {
				customers[i].enterprise = enterpriseId
				ok = customers[i].insertCustomer().Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var customer Customer
		json.Unmarshal(body, &customer)
		customer.enterprise = enterpriseId
		ok = customer.updateCustomer()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customer Customer
		customer.Id = int32(id)
		customer.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getSuppliers(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var supplier Supplier
			json.Unmarshal(body, &supplier)
			supplier.enterprise = enterpriseId
			ok = supplier.insertSupplier().Id > 0
		} else if string(body[0]) == "[" {
			var suppliers []Supplier
			json.Unmarshal(body, &suppliers)
			for i := 0; i < len(suppliers); i++ {
				suppliers[i].enterprise = enterpriseId
				ok = suppliers[i].insertSupplier().Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var supplier Supplier
		json.Unmarshal(body, &supplier)
		supplier.enterprise = enterpriseId
		ok = supplier.updateSupplier()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var supplier Supplier
		supplier.Id = int32(id)
		supplier.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getProduct(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var product Product
			json.Unmarshal(body, &product)
			product.enterprise = enterpriseId
			ok = product.insertProduct()
		} else if string(body[0]) == "[" {
			var products []Product
			json.Unmarshal(body, &products)
			for i := 0; i < len(products); i++ {
				products[i].enterprise = enterpriseId
				ok = products[i].insertProduct()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var product Product
		json.Unmarshal(body, &product)
		product.enterprise = enterpriseId
		ok = product.updateProduct()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var product Product
		product.Id = int32(id)
		product.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getCountries(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var country Country
		json.Unmarshal(body, &country)
		country.enterprise = enterpriseId
		ok = country.insertCountry()
	case "PUT":
		var country Country
		json.Unmarshal(body, &country)
		country.enterprise = enterpriseId
		ok = country.updateCountry()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var country Country
		country.Id = int32(id)
		country.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getStates(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var state State
		json.Unmarshal(body, &state)
		state.enterprise = enterpriseId
		ok = state.insertState()
	case "PUT":
		var state State
		json.Unmarshal(body, &state)
		state.enterprise = enterpriseId
		ok = state.updateState()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var state State
		state.Id = int32(id)
		state.enterprise = enterpriseId
		ok = state.deleteState()
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getColor(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var color Color
		json.Unmarshal(body, &color)
		color.enterprise = enterpriseId
		ok = color.insertColor()
	case "PUT":
		var color Color
		json.Unmarshal(body, &color)
		color.enterprise = enterpriseId
		ok = color.updateColor()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var color Color
		color.Id = int32(id)
		color.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getProductFamilies(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		productFamily.enterprise = enterpriseId
		ok = productFamily.insertProductFamily()
	case "PUT":
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		productFamily.enterprise = enterpriseId
		ok = productFamily.updateProductFamily()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var productFamily ProductFamily
		productFamily.Id = int32(id)
		productFamily.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		paginationQuery.Enterprise = enterpriseId
		data, _ := json.Marshal(paginationQuery.getAddresses())
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var address Address
			json.Unmarshal(body, &address)
			address.enterprise = enterpriseId
			ok = address.insertAddress().Id > 0
		} else if string(body[0]) == "[" {
			var addresses []Address
			json.Unmarshal(body, &addresses)
			for i := 0; i < len(addresses); i++ {
				addresses[i].enterprise = enterpriseId
				ok = addresses[i].insertAddress().Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		var address Address
		json.Unmarshal(body, &address)
		address.enterprise = enterpriseId
		ok = address.updateAddress()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var address Address
		address.Id = int32(id)
		address.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getCariers(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var carrier Carrier
		json.Unmarshal(body, &carrier)
		carrier.enterprise = enterpriseId
		ok = carrier.insertCarrier()
	case "PUT":
		var carrier Carrier
		json.Unmarshal(body, &carrier)
		carrier.enterprise = enterpriseId
		ok = carrier.updateCarrier()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var carrier Carrier
		carrier.Id = int32(id)
		carrier.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getBillingSeries(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		serie.enterprise = enterpriseId
		ok = serie.insertBillingSerie()
	case "PUT":
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		serie.enterprise = enterpriseId
		ok = serie.updateBillingSerie()
	case "DELETE":
		var serie BillingSerie
		serie.Id = string(body)
		serie.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getCurrencies(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var currency Currency
		json.Unmarshal(body, &currency)
		currency.enterprise = enterpriseId
		ok = currency.insertCurrency()
	case "PUT":
		var currency Currency
		json.Unmarshal(body, &currency)
		currency.enterprise = enterpriseId
		ok = currency.updateCurrency()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var currency Currency
		currency.Id = int32(id)
		currency.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPaymentMethods(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		paymentMethod.enterprise = enterpriseId
		ok = paymentMethod.insertPaymentMethod()
	case "PUT":
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		paymentMethod.enterprise = enterpriseId
		ok = paymentMethod.updatePaymentMethod()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var paymentMethod PaymentMethod
		paymentMethod.Id = int32(id)
		paymentMethod.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getLanguages(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var language Language
		json.Unmarshal(body, &language)
		language.enterprise = enterpriseId
		ok = language.insertLanguage()
	case "PUT":
		var language Language
		json.Unmarshal(body, &language)
		language.enterprise = enterpriseId
		ok = language.updateLanguage()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var language Language
		language.Id = int32(id)
		language.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPackages(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var packages Packages
		json.Unmarshal(body, &packages)
		packages.enterprise = enterpriseId
		ok = packages.insertPackage()
	case "PUT":
		var packages Packages
		json.Unmarshal(body, &packages)
		packages.enterprise = enterpriseId
		ok = packages.updatePackage()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var packages Packages
		packages.Id = int32(id)
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getIncoterm(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		incoterm.enterprise = enterpriseId
		ok = incoterm.insertIncoterm()
	case "PUT":
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		incoterm.enterprise = enterpriseId
		ok = incoterm.updateIncoterm()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var incoterm Incoterm
		incoterm.Id = int32(id)
		incoterm.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getWarehouses(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		warehouse.enterprise = enterpriseId
		ok = warehouse.insertWarehouse()
	case "PUT":
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		warehouse.enterprise = enterpriseId
		ok = warehouse.updateWarehouse()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var warehouse Warehouse
		warehouse.Id = string(body)
		warehouse.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		paginationQuery.Enterprise = enterpriseId
		data, _ := json.Marshal(paginationQuery.getWarehouseMovement())
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var warehouseMovement WarehouseMovement
			json.Unmarshal(body, &warehouseMovement)
			warehouseMovement.enterprise = enterpriseId
			ok = warehouseMovement.insertWarehouseMovement()
		} else if string(body[0]) == "[" {
			var warehouseMovements []WarehouseMovement
			json.Unmarshal(body, &warehouseMovements)
			for i := 0; i < len(warehouseMovements); i++ {
				warehouseMovements[i].enterprise = enterpriseId
				ok = warehouseMovements[i].insertWarehouseMovement()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		warehouseMovement.enterprise = enterpriseId
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
	// auth
	ok, userId, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(manufacturingPaginationQuery.getAllManufacturingOrders(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var manufacturingOrder ManufacturingOrder
		json.Unmarshal(body, &manufacturingOrder)
		manufacturingOrder.UserCreated = userId
		manufacturingOrder.enterprise = enterpriseId
		ok = manufacturingOrder.insertManufacturingOrder()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		manufacturingOrder.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getManufacturingOrderType(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "PUT":
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		manufacturingOrderType.enterprise = enterpriseId
		ok = manufacturingOrderType.updateManufacturingOrderType()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int32(id)
		manufacturingOrderType.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getShippings(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		shipping.enterprise = enterpriseId
		ok, _ = shipping.insertShipping()
	case "PUT":
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		shipping.enterprise = enterpriseId
		ok = shipping.updateShipping()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var shipping Shipping
		shipping.Id = int64(id)
		shipping.enterprise = enterpriseId
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
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getStock(int32(id), enterpriseId))
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusNotAcceptable)
}

func apiJournal(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getJournals(enterpriseId))
		w.Write(data)
		return
	case "POST":
		var journal Journal
		json.Unmarshal(body, &journal)
		journal.enterprise = enterpriseId
		ok = journal.insertJournal()
	case "PUT":
		var journal Journal
		json.Unmarshal(body, &journal)
		journal.enterprise = enterpriseId
		ok = journal.updateJournal()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var journal Journal
		journal.Id = int32(id)
		journal.enterprise = enterpriseId
		ok = journal.deleteJournal()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiAccount(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getAccounts(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var account Account
			json.Unmarshal(body, &account)
			account.enterprise = enterpriseId
			ok = account.insertAccount()
		} else if string(body[0]) == "[" {
			var account []Account
			json.Unmarshal(body, &account)
			for i := 0; i < len(account); i++ {
				account[i].enterprise = enterpriseId
				ok = account[i].insertAccount()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		var account Account
		json.Unmarshal(body, &account)
		account.enterprise = enterpriseId
		ok = account.updateAccount()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var account Account
		account.Id = int32(id)
		account.enterprise = enterpriseId
		ok = account.deleteAccount()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiAccountingMovement(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getAccountingMovement(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var accountingMovement AccountingMovement
			json.Unmarshal(body, &accountingMovement)
			accountingMovement.enterprise = enterpriseId
			ok = accountingMovement.insertAccountingMovement()
		} else if string(body[0]) == "[" {
			var accountingMovement []AccountingMovement
			json.Unmarshal(body, &accountingMovement)
			for i := 0; i < len(accountingMovement); i++ {
				accountingMovement[i].enterprise = enterpriseId
				ok = accountingMovement[i].insertAccountingMovement()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var accountingMovement AccountingMovement
		accountingMovement.Id = int64(id)
		accountingMovement.enterprise = enterpriseId
		ok = accountingMovement.deleteAccountingMovement()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiAccountingMovementDetail(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getAccountingMovementDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var accountingMovementDetail AccountingMovementDetail
			json.Unmarshal(body, &accountingMovementDetail)
			accountingMovementDetail.enterprise = enterpriseId
			ok = accountingMovementDetail.insertAccountingMovementDetail()
		} else if string(body[0]) == "[" {
			var accountingMovementDetail []AccountingMovementDetail
			json.Unmarshal(body, &accountingMovementDetail)
			for i := 0; i < len(accountingMovementDetail); i++ {
				accountingMovementDetail[i].enterprise = enterpriseId
				ok = accountingMovementDetail[i].insertAccountingMovementDetail()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		accountingMovementDetail.Id = int64(id)
		accountingMovementDetail.enterprise = enterpriseId
		ok = accountingMovementDetail.deleteAccountingMovementDetail()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCollectionOperation(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getColletionOperations(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var collectionOperation CollectionOperation
			json.Unmarshal(body, &collectionOperation)
			collectionOperation.enterprise = enterpriseId
			ok = collectionOperation.insertCollectionOperation()
		} else if string(body[0]) == "[" {
			var collectionOperation []CollectionOperation
			json.Unmarshal(body, &collectionOperation)
			for i := 0; i < len(collectionOperation); i++ {
				collectionOperation[i].enterprise = enterpriseId
				ok = collectionOperation[i].insertCollectionOperation()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var collectionOperation CollectionOperation
		collectionOperation.Id = int32(id)
		collectionOperation.enterprise = enterpriseId
		ok = collectionOperation.deleteCollectionOperation()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCharges(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getCharges(int32(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var charges Charges
			json.Unmarshal(body, &charges)
			charges.enterprise = enterpriseId
			ok = charges.insertCharges()
		} else if string(body[0]) == "[" {
			var charges []Charges
			json.Unmarshal(body, &charges)
			for i := 0; i < len(charges); i++ {
				charges[i].enterprise = enterpriseId
				ok = charges[i].insertCharges()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var charges Charges
		charges.Id = int32(id)
		charges.enterprise = enterpriseId
		ok = charges.deleteCharges()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPaymentTransaction(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPaymentTransactions(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var paymentTransaction PaymentTransaction
			json.Unmarshal(body, &paymentTransaction)
			paymentTransaction.enterprise = enterpriseId
			ok = paymentTransaction.insertPaymentTransaction()
		} else if string(body[0]) == "[" {
			var paymentTransaction []PaymentTransaction
			json.Unmarshal(body, &paymentTransaction)
			for i := 0; i < len(paymentTransaction); i++ {
				paymentTransaction[i].enterprise = enterpriseId
				ok = paymentTransaction[i].insertPaymentTransaction()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var paymentTransaction PaymentTransaction
		paymentTransaction.Id = int32(id)
		paymentTransaction.enterprise = enterpriseId
		ok = paymentTransaction.deletePaymentTransaction()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPayments(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		data, _ := json.Marshal(getPayments(int32(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if string(body[0]) == "{" {
			var Paymenp Payment
			json.Unmarshal(body, &Paymenp)
			Paymenp.enterprise = enterpriseId
			ok = Paymenp.insertPayment()
		} else if string(body[0]) == "[" {
			var payment []Payment
			json.Unmarshal(body, &payment)
			for i := 0; i < len(payment); i++ {
				payment[i].enterprise = enterpriseId
				ok = payment[i].insertPayment()
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var payment Payment
		payment.Id = int32(id)
		payment.enterprise = enterpriseId
		ok = payment.deletePayment()
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiPostSaleInvoices(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		var orderSearch OrderSearch
		json.Unmarshal(body, &orderSearch)
		orderSearch.Enterprise = enterpriseId
		orderSearch.NotPosted = true
		data, _ := json.Marshal(orderSearch.searchSalesInvoices())
		w.Write(data)
		return
	case "POST":
		var invoiceIds []int64
		json.Unmarshal(body, &invoiceIds)
		result := salesPostInvoices(invoiceIds, enterpriseId)
		resp, _ := json.Marshal(result)
		w.Write(resp)
		return
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusNotAcceptable)

}

func apiPostPurchaseInvoices(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId := checkApiKey(r)
	if !ok {
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
		var orderSearch OrderSearch
		json.Unmarshal(body, &orderSearch)
		orderSearch.Enterprise = enterpriseId
		orderSearch.NotPosted = true
		data, _ := json.Marshal(orderSearch.searchPurchaseInvoice())
		w.Write(data)
		return
	case "POST":
		var invoiceIds []int64
		json.Unmarshal(body, &invoiceIds)
		result := purchasePostInvoices(invoiceIds, enterpriseId)
		resp, _ := json.Marshal(result)
		w.Write(resp)
		return
	case "PUT":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case "DELETE":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusNotAcceptable)

}
