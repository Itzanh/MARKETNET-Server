package main

import (
	"gorm.io/gorm"
)

type Stock struct {
	ProductId                  int32     `json:"productId" gorm:"primaryKey;column:product;not null:true"`
	Product                    Product   `json:"-" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseId                string    `json:"warehouseId" gorm:"primaryKey;column:warehouse;not null:true;type:character(2)"`
	Warehouse                  Warehouse `json:"-" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity                   int32     `json:"quantity" gorm:"column:quantity;not null:true"`
	QuantityPendingReceived    int32     `json:"quantityPendingReceived" gorm:"column:quantity_pending_received;not null:true"`
	QuantityPendingServed      int32     `json:"quantityPendingServed" gorm:"column:quantity_pending_served;not null:true"`
	QuantityAvaiable           int32     `json:"quantityAvaiable" gorm:"column:quantity_available;not null:true"`
	QuantityPendingManufacture int32     `json:"quantityPendingManufacture" gorm:"column:quantity_pending_manufacture;not null:true"`
	EnterpriseId               int32     `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise                 Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *Stock) TableName() string {
	return "stock"
}

func getStock(productId int32, enterpriseId int32) []Stock {
	var stock []Stock = make([]Stock, 0)
	result := dbOrm.Model(&Stock{}).Where("stock.product = ? AND stock.enterprise = ?", productId, enterpriseId).Order("warehouse ASC").Joins("Warehouse").Find(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return stock
}

func getStockRow(productId int32, warehouseId string, enterpriseId int32) Stock {
	s := Stock{}
	result := dbOrm.Model(&Stock{}).Where("stock.product = ? AND stock.warehouse = ? AND stock.enterprise = ?", productId, warehouseId, enterpriseId).Joins("Warehouse").First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return Stock{}
	}
	return s
}

// Inserts a row with 0 stock in all columns
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func createStockRow(productId int32, warehouseId string, enterpriseId int32, trans gorm.DB) bool {
	var stock Stock = Stock{
		ProductId:    productId,
		WarehouseId:  warehouseId,
		EnterpriseId: enterpriseId,
	}
	result := trans.Create(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	return true
}

// Adds an amount to the quantity pending of serving, and substract the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingServing(productId int32, warehouseId string, quantity int32, enterpriseId int32, trans gorm.DB) bool {
	var stockRowCount int64
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).Count(&stockRowCount).First(&stock)
	if stockRowCount == 0 { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId, trans) { // we create the row, and retry the operation
			return addQuantityPendingServing(productId, warehouseId, quantity, enterpriseId, trans)
		} else {
			return false // the row could neither not be created or updated
		}
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.QuantityPendingServed += quantity
	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	ok := setQuantityAvailable(productId, warehouseId, enterpriseId, trans)
	return ok
}

// Adds an amount to the quantity pending of receiving, and add to the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingReveiving(productId int32, warehouseId string, quantity int32, enterpriseId int32, trans gorm.DB) bool {
	var stockRowCount int64
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).Count(&stockRowCount).First(&stock)
	if stockRowCount == 0 { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId, trans) { // we create the row, and retry the operation
			return addQuantityPendingReveiving(productId, warehouseId, quantity, enterpriseId, trans)
		} else {
			return false // the row could neither not be created or updated
		}
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.QuantityPendingReceived += quantity
	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	ok := setQuantityAvailable(productId, warehouseId, enterpriseId, trans)
	return ok
}

// Adds an amount to the quantity pending of manufacturing, and add to the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingManufacture(productId int32, warehouseId string, quantity int32, enterpriseId int32, trans gorm.DB) bool {
	var stockRowCount int64
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).Count(&stockRowCount).First(&stock)
	if stockRowCount == 0 { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId, trans) { // we create the row, and retry the operation
			return addQuantityPendingReveiving(productId, warehouseId, quantity, enterpriseId, trans)
		} else {
			return false // the row could neither not be created or updated
		}
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.QuantityPendingManufacture += quantity
	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return setQuantityAvailable(productId, warehouseId, enterpriseId, trans)
}

// Add an amount to the stock column on the stock row for this product.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityStock(productId int32, warehouseId string, quantity int32, enterpriseId int32, trans gorm.DB) bool {
	productRow := getProductRow(productId)
	if productRow.EnterpriseId != enterpriseId {
		return false
	}
	if !productRow.ControlStock {
		return true
	}

	var stockRowCount int64
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).Count(&stockRowCount).First(&stock)
	if stockRowCount == 0 { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId, trans) { // we create the row, and retry the operation
			return addQuantityStock(productId, warehouseId, quantity, enterpriseId, trans)
		} else {
			return false // the row could neither not be created or updated
		}
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.Quantity += quantity
	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return setQuantityAvailable(productId, warehouseId, enterpriseId, trans) && setProductStockAllWarehouses(productId, trans)
}

// Sets an amount to the stock column on the stock row for this product.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setQuantityStock(productId int32, warehouseId string, quantity int32, enterpriseId int32, trans gorm.DB) bool {
	productRow := getProductRow(productId)
	if productRow.EnterpriseId != enterpriseId {
		return false
	}
	if !productRow.ControlStock {
		return true
	}

	var stockRowCount int64
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).Count(&stockRowCount).First(&stock)
	if stockRowCount == 0 { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId, trans) { // we create the row, and retry the operation
			return setQuantityStock(productId, warehouseId, quantity, enterpriseId, trans)
		} else {
			return false // the row could neither not be created or updated
		}
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.Quantity = quantity
	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return setQuantityAvailable(productId, warehouseId, enterpriseId, trans) && setProductStockAllWarehouses(productId, trans)
}

// Sets the "Quantity available" field on the stock in the products.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setQuantityAvailable(productId int32, warehouseId string, enterpriseId int32, trans gorm.DB) bool {
	var stock Stock
	result := trans.Model(&Stock{}).Where("product = ? AND warehouse = ? AND enterprise = ?", productId, warehouseId, enterpriseId).First(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	stock.QuantityAvaiable = stock.Quantity + stock.QuantityPendingReceived - stock.QuantityPendingServed + stock.QuantityPendingManufacture

	result = trans.Save(&stock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return true
}

// Sets the "stock" field on the product row, sum of the stocks in all the warehouses.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setProductStockAllWarehouses(productId int32, trans gorm.DB) bool {
	var quantity int32
	result := trans.Model(&Stock{}).Where("product = ?", productId).Select("SUM(quantity) as quantity").Scan(&quantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&Product{}).Where("id = ?", productId).Update("stock", quantity)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return true
}
