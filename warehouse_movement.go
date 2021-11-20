package main

import (
	"strconv"
	"time"
)

type WarehouseMovement struct {
	Id                    int64     `json:"id"`
	Warehouse             string    `json:"warehouse"`
	Product               int32     `json:"product"`
	Quantity              int32     `json:"quantity"`
	DateCreated           time.Time `json:"dateCreated"`
	Type                  string    `json:"type"` // O = Out, I = In, R = Inventory regularization
	SalesOrder            *int64    `json:"salesOrder"`
	SalesOrderDetail      *int64    `json:"salesOrderDetail"`
	SalesInvoice          *int64    `json:"salesInvoice"`
	SalesInvoiceDetail    *int64    `json:"salesInvoiceDetail"`
	SalesDeliveryNote     *int64    `json:"salesDeliveryNote"`
	Description           string    `json:"description"`
	PurchaseOrder         *int64    `json:"purchaseOrder"`
	PurchaseOrderDetail   *int64    `json:"purchaseOrderDetail"`
	PurchaseInvoice       *int64    `json:"purchaseInvoice"`
	PurchaseInvoiceDetail *int64    `json:"purchaseInvoiceDetail"`
	PurchaseDeliveryNote  *int64    `json:"purchaseDeliveryNote"`
	DraggedStock          int32     `json:"draggedStock"`
	ProductName           string    `json:"productName"`
	Price                 float64   `json:"price"`
	VatPercent            float64   `json:"vatPercent"`
	TotalAmount           float64   `json:"totalAmount"`
	WarehouseName         string    `json:"warehouseName"`
	enterprise            int32
}

type WarehouseMovements struct {
	Rows      int32               `json:"rows"`
	Movements []WarehouseMovement `json:"movements"`
}

func (q *PaginationQuery) getWarehouseMovement() WarehouseMovements {
	wm := WarehouseMovements{}
	if !q.isValid() {
		return wm
	}

	wm.Movements = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM public.warehouse_movement WHERE enterprise=$3 ORDER BY id DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, q.enterprise)
	if err != nil {
		log("DB", err.Error())
		return wm
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		wm.Movements = append(wm.Movements, m)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.warehouse_movement WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, q.enterprise)
	row.Scan(&wm.Rows)

	return wm
}

type WarehouseMovementByWarehouse struct {
	PaginationQuery
	WarehouseId string `json:"warehouseId"`
}

func (w *WarehouseMovementByWarehouse) getWarehouseMovementByWarehouse() WarehouseMovements {
	wm := WarehouseMovements{}
	wm.Movements = make([]WarehouseMovement, 0)
	if len(w.WarehouseId) == 0 || len(w.WarehouseId) > 2 {
		return wm
	}

	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM public.warehouse_movement WHERE warehouse=$1 AND enterprise=$4 ORDER BY id DESC OFFSET $2 LIMIT $3`
	rows, err := db.Query(sqlStatement, w.WarehouseId, w.Offset, w.Limit, w.enterprise)
	if err != nil {
		log("DB", err.Error())
		return wm
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		wm.Movements = append(wm.Movements, m)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.warehouse_movement WHERE warehouse=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, w.WarehouseId, w.enterprise)
	if row.Err() != nil {
		return wm
	}
	row.Scan(&wm.Rows)

	return wm
}

func getWarehouseMovementRow(movementId int64) WarehouseMovement {
	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE id=$1`
	row := db.QueryRow(sqlStatement, movementId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return WarehouseMovement{}
	}

	m := WarehouseMovement{}
	row.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise)

	return m
}

