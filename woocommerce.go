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

func getWooCommerceAPI_URL(resourceName string) string {
	s := getSettingsRecord()

	return s.WooCommerceUrl + resourceName
}

func getWooCommerceJSON(URL string) ([]byte, error) {
	s := getSettingsRecord()

	// OAuth 1.0 request
	config := oauth1.NewConfig(s.WooCommerceConsumerKey, s.WooCommerceConsumerSecret)
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
	Id          int32          `json:"id"`
	DateCreated WcJsonDateTime `json:"date_created"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Billing     WCAddress      `json:"billing"`
	Shipping    WCAddress      `json:"shipping"`
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
	Id               int32               `json:"id"`
	Name             string              `json:"name"`
	DateCreated      WcJsonDateTime      `json:"date_created"`
	Description      string              `json:"description"`
	ShortDescription string              `json:"short_description"`
	Sku              string              `json:"sku"`
	Price            float32             `json:"price"`
	Weight           string              `json:"weight"`
	Dimensions       WCPRoductDimensions `json:"dimensions"`
	Images           []WCProductImage    `json:"images"`
	Variations       []int32             `json:"variations"`
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
	Id         int32                         `json:"id"`
	Sku        string                        `json:"sku"`
	Price      float32                       `json:"price"`
	Weight     string                        `json:"weight"`
	Dimensions WCPRoductDimensions           `json:"dimensions"`
	Attributes []WCProductVariationAttribute `json:"attributes"`
}

type WCProductVariationAttribute struct {
	Option string `json:"option"`
}

type WCOrder struct {
	Id            int32           `json:"id"`
	Status        string          `json:"status"`
	Currency      string          `json:"currency"`
	DateCreated   WcJsonDateTime  `json:"date_created"`
	DiscountTax   string          `json:"discount_tax"`
	ShippingTotal string          `json:"shipping_total"`
	ShippingTax   string          `json:"shipping_tax"`
	TotalTax      string          `json:"total_tax"`
	CustomerId    int32           `json:"customer_id"`
	OrderKey      string          `json:"order_key"`
	Billing       WCAddress       `json:"billing"`
	Shipping      WCAddress       `json:"shipping"`
	PaymentMethod string          `json:"payment_method"`
	LineItems     []WCOrderDetail `json:"line_items"`
}

type WCOrderDetail struct {
	Id          int32   `json:"id"`
	ProductId   int32   `json:"product_id"`
	VariationId int32   `json:"variation_id"`
	Quantity    int32   `json:"quantity"`
	TotalTax    string  `json:"total_tax"`
	Price       float32 `json:"price"`
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
func importFromWooCommerce() {
	s := getSettingsRecord()
	if s.Ecommerce != "W" {
		return
	}

	// get all data from WooCommerce, write it in tables like the ones that WooCommerce uses
	importWcCustomers()
	importWcProducts()
	importWcOrders()

	// trasnfer the data form the WooCommerce tables to the ERP
	copyWcCustomers()
	copyWcProducts()
	copyWcOrders()
}

// =====
// COPY THE DATA FROM WOOCOMMERCE TO THE WC MARKETNET TABLES
// =====

func importWcCustomers() {
	url := getWooCommerceAPI_URL("customers")
	jsonWC, err := getWooCommerceJSON(url)
	if err != nil {
		return
	}

	var customers []WCCustomer
	json.Unmarshal(jsonWC, &customers)

	sqlStatement := `UPDATE public.wc_customers SET wc_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(customers); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_customers WHERE id=$1`
		row := db.QueryRow(sqlStatement, customers[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.wc_customers(id, date_created, email, first_name, last_name, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, billing_company) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
			db.Exec(sqlStatement, customers[i].Id, customers[i].DateCreated.ToTime(), customers[i].Email, customers[i].FirstName, customers[i].LastName, customers[i].Billing.Address1, customers[i].Billing.Address2, customers[i].Billing.City, customers[i].Billing.PostCode, customers[i].Billing.Country, customers[i].Billing.State, customers[i].Billing.Phone, customers[i].Shipping.Address1, customers[i].Shipping.Address2, customers[i].Shipping.City, customers[i].Shipping.PostCode, customers[i].Shipping.Country, customers[i].Shipping.State, customers[i].Shipping.Phone, customers[i].Billing.Company)
		} else { // the row exists, update
			sqlStatement := `UPDATE public.wc_customers SET date_created=$2, email=$3, first_name=$4, last_name=$5, billing_address_1=$6, billing_address_2=$7, billing_city=$8, billing_postcode=$9, billing_country=$10, billing_state=$11, billing_phone=$12, shipping_address_1=$13, shipping_address_2=$14, shipping_city=$15, shipping_postcode=$16, shipping_country=$17, shipping_state=$18, shipping_phone=$19, billing_company=$20, wc_exists=true WHERE id=$1`
			db.Exec(sqlStatement, customers[i].Id, customers[i].DateCreated.ToTime(), customers[i].Email, customers[i].FirstName, customers[i].LastName, customers[i].Billing.Address1, customers[i].Billing.Address2, customers[i].Billing.City, customers[i].Billing.PostCode, customers[i].Billing.Country, customers[i].Billing.State, customers[i].Billing.Phone, customers[i].Shipping.Address1, customers[i].Shipping.Address2, customers[i].Shipping.City, customers[i].Shipping.PostCode, customers[i].Shipping.Country, customers[i].Shipping.State, customers[i].Shipping.Phone, customers[i].Billing.Company)
		}
	}

	sqlStatement = `DELETE FROM public.wc_customers WHERE wc_exists=false`
	db.Exec(sqlStatement)
}

func importWcProducts() {
	url := getWooCommerceAPI_URL("products")
	jsonWC, err := getWooCommerceJSON(url)
	if err != nil {
		return
	}

	var products []WCProduct
	json.Unmarshal(jsonWC, &products)

	sqlStatement := `UPDATE public.wc_products SET wc_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `UPDATE public.wc_product_variations SET wc_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(products); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_products WHERE id=$1`
		row := db.QueryRow(sqlStatement, products[i].Id)
		var rows int32
		row.Scan(&rows)

		var images []string = make([]string, 0)
		for j := 0; j < len(products[i].Images); j++ {
			images = append(images, products[i].Images[j].Src)
		}

		if rows == 0 { // the row does not exist, insert
			sqlStatement := `INSERT INTO public.wc_products(id, name, date_created, description, short_description, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, images, variations) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
			db.Exec(sqlStatement, products[i].Id, products[i].Name, products[i].DateCreated.ToTime(), products[i].Description, products[i].ShortDescription, products[i].Sku, products[i].Price, products[i].Weight, products[i].Dimensions.Length, products[i].Dimensions.Width, products[i].Dimensions.Height, pq.Array(images), pq.Array(products[i].Variations))
		} else { // the row exists, update
			sqlStatement := `UPDATE public.wc_products SET name=$2, date_created=$3, description=$4, short_description=$5, sku=$6, price=$7, weight=$8, dimensions_length=$9, dimensions_width=$10, dimensions_height=$11, images=$12, variations=$13, wc_exists=true WHERE id=$1`
			db.Exec(sqlStatement, products[i].Id, products[i].Name, products[i].DateCreated.ToTime(), products[i].Description, products[i].ShortDescription, products[i].Sku, products[i].Price, products[i].Weight, products[i].Dimensions.Length, products[i].Dimensions.Width, products[i].Dimensions.Height, pq.Array(images), pq.Array(products[i].Variations))
		}

		// get the variations
		url := getWooCommerceAPI_URL("products") + "/" + strconv.Itoa(int(products[i].Id)) + "/variations/"
		jsonWC, err := getWooCommerceJSON(url)
		if err != nil {
			continue
		}

		var variations []WCProductVariation
		json.Unmarshal(jsonWC, &variations)

		// insert/update variations
		for j := 0; j < len(variations); j++ {
			// ¿does the row exist?
			sqlStatement := `SELECT COUNT(*) FROM wc_product_variations WHERE id=$1`
			row := db.QueryRow(sqlStatement, variations[j].Id)
			var rows int32
			row.Scan(&rows)

			var attributes []string = make([]string, 0)
			for k := 0; k < len(variations[j].Attributes); k++ {
				attributes = append(attributes, variations[j].Attributes[k].Option)
			}

			if rows == 0 { // the row does not exist, insert
				sqlStatement := `INSERT INTO public.wc_product_variations(id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
				db.Exec(sqlStatement, variations[j].Id, variations[j].Sku, variations[j].Price, variations[j].Weight, variations[j].Dimensions.Length, variations[j].Dimensions.Width, variations[j].Dimensions.Height, pq.Array(attributes))
			} else { // the row exists, update
				sqlStatement := `UPDATE public.wc_product_variations SET sku=$2, price=$3, weight=$4, dimensions_length=$5, dimensions_width=$6, dimensions_height=$7, attributes=$8, wc_exists=true WHERE id=$1`
				db.Exec(sqlStatement, variations[j].Id, variations[j].Sku, variations[j].Price, variations[j].Weight, variations[j].Dimensions.Length, variations[j].Dimensions.Width, variations[j].Dimensions.Height, pq.Array(attributes))
			}
		}

	}

	sqlStatement = `DELETE FROM public.wc_products WHERE wc_exists=false`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.wc_product_variations WHERE wc_exists=false`
	db.Exec(sqlStatement)
}

func importWcOrders() {
	url := getWooCommerceAPI_URL("orders")
	jsonWC, err := getWooCommerceJSON(url)
	if err != nil {
		return
	}

	var orders []WCOrder
	json.Unmarshal(jsonWC, &orders)

	sqlStatement := `UPDATE public.wc_orders SET wc_exists=false`
	db.Exec(sqlStatement)

	for i := 0; i < len(orders); i++ {
		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM wc_orders WHERE id=$1`
		row := db.QueryRow(sqlStatement, orders[i].Id)
		var rows int32
		row.Scan(&rows)

		if rows == 0 {
			f_discount_tax, err := strconv.ParseFloat(orders[i].DiscountTax, 32)
			if err != nil {
				continue
			}
			f_shipping_total, err := strconv.ParseFloat(orders[i].ShippingTotal, 32)
			if err != nil {
				continue
			}
			f_shipping_tax, err := strconv.ParseFloat(orders[i].ShippingTax, 32)
			if err != nil {
				continue
			}
			f_total_tax, err := strconv.ParseFloat(orders[i].TotalTax, 32)
			if err != nil {
				continue
			}

			sqlStatement = `INSERT INTO public.wc_orders(id, status, currency, date_created, discount_tax, shipping_total, shipping_tax, total_tax, customer_id, order_key, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, payment_method, billing_company) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)`
			db.Exec(sqlStatement, orders[i].Id, orders[i].Status, orders[i].Currency, orders[i].DateCreated.ToTime(), f_discount_tax, f_shipping_total, f_shipping_tax, f_total_tax, orders[i].CustomerId, orders[i].OrderKey, orders[i].Billing.Address1, orders[i].Billing.Address2, orders[i].Billing.City, orders[i].Billing.PostCode, orders[i].Billing.Country, orders[i].Billing.State, orders[i].Billing.Phone, orders[i].Shipping.Address1, orders[i].Shipping.Address2, orders[i].Shipping.City, orders[i].Shipping.PostCode, orders[i].Shipping.Country, orders[i].Shipping.State, orders[i].Shipping.Phone, orders[i].PaymentMethod, orders[i].Billing.Company)

			// add order details
			for j := 0; j < len(orders[i].LineItems); j++ {
				f_total_tax, err := strconv.ParseFloat(orders[i].LineItems[j].TotalTax, 32)
				if err != nil {
					continue
				}

				sqlStatement := `INSERT INTO public.wc_order_details(id, "order", product_id, variation_id, quantity, total_tax, price) VALUES ($1, $2, $3, $4, $5, $6, $7)`
				db.Exec(sqlStatement, orders[i].LineItems[j].Id, orders[i].Id, orders[i].LineItems[j].ProductId, orders[i].LineItems[j].VariationId, orders[i].LineItems[j].Quantity, f_total_tax, orders[i].LineItems[j].Price)
			}

		} else { // if rows == 0
			sqlStatement := `UPDATE public.wc_orders SET wc_exists=true WHERE id=$1`
			db.Exec(sqlStatement, orders[i].Id)
		}
	} // for

	sqlStatement = `DELETE FROM wc_order_details WHERE NOT (SELECT wc_exists FROM wc_orders WHERE wc_orders.id=wc_order_details."order")`
	db.Exec(sqlStatement)
	sqlStatement = `DELETE FROM public.wc_orders WHERE wc_exists=false`
	db.Exec(sqlStatement)
}

