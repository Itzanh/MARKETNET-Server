package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesOrderDetailPackaged struct {
	OrderDetailId int64            `json:"orderDetailId" gorm:"primaryKey;column:order_detail;not null:true"`
	OrderDetail   SalesOrderDetail `json:"orderDetail" gorm:"foreignKey:OrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	PackagingId   int64            `json:"packagingId" gorm:"primaryKey;column:packaging;not null:true"`
	Packaging     Packaging        `json:"packaging" gorm:"foreignKey:PackagingId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity      int32            `json:"quantity" gorm:"column:quantity;not null:true"`
	EnterpriseId  int32            `json:"-" gorm:"column:enterprise;not null"`
	Enterprise    Settings         `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *SalesOrderDetailPackaged) TableName() string {
	return "sales_order_detail_packaged"
}

func getSalesOrderDetailPackaged(packagingId int64, enterpriseId int32) []SalesOrderDetailPackaged {
	var packaged []SalesOrderDetailPackaged = make([]SalesOrderDetailPackaged, 0)
	// get the packaged details for this packaging id and enterprise id using dbOrm
	result := dbOrm.Where("packaging = ? AND enterprise = ?", packagingId, enterpriseId).Preload(clause.Associations).Preload("OrderDetail.Product").Order("order_detail ASC").Find(&packaged)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return packaged
}

func getSalesOrderDetailPackagedRow(orderDetailId int64, packagingId int64) SalesOrderDetailPackaged {
	var detail SalesOrderDetailPackaged
	result := dbOrm.Where("order_detail = ? AND packaging = ?", orderDetailId, packagingId).Preload(clause.Associations).First(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return SalesOrderDetailPackaged{}
	}
	return detail
}

func (p *SalesOrderDetailPackaged) isValid() bool {
	return !(p.OrderDetailId <= 0 || p.PackagingId <= 0 || p.Quantity <= 0)
}

func (p *SalesOrderDetailPackaged) insertSalesOrderDetailPackaged(userId int32) bool {
	if !p.isValid() {
		return false
	}

	detail := getSalesOrderDetailRow(p.OrderDetailId)
	if detail.QuantityPendingPackaging <= 0 || p.Quantity > detail.QuantityPendingPackaging {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	var rowCount int64
	result := dbOrm.Model(&SalesOrderDetailPackaged{}).Where("order_detail = ? AND packaging = ?", p.OrderDetailId, p.PackagingId).Count(&rowCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	if rowCount == 0 {
		result := trans.Create(&p)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	} else {
		var detailPackaged SalesOrderDetailPackaged
		result := trans.Model(&SalesOrderDetailPackaged{}).Where("order_detail = ? AND packaging = ?", p.OrderDetailId, p.PackagingId).First(&detailPackaged)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		detailPackaged.Quantity += p.Quantity

		result = trans.Updates(&detailPackaged)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetailId, -p.Quantity, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	product := getProductRow(detail.ProductId)
	ok = addWeightPackaging(p.PackagingId, product.Weight*float64(p.Quantity), *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func (p *SalesOrderDetailPackaged) deleteSalesOrderDetailPackaged(userId int32, trans *gorm.DB) bool {
	if p.OrderDetailId <= 0 || p.PackagingId <= 0 {
		return false
	}

	inMemoryPackage := getSalesOrderDetailPackagedRow(p.OrderDetailId, p.PackagingId)
	if inMemoryPackage.OrderDetailId <= 0 || inMemoryPackage.EnterpriseId != p.EnterpriseId || inMemoryPackage.PackagingId <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	result := trans.Delete(&SalesOrderDetailPackaged{}, "order_detail = ? AND packaging = ? AND enterprise = ?", p.OrderDetailId, p.PackagingId, p.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetailId, inMemoryPackage.Quantity, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	detail := getSalesOrderDetailRow(p.OrderDetailId)
	product := getProductRow(detail.ProductId)
	ok = addWeightPackaging(p.PackagingId, -product.Weight*float64(inMemoryPackage.Quantity), *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		result := trans.Commit()
		if result.Error != nil {
			return false
		}
		///
	}
	return true
}

type SalesOrderDetailPackagedEAN13 struct {
	SalesOrder int64  `json:"salesOrder"`
	EAN13      string `json:"ean13"`
	Packaging  int64  `json:"packaging"`
	Quantity   int32  `json:"quantity"`
}

func (d *SalesOrderDetailPackagedEAN13) isValid() bool {
	return !(d.SalesOrder <= 0 || len(d.EAN13) != 13 || d.Packaging <= 0 || d.Quantity <= 0)
}

func (d *SalesOrderDetailPackagedEAN13) insertSalesOrderDetailPackagedEAN13(enterpriseId int32, userId int32) bool {
	if !d.isValid() {
		return false
	}

	product := getProductByBarcode(d.EAN13, enterpriseId)
	if product.Id <= 0 {
		return false
	}

	var detail SalesOrderDetail
	result := dbOrm.Where(`"order" = ? AND product = ?`, d.SalesOrder, product.Id).Preload(clause.Associations).First(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	p := SalesOrderDetailPackaged{}
	p.OrderDetailId = detail.Id
	p.PackagingId = d.Packaging
	p.Quantity = d.Quantity
	p.EnterpriseId = enterpriseId

	return p.insertSalesOrderDetailPackaged(userId)
}
