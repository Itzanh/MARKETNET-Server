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

func getShopifyAPI_URL(resourceName string) string {
	s := getSettingsRecord()

	return s.ShopifyUrl + resourceName + ".json"
}

func getShopifyJSON(URL string) ([]byte, error) {
	s := getSettingsRecord()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Set("X-Shopify-Access-Token", s.ShopifyToken)
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

func postShopifyJSON(URL string, data []byte) ([]byte, error) {
	s := getSettingsRecord()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))
	req.Header.Set("X-Shopify-Access-Token", s.ShopifyToken)
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

func putShopifyJSON(URL string, data []byte) ([]byte, error) {
	s := getSettingsRecord()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(data))
	req.Header.Set("X-Shopify-Access-Token", s.ShopifyToken)
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
	Id             int64       `json:"id"`
	Email          string      `json:"email"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	TaxExempt      bool        `json:"tax_exempt"`
	Phone          string      `json:"phone"`
	Currency       string      `json:"currency"`
	Addresses      []SYAddress `json:"addresses"`
	DefaultAddress SYAddress   `json:"default_address"`
}

type SYAddress struct {
	Id          int64  `json:"id"`
	CustomerId  int64  `json:"customer_id"`
	Company     string `json:"company"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Zip         string `json:"zip"`
	CountryCode string `json:"country_code"`
}

type SYProducts struct {
	Products []SYProduct `json:"products"`
}

type SYProduct struct {
	Id       int64       `json:"id"`
	Title    string      `json:"title"`
	BodyHtml string      `json:"body_html"`
	Variants []SYVariant `json:"variants"`
}

type SYVariant struct {
	ProductId int64   `json:"product_id"`
	Id        int64   `json:"id"`
	Title     string  `json:"title"`
	Price     string  `json:"price"`
	Sku       string  `json:"sku"`
	Option1   string  `json:"option1"`
	Option2   *string `json:"option2"`
	Option3   *string `json:"option3"`
	Taxable   bool    `json:"taxable"`
	Barcode   string  `json:"barcode"`
	Grams     int32   `json:"grams"`
}

type SYVariantDB struct {
	ProductId int64   `json:"product_id"`
	Id        int64   `json:"id"`
	Title     string  `json:"title"`
	Price     float32 `json:"price"`
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
	Id              int64                  `json:"id"`
	Currency        string                 `json:"currency"`
	TaxExempt       bool                   `json:"tax_exempt"`
	Name            string                 `json:"name"`
	LineItems       []SYDraftOrderLineItem `json:"line_items"`
	ShippingAddress SYAddress              `json:"shipping_address"`
	BillingAddress  SYAddress              `json:"billing_address"`
	TotalTax        string                 `json:"total_tax"`
	OrderId         *int64                 `json:"order_id"`
	Customer        SYCustomer             `json:"customer"`
}

type SYDraftOrderLineItem struct {
	Id        int64  `json:"id"`
	VariantId int64  `json:"variant_id"`
	ProductId int64  `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Taxable   bool   `json:"taxable"`
	Price     string `json:"price"`
}

type SYOrders struct {
	Orders []SYOrder `json:"orders"`
}

type SYOrder struct {
	Id                    int64                  `json:"id"`
	Currency              string                 `json:"currency"`
	CurrentTotalDiscounts string                 `json:"current_total_discounts"`
	TotalShippingPriceSet TotalShippingPriceSet  `json:"total_shipping_price_set"`
	Name                  string                 `json:"name"`
	LineItems             []SYDraftOrderLineItem `json:"line_items"`
	ShippingAddress       SYAddress              `json:"shipping_address"`
	BillingAddress        SYAddress              `json:"billing_address"`
	TotalTax              string                 `json:"total_tax"`
	Customer              SYCustomer             `json:"customer"`
	Gateway               string                 `json:"gateway"`
}

type TotalShippingPriceSet struct {
	ShopMoney ShopMoney `json:"shop_money"`
}

type ShopMoney struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type SYOrderLineItem struct {
	Id        int64  `json:"id"`
	VariantId int64  `json:"variant_id"`
	ProductId int64  `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Taxable   bool   `json:"taxable"`
	Price     string `json:"price"`
}

// main import function
func importFromShopify() {
	s := getSettingsRecord()
	if s.Ecommerce != "S" {
		return
	}

	// get all data from Shopify, write it in tables like the ones that Shopify uses
	importSyCustomers()
	importSyProducts()
	importSyDraftOrders()
	importSyOrders()

	// trasnfer the data form the Shopify tables to the ERP
	copySyCustomers()
	copySyProducts()
	copySyDraftOrders()
	copySyOrders()
}

// =====
// COPY THE DATA FROM WOOCOMMERCE TO THE WC MARKETNET TABLES
// =====

