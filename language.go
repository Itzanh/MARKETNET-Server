package main

import (
	"database/sql"
	"strings"

	"gorm.io/gorm"
)

type Language struct {
	Id           int32    `json:"id" gorm:"index:language_id_enterprise,unique:true,priority:1"`
	Name         string   `json:"name" gorm:"type:character varying(50);not null:true;index:language_name,type:gin"`
	Iso2         string   `json:"iso2" gorm:"column:iso_2;type:character(2);not null:true;index:language_iso_2,unique:true,priority:2"`
	Iso3         string   `json:"iso3" gorm:"column:iso_3;type:character(3);not null:true;index:language_iso_3,unique:true,priority:2,where:iso_3 <> ''::bpchar"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:language_id_enterprise,unique:true,priority:2;index:language_iso_2,unique:true,priority:1;index:language_iso_3,unique:true,priority:1,where:iso_3 <> ''::bpchar"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Language) TableName() string {
	return "language"
}

func getLanguages(enterpriseId int32) []Language {
	var languages []Language = make([]Language, 0)
	result := dbOrm.Model(&Language{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&languages)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return languages
	}

	return languages
}

func (l *Language) isValid() bool {
	l.Iso2 = strings.ToUpper(l.Iso2)
	l.Iso3 = strings.ToUpper(l.Iso3)
	return !(len(l.Name) == 0 || len(l.Name) > 50 || len(l.Iso2) != 2 || len(l.Iso3) != 3)
}

func searchLanguages(search string, enterpriseId int32) []Language {
	var languages []Language = make([]Language, 0)
	result := dbOrm.Model(&Language{}).Where("(name ILIKE @search_contains OR iso_2 = UPPER(@search) OR iso_3 = UPPER(@search)) AND enterprise=@enterprise_id",
		sql.Named("search_contains", "%"+search+"%"), sql.Named("search", search), sql.Named("enterprise_id", enterpriseId)).Order("id ASC").Find(&languages)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return languages
	}

	return languages
}

func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {
	var language Language
	tx.Model(&Language{}).Last(&language)
	l.Id = language.Id + 1
	return nil
}

func (l *Language) insertLanguage() bool {
	if !l.isValid() {
		return false
	}

	result := dbOrm.Create(&l)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (l *Language) updateLanguage() bool {
	if l.Id <= 0 || !l.isValid() {
		return false
	}

	var language Language
	result := dbOrm.Where("id = ? AND enterprise = ?", l.Id, l.EnterpriseId).First(&language)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	language.Name = l.Name
	language.Iso2 = l.Iso2
	language.Iso3 = l.Iso3

	result = dbOrm.Save(&language)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (l *Language) deleteLanguage() bool {
	if l.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", l.Id, l.EnterpriseId).Delete(&Language{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
