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
		Id:         "WA",
		Name:       "Test warehouse",
		enterprise: 1,
	}

	// insert
	ok := w.insertWarehouse()
	if !ok {
		t.Error("Insert error, warehouse not inserted")
		return
	}

	// update
	warehouses := getWarehouses(1)
	w = warehouses[len(warehouses)-1]

	w.Name = "Test test"
	ok = w.updateWarehouse()
	if !ok {
		t.Error("Update error, warehouse not updated")
		return
	}

	// check update
	warehouses = getWarehouses(1)
	w = warehouses[len(warehouses)-1]

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
		if s[i].Product <= 0 {
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
	if s.Product <= 0 {
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
		Name:                   "Glass Office Desk",
		Reference:              "OF-DSK",
		BarCode:                "1234067891234",
		ControlStock:           true,
		Weight:                 30,
		Family:                 &family,
		Width:                  160,
		Height:                 100,
		Depth:                  40,
		VatPercent:             21,
		Price:                  65,
		Manufacturing:          true,
		ManufacturingOrderType: &manufacturingOrderType,
		Supplier:               &supplier,
		TrackMinimumStock:      true,
		enterprise:             1,
	}

	ok := p.insertProduct()
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	w := Warehouse{
		Id:         "WA",
		Name:       "Test warehouse",
		enterprise: 1,
	}

	// insert
	ok = w.insertWarehouse()
	if !ok {
		t.Error("Insert error, warehouse not inserted")
		return
	}
	warehouses := getWarehouses(1)
	w = warehouses[len(warehouses)-1]

	ok = createStockRow(p.Id, w.Id, 1)
	if !ok {
		t.Error("Can't create stock rows")
	}

	ok = p.deleteProduct()
	if !ok {
		t.Error("Delete error, could not delete product")
		return
	}

	ok = w.deleteWarehouse()
	if !ok {
		t.Error("Delete error, warehouse not deleted")
		return
	}
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
		Name:                   "Glass Office Desk",
		Reference:              "OF-DSK",
		BarCode:                "1234067891234",
		ControlStock:           true,
		Weight:                 30,
		Family:                 &family,
		Width:                  160,
		Height:                 100,
		Depth:                  40,
		VatPercent:             21,
		Price:                  65,
		Manufacturing:          true,
		ManufacturingOrderType: &manufacturingOrderType,
		Supplier:               &supplier,
		TrackMinimumStock:      true,
		enterprise:             1,
	}
	p.insertProduct()
	products := getProduct(1)
	p = products[len(products)-1]

	w := Warehouse{
		Id:         "WA",
		Name:       "Test warehouse",
		enterprise: 1,
	}
	w.insertWarehouse()

	// test stock functionality
	// quantity pending serving
	ok := addQuantityPendingServing(p.Id, w.Id, 1, 1)
	if !ok {
		t.Error("Adding quantity pending serving has not worked")
		return
	}
	s := getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingServed != 1 {
		t.Error("Quantity pending serving not updated")
		return
	}
	ok = addQuantityPendingServing(p.Id, w.Id, -1, 1)
	if !ok {
		t.Error("Adding quantity pending serving has not worked")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingServed != 0 {
		t.Error("Quantity pending serving not updated")
		return
	}

	// quantity pending receiving
	ok = addQuantityPendingReveiving(p.Id, w.Id, 1, 1)
	if !ok {
		t.Error("Adding quantity pending receiving has not worked")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingReceived != 1 {
		t.Error("Quantity pending receiving not updated")
		return
	}
	ok = addQuantityPendingReveiving(p.Id, w.Id, -1, 1)
	if !ok {
		t.Error("Adding quantity pending receiving has not worked")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingReceived != 0 {
		t.Error("Quantity pending receiving not updated")
		return
	}

	// quantity pending manufacture
	ok = addQuantityPendingManufacture(p.Id, w.Id, 1, 1)
	if !ok {
		t.Error("Adding quantity pending manufacture has not worked")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingManufacture != 1 {
		t.Error("Quantity pending manufacture not updated")
		return
	}
	ok = addQuantityPendingManufacture(p.Id, w.Id, -1, 1)
	if !ok {
		t.Error("Adding quantity pending manufacture has not worked")
		return
	}
	s = getStockRow(p.Id, w.Id, 1)
	if s.QuantityPendingManufacture != 0 {
		t.Error("Quantity pending manufacture not updated")
		return
	}

	// stock
	ok = addQuantityStock(p.Id, w.Id, 1, 1)
	if !ok {
		t.Error("Adding quantity has not worked")
		return
	}
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
	ok = addQuantityStock(p.Id, w.Id, -1, 1)
	if !ok {
		t.Error("Adding quantity has not worked")
		return
	}
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

	// set stock
	ok = setQuantityStock(p.Id, w.Id, 1, 1)
	if !ok {
		t.Error("Setting quantity has not worked")
		return
	}
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
	p.deleteProduct()
	w.deleteWarehouse()
}

// WAREHOUSE MOVEMENT

/* GET */

func TestGetWarehouseMovement(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: 1, Enterprise: 1}
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

	q := WarehouseMovementByWarehouse{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1, Enterprise: 1}, WarehouseId: "W1"}
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

	s := WarehouseMovementSearch{PaginatedSearch: PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1, Enterprise: 1}, Search: ""}}
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
		Name:                   "Glass Office Desk",
		Reference:              "OF-DSK",
		BarCode:                "1234067891234",
		ControlStock:           true,
		Weight:                 30,
		Family:                 &family,
		Width:                  160,
		Height:                 100,
		Depth:                  40,
		VatPercent:             21,
		Price:                  65,
		Manufacturing:          true,
		ManufacturingOrderType: &manufacturingOrderType,
		Supplier:               &supplier,
		TrackMinimumStock:      true,
		enterprise:             1,
	}
	p.insertProduct()

	w := Warehouse{
		Id:         "WA",
		Name:       "Test warehouse",
		enterprise: 1,
	}
	w.insertWarehouse()
	warehouses := getWarehouses(1)
	w = warehouses[len(warehouses)-1]

	// test warehouse movements
	// create an input movement
	wm := WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok := wm.insertWarehouseMovement()
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
	q := WarehouseMovementByWarehouse{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, Enterprise: 1}, WarehouseId: w.Id}
	movements := q.getWarehouseMovementByWarehouse()
	wm = movements.Movements[len(movements.Movements)-1]
	ok = wm.deleteWarehouseMovement()
	if !ok {
		t.Error("Delete error, the warehouse movement could not be deleted")
		return
	}

	// create an input and then an output movement
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   -1,
		Type:       "O",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
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
		ok = movements.Movements[i].deleteWarehouseMovement()
		if !ok {
			t.Error("Delete error, the warehouse movement could not be deleted")
			return
		}
	}

	// test sales delivery note generation
	sn := SalesDeliveryNote{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		ShippingAddress: 1,
		Currency:        1,
		enterprise:      1,
	}

	_, noteId := sn.insertSalesDeliveryNotes()
	sn.Id = noteId

	wm = WarehouseMovement{
		Warehouse:         w.Id,
		Product:           p.Id,
		Quantity:          -1,
		Type:              "O",
		SalesDeliveryNote: &noteId,
		Price:             9.99,
		VatPercent:        21,
		enterprise:        1,
	}
	ok = wm.insertWarehouseMovement()
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
	if sn.TotalAmount != wm.TotalAmount || sn.TotalProducts != float32(abs(wm.Quantity))*wm.Price {
		t.Error("The totals in the sale delivery note has not updated successfully")
		return
	}

	wm.deleteWarehouseMovement()
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
	sn.deleteSalesDeliveryNotes()

	// test purchase delivery note generation
	pn := PurchaseDeliveryNote{
		Warehouse:       "W1",
		Supplier:        1,
		PaymentMethod:   1,
		BillingSeries:   "INT",
		Currency:        1,
		ShippingAddress: 3,
		enterprise:      1,
	}

	_, noteId = pn.insertPurchaseDeliveryNotes()
	pn.Id = noteId

	wm = WarehouseMovement{
		Warehouse:            w.Id,
		Product:              p.Id,
		Quantity:             -1,
		Type:                 "O",
		PurchaseDeliveryNote: &noteId,
		Price:                9.99,
		VatPercent:           21,
		enterprise:           1,
	}
	ok = wm.insertWarehouseMovement()
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
	if pn.TotalAmount != wm.TotalAmount || pn.TotalProducts != float32(abs(wm.Quantity))*wm.Price {
		t.Error("The totals in the purchase delivery note has not updated successfully")
		return
	}

	wm.deleteWarehouseMovement()
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
	pn.deletePurchaseDeliveryNotes()

	// test dragged stock
	// An input gets inserted, later an output, then a regularisation is made, and then an input is added. If this last row (input) gets deleted,
	// the stock of the product has to be equal to the stock set on the regularization again.
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   -1,
		Type:       "O",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   5,
		Type:       "R",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
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
	ok = wm.deleteWarehouseMovement()
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
		ok = movements.Movements[i].deleteWarehouseMovement()
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
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   -1,
		Type:       "O",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   5,
		Type:       "R",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
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
	ok = wm.deleteWarehouseMovement()
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
		ok = movements.Movements[i].deleteWarehouseMovement()
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
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   -1,
		Type:       "O",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   5,
		Type:       "R",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
	if !ok {
		t.Error("Insert error, the warehouse movement could not be inserted")
		return
	}
	wm = WarehouseMovement{
		Warehouse:  w.Id,
		Product:    p.Id,
		Quantity:   1,
		Type:       "I",
		enterprise: 1,
	}
	ok = wm.insertWarehouseMovement()
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
	ok = wm.deleteWarehouseMovement()
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
		ok = movements.Movements[i].deleteWarehouseMovement()
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
	p.deleteProduct()
	w.deleteWarehouse()
}