func importSyCustomers() {
	url := getShopifyAPI_URL("customers")
	jsonSY, err := getShopifyJSON(url)
	if err != nil {
		return
	}

	var customers SYCustomers
	json.Unmarshal(jsonSY, &customers)

	sqlStatement := `UPDATE public.sy_customers SET sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `UPDATE public.sy_addresses SET sy_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(customers.Customers); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_customers WHERE id=$1`
		row := db.QueryRow(sqlStatement, customers.Customers[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_customers(id, email, first_name, last_name, tax_exempt, phone, currency, default_address_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			db.Exec(sqlStatement, customers.Customers[i].Id, customers.Customers[i].Email, customers.Customers[i].FirstName, customers.Customers[i].LastName, customers.Customers[i].TaxExempt, customers.Customers[i].Phone, customers.Customers[i].Currency, customers.Customers[i].DefaultAddress.Id)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_customers SET email=$2, first_name=$3, last_name=$4, tax_exempt=$5, phone=$6, currency=$7, sy_exists=true, default_address_id=$8 WHERE id=$1`
			db.Exec(sqlStatement, customers.Customers[i].Id, customers.Customers[i].Email, customers.Customers[i].FirstName, customers.Customers[i].LastName, customers.Customers[i].TaxExempt, customers.Customers[i].Phone, customers.Customers[i].Currency, customers.Customers[i].DefaultAddress.Id)
		}

		for j := 0; j < len(customers.Customers[i].Addresses); j++ {
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_addresses WHERE id=$1`
			row := db.QueryRow(sqlStatement, customers.Customers[i].Addresses[j].Id)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_addresses(id, customer_id, company, address1, address2, city, province, zip, country_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
				db.Exec(sqlStatement, customers.Customers[i].Addresses[j].Id, customers.Customers[i].Addresses[j].CustomerId, customers.Customers[i].Addresses[j].Company, customers.Customers[i].Addresses[j].Address1, customers.Customers[i].Addresses[j].Address2, customers.Customers[i].Addresses[j].City, customers.Customers[i].Addresses[j].Province, customers.Customers[i].Addresses[j].Zip, customers.Customers[i].Addresses[j].CountryCode)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_addresses SET customer_id=$2, company=$3, address1=$4, address2=$5, city=$6, province=$7, zip=$8, country_code=$9, sy_exists=true WHERE id=$1`
				db.Exec(sqlStatement, customers.Customers[i].Addresses[j].Id, customers.Customers[i].Addresses[j].CustomerId, customers.Customers[i].Addresses[j].Company, customers.Customers[i].Addresses[j].Address1, customers.Customers[i].Addresses[j].Address2, customers.Customers[i].Addresses[j].City, customers.Customers[i].Addresses[j].Province, customers.Customers[i].Addresses[j].Zip, customers.Customers[i].Addresses[j].CountryCode)
			}

		}
	}

	sqlStatement = `DELETE FROM public.sy_customers WHERE sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.sy_addresses WHERE sy_exists=false`
	db.Exec(sqlStatement)
}

func importSyProducts() {
	url := getShopifyAPI_URL("products")
	jsonSY, err := getShopifyJSON(url)
	if err != nil {
		return
	}

	var products SYProducts
	json.Unmarshal(jsonSY, &products)

	sqlStatement := `UPDATE public.sy_products SET sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `UPDATE public.sy_variants SET sy_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(products.Products); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_products WHERE id=$1`
		row := db.QueryRow(sqlStatement, products.Products[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_products(id, title, body_html) VALUES ($1, $2, $3)`
			db.Exec(sqlStatement, products.Products[i].Id, products.Products[i].Title, strip.StripTags(products.Products[i].BodyHtml))
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_products SET title=$2, body_html=$3, sy_exists=true WHERE id=$1`
			db.Exec(sqlStatement, products.Products[i].Id, products.Products[i].Title, strip.StripTags(products.Products[i].BodyHtml))
		}

		// product variants
		for j := 0; j < len(products.Products[i].Variants); j++ {
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_variants WHERE id=$1`
			row := db.QueryRow(sqlStatement, products.Products[i].Variants[j].Id)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_variants(id, product_id, title, price, sku, option1, option2, option3, taxable, barcode, grams) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
				db.Exec(sqlStatement, products.Products[i].Variants[j].Id, products.Products[i].Variants[j].ProductId, products.Products[i].Variants[j].Title, products.Products[i].Variants[j].Price, products.Products[i].Variants[j].Sku, products.Products[i].Variants[j].Option1, products.Products[i].Variants[j].Option2, products.Products[i].Variants[j].Option3, products.Products[i].Variants[j].Taxable, products.Products[i].Variants[j].Barcode, products.Products[i].Variants[j].Grams)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_variants SET product_id=$2, title=$3, price=$4, sku=$5, option1=$6, option2=$7, option3=$8, taxable=$9, barcode=$10, grams=$11, sy_exists=true WHERE id=$1`
				db.Exec(sqlStatement, products.Products[i].Variants[j].Id, products.Products[i].Variants[j].ProductId, products.Products[i].Variants[j].Title, products.Products[i].Variants[j].Price, products.Products[i].Variants[j].Sku, products.Products[i].Variants[j].Option1, products.Products[i].Variants[j].Option2, products.Products[i].Variants[j].Option3, products.Products[i].Variants[j].Taxable, products.Products[i].Variants[j].Barcode, products.Products[i].Variants[j].Grams)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_products WHERE sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.sy_variants WHERE sy_exists=false`
	db.Exec(sqlStatement)
}

