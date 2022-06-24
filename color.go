/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"strings"

	"gorm.io/gorm"
)

type Color struct {
	Id           int32    `json:"id" gorm:"index:color_id_enterprise,unique:true,priority:1"`
	Name         string   `json:"name" gorm:"type:character varying(50);not null:true"`
	HexColor     string   `json:"hexColor" gorm:"type:character(6);not null:true"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:color_id_enterprise,unique:true,priority:2"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Color) TableName() string {
	return "color"
}

func getColor(enterpriseId int32) []Color {
	var color []Color = make([]Color, 0)
	dbOrm.Model(&Color{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&color)
	return color
}

func (c *Color) isValid() bool {
	c.HexColor = strings.ToUpper(c.HexColor)
	return !(len(c.Name) == 0 || len(c.Name) > 100 || len(c.HexColor) != 6 || !checkHex(c.HexColor))
}

func (c *Color) BeforeCreate(tx *gorm.DB) (err error) {
	var color Color
	tx.Model(&Color{}).Last(&color)
	c.Id = color.Id + 1
	return nil
}

func (c *Color) insertColor() bool {
	if !c.isValid() {
		return false
	}

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Color) updateColor() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	var color Color
	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&color)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	color.Name = c.Name
	color.HexColor = c.HexColor

	result = dbOrm.Save(&color)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Color) deleteColor() bool {
	if c.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Color{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func locateColor(enterpriseId int32) []NameInt32 {
	var color []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&Color{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&color)
	return color
}
