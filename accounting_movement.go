package main

import (
	"fmt"
	"sort"
	"time"
)

type AccountingMovement struct {
	Id               int64     `json:"id"`
	DateCreated      time.Time `json:"dateCreated"`
	AmountDebit      float64   `json:"amountDebit"`
	AmountCredit     float64   `json:"amountCredit"`
	FiscalYear       int16     `json:"fiscalYear"`
	Type             string    `json:"type"` // O: Opening, N: Normal, V: Variation of existences, R: Regularisation, C: Closing
	BillingSerie     string    `json:"billingSerie"`
	BillingSerieName string    `json:"billingSerieName"`
	enterprise       int32
}

func getAccountingMovement(enterpriseId int32) []AccountingMovement {
	accountingMovements := make([]AccountingMovement, 0)
	sqlStatement := `SELECT *,(SELECT name FROM billing_series WHERE billing_series.id=accounting_movement.billing_serie AND billing_series.enterprise=accounting_movement.enterprise) FROM public.accounting_movement WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		fmt.Println(err)
		log("DB", err.Error())
		return accountingMovements
	}

	for rows.Next() {
		a := AccountingMovement{}
		rows.Scan(&a.Id, &a.DateCreated, &a.AmountDebit, &a.AmountCredit, &a.FiscalYear, &a.Type, &a.BillingSerie, &a.enterprise, &a.BillingSerieName)
		accountingMovements = append(accountingMovements, a)
	}

	return accountingMovements
}

func searchAccountingMovements(search string, enterpriseId int32) []AccountingMovement {
	accountingMovements := make([]AccountingMovement, 0)
	sqlStatement := `SELECT DISTINCT accounting_movement.*,(SELECT name FROM billing_series WHERE billing_series.id=accounting_movement.billing_serie AND billing_series.enterprise=accounting_movement.enterprise) FROM accounting_movement INNER JOIN accounting_movement_detail ON accounting_movement_detail.movement=accounting_movement.id LEFT JOIN sales_invoice ON sales_invoice.accounting_movement=accounting_movement.id LEFT JOIN customer ON customer.id=sales_invoice.customer LEFT JOIN purchase_invoice ON purchase_invoice.accounting_movement=accounting_movement.id LEFT JOIN suppliers ON suppliers.id=purchase_invoice.supplier WHERE ((accounting_movement_detail.document_name ILIKE $1) OR (customer.name ILIKE $1) OR (suppliers.name ILIKE $1)) AND (accounting_movement.enterprise=$2) ORDER BY accounting_movement.date_created DESC`
	rows, err := db.Query(sqlStatement, "%"+search+"%", enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return accountingMovements
	}

	for rows.Next() {
		a := AccountingMovement{}
		rows.Scan(&a.Id, &a.DateCreated, &a.AmountDebit, &a.AmountCredit, &a.FiscalYear, &a.Type, &a.BillingSerie, &a.enterprise, &a.BillingSerieName)
		accountingMovements = append(accountingMovements, a)
	}

	return accountingMovements
}

func getAccountingMovementRow(accountingMovementId int64) AccountingMovement {
	sqlStatement := `SELECT * FROM public.accounting_movement WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, accountingMovementId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return AccountingMovement{}
	}

	a := AccountingMovement{}
	row.Scan(&a.Id, &a.DateCreated, &a.AmountDebit, &a.AmountCredit, &a.FiscalYear, &a.Type, &a.BillingSerie, &a.enterprise)

	return a
}

func (a *AccountingMovement) isValid() bool {
	return !((a.Type != "O" && a.Type != "N" && a.Type != "V" && a.Type != "R" && a.Type != "C") || len(a.BillingSerie) != 3)
}

