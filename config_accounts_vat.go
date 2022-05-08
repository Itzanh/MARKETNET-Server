package main

type ConfigAccountsVat struct {
	VatPercent            float64  `json:"vatPercent" gorm:"primaryKey;type:numeric(14,6)"`
	AccountSaleId         int32    `json:"accountSaleId" gorm:"column:account_sale;not null:true"`
	AccountSale           Account  `json:"accountSale" gorm:"foreignKey:AccountSaleId,EnterpriseId;references:Id,EnterpriseId"`
	AccountPurchaseId     int32    `json:"accountPurchaseId" gorm:"column:account_purchase;not null:true"`
	AccountPurchase       Account  `json:"accountPurchase" gorm:"foreignKey:AccountPurchaseId,EnterpriseId;references:Id,EnterpriseId"`
	AccountSaleNumber     int32    `json:"accountSaleNumber" gorm:"-"`
	JournalSale           int32    `json:"journalSale" gorm:"-"`
	AccountPurchaseNumber int32    `json:"accountPurchaseNumber" gorm:"-"`
	JournalPurchase       int32    `json:"journalPurchase" gorm:"-"`
	EnterpriseId          int32    `json:"-" gorm:"primaryKey;column:enterprise"`
	Enterprise            Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *ConfigAccountsVat) TableName() string {
	return "config_accounts_vat"
}

func getConfigAccountsVat(enterpriseId int32) []ConfigAccountsVat {
	configAccountsVat := make([]ConfigAccountsVat, 0)
	dbOrm.Model(&ConfigAccountsVat{}).Where("enterprise = ?", enterpriseId).Preload("AccountSale").Preload("AccountPurchase").Order("vat_percent ASC").Find(&configAccountsVat)
	return configAccountsVat
}

// Journal, Account
func getConfigAccountsVatSaleRow(vatPercent float64, enterpriseId int32) (int32, int32) {
	var configAccountVat ConfigAccountsVat
	dbOrm.Model(&ConfigAccountsVat{}).Where("vat_percent = ? AND enterprise = ?", vatPercent, enterpriseId).Preload("AccountSale").First(&configAccountVat)
	return configAccountVat.AccountSale.JournalId, configAccountVat.AccountSale.AccountNumber
}

// Journal, Account
func getConfigAccountsVatPurchaseRow(vatPercent float64, enterpriseId int32) (int32, int32) {
	var configAccountVat ConfigAccountsVat
	dbOrm.Model(&ConfigAccountsVat{}).Where("vat_percent = ? AND enterprise = ?", vatPercent, enterpriseId).Preload("AccountPurchase").First(&configAccountVat)
	return configAccountVat.AccountPurchase.JournalId, configAccountVat.AccountPurchase.AccountNumber
}

func (c *ConfigAccountsVat) isValid() bool {
	return !(c.VatPercent <= 0 || c.AccountSaleNumber <= 0 || c.JournalSale <= 0 || c.AccountPurchaseNumber <= 0 || c.JournalPurchase <= 0)
}

func (c *ConfigAccountsVat) insertConfigAccountsVat() bool {
	if !c.isValid() {
		return false
	}

	c.AccountSaleId = getAccountIdByAccountNumber(c.JournalSale, c.AccountSaleNumber, c.EnterpriseId)
	c.AccountPurchaseId = getAccountIdByAccountNumber(c.JournalPurchase, c.AccountPurchaseNumber, c.EnterpriseId)
	if c.AccountSaleId <= 0 || c.AccountPurchaseId <= 0 {
		return false
	}

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *ConfigAccountsVat) deleteConfigAccountsVat() bool {
	result := dbOrm.Where("vat_percent = ? AND enterprise = ?", c.VatPercent, c.EnterpriseId).Delete(&ConfigAccountsVat{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
