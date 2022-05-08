package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountingMovementDetail struct {
	Id              int64              `json:"id" gorm:"index:accounting_movement_detail_id_enterprise,unique:true,priority:1"`
	MovementId      int64              `json:"movementId" gorm:"column:movement;not null:true"`
	Movement        AccountingMovement `json:"movement" gorm:"foreignkey:MovementId,EnterpriseId;references:Id,EnterpriseId"`
	JournalId       int32              `json:"journalId" gorm:"column:journal;not null:true"`
	Journal         Journal            `json:"journal" gorm:"foreignkey:JournalId,EnterpriseId;references:Id,EnterpriseId"`
	AccountId       int32              `json:"accountId" gorm:"column:account;not null:true"`
	Account         Account            `json:"account" gorm:"foreignkey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	AccountNumber   int32              `json:"accountNumber" gorm:"-"`
	DateCreated     time.Time          `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Credit          float64            `json:"credit" gorm:"type:numeric(14,6);not null:true"`
	Debit           float64            `json:"debit" gorm:"type:numeric(14,6);not null:true"`
	Type            string             `json:"type" gorm:"type:character(1);not null:true"` // O: Opening, N: Normal, V: Variation of existences, R: Regularisation, C: Closing
	Note            string             `json:"note" gorm:"type:character varying(300);not null:true"`
	DocumentName    string             `json:"documentName" gorm:"type:character(15);not null:true;index:accounting_movement_detail_document_name,priority:1"`
	PaymentMethodId int32              `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod   PaymentMethod      `json:"paymentMethod" gorm:"foreignkey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId    int32              `json:"-" gorm:"column:enterprise;not null:true;index:accounting_movement_detail_id_enterprise,unique:true,priority:2"`
	Enterprise      Settings           `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (amd *AccountingMovementDetail) TableName() string {
	return "accounting_movement_detail"
}

func getAccountingMovementDetail(movementId int64, enterpriseId int32) []AccountingMovementDetail {
	accountingMovementDetail := make([]AccountingMovementDetail, 0)
	// get all accounting movement details from the database for the given movement and enterprise using dbOrm
	dbOrm.Where("movement = ? AND enterprise = ?", movementId, enterpriseId).Preload(clause.Associations).Order("id ASC").Find(&accountingMovementDetail)
	return accountingMovementDetail
}

func getAccountingMovementDetailRow(detailtId int64) AccountingMovementDetail {
	a := AccountingMovementDetail{}
	// get a single account movement detail from the database using dbOrm
	dbOrm.Where("id = ?", detailtId).Preload(clause.Associations).First(&a)
	return a
}

func (a *AccountingMovementDetail) isValid() bool {
	return !(a.MovementId <= 0 || a.JournalId <= 0 || a.AccountNumber <= 0 || (a.Credit == 0 && a.Debit == 0) || (a.Type != "O" && a.Type != "N" && a.Type != "V" && a.Type != "R" && a.Type != "C") || len(a.Note) > 300 || len(a.DocumentName) > 15 || a.PaymentMethodId <= 0)
}

func (amd *AccountingMovementDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var accountingMovementDetail AccountingMovementDetail
	tx.Model(&AccountingMovementDetail{}).Last(&accountingMovementDetail)
	amd.Id = accountingMovementDetail.Id + 1
	return nil
}

func (a *AccountingMovementDetail) insertAccountingMovementDetail(userId int32, trans *gorm.DB) bool {
	if !a.isValid() {
		return false
	}

	a.DateCreated = time.Now()

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	a.AccountId = getAccountIdByAccountNumber(a.JournalId, a.AccountNumber, a.EnterpriseId)
	if a.AccountId <= 0 {
		trans.Rollback()
		return false
	}

	/*// Round float to 2 decimal places (round to nearest)
	a.Credit = float64(math.Round(float64(a.Credit)*100) / 100)
	a.Debit = float64(math.Round(float64(a.Debit)*100) / 100)*/

	result := trans.Create(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(a.EnterpriseId, "accounting_movement_detail", int(a.Id), userId, "I")

	m := getAccountingMovementRow(a.MovementId)

	m.AmountCredit += a.Credit
	m.AmountDebit += a.Debit

	result = trans.Save(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	acc := getAccountRow(a.AccountId)
	acc.Debit = a.Debit
	acc.Credit = a.Credit

	result = trans.Save(&acc)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false
		}
		///
	}
	return true
}

func (a *AccountingMovementDetail) deleteAccountingMovementDetail(userId int32, trans *gorm.DB) bool {
	if a.Id <= 0 {
		return false
	}

	accountingMovementDetailInMemory := getAccountingMovementDetailRow(a.Id)
	settings := getSettingsRecordById(a.EnterpriseId)
	if accountingMovementDetailInMemory.Id <= 0 || accountingMovementDetailInMemory.EnterpriseId != a.EnterpriseId || (settings.LimitAccountingDate != nil && accountingMovementDetailInMemory.DateCreated.Before(*settings.LimitAccountingDate)) {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	inMemoryDetail := getAccountingMovementDetailRow(a.Id)
	if inMemoryDetail.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(a.EnterpriseId, "accounting_movement_detail", int(a.Id), userId, "I")

	result := trans.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).Delete(&AccountingMovementDetail{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	m := getAccountingMovementRow(inMemoryDetail.MovementId)

	m.AmountCredit -= inMemoryDetail.Credit
	m.AmountDebit -= inMemoryDetail.Debit

	result = trans.Save(&m)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	acc := getAccountRow(inMemoryDetail.AccountId)

	acc.Credit -= inMemoryDetail.Credit
	acc.Debit -= inMemoryDetail.Debit

	result = trans.Save(&acc)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false
		}
		///
	}
	return true
}
