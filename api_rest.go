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
	http.HandleFunc("/api/sale_order_details_digital_product_data", apiSaleOrderDetailsDigitalProductData)
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
	http.HandleFunc("/api/custom_fields", apiCustomFields)
	// warehouse
	http.HandleFunc("/api/warehouses", apiWarehouses)
	http.HandleFunc("/api/warehouse_movements", apiWarehouseMovements)
	// manufacturing
	http.HandleFunc("/api/manufacturing_orders", apiManufacturingOrders)
	http.HandleFunc("/api/manufacturing_order_types", apiManufacturingOrderTypes)
	http.HandleFunc("/api/complex_manufacturing_orders", apiComplexManufacturingOrders)
	http.HandleFunc("/api/complex_manufacturing_orders_components", apiComplexManufacturingOrdersComponents)
	http.HandleFunc("/api/manufacturing_order_type_components", apiManufacturingOrderTypesComponents)
	// preparation
	http.HandleFunc("/api/shippings", apiShipping)
	http.HandleFunc("/api/shipping_status_history", apiShippingStatusHistory)
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleOrders.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getSalesOrder(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.SaleOrders.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var saleOrder SaleOrder
			json.Unmarshal(body, &saleOrder)
			saleOrder.EnterpriseId = enterpriseId
			ok, _ = saleOrder.insertSalesOrder(userId)
		} else if string(body[0]) == "[" {
			var saleOrders []SaleOrder
			json.Unmarshal(body, &saleOrders)
			for i := 0; i < len(saleOrders); i++ {
				saleOrders[i].EnterpriseId = enterpriseId
				ok, _ = saleOrders[i].insertSalesOrder(userId)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.SaleOrders.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var saleOrder SaleOrder
		json.Unmarshal(body, &saleOrder)
		saleOrder.EnterpriseId = enterpriseId
		ok = saleOrder.updateSalesOrder(userId)
	case "DELETE":
		if !permission.SaleOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrder SaleOrder
		saleOrder.Id = int64(id)
		saleOrder.EnterpriseId = enterpriseId
		ok = saleOrder.deleteSalesOrder(userId).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleOrderDetails.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getSalesOrderDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.SaleOrderDetails.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var saleOrderDetail SalesOrderDetail
			json.Unmarshal(body, &saleOrderDetail)
			saleOrderDetail.EnterpriseId = enterpriseId
			ok = saleOrderDetail.insertSalesOrderDetail(userId).Ok
		} else if string(body[0]) == "[" {
			var saleOrderDetails []SalesOrderDetail
			json.Unmarshal(body, &saleOrderDetails)
			for i := 0; i < len(saleOrderDetails); i++ {
				saleOrderDetails[i].EnterpriseId = enterpriseId
				ok = saleOrderDetails[i].insertSalesOrderDetail(userId).Ok
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.SaleOrderDetails.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var salesOrderDetail SalesOrderDetail
		json.Unmarshal(body, &salesOrderDetail)
		salesOrderDetail.EnterpriseId = enterpriseId
		ok = salesOrderDetail.updateSalesOrderDetail(userId).Ok
	case "DELETE":
		if !permission.SaleOrderDetails.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var saleOrderDetail SalesOrderDetail
		saleOrderDetail.Id = int64(id)
		saleOrderDetail.EnterpriseId = enterpriseId
		ok = saleOrderDetail.deleteSalesOrderDetail(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiSaleOrderDetailsDigitalProductData(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleOrderDetailsDigitalProductData.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getSalesOrderDetailDigitalProductData(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.SaleOrderDetailsDigitalProductData.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var d SalesOrderDetailDigitalProductData
			json.Unmarshal(body, &d)
			d.EnterpriseId = enterpriseId
			ok = d.insertSalesOrderDetailDigitalProductData()
		} else if string(body[0]) == "[" {
			var d []SalesOrderDetailDigitalProductData
			json.Unmarshal(body, &d)
			for i := 0; i < len(d); i++ {
				d[i].EnterpriseId = enterpriseId
				ok = d[i].insertSalesOrderDetailDigitalProductData()
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.SaleOrderDetailsDigitalProductData.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var d SalesOrderDetailDigitalProductData
		json.Unmarshal(body, &d)
		d.EnterpriseId = enterpriseId
		ok = d.updateSalesOrderDetailDigitalProductData()
	case "DELETE":
		if !permission.SaleOrderDetailsDigitalProductData.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var d SalesOrderDetailDigitalProductData
		d.Id = int32(id)
		d.EnterpriseId = enterpriseId
		ok = d.deleteSalesOrderDetailDigitalProductData()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleInvoices.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		paginationQuery.enterprise = enterpriseId
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getSalesInvoices())
		w.Write(data)
		return
	case "POST":
		if !permission.SaleInvoices.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var saleInvoice SalesInvoice
			json.Unmarshal(body, &saleInvoice)
			saleInvoice.EnterpriseId = enterpriseId
			ok, _ = saleInvoice.insertSalesInvoice(userId, nil)
		} else if string(body[0]) == "[" {
			var saleInvoices []SalesInvoice
			json.Unmarshal(body, &saleInvoices)
			for i := 0; i < len(saleInvoices); i++ {
				saleInvoices[i].EnterpriseId = enterpriseId
				ok, _ = saleInvoices[i].insertSalesInvoice(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.SaleInvoices.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoice SalesInvoice
		salesInvoice.Id = int64(id)
		salesInvoice.EnterpriseId = enterpriseId
		ok = salesInvoice.deleteSalesInvoice(userId).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleInvoiceDetails.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getSalesInvoiceDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.SaleInvoiceDetails.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var salesInvoiceDetail SalesInvoiceDetail
			json.Unmarshal(body, &salesInvoiceDetail)
			salesInvoiceDetail.EnterpriseId = enterpriseId
			ok = salesInvoiceDetail.insertSalesInvoiceDetail(nil, userId).Ok
		} else if string(body[0]) == "[" {
			var salesInvoiceDetails []SalesInvoiceDetail
			json.Unmarshal(body, &salesInvoiceDetails)
			for i := 0; i < len(salesInvoiceDetails); i++ {
				salesInvoiceDetails[i].EnterpriseId = enterpriseId
				ok = salesInvoiceDetails[i].insertSalesInvoiceDetail(nil, userId).Ok
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.SaleInvoiceDetails.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesInvoiceDetail SalesInvoiceDetail
		salesInvoiceDetail.Id = int64(id)
		salesInvoiceDetail.EnterpriseId = enterpriseId
		ok = salesInvoiceDetail.deleteSalesInvoiceDetail(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.SaleDeliveryNotes.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		paginationQuery.enterprise = enterpriseId
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getSalesDeliveryNotes())
		w.Write(data)
		return
	case "POST":
		if !permission.SaleDeliveryNotes.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var salesDeliveryNote SalesDeliveryNote
			json.Unmarshal(body, &salesDeliveryNote)
			salesDeliveryNote.EnterpriseId = enterpriseId
			ok, _ = salesDeliveryNote.insertSalesDeliveryNotes(userId, nil)
		} else if string(body[0]) == "[" {
			var salesDeliveryNotes []SalesDeliveryNote
			json.Unmarshal(body, &salesDeliveryNotes)
			for i := 0; i < len(salesDeliveryNotes); i++ {
				salesDeliveryNotes[i].EnterpriseId = enterpriseId
				ok, _ = salesDeliveryNotes[i].insertSalesDeliveryNotes(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.SaleDeliveryNotes.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var salesDeliveryNote SalesDeliveryNote
		salesDeliveryNote.Id = int64(id)
		salesDeliveryNote.EnterpriseId = enterpriseId
		ok = salesDeliveryNote.deleteSalesDeliveryNotes(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PurchaseOrders.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getPurchaseOrder(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PurchaseOrders.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var purchaseOrder PurchaseOrder
			json.Unmarshal(body, &purchaseOrder)
			purchaseOrder.EnterpriseId = enterpriseId
			ok, _ = purchaseOrder.insertPurchaseOrder(userId, nil)
		} else if string(body[0]) == "[" {
			var purchaseOrders []PurchaseOrder
			json.Unmarshal(body, &purchaseOrders)
			for i := 0; i < len(purchaseOrders); i++ {
				purchaseOrders[i].EnterpriseId = enterpriseId
				ok, _ = purchaseOrders[i].insertPurchaseOrder(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.PurchaseOrders.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var purchaseOrder PurchaseOrder
		json.Unmarshal(body, &purchaseOrder)
		purchaseOrder.EnterpriseId = enterpriseId
		ok = purchaseOrder.updatePurchaseOrder(userId)
	case "DELETE":
		if !permission.PurchaseOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrder PurchaseOrder
		purchaseOrder.Id = int64(id)
		purchaseOrder.EnterpriseId = enterpriseId
		ok = purchaseOrder.deletePurchaseOrder(userId).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PurchaseOrderDetails.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPurchaseOrderDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PurchaseOrderDetails.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var purchaseOrderDetail PurchaseOrderDetail
			json.Unmarshal(body, &purchaseOrderDetail)
			purchaseOrderDetail.EnterpriseId = enterpriseId
			okAndErr, _ := purchaseOrderDetail.insertPurchaseOrderDetail(userId, nil)
			ok = okAndErr.Ok
		} else if string(body[0]) == "[" {
			var purchaseOrderDetails []PurchaseOrderDetail
			json.Unmarshal(body, &purchaseOrderDetails)
			for i := 0; i < len(purchaseOrderDetails); i++ {
				purchaseOrderDetails[i].EnterpriseId = enterpriseId
				ok, _ := purchaseOrderDetails[i].insertPurchaseOrderDetail(userId, nil)
				if !ok.Ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.PurchaseOrders.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		json.Unmarshal(body, &purchaseOrderDetail)
		purchaseOrderDetail.EnterpriseId = enterpriseId
		ok = purchaseOrderDetail.updatePurchaseOrderDetail(userId).Ok
	case "DELETE":
		if !permission.PurchaseOrderDetails.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseOrderDetail PurchaseOrderDetail
		purchaseOrderDetail.Id = int64(id)
		purchaseOrderDetail.EnterpriseId = enterpriseId
		ok = purchaseOrderDetail.deletePurchaseOrderDetail(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PurchaseInvoices.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getPurchaseInvoices(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PurchaseInvoices.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var purchaseInvoice PurchaseInvoice
			json.Unmarshal(body, &purchaseInvoice)
			purchaseInvoice.EnterpriseId = enterpriseId
			ok, _ = purchaseInvoice.insertPurchaseInvoice(userId, nil)
		} else if string(body[0]) == "[" {
			var purchaseInvoices []PurchaseInvoice
			json.Unmarshal(body, &purchaseInvoices)
			for i := 0; i < len(purchaseInvoices); i++ {
				purchaseInvoices[i].EnterpriseId = enterpriseId
				ok, _ = purchaseInvoices[i].insertPurchaseInvoice(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.PurchaseInvoices.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoice PurchaseInvoice
		purchaseInvoice.Id = int64(id)
		purchaseInvoice.EnterpriseId = enterpriseId
		ok = purchaseInvoice.deletePurchaseInvoice(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PurchaseInvoiceDetails.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPurchaseInvoiceDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PurchaseInvoiceDetails.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var purchaseInvoiceDetail PurchaseInvoiceDetail
			json.Unmarshal(body, &purchaseInvoiceDetail)
			purchaseInvoiceDetail.EnterpriseId = enterpriseId
			ok = purchaseInvoiceDetail.insertPurchaseInvoiceDetail(userId, nil).Ok
		} else if string(body[0]) == "[" {
			var purchaseInvoiceDetails []PurchaseInvoiceDetail
			json.Unmarshal(body, &purchaseInvoiceDetails)
			for i := 0; i < len(purchaseInvoiceDetails); i++ {
				purchaseInvoiceDetails[i].EnterpriseId = enterpriseId
				ok = purchaseInvoiceDetails[i].insertPurchaseInvoiceDetail(userId, nil).Ok
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.PurchaseInvoiceDetails.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseInvoiceDetail PurchaseInvoiceDetail
		purchaseInvoiceDetail.Id = int64(id)
		purchaseInvoiceDetail.EnterpriseId = enterpriseId
		ok = purchaseInvoiceDetail.deletePurchaseInvoiceDetail(userId, nil).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PurchaseDeliveryNotes.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getPurchaseDeliveryNotes(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PurchaseDeliveryNotes.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var purchaseDeliveryNote PurchaseDeliveryNote
			json.Unmarshal(body, &purchaseDeliveryNote)
			purchaseDeliveryNote.EnterpriseId = enterpriseId
			ok, _ = purchaseDeliveryNote.insertPurchaseDeliveryNotes(userId, nil)
		} else if string(body[0]) == "[" {
			var purchaseDeliveryNotes []PurchaseDeliveryNote
			json.Unmarshal(body, &purchaseDeliveryNotes)
			for i := 0; i < len(purchaseDeliveryNotes); i++ {
				purchaseDeliveryNotes[i].EnterpriseId = enterpriseId
				ok, _ = purchaseDeliveryNotes[i].insertPurchaseDeliveryNotes(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.PurchaseDeliveryNotes.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var purchaseDeliveryNote PurchaseDeliveryNote
		purchaseDeliveryNote.Id = int64(id)
		purchaseDeliveryNote.EnterpriseId = enterpriseId
		ok = purchaseDeliveryNote.deletePurchaseDeliveryNotes(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Customers.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		paginationQuery.enterprise = enterpriseId
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getCustomers())
		w.Write(data)
		return
	case "POST":
		if !permission.Customers.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var customer Customer
			json.Unmarshal(body, &customer)
			customer.EnterpriseId = enterpriseId
			ok = customer.insertCustomer(userId).Id > 0
		} else if string(body[0]) == "[" {
			var customers []Customer
			json.Unmarshal(body, &customers)
			for i := 0; i < len(customers); i++ {
				customers[i].EnterpriseId = enterpriseId
				ok = customers[i].insertCustomer(userId).Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.Customers.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var customer Customer
		json.Unmarshal(body, &customer)
		customer.EnterpriseId = enterpriseId
		ok = customer.updateCustomer(userId)
	case "DELETE":
		if !permission.Customers.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customer Customer
		customer.Id = int32(id)
		customer.EnterpriseId = enterpriseId
		ok = customer.deleteCustomer(userId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Suppliers.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getSuppliers(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Suppliers.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var supplier Supplier
			json.Unmarshal(body, &supplier)
			supplier.EnterpriseId = enterpriseId
			ok = supplier.insertSupplier(userId).Id > 0
		} else if string(body[0]) == "[" {
			var suppliers []Supplier
			json.Unmarshal(body, &suppliers)
			for i := 0; i < len(suppliers); i++ {
				suppliers[i].EnterpriseId = enterpriseId
				ok = suppliers[i].insertSupplier(userId).Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.Suppliers.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var supplier Supplier
		json.Unmarshal(body, &supplier)
		supplier.EnterpriseId = enterpriseId
		ok = supplier.updateSupplier(userId)
	case "DELETE":
		if !permission.Suppliers.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var supplier Supplier
		supplier.Id = int32(id)
		supplier.EnterpriseId = enterpriseId
		ok = supplier.deleteSupplier(userId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Products.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getProduct(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Suppliers.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var product Product
			json.Unmarshal(body, &product)
			product.EnterpriseId = enterpriseId
			ok = product.insertProduct(userId).Ok
		} else if string(body[0]) == "[" {
			var products []Product
			json.Unmarshal(body, &products)
			for i := 0; i < len(products); i++ {
				products[i].EnterpriseId = enterpriseId
				ok = products[i].insertProduct(userId).Ok
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.Suppliers.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var product Product
		json.Unmarshal(body, &product)
		product.EnterpriseId = enterpriseId
		ok = product.updateProduct(userId).Ok
	case "DELETE":
		if !permission.Suppliers.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var product Product
		product.Id = int32(id)
		product.EnterpriseId = enterpriseId
		ok = product.deleteProduct(userId).Ok
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Countries.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getCountries(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Countries.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var country Country
		json.Unmarshal(body, &country)
		country.EnterpriseId = enterpriseId
		ok = country.insertCountry()
	case "PUT":
		if !permission.Countries.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var country Country
		json.Unmarshal(body, &country)
		country.EnterpriseId = enterpriseId
		ok = country.updateCountry()
	case "DELETE":
		if !permission.Countries.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var country Country
		country.Id = int32(id)
		country.EnterpriseId = enterpriseId
		ok = country.deleteCountry()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.States.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getStates(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.States.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var state State
		json.Unmarshal(body, &state)
		state.EnterpriseId = enterpriseId
		ok = state.insertState()
	case "PUT":
		if !permission.States.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var state State
		json.Unmarshal(body, &state)
		state.EnterpriseId = enterpriseId
		ok = state.updateState()
	case "DELETE":
		if !permission.States.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var state State
		state.Id = int32(id)
		state.EnterpriseId = enterpriseId
		ok = state.deleteState()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Colors.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getColor(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Colors.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var color Color
		json.Unmarshal(body, &color)
		color.EnterpriseId = enterpriseId
		ok = color.insertColor()
	case "PUT":
		if !permission.Colors.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var color Color
		json.Unmarshal(body, &color)
		color.EnterpriseId = enterpriseId
		ok = color.updateColor()
	case "DELETE":
		if !permission.Colors.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var color Color
		color.Id = int32(id)
		color.EnterpriseId = enterpriseId
		ok = color.deleteColor()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.ProductFamilies.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getProductFamilies(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.ProductFamilies.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		productFamily.EnterpriseId = enterpriseId
		ok = productFamily.insertProductFamily()
	case "PUT":
		if !permission.ProductFamilies.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var productFamily ProductFamily
		json.Unmarshal(body, &productFamily)
		productFamily.EnterpriseId = enterpriseId
		ok = productFamily.updateProductFamily()
	case "DELETE":
		if !permission.ProductFamilies.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var productFamily ProductFamily
		productFamily.Id = int32(id)
		productFamily.EnterpriseId = enterpriseId
		ok = productFamily.deleteProductFamily()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Addresses.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		paginationQuery.enterprise = enterpriseId
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getAddresses())
		w.Write(data)
		return
	case "POST":
		if !permission.Addresses.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var address Address
			json.Unmarshal(body, &address)
			address.EnterpriseId = enterpriseId
			ok = address.insertAddress(userId).Id > 0
		} else if string(body[0]) == "[" {
			var addresses []Address
			json.Unmarshal(body, &addresses)
			for i := 0; i < len(addresses); i++ {
				addresses[i].EnterpriseId = enterpriseId
				ok = addresses[i].insertAddress(userId).Id > 0
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.Addresses.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var address Address
		json.Unmarshal(body, &address)
		address.EnterpriseId = enterpriseId
		ok = address.updateAddress()
	case "DELETE":
		if !permission.Addresses.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var address Address
		address.Id = int32(id)
		address.EnterpriseId = enterpriseId
		ok = address.deleteAddress()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Carriers.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getCariers(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Carriers.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var carrier Carrier
		json.Unmarshal(body, &carrier)
		carrier.EnterpriseId = enterpriseId
		ok = carrier.insertCarrier()
	case "PUT":
		if !permission.Carriers.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var carrier Carrier
		json.Unmarshal(body, &carrier)
		carrier.EnterpriseId = enterpriseId
		ok = carrier.updateCarrier()
	case "DELETE":
		if !permission.Carriers.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var carrier Carrier
		carrier.Id = int32(id)
		carrier.EnterpriseId = enterpriseId
		ok = carrier.deleteCarrier()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.BillingSeries.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getBillingSeries(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.BillingSeries.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		serie.EnterpriseId = enterpriseId
		ok = serie.insertBillingSerie()
	case "PUT":
		if !permission.BillingSeries.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var serie BillingSerie
		json.Unmarshal(body, &serie)
		serie.EnterpriseId = enterpriseId
		ok = serie.updateBillingSerie()
	case "DELETE":
		if !permission.BillingSeries.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var serie BillingSerie
		serie.Id = string(body)
		serie.EnterpriseId = enterpriseId
		ok = serie.deleteBillingSerie()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Currencies.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getCurrencies(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Currencies.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var currency Currency
		json.Unmarshal(body, &currency)
		currency.EnterpriseId = enterpriseId
		ok = currency.insertCurrency()
	case "PUT":
		if !permission.Currencies.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var currency Currency
		json.Unmarshal(body, &currency)
		currency.EnterpriseId = enterpriseId
		ok = currency.updateCurrency()
	case "DELETE":
		if !permission.Currencies.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var currency Currency
		currency.Id = int32(id)
		currency.EnterpriseId = enterpriseId
		ok = currency.deleteCurrency()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.PaymentMethods.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getPaymentMethods(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PaymentMethods.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		paymentMethod.EnterpriseId = enterpriseId
		ok = paymentMethod.insertPaymentMethod()
	case "PUT":
		if !permission.PaymentMethods.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paymentMethod PaymentMethod
		json.Unmarshal(body, &paymentMethod)
		paymentMethod.EnterpriseId = enterpriseId
		ok = paymentMethod.updatePaymentMethod()
	case "DELETE":
		if !permission.PaymentMethods.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var paymentMethod PaymentMethod
		paymentMethod.Id = int32(id)
		paymentMethod.EnterpriseId = enterpriseId
		ok = paymentMethod.deletePaymentMethod()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Languages.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getLanguages(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Languages.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var language Language
		json.Unmarshal(body, &language)
		language.EnterpriseId = enterpriseId
		ok = language.insertLanguage()
	case "PUT":
		if !permission.Languages.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var language Language
		json.Unmarshal(body, &language)
		language.EnterpriseId = enterpriseId
		ok = language.updateLanguage()
	case "DELETE":
		if !permission.Languages.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var language Language
		language.Id = int32(id)
		language.EnterpriseId = enterpriseId
		ok = language.deleteLanguage()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Packages.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getPackages(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Packages.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var packages Packages
		json.Unmarshal(body, &packages)
		packages.EnterpriseId = enterpriseId
		ok = packages.insertPackage()
	case "PUT":
		if !permission.Packages.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var packages Packages
		json.Unmarshal(body, &packages)
		packages.EnterpriseId = enterpriseId
		ok = packages.updatePackage()
	case "DELETE":
		if !permission.Packages.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var packages Packages
		packages.Id = int32(id)
		ok = packages.deletePackage()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Incoterms.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getIncoterm(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Incoterms.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		incoterm.EnterpriseId = enterpriseId
		ok = incoterm.insertIncoterm()
	case "PUT":
		if !permission.Incoterms.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var incoterm Incoterm
		json.Unmarshal(body, &incoterm)
		incoterm.EnterpriseId = enterpriseId
		ok = incoterm.updateIncoterm()
	case "DELETE":
		if !permission.Incoterms.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var incoterm Incoterm
		incoterm.Id = int32(id)
		incoterm.EnterpriseId = enterpriseId
		ok = incoterm.deleteIncoterm()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiCustomFields(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, _ := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		var query CustomFields
		json.Unmarshal(body, &query)
		query.EnterpriseId = enterpriseId
		data, _ := json.Marshal(query.getCustomFields())
		w.Write(data)
		return
	case "POST":
		var customField CustomFields
		json.Unmarshal(body, &customField)
		customField.EnterpriseId = enterpriseId
		ok = customField.insertCustomFields()
	case "PUT":
		var customField CustomFields
		json.Unmarshal(body, &customField)
		customField.EnterpriseId = enterpriseId
		ok = customField.updateCustomFields()
	case "DELETE":
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customField CustomFields
		customField.Id = int64(id)
		customField.EnterpriseId = enterpriseId
		ok = customField.deleteCustomFields()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Warehouses.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getWarehouses(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Warehouses.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		warehouse.EnterpriseId = enterpriseId
		ok = warehouse.insertWarehouse()
	case "PUT":
		if !permission.Warehouses.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var warehouse Warehouse
		json.Unmarshal(body, &warehouse)
		warehouse.EnterpriseId = enterpriseId
		ok = warehouse.updateWarehouse()
	case "DELETE":
		if !permission.Warehouses.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var warehouse Warehouse
		warehouse.Id = string(body)
		warehouse.EnterpriseId = enterpriseId
		ok = warehouse.deleteWarehouse()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.WarehouseMovements.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var paginationQuery PaginationQuery
		json.Unmarshal(body, &paginationQuery)
		paginationQuery.enterprise = enterpriseId
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && paginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(paginationQuery.getWarehouseMovement())
		w.Write(data)
		return
	case "POST":
		if !permission.WarehouseMovements.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var warehouseMovement WarehouseMovement
			json.Unmarshal(body, &warehouseMovement)
			warehouseMovement.EnterpriseId = enterpriseId
			ok = warehouseMovement.insertWarehouseMovement(userId, nil)
		} else if string(body[0]) == "[" {
			var warehouseMovements []WarehouseMovement
			json.Unmarshal(body, &warehouseMovements)
			for i := 0; i < len(warehouseMovements); i++ {
				warehouseMovements[i].EnterpriseId = enterpriseId
				ok = warehouseMovements[i].insertWarehouseMovement(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "DELETE":
		if !permission.WarehouseMovements.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var warehouseMovement WarehouseMovement
		warehouseMovement.Id = int64(id)
		warehouseMovement.EnterpriseId = enterpriseId
		ok = warehouseMovement.deleteWarehouseMovement(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.ManufacturingOrders.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingPaginationQuery ManufacturingPaginationQuery
		json.Unmarshal(body, &manufacturingPaginationQuery)
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && manufacturingPaginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(manufacturingPaginationQuery.getManufacturingOrder(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.ManufacturingOrders.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var manufacturingOrder ManufacturingOrder
			json.Unmarshal(body, &manufacturingOrder)
			manufacturingOrder.UserCreatedId = userId
			manufacturingOrder.EnterpriseId = enterpriseId
			ok = manufacturingOrder.insertManufacturingOrder(userId, nil).Ok
		} else if string(body[0]) == "[" {
			var manufacturingOrders []ManufacturingOrder
			json.Unmarshal(body, &manufacturingOrders)
			for i := 0; i < len(manufacturingOrders); i++ {
				manufacturingOrders[i].UserCreatedId = userId
				manufacturingOrders[i].EnterpriseId = enterpriseId
				ok = manufacturingOrders[i].insertManufacturingOrder(userId, nil).Ok
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.ManufacturingOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ok = toggleManufactuedManufacturingOrder(int64(id), userId, enterpriseId)
	case "DELETE":
		if !permission.ManufacturingOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrder ManufacturingOrder
		manufacturingOrder.Id = int64(id)
		manufacturingOrder.EnterpriseId = enterpriseId
		ok = manufacturingOrder.deleteManufacturingOrder(enterpriseId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.ManufacturingOrderTypes.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getManufacturingOrderType(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.ManufacturingOrderTypes.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		manufacturingOrderType.EnterpriseId = enterpriseId
		ok = manufacturingOrderType.insertManufacturingOrderType()
	case "PUT":
		if !permission.ManufacturingOrderTypes.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		json.Unmarshal(body, &manufacturingOrderType)
		manufacturingOrderType.EnterpriseId = enterpriseId
		ok = manufacturingOrderType.updateManufacturingOrderType()
	case "DELETE":
		if !permission.ManufacturingOrderTypes.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrderType ManufacturingOrderType
		manufacturingOrderType.Id = int32(id)
		manufacturingOrderType.EnterpriseId = enterpriseId
		ok = manufacturingOrderType.deleteManufacturingOrderType()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiComplexManufacturingOrders(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.ComplexManufacturingOrders.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingPaginationQuery ManufacturingPaginationQuery
		json.Unmarshal(body, &manufacturingPaginationQuery)
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && manufacturingPaginationQuery.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(manufacturingPaginationQuery.getAllComplexManufacturingOrders(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.ComplexManufacturingOrders.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var complexManufacturingOrder ComplexManufacturingOrder
			json.Unmarshal(body, &complexManufacturingOrder)
			complexManufacturingOrder.UserCreatedId = userId
			complexManufacturingOrder.EnterpriseId = enterpriseId
			ok, _ = complexManufacturingOrder.insertComplexManufacturingOrder(userId, nil)
		} else if string(body[0]) == "[" {
			var complexManufacturingOrders []ComplexManufacturingOrder
			json.Unmarshal(body, &complexManufacturingOrders)
			for i := 0; i < len(complexManufacturingOrders); i++ {
				complexManufacturingOrders[i].UserCreatedId = userId
				complexManufacturingOrders[i].EnterpriseId = enterpriseId
				ok, _ = complexManufacturingOrders[i].insertComplexManufacturingOrder(userId, nil)
				if !ok {
					break
				}
			}
		} else {
			ok = false
		}
	case "PUT":
		if !permission.ComplexManufacturingOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		toggleManufactuedComplexManufacturingOrder(int64(id), userId, enterpriseId)
	case "DELETE":
		if !permission.ComplexManufacturingOrders.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var complexManufacturingOrder ComplexManufacturingOrder
		complexManufacturingOrder.Id = int64(id)
		complexManufacturingOrder.EnterpriseId = enterpriseId
		ok = complexManufacturingOrder.deleteComplexManufacturingOrder(enterpriseId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiComplexManufacturingOrdersComponents(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.ComplexManufacturingOrdersComponents.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getComplexManufacturingOrderManufacturingOrder(int64(id), enterpriseId))
		w.Write(data)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func apiManufacturingOrderTypesComponents(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.ManufacturingOrderTypeComponents.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getManufacturingOrderTypeComponents(int32(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.ManufacturingOrderTypeComponents.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingOrderTypeComponent ManufacturingOrderTypeComponents
		json.Unmarshal(body, &manufacturingOrderTypeComponent)
		manufacturingOrderTypeComponent.EnterpriseId = enterpriseId
		ok, _ = manufacturingOrderTypeComponent.insertManufacturingOrderTypeComponents()
	case "PUT":
		if !permission.ManufacturingOrderTypeComponents.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var manufacturingOrderTypeComponent ManufacturingOrderTypeComponents
		json.Unmarshal(body, &manufacturingOrderTypeComponent)
		manufacturingOrderTypeComponent.EnterpriseId = enterpriseId
		ok, _ = manufacturingOrderTypeComponent.updateManufacturingOrderTypeComponents()
	case "DELETE":
		if !permission.ManufacturingOrderTypeComponents.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var manufacturingOrderTypeComponent ManufacturingOrderTypeComponents
		manufacturingOrderTypeComponent.Id = int32(id)
		manufacturingOrderTypeComponent.EnterpriseId = enterpriseId
		ok = manufacturingOrderTypeComponent.deleteManufacturingOrderTypeComponents()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	ok = false
	switch r.Method {
	case "GET":
		if !permission.Shippings.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getShippings(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Shippings.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		shipping.EnterpriseId = enterpriseId
		ok, _ = shipping.insertShipping(userId, nil)
	case "PUT":
		if !permission.Shippings.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var shipping Shipping
		json.Unmarshal(body, &shipping)
		shipping.EnterpriseId = enterpriseId
		ok = shipping.updateShipping(userId)
	case "DELETE":
		if !permission.Shippings.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var shipping Shipping
		shipping.Id = int64(id)
		shipping.EnterpriseId = enterpriseId
		ok = shipping.deleteShipping(userId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp, _ := json.Marshal(ok)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	w.Write(resp)
}

func apiShippingStatusHistory(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.ShippingStatusHistory.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getShippingStatusHistory(enterpriseId, int64(id)))
		w.Write(data)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func apiStock(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.Stock.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getStock(int32(id), enterpriseId))
		w.Write(data)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func apiJournal(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.Journal.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getJournals(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Journal.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var journal Journal
		json.Unmarshal(body, &journal)
		journal.EnterpriseId = enterpriseId
		ok = journal.insertJournal()
	case "PUT":
		if !permission.Journal.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var journal Journal
		json.Unmarshal(body, &journal)
		journal.EnterpriseId = enterpriseId
		ok = journal.updateJournal()
	case "DELETE":
		if !permission.Journal.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var journal Journal
		journal.Id = int32(id)
		journal.EnterpriseId = enterpriseId
		ok = journal.deleteJournal()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

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
	ok, _, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.Account.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getAccounts(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Account.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var account Account
			json.Unmarshal(body, &account)
			account.EnterpriseId = enterpriseId
			ok = account.insertAccount()
		} else if string(body[0]) == "[" {
			var account []Account
			json.Unmarshal(body, &account)
			for i := 0; i < len(account); i++ {
				account[i].EnterpriseId = enterpriseId
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
		if !permission.Account.Put {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var account Account
		json.Unmarshal(body, &account)
		account.EnterpriseId = enterpriseId
		ok = account.updateAccount()
	case "DELETE":
		if !permission.Account.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var account Account
		account.Id = int32(id)
		account.EnterpriseId = enterpriseId
		ok = account.deleteAccount()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.AccountingMovement.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, _ := json.Marshal(getAccountingMovement(enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.AccountingMovement.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var accountingMovement AccountingMovement
			json.Unmarshal(body, &accountingMovement)
			accountingMovement.EnterpriseId = enterpriseId
			ok = accountingMovement.insertAccountingMovement(userId, nil)
		} else if string(body[0]) == "[" {
			var accountingMovement []AccountingMovement
			json.Unmarshal(body, &accountingMovement)
			for i := 0; i < len(accountingMovement); i++ {
				accountingMovement[i].EnterpriseId = enterpriseId
				ok = accountingMovement[i].insertAccountingMovement(userId, nil)
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "DELETE":
		if !permission.AccountingMovement.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var accountingMovement AccountingMovement
		accountingMovement.Id = int64(id)
		accountingMovement.EnterpriseId = enterpriseId
		ok = accountingMovement.deleteAccountingMovement(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.AccountingMovementDetail.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getAccountingMovementDetail(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.AccountingMovementDetail.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var accountingMovementDetail AccountingMovementDetail
			json.Unmarshal(body, &accountingMovementDetail)
			accountingMovementDetail.EnterpriseId = enterpriseId
			ok = accountingMovementDetail.insertAccountingMovementDetail(userId, nil)
		} else if string(body[0]) == "[" {
			var accountingMovementDetail []AccountingMovementDetail
			json.Unmarshal(body, &accountingMovementDetail)
			for i := 0; i < len(accountingMovementDetail); i++ {
				accountingMovementDetail[i].EnterpriseId = enterpriseId
				ok = accountingMovementDetail[i].insertAccountingMovementDetail(userId, nil)
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "DELETE":
		if !permission.AccountingMovementDetail.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var accountingMovementDetail AccountingMovementDetail
		accountingMovementDetail.Id = int64(id)
		accountingMovementDetail.EnterpriseId = enterpriseId
		ok = accountingMovementDetail.deleteAccountingMovementDetail(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.CollectionOperation.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getColletionOperations(int64(id), enterpriseId))
		w.Write(data)
		return
	case "DELETE":
		if !permission.CollectionOperation.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var collectionOperation CollectionOperation
		collectionOperation.Id = int32(id)
		collectionOperation.EnterpriseId = enterpriseId
		ok = collectionOperation.deleteCollectionOperation(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.Charges.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getCharges(int32(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Charges.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var charges Charges
			json.Unmarshal(body, &charges)
			charges.EnterpriseId = enterpriseId
			ok = charges.insertCharges(userId)
		} else if string(body[0]) == "[" {
			var charges []Charges
			json.Unmarshal(body, &charges)
			for i := 0; i < len(charges); i++ {
				charges[i].EnterpriseId = enterpriseId
				ok = charges[i].insertCharges(userId)
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "DELETE":
		if !permission.Charges.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var charges Charges
		charges.Id = int32(id)
		charges.EnterpriseId = enterpriseId
		ok = charges.deleteCharges(userId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.PaymentTransaction.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPaymentTransactions(int64(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.PaymentTransaction.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var paymentTransaction PaymentTransaction
			json.Unmarshal(body, &paymentTransaction)
			paymentTransaction.EnterpriseId = enterpriseId
			ok = paymentTransaction.insertPaymentTransaction(userId, nil)
		} else if string(body[0]) == "[" {
			var paymentTransaction []PaymentTransaction
			json.Unmarshal(body, &paymentTransaction)
			for i := 0; i < len(paymentTransaction); i++ {
				paymentTransaction[i].EnterpriseId = enterpriseId
				ok = paymentTransaction[i].insertPaymentTransaction(userId, nil)
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "DELETE":
		if !permission.PaymentTransaction.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var paymentTransaction PaymentTransaction
		paymentTransaction.Id = int32(id)
		paymentTransaction.EnterpriseId = enterpriseId
		ok = paymentTransaction.deletePaymentTransaction(userId, nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.Payment.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(getPayments(int32(id), enterpriseId))
		w.Write(data)
		return
	case "POST":
		if !permission.Payment.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body[0]) == "{" {
			var Paymenp Payment
			json.Unmarshal(body, &Paymenp)
			Paymenp.EnterpriseId = enterpriseId
			ok = Paymenp.insertPayment(userId)
		} else if string(body[0]) == "[" {
			var payment []Payment
			json.Unmarshal(body, &payment)
			for i := 0; i < len(payment); i++ {
				payment[i].EnterpriseId = enterpriseId
				ok = payment[i].insertPayment(userId)
				if !ok {
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	case "DELETE":
		if !permission.Payment.Delete {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var payment Payment
		payment.Id = int32(id)
		payment.EnterpriseId = enterpriseId
		ok = payment.deletePayment(userId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.PostSaleInvoice.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal(body, &orderSearch)
		orderSearch.enterprise = enterpriseId
		orderSearch.NotPosted = true
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && orderSearch.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(orderSearch.searchSalesInvoices())
		w.Write(data)
		return
	case "POST":
		if !permission.PostSaleInvoice.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var invoiceIds []int64
		json.Unmarshal(body, &invoiceIds)
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && len(invoiceIds) > int(settings.Server.WebSecurity.MaxLimitApiQueries) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := salesPostInvoices(invoiceIds, enterpriseId, userId)
		resp, _ := json.Marshal(result)
		w.Write(resp)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func apiPostPurchaseInvoices(w http.ResponseWriter, r *http.Request) {
	// headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	w.Header().Add("Content-type", "application/json")
	// auth
	ok, userId, enterpriseId, permission := checkApiKey(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check body length
	if r.ContentLength > settings.Server.WebSecurity.MaxRequestBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, settings.Server.WebSecurity.MaxRequestBodyLength)
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	// methods
	switch r.Method {
	case "GET":
		if !permission.PostPurchaseInvoice.Get {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var orderSearch OrderSearch
		json.Unmarshal(body, &orderSearch)
		orderSearch.enterprise = enterpriseId
		orderSearch.NotPosted = true
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && orderSearch.Limit > settings.Server.WebSecurity.MaxLimitApiQueries {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(orderSearch.searchPurchaseInvoice())
		w.Write(data)
		return
	case "POST":
		if !permission.PostPurchaseInvoice.Post {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var invoiceIds []int64
		json.Unmarshal(body, &invoiceIds)
		if settings.Server.WebSecurity.MaxLimitApiQueries > 0 && len(invoiceIds) > int(settings.Server.WebSecurity.MaxLimitApiQueries) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := purchasePostInvoices(invoiceIds, enterpriseId, userId)
		resp, _ := json.Marshal(result)
		w.Write(resp)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
