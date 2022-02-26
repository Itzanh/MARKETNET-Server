package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/lib/pq"
)

// =====
// GENERIS FUNCTIONS
// =====

func getPrestaShopAPI_URL(resourceName string, enterpriseId int32) string {
	s := getSettingsRecordById(enterpriseId)

	return s.PrestaShopUrl + resourceName + "?ws_key=" + s.PrestaShopApiKey + "&output_format=JSON&language=" + strconv.Itoa(int(s.PrestaShopLanguageId))
}

func getPrestaShopJSON(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

// =====
// DEFINE PRESTASHOP CLASSES
// =====

type PSZones struct {
	Zones []PSZone `json:"zones"`
}

type PSZone struct {
	Id     int32  `json:"id"`
	Name   string `json:"name"`
	Active string `json:"active"`
}

type PSCurrencies struct {
	Currencies []PSCurrency `json:"currencies"`
}

type PSCurrency struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	IsoCode        string `json:"iso_code"`
	NumericIsoCode string `json:"numeric_iso_code"`
	ConversionRate string `json:"conversion_rate"`
	Deleted        string `json:"deleted"`
	Active         string `json:"active"`
}

type PSCountries struct {
	Countries []PSCountry `json:"countries"`
}

type PSCountry struct {
	Id         int32  `json:"id"`
	IdZone     string `json:"id_zone"`
	IdCurrency string `json:"id_currency"`
	IsoCode    string `json:"iso_code"`
	CallPrefix string `json:"call_prefix"`
	Active     string `json:"active"`
	Name       string `json:"name"`
}

type PSStates struct {
	States []PSState `json:"states"`
}

type PSState struct {
	Id        int32  `json:"id"`
	IdZone    string `json:"id_zone"`
	IdCountry string `json:"id_country"`
	IsoCode   string `json:"iso_code"`
	Name      string `json:"name"`
	Active    string `json:"active"`
}

type PSCustomers struct {
	Customers []PSCustomer `json:"customers"`
}

type PSCustomer struct {
	Id        int32  `json:"id"`
	IdLang    string `json:"id_lang"`
	Company   string `json:"company"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Note      string `json:"note"`
	Active    string `json:"active"`
	Deleted   string `json:"deleted"`
	DateAdd   string `json:"date_add"`
	DateUpd   string `json:"date_upd"`
}

type PSAddresses struct {
	Addresses []PSAddress `json:"addresses"`
}

type PSAddress struct {
	Id          int32  `json:"id"`
	IdCustomer  string `json:"id_customer"`
	IdCountry   string `json:"id_country"`
	IdState     string `json:"id_state"`
	Alias       string `json:"alias"`
	Company     string `json:"company"`
	Lastname    string `json:"lastname"`
	Firstname   string `json:"firstname"`
	VatNumber   string `json:"vat_number"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Postcode    string `json:"postcode"`
	City        string `json:"city"`
	Other       string `json:"other"`
	Phone       string `json:"phone"`
	PhoneMobile string `json:"phone_mobile"`
	Dni         string `json:"dni"`
	Deleted     string `json:"deleted"`
	DateAdd     string `json:"date_add"`
	DateUpd     string `json:"date_upd"`
}

type PSProducts struct {
	Products []PSProduct `json:"products"`
}

