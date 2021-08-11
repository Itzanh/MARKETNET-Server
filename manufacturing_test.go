package main

import "testing"

// ===== MANUFACTURING ORDERS

/* GET */

func TestGetAllManufacturingOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getAllManufacturingOrders()

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, manufacturing orders with ID 0.")
			return
		}
	}
}

func TestGetManufacturingOrdersByType(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	o := getManufacturingOrdersByType(1)

	for i := 0; i < len(o); i++ {
		if o[i].Id <= 0 {
			t.Error("Scan error, manufacturing orders with ID 0.")
			return
		}
	}
}

func TestGetManufacturingOrderRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getManufacturingOrderRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, manufacturing order row with ID 0.")
		return
	}

}

/* INSERT - UPDATE - DELETE */

func TestManufacturingOrderInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	mo := ManufacturingOrder{
		Product:     1,
		Type:        1,
		UserCreated: 1,
	}
	ok := mo.insertManufacturingOrder()
	if !ok {
		t.Error("Insert error, manufacturing order not inserted")
		return
	}

	o := getAllManufacturingOrders()
	mo = o[0]

	ok = toggleManufactuedManufacturingOrder(mo.Id, 1)
	if !ok {
		t.Error("The manufacturing order can't be toggled")
		return
	}
	ok = manufacturingOrderTagPrinted(mo.Id, 1)
	if !ok {
		t.Error("The manufacturing order can't be printed")
		return
	}
	ok = toggleManufactuedManufacturingOrder(mo.Id, 1)
	if !ok {
		t.Error("The manufacturing order can't be toggled")
		return
	}

	ok = mo.deleteManufacturingOrder()
	if !ok {
		t.Error("Delete error, can't delete manufacturing order")
		return
	}
}

func TestManufacturingOrderAllSaleOrder(t *testing.T) {
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
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
	}

	d.insertSalesOrderDetail()

	invoiceAllSaleOrder(orderId)
	ok := manufacturingOrderAllSaleOrder(orderId, 1)
	if !ok {
		t.Error("Could not manufacturing order all sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId)

	if len(r.ManufacturingOrders) == 0 {
		t.Error("The manufacturing order has not loaded from the sale order relations")
		return
	}

	// delete created manufacturing orders
	for i := 0; i < len(r.ManufacturingOrders); i++ {
		ok = r.ManufacturingOrders[0].deleteManufacturingOrder()
		if !ok {
			t.Error("Delete error, can't delete manufacturing orders")
			return
		}
	}

	// delete created sale invoice
	r.Invoices[0].deleteSalesInvoice()

	// delete created order
	details := getSalesOrderDetail(orderId)
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.deleteSalesOrder()
}

func TestManufacturingOrderPartiallySaleOrder(t *testing.T) {
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
	}

	_, orderId := o.insertSalesOrder()

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
	}

	d.insertSalesOrderDetail()
	invoiceAllSaleOrder(orderId)

	details := getSalesOrderDetail(orderId)
	odg := OrderDetailGenerate{OrderId: orderId, Selection: []OrderDetailGenerateSelection{{Id: details[0].Id, Quantity: 1}}}
	ok := odg.manufacturingOrderPartiallySaleOrder(1)
	if !ok {
		t.Error("Could not manufacturing order all sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId)

	if len(r.ManufacturingOrders) == 0 {
		t.Error("The manufacturing order has not loaded from the sale order relations")
		return
	}

	// delete created manufacturing orders
	for i := 0; i < len(r.ManufacturingOrders); i++ {
		ok = r.ManufacturingOrders[0].deleteManufacturingOrder()
		if !ok {
			t.Error("Delete error, can't delete manufacturing orders")
			return
		}
	}

	// delete created sale invoice
	r.Invoices[0].deleteSalesInvoice()

	// delete created order
	details[0].deleteSalesOrderDetail()
	o.Id = orderId
	o.deleteSalesOrder()
}

// ===== CUSTOMERS

/* GET */

func TestGetManufacturingOrderType(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	mot := getManufacturingOrderType()

	for i := 0; i < len(mot); i++ {
		if mot[i].Id <= 0 {
			t.Error("Scan error, manufacturing order type with ID 0.")
			return
		}
	}
}

/* INSERT - UPDATE - DELETE */

func TestManufacturingOrderTypeInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	mot := ManufacturingOrderType{
		Name: "Test",
	}

	// insert
	ok := mot.insertManufacturingOrderType()
	if !ok {
		t.Error("Insert error, manufacturing order type not inserted")
		return
	}

	// update
	types := getManufacturingOrderType()
	mot = types[len(types)-1]

	mot.Name = "Test test"
	ok = mot.updateManufacturingOrderType()
	if !ok {
		t.Error("Update error, manufacturing order type not updated")
		return
	}

	// check update
	types = getManufacturingOrderType()
	mot = types[len(types)-1]
	if mot.Name != "Test test" {
		t.Error("Update error, manufacturing order type not successfully updated")
		return
	}

	// delete
	ok = mot.deleteManufacturingOrderType()
	if !ok {
		t.Error("Delete error, manufacturing order type not deleted")
		return
	}
}
