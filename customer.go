package main

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Customer struct {
	Id                    int32          `json:"id" gorm:"index:customer_id_enterprise,unique:true,priority:1"`
	Name                  string         `json:"name" gorm:"type:character varying(303);not null:true;index:customer_name_trgm,type:gin"`
	Tradename             string         `json:"tradename" gorm:"type:character varying(150);not null:true"`
	FiscalName            string         `json:"fiscalName" gorm:"type:character varying(150);not null:true"`
	TaxId                 string         `json:"taxId" gorm:"type:character varying(25);not null:true;index:customer_tax_id,type:gin"`
	VatNumber             string         `json:"vatNumber" gorm:"type:character varying(25);not null:true"`
	Phone                 string         `json:"phone" gorm:"type:character varying(15);not null:true"`
	Email                 string         `json:"email" gorm:"type:character varying(150);not null:true;index:customer_email,type:gin"`
	MainAddressId         *int32         `json:"mainAddressId" gorm:"column:main_address"`
	MainAddress           *Address       `json:"mainAddress" gorm:"foreignKey:MainAddressId,EnterpriseId;references:Id,EnterpriseId"`
	CountryId             *int32         `json:"countryId" gorm:"column:country"`
	Country               *Country       `json:"country" gorm:"foreignKey:CountryId,EnterpriseId;references:Id,EnterpriseId"`
	StateId               *int32         `json:"stateId" gorm:"column:state"`
	State                 *State         `json:"state" gorm:"foreignKey:StateId,EnterpriseId;references:Id,EnterpriseId"`
	MainShippingAddressId *int32         `json:"mainShippingAddressId" gorm:"column:main_shipping_address"`
	MainShippingAddress   *Address       `json:"mainShippingAddress" gorm:"foreignKey:MainShippingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	MainBillingAddressId  *int32         `json:"mainBillingAddressId" gorm:"column:main_billing_address"`
	MainBillingAddress    *Address       `json:"mainBillingAddress" gorm:"foreignKey:MainBillingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	LanguageId            *int32         `json:"languageId" gorm:"column:language"`
	Language              *Language      `json:"language" gorm:"foreignKey:LanguageId,EnterpriseId;references:Id,EnterpriseId"`
	PaymentMethodId       *int32         `json:"paymentMethodId" gorm:"column:payment_method"`
	PaymentMethod         *PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId       *string        `json:"billingSeriesId" gorm:"type:character(3);column:billing_series"`
	BillingSeries         *BillingSerie  `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated           time.Time      `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	PrestaShopId          int32          `json:"-" gorm:"column:ps_id;not null:true;index:customer_ps_id,unique:true,priority:2,where:ps_id <> 0"`
	AccountId             *int32         `json:"accountId" gorm:"column:account"`
	Account               *Account       `json:"account" gorm:"foreignKey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	WooCommerceId         int32          `json:"-" gorm:"column:wc_id;not null:true;index:customer_wc_id,unique:true,priority:2,where:wc_id <> 0"`
	ShopifyId             int64          `json:"-" gorm:"column:sy_id;not null:true;index:customer_sy_id,unique:true,priority:2,where:sy_id <> 0"`
	EnterpriseId          int32          `json:"-" gorm:"column:enterprise;not null:true;index:customer_id_enterprise,unique:true,priority:2;index:customer_ps_id,unique:true,priority:1,where:ps_id <> 0;index:customer_wc_id,unique:true,priority:1,where:wc_id <> 0;index:customer_sy_id,unique:true,priority:1,where:sy_id <> 0"`
	Enterprise            Settings       `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Customer) TableName() string {
	return "customer"
}

type Customers struct {
	Rows      int64      `json:"rows"`
	Customers []Customer `json:"customers"`
}

func (q *PaginationQuery) getCustomers() Customers {
	ct := Customers{}
	if !q.isValid() {
		return ct
	}

	ct.Customers = make([]Customer, 0)
	// get all customers from the database using dbOrm
	dbOrm.Model(&Customer{}).Limit(int(q.Limit)).Offset(int(q.Offset)).Preload(clause.Associations).Order("customer.id DESC").Find(&ct.Customers)

	// get the total number of customers from the database using dbOrm
	dbOrm.Model(&Customer{}).Count(&ct.Rows)

	return ct
}