func (a *AccountingMovement) insertAccountingMovement() bool {
	if !a.isValid() {
		return false
	}

	a.FiscalYear = int16(time.Now().Year())

	sqlStatement := `INSERT INTO public.accounting_movement(fiscal_year, type, billing_serie, enterprise) VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(sqlStatement, a.FiscalYear, a.Type, a.BillingSerie, a.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var movementId int64
	row.Scan(&movementId)
	a.Id = movementId

	return movementId > 0
}

func (a *AccountingMovement) deleteAccountingMovement() bool {
	if a.Id <= 0 {
		return false
	}

	accountingMovementInMemory := getAccountingMovementRow(a.Id)
	settings := getSettingsRecordById(accountingMovementInMemory.enterprise)
	if accountingMovementInMemory.Id <= 0 || accountingMovementInMemory.enterprise != a.enterprise || (settings.LimitAccountingDate != nil && accountingMovementInMemory.DateCreated.Before(*settings.LimitAccountingDate)) {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	// cascade delete the collection operations
	c := getColletionOperations(a.Id, a.enterprise)
	for i := 0; i < len(c); i++ {
		ok := c[i].deleteCollectionOperation()
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// cascade delete the payment transactions
	p := getPaymentTransactions(a.Id, a.enterprise)
	for i := 0; i < len(p); i++ {
		ok := p[i].deletePaymentTransaction()
		if !ok {
			trans.Rollback()
			return false
		}
	}
	// cascade delete the details
	d := getAccountingMovementDetail(a.Id, a.enterprise)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteAccountingMovementDetail()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// set the relation null on the sale invoices
	sqlStatement := `UPDATE sales_invoice SET accounting_movement=NULL WHERE accounting_movement=$1`
	_, err = db.Exec(sqlStatement, a.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	// set the relation null on the purchase invoices
	sqlStatement = `UPDATE purchase_invoice SET accounting_movement=NULL WHERE accounting_movement=$1`
	_, err = db.Exec(sqlStatement, a.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// delete the movement
	sqlStatement = `DELETE FROM public.accounting_movement WHERE id=$1`
	_, err = db.Exec(sqlStatement, a.Id)
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

// Will add or take out credit and debit (if given a negative amount)
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (a *AccountingMovement) addCreditAndDebit(credit float64, debit float64) bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.accounting_movement SET amount_debit=amount_debit+$2, amount_credit=amount_credit+$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, a.Id, debit, credit)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

type PostInvoiceResult struct {
	Invoice int64 `json:"invoice"`
	Ok      bool  `json:"ok"`
	Result  int8  `json:"result"` // 0: Internal error, 1: The customer/supplier in the invoice has no account
}

// Transfer the sales invoices from management to accounting. Create the movements and the details for all the selected invoices.
func salesPostInvoices(invoiceIds []int64, enterpriseId int32) []PostInvoiceResult {
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
	if settings.SalesJournal == nil {
		return result
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return result
	}
	///

	for i := 0; i < len(invoiceIds); i++ {
		// get the selected invoice
		inv := getSalesInvoiceRow(invoiceIds[i])
		if inv.Id <= 0 || inv.enterprise != enterpriseId {
			trans.Rollback()
			return result
		}
		// get the invoice customer
		c := getCustomerRow(inv.Customer)
		if c.Id <= 0 {
			trans.Rollback()
			return result
		}
		if c.Account == nil {
			result[i].Result = 1
			continue
		}
		// get the account row
		a := getAccountRow(*c.Account)
		if a.Id <= 0 {
			trans.Rollback()
			return result
		}

		// create the accounting movement
		m := AccountingMovement{}
		m.Type = "N"
		m.BillingSerie = inv.BillingSeries
		m.enterprise = inv.enterprise
		ok := m.insertAccountingMovement()
		if !ok {
			trans.Rollback()
			return result
		}

		// create the details
		// 1. detail for the customer
		dCust := AccountingMovementDetail{}
		dCust.Movement = m.Id
		dCust.Journal = a.Journal
		dCust.AccountNumber = a.AccountNumber
		dCust.Debit = inv.TotalAmount
		dCust.Type = "N"
		dCust.DocumentName = inv.InvoiceName
		dCust.PaymentMethod = inv.PaymentMethod
		dCust.enterprise = enterpriseId
		ok = dCust.insertAccountingMovementDetail()
		if !ok {
			trans.Rollback()
			return result
		}

		// create the collection operation for this charge
		if inv.TotalAmount > 0 {
			p := getPaymentMethodRow(inv.PaymentMethod)

			co := CollectionOperation{}
			co.AccountingMovement = m.Id
			co.AccountingMovementDetail = dCust.Id
			co.Account = getAccountIdByAccountNumber(a.Journal, a.AccountNumber, a.enterprise)
			co.Total = inv.TotalAmount
			co.DocumentName = inv.InvoiceName
			co.PaymentMethod = inv.PaymentMethod
			co.Bank = p.Bank
			co.enterprise = enterpriseId
			ok := co.insertCollectionOperation()
			if !ok {
				trans.Rollback()
				return result
			}

			// paid in advance
			if p.PaidInAdvance {
				ch := Charges{}
				ch.CollectionOperation = co.Id
				ch.Amount = inv.TotalAmount
				ch.enterprise = enterpriseId
				ok = ch.insertCharges()
				if !ok {
					trans.Rollback()
					return result
				}
			}
		}

		// 2. details for the income
		dInc := AccountingMovementDetail{}
		dInc.Movement = m.Id
		dInc.Journal = *settings.SalesJournal
		dInc.AccountNumber = 1
		dInc.Credit = inv.TotalWithDiscount
		dInc.Type = "N"
		dInc.DocumentName = inv.InvoiceName
		dInc.PaymentMethod = inv.PaymentMethod
		dInc.enterprise = enterpriseId
		ok = dInc.insertAccountingMovementDetail()
		if !ok {
			trans.Rollback()
			return result
		}

		// 3. details for the VAT

		// get the details and sort
		det := getSalesInvoiceDetail(inv.Id, inv.enterprise)
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
				dVat.Movement = m.Id
				dVat.Journal = journal
				dVat.AccountNumber = accountNumber
				dVat.Credit = credit
				dVat.Type = "N"
				dVat.DocumentName = inv.InvoiceName
				dVat.PaymentMethod = inv.PaymentMethod
				dVat.enterprise = enterpriseId
				ok := dVat.insertAccountingMovementDetail()
				if !ok {
					trans.Rollback()
					return result
				}

				details = make([]SalesInvoiceDetail, 0)
				details = append(details, d[i])
			}
		}

		// set the accounting movement on the invoice
		sqlStatement := `UPDATE sales_invoice SET accounting_movement=$2 WHERE id=$1`
		_, err := db.Exec(sqlStatement, invoiceIds[i], m.Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return result
		}
		result[i].Ok = true
	}

	///
	err = trans.Commit()
	if err != nil {
		for i := 0; i < len(result); i++ {
			result[i].Ok = false
			result[i].Result = 0
		}
	}
	///
	return result
}

// Transfer the purchase invoices from management to accounting. Create the movements and the details for all the selected invoices.
func purchasePostInvoices(invoiceIds []int64, enterpriseId int32) []PostInvoiceResult {
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
	if settings.PurchaseJournal == nil {
		return result
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return result
	}
	///

	for i := 0; i < len(invoiceIds); i++ {
		// get the selected invoice
		inv := getPurchaseInvoiceRow(invoiceIds[i])
		if inv.Id <= 0 || inv.enterprise != enterpriseId {
			trans.Rollback()
			return result
		}
		// get the invoice customer
		s := getSupplierRow(inv.Supplier)
		if s.Id <= 0 {
			trans.Rollback()
			return result
		}
		if s.Account == nil {
			result[i].Result = 1
			continue
		}
		// get the account row
		a := getAccountRow(*s.Account)
		if a.Id <= 0 {
			trans.Rollback()
			return result
		}

		// create the accounting movement
		m := AccountingMovement{}
		m.Type = "N"
		m.BillingSerie = inv.BillingSeries
		m.enterprise = inv.enterprise
		ok := m.insertAccountingMovement()
		if !ok {
			trans.Rollback()
			return result
		}

		// create the details
		// 1. detail for the supplier
		dCust := AccountingMovementDetail{}
		dCust.Movement = m.Id
		dCust.Journal = a.Journal
		dCust.AccountNumber = a.AccountNumber
		dCust.Credit = inv.TotalAmount
		dCust.Type = "N"
		dCust.DocumentName = inv.InvoiceName
		dCust.PaymentMethod = inv.PaymentMethod
		dCust.enterprise = enterpriseId
		ok = dCust.insertAccountingMovementDetail()
		if !ok {
			trans.Rollback()
			return result
		}

		// create the payment transaction for this payment
		if inv.TotalAmount > 0 {
			p := getPaymentMethodRow(inv.PaymentMethod)

			pt := PaymentTransaction{}
			pt.AccountingMovement = m.Id
			pt.AccountingMovementDetail = dCust.Id
			pt.Account = getAccountIdByAccountNumber(a.Journal, a.AccountNumber, a.enterprise)
			pt.Total = inv.TotalAmount
			pt.DocumentName = inv.InvoiceName
			pt.PaymentMethod = *s.PaymentMethod
			pt.Bank = p.Bank
			pt.enterprise = enterpriseId
			ok := pt.insertPaymentTransaction()
			if !ok {
				trans.Rollback()
				return result
			}

			// paid in advance
			if p.PaidInAdvance {
				py := Payment{}
				py.PaymentTransaction = pt.Id
				py.Amount = inv.TotalAmount
				py.enterprise = enterpriseId
				ok = py.insertPayment()
				if !ok {
					trans.Rollback()
					return result
				}
			}
		}

		// 2. details for the income
		dInc := AccountingMovementDetail{}
		dInc.Movement = m.Id
		dInc.Journal = *settings.PurchaseJournal
		dInc.AccountNumber = 1
		dInc.Debit = inv.TotalWithDiscount
		dInc.Type = "N"
		dInc.DocumentName = inv.InvoiceName
		dInc.PaymentMethod = inv.PaymentMethod
		dInc.enterprise = enterpriseId
		ok = dInc.insertAccountingMovementDetail()
		if !ok {
			trans.Rollback()
			return result
		}

		// 3. details for the VAT

		// get the details and sort
		det := getPurchaseInvoiceDetail(inv.Id, inv.enterprise)
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
				dVat.Movement = m.Id
				dVat.Journal = journal
				dVat.AccountNumber = accountNumber
				dVat.Debit = debit
				dVat.Type = "N"
				dVat.DocumentName = inv.InvoiceName
				dVat.PaymentMethod = inv.PaymentMethod
				dVat.enterprise = enterpriseId
				ok := dVat.insertAccountingMovementDetail()
				if !ok {
					trans.Rollback()
					return result
				}

				details = make([]PurchaseInvoiceDetail, 0)
				details = append(details, d[i])
			}
		}

		// set the accounting movement on the invoice
		sqlStatement := `UPDATE purchase_invoice SET accounting_movement=$2 WHERE id=$1`
		_, err := db.Exec(sqlStatement, invoiceIds[i], m.Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return result
		}
		result[i].Ok = true
	}

	///
	err = trans.Commit()
	if err != nil {
		for i := 0; i < len(result); i++ {
			result[i].Ok = false
			result[i].Result = 0
		}
	}
	///
	return result
}

func getAccountingMovementSaleInvoices(movementId int64) []SalesInvoice {
	var invoices []SalesInvoice = make([]SalesInvoice, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice WHERE accounting_movement=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, movementId)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.SimplifiedInvoice, &i.Amending, &i.AmendedInvoice, &i.CustomerName)
		invoices = append(invoices, i)
	}

	return invoices
}

func getAccountingMovementPurchaseInvoices(movementId int64) []PurchaseInvoice {
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice WHERE accounting_movement=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, movementId)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	for rows.Next() {
		i := PurchaseInvoice{}
		rows.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice, &i.SupplierName)
		invoices = append(invoices, i)
	}

	return invoices
}
