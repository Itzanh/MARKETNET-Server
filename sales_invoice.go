package main

import (
	"fmt"
	"strconv"
	"time"
)

type SalesInvoice struct {
	Id                int32     `json:"id"`
	Customer          int32     `json:"customer"`
	DateCreated       time.Time `json:"dateCreated"`
	PaymentMethod     int16     `json:"paymentMethod"`
	BillingSeries     string    `json:"billingSeries"`
	Currency          int16     `json:"currency"`
	CurrencyChange    float32   `json:"currencyChange"`
	BillingAddress    int32     `json:"billingAddress"`
	TotalProducts     float32   `json:"totalProducts"`
	DiscountPercent   float32   `json:"discountPercent"`
	FixDiscount       float32   `json:"fixDiscount"`
	ShippingPrice     float32   `json:"shippingPrice"`
	ShippingDiscount  float32   `json:"shippingDiscount"`
	TotalWithDiscount float32   `json:"totalWithDiscount"`
	VatAmount         float32   `json:"vatAmount"`
	TotalAmount       float32   `json:"totalAmount"`
	LinesNumber       int16     `json:"linesNumber"`
	InvoiceNumber     int32     `json:"invoiceNumber"`
	InvoiceName       string    `json:"invoiceName"`
}

func getSalesInvoices() []SalesInvoice {
	var invoices []SalesInvoice = make([]SalesInvoice, 0)
	sqlStatement := `SELECT * FROM sales_invoice ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return invoices
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName)
		invoices = append(invoices, i)
	}

	return invoices
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

	sqlStatement := `DELETE FROM public.sales_invoice WHERE id=$1`
	res, err := db.Exec(sqlStatement, i.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesInvoice(invoiceId int32, totalAmount float32, vatPercent float32) bool {
	sqlStatement := `UPDATE sales_invoice SET total_products = total_products + $2, vat_amount = vat_amount + $3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, invoiceId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
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
		return false
	}

	sqlStatement = `UPDATE sales_invoice SET total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err = db.Exec(sqlStatement, invoiceId)
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

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type SalesOrderDetailInvoice struct {
	SaleOrderId int32                              `json:"saleOrderId"`
	Selection   []SalesOrderDetailInvoiceSelection `json:"selection"`
}

type SalesOrderDetailInvoiceSelection struct {
	Id       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (invoiceInfo *SalesOrderDetailInvoice) invoicePartiallySaleOrder() bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(invoiceInfo.SaleOrderId)
	if saleOrder.Id <= 0 || len(invoiceInfo.Selection) == 0 {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != invoiceInfo.SaleOrderId || invoiceInfo.Selection[i].Quantity == 0 || invoiceInfo.Selection[i].Quantity > orderDetail.Quantity {
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

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type SalesInvoiceRelations struct {
	Orders []SaleOrder `json:"orders"`
}

func getSalesInvoiceRelations(invoiceId int32) SalesInvoiceRelations {
	return SalesInvoiceRelations{Orders: getSalesInvoiceOrders(invoiceId)}
}

func getSalesInvoiceOrders(orderId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT DISTINCT sales_order.* FROM sales_invoice INNER JOIN sales_invoice_detail ON sales_invoice.id = sales_invoice_detail.invoice INNER JOIN sales_order_detail ON sales_invoice_detail.order_detail = sales_order_detail.id INNER JOIN sales_order ON sales_order_detail.order = sales_order.id WHERE sales_invoice.id = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return sales
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName)
		sales = append(sales, s)
	}

	return sales
}