type PSProduct struct {
	Id          int32  `json:"id"`
	OnSale      string `json:"on_sale"`
	Ean13       string `json:"ean13"`
	Price       string `json:"price"`
	Reference   string `json:"reference"`
	Active      string `json:"active"`
	DateAdd     string `json:"date_add"`
	DateUpd     string `json:"date_upd"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PSProductCombinations struct {
	Combinations []PSProductCombination `json:"combinations"`
}

type PSProductCombination struct {
	Id           int32                            `json:"id"`
	IdProduct    string                           `json:"id_product"`
	Ean13        string                           `json:"ean13"`
	Reference    string                           `json:"reference"`
	Price        string                           `json:"price"`
	Associations PSProductCombinationAssociations `json:"associations"`
}

type PSProductCombinationAssociations struct {
	ProductOptionValues []PSProductOptionValueCombination `json:"product_option_values"`
}

type PSProductOptionValueCombination struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PSProductOptionValues struct {
	ProductOptionValues []PSProductOptionValue `json:"product_option_values"`
}

type PSProductOptionValue struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type PSLanguages struct {
	Languages []PSLanguage `json:"languages"`
}

type PSLanguage struct {
	Id      int32  `json:"id"`
	Name    string `json:"name"`
	IsoCode string `json:"iso_code"`
	Active  string `json:"active"`
}

type PSCarriers struct {
	Carriers []PSCarrier `json:"carriers"`
}

type PSCarrier struct {
	Id        int32  `json:"id"`
	Deleted   string `json:"deleted"`
	Name      string `json:"name"`
	Active    string `json:"active"`
	Url       string `json:"url"`
	MaxWidth  string `json:"max_width"`
	MaxHeight string `json:"max_height"`
	MaxDepth  string `json:"max_depth"`
	MaxWeight string `json:"max_weight"`
}

type PSOrders struct {
	Orders []PSOrder `json:"orders"`
}

type PSOrder struct {
	Id                    int32  `json:"id"`
	IdAddressDelivery     string `json:"id_address_delivery"`
	IdAddressInvoice      string `json:"id_address_invoice"`
	IdCurrency            string `json:"id_currency"`
	IdLang                string `json:"id_lang"`
	IdCustomer            string `json:"id_customer"`
	IdCarrier             string `json:"id_carrier"`
	Module                string `json:"module"`
	DateAdd               string `json:"date_add"`
	DateUpd               string `json:"date_upd"`
	TotalDiscountsTaxExcl string `json:"total_discounts_tax_excl"`
	TotalShippingTaxExcl  string `json:"total_shipping_tax_excl"`
	Reference             string `json:"reference"`
	TotalPaidTaxIncl      string `json:"total_paid_tax_incl"`
	TotalPaidTaxExcl      string `json:"total_paid_tax_excl"`
}

type PSOrderDetails struct {
	OrderDetails []PSOrderDetail `json:"order_details"`
}

type PSOrderDetail struct {
	Id                 int32  `json:"id"`
	IdOrder            string `json:"id_order"`
	ProductId          string `json:"product_id"`
	ProductAttributeId string `json:"product_attribute_id"`
	ProductQuantity    string `json:"product_quantity"`
	ProductPrice       string `json:"product_price"`
}

// main import function
func importFromPrestaShop(enterpriseId int32) {
	s := getSettingsRecordById(enterpriseId)
	if s.Ecommerce != "P" {
		return
	}

	// get all data from PrestaShop, write it in tables like the ones that PrestaShop uses
	if !importPsZones(enterpriseId) {
		return
	}
	if !importPsCurrencies(enterpriseId) {
		return
	}
	if !importPsCountries(enterpriseId) {
		return
	}
	if !importPsStates(enterpriseId) {
		return
	}
	if !importPsCustomers(enterpriseId) {
		return
	}
	if !importPsAddresses(enterpriseId) {
		return
	}
	if !importPsProducts(enterpriseId) {
		return
	}
	if !importPsProductCombinations(enterpriseId) {
		return
	}
	if !importPsProductOptionValues(enterpriseId) {
		return
	}
	if !importPsLanguage(enterpriseId) {
		return
	}
	if !importPsCarriers(enterpriseId) {
		return
	}
	if !importPsOrders(enterpriseId) {
		return
	}
	if !importPsOrderDetails(enterpriseId) {
		return
	}

	// trasnfer the data form the PrestaShop tables to the ERP
	if !copyPsCurrencies(enterpriseId) {
		return
	}
	if !copyPsCountries(enterpriseId) {
		return
	}
	if !copyPsStates(enterpriseId) {
		return
	}
	if !copyPsCustomers(enterpriseId) {
		return
	}
	if !copyPsAddresses(enterpriseId) {
		return
	}
	if !copyPsLanguages(enterpriseId) {
		return
	}
	if !copyPsCarriers(enterpriseId) {
		return
	}
	if !copyPsProducts(enterpriseId) {
		return
	}
	if !copyPsOrders(enterpriseId) {
		return
	}
	if !copyPsOrderDetails(enterpriseId) {
		return
	}
}

// =====
// COPY THE DATA FROM PRESTASHOP TO THE PS MARKETNET TABLES
// =====

func importPsZones(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("zones", enterpriseId) + "&display=[id,name,active]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Zones</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var zones PSZones
	json.Unmarshal(jsonPS, &zones)

	sqlStatement := `UPDATE public.ps_zone SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(zones.Zones); i++ {
		zone := zones.Zones[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_zone WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, zone.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_zone(id, name, active, enterprise) VALUES ($1, $2, $3, $4)`
			db.Exec(sqlStatement, zone.Id, zone.Name, zone.Active)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_zone SET name=$2, active=$3, ps_exists=true WHERE id=$1 AND enterprise=$4`
			db.Exec(sqlStatement, zone.Id, zone.Name, zone.Active, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_zone WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsCurrencies(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("currencies", enterpriseId) + "&display=[id,name,symbol,iso_code,numeric_iso_code,conversion_rate,deleted,active]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Currencies</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var currencies PSCurrencies
	json.Unmarshal(jsonPS, &currencies)

	sqlStatement := `UPDATE public.ps_currency SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(currencies.Currencies); i++ {
		currency := currencies.Currencies[i]
		if currency.NumericIsoCode == "" {
			currency.NumericIsoCode = "0"
		}

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_currency WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, currency.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_currency(id, name, symbol, iso_code, numeric_iso_code, conversion_rate, deleted, active, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			db.Exec(sqlStatement, currency.Id, currency.Name, currency.Symbol, currency.IsoCode, currency.NumericIsoCode, currency.ConversionRate, currency.Deleted, currency.Active, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_currency SET name=$2, symbol=$3, iso_code=$4, numeric_iso_code=$5, conversion_rate=$6, deleted=$7, active=$8, ps_exists=true WHERE id=$1 AND enterprise=$9`
			db.Exec(sqlStatement, currency.Id, currency.Name, currency.Symbol, currency.IsoCode, currency.NumericIsoCode, currency.ConversionRate, currency.Deleted, currency.Active, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_currency WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsCountries(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("countries", enterpriseId) + "&display=[id,id_zone,id_currency,iso_code,call_prefix,name,active]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Countries</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var countries PSCountries
	json.Unmarshal(jsonPS, &countries)

	sqlStatement := `UPDATE public.ps_country SET ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(countries.Countries); i++ {
		country := countries.Countries[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_country WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, country.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_country(id, id_zone, id_currency, iso_code, call_prefix, active, name, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			db.Exec(sqlStatement, country.Id, country.IdZone, country.IdCurrency, country.IsoCode, country.CallPrefix, country.Active, country.Name, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_country SET id_zone=$2, id_currency=$3, iso_code=$4, call_prefix=$5, active=$6, name=$7, ps_exists=true WHERE id=$1 AND enterprise=$8`
			db.Exec(sqlStatement, country.Id, country.IdZone, country.IdCurrency, country.IsoCode, country.CallPrefix, country.Active, country.Name, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_country WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsStates(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("states", enterpriseId) + "&display=[id,id_zone,id_country,iso_code,name,active]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: States</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var states PSStates
	json.Unmarshal(jsonPS, &states)

	sqlStatement := `UPDATE public.ps_state SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(states.States); i++ {
		state := states.States[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_state WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, state.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_state(id, id_country, id_zone, name, iso_code, active, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			db.Exec(sqlStatement, state.Id, state.IdCountry, state.IdZone, state.Name, state.IsoCode, state.Active, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_state SET id_country=$2, id_zone=$3, name=$4, iso_code=$5, active=$6, ps_exists=true WHERE id=$1 AND enterprise=$7`
			db.Exec(sqlStatement, state.Id, state.IdCountry, state.IdZone, state.Name, state.IsoCode, state.Active, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_state WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsCustomers(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("customers", enterpriseId) + "&display=[id,id_lang,company,firstname,lastname,email,note,active,deleted,date_add,date_upd]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var customers PSCustomers
	json.Unmarshal(jsonPS, &customers)

	sqlStatement := `UPDATE public.ps_customer SET ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(customers.Customers); i++ {
		customer := customers.Customers[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_customer WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, customer.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_customer(id, id_lang, company, firstname, lastname, email, note, active, deleted, date_add, date_upd, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
			db.Exec(sqlStatement, customer.Id, customer.IdLang, customer.Company, customer.Firstname, customer.Lastname, customer.Email, customer.Note, customer.Active, customer.Deleted, customer.DateAdd, customer.DateUpd, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_customer SET id_lang=$2, company=$3, firstname=$4, lastname=$5, email=$6, note=$7, active=$8, deleted=$9, date_add=$10, date_upd=$11, ps_exists=true WHERE id=$1 AND enterprise=$12`
			db.Exec(sqlStatement, customer.Id, customer.IdLang, customer.Company, customer.Firstname, customer.Lastname, customer.Email, customer.Note, customer.Active, customer.Deleted, customer.DateAdd, customer.DateUpd, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_customer WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsAddresses(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("addresses", enterpriseId) + "&display=[id,id_customer,id_country,id_state,alias,company,lastname,firstname,vat_number,address1,address2,postcode,city,other,phone,phone_mobile,dni,date_add,date_upd,deleted]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Addresses</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var addresses PSAddresses
	json.Unmarshal(jsonPS, &addresses)

	sqlStatement := `UPDATE public.ps_address SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(addresses.Addresses); i++ {
		address := addresses.Addresses[i]
		if address.IdCustomer == "" || address.IdCustomer == "0" {
			continue
		}

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_address WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, address.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_address(id, id_country, id_state, id_customer, alias, company, lastname, firstname, address1, address2, postcode, city, other, phone, phone_mobile, vat_number, dni, date_add, date_upd, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
			db.Exec(sqlStatement, address.Id, address.IdCountry, address.IdState, address.IdCustomer, address.Alias, address.Company, address.Lastname, address.Firstname, address.Address1, address.Address2, address.Postcode, address.City, address.Other, address.Phone, address.PhoneMobile, address.VatNumber, address.Dni, address.DateAdd, address.DateUpd, address.Deleted)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_address SET id_country=$2, id_state=$3, id_customer=$4, alias=$5, company=$6, lastname=$7, firstname=$8, address1=$9, address2=$10, postcode=$11, city=$12, other=$13, phone=$14, phone_mobile=$15, vat_number=$16, dni=$17, date_add=$18, date_upd=$19, deleted=$20, ps_exists=true WHERE id=$1`
			db.Exec(sqlStatement, address.Id, address.IdCountry, address.IdState, address.IdCustomer, address.Alias, address.Company, address.Lastname, address.Firstname, address.Address1, address.Address2, address.Postcode, address.City, address.Other, address.Phone, address.PhoneMobile, address.VatNumber, address.Dni, address.DateAdd, address.DateUpd, address.Deleted)
		}
	}

	sqlStatement = `DELETE FROM public.ps_address WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsProducts(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("products", enterpriseId) + "&display=[id,name,description,on_sale,ean13,price,reference,active,date_add,date_upd]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Products</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var products PSProducts
	json.Unmarshal(jsonPS, &products)

	sqlStatement := `UPDATE public.ps_product SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(products.Products); i++ {
		product := products.Products[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_product WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, product.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_product(id, on_sale, ean13, price, reference, active, date_add, date_upd, name, description, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
			db.Exec(sqlStatement, product.Id, product.OnSale, product.Ean13, product.Price, product.Reference, product.Active, product.DateAdd, product.DateUpd, product.Name, product.Description, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_product SET on_sale=$2, ean13=$3, price=$4, reference=$5, active=$6, date_add=$7, date_upd=$8, name=$9, description=$10, ps_exists=true WHERE id=$1 AND enterprise=$11`
			db.Exec(sqlStatement, product.Id, product.OnSale, product.Ean13, product.Price, product.Reference, product.Active, product.DateAdd, product.DateUpd, product.Name, product.Description, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_product WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsProductCombinations(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("combinations", enterpriseId) + "&display=full"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Product combinations</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var combinations PSProductCombinations
	json.Unmarshal(jsonPS, &combinations)

	sqlStatement := `UPDATE public.ps_product_combination SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(combinations.Combinations); i++ {
		combination := combinations.Combinations[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_product_combination WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, combination.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		var productOptionValues []int32 = make([]int32, 0)
		for j := 0; j < len(combination.Associations.ProductOptionValues); j++ {
			value := combination.Associations.ProductOptionValues[j]
			id, _ := strconv.Atoi(value.Id)
			productOptionValues = append(productOptionValues, int32(id))
		}

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_product_combination(id, id_product, reference, ean13, price, product_option_values, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			db.Exec(sqlStatement, combination.Id, combination.IdProduct, combination.Reference, combination.Ean13, combination.Price, pq.Array(productOptionValues), enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_product_combination SET id_product=$2, reference=$3, ean13=$4, price=$5, product_option_values=$6, ps_exists=true WHERE id=$1 AND enterprise=$7`
			db.Exec(sqlStatement, combination.Id, combination.IdProduct, combination.Reference, combination.Ean13, combination.Price, pq.Array(productOptionValues), enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_product_combination WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsProductOptionValues(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("product_option_values", enterpriseId) + "&display=[id,name]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Product option values</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var values PSProductOptionValues
	json.Unmarshal(jsonPS, &values)

	sqlStatement := `UPDATE public.product_option_values SET ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(values.ProductOptionValues); i++ {
		value := values.ProductOptionValues[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM product_option_values WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, value.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_product_option_values(id, name, enterprise) VALUES ($1, $2, $3)`
			db.Exec(sqlStatement, value.Id, value.Name, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_product_option_values SET name=$2, ps_exists=true WHERE id=$1 AND enterprise=$3`
			db.Exec(sqlStatement, value.Id, value.Name, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.product_option_values WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsLanguage(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("languages", enterpriseId) + "&display=[id,name,iso_code,active]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Language</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var languages PSLanguages
	json.Unmarshal(jsonPS, &languages)

	sqlStatement := `UPDATE public.ps_language SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(languages.Languages); i++ {
		language := languages.Languages[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_language WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, language.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_language(id, name, iso_code, active, enterprise) VALUES ($1, $2, $3, $4, $5)`
			db.Exec(sqlStatement, language.Id, language.Name, language.IsoCode, language.Active, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_language SET name=$2, iso_code=$3, active=$4, ps_exists=true WHERE id=$1 AND enterprise$5`
			db.Exec(sqlStatement, language.Id, language.Name, language.IsoCode, language.Active, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_language WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsCarriers(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("carriers", enterpriseId) + "&display=[id,deleted,name,active,url,max_width,max_height,max_depth,max_weight]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Zones</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var carriers PSCarriers
	json.Unmarshal(jsonPS, &carriers)

	sqlStatement := `UPDATE public.ps_carrier SET ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(carriers.Carriers); i++ {
		carrier := carriers.Carriers[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_carrier WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, carrier.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_carrier(id, deleted, name, active, url, max_width, max_height, max_depth, max_weight, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
			db.Exec(sqlStatement, carrier.Id, carrier.Deleted, carrier.Name, carrier.Active, carrier.Url, carrier.MaxWidth, carrier.MaxHeight, carrier.MaxDepth, carrier.MaxWeight, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_carrier SET deleted=$2, name=$3, active=$4, url=$5, max_width=$6, max_height=$7, max_depth=$8, max_weight=$9, ps_exists=true WHERE id=$1 AND enterprise=$10`
			db.Exec(sqlStatement, carrier.Id, carrier.Deleted, carrier.Name, carrier.Active, carrier.Url, carrier.MaxWidth, carrier.MaxHeight, carrier.MaxDepth, carrier.MaxWeight, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_carrier WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsOrders(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("orders", enterpriseId) + "&display=[id,reference,id_carrier,id_lang,id_customer,id_currency,id_address_delivery,id_address_invoice,module,total_discounts_tax_excl,total_shipping_tax_excl,date_add,date_upd]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var orders PSOrders
	json.Unmarshal(jsonPS, &orders)

	sqlStatement := `UPDATE public.ps_order SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(orders.Orders); i++ {
		order := orders.Orders[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_order WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, order.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		taxIncluded := order.TotalPaidTaxExcl != order.TotalPaidTaxIncl

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_order(id, reference, id_carrier, id_lang, id_customer, id_currency, id_address_delivery, id_address_invoice, module, total_discounts_tax_excl, total_shipping_tax_excl, date_add, date_upd, tax_included, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
			db.Exec(sqlStatement, order.Id, order.Reference, order.IdCarrier, order.IdLang, order.IdCustomer, order.IdCurrency, order.IdAddressDelivery, order.IdAddressInvoice, order.Module, order.TotalDiscountsTaxExcl, order.TotalShippingTaxExcl, order.DateAdd, order.DateUpd, taxIncluded, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_order SET reference=$2, id_carrier=$3, id_lang=$4, id_customer=$5, id_currency=$6, id_address_delivery=$7, id_address_invoice=$8, module=$9, total_discounts_tax_excl=$10, total_shipping_tax_excl=$11, date_add=$12, date_upd=$13, tax_included=$14, ps_exists=true WHERE id=$1 AND enterprise=$15`
			db.Exec(sqlStatement, order.Id, order.Reference, order.IdCarrier, order.IdLang, order.IdCustomer, order.IdCurrency, order.IdAddressDelivery, order.IdAddressInvoice, order.Module, order.TotalDiscountsTaxExcl, order.TotalShippingTaxExcl, order.DateAdd, order.DateUpd, taxIncluded, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_order WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

func importPsOrderDetails(enterpriseId int32) bool {
	url := getPrestaShopAPI_URL("order_details", enterpriseId) + "&display=[id,id_order,product_id,product_attribute_id,product_quantity,product_price]"
	jsonPS, err := getPrestaShopJSON(url)
	if err != nil {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Order details</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var details PSOrderDetails
	json.Unmarshal(jsonPS, &details)

	sqlStatement := `UPDATE public.ps_order_detail SET ps_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(details.OrderDetails); i++ {
		detail := details.OrderDetails[i]

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM ps_order_detail WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, detail.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.ps_order_detail(id, id_order, product_id, product_attribute_id, product_quantity, product_price, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			db.Exec(sqlStatement, detail.Id, detail.IdOrder, detail.ProductId, detail.ProductAttributeId, detail.ProductQuantity, detail.ProductPrice, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.ps_order_detail SET id_order=$2, product_id=$3, product_attribute_id=$4, product_quantity=$5, product_price=$6, ps_exists=true WHERE id=$1 AND enterprise=$7`
			db.Exec(sqlStatement, detail.Id, detail.IdOrder, detail.ProductId, detail.ProductAttributeId, detail.ProductQuantity, detail.ProductPrice, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.ps_order_detail WHERE ps_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	return true
}

// =====
// TRANSFER THE DATA TO THE ERP TABLES
// =====

func copyPsCurrencies(enterpriseId int32) bool {
	sqlStatement := `SELECT iso_code FROM public.ps_currency WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Currencies</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var isoCode string
		rows.Scan(&isoCode)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM currency WHERE iso_code=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			sqlStatement := `SELECT name,conversion_rate,symbol,numeric_iso_code FROM public.ps_currency WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Currencies</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var name string
			var conversionRate float64
			var symbol string
			var numericIsoCode int32
			row.Scan(&name, &conversionRate, &symbol, &numericIsoCode)

			c := Currency{}
			c.IsoCode = isoCode
			c.Name = name
			c.Change = conversionRate
			c.Sign = symbol
			c.IsoNum = int16(numericIsoCode)
			c.ExchangeDate = time.Now()
			c.enterprise = enterpriseId
			c.insertCurrency()
		}
	}
	return true
}

func copyPsCountries(enterpriseId int32) bool {
	sqlStatement := `SELECT iso_code FROM public.ps_country WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Countries</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var isoCode string
		rows.Scan(&isoCode)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM country WHERE iso_2=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			sqlStatement := `SELECT name,id_zone,id_currency FROM public.ps_country WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Countries</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var name string
			var id_zone int32
			var call_prefix int32
			row.Scan(&name, &id_zone, &call_prefix)

			sqlStatement = `SELECT zone FROM public.ps_zone WHERE id=$1 AND enterprise=$2 LIMIT 1`
			row = db.QueryRow(sqlStatement, id_zone, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Contries</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var zone string
			row.Scan(&zone)

			c := Country{}
			c.Iso2 = isoCode
			c.Name = name
			c.PhonePrefix = int16(call_prefix)
			c.Zone = zone
			c.enterprise = enterpriseId
			c.insertCountry()
		}
	}
	return true
}

func copyPsStates(enterpriseId int32) bool {
	sqlStatement := `SELECT iso_code FROM public.ps_state WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: States</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var isoCode string
		rows.Scan(&isoCode)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM state WHERE iso_code=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			sqlStatement := `SELECT name,iso_code,id_country FROM public.ps_state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, isoCode, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: States</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var name string
			var iso_code string
			var id_country int32
			row.Scan(&name, &iso_code, &id_country)

			sqlStatement = `SELECT iso_code FROM ps_country WHERE id=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, id_country, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: States</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var iso_country string
			row.Scan(&iso_country)

			sqlStatement = `SELECT id FROM country WHERE iso_2=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, iso_country, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				s := getSettingsRecordById(enterpriseId)
				if len(s.EmailSendErrorEcommerce) > 0 {
					sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: States</p><p>Error data: "+row.Err().Error()+"</p>", enterpriseId)
				}
				return false
			}

			var country int32
			row.Scan(&country)

			s := State{}
			s.Country = country
			s.Name = name
			s.IsoCode = iso_code
			s.enterprise = enterpriseId
			s.insertState()
		}
	}
	return true
}

func copyPsCustomers(enterpriseId int32) bool {
	sqlStatement := `SELECT id FROM public.ps_customer WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var id int32
		rows.Scan(&id)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM customer WHERE ps_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, id, enterpriseId)
		if row.Err() != nil {
			log("PrestaShop", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			// get the customer data
			sqlStatement := `SELECT id_lang,company,firstname,lastname,email,date_add FROM public.ps_customer WHERE id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var id_lang int32
			var company string
			var firstname string
			var lastname string
			var email string
			var date_add time.Time
			row.Scan(&id_lang, &company, &firstname, &lastname, &email, &date_add)

			// get the customer language
			var lang int32
			if id_lang != 0 {
				sqlStatement := `SELECT iso_code FROM ps_language WHERE id=$1 AND enterprise=$2`
				row := db.QueryRow(sqlStatement, id_lang, enterpriseId)
				if row.Err() != nil {
					log("PrestaShop", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var iso_code string
				row.Scan(&iso_code)

				sqlStatement = `SELECT id FROM language WHERE iso_2=$1 AND enterprise=$2`
				row = db.QueryRow(sqlStatement, strings.ToUpper(iso_code), enterpriseId)
				if row.Err() != nil {
					log("PrestaShop", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				row.Scan(&lang)
			}

			// get the fiscal data from the address
			var taxId string
			var vatNumber string

			sqlStatement = `SELECT dni,vat_number FROM ps_address WHERE id_customer=$1 AND enterprise=$2 AND dni != '' ORDER BY id DESC LIMIT 1`
			row = db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&taxId, &vatNumber)

			c := Customer{}
			c.FiscalName = firstname + " " + lastname
			if len(company) == 0 {
				c.Tradename = c.FiscalName
			} else {
				c.Tradename = company
			}
			if len(c.FiscalName) == 0 {
				c.FiscalName = company
			}
			c.Name = c.Tradename + " / " + c.FiscalName
			c.Email = email
			c.DateCreated = date_add
			c.prestaShopId = id
			if lang != 0 {
				c.Language = &lang
			}
			c.TaxId = taxId
			c.VatNumber = vatNumber
			c.enterprise = enterpriseId
			c.insertCustomer(0)
		}
	}

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Customers</p>"+errorHtml, enterpriseId)
		}
	}

	return true
}

func copyPsAddresses(enterpriseId int32) bool {
	sqlStatement := `SELECT id FROM public.ps_address WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Addresses</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var id int32
		rows.Scan(&id)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM address WHERE ps_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			// get address data
			sqlStatement := `SELECT id_country,id_state,id_customer,address1,address2,postcode,city,other,phone,vat_number,dni FROM public.ps_address WHERE id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var id_country int32
			var id_state int32
			var id_customer int32
			var address1 string
			var address2 string
			var postcode string
			var city string
			var other string
			var phone string
			var vat_number string
			var dni string
			row.Scan(&id_country, &id_state, &id_customer, &address1, &address2, &postcode, &city, &other, &phone, &vat_number, &dni)

			// get customer
			sqlStatement = `SELECT id FROM customer WHERE ps_id=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, id_customer, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var customer int32
			row.Scan(&customer)

			// get country
			sqlStatement = `SELECT iso_code FROM ps_country WHERE id=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, id_country, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var iso_code string
			row.Scan(&iso_code)

			sqlStatement = `SELECT id FROM country WHERE iso_2=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, iso_code, enterpriseId)
			if row.Err() != nil {
				log("PrestaShop", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var country int32
			row.Scan(&country)

			// get state
			var state *int32
			if id_state != 0 {
				sqlStatement := `SELECT iso_code FROM ps_state WHERE id=$1 AND enterprise=$2`
				row = db.QueryRow(sqlStatement, id_state, enterpriseId)
				if row.Err() != nil {
					log("PrestaShop", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var iso_code string
				row.Scan(&iso_code)

				sqlStatement = `SELECT id FROM state WHERE iso_code=$1 AND enterprise=$2`
				row = db.QueryRow(sqlStatement, iso_code, enterpriseId)
				if row.Err() != nil {
					log("PrestaShop", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				row.Scan(&state)
			}

			a := Address{}
			a.Customer = &customer
			a.Country = country
			a.State = state
			a.City = city
			a.ZipCode = postcode
			a.Address = address1
			a.Address2 = address2
			a.Notes = other
			a.prestaShopId = id
			a.PrivateOrBusiness = "_"
			a.enterprise = enterpriseId
			a.insertAddress(0)

			// set the customer details if are empty
			c := getCustomerRow(customer)
			if c.TaxId == "" {
				c.TaxId = dni
			}
			if c.VatNumber == "" {
				c.VatNumber = vat_number
			}
			if c.Phone == "" {
				c.Phone = phone
			}
			if c.Country == nil {
				c.Country = &country
			}
			if c.State == nil {
				c.State = state
			}
			c.updateCustomer(0)
		}
	}

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Customers</p>"+errorHtml, enterpriseId)
		}
	}

	return true
}

func copyPsLanguages(enterpriseId int32) bool {
	sqlStatement := `SELECT iso_code,name FROM public.ps_language WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Languages</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var isoCode string
		var name string
		rows.Scan(&isoCode, &name)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM language WHERE iso_2=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, strings.ToUpper(isoCode), enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			l := Language{}
			l.Name = name
			l.Iso2 = strings.ToUpper(isoCode)
			l.enterprise = enterpriseId
			l.insertLanguage()
		}
	}
	return true
}

func copyPsCarriers(enterpriseId int32) bool {
	sqlStatement := `SELECT id,name,url,max_width,max_height,max_depth,max_weight FROM public.ps_carrier WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Carriers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		var name string
		var url string
		var max_width float64
		var max_height float64
		var max_depth float64
		var max_weight float64
		rows.Scan(&id, &name, &url, &max_width, &max_height, &max_depth, &max_weight)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM carrier WHERE ps_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			c := Carrier{}
			c.Name = name
			c.Web = url
			c.MaxWidth = max_width
			c.MaxHeight = max_height
			c.MaxDepth = max_depth
			c.MaxWeight = max_weight
			c.PrestaShopId = id
			c.enterprise = enterpriseId
			c.insertCarrier()
		}
	}
	return true
}

func copyPsProducts(enterpriseId int32) bool {
	sqlStatement := `SELECT id,name,ean13,reference,price,date_add,description FROM public.ps_product WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Prodocuts</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var ps_productId int32
		var name string
		var ean13 string
		var reference string
		var price float64
		var dateAdd time.Time
		var description string
		rows.Scan(&ps_productId, &name, &ean13, &reference, &price, &dateAdd, &description)
		description = strip.StripTags(description)

		sqlStatement := `SELECT COUNT(id) FROM ps_product_combination WHERE id_product=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, ps_productId, enterpriseId)

		var combinations int32
		row.Scan(&combinations)

		if combinations == 0 { // it's a simple product

			// does the product exist?
			sqlStatement := `SELECT id FROM product WHERE ps_id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, ps_productId, enterpriseId)

			var productId int32
			row.Scan(&productId)

			if productId <= 0 {
				p := Product{}
				p.Name = name
				p.BarCode = ean13
				p.Reference = reference
				p.Price = price
				p.DateCreated = dateAdd
				p.Description = description
				p.prestaShopId = ps_productId
				p.enterprise = enterpriseId
				result := p.insertProduct(0)
				if !result.Ok {
					errors = append(errors, "Error inserting a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}
			} else {
				p := getProductRow(productId)
				p.Name = name
				p.BarCode = ean13
				p.Reference = reference
				p.Price = price
				p.DateCreated = dateAdd
				p.Description = description
				result := p.updateProduct(0)
				if !result.Ok {
					errors = append(errors, "Error updating a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}
			}

		} else { // it's a product with combinations
			sqlStatement := `SELECT id,reference,ean13,product_option_values,price FROM ps_product_combination WHERE id_product=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, ps_productId, enterpriseId)
			if err != nil {
				log("PrestaShop", err.Error())
				return false
			}

			for rows.Next() {
				var combinationId int32
				var combinationReference string
				var combinationEan13 string
				var productOptionValues []int32
				var combinationPrice float64
				rows.Scan(&combinationId, &combinationReference, &combinationEan13, pq.Array(&productOptionValues), &combinationPrice)

				// does the product exist?
				sqlStatement := `SELECT id FROM product WHERE ps_id=$1 AND ps_combination_id=$2 AND enterprise=$3`
				row := db.QueryRow(sqlStatement, ps_productId, combinationId, enterpriseId)

				var productId int32
				row.Scan(&productId)

				// generate the product name
				combinationName := name
				for i := 0; i < len(productOptionValues); i++ {
					sqlStatement := `SELECT name FROM ps_product_option_values WHERE id=$1 AND enterprise=$2`
					row := db.QueryRow(sqlStatement, productOptionValues[i], enterpriseId)

					var name string
					row.Scan(&name)
					combinationName += " " + name
				}

				if combinationPrice == 0 {
					combinationPrice = price
				}

				if productId <= 0 {
					p := Product{}
					p.Name = combinationName
					p.BarCode = combinationEan13
					p.Reference = combinationReference
					p.Price = combinationPrice
					p.DateCreated = dateAdd
					p.Description = description
					p.prestaShopId = ps_productId
					p.prestaShopCombinationId = combinationId
					p.enterprise = enterpriseId
					result := p.insertProduct(0)
					if !result.Ok {
						errors = append(errors, "Error inserting a product with combinations into MARKETNET. Product name "+
							p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErorCode))+" extra data: "+stringArrayToString(result.ExtraData))
					}
				} else {
					p := getProductRow(productId)
					p.Name = combinationName
					p.BarCode = combinationEan13
					p.Reference = combinationReference
					p.Price = combinationPrice
					p.DateCreated = dateAdd
					p.Description = description
					result := p.updateProduct(0)
					if !result.Ok {
						errors = append(errors, "Error updating a product with combinations into MARKETNET. Product name "+
							p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErorCode))+" extra data: "+stringArrayToString(result.ExtraData))
					}
				}
			}

			rows.Close()
		}

	}

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Products</p>"+errorHtml, enterpriseId)
		}
	}

	return true
}

func copyPsOrders(enterpriseId int32) bool {
	settings := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id,reference,id_carrier,id_lang,id_customer,id_currency,id_address_delivery,id_address_invoice,module,total_discounts_tax_excl,total_shipping_tax_excl,tax_included FROM public.ps_order WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var orderId int32
		var reference string
		var idCarrier int32
		var idLang int32
		var idCustomer int32
		var idCurrency int32
		var idAddressDelivery int32
		var idAddressInvoice int32
		var module string
		var totalDiscountsTaxExcl float64
		var totalShippingTaxExcl float64
		var taxIncluded bool
		rows.Scan(&orderId, &reference, &idCarrier, &idLang, &idCustomer, &idCurrency, &idAddressDelivery, &idAddressInvoice, &module, &totalDiscountsTaxExcl, &totalShippingTaxExcl, &taxIncluded)

		// does the order exist?
		sqlStatement := `SELECT COUNT(id) FROM sales_order WHERE ps_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, orderId, enterpriseId)

		var orders int32
		row.Scan(&orders)

		if orders > 0 { // don't continue if the order exists
			continue
		}

		// get the carrier
		sqlStatement = `SELECT id FROM carrier WHERE ps_id=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, idCarrier, enterpriseId)

		var carrier int32
		row.Scan(&carrier)

		if carrier == 0 { // don't continue if the carrier doesn't exists
			errors = append(errors, "Can't import order: The carrier does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the language
		sqlStatement = `SELECT iso_code FROM ps_language WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, idLang, enterpriseId)

		var lang_iso_code string
		row.Scan(&lang_iso_code)

		sqlStatement = `SELECT id FROM language WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, strings.ToUpper(lang_iso_code), enterpriseId)

		var language int32
		row.Scan(&language)

		if language == 0 { // don't continue if the language doesn't exists
			errors = append(errors, "Can't import order. The language does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the customer
		sqlStatement = `SELECT id FROM customer WHERE ps_id=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, idCustomer, enterpriseId)

		var customer int32
		row.Scan(&customer)

		if customer == 0 { // don't continue if the customer doesn't exists
			errors = append(errors, "Can't import order. The customer does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the payment method
		sqlStatement = `SELECT id,paid_in_advance FROM payment_method WHERE prestashop_module_name=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, module, enterpriseId)

		var paymentMethod int32
		var paidInAdvance bool
		row.Scan(&paymentMethod, &paidInAdvance)

		if paymentMethod == 0 { // don't continue if the payment method doesn't exists
			errors = append(errors, "Can't import order. The payment method does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the currency
		sqlStatement = `SELECT iso_code FROM ps_currency WHERE id=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, idCurrency, enterpriseId)

		var currency_iso_code string
		row.Scan(&currency_iso_code)

		sqlStatement = `SELECT id FROM currency WHERE iso_code=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, currency_iso_code, enterpriseId)

		var currency int32
		row.Scan(&currency)

		if currency == 0 { // don't continue if the currency doesn't exists
			errors = append(errors, "Can't import order. The currency does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the billing address
		sqlStatement = `SELECT id,(SELECT zone FROM country WHERE country.id=address.country) FROM address WHERE ps_id=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, idAddressInvoice, enterpriseId)

		var billingAddress int32
		var billingZone string
		row.Scan(&billingAddress, &billingZone)

		if billingAddress == 0 { // don't continue if the billing address doesn't exists
			errors = append(errors, "Can't import order. The billing address does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the shipping address
		sqlStatement = `SELECT id FROM address WHERE ps_id=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, idAddressDelivery, enterpriseId)

		var shippingAddress int32
		row.Scan(&shippingAddress)

		if shippingAddress == 0 { // don't continue if the shipping address doesn't exists
			errors = append(errors, "Can't import order. The shipping address does not exist in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
			continue
		}

		s := SaleOrder{}
		s.Warehouse = settings.DefaultWarehouse
		s.Reference = reference
		s.Customer = customer
		s.PaymentMethod = paymentMethod
		s.Currency = currency
		s.BillingAddress = billingAddress
		s.ShippingAddress = shippingAddress
		s.prestaShopId = orderId

		if billingZone == "E" {
			s.BillingSeries = *settings.PrestaShopExportSerie
		} else if billingZone == "U" && !taxIncluded {
			s.BillingSeries = *settings.PrestaShopIntracommunitySerie
		} else {
			s.BillingSeries = *settings.PrestaShopInteriorSerie
		}

		s.enterprise = enterpriseId
		if ok, _ := s.insertSalesOrder(0); !ok {
			errors = append(errors, "Can't import order. Error creating the order in MARKETNET. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
		}

		// set the customer details if are empty
		c := getCustomerRow(customer)
		if c.PaymentMethod == nil {
			c.PaymentMethod = &paymentMethod
		}
		if c.BillingSeries == nil || *c.BillingSeries == "" {
			c.BillingSeries = &s.BillingSeries
		}
		if !c.updateCustomer(0) {
			errors = append(errors, "Can't update the customer data. Order reference "+reference+" order id "+strconv.Itoa(int(orderId)))
		}
	}

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Orders</p>"+errorHtml, enterpriseId)
		}
	}
	return true
}

func copyPsOrderDetails(enterpriseId int32) bool {
	sqlStatement := `SELECT id,id_order,product_id,product_attribute_id,product_quantity,product_price,(SELECT tax_included FROM ps_order WHERE ps_order.id=ps_order_detail.id_order),(SELECT vat_percent FROM product WHERE product.ps_id=ps_order_detail.product_id AND product.ps_combination_id=ps_order_detail.product_attribute_id) FROM public.ps_order_detail WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("PrestaShop", err.Error())
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	orderIds := make([]int64, 0)

	for rows.Next() {
		var detailId int32
		var orderId int32
		var productId int32
		var productAttributeId int32
		var ProductQuantity int32
		var productPrice float64
		var taxIncluded bool
		var vatPercent float64
		rows.Scan(&detailId, &orderId, &productId, &productAttributeId, &ProductQuantity, &productPrice, &taxIncluded, &vatPercent)

		sqlStatement := `SELECT COUNT(id) FROM sales_order_detail WHERE ps_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, detailId, enterpriseId)

		var details int32
		row.Scan(&details)

		if details > 0 { // the detail already exists
			continue
		}

		// get the sale order
		sqlStatement = `SELECT id FROM sales_order WHERE ps_id=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderId, enterpriseId)

		var order int64
		row.Scan(&order)

		if order <= 0 {
			errors = append(errors, "Can't import order detail. The order does not exist in MARKETNET. Order id "+strconv.Itoa(int(orderId)))
			continue
		}

		// get the product
		sqlStatement = `SELECT id FROM product WHERE ps_id=$1 AND ps_combination_id=$2 AND enterprise=$3`
		row = db.QueryRow(sqlStatement, productId, productAttributeId, enterpriseId)

		var product int32
		row.Scan(&product)

		if product <= 0 {
			errors = append(errors, "Can't import order detail. The product does not exist in MARKETNET. Order id "+strconv.Itoa(int(orderId)))
			continue
		}

		d := SalesOrderDetail{}
		d.Order = order
		d.Product = product
		d.Quantity = ProductQuantity
		d.Price = productPrice

		if !taxIncluded {
			d.VatPercent = 0
		} else {
			d.VatPercent = vatPercent
		}

		d.enterprise = enterpriseId
		d.prestaShopId = detailId
		ok := d.insertSalesOrderDetail(0).Ok

		if ok {
			found := false
			for i := 0; i < len(orderIds); i++ {
				if orderIds[i] == order {
					found = true
					break
				}
			}
			if !found {
				orderIds = append(orderIds, order)
			}
		}
	}

	// if the payment method is paid in advance, it means that this order is already paid (by VISA o PayPal etc)
	// automatically generate an invoice for this payment

	for i := 0; i < len(orderIds); i++ {
		sqlStatement = `SELECT paid_in_advance FROM payment_method WHERE id=(SELECT payment_method FROM sales_order WHERE id=$1) AND enterprise=$2`
		row := db.QueryRow(sqlStatement, orderIds[i], enterpriseId)

		var paidInAdvance bool
		row.Scan(&paidInAdvance)

		if paidInAdvance {
			invoiceAllSaleOrder(orderIds[i], enterpriseId, 0)
		}
	}

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.EmailSendErrorEcommerce, s.EmailSendErrorEcommerce, "PrestaShop import error", "<p>Error at: Customers</p>"+errorHtml, enterpriseId)
		}
	}

	return true
}

// =====
// CONFIGURE FROM THE ERP
// =====

type PSZoneWeb struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Zone       string `json:"zone"`
	enterprise int32
}

func getPSZones(enterpriseId int32) []PSZoneWeb {
	var zones []PSZoneWeb = make([]PSZoneWeb, 0)
	sqlStatement := `SELECT id,name,zone FROM public.ps_zone WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		return zones
	}
	defer rows.Close()

	for rows.Next() {
		z := PSZoneWeb{}
		rows.Scan(&z.Id, &z.Name, &z.Zone, &z.enterprise)
		zones = append(zones, z)
	}

	return zones
}

func (z *PSZoneWeb) updatePSZoneWeb() bool {
	if z.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.ps_zone SET zone=$2 WHERE id=$1 AND enterprise=$3`
	res, err := db.Exec(sqlStatement, z.Id, z.Zone, z.enterprise)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

//
// SET TRACKING NUMBER
//

func updateTrackingNumberPrestaShopOrder(salesOrderId int64, trackingNumber string, enterpriseId int32) bool {
	settings := getSettingsRecordById(enterpriseId)
	if settings.Ecommerce != "P" {
		return false
	}

	s := getSalesOrderRow(salesOrderId)
	if s.Id <= 0 || s.prestaShopId <= 0 {
		return false
	}

	url := settings.PrestaShopUrl + "orders/" + strconv.Itoa(int(s.prestaShopId)) + "/?ws_key=" + settings.PrestaShopApiKey

	xmlPs, err := getPrestaShopJSON(url)
	if err != nil {
		return false
	}

	index := strings.Index(string(xmlPs), "<shipping_number notFilterable=\"true\">")
	if index <= 0 {
		return false
	}
	index += len("<shipping_number notFilterable=\"true\">")
	indexEnd := strings.Index(string(xmlPs), "</shipping_number>")
	if indexEnd <= 0 {
		return false
	}

	xml := make([]byte, len(xmlPs))

	for i := 0; i < len(xmlPs); i++ {
		xml[i] = xmlPs[i]
	}

	xml = append(xml[:index], "<![CDATA["+trackingNumber+"]]>"...)
	xml = append(xml, xmlPs[indexEnd:]...)

	xmlSend := setStatusXmlOrderPrestaShop(xml, strconv.Itoa(int(settings.PrestashopStatusShipped)))
	if xmlSend == nil {
		return false
	}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(xmlSend))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	_, err = client.Do(req)

	return err == nil
}

func setStatusXmlOrderPrestaShop(xmlPs []byte, status string) []byte {
	index := strings.Index(string(xmlPs), "<current_state")
	if index <= 0 {
		return nil
	}
	index += len("<current_state")
	indexEnd := strings.Index(string(xmlPs), "</current_state>")
	if indexEnd <= 0 {
		return nil
	}

	xml := make([]byte, len(xmlPs))

	for i := 0; i < len(xmlPs); i++ {
		xml[i] = xmlPs[i]
	}

	xml = append(xml[:index], "><![CDATA["+status+"]]>"...)
	xml = append(xml, xmlPs[indexEnd:]...)

	return xml
}

func updateStatusPaymentAcceptedPrestaShop(orderId int64, enterpriseId int32) bool {
	settings := getSettingsRecordById(enterpriseId)
	if settings.Ecommerce != "P" {
		return false
	}

	s := getSalesOrderRow(orderId)
	if s.prestaShopId <= 0 {
		return true
	}

	sqlStatement := `SELECT paid_in_advance FROM payment_method WHERE id=(SELECT payment_method FROM sales_order WHERE id=$1) AND enterprise=$2`
	row := db.QueryRow(sqlStatement, orderId, enterpriseId)

	var paidInAdvance bool
	row.Scan(&paidInAdvance)

	if !paidInAdvance { // this is not an automatically generated invoice, someone accepted the payment, notify PrestaShop
		url := settings.PrestaShopUrl + "orders/" + strconv.Itoa(int(s.prestaShopId)) + "/?ws_key=" + settings.PrestaShopApiKey

		xmlPs, err := getPrestaShopJSON(url)
		if err != nil {
			return false
		}

		xml := setStatusXmlOrderPrestaShop(xmlPs, strconv.Itoa(int(settings.PrestashopStatusPaymentAccepted)))
		if xml == nil {
			return false
		}

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(xml))
		req.Header.Set("Content-Type", "text/xml; charset=utf-8")
		_, err = client.Do(req)

		return err == nil
	}

	return true
}
