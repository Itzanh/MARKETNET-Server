/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountingMovement struct {
	Id             int64        `json:"id" gorm:"index:accounting_movement_id_enterprise,unique:true,priority:1"`
	DateCreated    time.Time    `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true;index:accounting_movement_date_created,priority:1,sort:desc"`
	AmountDebit    float64      `json:"amountDebit" gorm:"type:numeric(14,6);not null:true"`
	AmountCredit   float64      `json:"amountCredit" gorm:"type:numeric(14,6);not null:true"`
	FiscalYear     int16        `json:"fiscalYear" gorm:"not null:true"`
	Type           string       `json:"type" gorm:"type:character(1);not null:true"` // O: Opening, N: Normal, V: Variation of existences, R: Regularisation, C: Closing
	BillingSerieId string       `json:"billingSerieId" gorm:"column:billing_serie;type:character(3);not null:true"`
	BillingSerie   BillingSerie `json:"billingSerie" gorm:"foreignkey:BillingSerieId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId   int32        `json:"-" gorm:"column:enterprise;not null:true;index:accounting_movement_id_enterprise,unique:true,priority:2"`
	Enterprise     Settings     `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (am *AccountingMovement) TableName() string {
	return "accounting_movement"
}

func getAccountingMovement(enterpriseId int32) []AccountingMovement {
	accountingMovements := make([]AccountingMovement, 0)
	// get all accounting movements from the database for the current enterprise sort by date created descending (newest first) using dbOrm
	dbOrm.Where("accounting_movement.enterprise = ?", enterpriseId).Order("accounting_movement.date_created DESC").Preload(clause.Associations).Find(&accountingMovements)
	return accountingMovements
}

type AccountingMovementSearch struct {
	Search         string     `json:"search"`
	Type           *string    `json:"type"`
	BillingSerieId *string    `json:"billingSerieId"`
	DateStart      *time.Time `json:"dateStart"`
	DateEnd        *time.Time `json:"dateEnd"`
}

func (query *AccountingMovementSearch) searchAccountingMovements(enterpriseId int32) []AccountingMovement {
	accountingMovements := make([]AccountingMovement, 0)
	// get all accounting movements from the database for the current enterprise sort by date created descending (newest first) using dbOrm
	cursor := dbOrm.Model(&AccountingMovement{}).Where(`(accounting_movement_detail.document_name ILIKE @search) AND (accounting_movement.enterprise = @enterpriseId)`, sql.Named("enterpriseId", enterpriseId), sql.Named("search", "%"+query.Search+"%"))
	if query.Type != nil {
		cursor = cursor.Where("accounting_movement.type = ?", *query.Type)
	}
	if query.BillingSerieId != nil {
		cursor = cursor.Where("accounting_movement.billing_serie = ?", *query.BillingSerieId)
	}
	if query.DateStart != nil {
		cursor = cursor.Where("accounting_movement.date_created >= ?", *query.DateStart)
	}
	if query.DateEnd != nil {
		cursor = cursor.Where("accounting_movement.date_created <= ?", *query.DateEnd)
	}
	result := cursor.Joins("INNER JOIN accounting_movement_detail ON accounting_movement_detail.movement=accounting_movement.id").Order("accounting_movement.date_created DESC").Preload(clause.Associations).Find(&accountingMovements)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return accountingMovements
	}
	return accountingMovements
}

func getAccountingMovementRow(accountingMovementId int64) AccountingMovement {
	// get a single accounting movement row by id using dbOrm
	a := AccountingMovement{}
	dbOrm.Where("id = ?", accountingMovementId).Preload("BillingSerie").First(&a)
	return a
}

func (a *AccountingMovement) isValid() bool {
	return !((a.Type != "O" && a.Type != "N" && a.Type != "V" && a.Type != "R" && a.Type != "C") || len(a.BillingSerieId) != 3)
}

func (a *AccountingMovement) BeforeCreate(tx *gorm.DB) (err error) {
	var accountingMovement AccountingMovement
	tx.Model(&AccountingMovement{}).Last(&accountingMovement)
	a.Id = accountingMovement.Id + 1
	return nil
}

func (am *AccountingMovement) insertAccountingMovement(userId int32, trans *gorm.DB) bool {
	if !am.isValid() {
		return false
	}

	am.DateCreated = time.Now()
	am.FiscalYear = int16(time.Now().Year())

	result := dbOrm.Create(&am)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(am.EnterpriseId, "accounting_movement", int(am.Id), userId, "I")

	return true
}

