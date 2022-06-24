/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManufacturingOrder struct {
	Id                   int64                  `json:"id" gorm:"index:manufacturing_order_id_enterprise,unique:true,priority:1"`
	OrderDetailId        *int64                 `json:"orderDetailId" gorm:"column:order_detail;index:manufacturing_order_for_stock_pending,priority:4,where:NOT manufactured AND order_detail IS NULL AND NOT complex"`
	OrderDetail          *SalesOrderDetail      `json:"orderDetail" gorm:"foreignKey:OrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId            int32                  `json:"productId" gorm:"column:product;not null:true;index:manufacturing_order_for_stock_pending,priority:2,where:NOT manufactured AND order_detail IS NULL AND NOT complex"`
	Product              Product                `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	TypeId               int32                  `json:"typeId" gorm:"column:type;not null:true"`
	Type                 ManufacturingOrderType `json:"type" gorm:"foreignKey:TypeId,EnterpriseId;references:Id,EnterpriseId"`
	Uuid                 string                 `json:"uuid" gorm:"column:uuid;not null:true;type:uuid;index:manufacturing_order_uuid,unique:true"`
	DateCreated          time.Time              `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone;index:manufacturing_order_date_created,sort:desc"`
	DateLastUpdate       time.Time              `json:"dateLastUpdate" gorm:"column:date_last_update;not null:true;type:timestamp(3) with time zone"`
	Manufactured         bool                   `json:"manufactured" gorm:"column:manufactured;not null:true;index:manufacturing_order_for_stock_pending,priority:3,where:NOT manufactured AND order_detail IS NULL AND NOT complex"`
	DateManufactured     *time.Time             `json:"dateManufactured" gorm:"column:date_manufactured;type:timestamp(3) with time zone"`
	UserManufacturedId   *int32                 `json:"userManufacturedId" gorm:"column:user_manufactured"`
	UserManufactured     *User                  `json:"userManufactured" gorm:"foreignKey:UserManufacturedId,EnterpriseId;references:Id,EnterpriseId"`
	UserCreatedId        int32                  `json:"userCreatedId" gorm:"column:user_created;not null:true"`
	UserCreated          User                   `json:"userCreated" gorm:"foreignKey:UserCreatedId,EnterpriseId;references:Id,EnterpriseId"`
	TagPrinted           bool                   `json:"tagPrinted" gorm:"column:tag_printed;not null:true"`
	DateTagPrinted       *time.Time             `json:"dateTagPrinted" gorm:"column:date_tag_printed;type:timestamp(3) with time zone"`
	OrderId              *int64                 `json:"orderId" gorm:"column:order"`
	Order                *SaleOrder             `json:"order" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	UserTagPrintedId     *int32                 `json:"userTagPrintedId" gorm:"column:user_tag_printed"`
	UserTagPrinted       *User                  `json:"userTagPrinted" gorm:"foreignKey:UserTagPrintedId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId         int32                  `json:"enterprise" gorm:"column:enterprise;not null:true;index:manufacturing_order_for_stock_pending,priority:1,where:NOT manufactured AND order_detail IS NULL AND NOT complex;index:manufacturing_order_id_enterprise,unique:true,priority:2"`
	Enterprise           Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	WarehouseId          string                 `json:"warehouseId" gorm:"column:warehouse;not null:true;type:character(2)"`
	Warehouse            Warehouse              `json:"warehouse" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseMovementId  *int64                 `json:"warehouseMovementId" gorm:"column:warehouse_movement"`
	WarehouseMovement    *WarehouseMovement     `json:"warehouseMovement" gorm:"foreignKey:WarehouseMovementId,EnterpriseId;references:Id,EnterpriseId"`
	QuantityManufactured int32                  `json:"quantityManufactured" gorm:"column:quantity_manufactured;not null:true"`
	Complex              bool                   `json:"-" gorm:"column:complex;not null:true;index:manufacturing_order_for_stock_pending,priority:5,where:NOT manufactured AND order_detail IS NULL AND NOT complex"`
}

func (mo *ManufacturingOrder) TableName() string {
	return "manufacturing_order"
}

type ManufacturingPaginationQuery struct {
	PaginationQuery
	OrderTypeId int32      `json:"orderTypeId"`
	DateStart   *time.Time `json:"dateStart"`
	DateEnd     *time.Time `json:"dateEnd"`
	Status      string     `json:"status"` // "" = All, "M" = Manufactured, "N" = Not manufactured
	Uuid        string     `json:"uuid"`
}

func (q *ManufacturingPaginationQuery) isDefault() bool {
	return q.OrderTypeId == 0 && q.DateStart == nil && q.DateEnd == nil && q.Status == "" && q.Uuid == ""
}

type ManufacturingOrders struct {
	Rows                int64                `json:"rows"`
	ManufacturingOrders []ManufacturingOrder `json:"manufacturingOrders"`
}

func (q *ManufacturingPaginationQuery) getManufacturingOrder(enterpriseId int32) ManufacturingOrders {
	if q.isDefault() {
		return (q.PaginationQuery).getAllManufacturingOrders(enterpriseId)
	} else {
		return q.getManufacturingOrdersByType(enterpriseId)
	}
}

func (q *PaginationQuery) getAllManufacturingOrders(enterpriseId int32) ManufacturingOrders {
	mo := ManufacturingOrders{}
	mo.ManufacturingOrders = make([]ManufacturingOrder, 0)
	result := dbOrm.Model(&ManufacturingOrder{}).Where("enterprise = ?", enterpriseId).Order("date_created DESC").Offset(int(q.Offset)).Limit(int(q.Limit)).Preload(clause.Associations).Find(&mo.ManufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}
	result = dbOrm.Model(&ManufacturingOrder{}).Where("enterprise = ?", enterpriseId).Count(&mo.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}
	return mo
}

func (q *ManufacturingPaginationQuery) getManufacturingOrdersByType(enterpriseId int32) ManufacturingOrders {
	mo := ManufacturingOrders{}
	mo.ManufacturingOrders = make([]ManufacturingOrder, 0)

	if len(q.Uuid) == 36 && checkUUID(q.Uuid) {
		manufacturingOrder := getManufacturingOrderByUUID(q.Uuid, enterpriseId)
		if manufacturingOrder.Id > 0 {
			mo.ManufacturingOrders = append(mo.ManufacturingOrders, manufacturingOrder)
			mo.Rows = 1
		}
		return mo
	} else if checkBase64(q.Uuid) {
		decodedUuid, err := base64ToUuid(q.Uuid)
		if err == nil {
			manufacturingOrder := getManufacturingOrderByUUID(decodedUuid, enterpriseId)
			if manufacturingOrder.Id > 0 {
				mo.ManufacturingOrders = append(mo.ManufacturingOrders, manufacturingOrder)
				mo.Rows = 1
			}
			return mo
		}
	}

	cursor := dbOrm.Model(&ManufacturingOrder{}).Where("enterprise = ?", enterpriseId)
	if q.OrderTypeId != 0 {
		cursor = cursor.Where("type = ?", q.OrderTypeId)
	}
	if q.DateStart != nil {
		cursor = cursor.Where("date_created >= ?", q.DateStart)
	}
	if q.DateEnd != nil {
		cursor = cursor.Where("date_created <= ?", q.DateEnd)
	}
	if q.Status == "M" {
		cursor = cursor.Where("manufactured = ?", true)
	} else if q.Status == "N" {
		cursor = cursor.Where("manufactured = ?", false)
	}
	result := cursor.Order("date_created DESC").Offset(int(q.Offset)).Limit(int(q.Limit)).Preload(clause.Associations).Count(&mo.Rows).Find(&mo.ManufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}

	return mo
}

func getManufacturingOrderRow(manufacturingOrderId int64) ManufacturingOrder {
	o := ManufacturingOrder{}
	result := dbOrm.Model(&ManufacturingOrder{}).Where("id = ?", manufacturingOrderId).Preload(clause.Associations).First(&o)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}
	return o
}

func getManufacturingOrderByUUID(manufacturingOrderUUID string, enterpriseId int32) ManufacturingOrder {
	o := ManufacturingOrder{}
	result := dbOrm.Model(&ManufacturingOrder{}).Where("uuid = ? AND enterprise = ?", manufacturingOrderUUID, enterpriseId).Preload(clause.Associations).First(&o)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}
	return o
}

func getManufacturingOrderRowTransaction(manufacturingOrderId int64, trans gorm.DB) ManufacturingOrder {
	o := ManufacturingOrder{}
	result := trans.Model(&ManufacturingOrder{}).Where("id = ?", manufacturingOrderId).Preload(clause.Associations).First(&o)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}
	return o
}

func getManufacturingOrdersForStockPending(enterpriseId int32, productId int32) []ManufacturingOrder {
	var orders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	result := dbOrm.Model(&ManufacturingOrder{}).Where("manufacturing_order.enterprise = ? AND manufacturing_order.product = ? AND NOT manufacturing_order.manufactured AND manufacturing_order.order_detail IS NULL AND NOT manufacturing_order.complex", enterpriseId, productId).Order("manufacturing_order.date_created DESC").Preload(clause.Associations).Find(&orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return orders
	}
	return orders
}

func (o *ManufacturingOrder) isValid() bool {
	return !((o.OrderDetailId != nil && *o.OrderDetailId <= 0) || o.ProductId <= 0 || (o.OrderId != nil && *o.OrderId <= 0))
}

func (c *ManufacturingOrder) BeforeCreate(tx *gorm.DB) (err error) {
	var manufacturingOrder ManufacturingOrder
	tx.Model(&ManufacturingOrder{}).Last(&manufacturingOrder)
	c.Id = manufacturingOrder.Id + 1
	return nil
}

// ERROR CODES:
// 1. There is no manufacturing order type in the product
// 2. The product is deactivated
func (o *ManufacturingOrder) insertManufacturingOrder(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if !o.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	// generate uuid
	o.Uuid = uuid.New().String()

	product := getProductRow(o.ProductId)
	if product.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	// don't allow deactivated products
	if product.Off {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}
	// get type if it's not specified
	if o.TypeId <= 0 {
		if !product.Manufacturing || product.ManufacturingOrderTypeId == nil {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
		o.TypeId = *product.ManufacturingOrderTypeId
	}

	// get quantity manufactured from the type if it's not specified
	mType := getManufacturingOrderTypeRow(o.TypeId)
	if mType.Id <= 0 || mType.EnterpriseId != o.EnterpriseId {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	o.QuantityManufactured = mType.QuantityManufactured

	// set the warehouse
	if len(o.WarehouseId) == 0 {
		s := getSettingsRecordById(o.EnterpriseId)
		o.WarehouseId = s.DefaultWarehouseId
	}

	o.DateCreated = time.Now()
	o.DateLastUpdate = time.Now()
	o.Manufactured = false
	o.DateManufactured = nil
	o.UserManufacturedId = nil
	o.TagPrinted = false
	o.DateTagPrinted = nil
	o.UserTagPrintedId = nil
	o.WarehouseMovementId = nil

	result := trans.Create(&o)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	result = trans.Model(&SalesOrderDetail{}).Where("id = ?", o.OrderDetailId).Update("status", "D")
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if o.OrderDetailId != nil {
		insertTransactionalLog(o.EnterpriseId, "sales_order_detail", int(*o.OrderDetailId), userId, "U")
	}

	if o.OrderId != nil {
		ok := setSalesOrderState(o.EnterpriseId, *o.OrderId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	ok := addQuantityPendingManufacture(o.ProductId, o.WarehouseId, 1, o.EnterpriseId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	insertTransactionalLog(o.EnterpriseId, "manufacturing_order", int(o.Id), userId, "I")

	return OkAndErrorCodeReturn{Ok: true}
}

func (o *ManufacturingOrder) deleteManufacturingOrder(userId int32, trans *gorm.DB) bool {
	if o.Id <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	inMemoryManufacturingOrder := getManufacturingOrderRow(o.Id)
	if inMemoryManufacturingOrder.Id <= 0 || inMemoryManufacturingOrder.EnterpriseId != o.EnterpriseId || inMemoryManufacturingOrder.Manufactured {
		return false
	}

	insertTransactionalLog(inMemoryManufacturingOrder.EnterpriseId, "manufacturing_order", int(o.Id), userId, "D")

	result := trans.Delete(&ManufacturingOrder{}, "id = ?", o.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if inMemoryManufacturingOrder.OrderDetailId != nil && *inMemoryManufacturingOrder.OrderDetailId > 0 && inMemoryManufacturingOrder.OrderId != nil && *inMemoryManufacturingOrder.OrderId > 0 {
		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", inMemoryManufacturingOrder.OrderDetailId).Update("status", "C")
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(o.EnterpriseId, "sales_order_detail", int(*inMemoryManufacturingOrder.OrderDetailId), userId, "U")

		ok := setSalesOrderState(inMemoryManufacturingOrder.EnterpriseId, *inMemoryManufacturingOrder.OrderId, userId, *trans)
		if !ok {
			return false
		}
	}

	ok := addQuantityPendingManufacture(inMemoryManufacturingOrder.ProductId, inMemoryManufacturingOrder.WarehouseId, -1, inMemoryManufacturingOrder.EnterpriseId, *trans)
	if !ok {
		return false
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

func toggleManufactuedManufacturingOrder(orderId int64, userId int32, enterpriseId int32) bool {
	if orderId <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	settings := getSettingsRecordById(enterpriseId)

	inMemoryManufacturingOrder := getManufacturingOrderRowTransaction(orderId, *trans)
	if inMemoryManufacturingOrder.EnterpriseId != enterpriseId {
		return false
	}

	// validation
	if inMemoryManufacturingOrder.Manufactured && inMemoryManufacturingOrder.DateManufactured != nil && int64(time.Since(*inMemoryManufacturingOrder.DateManufactured).Seconds()) > int64(settings.UndoManufacturingOrderSeconds) {
		return false
	}

	inMemoryManufacturingOrder.Manufactured = !inMemoryManufacturingOrder.Manufactured
	if inMemoryManufacturingOrder.Manufactured {
		now := time.Now()
		inMemoryManufacturingOrder.DateManufactured = &now
		inMemoryManufacturingOrder.UserManufacturedId = &userId
	} else {
		inMemoryManufacturingOrder.DateManufactured = nil
		inMemoryManufacturingOrder.UserManufacturedId = nil
	}

	result := trans.Model(&ManufacturingOrder{}).Where("id = ?", orderId).Updates(map[string]interface{}{
		"manufactured":      inMemoryManufacturingOrder.Manufactured,
		"date_manufactured": inMemoryManufacturingOrder.DateManufactured,
		"user_manufactured": inMemoryManufacturingOrder.UserManufacturedId,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(inMemoryManufacturingOrder.EnterpriseId, "manufacturing_order", int(orderId), userId, "U")

	inMemoryManufacturingOrder = getManufacturingOrderRowTransaction(orderId, *trans)
	if inMemoryManufacturingOrder.Id <= 0 {
		return false
	}
	if inMemoryManufacturingOrder.OrderDetailId != nil && *inMemoryManufacturingOrder.OrderDetailId > 0 {
		var status string
		if inMemoryManufacturingOrder.Manufactured {
			// are all the manufacturing orders manufactured?
			var manufacturedOrders int64
			result = trans.Model(&ManufacturingOrder{}).Where("order_detail = ? AND manufactured", inMemoryManufacturingOrder.OrderDetailId).Count(&manufacturedOrders)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			orderDetail := getSalesOrderDetailRow(*inMemoryManufacturingOrder.OrderDetailId)

			if int32(manufacturedOrders) >= orderDetail.Quantity {
				status = "E"
			} else {
				status = "D"
			}
		} else {
			status = "D"
		}

		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", inMemoryManufacturingOrder.OrderDetailId).Update("status", status)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.EnterpriseId, "sales_order_detail", int(*inMemoryManufacturingOrder.OrderDetailId), userId, "U")

		var orderId int64
		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", inMemoryManufacturingOrder.OrderDetailId).Pluck("\"order\"", &orderId)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		ok := setSalesOrderState(enterpriseId, orderId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// Create / delete warehouse movement
	if inMemoryManufacturingOrder.Manufactured {
		movement := WarehouseMovement{
			WarehouseId:  inMemoryManufacturingOrder.WarehouseId,
			ProductId:    inMemoryManufacturingOrder.ProductId,
			Quantity:     inMemoryManufacturingOrder.QuantityManufactured,
			Type:         "I", // Input
			EnterpriseId: enterpriseId,
		}
		ok := movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		result = trans.Model(&ManufacturingOrder{}).Where("id = ?", inMemoryManufacturingOrder.Id).Update("warehouse_movement", movement.Id)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.EnterpriseId, "manufacturing_order", int(inMemoryManufacturingOrder.Id), userId, "U")

		ok = addQuantityPendingManufacture(inMemoryManufacturingOrder.ProductId, inMemoryManufacturingOrder.WarehouseId, -1, inMemoryManufacturingOrder.EnterpriseId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	} else {
		result = trans.Model(&ManufacturingOrder{}).Where("id = ?", inMemoryManufacturingOrder.Id).Update("warehouse_movement", nil)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.EnterpriseId, "manufacturing_order", int(inMemoryManufacturingOrder.Id), userId, "U")

		movement := getWarehouseMovementRow(*inMemoryManufacturingOrder.WarehouseMovementId)
		ok := movement.deleteWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		ok = addQuantityPendingManufacture(inMemoryManufacturingOrder.ProductId, inMemoryManufacturingOrder.WarehouseId, 1, inMemoryManufacturingOrder.EnterpriseId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// manufacture / undo complex manufacturing orders
	setComplexManufacturingOrderManufacturingOrderManufactured(inMemoryManufacturingOrder.Id, inMemoryManufacturingOrder.Manufactured, inMemoryManufacturingOrder.EnterpriseId, userId, trans)

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func manufacturingOrderAllSaleOrder(saleOrderId int64, userId int32, enterpriseId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.EnterpriseId != enterpriseId {
		return false
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.EnterpriseId)

	if saleOrder.Id <= 0 || len(orderDetails) == 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	for i := 0; i < len(orderDetails); i++ {
		if orderDetails[i].Status == "C" {
			orderDetail := orderDetails[i]

			product := getProductRow(orderDetail.ProductId)
			if product.Id <= 0 || !product.Manufacturing || product.ManufacturingOrderTypeId == nil || *product.ManufacturingOrderTypeId == 0 {
				continue
			}
			manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderTypeId)
			if manufacturingOrderType.Id <= 0 || manufacturingOrderType.QuantityManufactured <= 0 || manufacturingOrderType.Complex {
				continue
			}

			for j := 0; j < int(orderDetail.Quantity); j += int(manufacturingOrderType.QuantityManufactured) {
				o := ManufacturingOrder{}
				o.ProductId = orderDetail.ProductId
				o.OrderDetailId = &orderDetail.Id
				o.OrderId = &saleOrder.Id
				o.UserCreatedId = userId
				o.EnterpriseId = enterpriseId
				o.WarehouseId = orderDetail.WarehouseId
				ok := o.insertManufacturingOrder(userId, trans).Ok
				if !ok {
					trans.Rollback()
					return false
				}
			}
		}
	}

	///
	result := trans.Commit()
	return result.Error == nil
	///
}

type ManufacturingOrderGenerate struct {
	Selection []ManufacturingOrderGenerateSelection `json:"selection"`
}

type ManufacturingOrderGenerateSelection struct {
	OrderId  int64 `json:"orderId"`
	Id       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (orderInfo *ManufacturingOrderGenerate) manufacturingOrderPartiallySaleOrder(userId int32, enterpriseId int32) bool {
	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	for i := 0; i < len(orderInfo.Selection); i++ {
		orderInfoSelection := orderInfo.Selection[i]
		// get the sale order and it's details
		saleOrder := getSalesOrderRow(orderInfoSelection.OrderId)
		if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId || len(orderInfo.Selection) == 0 {
			trans.Rollback()
			return false
		}

		// get the details
		orderDetail := getSalesOrderDetailRow(orderInfoSelection.Id)
		if orderDetail.Id <= 0 || orderDetail.OrderId != orderInfoSelection.OrderId || orderInfoSelection.Quantity == 0 || orderInfoSelection.Quantity > orderDetail.Quantity {
			trans.Rollback()
			return false
		}
		if orderDetail.Status == "C" {
			manufacturingOrderType := getManufacturingOrderTypeRow(orderDetail.ProductId)
			if manufacturingOrderType.Id <= 0 || manufacturingOrderType.QuantityManufactured <= 0 || manufacturingOrderType.Complex {
				continue
			}
			for j := 0; j < int(orderInfoSelection.Quantity); j += int(manufacturingOrderType.QuantityManufactured) {
				o := ManufacturingOrder{}
				o.ProductId = orderDetail.ProductId
				o.OrderDetailId = &orderDetail.Id
				o.OrderId = &orderDetail.OrderId
				o.UserCreatedId = userId
				o.EnterpriseId = enterpriseId
				o.WarehouseId = orderDetail.WarehouseId
				ok := o.insertManufacturingOrder(userId, trans).Ok
				if !ok {
					trans.Rollback()
					return false
				}
			}
		}
	}

	///
	result := trans.Commit()
	return result.Error == nil
	///
}

func manufacturingOrderTagPrinted(orderId int64, userId int32, enterpriseId int32) bool {
	if orderId <= 0 {
		return false
	}

	inMemoryManufacturingOrder := getManufacturingOrderRow(orderId)
	if inMemoryManufacturingOrder.Id <= 0 || inMemoryManufacturingOrder.EnterpriseId != enterpriseId {
		return false
	}
	if inMemoryManufacturingOrder.TagPrinted {
		return false
	}

	result := dbOrm.Model(&ManufacturingOrder{}).Where("id = ? AND enterprise = ?", orderId, enterpriseId).Updates(map[string]interface{}{
		"tag_printed":      true,
		"date_tag_printed": time.Now(),
		"user_tag_printed": userId,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "manufacturing_order", int(orderId), userId, "U")

	return true
}

type MultipleManufacturingOrders struct {
	Order   ManufacturingOrder `json:"order"`
	Quantiy int                `json:"quantity"`
}

func (o *MultipleManufacturingOrders) insertMultipleManufacturingOrders(userId int32) OkAndErrorCodeReturn {
	if !o.Order.isValid() || o.Quantiy <= 0 || o.Quantiy > 10000 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}

	for i := 0; i < o.Quantiy; i++ {
		ok := o.Order.insertManufacturingOrder(userId, trans)
		if !ok.Ok {
			trans.Rollback()
			return ok
		}
	}

	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}
}
