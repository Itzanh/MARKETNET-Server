/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "gorm.io/gorm"

type ManufacturingOrderType struct {
	Id                   int32    `json:"id" gorm:"index:manufacturing_order_type_id_enterprise,unique:true,priority:1"`
	Name                 string   `json:"name" gorm:"type:character varying(100);not null"`
	EnterpriseId         int32    `json:"-" gorm:"column:enterprise;not null:true;index:manufacturing_order_type_id_enterprise,unique:true,priority:2"`
	Enterprise           Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	QuantityManufactured int32    `json:"quantityManufactured" gorm:"not null:true"`
	Complex              bool     `json:"complex" gorm:"not null:true"`
}

func (t *ManufacturingOrderType) TableName() string {
	return "manufacturing_order_type"
}

func getManufacturingOrderType(enterpriseId int32) []ManufacturingOrderType {
	var types []ManufacturingOrderType = make([]ManufacturingOrderType, 0)
	dbOrm.Model(&ManufacturingOrderType{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&types)
	return types
}

func getManufacturingOrderTypeRow(typeId int32) ManufacturingOrderType {
	t := ManufacturingOrderType{}
	dbOrm.Model(&ManufacturingOrderType{}).Where("id = ?", typeId).First(&t)
	return t
}

func (t *ManufacturingOrderType) isValid() bool {
	return !(len(t.Name) == 0 || len(t.Name) > 100 || t.QuantityManufactured < 1)
}

func (t *ManufacturingOrderType) BeforeCreate(tx *gorm.DB) (err error) {
	var manufacturingOrderType ManufacturingOrderType
	tx.Model(&ManufacturingOrderType{}).Last(&manufacturingOrderType)
	t.Id = manufacturingOrderType.Id + 1
	return nil
}

func (t *ManufacturingOrderType) insertManufacturingOrderType() bool {
	if !t.isValid() {
		return false
	}

	result := dbOrm.Create(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (t *ManufacturingOrderType) updateManufacturingOrderType() bool {
	if t.Id <= 0 || !t.isValid() {
		return false
	}

	var manufacturingOrderType ManufacturingOrderType
	result := dbOrm.Where("id = ? AND enterprise = ?", t.Id, t.EnterpriseId).First(&manufacturingOrderType)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	manufacturingOrderType.Name = t.Name
	if manufacturingOrderType.Complex {
		manufacturingOrderType.QuantityManufactured = 0
	} else {
		manufacturingOrderType.QuantityManufactured = t.QuantityManufactured
	}

	result = dbOrm.Save(&manufacturingOrderType)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (t *ManufacturingOrderType) deleteManufacturingOrderType() bool {
	if t.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", t.Id, t.EnterpriseId).Delete(&ManufacturingOrderType{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
