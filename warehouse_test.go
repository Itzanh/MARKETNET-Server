package main

import "testing"

// ===== WAREHOUSE

/* GET */

func TestGetWarehouses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	w := getWarehouses(1)

	for i := 0; i < len(w); i++ {
		if len(w[i].Id) == 0 {
			t.Error("Scan error, warehouses with ID 0.")
			return
		}
	}
}

/* INSERT - UPDATE - DELETE */

func TestWarehouseInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	w := Warehouse{
		Id:           "WA",
		Name:         "Test warehouse",
		EnterpriseId: 1,
	}

	// insert
	ok := w.insertWarehouse()
	if !ok {
		t.Error("Insert error, warehouse not inserted")
		return
	}

	// update
	warehouses := getWarehouses(1)
	for i := 0; i < len(warehouses); i++ {
		if warehouses[i].Id == w.Id {
			w = warehouses[i]
			break
		}
	}

	w.Name = "Test test"
	ok = w.updateWarehouse()
	if !ok {
		t.Error("Update error, warehouse not updated")
		return
	}

	// check update
	warehouses = getWarehouses(1)
	for i := 0; i < len(warehouses); i++ {
		if warehouses[i].Id == w.Id {
			w = warehouses[i]
			break
		}
	}

	if w.Name != "Test test" {
		t.Error("Update error, warehouse not successfully updated")
		return
	}

	// delete
	ok = w.deleteWarehouse()
	if !ok {
		t.Error("Delete error, warehouse not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindWarehouseByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	warehouses := findWarehouseByName("", 1)
	if len(warehouses) == 0 || len(warehouses[0].Id) == 0 {
		t.Error("Can't scan warehouses")
		return
	}
}

func TestGetNameWarehouse(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	warehouseName := getNameWarehouse("W1", 1)
	if warehouseName == "" {
		t.Error("Can't get the name of the warehouse")
		return
	}
}

func TestRegenerateProductStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	ok := regenerateProductStock(1)
	if !ok {
		t.Error("There was an error regenerating the stock of the products")
		return
	}
}

// ===== STOCK

/* GET */

func TestGetStocks(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getStock(1, 1)

	for i := 0; i < len(s); i++ {
		if s[i].ProductId <= 0 {
			t.Error("Scan error, stock with ID 0.")
			return
		}
	}
}

func TestGetStockRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getStockRow(1, "W1", 1)
	if s.ProductId <= 0 {
		t.Error("Scan error, customer row with ID 0.")
		return
	}

}

/* INSERT */

func TestCreateStockRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(2)
	supplier := int32(1)

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
		SupplierId:               &supplier,
		TrackMinimumStock:        true,
		EnterpriseId:             1,
	}

	ok := p.insertProduct(0).Ok
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	w := Warehouse{
		Id:           "WA",
		Name:         "Test warehouse",
		EnterpriseId: 1,
	}

	// insert
	ok = w.insertWarehouse()
	if !ok {
		t.Error("Insert error, warehouse not inserted")
		return
	}
	warehouses := getWarehouses(1)
	for i := 0; i < len(warehouses); i++ {
		if warehouses[i].Id == w.Id {
			w = warehouses[i]
			break
		}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	ok = createStockRow(p.Id, w.Id, 1, *trans)
	if !ok {
		t.Error("Can't create stock rows")
		return
	}

	///
	result := trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	ok = p.deleteProduct(0).Ok
	if !ok {
		t.Error("Delete error, could not delete product")
		return
	}

	ok = w.deleteWarehouse()
	if !ok {
		t.Error("Delete error, warehouse not deleted")
		return
	}

	trans.Rollback()
}

/* FUNCTIONALITY */

func TestStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create
	family := int32(1)
	manufacturingOrderType := int32(2)
	supplier := int32(1)

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
		SupplierId:               &supplier,
		TrackMinimumStock:        true,
		EnterpriseId:             1,
	}
	p.insertProduct(0)
	products := getProduct(1)
	p = products[len(products)-1]

	w := Warehouse{
		Id:           "WA",
		Name:         "Test warehouse",
		EnterpriseId: 1,
	}
	w.insertWarehouse()

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	// test stock functionality
	// quantity pending serving
	ok := addQuantityPendingServing(p.Id, w.Id, 1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending serving has not worked")
		return
	}

	///
	result := trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s := getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingServed != 1 {
		t.Error("Quantity pending serving not updated", s.QuantityPendingServed)
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	ok = addQuantityPendingServing(p.Id, w.Id, -1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending serving has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingServed != 0 {
		t.Error("Quantity pending serving not updated", s.QuantityPendingServed)
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	// quantity pending receiving
	ok = addQuantityPendingReveiving(p.Id, w.Id, 1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending receiving has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingReceived != 1 {
		t.Error("Quantity pending receiving not updated")
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	ok = addQuantityPendingReveiving(p.Id, w.Id, -1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending receiving has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingReceived != 0 {
		t.Error("Quantity pending receiving not updated")
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	// quantity pending manufacture
	ok = addQuantityPendingManufacture(p.Id, w.Id, 1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending manufacture has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingManufacture != 1 {
		t.Error("Quantity pending manufacture not updated")
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	ok = addQuantityPendingManufacture(p.Id, w.Id, -1, 1, *trans)
	if !ok {
		t.Error("Adding quantity pending manufacture has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingManufacture != 0 {
		t.Error("Quantity pending manufacture not updated")
		return

	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	// stock
	ok = addQuantityStock(p.Id, w.Id, 1, 1, *trans)
	if !ok {
		t.Error("Adding quantity has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 1 {
		t.Error("Quantity not updated")
		return
	}
	p = getProductRow(p.Id)
	if p.Stock != 1 {
		t.Error("Quantity not updated")
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	ok = addQuantityStock(p.Id, w.Id, -1, 1, *trans)
	if !ok {
		t.Error("Adding quantity has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 0 {
		t.Error("Quantity not updated")
		return
	}
	p = getProductRow(p.Id)
	if p.Stock != 0 {
		t.Error("Quantity not updated")
		return
	}

	///
	trans = dbOrm.Begin()
	if trans.Error != nil {
		t.Error("Can't begin transaction")
		return
	}
	///

	// set stock
	ok = setQuantityStock(p.Id, w.Id, 1, 1, *trans)
	if !ok {
		t.Error("Setting quantity has not worked")
		return
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		t.Error(result.Error)
		return
	}
	///

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 1 {
		t.Error("Setting quantity not updated")
		return
	}
	p = getProductRow(p.Id)
	if p.Stock != 1 {
		t.Error("Quantity not updated")
		return
	}

	// delete
	p.deleteProduct(0)
	w.deleteWarehouse()
}

// WAREHOUSE MOVEMENT

/* GET */

func TestGetWarehouseMovement(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: 1, enterprise: 1}
	m := q.getWarehouseMovement()

	if len(m.Movements) == 0 || m.Movements[0].Id <= 0 {
		t.Error("Scan error, warehouse movement with ID 0.")
		return
	}
}

func TestGetWarehouseMovementByWarehouse(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := WarehouseMovementByWarehouse{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1, enterprise: 1}, WarehouseId: "W1"}
	m := q.getWarehouseMovementByWarehouse()

	if len(m.Movements) == 0 || m.Movements[0].Id <= 0 {
		t.Error("Scan error, warehouse movement with ID 0.")
		return
	}
}

func TestGetWarehouseMovementBySalesDeliveryNote(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := getWarehouseMovementBySalesDeliveryNote(1, 1)

	if len(m) == 0 || m[0].Id <= 0 {
		t.Error("Scan error, warehouse movement with ID 0.")
		return
	}
}

func TestGetWarehouseMovementByPurchaseDeliveryNote(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := getWarehouseMovementByPurchaseDeliveryNote(1, 1)

	if len(m) == 0 || m[0].Id <= 0 {
		t.Error("Scan error, warehouse movement with ID 0.")
		return
	}
}

func TestSearchWarehouseMovement(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := WarehouseMovementSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1, enterprise: 1}, Search: ""}}
	m := s.searchWarehouseMovement()

	if len(m.Movements) == 0 || m.Movements[0].Id <= 0 {
		t.Error("Scan error, warehouse movement with ID 0.")
		return
	}
}

func TestGetWarehouseMovementRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getWarehouseMovementRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, warehouse movement row with ID 0.")
		return
	}

}

/* INSERT - DELETE */

func TestWarehouseMovementInsertDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create
	family := int32(1)
	manufacturingOrderType := int32(2)
	supplier := int32(1)

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
		SupplierId:               &supplier,
		TrackMinimumStock:        true,
		EnterpriseId:             1,
	}
	p.insertProduct(0)

	w := Warehouse{
		Id:           "WA",
		Name:         "Test warehouse",
		EnterpriseId: 1,
	}
	w.insertWarehouse()

	// test warehouse movements
	// create an input movement
	wm := WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok := wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s := getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 1 {
		t.Errorf("The stock has not been updated %d", s.Quantity)
		return
	}
	// delete the warehouse movement
	q := WarehouseMovementByWarehouse{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}, WarehouseId: w.Id}
	movements := q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[len(movements.Movements)-1]
	ok = wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	// create an input and then an output movement
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     -1,
		Type:         "O",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 0 {
		t.Error("The stock has not been updated")
		return
	}
	// delete the warehouse movement
	movements = q.getWarehouseMovementByWarehouse()
	for i := 0; i < len(movements.Movements); i++ {
		ok = movements.Movements[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Delete error, the warehouse movement could not be deleted")
			return
		}
	}

	// test sales delivery note generation
	sn := SalesDeliveryNote{
		WarehouseId:       "W1",
		CustomerId:        1,
		PaymentMethodId:   3,
		BillingSeriesId:   "EXP",
		ShippingAddressId: 1,
		CurrencyId:        1,
		EnterpriseId:      1,
	}

	_, noteId := sn.insertSalesDeliveryNotes(0, nil)
	sn.Id = noteId

	wm = WarehouseMovement{
		WarehouseId:         w.Id,
		ProductId:           p.Id,
		Quantity:            -1,
		Type:                "O",
		SalesDeliveryNoteId: &noteId,
		Price:               9.99,
		VatPercent:          21,
		EnterpriseId:        1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != -1 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[len(movements.Movements)-1]
	sn = getSalesDeliveryNoteRow(sn.Id)
	if sn.TotalAmount != wm.TotalAmount || sn.TotalProducts != float64(abs(wm.Quantity))*wm.Price {
		t.Error("The totals in the sale delivery note has not updated successfully")
		return
	}

	wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	sn = getSalesDeliveryNoteRow(sn.Id)
	if sn.TotalAmount != 0 || sn.TotalProducts != 0 {
		t.Error("The totals in the sale delivery note has not updated successfully")
		return
	}

	sn.Id = noteId
	sn.deleteSalesDeliveryNotes(0, nil)

	// test purchase delivery note generation
	pn := PurchaseDeliveryNote{
		WarehouseId:       "W1",
		SupplierId:        1,
		PaymentMethodId:   1,
		BillingSeriesId:   "INT",
		CurrencyId:        1,
		ShippingAddressId: 3,
		EnterpriseId:      1,
	}

	_, noteId = pn.insertPurchaseDeliveryNotes(0, nil)
	pn.Id = noteId

	wm = WarehouseMovement{
		WarehouseId:            w.Id,
		ProductId:              p.Id,
		Quantity:               -1,
		Type:                   "O",
		PurchaseDeliveryNoteId: &noteId,
		Price:                  9.99,
		VatPercent:             21,
		EnterpriseId:           1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != -1 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[len(movements.Movements)-1]
	pn = getPurchaseDeliveryNoteRow(pn.Id)
	if pn.TotalAmount != wm.TotalAmount || pn.TotalProducts != float64(abs(wm.Quantity))*wm.Price {
		t.Error("The totals in the purchase delivery note has not updated successfully")
		return
	}

	wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	pn = getPurchaseDeliveryNoteRow(pn.Id)
	if pn.TotalAmount != 0 || pn.TotalProducts != 0 {
		t.Error("The totals in the purchase delivery note has not updated successfully")
		return
	}

	pn.Id = noteId
	pn.deletePurchaseDeliveryNotes(0, nil)

	// test dragged stock
	// An input gets inserted, later an output, then a regularisation is made, and then an input is added. If this last row (input) gets deleted,
	// the stock of the product has to be equal to the stock set on the regularization again.
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     -1,
		Type:         "O",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     5,
		Type:         "R",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 6 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[0]
	ok = wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 5 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	for i := 0; i < len(movements.Movements); i++ {
		ok = movements.Movements[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Delete error, the warehouse movement could not be deleted")
			return
		}
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 0 {
		t.Error("The stock has not been updated")
		return
	}

	// test dragged stock
	// An input gets inserted, later an output, again an input is inserted, then a regularisation is made, and then an input is added. If the last
	// row before the regularization (input) is deleted, it should not affect the product's quantity, or the dragged stock of the products.
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     -1,
		Type:         "O",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     5,
		Type:         "R",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 6 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[2]
	ok = wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 6 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	for i := 0; i < len(movements.Movements); i++ {
		ok = movements.Movements[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Delete error, the warehouse movement could not be deleted")
			return
		}
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 0 {
		t.Error("The stock has not been updated")
		return
	}

	// test dragged stock
	// An input gets inserted, later an output, again an input is inserted, then a regularisation is made, and then an input is added. If the
	// regularization is deleted, the dragged stock of the product should be recalculated.
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     -1,
		Type:         "O",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     5,
		Type:         "R",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		WarehouseId:  w.Id,
		ProductId:    p.Id,
		Quantity:     1,
		Type:         "I",
		EnterpriseId: 1,
	}
	ok = wm.insertWarehouseMovement(0, nil)
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 6 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[1]
	ok = wm.deleteWarehouseMovement(0, nil)
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 2 {
		t.Error("The stock has not been updated")
		return
	}

	movements = q.getWarehouseMovementByWarehouse()
	for i := 0; i < len(movements.Movements); i++ {
		ok = movements.Movements[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Delete error, the warehouse movement could not be deleted")
			return
		}
	}

	s = getStockRow(p.Id, w.Id, 1)
	if s.Quantity != 0 {
		t.Error("The stock has not been updated")
		return
	}

	// delete
	p.deleteProduct(0)
	w.deleteWarehouse()
}

// ===== INVENTORY

/* GET */

func TestGetInventories(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	inventories := getInventories(1)
	if len(inventories) > 0 && inventories[0].Id <= 0 {
		t.Error("Can't scan inventories")
		return
	}
}

func TestHetInventoryRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	inventory := getInventoryRow(1)
	if inventory.Id <= 0 {
		t.Error("Can't scan inventory row")
		return
	}
}

func TestGetInventoryProducts(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	inventoryProducts := getInventoryProducts(1, 1)
	if len(inventoryProducts) > 0 && inventoryProducts[0].InventoryId <= 0 {
		t.Error("Can't scan invenrory products")
		return
	}
}

func TestGetInventoryProductsRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	inventoryProducts := getInventoryProducts(1, 1)
	if len(inventoryProducts) > 0 {
		row := getInventoryProductsRow(1, inventoryProducts[0].ProductId, 1)
		if row.ProductId <= 0 {
			t.Error("Can't scan invenrory product row")
			return
		}
	}
}

/* INSERT - UPDATE - DELETE */

func TestInsertDeleteInventoryAndDetails(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert inventory

	i := Inventory{
		Name:         "Automatic test inventory",
		WarehouseId:  "W1",
		EnterpriseId: 1,
	}
	ok := i.insertInventory(1)
	if !ok {
		t.Error("Can't insert inventory")
		return
	}

	inventories := getInventories(1)
	i = inventories[0]

	// add single product
	input := InputInventoryProducts{
		Inventory: i.Id,
		InventoryProducts: []InventoryProducts{
			{
				ProductId:    1,
				InventoryId:  i.Id,
				Quantity:     1,
				EnterpriseId: 1,
			},
		},
	}
	ok = input.insertUpdateDeleteInventoryProducts(1)
	if !ok {
		t.Error("Can't save products", i.Id)
		return
	}
	inventoryProducts := getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) == 0 {
		t.Error("Can't add products")
		return
	}

	// delete a single product
	input = InputInventoryProducts{
		Inventory:         i.Id,
		InventoryProducts: []InventoryProducts{},
	}
	ok = input.insertUpdateDeleteInventoryProducts(1)
	if !ok {
		t.Error("Can't save products")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) != 0 {
		t.Error("Can't delete single product")
		return
	}

	// add a family
	input = InputInventoryProducts{
		Inventory: i.Id,
		FamilyId:  1,
	}
	ok = input.insertProductFamilyInventoryProducts(1)
	if !ok {
		t.Error("Can't insert product family")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) == 0 {
		t.Error("Can't add product family")
		return
	}

	// delete all
	input = InputInventoryProducts{
		Inventory: i.Id,
	}
	ok = input.deleteAllProductsInventoryProducts(1)
	if !ok {
		t.Error("Can't delete all")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) != 0 {
		t.Error("Can't delete all")
		return
	}

	// add all
	input = InputInventoryProducts{
		Inventory: i.Id,
	}
	ok = input.insertAllProductsInventoryProducts(1)
	if !ok {
		t.Error("Can't add all")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) == 0 {
		t.Error("Can't add all")
		return
	}

	// delete all
	input = InputInventoryProducts{
		Inventory: i.Id,
	}
	ok = input.deleteAllProductsInventoryProducts(1)
	if !ok {
		t.Error("Can't delete all")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) != 0 {
		t.Error("Can't delete all")
		return
	}

	// delete inventory
	okAndErr := i.deleteInventory(1)
	if !okAndErr.Ok {
		t.Error("Can't delete inventory", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}

/* FUNCTIONALITY */

func TestInventoryScanBarCode(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert inventory

	i := Inventory{
		Name:         "Automatic test inventory",
		WarehouseId:  "W1",
		EnterpriseId: 1,
	}
	ok := i.insertInventory(1)
	if !ok {
		t.Error("Can't insert inventory")
		return
	}

	inventories := getInventories(1)
	i = inventories[0]

	// add single product
	input := InputInventoryProducts{
		Inventory: i.Id,
		InventoryProducts: []InventoryProducts{
			{
				ProductId:    1,
				InventoryId:  i.Id,
				Quantity:     0,
				EnterpriseId: 1,
			},
		},
	}
	ok = input.insertUpdateDeleteInventoryProducts(1)
	if !ok {
		t.Error("Can't save products", i.Id)
		return
	}
	inventoryProducts := getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) == 0 {
		t.Error("Can't add products")
		return
	}

	// test scan barcode
	product := getProductRow(1)
	inputBarCode := BarCodeInputInventoryProducts{
		Inventory: i.Id,
		BarCode:   product.BarCode,
	}
	res := inputBarCode.insertOrCountInventoryProductsByBarcode(1)
	if !res.Ok {
		t.Error("Error scanning barcode")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) > 0 && inventoryProducts[0].Quantity != 1 {
		t.Error("Can't add quantity to products using barcode")
		return
	}

	// delete all
	input = InputInventoryProducts{
		Inventory: i.Id,
	}
	ok = input.deleteAllProductsInventoryProducts(1)
	if !ok {
		t.Error("Can't delete all")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) != 0 {
		t.Error("Can't delete all")
		return
	}

	// delete inventory
	okAndErr := i.deleteInventory(1)
	if !okAndErr.Ok {
		t.Error("Can't delete inventory", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}

func TestInventoryFinish(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert inventory

	i := Inventory{
		Name:         "Automatic test inventory",
		WarehouseId:  "W1",
		EnterpriseId: 1,
	}
	ok := i.insertInventory(1)
	if !ok {
		t.Error("Can't insert inventory")
		return
	}

	inventories := getInventories(1)
	i = inventories[0]

	// create a new product

	family := int32(1)
	manufacturingOrderType := int32(1)

	p := Product{
		Name:                     "Glass Office Desk",
		Reference:                "OF-DSK",
		BarCode:                  "",
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

	// add single product
	input := InputInventoryProducts{
		Inventory: i.Id,
		InventoryProducts: []InventoryProducts{
			{
				ProductId:    p.Id,
				InventoryId:  i.Id,
				Quantity:     15,
				EnterpriseId: 1,
			},
		},
	}
	ok = input.insertUpdateDeleteInventoryProducts(1)
	if !ok {
		t.Error("Can't save products", i.Id)
		return
	}
	inventoryProducts := getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) == 0 {
		t.Error("Can't add products")
		return
	}

	// finish inventory
	ok = i.finishInventory(0, 1)
	if !ok {
		t.Error("Error finishing inventory")
		return
	}
	stock := getStockRow(p.Id, i.WarehouseId, 1)
	if stock.Quantity != 15 {
		t.Error("Stock not updated!!!")
		return
	}

	// delete all
	// FORCE DELETE
	sqlStatement := `UPDATE inventory SET finished = false WHERE id = $1`
	db.Exec(sqlStatement, i.Id)

	input = InputInventoryProducts{
		Inventory: i.Id,
	}
	ok = input.deleteAllProductsInventoryProducts(1)
	if !ok {
		t.Error("Can't delete all")
		return
	}
	inventoryProducts = getInventoryProducts(i.Id, 1)
	if len(inventoryProducts) != 0 {
		t.Error("Can't delete all")
		return
	}

	// delete inventory
	okAndErr = i.deleteInventory(1)
	if !okAndErr.Ok {
		t.Error("Can't delete inventory", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}

	// delete warehouse movement
	wm := getProductWarehouseMovement(ProductPurchaseOrderDetailsQuery{
		ProductId: p.Id,
	}, 1)
	if len(wm) > 0 {
		ok = wm[0].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Can't delete warehouse movement")
			return
		}
	}

	// delete product
	okAndErr = p.deleteProduct(0)
	if !okAndErr.Ok {
		t.Error("Delete error, could not delete product", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}

// ===== TRANSFER BETWEEN WAREHOUSES

/* GET */

func TestSearchTransferBetweenWarehouses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := TransferBetweenWarehousesQuery{
		enterprise: 1,
	}
	transfers := q.searchTransferBetweenWarehouses()
	if len(transfers) == 0 || transfers[0].Id <= 0 {
		t.Error("Can't scan transfers")
		return
	}
}

func TestGetTransferBetweenWarehousesRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	transfer := getTransferBetweenWarehousesRow(1)
	if transfer.Id <= 0 {
		t.Error("Can't scan transfer row")
		return
	}
}

func TestGetTransferBetweenWarehousesDetails(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	details := getTransferBetweenWarehousesDetails(1, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't scan transfer details")
		return
	}
}

func TestGetTransferBetweenWarehousesDetailRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	detail := getTransferBetweenWarehousesDetailRow(1)
	if detail.Id <= 0 {
		t.Error("Can't scan transfer row")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestInsertDeleteTransferBetweenWarehouses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	transfer := TransferBetweenWarehouses{
		WarehouseOriginId:      "W1",
		WarehouseDestinationId: "WT",
		Name:                   "Automatic test",
		EnterpriseId:           1,
	}
	ok := transfer.insertTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't insert transfer between warehouses")
		return
	}

	q := TransferBetweenWarehousesQuery{
		enterprise: 1,
	}
	transfers := q.searchTransferBetweenWarehouses()
	transfer = transfers[0]

	d := TransferBetweenWarehousesDetail{
		TransferBetweenWarehousesId: transfer.Id,
		ProductId:                   1,
		Quantity:                    10,
		EnterpriseId:                1,
	}
	ok = d.insertTransferBetweenWarehousesDetail()
	if !ok {
		t.Error("Can't insert transfer between warehouses detail")
		return
	}

	details := getTransferBetweenWarehousesDetails(transfer.Id, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't scan transfer details")
		return
	}
	d = details[0]

	ok = d.deleteTransferBetweenWarehousesDetail(nil)
	if !ok {
		t.Error("Can't delete transfer between warehouses detail")
		return
	}

	// CLEAN UP

	ok = transfer.deleteTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't delete transfer between warehouses")
		return
	}
}

/* FUNCTIONALITY */

func TestTransferBetweenWarehousesUsingQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(1)

	p := Product{
		Name:                     "Glass Office Desk",
		Reference:                "OF-DSK",
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

	transfer := TransferBetweenWarehouses{
		WarehouseOriginId:      "W1",
		WarehouseDestinationId: "WT",
		Name:                   "Automatic test",
		EnterpriseId:           1,
	}
	ok := transfer.insertTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't insert transfer between warehouses")
		return
	}

	q := TransferBetweenWarehousesQuery{
		enterprise: 1,
	}
	transfers := q.searchTransferBetweenWarehouses()
	transfer = transfers[0]

	d := TransferBetweenWarehousesDetail{
		TransferBetweenWarehousesId: transfer.Id,
		ProductId:                   p.Id,
		Quantity:                    10,
		EnterpriseId:                1,
	}
	ok = d.insertTransferBetweenWarehousesDetail()
	if !ok {
		t.Error("Can't insert transfer between warehouses detail")
		return
	}

	details := getTransferBetweenWarehousesDetails(transfer.Id, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't scan transfer details")
		return
	}
	d = details[0]

	// add quantity
	query := TransferBetweenWarehousesDetailQuantityQuery{
		TransferBetweenWarehousesDetailId: d.Id,
		Quantity:                          5,
	}
	ok = query.transferBetweenWarehousesDetailQuantity(1, 0)
	if !ok {
		t.Error("Can't add quantity transfer between warehouses detail")
		return
	}

	details = getTransferBetweenWarehousesDetails(transfer.Id, 1)
	d = details[0]
	if d.Finished || d.QuantityTransferred != 5 {
		t.Error("Error tranfering quantity", d.Finished, d.QuantityTransferred)
		return
	}

	// add more quantity
	query = TransferBetweenWarehousesDetailQuantityQuery{
		TransferBetweenWarehousesDetailId: d.Id,
		Quantity:                          5,
	}
	ok = query.transferBetweenWarehousesDetailQuantity(1, 0)
	if !ok {
		t.Error("Can't add quantity transfer between warehouses detail")
		return
	}

	details = getTransferBetweenWarehousesDetails(transfer.Id, 1)
	d = details[0]
	if (!d.Finished) || d.QuantityTransferred != 10 {
		t.Error("Error tranfering quantity", d.Finished, d.QuantityTransferred)
		return
	}

	// CLEAN UP
	// force undo
	sqlStatement := `UPDATE public.transfer_between_warehouses_detail SET quantity_transferred=0, finished=false, warehouse_movement_out=NULL, warehouse_movement_in=NULL WHERE id=$1`
	db.Exec(sqlStatement, d.Id)

	ok = d.deleteTransferBetweenWarehousesDetail(nil)
	if !ok {
		t.Error("Can't delete transfer between warehouses detail")
		return
	}

	ok = transfer.deleteTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't delete transfer between warehouses")
		return
	}

	// delete warehouse movement
	wm := getProductWarehouseMovement(ProductPurchaseOrderDetailsQuery{
		ProductId: p.Id,
	}, 1)
	for i := 0; i < len(wm); i++ {
		ok = wm[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Can't delete warehouse movement")
			return
		}
	}

	// delete product
	okAndErr = p.deleteProduct(0)
	if !okAndErr.Ok {
		t.Error("Delete error, could not delete product", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}

func TestTransferBetweenWarehousesUsingBarCode(t *testing.T) {
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

	transfer := TransferBetweenWarehouses{
		WarehouseOriginId:      "W1",
		WarehouseDestinationId: "WT",
		Name:                   "Automatic test",
		EnterpriseId:           1,
	}
	ok := transfer.insertTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't insert transfer between warehouses")
		return
	}

	q := TransferBetweenWarehousesQuery{
		enterprise: 1,
	}
	transfers := q.searchTransferBetweenWarehouses()
	transfer = transfers[0]

	d := TransferBetweenWarehousesDetail{
		TransferBetweenWarehousesId: transfer.Id,
		ProductId:                   p.Id,
		Quantity:                    2,
		EnterpriseId:                1,
	}
	ok = d.insertTransferBetweenWarehousesDetail()
	if !ok {
		t.Error("Can't insert transfer between warehouses detail")
		return
	}

	details := getTransferBetweenWarehousesDetails(transfer.Id, 1)
	if len(details) == 0 || details[0].Id <= 0 {
		t.Error("Can't scan transfer details")
		return
	}
	d = details[0]

	// add quantity
	query := TransferBetweenWarehousesDetailBarCodeQuery{
		TransferBetweenWarehousesId: transfer.Id,
		BarCode:                     p.BarCode,
	}
	ok = query.transferBetweenWarehousesDetailBarCode(1, 0)
	if !ok {
		t.Error("Can't add quantity transfer between warehouses detail")
		return
	}

	details = getTransferBetweenWarehousesDetails(transfer.Id, 1)
	d = details[0]
	if d.Finished || d.QuantityTransferred != 1 {
		t.Error("Error tranfering quantity", d.Finished, d.QuantityTransferred)
	}

	// add more quantity
	query = TransferBetweenWarehousesDetailBarCodeQuery{
		TransferBetweenWarehousesId: transfer.Id,
		BarCode:                     p.BarCode,
	}
	ok = query.transferBetweenWarehousesDetailBarCode(1, 0)
	if !ok {
		t.Error("Can't add quantity transfer between warehouses detail")
		return
	}

	details = getTransferBetweenWarehousesDetails(transfer.Id, 1)
	d = details[0]
	if (!d.Finished) || d.QuantityTransferred != 2 {
		t.Error("Error tranfering quantity", d.Finished, d.QuantityTransferred)
	}

	// CLEAN UP
	// force undo
	sqlStatement := `UPDATE public.transfer_between_warehouses_detail SET quantity_transferred=0, finished=false, warehouse_movement_out=NULL, warehouse_movement_in=NULL WHERE id=$1`
	db.Exec(sqlStatement, d.Id)

	ok = d.deleteTransferBetweenWarehousesDetail(nil)
	if !ok {
		t.Error("Can't delete transfer between warehouses detail")
		return
	}

	ok = transfer.deleteTransferBetweenWarehouses()
	if !ok {
		t.Error("Can't delete transfer between warehouses")
		return
	}

	// delete warehouse movement
	wm := getProductWarehouseMovement(ProductPurchaseOrderDetailsQuery{
		ProductId: p.Id,
	}, 1)
	for i := 0; i < len(wm); i++ {
		ok = wm[i].deleteWarehouseMovement(0, nil)
		if !ok {
			t.Error("Can't delete warehouse movement")
			return
		}
	}

	// delete product
	okAndErr = p.deleteProduct(0)
	if !okAndErr.Ok {
		t.Error("Delete error, could not delete product", okAndErr.ErrorCode, okAndErr.ExtraData)
		return
	}
}
