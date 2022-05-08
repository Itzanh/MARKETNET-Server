package main

import (
	"time"

	"gorm.io/gorm"
)

type Charges struct {
	Id                               int32                    `json:"id"`
	AccountingMovementId             int64                    `json:"accountingMovementId" gorm:"column:accounting_movement;not null:true"`
	AccountingMovement               AccountingMovement       `json:"accountingMovement" gorm:"foreignkey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailDebitId  int64                    `json:"accountingMovementDetailDebitId" gorm:"column:accounting_movement_detail_debit;not null:true"`
	AccountingMovementDetailDebit    AccountingMovementDetail `json:"accountingMovementDetailDebit" gorm:"foreignkey:AccountingMovementDetailDebitId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailCreditId int64                    `json:"accountingMovementDetailCreditId" gorm:"column:accounting_movement_detail_credit;not null:true"`
	AccountingMovementDetailCredit   AccountingMovementDetail `json:"accountingMovementDetailCredit" gorm:"foreignkey:AccountingMovementDetailCreditId,EnterpriseId;references:Id,EnterpriseId"`
	AccountId                        int32                    `json:"accountId" gorm:"column:account;not null:true"`
	Account                          Account                  `json:"account" gorm:"foreignkey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated                      time.Time                `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Amount                           float64                  `json:"amount" gorm:"type:numeric(14,6);not null:true"`
	Concept                          string                   `json:"concept" gorm:"type:character varying(140);not null:true"`
	CollectionOperationId            int32                    `json:"collectionOperationId" gorm:"column:collection_operation;not null:true"`
	CollectionOperation              CollectionOperation      `json:"collectionOperation" gorm:"foreignkey:CollectionOperationId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId                     int32                    `json:"-" gorm:"column:enterprise;not null"`
	Enterprise                       Settings                 `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Charges) TableName() string {
	return "charges"
}

func getCharges(collectionOperation int32, enterpriseId int32) []Charges {
	charges := make([]Charges, 0)
	// get the charges for this enterprise and collection operation using dbOrm
	dbOrm.Model(&Charges{}).Where("collection_operation = ? AND enterprise = ?", collectionOperation, enterpriseId).Order("id ASC").Find(&charges)
	return charges
}

func getChargesRow(chargesId int32) Charges {
	c := Charges{}
	dbOrm.Model(&Charges{}).Where("id = ?", chargesId).First(&c)
	return c
}

func (c *Charges) isValid() bool {
	return !(c.CollectionOperationId <= 0 || len(c.Concept) > 140 || c.Amount <= 0)
}

func (c *Charges) BeforeCreate(tx *gorm.DB) (err error) {
	var charges Charges
	tx.Model(&Charges{}).Last(&charges)
	c.Id = charges.Id + 1
	return nil
}

func (c *Charges) insertCharges(userId int32) bool {
	// validation
	if !c.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	// get data from collection operation
	co := getColletionOperationRow(c.CollectionOperationId)
	if co.Id <= 0 || co.EnterpriseId != c.EnterpriseId || co.BankId == nil || co.Pending <= 0 {
		trans.Rollback()
		return false
	}

	c.AccountingMovementId = co.AccountingMovementId
	c.AccountingMovementDetailDebitId = co.AccountingMovementDetailId
	c.AccountId = co.AccountId

	ok := co.addQuantityCharges(c.Amount, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	am := getAccountingMovementRow(co.AccountingMovementId)
	if am.Id <= 0 {
		trans.Rollback()
		return false
	}

	// insert accounting movement for the payment
	m := AccountingMovement{}
	m.Type = "N"
	m.BillingSerieId = am.BillingSerieId
	m.EnterpriseId = c.EnterpriseId
	ok = m.insertAccountingMovement(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}

	// 1. debit detail for the bank
	bank := getAccountRow(*co.BankId)

	dInc := AccountingMovementDetail{}
	dInc.MovementId = m.Id
	dInc.JournalId = bank.JournalId
	dInc.AccountNumber = bank.AccountNumber
	dInc.Debit = c.Amount
	dInc.Type = "N"
	dInc.PaymentMethodId = co.PaymentMethodId
	dInc.EnterpriseId = c.EnterpriseId
	ok = dInc.insertAccountingMovementDetail(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}

	// 2. credit detail for the customer's account
	dCustDebit := getAccountingMovementDetailRow(c.AccountingMovementDetailDebitId)
	if dCustDebit.Id <= 0 {
		trans.Rollback()
		return false
	}

	dCust := AccountingMovementDetail{}
	dCust.MovementId = m.Id
	dCust.JournalId = dCustDebit.JournalId
	dCust.AccountNumber = dCustDebit.Account.AccountNumber
	dCust.Credit = c.Amount
	dCust.Type = "N"
	dCust.DocumentName = dCustDebit.DocumentName
	dCust.PaymentMethodId = co.PaymentMethodId
	dCust.EnterpriseId = c.EnterpriseId
	ok = dCust.insertAccountingMovementDetail(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}
	c.AccountingMovementDetailCreditId = dCust.Id

	// insert row
	c.DateCreated = time.Now()
	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if c.Id > 0 {
		insertTransactionalLog(c.EnterpriseId, "charges", int(c.Id), userId, "I")
	} else {
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	return true
	///
}

func (c *Charges) deleteCharges(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	inMemoryCharge := getChargesRow(c.Id)
	if inMemoryCharge.Id <= 0 || inMemoryCharge.EnterpriseId != c.EnterpriseId {
		trans.Rollback()
		return false
	}
	// get the collection operation
	co := getColletionOperationRow(inMemoryCharge.CollectionOperationId)
	if co.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "charges", int(c.Id), userId, "D")

	result := trans.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Charges{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// substract the paid amount
	ok := co.addQuantityCharges(-inMemoryCharge.Amount, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	// delete the associated account movement (credit)
	amd := getAccountingMovementDetailRow(inMemoryCharge.AccountingMovementDetailCreditId)
	am := getAccountingMovementRow(amd.MovementId)
	ok = am.deleteAccountingMovement(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	return true
	///
}
