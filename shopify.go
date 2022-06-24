/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	strip "github.com/grokify/html-strip-tags-go"
)

const SHOPIFY_TRANSACTION_ACCEPT_PAYMENT_KIND = "sale"

// =====
// GENERIC FUNCTIONS
// =====

func getShopifyAPI_URL(resourceName string, enterpriseId int32) string {
	s := getSettingsRecordById(enterpriseId)

	return s.SettingsEcommerce.ShopifyUrl + resourceName + ".json"
}

func getShopifyJSON(URL string, enterpriseId int32) ([]byte, error) {
	s := getSettingsRecordById(enterpriseId)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Set("X-Shopify-Access-Token", s.SettingsEcommerce.ShopifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func postShopifyJSON(URL string, data []byte, enterpriseId int32) ([]byte, error) {
	s := getSettingsRecordById(enterpriseId)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))
	req.Header.Set("X-Shopify-Access-Token", s.SettingsEcommerce.ShopifyToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func putShopifyJSON(URL string, data []byte, enterpriseId int32) ([]byte, error) {
	s := getSettingsRecordById(enterpriseId)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(data))
	req.Header.Set("X-Shopify-Access-Token", s.SettingsEcommerce.ShopifyToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
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
// DEFINE SHOPIFY CLASSES
// =====

type SYCustomers struct {
	Customers []SYCustomer `json:"customers"`
}

type SYCustomer struct {
	Id             int64       `json:"id" gorm:"primaryKey"`
	Email          string      `json:"email" gorm:"column:email;type:character varying(100);not null:true"`
	FirstName      string      `json:"first_name" gorm:"column:first_name;type:character varying(100);not null:true"`
	LastName       string      `json:"last_name" gorm:"column:last_name;type:character varying(100);not null:true"`
	TaxExempt      bool        `json:"tax_exempt" gorm:"column:tax_exempt;type:boolean;not null:true"`
	Phone          string      `json:"phone" gorm:"column:phone;type:character varying(25);not null:true"`
	Currency       string      `json:"currency" gorm:"column:currency;type:character varying(5);not null:true"`
	Addresses      []SYAddress `json:"addresses" gorm:"-"`
	SyExists       bool        `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	DefaultAddress SYAddress   `json:"default_address" gorm:"column:default_address_id;type:integer;not null:true"`
	EnterpriseId   int32       `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise     Settings    `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYCustomer) TableName() string {
	return "sy_customers"
}

type SYAddress struct {
	Id           int64    `json:"id" gorm:"primaryKey"`
	CustomerId   int64    `json:"customer_id" gorm:"column:customer_id;type:bigint;not null:true"`
	Company      string   `json:"company" gorm:"column:company;type:character varying(100);not null:true"`
	Address1     string   `json:"address1" gorm:"column:address1;type:character varying(100);not null:true"`
	Address2     string   `json:"address2" gorm:"column:address2;type:character varying(100);not null:true"`
	City         string   `json:"city" gorm:"column:city;type:character varying(50);not null:true"`
	Province     string   `json:"province" gorm:"column:province;type:character varying(50);not null:true"`
	Zip          string   `json:"zip" gorm:"column:zip;type:character varying(25);not null:true"`
	CountryCode  string   `json:"country_code" gorm:"column:country_code;type:character varying(5);not null:true"`
	SyExists     bool     `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYAddress) TableName() string {
	return "sy_addresses"
}

type SYProducts struct {
	Products []SYProduct `json:"products"`
}

type SYProduct struct {
	Id           int64       `json:"id" gorm:"primaryKey"`
	Title        string      `json:"title" gorm:"column:title;type:character varying(150);not null:true"`
	BodyHtml     string      `json:"body_html" gorm:"column:body_html;type:text;not null:true"`
	Variants     []SYVariant `json:"variants" gorm:"-"`
	SyExists     bool        `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	EnterpriseId int32       `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings    `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYProduct) TableName() string {
	return "sy_products"
}

type SYVariant struct {
	Id           int64    `json:"id" gorm:"primaryKey"`
	ProductId    int64    `json:"product_id" gorm:"column:product_id;type:bigint;not null:true"`
	Title        string   `json:"title" gorm:"column:title;type:character varying(150);not null:true"`
	Price        string   `json:"price" gorm:"column:price;type:numeric(12,6);not null:true"`
	Sku          string   `json:"sku" gorm:"column:sku;type:character varying(25);not null:true"`
	Option1      string   `json:"option1" gorm:"column:option1;type:character varying(150);not null:true"`
	Option2      *string  `json:"option2" gorm:"column:option2;type:character varying(150)"`
	Option3      *string  `json:"option3" gorm:"column:option3;type:character varying(150)"`
	Taxable      bool     `json:"taxable" gorm:"column:taxable;type:boolean;not null:true"`
	Barcode      string   `json:"barcode" gorm:"column:barcode;type:character varying(25);not null:true"`
	Grams        int32    `json:"grams" gorm:"column:grams;type:integer;not null:true"`
	SyExists     bool     `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYVariant) TableName() string {
	return "sy_variants"
}

type SYVariantDB struct {
	ProductId int64   `json:"product_id"`
	Id        int64   `json:"id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Sku       string  `json:"sku"`
	Option1   string  `json:"option1"`
	Option2   *string `json:"option2"`
	Option3   *string `json:"option3"`
	Taxable   bool    `json:"taxable"`
	Barcode   string  `json:"barcode"`
	Grams     int32   `json:"grams"`
}

type SYDraftOrders struct {
	DraftOrders []SYDraftOrder `json:"draft_orders"`
}

type SYDraftOrder struct {
	Id                         int64                  `json:"id" gorm:"primaryKey"`
	Currency                   string                 `json:"currency" gorm:"column:currency;type:character varying(5);not null:true"`
	TaxExempt                  bool                   `json:"tax_exempt" gorm:"column:tax_exempt;type:boolean;not null:true"`
	Name                       string                 `json:"name" gorm:"column:name;type:character varying(9);not null:true"`
	LineItems                  []SYDraftOrderLineItem `json:"line_items" gorm:"-"`
	ShippingAddress            SYAddress              `json:"shipping_address" gorm:"-"`
	BillingAddress             SYAddress              `json:"billing_address" gorm:"-"`
	ShippingAddress1           string                 `json:"-" gorm:"column:shipping_address_1;type:character varying(100);not null:true"`
	ShippingAddress2           string                 `json:"-" gorm:"column:shipping_address2;type:character varying(100);not null:true"`
	ShippingAddressCity        string                 `json:"-" gorm:"column:shipping_address_city;type:character varying(50);not null:true"`
	ShippingAddressZip         string                 `json:"-" gorm:"column:shipping_address_zip;type:character varying(25);not null:true"`
	ShippingAddressCountryCode string                 `json:"-" gorm:"column:shipping_address_country_code;type:character varying(5);not null:true"`
	BillingAddress1            string                 `json:"-" gorm:"column:billing_address_1;type:character varying(100);not null:true"`
	BillingAddress2            string                 `json:"-" gorm:"column:billing_address2;type:character varying(100);not null:true"`
	BillingAddressCity         string                 `json:"-" gorm:"column:billing_address_city;type:character varying(50);not null:true"`
	BillingAddressZip          string                 `json:"-" gorm:"column:billing_address_zip;type:character varying(25);not null:true"`
	BillingAddressCountryCode  string                 `json:"-" gorm:"column:billing_address_country_code;type:character varying(5);not null:true"`
	TotalTax                   string                 `json:"total_tax" gorm:"column:total_tax;type:numeric(14,6);not null:true"`
	Customer                   SYCustomer             `json:"customer" gorm:"-"`
	CustomerId                 int64                  `json:"-" gorm:"column:customer_id;type:bigint;not null:true"`
	SyExists                   bool                   `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	OrderId                    *int64                 `json:"order_id" gorm:"column:order_id;type:bigint;index:_sy_draft_orders_order_id,where:order_id IS NOT NULL"`
	EnterpriseId               int32                  `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise                 Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYDraftOrder) TableName() string {
	return "sy_draft_orders"
}

type SYDraftOrderLineItem struct {
	Id           int64    `json:"id" gorm:"primaryKey"`
	VariantId    int64    `json:"variant_id" gorm:"column:variant_id;type:bigint;not null:true"`
	ProductId    int64    `json:"product_id" gorm:"column:product_id;type:bigint;not null:true"`
	Quantity     int32    `json:"quantity" gorm:"column:quantity;type:integer;not null:true"`
	Taxable      bool     `json:"taxable" gorm:"column:taxable;type:boolean;not null:true"`
	Price        string   `json:"price" gorm:"column:price;type:numeric(12,6);not null:true"`
	DraftOrderId int64    `json:"-" gorm:"column:draft_order_id;type:bigint;not null:true"`
	SyExists     bool     `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYDraftOrderLineItem) TableName() string {
	return "sy_draft_order_line_item"
}

type SYOrders struct {
	Orders []SYOrder `json:"orders"`
}

type SYOrder struct {
	Id                            int64                  `json:"id" gorm:"primaryKey"`
	Currency                      string                 `json:"currency" gorm:"column:currency;type:character varying(5);not null:true"`
	CurrentTotalDiscounts         string                 `json:"current_total_discounts" gorm:"column:current_total_discounts;type:numeric(14,6);not null:true"`
	TotalShippingPriceSet         TotalShippingPriceSet  `json:"total_shipping_price_set" gorm:"-"`
	TotalShippingPriceSetAmount   float64                `json:"-" gorm:"column:total_shipping_price_set_amount;type:numeric(14,6);not null:true"`
	TotalShippingPriceSetCurrency string                 `json:"-" gorm:"column:total_shipping_price_set_currency_code;type:character varying(5);not null:true"`
	TaxExempt                     bool                   `json:"tax_exempt" gorm:"column:tax_exempt;type:boolean;not null:true"`
	Name                          string                 `json:"name" gorm:"column:name;type:character varying(9);not null:true"`
	LineItems                     []SYDraftOrderLineItem `json:"line_items" gorm:"-"`
	ShippingAddress               SYAddress              `json:"shipping_address" gorm:"-"`
	BillingAddress                SYAddress              `json:"billing_address" gorm:"-"`
	ShippingAddress1              string                 `json:"-" gorm:"column:shipping_address_1;type:character varying(100);not null:true"`
	ShippingAddress2              string                 `json:"-" gorm:"column:shipping_address2;type:character varying(100);not null:true"`
	ShippingAddressCity           string                 `json:"-" gorm:"column:shipping_address_city;type:character varying(50);not null:true"`
	ShippingAddressZip            string                 `json:"-" gorm:"column:shipping_address_zip;type:character varying(25);not null:true"`
	ShippingAddressCountryCode    string                 `json:"-" gorm:"column:shipping_address_country_code;type:character varying(5);not null:true"`
	BillingAddress1               string                 `json:"-" gorm:"column:billing_address_1;type:character varying(100);not null:true"`
	BillingAddress2               string                 `json:"-" gorm:"column:billing_address2;type:character varying(100);not null:true"`
	BillingAddressCity            string                 `json:"-" gorm:"column:billing_address_city;type:character varying(50);not null:true"`
	BillingAddressZip             string                 `json:"-" gorm:"column:billing_address_zip;type:character varying(25);not null:true"`
	BillingAddressCountryCode     string                 `json:"-" gorm:"column:billing_address_country_code;type:character varying(5);not null:true"`
	TotalTax                      string                 `json:"total_tax" gorm:"column:total_tax;type:numeric(14,6);not null:true"`
	Customer                      SYCustomer             `json:"customer" gorm:"-"`
	CustomerId                    int64                  `json:"-" gorm:"column:customer_id;type:bigint;not null:true"`
	SyExists                      bool                   `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	Gateway                       string                 `json:"gateway" gorm:"column:gateway;type:character varying(50);not null:true"`
	EnterpriseId                  int32                  `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise                    Settings               `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYOrder) TableName() string {
	return "sy_orders"
}

type TotalShippingPriceSet struct {
	ShopMoney ShopMoney `json:"shop_money"`
}

type ShopMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type SYOrderLineItem struct {
	Id           int64    `json:"id" gorm:"primaryKey"`
	VariantId    int64    `json:"variant_id" gorm:"column:variant_id;type:bigint;not null:true"`
	ProductId    int64    `json:"product_id" gorm:"column:product_id;type:bigint;not null:true"`
	Quantity     int32    `json:"quantity" gorm:"column:quantity;type:integer;not null:true"`
	Taxable      bool     `json:"taxable" gorm:"column:taxable;type:boolean;not null:true"`
	Price        string   `json:"price" gorm:"column:price;type:numeric(12,6);not null:true"`
	OrderId      int64    `json:"-" gorm:"column:order_id;type:bigint;not null:true"`
	SyExists     bool     `json:"-" gorm:"column:sy_exists;type:boolean;not null"`
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (SYOrderLineItem) TableName() string {
	return "sy_order_line_item"
}

// main import function
func importFromShopify(enterpriseId int32) {
	s := getSettingsRecordById(enterpriseId)
	if s.SettingsEcommerce.Ecommerce != "S" {
		return
	}

	// get all data from Shopify, write it in tables like the ones that Shopify uses
	if !importSyCustomers(enterpriseId) {
		return
	}
	if !importSyProducts(enterpriseId) {
		return
	}
	if !importSyDraftOrders(enterpriseId) {
		return
	}
	if !importSyOrders(enterpriseId) {
		return
	}

	// trasnfer the data form the Shopify tables to the ERP
	if !copySyCustomers(enterpriseId) {
		return
	}
	if !copySyProducts(enterpriseId) {
		return
	}
	if !copySyDraftOrders(enterpriseId) {
		return
	}
	if !copySyOrders(enterpriseId) {
		return
	}
}

// =====
// COPY THE DATA FROM WOOCOMMERCE TO THE WC MARKETNET TABLES
// =====

func importSyCustomers(enterpriseId int32) bool {
	url := getShopifyAPI_URL("customers", enterpriseId)
	jsonSY, err := getShopifyJSON(url, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var customers SYCustomers
	json.Unmarshal(jsonSY, &customers)

	sqlStatement := `UPDATE public.sy_customers SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `UPDATE public.sy_addresses SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(customers.Customers); i++ {
		customer := customers.Customers[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_customers WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, customer.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_customers(id, email, first_name, last_name, tax_exempt, phone, currency, default_address_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			db.Exec(sqlStatement, customer.Id, customer.Email, customer.FirstName, customer.LastName, customer.TaxExempt, customer.Phone, customer.Currency, customer.DefaultAddress.Id, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_customers SET email=$2, first_name=$3, last_name=$4, tax_exempt=$5, phone=$6, currency=$7, sy_exists=true, default_address_id=$8 WHERE id=$1 AND enterprise=$9`
			db.Exec(sqlStatement, customer.Id, customer.Email, customer.FirstName, customer.LastName, customer.TaxExempt, customer.Phone, customer.Currency, customer.DefaultAddress.Id, enterpriseId)
		}

		for j := 0; j < len(customer.Addresses); j++ {
			address := customer.Addresses[j]
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_addresses WHERE id=$1`
			row := db.QueryRow(sqlStatement, address.Id)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_addresses(id, customer_id, company, address1, address2, city, province, zip, country_code, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
				db.Exec(sqlStatement, address.Id, address.CustomerId, address.Company, address.Address1, address.Address2, address.City, address.Province, address.Zip, address.CountryCode, enterpriseId)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_addresses SET customer_id=$2, company=$3, address1=$4, address2=$5, city=$6, province=$7, zip=$8, country_code=$9, sy_exists=true WHERE id=$1 AND enterprise=$10`
				db.Exec(sqlStatement, address.Id, address.CustomerId, address.Company, address.Address1, address.Address2, address.City, address.Province, address.Zip, address.CountryCode, enterpriseId)
			}

		}
	}

	sqlStatement = `DELETE FROM public.sy_customers WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.sy_addresses WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	return true
}

func importSyProducts(enterpriseId int32) bool {
	url := getShopifyAPI_URL("products", enterpriseId)
	jsonSY, err := getShopifyJSON(url, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Products</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var products SYProducts
	json.Unmarshal(jsonSY, &products)

	sqlStatement := `UPDATE public.sy_products SET sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `UPDATE public.sy_variants SET sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(products.Products); i++ {
		product := products.Products[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_products WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, product.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_products(id, title, body_html, enterprise) VALUES ($1, $2, $3, $4)`
			db.Exec(sqlStatement, product.Id, product.Title, strip.StripTags(product.BodyHtml), enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_products SET title=$2, body_html=$3, sy_exists=true WHERE id=$1 AND enterprise=$4`
			db.Exec(sqlStatement, product.Id, product.Title, strip.StripTags(product.BodyHtml), enterpriseId)
		}

		// product variants
		for j := 0; j < len(product.Variants); j++ {
			variant := product.Variants[j]
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_variants WHERE id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, variant.Id, enterpriseId)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_variants(id, product_id, title, price, sku, option1, option2, option3, taxable, barcode, grams, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
				db.Exec(sqlStatement, variant.Id, variant.ProductId, variant.Title, variant.Price, variant.Sku, variant.Option1, variant.Option2, variant.Option3, variant.Taxable, variant.Barcode, variant.Grams, enterpriseId)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_variants SET product_id=$2, title=$3, price=$4, sku=$5, option1=$6, option2=$7, option3=$8, taxable=$9, barcode=$10, grams=$11, sy_exists=true WHERE id=$1 AND enterprise=$12`
				db.Exec(sqlStatement, variant.Id, variant.ProductId, variant.Title, variant.Price, variant.Sku, variant.Option1, variant.Option2, variant.Option3, variant.Taxable, variant.Barcode, variant.Grams, enterpriseId)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_products WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.sy_variants WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement)

	return true
}

