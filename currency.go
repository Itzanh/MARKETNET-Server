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
	Id           int32     `json:"id"`
	Name         string    `json:"name"`
	Sign         string    `json:"sign"`
	IsoCode      string    `json:"isoCode"`
	IsoNum       int16     `json:"isoNum"`
	Change       float64   `json:"change"`
	ExchangeDate time.Time `json:"exchangeDate"`
	enterprise   int32
}

func getCurrencies(enterpriseId int32) []Currency {
	var currencies []Currency = make([]Currency, 0)
	sqlStatement := `SELECT * FROM public.currency WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return currencies
	}
	for rows.Next() {
		c := Currency{}
		rows.Scan(&c.Id, &c.Name, &c.Sign, &c.IsoCode, &c.IsoNum, &c.Change, &c.ExchangeDate, &c.enterprise)
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

	sqlStatement := `INSERT INTO public.currency(name, sign, iso_code, iso_num, exchange, enterprise) VALUES ($1, $2, $3, $4, $5, $6)`
	res, err := db.Exec(sqlStatement, c.Name, c.Sign, c.IsoCode, c.IsoNum, c.Change, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Currency) updateCurrency() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.currency SET name=$2, sign=$3, iso_code=$4, iso_num=$5, exchange=$6 WHERE id=$1 AND enterprise=$7`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Sign, c.IsoCode, c.IsoNum, c.Change, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Currency) deleteCurrency() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.currency WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func getCurrencyExchange(currencyId int32) float64 {
	sqlStatement := `SELECT exchange FROM public.currency WHERE id=$1`
	row := db.QueryRow(sqlStatement, currencyId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}
	var change float64
	row.Scan(&change)
	return change
}

func findCurrencyByName(currencyName string, enterpriseId int32) []NameInt16 {
	var currencies []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.currency WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(currencyName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return currencies
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		currencies = append(currencies, c)
	}

	return currencies
}

func getNameCurrency(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.currency WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type LocateCurrency struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func locateCurrency(enterpriseId int32) []LocateCurrency {
	var currencies []LocateCurrency = make([]LocateCurrency, 0)
	sqlStatement := `SELECT id,name FROM public.currency WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return currencies
	}
	for rows.Next() {
		c := LocateCurrency{}
		rows.Scan(&c.Id, &c.Name)
		currencies = append(currencies, c)
	}

	return currencies
}

func updateCurrencyExchange(enterpriseId int32) {
	if getSettingsRecordById(enterpriseId).Currency != "E" {
		return
	}
	currencies := getCurrencies(enterpriseId)

	for i := 0; i < len(currencies); i++ {
		if len(currencies[i].IsoCode) == 0 || currencies[i].IsoCode == "EUR" {
			continue
		}

		now := time.Now()
		now = now.AddDate(0, 0, -1)
		currentDate := now.Format("2006-01-02")
		resp, err := http.Get(getSettingsRecordById(enterpriseId).CurrencyECBurl + "D." + currencies[i].IsoCode + ".EUR.SP00.A?startPeriod=" + currentDate + "&endPeriod=" + currentDate)
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
