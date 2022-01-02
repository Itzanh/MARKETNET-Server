package main

import (
	"database/sql"
	"math"
	"time"
)

type AccountingMovementDetail struct {
	Id                int64     `json:"id"`
	Movement          int64     `json:"movement"`
	Journal           int32     `json:"journal"`
	Account           int32     `json:"account"`
	DateCreated       time.Time `json:"dateCreated"`
	Credit            float64   `json:"credit"`
	Debit             float64   `json:"debit"`
	Type              string    `json:"type"` // O: Opening, N: Normal, V: Variation of existences, R: Regularisation, C: Closing
	Note              string    `json:"note"`
	DocumentName      string    `json:"documentName"`
	PaymentMethod     int32     `json:"paymentMethod"`
	AccountName       string    `json:"accountName"`
	AccountNumber     int32     `json:"accountNumber"`
	PaymentMethodName string    `json:"paymentMethodName"`
	enterprise        int32
}

func getAccountingMovementDetail(movementId int64, enterpriseId int32) []AccountingMovementDetail {
	accountingMovementDetail := make([]AccountingMovementDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=accounting_movement_detail.account),(SELECT account_number FROM account WHERE account.id=accounting_movement_detail.account),(SELECT name FROM payment_method WHERE payment_method.id=accounting_movement_detail.payment_method) FROM public.accounting_movement_detail WHERE movement=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, movementId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return accountingMovementDetail
	}
	defer rows.Close()

	for rows.Next() {
		a := AccountingMovementDetail{}
		rows.Scan(&a.Id, &a.Movement, &a.Journal, &a.Account, &a.DateCreated, &a.Credit, &a.Debit, &a.Type, &a.Note, &a.DocumentName, &a.PaymentMethod, &a.enterprise, &a.AccountName, &a.AccountNumber, &a.PaymentMethodName)
		accountingMovementDetail = append(accountingMovementDetail, a)
	}

	return accountingMovementDetail
}

func getAccountingMovementDetailRow(detailtId int64) AccountingMovementDetail {
	sqlStatement := `SELECT *,(SELECT account_number FROM account WHERE account.id=accounting_movement_detail.account) FROM public.accounting_movement_detail WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailtId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return AccountingMovementDetail{}
	}

	a := AccountingMovementDetail{}
	row.Scan(&a.Id, &a.Movement, &a.Journal, &a.Account, &a.DateCreated, &a.Credit, &a.Debit, &a.Type, &a.Note, &a.DocumentName, &a.PaymentMethod, &a.enterprise, &a.AccountNumber)

	return a
}

func (a *AccountingMovementDetail) isValid() bool {
	return !(a.Movement <= 0 || a.Journal <= 0 || a.AccountNumber <= 0 || (a.Credit == 0 && a.Debit == 0) || (a.Type != "O" && a.Type != "N" && a.Type != "V" && a.Type != "R" && a.Type != "C") || len(a.Note) > 300 || len(a.DocumentName) > 15 || a.PaymentMethod <= 0)
}

func (a *AccountingMovementDetail) insertAccountingMovementDetail(userId int32, trans *sql.Tx) bool {
	if !a.isValid() {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false
		}
		///
	}

	a.Account = getAccountIdByAccountNumber(a.Journal, a.AccountNumber, a.enterprise)
	if a.Account <= 0 {
		trans.Rollback()
		return false
	}

	// Round float to 2 decimal places (round to nearest)
	a.Credit = float64(math.Round(float64(a.Credit)*100) / 100)
	a.Debit = float64(math.Round(float64(a.Debit)*100) / 100)

	sqlStatement := `INSERT INTO public.accounting_movement_detail(movement, journal, account, credit, debit, type, note, document_name, payment_method, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	row := trans.QueryRow(sqlStatement, a.Movement, a.Journal, a.Account, a.Credit, a.Debit, a.Type, a.Note, a.DocumentName, a.PaymentMethod, a.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false
	}

	var detailId int64
	row.Scan(&detailId)
	if detailId <= 0 {
		trans.Rollback()
		return false
	}
	a.Id = detailId

	insertTransactionalLog(a.enterprise, "accounting_movement_detail", int(detailId), userId, "I")

	m := AccountingMovement{Id: a.Movement}
	ok := m.addCreditAndDebit(a.Credit, a.Debit, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}
	acc := Account{}
	acc.Id = a.Account
	ok = acc.addCreditAndDebit(a.Credit, a.Debit, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false
		}
		///
	}
	return true
}

func (a *AccountingMovementDetail) deleteAccountingMovementDetail(userId int32, trans *sql.Tx) bool {
	if a.Id <= 0 {
		return false
	}

	accountingMovementDetailInMemory := getAccountingMovementDetailRow(a.Id)
	settings := getSettingsRecordById(a.enterprise)
	if accountingMovementDetailInMemory.Id <= 0 || accountingMovementDetailInMemory.enterprise != a.enterprise || (settings.LimitAccountingDate != nil && accountingMovementDetailInMemory.DateCreated.Before(*settings.LimitAccountingDate)) {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false
		}
		///
	}

	inMemoryDetail := getAccountingMovementDetailRow(a.Id)
	if inMemoryDetail.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(a.enterprise, "accounting_movement_detail", int(a.Id), userId, "I")

	sqlStatement := `DELETE FROM public.accounting_movement_detail WHERE id=$1 AND enterprise=$2`
	_, err := trans.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	m := AccountingMovement{Id: inMemoryDetail.Movement}
	ok := m.addCreditAndDebit(-inMemoryDetail.Credit, -inMemoryDetail.Debit, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}
	acc := Account{}
	acc.Id = inMemoryDetail.Account
	ok = acc.addCreditAndDebit(-inMemoryDetail.Credit, -inMemoryDetail.Debit, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false
		}
		///
	}
	return true
}
