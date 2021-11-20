package main

import "testing"

// ===== PACKAGING

func TestPackaging(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create a sale order and two details (and a manual carrier)
	c := Carrier{
		Name:       "ACME Corp",
		MaxWeight:  35,
		MaxWidth:   150,
		MaxHeight:  150,
		MaxDepth:   150,
		Phone:      "987654321",
		Email:      "contact@acme.com",
		Web:        "acmecorp.com",
		Webservice: "_",
		enterprise: 1,
	}
	c.insertCarrier()
	carriers := getCariers(1)
	carrierId := carriers[len(carriers)-1].Id
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
		Carrier:         &carrierId,
		enterprise:      1,
	}
	_, orderId := o.insertSalesOrder(1)
	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}
	d.insertSalesOrderDetail(1)
	d = SalesOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      19.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}
	d.insertSalesOrderDetail(1)
	details := getSalesOrderDetail(orderId, 1)

	// create a package
	p := Packaging{
		Package:    1,
		SalesOrder: orderId,
		enterprise: 1,
	}
	ok := p.insertPackaging()
	if !ok {
		t.Error("Insert error, the packaging could not be inserted")
		return
	}

	// pack the detail
	orderPackaging := getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	detailPackaged := SalesOrderDetailPackaged{
		OrderDetail: details[0].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[0].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	// unpack de detail
	detailsPackaged := getSalesOrderDetailPackaged(orderPackaging[0].Id, 1)
	if len(detailsPackaged) == 0 || detailsPackaged[0].OrderDetail <= 0 {
		t.Error("Can't scan packed details")
		return
	}
	ok = detailsPackaged[0].deleteSalesOrderDetailPackaged(true, 0)
	if !ok {
		t.Error("Can't unpack a sale order detail")
		return
	}

	// create a second package, pack every detail in a separate package
	p = Packaging{
		Package:    1,
		SalesOrder: orderId,
		enterprise: 1,
	}
	ok = p.insertPackaging()
	if !ok {
		t.Error("Insert error, the packaging could not be inserted")
		return
	}
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}

	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[0].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[0].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}
	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[1].Id,
		Packaging:   orderPackaging[1].Id,
		Quantity:    details[1].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].QuantityPendingPackaging != 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}
	if details[1].QuantityPendingPackaging != 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}

	// check if deleting the package from the second detail unpacks de detail
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	ok = orderPackaging[1].deletePackaging(1, 0)
	if !ok {
		t.Error("Can't delete a package that contains a sale order detail")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[1].QuantityPendingPackaging == 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}

	// pack the second detail in the first package
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[1].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[1].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	// now that all the sale order is packaged, attempt to generate shipping
	ok = generateShippingFromSaleOrder(orderId, 1, 0)
	if !ok {
		t.Error("Could not generate shipping from sale order")
		return
	}

	orderRelations := getSalesOrderRelations(orderId, 1)
	if len(orderRelations.Shippings) == 0 || orderRelations.Shippings[0].Id <= 0 {
		t.Error("Can't scan shippings from the sale order relations")
		return
	}

	// ship the shipping! (set the shipping as sent)
	shipping := orderRelations.Shippings[0]
	ok = toggleShippingSent(shipping.Id, 1, 0).Ok
	if !ok {
		t.Error("Can't send shipping")
		return
	}

	// set the shipping as not sent
	ok = toggleShippingSent(shipping.Id, 1, 0).Ok
	if !ok {
		t.Error("Can't set shipping as not sent")
		return
	}

	// DELETE ALL
	// delete the shipping
	shipping.enterprise = 1
	ok = shipping.deleteShipping(0)
	if !ok {
		t.Error("Can't delete shipping")
		return
	}
	// delete the packages
	orderPackaging = getPackaging(orderId, 1)
	for i := 0; i < len(orderPackaging); i++ {
		ok = orderPackaging[i].deletePackaging(1, 0)
		if !ok {
			t.Error("Can't delete order packaging")
			return
		}
	}
	// delete the delivery note
	orderRelations = getSalesOrderRelations(orderId, 1)
	if len(orderRelations.DeliveryNotes) == 0 || orderRelations.DeliveryNotes[0].Id <= 0 {
		t.Error("Can't scan sale delivery notes")
		return
	}
	ok = orderRelations.DeliveryNotes[0].deleteSalesDeliveryNotes(0)
	if !ok {
		t.Error("Can't delete sale delivery note")
		return
	}
	// delete the details
	details = getSalesOrderDetail(orderId, 1)
	for i := 0; i < len(details); i++ {
		details[i].enterprise = 1
		ok = details[i].deleteSalesOrderDetail(1)
		if !ok {
			t.Error("Can't delete sale order detail")
			return
		}
	}
	// delete the sale order
	o.Id = orderId
	o.enterprise = 1
	ok = o.deleteSalesOrder(1)
	if !ok {
		t.Error("Can't delete sale order")
		return
	}
	// delete the carrier
	c.Id = carrierId
	ok = c.deleteCarrier()
	if !ok {
		t.Error("Can't delete carrier")
		return
	}
}

