package main

import (
	"strconv"
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
	Country             *int32    `json:"country"`
	State               *int32    `json:"state"`
	MainShippingAddress *int32    `json:"mainShippingAddress"`
	MainBillingAddress  *int32    `json:"mainBillingAddress"`
	Language            *int32    `json:"language"`
	PaymentMethod       *int32    `json:"paymentMethod"`
	BillingSeries       *string   `json:"billingSeries"`
	DateCreated         time.Time `json:"dateCreated"`
	Account             *int32    `json:"account"`
	CountryName         *string   `json:"countryName"`
	wooCommerceId       int32
	prestaShopId        int32
	shopifyId           int64
	enterprise          int32
}

type Customers struct {
	Rows      int32      `json:"rows"`
	Customers []Customer `json:"customers"`
}

func (q *PaginationQuery) getCustomers() Customers {
	ct := Customers{}
	if !q.isValid() {
		return ct
	}

	ct.Customers = make([]Customer, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=customer.country) FROM public.customer WHERE enterprise=$3 ORDER BY id DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, q.enterprise)
	if err != nil {
		log("DB", err.Error())
		return ct
	}
	defer rows.Close()

	for rows.Next() {
		c := Customer{}
		rows.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.prestaShopId, &c.Account, &c.wooCommerceId, &c.shopifyId, &c.enterprise, &c.CountryName)
		ct.Customers = append(ct.Customers, c)
	}

	sqlStatement = `SELECT COUNT(*) FROM customer WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, q.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ct
	}
	row.Scan(&ct.Rows)

	return ct
}

func (s *PaginatedSearch) searchCustomers() Customers {
	ct := Customers{}
	if !s.isValid() {
		return ct
	}

	ct.Customers = make([]Customer, 0)
	sqlStatement := `SELECT *,(SELECT name FROM country WHERE country.id=customer.country) FROM customer WHERE (name ILIKE $1 OR tax_id ILIKE $1 OR email ILIKE $1) AND enterprise=$4 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := db.Query(sqlStatement, "%"+s.Search+"%", s.Limit, s.Offset, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		return ct
	}
	defer rows.Close()

	for rows.Next() {
		c := Customer{}
		rows.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.prestaShopId, &c.Account, &c.wooCommerceId, &c.shopifyId, &c.enterprise, &c.CountryName)
		ct.Customers = append(ct.Customers, c)
	}

	sqlStatement = `SELECT COUNT(*) FROM customer WHERE (name ILIKE $1 OR tax_id ILIKE $1 OR email ILIKE $1) AND enterprise=$2`
	row := db.QueryRow(sqlStatement, "%"+s.Search+"%", s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ct
	}
	row.Scan(&ct.Rows)

	return ct
}

