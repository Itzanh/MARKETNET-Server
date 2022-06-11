package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductIncludedProduct struct {
	Id                int32    `json:"id" gorm:"column:id;index:product_included_products_id_enterprise,unique:true,priority:1"`
	ProductBaseId     int32    `json:"productBaseId" gorm:"column:product_base;type:integer;not null:true;index:product_included_products_product_base_product_included,unique:true,priority:1"`
	ProductBase       Product  `json:"productBase" gorm:"foreignKey:ProductBaseId,EnterpriseId;references:Id,EnterpriseId"`
	ProductIncludedId int32    `json:"productIncludedId" gorm:"column:product_included;type:integer;not null:true;index:product_included_products_product_base_product_included,unique:true,priority:2"`
	ProductIncluded   Product  `json:"productIncluded" gorm:"foreignKey:ProductIncludedId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity          int32    `json:"quantity" gorm:"column:quantity;type:integer;not null:true"`
	EnterpriseId      int32    `json:"-" gorm:"column:enterprise;not null:true;index:product_included_products_id_enterprise,unique:true,priority:2"`
	Enterprise        Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *ProductIncludedProduct) TableName() string {
	return "product_included_products"
}

func getProductIncludedProduct(productId int32, enterpriseId int32) []ProductIncludedProduct {
	var productIncludedProduct []ProductIncludedProduct = make([]ProductIncludedProduct, 0)
	result := dbOrm.Model(&ProductIncludedProduct{}).Where("product_base = ? AND enterprise = ?", productId, enterpriseId).Preload(clause.Associations).Order("id ASC").Find(&productIncludedProduct)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return productIncludedProduct
	}
	return productIncludedProduct
}

func (p *ProductIncludedProduct) isValid() bool {
	return !(p.ProductBaseId <= 0 || p.ProductIncludedId <= 0 || p.ProductBaseId == p.ProductIncludedId || p.Quantity <= 0)
}

func (p *ProductIncludedProduct) BeforeCreate(tx *gorm.DB) (err error) {
	var productIncludedProduct ProductIncludedProduct
	tx.Model(&ProductIncludedProduct{}).Last(&productIncludedProduct)
	p.Id = productIncludedProduct.Id + 1
	return nil
}

func (p *ProductIncludedProduct) insertProductIncludedProduct() bool {
	if !p.isValid() {
		return false
	}

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (p *ProductIncludedProduct) updateProductIncludedProduct() bool {
	if !p.isValid() {
		return false
	}

	var productIncludedProduct ProductIncludedProduct
	result := dbOrm.Where("id = ?", p.Id).First(&productIncludedProduct)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	if productIncludedProduct.Id <= 0 || productIncludedProduct.EnterpriseId != p.EnterpriseId {
		return false
	}

	productIncludedProduct.Quantity = p.Quantity

	result = dbOrm.Save(&productIncludedProduct)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (p *ProductIncludedProduct) deleteProductIncludedProduct() bool {
	result := dbOrm.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).Delete(&ProductIncludedProduct{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

type ProductIncludedProductSalesOrderDetail struct {
	Id                        int64                  `json:"id"`
	ProductIncludedProductsId int32                  `json:"-" gorm:"column:product_included_products;type:integer;not null:true;index:product_included_products_sales_order_details_product_ipsod,unique:true,priority:1;index:product_included_products_sales_order_details_product_ipso,unique:true,priority:1"`
	ProductIncludedProduct    ProductIncludedProduct `json:"productIncludedProduct" gorm:"foreignKey:ProductIncludedProductsId,EnterpriseId;references:Id,EnterpriseId"`
	SalesOrderDetailId        int64                  `json:"-" gorm:"column:sales_order_detail;type:bigint;not null:true;index:product_included_products_sales_order_details_product_ipsod,unique:true,priority:2"`
	SalesOrderDetail          SalesOrderDetail       `json:"-" gorm:"foreignKey:SalesOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	SalesOrderId              int64                  `json:"-" gorm:"column:sales_order;type:bigint;not null:true;index:product_included_products_sales_order_details_product_ipso,unique:true,priority:2"`
	SalesOrder                SaleOrder              `json:"-" gorm:"foreignKey:SalesOrderId,EnterpriseId;references:Id,EnterpriseId"`
	QuantityUnit              int32                  `json:"quantityUnit" gorm:"column:quantity_unit;type:integer;not null:true"`
	EnterpriseId              int32                  `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise                Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *ProductIncludedProductSalesOrderDetail) TableName() string {
	return "product_included_products_sales_order_details"
}

func getProductIncludedProductSalesOrderDetail(salesOrderDetailId int64, enterpriseId int32) []ProductIncludedProductSalesOrderDetail {
	var productIncludedProductSalesOrderDetail []ProductIncludedProductSalesOrderDetail = make([]ProductIncludedProductSalesOrderDetail, 0)
	result := dbOrm.Model(&ProductIncludedProductSalesOrderDetail{}).Where("sales_order_detail = ? AND enterprise = ?", salesOrderDetailId, enterpriseId).Preload(clause.Associations).Preload("ProductIncludedProduct.ProductBase").Order("id ASC").Find(&productIncludedProductSalesOrderDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return productIncludedProductSalesOrderDetail
}

// FOR INTERNAL USE ONLY
// This funcion is used to tell if there is already a sales order detail for the incluided product in the sale order or not
func getProductIncludedProductSalesOrder(productIncludedId int32, salesOrderId int64) *ProductIncludedProductSalesOrderDetail {
	var productIncludedProductSalesOrderDetail *ProductIncludedProductSalesOrderDetail
	result := dbOrm.Where("product_included_products = ? AND sales_order = ?", productIncludedId, salesOrderId).First(&productIncludedProductSalesOrderDetail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return productIncludedProductSalesOrderDetail
}

// A new line is inserted, check if the product of the line has included products. If it has, handle the included products on the order lines.
func (s *SalesOrderDetail) processProductIncludedProductOnNewInsertedLine(enterpriseId int32, userId int32) {
	productIncludedProduct := getProductIncludedProduct(s.ProductId, enterpriseId)
	if len(productIncludedProduct) == 0 { // The product does not have included products, skip this process
		return
	}

	for _, productIncluded := range productIncludedProduct {
		var productIncludedSalesOrderDetailId int64
		// Check if the included product is already in the sales order or not
		productIncludedProductSalesOrderDetail := getSalesOrderDetailRowByOrderAndProduct(s.OrderId, productIncluded.ProductIncludedId)
		if productIncludedProductSalesOrderDetail == nil {
			// The included product is not in the sales order, insert a new detail for it
			var newDetail = SalesOrderDetail{
				OrderId:          s.OrderId,
				ProductId:        productIncluded.ProductIncludedId,
				Price:            0,
				Quantity:         productIncluded.Quantity * s.Quantity,
				VatPercent:       0,
				EnterpriseId:     enterpriseId,
				IncludedProducts: true,
			}
			newDetail.insertSalesOrderDetail(userId)
			productIncludedSalesOrderDetailId = newDetail.Id
		} else {
			// The included product is already in the sales order, update the quantity of the detail
			oldDetail := getSalesOrderDetailRow(productIncludedProductSalesOrderDetail.Id)
			oldDetail.Quantity += productIncluded.Quantity * s.Quantity
			oldDetail.updateSalesOrderDetail(userId)
			productIncludedSalesOrderDetailId = oldDetail.Id
		}
		// Insert a new Product Included - Sale Order Detail to keep track of the included product
		var newProductIncludedProductSalesOrderDetail = ProductIncludedProductSalesOrderDetail{
			ProductIncludedProductsId: productIncluded.Id,
			SalesOrderDetailId:        productIncludedSalesOrderDetailId,
			SalesOrderId:              s.OrderId,
			QuantityUnit:              productIncluded.Quantity,
			EnterpriseId:              enterpriseId,
		}
		newProductIncludedProductSalesOrderDetail.insertProductIncludedProductSalesOrderDetail()
	}
}

// An existing line is updated, check if the product of the line has included products. If it has, handle the included products on the order lines.
func (s *SalesOrderDetail) processProductIncludedProductOnUpdatedLine(enterpriseId int32, userId int32, oldQuantity int32) {
	productIncludedProduct := getProductIncludedProduct(s.ProductId, enterpriseId)
	if len(productIncludedProduct) == 0 { // The product does not have included products, skip this process
		return
	}

	for _, productIncluded := range productIncludedProduct {
		productIncludedProductSalesOrderDetail := getProductIncludedProductSalesOrder(productIncluded.Id, s.OrderId)
		if productIncludedProductSalesOrderDetail == nil { // Unknown internal error, prevent the program from crashing (null pointer exception)
			continue
		}
		// Find the detail with the included product, and increase / decrease the quantity
		oldDetail := getSalesOrderDetailRow(productIncludedProductSalesOrderDetail.SalesOrderDetailId)
		// To change the quantity, the quantity in the included product will not be used,
		// insted it is going to use the quantity that was prevously saved in the Product Included - Sales Order Detail to keep track of the included products
		oldDetail.Quantity -= productIncludedProductSalesOrderDetail.QuantityUnit * oldQuantity
		oldDetail.Quantity += productIncludedProductSalesOrderDetail.QuantityUnit * s.Quantity
		oldDetail.updateSalesOrderDetail(userId)
	}
}

func (s *SalesOrderDetail) processProductIncludedProductOnDeletedLine(enterpriseId int32, userId int32) {
	productIncludedProduct := getProductIncludedProduct(s.ProductId, enterpriseId)
	if len(productIncludedProduct) == 0 { // The product does not have included products, skip this process
		return
	}

	for _, productIncluded := range productIncludedProduct {
		productIncludedProductSalesOrderDetail := getProductIncludedProductSalesOrder(productIncluded.Id, s.OrderId)
		if productIncludedProductSalesOrderDetail == nil { // Unknown internal error, prevent the program from crashing (null pointer exception)
			continue
		}
		// Find the detail with the included product, and increase / decrease the quantity
		oldDetail := getSalesOrderDetailRow(productIncludedProductSalesOrderDetail.SalesOrderDetailId)
		// To change the quantity, the quantity in the included product will not be used,
		// insted it is going to use the quantity that was prevously saved in the Product Included - Sales Order Detail to keep track of the included products
		oldDetail.Quantity -= productIncludedProductSalesOrderDetail.QuantityUnit * s.Quantity
		if oldDetail.Quantity <= 0 {
			// The quantity is 0 or less, delete the detail for the included product, as it's no longer needed
			productIncludedProductSalesOrderDetail.deleteProductIncludedProductSalesOrderDetail()
			oldDetail.deleteSalesOrderDetail(userId, nil)
		} else {
			// Update the detail for the included product with the new quantity
			productIncludedProductSalesOrderDetail.deleteProductIncludedProductSalesOrderDetail()
			oldDetail.updateSalesOrderDetail(userId)
		}
	}
}

func (p *ProductIncludedProductSalesOrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var productIncludedProductSalesOrderDetail ProductIncludedProductSalesOrderDetail
	tx.Model(&ProductIncludedProductSalesOrderDetail{}).Last(&productIncludedProductSalesOrderDetail)
	p.Id = productIncludedProductSalesOrderDetail.Id + 1
	return nil
}

func (p *ProductIncludedProductSalesOrderDetail) insertProductIncludedProductSalesOrderDetail() bool {
	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (p *ProductIncludedProductSalesOrderDetail) deleteProductIncludedProductSalesOrderDetail() bool {
	result := dbOrm.Where("id = ?", p.Id).Delete(&ProductIncludedProductSalesOrderDetail{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
