package main

import (
	"gorm.io/gorm"
)

type ProductAccount struct {
	Id           int32    `json:"id"`
	ProductId    int32    `json:"productId" gorm:"column:product;not null:true;index:product_account_product_type,unique:true,priority:1"`
	Product      Product  `json:"-" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,enterprise"`
	AccountId    int32    `json:"accountId" gorm:"column:account;not null:true"`
	Account      Account  `json:"account" gorm:"foreignKey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	JournalId    int32    `json:"journal" gorm:"column:jorunal;not null:true"`
	Journal      Journal  `json:"-" gorm:"foreignKey:JournalId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId int32    `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Type         string   `json:"type" gorm:"type:character(1);not null:true;index:product_account_product_type,unique:true,priority:2"` // S = Sale, P = Purchase
}

func (pa *ProductAccount) TableName() string {
	return "product_account"
}

func getProductAccounts(productId int32, enterpriseId int32) []ProductAccount {
	var accounts []ProductAccount = make([]ProductAccount, 0)
	result := dbOrm.Model(&ProductAccount{}).Where("product_account.product = ? AND product_account.enterprise = ?", productId, enterpriseId).Preload("Account").Order("product_account.id ASC").Find(&accounts)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return accounts
}

// returns the specified account for the product, or nil if not specified by the user
// accType: S = Sales, p = Purchase
func getProductAccount(productId int32, accType string) *ProductAccount {
	if productId <= 0 || (accType != "S" && accType != "P") {
		return nil
	}

	cursor := dbOrm.Model(&ProductAccount{}).Where("product = ? AND type = ?", productId, accType)

	var rowsNumber int64
	cursor.Count(&rowsNumber)

	if rowsNumber == 0 {
		return nil
	}

	var productAccount ProductAccount
	cursor.Joins("Account").First(&productAccount)
	return &productAccount
}

func (a *ProductAccount) isValid() bool {
	return !(a.ProductId <= 0 || a.AccountId <= 0 || a.EnterpriseId <= 0 || (a.Type != "S" && a.Type != "P"))
}

func (a *ProductAccount) BeforeCreate(tx *gorm.DB) (err error) {
	var productAccount ProductAccount
	tx.Model(&ProductAccount{}).Last(&productAccount)
	a.Id = productAccount.Id + 1
	return nil
}

func (a *ProductAccount) insertProductAccount() bool {
	if !a.isValid() {
		return false
	}

	account := getAccountRow(a.AccountId)
	if account.Id <= 0 {
		return false
	}
	a.JournalId = account.JournalId

	result := dbOrm.Create(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *ProductAccount) updateProductAccount() bool {
	if !a.isValid() || a.Id <= 0 {
		return false
	}

	var productAccount ProductAccount
	result := dbOrm.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).First(&productAccount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	productAccount.AccountId = a.AccountId

	account := getAccountRow(a.AccountId)
	if account.Id <= 0 {
		return false
	}
	productAccount.JournalId = account.JournalId

	result = dbOrm.Save(&productAccount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *ProductAccount) deleteProductAccount() bool {
	if a.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).Delete(&ProductAccount{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
