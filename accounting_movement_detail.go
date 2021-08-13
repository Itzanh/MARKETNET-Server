package main

import (
	"math"
	"time"
)

type AccountingMovementDetail struct {
	Id                int64     `json:"id"`
	Movement          int64     `json:"movement"`
	Journal           int16     `json:"journal"`
	Account           int32     `json:"account"`
	DateCreated       time.Time `json:"dateCreated"`
	Credit            float32   `json:"credit"`
	Debit             float32   `json:"debit"`
	Type              string    `json:"type"` // O: Opening, N: Normal, V: Variation of existences, R: Regularisation, C: Closing
	Note              string    `json:"note"`
	DocumentName      string    `json:"documentName"`
	PaymentMethod     int16     `json:"paymentMethod"`
	AccountName       string    `json:"accountName"`
	AccountNumber     int32     `json:"accountNumber"`
	PaymentMethodName string    `json:"paymentMethodName"`
}

func getAccountingMovementDetail(movementId int64) []AccountingMovementDetail {
	accountingMovementDetail := make([]AccountingMovementDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=accounting_movement_detail.account),(SELECT account_number FROM account WHERE account.id=accounting_movement_detail.account),(SELECT name FROM payment_method WHERE payment_method.id=accounting_movement_detail.payment_method) FROM public.accounting_movement_detail WHERE movement=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, movementId)
	if err != nil {
		log("DB", err.Error())
		return accountingMovementDetail
	}

	for rows.Next() {
		a := AccountingMovementDetail{}
		rows.Scan(&a.Id, &a.Movement, &a.Journal, &a.Account, &a.DateCreated, &a.Credit, &a.Debit, &a.Type, &a.Note, &a.DocumentName, &a.PaymentMethod, &a.AccountName, &a.AccountNumber, &a.PaymentMethodName)
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
	row.Scan(&a.Id, &a.Movement, &a.Journal, &a.Account, &a.DateCreated, &a.Credit, &a.Debit, &a.Type, &a.Note, &a.DocumentName, &a.PaymentMethod, &a.AccountNumber)

	return a
}

func (a *AccountingMovementDetail) isValid() bool {
	return !(a.Movement <= 0 || a.Journal <= 0 || a.AccountNumber <= 0 || (a.Credit == 0 && a.Debit == 0) || (a.Type != "O" && a.Type != "N" && a.Type != "V" && a.Type != "R" && a.Type != "C") || len(a.Note) > 300 || len(a.DocumentName) > 15 || a.PaymentMethod <= 0)
}

func (a *AccountingMovementDetail) insertAccountingMovementDetail() bool {
	if !a.isValid() {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	a.Account = getAccountIdByAccountNumber(a.Journal, a.AccountNumber)
	if a.Account <= 0 {
		trans.Rollback()
		return false
	}

	// Round float to 2 decimal places (round to nearest)
	a.Credit = float32(math.Round(float64(a.Credit)*100) / 100)
	a.Debit = float32(math.Round(float64(a.Debit)*100) / 100)

	sqlStatement := `INSERT INTO public.accounting_movement_detail(movement, journal, account, credit, debit, type, note, document_name, payment_method) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := db.QueryRow(sqlStatement, a.Movement, a.Journal, a.Account, a.Credit, a.Debit, a.Type, a.Note, a.DocumentName, a.PaymentMethod)
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

	m := AccountingMovement{Id: a.Movement}
	ok := m.addCreditAndDebit(a.Credit, a.Debit)
	if !ok {
		trans.Rollback()
		return false
	}
	acc := Account{}
	acc.Id = a.Account
	ok = acc.addCreditAndDebit(a.Credit, a.Debit)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}

func (a *AccountingMovementDetail) deleteAccountingMovementDetail() bool {
	if a.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	inMemoryDetail := getAccountingMovementDetailRow(a.Id)
	if inMemoryDetail.Id <= 0 {
		trans.Rollback()
		return false
	}

	sqlStatement := `DELETE FROM public.accounting_movement_detail WHERE id=$1`
	_, err = db.Exec(sqlStatement, a.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	m := AccountingMovement{Id: inMemoryDetail.Movement}
	ok := m.addCreditAndDebit(-inMemoryDetail.Credit, -inMemoryDetail.Debit)
	if !ok {
		trans.Rollback()
		return false
	}
	acc := Account{}
	acc.Id = inMemoryDetail.Account
	ok = acc.addCreditAndDebit(-inMemoryDetail.Credit, -inMemoryDetail.Debit)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}
