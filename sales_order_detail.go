package main

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesOrderDetail struct {
	Id                       int64                `json:"id" gorm:"index:sales_order_detail_id_enterprise,unique:true,priority:1"`
	OrderId                  int64                `json:"orderId" gorm:"column:order;not null:true;index:sales_order_detail_sales_order_product,unique:true,priority:1"`
	Order                    SaleOrder            `json:"-" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseId              string               `json:"warehouseId" gorm:"column:warehouse;type:character(2)"`
	Warehouse                Warehouse            `json:"warehouse" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId                int32                `json:"productId" gorm:"column:product;not null:true;index:sales_order_detail_sales_order_product,unique:true,priority:2"`
	Product                  Product              `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Price                    float64              `json:"price" gorm:"column:price;not null:true;type:numeric(14,6)"`
	Quantity                 int32                `json:"quantity" gorm:"column:quantity;not null:true"`
	VatPercent               float64              `json:"vatPercent" gorm:"column:vat_percent;not null:true;type:numeric(14,6)"`
	TotalAmount              float64              `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	QuantityInvoiced         int32                `json:"quantityInvoiced" gorm:"column:quantity_invoiced;not null:true"`
	QuantityDeliveryNote     int32                `json:"quantityDeliveryNote" gorm:"column:quantity_delivery_note;not null:true"`
	Status                   string               `json:"status" gorm:"type:character(1);not null:true"` // _ = Waiting for payment, A = Waiting for purchase order, B = Purchase order pending, C = Waiting for manufacturing orders, D = Manufacturing orders pending, E = Sent to preparation, F = Awaiting for shipping, G = Shipped, H = Receiced by the customer, Z = Cancelled
	QuantityPendingPackaging int32                `json:"quantityPendingPackaging" gorm:"column:quantity_pending_packaging;not null:true"`
	PurchaseOrderDetailId    *int64               `json:"purchaseOrderDetailId" gorm:"column:purchase_order_detail"`
	PurchaseOrderDetail      *PurchaseOrderDetail `json:"-" gorm:"foreignKey:PurchaseOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopId             int32                `json:"-" gorm:"column:ps_id;not null:true;index:sales_order_detail_ps_id,unique:true,priority:2,where:ps_id <> 0"`
	Cancelled                bool                 `json:"cancelled" gorm:"column:cancelled;not null:true"`
	WooCommerceId            int32                `json:"-" gorm:"column:wc_id;not null:true;index:sales_order_detail_wc_id,unique:true,priority:2,where:wc_id <> 0"`
	ShopifyId                int64                `json:"-" gorm:"column:sy_id;not null:true;index:sales_order_detail_sy_id,unique:true,priority:2,where:sy_id <> 0"`
	ShopifyDraftId           int64                `json:"-" gorm:"column:sy_draft_id;not null:true;index:sales_order_detail_sy_draft_id,unique:true,priority:2,where:sy_draft_id <> 0"`
	EnterpriseId             int32                `json:"-" gorm:"column:enterprise;not null:true;index:sales_order_detail_id_enterprise,unique:true,priority:2;index:sales_order_detail_ps_id,unique:true,priority:1,where:ps_id <> 0;index:sales_order_detail_sy_draft_id,unique:true,priority:1,where:sy_draft_id <> 0;;index:sales_order_detail_sy_id,unique:true,priority:1,where:sy_id <> 0;index:sales_order_detail_wc_id,unique:true,priority:1,where:wc_id <> 0"`
	Enterprise               Settings             `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *SalesOrderDetail) TableName() string {
	return "sales_order_detail"
}

func getSalesOrderDetail(orderId int64, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	// get all the sale order details from the database where the order id is the same as the one passed and the enterprise id is the same as the one passed order by id using dbOrm
	result := dbOrm.Where("\"order\" = ? AND enterprise = ?", orderId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return details
	}

	return details
}

type SalesOrderDetailWaitingForManufacturingOrders struct {
	SalesOrderDetail
	OrderName    string `json:"orderName"`
	CustomerName string `json:"customerName"`
}

func getSalesOrderDetailWaitingForManufacturingOrders(enterpriseId int32) []SalesOrderDetailWaitingForManufacturingOrders {
	var details []SalesOrderDetailWaitingForManufacturingOrders = make([]SalesOrderDetailWaitingForManufacturingOrders, 0)
	// get all the sale order details from the database where the order id is the same as the one passed and the enterprise id is the same as the one passed order by id using dbOrm
	result := dbOrm.Model(&SalesOrderDetailWaitingForManufacturingOrders{}).Where("enterprise = ? AND status = 'C'", enterpriseId).Order("ASC").Preload(clause.Associations).Limit(300).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return details
}

