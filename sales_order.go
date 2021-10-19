package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type SaleOrder struct {
	Id                 int64      `json:"id"`
	Warehouse          string     `json:"warehouse"`
	Reference          string     `json:"reference"`
	Customer           int32      `json:"customer"`
	DateCreated        time.Time  `json:"dateCreated"`
	DatePaymetAccepted *time.Time `json:"datePaymetAccepted"`
	PaymentMethod      int32      `json:"paymentMethod"`
	BillingSeries      string     `json:"billingSeries"`
	Currency           int32      `json:"currency"`
	CurrencyChange     float64    `json:"currencyChange"`
	BillingAddress     int32      `json:"billingAddress"`
	ShippingAddress    int32      `json:"shippingAddress"`
	LinesNumber        int16      `json:"linesNumber"`
	InvoicedLines      int16      `json:"invoicedLines"`
	DeliveryNoteLines  int16      `json:"deliveryNoteLines"`
	TotalProducts      float64    `json:"totalProducts"`
	DiscountPercent    float64    `json:"discountPercent"`
	FixDiscount        float64    `json:"fixDiscount"`
	ShippingPrice      float64    `json:"shippingPrice"`
	ShippingDiscount   float64    `json:"shippingDiscount"`
	TotalWithDiscount  float64    `json:"totalWithDiscount"`
	VatAmount          float64    `json:"vatAmount"`
	TotalAmount        float64    `json:"totalAmount"`
	Description        string     `json:"description"`
	Notes              string     `json:"notes"`
	Off                bool       `json:"off"`
	Cancelled          bool       `json:"cancelled"`
	Status             string     `json:"status"`
	OrderNumber        int32      `json:"orderNumber"`
	BillingStatus      string     `json:"billingStatus"`
	OrderName          string     `json:"orderName"`
	Carrier            *int32     `json:"carrier"`
	CustomerName       string     `json:"customerName"`
	prestaShopId       int32
	wooCommerceId      int32
	shopifyId          int64
	shopifyDraftId     int64
	enterprise         int32
}

type SaleOrders struct {
	Rows   int32       `json:"rows"`
	Orders []SaleOrder `json:"orders"`
}

