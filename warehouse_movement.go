package main

import (
	"time"
)

type WarehouseMovement struct {
	Id                    int64     `json:"id"`
	Warehouse             string    `json:"warehouse"`
	Product               int32     `json:"product"`
	Quantity              int32     `json:"quantity"`
	DateCreated           time.Time `json:"dateCreated"`
	Type                  string    `json:"type"` // O = Out, I = In
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
}

func getWarehouseMovement() []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT * FROM public.warehouse_movement ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func getWarehouseMovementByWarehouse(warehouseId string) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if len(warehouseId) == 0 || len(warehouseId) > 2 {
		return warehouseMovements
	}

	sqlStatement := `SELECT * FROM public.warehouse_movement WHERE warehouse=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, warehouseId)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)
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
	row.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)

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
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)
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
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func (m *WarehouseMovement) isValid() bool {
	return !(len(m.Warehouse) == 0 || len(m.Warehouse) > 2 || m.Product <= 0 || m.Quantity == 0 || len(m.Type) != 1)
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

	sqlStatement := `INSERT INTO public.warehouse_movement(warehouse, product, quantity, type, sales_order, sales_order_detail, sales_invoice, sales_invoice_detail, sales_delivery_note, dsc, purchase_order, purchase_order_detail, purchase_invoice, purchase_invoice_details, purchase_delivery_note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	res, err := db.Exec(sqlStatement, m.Warehouse, m.Product, m.Quantity, m.Type, m.SalesOrder, m.SalesOrderDetail, m.SalesInvoice, m.SalesInvoiceDetail, m.SalesDeliveryNote, m.Description, m.PurchaseOrder, m.PurchaseOrderDetail, m.PurchaseInvoice, m.PurchaseInvoiceDetail, m.PurchaseDeliveryNote)
	if err != nil {
		return false
	}
	ok := addQuantityStock(m.Product, m.Warehouse, m.Quantity)
	if !ok {
		trans.Rollback()
		return false
	}
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

	sqlStatement := `DELETE FROM public.warehouse_movement WHERE id=$1`
	res, err := db.Exec(sqlStatement, m.Id)
	if err != nil {
		return false
	}

	ok := addQuantityStock(inMemoryMovement.Product, inMemoryMovement.Warehouse, -inMemoryMovement.Quantity)
	if !ok {
		trans.Rollback()
		return false
	}
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

	rows, _ := res.RowsAffected()
	return rows > 0
}
