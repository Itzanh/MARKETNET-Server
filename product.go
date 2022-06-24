/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	Id                       int32                   `json:"id" gorm:"index:product_id_enterprise,unique:true,priority:1"`
	Name                     string                  `json:"name" gorm:"type:character varying(150);not null:true;index:product_name,type:gin"`
	Reference                string                  `json:"reference" gorm:"type:character varying(40);not null:true;index:product_reference,type:gin,where:reference::text <> ''::text"`
	BarCode                  string                  `json:"barCode" gorm:"column:barcode;type:character(13);not null:true;index:product_barcode,unique:true,priority:2,where:barcode <> ''::bpchar"`
	ControlStock             bool                    `json:"controlStock" gorm:"not null:true"`
	Weight                   float64                 `json:"weight" gorm:"type:numeric(14,6);not null:true"`
	FamilyId                 *int32                  `json:"familyId" gorm:"column:family"`
	Family                   *ProductFamily          `json:"family" gorm:"foreignKey:FamilyId,EnterpriseId;references:Id,EnterpriseId"`
	Width                    float64                 `json:"width" gorm:"type:numeric(14,6);not null:true"`
	Height                   float64                 `json:"height" gorm:"type:numeric(14,6);not null:true"`
	Depth                    float64                 `json:"depth" gorm:"type:numeric(14,6);not null:true"`
	Off                      bool                    `json:"off" gorm:"not null:true"`
	Stock                    int32                   `json:"stock" gorm:"not null:true"`
	VatPercent               float64                 `json:"vatPercent" gorm:"type:numeric(14,6);not null:true"`
	DateCreated              time.Time               `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Description              string                  `json:"description" gorm:"type:text;column:dsc"`
	ColorId                  *int32                  `json:"colorId" gorm:"column:color"`
	Color                    *Color                  `json:"color" gorm:"foreignKey:ColorId,EnterpriseId;references:Id,EnterpriseId"`
	Price                    float64                 `json:"price" gorm:"type:numeric(14,6);not null:true"`
	Manufacturing            bool                    `json:"manufacturing" gorm:"not null:true"`
	ManufacturingOrderTypeId *int32                  `json:"manufacturingOrderTypeId" gorm:"column:manufacturing_order_type"`
	ManufacturingOrderType   *ManufacturingOrderType `json:"-" gorm:"foreignKey:ManufacturingOrderTypeId,EnterpriseId;references:Id,EnterpriseId"`
	SupplierId               *int32                  `json:"supplierId" gorm:"column:supplier"`
	Supplier                 *Supplier               `json:"supplier" gorm:"foreignKey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopId             int32                   `json:"-" gorm:"column:ps_id;not null:true;index:product_ps_id,unique:true,priority:2,where:ps_id <> 0"`
	PrestaShopCombinationId  int32                   `json:"-" gorm:"column:ps_combination_id;not null:true;index:product_ps_id,unique:true,priority:3,where:ps_id <> 0"`
	MinimumStock             int32                   `json:"minimumStock" gorm:"not null:true"`
	TrackMinimumStock        bool                    `json:"trackMinimumStock" gorm:"not null:true;index:product_track_minimum_stock,where:track_minimum_stock = true"`
	WooCommerceId            int32                   `json:"-" gorm:"column:wc_id;not null:true;index:products_wc_id,unique:true,priority:2,where:wc_id <> 0"`
	WooCommerceVariationId   int32                   `json:"-" gorm:"column:wc_variation_id;not null:true;index:products_wc_id,unique:true,priority:3,where:wc_id <> 0"`
	ShopifyId                int64                   `json:"-" gorm:"column:sy_id;not null:true;index:product_sy_id,unique:true,priority:2,where:sy_id <> 0"`
	ShopifyVariantId         int64                   `json:"-" gorm:"column:sy_variant_id;not null:true;index:product_sy_id,unique:true,priority:3,where:sy_id <> 0"`
	EnterpriseId             int32                   `json:"-" gorm:"column:enterprise;not null:true;index:product_id_enterprise,unique:true,priority:2;index:product_barcode,unique:true,priority:1,where:barcode <> ''::bpchar;;index:product_ps_id,unique:true,priority:1,where:ps_id <> 0;;index:product_sy_id,unique:true,priority:1,where:sy_id <> 0;;index:products_wc_id,unique:true,priority:1,where:wc_id <> 0"`
	Enterprise               Settings                `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	DigitalProduct           bool                    `json:"digitalProduct" gorm:"not null:true"`
	PurchasePrice            float64                 `json:"purchasePrice" gorm:"type:numeric(14,6);not null:true"`
	MinimumPurchaseQuantity  int32                   `json:"minimumPurchaseQuantity" gorm:"not null:true"`
	OriginCountry            string                  `json:"originCountry" gorm:"type:character varying(2);not null:true"`
	HSCodeId                 *string                 `json:"HSCodeId" gorm:"column:hs_code;type:character varying(8)"`
	HSCode                   *HSCode                 `json:"HSCode" gorm:"foreignKey:HSCodeId;references:Id"`
	CostPrice                float64                 `json:"costPrice" gorm:"type:numeric(14,6);not null:true"`
}

func (p *Product) TableName() string {
	return "product"
}

func getProduct(enterpriseId int32) []Product {
	var products []Product = make([]Product, 0)
	result := dbOrm.Model(&Product{}).Where("product.enterprise = ?", enterpriseId).Preload(clause.Associations).Order("product.id ASC").Find(&products)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return products
}

type ProductSearch struct {
	Search            string `json:"search"`
	TrackMinimumStock bool   `json:"trackMinimumStock"`
	FamilyId          *int32 `json:"familyId"`
}

func (search *ProductSearch) searchProduct(enterpriseId int32) []Product {
	var products []Product = make([]Product, 0)

	// Check that the format for EAN13 barcodes is correct
	if len(search.Search) == 13 {
		if checkEan13(search.Search) {
			product := getProductByBarcode(search.Search, enterpriseId)
			if product.Id > 0 {
				products = append(products, product)
			}
			return products
		}
	}

	query := `((product.name ILIKE @search) OR (product.barcode = @text) OR (product.reference ILIKE @search) OR (product_family.name ILIKE @search)) AND (product.enterprise = @enterpriseId)`
	cursor := dbOrm.Model(&Product{}).Where(query, sql.Named("search", "%"+search.Search+"%"), sql.Named("enterpriseId", enterpriseId), sql.Named("text", search.Search))
	if search.TrackMinimumStock {
		cursor = cursor.Where("product.track_minimum_stock = true")
	}
	if search.FamilyId != nil {
		cursor = cursor.Where("product.family = ?", *search.FamilyId)
	}
	result := cursor.Joins("FULL JOIN product_family ON product.family = product_family.id").Preload(clause.Associations).Order("id ASC").Find(&products)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return products
}

func getProductRow(productId int32) Product {
	p := Product{}
	result := dbOrm.Model(&Product{}).Where("id = ?", productId).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return p
}

func getProductByBarcode(ean13 string, enterpriseId int32) Product {
	p := Product{}
	dbOrm.Model(&Product{}).Where("enterprise = ? AND barcode = ?", enterpriseId, ean13).First(&p)
	return p
}

func (p *Product) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 150 || len(p.Reference) > 40 || (len(p.BarCode) != 0 && len(p.BarCode) != 13) || p.VatPercent < 0 || p.Price < 0 || p.Weight < 0 || p.Width < 0 || p.Height < 0 || p.Depth < 0 || p.MinimumPurchaseQuantity < 0 || p.CostPrice < 0 || len(p.Description) > 3000)
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	var product Product
	tx.Model(&Product{}).Last(&product)
	p.Id = product.Id + 1
	return nil
}

func (p *Product) insertProduct(userId int32) OkAndErrorCodeReturn {
	if !p.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}
	p.BarCode = strings.Trim(p.BarCode, " ")

	// Check that the format for EAN13 barcodes is correct
	if len(p.BarCode) == 13 {
		if !checkEan13(p.BarCode) {
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	if len(p.BarCode) > 0 {
		var ean13AlreadyExists int64
		dbOrm.Model(&Product{}).Where("barcode = ? AND enterprise = ?", p.BarCode, p.EnterpriseId).Count(&ean13AlreadyExists)
		if ean13AlreadyExists > 0 {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
	}

	p.DateCreated = time.Now()

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(p.EnterpriseId, "product", int(p.Id), userId, "I")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "product", "POST", string(json))

	return OkAndErrorCodeReturn{Ok: true}
}

func (p *Product) updateProduct(userId int32) OkAndErrorCodeReturn {
	if p.Id <= 0 || !p.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}
	p.BarCode = strings.Trim(p.BarCode, " ")

	// Check that the format for EAN13 barcodes is correct
	if len(p.BarCode) == 13 {
		if !checkEan13(p.BarCode) {
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	if len(p.BarCode) > 0 {
		var ean13AlreadyExists int64
		dbOrm.Model(&Product{}).Where("barcode = ? AND enterprise = ? AND id != ?", p.BarCode, p.EnterpriseId, p.Id).Count(&ean13AlreadyExists)
		if ean13AlreadyExists > 0 {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
	}

	var product Product
	result := dbOrm.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).First(&product)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	product.Name = p.Name
	product.Reference = p.Reference
	product.BarCode = p.BarCode
	product.ControlStock = p.ControlStock
	product.FamilyId = p.FamilyId
	product.Weight = p.Weight
	product.Width = p.Width
	product.Height = p.Height
	product.Depth = p.Depth
	product.Off = p.Off
	product.Stock = p.Stock
	product.VatPercent = p.VatPercent
	product.Description = p.Description
	product.ColorId = p.ColorId
	product.Price = p.Price
	product.Manufacturing = p.Manufacturing
	product.ManufacturingOrderTypeId = p.ManufacturingOrderTypeId
	product.SupplierId = p.SupplierId
	product.MinimumStock = p.MinimumStock
	product.TrackMinimumStock = p.TrackMinimumStock
	product.DigitalProduct = p.DigitalProduct
	product.PurchasePrice = p.PurchasePrice
	product.MinimumPurchaseQuantity = p.MinimumPurchaseQuantity
	product.OriginCountry = p.OriginCountry
	product.HSCodeId = p.HSCodeId
	product.CostPrice = p.CostPrice

	result = dbOrm.Save(&product)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(p.EnterpriseId, "product", int(p.Id), userId, "U")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "product", "PUT", string(json))

	return OkAndErrorCodeReturn{Ok: true}
}

// ERROR CODES:
// 1: The product has plurals
func (p *Product) deleteProduct(userId int32) OkAndErrorCodeReturn {
	if p.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	// check plurals
	var plurals []string = make([]string, 0)

	var rowsCount int64
	result := dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "1")
	}

	result = dbOrm.Model(&ManufacturingOrder{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "2")
	}

	result = dbOrm.Model(&ManufacturingOrderTypeComponents{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "3")
	}

	result = dbOrm.Model(&Packages{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "4")
	}

	result = dbOrm.Model(&ProductAccount{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "5")
	}

	result = dbOrm.Model(&ProductImage{}).Where("product = ?", p.Id).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "6")
	}

	result = dbOrm.Model(&PurchaseInvoiceDetail{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "8")
	}

	result = dbOrm.Model(&PurchaseOrderDetail{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "9")
	}

	result = dbOrm.Model(&SalesInvoiceDetail{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "10")
	}

	result = dbOrm.Model(&SalesOrderDetail{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "11")
	}

	result = dbOrm.Model(&WarehouseMovement{}).Where("product = ? AND enterprise = ?", p.Id, p.EnterpriseId).Count(&rowsCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	if rowsCount > 0 {
		plurals = append(plurals, "12")
	}

	if len(plurals) > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1, ExtraData: plurals}
	}

	result = trans.Where("product = ?", p.Id).Delete(&Stock{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(p.EnterpriseId, "product", int(p.Id), userId, "D")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "product", "DELETE", string(json))

	result = trans.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).Delete(&Product{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

func findProductByName(productName string, enterpriseId int32) []NameInt32 {
	var products []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&Product{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise=? AND off=false", strings.ToUpper(productName), enterpriseId).Order("id ASC").Limit(10).Find(&products)
	return products
}

type OrderDetailDefaults struct {
	Price                   float64 `json:"price" gorm:"type:numeric(14,6);not null:true"`
	PurchasePrice           float64 `json:"purchasePrice" gorm:"type:numeric(14,6);not null:true"`
	VatPercent              float64 `json:"vatPercent" gorm:"type:numeric(14,6);not null:true"`
	MinimumPurchaseQuantity int32   `json:"minimumPurchaseQuantity" gorm:"not null:true"`
}

func getOrderDetailDefaults(productId int32, enterpriseId int32) OrderDetailDefaults {
	s := OrderDetailDefaults{}
	result := dbOrm.Model(&Product{}).Where("id = ? AND enterprise = ?", productId, enterpriseId).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

// Get the sales order details with pending status, with the product specified.
func getProductSalesOrderDetailsPending(productId int32, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	result := dbOrm.Model(&SalesOrderDetail{}).Where("sales_order_detail.product = ? AND sales_order_detail.quantity_delivery_note != sales_order_detail.quantity AND sales_order_detail.enterprise = ?", productId, enterpriseId).Joins("Product").Order("sales_order_detail.id DESC").Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return details
	}
	return details
}

// Get the purchase order details with pending status, with the product specified.
func getProductPurchaseOrderDetailsPending(productId int32, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	result := dbOrm.Model(&PurchaseOrderDetail{}).Where("purchase_order_detail.product = ? AND purchase_order_detail.quantity_delivery_note != purchase_order_detail.quantity AND purchase_order_detail.enterprise = ?", productId, enterpriseId).Joins("Product").Order("purchase_order_detail.id DESC").Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return details
	}
	return details
}

type ProductSalesOrderDetailsQuery struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Status    string     `json:"status"`
	ProductId int32      `json:"productId"`
}

func (q *ProductSalesOrderDetailsQuery) isValid() bool {
	return !(q.ProductId <= 0 || (len(q.Status) > 1))
}

// Get the sales order details with the product specified.
func getProductSalesOrderDetails(query ProductSalesOrderDetailsQuery, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	if !query.isValid() {
		return details
	}

	cursor := dbOrm.Model(&SalesOrderDetail{}).Where("sales_order_detail.product = ? AND sales_order_detail.enterprise = ?", query.ProductId, enterpriseId)
	if query.StartDate != nil {
		cursor.Where("(SELECT date_created FROM sales_order WHERE sales_order.id = sales_order_detail.\"order\") >= ?", query.StartDate)
	}
	if query.EndDate != nil {
		cursor.Where("(SELECT date_created FROM sales_order WHERE sales_order.id = sales_order_detail.\"order\") >= ?", query.EndDate)
	}
	if query.Status != "" {
		cursor.Where("sales_order_detail.status = ?", query.Status)
	}
	cursor.Joins("Product").Order("sales_order_detail.id DESC").Find(&details)
	return details
}

type ProductPurchaseOrderDetailsQuery struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	ProductId int32      `json:"productId"`
}

func (q *ProductPurchaseOrderDetailsQuery) isValid() bool {
	return !(q.ProductId <= 0)
}

// Get the purchase order details with the product specified.
func getProductPurchaseOrderDetails(query ProductPurchaseOrderDetailsQuery, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	if !query.isValid() {
		return details
	}

	cursor := dbOrm.Model(&PurchaseOrderDetail{}).Where("purchase_order_detail.product = ? AND purchase_order_detail.enterprise = ?", query.ProductId, enterpriseId)
	if query.StartDate != nil {
		cursor.Where("(SELECT date_created FROM purchase_order WHERE purchase_order.id = purchase_order_detail.\"order\") >= ?", query.StartDate)
	}
	if query.EndDate != nil {
		cursor.Where("(SELECT date_created FROM purchase_order WHERE purchase_order.id = purchase_order_detail.\"order\") >= ?", query.EndDate)
	}
	cursor.Joins("Product").Order("purchase_order_detail.id DESC").Find(&details)
	return details
}

// Get the warehouse movements with the product specified.
func getProductWarehouseMovement(query ProductPurchaseOrderDetailsQuery, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	if !query.isValid() {
		return warehouseMovements
	}

	cursor := dbOrm.Model(&WarehouseMovement{}).Where("warehouse_movement.product = ? AND warehouse_movement.enterprise = ?", query.ProductId, enterpriseId)
	if query.StartDate != nil {
		cursor.Where("warehouse_movement.date_created >= ?", query.StartDate)
	}
	if query.EndDate != nil {
		cursor.Where("warehouse_movement.date_created <= ?", query.EndDate)
	}
	cursor.Joins("Warehouse").Joins("Product").Order("warehouse_movement.id DESC").Find(&warehouseMovements)
	return warehouseMovements
}

type ProductManufacturingOrdersQuery struct {
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	Manufactured string     `json:"manufactured"` // Y = Yes, N = No, Empty to search all
	ProductId    int32      `json:"productId"`
}

func (q *ProductManufacturingOrdersQuery) isValid() bool {
	return !(q.ProductId <= 0 || (q.Manufactured != "" && q.Manufactured != "Y" && q.Manufactured != "N"))
}

// Get the manufacturing orders with the product specified.
func getProductManufacturingOrders(query ProductManufacturingOrdersQuery, enterpriseId int32) []ManufacturingOrder {
	manufacturingOrders := make([]ManufacturingOrder, 0)
	if !query.isValid() {
		return manufacturingOrders
	}

	cursor := dbOrm.Model(&ManufacturingOrder{}).Where("manufacturing_order.product = ? AND manufacturing_order.enterprise = ?", query.ProductId, enterpriseId)
	if query.StartDate != nil {
		cursor.Where("manufacturing_order.date_created >= ?", query.StartDate)
	}
	if query.EndDate != nil {
		cursor.Where("manufacturing_order.date_created <= ?", query.EndDate)
	}
	if query.Manufactured != "" {
		if query.Manufactured == "Y" {
			cursor.Where("manufacturing_order.manufactured = ?", true)
		} else {
			cursor.Where("manufacturing_order.manufactured = ?", false)
		}
	}
	cursor.Joins("Product").Joins("Type").Order("manufacturing_order.date_created DESC").Find(&manufacturingOrders)
	return manufacturingOrders
}

// Get the complex manufacturing orders with the product specified.
func getProductComplexManufacturingOrders(query ProductManufacturingOrdersQuery, enterpriseId int32) []ComplexManufacturingOrder {
	complexManufacturingOrders := make([]ComplexManufacturingOrder, 0)
	if !query.isValid() {
		return complexManufacturingOrders
	}

	cursor := dbOrm.Model(&ComplexManufacturingOrder{}).Joins("INNER JOIN complex_manufacturing_order_manufacturing_order ON complex_manufacturing_order_manufacturing_order.complex_manufacturing_order=complex_manufacturing_order.id").Where("complex_manufacturing_order_manufacturing_order.product = ? AND complex_manufacturing_order.enterprise = ?", query.ProductId, enterpriseId)
	if query.StartDate != nil {
		cursor.Where("complex_manufacturing_order.date_created >= ?", query.StartDate)
	}
	if query.EndDate != nil {
		cursor.Where("complex_manufacturing_order.date_created <= ?", query.EndDate)
	}
	if query.Manufactured != "" {
		if query.Manufactured == "Y" {
			cursor.Where("complex_manufacturing_order.manufactured = ?", true)
		} else {
			cursor.Where("complex_manufacturing_order.manufactured = ?", false)
		}
	}
	cursor.Joins("Type").Order("complex_manufacturing_order.date_created DESC").Find(&complexManufacturingOrders)
	return complexManufacturingOrders
}

func (p *Product) generateBarcode(enterpriseId int32) bool {
	var product Product
	result := dbOrm.Model(&Product{}).Where("enterprise = @enterprise_id AND SUBSTRING(barcode,0,LENGTH(@barcode)+1) = @barcode", sql.Named("enterprise_id", enterpriseId), sql.Named("barcode", getSettingsRecordById(enterpriseId).BarcodePrefix)).Order("barcode DESC").Limit(1).First(&product)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	barcode := product.BarCode[0:12]

	code, err := strconv.Atoi(barcode)
	if err != nil {
		log("EAN13", err.Error())
		return false
	}
	code++
	barcode = strconv.Itoa(code)

	checkCode := 0
	for i := 1; i < len(barcode); i += 2 {
		j, _ := strconv.Atoi(barcode[i : i+1])
		checkCode += j
	}
	checkCode *= 3
	for i := 0; i < len(barcode); i += 2 {
		j, _ := strconv.Atoi(barcode[i : i+1])
		checkCode += j
	}

	if 10-(checkCode%10) == 10 {
		checkCode = 0
	} else {
		checkCode = 10 - (checkCode % 10)
	}
	p.BarCode = barcode + strconv.Itoa(checkCode)
	return true
}

type ProductImage struct {
	Id        int32   `json:"id"`
	ProductId int32   `json:"productId" gorm:"column:product;not null:true"`
	Product   Product `json:"-" gorm:"foreignKey:ProductId;references:Id"`
	URL       string  `json:"url" gorm:"type:character varying(255);not null:true"`
}

func (pi *ProductImage) TableName() string {
	return "product_image"
}

func getProductImages(productId int32, enterpriseId int32) []ProductImage {
	var image []ProductImage = make([]ProductImage, 0)
	product := getProductRow(productId)
	if product.EnterpriseId != enterpriseId {
		return image
	}

	dbOrm.Model(&ProductImage{}).Where("product = ?", productId).Order("id ASC").Find(&image)
	return image
}

func (i *ProductImage) isValid() bool {
	return !(len(i.URL) == 0 || len(i.URL) > 255)
}

func (pi *ProductImage) BeforeCreate(tx *gorm.DB) (err error) {
	var productImage ProductImage
	result := tx.Model(&ProductImage{}).Last(&productImage)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return result.Error
	}
	pi.Id = productImage.Id + 1
	return nil
}

func (i *ProductImage) insertProductImage(enterpriseId int32) bool {
	if !i.isValid() || i.ProductId <= 0 {
		return false
	}

	p := getProductRow(i.ProductId)
	if p.EnterpriseId != enterpriseId {
		return false
	}

	result := dbOrm.Create(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (i *ProductImage) updateProductImage(enterpriseId int32) bool {
	if i.Id <= 0 || !i.isValid() {
		return false
	}

	var productImage ProductImage
	result := dbOrm.Model(&ProductImage{}).Where("id = ?", i.Id).First(&productImage)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	p := getProductRow(productImage.ProductId)
	if p.EnterpriseId != enterpriseId {
		return false
	}

	productImage.URL = i.URL

	result = dbOrm.Save(&productImage)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (i *ProductImage) deleteProductImage(enterpriseId int32) bool {
	if i.Id <= 0 {
		return false
	}

	var productImage ProductImage
	result := dbOrm.Model(&ProductImage{}).Where("id = ?", i.Id).First(&productImage)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	p := getProductRow(productImage.ProductId)
	if p.EnterpriseId != enterpriseId {
		return false
	}

	result = dbOrm.Where("id = ?", i.Id).Delete(&ProductImage{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func calculateMinimumStock(enterpriseId int32, userId int32) bool {
	s := getSettingsRecordById(enterpriseId)
	t := time.Now()
	if s.MinimumStockSalesPeriods <= 0 || s.MinimumStockSalesDays <= 0 {
		return false
	}
	t = t.AddDate(0, 0, -int(s.MinimumStockSalesDays))

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	var products []Product = make([]Product, 0)
	dbOrm.Model(&Product{}).Where("track_minimum_stock = true AND enterprise = ? AND off = false", enterpriseId).Find(&products)

	for i := 0; i < len(products); i++ {
		product := products[i]

		var quantitySold int32
		result := dbOrm.Model(&Product{}).Joins("INNER JOIN sales_order ON sales_order.id=sales_order_detail.order").Where("sales_order_detail.product = ? AND sales_order.date_created >= ?", product.Id, t).Select("SUM(sales_order_detail.quantity) AS quantity").Pluck("quantity", &quantitySold)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false
		}

		product.MinimumStock = quantitySold / int32(s.MinimumStockSalesPeriods)

		result = dbOrm.Save(&product)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false
		}

		insertTransactionalLog(enterpriseId, "product", int(product.Id), userId, "U")
		json, _ := json.Marshal(product)
		go fireWebHook(product.EnterpriseId, "product", "PUT", string(json))
	}

	///
	trans.Commit()
	return true
	///
}

type GenerateManufacturingOrPurchaseOrdersMinimumStock struct {
	Warehouse string `json:"warehouse"`
}

func (g *GenerateManufacturingOrPurchaseOrdersMinimumStock) generateManufacturingOrPurchaseOrdersMinimumStock(userId int32, enterpriseId int32) bool {
	if len(g.Warehouse) == 0 {
		s := getSettingsRecordById(enterpriseId)
		g.Warehouse = s.DefaultWarehouseId
	}
	var generadedPurchaseOrders map[int32]PurchaseOrder = make(map[int32]PurchaseOrder) // Key: supplier ID, Value: generated purchase order

	rows, err := dbOrm.Model(&Product{}).Joins("INNER JOIN stock ON stock.product=product.id").Where("product.track_minimum_stock = true AND stock.quantity_available < (product.minimum_stock * 2) AND product.enterprise = ? AND product.off = false", enterpriseId).Select("product.id, stock.quantity_available ,product.minimum_stock ,product.manufacturing ,product.manufacturing_order_type ,product.supplier").Rows()
	if err != nil {
		log("DB", err.Error())
		return false
	}
	defer rows.Close()

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	// iterate over the list of product that don't have covered the minimum stock
	for rows.Next() {
		var productId int32
		var quantityAvailable int32
		var minimumStock int32
		var manufacturing bool
		var manufacturingOrderType *int32
		var supplier *int32
		rows.Scan(&productId, &quantityAvailable, &minimumStock, &manufacturing, &manufacturingOrderType, &supplier)

		if manufacturing { // if the product is from manufacture, generate the manufacturing orders
			// generate manufacturing order or purchase orders until the available quantity is equal to the minimum stock * 2
			for i := quantityAvailable; i < (minimumStock * 2); i++ {

				o := ManufacturingOrder{ProductId: productId, TypeId: *manufacturingOrderType}
				o.UserCreatedId = userId
				o.EnterpriseId = enterpriseId
				o.WarehouseId = g.Warehouse
				ok := o.insertManufacturingOrder(userId, trans).Ok
				if !ok {
					trans.Rollback()
					return false
				}
			}
		} else { // if the product is not from manufacture, generate the purchase order to the supplier
			o, ok := generadedPurchaseOrders[*supplier]
			if !ok { // there is no purchase order generated for this supplier, create it and add to the map
				d := getSupplierDefaults(*supplier, enterpriseId)
				if d.BillingSeries == nil || d.Currency == nil || d.MainBillingAddress == nil || d.MainShippingAddress == nil || d.PaymentMethod == nil {
					continue
				}
				p := PurchaseOrder{}
				p.SupplierId = *supplier
				p.BillingSeriesId = *d.BillingSeries
				p.CurrencyId = *d.Currency
				p.BillingAddressId = *d.MainBillingAddress
				p.ShippingAddressId = *d.MainShippingAddress
				p.PaymentMethodId = *d.PaymentMethod

				p.EnterpriseId = enterpriseId
				ok, purchaseOrderId := p.insertPurchaseOrder(userId, trans)
				if !ok {
					trans.Rollback()
					return false
				}
				p.Id = purchaseOrderId
				generadedPurchaseOrders[*supplier] = p

				// generate the needs as a detail
				product := getProductRow(productId)
				det := PurchaseOrderDetail{
					OrderId:      p.Id,
					ProductId:    productId,
					Quantity:     (minimumStock * 2) - quantityAvailable,
					Price:        product.Price,
					VatPercent:   product.VatPercent,
					EnterpriseId: enterpriseId,
					WarehouseId:  g.Warehouse,
				}
				okAndErr, _ := det.insertPurchaseOrderDetail(userId, trans)
				if !okAndErr.Ok {
					trans.Rollback()
					return false
				}
			} else { // it already exists a purchase order for this supplier, add the needs as details
				// generate the needs as a detail
				product := getProductRow(productId)
				det := PurchaseOrderDetail{
					OrderId:      o.Id,
					ProductId:    productId,
					Quantity:     (minimumStock * 2) - quantityAvailable,
					Price:        product.Price,
					VatPercent:   product.VatPercent,
					EnterpriseId: enterpriseId,
					WarehouseId:  g.Warehouse,
				}
				okAndErr, _ := det.insertPurchaseOrderDetail(userId, trans)
				if !okAndErr.Ok {
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

type ProductLocate struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

type ProductLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name, 2 = Reference
	Value string `json:"value"`
}

func (q *ProductLocateQuery) locateProduct(enterpriseId int32) []ProductLocate {
	var products []ProductLocate = make([]ProductLocate, 0)
	query := ``
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		query = `enterprise = ? AND off = false`
		parameters = append(parameters, enterpriseId)
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			query = `enterprise = ? AND off = false`
			parameters = append(parameters, enterpriseId)
		} else {
			query = `id = ? AND enterprise = ? AND off = false`
			parameters = append(parameters, id)
			parameters = append(parameters, enterpriseId)
		}
	} else if q.Mode == 1 {
		query = `name ILIKE ? AND enterprise = ? AND off = false`
		parameters = append(parameters, "%"+q.Value+"%")
		parameters = append(parameters, enterpriseId)
	} else if q.Mode == 2 {
		query = `reference ILIKE ? AND enterprise = ? AND off = false`
		parameters = append(parameters, "%"+q.Value+"%")
		parameters = append(parameters, enterpriseId)
	}
	dbOrm.Model(&Product{}).Where(query, parameters...).Order("id ASC").Find(&products)
	return products
}

/* PRODUCT GENERATOR */

type ProductGenerator struct {
	Products                   []ProductGenerate `json:"products"`
	ManufacturingOrderTypeMode int16             `json:"manufacturingOrderTypeMode"` // 0 = No manufacturing, 1 = Create a manufacturing order type with the same name, 2 = Use the manufacturing order name
	ManufacturingOrderTypeName *string           `json:"manufacturingOrderTypeName"`
}

type ProductGenerate struct {
	Name            string  `json:"name"`
	Reference       string  `json:"reference"`
	GenerateBarCode bool    `json:"generateBarCode"`
	BarCode         string  `json:"barCode"`
	Weight          float64 `json:"weight"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	Depth           float64 `json:"depth"`
	Price           float64 `json:"price"`
	Manufacturing   bool    `json:"manufacturing"`
	InitialStock    int32   `json:"initialStock"`
}

