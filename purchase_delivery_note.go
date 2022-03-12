package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PurchaseDeliveryNote struct {
	Id                 int64     `json:"id"`
	Warehouse          string    `json:"warehouse"`
	Supplier           int32     `json:"supplier"`
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
	TotalVat           float64   `json:"totalVat"`
	TotalAmount        float64   `json:"totalAmount"`
	LinesNumber        int16     `json:"linesNumber"`
	DeliveryNoteNumber int32     `json:"deliveryNoteNumber"`
	DeliveryNoteName   string    `json:"deliveryNoteName"`
	Currency           int32     `json:"currency"`
	CurrencyChange     float64   `json:"currencyChange"`
	SupplierName       string    `json:"supplierName"`
	enterprise         int32
}

type PurchaseDeliveryNotes struct {
	Rows   int32                      `json:"rows"`
	Notes  []PurchaseDeliveryNote     `json:"notes"`
	Footer PurchaseDeliveryNoteFooter `json:"footer"`
}

type PurchaseDeliveryNoteFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseDeliveryNotes(enterpriseId int32) PurchaseDeliveryNotes {
	dn := PurchaseDeliveryNotes{}
	dn.Notes = make([]PurchaseDeliveryNote, 0)
	sqlStatement := `SELECT *,(SELECT name FROM suppliers WHERE suppliers.id=purchase_delivery_note.supplier) FROM public.purchase_delivery_note WHERE enterprise=$1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return dn
	}
	defer rows.Close()

	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise, &p.SupplierName)
		dn.Notes = append(dn.Notes, p)
	}

	sqlStatement = `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_delivery_note WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return dn
	}
	row.Scan(&dn.Rows, &dn.Footer.TotalProducts, &dn.Footer.TotalAmount)

	return dn
}

func (s *OrderSearch) searchPurchaseDeliveryNote() PurchaseDeliveryNotes {
	dn := PurchaseDeliveryNotes{}
	dn.Notes = make([]PurchaseDeliveryNote, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		sqlStatement := `SELECT purchase_delivery_note.*,(SELECT name FROM suppliers WHERE suppliers.id=purchase_delivery_note.supplier) FROM purchase_delivery_note WHERE delivery_note_number=$1 AND enterprise=$2 ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, orderNumber, s.enterprise)
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
		sqlStatement += ` AND purchase_delivery_note.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		sqlStatement += ` ORDER BY date_created DESC`
		rows, err = db.Query(sqlStatement, interfaces...)
	}
	if err != nil {
		log("DB", err.Error())
		return dn
	}
	defer rows.Close()

	for rows.Next() {
		p := PurchaseDeliveryNote{}
		rows.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise, &p.SupplierName)
		dn.Notes = append(dn.Notes, p)
	}

	var row *sql.Row
	if err == nil {
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_delivery_note WHERE delivery_note_number=$1 AND enterprise=$2`
		row = db.QueryRow(sqlStatement, orderNumber, s.enterprise)
	} else {
		var interfaces []interface{} = make([]interface{}, 0)
		interfaces = append(interfaces, "%"+s.Search+"%")
		sqlStatement := `SELECT COUNT(*),SUM(total_products),SUM(total_amount) FROM purchase_delivery_note INNER JOIN suppliers ON suppliers.id=purchase_delivery_note.supplier WHERE suppliers.name ILIKE $1`
		if s.DateStart != nil {
			sqlStatement += ` AND purchase_delivery_note.date_created >= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateStart)
		}
		if s.DateEnd != nil {
			sqlStatement += ` AND purchase_delivery_note.date_created <= $` + strconv.Itoa(len(interfaces)+1)
			interfaces = append(interfaces, s.DateEnd)
		}
		sqlStatement += ` AND purchase_delivery_note.enterprise = $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.enterprise)
		row = db.QueryRow(sqlStatement, interfaces...)
	}
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return dn
	}
	row.Scan(&dn.Rows, &dn.Footer.TotalProducts, &dn.Footer.TotalAmount)

	return dn
}

func getPurchaseDeliveryNoteRow(deliveryNoteId int64) PurchaseDeliveryNote {
	sqlStatement := `SELECT * FROM public.purchase_delivery_note WHERE id=$1`
	row := db.QueryRow(sqlStatement, deliveryNoteId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseDeliveryNote{}
	}

	p := PurchaseDeliveryNote{}
	row.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise)

	return p
}

func getPurchaseDeliveryNoteRowTransaction(deliveryNoteId int64, trans sql.Tx) PurchaseDeliveryNote {
	sqlStatement := `SELECT * FROM public.purchase_delivery_note WHERE id=$1`
	row := trans.QueryRow(sqlStatement, deliveryNoteId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseDeliveryNote{}
	}

	p := PurchaseDeliveryNote{}
	row.Scan(&p.Id, &p.Warehouse, &p.Supplier, &p.DateCreated, &p.PaymentMethod, &p.BillingSeries, &p.ShippingAddress, &p.TotalProducts, &p.DiscountPercent, &p.FixDiscount, &p.ShippingPrice, &p.ShippingDiscount, &p.TotalWithDiscount, &p.TotalVat, &p.TotalAmount, &p.LinesNumber, &p.DeliveryNoteName, &p.DeliveryNoteNumber, &p.Currency, &p.CurrencyChange, &p.enterprise)

	return p
}

func (n *PurchaseDeliveryNote) isValid() bool {
	return !(len(n.Warehouse) == 0 || len(n.Warehouse) > 2 || n.Supplier <= 0 || n.PaymentMethod <= 0 || len(n.BillingSeries) == 0 || len(n.BillingSeries) > 3 || n.ShippingAddress <= 0)
}

func (n *PurchaseDeliveryNote) insertPurchaseDeliveryNotes(userId int32, trans *sql.Tx) (bool, int64) {
	if !n.isValid() {
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

	n.DeliveryNoteNumber = getNextPurchaseDeliveryNoteNumber(n.BillingSeries, n.enterprise)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.Currency)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeries + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)

	sqlStatement := `INSERT INTO public.purchase_delivery_note(warehouse, supplier, payment_method, billing_series, shipping_address, discount_percent, fix_discount, shipping_price, shipping_discount, delivery_note_number, delivery_note_name, currency, currency_change, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	row := trans.QueryRow(sqlStatement, n.Warehouse, n.Supplier, n.PaymentMethod, n.BillingSeries, n.ShippingAddress, n.DiscountPercent, n.FixDiscount, n.ShippingPrice, n.ShippingDiscount, n.DeliveryNoteNumber, n.DeliveryNoteName, n.Currency, n.CurrencyChange, n.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false, 0
	}

	var noteId int64
	row.Scan(&noteId)
	n.Id = noteId

	if noteId > 0 {
		insertTransactionalLog(n.enterprise, "purchase_delivery_note", int(noteId), userId, "I")
		json, _ := json.Marshal(n)
		go fireWebHook(n.enterprise, "purchase_delivery_note", "POST", string(json))
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false, 0
		}
		///
	}

	return noteId > 0, noteId
}

