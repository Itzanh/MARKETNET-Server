/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func initialData(enterpriseId int32) {
	initialPaymentData(enterpriseId)
	initialLanguageData(enterpriseId)
	initialCurrenciesData(enterpriseId)
	initialCountriesData(enterpriseId)
	initialStatesData(enterpriseId)
	initialColorData(enterpriseId)
	initialIncotermData(enterpriseId)
	initialWarehouseData(enterpriseId)
	initiaBillingSeriesData(enterpriseId)
	initialJournals(enterpriseId)
	initialAccount(enterpriseId)
	initialReportTemplate(enterpriseId)
	initialPermissionDictionary(enterpriseId)
	initialHSCodes()
}

func initialPaymentData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM payment_method WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			paymentMethod[i].EnterpriseId = enterpriseId
			paymentMethod[i].insertPaymentMethod()
		}
		fmt.Println("INITIAL DATA: Generated payment methods data")
	}
}

func initialLanguageData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM language WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			language[i].EnterpriseId = enterpriseId
			language[i].insertLanguage()
		}
		fmt.Println("INITIAL DATA: Generated language data")
	}
}

func initialCurrenciesData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM currency WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			currencies[i].EnterpriseId = enterpriseId
			currencies[i].insertCurrency()
		}
		fmt.Println("INITIAL DATA: Generated currency data")
	}
}

type CountryInitialData struct {
	Id          int32   `json:"id"`
	Name        string  `json:"name"`
	Iso2        string  `json:"iso2"`
	Iso3        string  `json:"iso3"`
	UNCode      int16   `json:"unCode"`
	Zone        string  `json:"zone"` // N = National, U = European Union, E = Export
	PhonePrefix int16   `json:"phonePrefix"`
	Language    *string `json:"language"`
	Currency    *string `json:"currency"`
}

func initialCountriesData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM country WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/countries.json")
		if err != nil {
			return
		}

		var country []CountryInitialData
		json.Unmarshal(content, &country)
		for i := 0; i < len(country); i++ {
			c := Country{
				Name:         country[i].Name,
				Iso2:         country[i].Iso2,
				Iso3:         country[i].Iso3,
				UNCode:       country[i].UNCode,
				Zone:         country[i].Zone,
				PhonePrefix:  country[i].PhonePrefix,
				EnterpriseId: enterpriseId,
			}

			if country[i].Language != nil {
				sqlStatement := `SELECT id FROM language WHERE iso_2=$1 AND enterprise=$2`
				row := db.QueryRow(sqlStatement, country[i].Language, enterpriseId)

				row.Scan(&c.Language)
			}

			if country[i].Currency != nil {
				sqlStatement := `SELECT id FROM currency WHERE iso_code=$1 AND enterprise=$2`
				row := db.QueryRow(sqlStatement, country[i].Currency, enterpriseId)

				row.Scan(&c.Currency)
			}

			c.insertCountry()
		}
		fmt.Println("INITIAL DATA: Generated countries data")
	}
}

type StateInitialData struct {
	Id      int32  `json:"id"`
	Country string `json:"country"`
	Name    string `json:"name"`
	IsoCode string `json:"isoCode"`
}

func initialStatesData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM state WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/states.json")
		if err != nil {
			return
		}

		var states []StateInitialData
		json.Unmarshal(content, &states)

		sqlStatement := `SELECT id FROM country WHERE iso_2=$1 AND enterprise=$2`

		for i := 0; i < len(states); i++ {
			row := db.QueryRow(sqlStatement, states[i].Country, enterpriseId)

			var countryId int32
			row.Scan(&countryId)

			state := State{
				Name:         states[i].Name,
				IsoCode:      states[i].IsoCode,
				CountryId:    countryId,
				EnterpriseId: enterpriseId,
			}
			state.insertState()
		}
		fmt.Println("INITIAL DATA: Generated states data")
	}
}

func initialColorData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM color WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			color[i].EnterpriseId = enterpriseId
			color[i].insertColor()
		}
		fmt.Println("INITIAL DATA: Generated colors data")
	}
}

func initialIncotermData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM incoterm WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			incoterms[i].EnterpriseId = enterpriseId
			incoterms[i].insertIncoterm()
		}
		fmt.Println("INITIAL DATA: Generated incoterms data")
	}
}

func initialWarehouseData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM warehouse WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			warehouse[i].EnterpriseId = enterpriseId
			warehouse[i].insertWarehouse()
		}
		fmt.Println("INITIAL DATA: Generated warehouse data")
	}
}

