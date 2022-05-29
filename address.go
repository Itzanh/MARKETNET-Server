package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Address struct {
	Id                int32     `json:"id" gorm:"index:address_id_enterprise,unique:true,priority:1"`
	CustomerId        *int32    `json:"customerId" gorm:"column:customer"`
	Customer          *Customer `json:"customer" gorm:"foreignKey:CustomerId,EnterpriseId;references:Id,EnterpriseId"`
	Address           string    `json:"address" gorm:"type:character varying(200);not null:true;index:address_name,type:gin"`
	Address2          string    `json:"address2" gorm:"column:address_2;type:character varying(200);not null:true"`
	StateId           *int32    `json:"stateId" gorm:"column:state"`
	State             *State    `json:"state" gorm:"foreignKey:StateId,EnterpriseId;references:Id,EnterpriseId"`
	City              string    `json:"city" gorm:"type:character varying(100);not null:true"`
	CountryId         int32     `json:"countryId" gorm:"column:country;not null:true"`
	Country           Country   `json:"country" gorm:"foreignKey:CountryId,EnterpriseId;references:Id,EnterpriseId"`
	PrivateOrBusiness string    `json:"privateOrBusiness" gorm:"column:private_business;type:character(1);not null:true"` // P = Private, B = Business, _ = Not specified
	Notes             string    `json:"notes" gorm:"type:text;not null:true"`
	SupplierId        *int32    `json:"supplierId" gorm:"column:supplier"`
	Supplier          *Supplier `json:"supplier" gorm:"foreignKey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopId      int32     `json:"-" gorm:"column:ps_id;not null;index:address_ps_id,unique:true,priority:2,where:ps_id <> 0"`
	ZipCode           string    `json:"zipCode" gorm:"type:character varying(12);not null:true"`
	ShopifyId         int64     `json:"-" gorm:"column:sy_id;not null;index:address_sy_id,unique:true,priority:2,where:sy_id <> 0"`
	EnterpriseId      int32     `json:"-" gorm:"column:enterprise;not null:true;index:address_id_enterprise,unique:true,priority:2;index:address_ps_id,unique:true,priority:1,where:ps_id <> 0;index:address_sy_id,unique:true,priority:1,where:sy_id <> 0"`
	Enterprise        Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (a *Address) TableName() string {
	return "address"
}

type Addresses struct {
	Rows      int64     `json:"rows"`
	Addresses []Address `json:"addresses"`
}

func (q *PaginationQuery) getAddresses() Addresses {
	ad := Addresses{}
	if !q.isValid() {
		return ad
	}

	ad.Addresses = make([]Address, 0)
	result := dbOrm.Model(&Address{}).Where("address.enterprise = ?", q.enterprise).Preload(clause.Associations).Offset(int(q.Offset)).Limit(int(q.Limit)).Order("id ASC").Find(&ad.Addresses)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ad
	}
	ad.Rows = int64(len(ad.Addresses))
	return ad
}

func getAddressRow(addressId int32) Address {
	a := Address{}
	result := dbOrm.Model(&Address{}).Where("id = ?", addressId).Preload(clause.Associations).First(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return Address{}
	}
	return a
}

func (s *PaginatedSearch) searchAddresses() Addresses {
	ad := Addresses{}
	if !s.isValid() {
		return ad
	}

	ad.Addresses = make([]Address, 0)
	result := dbOrm.Model(&Address{}).Where(`(address ILIKE @search OR "Customer".name ILIKE @search OR "Supplier".name ILIKE @search) AND address.enterprise = @enterpriseId`, sql.Named("enterpriseId", s.enterprise), sql.Named("search", "%"+s.Search+"%")).Joins("Customer").Joins("Supplier").Preload(clause.Associations).Offset(int(s.Offset)).Limit(int(s.Limit)).Order("id ASC").Find(&ad.Addresses)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ad
	}

	// get the row count of addresses in the database using dbOrm
	result = dbOrm.Model(&Address{}).Where("address.enterprise = @enterpriseId", sql.Named("enterpriseId", s.enterprise)).Count(&ad.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ad
	}

	return ad
}

func (a *Address) isValid() bool {
	return !((a.CustomerId == nil && a.SupplierId == nil) || (a.CustomerId != nil && *a.CustomerId <= 0) || (a.SupplierId != nil && *a.SupplierId <= 0) || len(a.Address) == 0 || len(a.Address) > 200 || len(a.Address2) > 200 || len(a.City) == 0 || len(a.City) > 100 || a.CountryId <= 0 || (a.PrivateOrBusiness != "P" && a.PrivateOrBusiness != "B" && a.PrivateOrBusiness != "_") || len(a.Notes) > 1000 || len(a.ZipCode) > 12)
}

func (a *Address) BeforeCreate(tx *gorm.DB) (err error) {
	var address Address
	tx.Model(&Address{}).Last(&address)
	a.Id = address.Id + 1
	return nil
}

// 1 = Invalid
// 2 = Database error
func (a *Address) insertAddress(userId int32) OperationResult {
	if !a.isValid() {
		return OperationResult{Code: 1}
	}

	result := dbOrm.Create(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OperationResult{Code: 2}
	}

	if a.CustomerId != nil {
		c := getCustomerRow(*a.CustomerId)
		var ok bool = false
		if c.MainAddressId == nil {
			c.MainAddressId = &a.Id
			ok = true
		}
		if ok {
			c.updateCustomer(userId)
		}
	}
	if a.SupplierId != nil {
		s := getSupplierRow(*a.SupplierId)
		var ok bool = false
		if s.MainAddressId == nil {
			s.MainAddressId = &a.Id
			ok = true
		}
		if ok {
			s.updateSupplier(userId)
		}
	}

	return OperationResult{Id: int64(a.Id)}
}

func (a *Address) updateAddress() bool {
	if a.Id <= 0 || !a.isValid() {
		fmt.Println("INVALID")
		cAddress, _ := json.Marshal(a)
		fmt.Println(string(cAddress))
		return false
	}

	var address Address
	result := dbOrm.Model(&Address{}).Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).First(&address)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	address.CustomerId = a.CustomerId
	address.SupplierId = a.SupplierId
	address.Address = a.Address
	address.Address2 = a.Address2
	address.City = a.City
	address.StateId = a.StateId
	address.CountryId = a.CountryId
	address.PrivateOrBusiness = a.PrivateOrBusiness
	address.Notes = a.Notes
	address.ZipCode = a.ZipCode

	result = dbOrm.Save(&address)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *Address) deleteAddress() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE customer SET main_address=NULL WHERE main_address=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.EnterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE customer SET main_billing_address=NULL WHERE main_billing_address=$1 AND enterprise=$2`
	_, err = db.Exec(sqlStatement, a.Id, a.EnterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE customer SET main_shipping_address=NULL WHERE main_shipping_address=$1 AND enterprise=$2`
	_, err = db.Exec(sqlStatement, a.Id, a.EnterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `DELETE FROM address WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, a.Id, a.EnterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type AddressLocate struct {
	Id      int32  `json:"id"`
	Address string `json:"address"`
}

func locateAddressByCustomer(customerId int32, enterpriseId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	result := dbOrm.Model(&Address{}).Where("customer = ? AND enterprise = ?", customerId, enterpriseId).Select("id, address").Order("id ASC").Find(&addresses)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return addresses
}

func locateAddressBySupplier(supplierId int32, enterpriseId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	result := dbOrm.Model(&Address{}).Where("supplier = ? AND enterprise = ?", supplierId, enterpriseId).Select("id, address").Order("id ASC").Find(&addresses)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return addresses
}
