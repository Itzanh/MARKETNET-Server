package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Country struct {
	Id          int16  `json:"id"`
	Name        string `json:"name"`
	Iso2        string `json:"iso2"`
	Iso3        string `json:"iso3"`
	UNCode      int16  `json:"unCode"`
	Zone        string `json:"zone"`
	PhonePrefix int16  `json:"phonePrefix"`
	Language    *int16 `json:"language"`
	Currency    *int16 `json:"currency"`
}

func getCountries() []Country {
	var countries []Country = make([]Country, 0)
	sqlStatement := `SELECT * FROM public.country ORDER BY id ASC `
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return countries
	}
	for rows.Next() {
		c := Country{}
		rows.Scan(&c.Id, &c.Name, &c.Iso2, &c.Iso3, &c.UNCode, &c.Zone, &c.PhonePrefix, &c.Language, &c.Currency)
		countries = append(countries, c)
	}

	return countries
}

func (c *Country) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 75 || len(c.Iso2) != 2 || (len(c.Iso3) != 0 && len(c.Iso3) != 3) || c.UNCode < 0 || (c.Zone != "N" && c.Zone != "U" && c.Zone != "E") || c.PhonePrefix < 0)
}

func (c *Country) insertCountry() bool {
	if !c.isValid() {
		fmt.Println("INVALID")
		data, _ := json.Marshal(c)
		fmt.Println(string(data))
		return false
	}

	sqlStatement := `INSERT INTO public.country(name, iso_2, iso_3, un_code, zone, phone_prefix, language, currency) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	res, err := db.Exec(sqlStatement, c.Name, c.Iso2, c.Iso3, c.UNCode, c.Zone, c.PhonePrefix, c.Language, c.Currency)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Country) updateCountry() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.country SET name=$2, iso_2=$3, iso_3=$4, un_code=$5, zone=$6, phone_prefix=$7, language=$8, currency=$9 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Iso2, c.Iso3, c.UNCode, c.Zone, c.PhonePrefix, c.Language, c.Currency)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Country) deleteCountry() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.country WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCountryByName(languageName string) []NameInt16 {
	var countries []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.country WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return countries
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		countries = append(countries, c)
	}

	return countries
}

func getNameCountry(id int16) string {
	sqlStatement := `SELECT name FROM public.country WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