func (am *AccountingMovement) deleteAccountingMovement(userId int32, trans *gorm.DB) bool {
	if am.Id <= 0 {
		return false
	}

	accountingMovementInMemory := getAccountingMovementRow(am.Id)
	settings := getSettingsRecordById(accountingMovementInMemory.EnterpriseId)
	if accountingMovementInMemory.Id <= 0 || accountingMovementInMemory.EnterpriseId != am.EnterpriseId || (settings.LimitAccountingDate != nil && accountingMovementInMemory.DateCreated.Before(*settings.LimitAccountingDate)) {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		///
	}

	// cascade delete the collection operations
	c := getColletionOperations(am.Id, am.EnterpriseId)
	for i := 0; i < len(c); i++ {
		ok := c[i].deleteCollectionOperation(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// cascade delete the payment transactions
	p := getPaymentTransactions(am.Id, am.EnterpriseId)
	for i := 0; i < len(p); i++ {
		ok := p[i].deletePaymentTransaction(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// cascade delete the details
	d := getAccountingMovementDetail(am.Id, am.EnterpriseId)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteAccountingMovementDetail(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// set the relation null on the sale invoices
	var salesInvoiceId int32
	result := dbOrm.Model(&SalesInvoice{}).Where("accounting_movement = ?", am.Id).Select("id").Pluck("id", &salesInvoiceId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&SalesInvoice{}).Where("accounting_movement = ?", am.Id).Update("accounting_movement", nil)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(am.EnterpriseId, "sales_invoice", int(salesInvoiceId), userId, "U")

	// set the relation null on the purchase invoices
	var purchaseInvoiceId int32
	result = dbOrm.Model(&PurchaseInvoice{}).Where("accounting_movement = ?", am.Id).Select("id").Pluck("id", &purchaseInvoiceId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&PurchaseInvoice{}).Where("accounting_movement = ?", am.Id).Update("accounting_movement", nil)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(am.EnterpriseId, "purchase_invoice", int(purchaseInvoiceId), userId, "U")

	insertTransactionalLog(am.EnterpriseId, "accounting_movement", int(am.Id), userId, "D")

	// delete the movement
	result = trans.Where("id = ? AND enterprise = ?", am.Id, am.EnterpriseId).Delete(&AccountingMovement{})
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

type PostInvoiceResult struct {
	Invoice int64 `json:"invoice"`
	Ok      bool  `json:"ok"`
	Result  int8  `json:"result"` // 0: Internal error, 1: The customer/supplier in the invoice has no account
}

// Transfer the sales invoices from management to accounting. Create the movements and the details for all the selected invoices.
func salesPostInvoices(invoiceIds []int64, enterpriseId int32, userId int32) []PostInvoiceResult {
	result := make([]PostInvoiceResult, 0)
	// validation
	if len(invoiceIds) == 0 {
		return result
	}
	for i := 0; i < len(invoiceIds); i++ {
		result = append(result, PostInvoiceResult{Invoice: invoiceIds[i]})
	}
	for i := 0; i < len(invoiceIds); i++ {
		if invoiceIds[i] <= 0 {
			return result
		}
	}
	settings := getSettingsRecordById(enterpriseId)
	if settings.SalesJournalId == nil {
		return result
	}

	///
	trans := dbOrm.Begin()
	///

	for i := 0; i < len(invoiceIds); i++ {
		// get the selected invoice
		inv := getSalesInvoiceRow(invoiceIds[i])
		if inv.Id <= 0 || inv.EnterpriseId != enterpriseId {
			trans.Rollback()
			return result
		}
		// get the invoice customer
		c := getCustomerRow(inv.CustomerId)
		if c.Id <= 0 {
			trans.Rollback()
			return result
		}
		if c.AccountId == nil {
			result[i].Result = 1
			continue
		}
		// get the account row
		a := getAccountRow(*c.AccountId)
		if a.Id <= 0 {
			trans.Rollback()
			return result
		}

		// create the accounting movement
		m := AccountingMovement{}
		m.Type = "N"
		m.BillingSerieId = inv.BillingSeriesId
		m.EnterpriseId = inv.EnterpriseId
		ok := m.insertAccountingMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return result
		}

		// create the details
		// 1. detail for the customer
		dCust := AccountingMovementDetail{}
		dCust.MovementId = m.Id
		dCust.JournalId = a.JournalId
		dCust.AccountNumber = a.AccountNumber
		dCust.Debit = inv.TotalAmount
		dCust.Type = "N"
		dCust.DocumentName = inv.InvoiceName
		dCust.PaymentMethodId = inv.PaymentMethodId
		dCust.EnterpriseId = enterpriseId
		ok = dCust.insertAccountingMovementDetail(userId, trans)
		if !ok {
			trans.Rollback()
			return result
		}

		// create the collection operation for this charge
		if inv.TotalAmount > 0 {
			p := getPaymentMethodRow(inv.PaymentMethodId)

			co := CollectionOperation{}
			co.AccountingMovementId = m.Id
			co.AccountingMovementDetailId = dCust.Id
			co.AccountId = getAccountIdByAccountNumber(a.JournalId, a.AccountNumber, a.EnterpriseId)
			co.Total = inv.TotalAmount
			co.DocumentName = inv.InvoiceName
			co.PaymentMethodId = inv.PaymentMethodId
			co.BankId = p.Bank
			co.EnterpriseId = enterpriseId
			ok := co.insertCollectionOperation(userId, trans)
			if !ok {
				trans.Rollback()
				return result
			}

			// paid in advance
			if p.PaidInAdvance {
				ch := Charges{}
				ch.CollectionOperationId = co.Id
				ch.Amount = inv.TotalAmount
				ch.EnterpriseId = enterpriseId
				ok = ch.insertCharges(userId)
				if !ok {
					trans.Rollback()
					return result
				}
			}
		}

		// 2. details for the income
		// We can create a single line for the income, or we can split it in different accounts, it is used for income / spending per products or departments:
		// == 1. If there are no custom accounts in the products we create a single line:
		// 700.000001 Sales 100€
		// == 2. If some of the products have a custom account, and some other don't, we put the details price in the custom accounts, and the remaining in the generic sales account:
		// 700.000002 Wood sales 50€
		// 700.000001 Sales 50€
		// == 3. If none of the products are without a custom account, we don't create an income in the generic sales accounts:
		// 700.000002 Wood sales 50€
		// 700.000003 Shipping 25€
		// 700.000004 Software sales 25€
		det := getSalesInvoiceDetail(inv.Id, inv.EnterpriseId)
		var detailIncomeCredit float64 = inv.TotalWithDiscount
		for i := 0; i < len(det); i++ {
			if det[i].ProductId == nil {
				continue
			}
			customAccount := getProductAccount(*det[i].ProductId, "S") // Sales
			if customAccount == nil {
				continue
			}
			detailIncomeCredit -= det[i].Price * float64(det[i].Quantity)
			var accountingMovementDetailCount int64
			res := trans.Model(&AccountingMovementDetail{}).Where("movement = ? AND journal = ? AND account = ?", m.Id, customAccount.JournalId, customAccount.Account.Id).Count(&accountingMovementDetailCount)
			if res.Error != nil {
				log("DB", res.Error.Error())
				trans.Rollback()
				return result
			}
			if accountingMovementDetailCount == 0 {
				dInc := AccountingMovementDetail{}
				dInc.MovementId = m.Id
				dInc.JournalId = customAccount.JournalId
				dInc.AccountNumber = customAccount.Account.AccountNumber
				dInc.Credit = det[i].Price * float64(det[i].Quantity)
				dInc.Type = "N"
				dInc.DocumentName = inv.InvoiceName
				dInc.PaymentMethodId = inv.PaymentMethodId
				dInc.EnterpriseId = enterpriseId
				ok = dInc.insertAccountingMovementDetail(userId, trans)
				if !ok {
					trans.Rollback()
					return result
				}
			} else {
				var accountingMovementDetailId int64
				res = trans.Model(&AccountingMovementDetail{}).Where("movement = ? AND journal = ? AND account = ?", m.Id, customAccount.JournalId, customAccount.Account.Id).Select("id").Limit(1).Pluck("id", &accountingMovementDetailId)
				if res.Error != nil {
					log("DB", res.Error.Error())
					trans.Rollback()
					return result
				}

				accountingMovementDetail := getAccountingMovementDetailRow(accountingMovementDetailId)
				accountingMovementDetail.Credit += det[i].Price * float64(det[i].Quantity)
				res = trans.Model(&AccountingMovementDetail{}).Where("id = ?", accountingMovementDetailId).Update("credit", accountingMovementDetail.Credit)
				if res.Error != nil {
					log("DB", res.Error.Error())
					trans.Rollback()
					return result
				}
			}
		} // for i := 0; i < len(det); i++ {
		if detailIncomeCredit != 0 {
			dInc := AccountingMovementDetail{}
			dInc.MovementId = m.Id
			dInc.JournalId = *settings.SalesJournalId
			dInc.AccountNumber = 1
			dInc.Credit = detailIncomeCredit
			dInc.Type = "N"
			dInc.DocumentName = inv.InvoiceName
			dInc.PaymentMethodId = inv.PaymentMethodId
			dInc.EnterpriseId = enterpriseId
			ok = dInc.insertAccountingMovementDetail(userId, trans)
			if !ok {
				trans.Rollback()
				return result
			}
		}

		// 3. details for the VAT

		// get the details and sort
		det = getSalesInvoiceDetail(inv.Id, inv.EnterpriseId)
		d := make([]SalesInvoiceDetail, 0)
		for i := 0; i < len(det); i++ {
			if det[i].VatPercent > 0 {
				d = append(d, det[i])
			}
		}
		sort.Slice(d[:], func(i, j int) bool {
			return d[i].VatPercent < d[j].VatPercent
		})

		// multisplit for the vat percent
		details := make([]SalesInvoiceDetail, 0)
		for i := 0; i < len(d); i++ {
			if i == 0 || d[i].VatPercent == d[i-1].VatPercent {
				details = append(details, d[i])
			}
			if (i > 0 && d[i].VatPercent != d[i-1].VatPercent) || (i == len(d)-1) {
				if len(details) == 0 || details[0].VatPercent == 0 {
					details = make([]SalesInvoiceDetail, 0)
					continue
				}

				// get the account for this var percent
				journal, accountNumber := getConfigAccountsVatSaleRow(details[0].VatPercent, enterpriseId)

				// we have an array with all the same vat percent
				var credit float64
				for j := 0; j < len(details); j++ {
					credit += details[j].TotalAmount - (details[j].Price * float64(details[j].Quantity))
				}

				dVat := AccountingMovementDetail{}
				dVat.MovementId = m.Id
				dVat.JournalId = journal
				dVat.AccountNumber = accountNumber
				dVat.Credit = credit
				dVat.Type = "N"
				dVat.DocumentName = inv.InvoiceName
				dVat.PaymentMethodId = inv.PaymentMethodId
				dVat.EnterpriseId = enterpriseId
				ok := dVat.insertAccountingMovementDetail(userId, trans)
				if !ok {
					trans.Rollback()
					return result
				}

				details = make([]SalesInvoiceDetail, 0)
				details = append(details, d[i])
			}
		}

		// set the accounting movement on the invoice
		res := trans.Model(&SalesInvoice{}).Where("id = ?", invoiceIds[i]).Update("accounting_movement", m.Id)
		if res.Error != nil {
			log("DB", res.Error.Error())
			trans.Rollback()
			return result
		}

		insertTransactionalLog(a.EnterpriseId, "sales_invoice", int(invoiceIds[i]), userId, "U")
		result[i].Ok = true
	}

	///
	trans.Commit()
	///
	return result
}

// Transfer the purchase invoices from management to accounting. Create the movements and the details for all the selected invoices.
func purchasePostInvoices(invoiceIds []int64, enterpriseId int32, userId int32) []PostInvoiceResult {
	result := make([]PostInvoiceResult, 0)
	// validation
	if len(invoiceIds) == 0 {
		return result
	}
	for i := 0; i < len(invoiceIds); i++ {
		result = append(result, PostInvoiceResult{Invoice: invoiceIds[i]})
	}
	for i := 0; i < len(invoiceIds); i++ {
		if invoiceIds[i] <= 0 {
			return result
		}
	}
	settings := getSettingsRecordById(enterpriseId)
	if settings.PurchaseJournalId == nil {
		return result
	}

	///
	trans := dbOrm.Begin()
	///

	for i := 0; i < len(invoiceIds); i++ {
		// get the selected invoice
		inv := getPurchaseInvoiceRow(invoiceIds[i])
		if inv.Id <= 0 || inv.EnterpriseId != enterpriseId {
			trans.Rollback()
			return result
		}
		// get the invoice customer
		s := getSupplierRow(inv.SupplierId)
		if s.Id <= 0 {
			trans.Rollback()
			return result
		}
		if s.AccountId == nil {
			result[i].Result = 1
			continue
		}
		// get the account row
		a := getAccountRow(*s.AccountId)
		if a.Id <= 0 {
			trans.Rollback()
			return result
		}

		// create the accounting movement
		m := AccountingMovement{}
		m.Type = "N"
		m.BillingSerieId = inv.BillingSeriesId
		m.EnterpriseId = inv.EnterpriseId
		ok := m.insertAccountingMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return result
		}

		// create the details
		// 1. detail for the supplier
		dCust := AccountingMovementDetail{}
		dCust.MovementId = m.Id
		dCust.JournalId = a.JournalId
		dCust.AccountNumber = a.AccountNumber
		dCust.Credit = inv.TotalAmount
		dCust.Type = "N"
		dCust.DocumentName = inv.InvoiceName
		dCust.PaymentMethodId = inv.PaymentMethodId
		dCust.EnterpriseId = enterpriseId
		ok = dCust.insertAccountingMovementDetail(userId, trans)
		if !ok {
			trans.Rollback()
			return result
		}

		// create the payment transaction for this payment
		if inv.TotalAmount > 0 {
			p := getPaymentMethodRow(inv.PaymentMethodId)

			pt := PaymentTransaction{}
			pt.AccountingMovementId = m.Id
			pt.AccountingMovementDetailId = dCust.Id
			pt.AccountId = getAccountIdByAccountNumber(a.JournalId, a.AccountNumber, a.EnterpriseId)
			pt.Total = inv.TotalAmount
			pt.DocumentName = inv.InvoiceName
			pt.PaymentMethodId = *s.PaymentMethodId
			pt.BankId = p.Bank
			pt.EnterpriseId = enterpriseId
			ok := pt.insertPaymentTransaction(userId, trans)
			if !ok {
				trans.Rollback()
				return result
			}

			// paid in advance
			if p.PaidInAdvance {
				py := Payment{}
				py.PaymentTransactionId = pt.Id
				py.Amount = inv.TotalAmount
				py.EnterpriseId = enterpriseId
				ok = py.insertPayment(userId)
				if !ok {
					trans.Rollback()
					return result
				}
			}
		}

		// 2. details for the outcome
		// We can create a single line for the outcome, or we can split it in different accounts, it is used for income / spending per products or departments:
		// == 1. If there are no custom accounts in the products we create a single line:
		// 700.000001 Purchases 100€
		// == 2. If some of the products have a custom account, and some other don't, we put the details price in the custom accounts, and the remaining in the generic purchases account:
		// 700.000002 Wood purchase 50€
		// 700.000001 Purchases 50€
		// == 3. If none of the products are without a custom account, we don't create an income in the generic purchases accounts:
		// 700.000002 Wood purchase 50€
		// 700.000003 Shipping 25€
		// 700.000004 Software purchase 25€
		det := getPurchaseInvoiceDetail(inv.Id, inv.EnterpriseId)
		var detailIncomeCredit float64 = inv.TotalWithDiscount
		for i := 0; i < len(det); i++ {
			if det[i].ProductId == nil {
				continue
			}
			customAccount := getProductAccount(*det[i].ProductId, "P") // Purchases
			if customAccount == nil {
				continue
			}
			detailIncomeCredit -= det[i].Price * float64(det[i].Quantity)
			dInc := AccountingMovementDetail{}
			dInc.MovementId = m.Id
			dInc.JournalId = customAccount.JournalId
			dInc.AccountNumber = customAccount.Account.AccountNumber
			dInc.Credit = det[i].Price * float64(det[i].Quantity)
			dInc.Type = "N"
			dInc.DocumentName = inv.InvoiceName
			dInc.PaymentMethodId = inv.PaymentMethodId
			dInc.EnterpriseId = enterpriseId
			ok = dInc.insertAccountingMovementDetail(userId, trans)
			if !ok {
				trans.Rollback()
				return result
			}
		}
		if detailIncomeCredit != 0 {
			dInc := AccountingMovementDetail{}
			dInc.MovementId = m.Id
			dInc.JournalId = *settings.PurchaseJournalId
			dInc.AccountNumber = 1
			dInc.Debit = detailIncomeCredit
			dInc.Type = "N"
			dInc.DocumentName = inv.InvoiceName
			dInc.PaymentMethodId = inv.PaymentMethodId
			dInc.EnterpriseId = enterpriseId
			ok = dInc.insertAccountingMovementDetail(userId, trans)
			if !ok {
				trans.Rollback()
				return result
			}
		}

		// 3. details for the VAT

		// get the details and sort
		det = getPurchaseInvoiceDetail(inv.Id, inv.EnterpriseId)
		d := make([]PurchaseInvoiceDetail, 0)
		for i := 0; i < len(det); i++ {
			if det[i].VatPercent > 0 {
				d = append(d, det[i])
			}
		}
		sort.Slice(d[:], func(i, j int) bool {
			return d[i].VatPercent < d[j].VatPercent
		})

		// multisplit for the vat percent
		details := make([]PurchaseInvoiceDetail, 0)
		for i := 0; i < len(d); i++ {
			if i == 0 || d[i].VatPercent == d[i-1].VatPercent {
				details = append(details, d[i])
			}
			if (i > 0 && d[i].VatPercent != d[i-1].VatPercent) || (i == len(d)-1) {
				if len(details) == 0 || details[0].VatPercent == 0 {
					details = make([]PurchaseInvoiceDetail, 0)
					continue
				}

				// get the account for this var percent
				journal, accountNumber := getConfigAccountsVatPurchaseRow(details[0].VatPercent, enterpriseId)

				// we have an array with all the same vat percent
				var debit float64
				for j := 0; j < len(details); j++ {
					debit += details[j].TotalAmount - (details[j].Price * float64(details[j].Quantity))
				}

				dVat := AccountingMovementDetail{}
				dVat.MovementId = m.Id
				dVat.JournalId = journal
				dVat.AccountNumber = accountNumber
				dVat.Debit = debit
				dVat.Type = "N"
				dVat.DocumentName = inv.InvoiceName
				dVat.PaymentMethodId = inv.PaymentMethodId
				dVat.EnterpriseId = enterpriseId
				ok := dVat.insertAccountingMovementDetail(userId, trans)
				if !ok {
					trans.Rollback()
					return result
				}

				details = make([]PurchaseInvoiceDetail, 0)
				details = append(details, d[i])
			}
		}

		// set the accounting movement on the invoice
		res := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceIds[i]).Update("accounting_movement", m.Id)
		if res.Error != nil {
			log("DB", res.Error.Error())
			trans.Rollback()
			return result
		}
		insertTransactionalLog(a.EnterpriseId, "purchase_invoice", int(invoiceIds[i]), userId, "U")
		result[i].Ok = true
	}

	///
	trans.Commit()
	///
	return result
}

func getAccountingMovementSaleInvoices(movementId int64) []SalesInvoice {
	var invoices []SalesInvoice = make([]SalesInvoice, 0)
	result := dbOrm.Model(&SalesInvoice{}).Where("accounting_movement = ?", movementId).Order("date_created DESC").Preload(clause.Associations).Find(&invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return invoices
	}
	return invoices
}

func getAccountingMovementPurchaseInvoices(movementId int64) []PurchaseInvoice {
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	result := dbOrm.Model(&PurchaseInvoice{}).Where("accounting_movement = ?", movementId).Order("date_created DESC").Preload(clause.Associations).Find(&invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return invoices
	}
	return invoices
}

/* TRIAL BALANCE */

type TrialBalanceQuery struct {
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
	Journal   int32     `json:"journal"`
}

func (q *TrialBalanceQuery) isValid() bool {
	return !(q.Journal <= 0 || q.DateStart.IsZero() || q.DateEnd.IsZero())
}

type TrialBalanceAccount struct {
	JournalId int32   `json:"journalId" gorm:"column:journal"`
	AccountId int32   `json:"accountId" gorm:"column:account"`
	Account   Account `json:"account"`
	Credit    float64 `json:"credit"`
	Debit     float64 `json:"debit"`
	Balance   float64 `json:"balance"`
}

func (q *TrialBalanceQuery) getTrialBalance(enterpriseId int32) []TrialBalanceAccount {
	balance := make([]TrialBalanceAccount, 0)
	if !q.isValid() {
		return balance
	}

	result := dbOrm.Model(&AccountingMovementDetail{}).Where("accounting_movement_detail.journal = ? AND accounting_movement_detail.enterprise = ? AND accounting_movement_detail.date_created >= ? AND accounting_movement_detail.date_created <= ?", q.Journal, enterpriseId, q.DateStart, q.DateEnd).Select("accounting_movement_detail.journal, accounting_movement_detail.account, SUM(accounting_movement_detail.debit) AS debit, SUM(accounting_movement_detail.credit) AS credit").Group("accounting_movement_detail.journal,accounting_movement_detail.account").Find(&balance)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return balance
	}
	for i := 0; i < len(balance); i++ {
		balance[i].Balance = balance[i].Credit - balance[i].Debit
		balance[i].Account = getAccountRow(balance[i].AccountId)
	}
	return balance
}
