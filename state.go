package main

import (
	"database/sql"
	"strings"

	"gorm.io/gorm"
)

type State struct {
	Id           int32    `json:"id" gorm:"index:state_id_enterprise,unique:true,priority:1"`
	CountryId    int32    `json:"countryId" gorm:"column:country;not null:true"`
	Country      Country  `json:"country" gorm:"foreignKey:CountryId,EnterpriseId;references:Id,EnterpriseId"`
	Name         string   `json:"name" gorm:"type:character varying(100);not null:true;index:state_name,type:gin"`
	IsoCode      string   `json:"isoCode" gorm:"type:character varying(7);not null:true;index:state_iso_code,where:iso_code <> ''"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true;index:state_id_enterprise,unique:true,priority:2"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *State) TableName() string {
	return "state"
}

func getStates(enterpriseId int32) []State {
	var states []State = make([]State, 0)
	result := dbOrm.Model(&State{}).Where("enterprise = ?", enterpriseId).Preload("Country").Order("id ASC").Find(&states)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return states
}

func getStatesByCountry(countryId int32, enterpriseId int32) []State {
	var states []State = make([]State, 0)
	dbOrm.Model(&State{}).Where("country = ? AND enterprise = ?", countryId, enterpriseId).Preload("Country").Order("id ASC").Find(&states)
	result := dbOrm.Model(&State{}).Where("enterprise = ?", enterpriseId).Preload("Country").Order("id ASC").Find(&states)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return states
}

func getStateRow(id int32) State {
	s := State{}
	dbOrm.Model(&State{}).Where("id = ?", id).First(&s)
	return s
}

func searchStates(search string, enterpriseId int32) []State {
	var states []State = make([]State, 0)
	result := dbOrm.Model(&State{}).Where("state.name ILIKE @search AND state.enterprise = @enterpriseId", sql.Named("search", "%"+search+"%"), sql.Named("enterpriseId", enterpriseId)).Preload("Country").Order("state.id ASC").Find(&states)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return states
}

func (s *State) BeforeCreate(tx *gorm.DB) (err error) {
	var state State
	tx.Model(&State{}).Last(&state)
	s.Id = state.Id + 1
	return nil
}

func (c *State) isValid() bool {
	return !(c.CountryId <= 0 || len(c.Name) == 0 || len(c.Name) > 100 || len(c.IsoCode) > 7)
}

func (s *State) insertState() bool {
	if !s.isValid() {
		return false
	}

	result := dbOrm.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *State) updateState() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	var state State
	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&state)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	state.CountryId = c.CountryId
	state.Name = c.Name
	state.IsoCode = c.IsoCode

	result = dbOrm.Save(&state)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *State) deleteState() bool {
	if s.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).Delete(&State{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type StateNameQuery struct {
	CountryId int16  `json:"countryId"`
	Name      string `json:"cityName"`
}

func findStateByName(cityName StateNameQuery, enterpriseId int32) []NameInt32 {
	var states []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&State{}).Where("country=$1 AND UPPER(name) LIKE ($2 || '%') AND enterprise=$3", cityName.CountryId, strings.ToUpper(cityName.Name), enterpriseId).Limit(10).Find(&states)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return states
	}

	return states
}
