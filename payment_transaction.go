package main

import (
	"strconv"
	"time"
)

type PaymentTransaction struct {
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
	SupplierName             *string   `json:"supplierName"`
	enterprise               int32
}

func getPendingPaymentTransaction(enterpriseId int32) []PaymentTransaction {
	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)
	sqlStatement := `SELECT payment_transaction.*,(SELECT name FROM account WHERE account.id=payment_transaction.bank),(SELECT name FROM payment_method WHERE payment_method.id=payment_transaction.payment_method),(SELECT name FROM account WHERE account.id=payment_transaction.account),suppliers.name FROM public.payment_transaction FULL JOIN purchase_invoice ON purchase_invoice.accounting_movement=payment_transaction.accounting_movement FULL JOIN suppliers ON suppliers.id=purchase_invoice.supplier WHERE status='P' AND payment_transaction.enterprise=$1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return paymentTransaction
	}
	for rows.Next() {
		p := PaymentTransaction{}
		rows.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetail, &p.Account, &p.Bank, &p.Status, &p.DateCreated, &p.DateExpiration, &p.Total, &p.Paid, &p.Pending, &p.DocumentName, &p.PaymentMethod, &p.enterprise, &p.BankName, &p.PaymentMethodName, &p.AccountName, &p.SupplierName)
		paymentTransaction = append(paymentTransaction, p)
	}

	return paymentTransaction
}

func searchPaymentTransactions(search CollectionOperationPaymentTransactionSearch, enterpriseId int32) []PaymentTransaction {
	if search.isDefault() {
		return getPendingPaymentTransaction(enterpriseId)
	}

	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)
	sqlStatement := `SELECT payment_transaction.*,(SELECT name FROM account WHERE account.id=payment_transaction.bank),(SELECT name FROM payment_method WHERE payment_method.id=payment_transaction.payment_method),(SELECT name FROM account WHERE account.id=payment_transaction.account),suppliers.name FROM public.payment_transaction FULL JOIN purchase_invoice ON purchase_invoice.accounting_movement=payment_transaction.accounting_movement FULL JOIN suppliers ON suppliers.id=purchase_invoice.supplier WHERE payment_transaction.enterprise=$1`
	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, enterpriseId)

	if search.Mode != 0 {
		sqlStatement += ` AND payment_transaction.status=$2`
		if search.Mode == 1 {
			interfaces = append(interfaces, "P") // Pending
		} else if search.Mode == 2 {
			interfaces = append(interfaces, "C") // Paid
		} else if search.Mode == 3 {
			interfaces = append(interfaces, "U") // Unpaid
		}
	}

	if search.StartDate != nil {
		sqlStatement += ` AND payment_transaction.date_created >= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.StartDate)
	}

	if search.EndDate != nil {
		sqlStatement += ` AND payment_transaction.date_created <= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.EndDate)
	}

	if len(search.Search) > 0 {
		sqlStatement += ` AND suppliers.name ILIKE $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, "%"+search.Search+"%")
	}

	sqlStatement += ` ORDER BY payment_transaction.id DESC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return paymentTransaction
	}
	for rows.Next() {
		p := PaymentTransaction{}
		rows.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetail, &p.Account, &p.Bank, &p.Status, &p.DateCreated, &p.DateExpiration, &p.Total, &p.Paid, &p.Pending, &p.DocumentName, &p.PaymentMethod, &p.enterprise, &p.BankName, &p.PaymentMethodName, &p.AccountName, &p.SupplierName)
		paymentTransaction = append(paymentTransaction, p)
	}

	return paymentTransaction
}

func getPaymentTransactions(accountingMovement int64, enterpriseId int32) []PaymentTransaction {
	var paymentTransaction []PaymentTransaction = make([]PaymentTransaction, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=payment_transaction.bank),(SELECT name FROM payment_method WHERE payment_method.id=payment_transaction.payment_method),(SELECT name FROM account WHERE account.id=payment_transaction.account) FROM public.payment_transaction WHERE accounting_movement=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, accountingMovement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return paymentTransaction
	}
	for rows.Next() {
		p := PaymentTransaction{}
		rows.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetail, &p.Account, &p.Bank, &p.Status, &p.DateCreated, &p.DateExpiration, &p.Total, &p.Paid, &p.Pending, &p.DocumentName, &p.PaymentMethod, &p.enterprise, &p.BankName, &p.PaymentMethodName, &p.AccountName)
		paymentTransaction = append(paymentTransaction, p)
	}

	return paymentTransaction
}

func getPaymentTransactionRow(paymentTransactionId int32) PaymentTransaction {
	sqlStatement := `SELECT * FROM public.payment_transaction WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, paymentTransactionId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PaymentTransaction{}
	}

	p := PaymentTransaction{}
	row.Scan(&p.Id, &p.AccountingMovement, &p.AccountingMovementDetail, &p.Account, &p.Bank, &p.Status, &p.DateCreated, &p.DateExpiration, &p.Total, &p.Paid, &p.Pending, &p.DocumentName, &p.PaymentMethod, &p.enterprise)

	return p
}

func (c *PaymentTransaction) insertPaymentTransaction() bool {
	if c.Total <= 0 {
		return false
	}

	c.Pending = c.Total
	c.Paid = 0

	p := getPaymentMethodRow(c.PaymentMethod)
	c.DateExpiration = time.Now().AddDate(0, 0, int(p.DaysExpiration))

	sqlStatement := `INSERT INTO public.payment_transaction(accounting_movement, accounting_movement_detail, account, bank, date_expiration, total, paid, pending, document_name, payment_method, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	row := db.QueryRow(sqlStatement, c.AccountingMovement, c.AccountingMovementDetail, c.Account, c.Bank, c.DateExpiration, c.Total, c.Paid, c.Pending, c.DocumentName, c.PaymentMethod, c.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var collectionOperationId int32
	row.Scan(&collectionOperationId)
	c.Id = collectionOperationId

	return collectionOperationId > 0
}

// Adds or substracts the paid quantity on the payment transaction
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (c *PaymentTransaction) addQuantityCharges(charges float64) bool {
	sqlStatement := `UPDATE public.payment_transaction SET paid=paid+$2, pending=pending-$2, status=(CASE WHEN pending-$2=0 THEN 'C' ELSE 'P' END) WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, charges)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

func (c *PaymentTransaction) deletePaymentTransaction() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.payment_transaction WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, c.Id, c.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
