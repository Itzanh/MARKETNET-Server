package main

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectionOperation struct {
	Id                         int32                    `json:"id" gorm:"index:collection_operation_id_enterprise,unique:true,priority:1"`
	AccountingMovementId       int64                    `json:"accountingMovementId" gorm:"column:accounting_movement;not null:true"`
	AccountingMovement         AccountingMovement       `json:"accountingMovement" gorm:"foreignkey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	AccountingMovementDetailId int64                    `json:"accountingMovementDetailId" gorm:"column:accounting_movement_detail;not null:true"`
	AccountingMovementDetail   AccountingMovementDetail `json:"accountingMovementDetail" gorm:"foreignkey:AccountingMovementDetailId,EnterpriseId;references:Id,EnterpriseId"`
	AccountId                  int32                    `json:"accountId" gorm:"column:account;not null:true"`
	Account                    Account                  `json:"account" gorm:"foreignkey:AccountId,EnterpriseId;references:Id,EnterpriseId"`
	BankId                     *int32                   `json:"bankId" gorm:"column:bank"`
	Bank                       *Account                 `json:"bank" gorm:"foreignkey:BankId,EnterpriseId;references:Id,EnterpriseId"`
	Status                     string                   `json:"status" gorm:"type:character(1);not null:true;index:collection_operation_status_enterprise,priority:1"` // P = Pending, C = Paid, U = Unpaid
	DateCreated                time.Time                `json:"dateCreated" gorm:"type:timestamp(3) without time zone;not null:true;index:collection_operation_date_created"`
	DateExpiration             time.Time                `json:"dateExpiration" gorm:"type:timestamp(3) without time zone;not null:true"`
	Total                      float64                  `json:"total" gorm:"type:numeric(14,6);not null:true"`
	Paid                       float64                  `json:"paid" gorm:"type:numeric(14,6);not null:true"`
	Pending                    float64                  `json:"pending" gorm:"type:numeric(14,6);not null:true"`
	DocumentName               string                   `json:"documentName" gorm:"type:character(15);not null:true"`
	PaymentMethodId            int32                    `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod              PaymentMethod            `json:"paymentMethod" gorm:"foreignkey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId               int32                    `json:"-" gorm:"column:enterprise;not null:true;index:collection_operation_id_enterprise,unique:true,priority:2;index:collection_operation_status_enterprise,priority:2"`
	Enterprise                 Settings                 `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	CustomerName               *string                  `json:"customerName" gorm:"-"`
}

func (co *CollectionOperation) TableName() string {
	return "collection_operation"
}

type CollectionOperationPaymentTransactionSearch struct {
	Mode      uint8      `json:"mode"` // 0 = All, 1 = Pending, 2 = Paid, 3 = Unpaid
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

func (search *CollectionOperationPaymentTransactionSearch) isDefault() bool {
	return search.Mode == 1 && search.StartDate == nil && search.EndDate == nil
}

func getPendingColletionOperations(enterpriseId int32) []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	// get pending collection operations for the current enterprise where the status is 'P' using dbOrm
	result := dbOrm.Where("status = 'P' AND enterprise = ?", enterpriseId).Order("collection_operation.id DESC").Preload(clause.Associations).Find(&collectionOperation)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	for i := 0; i < len(collectionOperation); i++ {
		invoices := getAccountingMovementSaleInvoices(collectionOperation[i].AccountingMovementId)
		if len(invoices) > 0 {
			collectionOperation[i].CustomerName = &invoices[0].Customer.Name
		}
	}
	return collectionOperation
}

func searchCollectionOperations(search CollectionOperationPaymentTransactionSearch, enterpriseId int32) []CollectionOperation {
	if search.isDefault() {
		return getPendingColletionOperations(enterpriseId)
	}

	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)

	cursor := dbOrm.Model(&CollectionOperation{}).Where("collection_operation.enterprise = ?", enterpriseId)

	if search.Mode != 0 {
		var status string
		if search.Mode == 1 {
			status = "P" // Pending
		} else if search.Mode == 2 {
			status = "C" // Paid
		} else if search.Mode == 3 {
			status = "U" // Unpaid
		}
		cursor.Where("collection_operation.status = ?", status)
	}

	if search.StartDate != nil {
		cursor.Where("collection_operation.date_created >= ?", search.StartDate)
	}

	if search.EndDate != nil {
		cursor.Where("collection_operation.date_created <= ?", search.EndDate)
	}

	result := cursor.Order("collection_operation.id DESC").Preload(clause.Associations).Find(&collectionOperation)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	for i := 0; i < len(collectionOperation); i++ {
		invoices := getAccountingMovementSaleInvoices(collectionOperation[i].AccountingMovementId)
		if len(invoices) > 0 {
			collectionOperation[i].CustomerName = &invoices[0].Customer.Name
		}
	}

	return collectionOperation
}

func getColletionOperations(accountingMovement int64, enterpriseId int32) []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	// get collection operations for the current enterprise where the accounting_movement_id is equal to the accounting_movement_id using dbOrm
	result := dbOrm.Where("accounting_movement = ? AND enterprise = ?", accountingMovement, enterpriseId).Order("collection_operation.id DESC").Preload(clause.Associations).Find(&collectionOperation)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return collectionOperation
}

func getColletionOperationRow(collectionOperationId int32) CollectionOperation {
	o := CollectionOperation{}
	// get a single collection operation using dbOrm
	result := dbOrm.Where("id = ?", collectionOperationId).Preload(clause.Associations).First(&o)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return o
}

func (c *CollectionOperation) BeforeCreate(tx *gorm.DB) (err error) {
	var collectionOperation CollectionOperation
	tx.Model(&CollectionOperation{}).Last(&collectionOperation)
	c.Id = collectionOperation.Id + 1
	return nil
}

func (c *CollectionOperation) insertCollectionOperation(userId int32, trans *gorm.DB) bool {
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
		insertTransactionalLog(c.EnterpriseId, "collection_operation", int(c.Id), userId, "I")
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

// Adds or substracts the paid quantity on the collection operation
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (c *CollectionOperation) addQuantityCharges(charges float64, userId int32, trans gorm.DB) bool {
	var collectionOperation CollectionOperation
	result := dbOrm.Model(&CollectionOperation{}).Where("id = ?", c.Id).First(&collectionOperation)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	collectionOperation.Paid += charges
	collectionOperation.Pending -= charges
	if collectionOperation.Pending <= 0 {
		collectionOperation.Status = "C" // Paid
	} else {
		collectionOperation.Status = "P" // Pending
	}

	result = dbOrm.Save(&collectionOperation)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.EnterpriseId, "collection_operation", int(c.Id), userId, "U")

	return true
}

func (c *CollectionOperation) deleteCollectionOperation(userId int32, trans *gorm.DB) bool {
	if c.Id <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	insertTransactionalLog(c.EnterpriseId, "collection_operation", int(c.Id), userId, "D")

	result := trans.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&CollectionOperation{})
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
