package main

import (
	"strings"
)

type State struct {
	Id          int32  `json:"id"`
	Country     int16  `json:"country"`
	Name        string `json:"name"`
	IsoCode     string `json:"isoCode"`
	CountryName string `json:"countryName"`
}

func getStates() []State {
	var cities []State = make([]State, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=state.country) FROM public.state ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return cities
	}
	for rows.Next() {
		c := State{}
		rows.Scan(&c.Id, &c.Country, &c.Name, &c.IsoCode, &c.CountryName)
		cities = append(cities, c)
	}

	return cities
}

func getStatesByCountry(countryId int16) []State {
	var cities []State = make([]State, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=state.country) FROM public.state WHERE country=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, countryId)
	if err != nil {
		return cities
	}
	for rows.Next() {
		c := State{}
		rows.Scan(&c.Id, &c.Country, &c.Name, &c.IsoCode, &c.CountryName)
		cities = append(cities, c)
	}

	return cities
}

func (c *State) isValid() bool {
	return !(c.Country <= 0 || len(c.Name) == 0 || len(c.Name) > 100 || len(c.IsoCode) > 7)
}

func searchStates(search string) []State {
	var states []State = make([]State, 0)
	sqlStatement := `SELECT state.*,(SELECT name FROM country WHERE country.id=state.country) FROM state INNER JOIN country ON country.id=state.country WHERE state.name ILIKE $1 OR country.name ILIKE $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%")
	if err != nil {
		return states
	}
	for rows.Next() {
		c := State{}
		rows.Scan(&c.Id, &c.Country, &c.Name, &c.IsoCode, &c.CountryName)
		states = append(states, c)
	}

	return states
}

func (c *State) insertState() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.state(country, name, iso_code) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, c.Country, c.Name, c.IsoCode)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *State) updateState() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.state SET country=$2, name=$3, iso_code=$4 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Country, c.Name, c.IsoCode)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *State) deleteState() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.state WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type StateNameQuery struct {
	CountryId int16  `json:"countryId"`
	Name      string `json:"cityName"`
}

func findStateByName(cityName StateNameQuery) []NameInt32 {
	var cities []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.state WHERE country=$1 AND UPPER(name) LIKE $2 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, cityName.CountryId, strings.ToUpper(cityName.Name))
	if err != nil {
		return cities
	}
	for rows.Next() {
		c := NameInt32{}
		rows.Scan(&c.Id, &c.Name)
		cities = append(cities, c)
	}

	return cities
}

func getNameState(id int32) string {
	sqlStatement := `SELECT name FROM public.state WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
