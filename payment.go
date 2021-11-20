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
	Amount                         float64   `json:"amount"`
	Concept                        string    `json:"concept"`
	PaymentTransaction             int32     `json:"paymentTransaction"`
	enterprise                     int32
}

func getPayments(paymentTransaction int32, enterpriseId int32) []Payment {
	payments := make([]Payment, 0)
	sqlStatement := `SELECT * FROM public.payments WHERE payment_transaction=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, paymentTransaction, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return payments
	}

	for rows.Next() {
		p := Payment{}
		rows.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetailDebit, &p.AccountingMovementDetailCredit, &p.Account, &p.DateCreated, &p.Amount, &p.Concept, &p.PaymentTransaction, &p.enterprise)
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
	row.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetailDebit, &p.AccountingMovementDetailCredit, &p.Account, &p.DateCreated, &p.Amount, &p.Concept, &p.PaymentTransaction, &p.enterprise)

	return p
}

func (c *Payment) isValid() bool {
	return !(c.PaymentTransaction <= 0 || len(c.Concept) > 50 || c.Amount <= 0)
}

func (c *Payment) insertPayment(userId int32) bool {
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
	if pt.Id <= 0 || pt.Bank == nil || pt.enterprise != c.enterprise || pt.Pending <= 0 {
		trans.Rollback()
		return false
	}

	c.AccountingMovement = pt.AccountingMovement
	c.AccountingMovementDetailDebit = pt.AccountingMovementDetail
	c.Account = pt.Account

	ok := pt.addQuantityCharges(c.Amount, userId)
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
	m.enterprise = c.enterprise
	ok = m.insertAccountingMovement(userId)
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
	dInc.enterprise = c.enterprise
	ok = dInc.insertAccountingMovementDetail(userId)
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
	Supp.enterprise = c.enterprise
	ok = Supp.insertAccountingMovementDetail(userId)
	if !ok {
		trans.Rollback()
		return false
	}
	c.AccountingMovementDetailCredit = Supp.Id

	// insert row
	sqlStatement := `INSERT INTO public.payments(accounting_movement, accounting_movement_detail_debit, accounting_movement_detail_credit, account, amount, concept, payment_transaction, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	row := db.QueryRow(sqlStatement, c.AccountingMovement, c.AccountingMovementDetailDebit, c.AccountingMovementDetailCredit, c.Account, c.Amount, c.Concept, c.PaymentTransaction, c.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false
	}

	var paymentId int32
	row.Scan(&paymentId)
	c.Id = paymentId

	insertTransactionalLog(c.enterprise, "payments", int(c.Id), userId, "I")

	///
	err = trans.Commit()
	return err == nil
	///
}

func (c *Payment) deletePayment(userId int32) bool {
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
	if inMemoryPayment.Id <= 0 || inMemoryPayment.enterprise != c.enterprise {
		trans.Rollback()
		return false
	}
	// get the payment transaction
	pt := getPaymentTransactionRow(inMemoryPayment.PaymentTransaction)
	if pt.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(c.enterprise, "payments", int(c.Id), userId, "D")

	sqlStatement := `DELETE FROM public.payments WHERE id=$1 AND enterprise=$2`
	_, err = db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// substract the paid amount
	ok := pt.addQuantityCharges(-inMemoryPayment.Amount, userId)
	if !ok {
		trans.Rollback()
		return false
	}

	// delete the associated account movement (credit)
	amd := getAccountingMovementDetailRow(inMemoryPayment.AccountingMovementDetailCredit)
	am := getAccountingMovementRow(amd.Movement)
	ok = am.deleteAccountingMovement(userId)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}