func initiaBillingSeriesData(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM billing_series WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/billing_series.json")
		if err != nil {
			return
		}

		var billingSerie []BillingSerie
		json.Unmarshal(content, &billingSerie)
		for i := 0; i < len(billingSerie); i++ {
			billingSerie[i].EnterpriseId = enterpriseId
			billingSerie[i].insertBillingSerie()
		}
		fmt.Println("INITIAL DATA: Generated billing series data")
	}
}

func initialJournals(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM journal WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/journal.json")
		if err != nil {
			return
		}

		var journal []Journal
		json.Unmarshal(content, &journal)
		for i := 0; i < len(journal); i++ {
			journal[i].EnterpriseId = enterpriseId
			journal[i].insertJournal()
		}

		fmt.Println("INITIAL DATA: Generated journal data")
	}
}

func initialAccount(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM account WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/accounts.json")
		if err != nil {
			return
		}

		var account []Account
		json.Unmarshal(content, &account)
		for i := 0; i < len(account); i++ {
			account[i].EnterpriseId = enterpriseId
			account[i].insertAccount()
		}

		fmt.Println("INITIAL DATA: Generated accounts data")
	}
}

func initialConfigCreateEnterprise(enterpriseName string, enterpriseDescription string, enterpriseKey string) (bool, int32) {
	content, err := ioutil.ReadFile("./initial_data/config.json")
	if err != nil {
		fmt.Println(err)
		return false, 0
	}

	var config Settings
	json.Unmarshal(content, &config)

	config.EnterpriseName = enterpriseName
	config.EnterpriseDescription = enterpriseDescription
	config.EnterpriseKey = enterpriseKey

	dbOrm.Create(&config)

	return true, config.Id
}

func initialReportTemplate(enterpriseId int32) {
	content, err := ioutil.ReadFile("./reports/sales_order.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "SALES_ORDER", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/sales_invoice.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "SALES_INVOICE", Html: string(content)}.insertReportTemplate()
	ReportTemplate{EnterpriseId: enterpriseId, Key: "SALES_INVOICE_TICKET", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/sales_delivery_note.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "SALES_DELIVERY_NOTE", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/purchase_order.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "PURCHASE_ORDER", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/box_content.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "BOX_CONTENT", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/pallet_content.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "PALLET_CONTENT", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/carrier_pallet.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "CARRIER_PALLET", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/sales_order_digital_product_data.html")
	if err != nil {
		return
	}
	ReportTemplate{EnterpriseId: enterpriseId, Key: "SALES_ORDER_DIGITAL_PRODUCT_DATA", Html: string(content)}.insertReportTemplate()
}

// check every permission in the initial data file agains the ones in the database
// if a permission exists in both the database and the initial data file, do noghting
// if a permissions exists in the initial data file and not in the database, create it in the DDBB
// if a permission exists in the database but not in the initial data file, remove it from the database
func initialPermissionDictionary(enterpriseId int32) {
	content, err := ioutil.ReadFile("./initial_data/permission_dictionary.json")
	if err != nil {
		return
	}

	var initialPerm []PermissionDictionary = make([]PermissionDictionary, 0)
	json.Unmarshal(content, &initialPerm)
	perm := getPermissionDictionary(enterpriseId)

	for i := 0; i < len(initialPerm); i++ {
		// initial permission key
		ipk := initialPerm[i].Key
		var found bool = false

		for j := 0; j < len(perm); j++ {
			if perm[j].Key == ipk {
				found = true
				perm = append(perm[:j], perm[j+1:]...) // if one permission exists in the DB and not in the file, the final array wont't be empty
				break
			}
		}

		if !found {
			sqlStatement := `INSERT INTO public.permission_dictionary(enterprise, key, description) VALUES ($1, $2, $3)`
			db.Exec(sqlStatement, enterpriseId, initialPerm[i].Key, initialPerm[i].Description)
		}
	}

	for i := 0; i < len(perm); i++ {
		sqlStatement := `DELETE FROM public.permission_dictionary_group WHERE permission_key = $1 AND enterprise = $2`
		db.Exec(sqlStatement, perm[i].Key, enterpriseId)
		sqlStatement = `DELETE FROM public.permission_dictionary WHERE enterprise = $1 AND key = $2`
		db.Exec(sqlStatement, enterpriseId, perm[i].Key)
	}
}

func initialHSCodes() {
	sqlStatement := `SELECT COUNT(*) FROM hs_codes`
	row := db.QueryRow(sqlStatement)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return
	}

	var rowCount int
	row.Scan(&rowCount)

	if rowCount > 0 {
		return
	}
}
