package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/lib/pq"
)

const WOOCOMMERCE_PROCESSING = "processing"

// =====
// GENERIC FUNCTIONS
// =====

func getWooCommerceAPI_URL(resourceName string, enterpriseId int32) string {
	s := getSettingsRecordById(enterpriseId)

	return s.SettingsEcommerce.WooCommerceUrl + resourceName
}

func getWooCommerceJSON(URL string, enterpriseId int32) ([]byte, error) {
	s := getSettingsRecordById(enterpriseId)

	// OAuth 1.0 request
	config := oauth1.NewConfig(s.SettingsEcommerce.WooCommerceConsumerKey, s.SettingsEcommerce.WooCommerceConsumerSecret)
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	resp, err := httpClient.Get(URL)
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
// DEFINE WOOCOMMERCE CLASSES
// =====

type WcJsonDateTime time.Time

type WCCustomer struct {
	Id               int32          `json:"id" gorm:"primaryKey"`
	DateCreated      WcJsonDateTime `json:"date_created" gorm:"column:date_created;type:timestamp(0) with time zone;not null:true"`
	Email            string         `json:"email" gorm:"column:email;type:character varying(100);not null:true"`
	FirstName        string         `json:"first_name" gorm:"column:first_name;type:character varying(255);not null:true"`
	LastName         string         `json:"last_name" gorm:"column:last_name;type:character varying(255);not null:true"`
	Billing          WCAddress      `json:"billing" gorm:"-"`
	Shipping         WCAddress      `json:"shipping" gorm:"-"`
	BillingAddress1  string         `json:"-" gorm:"column:billing_address_1;type:character varying(255);not null:true"`
	BillingAddress2  string         `json:"-" gorm:"column:billing_address_2;type:character varying(255);not null:true"`
	BillingCity      string         `json:"-" gorm:"column:billing_city;type:character varying(255);not null:true"`
	BillingPostcode  string         `json:"-" gorm:"column:billing_postcode;type:character varying(255);not null:true"`
	BillingCountry   string         `json:"-" gorm:"column:billing_country;type:character varying(255);not null:true"`
	BillingState     string         `json:"-" gorm:"column:billing_state;type:character varying(255);not null:true"`
	BillingPhone     string         `json:"-" gorm:"column:billing_phone;type:character varying(255);not null:true"`
	ShippingAddress1 string         `json:"-" gorm:"column:shipping_address_1;type:character varying(255);not null:true"`
	ShippingAddress2 string         `json:"-" gorm:"column:shipping_address_2;type:character varying(255);not null:true"`
	ShippingCity     string         `json:"-" gorm:"column:shipping_city;type:character varying(255);not null:true"`
	ShippingPostcode string         `json:"-" gorm:"column:shipping_postcode;type:character varying(255);not null:true"`
	ShippingCountry  string         `json:"-" gorm:"column:shipping_country;type:character varying(255);not null:true"`
	ShippingState    string         `json:"-" gorm:"column:shipping_state;type:character varying(255);not null:true"`
	ShippingPhone    string         `json:"-" gorm:"column:shipping_phone;type:character varying(255);not null:true"`
	WcExists         bool           `json:"-" gorm:"column:wc_exists;type:boolean;not null"`
	BillingCompany   string         `json:"-" gorm:"column:billing_company;type:character varying(255);not null:true"`
	EnterpriseId     int32          `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise       Settings       `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (WCCustomer) TableName() string {
	return "wc_customers"
}

type WCAddress struct {
	Company  string `json:"company"`
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	PostCode string `json:"postcode"`
	Country  string `json:"country"`
	State    string `json:"state"`
	Phone    string `json:"phone"`
}

type WCProduct struct {
	Id               int32               `json:"id" gorm:"primaryKey"`
	Name             string              `json:"name" gorm:"column:name;type:character varying(255);not null:true"`
	DateCreated      WcJsonDateTime      `json:"date_created" gorm:"column:date_created;type:timestamp(0) with time zone;not null:true"`
	Description      string              `json:"description" gorm:"column:description;type:text;not null:true"`
	ShortDescription string              `json:"short_description" gorm:"column:short_description;type:character varying(255);not null:true"`
	Sku              string              `json:"sku" gorm:"column:sku;type:character varying(25);not null:true"`
	Price            float64             `json:"price" gorm:"column:price;type:numeric(10,6);not null:true"`
	Weight           string              `json:"weight" gorm:"column:weight;type:character varying(10);not null:true"`
	Dimensions       WCPRoductDimensions `json:"dimensions" gorm:"-"`
	DimensionsLength string              `json:"-" gorm:"column:dimensions_length;type:character varying(10);not null:true"`
	DimensionsWidth  string              `json:"-" gorm:"column:dimensions_width;type:character varying(10);not null:true"`
	DimensionsHeight string              `json:"-" gorm:"column:dimensions_height;type:character varying(10);not null:true"`
	Images           []WCProductImage    `json:"images" gorm:"-"`
	DBIMages         []string            `json:"-" gorm:"column:images;type:character varying(255)[];not null:true"`
	WcExists         bool                `json:"-" gorm:"column:wc_exists;type:boolean;not null"`
	Variations       []int32             `json:"variations" gorm:"column:variations;type:integer[];not null:true"`
	EnterpriseId     int32               `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise       Settings            `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (WCProduct) TableName() string {
	return "wc_products"
}

type WCPRoductDimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type WCProductImage struct {
	Src string `json:"src"`
}

type WCProductVariation struct {
	Id               int32                         `json:"id" gorm:"primaryKey"`
	Sku              string                        `json:"sku" gorm:"column:sku;type:character varying(25);not null:true"`
	Price            float64                       `json:"price" gorm:"column:price;type:numeric(12,6);not null:true"`
	Weight           string                        `json:"weight" gorm:"column:weight;type:character varying(10);not null:true"`
	Dimensions       WCPRoductDimensions           `json:"dimensions" gorm:"-"`
	DimensionsLength string                        `json:"-" gorm:"column:dimensions_length;type:character varying(10);not null:true"`
	DimensionsWidth  string                        `json:"-" gorm:"column:dimensions_width;type:character varying(10);not null:true"`
	DimensionsHeight string                        `json:"-" gorm:"column:dimensions_height;type:character varying(10);not null:true"`
	Attributes       []WCProductVariationAttribute `json:"attributes" gorm:"-"`
	DBAttributes     []string                      `json:"-" gorm:"column:attributes;type:character varying(255)[];not null:true"`
	WcExists         bool                          `json:"-" gorm:"column:wc_exists;type:boolean;not null"`
	EnterpriseId     int32                         `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise       Settings                      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (WCProductVariation) TableName() string {
	return "wc_product_variations"
}

type WCProductVariationAttribute struct {
	Option string `json:"option"`
}

type WCOrder struct {
	Id               int32           `json:"id" gorm:"primaryKey"`
	Status           string          `json:"status" gorm:"column:status;type:character varying(50);not null:true"`
	Currency         string          `json:"currency" gorm:"column:currency;type:character varying(3);not null:true"`
	DateCreated      WcJsonDateTime  `json:"date_created" gorm:"column:date_created;type:timestamp(0) with time zone;not null:true"`
	DiscountTax      string          `json:"discount_tax" gorm:"column:discount_tax;type:numeric(14,6);not null:true"`
	ShippingTotal    string          `json:"shipping_total" gorm:"column:shipping_total;type:numeric(14,6);not null:true"`
	ShippingTax      string          `json:"shipping_tax" gorm:"column:shipping_tax;type:numeric(14,6);not null:true"`
	TotalTax         string          `json:"total_tax" gorm:"column:total_tax;type:numeric(14,6);not null:true"`
	CustomerId       int32           `json:"customer_id" gorm:"column:customer_id;type:integer;not null:true"`
	OrderKey         string          `json:"order_key" gorm:"column:order_key;type:character varying(25);not null:true"`
	Billing          WCAddress       `json:"billing" gorm:"-"`
	Shipping         WCAddress       `json:"shipping" gorm:"-"`
	BillingAddress2  string          `json:"-" gorm:"column:billing_address_2;type:character varying(255);not null:true"`
	BillingCity      string          `json:"-" gorm:"column:billing_city;type:character varying(255);not null:true"`
	BillingPostcode  string          `json:"-" gorm:"column:billing_postcode;type:character varying(255);not null:true"`
	BillingCountry   string          `json:"-" gorm:"column:billing_country;type:character varying(255);not null:true"`
	BillingState     string          `json:"-" gorm:"column:billing_state;type:character varying(255);not null:true"`
	BillingPhone     string          `json:"-" gorm:"column:billing_phone;type:character varying(255);not null:true"`
	ShippingAddress1 string          `json:"-" gorm:"column:shipping_address_1;type:character varying(255);not null:true"`
	ShippingAddress2 string          `json:"-" gorm:"column:shipping_address_2;type:character varying(255);not null:true"`
	ShippingCity     string          `json:"-" gorm:"column:shipping_city;type:character varying(255);not null:true"`
	ShippingPostcode string          `json:"-" gorm:"column:shipping_postcode;type:character varying(255);not null:true"`
	ShippingCountry  string          `json:"-" gorm:"column:shipping_country;type:character varying(255);not null:true"`
	ShippingState    string          `json:"-" gorm:"column:shipping_state;type:character varying(255);not null:true"`
	ShippingPhone    string          `json:"-" gorm:"column:shipping_phone;type:character varying(255);not null:true"`
	PaymentMethod    string          `json:"payment_method" gorm:"column:payment_method;type:character varying(50);not null:true"`
	WcExists         bool            `json:"-" gorm:"column:wc_exists;type:boolean;not null"`
	BillingCompany   string          `json:"-" gorm:"column:billing_company;type:character varying(255);not null:true"`
	LineItems        []WCOrderDetail `json:"line_items" gorm:"-"`
	EnterpriseId     int32           `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise       Settings        `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (WCOrder) TableName() string {
	return "wc_orders"
}

type WCOrderDetail struct {
	Id           int32    `json:"id" gorm:"primaryKey"`
	Order        int32    `json:"order" gorm:"column:order;type:integer;not null:true"`
	ProductId    int32    `json:"product_id" gorm:"column:product_id;type:integer;not null:true"`
	VariationId  int32    `json:"variation_id" gorm:"column:variation_id;type:integer;not null:true"`
	Quantity     int32    `json:"quantity" gorm:"column:quantity;type:integer;not null:true"`
	TotalTax     string   `json:"total_tax" gorm:"column:total_tax;type:numeric(14,6);not null:true"`
	Price        float64  `json:"price" gorm:"column:price;type:numeric(14,6);not null:true"`
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (WCOrderDetail) TableName() string {
	return "wc_order_details"
}

func (j *WcJsonDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*j = WcJsonDateTime(t)
	return nil
}

func (j *WcJsonDateTime) ToTime() time.Time {
	return time.Time(*j)
}

// main import function
func importFromWooCommerce(enterpriseId int32) {
	s := getSettingsRecordById(enterpriseId)
	if s.SettingsEcommerce.Ecommerce != "W" {
		return
	}

	// get all data from WooCommerce, write it in tables like the ones that WooCommerce uses
	importWcCustomers(enterpriseId)
	importWcProducts(enterpriseId)
	importWcOrders(enterpriseId)

	// trasnfer the data form the WooCommerce tables to the ERP
	copyWcCustomers(enterpriseId)
	copyWcProducts(enterpriseId)
	copyWcOrders(enterpriseId)
}

// =====
// COPY THE DATA FROM WOOCOMMERCE TO THE WC MARKETNET TABLES
// =====

func importWcCustomers(enterpriseId int32) bool {
	url := getWooCommerceAPI_URL("customers", enterpriseId)
	jsonWC, err := getWooCommerceJSON(url, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	var customers []WCCustomer
	json.Unmarshal(jsonWC, &customers)

	sqlStatement := `UPDATE public.wc_customers SET wc_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(customers); i++ {
		customer := customers[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_customers WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, customer.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.wc_customers(id, date_created, email, first_name, last_name, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, billing_company, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`
			db.Exec(sqlStatement, customer.Id, customer.DateCreated.ToTime(), customer.Email, customer.FirstName, customer.LastName, customer.Billing.Address1, customer.Billing.Address2, customer.Billing.City, customer.Billing.PostCode, customer.Billing.Country, customer.Billing.State, customer.Billing.Phone, customer.Shipping.Address1, customer.Shipping.Address2, customer.Shipping.City, customer.Shipping.PostCode, customer.Shipping.Country, customer.Shipping.State, customer.Shipping.Phone, customer.Billing.Company, enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.wc_customers SET date_created=$2, email=$3, first_name=$4, last_name=$5, billing_address_1=$6, billing_address_2=$7, billing_city=$8, billing_postcode=$9, billing_country=$10, billing_state=$11, billing_phone=$12, shipping_address_1=$13, shipping_address_2=$14, shipping_city=$15, shipping_postcode=$16, shipping_country=$17, shipping_state=$18, shipping_phone=$19, billing_company=$20, wc_exists=true WHERE id=$1 AND enterprise=$21`
			db.Exec(sqlStatement, customer.Id, customer.DateCreated.ToTime(), customer.Email, customer.FirstName, customer.LastName, customer.Billing.Address1, customer.Billing.Address2, customer.Billing.City, customer.Billing.PostCode, customer.Billing.Country, customer.Billing.State, customer.Billing.Phone, customer.Shipping.Address1, customer.Shipping.Address2, customer.Shipping.City, customer.Shipping.PostCode, customer.Shipping.Country, customer.Shipping.State, customer.Shipping.Phone, customer.Billing.Company, enterpriseId)
		}
	}

	sqlStatement = `DELETE FROM public.wc_customers WHERE wc_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	return true
}

func importWcProducts(enterpriseId int32) bool {
	url := getWooCommerceAPI_URL("products", enterpriseId)
	jsonWC, err := getWooCommerceJSON(url, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Products</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	var errors []string = make([]string, 0)

	var products []WCProduct
	json.Unmarshal(jsonWC, &products)

	sqlStatement := `UPDATE public.wc_products SET wc_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `UPDATE public.wc_product_variations SET wc_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(products); i++ {
		product := products[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_products WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, product.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		var images []string = make([]string, 0)
		for j := 0; j < len(product.Images); j++ {
			images = append(images, product.Images[j].Src)
		}

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.wc_products(id, name, date_created, description, short_description, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, images, variations, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
			db.Exec(sqlStatement, product.Id, product.Name, product.DateCreated.ToTime(), product.Description, product.ShortDescription, product.Sku, product.Price, product.Weight, product.Dimensions.Length, product.Dimensions.Width, product.Dimensions.Height, pq.Array(images), pq.Array(product.Variations), enterpriseId)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.wc_products SET name=$2, date_created=$3, description=$4, short_description=$5, sku=$6, price=$7, weight=$8, dimensions_length=$9, dimensions_width=$10, dimensions_height=$11, images=$12, variations=$13, wc_exists=true WHERE id=$1 AND enterprise=$14`
			db.Exec(sqlStatement, product.Id, product.Name, product.DateCreated.ToTime(), product.Description, product.ShortDescription, product.Sku, product.Price, product.Weight, product.Dimensions.Length, product.Dimensions.Width, product.Dimensions.Height, pq.Array(images), pq.Array(product.Variations), enterpriseId)
		}

		// get the variations
		url := getWooCommerceAPI_URL("products", enterpriseId) + "/" + strconv.Itoa(int(product.Id)) + "/variations/"
		jsonWC, err := getWooCommerceJSON(url, enterpriseId)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		var variations []WCProductVariation
		json.Unmarshal(jsonWC, &variations)

		// insert/update variations
		for j := 0; j < len(variations); j++ {
			variation := variations[j]
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM wc_product_variations WHERE id=$1 AND enterprise=$2`
			row := db.QueryRow(sqlStatement, variation.Id, enterpriseId)
			var rows int32
			row.Scan(&rows)

			var attributes []string = make([]string, 0)
			for k := 0; k < len(variation.Attributes); k++ {
				attributes = append(attributes, variation.Attributes[k].Option)
			}

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.wc_product_variations(id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
				db.Exec(sqlStatement, variation.Id, variation.Sku, variation.Price, variation.Weight, variation.Dimensions.Length, variation.Dimensions.Width, variation.Dimensions.Height, pq.Array(attributes), enterpriseId)
			} else { // the row exists, update
				sqlStatement := `UPDATE public.wc_product_variations SET sku=$2, price=$3, weight=$4, dimensions_length=$5, dimensions_width=$6, dimensions_height=$7, attributes=$8, wc_exists=true WHERE id=$1 AND enterprise=$9`
				db.Exec(sqlStatement, variation.Id, variation.Sku, variation.Price, variation.Weight, variation.Dimensions.Length, variation.Dimensions.Width, variation.Dimensions.Height, pq.Array(attributes), enterpriseId)
			}
		}

	}

	sqlStatement = `DELETE FROM public.wc_products WHERE wc_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.wc_product_variations WHERE wc_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCOmmerce import error", "<p>Error at: Products</p>"+errorHtml, enterpriseId)
	}

	return true
}

func importWcOrders(enterpriseId int32) bool {
	url := getWooCommerceAPI_URL("orders", enterpriseId)
	jsonWC, err := getWooCommerceJSON(url, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	var errors []string = make([]string, 0)

	var orders []WCOrder
	json.Unmarshal(jsonWC, &orders)

	sqlStatement := `UPDATE public.wc_orders SET wc_exists=false WHERE enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	for i := 0; i < len(orders); i++ {
		order := orders[i]
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_orders WHERE id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, order.Id, enterpriseId)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			f_discount_tax, err := strconv.ParseFloat(order.DiscountTax, 32)
			if err != nil {
				errors = append(errors, "OrderId "+strconv.Itoa(int(order.Id))+" Error "+err.Error())
				continue
			}
			f_shipping_total, err := strconv.ParseFloat(order.ShippingTotal, 32)
			if err != nil {
				errors = append(errors, "OrderId "+strconv.Itoa(int(order.Id))+" Error "+err.Error())
				continue
			}
			f_shipping_tax, err := strconv.ParseFloat(order.ShippingTax, 32)
			if err != nil {
				errors = append(errors, "OrderId "+strconv.Itoa(int(order.Id))+" Error "+err.Error())
				continue
			}
			f_total_tax, err := strconv.ParseFloat(order.TotalTax, 32)
			if err != nil {
				errors = append(errors, "OrderId "+strconv.Itoa(int(order.Id))+" Error "+err.Error())
				continue
			}

			sqlStatement = `INSERT INTO public.wc_orders(id, status, currency, date_created, discount_tax, shipping_total, shipping_tax, total_tax, customer_id, order_key, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, payment_method, billing_company, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`
			db.Exec(sqlStatement, order.Id, order.Status, order.Currency, order.DateCreated.ToTime(), f_discount_tax, f_shipping_total, f_shipping_tax, f_total_tax, order.CustomerId, order.OrderKey, order.Billing.Address1, order.Billing.Address2, order.Billing.City, order.Billing.PostCode, order.Billing.Country, order.Billing.State, order.Billing.Phone, order.Shipping.Address1, order.Shipping.Address2, order.Shipping.City, order.Shipping.PostCode, order.Shipping.Country, order.Shipping.State, order.Shipping.Phone, order.PaymentMethod, order.Billing.Company, enterpriseId)

			// add order details
			for j := 0; j < len(order.LineItems); j++ {
				lineItem := order.LineItems[j]
				f_total_tax, err := strconv.ParseFloat(lineItem.TotalTax, 32)
				if err != nil {
					errors = append(errors, "OrderId "+strconv.Itoa(int(order.Id))+" Error "+err.Error())
					continue
				}

				sqlStatement := `INSERT INTO public.wc_order_details(id, "order", product_id, variation_id, quantity, total_tax, price, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
				db.Exec(sqlStatement, lineItem.Id, order.Id, lineItem.ProductId, lineItem.VariationId, lineItem.Quantity, f_total_tax, lineItem.Price, enterpriseId)
			}

		} else { // if rows == 0
			sqlStatement := `UPDATE public.wc_orders SET wc_exists=true WHERE id=$1 AND enterprise=$2`
			db.Exec(sqlStatement, order.Id, enterpriseId)
		}
	} // for

	sqlStatement = `DELETE FROM wc_order_details WHERE NOT (SELECT wc_exists FROM wc_orders WHERE wc_orders.id=wc_order_details."order") AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)
	sqlStatement = `DELETE FROM public.wc_orders WHERE wc_exists=false AND enterprise=$1`
	db.Exec(sqlStatement, enterpriseId)

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Orders</p>"+errorHtml, enterpriseId)
	}

	return true
}

// =====
// TRANSFER THE DATA TO THE ERP TABLES
// =====

func copyWcCustomers(enterpriseId int32) bool {
	sqlStatement := `SELECT id FROM public.wc_customers WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Customers</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var wcCustomerId int32
		rows.Scan(&wcCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM customer WHERE wc_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, wcCustomerId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, date_created, email, first_name, last_name, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, billing_company FROM public.wc_customers WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, wcCustomerId, enterpriseId)
		if row.Err() != nil {
			log("WooCommerce", row.Err().Error())
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int32
		var dateCreated time.Time
		var email string
		var firstName string
		var lastName string
		var billingAddress1 string
		var billingAddress2 string
		var billingCity string
		var billingPostcode string
		var billingCountry string
		var billingState string
		var billingPhone string
		var shippingAddress1 string
		var shippingAddress2 string
		var shippingCity string
		var shippingPostcode string
		var shippingCountry string
		var shippingState string
		var shippingPhone string
		var billingCompany string
		row.Scan(&id, &dateCreated, &email, &firstName, &lastName, &billingAddress1, &billingAddress2, &billingCity, &billingPostcode, &billingCountry, &billingState, &billingPhone, &shippingAddress1, &shippingAddress2, &shippingCity, &shippingPostcode, &shippingCountry, &shippingState, &shippingPhone, &billingCompany)

		if rows == 0 {
			// create customer
			c := Customer{}
			c.Tradename = firstName + " " + lastName
			if len(billingCompany) == 0 {
				c.FiscalName = firstName + " " + lastName
				c.Name = firstName + " " + lastName
			} else {
				c.FiscalName = billingCompany
				c.Name = billingCompany + " - " + firstName + " " + lastName
			}
			c.Email = email
			c.DateCreated = dateCreated
			c.WooCommerceId = id
			if len(billingPhone) > 0 {
				c.Phone = billingPhone
			} else if len(shippingPhone) > 0 {
				c.Phone = shippingPhone
			}
			c.EnterpriseId = enterpriseId
			res := c.insertCustomer(0)
			ok, customerId := res.Id > 0, int32(res.Id)
			if !ok {
				errors = append(errors, "Can't insert customer into MARKETNET. "+c.Tradename+" WooComemrce ID "+strconv.Itoa(int(id)))
				continue
			}

			// create billing address
			ba := Address{}
			ba.CustomerId = &customerId
			ba.Address = billingAddress1
			ba.Address2 = billingAddress2
			ba.City = billingCity
			ba.ZipCode = billingPostcode
			// search for the country by iso code
			if len(billingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					ba.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else if len(billingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					ba.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(billingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingState, enterpriseId)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						ba.StateId = &stateid
					}
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			}
			if len(billingCompany) == 0 {
				ba.PrivateOrBusiness = "P"
			} else {
				ba.PrivateOrBusiness = "B"
			}
			ba.EnterpriseId = enterpriseId
			ba.insertAddress(0)

			// create shipping address
			sa := Address{}
			sa.CustomerId = &customerId
			sa.Address = shippingAddress1
			sa.Address2 = shippingAddress2
			sa.City = shippingCity
			sa.ZipCode = shippingPostcode
			// search for the country by iso code
			if len(shippingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					sa.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else if len(shippingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					sa.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(shippingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingState, enterpriseId)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						sa.StateId = &stateid
					}
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			}
			if len(billingCompany) == 0 {
				sa.PrivateOrBusiness = "P"
			} else {
				sa.PrivateOrBusiness = "B"
			}
			sa.EnterpriseId = enterpriseId
			sa.insertAddress(0)

			// set as the main billing and shipping address
			customer := getCustomerRow(customerId)
			customer.MainBillingAddressId = &ba.Id
			customer.MainShippingAddressId = &sa.Id
			customer.updateCustomer(0)
		} else {
			// update the addresses
			sqlStatement := `SELECT id FROM public.customer WHERE wc_id=$1 AND enterprise=$2 LIMIT 1`
			row := db.QueryRow(sqlStatement, wcCustomerId, enterpriseId)
			if row.Err() != nil {
				log("WooCommerce", row.Err().Error())
				errors = append(errors, row.Err().Error())
				continue
			}

			var customerIdErp int32
			row.Scan(&customerIdErp)

			addresses := getCustomerAddresses(customerIdErp, enterpriseId)

			// compare both the billing and shipping address from woocommerce with the ones that are in the erp.
			// if the addresses are different, create new ones

			// billing address
			var billingAddressFound bool = false
			for i := 0; i < len(addresses); i++ {
				// search for the country by iso code
				var countryId int32
				var stateId *int32
				if len(billingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("WooComemrce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else if len(billingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(billingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingState, enterpriseId)
					if row.Err() == nil {
						var id int32
						row.Scan(&id)
						if id > 0 {
							stateId = &id
						}
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				}

				// are the addresses the same?
				if addresses[i].Address == billingAddress1 && addresses[i].Address2 == billingAddress2 && addresses[i].City == billingCity && addresses[i].ZipCode == billingPostcode && addresses[i].CountryId == countryId && addresses[i].StateId == stateId {
					billingAddressFound = true
					break
				} // if
			} // for
			if !billingAddressFound {
				// create billing address
				ba := Address{}
				ba.CustomerId = &customerIdErp
				ba.Address = billingAddress1
				ba.Address2 = billingAddress2
				ba.City = billingCity
				ba.ZipCode = billingPostcode
				if len(billingCompany) > 0 {
					ba.PrivateOrBusiness = "B"
				} else {
					ba.PrivateOrBusiness = "P"
				}
				// search for the country by iso code
				if len(billingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
					if row.Err() == nil {
						var countryId int32
						row.Scan(&countryId)
						ba.CountryId = countryId
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else if len(billingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
					if row.Err() == nil {
						var countryId int32
						row.Scan(&countryId)
						ba.CountryId = countryId
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(billingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingState, enterpriseId)
					if row.Err() == nil {
						var stateid int32
						row.Scan(&stateid)
						if stateid > 0 {
							ba.StateId = &stateid
						}
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				}
				ba.EnterpriseId = enterpriseId
				ba.insertAddress(0)
			}

			// shipping address
			var shippingAddressFound bool = false
			for i := 0; i < len(addresses); i++ {
				// search for the country by iso code
				var countryId int32
				var stateId *int32
				if len(shippingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else if len(shippingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(shippingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingState, enterpriseId)
					if row.Err() == nil {
						var id int32
						row.Scan(&id)
						if id > 0 {
							stateId = &id
						}
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				}

				// are the addresses the same?
				if addresses[i].Address == shippingAddress1 && addresses[i].Address2 == shippingAddress2 && addresses[i].City == shippingCity && addresses[i].ZipCode == shippingPostcode && addresses[i].CountryId == countryId && addresses[i].StateId == stateId {
					shippingAddressFound = true
					break
				} // if
			} // for
			if !shippingAddressFound {
				// create shipping address
				sa := Address{}
				sa.CustomerId = &customerIdErp
				sa.Address = shippingAddress1
				sa.Address2 = shippingAddress2
				sa.City = shippingCity
				sa.ZipCode = shippingPostcode
				if len(billingCompany) > 0 {
					sa.PrivateOrBusiness = "B"
				} else {
					sa.PrivateOrBusiness = "P"
				}
				// search for the country by iso code
				if len(shippingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
					if row.Err() == nil {
						var countryId int32
						row.Scan(&countryId)
						sa.CountryId = countryId
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else if len(shippingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
					if row.Err() == nil {
						var countryId int32
						row.Scan(&countryId)
						sa.CountryId = countryId
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(shippingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingState, enterpriseId)
					if row.Err() == nil {
						var stateid int32
						row.Scan(&stateid)
						if stateid > 0 {
							sa.StateId = &stateid
						}
					} else {
						log("WooCommerce", row.Err().Error())
						errors = append(errors, row.Err().Error())
					}
				}
				sa.EnterpriseId = enterpriseId
				sa.insertAddress(0)
			}

		} // else
	} // for rows.Next()

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Customers</p>"+errorHtml, enterpriseId)
	}

	return true
} // copyWcCustomers

func copyWcProducts(enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id FROM public.wc_products AND enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Products</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var wcCustomerId int32
		rows.Scan(&wcCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM product WHERE wc_id=$1 AND enterprise=$2`
		rowCount := db.QueryRow(sqlStatement, wcCustomerId, enterpriseId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, name, date_created, description, short_description, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, images, variations FROM public.wc_products WHERE id=$1 AND enterprise=$2 LIMIT 1`
		row := db.QueryRow(sqlStatement, wcCustomerId, enterpriseId)
		if row.Err() != nil {
			errors = append(errors, row.Err().Error())
			continue
		}

		var id int32
		var name string
		var dateCreated time.Time
		var description string
		var shortDescription string
		var sku string
		var price float64
		var weight string
		var dimensionsLength string
		var dimensionsWidth string
		var dimensionsHeight string
		var images []string
		var variations []int32
		row.Scan(&id, &name, &dateCreated, &description, &shortDescription, &sku, &price, &weight, &dimensionsLength, &dimensionsWidth, &dimensionsHeight, pq.Array(&images), pq.Array(&variations))

		if rows == 0 {
			// create the product
			// is it a simple product or a product with variations (red, blue, ...)?
			if len(variations) == 0 {
				p := Product{}
				p.Name = name
				p.DateCreated = dateCreated
				p.Description += shortDescription
				if len(p.Description) > 0 {
					p.Description += "\n\n"
				}
				p.Description = description
				p.BarCode = sku
				p.Price = price
				f_weight, err := strconv.ParseFloat(weight, 32)
				if err == nil {
					p.Weight = float64(f_weight)
				}
				f_length, err := strconv.ParseFloat(dimensionsLength, 32)
				if err == nil {
					p.Depth = float64(f_length)
				}
				f_width, err := strconv.ParseFloat(dimensionsWidth, 32)
				if err == nil {
					p.Width = float64(f_width)
				}
				f_height, err := strconv.ParseFloat(dimensionsHeight, 32)
				if err == nil {
					p.Height = float64(f_height)
				}
				p.VatPercent = s.DefaultVatPercent
				p.WooCommerceId = id
				p.EnterpriseId = enterpriseId
				result := p.insertProduct(0)
				if !result.Ok {
					errors = append(errors, "Error inserting a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}

				for i := 0; i < len(images); i++ {
					pi := ProductImage{
						ProductId: p.Id,
						URL:       images[i],
					}
					pi.insertProductImage(enterpriseId)
				}
			} else {
				// the product has variations. create a new product by each variation (combination of options)
				for i := 0; i < len(variations); i++ {
					sqlStatement := `SELECT id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes FROM public.wc_product_variations WHERE id=$1 AND enterprise=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, variations[i], enterpriseId)
					if row.Err() != nil {
						errors = append(errors, row.Err().Error())
						continue
					}

					var v_id int32
					var v_sku string
					var v_price float64
					var v_weight string
					var v_dimensionsLength string
					var v_dimensionsWidth string
					var v_dimensionsHeight string
					var attributes []string
					row.Scan(&v_id, &v_sku, &v_price, &v_weight, &v_dimensionsLength, &v_dimensionsWidth, &v_dimensionsHeight, pq.Array(&attributes))

					p := Product{}
					p.Name = name
					for j := 0; j < len(attributes); j++ {
						p.Name += " " + attributes[j]
					}
					p.DateCreated = dateCreated
					p.Description += shortDescription
					if len(p.Description) > 0 {
						p.Description += "\n\n"
					}
					p.Description = description
					p.BarCode = v_sku
					p.Price = v_price
					f_weight, err := strconv.ParseFloat(v_weight, 32)
					if err == nil {
						p.Weight = float64(f_weight)
					}
					f_length, err := strconv.ParseFloat(v_dimensionsLength, 32)
					if err == nil {
						p.Depth = float64(f_length)
					}
					f_width, err := strconv.ParseFloat(v_dimensionsWidth, 32)
					if err == nil {
						p.Width = float64(f_width)
					}
					f_height, err := strconv.ParseFloat(v_dimensionsHeight, 32)
					if err == nil {
						p.Height = float64(f_height)
					}
					p.VatPercent = s.DefaultVatPercent
					p.WooCommerceId = id
					p.WooCommerceVariationId = v_id
					p.EnterpriseId = enterpriseId
					result := p.insertProduct(0)
					if !result.Ok {
						errors = append(errors, "Error inserting a product with combinations into MARKETNET. Product name "+
							p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
					}

					for i := 0; i < len(images); i++ {
						pi := ProductImage{
							ProductId: p.Id,
							URL:       images[i],
						}
						pi.insertProductImage(enterpriseId)
					}
				}
			} // else
		} else { // if rows == 0 {

			if len(variations) == 0 {
				sqlStatement := `SELECT id FROM product WHERE wc_id=$1 AND wc_variation_id=$2 AND enterprise=$3 LIMIT 1`
				row := db.QueryRow(sqlStatement, id, 0, enterpriseId)
				if row.Err() != nil {
					errors = append(errors, row.Err().Error())
					continue
				}

				var productId int32
				row.Scan(&productId)
				p := getProductRow(productId)
				p.Name = name
				p.Description += shortDescription
				if len(p.Description) > 0 {
					p.Description += "\n\n"
				}
				p.Description = description
				p.BarCode = sku
				p.Price = price
				f_weight, err := strconv.ParseFloat(weight, 32)
				if err == nil {
					p.Weight = float64(f_weight)
				}
				f_length, err := strconv.ParseFloat(dimensionsLength, 32)
				if err == nil {
					p.Depth = float64(f_length)
				}
				f_width, err := strconv.ParseFloat(dimensionsWidth, 32)
				if err == nil {
					p.Width = float64(f_width)
				}
				f_height, err := strconv.ParseFloat(dimensionsHeight, 32)
				if err == nil {
					p.Height = float64(f_height)
				}
				result := p.updateProduct(0)
				if !result.Ok {
					errors = append(errors, "Error updating a simple product into MARKETNET. Product name "+
						p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
				}
			} else { // if len(variations) == 0 {
				for i := 0; i < len(variations); i++ {
					sqlStatement := `SELECT id FROM product WHERE wc_id=$1 AND wc_variation_id=$2 AND enterprise=$3 LIMIT 1`
					row := db.QueryRow(sqlStatement, id, variations[i], enterpriseId)
					if row.Err() != nil {
						errors = append(errors, row.Err().Error())
						continue
					}

					var productId int32
					row.Scan(&productId)
					p := getProductRow(productId)

					sqlStatement = `SELECT id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes FROM public.wc_product_variations WHERE id=$1 AND enterprise=$2 LIMIT 1`
					row = db.QueryRow(sqlStatement, variations[i], enterpriseId)
					if row.Err() != nil {
						errors = append(errors, row.Err().Error())
						continue
					}

					var v_id int32
					var v_sku string
					var v_price float64
					var v_weight string
					var v_dimensionsLength string
					var v_dimensionsWidth string
					var v_dimensionsHeight string
					var attributes []string
					row.Scan(&v_id, &v_sku, &v_price, &v_weight, &v_dimensionsLength, &v_dimensionsWidth, &v_dimensionsHeight, pq.Array(&attributes))

					p.Name = name
					p.Description += shortDescription
					if len(p.Description) > 0 {
						p.Description += "\n\n"
					}
					p.Description = description
					p.BarCode = v_sku
					p.Price = v_price
					f_weight, err := strconv.ParseFloat(v_weight, 32)
					if err == nil {
						p.Weight = float64(f_weight)
					}
					f_length, err := strconv.ParseFloat(v_dimensionsLength, 32)
					if err == nil {
						p.Depth = float64(f_length)
					}
					f_width, err := strconv.ParseFloat(v_dimensionsWidth, 32)
					if err == nil {
						p.Width = float64(f_width)
					}
					f_height, err := strconv.ParseFloat(v_dimensionsHeight, 32)
					if err == nil {
						p.Height = float64(f_height)
					}
					result := p.updateProduct(0)
					if !result.Ok {
						errors = append(errors, "Error updating a product with combinations into MARKETNET. Product name "+
							p.Name+" product reference "+p.Reference+" error code "+strconv.Itoa(int(result.ErrorCode))+" extra data: "+stringArrayToString(result.ExtraData))
					}
				} // for
			} // else
		} // else
	} // for

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Products</p>"+errorHtml, enterpriseId)
	}

	return true
} // copyWcProducts

func copyWcOrders(enterpriseId int32) bool {
	settings := getSettingsRecordById(enterpriseId)

	sqlStatement := `SELECT id, status, currency, date_created, discount_tax, shipping_total, shipping_tax, total_tax, customer_id, order_key, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, payment_method, billing_company FROM public.wc_orders WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("WooCommerce", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.SettingsEmail.EmailSendErrorEcommerce) > 0 {
			sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Orders</p><p>Error data: "+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	defer rows.Close()
	var errors []string = make([]string, 0)

	for rows.Next() {
		var id int32
		var status string
		var currency string
		var dateCreated time.Time
		var discountTax float64
		var shippingTotal float64
		var shippingTax float64
		var totalTax float64
		var customerId int32
		var orderKey string
		var billingAddress1 string
		var billingAddress2 string
		var billingCity string
		var billingPostcode string
		var billingCountry string
		var billingState string
		var shippingAddress1 string
		var shippingAddress2 string
		var shippingCity string
		var shippingPostcode string
		var shippingCountry string
		var shippingState string
		var paymentMethod string
		var billingCompany string
		rows.Scan(&id, &status, &currency, &dateCreated, &discountTax, &shippingTotal, &shippingTax, &totalTax, &customerId, &orderKey, &billingAddress1, &billingAddress2, &billingCity, &billingPostcode, &billingCountry, &billingState, &shippingAddress1, &shippingAddress2, &shippingCity, &shippingPostcode, &shippingCountry, &shippingState, &paymentMethod, &billingCompany)

		// does the order exist?
		sqlStatement := `SELECT COUNT(id) FROM sales_order WHERE wc_id=$1 AND enterprise=$2`
		row := db.QueryRow(sqlStatement, id, enterpriseId)

		var orders int32
		row.Scan(&orders)

		if orders > 0 { // don't continue if the order exists
			continue
		}

		// get the customer
		sqlStatement = `SELECT id FROM customer WHERE wc_id=$1 AND enterprise=$2 LIMIT 1`
		row = db.QueryRow(sqlStatement, customerId, enterpriseId)

		var customer int32
		row.Scan(&customer)

		if customer == 0 { // don't continue if the customer doesn't exists
			errors = append(errors, "Can't import order. The customer doesn't exists. Order id "+strconv.Itoa(int(id)))
			continue
		}

		// get the payment method
		sqlStatement = `SELECT id,paid_in_advance FROM payment_method WHERE woocommerce_module_name=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, paymentMethod, enterpriseId)

		var erpPaymentMethod int32
		var paidInAdvance bool
		row.Scan(&erpPaymentMethod, &paidInAdvance)

		if erpPaymentMethod <= 0 { // attempt the default one in the settings (no payment method = likely a manual order)
			if settings.SettingsEcommerce.WooCommerceDefaultPaymentMethodId == nil { // don't continue if the payment method doesn't exists
				errors = append(errors, "Can't import order. The payment method doesn't exists. Order id "+strconv.Itoa(int(id)))
				continue
			}
			erpPaymentMethod = *settings.SettingsEcommerce.WooCommerceDefaultPaymentMethodId
			paidInAdvance = getPaymentMethodRow(erpPaymentMethod).PaidInAdvance
		}

		// get the currency
		sqlStatement = `SELECT id FROM currency WHERE iso_code=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, currency, enterpriseId)

		var erpCurrency int32
		row.Scan(&erpCurrency)

		if erpCurrency == 0 { // don't continue if the currency doesn't exists
			errors = append(errors, "Can't import order. The currency doesn't exists. Order id "+strconv.Itoa(int(id)))
			continue
		}

		// get the billing address
		var billingAddress int32
		var billingZone string
		customerAddresses := getCustomerAddresses(customer, enterpriseId)
		for i := 0; i < len(customerAddresses); i++ {
			country := getCountryRow(customerAddresses[i].CountryId, enterpriseId)
			if customerAddresses[i].Address == billingAddress1 && customerAddresses[i].Address2 == billingAddress2 && customerAddresses[i].City == billingCity && customerAddresses[i].ZipCode == billingPostcode && (country.Iso2 == billingCountry || country.Iso3 == billingCountry) {
				billingAddress = customerAddresses[i].Id
				billingZone = country.Zone
				break
			}
		}
		if billingAddress == 0 {
			// create billing address
			ba := Address{}
			ba.CustomerId = &customer
			ba.Address = billingAddress1
			ba.Address2 = billingAddress2
			ba.City = billingCity
			ba.ZipCode = billingPostcode
			if len(billingCompany) > 0 {
				ba.PrivateOrBusiness = "B"
			} else {
				ba.PrivateOrBusiness = "P"
			}
			// search for the country by iso code
			if len(billingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					ba.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else if len(billingCountry) == 3 {
				sqlStatement := `SELECT id, zone FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId, &billingZone)
					ba.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(billingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingState, enterpriseId)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						ba.StateId = &stateid
					}
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			}
			ba.EnterpriseId = enterpriseId
			ba.insertAddress(0)
			billingAddress = ba.Id
		}

		// get the shipping address
		var shippingAddress int32
		for i := 0; i < len(customerAddresses); i++ {
			country := getCountryRow(customerAddresses[i].CountryId, enterpriseId)
			if customerAddresses[i].Address == shippingAddress1 && customerAddresses[i].Address2 == shippingAddress2 && customerAddresses[i].City == shippingCity && customerAddresses[i].ZipCode == shippingPostcode && (country.Iso2 == shippingCountry || country.Iso3 == shippingCountry) {
				shippingAddress = customerAddresses[i].Id
				break
			}
		}
		if shippingAddress == 0 {
			// create shipping address
			sa := Address{}
			sa.CustomerId = &customer
			sa.Address = shippingAddress1
			sa.Address2 = shippingAddress2
			sa.City = shippingCity
			sa.ZipCode = shippingPostcode
			if len(billingCompany) > 0 {
				sa.PrivateOrBusiness = "B"
			} else {
				sa.PrivateOrBusiness = "P"
			}
			// search for the country by iso code
			if len(shippingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					sa.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else if len(shippingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry, enterpriseId)
				if row.Err() == nil {
					var countryId int32
					row.Scan(&countryId)
					sa.CountryId = countryId
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(shippingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 AND enterprise=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingState, enterpriseId)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						sa.StateId = &stateid
					}
				} else {
					log("WooCommerce", row.Err().Error())
					errors = append(errors, row.Err().Error())
				}
			}
			sa.EnterpriseId = enterpriseId
			sa.insertAddress(0)
			shippingAddress = sa.Id
		}

		s := SaleOrder{}
		//s.Reference = orderKey
		s.CustomerId = customer
		s.PaymentMethodId = erpPaymentMethod
		s.CurrencyId = erpCurrency
		s.BillingAddressId = billingAddress
		s.ShippingAddressId = shippingAddress
		s.ShippingPrice = shippingTotal
		s.FixDiscount = discountTax
		s.WooCommerceId = id

		if billingZone == "E" {
			s.BillingSeriesId = *settings.SettingsEcommerce.WooCommerceExportSerieId
		} else if billingZone == "U" && totalTax == 0 {
			s.BillingSeriesId = *settings.SettingsEcommerce.WooCommerceIntracommunitySerieId
		} else {
			s.BillingSeriesId = *settings.SettingsEcommerce.WooCommerceInteriorSerieId
		}

		s.EnterpriseId = enterpriseId
		ok, orderId := s.insertSalesOrder(0)
		if !ok {
			errors = append(errors, "Can't import order. The order could not be created in MARKETNET. Order id "+strconv.Itoa(int(id)))
			continue
		}

		// set the customer details if are empty
		c := getCustomerRow(customer)
		if c.PaymentMethodId == nil {
			c.PaymentMethodId = &erpPaymentMethod
		}
		if c.BillingSeriesId == nil || *c.BillingSeriesId == "" {
			c.BillingSeriesId = &s.BillingSeriesId
		}
		c.updateCustomer(0)

		// insert the details
		sqlStatement = `SELECT id, product_id, variation_id, quantity, total_tax, price FROM public.wc_order_details WHERE "order"=$1 AND enterprise=$2`
		rows, err := db.Query(sqlStatement, id, enterpriseId)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		for rows.Next() {
			var id int32
			var productId int32
			var variationId int32
			var quantity int32
			var totalTax float64
			var price float64
			rows.Scan(&id, &productId, &variationId, &quantity, &totalTax, &price)

			// get the product
			sqlStatement = `SELECT id, vat_percent FROM product WHERE wc_id=$1 AND wc_variation_id=$2 AND enterprise=$3`
			row = db.QueryRow(sqlStatement, productId, variationId, enterpriseId)

			var product int32
			var vatPercent float64
			row.Scan(&product, &vatPercent)

			if product <= 0 {
				errors = append(errors, "Can't import order edtail. The product doesn't exists. Order id "+strconv.Itoa(int(id)))
				continue
			}

			d := SalesOrderDetail{}
			d.OrderId = orderId
			d.ProductId = product
			d.Quantity = quantity
			d.Price = price

			if totalTax == 0 {
				d.VatPercent = 0
			} else {
				d.VatPercent = vatPercent
			}

			d.WooCommerceId = id
			d.EnterpriseId = enterpriseId
			d.insertSalesOrderDetail(0)

		} // for rows.Next() {
		rows.Close()

		// if the payment method is paid in advance, it means that this order is already paid (by VISA o PayPal etc)
		// automatically generate an invoice for this payment
		if paidInAdvance {
			invoiceAllSaleOrder(orderId, enterpriseId, 0)
		}

	} // for rows.Next() {

	if len(errors) > 0 {
		var errorHtml string = ""
		for i := 0; i < len(errors); i++ {
			errorHtml += "<p>Error data:" + errors[i] + "</p>"
		}

		s := getSettingsRecordById(enterpriseId)
		sendEmail(s.SettingsEmail.EmailSendErrorEcommerce, s.SettingsEmail.EmailSendErrorEcommerce, "WooCommerce import error", "<p>Error at: Orders</p>"+errorHtml, enterpriseId)
	}

	return true
} // func copyWcOrders() {

type WooCommerceStatusUpdate struct {
	Id     int32  `json:"id"`
	Status string `json:"status"`
}

func updateTrackingNumberWooCommerceOrder(salesOrderId int64, trackingNumber string, enterpriseId int32) bool {
	// Currently not supported! (WooCommerce's web service does not have a tracking number field in the order)
	return true
}

func updateStatusPaymentAcceptedWooCommerce(salesOrderId int64, enterpriseId int32) bool {
	settings := getSettingsRecordById(enterpriseId)
	if settings.SettingsEcommerce.Ecommerce != "W" {
		return false
	}

	s := getSalesOrderRow(salesOrderId)
	if s.Id <= 0 || s.WooCommerceId <= 0 {
		return false
	}

	url := settings.SettingsEcommerce.WooCommerceUrl + "orders/" + strconv.Itoa(int(s.WooCommerceId))

	update := WooCommerceStatusUpdate{Id: s.WooCommerceId, Status: WOOCOMMERCE_PROCESSING}
	json, _ := json.Marshal(update)

	// OAuth 1.0 request
	config := oauth1.NewConfig(settings.SettingsEcommerce.WooCommerceConsumerKey, settings.SettingsEcommerce.WooCommerceConsumerSecret)
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	_, err := httpClient.Do(req)

	return err == nil
}
