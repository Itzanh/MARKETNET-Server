package main

import (
	"strings"
)

type City struct {
	Id      int32  `json:"id"`
	Country int16  `json:"country"`
	Name    string `json:"name"`
	ZipCode string `json:"zipCode"`
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
		rows.Scan(&c.Id, &c.Country, &c.Name, &c.ZipCode)
		cities = append(cities, c)
	}

	return cities
}

func (c *City) isValid() bool {
	return !(c.Country <= 0 || len(c.Name) == 0 || len(c.Name) > 100 || len(c.ZipCode) == 0 || len(c.ZipCode) > 15)
}

func (c *City) insertCity() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.city(country, name, zip_code) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, c.Country, c.Name, c.ZipCode)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *City) updateCity() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.city SET country=$2, name=$3, zip_code=$4 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Country, c.Name, c.ZipCode)
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

type CityName struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func findCityByName(cityName CityNameQuery) []CityName {
	var cities []CityName = make([]CityName, 0)
	sqlStatement := `SELECT id,name FROM public.city WHERE country = $1 AND UPPER(name) LIKE $2 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, cityName.CountryId, strings.ToUpper(cityName.Name))
	if err != nil {
		return cities
	}
	for rows.Next() {
		c := CityName{}
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
