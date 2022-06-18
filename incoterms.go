package main

import (
	"strings"

	"gorm.io/gorm"
)

type Incoterm struct {
	Id           int32    `json:"id" gorm:"index:incoterm_id_enterprise,unique:true,priority:1"`
	Key          string   `json:"key" gorm:"type:character(3);not null:true;index:incoterm_key,unique:true,priority:2"`
	Name         string   `json:"name" gorm:"type:character varying(50);not null:true"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:incoterm_id_enterprise,unique:true,priority:2;index:incoterm_key,unique:true,priority:1"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (i *Incoterm) TableName() string {
	return "incoterm"
}

func getIncoterm(enterpriseId int32) []Incoterm {
	var incoterms []Incoterm = make([]Incoterm, 0)
	dbOrm.Model(&Incoterm{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&incoterms)
	return incoterms
}

func (i *Incoterm) isValid() bool {
	i.Key = strings.ToUpper(i.Key)
	return !(len(i.Key) == 0 || len(i.Key) > 3 || len(i.Name) == 0 || len(i.Name) > 50)
}

func (i *Incoterm) BeforeCreate(tx *gorm.DB) (err error) {
	var incoterm Incoterm
	tx.Model(&Incoterm{}).Last(&incoterm)
	i.Id = incoterm.Id + 1
	return nil
}

func (i *Incoterm) insertIncoterm() bool {
	if !i.isValid() {
		return false
	}

	result := dbOrm.Create(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (i *Incoterm) updateIncoterm() bool {
	if i.Id <= 0 || !i.isValid() {
		return false
	}

	var incoterm Incoterm
	result := dbOrm.Where("id = ? AND enterprise = ?", i.Id, i.EnterpriseId).First(&incoterm)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	incoterm.Name = i.Name
	incoterm.Key = i.Key

	result = dbOrm.Save(&incoterm)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (i *Incoterm) deleteIncoterm() bool {
	if i.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", i.Id, i.EnterpriseId).Delete(&Incoterm{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
