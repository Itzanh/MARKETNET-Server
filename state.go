package main

import (
	"strings"
)

type State struct {
	Id          int32  `json:"id"`
	Country     int32  `json:"country"`
	Name        string `json:"name"`
	IsoCode     string `json:"isoCode"`
	CountryName string `json:"countryName"`
	enterprise  int32
}

func getStates(enterpriseId int32) []State {
	var cities []State = make([]State, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=state.country) FROM public.state WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return cities
	}
	for rows.Next() {
		s := State{}
		rows.Scan(&s.Id, &s.Country, &s.Name, &s.IsoCode, &s.enterprise, &s.CountryName)
		cities = append(cities, s)
	}

	return cities
}

func getStatesByCountry(countryId int32, enterpriseId int32) []State {
	var cities []State = make([]State, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=state.country) FROM public.state WHERE country=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, countryId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return cities
	}
	for rows.Next() {
		s := State{}
		rows.Scan(&s.Id, &s.Country, &s.Name, &s.IsoCode, &s.enterprise, &s.CountryName)
		cities = append(cities, s)
	}

	return cities
}

func getStateRow(id int32) State {
	sqlStatement := `SELECT * FROM public.state WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return State{}
	}

	s := State{}
	row.Scan(&s.Id, &s.Country, &s.Name, &s.IsoCode, &s.enterprise)

	return s
}

func (c *State) isValid() bool {
	return !(c.Country <= 0 || len(c.Name) == 0 || len(c.Name) > 100 || len(c.IsoCode) > 7)
}

func searchStates(search string, enterpriseId int32) []State {
	var states []State = make([]State, 0)
	sqlStatement := `SELECT state.*,(SELECT name FROM country WHERE country.id=state.country) FROM state INNER JOIN country ON country.id=state.country WHERE (state.name ILIKE $1 OR country.name ILIKE $1) AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%", enterpriseId)
	if err != nil {
		log("DB", err.Error())
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

	sqlStatement := `INSERT INTO public.state(country, name, iso_code, enterprise) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, c.Country, c.Name, c.IsoCode, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *State) updateState() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.state SET country=$2, name=$3, iso_code=$4 WHERE id=$1 AND enterprise=$5`
	res, err := db.Exec(sqlStatement, c.Id, c.Country, c.Name, c.IsoCode, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *State) deleteState() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.state WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type StateNameQuery struct {
	CountryId int16  `json:"countryId"`
	Name      string `json:"cityName"`
}

func findStateByName(cityName StateNameQuery, enterpriseId int32) []NameInt32 {
	var cities []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.state WHERE (country=$1 AND UPPER(name) LIKE $2 || '%') AND enterprise=$3 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, cityName.CountryId, strings.ToUpper(cityName.Name), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return cities
	}
	for rows.Next() {
		c := NameInt32{}
		rows.Scan(&c.Id, &c.Name)
		cities = append(cities, c)
	}

	return cities
}

func getNameState(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.state WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
