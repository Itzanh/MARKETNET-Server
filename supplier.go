package main

import (
	"strconv"
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
	State               *int32    `json:"state"`
	MainShippingAddress *int32    `json:"mainShippingAddress"`
	MainBillingAddress  *int32    `json:"mainBillingAddress"`
	Language            *int16    `json:"language"`
	PaymentMethod       *int16    `json:"paymentMethod"`
	BillingSeries       *string   `json:"billingSeries"`
	DateCreated         time.Time `json:"dateCreated"`
	Account             *int32    `json:"account"`
	CountryName         *string   `json:"countryName"`
}

func getSuppliers() []Supplier {
	var suppliers []Supplier = make([]Supplier, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=suppliers.country) FROM public.suppliers ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return suppliers
	}
	for rows.Next() {
		s := Supplier{}
		rows.Scan(&s.Id, &s.Name, &s.Tradename, &s.FiscalName, &s.TaxId, &s.VatNumber, &s.Phone, &s.Email, &s.MainAddress, &s.Country, &s.State, &s.MainShippingAddress, &s.MainBillingAddress, &s.Language, &s.PaymentMethod, &s.BillingSeries, &s.DateCreated, &s.Account, &s.CountryName)
		suppliers = append(suppliers, s)
	}

	return suppliers
}

func searchSuppliers(search string) []Supplier {
	var suppliers []Supplier = make([]Supplier, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=suppliers.country) FROM suppliers WHERE name ILIKE $1 OR tax_id ILIKE $1 OR email ILIKE $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%")
	if err != nil {
		log("DB", err.Error())
		return suppliers
	}
	for rows.Next() {
		s := Supplier{}
		rows.Scan(&s.Id, &s.Name, &s.Tradename, &s.FiscalName, &s.TaxId, &s.VatNumber, &s.Phone, &s.Email, &s.MainAddress, &s.Country, &s.State, &s.MainShippingAddress, &s.MainBillingAddress, &s.Language, &s.PaymentMethod, &s.BillingSeries, &s.DateCreated, &s.Account, &s.CountryName)
		suppliers = append(suppliers, s)
	}

	return suppliers
}

func getSupplierRow(supplierId int32) Supplier {
	sqlStatement := `SELECT * FROM public.suppliers WHERE id=$1`
	row := db.QueryRow(sqlStatement, supplierId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Supplier{}
	}

	s := Supplier{}
	row.Scan(&s.Id, &s.Name, &s.Tradename, &s.FiscalName, &s.TaxId, &s.VatNumber, &s.Phone, &s.Email, &s.MainAddress, &s.Country, &s.State, &s.MainShippingAddress, &s.MainBillingAddress, &s.Language, &s.PaymentMethod, &s.BillingSeries, &s.DateCreated, &s.Account)

	return s
}

func (s *Supplier) isValid() bool {
	return !(len(s.Name) == 0 || len(s.Name) > 303 || len(s.Tradename) == 0 || len(s.Tradename) > 150 || len(s.FiscalName) == 0 || len(s.FiscalName) > 150 || len(s.TaxId) > 25 || len(s.VatNumber) > 25 || len(s.Phone) > 25 || len(s.Email) > 100)
}

func (s *Supplier) insertSupplier() bool {
	if !s.isValid() {
		return false
	}

	// prevent error in the biling serie
	if s.BillingSeries != nil && *s.BillingSeries == "" {
		s.BillingSeries = nil
	}

	// set the accounting account
	if s.Country != nil && s.Account == nil {
		s.setSupplierAccount()
	}

	sqlStatement := `INSERT INTO public.suppliers(name, tradename, fiscal_name, tax_id, vat_number, phone, email, main_address, country, state, main_shipping_address, main_billing_address, language, payment_method, billing_series, account) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	res, err := db.Exec(sqlStatement, s.Name, s.Tradename, s.FiscalName, s.TaxId, s.VatNumber, s.Phone, s.Email, s.MainAddress, s.Country, s.State, s.MainShippingAddress, s.MainBillingAddress, s.Language, s.PaymentMethod, s.BillingSeries, s.Account)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Supplier) updateSupplier() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	// prevent error in the biling serie
	if s.BillingSeries != nil && *s.BillingSeries == "" {
		s.BillingSeries = nil
	}

	// set the accounting account
	if s.Country != nil && s.Account == nil {
		s.setSupplierAccount()
	}

	sqlStatement := `UPDATE public.suppliers SET name=$2, tradename=$3, fiscal_name=$4, tax_id=$5, vat_number=$6, phone=$7, email=$8, main_address=$9, country=$10, state=$11, main_shipping_address=$12, main_billing_address=$13, language=$14, payment_method=$15, billing_series=$16, account=$17 WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id, s.Name, s.Tradename, s.FiscalName, s.TaxId, s.VatNumber, s.Phone, s.Email, s.MainAddress, s.Country, s.State, s.MainShippingAddress, s.MainBillingAddress, s.Language, s.PaymentMethod, s.BillingSeries, s.Account)
	if err != nil {
		log("DB", err.Error())
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
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findSupplierByName(languageName string) []NameInt32 {
	var customers []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.suppliers WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		log("DB", err.Error())
		return customers
	}
	for rows.Next() {
		c := NameInt32{}
		rows.Scan(&c.Id, &c.Name)
		customers = append(customers, c)
	}

	return customers
}

