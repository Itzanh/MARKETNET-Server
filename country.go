package main

import (
	"database/sql"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Country struct {
	Id           int32     `json:"id" gorm:"index:country_id_enterprise,unique:true,priority:1"`
	Name         string    `json:"name" gorm:"type:character varying(50);not null:true;index:country_name,type:gin"`
	Iso2         string    `json:"iso2" gorm:"column:iso_2;type:character(2);not null:true;index:country_iso_2,unique:true,priority:2"`
	Iso3         string    `json:"iso3" gorm:"column:iso_3;type:character(3);not null:true;index:country_iso_3,unique:true,priority:2,where:iso_3 <> ''::bpchar"`
	UNCode       int16     `json:"unCode" gorm:"not null:true"`
	Zone         string    `json:"zone" gorm:"type:character(1);not null:true"` // N = National, U = European Union, E = Export
	PhonePrefix  int16     `json:"phonePrefix" gorm:"not null:true"`
	LanguageId   *int32    `json:"languageId" gorm:"column:language"`
	Language     *Language `json:"language" gorm:"foreignKey:LanguageId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyId   *int32    `json:"currencyId" gorm:"column:currency"`
	Currency     *Currency `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise;not null:true;index:country_id_enterprise,unique:true,priority:2;index:country_iso_2,unique:true,priority:1;index:country_iso_3,unique:true,priority:1,where:iso_3 <> ''::bpchar"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Country) TableName() string {
	return "country"
}

func getCountries(enterpriseId int32) []Country {
	var countries []Country = make([]Country, 0)
	result := dbOrm.Model(&Country{}).Where("country.enterprise = ?", enterpriseId).Joins("Language").Joins("Currency").Order("id ASC").Find(&countries)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return countries
	}

	return countries
}

func getCountryRow(id int32, enterpriseId int32) Country {
	var currency Country
	result := dbOrm.Where("id = ? AND enterprise = ?", id, enterpriseId).Preload(clause.Associations).First(&currency)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return Country{}
	}

	return currency
}

func (c *Country) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 75 || len(c.Iso2) != 2 || (len(c.Iso3) != 0 && len(c.Iso3) != 3) || c.UNCode < 0 || (c.Zone != "N" && c.Zone != "U" && c.Zone != "E") || c.PhonePrefix < 0)
}

func searchCountries(search string, enterpriseId int32) []Country {
	var countries []Country = make([]Country, 0)
	result := dbOrm.Model(&Country{}).Where("(country.name ILIKE @search_contains OR country.iso_2 = UPPER(@search) OR country.iso_3 = UPPER(@search)) AND country.enterprise=@enterprise_id",
		sql.Named("search_contains", "%"+search+"%"), sql.Named("search", search), sql.Named("enterprise_id", enterpriseId)).Joins("Language").Joins("Currency").Order("id ASC").Find(&countries)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return countries
	}

	return countries
}

func (c *Country) BeforeCreate(tx *gorm.DB) (err error) {
	var country Country
	tx.Model(&Country{}).Last(&country)
	c.Id = country.Id + 1
	return nil
}

func (c *Country) insertCountry() bool {
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

func (c *Country) updateCountry() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	var country Country
	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&country)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	country.Name = c.Name
	country.Iso2 = c.Iso2
	country.Iso3 = c.Iso3
	country.UNCode = c.UNCode
	country.Zone = c.Zone
	country.PhonePrefix = c.PhonePrefix
	country.LanguageId = c.LanguageId
	country.CurrencyId = c.CurrencyId

	result = dbOrm.Save(&country)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Country) deleteCountry() bool {
	if c.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Country{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findCountryByName(countryName string, enterpriseId int32) []NameInt32 {
	var countries []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Country{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(countryName), enterpriseId).Limit(10).Find(&countries)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return countries
	}

	return countries
}
