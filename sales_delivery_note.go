package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type SalesDeliveryNote struct {
	Id                 int64     `json:"id"`
	Warehouse          string    `json:"warehouse"`
	Customer           int32     `json:"customer"`
	DateCreated        time.Time `json:"dateCreated"`
	PaymentMethod      int32     `json:"paymentMethod"`
	BillingSeries      string    `json:"billingSeries"`
	ShippingAddress    int32     `json:"shippingAddress"`
	TotalProducts      float64   `json:"totalProducts"`
	DiscountPercent    float64   `json:"discountPercent"`
	FixDiscount        float64   `json:"fixDiscount"`
	ShippingPrice      float64   `json:"shippingPrice"`
	ShippingDiscount   float64   `json:"shippingDiscount"`
	TotalWithDiscount  float64   `json:"totalWithDiscount"`
	VatAmount          float64   `json:"vatAmount"`
	TotalAmount        float64   `json:"totalAmount"`
	LinesNumber        int16     `json:"linesNumber"`
	DeliveryNoteNumber int32     `json:"deliveryNoteNumber"`
	DeliveryNoteName   string    `json:"deliveryNoteName"`
	Currency           int32     `json:"currency"`
	CurrencyChange     float64   `json:"currencyChange"`
	CustomerName       string    `json:"customerName"`
	enterprise         int32
}

type SalesDeliveryNotes struct {
	Rows  int32               `json:"rows"`
	Notes []SalesDeliveryNote `json:"notes"`
}

func (q *PaginationQuery) getSalesDeliveryNotes() SalesDeliveryNotes {
	sd := SalesDeliveryNotes{}
	if !q.isValid() {
		return sd
	}

	sd.Notes = make([]SalesDeliveryNote, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=sales_delivery_note.customer) FROM public.sales_delivery_note WHERE enterprise=$3 ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, q.Enterprise)
	if err != nil {
		log("DB", err.Error())
		return sd
	}
	for rows.Next() {
		n := SalesDeliveryNote{}
		rows.Scan(&n.Id, &n.Warehouse, &n.Customer, &n.DateCreated, &n.PaymentMethod, &n.BillingSeries, &n.ShippingAddress, &n.TotalProducts, &n.DiscountPercent, &n.FixDiscount, &n.ShippingPrice, &n.ShippingDiscount, &n.TotalWithDiscount, &n.VatAmount, &n.TotalAmount, &n.LinesNumber, &n.DeliveryNoteName, &n.DeliveryNoteNumber, &n.Currency, &n.CurrencyChange, &n.enterprise, &n.CustomerName)
		sd.Notes = append(sd.Notes, n)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.sales_delivery_note WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, q.Enterprise)
	row.Scan(&sd.Rows)

	return sd
}

func getSalesDeliveryNoteRow(deliveryNoteId int64) SalesDeliveryNote {
	sqlStatement := `SELECT * FROM public.sales_delivery_note WHERE id=$1`
	row := db.QueryRow(sqlStatement, deliveryNoteId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesDeliveryNote{}
	}

	n := SalesDeliveryNote{}
	row.Scan(&n.Id, &n.Warehouse, &n.Customer, &n.DateCreated, &n.PaymentMethod, &n.BillingSeries, &n.ShippingAddress, &n.TotalProducts, &n.DiscountPercent, &n.FixDiscount, &n.ShippingPrice, &n.ShippingDiscount, &n.TotalWithDiscount, &n.VatAmount, &n.TotalAmount, &n.LinesNumber, &n.DeliveryNoteName, &n.DeliveryNoteNumber, &n.Currency, &n.CurrencyChange, &n.enterprise)

	return n
}