func (s *PaginatedSearch) searchCustomers() Customers {
	ct := Customers{}
	if !s.isValid() {
		return ct
	}

	ct.Customers = make([]Customer, 0)

	// get all customers from the database using dbOrm
	dbOrm.Model(&Customer{}).Where("(name ILIKE @search OR tax_id ILIKE @search OR email ILIKE @search) AND enterprise = @enterpriseId", sql.Named("search", "%"+s.Search+"%"), sql.Named("enterpriseId", s.enterprise)).Limit(int(s.Limit)).Offset(int(s.Offset)).Preload(clause.Associations).Order("customer.id DESC").Find(&ct.Customers)

	// get the total number of customers from the database using dbOrm
	dbOrm.Model(&Customer{}).Count(&ct.Rows)

	return ct
}

func getCustomerRow(customerId int32) Customer {
	// get a single customer from the database using dbOrm
	c := Customer{}
	result := dbOrm.Model(&Customer{}).Where("id = ?", customerId).First(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return c
}

func getCustomerEnterpriseRow(customerId int32, enterpriseId int32) Customer {
	// get a single customer from the database using dbOrm
	c := Customer{}
	result := dbOrm.Model(&Customer{}).Where("id = ? AND enterprise = ?", customerId, enterpriseId).Preload(clause.Associations).First(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return c
}

func (c *Customer) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 303 || len(c.Tradename) == 0 || len(c.Tradename) > 150 || len(c.FiscalName) == 0 || len(c.FiscalName) > 150 || len(c.TaxId) > 25 || len(c.VatNumber) > 25 || len(c.Phone) > 25 || len(c.Email) > 100 || (len(c.Email) > 0 && !emailIsValid(c.Email)) || (len(c.Phone) > 0 && !phoneIsValid(c.Phone)))
}

// set the new customer id before create in gorm
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	var customer Customer
	tx.Model(&Customer{}).Last(&customer)
	c.Id = customer.Id + 1
	return nil
}

// 1 = Invalid
// 2 = Database error
func (c *Customer) insertCustomer(userId int32) OperationResult {
	if !c.isValid() {
		return OperationResult{Code: 1}
	}

	// prevent error in the biling serie
	if c.BillingSeriesId != nil && *c.BillingSeriesId == "" {
		c.BillingSeriesId = nil
	}

	// set the accounting account
	if c.CountryId != nil && c.AccountId == nil {
		c.setCustomerAccount()
	}

	c.DateCreated = time.Now()

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OperationResult{Code: 2}
	}

	if c.Id > 0 {
		insertTransactionalLog(c.EnterpriseId, "customer", int(c.Id), userId, "I")
		json, _ := json.Marshal(c)
		go fireWebHook(c.EnterpriseId, "customer", "POST", string(json))
	}

	return OperationResult{Id: int64(c.Id)}
}

func (c *Customer) updateCustomer(userId int32) bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	// prevent error in the biling serie
	if c.BillingSeriesId != nil && *c.BillingSeriesId == "" {
		c.BillingSeriesId = nil
	}

	// set the accounting account
	if c.CountryId != nil && c.AccountId == nil {
		c.setCustomerAccount()
	}

	// get a single customer row from the database where id and enterprise are c.Id and c.EnterpriseId
	var customer Customer
	result := dbOrm.Model(&Customer{}).Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&customer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// copy all the fields from c to customer
	customer.Name = c.Name
	customer.Tradename = c.Tradename
	customer.FiscalName = c.FiscalName
	customer.TaxId = c.TaxId
	customer.VatNumber = c.VatNumber
	customer.Phone = c.Phone
	customer.Email = c.Email
	customer.MainAddressId = c.MainAddressId
	customer.CountryId = c.CountryId
	customer.StateId = c.StateId
	customer.MainShippingAddressId = c.MainShippingAddressId
	customer.MainBillingAddressId = c.MainBillingAddressId
	customer.LanguageId = c.LanguageId
	customer.PaymentMethodId = c.PaymentMethodId
	customer.BillingSeriesId = c.BillingSeriesId
	customer.AccountId = c.AccountId

	// update the customer in the database
	result = dbOrm.Save(&customer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "customer", int(c.Id), userId, "U")
	json, _ := json.Marshal(c)
	go fireWebHook(c.EnterpriseId, "customer", "PUT", string(json))

	return true
}

