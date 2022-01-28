package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PurchaseInvoice struct {
	Id                  int64     `json:"id"`
	Supplier            int32     `json:"supplier"`
	DateCreated         time.Time `json:"dateCreated"`
	PaymentMethod       int32     `json:"paymentMethod"`
	BillingSeries       string    `json:"billingSeries"`
	Currency            int32     `json:"currency"`
	CurrencyChange      float64   `json:"currencyChange"`
	BillingAddress      int32     `json:"billingAddress"`
	TotalProducts       float64   `json:"totalProducts"`
	DiscountPercent     float64   `json:"discountPercent"`
	FixDiscount         float64   `json:"fixDiscount"`
	ShippingPrice       float64   `json:"shippingPrice"`
	ShippingDiscount    float64   `json:"shippingDiscount"`
	TotalWithDiscount   float64   `json:"totalWithDiscount"`
	VatAmount           float64   `json:"vatAmount"`
	TotalAmount         float64   `json:"totalAmount"`
	LinesNumber         int16     `json:"linesNumber"`
	InvoiceNumber       int32     `json:"invoiceNumber"`
	InvoiceName         string    `json:"invoiceName"`
	AccountingMovement  *int64    `json:"accountingMovement"`
	SupplierName        string    `json:"supplierName"`
	Amending            bool      `json:"amending"`
	AmendedInvoice      *int64    `json:"amendedInvoice"`
	IncomeTax           bool      `json:"incomeTax"`
	IncomeTaxBase       float64   `json:"incomeTaxBase"`
	IncomeTaxPercentage float64   `json:"incomeTaxPercentage"`
	IncomeTaxValue      float64   `json:"incomeTaxValue"`
	Rent                bool      `json:"rent"`
	RentBase            float64   `json:"rentBase"`
	RentPercentage      float64   `json:"rentPercentage"`
	RentValue           float64   `json:"rentValue"`
	enterprise          int32
}

type PurchaseInvoices struct {
	Rows     int32                 `json:"rows"`
	Invoices []PurchaseInvoice     `json:"invoices"`
	Footer   PurchaseInvoiceFooter `json:"footer"`
}

type PurchaseInvoiceFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseInvoices(enterpriseId int32) PurchaseInvoices {
	in := PurchaseInvoices{}
	in.Invoices = make([]PurchaseInvoice, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return in
	}
	defer rows.Close()

	for rows.Next() {
		i := PurchaseInvoice{}
		rows.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice, &i.IncomeTax, &i.IncomeTaxBase, &i.IncomeTaxPercentage, &i.IncomeTaxValue, &i.Rent, &i.RentBase,
			&i.RentPercentage, &i.RentValue, &i.SupplierName)
		in.Invoices = append(in.Invoices, i)
	}

	sqlStatement = `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_invoice WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return in
	}
	row.Scan(&in.Rows, &in.Footer.TotalProducts, &in.Footer.TotalAmount)

	return in
}

func (s *OrderSearch) searchPurchaseInvoice() PurchaseInvoices {
	in := PurchaseInvoices{}
	in.Invoices = make([]PurchaseInvoice, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_invoice.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice WHERE invoice_number=$1 AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT purchase_invoice.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice INNER JOIN suppliers ON suppliers.id=purchase_invoice.supplier WHERE suppliers.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_invoice.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_invoice.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.NotPosted {
			sqlStatement += ` AND accounting_movement IS NULL`
		}
		sqlStatement += ` AND purchase_invoice.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return in
	}
	defer rows.Close()

	for rows.Next() {
		i := PurchaseInvoice{}
		rows.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice, &i.IncomeTax, &i.IncomeTaxBase, &i.IncomeTaxPercentage, &i.IncomeTaxValue, &i.Rent, &i.RentBase,
			&i.RentPercentage, &i.RentValue, &i.SupplierName)
		in.Invoices = append(in.Invoices, i)
	}

	var row *sql.Row
	if err == nil {
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_invoice WHERE invoice_number=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_invoice INNER JOIN suppliers ON suppliers.id=purchase_invoice.supplier WHERE suppliers.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_invoice.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_invoice.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.NotPosted {
			sqlStatement += ` AND accounting_movement IS NULL`
		}
		sqlStatement += ` AND purchase_invoice.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return in
	}
	row.Scan(&in.Rows, &in.Footer.TotalProducts, &in.Footer.TotalAmount)

	return in
}

func getPurchaseInvoiceRow(invoiceId int64) PurchaseInvoice {
	sqlStatement := `SELECT * FROM purchase_invoice WHERE id=$1`
	row := db.QueryRow(sqlStatement, invoiceId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseInvoice{}
	}

	i := PurchaseInvoice{}
	row.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
		&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
		&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice, &i.IncomeTax, &i.IncomeTaxBase, &i.IncomeTaxPercentage, &i.IncomeTaxValue, &i.Rent, &i.RentBase,
		&i.RentPercentage, &i.RentValue)

	return i
}

func getPurchaseInvoiceRowTransaction(invoiceId int64, trans sql.Tx) PurchaseInvoice {
	sqlStatement := `SELECT * FROM purchase_invoice WHERE id=$1`
	row := trans.QueryRow(sqlStatement, invoiceId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseInvoice{}
	}

	i := PurchaseInvoice{}
	row.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
		&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
		&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice, &i.IncomeTax, &i.IncomeTaxBase, &i.IncomeTaxPercentage, &i.IncomeTaxValue, &i.Rent, &i.RentBase,
		&i.RentPercentage, &i.RentValue)

	return i
}

func (i *PurchaseInvoice) isValid() bool {
	return !(i.Supplier <= 0 || i.PaymentMethod <= 0 || len(i.BillingSeries) == 0 || i.Currency <= 0 || i.BillingAddress <= 0 || i.IncomeTaxBase < 0 || i.IncomeTaxPercentage < 0 || i.RentBase < 0 || i.RentPercentage < 0)
}

func (i *PurchaseInvoice) insertPurchaseInvoice(userId int32, trans *sql.Tx) (bool, int64) {
	if !i.isValid() {
		return false, 0
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false, 0
		}
		///
	}

	i.InvoiceNumber = getNextPurchaseInvoiceNumber(i.BillingSeries, i.enterprise)
	if i.InvoiceNumber <= 0 {
		return false, 0
	}
	i.CurrencyChange = getCurrencyExchange(i.Currency)
	now := time.Now()
	i.InvoiceName = i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	sqlStatement := `INSERT INTO public.purchase_invoice(supplier, payment_method, billing_series, currency, currency_change, billing_address, discount_percent, fix_discount, shipping_price, shipping_discount, total_with_discount, total_amount, invoice_number, invoice_name, enterprise, income_tax, income_tax_percentage, rent, rent_percentage) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`
	row := trans.QueryRow(sqlStatement, i.Supplier, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, i.DiscountPercent, i.FixDiscount, i.ShippingPrice, i.ShippingDiscount, i.TotalWithDiscount, i.TotalAmount, i.InvoiceNumber, i.InvoiceName, i.enterprise, i.IncomeTax, i.IncomeTaxPercentage, i.Rent, i.RentPercentage)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false, 0
	}

	var invoiceId int64
	row.Scan(&invoiceId)

	if invoiceId > 0 {
		insertTransactionalLog(i.enterprise, "purchase_invoice", int(invoiceId), userId, "I")
		json, _ := json.Marshal(i)
		go fireWebHook(i.enterprise, "purchase_invoice", "POST", string(json))
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false, 0
		}
		///
	}

	return invoiceId > 0, invoiceId
}

// 1. can't delete details in posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (i *PurchaseInvoice) deletePurchaseInvoice(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
	if i.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice := getPurchaseInvoiceRow(i.Id)
	if invoice.enterprise != i.enterprise {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if invoice.AccountingMovement != nil {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(i.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoice := getPurchaseInvoiceRow(i.Id)
		invoiceNumber := getNextPurchaseInvoiceNumber(invoice.BillingSeries, invoice.enterprise)
		if invoiceNumber <= 0 || invoice.InvoiceNumber != (invoiceNumber-1) {
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
		}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	d := getPurchaseInvoiceDetail(i.Id, i.enterprise)
	for i := 0; i < len(d); i++ {
		ok := d[i].deletePurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(invoice.enterprise, "purchase_invoice", int(i.Id), userId, "D")
	json, _ := json.Marshal(i)
	go fireWebHook(i.enterprise, "purchase_invoice", "DELETE", string(json))

	sqlStatement := `DELETE FROM public.purchase_invoice WHERE id=$1`
	res, err := trans.Exec(sqlStatement, i.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	rows, _ := res.RowsAffected()
	return OkAndErrorCodeReturn{Ok: rows > 0}
}

func makeAmendingPurchaseInvoice(invoiceId int64, enterpriseId int32, quantity float64, description string, userId int32) bool {
	i := getPurchaseInvoiceRow(invoiceId)
	if i.Id <= 0 || i.enterprise != enterpriseId {
		return false
	}

	// we can't make an amending invoice the same day as the original invoice
	now := time.Now()
	if i.DateCreated.Year() == now.Year() && i.DateCreated.YearDay() == now.YearDay() {
		return false
	}
	// we can't make an amending invoice with a greater amount that the original invoice
	if quantity <= 0 || quantity > i.TotalAmount {
		return false
	}

	settings := getSettingsRecordById(enterpriseId)

	// get invoice name
	invoiceNumber := getNextPurchaseInvoiceNumber(i.BillingSeries, i.enterprise)
	if i.InvoiceNumber <= 0 {
		return false
	}
	invoiceName := i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	var detailAmount float64
	var vatPercent float64
	// VAT excluded invoice, when amending the invoice, we are not returning any tax money
	if i.VatAmount == 0 {
		detailAmount = quantity
		vatPercent = 0
	} else { // Invoice with VAT, we return the quantity without the tax to the customer, and then we add the tax percent, so the total of the invoice is the amount we want to return (in taxes and to the customer)
		detailAmount = quantity / (1 + (settings.DefaultVatPercent / 100))
		vatPercent = settings.DefaultVatPercent
	}
	var vatAmount float64 = (quantity / 100) * vatPercent

	sqlStatement := `INSERT INTO public.purchase_invoice(supplier, payment_method, billing_series, currency, currency_change, billing_address, enterprise, amending, amended_invoice, invoice_number, invoice_name, total_products, total_with_discount, vat_amount, total_amount, discount_percent, fix_discount, shipping_price, shipping_discount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 0, 0, 0, 0) RETURNING id`
	row := db.QueryRow(sqlStatement, i.Supplier, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, enterpriseId, true, i.Id, invoiceNumber, invoiceName, -detailAmount, -detailAmount, -vatAmount, -quantity)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var amendingInvoiceId int64
	row.Scan(&amendingInvoiceId)
	i.Id = amendingInvoiceId

	if amendingInvoiceId > 0 {
		insertTransactionalLog(enterpriseId, "purchase_invoice", int(amendingInvoiceId), userId, "I")
		json, _ := json.Marshal(i)
		go fireWebHook(i.enterprise, "purchase_invoice", "POST", string(json))
	}

	sqlStatement = `INSERT INTO public.purchase_invoice_details(invoice, description, price, quantity, vat_percent, total_amount, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	row = db.QueryRow(sqlStatement, amendingInvoiceId, description, -detailAmount, 1, vatPercent, -quantity, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var amendingInvoiceDetailId int64
	row.Scan(&amendingInvoiceDetailId)

	if amendingInvoiceDetailId > 0 {
		insertTransactionalLog(enterpriseId, "purchase_invoice_details", int(amendingInvoiceDetailId), userId, "I")
		d := getPurchaseInvoiceRow(amendingInvoiceDetailId)
		json, _ := json.Marshal(d)
		go fireWebHook(d.enterprise, "purchase_invoice_details", "POST", string(json))
	}

	return amendingInvoiceId > 0
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseInvoice(invoiceId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_invoice SET total_products=total_products+$2,vat_amount=vat_amount+$3 WHERE id = $1`
	_, err := trans.Exec(sqlStatement, invoiceId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// Applies the logic to calculate the totals of the purchase invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseInvoice(invoiceId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_invoice SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := trans.Exec(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE public.purchase_invoice SET income_tax_value = CASE WHEN income_tax THEN (income_tax_base/100)*income_tax_percentage ELSE 0 END, rent_value = CASE WHEN rent THEN (rent_base/100)*rent_percentage ELSE 0 END WHERE id = $1`
	_, err = trans.Exec(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE purchase_invoice SET total_amount = total_with_discount + vat_amount - income_tax_value - rent_value WHERE id = $1`
	_, err = trans.Exec(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_invoice", int(invoiceId), userId, "U")
	i := getPurchaseInvoiceRowTransaction(invoiceId, trans)
	json, _ := json.Marshal(i)
	go fireWebHook(i.enterprise, "purchase_invoice", "PUT", string(json))

	return err == nil
}

// Adds a income tax base amount to the income tax base. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addIncomeTaxBasePurchaseInvoice(invoiceId int64, totalAmount float64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_invoice SET income_tax_base=income_tax_base+$2 WHERE id = $1`
	_, err := trans.Exec(sqlStatement, invoiceId, totalAmount)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addRentBaseProductsPurchaseInvoice(invoiceId int64, totalAmount float64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_invoice SET rent_base=rent_base+$2 WHERE id = $1`
	_, err := trans.Exec(sqlStatement, invoiceId, totalAmount)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// ERROR CODES:
// 1. The order is already invoiced
// 2. There are no details to invoice
func invoiceAllPurchaseOrder(purchaseOrderId int64, enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	if purchaseOrder.enterprise != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.InvoicedLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}
	orderDetails := getPurchaseOrderDetail(purchaseOrderId, purchaseOrder.enterprise)
	filterPurchaseOrderDetails(orderDetails, func(pod PurchaseOrderDetail) bool { return pod.QuantityInvoiced < pod.Quantity })
	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	}

	// create an invoice for that order
	invoice := PurchaseInvoice{}
	invoice.Supplier = purchaseOrder.Supplier
	invoice.BillingAddress = purchaseOrder.BillingAddress
	invoice.BillingSeries = purchaseOrder.BillingSeries
	invoice.Currency = purchaseOrder.Currency
	invoice.PaymentMethod = purchaseOrder.PaymentMethod

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id, enterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.enterprise = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		invoiceDetail := PurchaseInvoiceDetail{}
		invoiceDetail.Invoice = invoiceId
		invoiceDetail.OrderDetail = &orderDetail.Id
		invoiceDetail.Price = orderDetail.Price
		invoiceDetail.Product = &orderDetail.Product
		invoiceDetail.Quantity = orderDetail.Quantity
		invoiceDetail.TotalAmount = orderDetail.TotalAmount
		invoiceDetail.VatPercent = orderDetail.VatPercent
		invoiceDetail.enterprise = enterpriseId
		ok = invoiceDetail.insertPurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	transErr = trans.Commit()
	return OkAndErrorCodeReturn{Ok: transErr == nil}
	///
}

// ERROR CODES:
// 1. The order is aleady invoiced
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail is already invoiced
// 4. The selected quantity is greater than the quantity pending of invoicing in the detail
func (invoiceInfo *OrderDetailGenerate) invoicePartiallyPurchaseOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the sale order and it's details
	purchaseOrder := getPurchaseOrderRow(invoiceInfo.OrderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.enterprise != enterpriseId || len(invoiceInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.InvoicedLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != invoiceInfo.OrderId || invoiceInfo.Selection[i].Quantity == 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if invoiceInfo.Selection[i].Quantity > orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 2, ExtraData: []string{product.Name}}
		}
		if orderDetail.QuantityInvoiced >= orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3, ExtraData: []string{product.Name}}
		}
		if (invoiceInfo.Selection[i].Quantity + orderDetail.QuantityInvoiced) > orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 4, ExtraData: []string{product.Name}}
		}
		purchaseOrderDetails = append(purchaseOrderDetails, orderDetail)
	}

	// create an invoice for that order
	invoice := PurchaseInvoice{}
	invoice.Supplier = purchaseOrder.Supplier
	invoice.BillingAddress = purchaseOrder.BillingAddress
	invoice.BillingSeries = purchaseOrder.BillingSeries
	invoice.Currency = purchaseOrder.Currency
	invoice.PaymentMethod = purchaseOrder.PaymentMethod

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id, enterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.enterprise = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(purchaseOrderDetails); i++ {
		orderDetail := purchaseOrderDetails[i]
		invoiceDetail := PurchaseInvoiceDetail{}
		invoiceDetail.Invoice = invoiceId
		invoiceDetail.OrderDetail = &orderDetail.Id
		invoiceDetail.Price = orderDetail.Price
		invoiceDetail.Product = &orderDetail.Product
		invoiceDetail.Quantity = invoiceInfo.Selection[i].Quantity
		invoiceDetail.TotalAmount = orderDetail.TotalAmount
		invoiceDetail.VatPercent = orderDetail.VatPercent
		invoiceDetail.enterprise = enterpriseId
		ok = invoiceDetail.insertPurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	transErr = trans.Commit()
	return OkAndErrorCodeReturn{Ok: transErr == nil}
	///
}

type PurchaseInvoiceRelations struct {
	Orders   []PurchaseOrder   `json:"orders"`
	Invoices []PurchaseInvoice `json:"invoices"`
}

func getPurchaseInvoiceRelations(invoiceId int64, enterpriseId int32) PurchaseInvoiceRelations {
	return PurchaseInvoiceRelations{
		Orders:   getPurchaseInvoiceOrders(invoiceId, enterpriseId),
		Invoices: getPurchaseInvoiceAmendingAmendedInvoices(invoiceId, enterpriseId),
	}
}

func getPurchaseInvoiceOrders(orderId int64, enterpriseId int32) []PurchaseOrder {
	var orders []PurchaseOrder = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT DISTINCT purchase_order.* FROM purchase_invoice INNER JOIN purchase_invoice_details ON purchase_invoice.id=purchase_invoice_details.invoice INNER JOIN purchase_order_detail ON purchase_invoice_details.order_detail=purchase_order_detail.id INNER JOIN purchase_order ON purchase_order_detail."order"=purchase_order.id WHERE purchase_invoice.id=$1 AND purchase_invoice.enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	defer rows.Close()

	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.enterprise)
		orders = append(orders, s)
	}

	return orders
}

func getPurchaseInvoiceAmendingAmendedInvoices(invoiceId int64, enterpriseId int32) []PurchaseInvoice {
	invoices := make([]PurchaseInvoice, 0)

	i := getPurchaseInvoiceRow(invoiceId)
	if i.enterprise != enterpriseId {
		return invoices
	}

	if i.Amending && i.AmendedInvoice != nil {
		invoices = append(invoices, getPurchaseInvoiceRow(*i.AmendedInvoice))
	}

	sqlStatement := `SELECT * FROM purchase_invoice WHERE amended_invoice=$1`
	rows, err := db.Query(sqlStatement, i.Id)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	defer rows.Close()

	for rows.Next() {
		inv := PurchaseInvoice{}
		rows.Scan(&inv.Id, &inv.Supplier, &inv.DateCreated, &inv.PaymentMethod, &inv.BillingSeries, &inv.Currency, &inv.CurrencyChange, &inv.BillingAddress, &inv.TotalProducts,
			&inv.DiscountPercent, &inv.FixDiscount, &inv.ShippingPrice, &inv.ShippingDiscount, &inv.TotalWithDiscount, &inv.VatAmount, &inv.TotalAmount, &inv.LinesNumber, &inv.InvoiceNumber, &inv.InvoiceName,
			&inv.AccountingMovement, &inv.enterprise, &inv.Amending, &inv.AmendedInvoice)
		invoices = append(invoices, inv)
	}

	return invoices
}