func (s *OrderSearch) searchSalesDelvieryNotes() SalesDeliveryNotes {
	sd := SalesDeliveryNotes{}
	if !s.isValid() {
		return sd
	}

	sd.Notes = make([]SalesDeliveryNote, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT sales_delivery_note.*,(SELECT name FROM customer WHERE customer.id=sales_delivery_note.customer) FROM sales_delivery_note WHERE delivery_note_number=$1 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT sales_delivery_note.*,(SELECT name FROM customer WHERE customer.id=sales_delivery_note.customer) FROM sales_delivery_note INNER JOIN customer ON customer.id=sales_delivery_note.customer WHERE customer.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_delivery_note.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_delivery_note.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` AND sales_delivery_note.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.Enterprise)
		sqlStatement += ` ORDER BY date_created DESC OFFSET $` + strconv.Itoa(len(interfaces)+1) + ` LIMIT $` + strconv.Itoa(len(interfaces)+2)
		interfaces = append(interfaces, s.Offset)
		interfaces = append(interfaces, s.Limit)
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return sd
	}
	for rows.Next() {
		n := SalesDeliveryNote{}
		rows.Scan(&n.Id, &n.Warehouse, &n.Customer, &n.DateCreated, &n.PaymentMethod, &n.BillingSeries, &n.ShippingAddress, &n.TotalProducts, &n.DiscountPercent, &n.FixDiscount, &n.ShippingPrice, &n.ShippingDiscount, &n.TotalWithDiscount, &n.VatAmount, &n.TotalAmount, &n.LinesNumber, &n.DeliveryNoteName, &n.DeliveryNoteNumber, &n.Currency, &n.CurrencyChange, &n.enterprise, &n.CustomerName)
		sd.Notes = append(sd.Notes, n)
	}

	var row *sql.Row
	orderNumber, err = strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT COUNT(*) FROM sales_delivery_note WHERE delivery_note_number=$1`
		row = db.QueryRow(sqlStatement, orderNumber)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*) FROM sales_delivery_note INNER JOIN customer ON customer.id=sales_delivery_note.customer WHERE customer.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND sales_delivery_note.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND sales_delivery_note.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` AND sales_delivery_note.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.Enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return sd
	}
	row.Scan(&sd.Rows)

	return sd
}

func (n *SalesDeliveryNote) isValid() bool {
	return !(len(n.Warehouse) == 0 || len(n.Warehouse) > 2 || n.Customer <= 0 || n.PaymentMethod <= 0 || len(n.BillingSeries) == 0 || len(n.BillingSeries) > 3 || n.ShippingAddress <= 0)
}

func (n *SalesDeliveryNote) insertSalesDeliveryNotes() (bool, int64) {
	if !n.isValid() {
		return false, 0
	}

	n.DeliveryNoteNumber = getNextSaleDeliveryNoteNumber(n.BillingSeries, n.enterprise)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.Currency)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)

	sqlStatement := `INSERT INTO public.sales_delivery_note(warehouse, customer, payment_method, billing_series, shipping_address, delivery_note_number, delivery_note_name, currency, currency_change, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	row := db.QueryRow(sqlStatement, n.Warehouse, n.Customer, n.PaymentMethod, n.BillingSeries, n.ShippingAddress, n.DeliveryNoteNumber, n.DeliveryNoteName, n.Currency, n.CurrencyChange, n.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0
	}

	var noteId int64
	row.Scan(&noteId)
	return noteId > 0, noteId
}

func (n *SalesDeliveryNote) deleteSalesDeliveryNotes() bool {
	if n.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	d := getWarehouseMovementBySalesDeliveryNote(n.Id, n.enterprise)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteWarehouseMovement()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.sales_delivery_note WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, n.Id, n.enterprise)
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

func deliveryNoteAllSaleOrder(saleOrderId int64, enterpriseId int32) (bool, int64) {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.enterprise != enterpriseId {
		return false, 0
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.enterprise)

	if saleOrder.Id <= 0 || len(orderDetails) == 0 {
		return false, 0
	}

	// create a delivery note for that order
	n := SalesDeliveryNote{}
	n.Customer = saleOrder.Customer
	n.ShippingAddress = saleOrder.ShippingAddress
	n.Currency = saleOrder.Currency
	n.PaymentMethod = saleOrder.PaymentMethod
	n.BillingSeries = saleOrder.BillingSeries
	n.Warehouse = saleOrder.Warehouse

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false, 0
	}
	///

	n.enterprise = enterpriseId
	ok, deliveryNoteId := n.insertSalesDeliveryNotes()
	if !ok {
		trans.Rollback()
		return false, 0
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "O"
		movement.Warehouse = saleOrder.Warehouse
		movement.Product = orderDetail.Product
		movement.Quantity = -orderDetail.Quantity
		movement.SalesDeliveryNote = &deliveryNoteId
		movement.SalesOrderDetail = &orderDetail.Id
		movement.SalesOrder = &saleOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.enterprise = enterpriseId
		ok = movement.insertWarehouseMovement()
		if !ok {
			trans.Rollback()
			return false, 0
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil, deliveryNoteId
	///
}

func (noteInfo *OrderDetailGenerate) deliveryNotePartiallySaleOrder(enterpriseId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(noteInfo.OrderId)
	if saleOrder.Id <= 0 || saleOrder.enterprise != enterpriseId || len(noteInfo.Selection) == 0 {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(noteInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != noteInfo.OrderId || noteInfo.Selection[i].Quantity == 0 || noteInfo.Selection[i].Quantity > orderDetail.Quantity {
			return false
		}
		saleOrderDetails = append(saleOrderDetails, orderDetail)
	}

	// create a delivery note for that order
	n := SalesDeliveryNote{}
	n.Customer = saleOrder.Customer
	n.ShippingAddress = saleOrder.ShippingAddress
	n.Currency = saleOrder.Currency
	n.PaymentMethod = saleOrder.PaymentMethod
	n.BillingSeries = saleOrder.BillingSeries
	n.Warehouse = saleOrder.Warehouse

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	n.enterprise = enterpriseId
	ok, deliveryNoteId := n.insertSalesDeliveryNotes()
	if !ok {
		trans.Rollback()
		return false
	}
	for i := 0; i < len(saleOrderDetails); i++ {
		orderDetail := saleOrderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "O"
		movement.Warehouse = saleOrder.Warehouse
		movement.Product = orderDetail.Product
		movement.Quantity = -noteInfo.Selection[i].Quantity
		movement.SalesDeliveryNote = &deliveryNoteId
		movement.SalesOrderDetail = &orderDetail.Id
		movement.SalesOrder = &saleOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.enterprise = enterpriseId
		ok = movement.insertWarehouseMovement()
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

type SalesDeliveryNoteLocate struct {
	Id               int32     `json:"id"`
	CustomerName     string    `json:"customerName"`
	DateCreated      time.Time `json:"dateCreated"`
	DeliveryNoteName string    `json:"deliveryNoteName"`
}

func locateSalesDeliveryNotesBySalesOrder(orderId int64, enterpriseId int32) []SalesDeliveryNoteLocate {
	var products []SalesDeliveryNoteLocate = make([]SalesDeliveryNoteLocate, 0)
	sqlStatement := `SELECT DISTINCT sales_delivery_note.id,(SELECT name FROM customer WHERE id=sales_delivery_note.customer),sales_delivery_note.date_created,sales_delivery_note.delivery_note_name FROM sales_order_detail INNER JOIN warehouse_movement ON warehouse_movement.sales_order_detail = sales_order_detail.id INNER JOIN sales_delivery_note ON warehouse_movement.sales_delivery_note = sales_delivery_note.id WHERE sales_order_detail."order" = $1 AND sales_delivery_note = $2`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := SalesDeliveryNoteLocate{}
		rows.Scan(&p.Id, &p.CustomerName, &p.DateCreated, &p.DeliveryNoteName)
		products = append(products, p)
	}

	return products
}

func getNameSalesDeliveryNote(id int64, enterpriseId int32) string {
	sqlStatement := `SELECT delivery_note_name FROM public.sales_delivery_note WHERE id = $1 AND enterprise =$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type SalesDeliveryNoteRelation struct {
	Orders    []SaleOrder `json:"orders"`
	Shippings []Shipping  `json:"shippings"`
}

func getSalesDeliveryNoteRelations(noteId int64, enterpriseId int32) SalesDeliveryNoteRelation {
	return SalesDeliveryNoteRelation{
		Orders:    getSalesDeliveryNoteOrders(noteId, enterpriseId),
		Shippings: getSalesDeliveryNoteShippings(noteId),
	}
}

func getSalesDeliveryNoteOrders(noteId int64, enterpriseId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT DISTINCT sales_order.* FROM sales_delivery_note INNER JOIN warehouse_movement ON sales_delivery_note.id=warehouse_movement.sales_delivery_note INNER JOIN sales_order ON sales_order.id=warehouse_movement.sales_order WHERE sales_delivery_note.id=$1 AND sales_delivery_note.enterprise=$2 ORDER BY sales_order.date_created DESC`
	rows, err := db.Query(sqlStatement, noteId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return sales
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier, &s.prestaShopId,
			&s.wooCommerceId, &s.shopifyId, &s.enterprise, &s.shopifyDraftId)
		sales = append(sales, s)
	}

	return sales
}

func getSalesDeliveryNoteShippings(noteId int64) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE shipping.delivery_note=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		log("DB", err.Error())
		return shippings
	}
	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.CustomerName, &s.SaleOrderName, &s.CarrierName)
		shippings = append(shippings, s)
	}

	return shippings
}

// Adds a total amount to the delivery note total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesDeliveryNote(noteId int64, totalAmount float64, vatPercent float64) bool {
	sqlStatement := `UPDATE sales_delivery_note SET total_products=total_products+$2, vat_amount=vat_amount+$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, noteId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return calcTotalsSaleDeliveryNote(noteId)
}

// Applies the logic to calculate the totals of the sales delivery note.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleDeliveryNote(noteId int64) bool {
	sqlStatement := `UPDATE sales_delivery_note SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err := db.Exec(sqlStatement, noteId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE sales_delivery_note SET total_amount=total_with_discount+vat_amount WHERE id = $1`
	_, err = db.Exec(sqlStatement, noteId)
	return err == nil
}
