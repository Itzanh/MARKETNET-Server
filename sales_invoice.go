package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type SalesInvoice struct {
	Id                 int64     `json:"id"`
	Customer           int32     `json:"customer"`
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
	SimplifiedInvoice  bool      `json:"simplifiedInvoice"`
	CustomerName       string    `json:"customerName"`
	Amending           bool      `json:"amending"`
	AmendedInvoice     *int64    `json:"amendedInvoice"`
	enterprise         int32
}

type SaleInvoices struct {
	Rows     int32          `json:"rows"`
	Invoices []SalesInvoice `json:"invoices"`
}

func (q *PaginationQuery) getSalesInvoices() SaleInvoices {
	si := SaleInvoices{}
	if !q.isValid() {
		return si
	}

	si.Invoices = make([]SalesInvoice, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice WHERE enterprise=$3 ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, q.enterprise)
	if err != nil {
		log("DB", err.Error())
		return si
	}
	defer rows.Close()

	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.SimplifiedInvoice, &i.Amending, &i.AmendedInvoice, &i.CustomerName)
		si.Invoices = append(si.Invoices, i)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.sales_invoice WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, q.enterprise)
	row.Scan(&si.Rows)

	return si
}

type OrderSearch struct {
	PaginatedSearch
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
	NotPosted bool       `json:"notPosted"`
}

func (s *OrderSearch) searchSalesInvoices() SaleInvoices {
	si := SaleInvoices{}
	if !s.isValid() {
		return si
	}

	si.Invoices = make([]SalesInvoice, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT sales_invoice.*,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice WHERE invoice_number=$1 AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT sales_invoice.*,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice INNER JOIN customer ON customer.id=sales_invoice.customer WHERE customer.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_invoice.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_invoice.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.NotPosted {
			sqlStatement += ` AND accounting_movement IS NULL`
		}
		sqlStatement += ` AND sales_invoice.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		sqlStatement += ` ORDER BY date_created DESC OFFSET $` + strconv.Itoa(len(interfaces)+1) + ` LIMIT $` + strconv.Itoa(len(interfaces)+2)
		interfaces = append(interfaces, s.Offset)
		interfaces = append(interfaces, s.Limit)
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return si
	}
	defer rows.Close()

	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.SimplifiedInvoice, &i.Amending, &i.AmendedInvoice, &i.CustomerName)
		si.Invoices = append(si.Invoices, i)
	}

	var row *sql.Row
	orderNumber, err = strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT COUNT(*) FROM sales_invoice WHERE invoice_number=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*) FROM sales_invoice INNER JOIN customer ON customer.id=sales_invoice.customer WHERE customer.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_invoice.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_invoice.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.NotPosted {
			sqlStatement += ` AND accounting_movement IS NULL`
		}
		sqlStatement += ` AND sales_invoice.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return si
	}
	row.Scan(&si.Rows)

	return si
}

func getSalesInvoiceRow(invoiceId int64) SalesInvoice {
	sqlStatement := `SELECT * FROM sales_invoice WHERE id=$1`
	row := db.QueryRow(sqlStatement, invoiceId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesInvoice{}
	}

	i := SalesInvoice{}
	row.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
		&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
		&i.AccountingMovement, &i.enterprise, &i.SimplifiedInvoice, &i.Amending, &i.AmendedInvoice)

	return i
}

func (i *SalesInvoice) isValid() bool {
	return !(i.Customer <= 0 || i.PaymentMethod <= 0 || len(i.BillingSeries) == 0 || i.Currency <= 0 || i.BillingAddress <= 0)
}