func getWarehouseMovementBySalesDeliveryNote(noteId int64, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM public.warehouse_movement WHERE sales_delivery_note=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func getWarehouseMovementByPurchaseDeliveryNote(noteId int64, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM public.warehouse_movement WHERE purchase_delivery_note=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

type WarehouseMovementSearch struct {
	PaginatedSearch
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (w *WarehouseMovementSearch) searchWarehouseMovement() WarehouseMovements {
	wm := WarehouseMovements{}
	if !w.isValid() {
		return wm
	}

	wm.Movements = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT warehouse_movement.*,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM warehouse_movement INNER JOIN product ON product.id=warehouse_movement.product WHERE product.name ILIKE $1`
	parameters := make([]interface{}, 0)
	parameters = append(parameters, "%"+w.Search+"%")
	if w.DateStart != nil {
		sqlStatement += ` AND warehouse_movement.date_created >= $2`
		parameters = append(parameters, w.DateStart)
	}
	if w.DateEnd != nil {
		sqlStatement += ` AND warehouse_movement.date_created <= $` + strconv.Itoa(len(parameters)+1)
		parameters = append(parameters, w.DateEnd)
	}
	sqlStatement += ` AND warehouse_movement.enterprise = $` + strconv.Itoa(len(parameters)+1)
	parameters = append(parameters, w.enterprise)
	sqlStatement += ` ORDER BY warehouse_movement.id DESC OFFSET $` + strconv.Itoa(len(parameters)+1) + ` LIMIT $` + strconv.Itoa(len(parameters)+2)
	parameters = append(parameters, w.Offset)
	parameters = append(parameters, w.Limit)
	rows, err := db.Query(sqlStatement, parameters...)
	if err != nil {
		log("DB", err.Error())
		return wm
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		wm.Movements = append(wm.Movements, m)
	}

	sqlStatement = `SELECT COUNT(warehouse_movement.*) FROM warehouse_movement INNER JOIN product ON product.id=warehouse_movement.product WHERE product.name ILIKE $1`
	parameters = make([]interface{}, 0)
	parameters = append(parameters, "%"+w.Search+"%")
	if w.DateStart != nil {
		sqlStatement += ` AND warehouse_movement.date_created >= $2`
		parameters = append(parameters, w.DateStart)
	}
	if w.DateEnd != nil {
		sqlStatement += ` AND warehouse_movement.date_created <= $` + strconv.Itoa(len(parameters)+1)
		parameters = append(parameters, w.DateEnd)
	}
	sqlStatement += ` AND warehouse_movement.enterprise = $` + strconv.Itoa(len(parameters)+1)
	parameters = append(parameters, w.enterprise)
	row := db.QueryRow(sqlStatement, parameters...)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return wm
	}
	row.Scan(&wm.Rows)

	return wm
}

func (m *WarehouseMovement) isValid() bool {
	return !(len(m.Warehouse) == 0 || len(m.Warehouse) > 2 || m.Product <= 0 || m.Quantity == 0 || len(m.Type) != 1 || (m.Type != "I" && m.Type != "O" && m.Type != "R"))
}

func (m *WarehouseMovement) insertWarehouseMovement(userId int32) bool {
	if !m.isValid() {
		return false
	}

	m.TotalAmount = absf((m.Price * float64(m.Quantity)) * (1 + (m.VatPercent / 100)))

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	// get the dragged stock
	if m.Type != "R" {
		var dragged_stock int32
		sqlStatement := `SELECT dragged_stock FROM warehouse_movement WHERE warehouse=$1 AND product=$2 ORDER BY date_created DESC LIMIT 1`
		row := db.QueryRow(sqlStatement, m.Warehouse, m.Product)
		row.Scan(&dragged_stock)
		m.DraggedStock = dragged_stock + m.Quantity
	} else { // Inventory regularization
		m.DraggedStock = m.Quantity
	}

	// insert the movement
	sqlStatement := `INSERT INTO public.warehouse_movement(warehouse, product, quantity, type, sales_order, sales_order_detail, sales_invoice, sales_invoice_detail, sales_delivery_note, dsc, purchase_order, purchase_order_detail, purchase_invoice, purchase_invoice_details, purchase_delivery_note, dragged_stock, price, vat_percent, total_amount, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20) RETURNING id`
	row := db.QueryRow(sqlStatement, m.Warehouse, m.Product, m.Quantity, m.Type, m.SalesOrder, m.SalesOrderDetail, m.SalesInvoice, m.SalesInvoiceDetail, m.SalesDeliveryNote, m.Description, m.PurchaseOrder, m.PurchaseOrderDetail, m.PurchaseInvoice, m.PurchaseInvoiceDetail, m.PurchaseDeliveryNote, m.DraggedStock, m.Price, m.VatPercent, m.TotalAmount, m.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var warehouseMovementId int64
	row.Scan(&warehouseMovementId)

	if warehouseMovementId <= 0 {
		return false
	}

	m.Id = warehouseMovementId

	insertTransactionalLog(m.enterprise, "warehouse_movement", int(warehouseMovementId), userId, "I")

	// update the product quantity
	ok := setQuantityStock(m.Product, m.Warehouse, m.DraggedStock, m.enterprise)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery notes generation
	if m.SalesOrderDetail != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*m.SalesOrderDetail, abs(m.Quantity), userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if m.PurchaseOrderDetail != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*m.PurchaseOrderDetail, abs(m.Quantity), m.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// sales delivery note price
	if m.SalesDeliveryNote != nil {
		ok = addTotalProductsSalesDeliveryNote(*m.SalesDeliveryNote, absf(m.Price*float64(m.Quantity)), m.VatPercent, m.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// purchase delivery note price
	if m.PurchaseDeliveryNote != nil {
		ok = addTotalProductsPurchaseDeliveryNote(*m.PurchaseDeliveryNote, absf(m.Price*float64(m.Quantity)), m.VatPercent, m.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	err := trans.Commit()
	return err == nil
	///
}

// Abs returns the absolute value of x.
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// Abs returns the absolute value of x.
func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (m *WarehouseMovement) deleteWarehouseMovement(userId int32) bool {
	if m.Id <= 0 {
		return false
	}

	inMemoryMovement := getWarehouseMovementRow(m.Id)
	if inMemoryMovement.Id <= 0 || inMemoryMovement.enterprise != m.enterprise {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	insertTransactionalLog(m.enterprise, "warehouse_movement", int(m.Id), userId, "D")

	// delete the warehouse movement
	sqlStatement := `DELETE FROM public.warehouse_movement WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, m.Id, m.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rowsCount, _ := res.RowsAffected()
	if rowsCount == 0 {
		return false
	}

	// update the dragged stock
	var draggedStock int32
	if inMemoryMovement.Type != "R" {
		draggedStock = inMemoryMovement.DraggedStock - inMemoryMovement.Quantity
	} else {
		sqlStatement := `SELECT dragged_stock FROM warehouse_movement WHERE warehouse=$1 AND product=$2 AND date_created<=$3 ORDER BY date_created DESC LIMIT 1`
		row := db.QueryRow(sqlStatement, inMemoryMovement.Warehouse, inMemoryMovement.Product, inMemoryMovement.DateCreated)
		row.Scan(&draggedStock)
	}

	sqlStatement = `SELECT id,quantity,type FROM warehouse_movement WHERE warehouse=$1 AND product=$2 AND date_created>=$3 ORDER BY date_created ASC, id ASC`
	rows, err := db.Query(sqlStatement, inMemoryMovement.Warehouse, inMemoryMovement.Product, inMemoryMovement.DateCreated)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	for rows.Next() {
		var movementId int64
		var quantity int32
		var movementType string
		rows.Scan(&movementId, &quantity, &movementType)

		if movementType == "R" {
			draggedStock = quantity
		} else {
			draggedStock += quantity
		}

		sqlStatement := `UPDATE warehouse_movement SET dragged_stock=$2 WHERE id=$1`
		_, err := db.Exec(sqlStatement, movementId, draggedStock)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
	}

	// update the product quantity
	ok := setQuantityStock(inMemoryMovement.Product, inMemoryMovement.Warehouse, draggedStock, m.enterprise)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery note generation
	if inMemoryMovement.SalesOrderDetail != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*inMemoryMovement.SalesOrderDetail, -abs(inMemoryMovement.Quantity), userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if inMemoryMovement.PurchaseOrderDetail != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*inMemoryMovement.PurchaseOrderDetail, -abs(inMemoryMovement.Quantity), m.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// sales delivery note price
	if inMemoryMovement.SalesDeliveryNote != nil {
		ok = addTotalProductsSalesDeliveryNote(*inMemoryMovement.SalesDeliveryNote, -absf(inMemoryMovement.Price*float64(inMemoryMovement.Quantity)), inMemoryMovement.VatPercent, m.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// purchase delivery note price
	if inMemoryMovement.PurchaseDeliveryNote != nil {
		ok = addTotalProductsPurchaseDeliveryNote(*inMemoryMovement.PurchaseDeliveryNote, -absf(inMemoryMovement.Price*float64(inMemoryMovement.Quantity)), inMemoryMovement.VatPercent, inMemoryMovement.enterprise, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	err = trans.Commit()
	if err != nil {
		return false
	}
	///

	return rowsCount > 0
}

func regenerateDraggedStock(warehouseId string, enterpriseId int32) bool {
	if len(warehouseId) == 0 || len(warehouseId) > 2 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	// select the list with the products with warehouse movements
	sqlStatement := `SELECT product FROM warehouse_movement WHERE warehouse=$1 AND enterprise=$2 GROUP BY product`
	rowsProducts, err := db.Query(sqlStatement, warehouseId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// for each product...
	for rowsProducts.Next() {
		// add the quantity for each row to drag the amount of stock
		var draggedStock int32 = 0

		var productId int32
		rowsProducts.Scan(&productId)

		sqlStatement := `SELECT id,quantity,type FROM warehouse_movement WHERE warehouse=$1 AND product=$2 ORDER BY date_created ASC, id ASC`
		rows, err := db.Query(sqlStatement, warehouseId, productId)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		// for each warehouse movement...
		for rows.Next() {
			var movementId int64
			var quantity int32
			var movementType string
			rows.Scan(&movementId, &quantity, &movementType)

			if movementType == "R" {
				draggedStock = quantity
			} else {
				draggedStock += quantity
			}

			sqlStatement := `UPDATE warehouse_movement SET dragged_stock=$2 WHERE id=$1`
			_, err := db.Exec(sqlStatement, movementId, draggedStock)
			if err != nil {
				log("DB", err.Error())
				trans.Rollback()
				return false
			}
		}
	}

	///
	err = trans.Commit()
	return err == nil
	///
}