func TestPackagingWithPallets(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create a sale order and two details (and a manual carrier)
	c := Carrier{
		Name:       "ACME Corp",
		MaxWeight:  35,
		MaxWidth:   150,
		MaxHeight:  150,
		MaxDepth:   150,
		Phone:      "987654321",
		Email:      "contact@acme.com",
		Web:        "acmecorp.com",
		Webservice: "_",
		Pallets:    true,
		enterprise: 1,
	}
	c.insertCarrier()
	carriers := getCariers(1)
	carrierId := carriers[len(carriers)-1].Id
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
		Carrier:         &carrierId,
		enterprise:      1,
	}
	_, orderId := o.insertSalesOrder(1)
	d := SalesOrderDetail{
		Order:      orderId,
		Product:    4,
		Price:      9.99,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}
	d.insertSalesOrderDetail(1)
	d = SalesOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      19.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}
	d.insertSalesOrderDetail(1)
	details := getSalesOrderDetail(orderId, 1)

	// create a pallet
	pallet := Pallet{
		SalesOrder: orderId,
		Name:       "Pallet 1",
		enterprise: 1,
	}
	ok := pallet.insertPallet()
	if !ok {
		t.Error("Insert error, pallet not inserted")
		return
	}
	pallets := getSalesOrderPallets(orderId, 1)
	if len(pallets.Pallets) == 0 || pallets.Pallets[0].Id <= 0 {
		t.Error("Can't scan pallets")
		return
	}
	pallet.Id = pallets.Pallets[0].Id

	// create a package
	p := Packaging{
		Package:    1,
		SalesOrder: orderId,
		Pallet:     &pallet.Id,
		enterprise: 1,
	}
	ok = p.insertPackaging()
	if !ok {
		t.Error("Insert error, the packaging could not be inserted")
		return
	}

	// pack the detail
	orderPackaging := getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	detailPackaged := SalesOrderDetailPackaged{
		OrderDetail: details[0].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[0].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	// unpack de detail
	detailsPackaged := getSalesOrderDetailPackaged(orderPackaging[0].Id, 1)
	if len(detailsPackaged) == 0 || detailsPackaged[0].OrderDetail <= 0 {
		t.Error("Can't scan packed details")
		return
	}
	ok = detailsPackaged[0].deleteSalesOrderDetailPackaged(true, 0)
	if !ok {
		t.Error("Can't unpack a sale order detail")
		return
	}

	// create a second package, pack every detail in a separate package
	p = Packaging{
		Package:    1,
		SalesOrder: orderId,
		Pallet:     &pallet.Id,
		enterprise: 1,
	}
	ok = p.insertPackaging()
	if !ok {
		t.Error("Insert error, the packaging could not be inserted")
		return
	}
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}

	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[0].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[0].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}
	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[1].Id,
		Packaging:   orderPackaging[1].Id,
		Quantity:    details[1].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[0].QuantityPendingPackaging != 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}
	if details[1].QuantityPendingPackaging != 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}

	// check if deleting the package from the second detail unpacks de detail
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	ok = orderPackaging[1].deletePackaging(1, 0)
	if !ok {
		t.Error("Can't delete a package that contains a sale order detail")
		return
	}

	details = getSalesOrderDetail(orderId, 1)
	if details[1].QuantityPendingPackaging == 0 {
		t.Error("Quantity pending packaging is not being updated")
		return
	}

	// pack the second detail in the first package
	orderPackaging = getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	detailPackaged = SalesOrderDetailPackaged{
		OrderDetail: details[1].Id,
		Packaging:   orderPackaging[0].Id,
		Quantity:    details[1].Quantity,
		enterprise:  1,
	}
	ok = detailPackaged.insertSalesOrderDetailPackaged(0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging")
		return
	}

	// now that all the sale order is packaged, attempt to generate shipping
	ok = generateShippingFromSaleOrder(orderId, 1, 0)
	if !ok {
		t.Error("Could not generate shipping from sale order")
		return
	}

	orderRelations := getSalesOrderRelations(orderId, 1)
	if len(orderRelations.Shippings) == 0 || orderRelations.Shippings[0].Id <= 0 {
		t.Error("Can't scan shippings from the sale order relations")
		return
	}

	// ship the shipping! (set the shipping as sent)
	shipping := orderRelations.Shippings[0]
	ok = toggleShippingSent(shipping.Id, 1, 0).Ok
	if !ok {
		t.Error("Can't send shipping")
		return
	}

	// set the shipping as not sent
	ok = toggleShippingSent(shipping.Id, 1, 0).Ok
	if !ok {
		t.Error("Can't set shipping as not sent")
		return
	}

	// delete the packages
	orderPackaging = getPackaging(orderId, 1)
	for i := 0; i < len(orderPackaging); i++ {
		ok = orderPackaging[i].deletePackaging(1, 0)
		if !ok {
			t.Error("Can't delete order packaging")
			return
		}
	}

	// delete the pallet to delete all the packages
	ok = pallets.Pallets[0].deletePallet()
	if !ok {
		t.Error("Can't delete pallet")
		return
	}

	// DELETE ALL
	// delete the shipping
	shipping.enterprise = 1
	ok = shipping.deleteShipping(0)
	if !ok {
		t.Error("Can't delete shipping")
		return
	}
	// delete the delivery note
	orderRelations = getSalesOrderRelations(orderId, 1)
	if len(orderRelations.DeliveryNotes) == 0 || orderRelations.DeliveryNotes[0].Id <= 0 {
		t.Error("Can't scan sale delivery notes")
		return
	}
	ok = orderRelations.DeliveryNotes[0].deleteSalesDeliveryNotes(0)
	if !ok {
		t.Error("Can't delete sale delivery note")
		return
	}
	// delete the details
	details = getSalesOrderDetail(orderId, 1)
	for i := 0; i < len(details); i++ {
		details[i].enterprise = 1
		ok = details[i].deleteSalesOrderDetail(1)
		if !ok {
			t.Error("Can't delete sale order detail")
			return
		}
	}
	// delete the sale order
	o.Id = orderId
	o.enterprise = 1
	ok = o.deleteSalesOrder(1)
	if !ok {
		t.Error("Can't delete sale order")
		return
	}
	// delete the carrier
	c.Id = carrierId
	ok = c.deleteCarrier()
	if !ok {
		t.Error("Can't delete carrier")
		return
	}
}

