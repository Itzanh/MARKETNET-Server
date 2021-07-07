package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type PurchaseDeliveryNote struct {
	Id                 int32     `json:"id"`
	Warehouse          string    `json:"warehouse"`
	Supplier           int32     `json:"supplier"`
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
	SupplierName       string    `json:"supplierName"`
}

func getPurchaseDeliveryNotes() []PurchaseDeliveryNote {
	var notes []PurchaseDeliveryNote = make([]PurchaseDeliveryNote, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_delivery_note.supplier) FROM public.purchase_delivery_note ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return notes
	}
	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.SupplierName)
		notes = append(notes, p)
	}

	return notes
}

func (s *OrderSearch) searchPurchaseDeliveryNote() []PurchaseDeliveryNote {
	var notes []PurchaseDeliveryNote = make([]PurchaseDeliveryNote, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_delivery_note.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_delivery_note.supplier) FROM purchase_delivery_note WHERE delivery_note_number=$1 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT purchase_delivery_note.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_delivery_note.supplier) FROM purchase_delivery_note INNER JOIN suppliers ON suppliers.id=purchase_delivery_note.supplier WHERE suppliers.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_delivery_note.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_delivery_note.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		fmt.Println(err)
		return notes
	}
	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.SupplierName)
		notes = append(notes, p)
	}

	return notes
}

func getPurchaseDeliveryNoteRow(deliveryNoteId int32) PurchaseDeliveryNote {
	sqlStatement := `SELECT * FROM public.purchase_delivery_note WHERE id=$1`
	row := db.QueryRow(sqlStatement, deliveryNoteId)
	if row.Err() != nil {
		return PurchaseDeliveryNote{}
	}

	p := PurchaseDeliveryNote{}
	row.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange)

	return p
}

func (n *PurchaseDeliveryNote) isValid() bool {
	return !(len(n.Warehouse) == 0 || len(n.Warehouse) > 2 || n.Supplier <= 0 || n.PaymentMethod <= 0 || len(n.BillingSeries) == 0 || len(n.BillingSeries) > 3 || n.ShippingAddress <= 0)
}

func (n *PurchaseDeliveryNote) insertPurchaseDeliveryNotes() (bool, int32) {
	if !n.isValid() {
		return false, 0
	}

	n.DeliveryNoteNumber = getNextPurchaseDeliveryNoteNumber(n.BillingSeries)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.Currency)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)

	sqlStatement := `INSERT INTO public.purchase_delivery_note(warehouse, supplier, payment_method, billing_series, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, delivery_note_number, delivery_note_name, currency, currency_change) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`
	row := db.QueryRow(sqlStatement, n.Warehouse, n.Supplier, n.PaymentMethod, n.BillingSeries, n.ShippingAddress, n.DiscountPercent, n.FixDiscount, n.ShippingPrice, n.ShippingDiscount, n.DeliveryNoteNumber, n.DeliveryNoteName, n.Currency, n.CurrencyChange)
	if row.Err() != nil {
		return false, 0
	}

	var noteId int32
	row.Scan(&noteId)
	return noteId > 0, noteId
}

