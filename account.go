package main

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	Id            int32    `json:"id" gorm:"primaryKey;index:account_id_enterprise,unique:true,priority:1"`
	JournalId     int32    `json:"journalId" gorm:"not null:true;column:journal;index:account_account_number_journal,unique:true,priority:3"`
	Journal       Journal  `json:"journal" gorm:"foreignKey:JournalId,EnterpriseId;references:Id,EnterpriseId"`
	Name          string   `json:"name" gorm:"type:character varying(150);not null:true;index:account_name,type:gin"`
	Credit        float64  `json:"credit" gorm:"type:numeric(14,6);not null:true"`
	Debit         float64  `json:"debit" gorm:"type:numeric(14,6);not null:true"`
	Balance       float64  `json:"balance" gorm:"type:numeric(14,6);not null:true"`
	AccountNumber int32    `json:"accountNumber" gorm:"not null:true;index:account_account_number_journal,unique:true,priority:2"`
	AccountName   string   `json:"accountName" gorm:"not null:true;type:character(10);default:''"`
	EnterpriseId  int32    `json:"-" gorm:"column:enterprise;not null:true;index:account_id_enterprise,unique:true,priority:2;index:account_account_number_journal,unique:true,priority:1"`
	Enterprise    Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (a *Account) TableName() string {
	return "account"
}

func getAccounts(enterpriseId int32) []Account {
	accounts := make([]Account, 0)
	result := dbOrm.Model(&Account{}).Where("account.enterprise = ?", enterpriseId).Order("account.account_name ASC").Preload(clause.Associations).Find(&accounts)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return nil
	}
	return accounts
}

type AccountSearch struct {
	Journal int32  `json:"journal"`
	Search  string `json:"search"`
}

func (s *AccountSearch) searchAccounts(enterpriseId int32) []Account {
	accounts := make([]Account, 0)
	var query string
	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, "%"+s.Search+"%")
	if s.Journal <= 0 {
		query = `(name ILIKE ?) AND (enterprise = ?)`
		interfaces = append(interfaces, enterpriseId)
	} else {
		query = `(name ILIKE ?) AND (journal = ?) AND (enterprise = ?)`
		interfaces = append(interfaces, s.Journal)
		interfaces = append(interfaces, enterpriseId)
	}

	result := dbOrm.Model(&Account{}).Where(query, interfaces...).Order("account.account_name ASC").Preload(clause.Associations).Find(&accounts)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return nil
	}
	return accounts
}

func getAccountRow(accountId int32) Account {
	a := Account{}
	dbOrm.Model(&Account{}).Where("id = ?", accountId).First(&a)
	return a
}

func (a *Account) isValid() bool {
	return !(a.JournalId <= 0 || len(a.Name) == 0 || len(a.Name) > 150)
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	var account Account
	tx.Model(&Account{}).Last(&account)
	a.Id = account.Id + 1
	return nil
}

func (a *Account) setAccountName() {
	a.AccountName = fmt.Sprintf("%03d", a.JournalId) + "." + fmt.Sprintf("%06d", a.AccountNumber)
}

func (a *Account) insertAccount() bool {
	if !a.isValid() {
		return false
	}

	if a.AccountNumber <= 0 {
		a.AccountNumber = a.getNextAccountNumber()
		if a.AccountNumber <= 0 {
			return false
		}
	}

	a.Credit = 0
	a.Debit = 0
	a.Balance = 0
	a.setAccountName()

	result := dbOrm.Create(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *Account) getNextAccountNumber() int32 {
	var accountNumber int32
	var rowsNumber int64
	result := dbOrm.Model(&Account{}).Where("journal = ?", a.JournalId).Order("account_number DESC").Limit(1).Count(&rowsNumber).Select("account_number").Pluck("account_number", &accountNumber)
	if rowsNumber == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 1
	}

	return accountNumber + 1
}

