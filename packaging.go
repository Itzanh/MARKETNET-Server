package main

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Packaging struct {
	Id              int64                      `json:"id" gorm:"index:_packaging_id_enterprise,unique:true,priority:1"`
	PackageId       int32                      `json:"packageId" gorm:"column:package;not null:true"`
	Package         Packages                   `json:"package" gorm:"foreignKey:PackageId,EnterpriseId;references:Id,EnterpriseId"`
	PackageName     string                     `json:"packageName" gorm:"-"` // Computed server-side
	SalesOrderId    int64                      `json:"salesOrderId" gorm:"column:sales_order;not null:true"`
	SalesOrder      SaleOrder                  `json:"salesOrder" gorm:"foreignKey:SalesOrderId,EnterpriseId;references:Id,EnterpriseId"`
	Weight          float64                    `json:"weight" gorm:"column:weight;not null:true;type:numeric(14,6)"`
	ShippingId      *int64                     `json:"shippingId" gorm:"column:shipping"`
	Shipping        *Shipping                  `json:"shipping" gorm:"foreignKey:ShippingId,EnterpriseId;references:Id,EnterpriseId"`
	DetailsPackaged []SalesOrderDetailPackaged `json:"detailsPackaged" gorm:"-"` // Computed server-side
	PalletId        *int32                     `json:"palletId" gorm:"column:pallet"`
	Pallet          *Pallet                    `json:"pallet" gorm:"foreignKey:PalletId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId    int32                      `json:"-" gorm:"column:enterprise;not null:true;index:_packaging_id_enterprise,unique:true,priority:2"`
	Enterprise      Settings                   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *Packaging) TableName() string {
	return "packaging"
}

func getPackaging(salesOrderId int64, enterpriseId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	result := dbOrm.Model(&Packaging{}).Where("sales_order = ? AND enterprise = ?", salesOrderId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&packaging)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return packaging
	}

	for i := 0; i < len(packaging); i++ {
		p := packaging[i]
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id, enterpriseId)

		_package := getPackagesRow(p.PackageId)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging[i] = p
	}

	return packaging
}

func getPackagingByShipping(shippingId int64, enterpriseId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	result := dbOrm.Model(&Packaging{}).Where("shipping = ? AND enterprise = ?", shippingId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&packaging)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return packaging
	}

	for i := 0; i < len(packaging); i++ {
		p := packaging[i]
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id, enterpriseId)

		_package := getPackagesRow(p.PackageId)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging[i] = p
	}

	return packaging
}

func getPackagingRow(packagingId int64) Packaging {
	p := Packaging{}
	result := dbOrm.Model(&Packaging{}).Where("id = ?", packagingId).Preload(clause.Associations).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return p
}

func (p *Packaging) isValid() bool {
	return !(p.PackageId <= 0 || p.SalesOrderId <= 0)
}

func (p *Packaging) BeforeCreate(tx *gorm.DB) (err error) {
	var packaging Packaging
	tx.Model(&Packaging{}).Last(&packaging)
	p.Id = packaging.Id + 1
	return nil
}

func (p *Packaging) insertPackaging() bool {
	if !p.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	_package := getPackagesRow(p.PackageId)
	if _package.Id <= 0 {
		trans.Rollback()
		return false
	}
	p.Weight = _package.Weight
	p.ShippingId = nil

	result := trans.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	s := getSalesOrderRow(p.SalesOrderId)
	addQuantityStock(_package.ProductId, s.WarehouseId, -1, p.EnterpriseId, *trans)

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func (p *Packaging) deletePackaging(enterpriseId int32, userId int32) bool {
	if p.Id <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	inMemoryPackaging := getPackagingRow(p.Id)
	if inMemoryPackaging.Id <= 0 || inMemoryPackaging.EnterpriseId != enterpriseId {
		trans.Rollback()
		return false
	}

	detailsPackaged := getSalesOrderDetailPackaged(p.Id, enterpriseId)
	for i := 0; i < len(detailsPackaged); i++ {
		detailsPackaged[i].EnterpriseId = enterpriseId
		ok := detailsPackaged[i].deleteSalesOrderDetailPackaged(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	result := trans.Delete(&Packaging{}, "id = ? AND enterprise = ?", p.Id, p.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	_package := getPackagesRow(inMemoryPackaging.PackageId)
	s := getSalesOrderRow(inMemoryPackaging.SalesOrderId)
	addQuantityStock(_package.ProductId, s.WarehouseId, 1, p.EnterpriseId, *trans)

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addWeightPackaging(packagingId int64, weight float64, trans gorm.DB) bool {
	var packagingWeight float64
	result := trans.Model(&Packaging{}).Where("id = ?", packagingId).Pluck("weight", &packagingWeight)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&Packaging{}).Where("id = ?", packagingId).Update("weight", packagingWeight+weight)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return true
}
