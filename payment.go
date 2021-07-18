package main

import (
	"time"
)

type Payment struct {
	Id                             int32     `json:"id"`
	AccountingMovement             int64     `json:"accountingMovement"`
	AccountingMovementDetailDebit  int64     `json:"accountingMovementDetailDebit"`
	AccountingMovementDetailCredit int64     `json:"accountingMovementDetailCredit"`
	Account                        int32     `json:"account"`
	DateCreated                    time.Time `json:"dateCreated"`
	Amount                         float32   `json:"amount"`
	Concept                        string    `json:"concept"`
	PaymentTransaction             int32     `json:"paymentTransaction"`
}

func getPayments(paymentTransaction int32) []Payment {
	payments := make([]Payment, 0)
	sqlStatement := `SELECT * FROM public.payments WHERE payment_transaction=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, paymentTransaction)
	if err != nil {
		log("DB", err.Error())
		return payments
	}

	for rows.Next() {
		p := Payment{}
		rows.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetailDebit, &p.AccountingMovementDetailCredit, &p.Account, &p.DateCreated, &p.Amount, &p.Concept, &p.PaymentTransaction)
		payments = append(payments, p)
	}
	return payments
}

func getPaymentsRow(chargesId int32) Payment {
	sqlStatement := `SELECT * FROM public.payments WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, chargesId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Payment{}
	}

	p := Payment{}
	row.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetailDebit, &p.AccountingMovementDetailCredit, &p.Account, &p.DateCreated, &p.Amount, &p.Concept, &p.PaymentTransaction)

	return p
}

func (c *Payment) isValid() bool {
	return !(c.PaymentTransaction <= 0 || len(c.Concept) > 50 || c.Amount <= 0)
}

func (c *Payment) insertPayment() bool {
	// validation
	if !c.isValid() {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	// get data from payment transaction
	pt := getPaymentTransactionRow(c.PaymentTransaction)
	if pt.Id <= 0 || pt.Bank == nil || pt.Pending <= 0 {
		trans.Rollback()
		return false
	}

	c.AccountingMovement = pt.AccountingMovement
	c.AccountingMovementDetailDebit = pt.AccountingMovementDetail
	c.Account = pt.Account

	ok := pt.addQuantityCharges(c.Amount)
	if !ok {
		trans.Rollback()
		return false
	}

	am := getAccountingMovementRow(pt.AccountingMovement)
	if am.Id <= 0 {
		trans.Rollback()
		return false
	}

	// insert accounting movement for the payment
	m := AccountingMovement{}
	m.Type = "N"
	m.BillingSerie = am.BillingSerie
	ok = m.insertAccountingMovement()
	if !ok {
		trans.Rollback()
		return false
	}

	// 1. debit detail for the bank
	bank := getAccountRow(*pt.Bank)

	dInc := AccountingMovementDetail{}
	dInc.Movement = m.Id
	dInc.Journal = bank.Journal
	dInc.AccountNumber = bank.AccountNumber
	dInc.Credit = c.Amount
	dInc.Type = "N"
	dInc.PaymentMethod = pt.PaymentMethod
	ok = dInc.insertAccountingMovementDetail()
	if !ok {
		trans.Rollback()
		return false
	}

	// 2. credit detail for the suppliers's account
	dSuppDebit := getAccountingMovementDetailRow(c.AccountingMovementDetailDebit)
	if dSuppDebit.Id <= 0 {
		trans.Rollback()
		return false
	}

	Supp := AccountingMovementDetail{}
	Supp.Movement = m.Id
	Supp.Journal = dSuppDebit.Journal
	Supp.AccountNumber = dSuppDebit.AccountNumber
	Supp.Debit = c.Amount
	Supp.Type = "N"
	Supp.DocumentName = dSuppDebit.DocumentName
	Supp.PaymentMethod = pt.PaymentMethod
	ok = Supp.insertAccountingMovementDetail()
	if !ok {
		trans.Rollback()
		return false
	}
	c.AccountingMovementDetailCredit = Supp.Id

	// insert row
	sqlStatement := `INSERT INTO public.payments(accounting_movement, accounting_movement_detail_debit, accounting_movement_detail_credit, account, amount, concept, payment_transaction) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(sqlStatement, c.AccountingMovement, c.AccountingMovementDetailDebit, c.AccountingMovementDetailCredit, c.Account, c.Amount, c.Concept, c.PaymentTransaction)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}

func (c *Payment) deletePayment() bool {
	if c.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	inMemoryPayment := getPaymentsRow(c.Id)
	if inMemoryPayment.Id <= 0 {
		trans.Rollback()
		return false
	}
	// get the payment transaction
	pt := getPaymentTransactionRow(inMemoryPayment.PaymentTransaction)
	if pt.Id <= 0 {
		trans.Rollback()
		return false
	}

	sqlStatement := `DELETE FROM public.payments WHERE id=$1`
	_, err = db.Exec(sqlStatement, c.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// substract the paid amount
	ok := pt.addQuantityCharges(-inMemoryPayment.Amount)
	if !ok {
		trans.Rollback()
		return false
	}

	// delete the associated account movement (credit)
	amd := getAccountingMovementDetailRow(inMemoryPayment.AccountingMovementDetailCredit)
	am := getAccountingMovementRow(amd.Movement)
	ok = am.deleteAccountingMovement()
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}
