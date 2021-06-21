package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type PurchaseOrder struct {
	Id                int32      `json:"id"`
	Warehouse         string     `json:"warehouse"`
	SupplierReference string     `json:"supplierReference"`
	Supplier          int32      `json:"supplier"`
	DateCreated       time.Time  `json:"dateCreated"`
	DatePaid          *time.Time `json:"datePaid"`
	PaymentMethod     int16      `json:"paymentMethod"`
	BillingSeries     string     `json:"billingSeries"`
	Currency          int16      `json:"currency"`
	CurrencyChange    float32    `json:"currencyChange"`
	BillingAddress    int32      `json:"billingAddress"`
	ShippingAddress   int32      `json:"shippingAddress"`
	LinesNumber       int16      `json:"linesNumber"`
	InvoicedLines     int16      `json:"invoicedLines"`
	DeliveryNoteLines int16      `json:"deliveryNoteLines"`
	TotalProducts     float32    `json:"totalProducts"`
	DiscountPercent   float32    `json:"discountPercent"`
	FixDiscount       float32    `json:"fixDiscount"`
	ShippingPrice     float32    `json:"shippingPrice"`
	ShippingDiscount  float32    `json:"shippingDiscount"`
	TotalWithDiscount float32    `json:"totalWithDiscount"`
	VatAmount         float32    `json:"vatAmount"`
	TotalAmount       float32    `json:"totalAmount"`
	Description       string     `json:"description"`
	Notes             string     `json:"notes"`
	Off               bool       `json:"off"`
	Cancelled         bool       `json:"cancelled"`
	OrderNumber       int32      `json:"orderNumber"`
	BillingStatus     string     `json:"billingStatus"`
	OrderName         string     `json:"orderName"`
}

func getPurchaseOrder() []PurchaseOrder {
	var purchases []PurchaseOrder = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT * FROM purchase_order ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return purchases
	}
	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName)
		purchases = append(purchases, s)
	}

	return purchases
}

func (s *OrderSearch) searchPurchaseOrder() []PurchaseOrder {
	var purchases []PurchaseOrder = make([]PurchaseOrder, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_order.* FROM purchase_order WHERE order_number=$1 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT purchase_order.* FROM purchase_order INNER JOIN suppliers ON suppliers.id=purchase_order.supplier WHERE (suppliers.name ILIKE $1 OR purchase_order.supplier_reference ILIKE $1)`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_order.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_order.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		return purchases
	}
	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName)
		purchases = append(purchases, s)
	}

	return purchases
}

func getPurchaseOrderRow(orderId int32) PurchaseOrder {
	sqlStatement := `SELECT * FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		return PurchaseOrder{}
	}

	p := PurchaseOrder{}
	row.Scan(&p.Id, &p.Warehouse, &p.SupplierReference, &p.Supplier, &p.DateCreated, &p.DatePaid, &p.PaymentMethod, &p.BillingSeries, &p.Currency, &p.CurrencyChange,
		&p.BillingAddress, &p.ShippingAddress, &p.LinesNumber, &p.InvoicedLines, &p.DeliveryNoteLines, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount,
		&p.TotalWithDiscount, &p.VatAmount, &p.TotalAmount, &p.Description, &p.Notes, &p.Off, &p.Cancelled, &p.OrderNumber, &p.BillingStatus, &p.OrderName)

	return p
}

func (p *PurchaseOrder) isValid() bool {
	return !(len(p.Warehouse) == 0 || len(p.SupplierReference) > 40 || p.Supplier <= 0 || p.PaymentMethod <= 0 || len(p.BillingSeries) == 0 || p.Currency <= 0 || p.BillingAddress <= 0 || p.ShippingAddress <= 0 || len(p.Notes) > 250)
}

func (p *PurchaseOrder) insertPurchaseOrder() (bool, int32) {
	if !p.isValid() {
		return false, 0
	}

	p.OrderNumber = getNextPurchaseOrderNumber(p.BillingSeries)
	if p.OrderNumber <= 0 {
		return false, 0
	}
	p.CurrencyChange = getCurrencyExchange(p.Currency)
	now := time.Now()
	p.OrderName = p.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", p.OrderNumber)

	sqlStatement := `INSERT INTO public.purchase_order(warehouse, supplier_reference, supplier, payment_method, billing_series, currency, currency_change, billing_address, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, dsc, notes, order_number, order_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`
	row := db.QueryRow(sqlStatement, p.Warehouse, p.SupplierReference, p.Supplier, p.PaymentMethod, p.BillingSeries, p.Currency, p.CurrencyChange, p.BillingAddress, p.ShippingAddress, p.DiscountPercent, p.FixDiscount, p.ShippingPrice, p.ShippingDiscount, p.Description, p.Notes, p.OrderNumber, p.OrderName)
	if row.Err() != nil {
		return false, 0
	}

	var invoiceId int32
	row.Scan(&invoiceId)
	return invoiceId > 0, invoiceId
}