func importSyDraftOrders(enterpriseId int32) bool {
	url := getShopifyAPI_URL("draft_orders", enterpriseId)
	jsonSY, err := getShopifyJSON(url, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Draft orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var orders SYDraftOrders
	json.Unmarshal(jsonSY, &orders)

	sqlStatement := `UPDATE public.sy_draft_orders SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `UPDATE public.sy_draft_order_line_item SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(orders.DraftOrders); i++ {
		draftOrder := orders.DraftOrders[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_draft_orders WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, draftOrder.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_draft_orders(id, currency, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, order_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
			db.Exec(sqlStatement, draftOrder.Id, draftOrder.Currency, draftOrder.TaxExempt, draftOrder.Name, draftOrder.ShippingAddress.Address1, draftOrder.ShippingAddress.Address2, draftOrder.ShippingAddress.City, draftOrder.ShippingAddress.Zip, draftOrder.ShippingAddress.CountryCode, draftOrder.BillingAddress.Address1, draftOrder.BillingAddress.Address2, draftOrder.BillingAddress.City, draftOrder.BillingAddress.Zip, draftOrder.BillingAddress.CountryCode, draftOrder.TotalTax, draftOrder.Customer.Id, draftOrder.OrderId, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_draft_orders SET currency=$2, tax_exempt=$3, name=$4, shipping_address_1=$5, shipping_address2=$6, shipping_address_city=$7, shipping_address_zip=$8, shipping_address_country_code=$9, billing_address_1=$10, billing_address2=$11, billing_address_city=$12, billing_address_zip=$13, billing_address_country_code=$14, total_tax=$15, customer_id=$16, sy_exists=true, order_id=$17 WHERE id=$1 AND enterprise=$18`
			db.Exec(sqlStatement, draftOrder.Id, draftOrder.Currency, draftOrder.TaxExempt, draftOrder.Name, draftOrder.ShippingAddress.Address1, draftOrder.ShippingAddress.Address2, draftOrder.ShippingAddress.City, draftOrder.ShippingAddress.Zip, draftOrder.ShippingAddress.CountryCode, draftOrder.BillingAddress.Address1, draftOrder.BillingAddress.Address2, draftOrder.BillingAddress.City, draftOrder.BillingAddress.Zip, draftOrder.BillingAddress.CountryCode, draftOrder.TotalTax, draftOrder.Customer.Id, draftOrder.OrderId, enterpriseId)
		}

		for j := 0; j < len(draftOrder.LineItems); j++ {
			lineItem := draftOrder.LineItems[j]
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_draft_order_line_item WHERE id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, lineItem.Id, enterpriseId)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_draft_order_line_item(id, variant_id, product_id, quantity, taxable, price, draft_order_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
				db.Exec(sqlStatement, lineItem.Id, lineItem.VariantId, lineItem.ProductId, lineItem.Quantity, lineItem.Taxable, lineItem.Price, draftOrder.Id, enterpriseId)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_draft_order_line_item SET variant_id=$2, product_id=$3, quantity=$4, taxable=$5, price=$6, draft_order_id=$7, sy_exists=true WHERE id=$1 AND enterprise=$8`
				db.Exec(sqlStatement, lineItem.Id, lineItem.VariantId, lineItem.ProductId, lineItem.Quantity, lineItem.Taxable, lineItem.Price, draftOrder.Id, enterpriseId)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_draft_orders WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.sy_draft_order_line_item WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	return true
}

func importSyOrders(enterpriseId int32) bool {
	url := getShopifyAPI_URL("orders", enterpriseId)
	jsonSY, err := getShopifyJSON(url, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var orders SYOrders
	json.Unmarshal(jsonSY, &orders)

	sqlStatement := `UPDATE public.sy_orders SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `UPDATE public.sy_order_line_item SET sy_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(orders.Orders); i++ {
		order := orders.Orders[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_orders WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, order.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_orders(id, currency, current_total_discounts, total_shipping_price_set_amount, total_shipping_price_set_currency_code, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, gateway, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`
			db.Exec(sqlStatement, order.Id, order.Currency, order.CurrentTotalDiscounts, order.TotalShippingPriceSet.ShopMoney.Amount, order.TotalShippingPriceSet.ShopMoney.CurrencyCode, order.Customer.TaxExempt, order.Name, order.ShippingAddress.Address1, order.ShippingAddress.Address2, order.ShippingAddress.City, order.ShippingAddress.Zip, order.ShippingAddress.CountryCode, order.BillingAddress.Address1, order.BillingAddress.Address2, order.BillingAddress.City, order.BillingAddress.Zip, order.BillingAddress.CountryCode, order.TotalTax, order.Customer.Id, order.Gateway, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_orders SET currency=$2, current_total_discounts=$3, total_shipping_price_set_amount=$4, total_shipping_price_set_currency_code=$5, tax_exempt=$6, name=$7, shipping_address_1=$8, shipping_address2=$9, shipping_address_city=$10, shipping_address_zip=$11, shipping_address_country_code=$12, billing_address_1=$13, billing_address2=$14, billing_address_city=$15, billing_address_zip=$16, billing_address_country_code=$17, total_tax=$18, customer_id=$19, sy_exists=true, gateway=$20 WHERE id=$1 AND enterprise=$21`
			db.Exec(sqlStatement, order.Id, order.Currency, order.CurrentTotalDiscounts, order.TotalShippingPriceSet.ShopMoney.Amount, order.TotalShippingPriceSet.ShopMoney.CurrencyCode, order.Customer.TaxExempt, order.Name, order.ShippingAddress.Address1, order.ShippingAddress.Address2, order.ShippingAddress.City, order.ShippingAddress.Zip, order.ShippingAddress.CountryCode, order.BillingAddress.Address1, order.BillingAddress.Address2, order.BillingAddress.City, order.BillingAddress.Zip, order.BillingAddress.CountryCode, order.TotalTax, order.Customer.Id, order.Gateway, enterpriseId)
		}

		for j := 0; j < len(order.LineItems); j++ {
			lineItem := order.LineItems[j]
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_order_line_item WHERE id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, lineItem.Id, enterpriseId)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_order_line_item(id, variant_id, product_id, quantity, taxable, price, order_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
				db.Exec(sqlStatement, lineItem.Id, lineItem.VariantId, lineItem.ProductId, lineItem.Quantity, lineItem.Taxable, lineItem.Price, order.Id, enterpriseId)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_order_line_item SET variant_id=$2, product_id=$3, quantity=$4, taxable=$5, price=$6, order_id=$7, sy_exists=true WHERE id=$1 AND enterprise=$8`
				db.Exec(sqlStatement, lineItem.Id, lineItem.VariantId, lineItem.ProductId, lineItem.Quantity, lineItem.Taxable, lineItem.Price, order.Id, enterpriseId)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_orders WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.sy_order_line_item WHERE sy_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	return true
}

// =====
// TRANSFER THE DATA TO THE ERP TABLES
// =====

func copySyCustomers(enterpriseId int32) bool {
	sqlStatement := `SELECT id FROM public.sy_customers WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var syCustomerId int64
		rows.Scan(&syCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM customer WHERE sy_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, syCustomerId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, email, first_name, last_name, phone, default_address_id FROM public.sy_customers WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, syCustomerId, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int64
		var email string
		var firstName string
		var lastName string
		var phone string
		var defaultAddressId int64
		row.Scan(&id, &email, &firstName, &lastName, &phone, &defaultAddressId)

		if rows == 0 {
			// create customer
			c := Customer{}
			c.Email = email
			c.ShopifyId = id
			c.Phone = phone
			c.Tradename = firstName + " " + lastName

			sqlStatement := `SELECT company FROM public.sy_addresses WHERE id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, defaultAddressId, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}
			var company string

			if len(company) > 0 {
				c.FiscalName = company
				c.Name = c.FiscalName + " - " + c.Tradename
			} else {
				c.FiscalName = c.Tradename
				c.Name = c.Tradename
			}

			c.EnterpriseId = enterpriseId
			res := c.insertCustomer(0)
			ok, customerId := res.Id > 0, int32(res.Id)
			if !ok {
				continue
			}

			// add addresses
			sqlStatement = `SELECT id, address1, address2, city, province, zip, country_code FROM public.sy_addresses WHERE customer_id=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, id, enterpriseId)
			if err != nil {
				log("Shopify", rows.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id int64
				var address1 string
				var address2 string
				var city string
				var province string
				var zip string
				var countryCode string
				rows.Scan(&id, &address1, &address2, &city, &province, &zip, &countryCode)

				// get country
				var countryId int32
				sqlStatement := `SELECT id FROM public.country WHERE (iso_2=$1 OR iso_3=$1) AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, countryCode, enterpriseId)
				if row.Err() != nil {
					log("Shopify", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}
				row.Scan(&countryId)

				// get province
				var provinceId *int32
				sqlStatement = `SELECT id FROM public.state WHERE name=$1 AND enterprise=$2 LIMIT 1`
				row = db.QueryRow(sqlStatement, province, enterpriseId)
				if row.Err() != nil {
					log("Shopify", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var stateId int32
				row.Scan(&stateId)
				if stateId > 0 {
					provinceId = &stateId
				}

				a := Address{}
				a.CustomerId = &customerId
				a.Address = address1
				a.Address2 = address2
				a.City = city
				a.ZipCode = zip
				a.CountryId = countryId
				a.StateId = provinceId
				a.ShopifyId = id
				if len(company) > 0 {
					a.PrivateOrBusiness = "B"
				} else {
					a.PrivateOrBusiness = "P"
				}
				a.EnterpriseId = enterpriseId
				a.insertAddress(0)
			} // for rows.Next()
		} else { // if rows == 0
			// update the customer
			sqlStatement := `SELECT id FROM customer WHERE sy_id=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var erpCustomerId int32
			row.Scan(&erpCustomerId)
			if erpCustomerId <= 0 {
				continue
			}

			c := getCustomerRow(erpCustomerId)
			c.Email = email
			c.ShopifyId = id
			c.Phone = phone
			c.Tradename = firstName + " " + lastName

			sqlStatement = `SELECT company FROM public.sy_addresses WHERE id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, defaultAddressId, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}
			var company string

			if len(company) > 0 {
				c.FiscalName = company
				c.Name = c.FiscalName + " - " + c.Tradename
			} else {
				c.FiscalName = c.Tradename
				c.Name = c.Tradename
			}
			c.updateCustomer(0)

			// add/update addresses
			sqlStatement = `SELECT id, address1, address2, city, province, zip, country_code FROM public.sy_addresses WHERE customer_id=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, id, enterpriseId)
			if err != nil {
				log("Shopify", rows.Err().Error())
				errors = append(errors, rows.Err().Error())
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id int64
				var address1 string
				var address2 string
				var city string
				var province string
				var zip string
				var countryCode string
				rows.Scan(&id, &address1, &address2, &city, &province, &zip, &countryCode)

				// get country
				var countryId int32
				sqlStatement = `SELECT id FROM public.country WHERE (iso_2=$1 OR iso_3=$1) AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, countryCode, enterpriseId)
				if row.Err() != nil {
					log("Shopify", rows.Err().Error())
					errors = append(errors, rows.Err().Error())
					continue
				}
				row.Scan(&countryId)

				// get province
				var provinceId *int32
				sqlStatement = `SELECT id FROM public.state WHERE name=$1 AND enterprise=$2 LIMIT 1`
				row = db.QueryRow(sqlStatement, province, enterpriseId)
				if row.Err() != nil {
					log("Shopify", rows.Err().Error())
					errors = append(errors, rows.Err().Error())
					continue
				}

				var stateId int32
				row.Scan(&stateId)
				if stateId > 0 {
					provinceId = &stateId
				}

				// ¿does the row exist?
				sqlStatement := `SELECT COUNT(*) FROM address WHERE sy_id=$1 AND enterprise=$2`
				rowCount := db.QueryRow(sqlStatement, id, enterpriseId)
				var rows int32
				rowCount.Scan(&rows)

				if rows == 0 {
					a := Address{}
					a.CustomerId = &c.Id
					a.Address = address1
					a.Address2 = address2
					a.City = city
					a.ZipCode = zip
					a.CountryId = countryId
					a.StateId = provinceId
					a.ShopifyId = id
					if len(company) > 0 {
						a.PrivateOrBusiness = "B"
					} else {
						a.PrivateOrBusiness = "P"
					}
					a.EnterpriseId = enterpriseId
					a.insertAddress(0)
				} else {
					sqlStatement := `SELECT id FROM address WHERE sy_id=$1 AND enterprise=$2`
					rowCount := db.QueryRow(sqlStatement, id, enterpriseId)
					var addressId int32
					rowCount.Scan(&addressId)

					a := getAddressRow(addressId)
					a.CustomerId = &c.Id
					a.Address = address1
					a.Address2 = address2
					a.City = city
					a.ZipCode = zip
					a.CountryId = countryId
					a.StateId = provinceId
					if len(company) > 0 {
						a.PrivateOrBusiness = "B"
					} else {
						a.PrivateOrBusiness = "P"
					}
					a.updateAddress()
				}
			} // for rows.Next()
		}
	} // for rows.Next()

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Customers</p>"+errorHtml, enterpriseId)
	}

	return true
} // copySyCustomers()

func copySyProducts(enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id FROM public.sy_products WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Products</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var syProductId int64
		rows.Scan(&syProductId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM product WHERE sy_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, syProductId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, title, body_html FROM public.sy_products WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, syProductId, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int64
		var title string
		var bodyHtml string
		row.Scan(&id, &title, &bodyHtml)

		var variants []SYVariantDB = make([]SYVariantDB, 0)
		sqlStatement = `SELECT id, product_id, title, price, sku, option1, option2, option3, taxable, barcode, grams FROM public.sy_variants WHERE product_id=$1 AND enterprise=$2`
		rowsVariants, err := db.Query(sqlStatement, id, enterpriseId)
		if err != nil {
			log("Shopify", err.Error())
			errors = append(errors, err.Error())
			continue
		}
		defer rowsVariants.Close()

		for rowsVariants.Next() {
			v := SYVariantDB{}
			rowsVariants.Scan(&v.Id, &v.ProductId, &v.Title, &v.Price, &v.Sku, &v.Option1, &v.Option2, &v.Option3, &v.Taxable, &v.Barcode, &v.Grams)
			variants = append(variants, v)
		}

		if rows == 0 {
			// if the product uses variants, crate a product on the ERP for every single variant, or, if there is only one variant, create a single product on the ERP
			if len(variants) == 1 {
				p := Product{}
				p.ShopifyId = id
				p.ShopifyVariantId = variants[0].Id
				p.Name = title
				p.Description = bodyHtml
				p.Price = variants[0].Price
				p.Weight = float64(variants[0].Grams / 1000)
				if len(variants[0].Sku) > 0 {
					p.BarCode = variants[0].Sku
				} else {
					p.BarCode = variants[0].Barcode
				}
				if variants[0].Taxable {
					p.VatPercent = s.DefaultVatPercent
				} else {
					p.VatPercent = 0
				}
				p.EnterpriseId = enterpriseId
				result := p.insertProduct(0)
				if !result.Ok {
					errors = append(errors, "Error inserting a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}
			} else {
				for i := 0; i < len(variants); i++ {
					p := Product{}
					p.ShopifyId = id
					p.ShopifyVariantId = variants[i].Id
					p.Name = title + " " + variants[i].Option1
					if variants[i].Option2 != nil {
						p.Name += " " + *variants[i].Option2
					}
					if variants[i].Option3 != nil {
						p.Name += " " + *variants[i].Option3
					}
					p.Description = bodyHtml
					p.Price = variants[i].Price
					p.Weight = float64(variants[i].Grams / 1000)
					if len(variants[i].Sku) > 0 {
						p.BarCode = variants[i].Sku
					} else {
						p.BarCode = variants[i].Barcode
					}
					if variants[i].Taxable {
						p.VatPercent = s.DefaultVatPercent
					} else {
						p.VatPercent = 0
					}
					p.EnterpriseId = enterpriseId
					result := p.insertProduct(0)
					if !result.Ok {
						errors = append(errors, "Error inserting a product with combinations into MARKETNET. Product name "+
							p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
					}
				}
			}
		} else { // if rows == 0
			// if the product uses variants, crate a product on the ERP for every single variant, or, if there is only one variant, create a single product on the ERP
			if len(variants) == 1 {
				sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 AND enterprise=$3 LIMIT 1`
				row := db.QueryRow(sqlStatement, id, variants[0].Id, enterpriseId)
				if row.Err() != nil {
					log("Shopify", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var productId int32
				row.Scan(&productId)

				if productId <= 0 {
					continue
				}

				p := getProductRow(productId)
				p.Name = title
				p.Description = bodyHtml
				p.Price = variants[0].Price
				p.Weight = float64(variants[0].Grams / 1000)
				if len(variants[0].Sku) > 0 {
					p.BarCode = variants[0].Sku
				} else {
					p.BarCode = variants[0].Barcode
				}
				if variants[0].Taxable {
					p.VatPercent = s.DefaultVatPercent
				} else {
					p.VatPercent = 0
				}
				result := p.updateProduct(0)
				if !result.Ok {
					errors = append(errors, "Error updating a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}
			} else {
				for i := 0; i < len(variants); i++ {
					sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 AND enterprise=$3 LIMIT 1`
					row := db.QueryRow(sqlStatement, id, variants[i].Id, enterpriseId)
					if row.Err() != nil {
						log("Shopify", row.Err().Error())
						errors = append(errors, row.Err().Error())
						continue
					}

					var productId int32
					row.Scan(&productId)

					if productId > 0 { // the variant already exists
						p := getProductRow(productId)
						p.Name = title + " " + variants[i].Option1
						if variants[i].Option2 != nil {
							p.Name += " " + *variants[i].Option2
						}
						if variants[i].Option3 != nil {
							p.Name += " " + *variants[i].Option3
						}
						p.Description = bodyHtml
						p.Price = variants[i].Price
						p.Weight = float64(variants[i].Grams / 1000)
						if len(variants[i].Sku) > 0 {
							p.BarCode = variants[i].Sku
						} else {
							p.BarCode = variants[i].Barcode
						}
						if variants[i].Taxable {
							p.VatPercent = s.DefaultVatPercent
						} else {
							p.VatPercent = 0
						}
						result := p.updateProduct(0)
						if !result.Ok {
							errors = append(errors, "Error updating a product with combinations into MARKETNET. Product name "+
								p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
						}
					} else { // the variant does not exist
						p := Product{}
						p.ShopifyId = id
						p.ShopifyVariantId = variants[i].Id
						p.Name = title + " " + variants[i].Option1
						if variants[i].Option2 != nil {
							p.Name += " " + *variants[i].Option2
						}
						if variants[i].Option3 != nil {
							p.Name += " " + *variants[i].Option3
						}
						p.Description = bodyHtml
						p.Price = variants[i].Price
						p.Weight = float64(variants[i].Grams / 1000)
						if len(variants[i].Sku) > 0 {
							p.BarCode = variants[i].Sku
						} else {
							p.BarCode = variants[i].Barcode
						}
						if variants[i].Taxable {
							p.VatPercent = s.DefaultVatPercent
						} else {
							p.VatPercent = 0
						}
						p.EnterpriseId = enterpriseId
						result := p.insertProduct(0)
						if !result.Ok {
							errors = append(errors, "Error inserting a product with combinations into MARKETNET. Product name "+
								p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
						}
					}
				} // for
			} // else
		} // else
	} // for rows.Next()

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Products</p>"+errorHtml, enterpriseId)
	}

	return true
} // copySyProducts

func copySyDraftOrders(enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id FROM public.sy_draft_orders WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Draft orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var syDraftOrderId int64
		rows.Scan(&syDraftOrderId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sales_order WHERE sy_draft_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, syDraftOrderId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, currency, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id FROM public.sy_draft_orders WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, syDraftOrderId, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int64
		var currency string
		var taxExempt bool
		var name string
		var shippingAddress1 string
		var shippingAddress2 string
		var shippingAddressCity string
		var shippingAddressZip string
		var shippingAddressCountryCode string
		var billingAddress1 string
		var billingAddress2 string
		var billingAddressCity string
		var billingAddressZip string
		var billingAddressCountryCode string
		var totalTax float64
		var customerId int64
		row.Scan(&id, &currency, &taxExempt, &name, &shippingAddress1, &shippingAddress2, &shippingAddressCity, &shippingAddressZip, &shippingAddressCountryCode, &billingAddress1, &billingAddress2, &billingAddressCity, &billingAddressZip, &billingAddressCountryCode, &totalTax, &customerId)

		// get the currency
		var currencyId int32
		sqlStatement = `SELECT id FROM public.currency WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, currency, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		row.Scan(&currencyId)
		if currencyId <= 0 {
			errors = append(errors, "Can't import draft order. The currency does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
			continue
		}

		// get the customer
		var customerIdErp int32
		sqlStatement = `SELECT id FROM customer WHERE sy_id=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, customerId, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		row.Scan(&customerIdErp)
		if customerIdErp <= 0 {
			errors = append(errors, "Can't import draft order. The customer does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
			continue
		}

		// get the billing address
		var billingAddressId int32
		var billingZone string
		sqlStatement = `SELECT id,(SELECT zone FROM country WHERE country.id=address.country) FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 AND enterprise=$6 LIMIT 1`
		row = db.QueryRow(sqlStatement, billingAddress1, billingAddress2, billingAddressCity, billingAddressZip, billingAddressCountryCode, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		row.Scan(&billingAddressId, &billingZone)
		if billingAddressId <= 0 {
			errors = append(errors, "Can't import draft order. The billing address does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
			continue
		}

		// get the shipping address
		var shippingAddressId int32
		sqlStatement = `SELECT id FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 AND enterprise=$6 LIMIT 1`
		row = db.QueryRow(sqlStatement, shippingAddress1, shippingAddress2, shippingAddressCity, shippingAddressZip, shippingAddressCountryCode, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		row.Scan(&shippingAddressId)
		if shippingAddressId <= 0 {
			errors = append(errors, "Can't import draft order. The shipping address does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
			continue
		}

		if rows == 0 {
			o := SaleOrder{}
			o.BillingAddressId = billingAddressId
			o.ShippingAddressId = shippingAddressId
			o.CustomerId = customerIdErp
			o.Reference = name
			o.CurrencyId = currencyId
			o.PaymentMethodId = *s.SettingsEcommerce.ShopifyDefaultPaymentMethodId

			if billingZone == "E" {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyExportSerieId
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyIntracommunitySerieId
			} else {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyInteriorSerieId
			}

			o.ShopifyDraftId = id
			o.EnterpriseId = enterpriseId
			ok, orderId := o.insertSalesOrder(0)
			if !ok {
				errors = append(errors, "Can't import draft order. The order could not be created in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeriesId == nil || *c.BillingSeriesId == "" {
				c.BillingSeriesId = &o.BillingSeriesId
			}
			c.updateCustomer(0)

			// insert the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_draft_order_line_item WHERE draft_order_id=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, id, enterpriseId)
			if err != nil {
				log("Shopify", err.Error())
				errors = append(errors, err.Error())
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float64
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 AND enterprise=$3 LIMIT 1`
				row := db.QueryRow(sqlStatement, productId, variantId, enterpriseId)
				if row.Err() != nil {
					log("Shopify", err.Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					errors = append(errors, "Can't import draft order detail. The product does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
					continue
				}

				d := SalesOrderDetail{}
				d.OrderId = orderId
				d.Quantity = quantity
				d.Price = price
				if taxable && !taxExempt {
					d.VatPercent = s.DefaultVatPercent
				} else {
					d.VatPercent = 0
				}
				d.ProductId = productIdErp
				d.ShopifyDraftId = id
				d.EnterpriseId = enterpriseId
				d.insertSalesOrderDetail(0)
			} // for rows.Next()
		} else { // if rows == 0
			var orderIdErp int64
			sqlStatement := `SELECT id FROM sales_order WHERE sy_draft_id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&orderIdErp)
			if orderIdErp <= 0 {
				errors = append(errors, "Can't import draft order. The draft order does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			o := getSalesOrderRow(orderIdErp)
			o.BillingAddressId = billingAddressId
			o.ShippingAddressId = shippingAddressId
			o.CustomerId = customerIdErp
			o.Reference = name
			o.CurrencyId = currencyId
			o.PaymentMethodId = *s.SettingsEcommerce.ShopifyDefaultPaymentMethodId

			if billingZone == "E" {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyExportSerieId
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyIntracommunitySerieId
			} else {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyInteriorSerieId
			}

			o.EnterpriseId = enterpriseId
			ok := o.updateSalesOrder(0)
			if !ok {
				errors = append(errors, "Can't import draft order. Can't update the order in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeriesId == nil || *c.BillingSeriesId == "" {
				c.BillingSeriesId = &o.BillingSeriesId
			}
			c.updateCustomer(0)

			// insert/update the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_draft_order_line_item WHERE draft_order_id=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, id, enterpriseId)
			if err != nil {
				log("Shopify", err.Error())
				errors = append(errors, err.Error())
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float64
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				var salesOrderDetailId int64
				sqlStatement := `SELECT id FROM sales_order_detail WHERE sy_draft_id=$1 AND enterprise=$2`
				row := db.QueryRow(sqlStatement, id, enterpriseId)
				if row.Err() != nil {
					log("Shopify", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				row.Scan(&salesOrderDetailId)

				sqlStatement = `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 AND enterprise=$3 LIMIT 1`
				row = db.QueryRow(sqlStatement, productId, variantId, enterpriseId)
				if row.Err() != nil {
					log("Shopify", row.Err().Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					errors = append(errors, "Can't import draft order detail. The product does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
					continue
				}

				if salesOrderDetailId <= 0 {
					d := SalesOrderDetail{}
					d.OrderId = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.ProductId = productIdErp
					d.ShopifyDraftId = id
					d.EnterpriseId = enterpriseId
					d.insertSalesOrderDetail(0)
				} else { // if salesOrderDetailId <= 0
					d := getSalesOrderDetailRow(salesOrderDetailId)
					d.OrderId = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.ProductId = productIdErp
					d.EnterpriseId = enterpriseId
					d.updateSalesOrderDetail(0)
				}
			} // for rows.Next()
		} // else
	} // for rows.Next()

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Draft orders</p>"+errorHtml, enterpriseId)
	}

	return true
} // copySyDraftOrders

func copySyOrders(enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id FROM public.sy_orders WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("Shopify", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var syOrderId int64
		rows.Scan(&syOrderId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sales_order WHERE sy_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, syOrderId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, currency, current_total_discounts, total_shipping_price_set_amount, total_shipping_price_set_currency_code, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, gateway FROM public.sy_orders WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, syOrderId, enterpriseId)
		if row.Err() != nil {
			log("Shopify", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int64
		var currency string
		var currentTotalDiscounts float64
		var totalShippingPriceSetAmount float64
		var totalShippingPriceSetCurrencyCode string
		var taxExempt bool
		var name string
		var shippingAddress1 string
		var shippingAddress2 string
		var shippingAddressCity string
		var shippingAddressZip string
		var shippingAddressCountryCode string
		var billingAddress1 string
		var billingAddress2 string
		var billingAddressCity string
		var billingAddressZip string
		var billingAddressCountryCode string
		var totalTax float64
		var customerId int64
		var gateway string
		row.Scan(&id, &currency, &currentTotalDiscounts, &totalShippingPriceSetAmount, &totalShippingPriceSetCurrencyCode, &taxExempt, &name, &shippingAddress1, &shippingAddress2, &shippingAddressCity, &shippingAddressZip, &shippingAddressCountryCode, &billingAddress1, &billingAddress2, &billingAddressCity, &billingAddressZip, &billingAddressCountryCode, &totalTax, &customerId, &gateway)

		if rows == 0 {
			// the order has been paid in Shopify, update one last time the order and set invoice it

			// get the currency
			var currencyId int32
			sqlStatement = `SELECT id FROM public.currency WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
			row = db.QueryRow(sqlStatement, currency, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&currencyId)
			if currencyId <= 0 {
				errors = append(errors, "Can't import order. The currency does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// get the customer
			var customerIdErp int32
			sqlStatement = `SELECT id FROM customer WHERE sy_id=$1 AND enterprise=$2 LIMIT 1`
			row = db.QueryRow(sqlStatement, customerId, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&customerIdErp)
			if customerIdErp <= 0 {
				errors = append(errors, "Can't import order. The customer does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// get the billing address
			var billingAddressId int32
			var billingZone string
			sqlStatement = `SELECT id,(SELECT zone FROM country WHERE country.id=address.country) FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 AND enterprise=$6 LIMIT 1`
			row = db.QueryRow(sqlStatement, billingAddress1, billingAddress2, billingAddressCity, billingAddressZip, billingAddressCountryCode, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&billingAddressId, &billingZone)
			if billingAddressId <= 0 {
				errors = append(errors, "Can't import order. The billiong addresses does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// get the shipping address
			var shippingAddressId int32
			sqlStatement = `SELECT id FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 AND enterprise=$6 LIMIT 1`
			row = db.QueryRow(sqlStatement, shippingAddress1, shippingAddress2, shippingAddressCity, shippingAddressZip, shippingAddressCountryCode, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			row.Scan(&shippingAddressId)
			if shippingAddressId <= 0 {
				errors = append(errors, "Can't import order. The shipping address does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// get the order id
			sqlStatement := `SELECT id FROM sy_draft_orders WHERE order_id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, id, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var draftOrderId int64
			row.Scan(&draftOrderId)
			if draftOrderId <= 0 {
				errors = append(errors, "Can't import order. The draft order does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			sqlStatement = `SELECT id FROM sales_order WHERE sy_draft_id=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, draftOrderId, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var saleOrderIdErp int64
			row.Scan(&saleOrderIdErp)
			if saleOrderIdErp <= 0 {
				errors = append(errors, "Can't import order. The sale order does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// get the payment method
			var paymentMethod int32
			sqlStatement = `SELECT id FROM payment_method WHERE shopify_module_name=$1 AND enterprise=$2`
			row = db.QueryRow(sqlStatement, gateway, enterpriseId)
			if row.Err() != nil {
				log("Shopify", row.Err().Error())
				errors = append(errors, row.Err().Error())
				paymentMethod = *s.SettingsEcommerce.ShopifyDefaultPaymentMethodId
			} else {
				row.Scan(&paymentMethod)
				if paymentMethod <= 0 {
					paymentMethod = *s.SettingsEcommerce.ShopifyDefaultPaymentMethodId
				}
			}

			// update the order
			o := getSalesOrderRow(saleOrderIdErp)
			o.BillingAddressId = billingAddressId
			o.ShippingAddressId = shippingAddressId
			o.CustomerId = customerIdErp
			o.Reference = name
			o.CurrencyId = currencyId
			o.PaymentMethodId = paymentMethod

			if billingZone == "E" {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyExportSerieId
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyIntracommunitySerieId
			} else {
				o.BillingSeriesId = *s.SettingsEcommerce.ShopifyInteriorSerieId
			}

			o.ShopifyId = id
			o.EnterpriseId = enterpriseId
			ok := o.updateSalesOrder(0)
			if !ok {
				errors = append(errors, "Can't import order. Can't update the existing sale order in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeriesId == nil || *c.BillingSeriesId == "" {
				c.BillingSeriesId = &o.BillingSeriesId
			}
			c.updateCustomer(0)

			// update the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_order_line_item WHERE order_id=$1 AND enterprise=$2`
			rows, err := db.Query(sqlStatement, id, enterpriseId)
			if err != nil {
				log("Shopify", err.Error())
				errors = append(errors, err.Error())
				continue
			}
			defer rows.Close()

			details := getSalesOrderDetail(o.Id, enterpriseId)

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float64
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				sqlStatement = `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 AND enterprise=$3 LIMIT 1`
				row = db.QueryRow(sqlStatement, productId, variantId, enterpriseId)
				if row.Err() != nil {
					log("Shopify", err.Error())
					errors = append(errors, row.Err().Error())
					continue
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					errors = append(errors, "Can't import order detail. The product does not exist in MARKETNET. Order id + "+strconv.Itoa(int(id))+" name "+name)
					continue
				}

				var salesOrderDetailId int64
				for i := 0; i < len(details); i++ {
					if details[i].ProductId == productIdErp {
						salesOrderDetailId = details[i].Id
						break
					}
				}

				if salesOrderDetailId > 0 {
					d := getSalesOrderDetailRow(salesOrderDetailId)
					d.OrderId = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.ProductId = productIdErp
					d.ShopifyId = id
					d.EnterpriseId = enterpriseId
					d.updateSalesOrderDetail(0)
				}
			}

			// create the invoice for the order
			invoiceAllSaleOrder(o.Id, enterpriseId, 0)
		} // if rows == 0
	} // for rows.Next()

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "Shopify import error", "<p>Error at: Orders</p>"+errorHtml, enterpriseId)
	}

	return true
} // copySyOrders

type SYFulfillmentContainer struct {
	Fulfillment SYFulfillment `json:"fulfillment"`
}

type SYFulfillment struct {
	LocationId      int64                   `json:"location_id"`
	TrackingNumber  string                  `json:"tracking_number"`
	TrackingCompany string                  `json:"tracking_company"`
	LineItems       []SYFulfillmentLineItem `json:"line_items"`
}

type SYFulfillmentLineItem struct {
	Id int64 `json:"id"`
}

type SYLocations struct {
	Locations []SYLocation `json:"locations"`
}

type SYLocation struct {
	Id int64 `json:"id"`
}

func updateTrackingNumberShopifyOrder(salesOrderId int64, trackingNumber string, enterpriseId int32) bool {
	fulfillment := SYFulfillmentContainer{}
	fulfillment.Fulfillment = SYFulfillment{}
	fulfillment.Fulfillment.TrackingNumber = trackingNumber
	order := getSalesOrderRow(salesOrderId)
	settings := getSettingsRecordById(enterpriseId)
	// we need to obtain the location_id for the fulfillment
	if settings.SettingsEcommerce.ShopifyShopLocationId <= 0 {
		url := getShopifyAPI_URL("locations", enterpriseId)
		jsonSY, err := getShopifyJSON(url, enterpriseId)
		if err != nil {
			return false
		}

		var locations SYLocations
		json.Unmarshal(jsonSY, &locations)

		if len(locations.Locations) == 1 {
			fulfillment.Fulfillment.LocationId = locations.Locations[0].Id
		} else {
			return false
		}
	} else {
		fulfillment.Fulfillment.LocationId = settings.SettingsEcommerce.ShopifyShopLocationId
	}

	// get the name of the carrier
	if order.CarrierId == nil {
		return false
	}
	carrier := getCarierRow(*order.CarrierId)
	fulfillment.Fulfillment.TrackingCompany = carrier.Name

	// line items
	details := getSalesOrderDetail(salesOrderId, enterpriseId)
	for i := 0; i < len(details); i++ {
		fulfillment.Fulfillment.LineItems = append(fulfillment.Fulfillment.LineItems, SYFulfillmentLineItem{
			Id: details[i].ShopifyId,
		})
	}

	// send data
	data, _ := json.Marshal(fulfillment)
	url := getShopifyAPI_URL("orders/"+strconv.Itoa(int(order.ShopifyId))+"/fulfillments", enterpriseId)
	postShopifyJSON(url, data, enterpriseId)
	return true
}

type SYDraftOrderComplete struct {
	PaymentPending bool `json:"payment_pending"`
}

func updateStatusPaymentAcceptedShopify(salesOrderId int64, enterpriseId int32) bool {
	complete := SYDraftOrderComplete{PaymentPending: false}
	order := getSalesOrderRow(salesOrderId)
	// send data
	data, _ := json.Marshal(complete)
	url := getShopifyAPI_URL("draft_orders/"+strconv.Itoa(int(order.ShopifyDraftId))+"/complete", enterpriseId)
	putShopifyJSON(url, data, enterpriseId)
	return true
}
