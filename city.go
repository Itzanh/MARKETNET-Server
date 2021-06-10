package main

import (
	"fmt"
	"strings"
)

type City struct {
	Id        int32  `json:"id"`
	Country   int16  `json:"country"`
	Name      string `json:"name"`
	ZipCode   string `json:"zipCode"`
	NameAscii string `json:"nameAscii"`
}

func getCities() []City {
	var cities []City = make([]City, 0)
	sqlStatement := `SELECT * FROM public.city ORDER BY id ASC `
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return cities
	}
	for rows.Next() {
		c := City{}
		rows.Scan(&c.Id, &c.Country, &c.Name, &c.ZipCode, &c.NameAscii)
		cities = append(cities, c)
	}

	return cities
}

func (c *City) isValid() bool {
	return !(c.Country <= 0 || len(c.Name) == 0 || len(c.Name) > 100 || len(c.ZipCode) > 15 || len(c.NameAscii) == 0 || len(c.NameAscii) > 100)
}

func (c *City) insertCity() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.city(country, name, zip_code, name_ascii) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, c.Country, c.Name, c.ZipCode, c.NameAscii)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *City) updateCity() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.city SET country=$2, name=$3, zip_code=$4, name_ascii=$5 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Country, c.Name, c.ZipCode, c.NameAscii)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *City) deleteCity() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.city WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type CityNameQuery struct {
	CountryId int16  `json:"countryId"`
	Name      string `json:"cityName"`
}

func findCityByName(cityName CityNameQuery) []NameInt32 {
	var cities []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.city WHERE country = $1 AND UPPER(name) LIKE $2 || '%' ORDER BY id ASC LIMIT 10`
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

func getNameCity(id int32) string {
	sqlStatement := `SELECT name FROM public.city WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
