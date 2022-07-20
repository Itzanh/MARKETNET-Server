/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeprecatedProducts struct {
	Id            int64      `json:"id" gorm:"primaryKey;index:deprecated_products_id_enterprise,unique:true"`
	ProductId     int32      `json:"productId" gorm:"column:product;not null:true"`
	Product       Product    `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated   time.Time  `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	DateDrop      time.Time  `json:"dateDrop" gorm:"type:timestamp(3) with time zone;not null:true"`
	DateDropped   *time.Time `json:"dateDropped" gorm:"type:timestamp(3) with time zone"`
	Dropped       bool       `json:"dropped"`
	Reason        string     `json:"reason" gorm:"type:character varying(255);not null:true"`
	UserCreatedId int32      `json:"userCreatedId" gorm:"column:user_created;not null:true"`
	UserCreated   User       `json:"userCreated" gorm:"foreignKey:UserCreatedId,EnterpriseId;references:Id,EnterpriseId"`
	TotalChecks   int16      `json:"totalChecks" gorm:"type:smallint;not null:true"`
	ChecksDone    int16      `json:"checksDone" gorm:"type:smallint;not null:true"`
	EnterpriseId  int32      `json:"-" gorm:"column:enterprise;not null:true;index:deprecated_products_id_enterprise,unique:true"`
	Enterprise    Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (DeprecatedProducts) TableName() string {
	return "deprecated_products"
}

type DeprecatedProductsQuery struct {
	Dropped     *bool  `json:"dropped"`
	ProductName string `json:"productName"`
}

func (q *DeprecatedProductsQuery) searchDeprecatedProducts(enterpriseId int32) []DeprecatedProducts {
	var deprecatedProducts []DeprecatedProducts = make([]DeprecatedProducts, 0)

	cursor := dbOrm.Model(&deprecatedProducts).Where("enterprise = ?", enterpriseId)
	if q.Dropped != nil {
		cursor.Where("dropped = ?", q.Dropped)
	}
	if q.ProductName != "" {
		cursor.Where("product.name ILIKE ?", "%"+q.ProductName+"%")
	}
	result := cursor.Preload(clause.Associations).Order("date_created DESC").Find(&deprecatedProducts)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return deprecatedProducts
}

// For internal use only
func getDeprecatedProductRow(deprecatedProductId int64) DeprecatedProducts {
	dp := DeprecatedProducts{}
	result := dbOrm.Model(&DeprecatedProducts{}).Where("id = ?", deprecatedProductId).First(&dp)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return dp
}

func (d *DeprecatedProducts) isValid() bool {
	return !(d.ProductId <= 0 || len(d.Reason) == 0 || len(d.Reason) > 255 || d.DateDrop.IsZero())
}

func (d *DeprecatedProducts) insertDeprecatedProduct(userId int32) bool {
	if !d.isValid() {
		return false
	}

	d.DateCreated = time.Now()
	d.DateDropped = nil
	d.Dropped = false
	d.UserCreatedId = userId

	result := dbOrm.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (d *DeprecatedProducts) dropDeprecatedProduct() bool {
	if d.Id <= 0 {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(d.Id)
	if deprecatedProduct.Id <= 0 || deprecatedProduct.EnterpriseId != d.EnterpriseId || deprecatedProduct.Dropped {
		return false
	}

	// check if the product is still in use
	uses := calcDeprecatedProductUses(d.Id, d.EnterpriseId)
	if !uses.NoUses {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	// deactivate the product
	result := trans.Model(&Product{}).Where("id = ?", deprecatedProduct.ProductId).Update("off", true)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	deprecatedProduct.Dropped = true
	now := time.Now()
	deprecatedProduct.DateDropped = &now

	result = trans.Save(&deprecatedProduct)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///

	return true
}

func (d *DeprecatedProducts) deleteDeprecatedProduct() bool {
	if d.Id <= 0 {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(d.Id)
	if deprecatedProduct.Id <= 0 || deprecatedProduct.EnterpriseId != d.EnterpriseId || deprecatedProduct.Dropped {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	result := trans.Delete(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ?", d.Id, d.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Delete(&DeprecatedProducts{}).Where("id = ? AND enterprise = ?", d.Id, d.EnterpriseId).Delete(&DeprecatedProducts{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///

	return true
}

type DeprecatedProductCheckList struct {
	Id                  int64              `json:"id" gorm:"priamryKey;index:deprecated_product_check_list_id_enterprise,unique:true"`
	DeprecatedProductId int64              `json:"deprecatedProductId" gorm:"column:deprecated_product;not null:true;index:deprecated_product_check_list_dp_position,unique:true"`
	DeprecatedProduct   DeprecatedProducts `json:"-" gorm:"foreignKey:DeprecatedProductId,EnterpriseId;references:Id,EnterpriseId"`
	Text                string             `json:"text" gorm:"type:character varying(255);not null:true"`
	Position            int32              `json:"position" gorm:"not null:true;index:deprecated_product_check_list_dp_position,unique:true"`
	Checked             bool               `json:"checked" gorm:"not null:true"`
	DateChecked         *time.Time         `json:"dateChecked" gorm:"type:timestamp(0) with time zone"`
	UserCheckedId       *int32             `json:"userCheckedId" gorm:"column:user_checked"`
	UserChecked         *User              `json:"userChecked" gorm:"foreignKey:UserCheckedId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId        int32              `json:"-" gorm:"column:enterprise;not null:true;index:deprecated_product_check_list_id_enterprise,unique:true"`
	Enterprise          Settings           `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (DeprecatedProductCheckList) TableName() string {
	return "deprecated_product_check_list"
}

func getDeprecatedProductCheckList(deprecatedProductId int64, enterpriseId int32) []DeprecatedProductCheckList {
	var deprecatedProductCheckList []DeprecatedProductCheckList = make([]DeprecatedProductCheckList, 0)

	result := dbOrm.Model(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ?", deprecatedProductId, enterpriseId).Order("position ASC").Preload("UserChecked").Find(&deprecatedProductCheckList)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return deprecatedProductCheckList
}

// For internal use only
func getDeprecatedProductCheckListRow(deprecatedProductId int64, enterpriseId int32) *DeprecatedProductCheckList {
	dp := &DeprecatedProductCheckList{}
	result := dbOrm.Model(&DeprecatedProductCheckList{}).Where("id = ? AND enterprise = ?", deprecatedProductId, enterpriseId).First(&dp)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return dp
}

func (c *DeprecatedProductCheckList) isValid() bool {
	return !(c.DeprecatedProductId <= 0 || len(c.Text) == 0 || len(c.Text) > 255)
}

func (d *DeprecatedProducts) getNextPosition() int32 {
	var lastPosition int32 = 0

	result := dbOrm.Model(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ?", d.Id, d.EnterpriseId).Order("position DESC").Select("position").Limit(1).Pluck("position", &lastPosition)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return lastPosition + 1
}

func (c *DeprecatedProductCheckList) BeforeCreate(tx *gorm.DB) (err error) {
	var deprecatedProductCheckList DeprecatedProductCheckList
	tx.Model(&DeprecatedProductCheckList{}).Last(&deprecatedProductCheckList)
	c.Id = deprecatedProductCheckList.Id + 1
	return nil
}

func (c *DeprecatedProductCheckList) insertDeprecatedProductCheckList() bool {
	if !c.isValid() {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(c.DeprecatedProductId)
	if deprecatedProduct.Dropped {
		return false
	}

	c.Position = deprecatedProduct.getNextPosition()
	c.Checked = false
	c.DateChecked = nil
	c.UserCheckedId = nil

	///
	trans := dbOrm.Begin()
	///

	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&DeprecatedProducts{}).Where("id = ?", c.DeprecatedProductId).Updates(map[string]interface{}{
		"total_checks": deprecatedProduct.TotalChecks + 1,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///
	return true
}

func (c *DeprecatedProductCheckList) toggleDeprecatedProductCheckList(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	c = getDeprecatedProductCheckListRow(c.Id, c.EnterpriseId)
	if c == nil {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(c.DeprecatedProductId)
	if deprecatedProduct.Dropped {
		return false
	}

	c.Checked = !c.Checked
	if c.Checked {
		c.UserCheckedId = &userId
		now := time.Now()
		c.DateChecked = &now
	} else {
		c.UserCheckedId = nil
		c.DateChecked = nil
	}

	///
	trans := dbOrm.Begin()
	///

	result := trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", c.Id).Save(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	if c.Checked {
		deprecatedProduct.ChecksDone++
	} else {
		deprecatedProduct.ChecksDone--
	}
	result = trans.Model(&DeprecatedProducts{}).Where("id = ?", c.DeprecatedProductId).Updates(map[string]interface{}{
		"checks_done": deprecatedProduct.ChecksDone,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///
	return true
}

func (c *DeprecatedProductCheckList) deleteDeprecatedProductCheckList() bool {
	if c.Id <= 0 {
		return false
	}

	// verifications
	c = getDeprecatedProductCheckListRow(c.Id, c.EnterpriseId)
	if c == nil {
		return false
	}

	if c.Checked {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(c.DeprecatedProductId)
	if deprecatedProduct.Dropped {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	// delete
	result := trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", c.Id).Delete(&DeprecatedProductCheckList{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	result = trans.Model(&DeprecatedProducts{}).Where("id = ?", c.DeprecatedProductId).Updates(map[string]interface{}{
		"total_checks": deprecatedProduct.TotalChecks - 1,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// re-enumerate positions in the list
	var deprecatedProductCheckList []DeprecatedProductCheckList = make([]DeprecatedProductCheckList, 0)
	result = trans.Model(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ?", c.DeprecatedProductId, c.EnterpriseId).Order("position ASC").Find(&deprecatedProductCheckList)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	for i := 0; i < len(deprecatedProductCheckList); i++ {
		result = trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", deprecatedProductCheckList[i].Id).Update("position", i+1)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	///

	return true
}

type DeprecatedProductCheckListChangePosition struct {
	Id                  int64 `json:"id"`
	DeprecatedProductId int64 `json:"deprecatedProductId"`
	NewPosition         int32 `json:"newPosition"`
}

func (p *DeprecatedProductCheckListChangePosition) isValid() bool {
	return !(p.Id <= 0 || p.DeprecatedProductId <= 0 || p.NewPosition <= 0)
}

// This funcion is used to swap the positions of the list, so the checks can be sorted from the front-end once they are created
func (p *DeprecatedProductCheckListChangePosition) movePositionDeprecatedProductCheckList(enterpriseId int32) bool {
	if !p.isValid() {
		return false
	}

	deprecatedProduct := getDeprecatedProductRow(p.DeprecatedProductId)
	if deprecatedProduct.Dropped || deprecatedProduct.EnterpriseId != enterpriseId {
		return false
	}

	oldDeprecatedProductCheckList := getDeprecatedProductCheckListRow(p.Id, enterpriseId)
	if oldDeprecatedProductCheckList == nil {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	// set the position of the check to be changed to 0, so it doesn't collide when the other positions are slided
	result := trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", p.Id).Update("position", 0)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	var maxPosition int32
	result = trans.Model(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ?", p.DeprecatedProductId, enterpriseId).Order("position DESC").Limit(1).Select("position").Pluck("position", &maxPosition)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	// if there are no more positions below this one, don't allow to go down any deeper
	if p.NewPosition > maxPosition {
		trans.Rollback()
		return false
	}

	// re-enumerate positions in the list, leaving a blank position for the changed check
	var deprecatedProductCheckList []DeprecatedProductCheckList = make([]DeprecatedProductCheckList, 0)
	result = trans.Model(&DeprecatedProductCheckList{}).Where("deprecated_product = ? AND enterprise = ? AND position > 0", p.DeprecatedProductId, enterpriseId).Order("position ASC").Find(&deprecatedProductCheckList)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	var currentPosition int32 = 0
	for i := 0; i < len(deprecatedProductCheckList); i++ {
		currentPosition++
		if currentPosition == p.NewPosition {
			currentPosition++
		}
		result = trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", deprecatedProductCheckList[i].Id).Update("position", currentPosition)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	// set the new position
	result = trans.Model(&DeprecatedProductCheckList{}).Where("id = ?", p.Id).Update("position", p.NewPosition)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///

	return true
}

type DeprecatedProductInUse struct {
	SalesOrders                        []DeprecatedProductInUseOrders `json:"salesOrders"`
	PurchaseOrders                     []DeprecatedProductInUseOrders `json:"purchaseOrder"`
	ManufacturingOrdersQuantity        int64                          `json:"manufacturingOrdersQuantity"`
	ComplexManufacturingOrdersQuantity int64                          `json:"complexManufacturingOrdersQuantity"`
	UnitsInStock                       int32                          `json:"unitsInStock"`
	NoUses                             bool                           `json:"noUses"`
}

type DeprecatedProductInUseOrders struct {
	OrderName string `json:"orderName"`
	Quantity  int32  `json:"quantity"`
}

func (d *DeprecatedProductInUse) deprecatedProductHasNoUses() bool {
	return len(d.SalesOrders) == 0 && len(d.PurchaseOrders) == 0 && d.ManufacturingOrdersQuantity == 0 && d.ComplexManufacturingOrdersQuantity == 0
}

func calcDeprecatedProductUses(deprecatedProductId int64, enterpriseId int32) DeprecatedProductInUse {
	usage := DeprecatedProductInUse{}
	usage.SalesOrders = make([]DeprecatedProductInUseOrders, 0)
	usage.PurchaseOrders = make([]DeprecatedProductInUseOrders, 0)
	deprecatedProduct := getDeprecatedProductRow(deprecatedProductId)
	if deprecatedProduct.EnterpriseId != enterpriseId {
		return usage
	}

	// get the sales order details
	result := dbOrm.Model(&SalesOrderDetail{}).Where(`product = ? AND "sales_order_detail".enterprise = ? AND quantity_delivery_note < quantity`, deprecatedProduct.ProductId, deprecatedProduct.EnterpriseId).Joins("Order").Select(`"Order"."order_name", quantity`).Scan(&usage.SalesOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return usage
	}

	// get the purchases order details
	result = dbOrm.Model(&PurchaseOrderDetail{}).Where(`product = ? AND "purchase_order_detail".enterprise = ? AND quantity_delivery_note < quantity`, deprecatedProduct.ProductId, deprecatedProduct.EnterpriseId).Joins("Order").Select(`"Order"."order_name", quantity`).Scan(&usage.PurchaseOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return usage
	}

	// get the manufacturing orders quantity
	result = dbOrm.Model(&ManufacturingOrder{}).Where("product = ? AND enterprise = ? AND NOT manufactured", deprecatedProduct.ProductId, deprecatedProduct.EnterpriseId).Count(&usage.ManufacturingOrdersQuantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return usage
	}

	// get the complex manufacturing orders quantity
	result = dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("product = ? AND enterprise = ? AND NOT manufactured", deprecatedProduct.ProductId, deprecatedProduct.EnterpriseId).Count(&usage.ComplexManufacturingOrdersQuantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return usage
	}

	// get the units in stock
	product := getProductRow(deprecatedProduct.ProductId)
	usage.UnitsInStock = product.Stock

	usage.NoUses = usage.deprecatedProductHasNoUses()
	return usage
}
