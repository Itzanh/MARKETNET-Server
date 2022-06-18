package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Packages struct {
	Id           int32    `json:"id" gorm:"index:_packages_id_enterprise,unique:true,priority:1"`
	Name         string   `json:"name" gorm:"type:character varying(50);not null:true"`
	Weight       float64  `json:"weight" gorm:"type:numeric(14,6);not null:true"`
	Width        float64  `json:"width" gorm:"type:numeric(14,6);not null:true"`
	Height       float64  `json:"height" gorm:"type:numeric(14,6);not null:true"`
	Depth        float64  `json:"depth" gorm:"type:numeric(14,6);not null:true"`
	ProductId    int32    `json:"productId" gorm:"column:product;not null:true"`
	Product      Product  `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:_packages_id_enterprise,unique:true,priority:2"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *Packages) TableName() string {
	return "packages"
}

func getPackages(enterpriseId int32) []Packages {
	var packages []Packages = make([]Packages, 0)
	dbOrm.Model(&Packages{}).Where("enterprise = ?", enterpriseId).Preload(clause.Associations).Order("id ASC").Find(&packages)
	return packages
}

func getPackagesRow(packageId int32) Packages {
	p := Packages{}
	dbOrm.Model(&Packages{}).Where("id = ?", packageId).Find(&p)
	return p
}

func (p *Packages) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 50 || p.Weight < 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 || p.ProductId <= 0)
}

func (p *Packages) BeforeCreate(tx *gorm.DB) (err error) {
	var packages Packages
	tx.Model(&Packages{}).Last(&packages)
	p.Id = packages.Id + 1
	return nil
}

func (p *Packages) insertPackage() bool {
	if !p.isValid() {
		return false
	}

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *Packages) updatePackage() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	var packages Packages
	result := dbOrm.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).First(&packages)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	packages.Name = p.Name
	packages.Weight = p.Weight
	packages.Width = p.Width
	packages.Height = p.Height
	packages.Depth = p.Depth
	packages.ProductId = p.ProductId

	result = dbOrm.Save(&packages)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *Packages) deletePackage() bool {
	if p.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).Delete(&Packages{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
