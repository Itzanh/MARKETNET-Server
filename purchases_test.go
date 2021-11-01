package main

import (
	"testing"
	"time"
)

// ===== PURCHASE ORDER

/* GET */

func TestGetPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrder(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase orders with ID 0.")
			return
		}
	}
}

func TestSearchPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}}
	o := q.searchPurchaseOrder()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase orders with ID 0.")
			return
		}
	}

	// search for ID
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "1"}}
	o = q.searchPurchaseOrder()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase orders with ID 0.")
			return
		}
	}

	// search for customer name
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "Itzan"}}
	o = q.searchPurchaseOrder()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase orders with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}, DateStart: &start, DateEnd: &end}
	o = q.searchPurchaseOrder()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase orders with ID 0.")
			return
		}
	}
}

func TestGetPurchaseOrderRowr(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrderRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, purchase order row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestPurchaseOrderInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	ok, orderId := o.insertPurchaseOrder()
	if !ok || orderId <= 0 {
		t.Error("Insert error, purchase order not inserted")
		return
	}

	o.Id = orderId
	ok = o.deletePurchaseOrder()
	if !ok {
		t.Error("Delete error, purchase order not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestGetPurchaseOrderDefaults(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	d := getPurchaseOrderDefaults(1)
	if len(d.Warehouse) == 0 || len(d.WarehouseName) == 0 {
		t.Error("Purchase order defaults not loaded")
		return
	}
}

func TestGetPurchaseOrderRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrder(1)

	var checkInvoices int8 = 0      // 0 = Not checked, 1 = OK, 2 = Error
	var checkDeliveryNotes int8 = 0 // 0 = Not checked, 1 = OK, 2 = Error

	for i := 0; i < len(o); i++ {
		r := getPurchaseOrderRelations(o[i].Id, 1)

		if checkInvoices == 0 && len(r.Invoices) > 0 {
			if r.Invoices[0].Id <= 0 {
				checkInvoices = 2
			} else {
				checkInvoices = 1
			}
		}

		if checkDeliveryNotes == 0 && len(r.DeliveryNotes) > 0 {
			if r.DeliveryNotes[0].Id <= 0 {
				checkDeliveryNotes = 2
			} else {
				checkDeliveryNotes = 1
			}
		}

		if checkInvoices != 0 || checkDeliveryNotes != 0 {
			break
		}
	}

	if checkInvoices == 2 || checkDeliveryNotes == 2 {
		t.Errorf("Error scanning purchase order relations checkInvoices %q checkDeliveryNotes %q", checkInvoices, checkDeliveryNotes)
	}
}

// ===== PURCHASE ORDER DETAILS

/* GET */

func TestGetPurchaseOrderDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrderDetail(1, 1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase order detail with ID 0.")
			return
		}
	}
}

func TestGetPurchaseOrderDetailRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrderDetailRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, purchase order detail row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestPurchaseOrderDetailInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, orderId := o.insertPurchaseOrder()

	d := PurchaseOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      15,
		Quantity:   15,
		VatPercent: 21,
		enterprise: 1,
	}

	// insert
	ok, detailId := d.insertPurchaseOrderDetail(true)
	if !ok || detailId <= 0 {
		t.Error("Insert error, purchase order detail not inserted")
		return
	}

	// check totals
	order := getPurchaseOrderRow(orderId)
	detail := getPurchaseOrderDetailRow(detailId)

	if order.TotalAmount != detail.TotalAmount || order.TotalProducts != detail.Price*float64(detail.Quantity) {
		t.Error("Purchase order totals not updated succlessfully")
		return
	}

	// check loading order details
	details := getPurchaseOrderDetail(orderId, 1)
	if details[0].Id <= 0 {
		t.Error("Purchase order detail not scanneed successfully")
		return
	}

	// delete
	d.Id = detailId
	ok = d.deletePurchaseOrderDetail()
	if !ok {
		t.Error("Delete error, purchase order not deleted")
		return
	}

	o.Id = orderId
	o.deletePurchaseOrder()
}

/* FUNCTIONALITY */

