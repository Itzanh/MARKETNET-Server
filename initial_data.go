package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func initialData() {
	initialPaymentData()
	initialLanguageData()
	initialCurrenciesData()
	initialCountriesData()
	//initialCitiesData()
	initialColorData()
	initialIncotermData()
	initialWarehouseData()
}

func initialPaymentData() {
	sqlStatement := `SELECT COUNT(*) FROM payment_method`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/payment_methods.json")
		if err != nil {
			return
		}

		var paymentMethod []PaymentMethod
		json.Unmarshal(content, &paymentMethod)
		for i := 0; i < len(paymentMethod); i++ {
			paymentMethod[i].insertPaymentMethod()
		}
		fmt.Println("INITIAL DATA: Generated payment methods data")
	}
}

func initialLanguageData() {
	sqlStatement := `SELECT COUNT(*) FROM language`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/languages.json")
		if err != nil {
			return
		}

		var language []Language
		json.Unmarshal(content, &language)
		for i := 0; i < len(language); i++ {
			language[i].insertLanguage()
		}
		fmt.Println("INITIAL DATA: Generated language data")
	}
}

func initialCurrenciesData() {
	sqlStatement := `SELECT COUNT(*) FROM currency`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/currencies.json")
		if err != nil {
			return
		}

		var currencies []Currency
		json.Unmarshal(content, &currencies)
		for i := 0; i < len(currencies); i++ {
			currencies[i].insertCurrency()
		}
		fmt.Println("INITIAL DATA: Generated currency data")
	}
}

func initialCountriesData() {
	sqlStatement := `SELECT COUNT(*) FROM country`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/countries.json")
		if err != nil {
			return
		}

		var country []Country
		json.Unmarshal(content, &country)
		for i := 0; i < len(country); i++ {
			country[i].insertCountry()
		}
		fmt.Println("INITIAL DATA: Generated countries data")
	}
}

/*func initialCitiesData() {
	sqlStatement := `SELECT COUNT(*) FROM city`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/cities.json")
		if err != nil {
			return
		}

		var cities []State
		json.Unmarshal(content, &cities)
		for i := 0; i < len(cities); i++ {
			cities[i].insertCity()
		}
		fmt.Println("INITIAL DATA: Generated cities data")
	}
}*/

func initialColorData() {
	sqlStatement := `SELECT COUNT(*) FROM color`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/colors.json")
		if err != nil {
			return
		}

		var color []Color
		json.Unmarshal(content, &color)
		for i := 0; i < len(color); i++ {
			color[i].insertColor()
		}
		fmt.Println("INITIAL DATA: Generated colors data")
	}
}

func initialIncotermData() {
	sqlStatement := `SELECT COUNT(*) FROM incoterm`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/incoterms.json")
		if err != nil {
			return
		}

		var incoterms []Incoterm
		json.Unmarshal(content, &incoterms)
		for i := 0; i < len(incoterms); i++ {
			incoterms[i].insertIncoterm()
		}
		fmt.Println("INITIAL DATA: Generated incoterms data")
	}
}

func initialWarehouseData() {
	sqlStatement := `SELECT COUNT(*) FROM warehouse`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/warehouse.json")
		if err != nil {
			return
		}

		var warehouse []Warehouse
		json.Unmarshal(content, &warehouse)
		for i := 0; i < len(warehouse); i++ {
			warehouse[i].insertWarehouse()
		}
		fmt.Println("INITIAL DATA: Generated warehouse data")
	}
}