func getNameSupplier(id int32) string {
	sqlStatement := `SELECT name FROM public.suppliers WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

func getSupplierDefaults(customerId int32) ContactDefauls {
	sqlStatement := `SELECT main_shipping_address, (SELECT address AS main_shipping_address_name FROM address WHERE address.id = suppliers.main_shipping_address), main_billing_address, (SELECT address AS main_billing_address_name FROM address WHERE address.id = suppliers.main_billing_address), payment_method, (SELECT name AS payment_method_name FROM payment_method WHERE payment_method.id = suppliers.payment_method), billing_series, (SELECT name AS billing_series_name FROM billing_series WHERE billing_series.id = suppliers.billing_series), (SELECT currency FROM country WHERE country.id = suppliers.country), (SELECT name AS currency_name FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = suppliers.country)), (SELECT exchange FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = suppliers.country)) FROM public.suppliers WHERE id = $1`
	row := db.QueryRow(sqlStatement, customerId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ContactDefauls{}
	}
	s := ContactDefauls{}
	row.Scan(&s.MainShippingAddress, &s.MainShippingAddressName, &s.MainBillingAddress, &s.MainBillingAddressName, &s.PaymentMethod, &s.PaymentMethodName, &s.BillingSeries, &s.BillingSeriesName, &s.Currency, &s.CurrencyName, &s.CurrencyChange)
	return s
}

func getSupplierAddresses(supplierId int32) []Address {
	var addresses []Address = make([]Address, 0)
	sqlStatement := `SELECT *,CASE WHEN address.customer IS NOT NULL THEN (SELECT name FROM customer WHERE customer.id=address.customer) ELSE (SELECT name FROM suppliers WHERE suppliers.id=address.supplier) END,(SELECT name FROM country WHERE country.id=address.country),(SELECT name FROM state WHERE state.id=address.state) FROM address WHERE supplier=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, supplierId)
	if err != nil {
		log("DB", err.Error())
		return addresses
	}
	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode, &a.ShopifyId, &a.ContactName, &a.CountryName, &a.StateName)
		addresses = append(addresses, a)
	}

	return addresses
}

func getSupplierPurchaseOrders(supplierId int32) []PurchaseOrder {
	var purchases []PurchaseOrder = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_order.supplier) FROM purchase_order WHERE supplier=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, supplierId)
	if err != nil {
		log("DB", err.Error())
		return purchases
	}
	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.SupplierName)
		purchases = append(purchases, s)
	}

	return purchases
}

func (c *Supplier) setSupplierAccount() {
	sqlStatement := `SELECT un_code FROM country WHERE id=$1`
	row := db.QueryRow(sqlStatement, c.Country)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return
	}

	var unCode int16
	row.Scan(&unCode)
	if unCode <= 0 {
		return
	}

	s := getSettingsRecord()
	if s.SupplierJournal == nil {
		return
	}

	a := Account{}
	a.Journal = *s.SupplierJournal
	a.Name = c.FiscalName
	ok := a.insertAccount()
	if ok {
		c.Account = &a.Id
	}
}

type SupplierLocate struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type SupplierLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name
	Value string `json:"value"`
}

func (q *SupplierLocateQuery) locateSuppliers() []SupplierLocate {
	var suppliers []SupplierLocate = make([]SupplierLocate, 0)
	sqlStatement := ``
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		sqlStatement = `SELECT id,name FROM public.suppliers ORDER BY id ASC`
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			sqlStatement = `SELECT id,name FROM public.suppliers ORDER BY id ASC`
		} else {
			sqlStatement = `SELECT id,name FROM public.suppliers WHERE id=$1`
			parameters = append(parameters, id)
		}
	} else if q.Mode == 1 {
		sqlStatement = `SELECT id,name FROM public.suppliers WHERE name ILIKE $1 ORDER BY id ASC`
		parameters = append(parameters, "%"+q.Value+"%")
	}
	rows, err := db.Query(sqlStatement, parameters...)
	if err != nil {
		log("DB", err.Error())
		return suppliers
	}
	for rows.Next() {
		s := SupplierLocate{}
		rows.Scan(&s.Id, &s.Name)
		suppliers = append(suppliers, s)
	}

	return suppliers
}
