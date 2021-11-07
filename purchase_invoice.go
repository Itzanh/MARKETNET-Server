package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type PurchaseInvoice struct {
	Id                 int64     `json:"id"`
	Supplier           int32     `json:"supplier"`
	DateCreated        time.Time `json:"dateCreated"`
	PaymentMethod      int32     `json:"paymentMethod"`
	BillingSeries      string    `json:"billingSeries"`
	Currency           int32     `json:"currency"`
	CurrencyChange     float64   `json:"currencyChange"`
	BillingAddress     int32     `json:"billingAddress"`
	TotalProducts      float64   `json:"totalProducts"`
	DiscountPercent    float64   `json:"discountPercent"`
	FixDiscount        float64   `json:"fixDiscount"`
	ShippingPrice      float64   `json:"shippingPrice"`
	ShippingDiscount   float64   `json:"shippingDiscount"`
	TotalWithDiscount  float64   `json:"totalWithDiscount"`
	VatAmount          float64   `json:"vatAmount"`
	TotalAmount        float64   `json:"totalAmount"`
	LinesNumber        int16     `json:"linesNumber"`
	InvoiceNumber      int32     `json:"invoiceNumber"`
	InvoiceName        string    `json:"invoiceName"`
	AccountingMovement *int64    `json:"accountingMovement"`
	SupplierName       string    `json:"supplierName"`
	Amending           bool      `json:"amending"`
	AmendedInvoice     *int64    `json:"amendedInvoice"`
	enterprise         int32
}

func getPurchaseInvoices(enterpriseId int32) []PurchaseInvoice {
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
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

func (s *OrderSearch) searchPurchaseInvoice() []PurchaseInvoice {
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_invoice.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_invoice.supplier) FROM purchase_invoice WHERE invoice_number=$1 AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.Enterprise)
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
		interfaces = append(interfaces, s.Enterprise)
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
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
		&i.AccountingMovement, &i.enterprise, &i.Amending, &i.AmendedInvoice)

	return i
}

func (i *PurchaseInvoice) isValid() bool {
	return !(i.Supplier <= 0 || i.PaymentMethod <= 0 || len(i.BillingSeries) == 0 || i.Currency <= 0 || i.BillingAddress <= 0)
}

func (i *PurchaseInvoice) insertPurchaseInvoice() (bool, int64) {
	if !i.isValid() {
		return false, 0
	}

	i.InvoiceNumber = getNextPurchaseInvoiceNumber(i.BillingSeries, i.enterprise)
	if i.InvoiceNumber <= 0 {
		return false, 0
	}
	i.CurrencyChange = getCurrencyExchange(i.Currency)
	now := time.Now()
	i.InvoiceName = i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	sqlStatement := `INSERT INTO public.purchase_invoice(supplier, payment_method, billing_series, currency, currency_change, billing_address, discount_percent, fix_discount, shipping_price, shipping_discount, total_with_discount, total_amount, invoice_number, invoice_name, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	row := db.QueryRow(sqlStatement, i.Supplier, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, i.DiscountPercent, i.FixDiscount, i.ShippingPrice, i.ShippingDiscount, i.TotalWithDiscount, i.TotalAmount, i.InvoiceNumber, i.InvoiceName, i.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0
	}

	var invoiceId int64
	row.Scan(&invoiceId)
	return invoiceId > 0, invoiceId
}

func (i *PurchaseInvoice) deletePurchaseInvoice() bool {
	if i.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(i.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return false
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoice := getPurchaseInvoiceRow(i.Id)
		invoiceNumber := getNextPurchaseInvoiceNumber(invoice.BillingSeries, invoice.enterprise)
		if invoiceNumber <= 0 || invoice.InvoiceNumber != (invoiceNumber-1) {
			return false
		}
	}

	inMemoryInvoice := getPurchaseInvoiceRow(i.Id)
	if inMemoryInvoice.enterprise != i.enterprise {
		return false
	}

	d := getPurchaseInvoiceDetail(i.Id, i.enterprise)
	for i := 0; i < len(d); i++ {
		ok := d[i].deletePurchaseInvoiceDetail()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.purchase_invoice WHERE id=$1`
	res, err := db.Exec(sqlStatement, i.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	if err != nil {
		return false
	}
	///

	rows, _ := res.RowsAffected()
	return rows > 0
}

func makeAmendingPurchaseInvoice(invoiceId int64, enterpriseId int32, quantity float64, description string) bool {
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

	sqlStatement = `INSERT INTO public.purchase_invoice_details(invoice, description, price, quantity, vat_percent, total_amount, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, amendingInvoiceId, description, -detailAmount, 1, vatPercent, -quantity, enterpriseId)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseInvoice(invoiceId int64, totalAmount float64, vatPercent float64) bool {
	sqlStatement := `UPDATE purchase_invoice SET total_products=total_products+$2,vat_amount=vat_amount+$3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId)
}

// Applies the logic to calculate the totals of the purchase invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseInvoice(invoiceId int64) bool {
	sqlStatement := `UPDATE purchase_invoice SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE purchase_invoice SET total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err = db.Exec(sqlStatement, invoiceId)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func invoiceAllPurchaseOrder(purchaseOrderId int64, enterpriseId int32) bool {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	if purchaseOrder.enterprise != enterpriseId {
		return false
	}
	orderDetails := getPurchaseOrderDetail(purchaseOrderId, purchaseOrder.enterprise)

	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return false
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
		return false
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id)
	if !ok {
		trans.Rollback()
		return false
	}

	invoice.enterprise = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice()
	if !ok {
		trans.Rollback()
		return false
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
		ok = invoiceDetail.insertPurchaseInvoiceDetail(false)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func (invoiceInfo *OrderDetailGenerate) invoicePartiallyPurchaseOrder(enterpriseId int32) bool {
	// get the sale order and it's details
	purchaseOrder := getPurchaseOrderRow(invoiceInfo.OrderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.enterprise != enterpriseId || len(invoiceInfo.Selection) == 0 {
		return false
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != invoiceInfo.OrderId || invoiceInfo.Selection[i].Quantity == 0 || invoiceInfo.Selection[i].Quantity > orderDetail.Quantity {
			return false
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
		return false
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id)
	if !ok {
		trans.Rollback()
		return false
	}

	invoice.enterprise = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice()
	if !ok {
		trans.Rollback()
		return false
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
		ok = invoiceDetail.insertPurchaseInvoiceDetail(false)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
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

	for rows.Next() {
		inv := PurchaseInvoice{}
		rows.Scan(&inv.Id, &inv.Supplier, &inv.DateCreated, &inv.PaymentMethod, &inv.BillingSeries, &inv.Currency, &inv.CurrencyChange, &inv.BillingAddress, &inv.TotalProducts,
			&inv.DiscountPercent, &inv.FixDiscount, &inv.ShippingPrice, &inv.ShippingDiscount, &inv.TotalWithDiscount, &inv.VatAmount, &inv.TotalAmount, &inv.LinesNumber, &inv.InvoiceNumber, &inv.InvoiceName,
			&inv.AccountingMovement, &inv.enterprise, &inv.Amending, &inv.AmendedInvoice)
		invoices = append(invoices, inv)
	}

	return invoices
}
