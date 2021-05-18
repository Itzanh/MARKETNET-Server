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
	var invoiced []SalesInvoice
	sqlStatement := `SELECT * FROM sales_invoice ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return invoiced
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName)
		invoiced = append(invoiced, i)
	}

	return invoiced
}

func (i *SalesInvoice) isValid() bool {
	return !(i.Customer <= 0 || i.PaymentMethod <= 0 || len(i.BillingSeries) == 0 || i.Currency <= 0 || i.BillingAddress <= 0)
}

func (i *SalesInvoice) insertSalesInvoice() bool {
	if !i.isValid() {
		return false
	}

	i.InvoiceNumber = getNextInvoiceNumber(i.BillingSeries)
	if i.InvoiceNumber <= 0 {
		return false
	}
	i.CurrencyChange = getCurrencyExchange(i.Currency)
	now := time.Now()
	i.InvoiceName = i.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	sqlStatement := `INSERT INTO public.sales_invoice(customer, payment_method, billing_series, currency, currency_change, billing_address, discount_percent, fix_discount, shipping_price, shipping_discount, total_with_discount, total_amount, invoice_number, invoice_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	res, err := db.Exec(sqlStatement, i.Customer, i.PaymentMethod, i.BillingSeries, i.Currency, i.CurrencyChange, i.BillingAddress, i.DiscountPercent, i.FixDiscount, i.ShippingPrice, i.ShippingDiscount, i.TotalWithDiscount, i.TotalAmount, i.InvoiceNumber, i.InvoiceName)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
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
