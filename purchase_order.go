package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PurchaseOrder struct {
	Id                int64      `json:"id"`
	Warehouse         string     `json:"warehouse"`
	SupplierReference string     `json:"supplierReference"`
	Supplier          int32      `json:"supplier"`
	DateCreated       time.Time  `json:"dateCreated"`
	DatePaid          *time.Time `json:"datePaid"`
	PaymentMethod     int32      `json:"paymentMethod"`
	BillingSeries     string     `json:"billingSeries"`
	Currency          int32      `json:"currency"`
	CurrencyChange    float64    `json:"currencyChange"`
	BillingAddress    int32      `json:"billingAddress"`
	ShippingAddress   int32      `json:"shippingAddress"`
	LinesNumber       int16      `json:"linesNumber"`
	InvoicedLines     int16      `json:"invoicedLines"`
	DeliveryNoteLines int16      `json:"deliveryNoteLines"`
	TotalProducts     float64    `json:"totalProducts"`
	DiscountPercent   float64    `json:"discountPercent"`
	FixDiscount       float64    `json:"fixDiscount"`
	ShippingPrice     float64    `json:"shippingPrice"`
	ShippingDiscount  float64    `json:"shippingDiscount"`
	TotalWithDiscount float64    `json:"totalWithDiscount"`
	VatAmount         float64    `json:"vatAmount"`
	TotalAmount       float64    `json:"totalAmount"`
	Description       string     `json:"description"`
	Notes             string     `json:"notes"`
	Off               bool       `json:"off"`
	Cancelled         bool       `json:"cancelled"`
	OrderNumber       int32      `json:"orderNumber"`
	BillingStatus     string     `json:"billingStatus"`
	OrderName         string     `json:"orderName"`
	SupplierName      string     `json:"supplierName"`
	enterprise        int32
}

type PurchaseOrders struct {
	Rows   int32               `json:"rows"`
	Orders []PurchaseOrder     `json:"orders"`
	Footer PurchaseOrderFooter `json:"footer"`
}

type PurchaseOrderFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseOrder(enterpriseId int32) PurchaseOrders {
	o := PurchaseOrders{}
	o.Orders = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_order.supplier) FROM purchase_order WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return o
	}
	defer rows.Close()

	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.enterprise, &s.SupplierName)
		o.Orders = append(o.Orders, s)
	}

	sqlStatement = `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_order WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return o
	}
	row.Scan(&o.Rows, &o.Footer.TotalProducts, &o.Footer.TotalAmount)

	return o
}

func (s *OrderSearch) searchPurchaseOrder() PurchaseOrders {
	o := PurchaseOrders{}
	o.Orders = make([]PurchaseOrder, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_order.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_order.supplier) FROM purchase_order WHERE order_number=$1 AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT purchase_order.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_order.supplier) FROM purchase_order INNER JOIN suppliers ON suppliers.id=purchase_order.supplier WHERE (suppliers.name ILIKE $1 OR purchase_order.supplier_reference ILIKE $1)`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_order.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_order.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` AND purchase_order.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return o
	}
	defer rows.Close()

	for rows.Next() {
		s := PurchaseOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.SupplierReference, &s.Supplier, &s.DateCreated, &s.DatePaid, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.enterprise, &s.SupplierName)
		o.Orders = append(o.Orders, s)
	}

	var row *sql.Row
	if err == nil {
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_order WHERE order_number=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_order INNER JOIN suppliers ON suppliers.id=purchase_order.supplier WHERE (suppliers.name ILIKE $1 OR purchase_order.supplier_reference ILIKE $1)`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_order.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_order.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` AND purchase_order.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return o
	}
	row.Scan(&o.Rows, &o.Footer.TotalProducts, &o.Footer.TotalAmount)

	return o
}

func getPurchaseOrderRow(orderId int64) PurchaseOrder {
	sqlStatement := `SELECT * FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseOrder{}
	}

	p := PurchaseOrder{}
	row.Scan(&p.Id, &p.Warehouse, &p.SupplierReference, &p.Supplier, &p.DateCreated, &p.DatePaid, &p.PaymentMethod, &p.BillingSeries, &p.Currency, &p.CurrencyChange,
		&p.BillingAddress, &p.ShippingAddress, &p.LinesNumber, &p.InvoicedLines, &p.DeliveryNoteLines, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount,
		&p.TotalWithDiscount, &p.VatAmount, &p.TotalAmount, &p.Description, &p.Notes, &p.Off, &p.Cancelled, &p.OrderNumber, &p.BillingStatus, &p.OrderName, &p.enterprise)

	return p
}

