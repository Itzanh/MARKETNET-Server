package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type SaleOrder struct {
	Id                 int32      `json:"id"`
	Warehouse          string     `json:"warehouse"`
	Reference          string     `json:"reference"`
	Customer           int32      `json:"customer"`
	DateCreated        time.Time  `json:"dateCreated"`
	DatePaymetAccepted *time.Time `json:"datePaymetAccepted"`
	PaymentMethod      int16      `json:"paymentMethod"`
	BillingSeries      string     `json:"billingSeries"`
	Currency           int16      `json:"currency"`
	CurrencyChange     float32    `json:"currencyChange"`
	BillingAddress     int32      `json:"billingAddress"`
	ShippingAddress    int32      `json:"shippingAddress"`
	LinesNumber        int16      `json:"linesNumber"`
	InvoicedLines      int16      `json:"invoicedLines"`
	DeliveryNoteLines  int16      `json:"deliveryNoteLines"`
	TotalProducts      float32    `json:"totalProducts"`
	DiscountPercent    float32    `json:"discountPercent"`
	FixDiscount        float32    `json:"fixDiscount"`
	ShippingPrice      float32    `json:"shippingPrice"`
	ShippingDiscount   float32    `json:"shippingDiscount"`
	TotalWithDiscount  float32    `json:"totalWithDiscount"`
	VatAmount          float32    `json:"vatAmount"`
	TotalAmount        float32    `json:"totalAmount"`
	Description        string     `json:"description"`
	Notes              string     `json:"notes"`
	Off                bool       `json:"off"`
	Cancelled          bool       `json:"cancelled"`
	Status             string     `json:"status"`
	OrderNumber        int32      `json:"orderNumber"`
	BillingStatus      string     `json:"billingStatus"`
	OrderName          string     `json:"orderName"`
}

func getSalesOrder() []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT * FROM sales_order ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
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

func getSalesOrderPreparation() []SaleOrder {
	return getSalesOrderStatus("E")
}

func getSalesOrderAwaitingShipping() []SaleOrder {
	return getSalesOrderStatus("F")
}

func getSalesOrderStatus(status string) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT * FROM sales_order WHERE status = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, status)
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

func getSalesOrderRow(id int32) SaleOrder {
	sqlStatement := `SELECT * FROM sales_order WHERE id = $1 ORDER BY date_created DESC`
	rows := db.QueryRow(sqlStatement, id)
	if rows.Err() != nil {
		return SaleOrder{}
	}

	s := SaleOrder{}
	rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
		&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
		&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName)

	return s
}

func (s *SaleOrder) isValid() bool {
	return !(len(s.Warehouse) == 0 || len(s.Reference) > 9 || s.Customer <= 0 || s.PaymentMethod <= 0 || len(s.BillingSeries) == 0 || s.Currency <= 0 || s.BillingAddress <= 0 || s.ShippingAddress <= 0 || len(s.Notes) > 250)
}

