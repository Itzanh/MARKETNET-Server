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
	SupplierName       string    `json:"supplierName"`
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
			&i.AccountingMovement, &i.enterprise, &i.SupplierName)
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
			&i.AccountingMovement, &i.enterprise, &i.SupplierName)
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
		&i.AccountingMovement, &i.enterprise)

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

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseInvoice(invoiceId int64, totalAmount float32, vatPercent float32) bool {
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
		invoiceDetail.Product = orderDetail.Product
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
		invoiceDetail.Product = orderDetail.Product
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
	Orders []PurchaseOrder `json:"orders"`
}

func getPurchaseInvoiceRelations(invoiceId int64, enterpriseId int32) PurchaseInvoiceRelations {
	return PurchaseInvoiceRelations{Orders: getPurchaseInvoiceOrders(invoiceId, enterpriseId)}
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
