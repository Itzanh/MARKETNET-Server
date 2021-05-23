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

	n.DeliveryNoteNumber = getNextSaleInvoiceNumber(n.BillingSeries)
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

func deliveryNoteAllSaleOrder(saleOrderId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	orderDetails := getSalesOrderDetail(saleOrderId)

	if saleOrder.Id <= 0 || len(orderDetails) == 0 {
		return false
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
			return false
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
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