func getSalesOrderDetailRow(detailId int64) SalesOrderDetail {
	var detail SalesOrderDetail = SalesOrderDetail{}
	// get a single sale order detail from the database where the id is the same as the one passed using dbOrm
	result := dbOrm.Where("id = ?", detailId).Preload(clause.Associations).First(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return detail
	}

	return detail
}

func getSalesOrderDetailRowTransaction(detailId int64, trans gorm.DB) SalesOrderDetail {
	// get a single sale order detail from the database where the id is the same as the one passed using dbOrm
	var detail SalesOrderDetail = SalesOrderDetail{}
	result := trans.Where("id = ?", detailId).Preload(clause.Associations).First(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return detail
	}

	return detail
}

// Used for purchases
func getSalesOrderDetailWaitingForPurchaseOrder(productId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	// get all the sale order details from the database where the order id is the same as the one passed and the enterprise id is the same as the one passed order by id using dbOrm
	result := dbOrm.Where("product = ? AND status = 'A'", productId).Order("id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return details
	}

	return details
}

// Used for purchases
func getSalesOrderDetailPurchaseOrderPending(purchaseOrderDetail int64) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	// get all the sale order details from the database where the order id is the same as the one passed and the enterprise id is the same as the one passed order by id using dbOrm
	result := dbOrm.Where("purchase_order_detail = ? AND status = 'B'", purchaseOrderDetail).Order("id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return details
	}

	return details
}

func (s *SalesOrderDetail) isValid() bool {
	return !(s.OrderId <= 0 || s.ProductId <= 0 || s.Quantity <= 0 || s.VatPercent < 0)
}

func (s *SalesOrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var salesOrderDetail SalesOrderDetail
	tx.Model(&SalesOrderDetail{}).Last(&salesOrderDetail)
	s.Id = salesOrderDetail.Id + 1
	return nil
}