func getCustomerRow(customerId int32) Customer {
	sqlStatement := `SELECT * FROM public.customer WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, customerId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Customer{}
	}

	c := Customer{}
	row.Scan(&c.Id, &c.Name, &c.Tradename, &c.FiscalName, &c.TaxId, &c.VatNumber, &c.Phone, &c.Email, &c.MainAddress, &c.Country, &c.State, &c.MainShippingAddress, &c.MainBillingAddress, &c.Language, &c.PaymentMethod, &c.BillingSeries, &c.DateCreated, &c.prestaShopId, &c.Account, &c.wooCommerceId, &c.shopifyId, &c.enterprise)

	return c
}

func (c *Customer) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 303 || len(c.Tradename) == 0 || len(c.Tradename) > 150 || len(c.FiscalName) == 0 || len(c.FiscalName) > 150 || len(c.TaxId) > 25 || len(c.VatNumber) > 25 || len(c.Phone) > 25 || len(c.Email) > 100 || (len(c.Email) > 0 && !emailIsValid(c.Email)) || (len(c.Phone) > 0 && !phoneIsValid(c.Phone)))
}

// 1 = Invalid
// 2 = Database error
func (c *Customer) insertCustomer(userId int32) OperationResult {
	if !c.isValid() {
		return OperationResult{Code: 1}
	}

	// prevent error in the biling serie
	if c.BillingSeries != nil && *c.BillingSeries == "" {
		c.BillingSeries = nil
	}

	// set the accounting account
	if c.Country != nil && c.Account == nil {
		c.setCustomerAccount()
	}

	sqlStatement := `INSERT INTO public.customer(name, tradename, fiscal_name, tax_id, vat_number, phone, email, main_address, country, state, main_shipping_address, main_billing_address, language, payment_method, billing_series, ps_id, account, wc_id, sy_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20) RETURNING id`
	row := db.QueryRow(sqlStatement, c.Name, c.Tradename, c.FiscalName, c.TaxId, c.VatNumber, c.Phone, c.Email, c.MainAddress, c.Country, c.State, c.MainShippingAddress, c.MainBillingAddress, c.Language, c.PaymentMethod, c.BillingSeries, c.prestaShopId, c.Account, c.wooCommerceId, c.shopifyId, c.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return OperationResult{Code: 2}
	}

	var customerId int32
	row.Scan(&customerId)
	c.Id = customerId

	if customerId > 0 {
		insertTransactionalLog(c.enterprise, "customer", int(c.Id), userId, "I")
	}

	return OperationResult{Id: int64(customerId)}
}

func (c *Customer) updateCustomer(userId int32) bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	// prevent error in the biling serie
	if c.BillingSeries != nil && *c.BillingSeries == "" {
		c.BillingSeries = nil
	}

	// set the accounting account
	if c.Country != nil && c.Account == nil {
		c.setCustomerAccount()
	}

	sqlStatement := `UPDATE public.customer SET name=$2, tradename=$3, fiscal_name=$4, tax_id=$5, vat_number=$6, phone=$7, email=$8, main_address=$9, country=$10, state=$11, main_shipping_address=$12, main_billing_address=$13, language=$14, payment_method=$15, billing_series=$16, account=$17 WHERE id=$1 AND enterprise=$18`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.Tradename, c.FiscalName, c.TaxId, c.VatNumber, c.Phone, c.Email, c.MainAddress, c.Country, c.State, c.MainShippingAddress, c.MainBillingAddress, c.Language, c.PaymentMethod, c.BillingSeries, c.Account, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	insertTransactionalLog(c.enterprise, "customer", int(c.Id), userId, "U")

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Customer) deleteCustomer(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	insertTransactionalLog(c.enterprise, "customer", int(c.Id), userId, "D")

	sqlStatement := `DELETE FROM public.customer WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCustomerByName(customerName string, enterpriseId int32) []NameInt32 {
	var customers []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.customer WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(customerName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return customers
	}
	defer rows.Close()

	for rows.Next() {
		c := NameInt32{}
		rows.Scan(&c.Id, &c.Name)
		customers = append(customers, c)
	}

	return customers
}

func getNameCustomer(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.customer WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
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
	PaymentMethod           *int32   `json:"paymentMethod"`
	PaymentMethodName       *string  `json:"paymentMethodName"`
	BillingSeries           *string  `json:"billingSeries"`
	BillingSeriesName       *string  `json:"billingSeriesName"`
	Currency                *int32   `json:"currency"`
	CurrencyName            *string  `json:"currencyName"`
	CurrencyChange          *float64 `json:"currencyChange"`
}

func getCustomerDefaults(customerId int32, enterpriseId int32) ContactDefauls {
	sqlStatement := `SELECT main_shipping_address, (SELECT address AS main_shipping_address_name FROM address WHERE address.id = customer.main_shipping_address), main_billing_address, (SELECT address AS main_billing_address_name FROM address WHERE address.id = customer.main_billing_address), payment_method, (SELECT name AS payment_method_name FROM payment_method WHERE payment_method.id = customer.payment_method), billing_series, (SELECT name AS billing_series_name FROM billing_series WHERE billing_series.id = customer.billing_series AND billing_series.enterprise=customer.enterprise), (SELECT currency FROM country WHERE country.id = customer.country), (SELECT name AS currency_name FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)), (SELECT exchange FROM currency WHERE currency.id = (SELECT currency FROM country WHERE country.id = customer.country)) FROM public.customer WHERE id = $1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, customerId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ContactDefauls{}
	}
	c := ContactDefauls{}
	row.Scan(&c.MainShippingAddress, &c.MainShippingAddressName, &c.MainBillingAddress, &c.MainBillingAddressName, &c.PaymentMethod, &c.PaymentMethodName, &c.BillingSeries, &c.BillingSeriesName, &c.Currency, &c.CurrencyName, &c.CurrencyChange)

	return c
}

func getCustomerAddresses(customerId int32, enterpriseId int32) []Address {
	var addresses []Address = make([]Address, 0)
	sqlStatement := `SELECT *,CASE WHEN address.customer IS NOT NULL THEN (SELECT name FROM customer WHERE customer.id=address.customer) ELSE (SELECT name FROM suppliers WHERE suppliers.id=address.supplier) END,(SELECT name FROM country WHERE country.id=address.country),(SELECT name FROM state WHERE state.id=address.state) FROM address WHERE customer=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, customerId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return addresses
	}
	defer rows.Close()

	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.prestaShopId, &a.ZipCode, &a.shopifyId, &a.enterprise, &a.ContactName, &a.CountryName, &a.StateName)
		addresses = append(addresses, a)
	}

	return addresses
}

func getCustomerSaleOrders(customerId int32, enterpriseId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order WHERE customer=$1 AND enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, customerId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return sales
	}
	defer rows.Close()

	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.shopifyDraftId, &s.enterprise, &s.CustomerName)
		sales = append(sales, s)
	}

	return sales
}

func (c *Customer) setCustomerAccount() {
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

	s := getSettingsRecordById(c.enterprise)
	if s.CustomerJournal == nil {
		return
	}

	aId := getAccountIdByAccountNumber(*s.CustomerJournal, int32(unCode), c.enterprise) // 430 -> standard for "Customers"
	if aId > 0 {
		c.Account = &aId
	}
}

type CustomerLocate struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type CustomerLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name
	Value string `json:"value"`
}

func (q *CustomerLocateQuery) locateCustomers(enterpriseId int32) []CustomerLocate {
	var customers []CustomerLocate = make([]CustomerLocate, 0)
	sqlStatement := ``
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		sqlStatement = `SELECT id,name FROM public.customer ORDER BY id ASC LIMIT 100`
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			sqlStatement = `SELECT id,name FROM public.customer ORDER BY id ASC LIMIT 100`
		} else {
			sqlStatement = `SELECT id,name FROM public.customer WHERE id=$1`
			parameters = append(parameters, id)
		}
	} else if q.Mode == 1 {
		sqlStatement = `SELECT id,name FROM public.customer WHERE name ILIKE $1 ORDER BY id ASC LIMIT 100`
		parameters = append(parameters, "%"+q.Value+"%")
	}
	rows, err := db.Query(sqlStatement, parameters...)
	if err != nil {
		log("DB", err.Error())
		return customers
	}
	for rows.Next() {
		c := CustomerLocate{}
		rows.Scan(&c.Id, &c.Name)
		customers = append(customers, c)
	}

	return customers
}
