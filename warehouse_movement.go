/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WarehouseMovement struct {
	Id                     int64                 `json:"id" gorm:"index:warehouse_movement_id_enterprise,unique:true,priority:1"`
	WarehouseId            string                `json:"warehouseId" gorm:"column:warehouse;type:character(2);not null"`
	Warehouse              Warehouse             `json:"warehouse" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId              int32                 `json:"productId" gorm:"column:product;not null"`
	Product                Product               `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity               int32                 `json:"quantity" gorm:"column:quantity;not null"`
	DateCreated            time.Time             `json:"dateCreated" gorm:"column:date_created;not null;type:timestamp(3) with time zone"`
	Type                   string                `json:"type" gorm:"column:type;type:character(1);not null:true"` // O = Out, I = In, R = Inventory regularization
	SalesOrderId           *int64                `json:"salesOrderId" gorm:"column:sales_order"`
	SalesOrder             *SaleOrder            `json:"salesOrder" gorm:"foreignKey:SalesOrderId,EnterpriseId;references:Id,EnterpriseId"`
	SalesOrderDetailId     *int64                `json:"salesOrderDetailId" gorm:"column:sales_order_detail"`
	SalesOrderDetail       *SalesOrderDetail     `json:"salesOrderDetail" gorm:"foreignKey:SalesOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	SalesDeliveryNoteId    *int64                `json:"salesDeliveryNoteId" gorm:"column:sales_delivery_note"`
	SalesDeliveryNote      *SalesDeliveryNote    `json:"salesDeliveryNote" gorm:"foreignKey:SalesDeliveryNoteId,EnterpriseId;references:Id,EnterpriseId"`
	Description            string                `json:"description" gorm:"column:dsc;type:text;not null"`
	PurchaseOrderId        *int64                `json:"purchaseOrderId" gorm:"column:purchase_order"`
	PurchaseOrder          *PurchaseOrder        `json:"purchaseOrder" gorm:"foreignKey:PurchaseOrderId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseOrderDetailId  *int64                `json:"purchaseOrderDetailId" gorm:"column:purchase_order_detail"`
	PurchaseOrderDetail    *PurchaseOrderDetail  `json:"purchaseOrderDetail" gorm:"foreignKey:PurchaseOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseDeliveryNoteId *int64                `json:"purchaseDeliveryNoteId" gorm:"column:purchase_delivery_note"`
	PurchaseDeliveryNote   *PurchaseDeliveryNote `json:"purchaseDeliveryNote" gorm:"foreignKey:PurchaseDeliveryNoteId,EnterpriseId;references:Id,EnterpriseId"`
	DraggedStock           int32                 `json:"draggedStock" gorm:"column:dragged_stock;not null:true"`
	Price                  float64               `json:"price" gorm:"column:price;not null:true;type:numeric(14,6)"`
	VatPercent             float64               `json:"vatPercent" gorm:"column:vat_percent;not null:true;type:numeric(14,6)"`
	TotalAmount            float64               `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	EnterpriseId           int32                 `json:"-" gorm:"column:enterprise;not null:true;index:warehouse_movement_id_enterprise,unique:true,priority:2"`
	Enterprise             Settings              `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Manual                 bool                  `json:"manual" gorm:"column:manual;not null:true;type:boolean;default:false"`
}

func (w *WarehouseMovement) TableName() string {
	return "warehouse_movement"
}

type WarehouseMovements struct {
	Rows      int64               `json:"rows"`
	Movements []WarehouseMovement `json:"movements"`
}

func (q *PaginationQuery) getWarehouseMovement() WarehouseMovements {
	wm := WarehouseMovements{}
	if !q.isValid() {
		return wm
	}

	wm.Movements = make([]WarehouseMovement, 0)
	cursor := dbOrm.Model(&WarehouseMovement{}).Where("enterprise = ?", q.enterprise).Order("id DESC").Offset(int(q.Offset)).Limit(int(q.Limit))
	if cursor.Error != nil {
		log("DB", cursor.Error.Error())
		return wm
	}
	result := cursor.Preload(clause.Associations).Find(&wm.Movements)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return wm
	}
	result = cursor.Count(&wm.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return wm
	}
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

	cursor := dbOrm.Model(&WarehouseMovement{}).Where("warehouse = ? AND enterprise = ?", w.WarehouseId, w.enterprise).Order("id DESC").Offset(int(w.Offset)).Limit(int(w.Limit))
	if cursor.Error != nil {
		log("DB", cursor.Error.Error())
		return wm
	}
	result := cursor.Preload(clause.Associations).Find(&wm.Movements)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return wm
	}
	result = cursor.Count(&wm.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return wm
	}
	return wm
}

func getWarehouseMovementRow(movementId int64) WarehouseMovement {
	m := WarehouseMovement{}
	result := dbOrm.Model(&WarehouseMovement{}).Where("id = ?", movementId).Preload(clause.Associations).First(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return m
	}
	return m
}

func getWarehouseMovementBySalesDeliveryNote(noteId int64, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	result := dbOrm.Model(&WarehouseMovement{}).Where("sales_delivery_note = ? AND enterprise = ?", noteId, enterpriseId).Preload(clause.Associations).Find(&warehouseMovements)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return warehouseMovements
	}
	return warehouseMovements
}

func getWarehouseMovementByPurchaseDeliveryNote(noteId int64, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if noteId <= 0 {
		return warehouseMovements
	}

	result := dbOrm.Model(&WarehouseMovement{}).Where("purchase_delivery_note = ? AND enterprise = ?", noteId, enterpriseId).Preload(clause.Associations).Find(&warehouseMovements)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return warehouseMovements
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
	cursor := dbOrm.Model(&WarehouseMovement{}).Where("warehouse_movement.enterprise = ?", w.enterprise).Joins("INNER JOIN product ON product.id=warehouse_movement.product").Where("product.name ILIKE ?", "%"+w.Search+"%")
	if w.DateStart != nil {
		cursor = cursor.Where("warehouse_movement.date_created >= ?", w.DateStart)
	}
	if w.DateEnd != nil {
		cursor = cursor.Where("warehouse_movement.date_created <= ?", w.DateEnd)
	}
	result := cursor.Order("warehouse_movement.id DESC").Offset(int(w.Offset)).Limit(int(w.Limit)).Preload(clause.Associations).Find(&wm.Movements).Count(&wm.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return wm
	}

	return wm
}

func (m *WarehouseMovement) isValid() bool {
	return !(len(m.WarehouseId) == 0 || len(m.WarehouseId) > 2 || m.ProductId <= 0 || m.Quantity == 0 || len(m.Type) != 1 || (m.Type != "I" && m.Type != "O" && m.Type != "R") || len(m.Description) > 3000)
}

func (wm *WarehouseMovement) BeforeCreate(tx *gorm.DB) (err error) {
	var warehouseMovement WarehouseMovement
	tx.Model(&WarehouseMovement{}).Last(&warehouseMovement)
	wm.Id = warehouseMovement.Id + 1
	return nil
}

func (m *WarehouseMovement) insertWarehouseMovement(userId int32, trans *gorm.DB) bool {
	if !m.isValid() {
		return false
	}

	m.TotalAmount = absf((m.Price * float64(m.Quantity)) * (1 + (m.VatPercent / 100)))

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	// get the dragged stock
	if m.Type != "R" {
		var dragged_stock int32
		result := trans.Model(&WarehouseMovement{}).Where("warehouse = ? AND product = ?", m.WarehouseId, m.ProductId).Order("date_created DESC").Limit(1).Pluck("dragged_stock", &dragged_stock)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
		m.DraggedStock = dragged_stock + m.Quantity
	} else { // Inventory regularization
		m.DraggedStock = m.Quantity
	}

	m.DateCreated = time.Now()

	// insert the movement
	result := trans.Create(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(m.EnterpriseId, "warehouse_movement", int(m.Id), userId, "I")

	// update the product quantity
	ok := setQuantityStock(m.ProductId, m.WarehouseId, m.DraggedStock, m.EnterpriseId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery notes generation
	if m.SalesOrderDetailId != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*m.SalesOrderDetailId, abs(m.Quantity), userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if m.PurchaseOrderDetailId != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*m.PurchaseOrderDetailId, abs(m.Quantity), m.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// sales delivery note price
	if m.SalesDeliveryNoteId != nil {
		ok = addTotalProductsSalesDeliveryNote(*m.SalesDeliveryNoteId, absf(m.Price*float64(m.Quantity)), m.VatPercent, m.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// purchase delivery note price
	if m.PurchaseDeliveryNoteId != nil {
		ok = addTotalProductsPurchaseDeliveryNote(*m.PurchaseDeliveryNoteId, absf(m.Price*float64(m.Quantity)), m.VatPercent, m.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	if beginTransaction {
		///
		result = trans.Commit()
		return result.Error == nil
		///
	}
	return true
}

func (m *WarehouseMovement) deleteWarehouseMovement(userId int32, trans *gorm.DB) bool {
	if m.Id <= 0 {
		return false
	}

	inMemoryMovement := getWarehouseMovementRow(m.Id)
	if inMemoryMovement.Id <= 0 || inMemoryMovement.EnterpriseId != m.EnterpriseId {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	insertTransactionalLog(m.EnterpriseId, "warehouse_movement", int(m.Id), userId, "D")

	// delete the warehouse movement
	result := trans.Delete(&WarehouseMovement{}, "id = ? AND enterprise = ?", m.Id, m.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// update the dragged stock
	var draggedStock int32
	if inMemoryMovement.Type != "R" {
		draggedStock = inMemoryMovement.DraggedStock - inMemoryMovement.Quantity
	} else {
		result = trans.Model(&WarehouseMovement{}).Where("warehouse = ? AND product = ? AND date_created <= ?", inMemoryMovement.WarehouseId, inMemoryMovement.ProductId, inMemoryMovement.DateCreated).Order("date_created DESC").Limit(1).Pluck("dragged_stock", &draggedStock)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	var draggedStocks []WarehouseMovementDraggedStock = make([]WarehouseMovementDraggedStock, 0)
	result = trans.Model(&WarehouseMovement{}).Where("warehouse = ? AND product = ? AND date_created >= ?", inMemoryMovement.WarehouseId, inMemoryMovement.ProductId, inMemoryMovement.DateCreated).Order("date_created ASC, id ASC").Find(&draggedStocks)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	for i := 0; i < len(draggedStocks); i++ {
		d := draggedStocks[i]

		if d.Type == "R" {
			draggedStock = d.Quantity
		} else {
			draggedStock += d.Quantity
		}

		result = trans.Model(&WarehouseMovement{}).Where("id = ?", d.Id).Update("dragged_stock", draggedStock)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///

	// update the product quantity
	ok := setQuantityStock(inMemoryMovement.ProductId, inMemoryMovement.WarehouseId, draggedStock, m.EnterpriseId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}
	// delivery note generation
	if inMemoryMovement.SalesOrderDetailId != nil {
		ok = addQuantityDeliveryNoteSalesOrderDetail(*inMemoryMovement.SalesOrderDetailId, -abs(inMemoryMovement.Quantity), userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	if inMemoryMovement.PurchaseOrderDetailId != nil {
		ok = addQuantityDeliveryNotePurchaseOrderDetail(*inMemoryMovement.PurchaseOrderDetailId, -abs(inMemoryMovement.Quantity), m.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// sales delivery note price
	if inMemoryMovement.SalesDeliveryNoteId != nil {
		ok = addTotalProductsSalesDeliveryNote(*inMemoryMovement.SalesDeliveryNoteId, -absf(inMemoryMovement.Price*float64(inMemoryMovement.Quantity)), inMemoryMovement.VatPercent, m.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// purchase delivery note price
	if inMemoryMovement.PurchaseDeliveryNoteId != nil {
		ok = addTotalProductsPurchaseDeliveryNote(*inMemoryMovement.PurchaseDeliveryNoteId, -absf(inMemoryMovement.Price*float64(inMemoryMovement.Quantity)), inMemoryMovement.VatPercent, inMemoryMovement.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false
		}
		///
	}

	return true
}

type WarehouseMovementDraggedStock struct {
	Id       int64
	Quantity int32
	Type     string
}

func regenerateDraggedStock(warehouseId string, enterpriseId int32) bool {
	if len(warehouseId) == 0 || len(warehouseId) > 2 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	// select the list with the products with warehouse movements
	var productIds []int32 = make([]int32, 0)
	result := dbOrm.Model(&WarehouseMovement{}).Where("warehouse = ? AND enterprise = ?", warehouseId, enterpriseId).Group("product").Select("product").Find(&productIds)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// for each product...
	var draggedStock int32
	var productId int32
	for i := 0; i < len(productIds); i++ {
		draggedStock = 0
		// add the quantity for each row to drag the amount of stock
		productId = productIds[i]

		rows, err := dbOrm.Model(&WarehouseMovement{}).Where("warehouse = ? AND product = ?", warehouseId, productId).Order("date_created ASC, id ASC").Rows()
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

			result = trans.Model(&WarehouseMovement{}).Where("id = ?", movementId).Update("dragged_stock", draggedStock)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}
		}

		rows.Close()
	}

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func regenerateDraggedStockAllWarehouses(enterpriseId int32) bool {
	warehouses := getWarehouses(enterpriseId)
	for i := 0; i < len(warehouses); i++ {
		if !regenerateDraggedStock(warehouses[i].Id, enterpriseId) {
			return false
		}
	}
	return true
}

type WarehouseMovementRelations struct {
	PurchaseDeliveryNoteName   *string                     `json:"purchaseDeliveryNoteName"`
	PurchaseOrderName          *string                     `json:"purchaseOrderName"`
	SalesDeliveryNoteName      *string                     `json:"saleDeliveryNoteName"`
	SalesOrderName             *string                     `json:"saleOrderName"`
	ManufacturingOrders        []ManufacturingOrder        `json:"manufacturingOrders"`
	ComplexManufacturingOrders []ComplexManufacturingOrder `json:"complexManufacturingOrders"`
	TransferBetweenWarehouses  []TransferBetweenWarehouses `json:"transferBetweenWarehouses"`
}

func getWarehouseMovementRelations(warehouseMovementId int64, enterpriseId int32) WarehouseMovementRelations {
	r := WarehouseMovementRelations{}
	r.ManufacturingOrders = make([]ManufacturingOrder, 0)
	r.ComplexManufacturingOrders = make([]ComplexManufacturingOrder, 0)
	r.TransferBetweenWarehouses = make([]TransferBetweenWarehouses, 0)

	movement := getWarehouseMovementRow(warehouseMovementId)

	if movement.PurchaseDeliveryNoteId != nil {
		purchaseDeliveryNoteName := getPurchaseDeliveryNoteRow(*movement.PurchaseDeliveryNoteId).DeliveryNoteName
		r.PurchaseDeliveryNoteName = &purchaseDeliveryNoteName
	}
	if movement.PurchaseOrderId != nil {
		purchaseOrderName := getPurchaseOrderRow(*movement.PurchaseOrderId).OrderName
		r.PurchaseOrderName = &purchaseOrderName
	}
	if movement.SalesDeliveryNoteId != nil {
		salesDeliveryNoteName := getSalesDeliveryNoteRow(*movement.SalesDeliveryNoteId).DeliveryNoteName
		r.SalesDeliveryNoteName = &salesDeliveryNoteName
	}
	if movement.SalesOrderId != nil {
		salesOrderName := getSalesOrderRow(*movement.SalesOrderId).OrderName
		r.SalesOrderName = &salesOrderName
	}

	// complex manufacturing orders
	var complexManufacturingOrderDetails []ComplexManufacturingOrderManufacturingOrder
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("warehouse_movement = ? AND enterprise = ?", warehouseMovementId, enterpriseId).Find(&complexManufacturingOrderDetails)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	for i := 0; i < len(complexManufacturingOrderDetails); i++ {
		cmo := getComplexManufacturingOrderRow(complexManufacturingOrderDetails[i].ComplexManufacturingOrderId)
		r.ComplexManufacturingOrders = append(r.ComplexManufacturingOrders, cmo)

	}

	// manufacturing orders
	result = dbOrm.Model(&ManufacturingOrder{}).Where("warehouse_movement = ? AND enterprise = ?", warehouseMovementId, enterpriseId).Find(&r.ManufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	// transfer between warehouses
	var transferBetweenWarehousesIds []int64
	result = dbOrm.Model(&TransferBetweenWarehousesDetail{}).Where("(warehouse_movement_out = @warehouseMovementId OR warehouse_movement_in = @warehouseMovementId) AND enterprise = @enterpriseId", sql.Named("warehouseMovementId", warehouseMovementId), sql.Named("enterpriseId", enterpriseId)).Select("transfer_between_warehouses").Distinct().Scan(&transferBetweenWarehousesIds)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	for i := 0; i < len(transferBetweenWarehousesIds); i++ {
		r.TransferBetweenWarehouses = append(r.TransferBetweenWarehouses, getTransferBetweenWarehousesRow(transferBetweenWarehousesIds[i]))
	}

	return r
}
