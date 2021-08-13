package main

import "time"

type CollectionOperation struct {
	Id                       int32     `json:"id"`
	AccountingMovement       int64     `json:"accountingMovement"`
	AccountingMovementDetail int64     `json:"accountingMovementDetail"`
	Account                  int32     `json:"account"`
	Bank                     *int32    `json:"bank"`
	Status                   string    `json:"status"` // P = Pending, C = Paid, U = Unpaid
	DateCreated              time.Time `json:"dateCreated"`
	DateExpiration           time.Time `json:"dateExpiration"`
	Total                    float32   `json:"total"`
	Paid                     float32   `json:"paid"`
	Pending                  float32   `json:"pending"`
	DocumentName             string    `json:"documentName"`
	PaymentMethod            int16     `json:"paymentMethod"`
	BankName                 string    `json:"bankName"`
	PaymentMethodName        string    `json:"paymentMethodName"`
	AccountName              string    `json:"accountName"`
}

func getPendingColletionOperations() []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	sqlStatement := `SELECT collection_operation.*,(SELECT name FROM account WHERE account.id=collection_operation.bank),(SELECT name FROM payment_method WHERE payment_method.id=collection_operation.payment_method),(SELECT name FROM account WHERE account.id=collection_operation.account) FROM public.collection_operation WHERE status='P' ORDER BY id DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return collectionOperation
	}
	for rows.Next() {
		o := CollectionOperation{}
		rows.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.BankName, &o.PaymentMethodName, &o.AccountName)
		collectionOperation = append(collectionOperation, o)
	}

	return collectionOperation
}

func getColletionOperations(accountingMovement int64) []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=collection_operation.bank),(SELECT name FROM payment_method WHERE payment_method.id=collection_operation.payment_method),(SELECT name FROM account WHERE account.id=collection_operation.account) FROM public.collection_operation WHERE accounting_movement=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, accountingMovement)
	if err != nil {
		log("DB", err.Error())
		return collectionOperation
	}
	for rows.Next() {
		o := CollectionOperation{}
		rows.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.BankName, &o.PaymentMethodName, &o.AccountName)
		collectionOperation = append(collectionOperation, o)
	}

	return collectionOperation
}

func getColletionOperationRow(collectionOperationId int32) CollectionOperation {
	sqlStatement := `SELECT * FROM public.collection_operation WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, collectionOperationId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return CollectionOperation{}
	}

	o := CollectionOperation{}
	row.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod)

	return o
}

func (c *CollectionOperation) insertCollectionOperation() bool {
	if c.Total <= 0 {
		return false
	}

	c.Pending = c.Total
	c.Paid = 0

	p := getPaymentMethodRow(c.PaymentMethod)
	c.DateExpiration = time.Now().AddDate(0, 0, int(p.DaysExpiration))

	sqlStatement := `INSERT INTO public.collection_operation(accounting_movement, accounting_movement_detail, account, bank, date_expiration, total, paid, pending, document_name, payment_method) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	row := db.QueryRow(sqlStatement, c.AccountingMovement, c.AccountingMovementDetail, c.Account, c.Bank, c.DateExpiration, c.Total, c.Paid, c.Pending, c.DocumentName, c.PaymentMethod)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var collectionOperationId int32
	row.Scan(&collectionOperationId)
	c.Id = collectionOperationId

	return collectionOperationId > 0
}

// Adds or substracts the paid quantity on the collection operation
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (c *CollectionOperation) addQuantityCharges(charges float32) bool {
	sqlStatement := `UPDATE public.collection_operation SET paid=paid+$2, pending=pending-$2, status=(CASE WHEN pending-$2=0 THEN 'C' ELSE 'P' END) WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, charges)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

func (c *CollectionOperation) deleteCollectionOperation() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.collection_operation WHERE id=$1`
	_, err := db.Exec(sqlStatement, c.Id)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