func (n *PurchaseDeliveryNote) deletePurchaseDeliveryNotes(userId int32, trans *sql.Tx) bool {
	if n.Id <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false
		}
		///
	}

	inMemoryNote := getPurchaseDeliveryNoteRow(n.Id)
	if inMemoryNote.enterprise != n.enterprise {
		return false
	}

	d := getWarehouseMovementByPurchaseDeliveryNote(n.Id, n.enterprise)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	insertTransactionalLog(n.enterprise, "purchase_delivery_note", int(n.Id), userId, "D")
	json, _ := json.Marshal(n)
	go fireWebHook(n.enterprise, "purchase_delivery_note", "DELETE", string(json))

	sqlStatement := `DELETE FROM public.purchase_delivery_note WHERE id=$1 AND enterprise=$2`
	res, err := trans.Exec(sqlStatement, n.Id, n.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false
		}
		///
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. There are no details to generate the delivery note
func deliveryNoteAllPurchaseOrder(purchaseOrderId int64, enterpriseId int32, userId int32) (OkAndErrorCodeReturn, int64) {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	if purchaseOrder.enterprise != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if purchaseOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if purchaseOrder.DeliveryNoteLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}, 0
	}
	orderDetails := getPurchaseOrderDetail(purchaseOrderId, purchaseOrder.enterprise)
	filterPurchaseOrderDetails(orderDetails, func(pod PurchaseOrderDetail) bool { return pod.QuantityDeliveryNote < pod.Quantity })
	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
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
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	///

	n.enterprise = enterpriseId
	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
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
		movement.enterprise = enterpriseId
		ok = movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
	}

	///
	transErr = trans.Commit()
	return OkAndErrorCodeReturn{Ok: transErr == nil}, deliveryNoteId
	///
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail has a delivery note generated
// 4. The selected quantity is greater than the quantity pending of delivery note generation in the detail
func (noteInfo *OrderDetailGenerate) deliveryNotePartiallyPurchaseOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(noteInfo.OrderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.enterprise != enterpriseId || len(noteInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.DeliveryNoteLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(noteInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != noteInfo.OrderId || noteInfo.Selection[i].Quantity == 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if noteInfo.Selection[i].Quantity > orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2, ExtraData: []string{product.Name}}
		}
		if orderDetail.QuantityDeliveryNote >= orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{product.Name}}
		}
		if (noteInfo.Selection[i].Quantity + orderDetail.QuantityDeliveryNote) > orderDetail.Quantity {
			product := getProductRow(orderDetail.Product)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 4, ExtraData: []string{product.Name}}
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
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	n.enterprise = enterpriseId
	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(purchaseOrderDetails); i++ {
		orderDetail := purchaseOrderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "I"
		movement.Warehouse = purchaseOrder.Warehouse
		movement.Product = orderDetail.Product
		movement.Quantity = noteInfo.Selection[i].Quantity
		movement.PurchaseDeliveryNote = &deliveryNoteId
		movement.PurchaseOrderDetail = &orderDetail.Id
		movement.PurchaseOrder = &purchaseOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.enterprise = enterpriseId
		ok = movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	transErr = trans.Commit()
	return OkAndErrorCodeReturn{Ok: transErr == nil}
	///
}