func (p *PurchaseOrder) updatePurchaseOrder() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	inMemoryOrder := getPurchaseOrderRow(p.Id)
	if inMemoryOrder.Id <= 0 {
		trans.Rollback()
		return false
	}

	var res sql.Result
	var err error
	if inMemoryOrder.InvoicedLines == 0 { // if the payment is pending, we allow to change more fields
		if p.Currency != inMemoryOrder.Currency {
			p.CurrencyChange = getCurrencyExchange(p.Currency)
		}

		sqlStatement := `UPDATE purchase_order SET supplier=$2, payment_method=$3, currency=$4, currency_change=$5, billing_address=$6, shipping_address=$7, discount_percent=$8, fix_discount=$9, shipping_price=$10, shipping_discount=$11, dsc=$12, notes=$13, supplier_reference=$14 WHERE id = $1`
		res, err = db.Exec(sqlStatement, p.Id, p.Supplier, p.PaymentMethod, p.Currency, p.CurrencyChange, p.BillingAddress, p.ShippingAddress, p.DiscountPercent, p.FixDiscount, p.ShippingPrice, p.ShippingDiscount, p.Description, p.Notes, p.SupplierReference)

		if p.DiscountPercent != inMemoryOrder.DiscountPercent || p.FixDiscount != inMemoryOrder.FixDiscount || p.ShippingPrice != inMemoryOrder.ShippingPrice || p.ShippingDiscount != inMemoryOrder.ShippingDiscount {
			ok := calcTotalsPurchaseOrder(p.Id)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	} else {
		sqlStatement := `UPDATE purchase_order SET supplier=$2, billing_address=$3, shipping_address=$4, dsc=$5, notes=$6, supplier_reference=$7 WHERE id = $1`
		res, err = db.Exec(sqlStatement, p.Id, p.Supplier, p.BillingAddress, p.ShippingAddress, p.Description, p.Notes, p.SupplierReference)
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

func (p *PurchaseOrder) deletePurchaseOrder() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.purchase_order WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseOrder(orderId int32, totalAmount float32, vatPercent float32) bool {
	sqlStatement := `UPDATE purchase_order SET total_products=total_products+$2,total_vat=total_vat+$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, orderId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return calcTotalsPurchaseOrder(orderId)
}

// If the payment accepted date is null, sets it to the current date and time.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setDatePaymentAcceptedPurchaseOrder(orderId int32) bool {
	sqlStatement := `UPDATE purchase_order SET date_paid=CASE WHEN date_paid IS NOT NULL THEN date_paid ELSE CURRENT_TIMESTAMP(3) END WHERE id=$1`
	_, err := db.Exec(sqlStatement, orderId)
	return err == nil
}

// Applies the logic to calculate the totals of the purchase order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseOrder(orderId int32) bool {
	sqlStatement := `UPDATE purchase_order SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err := db.Exec(sqlStatement, orderId)
	if err != nil {
		return false
	}

	sqlStatement = `UPDATE purchase_order SET total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err = db.Exec(sqlStatement, orderId)
	return err == nil
}

type PurchaseOrderDefaults struct {
	Warehouse     string `json:"warehouse"`
	WarehouseName string `json:"warehouseName"`
}

func getPurchaseOrderDefaults() PurchaseOrderDefaults {
	s := getSettingsRecord()
	return PurchaseOrderDefaults{Warehouse: s.DefaultWarehouse, WarehouseName: s.DefaultWarehouseName}
}

type PurchaseOrderRelations struct {
	Invoices      []PurchaseInvoice      `json:"invoices"`
	DeliveryNotes []PurchaseDeliveryNote `json:"deliveryNotes"`
}

func getPurchaseOrderRelations(orderId int32) PurchaseOrderRelations {
	return PurchaseOrderRelations{
		Invoices:      getPurchaseOrderInvoices(orderId),
		DeliveryNotes: getPurchaseOrderDeliveryNotes(orderId),
	}
}

func getPurchaseOrderInvoices(orderId int32) []PurchaseInvoice {
	// INVOICE
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	sqlStatement := `SELECT DISTINCT purchase_invoice.* FROM purchase_order INNER JOIN purchase_order_detail ON purchase_order.id=purchase_order_detail.order INNER JOIN purchase_invoice_details ON purchase_order_detail.id=purchase_invoice_details.order_detail INNER JOIN purchase_invoice ON purchase_invoice.id=purchase_invoice_details.invoice WHERE purchase_order.id=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return invoices
	}
	for rows.Next() {
		i := PurchaseInvoice{}
		rows.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName)
		invoices = append(invoices, i)
	}

	return invoices
}

func getPurchaseOrderDeliveryNotes(orderId int32) []PurchaseDeliveryNote {
	// DELIVERY NOTES
	var products []PurchaseDeliveryNote = make([]PurchaseDeliveryNote, 0)
	sqlStatement := `SELECT DISTINCT purchase_delivery_note.* FROM purchase_order_detail INNER JOIN warehouse_movement ON warehouse_movement.purchase_order_detail=purchase_order_detail.id INNER JOIN purchase_delivery_note ON warehouse_movement.purchase_delivery_note=purchase_delivery_note.id WHERE purchase_order_detail."order"=$1`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return products
	}
	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange)
		products = append(products, p)
	}

	return products
}