func TestPackWithjEAN13(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create a sale order and two details
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
	_, orderId := o.insertSalesOrder(1)
	d := SalesOrderDetail{
		Order:      orderId,
		Product:    1,
		Price:      9.99,
		Quantity:   4,
		VatPercent: 21,
		enterprise: 1,
	}
	d.insertSalesOrderDetail(1)
	details := getSalesOrderDetail(orderId, 1)

	// create a package
	p := Packaging{
		Package:    1,
		SalesOrder: orderId,
		enterprise: 1,
	}
	ok := p.insertPackaging()
	if !ok {
		t.Error("Insert error, the packaging could not be inserted")
		return
	}

	// pack the detail
	orderPackaging := getPackaging(orderId, 1)
	if len(orderPackaging) == 0 || orderPackaging[0].Id <= 0 {
		t.Error("Can't scan sales order packaging")
		return
	}
	product := getProductRow(1)
	detailPackaged := SalesOrderDetailPackagedEAN13{
		SalesOrder: orderId,
		EAN13:      product.BarCode,
		Packaging:  orderPackaging[0].Id,
		Quantity:   details[0].Quantity,
	}
	ok = detailPackaged.insertSalesOrderDetailPackagedEAN13(1, 0)
	if !ok {
		t.Error("Can't pack a sale order detail inside a packaging by EAN13")
		return
	}

	// unpack de detail
	detailsPackaged := getSalesOrderDetailPackaged(orderPackaging[0].Id, 1)
	if len(detailsPackaged) == 0 || detailsPackaged[0].OrderDetail <= 0 {
		t.Error("Can't scan packed details")
		return
	}
	ok = detailsPackaged[0].deleteSalesOrderDetailPackaged(true, 0)
	if !ok {
		t.Error("Can't unpack a sale order detail")
		return
	}

	// delete the packages
	orderPackaging = getPackaging(orderId, 1)
	for i := 0; i < len(orderPackaging); i++ {
		ok = orderPackaging[i].deletePackaging(1, 0)
		if !ok {
			t.Error("Can't delete order packaging")
			return
		}
	}
	// delete the details
	details = getSalesOrderDetail(orderId, 1)
	for i := 0; i < len(details); i++ {
		details[i].enterprise = 1
		ok = details[i].deleteSalesOrderDetail(1)
		if !ok {
			t.Error("Can't delete sale order detail")
			return
		}
	}
	// delete the sale order
	o.Id = orderId
	o.enterprise = 1
	ok = o.deleteSalesOrder(1)
	if !ok {
		t.Error("Can't delete sale order")
		return
	}
}