type PurchaseDeliveryNoteRelation struct {
	Orders []PurchaseOrder `json:"orders"`
}

func getPurchaseDeliveryNoteRelations(noteId int32, enterpriseId int32) PurchaseDeliveryNoteRelation {
	return PurchaseDeliveryNoteRelation{
		Orders: getPurchaseDeliveryNoteOrders(noteId, enterpriseId),
	}
}

func getPurchaseDeliveryNoteOrders(noteId int32, enterpriseId int32) []PurchaseOrder {
	var orders []PurchaseOrder = make([]PurchaseOrder, 0)
	sqlStatement := `SELECT DISTINCT purchase_order.* FROM purchase_delivery_note INNER JOIN warehouse_movement ON purchase_delivery_note.id=warehouse_movement.purchase_delivery_note INNER JOIN purchase_order ON purchase_order.id=warehouse_movement.purchase_order WHERE purchase_delivery_note.id=$1 AND purchase_delivery_note.enterprise=$2 ORDER BY purchase_order.date_created DESC`
	rows, err := db.Query(sqlStatement, noteId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	defer rows.Close()

	for rows.Next() {
		o := PurchaseOrder{}
		rows.Scan(&o.Id, &o.Warehouse, &o.SupplierReference, &o.Supplier, &o.DateCreated, &o.DatePaid, &o.PaymentMethod, &o.BillingSeries, &o.Currency, &o.CurrencyChange,
			&o.BillingAddress, &o.ShippingAddress, &o.LinesNumber, &o.InvoicedLines, &o.DeliveryNoteLines, &o.TotalProducts, &o.DiscountPercent, &o.FixDiscount, &o.ShippingPrice, &o.ShippingDiscount,
			&o.TotalWithDiscount, &o.VatAmount, &o.TotalAmount, &o.Description, &o.Notes, &o.Off, &o.Cancelled, &o.OrderNumber, &o.BillingStatus, &o.OrderName, &o.enterprise)
		orders = append(orders, o)
	}

	return orders
}

// Adds a total amount to the delivery note total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseDeliveryNote(noteId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_delivery_note SET total_products=total_products+$2,total_vat=total_vat+$3 WHERE id=$1`
	_, err := trans.Exec(sqlStatement, noteId, totalAmount, (totalAmount/100)*vatPercent)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseDeliveryNote(noteId, enterpriseId, userId, trans)
}

// Applies the logic to calculate the totals of the delivery note.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseDeliveryNote(noteId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_delivery_note SET total_with_discount=(total_products-total_products*(discount_percent/100))-fix_discount+shipping_price-shipping_discount,total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err := trans.Exec(sqlStatement, noteId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE purchase_delivery_note SET total_amount=total_with_discount+total_vat WHERE id=$1`
	_, err = trans.Exec(sqlStatement, noteId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_delivery_note", int(noteId), userId, "U")
	n := getPurchaseDeliveryNoteRowTransaction(noteId, trans)
	json, _ := json.Marshal(n)
	go fireWebHook(n.enterprise, "purchase_delivery_note", "PUT", string(json))

	return err == nil
}