func (s *SaleOrder) insertSalesOrder() bool {
	if !s.isValid() {
		return false
	}

	s.OrderNumber = getNextOrderNumber(s.BillingSeries)
	if s.OrderNumber <= 0 {
		return false
	}
	s.CurrencyChange = getCurrencyExchange(s.Currency)
	now := time.Now()
	s.OrderName = s.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", s.OrderNumber)

	sqlStatement := `INSERT INTO public.sales_order(warehouse, reference, customer, payment_method, billing_series, currency, currency_change, billing_address, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, dsc, notes, order_number, order_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	res, err := db.Exec(sqlStatement, s.Warehouse, s.Reference, s.Customer, s.PaymentMethod, s.BillingSeries, s.Currency, s.CurrencyChange, s.BillingAddress, s.ShippingAddress, s.DiscountPercent, s.FixDiscount, s.ShippingPrice, s.ShippingDiscount, s.Description, s.Notes, s.OrderNumber, s.OrderName)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *SaleOrder) updateSalesOrder() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	inMemoryOrder := getSalesOrderRow(s.Id)
	if inMemoryOrder.Id <= 0 {
		trans.Rollback()
		return false
	}

	var res sql.Result
	var err error
	if inMemoryOrder.Status == "_" { // if the payment is pending, we allow to change more fields
		if s.Currency != inMemoryOrder.Currency {
			s.CurrencyChange = getCurrencyExchange(s.Currency)
		}

		sqlStatement := `UPDATE sales_order SET customer=$2, payment_method=$3, currency=$4, currency_change=$5, billing_address=$6, shipping_address=$7, discount_percent=$8, fix_discount=$9, shipping_price=$10, shipping_discount=$11, dsc=$12, notes=$13, reference=$14 WHERE id = $1`
		res, err = db.Exec(sqlStatement, s.Id, s.Customer, s.PaymentMethod, s.Currency, s.CurrencyChange, s.BillingAddress, s.ShippingAddress, s.DiscountPercent, s.FixDiscount, s.ShippingPrice, s.ShippingDiscount, s.Description, s.Notes, s.Reference)

		if s.DiscountPercent != inMemoryOrder.DiscountPercent || s.FixDiscount != inMemoryOrder.FixDiscount || s.ShippingPrice != inMemoryOrder.ShippingPrice || s.ShippingDiscount != inMemoryOrder.ShippingDiscount {
			ok := calcTotalsSaleOrder(s.Id)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	} else {
		sqlStatement := `UPDATE sales_order SET customer=$2, billing_address=$3, shipping_address=$4, dsc=$5, notes=$6, reference=$7 WHERE id = $1`
		res, err = db.Exec(sqlStatement, s.Id, s.Customer, s.BillingAddress, s.ShippingAddress, s.Description, s.Notes, s.Reference)
	}

	if err != nil {
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

func (s *SaleOrder) deleteSalesOrder() bool {
	if s.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.sales_order WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesOrder(orderId int32, totalAmount float32, vatPercent float32) bool {
	sqlStatement := `UPDATE sales_order SET total_products = total_products + $2, vat_amount = vat_amount + $3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		return false
	}

	return calcTotalsSaleOrder(orderId)
}

// Adds the discounts to the fix discount of the order. This function will substract if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addDiscountsSalesOrder(orderId int32, amountTaxExcluded float32) bool {
	sqlStatement := `UPDATE sales_order SET fix_discount=fix_discount+$2 WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId, amountTaxExcluded)
	if err != nil {
		return false
	}

	return calcTotalsSaleOrder(orderId)
}

// Applies the logic to calculate the totals of the sales order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleOrder(orderId int32) bool {
	sqlStatement := `UPDATE sales_order SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId)
	if err != nil {
		return false
	}

	sqlStatement = `UPDATE sales_order SET total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err = db.Exec(sqlStatement, orderId)
	return err == nil
}

type SaleOrderDefaults struct {
	Warehouse     string `json:"warehouse"`
	WarehouseName string `json:"warehouseName"`
}

func getSaleOrderDefaults() SaleOrderDefaults {
	return SaleOrderDefaults{Warehouse: "W1", WarehouseName: "Main Warehouse"}
}

type SalesOrderRelations struct {
	Invoices            []SalesInvoice       `json:"invoices"`
	ManufacturingOrders []ManufacturingOrder `json:"manufacturingOrders"`
}

func getSalesOrderRelations(orderId int32) SalesOrderRelations {
	return SalesOrderRelations{
		Invoices:            getSalesOrderInvoices(orderId),
		ManufacturingOrders: getSalesOrderManufacturingOrders(orderId),
	}
}

func getSalesOrderInvoices(orderId int32) []SalesInvoice {
	// INVOICE
	var invoiced []SalesInvoice = make([]SalesInvoice, 0)
	sqlStatement := `SELECT DISTINCT sales_invoice.* FROM sales_order INNER JOIN sales_order_detail ON sales_order.id = sales_order_detail.order INNER JOIN sales_invoice_detail ON sales_order_detail.id = sales_invoice_detail.order_detail INNER JOIN sales_invoice ON sales_invoice.id = sales_invoice_detail.invoice WHERE sales_order.id = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId)
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

func getSalesOrderManufacturingOrders(orderId int32) []ManufacturingOrder {
	// MANUFACTURING ORDERS
	var manufacturingOrders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE "order" = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return manufacturingOrders
	}
	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted)
		manufacturingOrders = append(manufacturingOrders, o)
	}

	return manufacturingOrders
}

func setSalesOrderState(orderId int32) bool {
	sqlStatement := `SELECT status FROM sales_order_detail WHERE "order" = $1 ORDER BY status ASC LIMIT 1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		return false
	}

	var status string
	row.Scan(&status)
	if status == "" {
		return false
	}

	sqlStatement = `UPDATE sales_order SET status = $2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, orderId, status)
	rows, _ := res.RowsAffected()

	return rows > 0 && err == nil
}
