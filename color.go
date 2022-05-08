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
	return !(len(c.Name) == 0 || len(c.Name) > 100 || len(c.HexColor) > 6)
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

func findColorByName(colorName string, enterpriseId int32) []NameInt32 {
	var colors []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&Color{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(colorName), enterpriseId).Order("id ASC").Limit(10).Find(&colors)
	return colors
}

func getNameColor(id int32, enterpriseId int32) string {
	var color Color
	result := dbOrm.Where("id = ? AND enterprise = ?", id, enterpriseId).First(&color)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	return color.Name
}

func locateColor(enterpriseId int32) []NameInt32 {
	var color []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&Color{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&color)
	return color
}
