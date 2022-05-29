package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ComplexManufacturingOrder struct {
	Id                         int64                  `json:"id" gorm:"index:complex_manufacturing_order_id_enterprise,unique:true,priotity:1"`
	TypeId                     int32                  `json:"typeId" gorm:"column:type;not null:tre"`
	Type                       ManufacturingOrderType `json:"type" gorm:"foreignKey:TypeId,EnterpriseId;references:Id,EnterpriseId"`
	Manufactured               bool                   `json:"manufactured" gorm:"column:manufactured;not null:true"`
	DateManufactured           *time.Time             `json:"dateManufactured" gorm:"column:date_manufactured;type:timestamp(3) with time zone"`
	UserManufacturedId         *int32                 `json:"userManufacturedId" gorm:"column:user_manufactured"`
	UserManufactured           *User                  `json:"userManufactured" gorm:"foreignKey:UserManufacturedId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId               int32                  `json:"enterprise" gorm:"column:enterprise;not null:true;index:complex_manufacturing_order_id_enterprise,unique:true,priotity:2"`
	Enterprise                 Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	QuantityPendingManufacture int32                  `json:"quantityPendingManufacture" gorm:"column:quantity_pending_manufacture;not null:true"`
	QuantityManufactured       int32                  `json:"quantityManufactured" gorm:"column:quantity_manufactured;not null:true"`
	WarehouseId                string                 `json:"warehouseId" gorm:"column:warehouse;not null:true;type:character(2)"`
	Warehouse                  Warehouse              `json:"warehouse" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated                time.Time              `json:"dateCreated" gorm:"column:date_created;type:timestamp(3) with time zone;not null:true"`
	Uuid                       string                 `json:"uuid" gorm:"column:uuid;not null:true;type:uuid"`
	UserCreatedId              int32                  `json:"userCreatedId" gorm:"column:user_created;not null:true"`
	UserCreated                User                   `json:"userCreated" gorm:"foreignKey:UserCreatedId,EnterpriseId;references:Id,EnterpriseId"`
	TagPrinted                 bool                   `json:"tagPrinted" gorm:"column:tag_printed;not null:true"`
	DateTagPrinted             *time.Time             `json:"dateTagPrinted" gorm:"column:date_tag_printed;type:timestamp(3) with time zone"`
	UserTagPrintedId           *int32                 `json:"userTagPrintedId" gorm:"column:user_tag_printed"`
	UserTagPrinted             *User                  `json:"userTagPrinted" gorm:"foreignKey:UserTagPrintedId,EnterpriseId;references:Id,EnterpriseId"`
}

func (c *ComplexManufacturingOrder) TableName() string {
	return "complex_manufacturing_order"
}

type ComplexManufacturingOrders struct {
	Rows                       int64                       `json:"rows"`
	ComplexManufacturingOrders []ComplexManufacturingOrder `json:"complexManufacturingOrder"`
}

func (q *ManufacturingPaginationQuery) getComplexManufacturingOrder(enterpriseId int32) ComplexManufacturingOrders {
	if q.isDefault() {
		return (q.PaginationQuery).getAllComplexManufacturingOrders(enterpriseId)
	} else {
		return q.getComplexManufacturingOrdersByType(enterpriseId)
	}
}

func (q *PaginationQuery) getAllComplexManufacturingOrders(enterpriseId int32) ComplexManufacturingOrders {
	mo := ComplexManufacturingOrders{}
	mo.ComplexManufacturingOrders = make([]ComplexManufacturingOrder, 0)
	result := dbOrm.Model(&ComplexManufacturingOrder{}).Where("enterprise = ?", enterpriseId).Order("date_created DESC").Offset(int(q.Offset)).Limit(int(q.Limit)).Preload(clause.Associations).Find(&mo.ComplexManufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}
	result = dbOrm.Model(&ComplexManufacturingOrder{}).Where("enterprise = ?", enterpriseId).Count(&mo.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}
	return mo
}

func (q *ManufacturingPaginationQuery) getComplexManufacturingOrdersByType(enterpriseId int32) ComplexManufacturingOrders {
	mo := ComplexManufacturingOrders{}
	mo.ComplexManufacturingOrders = make([]ComplexManufacturingOrder, 0)

	cursor := dbOrm.Model(&ComplexManufacturingOrder{}).Where("enterprise = ?", enterpriseId)
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
	result := cursor.Order("date_created DESC").Offset(int(q.Offset)).Limit(int(q.Limit)).Preload(clause.Associations).Count(&mo.Rows).Find(&mo.ComplexManufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return mo
	}

	return mo
}

func getComplexManufacturingOrderRow(complexManufacturingOrderId int64) ComplexManufacturingOrder {
	c := ComplexManufacturingOrder{}
	result := dbOrm.Model(&ComplexManufacturingOrder{}).Where("id = ?", complexManufacturingOrderId).Preload(clause.Associations).First(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return c
	}
	return c
}

func getComplexManufacturingOrderRowTransaction(complexManufacturingOrderId int64, trans gorm.DB) ComplexManufacturingOrder {
	c := ComplexManufacturingOrder{}
	result := trans.Model(&ComplexManufacturingOrder{}).Where("id = ?", complexManufacturingOrderId).Preload(clause.Associations).First(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return c
	}
	return c
}

