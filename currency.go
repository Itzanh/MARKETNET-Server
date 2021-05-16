package main

import "strings"

type Currency struct {
	Id      int16   `json:"id"`
	Name    string  `json:"name"`
	Sign    string  `json:"sign"`
	IsoCode string  `json:"isoCode"`
	IsoNum  int16   `json:"isoNum"`
	Change  float32 `json:"change"`
}

func getCurrencies() []Currency {
	var currencies []Currency = make([]Currency, 0)
	sqlStatement := `SELECT * FROM public.currency ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return currencies
	}
	for rows.Next() {
		c := Currency{}
		rows.Scan(&c.Id, &c.Name, &c.Sign, &c.IsoCode, &c.IsoNum, &c.Change)
		currencies = append(currencies, c)
	}

	return currencies
}

func (c *Currency) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || len([]rune(c.Sign)) > 3 || len(c.IsoCode) > 3 || c.IsoNum <= 0 || c.Change <= 0)
}

func (c *Currency) insertCurrency() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.currency(name, sign, iso_code, iso_num, change) VALUES ($1, $2, $3, $4, $5)`
	res, err := db.Exec(sqlStatement, c.Name, c.Sign, c.IsoCode, c.IsoNum, c.Change)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Currency) updateCurrency() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.currency SET name=$2, sign=$3, iso_code=$4, iso_num=$5, change=$6 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Sign, c.IsoCode, c.IsoNum, c.Change)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Currency) deleteCurrency() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.currency WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func getCurrencyExchange(currencyId int16) float32 {
	sqlStatement := `SELECT change FROM public.currency WHERE id = $1`
	row := db.QueryRow(sqlStatement, currencyId)
	if row.Err() != nil {
		return 0
	}
	var change float32
	row.Scan(&change)
	return change
}

type CurrencyName struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

func findCurrencyByName(currencyName string) []CurrencyName {
	var currencies []CurrencyName = make([]CurrencyName, 0)
	sqlStatement := `SELECT id,name FROM public.currency WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(currencyName))
	if err != nil {
		return currencies
	}
	for rows.Next() {
		c := CurrencyName{}
		rows.Scan(&c.Id, &c.Name)
		currencies = append(currencies, c)
	}

	return currencies
}

func getNameCurrency(id int16) string {
	sqlStatement := `SELECT name FROM public.currency WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
