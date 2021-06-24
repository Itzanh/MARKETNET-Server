package main

import (
	"fmt"
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
	SalesOrder            *int32    `json:"salesOrder"`
	SalesOrderDetail      *int32    `json:"salesOrderDetail"`
	SalesInvoice          *int32    `json:"salesInvoice"`
	SalesInvoiceDetail    *int32    `json:"salesInvoiceDetail"`
	SalesDeliveryNote     *int32    `json:"salesDeliveryNote"`
	Description           string    `json:"description"`
	PurchaseOrder         *int32    `json:"purchaseOrder"`
	PurchaseOrderDetail   *int32    `json:"purchaseOrderDetail"`
	PurchaseInvoice       *int32    `json:"purchaseInvoice"`
	PurchaseInvoiceDetail *int32    `json:"purchaseInvoiceDetail"`
	PurchaseDeliveryNote  *int32    `json:"purchaseDeliveryNote"`
	DraggedStock          int32     `json:"draggedStock"`
}

func getWarehouseMovement() []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT * FROM public.warehouse_movement ORDER BY id DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func getWarehouseMovementByWarehouse(warehouseId string) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if len(warehouseId) == 0 || len(warehouseId) > 2 {
		return warehouseMovements
	}

	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE warehouse=$1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, warehouseId)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func getWarehouseMovementRow(movementId int64) WarehouseMovement {
	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE id=$1`
	row := db.QueryRow(sqlStatement, movementId)
	if row.Err() != nil {
		return WarehouseMovement{}
	}

	m := WarehouseMovement{}
	row.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)

	return m
}

func getWarehouseMovementBySalesDeliveryNote(noteId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE sales_delivery_note=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func getWarehouseMovementByPurchaseDeliveryNote(noteId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE purchase_delivery_note=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

type WarehouseMovementSearch struct {
	Search    string     `json:"search"`
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (w *WarehouseMovementSearch) searchWarehouseMovement() []WarehouseMovement {
	fmt.Println("searchWarehouseMovement")
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT warehouse_movement.* FROM warehouse_movement INNER JOIN product ON product.id=warehouse_movement.product WHERE product.name ILIKE $1`
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
	sqlStatement += ` ORDER BY warehouse_movement.id DESC`
	rows, err := db.Query(sqlStatement, parameters...)
	if err != nil {
		fmt.Println(err)
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func (m *WarehouseMovement) isValid() bool {
	return !(len(m.Warehouse) == 0 || len(m.Warehouse) > 2 || m.Product <= 0 || m.Quantity == 0 || len(m.Type) != 1 || (m.Type != "I" && m.Type != "O" && m.Type != "R"))
}

func (m *WarehouseMovement) insertWarehouseMovement() bool {
	if !m.isValid() {
		return false
	}

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
	sqlStatement := `INSERT INTO public.warehouse_movement(warehouse, product, quantity, type, sales_order, sales_order_detail, sales_invoice, sales_invoice_detail, sales_delivery_note, dsc, purchase_order, purchase_order_detail, purchase_invoice, purchase_invoice_details, purchase_delivery_note, dragged_stock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	res, err := db.Exec(sqlStatement, m.Warehouse, m.Product, m.Quantity, m.Type, m.SalesOrder, m.SalesOrderDetail, m.SalesInvoice, m.SalesInvoiceDetail, m.SalesDeliveryNote, m.Description, m.PurchaseOrder, m.PurchaseOrderDetail, m.PurchaseInvoice, m.PurchaseInvoiceDetail, m.PurchaseDeliveryNote, m.DraggedStock)
	if err != nil {
		return false
	}
	// update the product quantity
	ok := setQuantityStock(m.Product, m.Warehouse, m.DraggedStock)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery notes generation
	if m.SalesOrderDetail != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*m.SalesOrderDetail, abs(m.Quantity))
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if m.PurchaseOrderDetail != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*m.PurchaseOrderDetail, abs(m.Quantity))
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

	rows, _ := res.RowsAffected()
	return rows > 0
}

// Abs returns the absolute value of x.
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (m *WarehouseMovement) deleteWarehouseMovement() bool {
	if m.Id <= 0 {
		return false
	}

	inMemoryMovement := getWarehouseMovementRow(m.Id)
	if inMemoryMovement.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	// delete the warehouse movement
	sqlStatement := `DELETE FROM public.warehouse_movement WHERE id=$1`
	res, err := db.Exec(sqlStatement, m.Id)
	if err != nil {
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

	sqlStatement = `SELECT id,quantity,type FROM warehouse_movement WHERE warehouse=$1 AND product=$2 AND date_created>=$3 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, inMemoryMovement.Warehouse, inMemoryMovement.Product, inMemoryMovement.DateCreated)
	if err != nil {
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
			trans.Rollback()
			return false
		}
	}

	// update the product quantity
	ok := setQuantityStock(inMemoryMovement.Product, inMemoryMovement.Warehouse, draggedStock)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery note generation
	if inMemoryMovement.SalesOrderDetail != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*inMemoryMovement.SalesOrderDetail, -abs(inMemoryMovement.Quantity))
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if inMemoryMovement.PurchaseOrderDetail != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*inMemoryMovement.PurchaseOrderDetail, -abs(inMemoryMovement.Quantity))
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

	rowsCount, _ := res.RowsAffected()
	return rowsCount > 0
}

func regenerateDraggedStock(warehouseId string) bool {
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
	sqlStatement := `SELECT product FROM warehouse_movement WHERE warehouse=$1 GROUP BY product`
	rowsProducts, err := db.Query(sqlStatement, warehouseId)
	if err != nil {
		trans.Rollback()
		return false
	}

	// for each product...
	for rowsProducts.Next() {
		// add the quantity for each row to drag the amount of stock
		var draggedStock int32 = 0

		var productId int32
		rowsProducts.Scan(&productId)

		sqlStatement := `SELECT id,quantity,type FROM warehouse_movement WHERE warehouse=$1 AND product=$2 ORDER BY date_created ASC`
		rows, err := db.Query(sqlStatement, warehouseId, productId)
		if err != nil {
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