func (a *Account) updateAccount() bool {
	if a.Id <= 0 || !a.isValid() || a.AccountNumber <= 0 {
		return false
	}

	var account Account
	result := dbOrm.Model(&Account{}).Where("id = ? AND enterprise=?", a.Id, a.EnterpriseId).First(&account)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	account.Name = a.Name
	account.AccountNumber = a.AccountNumber
	account.setAccountName()

	result = dbOrm.Save(&account)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *Account) deleteAccount() bool {
	if a.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).Delete(&Account{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func getAccountIdByAccountNumber(journal int32, accountNumber int32, enterpriseId int32) int32 {
	if journal <= 0 || accountNumber <= 0 || enterpriseId <= 0 {
		return 0
	}

	var accountId int32
	result := dbOrm.Model(&Account{}).Where("account_number = ? AND journal = ? AND enterprise = ?", accountNumber, journal, enterpriseId).Select("id").Limit(1).Pluck("id", &accountId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return accountId
}

type AccountLocate struct {
	Id            int32
	Journal       int32
	Name          string
	AccountNumber int32
}

func locateAccountForCustomer(enterpriseId int32) []NameInt32 {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	names := make([]NameInt32, 0)
	dbOrm.Model(&Account{}).Where("journal = ? AND enterprise = ?", s.CustomerJournalId, enterpriseId).Order("account_number ASC").Find(&accounts)

	var account AccountLocate
	for i := 0; i < len(accounts); i++ {
		account = accounts[i]
		names = append(names, NameInt32{
			Id:   account.Id,
			Name: fmt.Sprintf("%d.%06d - %s", account.Journal, account.AccountNumber, account.Name),
		})
	}

	return names
}

func locateAccountForSupplier(enterpriseId int32) []NameInt32 {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	names := make([]NameInt32, 0)
	dbOrm.Model(&Account{}).Where("journal = ? AND enterprise = ?", s.SupplierJournalId, enterpriseId).Order("account_number ASC").Find(&accounts)

	var account AccountLocate
	for i := 0; i < len(accounts); i++ {
		account = accounts[i]
		names = append(names, NameInt32{
			Id:   account.Id,
			Name: fmt.Sprintf("%d.%06d - %s", account.Journal, account.AccountNumber, account.Name),
		})
	}

	return names
}

func locateAccountForSales(enterpriseId int32) []NameInt32 {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	names := make([]NameInt32, 0)
	result := dbOrm.Model(&Account{}).Where("journal = ? AND enterprise = ?", s.SalesJournalId, enterpriseId).Order("account_number ASC").Find(&accounts)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	var account AccountLocate
	for i := 0; i < len(accounts); i++ {
		account = accounts[i]
		names = append(names, NameInt32{
			Id:   account.Id,
			Name: fmt.Sprintf("%d.%06d - %s", account.Journal, account.AccountNumber, account.Name),
		})
	}

	return names
}

func locateAccountForPurchases(enterpriseId int32) []NameInt32 {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	names := make([]NameInt32, 0)
	result := dbOrm.Model(&Account{}).Where("journal = ? AND enterprise = ?", s.PurchaseJournalId, enterpriseId).Order("account_number ASC").Find(&accounts)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	var account AccountLocate
	for i := 0; i < len(accounts); i++ {
		account = accounts[i]
		names = append(names, NameInt32{
			Id:   account.Id,
			Name: fmt.Sprintf("%d.%06d - %s", account.Journal, account.AccountNumber, account.Name),
		})
	}

	return names
}

func locateAccountForBanks(enterpriseId int32) []NameInt32 {
	accounts := make([]AccountLocate, 0)
	names := make([]NameInt32, 0)

	var journals []Journal = make([]Journal, 0)
	dbOrm.Model(&Journal{}).Where("type = 'B'").Find(&journals)

	for i := 0; i < len(journals); i++ {
		result := dbOrm.Model(&Account{}).Where("journal = ? AND enterprise = ?", journals[i].Id, enterpriseId).Order("account_number ASC").Find(&accounts)
		if result.Error != nil {
			log("DB", result.Error.Error())
		}

		var account AccountLocate
		for i := 0; i < len(accounts); i++ {
			account = accounts[i]
			names = append(names, NameInt32{
				Id:   account.Id,
				Name: fmt.Sprintf("%d.%06d - %s", account.Journal, account.AccountNumber, account.Name),
			})
		}
	}

	return names
}