func (i *SalesInvoice) insertSalesInvoice(userId int32) (bool, int64) {
	if !i.isValid() {
		return false, 0
	}

	// get invoice name
	i.InvoiceNumber = getNextSaleInvoiceNumber(i.BillingSeries, i.enterprise)
	if i.InvoiceNumber <= 0 {
		return false, 0
	}
	now := time.Now()
	i.InvoiceName = i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	// get currency exchange
	i.CurrencyChange = getCurrencyExchange(i.Currency)

	// simplified invoice
	address := getAddressRow(i.BillingAddress)
	if address.Id <= 0 {
		return false, 0
	}
	country := getCountryRow(address.Country, i.enterprise)
	if country.Id <= 0 {
		return false, 0
	}
	if country.Zone == "E" { // Export
		i.SimplifiedInvoice = false
	} else {
		customer := getCustomerRow(i.Customer)
		if country.Zone == "N" { // National
			i.SimplifiedInvoice = len(customer.TaxId) == 0
		} else { // European Union
			i.SimplifiedInvoice = len(customer.TaxId) == 0 && len(customer.VatNumber) == 0
		}
	}

	sqlStatement := `INSERT INTO public.sales_invoice(customer, payment_method, billing_series, currency, currency_change, billing_address, discount_percent, fix_discount, shipping_price, shipping_discount, total_with_discount, total_amount, invoice_number, invoice_name, enterprise, simplified_invoice) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`
	row := db.QueryRow(sqlStatement, i.Customer, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, i.DiscountPercent, i.FixDiscount, i.ShippingPrice, i.ShippingDiscount, i.TotalWithDiscount, i.TotalAmount, i.InvoiceNumber, i.InvoiceName, i.enterprise, i.SimplifiedInvoice)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0
	}

	var invoiceId int64
	row.Scan(&invoiceId)

	if invoiceId > 0 {
		insertTransactionalLog(i.enterprise, "sales_invoice", int(invoiceId), userId, "I")
	}

	return invoiceId > 0, invoiceId
}

