package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type SalesInvoice struct {
	Id                 int32     `json:"id"`
	Customer           int32     `json:"customer"`
	DateCreated        time.Time `json:"dateCreated"`
	PaymentMethod      int16     `json:"paymentMethod"`
	BillingSeries      string    `json:"billingSeries"`
	Currency           int16     `json:"currency"`
	CurrencyChange     float32   `json:"currencyChange"`
	BillingAddress     int32     `json:"billingAddress"`
	TotalProducts      float32   `json:"totalProducts"`
	DiscountPercent    float32   `json:"discountPercent"`
	FixDiscount        float32   `json:"fixDiscount"`
	ShippingPrice      float32   `json:"shippingPrice"`
	ShippingDiscount   float32   `json:"shippingDiscount"`
	TotalWithDiscount  float32   `json:"totalWithDiscount"`
	VatAmount          float32   `json:"vatAmount"`
	TotalAmount        float32   `json:"totalAmount"`
	LinesNumber        int16     `json:"linesNumber"`
	InvoiceNumber      int32     `json:"invoiceNumber"`
	InvoiceName        string    `json:"invoiceName"`
	AccountingMovement *int64    `json:"accountingMovement"`
	CustomerName       string    `json:"customerName"`
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
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit)
	if err != nil {
		log("DB", err.Error())
		return si
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.CustomerName)
		si.Invoices = append(si.Invoices, i)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.sales_invoice`
	row := db.QueryRow(sqlStatement)
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
		sqlStatement := `SELECT sales_invoice.*,(SELECT name FROM customer WHERE customer.id=sales_invoice.customer) FROM sales_invoice WHERE invoice_number=$1 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber)
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
		sqlStatement += ` ORDER BY date_created DESC OFFSET $` + strconv.Itoa(len(interfaces)+1) + ` LIMIT $` + strconv.Itoa(len(interfaces)+2)
		interfaces = append(interfaces, s.Offset)
		interfaces = append(interfaces, s.Limit)
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return si
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.CustomerName)
		si.Invoices = append(si.Invoices, i)
	}

	var row *sql.Row
	orderNumber, err = strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT COUNT(*) FROM sales_invoice WHERE invoice_number=$1 ORDER BY date_created DESC`
		row = db.QueryRow(sqlStatement, orderNumber)
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
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return si
	}
	row.Scan(&si.Rows)

	return si
}

func getSalesInvoiceRow(invoiceId int32) SalesInvoice {
	sqlStatement := `SELECT * FROM sales_invoice WHERE id=$1`
	row := db.QueryRow(sqlStatement, invoiceId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesInvoice{}
	}

	i := SalesInvoice{}
	row.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
		&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
		&i.AccountingMovement)

	return i
}

func (i *SalesInvoice) isValid() bool {
	return !(i.Customer <= 0 || i.PaymentMethod <= 0 || len(i.BillingSeries) == 0 || i.Currency <= 0 || i.BillingAddress <= 0)
}

func (i *SalesInvoice) insertSalesInvoice() (bool, int32) {
	if !i.isValid() {
		return false, 0
	}

	i.InvoiceNumber = getNextSaleInvoiceNumber(i.BillingSeries)
	if i.InvoiceNumber <= 0 {
		return false, 0
	}
	i.CurrencyChange = getCurrencyExchange(i.Currency)
	now := time.Now()
	i.InvoiceName = i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	sqlStatement := `INSERT INTO public.sales_invoice(customer, payment_method, billing_series, currency, currency_change, billing_address, discount_percent, fix_discount, shipping_price, shipping_discount, total_with_discount, total_amount, invoice_number, invoice_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	row := db.QueryRow(sqlStatement, i.Customer, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, i.DiscountPercent, i.FixDiscount, i.ShippingPrice, i.ShippingDiscount, i.TotalWithDiscount, i.TotalAmount, i.InvoiceNumber, i.InvoiceName)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0
	}

	var invoiceId int32
	row.Scan(&invoiceId)
	return invoiceId > 0, invoiceId
}

func (i *SalesInvoice) deleteSalesInvoice() bool {
	if i.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	d := getSalesInvoiceDetail(i.Id)

	for i := 0; i < len(d); i++ {
		ok := d[i].deleteSalesInvoiceDetail()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.sales_invoice WHERE id=$1`
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

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesInvoice(invoiceId int32, totalAmount float32, vatPercent float32) bool {
	sqlStatement := `UPDATE sales_invoice SET total_products = total_products + $2, vat_amount = vat_amount + $3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsSaleInvoice(invoiceId)
}

// Applies the logic to calculate the totals of the sales invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleInvoice(invoiceId int32) bool {
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

	return err == nil
}