func (c *Customer) deleteCustomer(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "customer", int(c.Id), userId, "D")
	json, _ := json.Marshal(c)
	go fireWebHook(c.EnterpriseId, "customer", "DELETE", string(json))

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Customer{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findCustomerByName(customerName string, enterpriseId int32) []NameInt32 {
	var customers []NameInt32 = make([]NameInt32, 0)
	// get 10 customers from the database where name like customerName and enterprise is enterpriseId
	result := dbOrm.Model(&Customer{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(customerName), enterpriseId).Limit(10).Scan(&customers)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return customers
	}
	return customers
}

// Used both in customers and suppliers
type ContactDefauls struct {
	MainShippingAddress     *int32  `json:"mainShippingAddress"`
	MainShippingAddressName *string `json:"mainShippingAddressName"`
	MainBillingAddress      *int32  `json:"mainBillingAddress"`
	MainBillingAddressName  *string `json:"mainBillingAddressName"`
	PaymentMethod           *int32  `json:"paymentMethod"`
	PaymentMethodName       *string `json:"paymentMethodName"`
	BillingSeries           *string `json:"billingSeries"`
	BillingSeriesName       *string `json:"billingSeriesName"`
	Currency                *int32  `json:"currency"`
	CurrencyName            *string `json:"currencyName"`
	CurrencyChange          float64 `json:"currencyChange"`
}

func getCustomerDefaults(customerId int32, enterpriseId int32) ContactDefauls {
	// get a single customer from the database where id and enterprise are customerId and enterpriseId
	var customer Customer
	result := dbOrm.Model(&Customer{}).Where("id = ? AND enterprise = ?", customerId, enterpriseId).Preload(clause.Associations).First(&customer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ContactDefauls{}
	}

	var contactDefauls ContactDefauls = ContactDefauls{}

	contactDefauls.MainShippingAddress = customer.MainShippingAddressId
	if customer.MainShippingAddressId != nil {
		contactDefauls.MainShippingAddressName = &customer.MainShippingAddress.Address
	}

	contactDefauls.MainBillingAddress = customer.MainBillingAddressId
	if customer.MainBillingAddressId != nil {
		contactDefauls.MainBillingAddressName = &customer.MainBillingAddress.Address
	}

	contactDefauls.PaymentMethod = customer.PaymentMethodId
	if customer.PaymentMethodId != nil {
		contactDefauls.PaymentMethodName = &customer.PaymentMethod.Name
	}

	contactDefauls.BillingSeries = customer.BillingSeriesId
	if customer.BillingSeriesId != nil {
		contactDefauls.BillingSeriesName = &customer.BillingSeries.Name
	}

	if customer.CountryId != nil {
		country := getCountryRow(*customer.CountryId, enterpriseId)
		if country.Currency != nil {
			contactDefauls.Currency = &country.Currency.Id
			contactDefauls.CurrencyName = &country.Currency.Name
			contactDefauls.CurrencyChange = country.Currency.Change
		}
	}

	return contactDefauls
}

func getCustomerAddresses(customerId int32, enterpriseId int32) []Address {
	var addresses []Address = make([]Address, 0)
	// get all the addresses from the database where customer is customerId and enterprise is enterpriseId
	result := dbOrm.Model(&Address{}).Where("customer = ? AND enterprise = ?", customerId, enterpriseId).Preload(clause.Associations).Find(&addresses)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return addresses
	}
	return addresses
}

func getCustomerSaleOrders(customerId int32, enterpriseId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	// get all the sale orders from the database where customer is customerId and enterprise is enterpriseId
	result := dbOrm.Model(&SaleOrder{}).Where("customer = ? AND enterprise = ?", customerId, enterpriseId).Preload(clause.Associations).Order("date_created DESC").Find(&sales)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return sales
	}
	return sales
}

func (c *Customer) setCustomerAccount() {
	if c.CountryId != nil {
		country := getCountryRow(*c.CountryId, c.EnterpriseId)
		if country.UNCode <= 0 {
			return
		}
	}

	settings := getSettingsRecordById(c.EnterpriseId)
	if settings.CustomerJournalId == nil {
		return
	}

	a := Account{}
	a.JournalId = *settings.CustomerJournalId
	a.Name = c.FiscalName
	a.EnterpriseId = c.EnterpriseId
	ok := a.insertAccount()
	if ok {
		c.AccountId = &a.Id
	}
}

type CustomerLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name
	Value string `json:"value"`
}

func (q *CustomerLocateQuery) locateCustomers(enterpriseId int32) []NameInt32 {
	var customers []NameInt32 = make([]NameInt32, 0)
	var query string
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		query = ""
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			query = ""
		} else {
			query = `id = ?`
			parameters = append(parameters, id)
		}
	} else if q.Mode == 1 {
		query = `name LIKE ?`
		parameters = append(parameters, "%"+q.Value+"%")
	}
	dbOrm.Model(&Customer{}).Where(query, parameters...).Limit(100).Order("id ASC").Find(&customers)

	return customers
}
