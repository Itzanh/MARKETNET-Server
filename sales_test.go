package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func ConnectTestWithDB(t *testing.T) {
	// read settings
	var ok bool
	settings, ok = getBackendSettings()
	if !ok {
		t.Error("ERROR READING SETTINGS FILE")
		return
	}

	// connect to PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", settings.Db.Host, settings.Db.Port, settings.Db.User, settings.Db.Password, settings.Db.Dbname)
	db, _ = sql.Open("postgres", psqlInfo) // control error
	err := db.Ping()
	if err != nil {
		t.Error(err)
		return
	}
}

/* == SALES == */

/* GET */

// ===== SALE ORDERS

func TestGetSalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32}
	o := q.getSalesOrder(1)

	for i := 0; i < len(o.Orders); i++ {
		if o.Orders[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}
}

func TestSearchSalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := SalesOrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}}}
	o := q.searchSalesOrder()

	for i := 0; i < len(o.Orders); i++ {
		if o.Orders[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

	// search for ID
	q = SalesOrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}, Search: "1"}}
	o = q.searchSalesOrder()

	for i := 0; i < len(o.Orders); i++ {
		if o.Orders[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

	// search for customer name
	q = SalesOrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}, Search: "Itzan"}}
	o = q.searchSalesOrder()

	for i := 0; i < len(o.Orders); i++ {
		if o.Orders[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = SalesOrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}}, DateStart: &start, DateEnd: &end}
	o = q.searchSalesOrder()

	for i := 0; i < len(o.Orders); i++ {
		if o.Orders[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}
}

func TestGetStatusSalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesOrderPreparation(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

	o = getSalesOrderAwaitingShipping(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

	o = getSalesOrderStatus("A", 1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}
}

func TestGetRowStatusSalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesOrderRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, sale order row with ID 0.")
		return
	}

}

func TestLocateSaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := locateSaleOrder(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, sale orders with ID 0.")
			return
		}
	}

}

/* INSERT - UPDATE - DELETE */

func TestIsValidSalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}
	if !o.isValid() {
		t.Error("Incorrect is valid in sale order.")
		return
	}
}

func TestSalesOrderInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	ok, orderId := o.insertSalesOrder()
	if !ok || orderId <= 0 {
		t.Error("Insert error, sale order not inserted.")
		return
	}

	o.Id = orderId
	carrer := int32(1)
	o.Carrier = &carrer

	o.enterprise = 1
	ok = o.updateSalesOrder()
	if !ok || orderId <= 0 {
		t.Error("Update error, sale order not updated.")
		return
	}

	orderInMemory := getSalesOrderRow(orderId)
	if *orderInMemory.Carrier != carrer {
		t.Error("Update not successful, sale order not updated.")
		return
	}

	o.enterprise = 1
	ok = o.deleteSalesOrder()
	if !ok {
		t.Error("Delete error, sale order not deleted.")
		return
	}
}

/* FUNCTIONALITY */

func TestGetSaleOrderDefaults(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	d := getSaleOrderDefaults(1)
	if len(d.Warehouse) == 0 || len(d.WarehouseName) == 0 {
		t.Error("Can't get sales order defaults, empty defaults.")
	}
}

func TestGetSalesOrderRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32}
	o := q.getSalesOrder(1)

	var checkInvoices int8 = 0            // 0 = Not checked, 1 = OK, 2 = Error
	var checkManufacturingOrders int8 = 0 // 0 = Not checked, 1 = OK, 2 = Error
	var checkDeliveryNotes int8 = 0       // 0 = Not checked, 1 = OK, 2 = Error
	var checkShippings int8 = 0           // 0 = Not checked, 1 = OK, 2 = Error

	for i := 0; i < len(o.Orders); i++ {
		r := getSalesOrderRelations(o.Orders[i].Id, 1)

		if checkInvoices == 0 && len(r.Invoices) > 0 {
			if r.Invoices[0].Id <= 0 {
				checkInvoices = 2
			} else {
				checkInvoices = 1
			}
		}

		if checkManufacturingOrders == 0 && len(r.ManufacturingOrders) > 0 {
			if r.ManufacturingOrders[0].Id <= 0 {
				checkManufacturingOrders = 2
			} else {
				checkManufacturingOrders = 1
			}
		}

		if checkDeliveryNotes == 0 && len(r.DeliveryNotes) > 0 {
			if r.DeliveryNotes[0].Id <= 0 {
				checkDeliveryNotes = 2
			} else {
				checkDeliveryNotes = 1
			}
		}

		if checkShippings == 0 && len(r.Shippings) > 0 {
			if r.Shippings[0].Id <= 0 {
				checkShippings = 2
			} else {
				checkShippings = 1
			}
		}

		if checkInvoices != 0 || checkManufacturingOrders != 0 || checkDeliveryNotes != 0 || checkShippings != 0 {
			break
		}
	}

	if checkInvoices == 2 || checkManufacturingOrders == 2 || checkDeliveryNotes == 2 || checkShippings == 2 {
		t.Errorf("Error scanning sale order relations checkInvoices %q checkManufacturingOrders %q checkDeliveryNotes %q checkShippings %q", checkInvoices, checkManufacturingOrders, checkDeliveryNotes, checkShippings)
	}
}

// ===== SALE ORDER DETAILS

/* GET */

func TestGetRowStatusSalesOrderDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	d := getSalesOrderDetail(1, 1)
	for i := 0; i < len(d); i++ {
		if d[i].Id <= 0 {
			t.Error("Scan error, sale order details with ID 0.")
			return
		}
	}

}

func TestGetRowStatusSalesOrderDetailRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesOrderDetailRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, sale order detail row with ID 0.")
		return
	}

}

func TestGetSalesOrderDetailWaitingForPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	for i := 1; i < MAX_INT32; i++ {
		d := getSalesOrderDetailWaitingForPurchaseOrder(int32(i))
		if len(d) > 0 {
			if d[0].Id <= 0 {
				t.Error("Scan error, sale order details with ID 0.")
				return
			} else {
				return
			}
		}
	}

}

func TestGetSalesOrderDetailPurchaseOrderPending(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	for i := 1; i < MAX_INT32; i++ {
		d := getSalesOrderDetailPurchaseOrderPending(int64(i))
		if len(d) > 0 {
			if d[0].Id <= 0 {
				t.Error("Scan error, sale order details with ID 0.")
				return
			} else {
				return
			}
		}
	}
}

/* INSERT - UPDATE - DELETE */

func TestIsValidSaleOrderDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	d := SalesOrderDetail{
		Order:                    55042,
		Product:                  4,
		Price:                    9.99,
		Quantity:                 2,
		VatPercent:               21,
		TotalAmount:              24.1758,
		QuantityInvoiced:         0,
		QuantityDeliveryNote:     0,
		Status:                   "_",
		QuantityPendingPackaging: 2,
		PurchaseOrderDetail:      nil,
		prestaShopId:             0,
		ProductName:              "",
		Cancelled:                false}

	ok := d.isValid()
	if !ok {
		t.Error("Sale order detail not valid")
	}
}

func TestSaleOrderDetailInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	// test insert
	ok := d.insertSalesOrderDetail()
	if !ok {
		t.Error("Insert error, sale order detail not inserted")
		return
	}

	// check the total amount is correct
	details := getSalesOrderDetail(orderId, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't retrieve the sale order details")
		return
	}

	if details[0].TotalAmount != (details[0].Price*float64(details[0].Quantity))*(1+(details[0].VatPercent/100)) {
		t.Error("Incorrect total amount creating sale order detail")
		return
	}

	// check that the sale order has been updated correctly
	inMemoryOrder := getSalesOrderRow(orderId)
	if inMemoryOrder.TotalProducts != float64(details[0].Quantity)*details[0].Price || inMemoryOrder.TotalAmount != details[0].TotalAmount {
		t.Error("The sale order totals has not been updated")
		return
	}

	if inMemoryOrder.LinesNumber != 1 {
		t.Error("The sale order number of lines is not updated upon insert")
		return
	}

	// attemp to update quantity
	details[0].Quantity = 4
	details[0].enterprise = 1
	ok = details[0].updateSalesOrderDetail()
	if !ok {
		t.Error("Update error, sale order detail not updated")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't retrieve the sale order details")
		return
	}

	// check the quantity has been updated
	if details[0].Quantity != 4 {
		t.Error("Update not successful, the sale order has not been updated")
		return
	}

	// check the total amount is correct
	if details[0].TotalAmount != (details[0].Price*float64(details[0].Quantity))*(1+(details[0].VatPercent/100)) {
		t.Error("Incorrect total amount updating sale order detail")
		return
	}

	// check that the sale order has been updated correctly
	inMemoryOrder = getSalesOrderRow(orderId)
	if inMemoryOrder.TotalProducts != float64(details[0].Quantity)*details[0].Price || inMemoryOrder.TotalAmount != details[0].TotalAmount {
		t.Error("The sale order totals has not been updated")
		return
	}

	// cancel the detail, check if cancelled
	ok = cancelSalesOrderDetail(details[0].Id, 1)
	if !ok {
		t.Error("Sale order detail could not be cancelled")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't retrieve the sale order details")
		return
	}

	if !details[0].Cancelled || details[0].QuantityInvoiced != details[0].Quantity || details[0].QuantityDeliveryNote != details[0].Quantity || details[0].Status != "Z" {
		t.Error("Cancelling the sale order detail has not updated the detail")
	}

	// uncancel the detail, check if uncancelled
	ok = cancelSalesOrderDetail(details[0].Id, 1)
	if !ok {
		t.Error("Sale order detail could not be uncancelled")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't retrieve the sale order details")
		return
	}

	if details[0].Cancelled || details[0].QuantityInvoiced != 0 || details[0].QuantityDeliveryNote != 0 || details[0].Status == "Z" {
		t.Error("Uncancelling the sale order detail has not updated the detail")
	}

	// attempt delete
	details[0].enterprise = 1
	ok = details[0].deleteSalesOrderDetail()
	if !ok {
		t.Error("Delete error, sale order detail not deleted")
		return
	}

	// check that the sale order has been updated correctly
	inMemoryOrder = getSalesOrderRow(orderId)
	if inMemoryOrder.LinesNumber != 0 {
		t.Error("The sale order number of lines is not updated upon delete")
		return
	}

	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder()
}

// ===== SALE INVOICES

/* GET */

func TestGetSalesInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}
	o := q.getSalesInvoices()

	for i := 0; i < len(o.Invoices); i++ {
		if o.Invoices[i].Id <= 0 {
			t.Error("Scan error, sale invoices with ID 0")
			return
		}
	}
}

func TestSearchSalesInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}}
	si := q.searchSalesInvoices()

	for i := 0; i < len(si.Invoices); i++ {
		if si.Invoices[i].Id <= 0 {
			t.Error("Scan error, sale invoices with ID 0.")
			return
		}
	}

	// search for ID
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "1"}}
	si = q.searchSalesInvoices()

	for i := 0; i < len(si.Invoices); i++ {
		if si.Invoices[i].Id <= 0 {
			t.Error("Scan error, sale invoices with ID 0.")
			return
		}
	}

	// search for customer name
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "Itzan"}}
	si = q.searchSalesInvoices()

	for i := 0; i < len(si.Invoices); i++ {
		if si.Invoices[i].Id <= 0 {
			t.Error("Scan error, sale invoices with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}, DateStart: &start, DateEnd: &end}
	si = q.searchSalesInvoices()

	for i := 0; i < len(si.Invoices); i++ {
		if si.Invoices[i].Id <= 0 {
			t.Error("Scan error, sale invoices with ID 0.")
			return
		}
	}
}

func TestGetSalesInvoiceRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesInvoiceRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, sale invoices row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestSaleInvoiceInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	i := SalesInvoice{
		Customer:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 1,
		enterprise:     1,
	}

	ok, invoiceId := i.insertSalesInvoice()
	if !ok {
		t.Error("Insert error, the invoice could not be inserted")
		return
	}

	i.Id = invoiceId
	ok = i.deleteSalesInvoice()
	if !ok {
		t.Error("Delete error, the invoice could not be deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestInvoiceAllSaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertSalesOrderDetail()

	ok := invoiceAllSaleOrder(orderId, 1)
	if !ok {
		t.Error("Could not invoice all sale order")
		return
	}

	// get invoice from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.Invoices) == 0 {
		t.Error("The invoice has not loaded from the sale order relations")
		return
	}

	invoice := getSalesInvoiceRow(r.Invoices[0].Id)
	order := getSalesOrderRow(orderId)

	// the totals in the order and in the invoice can not be different
	if invoice.TotalProducts != order.TotalProducts || invoice.TotalAmount != order.TotalAmount {
		t.Error("The totals in the order and in the invoice can not be different")
		return
	}

	if order.InvoicedLines == 0 {
		t.Error("The invoiced lines number in the order has not been updated")
		return
	}

	details := getSalesOrderDetail(orderId, 1)
	if details[0].QuantityInvoiced == 0 {
		t.Error("The quantity invoiced in the sale order detail has not been updated")
		return
	}

	// delete created invoice
	ok = invoice.deleteSalesInvoice()
	if !ok {
		t.Error("The invoice creted could not be deleted")
		return
	}

	// delete created order
	details[0].enterprise = 1
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder()
}

func TestIInvoicePartiallySaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertSalesOrderDetail()
	details := getSalesOrderDetail(orderId, 1)

	invoiceInfo := OrderDetailGenerate{
		OrderId: orderId,
		Selection: []OrderDetailGenerateSelection{
			{
				Id:       details[0].Id,
				Quantity: d.Quantity / 2,
			},
		},
	}
	ok := invoiceInfo.invoicePartiallySaleOrder(1)
	if !ok {
		t.Error("Could not invoice partially sale order")
		return
	}

	// get invoice from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.Invoices) == 0 {
		t.Error("The invoice has not loaded from the sale order relations")
		return
	}

	invoice := getSalesInvoiceRow(r.Invoices[0].Id)
	order := getSalesOrderRow(orderId)

	// the totals in the order and in the invoice can not be different
	if invoice.TotalProducts != order.TotalProducts/2 || invoice.TotalAmount != order.TotalAmount/2 {
		t.Error("The totals in the order and in the invoice can not be different")
		return
	}

	if order.InvoicedLines != 0 {
		t.Error("The invoiced lines number in the order has not been incorrectly updated")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].QuantityInvoiced == 0 {
		t.Error("The quantity invoiced in the sale order detail has not been updated")
		return
	}

	// delete created invoice
	ok = invoice.deleteSalesInvoice()
	if !ok {
		t.Error("The invoice creted could not be deleted")
		return
	}

	// delete created order
	details[0].enterprise = 1
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder()
}

func TestGetSalesInvoiceRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}
	o := q.getSalesInvoices()

	var checkOrders int8 = 0        // 0 = Not checked, 1 = OK, 2 = Error
	var checkDeliveryNotes int8 = 0 // 0 = Not checked, 1 = OK, 2 = Error

	for i := 0; i < len(o.Invoices); i++ {
		r := getSalesInvoiceRelations(o.Invoices[i].Id, 1)

		if checkOrders == 0 && len(r.Orders) > 0 {
			if r.Orders[0].Id <= 0 {
				checkOrders = 2
			} else {
				checkOrders = 1
			}
		}

		if checkDeliveryNotes == 0 && len(r.DeliveryNotes) > 0 {
			if r.DeliveryNotes[0].Id <= 0 {
				checkDeliveryNotes = 2
			} else {
				checkDeliveryNotes = 1
			}
		}

		if checkOrders != 0 || checkDeliveryNotes != 0 {
			break
		}
	}

	if checkOrders == 2 || checkDeliveryNotes == 2 {
		t.Errorf("Error scanning sale order relations checkOrders %q checkShippings %q", checkOrders, checkDeliveryNotes)
	}
}

