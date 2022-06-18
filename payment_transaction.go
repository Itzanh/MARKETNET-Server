package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentTransaction struct {
	Id                         int32                    `json:"id" gorm:"index:payment_transaction_id_enterprise,unique:true,priority:1"`
	AccountingMovementId       int64                    `json:"accountingMovementId" gorm:"column:accounting_movement;not null:true"`
	AccountingMovement         AccountingMovement       `json:"accountingMovement" gorm:"foreignkey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailId int64                    `json:"accountingMovementDetailId" gorm:"column:accounting_movement_detail;not null:true"`
	AccountingMovementDetail   AccountingMovementDetail `json:"accountingMovementDetail" gorm:"foreignkey:AccountingMovementDetailId,EnterpriseId;references:Id,EnterpriseId"`
	AccountId                  int32                    `json:"accountId" gorm:"column:account;not null:true"`
	Account                    Account                  `json:"account" gorm:"foreignkey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	BankId                     *int32                   `json:"bankId" gorm:"column:bank"`
	Bank                       *Account                 `json:"bank" gorm:"foreignkey:BankId,EnterpriseId;references:Id,EnterpriseId"`
	Status                     string                   `json:"status" gorm:"type:character(1);not null:true"` // P = Pending, C = Paid, U = Unpaid
	DateCreated                time.Time                `json:"dateCreated" gorm:"type:timestamp(3) without time zone;not null:true"`
	DateExpiration             time.Time                `json:"dateExpiration" gorm:"type:timestamp(3) without time zone;not null:true"`
	Total                      float64                  `json:"total" gorm:"type:numeric(14,6);not null:true"`
	Paid                       float64                  `json:"paid" golorm:"type:numeric(14,6);not null:true"`
	Pending                    float64                  `json:"pending" gorm:"type:numeric(14,6);not null:true"`
	DocumentName               string                   `json:"documentName" gorm:"not null:true"`
	PaymentMethodId            int32                    `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod              PaymentMethod            `json:"paymentMethod" gorm:"foreignkey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId               int32                    `json:"-" gorm:"column:enterprise;not null:true;index:payment_transaction_id_enterprise,unique:true,priority:2"`
	Enterprise                 Settings                 `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	SupplierName               *string                  `json:"supplierName" gorm:"-"`
}

func (pt *PaymentTransaction) TableName() string {
	return "payment_transaction"
}

func getPendingPaymentTransaction(enterpriseId int32) []PaymentTransaction {
	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)
	// get pending payment transactions for the current enterprise where the status is 'P' using dbOrm
	result := dbOrm.Model(&PaymentTransaction{}).Where("payment_transaction.status = ? AND payment_transaction.enterprise = ?", "P", enterpriseId).Preload(clause.Associations).Preload("AccountingMovement.PurchaseInvoice.Supplier").Order("payment_transaction.id DESC").Find(&paymentTransaction)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	for i := 0; i < len(paymentTransaction); i++ {
		invoices := getAccountingMovementPurchaseInvoices(paymentTransaction[i].AccountingMovementId)
		if len(invoices) > 0 {
			paymentTransaction[i].SupplierName = &invoices[0].Supplier.Name
		}
	}
	return paymentTransaction
}

func searchPaymentTransactions(search CollectionOperationPaymentTransactionSearch, enterpriseId int32) []PaymentTransaction {
	if search.isDefault() {
		return getPendingPaymentTransaction(enterpriseId)
	}

	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)

	// get records using dbOrm
	cursor := dbOrm.Model(&PaymentTransaction{}).Where("payment_transaction.enterprise = ?", enterpriseId)

	if search.Mode != 0 {
		var status string
		if search.Mode == 1 {
			status = "P" // Pending
		} else if search.Mode == 2 {
			status = "C" // Paid
		} else if search.Mode == 3 {
			status = "U" // Unpaid
		}
		cursor.Where("payment_transaction.status = ?", status)
	}

	if search.StartDate != nil {
		cursor.Where("payment_transaction.date_created >= ?", search.StartDate)
	}

	if search.EndDate != nil {
		cursor.Where("payment_transaction.date_created <= ?", search.EndDate)
	}

	result := cursor.Preload(clause.Associations).Preload("AccountingMovement.PurchaseInvoice.Supplier").Order("payment_transaction.id DESC").Find(&paymentTransaction)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	for i := 0; i < len(paymentTransaction); i++ {
		invoices := getAccountingMovementPurchaseInvoices(paymentTransaction[i].AccountingMovementId)
		if len(invoices) > 0 {
			paymentTransaction[i].SupplierName = &invoices[0].Supplier.Name
		}
	}

	return paymentTransaction
}

func getPaymentTransactions(accountingMovement int64, enterpriseId int32) []PaymentTransaction {
	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)
	// get payment transactions for the current enterprise where the accountingMovementId is equal to the accountingMovementId using dbOrm
	result := dbOrm.Where("payment_transaction.accounting_movement = ? AND payment_transaction.enterprise = ?", accountingMovement, enterpriseId).Preload(clause.Associations).Order("id ASC").Find(&paymentTransaction)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return paymentTransaction
}

func getPaymentTransactionRow(paymentTransactionId int32) PaymentTransaction {
	p := PaymentTransaction{}
	// get a single payment transaction using dbOrm
	result := dbOrm.Model(&PaymentTransaction{}).Where("payment_transaction.id = ?", paymentTransactionId).Preload(clause.Associations).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return p
}

func (c *PaymentTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	var paymentTransaction PaymentTransaction
	tx.Model(&PaymentTransaction{}).Last(&paymentTransaction)
	c.Id = paymentTransaction.Id + 1
	return nil
}

func (c *PaymentTransaction) insertPaymentTransaction(userId int32, trans *gorm.DB) bool {
	if c.Total <= 0 {
		return false
	}

	c.Pending = c.Total
	c.Paid = 0
	c.Status = "P"

	p := getPaymentMethodRow(c.PaymentMethodId)
	c.DateCreated = time.Now()
	c.DateExpiration = time.Now().AddDate(0, 0, int(p.DaysExpiration))

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	result := trans.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if c.Id > 0 {
		insertTransactionalLog(c.EnterpriseId, "payment_transaction", int(c.Id), userId, "I")
	} else {
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

// Adds or substracts the paid quantity on the payment transaction
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (pt *PaymentTransaction) addQuantityCharges(charges float64, userId int32, trans gorm.DB) bool {
	var paymentTransaction PaymentTransaction
	result := dbOrm.Model(&PaymentTransaction{}).Where("id = ?", pt.Id).First(&paymentTransaction)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	paymentTransaction.Paid += charges
	paymentTransaction.Pending -= charges
	if paymentTransaction.Pending == 0 {
		paymentTransaction.Status = "C"
	} else {
		paymentTransaction.Status = "P"
	}

	result = trans.Save(&paymentTransaction)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(pt.EnterpriseId, "payment_transaction", int(pt.Id), userId, "U")

	return true
}

func (c *PaymentTransaction) deletePaymentTransaction(userId int32, trans *gorm.DB) bool {
	if c.Id <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	insertTransactionalLog(c.EnterpriseId, "payment_transaction", int(c.Id), userId, "D")

	result := trans.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&PaymentTransaction{})
	if result.Error != nil {
		log("DB", result.Error.Error())
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
