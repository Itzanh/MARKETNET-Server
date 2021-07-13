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
	initialStatesData()
	initialColorData()
	initialIncotermData()
	initialWarehouseData()
	initiaBillingSeriesData()
	initialJournals()
	initialAccount()
	initialConfig()
	initialUser()
	initialGroup()
	initialUserGroup()
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

func initialStatesData() {
	sqlStatement := `SELECT COUNT(*) FROM state`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/states.json")
		if err != nil {
			return
		}

		var state []State
		json.Unmarshal(content, &state)
		for i := 0; i < len(state); i++ {
			state[i].insertState()
		}
		fmt.Println("INITIAL DATA: Generated states data")
	}
}

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

func initiaBillingSeriesData() {
	sqlStatement := `SELECT COUNT(*) FROM billing_series`
	row := db.QueryRow(sqlStatement)
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
			billingSerie[i].insertBillingSerie()
		}
		fmt.Println("INITIAL DATA: Generated billing series data")
	}
}

func initialJournals() {
	sqlStatement := `SELECT COUNT(*) FROM journal`
	row := db.QueryRow(sqlStatement)
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
			journal[i].insertJournal()
		}

		fmt.Println("INITIAL DATA: Generated journal data")
	}
}

func initialAccount() {
	sqlStatement := `SELECT COUNT(*) FROM account`
	row := db.QueryRow(sqlStatement)
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
			account[i].insertAccount()
		}

		fmt.Println("INITIAL DATA: Generated accounts data")
	}
}

func initialConfig() {
	sqlStatement := `SELECT COUNT(*) FROM config WHERE id=1`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		content, err := ioutil.ReadFile("./initial_data/config.json")
		if err != nil {
			return
		}

		var config Settings
		json.Unmarshal(content, &config)

		var salesAccount *int32
		if config.SalesJournal != nil && *config.SalesJournal > 0 {
			acc := getAccountIdByAccountNumber(*config.SalesJournal, 1)
			if acc > 0 {
				salesAccount = &acc
			}
		}
		var purchaseAccount *int32
		if config.PurchaseJournal != nil && *config.PurchaseJournal > 0 {
			acc := getAccountIdByAccountNumber(*config.PurchaseJournal, 1)
			if acc > 0 {
				purchaseAccount = &acc
			}
		}

		sqlStatement := `INSERT INTO public.config(id, default_vat_percent, default_warehouse, date_format, enterprise_name, enterprise_description, ecommerce, email, currency, currency_ecb_url, barcode_prefix, prestashop_url, prestashop_api_key, prestashop_language_id, prestashop_export_serie, prestashop_intracommunity_serie, prestashop_interior_serie, cron_currency, cron_prestashop, sendgrid_key, email_from, name_from, pallet_weight, pallet_width, pallet_height, pallet_depth, max_connections, customer_journal, sales_journal, sales_account, supplier_journal, purchase_journal, purchase_account) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)`
		_, err = db.Exec(sqlStatement, config.DefaultVatPercent, config.DefaultWarehouse, config.DateFormat, config.EnterpriseName, config.EnterpriseDescription, config.Ecommerce, config.Email, config.Currency, config.CurrencyECBurl, config.BarcodePrefix, config.PrestaShopUrl, config.PrestaShopApiKey, config.PrestaShopLanguageId, config.PrestaShopExportSerie, config.PrestaShopIntracommunitySerie, config.PrestaShopInteriorSerie, config.CronCurrency, config.CronPrestaShop, config.SendGridKey, config.EmailFrom, config.NameFrom, config.PalletWeight, config.PalletWidth, config.PalletHeight, config.PalletDepth, config.MaxConnections, config.CustomerJournal, config.SalesJournal, salesAccount, config.SupplierJournal, config.PurchaseJournal, purchaseAccount)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("INITIAL DATA: Generated config data")
	}
}

func initialUser() {
	sqlStatement := `SELECT COUNT(*) FROM "user" WHERE username='marketnet'`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		u := UserInsert{Username: "marketnet", FullName: "MARKETNET ADMINISTRATOR", Password: "admin1234", Language: "en"} // INITIAL PASSWORD, USER MUST CHANGE THIS!!!
		u.insertUser()

		fmt.Println("INITIAL DATA: Generated admin user")
	}
}

func initialGroup() {
	sqlStatement := `SELECT COUNT(*) FROM "group" WHERE name='Administrators'`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		g := Group{Name: "Administrators", Sales: true, Purchases: true, Masters: true, Warehouse: true, Manufacturing: true, Preparation: true, Admin: true, PrestaShop: true, Accounting: true}
		g.insertGroup()

		fmt.Println("INITIAL DATA: Generated admin group")
	}
}

func initialUserGroup() {
	sqlStatement := `SELECT COUNT(*) FROM user_group INNER JOIN "user" ON "user".id=user_group."user" INNER JOIN "group" ON "group".id=user_group."group" WHERE "user".username = 'marketnet' AND "group".name = 'Administrators'`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {

		sqlStatement := `SELECT id FROM "user" WHERE username='marketnet'`
		row := db.QueryRow(sqlStatement)
		var userId int16
		row.Scan(&userId)

		sqlStatement = `SELECT id FROM "group" WHERE name='Administrators'`
		row = db.QueryRow(sqlStatement)
		var groupId int16
		row.Scan(&groupId)

		ug := UserGroup{}
		ug.User = userId
		ug.Group = groupId
		ug.insertUserGroup()

		fmt.Println("INITIAL DATA: Added the admin user to the admin group")
	}
}
