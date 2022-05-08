package main

import (
	"strings"
)

type Warehouse struct {
	Id           string   `json:"id" gorm:"primaryKey;type:character(2)"`
	Name         string   `json:"name" gorm:"column:name;type:character varying(50);not null:true"`
	EnterpriseId int32    `gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (w *Warehouse) TableName() string {
	return "warehouse"
}

func getWarehouses(enterpriseId int32) []Warehouse {
	var warehouses []Warehouse = make([]Warehouse, 0)
	result := dbOrm.Model(&Warehouse{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&warehouses)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return warehouses
	}

	return warehouses
}

func (w *Warehouse) isValid() bool {
	return !(len(w.Id) == 0 || len(w.Id) > 2 || len(w.Name) == 0 || len(w.Name) > 50)
}

func (w *Warehouse) insertWarehouse() bool {
	if !w.isValid() {
		return false
	}

	result := dbOrm.Create(&w)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (w *Warehouse) updateWarehouse() bool {
	if !w.isValid() {
		return false
	}

	var warehouse Warehouse
	result := dbOrm.Where("id = ? AND enterprise = ?", w.Id, w.EnterpriseId).First(&warehouse)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	warehouse.Name = w.Name

	result = dbOrm.Save(&warehouse)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (w *Warehouse) deleteWarehouse() bool {
	if w.Id == "" || len(w.Id) != 2 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", w.Id, w.EnterpriseId).Delete(&Warehouse{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findWarehouseByName(warehouseName string, enterpriseId int32) []NameString {
	var warehouses []NameString = make([]NameString, 0)
	result := dbOrm.Model(&Currency{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(warehouseName), enterpriseId).Limit(10).Find(&warehouses)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return warehouses
	}

	return warehouses
}

func getNameWarehouse(id string, enterpriseId int32) string {
	var warehouse Warehouse
	result := dbOrm.Where("id = ? AND enterprise = ?", id, enterpriseId).First(&warehouse)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	return warehouse.Name
}

// Regenerates the stock of the product for all the products in the database.
// This "stock" field is the sum of the stock in all the warehouses.
func regenerateProductStock(enterpriseId int32) bool {
	sqlStatement := `UPDATE product SET stock = CASE WHEN (SELECT SUM(quantity) FROM stock WHERE stock.product=product.id) IS NULL THEN 0 ELSE (SELECT SUM(quantity) FROM stock WHERE stock.product=product.id) END WHERE enterprise=$1`
	result := dbOrm.Exec(sqlStatement, enterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