func (q *PaginationQuery) getSalesOrder(enterpriseId int32) SaleOrders {
	so := SaleOrders{}
	if !q.isValid() {
		return so
	}

	so.Orders = make([]SaleOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order WHERE enterprise=$1 ORDER BY date_created DESC OFFSET $2 LIMIT $3`
	rows, err := db.Query(sqlStatement, enterpriseId, q.Offset, q.Limit)
	if err != nil {
		log("DB", err.Error())
		return so
	}

	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.shopifyDraftId, &s.enterprise, &s.CustomerName)
		so.Orders = append(so.Orders, s)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.sales_order WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	row.Scan(&so.Rows)

	return so
}

type SalesOrderSearch struct {
	PaginatedSearch
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
	Status    string     `json:"status"`
}

func (s *SalesOrderSearch) searchSalesOrder() SaleOrders {
	so := SaleOrders{}
	if !s.isValid() {
		return so
	}

	so.Orders = make([]SaleOrder, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT sales_order.*,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order WHERE (order_number=$1 OR id=$1) AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.Enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT sales_order.*,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order INNER JOIN customer ON customer.id=sales_order.customer WHERE (reference ILIKE $1 OR customer.name ILIKE $1)`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_order.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_order.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.Status != "" {
			sqlStatement += ` AND status = $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.Status)
		}
		sqlStatement += ` AND sales_order.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.Enterprise)
		sqlStatement += ` ORDER BY date_created DESC OFFSET $` + strconv.Itoa(len(interfaces)+1) + ` LIMIT $` + strconv.Itoa(len(interfaces)+2)
		interfaces = append(interfaces, s.Offset)
		interfaces = append(interfaces, s.Limit)
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return so
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.shopifyDraftId, &s.enterprise, &s.CustomerName)
		so.Orders = append(so.Orders, s)
	}

	var row *sql.Row
	orderNumber, err = strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT COUNT(*) FROM sales_order WHERE order_number=$1 OR id=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderNumber, s.Enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*) FROM sales_order INNER JOIN customer ON customer.id=sales_order.customer WHERE (reference ILIKE $1 OR customer.name ILIKE $1)`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_order.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_order.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		if s.Status != "" {
			sqlStatement += ` AND status = $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.Status)
		}
		sqlStatement += ` AND sales_order.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.Enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return so
	}
	row.Scan(&so.Rows)

	return so
}

func getSalesOrderPreparation(enterpriseId int32) []SaleOrder {
	return getSalesOrderStatus("E", enterpriseId)
}

func getSalesOrderAwaitingShipping(enterpriseId int32) []SaleOrder {
	return getSalesOrderStatus("F", enterpriseId)
}

func getSalesOrderStatus(status string, enterpriseId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_order.customer) FROM sales_order WHERE status = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, status)
	if err != nil {
		return sales
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.shopifyDraftId, &s.enterprise, &s.CustomerName)
		sales = append(sales, s)
	}

	return sales
}

func getSalesOrderRow(id int64) SaleOrder {
	sqlStatement := `SELECT * FROM sales_order WHERE id = $1 ORDER BY date_created DESC`
	rows := db.QueryRow(sqlStatement, id)
	if rows.Err() != nil {
		return SaleOrder{}
	}

	s := SaleOrder{}
	rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
		&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
		&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
		&s.wooCommerceId, &s.shopifyId, &s.shopifyDraftId, &s.enterprise)

	return s
}

func (s *SaleOrder) isValid() bool {
	return !(len(s.Warehouse) == 0 || len(s.Reference) > 9 || s.Customer <= 0 || s.PaymentMethod <= 0 || len(s.BillingSeries) == 0 || s.Currency <= 0 || s.BillingAddress <= 0 || s.ShippingAddress <= 0 || len(s.Notes) > 250)
}

func (s *SaleOrder) insertSalesOrder() (bool, int64) {
	if !s.isValid() {
		return false, 0
	}

	s.OrderNumber = getNextSaleOrderNumber(s.BillingSeries, s.enterprise)
	if s.OrderNumber <= 0 {
		return false, 0
	}
	s.CurrencyChange = getCurrencyExchange(s.Currency)
	now := time.Now()
	s.OrderName = s.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", s.OrderNumber)

	sqlStatement := `INSERT INTO public.sales_order(warehouse, reference, customer, payment_method, billing_series, currency, currency_change, billing_address, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, dsc, notes, order_number, order_name, carrier, ps_id, wc_id, sy_draft_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Warehouse, s.Reference, s.Customer, s.PaymentMethod, s.BillingSeries, s.Currency, s.CurrencyChange, s.BillingAddress, s.ShippingAddress, s.DiscountPercent, s.FixDiscount, s.ShippingPrice, s.ShippingDiscount, s.Description, s.Notes, s.OrderNumber, s.OrderName, s.Carrier, s.prestaShopId, s.wooCommerceId, s.shopifyDraftId, s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0
	}

	var orderId int64
	row.Scan(&orderId)
	return orderId > 0, orderId
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

		sqlStatement := `UPDATE sales_order SET customer=$2, payment_method=$3, currency=$4, currency_change=$5, billing_address=$6, shipping_address=$7, discount_percent=$8, fix_discount=$9, shipping_price=$10, shipping_discount=$11, dsc=$12, notes=$13, reference=$14, carrier=$15, sy_id=$16 WHERE id=$1 AND enterprise=$17`
		res, err = db.Exec(sqlStatement, s.Id, s.Customer, s.PaymentMethod, s.Currency, s.CurrencyChange, s.BillingAddress, s.ShippingAddress, s.DiscountPercent, s.FixDiscount, s.ShippingPrice, s.ShippingDiscount, s.Description, s.Notes, s.Reference, s.Carrier, s.shopifyId, s.enterprise)

		if s.DiscountPercent != inMemoryOrder.DiscountPercent || s.FixDiscount != inMemoryOrder.FixDiscount || s.ShippingPrice != inMemoryOrder.ShippingPrice || s.ShippingDiscount != inMemoryOrder.ShippingDiscount {
			ok := calcTotalsSaleOrder(s.Id)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	} else {
		sqlStatement := `UPDATE sales_order SET customer=$2, billing_address=$3, shipping_address=$4, dsc=$5, notes=$6, reference=$7, carrier=$8 WHERE id=$1 AND enterprise=$9`
		res, err = db.Exec(sqlStatement, s.Id, s.Customer, s.BillingAddress, s.ShippingAddress, s.Description, s.Notes, s.Reference, s.Carrier, s.enterprise)
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
	return rows > 0
}

