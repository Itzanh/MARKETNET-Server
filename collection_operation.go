package main

import (
	"strconv"
	"time"
)

type CollectionOperation struct {
	Id                       int32     `json:"id"`
	AccountingMovement       int64     `json:"accountingMovement"`
	AccountingMovementDetail int64     `json:"accountingMovementDetail"`
	Account                  int32     `json:"account"`
	Bank                     *int32    `json:"bank"`
	Status                   string    `json:"status"` // P = Pending, C = Paid, U = Unpaid
	DateCreated              time.Time `json:"dateCreated"`
	DateExpiration           time.Time `json:"dateExpiration"`
	Total                    float64   `json:"total"`
	Paid                     float64   `json:"paid"`
	Pending                  float64   `json:"pending"`
	DocumentName             string    `json:"documentName"`
	PaymentMethod            int32     `json:"paymentMethod"`
	BankName                 string    `json:"bankName"`
	PaymentMethodName        string    `json:"paymentMethodName"`
	AccountName              string    `json:"accountName"`
	CustomerName             *string   `json:"customerName"`
	enterprise               int32
}

type CollectionOperationPaymentTransactionSearch struct {
	Mode      uint8      `json:"mode"` // 0 = All, 1 = Pending, 2 = Paid, 3 = Unpaid
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Search    string     `json:"search"` // Customer / supplier name
}

func (search *CollectionOperationPaymentTransactionSearch) isDefault() bool {
	return search.Mode == 1 && search.StartDate == nil && search.EndDate == nil && len(search.Search) == 0
}

func getPendingColletionOperations(enterpriseId int32) []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	sqlStatement := `SELECT collection_operation.*,(SELECT name FROM account WHERE account.id=collection_operation.bank),(SELECT name FROM payment_method WHERE payment_method.id=collection_operation.payment_method),(SELECT name FROM account WHERE account.id=collection_operation.account),customer.name FROM public.collection_operation FULL JOIN sales_invoice ON sales_invoice.accounting_movement=collection_operation.accounting_movement FULL JOIN customer ON customer.id=sales_invoice.customer WHERE status='P' AND collection_operation.enterprise=$1 ORDER BY collection_operation.id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return collectionOperation
	}
	defer rows.Close()

	for rows.Next() {
		o := CollectionOperation{}
		rows.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.enterprise, &o.BankName, &o.PaymentMethodName, &o.AccountName, &o.CustomerName)
		collectionOperation = append(collectionOperation, o)
	}

	return collectionOperation
}

func searchCollectionOperations(search CollectionOperationPaymentTransactionSearch, enterpriseId int32) []CollectionOperation {
	if search.isDefault() {
		return getPendingColletionOperations(enterpriseId)
	}

	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	sqlStatement := `SELECT collection_operation.*,(SELECT name FROM account WHERE account.id=collection_operation.bank),(SELECT name FROM payment_method WHERE payment_method.id=collection_operation.payment_method),(SELECT name FROM account WHERE account.id=collection_operation.account),customer.name FROM public.collection_operation FULL JOIN sales_invoice ON sales_invoice.accounting_movement=collection_operation.accounting_movement FULL JOIN customer ON customer.id=sales_invoice.customer WHERE collection_operation.enterprise=$1`
	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, enterpriseId)

	if search.Mode != 0 {
		sqlStatement += ` AND collection_operation.status=$2`
		if search.Mode == 1 {
			interfaces = append(interfaces, "P") // Pending
		} else if search.Mode == 2 {
			interfaces = append(interfaces, "C") // Paid
		} else if search.Mode == 3 {
			interfaces = append(interfaces, "U") // Unpaid
		}
	}

	if search.StartDate != nil {
		sqlStatement += ` AND collection_operation.date_created >= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.StartDate)
	}

	if search.EndDate != nil {
		sqlStatement += ` AND collection_operation.date_created <= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.EndDate)
	}

	if len(search.Search) > 0 {
		sqlStatement += ` AND customer.name ILIKE $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, "%"+search.Search+"%")
	}

	sqlStatement += ` ORDER BY collection_operation.id DESC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return collectionOperation
	}
	defer rows.Close()

	for rows.Next() {
		o := CollectionOperation{}
		rows.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.enterprise, &o.BankName, &o.PaymentMethodName, &o.AccountName, &o.CustomerName)
		collectionOperation = append(collectionOperation, o)
	}

	return collectionOperation
}

func getColletionOperations(accountingMovement int64, enterpriseId int32) []CollectionOperation {
	var collectionOperation []CollectionOperation = make([]CollectionOperation, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=collection_operation.bank),(SELECT name FROM payment_method WHERE payment_method.id=collection_operation.payment_method),(SELECT name FROM account WHERE account.id=collection_operation.account) FROM public.collection_operation WHERE accounting_movement=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, accountingMovement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return collectionOperation
	}
	defer rows.Close()

	for rows.Next() {
		o := CollectionOperation{}
		rows.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.enterprise, &o.BankName, &o.PaymentMethodName, &o.AccountName)
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
	row.Scan(&o.Id, &o.AccountingMovement, &o.AccountingMovementDetail, &o.Account, &o.Bank, &o.Status, &o.DateCreated, &o.DateExpiration, &o.Total, &o.Paid, &o.Pending, &o.DocumentName, &o.PaymentMethod, &o.enterprise)

	return o
}

func (c *CollectionOperation) insertCollectionOperation(userId int32) bool {
	if c.Total <= 0 {
		return false
	}

	c.Pending = c.Total
	c.Paid = 0

	p := getPaymentMethodRow(c.PaymentMethod)
	c.DateExpiration = time.Now().AddDate(0, 0, int(p.DaysExpiration))

	sqlStatement := `INSERT INTO public.collection_operation(accounting_movement, accounting_movement_detail, account, bank, date_expiration, total, paid, pending, document_name, payment_method, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	row := db.QueryRow(sqlStatement, c.AccountingMovement, c.AccountingMovementDetail, c.Account, c.Bank, c.DateExpiration, c.Total, c.Paid, c.Pending, c.DocumentName, c.PaymentMethod, c.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var collectionOperationId int32
	row.Scan(&collectionOperationId)
	c.Id = collectionOperationId

	if collectionOperationId > 0 {
		insertTransactionalLog(c.enterprise, "collection_operation", int(c.Id), userId, "I")
	}

	return collectionOperationId > 0
}

// Adds or substracts the paid quantity on the collection operation
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (c *CollectionOperation) addQuantityCharges(charges float64, userId int32) bool {
	sqlStatement := `UPDATE public.collection_operation SET paid=paid+$2, pending=pending-$2, status=(CASE WHEN pending-$2=0 THEN 'C' ELSE 'P' END) WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, charges)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	insertTransactionalLog(c.enterprise, "collection_operation", int(c.Id), userId, "U")

	return err == nil && rows > 0
}

func (c *CollectionOperation) deleteCollectionOperation(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	insertTransactionalLog(c.enterprise, "collection_operation", int(c.Id), userId, "D")

	sqlStatement := `DELETE FROM public.collection_operation WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, c.Id, c.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