func (g *ProductGenerator) productGenerator(enterpriseId int32, userId int32) bool {
	var manufacturingOrderTypeId int32
	if g.ManufacturingOrderTypeMode == 2 {
		if g.ManufacturingOrderTypeName == nil {
			return false
		}

		mot := ManufacturingOrderType{
			Name:         *g.ManufacturingOrderTypeName,
			EnterpriseId: enterpriseId,
		}
		mot.insertManufacturingOrderType()
		manufacturingOrderTypeId = mot.Id
	}

	for i := 0; i < len(g.Products); i++ {
		product := g.Products[i]

		p := Product{
			Name:          product.Name,
			Reference:     product.Reference,
			BarCode:       product.BarCode,
			Weight:        product.Weight,
			Width:         product.Width,
			Height:        product.Height,
			Depth:         product.Depth,
			Price:         product.Price,
			Manufacturing: product.Manufacturing,
			EnterpriseId:  enterpriseId,
		}

		if product.Manufacturing && g.ManufacturingOrderTypeMode == 1 {
			mot := ManufacturingOrderType{
				Name:                 product.Name,
				EnterpriseId:         enterpriseId,
				QuantityManufactured: 1,
			}
			mot.insertManufacturingOrderType()
			p.ManufacturingOrderTypeId = &mot.Id
		} else if product.Manufacturing && g.ManufacturingOrderTypeMode == 2 {
			p.ManufacturingOrderTypeId = &manufacturingOrderTypeId
		}

		ok := p.insertProduct(userId).Ok
		if !ok {
			return false
		}

		if product.GenerateBarCode {
			p := getProductRow(p.Id)
			p.generateBarcode(enterpriseId)
			p.updateProduct(userId)
		}

		if product.InitialStock != 0 {
			s := getSettingsRecordById(enterpriseId)
			wm := WarehouseMovement{
				WarehouseId:  s.DefaultWarehouseId,
				ProductId:    p.Id,
				Quantity:     product.InitialStock,
				Type:         "R",
				EnterpriseId: enterpriseId,
			}
			wm.insertWarehouseMovement(userId, nil)
		}
	}
	return true
}

func getProductsByManufacturingOrderType(manufacturingOrderTypeId int32, enterpriseId int32) []Product {
	var products []Product = make([]Product, 0)
	dbOrm.Model(&Product{}).Where("enterprise = ? AND manufacturing_order_type = ?", enterpriseId, manufacturingOrderTypeId).Order("id ASC").Find(&products)
	return products
}
