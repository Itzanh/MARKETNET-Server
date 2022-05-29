package main

import (
	"gorm.io/gorm/clause"
)

type TransferBetweenWarehousesMinimumStock struct {
	Id                           int64      `json:"id"`
	ProductId                    int32      `json:"productId" gorm:"column:product;not null:true"`
	Product                      Product    `json:"-" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseDestinationId       string     `json:"warehouseDestinationId" gorm:"column:warehouse_destination;type:character(2);not null"`
	WarehouseDestination         Warehouse  `json:"warehouseDestination" gorm:"foreignKey:WarehouseDestinationId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity                     int32      `json:"quantity" gorm:"column:quantity;type:integer;not null:true"`
	OriginWarehouseWithMoreStock bool       `json:"originWarehouseWithMoreStock" gorm:"column:origin_warehouse_with_more_stock;type:boolean;not null"`
	WarehouseOriginId            *string    `json:"warehouseOriginId" gorm:"column:warehouse_origin;type:character(2)"` // only different from null when the "originWarehouseWithMoreStock" field is false
	WarehouseOrigin              *Warehouse `json:"warehouseOrigin" gorm:"foreignKey:WarehouseOriginId,EnterpriseId;references:Id,EnterpriseId"`
	UseOtherWarehousesFallback   *bool      `json:"useOtherWarehousesFallback" gorm:"column:use_other_warehouses_fallback;type:boolean"` // only different from null when the "originWarehouseWithMoreStock" field is false
	EnterpriseId                 int32      `json:"enterprise" gorm:"column:enterprise;type:integer;not null"`
	Enterprise                   Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (TransferBetweenWarehousesMinimumStock) TableName() string {
	return "transfer_between_warehouses_minimum_stock"
}

func getTransferBetweenWarehousesMinimumStock(productId int32, enterpriseId int32) []TransferBetweenWarehousesMinimumStock {
	var transferBetweenWarehousesMinimumStock []TransferBetweenWarehousesMinimumStock
	result := dbOrm.Where("enterprise = ?", enterpriseId).Where("product = ?", productId).Preload(clause.Associations).Find(&transferBetweenWarehousesMinimumStock)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return transferBetweenWarehousesMinimumStock
}

func (m *TransferBetweenWarehousesMinimumStock) isValid() bool {
	if len(m.WarehouseDestinationId) != 2 || m.Quantity <= 0 || m.ProductId <= 0 {
		return false
	}

	if m.OriginWarehouseWithMoreStock {
		m.WarehouseOriginId = nil
		m.UseOtherWarehousesFallback = nil
		return true
	} else {
		return !(m.WarehouseOriginId == nil || len(*m.WarehouseOriginId) != 2 || m.UseOtherWarehousesFallback == nil || *m.WarehouseOriginId == m.WarehouseDestinationId)
	}
}

func (m *TransferBetweenWarehousesMinimumStock) insertTransferBetweenWarehousesMinimumStock() bool {
	if !m.isValid() {
		return false
	}

	result := dbOrm.Create(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (m *TransferBetweenWarehousesMinimumStock) updateTransferBetweenWarehousesMinimumStock() bool {
	if !m.isValid() {
		return false
	}

	var transferBetweenWarehousesMinimumStock TransferBetweenWarehousesMinimumStock
	result := dbOrm.Where("id = ? AND enterprise = ?", m.Id, m.EnterpriseId).Find(&transferBetweenWarehousesMinimumStock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	transferBetweenWarehousesMinimumStock.WarehouseDestinationId = m.WarehouseDestinationId
	transferBetweenWarehousesMinimumStock.Quantity = m.Quantity
	transferBetweenWarehousesMinimumStock.OriginWarehouseWithMoreStock = m.OriginWarehouseWithMoreStock
	transferBetweenWarehousesMinimumStock.WarehouseOriginId = m.WarehouseOriginId
	transferBetweenWarehousesMinimumStock.UseOtherWarehousesFallback = m.UseOtherWarehousesFallback

	result = dbOrm.Model(&TransferBetweenWarehousesMinimumStock{}).Where("id = ?", m.Id).Save(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (m *TransferBetweenWarehousesMinimumStock) deleteTransferBetweenWarehousesMinimumStock() bool {
	result := dbOrm.Model(&TransferBetweenWarehousesMinimumStock{}).Where("id = ? AND enterprise = ?", m.Id, m.EnterpriseId).Delete(&TransferBetweenWarehousesMinimumStock{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func getStockLowerThanTransferBetweenWarehousesForMinimumStock(enterpriseId int32) []TransferBetweenWarehousesMinimumStock {
	var transferBetweenWarehousesMinimumStock []TransferBetweenWarehousesMinimumStock
	result := dbOrm.Model(&TransferBetweenWarehousesMinimumStock{}).Joins("INNER JOIN stock ON stock.warehouse = transfer_between_warehouses_minimum_stock.warehouse_destination AND stock.product = transfer_between_warehouses_minimum_stock.product AND stock.enterprise = transfer_between_warehouses_minimum_stock.enterprise").Where("stock.quantity < transfer_between_warehouses_minimum_stock.quantity").Find(&transferBetweenWarehousesMinimumStock)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return transferBetweenWarehousesMinimumStock
}

func (m *TransferBetweenWarehousesMinimumStock) getTransferBetweenWarehousesMinimumStockWarehouse(enterpriseId int32) string {
	if m.OriginWarehouseWithMoreStock {
		stocks := getStock(m.ProductId, enterpriseId)
		if len(stocks) == 0 {
			return ""
		}

		var moreStock Stock = stocks[0]
		for _, stock := range stocks {
			if stock.QuantityAvaiable > moreStock.QuantityAvaiable {
				moreStock = stock
			}
		}
		return moreStock.WarehouseId
	} else { // if m.OriginWarehouseWithMoreStock {
		stock := getStockRow(m.ProductId, *m.WarehouseOriginId, enterpriseId)
		if stock.QuantityAvaiable > 0 {
			return stock.WarehouseId
		} else {
			if m.UseOtherWarehousesFallback != nil && *m.UseOtherWarehousesFallback {
				stocks := getStock(m.ProductId, enterpriseId)
				if len(stocks) == 0 {
					return ""
				}

				var moreStock Stock = stocks[0]
				for _, stock := range stocks {
					if stock.QuantityAvaiable > moreStock.QuantityAvaiable {
						moreStock = stock
					}
				}
				return moreStock.WarehouseId
			} else {
				return stock.WarehouseId
			}
		}
	}
}

func generateTransferBetweenWarehousesForMinimumStock(enterpriseId int32) bool {
	var transferBetweenWarehousesMinimumStock = getStockLowerThanTransferBetweenWarehousesForMinimumStock(enterpriseId)
	if len(transferBetweenWarehousesMinimumStock) == 0 {
		return false
	}

	var transfersBetweenWarehouses map[string]TransferBetweenWarehouses = make(map[string]TransferBetweenWarehouses)

	///
	trans := dbOrm.Begin()
	///

	for _, transferBetweenWarehousesMinimumStock := range transferBetweenWarehousesMinimumStock {
		warehouse := transferBetweenWarehousesMinimumStock.getTransferBetweenWarehousesMinimumStockWarehouse(enterpriseId)
		if warehouse == "" {
			trans.Rollback()
			return false
		}

		_, ok := transfersBetweenWarehouses[warehouse]
		if !ok {
			var transferBetweenWarehouses = TransferBetweenWarehouses{
				WarehouseOriginId:      warehouse,
				WarehouseDestinationId: transferBetweenWarehousesMinimumStock.WarehouseDestinationId,
				EnterpriseId:           enterpriseId,
				Name:                   warehouse + " -> " + transferBetweenWarehousesMinimumStock.WarehouseDestinationId,
			}
			ok = transferBetweenWarehouses.insertTransferBetweenWarehouses()
			if !ok {
				trans.Rollback()
				return false
			}
			transfersBetweenWarehouses[transferBetweenWarehousesMinimumStock.WarehouseDestinationId] = transferBetweenWarehouses
		}

		var transferBetweenWarehousesDetail = TransferBetweenWarehousesDetail{
			TransferBetweenWarehousesId: transfersBetweenWarehouses[transferBetweenWarehousesMinimumStock.WarehouseDestinationId].Id,
			EnterpriseId:                enterpriseId,
			ProductId:                   transferBetweenWarehousesMinimumStock.ProductId,
			Quantity:                    minInt32(transferBetweenWarehousesMinimumStock.Quantity, getStockRow(transferBetweenWarehousesMinimumStock.ProductId, warehouse, enterpriseId).QuantityAvaiable),
		}

		ok = transferBetweenWarehousesDetail.insertTransferBetweenWarehousesDetail()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	///

	return true
}