func importSyDraftOrders() {
	url := getShopifyAPI_URL("draft_orders")
	jsonSY, err := getShopifyJSON(url)
	if err != nil {
		return
	}

	var orders SYDraftOrders
	json.Unmarshal(jsonSY, &orders)

	sqlStatement := `UPDATE public.sy_draft_orders SET sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `UPDATE public.sy_draft_order_line_item SET sy_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(orders.DraftOrders); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_draft_orders WHERE id=$1`
		row := db.QueryRow(sqlStatement, orders.DraftOrders[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_draft_orders(id, currency, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, order_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
			db.Exec(sqlStatement, orders.DraftOrders[i].Id, orders.DraftOrders[i].Currency, orders.DraftOrders[i].TaxExempt, orders.DraftOrders[i].Name, orders.DraftOrders[i].ShippingAddress.Address1, orders.DraftOrders[i].ShippingAddress.Address2, orders.DraftOrders[i].ShippingAddress.City, orders.DraftOrders[i].ShippingAddress.Zip, orders.DraftOrders[i].ShippingAddress.CountryCode, orders.DraftOrders[i].BillingAddress.Address1, orders.DraftOrders[i].BillingAddress.Address2, orders.DraftOrders[i].BillingAddress.City, orders.DraftOrders[i].BillingAddress.Zip, orders.DraftOrders[i].BillingAddress.CountryCode, orders.DraftOrders[i].TotalTax, orders.DraftOrders[i].Customer.Id, orders.DraftOrders[i].OrderId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_draft_orders SET currency=$2, tax_exempt=$3, name=$4, shipping_address_1=$5, shipping_address2=$6, shipping_address_city=$7, shipping_address_zip=$8, shipping_address_country_code=$9, billing_address_1=$10, billing_address2=$11, billing_address_city=$12, billing_address_zip=$13, billing_address_country_code=$14, total_tax=$15, customer_id=$16, sy_exists=true, order_id=$17 WHERE id=$1`
			db.Exec(sqlStatement, orders.DraftOrders[i].Id, orders.DraftOrders[i].Currency, orders.DraftOrders[i].TaxExempt, orders.DraftOrders[i].Name, orders.DraftOrders[i].ShippingAddress.Address1, orders.DraftOrders[i].ShippingAddress.Address2, orders.DraftOrders[i].ShippingAddress.City, orders.DraftOrders[i].ShippingAddress.Zip, orders.DraftOrders[i].ShippingAddress.CountryCode, orders.DraftOrders[i].BillingAddress.Address1, orders.DraftOrders[i].BillingAddress.Address2, orders.DraftOrders[i].BillingAddress.City, orders.DraftOrders[i].BillingAddress.Zip, orders.DraftOrders[i].BillingAddress.CountryCode, orders.DraftOrders[i].TotalTax, orders.DraftOrders[i].Customer.Id, orders.DraftOrders[i].OrderId)
		}

		for j := 0; j < len(orders.DraftOrders[i].LineItems); j++ {
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_draft_order_line_item WHERE id=$1`
			row := db.QueryRow(sqlStatement, orders.DraftOrders[i].LineItems[j].Id)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_draft_order_line_item(id, variant_id, product_id, quantity, taxable, price, draft_order_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
				db.Exec(sqlStatement, orders.DraftOrders[i].LineItems[j].Id, orders.DraftOrders[i].LineItems[j].VariantId, orders.DraftOrders[i].LineItems[j].ProductId, orders.DraftOrders[i].LineItems[j].Quantity, orders.DraftOrders[i].LineItems[j].Taxable, orders.DraftOrders[i].LineItems[j].Price, orders.DraftOrders[i].Id)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_draft_order_line_item SET variant_id=$2, product_id=$3, quantity=$4, taxable=$5, price=$6, draft_order_id=$7, sy_exists=true WHERE id=$1`
				db.Exec(sqlStatement, orders.DraftOrders[i].LineItems[j].Id, orders.DraftOrders[i].LineItems[j].VariantId, orders.DraftOrders[i].LineItems[j].ProductId, orders.DraftOrders[i].LineItems[j].Quantity, orders.DraftOrders[i].LineItems[j].Taxable, orders.DraftOrders[i].LineItems[j].Price, orders.DraftOrders[i].Id)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_draft_orders WHERE sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.sy_draft_order_line_item WHERE sy_exists=false`
	db.Exec(sqlStatement)
}

func importSyOrders() {
	url := getShopifyAPI_URL("orders")
	jsonSY, err := getShopifyJSON(url)
	if err != nil {
		return
	}

	var orders SYOrders
	json.Unmarshal(jsonSY, &orders)

	sqlStatement := `UPDATE public.sy_orders SET sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `UPDATE public.sy_order_line_item SET sy_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(orders.Orders); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sy_orders WHERE id=$1`
		row := db.QueryRow(sqlStatement, orders.Orders[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.sy_orders(id, currency, current_total_discounts, total_shipping_price_set_amount, total_shipping_price_set_currency_code, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, gateway) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
			db.Exec(sqlStatement, orders.Orders[i].Id, orders.Orders[i].Currency, orders.Orders[i].CurrentTotalDiscounts, orders.Orders[i].TotalShippingPriceSet.ShopMoney.Amount, orders.Orders[i].TotalShippingPriceSet.ShopMoney.CurrencyCode, orders.Orders[i].Customer.TaxExempt, orders.Orders[i].Name, orders.Orders[i].ShippingAddress.Address1, orders.Orders[i].ShippingAddress.Address2, orders.Orders[i].ShippingAddress.City, orders.Orders[i].ShippingAddress.Zip, orders.Orders[i].ShippingAddress.CountryCode, orders.Orders[i].BillingAddress.Address1, orders.Orders[i].BillingAddress.Address2, orders.Orders[i].BillingAddress.City, orders.Orders[i].BillingAddress.Zip, orders.Orders[i].BillingAddress.CountryCode, orders.Orders[i].TotalTax, orders.Orders[i].Customer.Id, orders.Orders[i].Gateway)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.sy_orders SET currency=$2, current_total_discounts=$3, total_shipping_price_set_amount=$4, total_shipping_price_set_currency_code=$5, tax_exempt=$6, name=$7, shipping_address_1=$8, shipping_address2=$9, shipping_address_city=$10, shipping_address_zip=$11, shipping_address_country_code=$12, billing_address_1=$13, billing_address2=$14, billing_address_city=$15, billing_address_zip=$16, billing_address_country_code=$17, total_tax=$18, customer_id=$19, sy_exists=true, gateway=$20 WHERE id=$1`
			db.Exec(sqlStatement, orders.Orders[i].Id, orders.Orders[i].Currency, orders.Orders[i].CurrentTotalDiscounts, orders.Orders[i].TotalShippingPriceSet.ShopMoney.Amount, orders.Orders[i].TotalShippingPriceSet.ShopMoney.CurrencyCode, orders.Orders[i].Customer.TaxExempt, orders.Orders[i].Name, orders.Orders[i].ShippingAddress.Address1, orders.Orders[i].ShippingAddress.Address2, orders.Orders[i].ShippingAddress.City, orders.Orders[i].ShippingAddress.Zip, orders.Orders[i].ShippingAddress.CountryCode, orders.Orders[i].BillingAddress.Address1, orders.Orders[i].BillingAddress.Address2, orders.Orders[i].BillingAddress.City, orders.Orders[i].BillingAddress.Zip, orders.Orders[i].BillingAddress.CountryCode, orders.Orders[i].TotalTax, orders.Orders[i].Customer.Id, orders.Orders[i].Gateway)
		}

		for j := 0; j < len(orders.Orders[i].LineItems); j++ {
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM sy_order_line_item WHERE id=$1`
			row := db.QueryRow(sqlStatement, orders.Orders[i].LineItems[j].Id)
			var rows int32
			row.Scan(&rows)

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.sy_order_line_item(id, variant_id, product_id, quantity, taxable, price, order_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
				db.Exec(sqlStatement, orders.Orders[i].LineItems[j].Id, orders.Orders[i].LineItems[j].VariantId, orders.Orders[i].LineItems[j].ProductId, orders.Orders[i].LineItems[j].Quantity, orders.Orders[i].LineItems[j].Taxable, orders.Orders[i].LineItems[j].Price, orders.Orders[i].Id)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.sy_order_line_item SET variant_id=$2, product_id=$3, quantity=$4, taxable=$5, price=$6, order_id=$7, sy_exists=true WHERE id=$1`
				db.Exec(sqlStatement, orders.Orders[i].LineItems[j].Id, orders.Orders[i].LineItems[j].VariantId, orders.Orders[i].LineItems[j].ProductId, orders.Orders[i].LineItems[j].Quantity, orders.Orders[i].LineItems[j].Taxable, orders.Orders[i].LineItems[j].Price, orders.Orders[i].Id)
			}
		}
	}

	sqlStatement = `DELETE FROM public.sy_orders WHERE sy_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.sy_order_line_item WHERE sy_exists=false`
	db.Exec(sqlStatement)
}

// =====
// TRANSFER THE DATA TO THE ERP TABLES
// =====

func copySyCustomers() {
	sqlStatement := `SELECT id FROM public.sy_customers`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var syCustomerId int64
		rows.Scan(&syCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM customer WHERE sy_id=$1`
		rowCount := db.QueryRow(sqlStatement, syCustomerId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, email, first_name, last_name, phone, default_address_id FROM public.sy_customers WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, syCustomerId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			return
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

			sqlStatement := `SELECT company FROM public.sy_addresses WHERE id=$1 LIMIT 1`
			row := db.QueryRow(sqlStatement, defaultAddressId)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				return
			}
			var company string

			if len(company) > 0 {
				c.FiscalName = company
				c.Name = c.FiscalName + " - " + c.Tradename
			} else {
				c.FiscalName = c.Tradename
				c.Name = c.Tradename
			}

			ok, customerId := c.insertCustomer()
			if !ok {
				continue
			}

			// add addresses
			sqlStatement = `SELECT id, address1, address2, city, province, zip, country_code FROM public.sy_addresses WHERE customer_id=$1`
			rows, err := db.Query(sqlStatement, id)
			if err != nil {
				log("DB", rows.Err().Error())
				return
			}

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
				var countryId int16
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 OR iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, countryCode)
				if row.Err() != nil {
					continue
				}
				row.Scan(&countryId)

				// get province
				var provinceId *int32
				sqlStatement = `SELECT id FROM public.state WHERE name=$1 LIMIT 1`
				row = db.QueryRow(sqlStatement, province)
				if row.Err() != nil {
					continue
				}

				var stateId int32
				row.Scan(&stateId)
				if stateId > 0 {
					provinceId = &stateId
				}

				a := Address{}
				a.Customer = &customerId
				a.Address = address1
				a.Address2 = address2
				a.City = city
				a.ZipCode = zip
				a.Country = countryId
				a.State = provinceId
				a.ShopifyId = id
				if len(company) > 0 {
					a.PrivateOrBusiness = "B"
				} else {
					a.PrivateOrBusiness = "P"
				}
				a.insertAddress()
			} // for rows.Next()
		} else { // if rows == 0
			// update the customer
			sqlStatement := `SELECT id FROM customer WHERE sy_id=$1`
			row = db.QueryRow(sqlStatement, id)
			if row.Err() != nil {
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

			sqlStatement = `SELECT company FROM public.sy_addresses WHERE id=$1 LIMIT 1`
			row := db.QueryRow(sqlStatement, defaultAddressId)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				return
			}
			var company string

			if len(company) > 0 {
				c.FiscalName = company
				c.Name = c.FiscalName + " - " + c.Tradename
			} else {
				c.FiscalName = c.Tradename
				c.Name = c.Tradename
			}
			c.updateCustomer()

			// add/update addresses
			sqlStatement = `SELECT id, address1, address2, city, province, zip, country_code FROM public.sy_addresses WHERE customer_id=$1`
			rows, err := db.Query(sqlStatement, id)
			if err != nil {
				log("DB", rows.Err().Error())
				return
			}

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
				var countryId int16
				sqlStatement = `SELECT id FROM public.country WHERE iso_2=$1 OR iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, countryCode)
				if row.Err() != nil {
					continue
				}
				row.Scan(&countryId)

				// get province
				var provinceId *int32
				sqlStatement = `SELECT id FROM public.state WHERE name=$1 LIMIT 1`
				row = db.QueryRow(sqlStatement, province)
				if row.Err() != nil {
					continue
				}

				var stateId int32
				row.Scan(&stateId)
				if stateId > 0 {
					provinceId = &stateId
				}

				// ¿does the row exist?
				sqlStatement := `SELECT COUNT(*) FROM address WHERE sy_id=$1`
				rowCount := db.QueryRow(sqlStatement, id)
				var rows int32
				rowCount.Scan(&rows)

				if rows == 0 {
					a := Address{}
					a.Customer = &c.Id
					a.Address = address1
					a.Address2 = address2
					a.City = city
					a.ZipCode = zip
					a.Country = countryId
					a.State = provinceId
					a.ShopifyId = id
					if len(company) > 0 {
						a.PrivateOrBusiness = "B"
					} else {
						a.PrivateOrBusiness = "P"
					}
					a.insertAddress()
				} else {
					sqlStatement := `SELECT id FROM address WHERE sy_id=$1`
					rowCount := db.QueryRow(sqlStatement, id)
					var addressId int32
					rowCount.Scan(&addressId)

					a := getAddressRow(addressId)
					a.Customer = &c.Id
					a.Address = address1
					a.Address2 = address2
					a.City = city
					a.ZipCode = zip
					a.Country = countryId
					a.State = provinceId
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
} // copySyCustomers()

func copySyProducts() {
	s := getSettingsRecord()

	sqlStatement := `SELECT id FROM public.sy_products`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var syProductId int64
		rows.Scan(&syProductId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM product WHERE sy_id=$1`
		rowCount := db.QueryRow(sqlStatement, syProductId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, title, body_html FROM public.sy_products WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, syProductId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			return
		}

		var id int64
		var title string
		var bodyHtml string
		row.Scan(&id, &title, &bodyHtml)

		var variants []SYVariantDB = make([]SYVariantDB, 0)
		sqlStatement = `SELECT id, product_id, title, price, sku, option1, option2, option3, taxable, barcode, grams FROM public.sy_variants WHERE product_id=$1`
		rowsVariants, err := db.Query(sqlStatement, id)
		if err != nil {
			log("DB", err.Error())
			return
		}

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
				p.Weight = float32(variants[0].Grams / 1000)
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
				p.insertProduct()
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
					p.Weight = float32(variants[i].Grams / 1000)
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
					p.insertProduct()
				}
			}
		} else { // if rows == 0
			// if the product uses variants, crate a product on the ERP for every single variant, or, if there is only one variant, create a single product on the ERP
			if len(variants) == 1 {
				sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variation_id=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, id, variants[0].Id)
				if row.Err() != nil {
					log("DB", row.Err().Error())
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
				p.Weight = float32(variants[0].Grams / 1000)
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
				p.updateProduct()
			} else {
				for i := 0; i < len(variants); i++ {
					sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variation_id=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, id, variants[i].Id)
					if row.Err() != nil {
						log("DB", row.Err().Error())
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
						p.Weight = float32(variants[i].Grams / 1000)
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
						p.updateProduct()
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
						p.Weight = float32(variants[i].Grams / 1000)
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
						p.insertProduct()
					}
				} // for
			} // else
		} // else
	} // for rows.Next()
} // copySyProducts

