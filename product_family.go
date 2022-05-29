package main

import (
	"gorm.io/gorm"
)

type ProductFamily struct {
	Id           int32    `json:"id" gorm:"index:product_family_id_enterprise,unique:true,priority:1"`
	Name         string   `json:"name" gorm:"type:character varying(100);not null:true"`
	Reference    string   `json:"reference" gorm:"type:character varying(40);not null:true;index:product_family_reference,unique:true,priority:2"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:product_family_id_enterprise,unique:true,priority:2;index:product_family_reference,unique:true,priority:1"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (pf *ProductFamily) TableName() string {
	return "product_family"
}

func getProductFamilies(enterpriseId int32) []ProductFamily {
	var families []ProductFamily = make([]ProductFamily, 0)
	dbOrm.Model(&ProductFamily{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&families)
	return families
}

func (f *ProductFamily) isValid() bool {
	return !(len(f.Name) == 0 || len(f.Name) > 100 || len(f.Reference) == 0 || len(f.Reference) > 40)
}

func (f *ProductFamily) BeforeCreate(tx *gorm.DB) (err error) {
	var productFamily ProductFamily
	tx.Model(&ProductFamily{}).Last(&productFamily)
	f.Id = productFamily.Id + 1
	return nil
}

func (f *ProductFamily) insertProductFamily() bool {
	if !f.isValid() {
		return false
	}

	result := dbOrm.Create(&f)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (f *ProductFamily) updateProductFamily() bool {
	if f.Id <= 0 || !f.isValid() {
		return false
	}

	var productFamily ProductFamily
	result := dbOrm.Where("id = ? AND enterprise = ?", f.Id, f.EnterpriseId).First(&productFamily)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	productFamily.Name = f.Name
	productFamily.Reference = f.Reference

	result = dbOrm.Save(&productFamily)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (f *ProductFamily) deleteProductFamily() bool {
	if f.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", f.Id, f.EnterpriseId).Delete(&ProductFamily{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func locateProductFamilies(enterpriseId int32) []NameInt32 {
	var families []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&ProductFamily{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&families)
	return families
}
