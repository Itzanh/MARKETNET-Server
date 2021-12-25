package main

import (
	"strings"
)

type Country struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Iso2        string `json:"iso2"`
	Iso3        string `json:"iso3"`
	UNCode      int16  `json:"unCode"`
	Zone        string `json:"zone"` // N = National, U = European Union, E = Export
	PhonePrefix int16  `json:"phonePrefix"`
	Language    *int32 `json:"language"`
	Currency    *int32 `json:"currency"`
	enterprise  int32
}

func getCountries(enterpriseId int32) []Country {
	var countries []Country = make([]Country, 0)
	sqlStatement := `SELECT * FROM public.country WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return countries
	}
	defer rows.Close()

	for rows.Next() {
		c := Country{}
		rows.Scan(&c.Id, &c.Name, &c.Iso2, &c.Iso3, &c.UNCode, &c.Zone, &c.PhonePrefix, &c.Language, &c.Currency, &c.enterprise)
		countries = append(countries, c)
	}

	return countries
}

func getCountryRow(id int32, enterpriseId int32) Country {
	sqlStatement := `SELECT * FROM public.country WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Country{}
	}

	c := Country{}
	row.Scan(&c.Id, &c.Name, &c.Iso2, &c.Iso3, &c.UNCode, &c.Zone, &c.PhonePrefix, &c.Language, &c.Currency, &c.enterprise)

	return c
}

func (c *Country) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 75 || len(c.Iso2) != 2 || (len(c.Iso3) != 0 && len(c.Iso3) != 3) || c.UNCode < 0 || (c.Zone != "N" && c.Zone != "U" && c.Zone != "E") || c.PhonePrefix < 0)
}

func searchCountries(search string, enterpriseId int32) []Country {
	var countries []Country = make([]Country, 0)
	sqlStatement := `SELECT * FROM public.country WHERE (name ILIKE $1 OR iso_2 = UPPER($2) OR iso_3 = UPPER($2)) AND enterprise=$3 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%", search, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return countries
	}
	defer rows.Close()

	for rows.Next() {
		c := Country{}
		rows.Scan(&c.Id, &c.Name, &c.Iso2, &c.Iso3, &c.UNCode, &c.Zone, &c.PhonePrefix, &c.Language, &c.Currency, &c.enterprise)
		countries = append(countries, c)
	}

	return countries
}

func (c *Country) insertCountry() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.country(name, iso_2, iso_3, un_code, zone, phone_prefix, language, currency, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	res, err := db.Exec(sqlStatement, c.Name, c.Iso2, c.Iso3, c.UNCode, c.Zone, c.PhonePrefix, c.Language, c.Currency, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Country) updateCountry() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.country SET name=$2, iso_2=$3, iso_3=$4, un_code=$5, zone=$6, phone_prefix=$7, language=$8, currency=$9 WHERE id=$1 AND enterprise=$10`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Iso2, c.Iso3, c.UNCode, c.Zone, c.PhonePrefix, c.Language, c.Currency, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Country) deleteCountry() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.country WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCountryByName(languageName string, enterpriseId int32) []NameInt16 {
	var countries []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.country WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return countries
	}
	defer rows.Close()

	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		countries = append(countries, c)
	}

	return countries
}

func getNameCountry(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.country WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
