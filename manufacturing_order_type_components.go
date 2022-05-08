package main

import (
	"gorm.io/gorm"
)

type ManufacturingOrderTypeComponents struct {
	Id                       int32                  `json:"id" gorm:"index:manufacturing_order_type_components_id_enterprise,unique:true,priority:1"`
	ManufacturingOrderTypeId int32                  `json:"manufacturingOrderTypeId" gorm:"column:manufacturing_order_type;not null:true;index:manufacturing_order_type_components_component,unique:true,priority:1;index:manufacturing_order_type_components_manufacturing_order_type_ty,unique:true,priority:1"`
	ManufacturingOrderType   ManufacturingOrderType `json:"manufacturingOrderType" gorm:"foreignKey:ManufacturingOrderTypeId,EnterpriseId;references:Id,EnterpriseId"`
	Type                     string                 `json:"type" gorm:"type:character(1);not null:true;index:manufacturing_order_type_components_manufacturing_order_type_ty,unique:true,priority:2"` // I = Input, O = Output
	ProductId                int32                  `json:"productId" gorm:"column:product;not null:true;index:manufacturing_order_type_components_component,unique:true,priority:2;index:manufacturing_order_type_components_manufacturing_order_type_ty,unique:true,priority:3"`
	Product                  Product                `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity                 int32                  `json:"quantity" gorm:"not null:true"`
	EnterpriseId             int32                  `json:"-" gorm:"column:enterprise;not null:true;index:manufacturing_order_type_components_id_enterprise,unique:true,priority:2"`
	Enterprise               Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *ManufacturingOrderTypeComponents) TableName() string {
	return "manufacturing_order_type_components"
}

func getManufacturingOrderTypeComponents(manfuacturingOrderTypeId int32, enterpriserId int32) []ManufacturingOrderTypeComponents {
	var components []ManufacturingOrderTypeComponents = make([]ManufacturingOrderTypeComponents, 0)
	manufacturingOrderType := getManufacturingOrderTypeRow(manfuacturingOrderTypeId)
	if manufacturingOrderType.EnterpriseId != enterpriserId {
		return components
	}

	dbOrm.Model(&ManufacturingOrderTypeComponents{}).Where("manufacturing_order_type_components.manufacturing_order_type = ?", manfuacturingOrderTypeId).Joins("ManufacturingOrderType").Joins("Product").Order("manufacturing_order_type_components.product ASC").Find(&components)
	return components
}

func getManufacturingOrderTypeComponentRow(manfuacturingOrderTypeId int32) ManufacturingOrderTypeComponents {
	c := ManufacturingOrderTypeComponents{}
	dbOrm.Model(&ManufacturingOrderTypeComponents{}).Where("id = ?", manfuacturingOrderTypeId).First(&c)
	return c
}

func getManufacturingOrderTypeComponentRowTransaction(manfuacturingOrderTypeId int32, trans gorm.DB) ManufacturingOrderTypeComponents {
	c := ManufacturingOrderTypeComponents{}
	trans.Model(&ManufacturingOrderTypeComponents{}).Where("id = ?", manfuacturingOrderTypeId).First(&c)
	return c
}

// returns:
// ok
// code
// 0 = parameter error / ok
// 1 = the input product has the same manufacturing order type as the component
// 2 = the output product doesn't have the same manufacturing order type as the component
// 3 = the product already exist in one of the components
func (c *ManufacturingOrderTypeComponents) isValid() (bool, uint8) {
	if c.ProductId <= 0 {
		return false, 0
	}
	// the manufacturing order type has to be the same as this one for the output, and different on the input to make sure that there are no recursivity errors
	product := getProductRow(c.ProductId)
	if product.Id <= 0 {
		return false, 0
	}
	if c.Type == "I" {
		if product.ManufacturingOrderTypeId != nil && *product.ManufacturingOrderTypeId == c.ManufacturingOrderTypeId {
			return false, 1
		}
	} else if c.Type == "O" {
		if product.ManufacturingOrderTypeId == nil || *product.ManufacturingOrderTypeId != c.ManufacturingOrderTypeId {
			return false, 2
		}
	} else {
		return false, 0
	}

	if c.Id > 0 { // update
		// check that the product has not been associated yet
		components := getManufacturingOrderTypeComponents(c.ManufacturingOrderTypeId, c.EnterpriseId)
		for i := 0; i < len(components); i++ {
			if components[i].Id != c.Id && components[i].ProductId == c.ProductId {
				return false, 3
			}
		}
	} else { // insert
		// check that the product has not been associated yet
		components := getManufacturingOrderTypeComponents(c.ManufacturingOrderTypeId, c.EnterpriseId)
		for i := 0; i < len(components); i++ {
			if components[i].ProductId == c.ProductId {
				return false, 3
			}
		}
	}

	return !(c.ManufacturingOrderTypeId <= 0 || (c.Type != "I" && c.Type != "O") || c.Quantity <= 0), 0
}

func (c *ManufacturingOrderTypeComponents) BeforeCreate(tx *gorm.DB) (err error) {
	var manufacturingOrderTypeComponents ManufacturingOrderTypeComponents
	tx.Model(&ManufacturingOrderTypeComponents{}).Last(&manufacturingOrderTypeComponents)
	c.Id = manufacturingOrderTypeComponents.Id + 1
	return nil
}

func (c *ManufacturingOrderTypeComponents) insertManufacturingOrderTypeComponents() (bool, uint8) {
	ok, errorCode := c.isValid()
	if !ok {
		return false, errorCode
	}

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0
	}

	return true, 0
}

func (c *ManufacturingOrderTypeComponents) updateManufacturingOrderTypeComponents() (bool, uint8) {
	ok, errorCode := c.isValid()
	if c.Id <= 0 || !ok {
		return false, errorCode
	}

	var manufacturingOrderTypeComponents ManufacturingOrderTypeComponents
	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&manufacturingOrderTypeComponents)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0
	}

	manufacturingOrderTypeComponents.Type = c.Type
	manufacturingOrderTypeComponents.ProductId = c.ProductId
	manufacturingOrderTypeComponents.Quantity = c.Quantity

	result = dbOrm.Save(&manufacturingOrderTypeComponents)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0
	}

	return true, 0
}

func (c *ManufacturingOrderTypeComponents) deleteManufacturingOrderTypeComponents() bool {
	if c.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&ManufacturingOrderTypeComponents{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