func (s *SaleOrder) deleteSalesOrder() bool {
	if s.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	d := getSalesOrderDetail(s.Id, s.enterprise)

	// prevent the order to be deleted if there is an invoice or a delivery note
	for i := 0; i < len(d); i++ {
		if d[i].QuantityInvoiced > 0 || d[i].QuantityDeliveryNote > 0 {
			trans.Rollback()
			return false
		}
	}

	// delete sales order detail packaged
	sqlStatement := `DELETE FROM sales_order_detail_packaged WHERE order_detail = $1 AND enterprise = $2`
	for i := 0; i < len(d); i++ {
		_, err := db.Exec(sqlStatement, d[i].Id, s.enterprise)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
	}

	// delete packaging
	sqlStatement = `DELETE FROM packaging WHERE sales_order = $1 AND enterprise = $2`
	_, err := db.Exec(sqlStatement, s.Id, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// delete pallets
	sqlStatement = `DELETE FROM pallets WHERE sales_order = $1 AND enterprise = $2`
	_, err = db.Exec(sqlStatement, s.Id, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	// delete details
	for i := 0; i < len(d); i++ {
		d[i].enterprise = s.enterprise
		ok := d[i].deleteSalesOrderDetail()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// delete sale order
	sqlStatement = `DELETE FROM public.sales_order WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, s.Id, s.enterprise)
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

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesOrder(orderId int64, totalAmount float64, vatPercent float64) bool {
	sqlStatement := `UPDATE sales_order SET total_products = total_products + $2, vat_amount = vat_amount + $3 WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsSaleOrder(orderId)
}

// Adds the discounts to the fix discount of the order. This function will substract if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addDiscountsSalesOrder(orderId int64, amountTaxExcluded float64) bool {
	sqlStatement := `UPDATE sales_order SET fix_discount=fix_discount+$2 WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId, amountTaxExcluded)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsSaleOrder(orderId)
}

// If the payment accepted date is null, sets it to the current date and time.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setDatePaymentAcceptedSalesOrder(orderId int64) bool {
	sqlStatement := `UPDATE sales_order SET date_payment_accepted=CASE WHEN date_payment_accepted IS NOT NULL THEN date_payment_accepted ELSE CURRENT_TIMESTAMP(3) END WHERE id=$1`
	_, err := db.Exec(sqlStatement, orderId)
	return err == nil
}

// Applies the logic to calculate the totals of the sales order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleOrder(orderId int64) bool {
	sqlStatement := `UPDATE sales_order SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := db.Exec(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
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

func getSaleOrderDefaults(enterpriseId int32) SaleOrderDefaults {
	s := getSettingsRecordById(enterpriseId)
	return SaleOrderDefaults{Warehouse: s.DefaultWarehouse, WarehouseName: s.DefaultWarehouseName}
}

type SalesOrderRelations struct {
	Invoices            []SalesInvoice       `json:"invoices"`
	ManufacturingOrders []ManufacturingOrder `json:"manufacturingOrders"`
	DeliveryNotes       []SalesDeliveryNote  `json:"deliveryNotes"`
	Shippings           []Shipping           `json:"shippings"`
}

func getSalesOrderRelations(orderId int64, enterpriseId int32) SalesOrderRelations {
	return SalesOrderRelations{
		Invoices:            getSalesOrderInvoices(orderId, enterpriseId),
		ManufacturingOrders: getSalesOrderManufacturingOrders(orderId, enterpriseId),
		DeliveryNotes:       getSalesOrderDeliveryNotes(orderId, enterpriseId),
		Shippings:           getSalesOrderShippings(orderId, enterpriseId),
	}
}

func getSalesOrderInvoices(orderId int64, enterpriseId int32) []SalesInvoice {
	// INVOICE
	var invoices []SalesInvoice = make([]SalesInvoice, 0)
	sqlStatement := `SELECT DISTINCT sales_invoice.* FROM sales_order INNER JOIN sales_order_detail ON sales_order.id = sales_order_detail.order INNER JOIN sales_invoice_detail ON sales_order_detail.id = sales_invoice_detail.order_detail INNER JOIN sales_invoice ON sales_invoice.id = sales_invoice_detail.invoice WHERE sales_order.id = $1 AND sales_order.enterprise = $2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return invoices
	}
	for rows.Next() {
		i := SalesInvoice{}
		rows.Scan(&i.Id, &i.Customer, &i.DateCreated, &i.PaymentMethod, &i.BillingSeries, &i.Currency, &i.CurrencyChange, &i.BillingAddress, &i.TotalProducts,
			&i.DiscountPercent, &i.FixDiscount, &i.ShippingPrice, &i.ShippingDiscount, &i.TotalWithDiscount, &i.VatAmount, &i.TotalAmount, &i.LinesNumber, &i.InvoiceNumber, &i.InvoiceName,
			&i.AccountingMovement, &i.enterprise, &i.SimplifiedInvoice, &i.Amending, &i.AmendedInvoice)
		invoices = append(invoices, i)
	}

	return invoices
}

func getSalesOrderManufacturingOrders(orderId int64, enterpriseId int32) []ManufacturingOrder {
	// MANUFACTURING ORDERS
	var manufacturingOrders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE "order"=$1 AND enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return manufacturingOrders
	}
	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise)
		manufacturingOrders = append(manufacturingOrders, o)
	}

	return manufacturingOrders
}

func getSalesOrderDeliveryNotes(orderId int64, enterpriseId int32) []SalesDeliveryNote {
	// DELIVERY NOTES
	var notes []SalesDeliveryNote = make([]SalesDeliveryNote, 0)
	sqlStatement := `SELECT DISTINCT sales_delivery_note.* FROM sales_order_detail INNER JOIN warehouse_movement ON warehouse_movement.sales_order_detail = sales_order_detail.id INNER JOIN sales_delivery_note ON warehouse_movement.sales_delivery_note = sales_delivery_note.id WHERE sales_order_detail."order" = $1 AND sales_delivery_note.enterprise=$2`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return notes
	}
	for rows.Next() {
		p := SalesDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Customer, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.VatAmount, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise)
		notes = append(notes, p)
	}

	return notes
}

func getSalesOrderShippings(orderId int64, enterpriseId int32) []Shipping {
	// SHIPPINGS
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE "order"=$1 AND shipping.enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return shippings
	}
	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.CustomerName, &s.SaleOrderName, &s.CarrierName, &s.CarrierWebService, &s.enterprise)
		shippings = append(shippings, s)
	}

	return shippings
}

func setSalesOrderState(orderId int64) bool {
	sqlStatement := `SELECT status FROM sales_order_detail WHERE "order" = $1 ORDER BY status ASC LIMIT 1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var status string
	row.Scan(&status)
	if status == "" {
		status = "_"
	}

	sqlStatement = `UPDATE sales_order SET status = $2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, orderId, status)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return rows > 0 && err == nil
}

type SaleOrderLocate struct {
	Id           int32     `json:"id"`
	Customer     int32     `json:"customer"`
	CustomerName string    `json:"customerName"`
	OrderName    string    `json:"orderName"`
	DateCreated  time.Time `json:"dateCreated"`
}

func locateSaleOrder(enterpriseId int32) []SaleOrderLocate {
	var sales []SaleOrderLocate = make([]SaleOrderLocate, 0)
	sqlStatement := `SELECT id,customer,(SELECT name FROM customer WHERE id=sales_order.customer),order_name,date_created FROM sales_order WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return sales
	}
	for rows.Next() {
		s := SaleOrderLocate{}
		rows.Scan(&s.Id, &s.Customer, &s.CustomerName, &s.OrderName, &s.DateCreated)
		sales = append(sales, s)
	}

	return sales
}

// Add an amount to the lines_number field in the sale order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderLinesNumber(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET lines_number=lines_number+1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

// Takes out an amount to the lines_number field in the sale order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderLinesNumber(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET lines_number=lines_number-1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

// Add an amount to the invoiced_lines field in the sale order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderInvoicedLines(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET invoiced_lines=invoiced_lines+1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

// Takes out an amount to the invoiced_lines field in the sale order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderInvoicedLines(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET invoiced_lines=invoiced_lines-1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

// Add an amount to the delivery_note_lines field in the sale order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderDeliveryNoteLines(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET delivery_note_lines=delivery_note_lines+1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}

// Takes out an amount to the delivery_note_lines field in the sale order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderDeliveryNoteLines(orderId int64) bool {
	sqlStatement := `UPDATE public.sales_order SET delivery_note_lines=delivery_note_lines-1 WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderId)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil && rows > 0
}
