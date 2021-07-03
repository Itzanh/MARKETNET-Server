package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Currency struct {
	Id           int16     `json:"id"`
	Name         string    `json:"name"`
	Sign         string    `json:"sign"`
	IsoCode      string    `json:"isoCode"`
	IsoNum       int16     `json:"isoNum"`
	Change       float32   `json:"change"`
	ExchangeDate time.Time `json:"exchangeDate"`
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
		rows.Scan(&c.Id, &c.Name, &c.Sign, &c.IsoCode, &c.IsoNum, &c.Change, &c.ExchangeDate)
		currencies = append(currencies, c)
	}

	return currencies
}

func (c *Currency) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 75 || len([]rune(c.Sign)) > 3 || len(c.IsoCode) > 3 || c.IsoNum < 0 || c.Change < 0)
}

func (c *Currency) insertCurrency() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.currency(name, sign, iso_code, iso_num, exchange) VALUES ($1, $2, $3, $4, $5)`
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

	sqlStatement := `UPDATE public.currency SET name=$2, sign=$3, iso_code=$4, iso_num=$5, exchange=$6 WHERE id=$1`
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
	sqlStatement := `SELECT exchange FROM public.currency WHERE id=$1`
	row := db.QueryRow(sqlStatement, currencyId)
	if row.Err() != nil {
		return 0
	}
	var change float32
	row.Scan(&change)
	return change
}

func findCurrencyByName(currencyName string) []NameInt16 {
	var currencies []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.currency WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(currencyName))
	if err != nil {
		return currencies
	}
	for rows.Next() {
		c := NameInt16{}
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

func updateCurrencyExchange() {
	if getSettingsRecord().Currency != "E" {
		return
	}
	currencies := getCurrencies()

	for i := 0; i < len(currencies); i++ {
		if len(currencies[i].IsoCode) == 0 || currencies[i].IsoCode == "EUR" {
			continue
		}

		now := time.Now()
		now = now.AddDate(0, 0, -1)
		currentDate := now.Format("2006-01-02")
		resp, err := http.Get(getSettingsRecord().CurrencyECBurl + "D." + currencies[i].IsoCode + ".EUR.SP00.A?startPeriod=" + currentDate + "&endPeriod=" + currentDate)
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		if len(body) == 0 { // there is no data for this date
			continue
		}
		xml := string(body)

		// check that there is a field with the date
		dateIndexOf := strings.Index(xml, "<generic:ObsDimension value=\"")
		if dateIndexOf <= 0 {
			continue
		}
		dateIndexOf += len("<generic:ObsDimension value=\"")

		// chech that there is a field with the exchange rate
		valueIndexOf := strings.Index(xml, "<generic:ObsValue value=\"")
		if valueIndexOf <= 0 {
			continue
		}
		valueIndexOf += len("<generic:ObsValue value=\"")

		// extract the field and parse the data
		dateString := xml[dateIndexOf : dateIndexOf+10]
		valueString := xml[valueIndexOf : valueIndexOf+strings.Index(xml[valueIndexOf:], "\"")]

		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			continue
		}
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			continue
		}

		sqlStatement := `UPDATE currency SET exchange=$2,exchange_date=$3 WHERE id=$1`
		db.Exec(sqlStatement, currencies[i].Id, value, date)
	}
}