func TestGetSalesOrderDetailsFromPurchaseOrderDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseOrder(1)
	for i := 0; i < len(o); i++ {
		d := getPurchaseOrderDetail(o[i].Id, 1)
		for j := 0; j < len(d); j++ {
			salesDetails := getSalesOrderDetailsFromPurchaseOrderDetail(d[j].Id, 1)
			if len(salesDetails) > 0 {
				if salesDetails[0].Id <= 0 {
					t.Error("Sales order details from purchase order details not scanned correctyl")
					return
				} else {
					return
				}
			}
		}
	}
}

// ===== PURCHASE INVOICE

/* GET */

func TestGetPurchaseInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseInvoices(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoices with ID 0.")
			return
		}
	}
}

func TestSearchPurchaseInvoice(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}}
	o := q.searchPurchaseInvoice()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoices with ID 0.")
			return
		}
	}

	// search for ID
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "1"}}
	o = q.searchPurchaseInvoice()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoices with ID 0.")
			return
		}
	}

	// search for customer name
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "Itzan"}}
	o = q.searchPurchaseInvoice()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoices with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}, DateStart: &start, DateEnd: &end}
	o = q.searchPurchaseInvoice()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoices with ID 0.")
			return
		}
	}
}

func TestGetPurchaseInvoiceRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseInvoiceRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, purchase invoice row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestPurchaseInvoiceInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	i := PurchaseInvoice{
		Supplier:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 3,
		enterprise:     1,
	}

	ok, invoiceId := i.insertPurchaseInvoice()
	if !ok || invoiceId <= 0 {
		t.Error("Insert error, can't insert purchase invocice")
		return
	}

	i.Id = invoiceId
	ok = i.deletePurchaseInvoice()
	if !ok {
		t.Error("Delete error, can't delete purchase invocice")
		return
	}
}

/* FUNCTIONALITY */

func TestInvoiceAllPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, orderId := o.insertPurchaseOrder()

	d := PurchaseOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      15,
		Quantity:   15,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertPurchaseOrderDetail(true)

	ok := invoiceAllPurchaseOrder(orderId, 1)
	if !ok {
		t.Error("Could not invoice all purchase order")
		return
	}

	// get invoice from the purchase order relations
	r := getPurchaseOrderRelations(orderId, 1)

	if len(r.Invoices) == 0 {
		t.Error("The invoice has not loaded from the purchase order relations")
		return
	}

	invoice := getPurchaseInvoiceRow(r.Invoices[0].Id)
	order := getPurchaseOrderRow(orderId)

	// the totals in the order and in the invoice can not be different
	if invoice.TotalProducts != order.TotalProducts || invoice.TotalAmount != order.TotalAmount {
		t.Error("The totals in the order and in the invoice can not be different")
		return
	}

	if order.InvoicedLines == 0 {
		t.Error("The invoiced lines number in the order has not been updated")
		return
	}

	details := getPurchaseOrderDetail(orderId, 1)
	if details[0].QuantityInvoiced == 0 {
		t.Error("The quantity invoiced in the sale order detail has not been updated")
		return
	}

	// delete created invoice
	ok = invoice.deletePurchaseInvoice()
	if !ok {
		t.Error("The invoice creted could not be deleted")
		return
	}

	// delete created order
	details[0].deletePurchaseOrderDetail()
	o.Id = orderId
	o.deletePurchaseOrder()
}

func TestIInvoicePartiallyPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, orderId := o.insertPurchaseOrder()

	d := PurchaseOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      15,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertPurchaseOrderDetail(true)
	details := getPurchaseOrderDetail(orderId, 1)

	invoiceInfo := OrderDetailGenerate{
		OrderId: orderId,
		Selection: []OrderDetailGenerateSelection{
			{
				Id:       details[0].Id,
				Quantity: d.Quantity / 2,
			},
		},
	}
	ok := invoiceInfo.invoicePartiallyPurchaseOrder(1)
	if !ok {
		t.Error("Could not invoice partially purchase order")
		return
	}

	// get invoice from the purchase order relations
	r := getPurchaseOrderRelations(orderId, 1)

	if len(r.Invoices) == 0 {
		t.Error("The invoice has not loaded from the purchase order relations")
		return
	}

	invoice := getPurchaseInvoiceRow(r.Invoices[0].Id)
	order := getPurchaseOrderRow(orderId)

	// the totals in the order and in the invoice can not be different
	if invoice.TotalProducts != order.TotalProducts/2 || invoice.TotalAmount != order.TotalAmount/2 {
		t.Error("The totals in the order and in the invoice can not be different")
		return
	}

	if order.InvoicedLines != 0 {
		t.Error("The invoiced lines number in the order has not been incorrectly updated")
		return
	}

	details = getPurchaseOrderDetail(orderId, 1)
	if details[0].QuantityInvoiced == 0 {
		t.Error("The quantity invoiced in the purchase order detail has not been updated")
		return
	}

	// delete created invoice
	ok = invoice.deletePurchaseInvoice()
	if !ok {
		t.Error("The invoice creted could not be deleted")
		return
	}

	// delete created order
	details[0].deletePurchaseOrderDetail()
	o.Id = orderId
	o.deletePurchaseOrder()
}

func TestGetPurchaseInvoiceRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseInvoices(1)

	var checkOrders int8 = 0   // 0 = Not checked, 1 = OK, 2 = Error
	var checkInvoices int8 = 0 // 0 = Not checked, 1 = OK, 2 = Error

	for i := 0; i < len(o); i++ {
		r := getPurchaseInvoiceRelations(o[i].Id, 1)

		if checkOrders == 0 && len(r.Orders) > 0 {
			if r.Orders[0].Id <= 0 {
				checkOrders = 2
			} else {
				checkOrders = 1
			}
		}

		if checkInvoices == 0 && len(r.Invoices) > 0 {
			if r.Invoices[0].Id <= 0 {
				checkInvoices = 2
			} else {
				checkInvoices = 1
			}
		}

		if checkOrders != 0 || checkInvoices != 0 {
			break
		}
	}

	if checkOrders == 2 || checkInvoices == 2 {
		t.Errorf("Error scanning sale order relations checkOrders %q", checkOrders)
	}
}

// ===== PURCHASE INVOICE DETAILS

/* GET */

func TestGetPurchaseInvoiceDetail(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseInvoiceDetail(1, 1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase invoice details with ID 0.")
			return
		}
	}
}

func TestGetPurchaseInvoiceDetailRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseInvoiceDetailRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, purchase invoice detail row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestPurchaseInvoiceDetailInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	var product int32 = 3

	i := PurchaseInvoice{
		Supplier:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 3,
		enterprise:     1,
	}

	_, invoiceId := i.insertPurchaseInvoice()

	d := PurchaseInvoiceDetail{
		Invoice:    invoiceId,
		Product:    &product,
		Price:      15,
		Quantity:   3,
		VatPercent: 21,
		enterprise: 1,
	}

	// insert
	ok := d.insertPurchaseInvoiceDetail(true)
	if !ok {
		t.Error("Insert error, can't insert purchase invocice detail")
		return
	}

	// check the totals
	invoice := getPurchaseInvoiceRow(invoiceId)
	details := getPurchaseInvoiceDetail(invoiceId, 1)
	if invoice.TotalAmount != details[0].TotalAmount || invoice.TotalProducts != float64(details[0].Quantity)*details[0].Price {
		t.Error("The total of the invoice has not been updated successfully")
		return
	}

	// delete
	i.Id = invoiceId
	ok = i.deletePurchaseInvoice()
	if !ok {
		t.Error("Delete error, can't delete purchase invocice detail")
		return
	}
}

// ===== PURCHASE DELIVERY NOTE

/* GET */

func TestGetPurchaseDeliveryNote(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseDeliveryNotes(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase delivery note with ID 0.")
			return
		}
	}
}

func TestSearchPurchaseDeliveryNote(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}}
	o := q.searchPurchaseDeliveryNote()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase delivery note with ID 0.")
			return
		}
	}

	// search for ID
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "1"}}
	o = q.searchPurchaseDeliveryNote()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase delivery note with ID 0.")
			return
		}
	}

	// search for customer name
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, Search: "Itzan"}}
	o = q.searchPurchaseDeliveryNote()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase delivery note with ID 0.")
			return
		}
	}

	// search with date
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	q = OrderSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}}, DateStart: &start, DateEnd: &end}
	o = q.searchPurchaseDeliveryNote()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, purchase delivery note with ID 0.")
			return
		}
	}
}

func TestGetPurchaseDeliveryNoteRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getPurchaseDeliveryNoteRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, purchase delivery note row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestPurchaseDeliveryNoteInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	i := PurchaseDeliveryNote{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		ShippingAddress: 3,
		enterprise:      1,
	}

	ok, noteId := i.insertPurchaseDeliveryNotes()
	if !ok || noteId <= 0 {
		t.Error("Insert error, can't insert purchase delivery note")
		return
	}

	i.Id = noteId
	ok = i.deletePurchaseDeliveryNotes()
	if !ok {
		t.Error("Delete error, can't delete purchase delivery note")
		return
	}
}

/* FUNCTIONALITY */

func TestDeliveryNoteAllPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, orderId := o.insertPurchaseOrder()

	d := PurchaseOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      15,
		Quantity:   15,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertPurchaseOrderDetail(true)

	ok, noteId := deliveryNoteAllPurchaseOrder(orderId, 1)
	if !ok || noteId <= 0 {
		t.Error("Could not delivery note all purchase order")
		return
	}

	// get delivery note from the purchase order relations
	r := getPurchaseOrderRelations(orderId, 1)

	if len(r.DeliveryNotes) == 0 {
		t.Error("The delivery note has not loaded from the purchase order relations")
		return
	}

	note := getPurchaseDeliveryNoteRow(noteId)
	order := getPurchaseOrderRow(orderId)

	// the totals in the order and in the delivery note can not be different
	if note.TotalProducts != order.TotalProducts || note.TotalAmount != order.TotalAmount {
		t.Error("The totals in the order and in the delivery note can not be different")
		return
	}

	if order.DeliveryNoteLines == 0 {
		t.Error("The delivery note lines number in the order has not been updated")
		return
	}

	details := getPurchaseOrderDetail(orderId, 1)
	if details[0].QuantityDeliveryNote == 0 {
		t.Error("The quantity delivery note in the sale order detail has not been updated")
		return
	}

	// delete created delivery note
	ok = note.deletePurchaseDeliveryNotes()
	if !ok {
		t.Error("The delivery note creted could not be deleted")
		return
	}

	// delete created order
	details[0].deletePurchaseOrderDetail()
	o.Id = orderId
	o.deletePurchaseOrder()
}

func TestDeliveryNotePartiallyPurchaseOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := PurchaseOrder{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		BillingAddress:  3,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, orderId := o.insertPurchaseOrder()

	d := PurchaseOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      15,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}

	d.insertPurchaseOrderDetail(true)
	details := getPurchaseOrderDetail(orderId, 1)

	invoiceInfo := OrderDetailGenerate{
		OrderId: orderId,
		Selection: []OrderDetailGenerateSelection{
			{
				Id:       details[0].Id,
				Quantity: d.Quantity / 2,
			},
		},
	}
	ok := invoiceInfo.deliveryNotePartiallyPurchaseOrder(1)
	if !ok {
		t.Error("Could not deliveryNote partially purchase order")
		return
	}

	// get delivery note from the purchase order relations
	r := getPurchaseOrderRelations(orderId, 1)

	if len(r.DeliveryNotes) == 0 {
		t.Error("The delivery note has not loaded from the purchase order relations")
		return
	}

	note := getPurchaseDeliveryNoteRow(r.DeliveryNotes[0].Id)
	order := getPurchaseOrderRow(orderId)

	// the totals in the order and in the delivery note can not be different
	if note.TotalProducts != order.TotalProducts/2 || note.TotalAmount != order.TotalAmount/2 {
		t.Error("The totals in the order and in the delivery note can not be different")
		return
	}

	if order.DeliveryNoteLines != 0 {
		t.Error("The delivery note lines number in the order has not been incorrectly updated")
		return
	}

	details = getPurchaseOrderDetail(orderId, 1)
	if details[0].QuantityDeliveryNote == 0 {
		t.Error("The quantity delivery note in the purchase order detail has not been updated")
		return
	}

	// delete created delivery note
	ok = note.deletePurchaseDeliveryNotes()
	if !ok {
		t.Error("The delivery note creted could not be deleted")
		return
	}

	// delete created order
	details[0].deletePurchaseOrderDetail()
	o.Id = orderId
	o.deletePurchaseOrder()
}