func (n *PurchaseDeliveryNote) deletePurchaseDeliveryNotes() bool {
	if n.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	d := getWarehouseMovementByPurchaseDeliveryNote(n.Id)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteWarehouseMovement()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.purchase_delivery_note WHERE id=$1`
	res, err := db.Exec(sqlStatement, n.Id)
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

func deliveryNoteAllPurchaseOrder(purchaseOrderId int32) (bool, int32) {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	orderDetails := getPurchaseOrderDetail(purchaseOrderId)

	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return false, 0
	}

	// create a delivery note for that order
	n := PurchaseDeliveryNote{}
	n.Supplier = purchaseOrder.Supplier
	n.ShippingAddress = purchaseOrder.ShippingAddress
	n.Currency = purchaseOrder.Currency
	n.PaymentMethod = purchaseOrder.PaymentMethod
	n.BillingSeries = purchaseOrder.BillingSeries
	n.Warehouse = purchaseOrder.Warehouse

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false, 0
	}
	///

	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes()
	if !ok {
		trans.Rollback()
		return false, 0
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "I"
		movement.Warehouse = purchaseOrder.Warehouse
		movement.Product = orderDetail.Product
		movement.Quantity = orderDetail.Quantity
		movement.PurchaseDeliveryNote = &deliveryNoteId
		movement.PurchaseOrderDetail = &orderDetail.Id
		movement.PurchaseOrder = &purchaseOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
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

func (noteInfo *OrderDetailGenerate) deliveryNotePartiallyPurchaseOrder() bool {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(noteInfo.OrderId)
	if purchaseOrder.Id <= 0 || len(noteInfo.Selection) == 0 {
		return false
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(noteInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != noteInfo.OrderId || noteInfo.Selection[i].Quantity == 0 || noteInfo.Selection[i].Quantity > orderDetail.Quantity {
			return false
		}
		purchaseOrderDetails = append(purchaseOrderDetails, orderDetail)
	}

	// create a delivery note for that order
	n := PurchaseDeliveryNote{}
	n.Supplier = purchaseOrder.Supplier
	n.ShippingAddress = purchaseOrder.ShippingAddress
	n.Currency = purchaseOrder.Currency
	n.PaymentMethod = purchaseOrder.PaymentMethod
	n.BillingSeries = purchaseOrder.BillingSeries
	n.Warehouse = purchaseOrder.Warehouse

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes()
	if !ok {
		trans.Rollback()
		return false
	}
	for i := 0; i < len(purchaseOrderDetails); i++ {
		orderDetail := purchaseOrderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "I"
		movement.Warehouse = purchaseOrder.Warehouse
		movement.Product = orderDetail.Product
		movement.Quantity = orderDetail.Quantity
		movement.PurchaseDeliveryNote = &deliveryNoteId
		movement.PurchaseOrderDetail = &orderDetail.Id
		movement.PurchaseOrder = &purchaseOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
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

type PurchaseDeliveryNoteRelation struct {
	Orders []PurchaseOrder `json:"orders"`
}

func getPurchaseDeliveryNoteRelations(noteId int32) PurchaseDeliveryNoteRelation {
	return PurchaseDeliveryNoteRelation{
		Orders: getPurchaseDeliveryNoteOrders(noteId),
	}
}

func getPurchaseDeliveryNoteOrders(noteId int32) []PurchaseOrder {
	var orders []PurchaseOrder = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT DISTINCT purchase_order.* FROM purchase_delivery_note INNER JOIN warehouse_movement ON purchase_delivery_note.id=warehouse_movement.purchase_delivery_note INNER JOIN purchase_order ON purchase_order.id=warehouse_movement.purchase_order WHERE purchase_delivery_note.id=$1 ORDER BY purchase_order.date_created DESC`
	rows, err := db.Query(sqlStatement, noteId)
	if err != nil {
		return orders
	}
	for rows.Next() {
		o := PurchaseOrder{}
		rows.Scan(&o.Id, &o.Warehouse, &o.SupplierReference, &o.Supplier, &o.DateCreated, &o.DatePaid, &o.PaymentMethod, &o.BillingSeries, &o.Currency, &o.CurrencyChange,
			&o.BillingAddress, &o.ShippingAddress, &o.LinesNumber, &o.InvoicedLines, &o.DeliveryNoteLines, &o.TotalProducts, &o.DiscountPercent, &o.FixDiscount, &o.ShippingPrice, &o.ShippingDiscount,
			&o.TotalWithDiscount, &o.VatAmount, &o.TotalAmount, &o.Description, &o.Notes, &o.Off, &o.Cancelled, &o.OrderNumber, &o.BillingStatus, &o.OrderName)
		orders = append(orders, o)
	}

	return orders
}

// Adds a total amount to the delivery note total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseDeliveryNote(noteId int32, totalAmount float32, vatPercent float32) bool {
	sqlStatement := `UPDATE purchase_delivery_note SET total_products=total_products+$2,total_vat=total_vat+$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, noteId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return calcTotalsPurchaseDeliveryNote(noteId)
}

// Applies the logic to calculate the totals of the delivery note.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseDeliveryNote(noteId int32) bool {
	sqlStatement := `UPDATE purchase_delivery_note SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err := db.Exec(sqlStatement, noteId)
	if err != nil {
		return false
	}

	sqlStatement = `UPDATE purchase_delivery_note SET total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err = db.Exec(sqlStatement, noteId)
	return err == nil
}