func invoiceAllSaleOrder(saleOrderId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	orderDetails := getSalesOrderDetail(saleOrderId)

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

	ok := setDatePaymentAcceptedSalesOrder(saleOrder.Id)
	if !ok {
		trans.Rollback()
		return false
	}

	ok, invoiceId := invoice.insertSalesInvoice()
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
		invoiceDetal.Product = orderDetail.Product
		invoiceDetal.Quantity = orderDetail.Quantity
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		ok = invoiceDetal.insertSalesInvoiceDetail(false)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	go ecommerceControllerupdateStatusPaymentAccepted(saleOrderId)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type OrderDetailGenerate struct {
	OrderId   int32                          `json:"orderId"`
	Selection []OrderDetailGenerateSelection `json:"selection"`
}

type OrderDetailGenerateSelection struct {
	Id       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (invoiceInfo *OrderDetailGenerate) invoicePartiallySaleOrder() bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(invoiceInfo.OrderId)
	if saleOrder.Id <= 0 || len(invoiceInfo.Selection) == 0 {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != invoiceInfo.OrderId || invoiceInfo.Selection[i].Quantity == 0 || invoiceInfo.Selection[i].Quantity > orderDetail.Quantity {
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

	ok := setDatePaymentAcceptedSalesOrder(saleOrder.Id)
	if !ok {
		trans.Rollback()
		return false
	}

	ok, invoiceId := invoice.insertSalesInvoice()
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
		invoiceDetal.Product = orderDetail.Product
		invoiceDetal.Quantity = invoiceInfo.Selection[i].Quantity
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		ok = invoiceDetal.insertSalesInvoiceDetail(false)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	go ecommerceControllerupdateStatusPaymentAccepted(invoiceInfo.OrderId)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type SalesInvoiceRelations struct {
	Orders        []SaleOrder         `json:"orders"`
	DeliveryNotes []SalesDeliveryNote `json:"notes"`
}

func getSalesInvoiceRelations(invoiceId int32) SalesInvoiceRelations {
	return SalesInvoiceRelations{
		Orders:        getSalesInvoiceOrders(invoiceId),
		DeliveryNotes: getSalesInvoiceDeliveryNotes(invoiceId),
	}
}

func getSalesInvoiceOrders(invoiceId int32) []SaleOrder {
	var orders []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT DISTINCT sales_order.* FROM sales_invoice INNER JOIN sales_invoice_detail ON sales_invoice.id = sales_invoice_detail.invoice INNER JOIN sales_order_detail ON sales_invoice_detail.order_detail = sales_order_detail.id INNER JOIN sales_order ON sales_order_detail.order = sales_order.id WHERE sales_invoice.id = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.PrestaShopId,
			&s.WooCommerceId)
		orders = append(orders, s)
	}

	return orders
}

func getSalesInvoiceDeliveryNotes(invoiceId int32) []SalesDeliveryNote {
	var notes []SalesDeliveryNote = make([]SalesDeliveryNote, 0)
	sqlStatement := `SELECT DISTINCT sales_delivery_note.* FROM sales_invoice INNER JOIN sales_invoice_detail ON sales_invoice.id = sales_invoice_detail.invoice INNER JOIN sales_order_detail ON sales_invoice_detail.order_detail = sales_order_detail.id INNER JOIN warehouse_movement ON warehouse_movement.sales_order_detail=sales_order_detail.id INNER JOIN sales_delivery_note ON sales_delivery_note.id=warehouse_movement.sales_delivery_note WHERE sales_invoice.id=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, invoiceId)
	if err != nil {
		log("DB", err.Error())
		return notes
	}
	for rows.Next() {
		n := SalesDeliveryNote{}
		rows.Scan(&n.Id, &n.Warehouse, &n.Customer, &n.DateCreated, &n.PaymentMethod, &n.BillingSeries, &n.ShippingAddress, &n.TotalProducts, &n.DiscountPercent, &n.FixDiscount, &n.ShippingPrice, &n.ShippingDiscount, &n.TotalWithDiscount, &n.VatAmount, &n.TotalAmount, &n.LinesNumber, &n.DeliveryNoteName, &n.DeliveryNoteNumber, &n.Currency, &n.CurrencyChange)
		notes = append(notes, n)
	}

	return notes
}