// ===== SALE INVOICE DETAILS

/* GET */

func TestGetSalesInvoiceDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	d := getSalesInvoiceDetail(1, 1)

	for i := 0; i < len(d); i++ {
		if d[i].Id <= 0 {
			t.Error("Scan error, sale invoice details with ID 0.")
			return
		}
	}
}

func TestGetSalesInvoiceDetailRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesInvoiceDetailRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, sale invoice detail row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestSalesInvoiceDetailInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	i := SalesInvoice{
		Customer:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 1,
		enterprise:     1,
	}

	_, invoiceId := i.insertSalesInvoice()

	d := SalesInvoiceDetail{
		Invoice:    invoiceId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	ok := d.insertSalesInvoiceDetail(true)
	if !ok {
		t.Error("Insert error, sale invoice detail not inserted")
		return
	}

	// check if the totals are the same
	details := getSalesInvoiceDetail(invoiceId, 1)
	invoice := getSalesInvoiceRow(invoiceId)
	if invoice.TotalAmount != details[0].TotalAmount || invoice.TotalProducts != float64(details[0].Quantity)*details[0].Price {
		t.Error("The totals of the invoice are not updated correctly")
		return
	}

	// delete detail
	details = getSalesInvoiceDetail(invoiceId, 1)
	ok = details[0].deleteSalesInvoiceDetail()
	if !ok {
		t.Error("Delete error, sale invoice detail not deleted")
		return
	}

	// delete invoice
	i.Id = invoiceId
	i.deleteSalesInvoice()
}

// ===== SALE DELIVERY NOTE

/* GET */

func TestGetSalesDeliveryNotes(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}
	o := q.getSalesDeliveryNotes()

	for i := 0; i < len(o.Notes); i++ {
		if o.Notes[i].Id <= 0 {
			t.Error("Scan error, sale delivery note with ID 0.")
			return
		}
	}
}

func TestSearchSalesDelvieryNotes(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}}
	o := q.searchSalesDelvieryNotes()

	for i := 0; i < len(o.Notes); i++ {
		if o.Notes[i].Id <= 0 {
			t.Error("Scan error, sale delivery note with ID 0.")
			return
		}
	}

	// search for ID
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "1"}}
	o = q.searchSalesDelvieryNotes()

	for i := 0; i < len(o.Notes); i++ {
		if o.Notes[i].Id <= 0 {
			t.Error("Scan error, sale delivery note with ID 0.")
			return
		}
	}

	// search for customer name
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "Itzan"}}
	o = q.searchSalesDelvieryNotes()

	for i := 0; i < len(o.Notes); i++ {
		if o.Notes[i].Id <= 0 {
			t.Error("Scan error, sale delivery note with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}, DateStart: &start, DateEnd: &end}
	o = q.searchSalesDelvieryNotes()

	for i := 0; i < len(o.Notes); i++ {
		if o.Notes[i].Id <= 0 {
			t.Error("Scan error, sale delivery note with ID 0.")
			return
		}
	}
}

func TestGetSalesDeliveryNoteRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSalesDeliveryNoteRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, sale delivery note row with ID 0.")
		return
	}

}

func TestLocateSalesDeliveryNotesBySalesOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32}
	o := q.getSalesOrder(1)
	for i := len(o.Orders) - 1; i >= 0; i++ {
		notes := locateSalesDeliveryNotesBySalesOrder(o.Orders[i].Id, 1)
		if len(notes) > 0 {
			if notes[0].Id == 0 {
				t.Error("Scan error locating sale delivery notes")
				return
			} else {
				return
			}
		}
	}
}