func (i *SalesInvoice) deleteSalesInvoice(userId int32) bool {
	if i.Id <= 0 {
		return false
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(i.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return false
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoice := getSalesInvoiceRow(i.Id)
		invoiceNumber := getNextSaleInvoiceNumber(invoice.BillingSeries, invoice.enterprise)
		if invoiceNumber <= 0 || invoice.InvoiceNumber != (invoiceNumber-1) {
			return false
		}
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	d := getSalesInvoiceDetail(i.Id, i.enterprise)

	for i := 0; i < len(d); i++ {
		ok := d[i].deleteSalesInvoiceDetail(userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	insertTransactionalLog(i.enterprise, "sales_invoice", int(i.Id), userId, "D")

	sqlStatement := `DELETE FROM public.sales_invoice WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, i.Id, i.enterprise)
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

func toggleSimplifiedInvoiceSalesInvoice(invoiceId int64, enterpriseId int32, userId int32) bool {
	sqlStatement := `UPDATE sales_invoice SET simplified_invoice = NOT simplified_invoice WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, invoiceId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_invoice", int(invoiceId), userId, "U")

	rows, _ := res.RowsAffected()
	return rows > 0
}

type MakeAmendingInvoice struct {
	InvoiceId   int64   `json:"invoiceId"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
}

func makeAmendingSaleInvoice(invoiceId int64, enterpriseId int32, quantity float64, description string) bool {
	i := getSalesInvoiceRow(invoiceId)
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
	invoiceNumber := getNextSaleInvoiceNumber(i.BillingSeries, i.enterprise)
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

	sqlStatement := `INSERT INTO public.sales_invoice(customer, payment_method, billing_series, currency, currency_change, billing_address, enterprise, simplified_invoice, amending, amended_invoice, invoice_number, invoice_name, total_products, total_with_discount, vat_amount, total_amount, discount_percent, fix_discount, shipping_price, shipping_discount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 0, 0, 0, 0) RETURNING id`
	row := db.QueryRow(sqlStatement, i.Customer, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, enterpriseId, i.SimplifiedInvoice, true, i.Id, invoiceNumber, invoiceName, -detailAmount, -detailAmount, -vatAmount, -quantity)

	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var amendingInvoiceId int64
	row.Scan(&amendingInvoiceId)

	sqlStatement = `INSERT INTO public.sales_invoice_detail(invoice, description, price, quantity, vat_percent, total_amount, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, amendingInvoiceId, description, -detailAmount, 1, vatPercent, -quantity, enterpriseId)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesInvoice(invoiceId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32) bool {
	sqlStatement := `UPDATE sales_invoice SET total_products = total_products + $2, vat_amount = vat_amount + $3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsSaleInvoice(enterpriseId, invoiceId, userId)
}

// Applies the logic to calculate the totals of the sales invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleInvoice(enterpriseId int32, invoiceId int64, userId int32) bool {
	sqlStatement := `UPDATE sales_invoice SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE sales_invoice SET total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err = db.Exec(sqlStatement, invoiceId)

	if err != nil {
		log("DB", err.Error())
	}

	insertTransactionalLog(enterpriseId, "sales_invoice", int(invoiceId), userId, "U")

	return err == nil
}

func invoiceAllSaleOrder(saleOrderId int64, enterpriseId int32, userId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.enterprise != enterpriseId {
		return false
	}
	if saleOrder.InvoicedLines >= saleOrder.LinesNumber {
		return false
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.enterprise)

	if saleOrder.Id <= 0 || len(orderDetails) == 0 {
		return false
	}

	// create an invoice for that order
	invoice := SalesInvoice{}
	invoice.Customer = saleOrder.Customer
	invoice.BillingAddress = saleOrder.BillingAddress
	invoice.BillingSeries = saleOrder.BillingSeries
	invoice.Currency = saleOrder.Currency
	invoice.PaymentMethod = saleOrder.PaymentMethod

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	ok := setDatePaymentAcceptedSalesOrder(enterpriseId, saleOrder.Id, userId)
	if !ok {
		trans.Rollback()
		return false
	}

	invoice.enterprise = saleOrder.enterprise
	ok, invoiceId := invoice.insertSalesInvoice(userId)
	if !ok {
		trans.Rollback()
		return false
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		invoiceDetal := SalesInvoiceDetail{}
		invoiceDetal.Invoice = invoiceId
		invoiceDetal.OrderDetail = &orderDetail.Id
		invoiceDetal.Price = orderDetail.Price
		invoiceDetal.Product = &orderDetail.Product
		invoiceDetal.Quantity = orderDetail.Quantity
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		invoiceDetal.enterprise = invoice.enterprise
		ok = invoiceDetal.insertSalesInvoiceDetail(false, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	go ecommerceControllerupdateStatusPaymentAccepted(saleOrderId, invoice.enterprise)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type OrderDetailGenerate struct {
	OrderId   int64                          `json:"orderId"`
	Selection []OrderDetailGenerateSelection `json:"selection"`
}

type OrderDetailGenerateSelection struct {
	Id       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (invoiceInfo *OrderDetailGenerate) invoicePartiallySaleOrder(enterpriseId int32, userId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(invoiceInfo.OrderId)
	if saleOrder.Id <= 0 || saleOrder.enterprise != enterpriseId || len(invoiceInfo.Selection) == 0 {
		return false
	}
	if saleOrder.InvoicedLines >= saleOrder.LinesNumber {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != invoiceInfo.OrderId || invoiceInfo.Selection[i].Quantity == 0 || invoiceInfo.Selection[i].Quantity > orderDetail.Quantity || orderDetail.QuantityInvoiced >= orderDetail.Quantity {
			return false
		}
		saleOrderDetails = append(saleOrderDetails, orderDetail)
	}

	// create an invoice for that order
	invoice := SalesInvoice{}
	invoice.Customer = saleOrder.Customer
	invoice.BillingAddress = saleOrder.BillingAddress
	invoice.BillingSeries = saleOrder.BillingSeries
	invoice.Currency = saleOrder.Currency
	invoice.PaymentMethod = saleOrder.PaymentMethod

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	ok := setDatePaymentAcceptedSalesOrder(enterpriseId, saleOrder.Id, userId)
	if !ok {
		trans.Rollback()
		return false
	}

	invoice.enterprise = saleOrder.enterprise
	ok, invoiceId := invoice.insertSalesInvoice(userId)
	if !ok {
		trans.Rollback()
		return false
	}
	for i := 0; i < len(saleOrderDetails); i++ {
		orderDetail := saleOrderDetails[i]
		invoiceDetal := SalesInvoiceDetail{}
		invoiceDetal.Invoice = invoiceId
		invoiceDetal.OrderDetail = &orderDetail.Id
		invoiceDetal.Price = orderDetail.Price
		invoiceDetal.Product = &orderDetail.Product
		invoiceDetal.Quantity = invoiceInfo.Selection[i].Quantity
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		invoiceDetal.enterprise = invoice.enterprise
		ok = invoiceDetal.insertSalesInvoiceDetail(false, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	go ecommerceControllerupdateStatusPaymentAccepted(invoiceInfo.OrderId, invoice.enterprise)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type SalesInvoiceRelations struct {
	Orders        []SaleOrder         `json:"orders"`
	DeliveryNotes []SalesDeliveryNote `json:"notes"`
	Invoices      []SalesInvoice      `json:"invoices"`
}

func getSalesInvoiceRelations(invoiceId int64, enterpriseId int32) SalesInvoiceRelations {
	return SalesInvoiceRelations{
		Orders:        getSalesInvoiceOrders(invoiceId, enterpriseId),
		DeliveryNotes: getSalesInvoiceDeliveryNotes(invoiceId, enterpriseId),
		Invoices:      getSalesInvoiceAmendingAmendedInvoices(invoiceId, enterpriseId),
	}
}

func getSalesInvoiceOrders(invoiceId int64, enterpriseId int32) []SaleOrder {
	var orders []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT DISTINCT sales_order.* FROM sales_invoice INNER JOIN sales_invoice_detail ON sales_invoice.id = sales_invoice_detail.invoice INNER JOIN sales_order_detail ON sales_invoice_detail.order_detail = sales_order_detail.id INNER JOIN sales_order ON sales_order_detail.order = sales_order.id WHERE sales_invoice.id = $1 AND sales_invoice.enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, invoiceId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	defer rows.Close()

	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.enterprise, &s.shopifyDraftId)
		orders = append(orders, s)
	}

	return orders
}

func getSalesInvoiceDeliveryNotes(invoiceId int64, enterpriseId int32) []SalesDeliveryNote {
	var notes []SalesDeliveryNote = make([]SalesDeliveryNote, 0)
	sqlStatement := `SELECT DISTINCT sales_delivery_note.* FROM sales_invoice INNER JOIN sales_invoice_detail ON sales_invoice.id = sales_invoice_detail.invoice INNER JOIN sales_order_detail ON sales_invoice_detail.order_detail = sales_order_detail.id INNER JOIN warehouse_movement ON warehouse_movement.sales_order_detail=sales_order_detail.id INNER JOIN sales_delivery_note ON sales_delivery_note.id=warehouse_movement.sales_delivery_note WHERE sales_invoice.id=$1 AND sales_delivery_note.enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, invoiceId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		n := SalesDeliveryNote{}
		rows.Scan(&n.Id, &n.Warehouse, &n.Customer, &n.DateCreated, &n.PaymentMethod, &n.BillingSeries, &n.ShippingAddress, &n.TotalProducts, &n.DiscountPercent, &n.FixDiscount, &n.ShippingPrice, &n.ShippingDiscount, &n.TotalWithDiscount, &n.VatAmount, &n.TotalAmount, &n.LinesNumber, &n.DeliveryNoteName, &n.DeliveryNoteNumber, &n.Currency, &n.CurrencyChange, &n.enterprise)
		notes = append(notes, n)
	}

	return notes
}

func getSalesInvoiceAmendingAmendedInvoices(invoiceId int64, enterpriseId int32) []SalesInvoice {
	invoices := make([]SalesInvoice, 0)

	i := getSalesInvoiceRow(invoiceId)
	if i.enterprise != enterpriseId {
		return invoices
	}

	if i.Amending && i.AmendedInvoice != nil {
		invoices = append(invoices, getSalesInvoiceRow(*i.AmendedInvoice))
	}

	sqlStatement := `SELECT * FROM sales_invoice WHERE amended_invoice=$1`
	rows, err := db.Query(sqlStatement, i.Id)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	defer rows.Close()

	for rows.Next() {
		inv := SalesInvoice{}
		rows.Scan(&inv.Id, &inv.Customer, &inv.DateCreated, &inv.PaymentMethod, &inv.BillingSeries, &inv.Currency, &inv.CurrencyChange, &inv.BillingAddress, &inv.TotalProducts,
			&inv.DiscountPercent, &inv.FixDiscount, &inv.ShippingPrice, &inv.ShippingDiscount, &inv.TotalWithDiscount, &inv.VatAmount, &inv.TotalAmount, &inv.LinesNumber, &inv.InvoiceNumber, &inv.InvoiceName,
			&inv.AccountingMovement, &inv.enterprise, &inv.SimplifiedInvoice, &inv.Amending, &inv.AmendedInvoice)
		invoices = append(invoices, inv)
	}

	return invoices
}
