package main

import (
	"strings"
	"time"
)

type Supplier struct {
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
	City                *int32    `json:"city"`
	MainShippingAddress *int32    `json:"mainShippingAddress"`
	MainBillingAddress  *int32    `json:"mainBillingAddress"`
	Language            *int16    `json:"language"`
	PaymentMethod       *int32    `json:"paymentMethod"`
	BillingSeries       *string   `json:"billingSeries"`
	DateCreated         time.Time `json:"dateCreated"`
}

func getSuppliers() []Supplier {
	var customers []Supplier = make([]Supplier, 0)
	sqlStatement := `SELECT * FROM public.suppliers ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return customers
	}
	for rows.Next() {
		c := Supplier{}
		rows.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.City, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated)
		customers = append(customers, c)
	}

	return customers
}

func (s *Supplier) isValid() bool {
	return !(len(s.Name) == 0 || len(s.Name) > 303 || len(s.Tradename) == 0 || len(s.Tradename) > 150 || len(s.FiscalName) == 0 || len(s.FiscalName) > 150 || len(s.TaxId) > 25 || len(s.VatNumber) > 25 || len(s.Phone) > 25 || len(s.Email) > 100)
}

func (s *Supplier) insertSupplier() bool {
	if !s.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.suppliers(name, tradename, fiscal_name, tax_id, vat_number, phone, email, main_address, country, city, main_shipping_address, main_billing_address, language, payment_method, billing_series) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	res, err := db.Exec(sqlStatement, s.Name, s.Tradename, s.FiscalName, s.TaxId, s.VatNumber, s.Phone, s.Email, s.MainAddress, s.Country, s.City, s.MainShippingAddress, s.MainBillingAddress, s.Language, s.PaymentMethod, s.BillingSeries)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Supplier) updateSupplier() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.suppliers SET name=$2, tradename=$3, fiscal_name=$4, tax_id=$5, vat_number=$6, phone=$7, email=$8, main_address=$9, country=$10, city=$11, main_shipping_address=$12, main_billing_address=$13, language=$14, payment_method=$15, billing_series=$16 WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id, s.Name, s.Tradename, s.FiscalName, s.TaxId, s.VatNumber, s.Phone, s.Email, s.MainAddress, s.Country, s.City, s.MainShippingAddress, s.MainBillingAddress, s.Language, s.PaymentMethod, s.BillingSeries)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Supplier) deleteSupplier() bool {
	if s.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.suppliers WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type SupplierName struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func findSupplierByName(languageName string) []CustomerName {
	var customers []CustomerName = make([]CustomerName, 0)
	sqlStatement := `SELECT id,name FROM public.suppliers WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return customers
	}
	for rows.Next() {
		c := CustomerName{}
		rows.Scan(&c.Id, &c.Name)
		customers = append(customers, c)
	}

	return customers
}

func getNameSupplier(id int32) string {
	sqlStatement := `SELECT name FROM public.suppliers WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type SupplierDefauls struct {
	MainShippingAddress     *int32   `json:"mainShippingAddress"`
	MainShippingAddressName *string  `json:"mainShippingAddressName"`
	MainBillingAddress      *int32   `json:"mainBillingAddress"`
	MainBillingAddressName  *string  `json:"mainBillingAddressName"`
	PaymentMethod           *int32   `json:"paymentMethod"`
	PaymentMethodName       *string  `json:"paymentMethodName"`
	BillingSeries           *string  `json:"billingSeries"`
	BillingSeriesName       *string  `json:"billingSeriesName"`
	Currency                *int16   `json:"currency"`
	CurrencyName            *string  `json:"currencyName"`
	CurrencyChange          *float32 `json:"currencyChange"`
}

func getSupplierDefaults(customerId int32) SupplierDefauls {
	sqlStatement := `SELECT main_shipping_address, (SELECT address AS main_shipping_address_name FROM address WHERE address.id = customer.main_shipping_address), main_billing_address, (SELECT address AS main_billing_address_name FROM address WHERE address.id = customer.main_billing_address), payment_method, (SELECT name AS payment_method_name FROM payment_method WHERE payment_method.id = customer.payment_method), billing_series, (SELECT name AS billing_series_name FROM billing_series WHERE billing_series.id = customer.billing_series), (SELECT currency FROM country WHERE country.id = customer.country), (SELECT name AS currency_name FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)), (SELECT change FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)) FROM public.suppliers WHERE id = $1`
	row := db.QueryRow(sqlStatement, customerId)
	if row.Err() != nil {
		return SupplierDefauls{}
	}
	s := SupplierDefauls{}
	row.Scan(&s.MainShippingAddress, &s.MainShippingAddressName, &s.MainBillingAddress, &s.MainBillingAddressName, &s.PaymentMethod, &s.PaymentMethodName, &s.BillingSeries, &s.BillingSeriesName, &s.Currency, &s.CurrencyName, &s.CurrencyChange)
	return s
}
