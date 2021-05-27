package main

import (
	"fmt"
	"strconv"
	"time"
)

type SalesDeliveryNote struct {
	Id                 int32     `json:"id"`
	Warehouse          string    `json:"warehouse"`
	Customer           int32     `json:"customer"`
	DateCreated        time.Time `json:"dateCreated"`
	PaymentMethod      int16     `json:"paymentMethod"`
	BillingSeries      string    `json:"billingSeries"`
	ShippingAddress    int32     `json:"shippingAddress"`
	TotalProducts      float32   `json:"totalProducts"`
	DiscountPercent    float32   `json:"discountPercent"`
	FixDiscount        float32   `json:"fixDiscount"`
	ShippingPrice      float32   `json:"shippingPrice"`
	ShippingDiscount   float32   `json:"shippingDiscount"`
	TotalWithDiscount  float32   `json:"totalWithDiscount"`
	TotalVat           float32   `json:"totalVat"`
	TotalAmount        float32   `json:"totalAmount"`
	LinesNumber        int16     `json:"linesNumber"`
	DeliveryNoteNumber int32     `json:"deliveryNoteNumber"`
	DeliveryNoteName   string    `json:"deliveryNoteName"`
	Currency           int16     `json:"currency"`
	CurrencyChange     float32   `json:"currencyChange"`
}

func getSalesDeliveryNotes() []SalesDeliveryNote {
	var products []SalesDeliveryNote = make([]SalesDeliveryNote, 0)
	sqlStatement := `SELECT * FROM public.sales_delivery_note ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return products
	}
	for rows.Next() {
		p := SalesDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Customer, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange)
		products = append(products, p)
	}

	return products
}

func (n *SalesDeliveryNote) isValid() bool {
	return !(len(n.Warehouse) == 0 || len(n.Warehouse) > 2 || n.Customer <= 0 || n.PaymentMethod <= 0 || len(n.BillingSeries) == 0 || len(n.BillingSeries) > 3 || n.ShippingAddress <= 0)
}

func (n *SalesDeliveryNote) insertSalesDeliveryNotes() (bool, int32) {
	if !n.isValid() {
		return false, 0
	}

	n.DeliveryNoteNumber = getNextSaleDeliveryNoteNumber(n.BillingSeries)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.Currency)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)

	sqlStatement := `INSERT INTO public.sales_delivery_note(warehouse, customer, payment_method, billing_series, shipping_address, delivery_note_number, delivery_note_name, currency, currency_change) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := db.QueryRow(sqlStatement, n.Warehouse, n.Customer, n.PaymentMethod, n.BillingSeries, n.ShippingAddress, n.DeliveryNoteNumber, n.DeliveryNoteName, n.Currency, n.CurrencyChange)
	if row.Err() != nil {
		return false, 0
	}

	var invoiceId int32
	row.Scan(&invoiceId)
	return invoiceId > 0, invoiceId
}

func (n *SalesDeliveryNote) deleteSalesDeliveryNotes() bool {
	if n.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.sales_delivery_note WHERE id=$1`
	res, err := db.Exec(sqlStatement, n.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func deliveryNoteAllSaleOrder(saleOrderId int32) (bool, int32) {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	orderDetails := getSalesOrderDetail(saleOrderId)

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

type SalesOrderDetailDeliveryNote struct {
	SaleOrderId int32                                   `json:"saleOrderId"`
	Selection   []SalesOrderDetailDeliveryNoteSelection `json:"selection"`
}

type SalesOrderDetailDeliveryNoteSelection struct {
	Id       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (noteInfo *SalesOrderDetailDeliveryNote) deliveryNotePartiallySaleOrder() bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(noteInfo.SaleOrderId)
	if saleOrder.Id <= 0 || len(noteInfo.Selection) == 0 {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(noteInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != noteInfo.SaleOrderId || noteInfo.Selection[i].Quantity == 0 || noteInfo.Selection[i].Quantity > orderDetail.Quantity {
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
		movement.Quantity = -orderDetail.Quantity
		movement.SalesDeliveryNote = &deliveryNoteId
		movement.SalesOrderDetail = &orderDetail.Id
		movement.SalesOrder = &saleOrder.Id
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

func locateSalesDeliveryNotesBySalesOrder(orderId int32) []SalesDeliveryNoteLocate {
	var products []SalesDeliveryNoteLocate = make([]SalesDeliveryNoteLocate, 0)
	sqlStatement := `SELECT DISTINCT sales_delivery_note.id,(SELECT name FROM customer WHERE id=sales_delivery_note.customer),sales_delivery_note.date_created,sales_delivery_note.delivery_note_name FROM sales_order_detail INNER JOIN warehouse_movement ON warehouse_movement.sales_order_detail = sales_order_detail.id INNER JOIN sales_delivery_note ON warehouse_movement.sales_delivery_note = sales_delivery_note.id WHERE sales_order_detail."order" = $1`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return products
	}
	for rows.Next() {
		p := SalesDeliveryNoteLocate{}
		rows.Scan(&p.Id, &p.CustomerName, &p.DateCreated, &p.DeliveryNoteName)
		products = append(products, p)
	}

	return products
}

func getNameSalesDeliveryNote(id int32) string {
	sqlStatement := `SELECT delivery_note_name FROM public.sales_delivery_note WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
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

func getSalesDeliveryNoteRelations(noteId int32) SalesDeliveryNoteRelation {
	return SalesDeliveryNoteRelation{
		Orders:    getSalesDeliveryNoteOrders(noteId),
		Shippings: getSalesDeliveryNoteShippings(noteId),
	}
}

func getSalesDeliveryNoteOrders(noteId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	sqlStatement := `SELECT DISTINCT sales_order.* FROM sales_delivery_note INNER JOIN warehouse_movement ON sales_delivery_note.id=warehouse_movement.sales_delivery_note INNER JOIN sales_order ON sales_order.id=warehouse_movement.sales_order WHERE sales_delivery_note.id=$1 ORDER BY sales_order.date_created DESC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		return sales
	}
	for rows.Next() {
		s := SaleOrder{}
		rows.Scan(&s.Id, &s.Warehouse, &s.Reference, &s.Customer, &s.DateCreated, &s.DatePaymetAccepted, &s.PaymentMethod, &s.BillingSeries, &s.Currency, &s.CurrencyChange,
			&s.BillingAddress, &s.ShippingAddress, &s.LinesNumber, &s.InvoicedLines, &s.DeliveryNoteLines, &s.TotalProducts, &s.DiscountPercent, &s.FixDiscount, &s.ShippingPrice, &s.ShippingDiscount,
			&s.TotalWithDiscount, &s.VatAmount, &s.TotalAmount, &s.Description, &s.Notes, &s.Off, &s.Cancelled, &s.Status, &s.OrderNumber, &s.BillingStatus, &s.OrderName, &s.Carrier)
		sales = append(sales, s)
	}

	return sales
}

func getSalesDeliveryNoteShippings(noteId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE shipping.delivery_note=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		return shippings
	}
	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.CustomerName, &s.SaleOrderName, &s.CarrierName)
		shippings = append(shippings, s)
	}

	return shippings
}