// Specify a negative number to substract
// DOES NOT OPEN A TRANSACTION
func addQuantityPendingManufactureComplexManufacturingOrder(complexManufacturingOrderId int64, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) bool {
	result := trans.Model(&ComplexManufacturingOrder{}).Where("id = ?", complexManufacturingOrderId).Update("quantity_pending_manufacture", quantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "complex_manufacturing_order", int(complexManufacturingOrderId), userId, "U")

	return true
}

// Specify a negative number to substract
// DOES NOT OPEN A TRANSACTION
func addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId int64, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) bool {
	result := trans.Model(&ComplexManufacturingOrder{}).Where("id = ?", complexManufacturingOrderId).Update("quantity_manufactured", quantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "complex_manufacturing_order", int(complexManufacturingOrderId), userId, "U")

	return true
}

func complexManufacturingOrderAllSaleOrder(saleOrderId int64, userId int32, enterpriseId int32) bool {
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId {
		return false
	}
	details := getSalesOrderDetail(saleOrderId, enterpriseId)
	if len(details) == 0 {
		return false
	}

	return complexManufacturingOrerGeneration(userId, enterpriseId, details)
}

func (orderInfo *ManufacturingOrderGenerate) complexManufacturingOrderPartiallySaleOrder(userId int32, enterpriseId int32) bool {
	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(orderInfo.Selection); i++ {
		orderInfoSelection := orderInfo.Selection[i]
		// get the sale order and it's details
		saleOrder := getSalesOrderRow(orderInfoSelection.OrderId)
		if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId || len(orderInfo.Selection) == 0 {
			return false
		}

		orderDetail := getSalesOrderDetailRow(orderInfoSelection.Id)
		if orderDetail.Id <= 0 || orderDetail.OrderId != orderInfoSelection.OrderId || orderInfoSelection.Quantity == 0 || orderInfoSelection.Quantity > orderDetail.Quantity {
			return false
		}
		if orderDetail.Status == "C" {
			saleOrderDetails = append(saleOrderDetails, orderDetail)
		}
	}

	return complexManufacturingOrerGeneration(userId, enterpriseId, saleOrderDetails)
}

func complexManufacturingOrerGeneration(userId int32, enterpriseId int32, details []SalesOrderDetail) bool {
	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	for i := 0; i < len(details); i++ {
		if details[i].Status != "C" {
			continue
		}
		orderDetail := details[i]

		product := getProductRow(orderDetail.ProductId)
		if product.Id <= 0 || !product.Manufacturing || product.ManufacturingOrderTypeId == nil || *product.ManufacturingOrderTypeId == 0 {
			continue
		}
		manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderTypeId)
		if manufacturingOrderType.Id <= 0 || manufacturingOrderType.QuantityManufactured <= 0 || !manufacturingOrderType.Complex {
			continue
		}

		components := getManufacturingOrderTypeComponents(manufacturingOrderType.Id, enterpriseId)
		var component *ManufacturingOrderTypeComponents = nil
		for i := 0; i < len(components); i++ {
			if components[i].Type == "O" && components[i].ProductId == product.Id {
				component = &components[i]
				break
			}
		}
		if component == nil {
			trans.Rollback()
			return false
		}

		for j := 0; j < int(orderDetail.Quantity); j += int(component.Quantity) {
			cmo := ComplexManufacturingOrder{
				TypeId:       manufacturingOrderType.Id,
				EnterpriseId: enterpriseId,
				WarehouseId:  orderDetail.WarehouseId,
			}
			ok, _ := cmo.insertComplexManufacturingOrder(1, trans)
			if !ok {
				trans.Rollback()
				return false
			}

			id := getPendingComplexManufacturingOrderOutputsWithoutSaleOrderDetail(product.Id)
			if id == nil || *id <= 0 {
				continue
			}

			result := trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("id = ?", id).Update("sale_order_detail", orderDetail.Id)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(enterpriseId, "complex_manufacturing_order_manufacturing_order", int(*id), userId, "U")

			result = trans.Model(&SalesOrderDetail{}).Where("id = ?", orderDetail.Id).Update("status", "D")
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(enterpriseId, "sales_order_detail", int(orderDetail.Id), userId, "U")

			ok = setSalesOrderState(enterpriseId, orderDetail.OrderId, userId, *trans)
			if !ok {
				trans.Rollback()
				return false
			}
		} // for j := 0; j < int(orderDetail.Quantity); j += int(component.Quantity) {
	} // for i := 0; i < len(details); i++

	///
	result := trans.Commit()
	return result.Error == nil
	///
}

func (c *ComplexManufacturingOrder) isValid() bool {
	return !(c.TypeId <= 0 || c.EnterpriseId == 0)
}

func (c *ComplexManufacturingOrder) BeforeCreate(tx *gorm.DB) (err error) {
	var complexManufacturingOrder ComplexManufacturingOrder
	tx.Model(&ComplexManufacturingOrder{}).Last(&complexManufacturingOrder)
	c.Id = complexManufacturingOrder.Id + 1
	return nil
}