func TestGetNameSalesDeliveryNote(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	name := getNameSalesDeliveryNote(1, 1)
	if name == "" {
		t.Error("Can't get the name of the sale delivery note")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestSalesDeliveryNoteInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	n := SalesDeliveryNote{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		ShippingAddress: 1,
		Currency:        1,
		enterprise:      1,
	}

	ok, noteId := n.insertSalesDeliveryNotes()
	if !ok || noteId <= 0 {
		t.Error("Insert error, delivey note not inserted")
		return
	}

	n.Id = noteId
	ok = n.deleteSalesDeliveryNotes()
	if !ok {
		t.Error("Delete error, delivey note not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestDeliveryNoteAllSaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertSalesOrderDetail()

	ok, noteId := deliveryNoteAllSaleOrder(orderId, 1)
	if !ok {
		t.Error("Could not delivery note all sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.DeliveryNotes) == 0 {
		t.Error("The note has not loaded from the sale order relations")
		return
	}

	note := getSalesDeliveryNoteRow(noteId)
	order := getSalesOrderRow(orderId)

	// the totals in the order and in the note can not be different
	if note.TotalProducts != order.TotalProducts || note.TotalAmount != order.TotalAmount {
		t.Error("The totals in the order and in the invoice can not be different")
		return
	}

	if order.DeliveryNoteLines == 0 {
		t.Error("The delivery note lines number in the order has not been updated")
		return
	}

	details := getSalesOrderDetail(orderId, 1)
	if details[0].QuantityDeliveryNote == 0 {
		t.Error("The quantity delivery note in the sale order detail has not been updated")
		return
	}

	// delete created delivery note
	ok = note.deleteSalesDeliveryNotes()
	if !ok {
		t.Error("The delivery note creted could not be deleted")
		return
	}

	// delete created order
	details[0].enterprise = 1
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder()
}

func TestDeliveryNotePartiallySaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertSalesOrderDetail()
	details := getSalesOrderDetail(orderId, 1)

	invoiceInfo := OrderDetailGenerate{
		OrderId: orderId,
		Selection: []OrderDetailGenerateSelection{
			{
				Id:       details[0].Id,
				Quantity: d.Quantity / 2,
			},
		},
	}
	ok := invoiceInfo.deliveryNotePartiallySaleOrder(1)
	if !ok {
		t.Error("Could not delivery note partially sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.DeliveryNotes) == 0 {
		t.Error("The delivery note has not loaded from the sale order relations")
		return
	}

	note := getSalesDeliveryNoteRow(r.DeliveryNotes[0].Id)
	order := getSalesOrderRow(orderId)

	// the totals in the order and in the note can not be different
	if note.TotalProducts != order.TotalProducts/2 || note.TotalAmount != order.TotalAmount/2 {
		t.Error("The totals in the order and in the note can not be different")
		return
	}

	if order.DeliveryNoteLines != 0 {
		t.Error("The delivery note lines number in the order has not been incorrectly updated")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].QuantityDeliveryNote == 0 {
		t.Error("The quantity delivery note in the sale order detail has not been updated")
		return
	}

	// delete created delivery note
	ok = note.deleteSalesDeliveryNotes()
	if !ok {
		t.Error("The delivery note creted could not be deleted")
		return
	}

	// delete created order
	details[0].enterprise = 1
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder()
}

func TestGetSalesDeliveryNoteRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}
	o := q.getSalesDeliveryNotes()

	var checkOrders int8 = 0    // 0 = Not checked, 1 = OK, 2 = Error
	var checkShippings int8 = 0 // 0 = Not checked, 1 = OK, 2 = Error

	for i := 0; i < len(o.Notes); i++ {
		r := getSalesDeliveryNoteRelations(o.Notes[i].Id, 1)

		if checkOrders == 0 && len(r.Orders) > 0 {
			if r.Orders[0].Id <= 0 {
				checkOrders = 2
			} else {
				checkOrders = 1
			}
		}

		if checkShippings == 0 && len(r.Shippings) > 0 {
			if r.Shippings[0].Id <= 0 {
				checkShippings = 2
			} else {
				checkShippings = 1
			}
		}

		if checkOrders != 0 || checkShippings != 0 {
			break
		}
	}

	if checkOrders == 2 || checkShippings == 2 {
		t.Errorf("Error scanning sale delivery note relations checkOrders %q checkShippings %q", checkOrders, checkShippings)
	}
}
