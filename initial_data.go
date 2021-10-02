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
	initialConfig(enterpriseId)
	initialUser(enterpriseId)
	initialGroup(enterpriseId)
	initialUserGroup()
	initialReportTemplate(enterpriseId)
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
			paymentMethod[i].enterprise = enterpriseId
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
			language[i].enterprise = enterpriseId
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
			currencies[i].enterprise = enterpriseId
			currencies[i].insertCurrency()
		}
		fmt.Println("INITIAL DATA: Generated currency data")
	}
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

		var country []Country
		json.Unmarshal(content, &country)
		for i := 0; i < len(country); i++ {
			country[i].enterprise = enterpriseId
			country[i].insertCountry()
		}
		fmt.Println("INITIAL DATA: Generated countries data")
	}
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

		var state []State
		json.Unmarshal(content, &state)
		for i := 0; i < len(state); i++ {
			state[i].enterprise = enterpriseId
			state[i].insertState()
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
			color[i].enterprise = enterpriseId
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
			incoterms[i].enterprise = enterpriseId
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
			warehouse[i].enterprise = enterpriseId
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
			billingSerie[i].enterprise = enterpriseId
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
			journal[i].enterprise = enterpriseId
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
			account[i].enterprise = enterpriseId
			account[i].insertAccount()
		}

		fmt.Println("INITIAL DATA: Generated accounts data")
	}
}

func initialConfig(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM config WHERE id=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
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
			acc := getAccountIdByAccountNumber(*config.SalesJournal, 1, enterpriseId)
			if acc > 0 {
				salesAccount = &acc
			}
		}
		var purchaseAccount *int32
		if config.PurchaseJournal != nil && *config.PurchaseJournal > 0 {
			acc := getAccountIdByAccountNumber(*config.PurchaseJournal, 1, enterpriseId)
			if acc > 0 {
				purchaseAccount = &acc
			}
		}

		sqlStatement := `INSERT INTO public.config(id, default_vat_percent, default_warehouse, date_format, enterprise_name, enterprise_description, ecommerce, email, currency, currency_ecb_url, barcode_prefix, prestashop_url, prestashop_api_key, prestashop_language_id, prestashop_export_serie, prestashop_intracommunity_serie, prestashop_interior_serie, cron_currency, cron_prestashop, sendgrid_key, email_from, name_from, pallet_weight, pallet_width, pallet_height, pallet_depth, max_connections, customer_journal, sales_journal, sales_account, supplier_journal, purchase_journal, purchase_account, enterprise_key) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33)`
		_, err = db.Exec(sqlStatement, config.DefaultVatPercent, config.DefaultWarehouse, config.DateFormat, config.EnterpriseName, config.EnterpriseDescription, config.Ecommerce, config.Email, config.Currency, config.CurrencyECBurl, config.BarcodePrefix, config.PrestaShopUrl, config.PrestaShopApiKey, config.PrestaShopLanguageId, config.PrestaShopExportSerie, config.PrestaShopIntracommunitySerie, config.PrestaShopInteriorSerie, config.CronCurrency, config.CronPrestaShop, config.SendGridKey, config.EmailFrom, config.NameFrom, config.PalletWeight, config.PalletWidth, config.PalletHeight, config.PalletDepth, config.MaxConnections, config.CustomerJournal, config.SalesJournal, salesAccount, config.SupplierJournal, config.PurchaseJournal, purchaseAccount, config.EnterpriseKey)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("INITIAL DATA: Generated config data")
	}
}

func initialConfigCreateEnterprise(enterpriseName string, enterpriseDescription string, enterpriseKey string) (bool, int32) {
	content, err := ioutil.ReadFile("./initial_data/config.json")
	if err != nil {
		return false, 0
	}

	var config Settings
	json.Unmarshal(content, &config)

	var salesAccount *int32
	var purchaseAccount *int32

	sqlStatement := `INSERT INTO public.config(id, default_vat_percent, default_warehouse, date_format, enterprise_name, enterprise_description, ecommerce, email, currency, currency_ecb_url, barcode_prefix, prestashop_url, prestashop_api_key, prestashop_language_id, prestashop_export_serie, prestashop_intracommunity_serie, prestashop_interior_serie, cron_currency, cron_prestashop, sendgrid_key, email_from, name_from, pallet_weight, pallet_width, pallet_height, pallet_depth, max_connections, customer_journal, sales_journal, sales_account, supplier_journal, purchase_journal, purchase_account, enterprise_key) VALUES (3, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33) RETURNING id`
	row := db.QueryRow(sqlStatement, config.DefaultVatPercent, nil, config.DateFormat, enterpriseName, enterpriseDescription, config.Ecommerce, config.Email, config.Currency, config.CurrencyECBurl, config.BarcodePrefix, config.PrestaShopUrl, config.PrestaShopApiKey, config.PrestaShopLanguageId, nil, nil, nil, config.CronCurrency, config.CronPrestaShop, config.SendGridKey, config.EmailFrom, config.NameFrom, config.PalletWeight, config.PalletWidth, config.PalletHeight, config.PalletDepth, config.MaxConnections, nil, nil, salesAccount, nil, nil, purchaseAccount, enterpriseKey)
	if row.Err() != nil {
		fmt.Println(row.Err())
	}

	var enterpriseId int32
	row.Scan(&enterpriseId)
	return true, enterpriseId
}

func initialUser(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM "user" WHERE username='marketnet'`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		u := UserInsert{Username: "marketnet", FullName: "MARKETNET ADMINISTRATOR", Password: "admin1234", Language: "en"} // INITIAL PASSWORD, USER MUST CHANGE THIS!!!
		u.insertUser(enterpriseId)

		fmt.Println("INITIAL DATA: Generated admin user")
	}
}

func initialGroup(enterpriseId int32) {
	sqlStatement := `SELECT COUNT(*) FROM "group" WHERE name='Administrators'`
	row := db.QueryRow(sqlStatement)
	var rows int32
	row.Scan(&rows)

	if rows == 0 {
		g := Group{Name: "Administrators", Sales: true, Purchases: true, Masters: true, Warehouse: true, Manufacturing: true, Preparation: true, Admin: true, PrestaShop: true, Accounting: true, enterprise: enterpriseId}
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
		var userId int32
		row.Scan(&userId)

		sqlStatement = `SELECT id FROM "group" WHERE name='Administrators'`
		row = db.QueryRow(sqlStatement)
		var groupId int32
		row.Scan(&groupId)

		ug := UserGroup{}
		ug.User = userId
		ug.Group = groupId
		ug.insertUserGroup()

		fmt.Println("INITIAL DATA: Added the admin user to the admin group")
	}
}

func initialReportTemplate(enterpriseId int32) {
	templates := getReportTemplates(enterpriseId)

	if len(templates) > 0 {
		return
	}

	content, err := ioutil.ReadFile("./reports/sales_order.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "SALES_ORDER", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/sales_invoice.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "SALES_INVOICE", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/sales_delivery_note.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "SALES_DELIVERY_NOTE", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/purchase_order.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "PURCHASE_ORDER", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/box_content.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "BOX_CONTENT", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/pallet_content.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "PALLET_CONTENT", Html: string(content)}.insertReportTemplate()

	content, err = ioutil.ReadFile("./reports/carrier_pallet.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "CARRIER_PALLET", Html: string(content)}.insertReportTemplate()
}