func (c *ComplexManufacturingOrder) insertComplexManufacturingOrder(userId int32, trans *gorm.DB) (bool, *int64) {
	if !c.isValid() {
		return false, nil
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false, nil
		}
		///
	}

	// generate uuid
	c.Uuid = uuid.New().String()

	// set the warehouse
	if len(c.WarehouseId) == 0 {
		s := getSettingsRecordById(c.EnterpriseId)
		c.WarehouseId = s.DefaultWarehouseId
	}

	c.Manufactured = false
	c.DateManufactured = nil
	c.UserManufacturedId = nil
	c.QuantityPendingManufacture = 0
	c.QuantityManufactured = 0
	c.DateCreated = time.Now()
	c.TagPrinted = false
	c.DateTagPrinted = nil
	c.UserTagPrintedId = nil

	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, nil
	}

	complexManufacturingOrder := getComplexManufacturingOrderRowTransaction(c.Id, *trans)

	insertTransactionalLog(c.EnterpriseId, "complex_manufacturing_order", int(c.Id), userId, "I")

	components := getManufacturingOrderTypeComponents(c.TypeId, c.EnterpriseId)

	var subOrders []ComplexManufacturingOrderManufacturingOrder = make([]ComplexManufacturingOrderManufacturingOrder, 0)
	for i := 0; i < len(components); i++ {
		if components[i].Type != "I" { // Only Input
			continue
		}

		manufacturingOrderTypeComponent := components[i]
		if manufacturingOrderTypeComponent.Id <= 0 || manufacturingOrderTypeComponent.EnterpriseId != c.EnterpriseId || manufacturingOrderTypeComponent.Quantity <= 0 {
			trans.Rollback()
			return false, nil
		}

		stock := getStockRow(manufacturingOrderTypeComponent.ProductId, c.WarehouseId, c.EnterpriseId)
		if stock.QuantityAvaiable >= manufacturingOrderTypeComponent.Quantity {
			// there is stock for the manufacturing, we make a manufacturing order to reserve the stock
			wm := WarehouseMovement{
				WarehouseId:  c.WarehouseId,
				ProductId:    manufacturingOrderTypeComponent.ProductId,
				Quantity:     manufacturingOrderTypeComponent.Quantity,
				Type:         "O",
				EnterpriseId: c.EnterpriseId,
			}
			ok := wm.insertWarehouseMovement(userId, trans)
			if !ok {
				if beginTransaction {
					trans.Rollback()
				}
				return false, nil
			}

			c := ComplexManufacturingOrderManufacturingOrder{
				Type:                              "I",
				ComplexManufacturingOrderId:       complexManufacturingOrder.Id,
				EnterpriseId:                      complexManufacturingOrder.EnterpriseId,
				WarehouseMovementId:               &wm.Id,
				ProductId:                         manufacturingOrderTypeComponent.ProductId,
				ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
				Manufactured:                      true,
			}
			subOrders = append(subOrders, c)
		} else { // if stock.QuantityAvaialbe >= manufacturingOrderTypeComponent.Quantity {
			// the product is from a supplier or from manufacturing?
			product := getProductRow(manufacturingOrderTypeComponent.ProductId)
			if product.Manufacturing {
				if product.ManufacturingOrderTypeId == nil {
					continue
				}
				manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderTypeId)
				if manufacturingOrderType.Complex {
					cmo := ComplexManufacturingOrder{
						TypeId:       manufacturingOrderType.Id,
						WarehouseId:  c.WarehouseId,
						EnterpriseId: c.EnterpriseId,
					}
					ok, recursiveComplexManufacturingOrderId := cmo.insertComplexManufacturingOrder(userId, trans) // RECURSIVITY
					if !ok || recursiveComplexManufacturingOrderId == nil {
						trans.Rollback()
						return false, nil
					}

					recursiveComponents := getComplexManufacturingOrderManufacturingOrder(*recursiveComplexManufacturingOrderId, c.EnterpriseId)
					var recursiveComponent *ComplexManufacturingOrderManufacturingOrder

					for j := 0; j < len(recursiveComponents); j++ {
						if recursiveComponents[i].Type == "O" && recursiveComponents[i].ProductId == product.Id {
							recursiveComponent = &recursiveComponents[i]
							break
						}
					}

					if recursiveComponent == nil {
						trans.Rollback()
						return false, nil
					}

					c := ComplexManufacturingOrderManufacturingOrder{
						Type:                        "I",
						ComplexManufacturingOrderId: complexManufacturingOrder.Id,
						EnterpriseId:                complexManufacturingOrder.EnterpriseId,
						ComplexManufacturingOrderManufacturingOrderOutputId: &recursiveComponent.Id,
						ProductId:                         manufacturingOrderTypeComponent.ProductId,
						ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
						Manufactured:                      false,
					}
					subOrders = append(subOrders, c)
				} else { // if product.ManufacturingOrderType == nil {
					// we search existing orders to make stock (without order and without complex order)
					manufacturingOrders := getManufacturingOrdersForStockPending(c.EnterpriseId, manufacturingOrderTypeComponent.ProductId)
					var quantityManufacturedForStock int32 = 0
					for i := 0; i < len(manufacturingOrders); i++ {
						quantityManufacturedForStock += manufacturingOrders[0].QuantityManufactured
					}

					// associate with the existing orders
					if quantityManufacturedForStock >= manufacturingOrderTypeComponent.Quantity {
						var quantityAdded int32 = 0
						// the orders come sorted by date_created ASC, so the ones that are older are first (the ones we expect to manufacture before)
						for i := 0; i < len(manufacturingOrders); i++ {
							c := ComplexManufacturingOrderManufacturingOrder{
								Type:                              "I",
								ComplexManufacturingOrderId:       complexManufacturingOrder.Id,
								EnterpriseId:                      complexManufacturingOrder.EnterpriseId,
								ManufacturingOrderId:              &manufacturingOrders[i].Id,
								ProductId:                         manufacturingOrderTypeComponent.ProductId,
								ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
								Manufactured:                      false,
							}
							subOrders = append(subOrders, c)
							// set the manufacturing order as complex, so it doesn't count as stock
							result := trans.Model(&ManufacturingOrder{}).Where("id = ?", manufacturingOrders[i].Id).Update("complex", true)
							if result.Error != nil {
								log("DB", result.Error.Error())
								trans.Rollback()
								return false, nil
							}
							insertTransactionalLog(c.EnterpriseId, "manufacturing_order", int(manufacturingOrders[i].Id), userId, "U")
							// stop the loop as soon as we get enought quantity
							quantityAdded += manufacturingOrders[i].QuantityManufactured
							if quantityAdded >= manufacturingOrderTypeComponent.Quantity {
								break
							}
						}
					} else { // there are no stock orders, create a new one
						manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderTypeId)
						for i := 0; i < int(manufacturingOrderTypeComponent.Quantity); i += int(manufacturingOrderType.QuantityManufactured) {
							mo := ManufacturingOrder{
								ProductId:    manufacturingOrderTypeComponent.ProductId,
								TypeId:       manufacturingOrderTypeComponent.ManufacturingOrderTypeId,
								EnterpriseId: complexManufacturingOrder.EnterpriseId,
								WarehouseId:  c.WarehouseId,
								Complex:      true,
							}
							mo.insertManufacturingOrder(userId, trans)
							c := ComplexManufacturingOrderManufacturingOrder{
								Type:                              "I",
								ComplexManufacturingOrderId:       complexManufacturingOrder.Id,
								EnterpriseId:                      complexManufacturingOrder.EnterpriseId,
								ManufacturingOrderId:              &mo.Id,
								ProductId:                         manufacturingOrderTypeComponent.ProductId,
								ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
								Manufactured:                      false,
							}
							subOrders = append(subOrders, c)
						}
					}
				} // } else { // if product.ManufacturingOrderType == nil {
			} else { // if product.Manufacturing
				var purchaseDetailId int64 = 0
				// search for a pending purchase order detail
				result := dbOrm.Model(&PurchaseOrderDetail{}).Where("product = ? AND quantity_delivery_note = 0 AND quantity - quantity_assigned_sale >= ?", manufacturingOrderTypeComponent.ProductId, manufacturingOrderTypeComponent.Quantity).Select("id").Pluck("id", &purchaseDetailId)
				if result.Error != nil {
					log("DB", result.Error.Error())
					trans.Rollback()
					return false, nil
				}

				c := ComplexManufacturingOrderManufacturingOrder{
					Type:                              "I",
					ComplexManufacturingOrderId:       complexManufacturingOrder.Id,
					EnterpriseId:                      complexManufacturingOrder.EnterpriseId,
					PurchaseOrderDetailId:             &purchaseDetailId,
					ProductId:                         manufacturingOrderTypeComponent.ProductId,
					ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
					Manufactured:                      false,
				}
				subOrders = append(subOrders, c)

				// there are no pending purchase order details, return error
				if purchaseDetailId == 0 {
					trans.Rollback()
					return false, nil
				} else {
					// add quantity assigned to sale orders
					ok := addQuantityAssignedSalePurchaseOrder(purchaseDetailId, manufacturingOrderTypeComponent.Quantity, complexManufacturingOrder.EnterpriseId, userId, *trans)
					if !ok {
						trans.Rollback()
						return false, nil
					}
				}
			}
		} // if stock.QuantityAvaialbe >= manufacturingOrderTypeComponent.Quantity {

	} // for i := 0; i < len(components); i++ {

	for i := 0; i < len(components); i++ {
		if components[i].Type != "O" { // Only Output
			continue
		}

		manufacturingOrderTypeComponent := components[i]
		if manufacturingOrderTypeComponent.Id <= 0 || manufacturingOrderTypeComponent.EnterpriseId != c.EnterpriseId || manufacturingOrderTypeComponent.Quantity <= 0 {
			trans.Rollback()
			return false, nil
		}

		c := ComplexManufacturingOrderManufacturingOrder{
			Type:                              "O",
			ComplexManufacturingOrderId:       complexManufacturingOrder.Id,
			EnterpriseId:                      complexManufacturingOrder.EnterpriseId,
			ProductId:                         manufacturingOrderTypeComponent.ProductId,
			ManufacturingOrderTypeComponentId: manufacturingOrderTypeComponent.Id,
			Manufactured:                      false,
		}
		subOrders = append(subOrders, c)
	} // for i := 0; i < len(components); i++ {

	for i := 0; i < len(subOrders); i++ {
		ok := subOrders[i].insertComplexManufacturingOrderManufacturingOrder(userId, *trans)
		if !ok {
			trans.Rollback()
			return false, nil
		}
	} // for i := 0; i < len(subOrders); i++ {

	if beginTransaction {
		///
		result := trans.Commit()
		if result.Error != nil {
			return false, nil
		}
		///
	}

	return true, &c.Id
}