// 1. the product is deactivated
// 2. there is aleady a detail with this product
func (s *SalesOrderDetail) insertSalesOrderDetail(userId int32) OkAndErrorCodeReturn {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	p := getProductRow(s.ProductId)
	if p.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if p.Off {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}
	config := getSettingsRecordById(s.EnterpriseId)

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))
	s.Status = "_"

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	var countProductInSaleOrder int64
	result := dbOrm.Model(&SalesOrderDetail{}).Where("product = ? AND \"order\" = ?", s.ProductId, s.OrderId).Count(&countProductInSaleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if countProductInSaleOrder > 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	s.QuantityInvoiced = 0
	s.QuantityDeliveryNote = 0
	s.QuantityPendingPackaging = s.Quantity
	s.PurchaseOrderDetail = nil
	s.Cancelled = false
	s.WarehouseId = config.DefaultWarehouseId

	result = trans.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsSalesOrder(s.EnterpriseId, s.OrderId, userId, s.Price*float64(s.Quantity), s.VatPercent, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = setSalesOrderState(s.EnterpriseId, s.OrderId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = addSalesOrderLinesNumber(s.EnterpriseId, s.OrderId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	insertTransactionalLog(s.EnterpriseId, "sales_order_detail", int(s.Id), userId, "I")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order_detail", "POST", string(json))

	return OkAndErrorCodeReturn{Ok: true}
}

// 1. the product is deactivated
// 2. there is aleady a detail with this product
// 3. can't update an invoiced sale order detail
func (s *SalesOrderDetail) updateSalesOrderDetail(userId int32) OkAndErrorCodeReturn {
	if s.Id <= 0 || !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	p := getProductRow(s.ProductId)
	if p.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if p.Off {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	var countProductInSaleOrder int64
	result := dbOrm.Model(&SalesOrderDetail{}).Where("product = ? AND \"order\" = ? AND id != ?", s.ProductId, s.OrderId, s.Id).Count(&countProductInSaleOrder) // don't count the existing detail
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if countProductInSaleOrder > 0 { // we are not counting this existing detail
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	// get a single inMemoryDetail from the database by id and enterprise using dbOrm
	var inMemoryDetail SalesOrderDetail
	result = trans.Model(&SalesOrderDetail{}).Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).First(&inMemoryDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if inMemoryDetail.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if inMemoryDetail.QuantityInvoiced > 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	// take out the old value
	ok := addTotalProductsSalesOrder(s.EnterpriseId, inMemoryDetail.OrderId, userId, -(inMemoryDetail.Price * float64(inMemoryDetail.Quantity)), inMemoryDetail.VatPercent, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	inMemoryDetail.ProductId = s.ProductId
	inMemoryDetail.Price = s.Price
	inMemoryDetail.Quantity = s.Quantity
	inMemoryDetail.VatPercent = s.VatPercent
	inMemoryDetail.TotalAmount = s.TotalAmount
	inMemoryDetail.ShopifyId = s.ShopifyId

	// save the detail in the database using dbOrm
	result = dbOrm.Save(&inMemoryDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	// add the new value
	ok = addTotalProductsSalesOrder(s.EnterpriseId, s.OrderId, userId, s.Price*float64(s.Quantity), s.VatPercent, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	insertTransactionalLog(s.EnterpriseId, "sales_order_detail", int(s.Id), userId, "U")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))

	return OkAndErrorCodeReturn{Ok: true}
}

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a single transaction.
//
// ERROR CODES:
// 1. the detail is already invoiced
// 2. the detail has a delivery note generated
// 3. there are complex manufacturing orders already created
// 4. there are manufacturing orders already created
// 5. there is digital product data that must be deleted first
// 6. the product has been packaged
func (s *SalesOrderDetail) deleteSalesOrderDetail(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		var err error
		trans = dbOrm.Begin()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	detailInMemory := getSalesOrderDetailRow(s.Id)
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

	// check for complex_manufacturing_order_manufacturing_order using dbOrm
	var complexManufacturingOrderManufacturingOrderRows int64
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("sale_order_detail = ?", s.Id).Count(&complexManufacturingOrderManufacturingOrderRows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if complexManufacturingOrderManufacturingOrderRows > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
	}

	// check for manufacturing_order
	var manufacturingOrderRows int64
	result = dbOrm.Model(&ManufacturingOrder{}).Where("order_detail = ?", s.Id).Count(&manufacturingOrderRows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if manufacturingOrderRows > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 4}
	}

	// check for sales_order_detail_digital_product_data
	var salesOrderDetailDigitalProductDataRows int64
	result = dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Where("detail = ?", s.Id).Count(&salesOrderDetailDigitalProductDataRows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if salesOrderDetailDigitalProductDataRows > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 5}
	}

	// check for sales_order_detail_packaged
	var salesOrderDetailPackagedRows int64
	result = dbOrm.Model(&SalesOrderDetailPackaged{}).Where("order_detail = ?", s.Id).Count(&salesOrderDetailPackagedRows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if salesOrderDetailPackagedRows > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 6}
	}

	insertTransactionalLog(s.EnterpriseId, "sales_order_detail", int(s.Id), userId, "D")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order_detail", "DELETE", string(json))

	result = trans.Model(&SalesOrderDetail{}).Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).Delete(&SalesOrderDetail{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsSalesOrder(s.EnterpriseId, detailInMemory.OrderId, userId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = setSalesOrderState(detailInMemory.EnterpriseId, detailInMemory.OrderId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = removeSalesOrderLinesNumber(detailInMemory.EnterpriseId, detailInMemory.OrderId, userId, *trans)
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

	return OkAndErrorCodeReturn{Ok: true}
}

// Adds an invoiced quantity to the sale order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvociedSalesOrderDetail(detailId int64, quantity int32, userId int32, trans gorm.DB) bool {
	detailBefore := getSalesOrderDetailRowTransaction(detailId, trans)
	if detailBefore.Id <= 0 {
		return false
	}
	salesOrder := getSalesOrderRow(detailBefore.OrderId)
	if salesOrder.Id <= 0 {
		return false
	}

	detailAfter := getSalesOrderDetailRowTransaction(detailId, trans)
	if detailAfter.Id <= 0 {
		return false
	}

	detailAfter.QuantityInvoiced += quantity

	result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detailId).Update("quantity_invoiced", detailAfter.QuantityInvoiced)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	var ok bool
	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity { // set as invoced
		ok = addQuantityPendingServing(detailBefore.ProductId, detailBefore.WarehouseId, detailBefore.Quantity, detailBefore.EnterpriseId, trans)
		// set the order detail state applying the workflow logic
		if ok {
			status, purchaseOrderDetail, warehouseId := detailBefore.computeStatus(userId, trans)
			detailAfter.Status = status
			detailAfter.PurchaseOrderDetailId = purchaseOrderDetail
			if len(warehouseId) == 0 {
				config := getSettingsRecordById(detailBefore.EnterpriseId)
				warehouseId = config.DefaultWarehouseId
			}
			result = trans.Model(&SalesOrderDetail{}).Where("id = ?", detailId).Updates(map[string]interface{}{
				"status":                detailAfter.Status,
				"purchase_order_detail": detailAfter.PurchaseOrderDetailId,
				"warehouse":             warehouseId,
			})
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}
		}
		if !ok {
			return false
		}
		ok = addSalesOrderInvoicedLines(salesOrder.EnterpriseId, detailBefore.OrderId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok = addQuantityPendingServing(detailBefore.ProductId, detailBefore.WarehouseId, -detailBefore.Quantity, detailBefore.EnterpriseId, trans)
		// reset order detail state to "Waiting for Payment"
		if ok {
			detailAfter.Status = "_"
			detailAfter.PurchaseOrderDetailId = nil
			result = trans.Model(&SalesOrderDetail{}).Where("id = ?", detailId).Updates(map[string]interface{}{
				"status":                detailAfter.Status,
				"purchase_order_detail": detailAfter.PurchaseOrderDetailId,
			})
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}
		}
		if !ok {
			return false
		}
		ok = removeSalesOrderInvoicedLines(salesOrder.EnterpriseId, detailBefore.OrderId, userId, trans)
		if !ok {
			return false
		}

		// reset relations
		if detailBefore.PurchaseOrderDetailId != nil {
			ok := addQuantityAssignedSalePurchaseOrder(*detailBefore.PurchaseOrderDetailId, detailBefore.Quantity, detailBefore.EnterpriseId, userId, trans)
			if !ok {
				return false
			}
		}
		orders := getSalesOrderManufacturingOrders(salesOrder.Id, salesOrder.EnterpriseId)
		for i := 0; i < len(orders); i++ {
			if orders[i].OrderDetailId != nil || *orders[i].OrderDetailId != detailBefore.Id {
				continue
			}

			ok := orders[i].deleteManufacturingOrder(userId, &trans)
			if !ok {
				return false
			}
		}
		result = trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("sale_order_detail = ?", detailBefore.Id).Update("sale_order_detail", nil)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
		// -- reset relations
	}

	ok = setSalesOrderState(salesOrder.EnterpriseId, salesOrder.Id, userId, trans)
	if !ok {
		return false
	}

	insertTransactionalLog(detailBefore.EnterpriseId, "sales_order_detail", int(detailId), userId, "U")
	s := getSalesOrderRowTransaction(detailId, trans)
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))

	return true
}

// returns: status, purchase order detail id, warehouse id
func (s *SalesOrderDetail) computeStatus(userId int32, trans gorm.DB) (string, *int64, string) {
	product := getProductRow(s.ProductId)
	if product.Id <= 0 {
		return "", nil, ""
	}

	order := getSalesOrderRow(s.OrderId)
	stock := getStockRowAvailable(s.ProductId, s.EnterpriseId)
	if !product.ControlStock {
		return "E", nil, ""
	} else if stock.QuantityAvaiable > 0 { // the product is in stock, send to preparation
		return "E", nil, stock.WarehouseId
	} else { // the product is not in stock, purchase or manufacture
		if product.Manufacturing {
			// search for pending manufacturing order for stock
			manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderTypeId)
			if manufacturingOrderType.Complex {
				rows, err := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("product = ? AND type = 'O' AND manufactured = false AND sale_order_detail IS NULL", s.ProductId).Order("id ASC").Select(" id, manufacturing_order_type_component, warehouse").Rows()
				if err != nil {
					log("DB", err.Error())
					// fallback
					return "C", nil, ""
				}
				defer rows.Close()

				var orders []int64 = make([]int64, 0)
				var quantities []int32 = make([]int32, 0)
				var totalQuantityManufactured int32 = 0
				var warehouseId string

				for rows.Next() {
					var complexManufacturingOrderForStockId int64
					var manufacturingOrderTypeComponentId int32
					rows.Scan(&complexManufacturingOrderForStockId, &manufacturingOrderTypeComponentId, &warehouseId)
					orders = append(orders, complexManufacturingOrderForStockId)

					com := getManufacturingOrderTypeComponentRow(manufacturingOrderTypeComponentId)
					quantities = append(quantities, com.Quantity)
					totalQuantityManufactured += com.Quantity
				}

				if totalQuantityManufactured >= s.Quantity {
					var quantityAssigned int32 = 0
					for i := 0; i < len(orders); i++ {
						result := trans.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("sale_order_detail = ?", orders[i]).Update("sale_order_detail", s.Id)
						if result.Error != nil {
							log("DB", result.Error.Error())
							// fallback
							return "C", nil, ""
						}

						quantityAssigned += quantities[i]
						if quantityAssigned >= s.Quantity {
							break
						}
					}

					return "D", nil, warehouseId
				} else {
					return "C", nil, ""
				}
			} else { // if manufacturingOrderType.Complex {
				rows, err := dbOrm.Model(&ManufacturingOrder{}).Where("manufactured = false AND product = $1 AND complex = false", product.Id).Order("date_created ASC").Select("id, quantity_manufactured, warehouse").Rows()
				if err != nil {
					log("DB", err.Error())
					// fallback
					return "C", nil, ""
				}
				defer rows.Close()
				var totalQuantityManufactured int32 = 0
				var orders []int64 = make([]int64, 0)
				var quantities []int32 = make([]int32, 0)
				var warehouseId string

				for rows.Next() {
					var manufacturingOrderForStockId int64
					var quantityManufactured int32
					rows.Scan(&manufacturingOrderForStockId, &quantityManufactured, &warehouseId)
					totalQuantityManufactured += quantityManufactured
					orders = append(orders, manufacturingOrderForStockId)
					quantities = append(quantities, quantityManufactured)
				}

				if totalQuantityManufactured < s.Quantity {
					return "C", nil, ""
				} else {
					var quantityAssigned int32 = 0
					for i := 0; i < len(orders); i++ {
						result := trans.Model(&ManufacturingOrder{}).Where("id = ?", orders[i]).Save(map[string]interface{}{
							"order_detail": s.Id,
							"order":        s.OrderId,
						})
						if result.Error != nil {
							log("DB", result.Error.Error())
							// fallback
							return "C", nil, ""
						}

						quantityAssigned += quantities[i]
						if quantityAssigned >= s.Quantity {
							break
						}
					}
					return "D", nil, warehouseId
				}
			}
		} else {
			// search for pending purchases using dbOrm
			var purchaseDetail PurchaseOrderDetail
			result := dbOrm.Model(&PurchaseOrderDetail{}).Where("product = ? AND quantity_delivery_note = 0 AND quantity - quantity_assigned_sale >= ?", s.ProductId, s.Quantity).Order(`(SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order") ASC`).Limit(1).First(&purchaseDetail)
			if result.Error != nil {
				log("DB", result.Error.Error())
				// fallback
				return "A", nil, ""
			}

			if purchaseDetail.Id <= 0 {
				return "A", nil, ""
			}

			// add quantity assigned to sale orders
			ok := addQuantityAssignedSalePurchaseOrder(purchaseDetail.Id, s.Quantity, order.EnterpriseId, userId, trans)
			if !ok {
				return "A", nil, ""
			}

			// set the purchase order detail
			return "B", &purchaseDetail.Id, purchaseDetail.WarehouseId
		}
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityPendingPackagingSaleOrderDetail(detailId int64, quantity int32, userId int32, trans gorm.DB) bool {
	detail := getSalesOrderDetailRow(detailId)
	if detail.Id <= 0 {
		trans.Rollback()
		return false
	}
	detail.QuantityPendingPackaging += quantity

	if detail.QuantityPendingPackaging <= 0 {
		detail.Status = "F"
	} else {
		detail.Status = "E"
	}

	result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detailId).Updates(map[string]interface{}{
		"quantity_pending_packaging": detail.QuantityPendingPackaging,
		"status":                     detail.Status,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if setSalesOrderState(detail.EnterpriseId, detail.OrderId, userId, trans) {
		insertTransactionalLog(detail.EnterpriseId, "sales_order_detail", int(detailId), userId, "U")
		json, _ := json.Marshal(detail)
		go fireWebHook(detail.EnterpriseId, "sales_order_detail", "PUT", string(json))
		return true
	} else {
		return false
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNoteSalesOrderDetail(detailId int64, quantity int32, userId int32, trans gorm.DB) bool {
	detailBefore := getSalesOrderDetailRowTransaction(detailId, trans)
	if detailBefore.Id <= 0 {
		return false
	}

	detailAfter := getSalesOrderDetailRowTransaction(detailId, trans)
	if detailAfter.Id <= 0 {
		return false
	}

	detailAfter.QuantityDeliveryNote += quantity

	result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detailId).Update("quantity_delivery_note", detailAfter.QuantityDeliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	var ok bool
	if detailBefore.QuantityDeliveryNote != detailBefore.Quantity && detailAfter.QuantityDeliveryNote == detailAfter.Quantity { // set as delivery note generated
		ok = addSalesOrderDeliveryNoteLines(detailBefore.EnterpriseId, detailBefore.OrderId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok = removeSalesOrderDeliveryNoteLines(detailBefore.EnterpriseId, detailBefore.OrderId, userId, trans)
		if !ok {
			return false
		}
	}

	insertTransactionalLog(detailBefore.EnterpriseId, "sales_order_detail", int(detailId), userId, "U")
	s := getSalesOrderDetailRow(detailId)
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))
	return true
}

func cancelSalesOrderDetail(detailId int64, enterpriseId int32, userId int32) bool {
	detail := getSalesOrderDetailRow(detailId)
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

		detail.QuantityInvoiced += detail.Quantity
		detail.QuantityDeliveryNote += detail.Quantity
		detail.Status = "Z"
		detail.Cancelled = true

		result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detail.Id).Updates(map[string]interface{}{
			"quantity_invoiced":      detail.QuantityInvoiced,
			"quantity_delivery_note": detail.QuantityDeliveryNote,
			"status":                 detail.Status,
			"cancelled":              detail.Cancelled,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		ok := setSalesOrderState(enterpriseId, detail.OrderId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}

		insertTransactionalLog(detail.EnterpriseId, "sales_order_detail", int(detailId), userId, "U")
		json, _ := json.Marshal(detail)
		go fireWebHook(enterpriseId, "sales_order_detail", "PUT", string(json))

		///
		result = trans.Commit()
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false
		}
		///

		return true
	} else {
		if detail.Quantity <= 0 || detail.QuantityInvoiced == 0 || detail.QuantityDeliveryNote == 0 {
			return false
		}

		detail.QuantityInvoiced = 0
		detail.QuantityDeliveryNote = 0
		detail.Cancelled = false

		result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detail.Id).Updates(map[string]interface{}{
			"quantity_invoiced":      detail.QuantityInvoiced,
			"quantity_delivery_note": detail.QuantityDeliveryNote,
			"cancelled":              detail.Cancelled,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		status, purchaseOrderDetail, warehouseId := detail.computeStatus(userId, *trans)
		if len(warehouseId) == 0 {
			config := getSettingsRecordById(detail.EnterpriseId)
			warehouseId = config.DefaultWarehouseId
		}

		detail.Status = status
		detail.PurchaseOrderDetailId = purchaseOrderDetail

		result = trans.Model(&SalesOrderDetail{}).Where("id = ?", detail.Id).Updates(map[string]interface{}{
			"status":                detail.Status,
			"purchase_order_detail": detail.PurchaseOrderDetailId,
			"warehouse":             warehouseId,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
		}

		ok := setSalesOrderState(enterpriseId, detail.OrderId, userId, *trans)
		if !ok {
			return false
		}

		insertTransactionalLog(detail.EnterpriseId, "sales_order_detail", int(detailId), userId, "U")
		s := getSalesOrderDetailRow(detailId)
		json, _ := json.Marshal(s)
		go fireWebHook(s.EnterpriseId, "sales_order_detail", "PUT", string(json))

		///
		result = trans.Commit()
		return result.Error == nil
		///
	}
}

type SalePurchasesOrderDetail struct {
	Id          int32     `json:"id"`
	OrderId     int64     `json:"orderId" gorm:"column:order;not null:true;index:sales_order_detail_sales_order_product,unique:true,priority:1"`
	Order       SaleOrder `json:"-" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity    int32     `json:"quantity"`
	TotalAmount float64   `json:"totalAmount"`
}

func getPurchasesOrderDetailsFromSaleOrderDetail(detailId int32, enterpriseId int32) []SalesOrderDetail {
	saleOrderDetails := make([]SalesOrderDetail, 0)
	result := dbOrm.Where("sales_order_detail.enterprise = ? AND sales_order_detail.purchase_order_detail = ?", enterpriseId, detailId).Preload(clause.Associations).Find(&saleOrderDetails)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return saleOrderDetails
}

func filterSalesOrderDetails(input []SalesOrderDetail, test func(SalesOrderDetail) bool) (output []SalesOrderDetail) {
	for _, s := range input {
		if test(s) {
			output = append(output, s)
		}
	}
	return
}
