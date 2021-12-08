package main

import (
	"database/sql"
	"time"
)

type PurchaseOrderDetail struct {
	Id                       int64   `json:"id"`
	Order                    int64   `json:"order"`
	Product                  int32   `json:"product"`
	Price                    float64 `json:"price"`
	Quantity                 int32   `json:"quantity"`
	VatPercent               float64 `json:"vatPercent"`
	TotalAmount              float64 `json:"totalAmount"`
	QuantityInvoiced         int32   `json:"quantityInvoiced"`
	QuantityDeliveryNote     int32   `json:"quantityDeliveryNote"`
	QuantityPendingPackaging int32   `json:"quantityPendingPackaging"`
	QuantityAssignedSale     int32   `json:"quantityAssignedSale"`
	ProductName              string  `json:"productName"`
	enterprise               int32
}

func getPurchaseOrderDetail(orderId int64, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=purchase_order_detail.product) FROM purchase_order_detail WHERE "order"=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.ProductName)
		details = append(details, d)
	}

	return details
}

func getPurchaseOrderDetailRow(detailId int64) PurchaseOrderDetail {
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseOrderDetail{}
	}

	d := PurchaseOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise)

	return d
}

func (d *PurchaseOrderDetail) isValid() bool {
	return !(d.Order <= 0 || d.Product <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (s *PurchaseOrderDetail) insertPurchaseOrderDetail(beginTrans bool, userId int32) (bool, int64) {
	if !s.isValid() {
		return false, 0
	}

	p := getProductRow(s.Product)
	if p.Id <= 0 || p.Off {
		return false, 0
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	///
	var trans *sql.Tx
	if beginTrans {
		var err error
		trans, err = db.Begin()
		if err != nil {
			return false, 0
		}
	}
	///

	sqlStatement := `INSERT INTO public.purchase_order_detail("order", product, price, quantity, vat_percent, total_amount, quantity_pending_packaging, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Quantity, s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}

	var detailId int64
	row.Scan(&detailId)
	s.Id = detailId

	insertTransactionalLog(s.enterprise, "purchase_order_detail", int(detailId), userId, "I")

	if detailId <= 0 {
		return false, 0
	}

	ok := addTotalProductsPurchaseOrder(s.Order, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId)
	if !ok {
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}
	ok = addPurchaseOrderLinesNumber(s.Order)
	if !ok {
		trans.Rollback()
		return false, 0
	}

	quantityAssignedSale := associatePurchaseOrderWithPendingSalesOrders(detailId, s.Product, s.Quantity, s.enterprise, userId)
	if quantityAssignedSale < 0 {
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}

	sqlStatement = `UPDATE purchase_order_detail SET quantity_assigned_sale=$2 WHERE id=$1`
	_, err := db.Exec(sqlStatement, detailId, quantityAssignedSale)
	if err != nil {
		log("DB", err.Error())
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}

	insertTransactionalLog(s.enterprise, "purchase_order_detail", int(detailId), userId, "U")

	// add quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row = db.QueryRow(sqlStatement, s.Order)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(s.Product, warehouse, s.Quantity, s.enterprise)
	if !ok {
		if beginTrans {
			trans.Rollback()
		}
		return false, 0
	}

	if beginTrans {
		///
		err = trans.Commit()
		return err == nil, detailId
		///
	} else {
		return true, detailId
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func associatePurchaseOrderWithPendingSalesOrders(purchaseDetailId int64, productId int32, quantity int32, enterpriseId int32, userId int32) int32 {
	// associate pending sales order detail until there are no more quantity pending to be assigned, or there are no more pending sales order details
	sqlStatement := `SELECT id,quantity,"order" FROM sales_order_detail WHERE product=$1 AND status='A' ORDER BY (SELECT date_created FROM sales_order WHERE sales_order.id=sales_order_detail."order") ASC`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		log("DB", err.Error())
		return -1
	}

	var quantityAssignedSale int32
	for quantityAssignedSale < quantity {
		if rows.Next() {
			var salesDetailId int32
			var salesQuantity int32
			var orderId int64
			rows.Scan(&salesDetailId, &salesQuantity, &orderId)

			if quantityAssignedSale+salesQuantity > quantity { // no more rows to proecss
				return quantityAssignedSale
			}

			sqlStatement := `UPDATE sales_order_detail SET status='B',purchase_order_detail=$2 WHERE id=$1`
			_, err := db.Exec(sqlStatement, salesDetailId, purchaseDetailId)
			if err == nil {
				quantityAssignedSale += salesQuantity
				setSalesOrderState(enterpriseId, orderId, userId)
				insertTransactionalLog(enterpriseId, "sales_order_detail", int(salesDetailId), userId, "U")
			} else {
				log("DB", err.Error())
			}
		} else { // no more rows to proecss
			return quantityAssignedSale
		}
	}
	return quantityAssignedSale
}

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a transaction.
func (s *PurchaseOrderDetail) deletePurchaseOrderDetail(userId int32) bool {
	if s.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	detailInMemory := getPurchaseOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.enterprise != s.enterprise || detailInMemory.QuantityInvoiced > 0 || detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return false
	}

	// roll back the state of the sale order details
	details := getSalesOrderDetailPurchaseOrderPending(s.Id)
	for i := 0; i < len(details); i++ {
		sqlStatement := `UPDATE sales_order_detail SET status='A',purchase_order_detail=NULL WHERE id=$1`
		_, err = db.Exec(sqlStatement, details[i].Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
		ok := setSalesOrderState(detailInMemory.enterprise, details[i].Order, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	insertTransactionalLog(detailInMemory.enterprise, "purchase_order_detail", int(s.Id), userId, "D")

	sqlStatement := `DELETE FROM public.purchase_order_detail WHERE id=$1`
	_, err = db.Exec(sqlStatement, s.Id)
	if err != nil {
		trans.Rollback()
		return false
	}

	ok := addTotalProductsPurchaseOrder(detailInMemory.Order, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, detailInMemory.enterprise, userId)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = removePurchaseOrderLinesNumber(detailInMemory.Order)
	if !ok {
		trans.Rollback()
		return false
	}

	// substract quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailInMemory.Order)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(detailInMemory.Product, warehouse, -detailInMemory.Quantity, s.enterprise)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}

// Adds quantity to the field to prevent from other sale orders to use the quantity that is already reserved for order that are already waiting a purchase order.
// This function will substract if a negative quantity is given.
// THIS FUNCION DOES NOT OPEN A TRANSACTION
func addQuantityAssignedSalePurchaseOrder(detailId int64, quantity int32, enterpriseId int32, userId int32) bool {
	sqlStatement := `UPDATE purchase_order_detail SET quantity_assigned_sale=quantity_assigned_sale+$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()
	if err != nil || rows == 0 {
		if err != nil {
			log("DB", err.Error())
		}

		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")

	return true
}

// Adds an invoiced quantity to the purchase order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvoicedPurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET quantity_invoiced=quantity_invoiced+$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")

	detailAfter := getPurchaseOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity { // set as invoced
		ok := addPurchaseOrderInvoicedLines(detailBefore.Order, enterpriseId, userId)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok := removePurchaseOrderInvoicedLines(detailBefore.Order, enterpriseId, userId)
		if !ok {
			return false
		}
	}

	return err == nil && rows > 0
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNotePurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET quantity_delivery_note=quantity_delivery_note+$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()
	if err != nil && rows == 0 {
		log("DB", err.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")

	detailAfter := getPurchaseOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	if detailBefore.QuantityDeliveryNote != detailBefore.Quantity && detailAfter.QuantityDeliveryNote == detailAfter.Quantity { // set as delivery note generated
		ok := addPurchaseOrderDeliveryNoteLines(detailBefore.Order, enterpriseId, userId)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok := removePurchaseOrderDeliveryNoteLines(detailBefore.Order, enterpriseId, userId)
		if !ok {
			return false
		}
	}

	if quantity > 0 { // the purchase order has been added to a delivery note, advance the status from the pending sales order details
		return setSalesOrderDetailStateAllPendingPurchaseOrder(detailId, enterpriseId, userId)
	} else { // the delivery note details has been removed, roll back the sales order detail status
		return undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId, enterpriseId, userId)
	}
}

// All the purchase order detail has been added to a delivery note. Advance the status from all the pending sales details to "Sent to preparation".
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setSalesOrderDetailStateAllPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32) bool {
	sqlStatement := `SELECT "order" FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := db.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE sales_order_detail SET status='E' WHERE purchase_order_detail=$1 AND status='B'`
	_, err = db.Exec(sqlStatement, detailId)

	for rows.Next() {
		var orderId int64
		rows.Scan(&orderId)
		setSalesOrderState(enterpriseId, orderId, userId)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailId), userId, "U")
	}

	return err == nil
}

// The purchase order detail was added to a delivery note and it advanced the status from the sales details, but the delivery note was deleted.
// Roll back the status change.
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32) bool {
	sqlStatement := `SELECT "order" FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := db.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `UPDATE sales_order_detail SET status='B' WHERE purchase_order_detail=$1 AND status='E'`
	_, err = db.Exec(sqlStatement, detailId)

	for rows.Next() {
		var orderId int64
		rows.Scan(&orderId)
		setSalesOrderState(enterpriseId, orderId, userId)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailId), userId, "U")
	}

	return err == nil
}

type PurchaseSalesOrderDetail struct {
	Id           int32     `json:"id"`
	Order        int32     `json:"order"`
	OrderName    string    `json:"orderName"`
	DateCreated  time.Time `json:"dateCreated"`
	CustomerName string    `json:"customerName"`
	Quantity     int32     `json:"quantity"`
	TotalAmount  float64   `json:"totalAmount"`
}

func getSalesOrderDetailsFromPurchaseOrderDetail(detailId int64, enterpriseId int32) []PurchaseSalesOrderDetail {
	purchaseSalesOrderDetail := make([]PurchaseSalesOrderDetail, 0)
	sqlStatement := `SELECT sales_order_detail.id,"order",(SELECT order_name FROM sales_order WHERE sales_order.id=sales_order_detail."order"),(SELECT date_created FROM sales_order WHERE sales_order.id=sales_order_detail."order"),(SELECT name FROM customer WHERE customer.id=(SELECT customer FROM sales_order WHERE sales_order.id=sales_order_detail."order")),quantity,total_amount FROM sales_order_detail WHERE purchase_order_detail=$1 AND enterprise=$2`
	rows, err := db.Query(sqlStatement, detailId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return purchaseSalesOrderDetail
	}

	for rows.Next() {
		p := PurchaseSalesOrderDetail{}
		rows.Scan(&p.Id, &p.Order, &p.OrderName, &p.DateCreated, &p.CustomerName, &p.Quantity, &p.TotalAmount)
		purchaseSalesOrderDetail = append(purchaseSalesOrderDetail, p)
	}

	return purchaseSalesOrderDetail
}
