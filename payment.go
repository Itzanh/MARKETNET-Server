package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Payment struct {
	Id                               int32                    `json:"id"`
	AccountingMovementId             int64                    `json:"accountingMovementId" gorm:"column:accounting_movement;not null:true"`
	AccountingMovement               AccountingMovement       `json:"-" gorm:"foreignkey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailDebitId  int64                    `json:"accountingMovementDetailDebitId" gorm:"column:accounting_movement_detail_debit;not null:true"`
	AccountingMovementDetailDebit    AccountingMovementDetail `json:"-" gorm:"foreignkey:AccountingMovementDetailDebitId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailCreditId int64                    `json:"accountingMovementDetailCreditId" gorm:"column:accounting_movement_detail_credit;not null:true"`
	AccountingMovementDetailCredit   AccountingMovementDetail `json:"-" gorm:"foreignkey:AccountingMovementDetailCreditId,EnterpriseId;references:Id,EnterpriseId"`
	AccountId                        int32                    `json:"accountId" gorm:"column:account;not null:true"`
	Account                          Account                  `json:"account" gorm:"foreignkey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated                      time.Time                `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Amount                           float64                  `json:"amount" gorm:"type:numeric(14,6);not null:true"`
	Concept                          string                   `json:"concept" gorm:"type:character varying(140);not null:true"`
	PaymentTransactionId             int32                    `json:"paymentTransactionId" gorm:"column:payment_transaction;not null:true"`
	PaymentTransaction               PaymentTransaction       `json:"-" gorm:"foreignkey:PaymentTransactionId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId                     int32                    `json:"-" gorm:"column:enterprise;not null"`
	Enterprise                       Settings                 `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *Payment) TableName() string {
	return "payments"
}

func getPayments(paymentTransaction int32, enterpriseId int32) []Payment {
	payments := make([]Payment, 0)
	// get payments for this enterprise and payment transacion
	dbOrm.Model(&Payment{}).Where("payment_transaction = ? AND enterprise = ?", paymentTransaction, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&payments)
	return payments
}

func getPaymentsRow(chargesId int32) Payment {
	p := Payment{}
	dbOrm.Model(&Payment{}).Where("id = ?", chargesId).First(&p)
	return p
}

func (c *Payment) isValid() bool {
	return !(c.PaymentTransactionId <= 0 || len(c.Concept) > 140 || c.Amount <= 0)
}

func (c *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	var payment Payment
	tx.Model(&Payment{}).Last(&payment)
	c.Id = payment.Id + 1
	return nil
}

func (c *Payment) insertPayment(userId int32) bool {
	// validation
	if !c.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	// get data from payment transaction
	pt := getPaymentTransactionRow(c.PaymentTransactionId)
	if pt.Id <= 0 || pt.BankId == nil || pt.EnterpriseId != c.EnterpriseId || pt.Pending <= 0 {
		trans.Rollback()
		return false
	}

	c.AccountingMovementId = pt.AccountingMovementId
	c.AccountingMovementDetailDebitId = pt.AccountingMovementDetailId
	c.AccountId = pt.AccountId

	ok := pt.addQuantityCharges(c.Amount, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	am := getAccountingMovementRow(pt.AccountingMovementId)
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
	bank := getAccountRow(*pt.BankId)

	dInc := AccountingMovementDetail{}
	dInc.MovementId = m.Id
	dInc.JournalId = bank.JournalId
	dInc.AccountNumber = bank.AccountNumber
	dInc.Credit = c.Amount
	dInc.Type = "N"
	dInc.PaymentMethodId = pt.PaymentMethodId
	dInc.EnterpriseId = c.EnterpriseId
	ok = dInc.insertAccountingMovementDetail(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}

	// 2. credit detail for the suppliers's account
	dSuppDebit := getAccountingMovementDetailRow(c.AccountingMovementDetailDebitId)
	if dSuppDebit.Id <= 0 {
		trans.Rollback()
		return false
	}

	Supp := AccountingMovementDetail{}
	Supp.MovementId = m.Id
	Supp.JournalId = dSuppDebit.JournalId
	Supp.AccountNumber = dSuppDebit.Account.AccountNumber
	Supp.Debit = c.Amount
	Supp.Type = "N"
	Supp.DocumentName = dSuppDebit.DocumentName
	Supp.PaymentMethodId = pt.PaymentMethodId
	Supp.EnterpriseId = c.EnterpriseId
	ok = Supp.insertAccountingMovementDetail(userId, trans)
	if !ok {
		trans.Rollback()
		return false
	}
	c.AccountingMovementDetailCreditId = Supp.Id

	// insert row
	c.DateCreated = time.Now()
	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "payments", int(c.Id), userId, "I")

	///
	trans.Commit()
	return true
	///
}

func (c *Payment) deletePayment(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	inMemoryPayment := getPaymentsRow(c.Id)
	if inMemoryPayment.Id <= 0 || inMemoryPayment.EnterpriseId != c.EnterpriseId {
		trans.Rollback()
		return false
	}
	// get the payment transaction
	pt := getPaymentTransactionRow(inMemoryPayment.PaymentTransactionId)
	if pt.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "payments", int(c.Id), userId, "D")

	result := trans.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Payment{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// substract the paid amount
	ok := pt.addQuantityCharges(-inMemoryPayment.Amount, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	// delete the associated account movement (credit)
	amd := getAccountingMovementDetailRow(inMemoryPayment.AccountingMovementDetailCreditId)
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