func getPendingComplexManufacturingOrderOutputsWithoutSaleOrderDetail(productId int32) *int64 {
	var id int64
	result := dbOrm.Model(&ComplexManufacturingOrder{}).Where("(product = ?) AND (NOT manufactured) AND (type = 'O') AND (sale_order_detail IS NULL)", productId).Order("id ASC").Limit(1).Select("id").Pluck("id", &id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return &id
}

func (c *ComplexManufacturingOrder) deleteComplexManufacturingOrder(userId int32, trans *gorm.DB) bool {
	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	orderInMemory := getComplexManufacturingOrderRow(c.Id)
	if orderInMemory.Id <= 0 || orderInMemory.EnterpriseId != c.EnterpriseId {
		return false
	}

	components := getComplexManufacturingOrderManufacturingOrder(c.Id, c.EnterpriseId)

	for i := 0; i < len(components); i++ {
		ok := components[i].deleteComplexManufacturingOrderManufacturingOrder(userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	result := trans.Where("id = ?", c.Id).Delete(&ComplexManufacturingOrder{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "complex_manufacturing_order", int(c.Id), userId, "D")

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

func toggleManufactuedComplexManufacturingOrder(orderid int64, userId int32, enterpriseId int32) bool {
	if orderid <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	settings := getSettingsRecordById(enterpriseId)

	inMemoryComplexManufacturingOrder := getComplexManufacturingOrderRow(orderid)
	if inMemoryComplexManufacturingOrder.EnterpriseId != enterpriseId {
		trans.Rollback()
		return false
	}

	// validation
	if inMemoryComplexManufacturingOrder.Manufactured && inMemoryComplexManufacturingOrder.DateManufactured != nil && int64(time.Since(*inMemoryComplexManufacturingOrder.DateManufactured).Seconds()) > int64(settings.UndoManufacturingOrderSeconds) {
		trans.Rollback()
		return false
	}
	if !inMemoryComplexManufacturingOrder.Manufactured && inMemoryComplexManufacturingOrder.QuantityManufactured != inMemoryComplexManufacturingOrder.QuantityPendingManufacture {
		trans.Rollback()
		return false
	}

	result := trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("complex_manufacturing_order = ? AND type = 'O'", orderid).Update("manufactured", !inMemoryComplexManufacturingOrder.Manufactured)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "complex_manufacturing_order_manufacturing_order", int(orderid), userId, "U")

	cmomo := getComplexManufacturingOrderManufacturingOrder(orderid, enterpriseId)
	if !inMemoryComplexManufacturingOrder.Manufactured {
		for i := 0; i < len(cmomo); i++ {
			if cmomo[i].Type == "I" {
				continue
			}

			com := getManufacturingOrderTypeComponentRow(cmomo[i].ManufacturingOrderTypeComponentId)

			wm := WarehouseMovement{
				ProductId:    cmomo[i].ProductId,
				WarehouseId:  inMemoryComplexManufacturingOrder.WarehouseId,
				Quantity:     com.Quantity,
				Type:         "O",
				EnterpriseId: enterpriseId,
			}
			wm.insertWarehouseMovement(userId, trans)

			result = trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("id = ?", cmomo[i].Id).Update("warehouse_movement", wm.Id)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(enterpriseId, "complex_manufacturing_order_manufacturing_order", int(cmomo[i].Id), userId, "U")

			if cmomo[i].SaleOrderDetailId != nil {
				sqlStatement := `SELECT COUNT(*) FROM public.complex_manufacturing_order_manufacturing_order WHERE sale_order_detail=$1 AND NOT manufactured`
				row := db.QueryRow(sqlStatement, cmomo[i].SaleOrderDetailId)
				if row.Err() != nil {
					log("DB", row.Err().Error())
					return false
				}

				var ordersPending int32
				row.Scan(&ordersPending)

				if ordersPending == 0 {
					result = trans.Model(&SalesOrderDetail{}).Where("id = ?", cmomo[i].SaleOrderDetailId).Update("status", "E")
					if result.Error != nil {
						log("DB", result.Error.Error())
						trans.Rollback()
						return false
					}
					ok := setSalesOrderState(enterpriseId, cmomo[i].SaleOrderDetail.OrderId, userId, *trans)
					if !ok {
						trans.Rollback()
						return false
					}

					insertTransactionalLog(enterpriseId, "sales_order_detail", int(*cmomo[i].SaleOrderDetailId), userId, "U")
				}
			}

			ok := addQuantityPendingManufacture(cmomo[i].ProductId, inMemoryComplexManufacturingOrder.WarehouseId, -com.Quantity, inMemoryComplexManufacturingOrder.EnterpriseId, *trans)
			if !ok {
				return false
			}
		} // for i := 0; i < len(cmomo); i++ {

		result = trans.Model(&ComplexManufacturingOrder{}).Where("id = ?", inMemoryComplexManufacturingOrder.Id).Updates(map[string]interface{}{
			"manufactured":      true,
			"date_manufactured": time.Now(),
			"user_manufactured": userId,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(enterpriseId, "complex_manufacturing_order", int(inMemoryComplexManufacturingOrder.Id), userId, "U")
	} else { // if !inMemoryComplexManufacturingOrder.Manufactured {
		for i := 0; i < len(cmomo); i++ {
			if cmomo[i].Type == "I" {
				continue
			}

			result = trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("id = ?", cmomo[i].Id).Update("warehouse_movement", nil)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(enterpriseId, "complex_manufacturing_order_manufacturing_order", int(cmomo[i].Id), userId, "U")

			if cmomo[i].WarehouseMovementId != nil {
				wm := getWarehouseMovementRow(*cmomo[i].WarehouseMovementId)
				ok := wm.deleteWarehouseMovement(userId, trans)
				if !ok {
					return false
				}
			}

			if cmomo[i].SaleOrderDetailId != nil {
				result = trans.Model(&SalesOrderDetail{}).Where("id = ?", cmomo[i].SaleOrderDetailId).Update("status", "D")
				if result.Error != nil {
					log("DB", result.Error.Error())
					trans.Rollback()
					return false
				}
				ok := setSalesOrderState(enterpriseId, cmomo[i].SaleOrderDetail.OrderId, userId, *trans)
				if !ok {
					trans.Rollback()
					return false
				}

				insertTransactionalLog(enterpriseId, "sales_order_detail", int(*cmomo[i].SaleOrderDetailId), userId, "U")
			}

			com := getManufacturingOrderTypeComponentRow(cmomo[i].ManufacturingOrderTypeComponentId)
			ok := addQuantityPendingManufacture(cmomo[i].ProductId, inMemoryComplexManufacturingOrder.WarehouseId, com.Quantity, inMemoryComplexManufacturingOrder.EnterpriseId, *trans)
			if !ok {
				return false
			}
		} // for i := 0; i < len(cmomo); i++ {

		result = trans.Model(&ComplexManufacturingOrder{}).Where("id = ?", inMemoryComplexManufacturingOrder.Id).Updates(map[string]interface{}{
			"manufactured":      false,
			"date_manufactured": nil,
			"user_manufactured": nil,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(enterpriseId, "complex_manufacturing_order", int(inMemoryComplexManufacturingOrder.Id), userId, "U")

	} // } else { // if !inMemoryComplexManufacturingOrder.Manufactured {

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

type ComplexManufacturingOrderManufacturingOrder struct {
	Id                                                  int64                                        `json:"id" gorm:"index:complex_manufacturing_order_manufacturing_order_id_enterprise,unique:true,priority:1"`
	ManufacturingOrderId                                *int64                                       `json:"manufacturingOrderId" gorm:"column:manufacturing_order;index:complex_manufacturing_order_complex_manufacturing_order_manufac,unique:true,priority:2"`
	ManufacturingOrder                                  *ManufacturingOrder                          `json:"manufacturingOrder" gorm:"foreignKey:ManufacturingOrderId,EnterpriseId;references:Id,EnterpriseId"`
	Type                                                string                                       `json:"type" gorm:"column:type;type:character(1);not null:true;index:complex_manufacturing_order_complex_manufacturing_order_type,priority:2;index:complex_manufacturing_order_manufacturing_order_pending_product,priority:3,where:NOT manufactured AND type = 'O' AND sale_order_detail IS NULL"` // I = Input, O = Output
	ComplexManufacturingOrderId                         int64                                        `json:"complexManufacturingOrderId" gorm:"column:complex_manufacturing_order;not null:true;index:complex_manufacturing_order_complex_manufacturing_order_manufac,unique:true,priority:1;index:complex_manufacturing_order_complex_manufacturing_order_type,priority:1"`
	ComplexManufacturingOrder                           ComplexManufacturingOrder                    `json:"complexManufacturingOrder" gorm:"foreignKey:ComplexManufacturingOrderId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId                                        int32                                        `json:"-" gorm:"column:enterprise;not null:true;index:complex_manufacturing_order_manufacturing_order_id_enterprise,unique:true,priority:2"`
	Enterprise                                          Settings                                     `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	WarehouseMovementId                                 *int64                                       `json:"warehouseMovementId" gorm:"column:warehouse_movement"`
	WarehouseMovement                                   *WarehouseMovement                           `json:"warehouseMovement" gorm:"foreignKey:WarehouseMovementId,EnterpriseId;references:Id,EnterpriseId"`
	Manufactured                                        bool                                         `json:"manufactured" gorm:"column:manufactured;not null:true;index:complex_manufacturing_order_manufacturing_order_pending_product,priority:2,where:NOT manufactured AND type = 'O' AND sale_order_detail IS NULL"`
	ProductId                                           int32                                        `json:"productId" gorm:"column:product;not null:true;index:complex_manufacturing_order_manufacturing_order_pending_product,priority:1,where:NOT manufactured AND type = 'O' AND sale_order_detail IS NULL"`
	Product                                             Product                                      `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	ManufacturingOrderTypeComponentId                   int32                                        `json:"manufacturingOrderTypeComponentId" gorm:"column:manufacturing_order_type_component;not null:true"`
	ManufacturingOrderTypeComponent                     ManufacturingOrderTypeComponents             `json:"manufacturingOrderTypeComponent" gorm:"foreignKey:ManufacturingOrderTypeComponentId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseOrderDetailId                               *int64                                       `json:"purchaseOrderDetailId" gorm:"column:purchase_order_detail"`
	PurchaseOrderDetail                                 *PurchaseOrderDetail                         `json:"purchaseOrderDetail" gorm:"foreignKey:PurchaseOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	SaleOrderDetailId                                   *int64                                       `json:"saleOrderDetailId" gorm:"column:sale_order_detail;index:complex_manufacturing_order_manufacturing_order_pending_product,priority:5,where:NOT manufactured AND type = 'O' AND sale_order_detail IS NULL"`
	SaleOrderDetail                                     *SalesOrderDetail                            `json:"saleOrderDetail" gorm:"foreignKey:SaleOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	ComplexManufacturingOrderManufacturingOrderOutputId *int64                                       `json:"complexManufacturingOrderManufacturingOrderOutputId" gorm:"column:complex_manufacturing_order_manufacturing_order_output"`
	ComplexManufacturingOrderManufacturingOrderOutput   *ComplexManufacturingOrderManufacturingOrder `json:"complexManufacturingOrderManufacturingOrderOutput" gorm:"foreignKey:ComplexManufacturingOrderManufacturingOrderOutputId;references:Id"`
	SaleOrderName                                       *string                                      `json:"saleOrderName" gorm:"-"`
	PurchaseOrderName                                   *string                                      `json:"purchaseOrderName" gorm:"-"`
}

func (c *ComplexManufacturingOrderManufacturingOrder) TableName() string {
	return "complex_manufacturing_order_manufacturing_order"
}

func getComplexManufacturingOrderManufacturingOrder(complexManufacturingOrderId int64, enterpriseId int32) []ComplexManufacturingOrderManufacturingOrder {
	var orders []ComplexManufacturingOrderManufacturingOrder = make([]ComplexManufacturingOrderManufacturingOrder, 0)
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("complex_manufacturing_order = ? AND enterprise = ?", complexManufacturingOrderId, enterpriseId).Preload(clause.Associations).Preload("SaleOrderDetail.Order").Preload("PurchaseOrderDetail.Order").Find(&orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}

	for i := 0; i < len(orders); i++ {
		if orders[i].SaleOrderDetailId != nil {
			orders[i].SaleOrderName = &orders[i].SaleOrderDetail.Order.OrderName
		}
		if orders[i].PurchaseOrderDetailId != nil {
			orders[i].PurchaseOrderName = &orders[i].PurchaseOrderDetail.Order.OrderName
		}
	}

	return orders
}

func getComplexManufacturingOrderManufacturingOrderRow(complexManufacturingOrderManufacturingOrderId int64) ComplexManufacturingOrderManufacturingOrder {
	c := ComplexManufacturingOrderManufacturingOrder{}
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("id = ?", complexManufacturingOrderManufacturingOrderId).Preload(clause.Associations).First(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ComplexManufacturingOrderManufacturingOrder{}
	}
	return c
}

func (c *ComplexManufacturingOrderManufacturingOrder) isValid() bool {
	return !(c.ProductId <= 0 || (c.Type != "I" && c.Type != "O") || c.ComplexManufacturingOrderId <= 0 || c.ManufacturingOrderTypeComponentId <= 0)
}

func (c *ComplexManufacturingOrderManufacturingOrder) BeforeCreate(tx *gorm.DB) (err error) {
	var complexManufacturingOrderManufacturingOrder ComplexManufacturingOrderManufacturingOrder
	tx.Model(&ComplexManufacturingOrderManufacturingOrder{}).Last(&complexManufacturingOrderManufacturingOrder)
	c.Id = complexManufacturingOrderManufacturingOrder.Id + 1
	return nil
}

// DOES NOT OPEN A TRANSACTION
func (c *ComplexManufacturingOrderManufacturingOrder) insertComplexManufacturingOrderManufacturingOrder(userId int32, trans gorm.DB) bool {
	if !c.isValid() {
		return false
	}

	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "complex_manufacturing_order_manufacturing_order", int(c.Id), userId, "I")

	ok := addQuantityPendingManufactureComplexManufacturingOrder(c.ComplexManufacturingOrderId, 1, c.EnterpriseId, userId, trans)
	if ok && c.WarehouseMovementId != nil {
		return addQuantityManufacturedComplexManufacturingOrder(c.ComplexManufacturingOrderId, 1, c.EnterpriseId, userId, trans)
	}
	if ok {
		order := getComplexManufacturingOrderRowTransaction(c.ComplexManufacturingOrderId, trans)
		com := getManufacturingOrderTypeComponentRowTransaction(c.ManufacturingOrderTypeComponentId, trans)
		return addQuantityPendingManufacture(c.ProductId, order.WarehouseId, com.Quantity, c.EnterpriseId, trans)
	}
	return ok
}

// DOES NOT OPEN A TRANSACTION
func (c *ComplexManufacturingOrderManufacturingOrder) deleteComplexManufacturingOrderManufacturingOrder(userId int32, trans gorm.DB) bool {
	if c.Id <= 0 {
		return false
	}

	comInMemory := getComplexManufacturingOrderManufacturingOrderRow(c.Id)
	if comInMemory.Id <= 0 || comInMemory.EnterpriseId != c.EnterpriseId {
		return false
	}

	result := trans.Delete(&ComplexManufacturingOrderManufacturingOrder{}, "id = ?", c.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "complex_manufacturing_order_manufacturing_order", int(c.Id), userId, "D")

	ok := addQuantityPendingManufactureComplexManufacturingOrder(c.ComplexManufacturingOrderId, -1, c.EnterpriseId, userId, trans)
	if !ok {
		return false
	}

	if comInMemory.ManufacturingOrderId != nil {
		mo := getManufacturingOrderRow(*comInMemory.ManufacturingOrderId)
		ok := mo.deleteManufacturingOrder(userId, &trans)
		if !ok {
			return false
		}
	}

	if comInMemory.WarehouseMovementId != nil {
		wm := getWarehouseMovementRow(*comInMemory.WarehouseMovementId)
		ok := wm.deleteWarehouseMovement(userId, &trans)
		if !ok {
			return false
		}
	}

	if comInMemory.PurchaseOrderDetailId != nil {
		component := getManufacturingOrderTypeComponentRow(comInMemory.ManufacturingOrderTypeComponentId)
		ok := addQuantityAssignedSalePurchaseOrder(*comInMemory.PurchaseOrderDetailId, component.Quantity, comInMemory.EnterpriseId, userId, trans)
		if !ok {
			return false
		}
	}

	if comInMemory.SaleOrderDetailId != nil {
		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", comInMemory.SaleOrderDetailId).Update("status", "C")
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(c.EnterpriseId, "sales_order_detail", int(*comInMemory.SaleOrderDetailId), userId, "U")

		ok := setSalesOrderState(c.EnterpriseId, comInMemory.SaleOrderDetail.OrderId, userId, trans)
		if !ok {
			return false
		}
	}

	return true
}

func setComplexManufacturingOrderManufacturingOrderManufactured(manufacturingOrderId int64, manufactured bool, enterpriseId int32, userId int32, trans *gorm.DB) bool {
	var id int64
	var complexManufacturingOrderId int64
	var orderManufactured bool
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("manufacturing_order = ? AND type = 'O'", manufacturingOrderId).Pluck("id", &id).Pluck("complex_manufacturing_order", &complexManufacturingOrderId).Pluck("order_manufactured", &orderManufactured)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if complexManufacturingOrderId <= 0 || manufactured == orderManufactured {
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

	// update the sub-order
	result = dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("manufacturing_order = ? AND type = 'O'", manufacturingOrderId).Update("manufactured", manufactured)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "complex_manufacturing_order_manufacturing_order", int(manufacturingOrderId), userId, "U")

	// update the quantities
	if !orderManufactured == manufactured {
		ok := addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId, 1, enterpriseId, userId, *trans)
		if !ok {
			return false
		}
	} else if orderManufactured == !manufactured {
		ok := addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId, -1, enterpriseId, userId, *trans)
		if !ok {
			return false
		}
	}

	// recursivity
	cmomo := getComplexManufacturingOrderManufacturingOrderRow(id)
	if cmomo.ComplexManufacturingOrderManufacturingOrderOutputId != nil {
		ok := setComplexManufacturingOrderManufacturingOrderManufactured(manufacturingOrderId, manufactured, enterpriseId, userId, trans)
		if !ok {
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

func complexManufacturingOrderTagPrinted(orderId int64, userId int32, enterpriseId int32) bool {
	if orderId <= 0 {
		return false
	}

	result := dbOrm.Model(&ComplexManufacturingOrder{}).Where("id = ?", orderId).Updates(map[string]interface{}{
		"tag_printed":      true,
		"date_tag_printed": time.Now(),
		"user_tag_printed": userId,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "complex_manufacturing_order", int(orderId), userId, "U")

	return true
}

type MultipleComplexManufacturingOrders struct {
	Order   ComplexManufacturingOrder `json:"order"`
	Quantiy int                       `json:"quantity"`
}

func (o *MultipleComplexManufacturingOrders) insertMultipleComplexManufacturingOrders(userId int32) bool {
	if !o.Order.isValid() || o.Quantiy <= 0 || o.Quantiy > 10000 {
		return false
	}

	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}

	for i := 0; i < o.Quantiy; i++ {
		ok, _ := o.Order.insertComplexManufacturingOrder(userId, trans)
		if !ok {
			trans.Rollback()
			return ok
		}
	}

	trans.Commit()
	return true
}