// =====
// TRANSFER THE DATA TO THE ERP TABLES
// =====

func copyWcCustomers() {
	sqlStatement := `SELECT id FROM public.wc_customers`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var wcCustomerId int32
		rows.Scan(&wcCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM customer WHERE wc_id=$1`
		rowCount := db.QueryRow(sqlStatement, wcCustomerId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, date_created, email, first_name, last_name, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, billing_phone, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, shipping_phone, billing_company FROM public.wc_customers WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, wcCustomerId)
		if row.Err() != nil {
			return
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
			ok, customerId := c.insertCustomer()
			if !ok {
				continue
			}

			// create billing address
			ba := Address{}
			ba.Customer = &customerId
			ba.Address = billingAddress1
			ba.Address2 = billingAddress2
			ba.City = billingCity
			ba.ZipCode = billingPostcode
			// search for the country by iso code
			if len(billingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					ba.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else if len(billingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					ba.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(billingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingState)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						ba.State = &stateid
					}
				} else {
					log("DB", row.Err().Error())
				}
			}
			if len(billingCompany) == 0 {
				ba.PrivateOrBusiness = "P"
			} else {
				ba.PrivateOrBusiness = "B"
			}
			ba.insertAddress()

			// create shipping address
			sa := Address{}
			sa.Customer = &customerId
			sa.Address = shippingAddress1
			sa.Address2 = shippingAddress2
			sa.City = shippingCity
			sa.ZipCode = shippingPostcode
			// search for the country by iso code
			if len(shippingCountry) == 2 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					sa.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else if len(shippingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					sa.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(shippingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingState)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						sa.State = &stateid
					}
				} else {
					log("DB", row.Err().Error())
				}
			}
			if len(billingCompany) == 0 {
				sa.PrivateOrBusiness = "P"
			} else {
				sa.PrivateOrBusiness = "B"
			}
			sa.insertAddress()

			// set as the main billing and shipping address
			customer := getCustomerRow(customerId)
			customer.MainBillingAddress = &ba.Id
			customer.MainShippingAddress = &sa.Id
			customer.updateCustomer()
		} else {
			// update the addresses
			sqlStatement := `SELECT id FROM public.customer WHERE wc_id=$1 LIMIT 1`
			row := db.QueryRow(sqlStatement, wcCustomerId)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				continue
			}

			var customerIdErp int32
			row.Scan(&customerIdErp)

			addresses := getCustomerAddresses(customerIdErp)

			// compare both the billing and shipping address from woocommerce with the ones that are in the erp.
			// if the addresses are different, create new ones

			// billing address
			var billingAddressFound bool = false
			for i := 0; i < len(addresses); i++ {
				// search for the country by iso code
				var countryId int16
				var stateId *int32
				if len(billingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("DB", row.Err().Error())
					}
				} else if len(billingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("DB", row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(billingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingState)
					if row.Err() == nil {
						var id int32
						row.Scan(&id)
						if id > 0 {
							stateId = &id
						}
					} else {
						log("DB", row.Err().Error())
					}
				}

				// are the addresses the same?
				if addresses[i].Address == billingAddress1 && addresses[i].Address2 == billingAddress2 && addresses[i].City == billingCity && addresses[i].ZipCode == billingPostcode && addresses[i].Country == countryId && addresses[i].State == stateId {
					billingAddressFound = true
					break
				} // if
			} // for
			if !billingAddressFound {
				// create billing address
				ba := Address{}
				ba.Customer = &customerIdErp
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
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry)
					if row.Err() == nil {
						var countryId int16
						row.Scan(&countryId)
						ba.Country = countryId
					} else {
						log("DB", row.Err().Error())
					}
				} else if len(billingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingCountry)
					if row.Err() == nil {
						var countryId int16
						row.Scan(&countryId)
						ba.Country = countryId
					} else {
						log("DB", row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(billingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, billingState)
					if row.Err() == nil {
						var stateid int32
						row.Scan(&stateid)
						if stateid > 0 {
							ba.State = &stateid
						}
					} else {
						log("DB", row.Err().Error())
					}
				}
				ba.insertAddress()
			}

			// shipping address
			var shippingAddressFound bool = false
			for i := 0; i < len(addresses); i++ {
				// search for the country by iso code
				var countryId int16
				var stateId *int32
				if len(shippingCountry) == 2 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("DB", row.Err().Error())
					}
				} else if len(shippingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry)
					if row.Err() == nil {
						row.Scan(&countryId)
					} else {
						log("DB", row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(shippingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingState)
					if row.Err() == nil {
						var id int32
						row.Scan(&id)
						if id > 0 {
							stateId = &id
						}
					} else {
						log("DB", row.Err().Error())
					}
				}

				// are the addresses the same?
				if addresses[i].Address == shippingAddress1 && addresses[i].Address2 == shippingAddress2 && addresses[i].City == shippingCity && addresses[i].ZipCode == shippingPostcode && addresses[i].Country == countryId && addresses[i].State == stateId {
					shippingAddressFound = true
					break
				} // if
			} // for
			if !shippingAddressFound {
				// create shipping address
				sa := Address{}
				sa.Customer = &customerIdErp
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
					sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry)
					if row.Err() == nil {
						var countryId int16
						row.Scan(&countryId)
						sa.Country = countryId
					} else {
						log("DB", row.Err().Error())
					}
				} else if len(shippingCountry) == 3 {
					sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingCountry)
					if row.Err() == nil {
						var countryId int16
						row.Scan(&countryId)
						sa.Country = countryId
					} else {
						log("DB", row.Err().Error())
					}
				} else {
					continue
				}
				// search for the state
				if len(shippingState) > 0 {
					sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, shippingState)
					if row.Err() == nil {
						var stateid int32
						row.Scan(&stateid)
						if stateid > 0 {
							sa.State = &stateid
						}
					} else {
						log("DB", row.Err().Error())
					}
				}
				sa.insertAddress()
			}

		} // else
	} // for rows.Next()
} // copyWcCustomers

func copyWcProducts() {
	s := getSettingsRecord()

	sqlStatement := `SELECT id FROM public.wc_products`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var wcCustomerId int32
		rows.Scan(&wcCustomerId)

		// ¿does the row exist?
		sqlStatement := `SELECT COUNT(*) FROM product WHERE wc_id=$1`
		rowCount := db.QueryRow(sqlStatement, wcCustomerId)
		var rows int32
		rowCount.Scan(&rows)

		sqlStatement = `SELECT id, name, date_created, description, short_description, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, images, variations FROM public.wc_products WHERE id=$1 LIMIT 1`
		row := db.QueryRow(sqlStatement, wcCustomerId)
		if row.Err() != nil {
			return
		}

		var id int32
		var name string
		var dateCreated time.Time
		var description string
		var shortDescription string
		var sku string
		var price float32
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
					p.Weight = float32(f_weight)
				}
				f_length, err := strconv.ParseFloat(dimensionsLength, 32)
				if err == nil {
					p.Depth = float32(f_length)
				}
				f_width, err := strconv.ParseFloat(dimensionsWidth, 32)
				if err == nil {
					p.Width = float32(f_width)
				}
				f_height, err := strconv.ParseFloat(dimensionsHeight, 32)
				if err == nil {
					p.Height = float32(f_height)
				}
				p.VatPercent = s.DefaultVatPercent
				p.WooCommerceId = id
				p.insertProduct()

				for i := 0; i < len(images); i++ {
					pi := ProductImage{
						Product: p.Id,
						URL:     images[i],
					}
					pi.insertProductImage()
				}
			} else {
				// the product has variations. create a new product by each variation (combination of options)
				for i := 0; i < len(variations); i++ {
					sqlStatement := `SELECT id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes FROM public.wc_product_variations WHERE id=$1 LIMIT 1`
					row := db.QueryRow(sqlStatement, variations[i])
					if row.Err() != nil {
						continue
					}

					var v_id int32
					var v_sku string
					var v_price float32
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
						p.Weight = float32(f_weight)
					}
					f_length, err := strconv.ParseFloat(v_dimensionsLength, 32)
					if err == nil {
						p.Depth = float32(f_length)
					}
					f_width, err := strconv.ParseFloat(v_dimensionsWidth, 32)
					if err == nil {
						p.Width = float32(f_width)
					}
					f_height, err := strconv.ParseFloat(v_dimensionsHeight, 32)
					if err == nil {
						p.Height = float32(f_height)
					}
					p.VatPercent = s.DefaultVatPercent
					p.WooCommerceId = id
					p.WooCommerceVariationId = v_id
					p.insertProduct()

					for i := 0; i < len(images); i++ {
						pi := ProductImage{
							Product: p.Id,
							URL:     images[i],
						}
						pi.insertProductImage()
					}
				}
			} // else
		} else { // if rows == 0 {

			if len(variations) == 0 {
				sqlStatement := `SELECT id FROM product WHERE wc_id=$1 AND wc_variation_id=$2 LIMIT 1`
				row := db.QueryRow(sqlStatement, id, 0)
				if row.Err() != nil {
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
					p.Weight = float32(f_weight)
				}
				f_length, err := strconv.ParseFloat(dimensionsLength, 32)
				if err == nil {
					p.Depth = float32(f_length)
				}
				f_width, err := strconv.ParseFloat(dimensionsWidth, 32)
				if err == nil {
					p.Width = float32(f_width)
				}
				f_height, err := strconv.ParseFloat(dimensionsHeight, 32)
				if err == nil {
					p.Height = float32(f_height)
				}
				p.updateProduct()
			} else { // if len(variations) == 0 {
				for i := 0; i < len(variations); i++ {
					sqlStatement := `SELECT id FROM product WHERE wc_id=$1 AND wc_variation_id=$2 LIMIT 1`
					row := db.QueryRow(sqlStatement, id, variations[i])
					if row.Err() != nil {
						continue
					}

					var productId int32
					row.Scan(&productId)
					p := getProductRow(productId)

					sqlStatement = `SELECT id, sku, price, weight, dimensions_length, dimensions_width, dimensions_height, attributes FROM public.wc_product_variations WHERE id=$1 LIMIT 1`
					row = db.QueryRow(sqlStatement, variations[i])
					if row.Err() != nil {
						continue
					}

					var v_id int32
					var v_sku string
					var v_price float32
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
						p.Weight = float32(f_weight)
					}
					f_length, err := strconv.ParseFloat(v_dimensionsLength, 32)
					if err == nil {
						p.Depth = float32(f_length)
					}
					f_width, err := strconv.ParseFloat(v_dimensionsWidth, 32)
					if err == nil {
						p.Width = float32(f_width)
					}
					f_height, err := strconv.ParseFloat(v_dimensionsHeight, 32)
					if err == nil {
						p.Height = float32(f_height)
					}
					p.updateProduct()
				} // for
			} // else
		} // else
	} // for
} // copyWcProducts

func copyWcOrders() {
	settings := getSettingsRecord()

	sqlStatement := `SELECT id, status, currency, date_created, discount_tax, shipping_total, shipping_tax, total_tax, customer_id, order_key, billing_address_1, billing_address_2, billing_city, billing_postcode, billing_country, billing_state, shipping_address_1, shipping_address_2, shipping_city, shipping_postcode, shipping_country, shipping_state, payment_method, billing_company FROM public.wc_orders;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}

	for rows.Next() {
		var id int32
		var status string
		var currency string
		var dateCreated time.Time
		var discountTax float32
		var shippingTotal float32
		var shippingTax float32
		var totalTax float32
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
		sqlStatement := `SELECT COUNT(id) FROM sales_order WHERE wc_id=$1`
		row := db.QueryRow(sqlStatement, id)

		var orders int32
		row.Scan(&orders)

		if orders > 0 { // don't continue if the order exists
			continue
		}

		// get the customer
		sqlStatement = `SELECT id FROM customer WHERE wc_id=$1 LIMIT 1`
		row = db.QueryRow(sqlStatement, customerId)

		var customer int32
		row.Scan(&customer)

		if customer == 0 { // don't continue if the customer doesn't exists
			continue
		}

		// get the payment method
		sqlStatement = `SELECT id,paid_in_advance FROM payment_method WHERE woocommerce_module_name=$1`
		row = db.QueryRow(sqlStatement, paymentMethod)

		var erpPaymentMethod int16
		var paidInAdvance bool
		row.Scan(&erpPaymentMethod, &paidInAdvance)

		if erpPaymentMethod == 0 { // attemp the default one in the settings (no payment method = likely a manual order)
			if settings.WooCommerceDefaultPaymentMethod == nil { // don't continue if the payment method doesn't exists
				continue
			}
			erpPaymentMethod = *settings.WooCommerceDefaultPaymentMethod
		}

		// get the currency
		sqlStatement = `SELECT id FROM currency WHERE iso_code=$1`
		row = db.QueryRow(sqlStatement, currency)

		var erpCurrency int16
		row.Scan(&erpCurrency)

		if erpCurrency == 0 { // don't continue if the currency doesn't exists
			continue
		}

		// get the billing address
		var billingAddress int32
		var billingZone string
		customerAddresses := getCustomerAddresses(customer)
		for i := 0; i < len(customerAddresses); i++ {
			country := getCountryRow(customerAddresses[i].Country)
			if customerAddresses[i].Address == billingAddress1 && customerAddresses[i].Address2 == billingAddress2 && customerAddresses[i].City == billingCity && customerAddresses[i].ZipCode == billingPostcode && (country.Iso2 == billingCountry || country.Iso3 == billingCountry) {
				billingAddress = customerAddresses[i].Id
				billingZone = country.Zone
				break
			}
		}
		if billingAddress == 0 {
			// create billing address
			ba := Address{}
			ba.Customer = &customer
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
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					ba.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else if len(billingCountry) == 3 {
				sqlStatement := `SELECT id, zone FROM public.country WHERE iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId, &billingZone)
					ba.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(billingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, billingState)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						ba.State = &stateid
					}
				} else {
					log("DB", row.Err().Error())
				}
			}
			ba.insertAddress()
			billingAddress = ba.Id
		}

		// get the shipping address
		var shippingAddress int32
		for i := 0; i < len(customerAddresses); i++ {
			country := getCountryRow(customerAddresses[i].Country)
			if customerAddresses[i].Address == shippingAddress1 && customerAddresses[i].Address2 == shippingAddress2 && customerAddresses[i].City == shippingCity && customerAddresses[i].ZipCode == shippingPostcode && (country.Iso2 == shippingCountry || country.Iso3 == shippingCountry) {
				shippingAddress = customerAddresses[i].Id
				break
			}
		}
		if shippingAddress == 0 {
			// create shipping address
			sa := Address{}
			sa.Customer = &customer
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
				sqlStatement := `SELECT id FROM public.country WHERE iso_2=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					sa.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else if len(shippingCountry) == 3 {
				sqlStatement := `SELECT id FROM public.country WHERE iso_3=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingCountry)
				if row.Err() == nil {
					var countryId int16
					row.Scan(&countryId)
					sa.Country = countryId
				} else {
					log("DB", row.Err().Error())
				}
			} else {
				continue
			}
			// search for the state
			if len(shippingState) > 0 {
				sqlStatement := `SELECT id FROM public.state WHERE iso_code=$1 LIMIT 1`
				row := db.QueryRow(sqlStatement, shippingState)
				if row.Err() == nil {
					var stateid int32
					row.Scan(&stateid)
					if stateid > 0 {
						sa.State = &stateid
					}
				} else {
					log("DB", row.Err().Error())
				}
			}
			sa.insertAddress()
			shippingAddress = sa.Id
		}

		s := SaleOrder{}
		s.Warehouse = settings.DefaultWarehouse
		//s.Reference = orderKey
		s.Customer = customer
		s.PaymentMethod = erpPaymentMethod
		s.Currency = erpCurrency
		s.BillingAddress = billingAddress
		s.ShippingAddress = shippingAddress
		s.ShippingPrice = shippingTotal
		s.FixDiscount = discountTax
		s.WooCommerceId = id

		if billingZone == "E" {
			s.BillingSeries = *settings.WooCommerceExportSerie
		} else if billingZone == "U" && totalTax == 0 {
			s.BillingSeries = *settings.WooCommerceIntracommunitySerie
		} else {
			s.BillingSeries = *settings.WooCommerceInteriorSerie
		}

		ok, orderId := s.insertSalesOrder()
		if !ok {
			continue
		}

		// set the customer details if are empty
		c := getCustomerRow(customer)
		if c.PaymentMethod == nil {
			c.PaymentMethod = &erpPaymentMethod
		}
		if c.BillingSeries == nil || *c.BillingSeries == "" {
			c.BillingSeries = &s.BillingSeries
		}
		c.updateCustomer()

		// insert the details
		sqlStatement = `SELECT id, product_id, variation_id, quantity, total_tax, price FROM public.wc_order_details WHERE "order"=$1`
		rows, err := db.Query(sqlStatement, id)
		if err != nil {
			return
		}

		for rows.Next() {
			var id int32
			var productId int32
			var variationId int32
			var quantity int32
			var totalTax float32
			var price float32
			rows.Scan(&id, &productId, &variationId, &quantity, &totalTax, &price)

			// get the product
			sqlStatement = `SELECT id, vat_percent FROM product WHERE wc_id=$1 AND wc_variation_id=$2`
			row = db.QueryRow(sqlStatement, productId, variationId)

			var product int32
			var vatPercent float32
			row.Scan(&product, &vatPercent)

			if product <= 0 {
				continue
			}

			d := SalesOrderDetail{}
			d.Order = orderId
			d.Product = product
			d.Quantity = quantity
			d.Price = price

			if totalTax == 0 {
				d.VatPercent = 0
			} else {
				d.VatPercent = vatPercent
			}

			d.WooCommerceId = id
			d.insertSalesOrderDetail()

		} // for rows.Next() {

		// if the payment method is paid in advance, it means that this order is already paid (by VISA o PayPal etc)
		// automatically generate an invoice for this payment
		if paidInAdvance {
			invoiceAllSaleOrder(orderId)
		}

	} // for rows.Next() {
} // func copyWcOrders() {

type WooCommerceStatusUpdate struct {
	Id     int32  `json:"id"`
	Status string `json:"status"`
}

func updateTrackingNumberWooCommerceOrder(salesOrderId int32, trackingNumber string) bool {
	// Currently not supported! (WooCommerce's web service does not have a tracking number field in the order)
	return true
}

func updateStatusPaymentAcceptedWooCommerce(salesOrderId int32) bool {
	settings := getSettingsRecord()
	if settings.Ecommerce != "W" {
		return false
	}

	s := getSalesOrderRow(salesOrderId)
	if s.Id <= 0 || s.WooCommerceId <= 0 {
		return false
	}

	url := settings.WooCommerceUrl + "orders/" + strconv.Itoa(int(s.WooCommerceId))

	update := WooCommerceStatusUpdate{Id: s.WooCommerceId, Status: WOOCOMMERCE_PROCESSING}
	json, _ := json.Marshal(update)

	// OAuth 1.0 request
	config := oauth1.NewConfig(settings.WooCommerceConsumerKey, settings.WooCommerceConsumerSecret)
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	_, err := httpClient.Do(req)

	return err == nil
}