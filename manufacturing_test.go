/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"testing"
)

// ===== MANUFACTURING ORDERS

/* GET */

func TestGetAllManufacturingOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := ManufacturingPaginationQuery{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}, OrderTypeId: 0}
	o := q.getAllManufacturingOrders(1)

	for i := 0; i < len(o.ManufacturingOrders); i++ {
		if o.ManufacturingOrders[i].Id <= 0 {
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
	q := ManufacturingPaginationQuery{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}, OrderTypeId: 1}
	o := q.getManufacturingOrdersByType(1)

	for i := 0; i < len(o.ManufacturingOrders); i++ {
		if o.ManufacturingOrders[i].Id <= 0 {
			t.Error("Scan error, manufacturing orders with ID 0.")
			return
		}
	}
}

func TestGetManufacturingOrderRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getManufacturingOrderRow(2)
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
		ProductId:     1,
		TypeId:        1,
		UserCreatedId: 1,
		EnterpriseId:  1,
	}
	ok := mo.insertManufacturingOrder(0, nil).Ok
	if !ok {
		t.Error("Insert error, manufacturing order not inserted")
		return
	}

	q := ManufacturingPaginationQuery{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1}, OrderTypeId: 0}
	o := q.getAllManufacturingOrders(1)
	mo = o.ManufacturingOrders[0]

	ok = toggleManufactuedManufacturingOrder(mo.Id, 1, 1)
	if !ok {
		t.Error("The manufacturing order can't be toggled")
		return
	}
	ok = manufacturingOrderTagPrinted(mo.Id, 1, 1)
	if !ok {
		t.Error("The manufacturing order can't be printed")
		return
	}
	ok = toggleManufactuedManufacturingOrder(mo.Id, 1, 1)
	if !ok {
		t.Error("The manufacturing order can't be toggled")
		return
	}

	ok = mo.deleteManufacturingOrder(0, nil)
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
		CustomerId:        1,
		PaymentMethodId:   3,
		BillingSeriesId:   "EXP",
		CurrencyId:        1,
		BillingAddressId:  1,
		ShippingAddressId: 1,
		Description:       "",
		Notes:             "",
		EnterpriseId:      1,
	}

	_, orderId := o.insertSalesOrder(1)

	d := SalesOrderDetail{
		OrderId:      orderId,
		ProductId:    1,
		Price:        9.99,
		Quantity:     2,
		VatPercent:   21,
		EnterpriseId: 1,
	}

	d.insertSalesOrderDetail(1)

	invoiceAllSaleOrder(orderId, 1, 0)
	ok := manufacturingOrderAllSaleOrder(orderId, 1, 1)
	if !ok {
		t.Error("Could not manufacturing order all sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.ManufacturingOrders) == 0 {
		t.Error("The manufacturing order has not loaded from the sale order relations")
		return
	}

	// delete created manufacturing orders
	for i := 0; i < len(r.ManufacturingOrders); i++ {
		ok = r.ManufacturingOrders[i].deleteManufacturingOrder(0, nil)
		if !ok {
			t.Error("Delete error, can't delete manufacturing orders")
			return
		}
	}

	// delete created sale invoice
	r.Invoices[0].deleteSalesInvoice(0)

	// delete created order
	details := getSalesOrderDetail(orderId, 1)
	details[0].EnterpriseId = 1
	details[0].deleteSalesOrderDetail(1, nil)
	o.Id = orderId
	o.EnterpriseId = 1
	o.deleteSalesOrder(1)
}

func TestManufacturingOrderPartiallySaleOrder(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		CustomerId:        1,
		PaymentMethodId:   3,
		BillingSeriesId:   "EXP",
		CurrencyId:        1,
		BillingAddressId:  1,
		ShippingAddressId: 1,
		Description:       "",
		Notes:             "",
		EnterpriseId:      1,
	}

	_, orderId := o.insertSalesOrder(1)

	d := SalesOrderDetail{
		OrderId:      orderId,
		ProductId:    1,
		Price:        9.99,
		Quantity:     2,
		VatPercent:   21,
		EnterpriseId: 1,
	}

	d.insertSalesOrderDetail(1)
	invoiceAllSaleOrder(orderId, 1, 0)

	details := getSalesOrderDetail(orderId, 1)
	odg := ManufacturingOrderGenerate{Selection: []ManufacturingOrderGenerateSelection{
		{Id: details[0].Id, Quantity: 1, OrderId: orderId},
	}}
	ok := odg.manufacturingOrderPartiallySaleOrder(1, 1)
	if !ok {
		t.Error("Could not manufacturing order all sale order")
		return
	}

	// get note from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.ManufacturingOrders) == 0 {
		t.Error("The manufacturing order has not loaded from the sale order relations")
		return
	}

	// delete created manufacturing orders
	for i := 0; i < len(r.ManufacturingOrders); i++ {
		ok = r.ManufacturingOrders[i].deleteManufacturingOrder(0, nil)
		if !ok {
			t.Error("Delete error, can't delete manufacturing orders")
			return
		}
	}

	// delete created sale invoice
	r.Invoices[0].deleteSalesInvoice(0)

	// delete created order
	details[0].EnterpriseId = 1
	details[0].deleteSalesOrderDetail(1, nil)
	o.Id = orderId
	o.EnterpriseId = 1
	o.deleteSalesOrder(1)
}

func TestManufacturingOrderQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(1)

	p := Product{
		Name:                     "Glass Office Desk",
		Reference:                "OF-DSK",
		BarCode:                  "1234067891236",
		ControlStock:             true,
		Weight:                   30,
		FamilyId:                 &family,
		Width:                    160,
		Height:                   100,
		Depth:                    40,
		VatPercent:               21,
		Price:                    65,
		Manufacturing:            true,
		ManufacturingOrderTypeId: &manufacturingOrderType,
		TrackMinimumStock:        true,
		EnterpriseId:             1,
	}

	okAndErr := p.insertProduct(0)
	if !okAndErr.Ok {
		t.Error("Insert error, could not insert product", okAndErr.ErrorCode)
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	o := SaleOrder{
		CustomerId:        1,
		PaymentMethodId:   3,
		BillingSeriesId:   "EXP",
		CurrencyId:        1,
		BillingAddressId:  1,
		ShippingAddressId: 1,
		Description:       "",
		Notes:             "",
		EnterpriseId:      1,
	}

	_, orderId := o.insertSalesOrder(0)

	d := SalesOrderDetail{
		OrderId:      orderId,
		ProductId:    p.Id,
		Price:        9.99,
		Quantity:     2,
		VatPercent:   21,
		EnterpriseId: 1,
	}

	d.insertSalesOrderDetail(0)

	invoiceAllSaleOrder(orderId, 1, 0)

	details := getSalesOrderDetail(orderId, 1)
	if details[0].Status != "C" {
		t.Error("The status is not correct when manufacturing orders are not generated yet")
		return
	}

	ok := manufacturingOrderAllSaleOrder(orderId, 1, 1)
	if !ok {
		t.Error("Could not manufacturing order all sale order")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "D" {
		t.Error("The status is not correct when manufactured 0/2")
	}

	// get orders from the sale order relations
	r := getSalesOrderRelations(orderId, 1)

	if len(r.ManufacturingOrders) == 0 {
		t.Error("The manufacturing order has not loaded from the sale order relations")
		return
	}

	// set the first as manufactured
	if !toggleManufactuedManufacturingOrder(r.ManufacturingOrders[0].Id, 1, 1) {
		t.Error("Can't set a manufacturing order as manufactured: toggle error")
		return
	}
	r = getSalesOrderRelations(orderId, 1)
	if r.ManufacturingOrders[0].Manufactured == false {
		t.Error("Can't set a manufacturing order as manufactured")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "D" {
		t.Error("The status is not correct when manufactured 1/2")
	}

	// set the second as manufactured
	toggleManufactuedManufacturingOrder(r.ManufacturingOrders[1].Id, 1, 1)
	r = getSalesOrderRelations(orderId, 1)
	if r.ManufacturingOrders[1].Manufactured == false {
		t.Error("Can't set a manufacturing order as manufactured")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "E" {
		t.Error("The status is not correct when manufactured 2/2")
	}

	// set the second as NOT manufactured
	toggleManufactuedManufacturingOrder(r.ManufacturingOrders[1].Id, 1, 1)
	r = getSalesOrderRelations(orderId, 1)
	if r.ManufacturingOrders[1].Manufactured == true {
		t.Error("Can't set a manufacturing order as NOT manufactured")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "D" {
		t.Error("The status is not correct when manufactured 1/2")
	}

	// set the first as NOT manufactured
	toggleManufactuedManufacturingOrder(r.ManufacturingOrders[0].Id, 1, 1)
	r = getSalesOrderRelations(orderId, 1)
	if r.ManufacturingOrders[0].Manufactured == true {
		t.Error("Can't set a manufacturing order as NOT manufactured")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "D" {
		t.Error("The status is not correct when manufactured 0/2")
	}

	// delete created manufacturing orders
	for i := 0; i < len(r.ManufacturingOrders); i++ {
		ok = r.ManufacturingOrders[i].deleteManufacturingOrder(0, nil)
		if !ok {
			t.Error("Delete error, can't delete manufacturing orders")
			return
		}
	}

	// delete created sale invoice
	r.Invoices[0].deleteSalesInvoice(0)

	// delete created order
	details = getSalesOrderDetail(orderId, 1)
	details[0].EnterpriseId = 1
	details[0].deleteSalesOrderDetail(0, nil)
	o.Id = orderId
	o.EnterpriseId = 1
	o.deleteSalesOrder(0)

	ok = p.deleteProduct(0).Ok
	if !ok {
		t.Error("Delete error, could not delete product")
		return
	}
}

func TestAsignSaleOrderToManufacturingOrderForStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(1)

	p := Product{
		Name:                     "Glass Office Desk",
		Reference:                "OF-DSK",
		BarCode:                  "1234067891236",
		ControlStock:             true,
		Weight:                   30,
		FamilyId:                 &family,
		Width:                    160,
		Height:                   100,
		Depth:                    40,
		VatPercent:               21,
		Price:                    65,
		Manufacturing:            true,
		ManufacturingOrderTypeId: &manufacturingOrderType,
		TrackMinimumStock:        true,
		EnterpriseId:             1,
	}

	okAndErr := p.insertProduct(0)
	if !okAndErr.Ok {
		t.Error("Insert error, could not insert product", okAndErr.ErrorCode)
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	// creare manufacturing order
	mo := ManufacturingOrder{
		ProductId:     p.Id,
		UserCreatedId: 1,
		Complex:       false,
		EnterpriseId:  1,
	}
	okAndErr = mo.insertManufacturingOrder(0, nil)
	if !okAndErr.Ok {
		t.Error("Insert error, could not insert manufacturing order", okAndErr.ErrorCode)
		return
	}

	// create sale order
	o := SaleOrder{
		CustomerId:        1,
		PaymentMethodId:   3,
		BillingSeriesId:   "EXP",
		CurrencyId:        1,
		BillingAddressId:  1,
		ShippingAddressId: 1,
		Description:       "",
		Notes:             "",
		EnterpriseId:      1,
	}

	_, orderId := o.insertSalesOrder(0)

	d := SalesOrderDetail{
		OrderId:      orderId,
		ProductId:    p.Id,
		Price:        9.99,
		Quantity:     1,
		VatPercent:   21,
		EnterpriseId: 1,
	}

	d.insertSalesOrderDetail(0)

	invoiceAllSaleOrder(orderId, 1, 0)

	// has the order associated with the manufacturing order
	details := getSalesOrderDetail(orderId, 1)
	if details[0].Status != "D" {
		t.Error("Error status ", details[0].Status)
		return
	}
	mo = getManufacturingOrderRow(mo.Id)
	if mo.OrderDetailId == nil {
		t.Error("Order not associated")
		return
	}

	// delete manufacturing order
	ok := mo.deleteManufacturingOrder(1, nil)
	if !ok {
		t.Error("Delete error, could not delete manufacturing order")
		return
	}

	// delete created sale invoice
	r := getSalesOrderRelations(orderId, 1)
	r.Invoices[0].deleteSalesInvoice(0)

	// delete created order
	details = getSalesOrderDetail(orderId, 1)
	details[0].EnterpriseId = 1
	details[0].deleteSalesOrderDetail(0, nil)
	o.Id = orderId
	o.EnterpriseId = 1
	o.deleteSalesOrder(0)

	// delete product
	okAndErr = p.deleteProduct(0)
	if !okAndErr.Ok {
		t.Error("Delete error, could not delete product", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}

// ===== MANUFACTURING ORDER TYPE

/* GET */

func TestGetManufacturingOrderType(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	mot := getManufacturingOrderType(1)

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
		Name:                 "Test",
		EnterpriseId:         1,
		QuantityManufactured: 1,
	}

	// insert
	ok := mot.insertManufacturingOrderType()
	if !ok {
		t.Error("Insert error, manufacturing order type not inserted")
		return
	}

	// update
	types := getManufacturingOrderType(1)
	mot = types[len(types)-1]

	mot.Name = "Test test"
	ok = mot.updateManufacturingOrderType()
	if !ok {
		t.Error("Update error, manufacturing order type not updated")
		return
	}

	// check update
	types = getManufacturingOrderType(1)
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