func getPurchaseOrderRowTransaction(orderId int64, trans sql.Tx) PurchaseOrder {
	sqlStatement := `SELECT * FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseOrder{}
	}

	p := PurchaseOrder{}
	row.Scan(&p.Id, &p.Warehouse, &p.SupplierReference, &p.Supplier, &p.DateCreated, &p.DatePaid, &p.PaymentMethod, &p.BillingSeries, &p.Currency, &p.CurrencyChange,
		&p.BillingAddress, &p.ShippingAddress, &p.LinesNumber, &p.InvoicedLines, &p.DeliveryNoteLines, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount,
		&p.TotalWithDiscount, &p.VatAmount, &p.TotalAmount, &p.Description, &p.Notes, &p.Off, &p.Cancelled, &p.OrderNumber, &p.BillingStatus, &p.OrderName, &p.enterprise)

	return p
}

func (p *PurchaseOrder) isValid() bool {
	return !(len(p.Warehouse) == 0 || len(p.SupplierReference) > 40 || p.Supplier <= 0 || p.PaymentMethod <= 0 || len(p.BillingSeries) == 0 || p.Currency <= 0 || p.BillingAddress <= 0 || p.ShippingAddress <= 0 || len(p.Notes) > 250)
}

func (p *PurchaseOrder) insertPurchaseOrder(userId int32, trans *sql.Tx) (bool, int64) {
	if !p.isValid() {
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

	p.OrderNumber = getNextPurchaseOrderNumber(p.BillingSeries, p.enterprise)
	if p.OrderNumber <= 0 {
		return false, 0
	}
	p.CurrencyChange = getCurrencyExchange(p.Currency)
	now := time.Now()
	p.OrderName = p.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", p.OrderNumber)

	sqlStatement := `INSERT INTO public.purchase_order(warehouse, supplier_reference, supplier, payment_method, billing_series, currency, currency_change, billing_address, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, dsc, notes, order_number, order_name, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`
	row := trans.QueryRow(sqlStatement, p.Warehouse, p.SupplierReference, p.Supplier, p.PaymentMethod, p.BillingSeries, p.Currency, p.CurrencyChange, p.BillingAddress, p.ShippingAddress, p.DiscountPercent, p.FixDiscount, p.ShippingPrice, p.ShippingDiscount, p.Description, p.Notes, p.OrderNumber, p.OrderName, p.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false, 0
	}

	var invoiceId int64
	row.Scan(&invoiceId)
	p.Id = invoiceId

	if invoiceId > 0 {
		insertTransactionalLog(p.enterprise, "purchase_order", int(invoiceId), userId, "I")
		json, _ := json.Marshal(p)
		go fireWebHook(p.enterprise, "purchase_order", "POST", string(json))
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

func (p *PurchaseOrder) updatePurchaseOrder(userId int32) bool {
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
	if inMemoryOrder.Id <= 0 || inMemoryOrder.enterprise != p.enterprise {
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
		res, err = trans.Exec(sqlStatement, p.Id, p.Supplier, p.PaymentMethod, p.Currency, p.CurrencyChange, p.BillingAddress, p.ShippingAddress, p.DiscountPercent, p.FixDiscount, p.ShippingPrice, p.ShippingDiscount, p.Description, p.Notes, p.SupplierReference)

		if p.DiscountPercent != inMemoryOrder.DiscountPercent || p.FixDiscount != inMemoryOrder.FixDiscount || p.ShippingPrice != inMemoryOrder.ShippingPrice || p.ShippingDiscount != inMemoryOrder.ShippingDiscount {
			ok := calcTotalsPurchaseOrder(p.Id, p.enterprise, userId, *trans)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	} else {
		sqlStatement := `UPDATE purchase_order SET supplier=$2, billing_address=$3, shipping_address=$4, dsc=$5, notes=$6, supplier_reference=$7 WHERE id = $1`
		res, err = trans.Exec(sqlStatement, p.Id, p.Supplier, p.BillingAddress, p.ShippingAddress, p.Description, p.Notes, p.SupplierReference)
	}

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

	if rows > 0 {
		insertTransactionalLog(p.enterprise, "purchase_order", int(p.Id), userId, "U")
		json, _ := json.Marshal(p)
		go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))
	}

	return rows > 0
}

// ERROR CODES:
// 1. Alerady invoiced
// 2. Delivery note generated
// 3. Error deleting detail <product>: <error>
func (p *PurchaseOrder) deletePurchaseOrder(userId int32) OkAndErrorCodeReturn {
	if p.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	inMemoryOrder := getPurchaseOrderRow(p.Id)
	if inMemoryOrder.enterprise != p.enterprise {
		return OkAndErrorCodeReturn{Ok: false}
	}

	d := getPurchaseOrderDetail(p.Id, p.enterprise)

	// prevent the order to be deleted if there is an invoice or a delivery note
	for i := 0; i < len(d); i++ {
		if d[i].QuantityInvoiced > 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
		}
		if d[i].QuantityDeliveryNote > 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
		}
	}

	// delete details
	for i := 0; i < len(d); i++ {
		d[i].enterprise = p.enterprise
		ok := d[i].deletePurchaseOrderDetail(userId, trans)
		if !ok.Ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3, ExtraData: []string{strconv.Itoa(int(ok.ErorCode)), d[i].ProductName}}
		}
	}

	insertTransactionalLog(inMemoryOrder.enterprise, "purchase_order", int(p.Id), userId, "D")
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "DELETE", string(json))

	sqlStatement := `DELETE FROM public.purchase_order WHERE id=$1`
	res, err := trans.Exec(sqlStatement, p.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	err = trans.Commit()
	if err != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	rows, _ := res.RowsAffected()
	return OkAndErrorCodeReturn{Ok: rows > 0}
}

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseOrder(orderId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_order SET total_products=total_products+$2,total_vat=total_vat+$3 WHERE id=$1`
	_, err := trans.Exec(sqlStatement, orderId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseOrder(orderId, enterpriseId, userId, trans)
}

// If the payment accepted date is null, sets it to the current date and time.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setDatePaymentAcceptedPurchaseOrder(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_order SET date_paid=CASE WHEN date_paid IS NOT NULL THEN date_paid ELSE CURRENT_TIMESTAMP(3) END WHERE id=$1`
	_, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return true
}

// Applies the logic to calculate the totals of the purchase order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseOrder(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_order SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE purchase_order SET total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err = trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return true
}

