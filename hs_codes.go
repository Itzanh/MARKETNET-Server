package main

import "database/sql"

type HSCode struct {
	Id   string `json:"id" gorm:"primaryKey;type:character varying(8);not null:true"`
	Name string `json:"name" gorm:"type:character varying(255);not null:true"`
}

func (h *HSCode) TableName() string {
	return "hs_codes"
}

type HSCodeQuery struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (q *HSCodeQuery) getHSCodes() []HSCode {
	var codes []HSCode = make([]HSCode, 0)
	// get the HS Codes from the database where id or name matches the query using named arguments order by id ascending using dbOrm
	dbOrm.Model(HSCode{}).Where("id LIKE @search OR name LIKE @search", sql.Named("search", "%"+q.Id+"%")).Order("id ASC").Find(&codes)
	return codes
}
