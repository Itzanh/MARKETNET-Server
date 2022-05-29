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

type Supplier struct {
	Id                    int32          `json:"id" gorm:"index:suppliers_id_enterprise,unique:true,priority:1"`
	Name                  string         `json:"name" gorm:"type:character varying(303);not null:true;index:supplier_name_trgm,type:gin"`
	Tradename             string         `json:"tradename" gorm:"type:character varying(150);not null:true"`
	FiscalName            string         `json:"fiscalName" gorm:"type:character varying(150);not null:true"`
	TaxId                 string         `json:"taxId" gorm:"type:character varying(25);not null:true;index:supplier_tax_id,type:gin"`
	VatNumber             string         `json:"vatNumber" gorm:"type:character varying(25);not null:true"`
	Phone                 string         `json:"phone" gorm:"type:character varying(15);not null:true"`
	Email                 string         `json:"email" gorm:"type:character varying(150);not null:true;index:supplier_email,type:gin"`
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
	AccountId             *int32         `json:"accountId" gorm:"column:account"`
	Account               *Account       `json:"account" gorm:"foreignKey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId          int32          `json:"-" gorm:"column:enterprise;not null:true;index:suppliers_id_enterprise,unique:true,priority:2"`
	Enterprise            Settings       `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *Supplier) TableName() string {
	return "suppliers"
}

func getSuppliers(enterpriseId int32) []Supplier { // .Joins("State")
	var suppliers []Supplier = make([]Supplier, 0)
	result := dbOrm.Model(&Supplier{}).Where("suppliers.enterprise = ?", enterpriseId).Preload(clause.Associations).Order("suppliers.id ASC").Find(&suppliers)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return suppliers
}

func searchSuppliers(search string, enterpriseId int32) []Supplier {
	var suppliers []Supplier = make([]Supplier, 0)
	result := dbOrm.Model(&Supplier{}).Where("(suppliers.name ILIKE @search OR suppliers.tax_id ILIKE @search OR suppliers.email ILIKE @search) AND suppliers.enterprise = @enterpriseId", sql.Named("search", "%"+search+"%"), sql.Named("enterpriseId", enterpriseId)).Preload(clause.Associations).Order("suppliers.id ASC").Find(&suppliers)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return suppliers
}

func getSupplierRow(supplierId int32) Supplier {
	s := Supplier{}
	result := dbOrm.Model(&Supplier{}).Where("id = ?", supplierId).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

func getSupplierEnterpriseRow(supplierId int32, enterpriseId int32) Supplier {
	s := Supplier{}
	result := dbOrm.Model(&Supplier{}).Where("id = ? AND enterprise = ?", supplierId, enterpriseId).Preload(clause.Associations).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

func (s *Supplier) isValid() bool {
	return !(len(s.Name) == 0 || len(s.Name) > 303 || len(s.Tradename) == 0 || len(s.Tradename) > 150 || len(s.FiscalName) == 0 || len(s.FiscalName) > 150 || len(s.TaxId) > 25 || len(s.VatNumber) > 25 || len(s.Phone) > 25 || len(s.Email) > 100 || (len(s.Email) > 0 && !emailIsValid(s.Email)) || (len(s.Phone) > 0 && !phoneIsValid(s.Phone)))
}

func (s *Supplier) BeforeCreate(tx *gorm.DB) (err error) {
	var supplier Supplier
	tx.Model(&Supplier{}).Last(&supplier)
	s.Id = supplier.Id + 1
	return nil
}

// 1 = Invalid
// 2 = Database error
func (s *Supplier) insertSupplier(userId int32) OperationResult {
	if !s.isValid() {
		return OperationResult{Code: 1}
	}

	// prevent error in the biling serie
	if s.BillingSeriesId != nil && *s.BillingSeriesId == "" {
		s.BillingSeriesId = nil
	}

	// set the accounting account
	if s.CountryId != nil && s.AccountId == nil {
		s.setSupplierAccount()
	}

	s.DateCreated = time.Now()

	result := dbOrm.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OperationResult{Code: 2}
	}

	insertTransactionalLog(s.EnterpriseId, "suppliers", int(s.Id), userId, "I")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "suppliers", "POST", string(json))

	return OperationResult{Id: int64(s.Id)}
}

func (s *Supplier) updateSupplier(userId int32) bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	// prevent error in the biling serie
	if s.BillingSeriesId != nil && *s.BillingSeriesId == "" {
		s.BillingSeriesId = nil
	}

	// set the accounting account
	if s.CountryId != nil && s.AccountId == nil {
		s.setSupplierAccount()
	}

	var supplier Supplier
	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).First(&supplier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	supplier.Name = s.Name
	supplier.Tradename = s.Tradename
	supplier.FiscalName = s.FiscalName
	supplier.TaxId = s.TaxId
	supplier.VatNumber = s.VatNumber
	supplier.Phone = s.Phone
	supplier.Email = s.Email
	supplier.MainAddress = s.MainAddress
	supplier.CountryId = s.CountryId
	supplier.StateId = s.StateId
	supplier.MainShippingAddressId = s.MainShippingAddressId
	supplier.MainBillingAddressId = s.MainBillingAddressId
	supplier.LanguageId = s.LanguageId
	supplier.PaymentMethodId = s.PaymentMethodId
	supplier.BillingSeriesId = s.BillingSeriesId
	supplier.AccountId = s.AccountId

	result = dbOrm.Save(&supplier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(s.EnterpriseId, "suppliers", int(s.Id), userId, "U")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "suppliers", "PUT", string(json))

	return true
}