func copySyDraftOrders() {
	s := getSettingsRecord()

	sqlStatement := `SELECT id FROM public.sy_draft_orders`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var syDraftOrderId int64
		rows.Scan(&syDraftOrderId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sales_order WHERE sy_draft_id=$1`
		rowCount := db.QueryRow(sqlStatement, syDraftOrderId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, currency, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id FROM public.sy_draft_orders WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, syDraftOrderId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			return
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
		var totalTax float32
		var customerId int64
		row.Scan(&id, &currency, &taxExempt, &name, &shippingAddress1, &shippingAddress2, &shippingAddressCity, &shippingAddressZip, &shippingAddressCountryCode, &billingAddress1, &billingAddress2, &billingAddressCity, &billingAddressZip, &billingAddressCountryCode, &totalTax, &customerId)

		// get the currency
		var currencyId int16
		sqlStatement = `SELECT id FROM public.currency WHERE iso_code=$1 LIMIT 1`
		row = db.QueryRow(sqlStatement, currency)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			continue
		}

		row.Scan(&currencyId)
		if currencyId <= 0 {
			continue
		}

		// get the customer
		var customerIdErp int32
		sqlStatement = `SELECT id FROM customer WHERE sy_id=$1 LIMIT 1`
		row = db.QueryRow(sqlStatement, customerId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			continue
		}

		row.Scan(&customerIdErp)
		if customerIdErp <= 0 {
			continue
		}

		// get the billing address
		var billingAddressId int32
		var billingZone string
		sqlStatement = `SELECT id,(SELECT zone FROM country WHERE country.id=address.country) FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 LIMIT 1`
		row = db.QueryRow(sqlStatement, billingAddress1, billingAddress2, billingAddressCity, billingAddressZip, billingAddressCountryCode)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			continue
		}

		row.Scan(&billingAddressId, &billingZone)
		if billingAddressId <= 0 {
			continue
		}

		// get the shipping address
		var shippingAddressId int32
		sqlStatement = `SELECT id FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 LIMIT 1`
		row = db.QueryRow(sqlStatement, shippingAddress1, shippingAddress2, shippingAddressCity, shippingAddressZip, shippingAddressCountryCode)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			continue
		}

		row.Scan(&shippingAddressId)
		if shippingAddressId <= 0 {
			continue
		}

		if rows == 0 {
			o := SaleOrder{}
			o.BillingAddress = billingAddressId
			o.ShippingAddress = shippingAddressId
			o.Customer = customerIdErp
			o.Reference = name
			o.Currency = currencyId
			o.Warehouse = s.DefaultWarehouse
			o.PaymentMethod = *s.ShopifyDefaultPaymentMethod

			if billingZone == "E" {
				o.BillingSeries = *s.ShopifyExportSerie
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeries = *s.ShopifyIntracommunitySerie
			} else {
				o.BillingSeries = *s.ShopifyInteriorSerie
			}

			o.ShopifyDraftId = id
			ok, orderId := o.insertSalesOrder()
			if !ok {
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeries == nil || *c.BillingSeries == "" {
				c.BillingSeries = &o.BillingSeries
			}
			c.updateCustomer()

			// insert the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_draft_order_line_item WHERE draft_order_id=$1`
			rows, err := db.Query(sqlStatement, id)
			if err != nil {
				log("DB", err.Error())
				return
			}

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float32
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				sqlStatement := `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, productId, variantId)
				if row.Err() != nil {
					log("DB", err.Error())
					return
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					continue
				}

				d := SalesOrderDetail{}
				d.Order = orderId
				d.Quantity = quantity
				d.Price = price
				if taxable && !taxExempt {
					d.VatPercent = s.DefaultVatPercent
				} else {
					d.VatPercent = 0
				}
				d.Product = productIdErp
				d.ShopifyDraftId = id
				d.insertSalesOrderDetail()
			} // for rows.Next()
		} else { // if rows == 0
			var orderIdErp int32
			sqlStatement := `SELECT id FROM sales_order WHERE sy_draft_id=$1 LIMIT 1`
			row := db.QueryRow(sqlStatement, id)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			row.Scan(&orderIdErp)
			if orderIdErp <= 0 {
				continue
			}

			o := getSalesOrderRow(orderIdErp)
			o.BillingAddress = billingAddressId
			o.ShippingAddress = shippingAddressId
			o.Customer = customerIdErp
			o.Reference = name
			o.Currency = currencyId
			o.Warehouse = s.DefaultWarehouse
			o.PaymentMethod = *s.ShopifyDefaultPaymentMethod

			if billingZone == "E" {
				o.BillingSeries = *s.ShopifyExportSerie
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeries = *s.ShopifyIntracommunitySerie
			} else {
				o.BillingSeries = *s.ShopifyInteriorSerie
			}

			ok := o.updateSalesOrder()
			if !ok {
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeries == nil || *c.BillingSeries == "" {
				c.BillingSeries = &o.BillingSeries
			}
			c.updateCustomer()

			// insert/update the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_draft_order_line_item WHERE draft_order_id=$1`
			rows, err := db.Query(sqlStatement, id)
			if err != nil {
				log("DB", err.Error())
				return
			}

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float32
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				var salesOrderDetailId int32
				sqlStatement := `SELECT id FROM sales_order_detail WHERE sy_draft_id=$1`
				row := db.QueryRow(sqlStatement, id)
				if row.Err() != nil {
					log("DB", row.Err().Error())
					return
				}

				row.Scan(&salesOrderDetailId)

				sqlStatement = `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 LIMIT 1`
				row = db.QueryRow(sqlStatement, productId, variantId)
				if row.Err() != nil {
					log("DB", row.Err().Error())
					return
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					continue
				}

				if salesOrderDetailId <= 0 {
					d := SalesOrderDetail{}
					d.Order = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.Product = productIdErp
					d.ShopifyDraftId = id
					d.insertSalesOrderDetail()
				} else { // if salesOrderDetailId <= 0
					d := getSalesOrderDetailRow(salesOrderDetailId)
					d.Order = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.Product = productIdErp
					d.updateSalesOrderDetail()
				}
			} // for rows.Next()
		} // else
	} // for rows.Next()
} // copySyDraftOrders

