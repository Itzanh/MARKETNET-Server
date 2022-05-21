package main

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseOrderDetail struct {
	Id                   int64         `json:"id" gorm:"index:purchase_order_detail_id_enterprise,unique:true,priority:1"`
	OrderId              int64         `json:"orderId" gorm:"column:order;not null:true;index:purchase_order_detail_purchase_order_product,unique:true,priority:1"`
	Order                PurchaseOrder `json:"-" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId            int32         `json:"productId" gorm:"column:product;not null:true;index:purchase_order_detail_purchase_order_product,unique:true,priority:2"`
	Product              Product       `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Price                float64       `json:"price" gorm:"column:price;not null:true;type:numeric(14,6)"`
	Quantity             int32         `json:"quantity" gorm:"column:quantity;not null:true"`
	VatPercent           float64       `json:"vatPercent" gorm:"column:vat_percent;not null:true;type:numeric(14,6)"`
	TotalAmount          float64       `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	QuantityInvoiced     int32         `json:"quantityInvoiced" gorm:"column:quantity_invoiced;not null:true"`
	QuantityDeliveryNote int32         `json:"quantityDeliveryNote" gorm:"column:quantity_delivery_note;not null:true"`
	QuantityAssignedSale int32         `json:"quantityAssignedSale" gorm:"column:quantity_assigned_sale;not null:true"`
	EnterpriseId         int32         `json:"-" gorm:"column:enterprise;not null:true;index:purchase_order_detail_id_enterprise,unique:true,priority:2"`
	Enterprise           Settings      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Cancelled            bool          `json:"cancelled" gorm:"column:cancelled;not null:true"`
}

func (pod *PurchaseOrderDetail) TableName() string {
	return "purchase_order_detail"
}

func getPurchaseOrderDetail(orderId int64, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	result := dbOrm.Where("purchase_order_detail.order = ? AND purchase_order_detail.enterprise = ?", orderId, enterpriseId).Order("purchase_order_detail.id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return details
}

func getPurchaseOrderDetailRow(detailId int64) PurchaseOrderDetail {
	var d PurchaseOrderDetail
	result := dbOrm.Where("purchase_order_detail.id = ?", detailId).Preload(clause.Associations).First(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return PurchaseOrderDetail{}
	}
	return d
}

func getPurchaseOrderDetailRowTransaction(detailId int64, trans gorm.DB) PurchaseOrderDetail {
	var d PurchaseOrderDetail
	result := trans.Where("purchase_order_detail.id = ?", detailId).Preload(clause.Associations).First(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return PurchaseOrderDetail{}
	}
	return d
}

func (d *PurchaseOrderDetail) isValid() bool {
	return !(d.OrderId <= 0 || d.ProductId <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (d *PurchaseOrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var purchaseOrderDetail PurchaseOrderDetail
	tx.Model(&PurchaseOrderDetail{}).Last(&purchaseOrderDetail)
	d.Id = purchaseOrderDetail.Id + 1
	return nil
}

// 1. the product is deactivated
// 2. there is aleady a detail with this product
func (s *PurchaseOrderDetail) insertPurchaseOrderDetail(userId int32, trans *gorm.DB) (OkAndErrorCodeReturn, int64) {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	product := getProductRow(s.ProductId)
	if product.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if product.Off {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}, 0
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	var countProductInSaleOrder int64
	result := dbOrm.Model(&PurchaseOrderDetail{}).Where("\"order\" = ? AND product = ?", s.OrderId, s.ProductId).Count(&countProductInSaleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if countProductInSaleOrder > 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}, 0
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	///
	var beginTrans bool = (trans == nil)
	if beginTrans {
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
	}
	///

	s.QuantityInvoiced = 0
	s.QuantityDeliveryNote = 0
	s.QuantityAssignedSale = 0
	s.Cancelled = false

	result = trans.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	insertTransactionalLog(s.EnterpriseId, "purchase_order_detail", int(s.Id), userId, "I")
	jsn, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "purchase_order_detail", "POST", string(jsn))

	ok := addTotalProductsPurchaseOrder(s.OrderId, s.Price*float64(s.Quantity), s.VatPercent, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	ok = addPurchaseOrderLinesNumber(s.OrderId, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	quantityAssignedSale := associatePurchaseOrderWithPendingSalesOrders(s.Id, s.ProductId, s.Quantity, s.EnterpriseId, userId, *trans)
	if quantityAssignedSale < 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	result = trans.Model(&PurchaseOrderDetail{}).Where("id = ?", s.Id).Update("quantity_assigned_sale", quantityAssignedSale)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	insertTransactionalLog(s.EnterpriseId, "purchase_order_detail", int(s.Id), userId, "U")
	jsn, _ = json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "purchase_order_detail", "POST", string(jsn))

	// add quantity pending receiving
	purchaseOrder := getPurchaseOrderRow(s.OrderId)
	ok = addQuantityPendingReveiving(s.ProductId, purchaseOrder.WarehouseId, s.Quantity, s.EnterpriseId, *trans)
	if !ok {
		if beginTrans {
			trans.Rollback()
		}
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	if beginTrans {
		///
		result = trans.Commit()
		return OkAndErrorCodeReturn{Ok: result.Error == nil}, s.Id
		///
	} else {
		return OkAndErrorCodeReturn{Ok: true}, s.Id
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func deassociatePurchaseOrderWithPendingSalesOrders(purchaseDetailId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	result := trans.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ?", purchaseDetailId).Updates(map[string]interface{}{
		"status":                "A",
		"purchase_order_detail": nil,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	return true
}

type AssociatePurchaseOrderWithPendingSalesOrders struct {
	SalesDetailId int64
	SalesQuantity int32
	OrderId       int64
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func associatePurchaseOrderWithPendingSalesOrders(purchaseDetailId int64, productId int32, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) int32 {
	// associate pending sales order detail until there are no more quantity pending to be assigned, or there are no more pending sales order details
	var salesOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	result := trans.Model(&SalesOrderDetail{}).Where("product = ? AND status = 'A'", productId).Order("(SELECT date_created FROM sales_order WHERE sales_order.id=sales_order_detail.\"order\") ASC").Find(&salesOrderDetails)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return -1
	}
	var associations []AssociatePurchaseOrderWithPendingSalesOrders = make([]AssociatePurchaseOrderWithPendingSalesOrders, 0)

	var i int
	var quantityAssignedSale int32
	for quantityAssignedSale < quantity {
		if i < len(salesOrderDetails) {
			saleDetail := salesOrderDetails[i]
			i++

			if quantityAssignedSale+saleDetail.Quantity > quantity { // no more rows to proecss
				break
			}

			quantityAssignedSale += saleDetail.Quantity
			associations = append(associations, AssociatePurchaseOrderWithPendingSalesOrders{
				SalesDetailId: saleDetail.Id,
				SalesQuantity: saleDetail.Quantity,
				OrderId:       saleDetail.OrderId,
			})
		} else { // no more rows to process
			break
		}
	}

	for i := 0; i < len(associations); i++ {
		association := associations[i]
		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", association.SalesDetailId).Updates(map[string]interface{}{
			"status":                "B",
			"purchase_order_detail": purchaseDetailId,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return -1
		}
		setSalesOrderState(enterpriseId, association.OrderId, userId, trans)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(association.SalesDetailId), userId, "U")
	}
	return quantityAssignedSale
}

// ERROR CODES:
// 1. the detail is already invoiced
// 2. the detail has a delivery note generated
func (s *PurchaseOrderDetail) updatePurchaseOrderDetail(userId int32) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	detailInMemory := getPurchaseOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.EnterpriseId != s.EnterpriseId {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.QuantityInvoiced > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}
	if detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	if detailInMemory.Quantity != s.Quantity {
		if !deassociatePurchaseOrderWithPendingSalesOrders(s.Id, s.EnterpriseId, userId, *trans) {
			return OkAndErrorCodeReturn{Ok: false}
		}
		s.QuantityAssignedSale = associatePurchaseOrderWithPendingSalesOrders(s.Id, s.ProductId, s.Quantity, s.EnterpriseId, userId, *trans)
		if s.QuantityAssignedSale < 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	} else {
		s.QuantityAssignedSale = detailInMemory.QuantityAssignedSale
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	ok := addTotalProductsPurchaseOrder(s.OrderId, -detailInMemory.Price*float64(detailInMemory.Quantity), detailInMemory.VatPercent, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok = addTotalProductsPurchaseOrder(s.OrderId, s.Price*float64(s.Quantity), s.VatPercent, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	detailInMemory.Price = s.Price
	detailInMemory.Quantity = s.Quantity
	detailInMemory.VatPercent = s.VatPercent
	detailInMemory.TotalAmount = s.TotalAmount
	detailInMemory.QuantityAssignedSale = s.QuantityAssignedSale

	result := trans.Model(&PurchaseOrderDetail{}).Where("id = ?", s.Id).Updates(map[string]interface{}{
		"price":                  detailInMemory.Price,
		"quantity":               detailInMemory.Quantity,
		"vat_percent":            detailInMemory.VatPercent,
		"total_amount":           detailInMemory.TotalAmount,
		"quantity_assigned_sale": detailInMemory.QuantityAssignedSale,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a transaction.
//
// ERROR CODES:
// 1. the detail is already invoiced
// 2. the detail has a delivery note generated
func (s *PurchaseOrderDetail) deletePurchaseOrderDetail(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	detailInMemory := getPurchaseOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.EnterpriseId != s.EnterpriseId {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.QuantityInvoiced > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}
	if detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	// roll back the state of the sale order details
	details := getSalesOrderDetailPurchaseOrderPending(s.Id)
	for i := 0; i < len(details); i++ {
		result := trans.Model(&SalesOrderDetail{}).Where("id = ?", details[i].Id).Updates(map[string]interface{}{
			"status":                "A",
			"purchase_order_detail": nil,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
		ok := setSalesOrderState(detailInMemory.EnterpriseId, details[i].OrderId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(detailInMemory.EnterpriseId, "purchase_order_detail", int(s.Id), userId, "D")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "purchase_order_detail", "DELETE", string(json))

	result := trans.Delete(&PurchaseOrderDetail{}, "id = ?", s.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsPurchaseOrder(detailInMemory.OrderId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, detailInMemory.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = removePurchaseOrderLinesNumber(detailInMemory.OrderId, detailInMemory.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	// substract quantity pending receiving
	purchaseOrder := getPurchaseOrderRow(detailInMemory.OrderId)
	ok = addQuantityPendingReveiving(detailInMemory.ProductId, purchaseOrder.WarehouseId, -detailInMemory.Quantity, s.EnterpriseId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		trans.Commit()
		///
	}
	return OkAndErrorCodeReturn{Ok: true}
}

// Adds quantity to the field to prevent from other sale orders to use the quantity that is already reserved for order that are already waiting a purchase order.
// This function will substract if a negative quantity is given.
// THIS FUNCION DOES NOT OPEN A TRANSACTION
func addQuantityAssignedSalePurchaseOrder(detailId int64, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrderDetail PurchaseOrderDetail
	result := trans.Model(&PurchaseOrderDetail{}).Where("id = ?", detailId).First(&purchaseOrderDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrderDetail.QuantityAssignedSale += quantity

	result = trans.Save(&purchaseOrderDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	json, _ := json.Marshal(purchaseOrderDetail)
	go fireWebHook(enterpriseId, "purchase_order_detail", "PUT", string(json))

	return true
}

// Adds an invoiced quantity to the purchase order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvoicedPurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	var detailAfter PurchaseOrderDetail
	result := trans.Model(&PurchaseOrderDetail{}).Where("id = ?", detailId).First(&detailAfter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	detailAfter.QuantityInvoiced += quantity

	result = trans.Save(&detailAfter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	json, _ := json.Marshal(detailAfter)
	go fireWebHook(enterpriseId, "purchase_order_detail", "PUT", string(json))

	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity { // set as invoced
		ok := addPurchaseOrderInvoicedLines(detailBefore.OrderId, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok := removePurchaseOrderInvoicedLines(detailBefore.OrderId, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	}

	return true
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNotePurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32, trans gorm.DB) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	var detailAfter PurchaseOrderDetail
	result := trans.Model(&PurchaseOrderDetail{}).Where("id = ?", detailId).First(&detailAfter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	detailAfter.QuantityDeliveryNote += quantity

	result = trans.Save(&detailAfter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	json, _ := json.Marshal(detailAfter)
	go fireWebHook(enterpriseId, "purchase_order_detail", "PUT", string(json))

	if detailBefore.QuantityDeliveryNote != detailBefore.Quantity && detailAfter.QuantityDeliveryNote == detailAfter.Quantity { // set as delivery note generated
		ok := addPurchaseOrderDeliveryNoteLines(detailBefore.OrderId, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok := removePurchaseOrderDeliveryNoteLines(detailBefore.OrderId, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	}

	if quantity > 0 { // the purchase order has been added to a delivery note, advance the status from the pending sales order details
		return setSalesOrderDetailStateAllPendingPurchaseOrder(detailId, enterpriseId, userId, trans) && setComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId, enterpriseId, userId, trans)
	} else { // the delivery note details has been removed, roll back the sales order detail status
		return undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId, enterpriseId, userId, trans) && undoComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId, enterpriseId, userId, trans)
	}
}

func cancelPurchaseOrderDetail(detailId int64, enterpriseId int32, userId int32) bool {
	detail := getPurchaseOrderDetailRow(detailId)
	if detail.Id <= 0 || detail.EnterpriseId != enterpriseId {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	if !detail.Cancelled {
		if detail.Quantity <= 0 || detail.QuantityInvoiced < 0 || detail.QuantityDeliveryNote > 0 {
			return false
		}

		detail.QuantityInvoiced = detail.Quantity
		detail.QuantityDeliveryNote = detail.Quantity
		detail.Cancelled = true

		result := trans.Save(&detail)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(detail.EnterpriseId, "purchase_order_detail", int(detailId), userId, "U")
		json, _ := json.Marshal(detail)
		go fireWebHook(enterpriseId, "purchase_order_detail", "PUT", string(json))

		///
		result = trans.Commit()
		return result.Error == nil
		///
	} else {
		if detail.Quantity <= 0 || detail.QuantityInvoiced == 0 || detail.QuantityDeliveryNote == 0 {
			return false
		}

		detail.QuantityInvoiced = 0
		detail.QuantityDeliveryNote = 0
		detail.Cancelled = false
		detail.QuantityAssignedSale = 0

		result := trans.Save(&detail)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		salesDetails := getSalesOrderDetailsFromPurchaseOrderDetail(detail.Id, detail.EnterpriseId)

		for i := 0; i < len(salesDetails); i++ {
			result := trans.Model(&SalesOrderDetail{}).Where("id = ?", salesDetails[i].Id).Updates(map[string]interface{}{
				"status":                "A",
				"purchase_order_detail": nil,
			})
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}
		}

		insertTransactionalLog(detail.EnterpriseId, "purchase_order_detail", int(detailId), userId, "U")
		json, _ := json.Marshal(detail)
		go fireWebHook(enterpriseId, "purchase_order_detail", "PUT", string(json))

		///
		result = trans.Commit()
		return result.Error == nil
		///
	}
}

// All the purchase order detail has been added to a delivery note. Advance the status from all the pending sales details to "Sent to preparation".
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setSalesOrderDetailStateAllPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	purchaseOrderDetail := getPurchaseOrderDetailRowTransaction(detailId, trans)

	// Get the quantity that the orders are currently using
	var quantityUsedDeliveryNote *int32
	result := dbOrm.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ? AND status != 'B'", detailId).Select("SUM(quantity) AS quantity").Scan(&quantityUsedDeliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if quantityUsedDeliveryNote == nil {
		zero := int32(0)
		quantityUsedDeliveryNote = &zero
	}

	// Get the quantity that the orders are not currently using + the added quantity
	quantityAddedToDeliveryNote := purchaseOrderDetail.QuantityDeliveryNote - *quantityUsedDeliveryNote
	rows, err := dbOrm.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ? AND status = 'B'", detailId).Select(`id,"order",quantity`).Order("quantity ASC").Rows()
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	defer rows.Close()

	var quantityUsed int32

	var salesOrderDetailId int64
	var saleOrderId int64
	var quantity int32
	for rows.Next() {
		rows.Scan(&salesOrderDetailId, &saleOrderId, &quantity)

		if quantityUsed+quantity > quantityAddedToDeliveryNote {
			return true
		}

		quantityUsed += quantity

		result := trans.Model(&SalesOrderDetail{}).Where("id = ?", salesOrderDetailId).Update("status", "E")
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		setSalesOrderState(enterpriseId, saleOrderId, userId, trans)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(salesOrderDetailId), userId, "U")
		s := getSalesOrderDetailRowTransaction(salesOrderDetailId, trans)
		json, _ := json.Marshal(s)

		if quantityUsed >= quantityAddedToDeliveryNote {
			return true
		}

		go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))
	}

	return true
}

// The purchase order detail was added to a delivery note and it advanced the status from the sales details, but the delivery note was deleted.
// Roll back the status change.
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	purchaseOrderDetail := getPurchaseOrderDetailRowTransaction(detailId, trans)

	var quantityUsedDeliveryNote *int32
	result := dbOrm.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ? AND status = 'E'", detailId).Select("SUM(quantity) AS quantity").Scan(&quantityUsedDeliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	if quantityUsedDeliveryNote == nil {
		zero := int32(0)
		quantityUsedDeliveryNote = &zero
	}

	quantityToRemoveFromDeliveryNote := *quantityUsedDeliveryNote - purchaseOrderDetail.QuantityDeliveryNote
	// The sale orders are using less quantity that the one remaining in the purchase delivery note, do nothing.
	if quantityToRemoveFromDeliveryNote <= 0 {
		return true
	}

	var quantityDeleted int32

	rows, err := dbOrm.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ? AND status = 'E'", detailId).Select(`id,"order",quantity`).Order("quantity DESC").Rows()
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	defer rows.Close()

	var salesOrderDetailId int64
	var saleOrderId int64
	var quantity int32
	for rows.Next() {
		rows.Scan(&salesOrderDetailId, &saleOrderId, &quantity)

		result := trans.Model(&SalesOrderDetail{}).Where("id = ?", salesOrderDetailId).Update("status", "B")
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		setSalesOrderState(enterpriseId, saleOrderId, userId, trans)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(salesOrderDetailId), userId, "U")
		s := getSalesOrderDetailRowTransaction(salesOrderDetailId, trans)
		json, _ := json.Marshal(s)
		go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))

		quantityDeleted += quantity

		if quantityDeleted >= quantityToRemoveFromDeliveryNote {
			return true
		}
	}

	return err == nil
}

// Gets all the sub orders from complex manufacturing orders that are waiting for this pending purchase order,
// and sets them to manufactured (and the parent order as manufacturable).
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	rows, err := trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("purchase_order_detail = ?", detailId).Distinct("complex_manufacturing_order").Select("complex_manufacturing_order").Rows()
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	defer rows.Close()

	var complexManufacturingOrders []int64 = make([]int64, 0)

	for rows.Next() {
		var complexManufacturingOrder int64
		rows.Scan(&complexManufacturingOrder)
	}

	for i := 0; i < len(complexManufacturingOrders); i++ {
		ok := setComplexManufacturingOrderManufacturingOrderManufactured(complexManufacturingOrders[i], true, enterpriseId, userId, &trans)
		if !ok {
			return false
		}
	}

	return true
}

// The purchase order detail was added to a delivery note and it advanced the status from the complex manufacturing orders, but the delivery note was deleted.
// Roll back the status change.
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func undoComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	rows, err := trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("purchase_order_detail = ?", detailId).Distinct("complex_manufacturing_order").Select("complex_manufacturing_order").Rows()
	if err != nil {
		log("DB", err.Error())
		return false
	}
	defer rows.Close()

	var complexManufacturingOrders []int64 = make([]int64, 0)

	for rows.Next() {
		var complexManufacturingOrder int64
		rows.Scan(&complexManufacturingOrder)
	}

	for i := 0; i < len(complexManufacturingOrders); i++ {
		ok := setComplexManufacturingOrderManufacturingOrderManufactured(complexManufacturingOrders[i], false, enterpriseId, userId, &trans)
		if !ok {
			return false
		}
	}

	return true
}

type PurchaseSalesOrderDetail struct {
	Id           int32     `json:"id"`
	Order        int64     `json:"order"`
	OrderName    string    `json:"orderName"`
	DateCreated  time.Time `json:"dateCreated"`
	Customer     int32     `json:"customer"`
	CustomerName string    `json:"customerName"`
	Quantity     int32     `json:"quantity"`
	TotalAmount  float64   `json:"totalAmount"`
}

func getSalesOrderDetailsFromPurchaseOrderDetail(detailId int64, enterpriseId int32) []PurchaseSalesOrderDetail {
	purchaseSalesOrderDetail := make([]PurchaseSalesOrderDetail, 0)
	result := dbOrm.Model(&SalesOrderDetail{}).Where("purchase_order_detail = ? AND enterprise = ?", detailId, enterpriseId).Order("id DESC").Scan(&purchaseSalesOrderDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	for i := 0; i < len(purchaseSalesOrderDetail); i++ {
		order := getSalesOrderRow(purchaseSalesOrderDetail[i].Order)
		purchaseSalesOrderDetail[i].OrderName = order.OrderName
		purchaseSalesOrderDetail[i].DateCreated = order.DateCreated
		purchaseSalesOrderDetail[i].CustomerName = order.Customer.Name
	}
	return purchaseSalesOrderDetail
}

type PurchaseComplexManufacturingOrder struct {
	Id           int64     `json:"id"`
	Type         int32     `json:"type"`
	Manufactured bool      `json:"manufactured"`
	DateCreated  time.Time `json:"dateCreated"`
	TypeName     string    `json:"typeName"`
}

func getComplexManufacturingOrdersFromPurchaseOrderDetail(detailId int64, enterpriseId int32) []PurchaseComplexManufacturingOrder {
	purchaseComplexManufacturingOrder := make([]PurchaseComplexManufacturingOrder, 0)
	var complexManufacturingOrderDetails []ComplexManufacturingOrderManufacturingOrder = make([]ComplexManufacturingOrderManufacturingOrder, 0)
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("purchase_order_detail = ?", detailId).Select("complex_manufacturing_order").Distinct("complex_manufacturing_order").Scan(&complexManufacturingOrderDetails)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}

	for i := 0; i < len(complexManufacturingOrderDetails); i++ {
		var complexManufacturingOrder PurchaseComplexManufacturingOrder
		result := dbOrm.Model(&ComplexManufacturingOrder{}).Where("id = ?", complexManufacturingOrderDetails[i].ComplexManufacturingOrderId).Scan(&complexManufacturingOrder)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return nil
		}
		complexManufacturingOrder.TypeName = getManufacturingOrderTypeRow(complexManufacturingOrder.Type).Name
		purchaseComplexManufacturingOrder = append(purchaseComplexManufacturingOrder, complexManufacturingOrder)
	}

	return purchaseComplexManufacturingOrder
}

func filterPurchaseOrderDetails(input []PurchaseOrderDetail, test func(PurchaseOrderDetail) bool) (output []PurchaseOrderDetail) {
	for _, s := range input {
		if test(s) {
			output = append(output, s)
		}
	}
	return
}
