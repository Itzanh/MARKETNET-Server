package main

import (
	"strings"
	"time"
)

type Customer struct {
	Id                  int32     `json:"id"`
	Name                string    `json:"name"`
	Tradename           string    `json:"tradename"`
	FiscalName          string    `json:"fiscalName"`
	TaxId               string    `json:"taxId"`
	VatNumber           string    `json:"vatNumber"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email"`
	MainAddress         *int32    `json:"mainAddress"`
	Country             *int16    `json:"country"`
	State               *int32    `json:"state"`
	MainShippingAddress *int32    `json:"mainShippingAddress"`
	MainBillingAddress  *int32    `json:"mainBillingAddress"`
	Language            *int16    `json:"language"`
	PaymentMethod       *int16    `json:"paymentMethod"`
	BillingSeries       *string   `json:"billingSeries"`
	DateCreated         time.Time `json:"dateCreated"`
	PrestaShopId        int32     `json:"prestaShopId"`
	CountryName         *string   `json:"countryName"`
}

func getCustomers() []Customer {
	var customers []Customer = make([]Customer, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=customer.country) FROM public.customer ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return customers
	}
	for rows.Next() {
		c := Customer{}
		rows.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.PrestaShopId, &c.CountryName)
		customers = append(customers, c)
	}

	return customers
}

func searchCustomers(search string) []Customer {
	var customers []Customer = make([]Customer, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=customer.country) FROM customer WHERE name ILIKE $1 OR tax_id ILIKE $1 OR email ILIKE $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%")
	if err != nil {
		return customers
	}
	for rows.Next() {
		c := Customer{}
		rows.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.PrestaShopId, &c.CountryName)
		customers = append(customers, c)
	}

	return customers
}

func getCustomerRow(customerId int32) Customer {
	sqlStatement := `SELECT * FROM public.customer WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, customerId)
	if row.Err() != nil {
		return Customer{}
	}

	c := Customer{}
	row.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.PrestaShopId)

	return c
}

func (c *Customer) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 303 || len(c.Tradename) == 0 || len(c.Tradename) > 150 || len(c.FiscalName) == 0 || len(c.FiscalName) > 150 || len(c.TaxId) > 25 || len(c.VatNumber) > 25 || len(c.Phone) > 25 || len(c.Email) > 100)
}

func (c *Customer) insertCustomer() bool {
	if !c.isValid() {
		return false
	}

	if c.BillingSeries != nil && *c.BillingSeries == "" {
		c.BillingSeries = nil
	}

	sqlStatement := `INSERT INTO public.customer(name, tradename, fiscal_name, tax_id, vat_number, phone, email, main_address, country, state, main_shipping_address, main_billing_address, language, payment_method, billing_series, ps_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	res, err := db.Exec(sqlStatement, c.Name, c.Tradename, c.FiscalName, c.TaxId, c.VatNumber, c.Phone, c.Email, c.MainAddress, c.Country, c.State, c.MainShippingAddress, c.MainBillingAddress, c.Language, c.PaymentMethod, c.BillingSeries, c.PrestaShopId)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Customer) updateCustomer() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	if c.BillingSeries != nil && *c.BillingSeries == "" {
		c.BillingSeries = nil
	}

	sqlStatement := `UPDATE public.customer SET name=$2, tradename=$3, fiscal_name=$4, tax_id=$5, vat_number=$6, phone=$7, email=$8, main_address=$9, country=$10, state=$11, main_shipping_address=$12, main_billing_address=$13, language=$14, payment_method=$15, billing_series=$16 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Tradename, c.FiscalName, c.TaxId, c.VatNumber, c.Phone, c.Email, c.MainAddress, c.Country, c.State, c.MainShippingAddress, c.MainBillingAddress, c.Language, c.PaymentMethod, c.BillingSeries)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Customer) deleteCustomer() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.customer WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCustomerByName(languageName string) []NameInt32 {
	var customers []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.customer WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return customers
	}
	for rows.Next() {
		c := NameInt32{}
		rows.Scan(&c.Id, &c.Name)
		customers = append(customers, c)
	}

	return customers
}

func getNameCustomer(id int32) string {
	sqlStatement := `SELECT name FROM public.customer WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

// Used both in customers and suppliers
type ContactDefauls struct {
	MainShippingAddress     *int32   `json:"mainShippingAddress"`
	MainShippingAddressName *string  `json:"mainShippingAddressName"`
	MainBillingAddress      *int32   `json:"mainBillingAddress"`
	MainBillingAddressName  *string  `json:"mainBillingAddressName"`
	PaymentMethod           *int16   `json:"paymentMethod"`
	PaymentMethodName       *string  `json:"paymentMethodName"`
	BillingSeries           *string  `json:"billingSeries"`
	BillingSeriesName       *string  `json:"billingSeriesName"`
	Currency                *int16   `json:"currency"`
	CurrencyName            *string  `json:"currencyName"`
	CurrencyChange          *float32 `json:"currencyChange"`
}

func getCustomerDefaults(customerId int32) ContactDefauls {
	sqlStatement := `SELECT main_shipping_address, (SELECT address AS main_shipping_address_name FROM address WHERE address.id = customer.main_shipping_address), main_billing_address, (SELECT address AS main_billing_address_name FROM address WHERE address.id = customer.main_billing_address), payment_method, (SELECT name AS payment_method_name FROM payment_method WHERE payment_method.id = customer.payment_method), billing_series, (SELECT name AS billing_series_name FROM billing_series WHERE billing_series.id = customer.billing_series), (SELECT currency FROM country WHERE country.id = customer.country), (SELECT name AS currency_name FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)), (SELECT exchange FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)) FROM public.customer WHERE id = $1`
	row := db.QueryRow(sqlStatement, customerId)
	if row.Err() != nil {
		return ContactDefauls{}
	}
	c := ContactDefauls{}
	row.Scan(&c.MainShippingAddress, &c.MainShippingAddressName, &c.MainBillingAddress, &c.MainBillingAddressName, &c.PaymentMethod, &c.PaymentMethodName, &c.BillingSeries, &c.BillingSeriesName, &c.Currency, &c.CurrencyName, &c.CurrencyChange)
	return c
}

func getCustomerAddresses(customerId int32) []Address {
	var addresses []Address = make([]Address, 0)
	sqlStatement := `SELECT *,CASE WHEN address.customer IS NOT NULL THEN (SELECT name FROM customer WHERE customer.id=address.customer) ELSE (SELECT name FROM suppliers WHERE suppliers.id=address.supplier) END,(SELECT name FROM country WHERE country.id=address.country),(SELECT name FROM state WHERE state.id=address.state) FROM address WHERE customer=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, customerId)
	if err != nil {
		return addresses
	}
	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode, &a.ContactName, &a.CountryName, &a.StateName)
		addresses = append(addresses, a)
	}

	return addresses
}

func getCustomerSaleOrders(customerId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order WHERE customer=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, customerId)
	if err != nil {
		return sales
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.PrestaShopId,
			&s.CustomerName)
		sales = append(sales, s)
	}

	return sales
}