func copySyOrders() {
	s := getSettingsRecord()

	sqlStatement := `SELECT id FROM public.sy_orders`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var syOrderId int64
		rows.Scan(&syOrderId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM sales_order WHERE sy_id=$1`
		rowCount := db.QueryRow(sqlStatement, syOrderId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, currency, current_total_discounts, total_shipping_price_set_amount, total_shipping_price_set_currency_code, tax_exempt, name, shipping_address_1, shipping_address2, shipping_address_city, shipping_address_zip, shipping_address_country_code, billing_address_1, billing_address2, billing_address_city, billing_address_zip, billing_address_country_code, total_tax, customer_id, gateway FROM public.sy_orders WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, syOrderId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			return
		}

		var id int64
		var currency string
		var currentTotalDiscounts float32
		var totalShippingPriceSetAmount float32
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
		var totalTax float32
		var customerId int64
		var gateway string
		row.Scan(&id, &currency, &currentTotalDiscounts, &totalShippingPriceSetAmount, &totalShippingPriceSetCurrencyCode, &taxExempt, &name, &shippingAddress1, &shippingAddress2, &shippingAddressCity, &shippingAddressZip, &shippingAddressCountryCode, &billingAddress1, &billingAddress2, &billingAddressCity, &billingAddressZip, &billingAddressCountryCode, &totalTax, &customerId, &gateway)

		if rows == 0 {
			// the order has been paid in Shopify, update one last time the order and set invoice it

			// get the currency
			var currencyId int16
			sqlStatement = `SELECT id FROM public.currency WHERE iso_code=$1 LIMIT 1`
			row = db.QueryRow(sqlStatement, currency)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			row.Scan(&currencyId)
			if currencyId <= 0 {
				continue
			}

			// get the customer
			var customerIdErp int32
			sqlStatement = `SELECT id FROM customer WHERE sy_id=$1 LIMIT 1`
			row = db.QueryRow(sqlStatement, customerId)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			row.Scan(&customerIdErp)
			if customerIdErp <= 0 {
				continue
			}

			// get the billing address
			var billingAddressId int32
			var billingZone string
			sqlStatement = `SELECT id,(SELECT zone FROM country WHERE country.id=address.country) FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 LIMIT 1`
			row = db.QueryRow(sqlStatement, billingAddress1, billingAddress2, billingAddressCity, billingAddressZip, billingAddressCountryCode)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			row.Scan(&billingAddressId, &billingZone)
			if billingAddressId <= 0 {
				continue
			}

			// get the shipping address
			var shippingAddressId int32
			sqlStatement = `SELECT id FROM public.address WHERE address=$1 AND address_2=$2 AND city=$3 AND zip_code=$4 AND (SELECT iso_2 FROM country WHERE country.id=address.country)=$5 LIMIT 1`
			row = db.QueryRow(sqlStatement, shippingAddress1, shippingAddress2, shippingAddressCity, shippingAddressZip, shippingAddressCountryCode)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			row.Scan(&shippingAddressId)
			if shippingAddressId <= 0 {
				continue
			}

			// get the order id
			sqlStatement := `SELECT id FROM sy_draft_orders WHERE order_id=$1`
			row := db.QueryRow(sqlStatement, id)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			var draftOrderId int64
			row.Scan(&draftOrderId)
			if draftOrderId <= 0 {
				continue
			}

			sqlStatement = `SELECT id FROM sales_order WHERE sy_draft_id=$1`
			row = db.QueryRow(sqlStatement, draftOrderId)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			var saleOrderIdErp int32
			row.Scan(&saleOrderIdErp)
			if saleOrderIdErp <= 0 {
				continue
			}

			// get the payment method
			var paymentMethod int16
			sqlStatement = `SELECT id FROM payment_method WHERE shopify_module_name=$1`
			row = db.QueryRow(sqlStatement, gateway)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				paymentMethod = *s.ShopifyDefaultPaymentMethod
			} else {
				row.Scan(&paymentMethod)
				if paymentMethod <= 0 {
					paymentMethod = *s.ShopifyDefaultPaymentMethod
				}
			}

			// update the order
			o := getSalesOrderRow(saleOrderIdErp)
			o.BillingAddress = billingAddressId
			o.ShippingAddress = shippingAddressId
			o.Customer = customerIdErp
			o.Reference = name
			o.Currency = currencyId
			o.Warehouse = s.DefaultWarehouse
			o.PaymentMethod = paymentMethod

			if billingZone == "E" {
				o.BillingSeries = *s.ShopifyExportSerie
			} else if billingZone == "U" && totalTax == 0 {
				o.BillingSeries = *s.ShopifyIntracommunitySerie
			} else {
				o.BillingSeries = *s.ShopifyInteriorSerie
			}

			o.ShopifyId = id
			ok := o.updateSalesOrder()
			if !ok {
				continue
			}

			// set the customer details if are empty
			c := getCustomerRow(customerIdErp)
			if c.BillingSeries == nil || *c.BillingSeries == "" {
				c.BillingSeries = &o.BillingSeries
			}
			c.updateCustomer()

			// update the details
			sqlStatement = `SELECT id, variant_id, product_id, quantity, taxable, price FROM public.sy_order_line_item WHERE order_id=$1`
			rows, err := db.Query(sqlStatement, id)
			if err != nil {
				log("DB", err.Error())
				return
			}

			details := getSalesOrderDetail(o.Id)

			for rows.Next() {
				var id int64
				var variantId int64
				var productId int64
				var quantity int32
				var taxable bool
				var price float32
				rows.Scan(&id, &variantId, &productId, &quantity, &taxable, &price)

				sqlStatement = `SELECT id FROM product WHERE sy_id=$1 AND sy_variant_id=$2 LIMIT 1`
				row = db.QueryRow(sqlStatement, productId, variantId)
				if row.Err() != nil {
					log("DB", err.Error())
					return
				}

				var productIdErp int32
				row.Scan(&productIdErp)
				if productIdErp <= 0 {
					continue
				}

				var salesOrderDetailId int32
				for i := 0; i < len(details); i++ {
					if details[i].Product == productIdErp {
						salesOrderDetailId = details[i].Id
						break
					}
				}

				if salesOrderDetailId > 0 {
					d := getSalesOrderDetailRow(salesOrderDetailId)
					d.Order = o.Id
					d.Quantity = quantity
					d.Price = price
					if taxable && !taxExempt {
						d.VatPercent = s.DefaultVatPercent
					} else {
						d.VatPercent = 0
					}
					d.Product = productIdErp
					d.ShopifyId = id
					d.updateSalesOrderDetail()
				}
			}

			// create the invoice for the order
			invoiceAllSaleOrder(o.Id)
		} // if rows == 0
	} // for rows.Next()
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

func updateTrackingNumberShopifyOrder(salesOrderId int32, trackingNumber string) bool {
	fulfillment := SYFulfillmentContainer{}
	fulfillment.Fulfillment = SYFulfillment{}
	fulfillment.Fulfillment.TrackingNumber = trackingNumber
	order := getSalesOrderRow(salesOrderId)
	settings := getSettingsRecord()
	// we need to obtain the location_id for the fulfillment
	if settings.ShopifyShopLocationId <= 0 {
		url := getShopifyAPI_URL("locations")
		jsonSY, err := getShopifyJSON(url)
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
		fulfillment.Fulfillment.LocationId = settings.ShopifyShopLocationId
	}

	// get the name of the carrier
	if order.Carrier == nil {
		return false
	}
	carrier := getCarierRow(*order.Carrier)
	fulfillment.Fulfillment.TrackingCompany = carrier.Name

	// line items
	details := getSalesOrderDetail(salesOrderId)
	for i := 0; i < len(details); i++ {
		fulfillment.Fulfillment.LineItems = append(fulfillment.Fulfillment.LineItems, SYFulfillmentLineItem{
			Id: details[i].ShopifyId,
		})
	}

	// send data
	data, _ := json.Marshal(fulfillment)
	url := getShopifyAPI_URL("orders/" + strconv.Itoa(int(order.ShopifyId)) + "/fulfillments")
	postShopifyJSON(url, data)
	return true
}

type SYDraftOrderComplete struct {
	PaymentPending bool `json:"payment_pending"`
}

func updateStatusPaymentAcceptedShopify(salesOrderId int32) bool {
	complete := SYDraftOrderComplete{PaymentPending: false}
	order := getSalesOrderRow(salesOrderId)
	// send data
	data, _ := json.Marshal(complete)
	url := getShopifyAPI_URL("draft_orders/" + strconv.Itoa(int(order.ShopifyDraftId)) + "/complete")
	putShopifyJSON(url, data)
	return true
}