func (s *Supplier) deleteSupplier(userId int32) bool {
	if s.Id <= 0 {
		return false
	}

	insertTransactionalLog(s.EnterpriseId, "suppliers", int(s.Id), userId, "D")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "suppliers", "DELETE", string(json))

	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).Delete(&Supplier{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findSupplierByName(languageName string, enterpriseId int32) []NameInt32 {
	var suppliers []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&Supplier{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(languageName), enterpriseId).Order("id ASC").Limit(10).Find(&suppliers)
	return suppliers
}

func getSupplierDefaults(customerId int32, enterpriseId int32) ContactDefauls {
	// get a single supplier from the database where id and enterprise are customerId and enterpriseId
	var supplier Supplier
	result := dbOrm.Model(&Supplier{}).Where("id = ? AND enterprise = ?", customerId, enterpriseId).Preload(clause.Associations).First(&supplier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ContactDefauls{}
	}

	var contactDefauls ContactDefauls = ContactDefauls{}

	contactDefauls.MainShippingAddress = supplier.MainShippingAddressId
	if supplier.MainShippingAddressId != nil {
		contactDefauls.MainShippingAddressName = &supplier.MainShippingAddress.Address
	}

	contactDefauls.MainBillingAddress = supplier.MainBillingAddressId
	if supplier.MainBillingAddressId != nil {
		contactDefauls.MainBillingAddressName = &supplier.MainBillingAddress.Address
	}

	contactDefauls.PaymentMethod = supplier.PaymentMethodId
	if supplier.PaymentMethodId != nil {
		contactDefauls.PaymentMethodName = &supplier.PaymentMethod.Name
	}

	contactDefauls.BillingSeries = supplier.BillingSeriesId
	if supplier.BillingSeriesId != nil {
		contactDefauls.BillingSeriesName = &supplier.BillingSeries.Name
	}

	if supplier.CountryId != nil {
		country := getCountryRow(*supplier.CountryId, enterpriseId)
		if country.Currency != nil {
			contactDefauls.Currency = &country.Currency.Id
			contactDefauls.CurrencyName = &country.Currency.Name
			contactDefauls.CurrencyChange = country.Currency.Change
		}
	}

	return contactDefauls
}

func getSupplierAddresses(supplierId int32, enterpriseId int32) []Address {
	var addresses []Address = make([]Address, 0)
	dbOrm.Model(&Address{}).Where("supplier = ? AND enterprise = ?", supplierId, enterpriseId).Preload(clause.Associations).Order("address.id ASC").Find(&addresses)
	return addresses
}

func getSupplierPurchaseOrders(supplierId int32, enterpriseId int32) []PurchaseOrder {
	var purchases []PurchaseOrder = make([]PurchaseOrder, 0)
	dbOrm.Model(&PurchaseOrder{}).Where("supplier = ? AND enterprise = ?", supplierId, enterpriseId).Preload(clause.Associations).Order("id ASC").Find(&purchases)
	return purchases
}

func (s *Supplier) setSupplierAccount() {
	if s.CountryId != nil {
		c := getCountryRow(*s.CountryId, s.EnterpriseId)
		if c.UNCode <= 0 {
			return
		}
	}

	settings := getSettingsRecordById(s.EnterpriseId)
	if settings.SupplierJournalId == nil {
		return
	}

	a := Account{}
	a.JournalId = *settings.SupplierJournalId
	a.Name = s.FiscalName
	a.EnterpriseId = s.EnterpriseId
	ok := a.insertAccount()
	if ok {
		s.AccountId = &a.Id
	}
}

type SupplierLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name
	Value string `json:"value"`
}

func (q *SupplierLocateQuery) locateSuppliers(enterpriseId int32) []NameInt32 {
	var suppliers []NameInt32 = make([]NameInt32, 0)
	var query string
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		query = `enterprise = ?`
		parameters = append(parameters, enterpriseId)
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			query = `enterprise = ?`
			parameters = append(parameters, enterpriseId)
		} else {
			query = `id = ? AND enterprise = ?`
			parameters = append(parameters, id)
			parameters = append(parameters, enterpriseId)
		}
	} else if q.Mode == 1 {
		query = `name ILIKE ? AND enterprise = ?`
		parameters = append(parameters, "%"+q.Value+"%")
		parameters = append(parameters, enterpriseId)
	}
	dbOrm.Model(&Supplier{}).Where(query, parameters...).Limit(100).Order("id ASC").Find(&suppliers)
	return suppliers
}
