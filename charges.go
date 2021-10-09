package main

import (
	"time"
)

type Charges struct {
	Id                             int32     `json:"id"`
	AccountingMovement             int64     `json:"accountingMovement"`
	AccountingMovementDetailDebit  int64     `json:"accountingMovementDetailDebit"`
	AccountingMovementDetailCredit int64     `json:"accountingMovementDetailCredit"`
	Account                        int32     `json:"account"`
	DateCreated                    time.Time `json:"dateCreated"`
	Amount                         float64   `json:"amount"`
	Concept                        string    `json:"concept"`
	CollectionOperation            int32     `json:"collectionOperation"`
	enterprise                     int32
}

func getCharges(collectionOperation int32, enterpriseId int32) []Charges {
	charges := make([]Charges, 0)
	sqlStatement := `SELECT * FROM public.charges WHERE collection_operation=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, collectionOperation, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return charges
	}

	for rows.Next() {
		c := Charges{}
		rows.Scan(&c.Id, &c.AccountingMovement, &c.AccountingMovementDetailDebit, &c.AccountingMovementDetailCredit, &c.Account, &c.DateCreated, &c.Amount, &c.Concept, &c.CollectionOperation, &c.enterprise)
		charges = append(charges, c)
	}
	return charges
}

func getChargesRow(chargesId int32) Charges {
	sqlStatement := `SELECT * FROM public.charges WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, chargesId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Charges{}
	}

	c := Charges{}
	row.Scan(&c.Id, &c.AccountingMovement, &c.AccountingMovementDetailDebit, &c.AccountingMovementDetailCredit, &c.Account, &c.DateCreated, &c.Amount, &c.Concept, &c.CollectionOperation, &c.enterprise)

	return c
}

func (c *Charges) isValid() bool {
	return !(c.CollectionOperation <= 0 || len(c.Concept) > 50 || c.Amount <= 0)
}

func (c *Charges) insertCharges() bool {
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

	// get data from collection operation
	co := getColletionOperationRow(c.CollectionOperation)
	if co.Id <= 0 || co.enterprise != c.enterprise || co.Bank == nil || co.Pending <= 0 {
		trans.Rollback()
		return false
	}

	c.AccountingMovement = co.AccountingMovement
	c.AccountingMovementDetailDebit = co.AccountingMovementDetail
	c.Account = co.Account

	ok := co.addQuantityCharges(c.Amount)
	if !ok {
		trans.Rollback()
		return false
	}

	am := getAccountingMovementRow(co.AccountingMovement)
	if am.Id <= 0 {
		trans.Rollback()
		return false
	}

	// insert accounting movement for the payment
	m := AccountingMovement{}
	m.Type = "N"
	m.BillingSerie = am.BillingSerie
	m.enterprise = c.enterprise
	ok = m.insertAccountingMovement()
	if !ok {
		trans.Rollback()
		return false
	}

	// 1. debit detail for the bank
	bank := getAccountRow(*co.Bank)

	dInc := AccountingMovementDetail{}
	dInc.Movement = m.Id
	dInc.Journal = bank.Journal
	dInc.AccountNumber = bank.AccountNumber
	dInc.Debit = c.Amount
	dInc.Type = "N"
	dInc.PaymentMethod = co.PaymentMethod
	dInc.enterprise = c.enterprise
	ok = dInc.insertAccountingMovementDetail()
	if !ok {
		trans.Rollback()
		return false
	}

	// 2. credit detail for the customer's account
	dCustDebit := getAccountingMovementDetailRow(c.AccountingMovementDetailDebit)
	if dCustDebit.Id <= 0 {
		trans.Rollback()
		return false
	}

	dCust := AccountingMovementDetail{}
	dCust.Movement = m.Id
	dCust.Journal = dCustDebit.Journal
	dCust.AccountNumber = dCustDebit.AccountNumber
	dCust.Credit = c.Amount
	dCust.Type = "N"
	dCust.DocumentName = dCustDebit.DocumentName
	dCust.PaymentMethod = co.PaymentMethod
	dCust.enterprise = c.enterprise
	ok = dCust.insertAccountingMovementDetail()
	if !ok {
		trans.Rollback()
		return false
	}
	c.AccountingMovementDetailCredit = dCust.Id

	// insert row
	sqlStatement := `INSERT INTO public.charges(accounting_movement, accounting_movement_detail_debit, accounting_movement_detail_credit, account, amount, concept, collection_operation, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(sqlStatement, c.AccountingMovement, c.AccountingMovementDetailDebit, c.AccountingMovementDetailCredit, c.Account, c.Amount, c.Concept, c.CollectionOperation, c.enterprise)
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

func (c *Charges) deleteCharges() bool {
	if c.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	inMemoryCharge := getChargesRow(c.Id)
	if inMemoryCharge.Id <= 0 || inMemoryCharge.enterprise != c.enterprise {
		trans.Rollback()
		return false
	}
	// get the collection operation
	co := getColletionOperationRow(inMemoryCharge.CollectionOperation)
	if co.Id <= 0 {
		trans.Rollback()
		return false
	}

	sqlStatement := `DELETE FROM public.charges WHERE id=$1 AND enterprise=$2`
	_, err = db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// substract the paid amount
	ok := co.addQuantityCharges(-inMemoryCharge.Amount)
	if !ok {
		trans.Rollback()
		return false
	}

	// delete the associated account movement (credit)
	amd := getAccountingMovementDetailRow(inMemoryCharge.AccountingMovementDetailCredit)
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
