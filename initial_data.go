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
	initialPermissionDictionary(enterpriseId)
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
				Name:        country[i].Name,
				Iso2:        country[i].Iso2,
				Iso3:        country[i].Iso3,
				UNCode:      country[i].UNCode,
				Zone:        country[i].Zone,
				PhonePrefix: country[i].PhonePrefix,
				enterprise:  enterpriseId,
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
				Name:       states[i].Name,
				IsoCode:    states[i].IsoCode,
				Country:    countryId,
				enterprise: enterpriseId,
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
		g := Group{Name: "Administrators", Sales: true, Purchases: true, Masters: true, Warehouse: true, Manufacturing: true, Preparation: true, Admin: true, PrestaShop: true, Accounting: true, enterprise: enterpriseId, PointOfSale: true}
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
	ReportTemplate{enterprise: enterpriseId, Key: "SALES_INVOICE_TICKET", Html: string(content)}.insertReportTemplate()

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

	content, err = ioutil.ReadFile("./reports/sales_order_digital_product_data.html")
	if err != nil {
		return
	}
	ReportTemplate{enterprise: enterpriseId, Key: "SALES_ORDER_DIGITAL_PRODUCT_DATA", Html: string(content)}.insertReportTemplate()
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