type PurchaseOrderDefaults struct {
	Warehouse     string `json:"warehouse"`
	WarehouseName string `json:"warehouseName"`
}

func getPurchaseOrderDefaults(enterpriseId int32) PurchaseOrderDefaults {
	s := getSettingsRecordById(enterpriseId)
	return PurchaseOrderDefaults{Warehouse: s.DefaultWarehouse, WarehouseName: s.DefaultWarehouseName}
}

type PurchaseOrderRelations struct {
	Invoices      []PurchaseInvoice      `json:"invoices"`
	DeliveryNotes []PurchaseDeliveryNote `json:"deliveryNotes"`
}

func getPurchaseOrderRelations(orderId int64, enterpriseId int32) PurchaseOrderRelations {
	return PurchaseOrderRelations{
		Invoices:      getPurchaseOrderInvoices(orderId, enterpriseId),
		DeliveryNotes: getPurchaseOrderDeliveryNotes(orderId, enterpriseId),
	}
}

func getPurchaseOrderInvoices(orderId int64, enterpriseId int32) []PurchaseInvoice {
	// INVOICE
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	sqlStatement := `SELECT DISTINCT purchase_invoice.* FROM purchase_order INNER JOIN purchase_order_detail ON purchase_order.id=purchase_order_detail.order INNER JOIN purchase_invoice_details ON purchase_order_detail.id=purchase_invoice_details.order_detail INNER JOIN purchase_invoice ON purchase_invoice.id=purchase_invoice_details.invoice WHERE purchase_order.id=$1 AND purchase_order.enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	defer rows.Close()

	for rows.Next() {
		i := PurchaseInvoice{}
		rows.Scan(&i.Id, &i.Supplier, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.Amending, &i.AmendedInvoice, &i.enterprise)
		invoices = append(invoices, i)
	}

	return invoices
}

func getPurchaseOrderDeliveryNotes(orderId int64, enterpriseId int32) []PurchaseDeliveryNote {
	// DELIVERY NOTES
	var products []PurchaseDeliveryNote = make([]PurchaseDeliveryNote, 0)
	sqlStatement := `SELECT DISTINCT purchase_delivery_note.* FROM purchase_order_detail INNER JOIN warehouse_movement ON warehouse_movement.purchase_order_detail=purchase_order_detail.id INNER JOIN purchase_delivery_note ON warehouse_movement.purchase_delivery_note=purchase_delivery_note.id WHERE purchase_order_detail."order"=$1 AND purchase_order_detail.enterprise=$2`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	defer rows.Close()

	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise)
		products = append(products, p)
	}

	return products
}

// Add an amount to the lines_number field in the purchase order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderLinesNumber(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET lines_number=lines_number+1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}

// Takes out an amount to the lines_number field in the purchase order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderLinesNumber(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET lines_number=lines_number-1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}

// Add an amount to the invoiced_lines field in the purchase order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderInvoicedLines(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET invoiced_lines=invoiced_lines+1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}

// Takes out an amount to the invoiced_lines field in the purchase order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderInvoicedLines(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET invoiced_lines=invoiced_lines-1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}

// Add an amount to the delivery_note_lines field in the purchase order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderDeliveryNoteLines(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET delivery_note_lines=delivery_note_lines+1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}

// Takes out an amount to the delivery_note_lines field in the purchase order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderDeliveryNoteLines(orderId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.purchase_order SET delivery_note_lines=delivery_note_lines-1 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	p := getPurchaseOrderRowTransaction(orderId, trans)
	json, _ := json.Marshal(p)
	go fireWebHook(p.enterprise, "purchase_order", "PUT", string(json))

	return err == nil && rows > 0
}
