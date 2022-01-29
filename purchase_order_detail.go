package main

import (
	"database/sql"
	"encoding/json"
	"sort"
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
	Cancelled                bool    `json:"cancelled"`
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
	defer rows.Close()

	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.Cancelled, &d.ProductName)
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
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.Cancelled)

	return d
}

func getPurchaseOrderDetailRowTransaction(detailId int64, trans sql.Tx) PurchaseOrderDetail {
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE id=$1`
	row := trans.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseOrderDetail{}
	}

	d := PurchaseOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.Cancelled)

	return d
}

func (d *PurchaseOrderDetail) isValid() bool {
	return !(d.Order <= 0 || d.Product <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

// 1. the product is deactivated
// 2. there is aleady a detail with this product
func (s *PurchaseOrderDetail) insertPurchaseOrderDetail(userId int32, trans *sql.Tx) (OkAndErrorCodeReturn, int64) {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	product := getProductRow(s.Product)
	if product.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if product.Off {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}, 0
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	sqlStatement := `SELECT COUNT(purchase_order_detail) FROM public.purchase_order_detail WHERE "order" = $1 AND product = $2`
	row := db.QueryRow(sqlStatement, s.Order, s.Product)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	var countProductInSaleOrder int16
	row.Scan(&countProductInSaleOrder)
	if countProductInSaleOrder > 0 {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}, 0
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	///
	var beginTrans bool = (trans == nil)
	if beginTrans {
		var err error
		trans, err = db.Begin()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
	}
	///

	sqlStatement = `INSERT INTO public.purchase_order_detail("order", product, price, quantity, vat_percent, total_amount, quantity_pending_packaging, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	row = trans.QueryRow(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Quantity, s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	var detailId int64
	row.Scan(&detailId)
	s.Id = detailId

	if detailId <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	insertTransactionalLog(s.enterprise, "purchase_order_detail", int(detailId), userId, "I")
	jsn, _ := json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "POST", string(jsn))

	ok := addTotalProductsPurchaseOrder(s.Order, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	ok = addPurchaseOrderLinesNumber(s.Order, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	quantityAssignedSale := associatePurchaseOrderWithPendingSalesOrders(detailId, s.Product, s.Quantity, s.enterprise, userId, *trans)
	if quantityAssignedSale < 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	sqlStatement = `UPDATE purchase_order_detail SET quantity_assigned_sale=$2 WHERE id=$1`
	_, err := trans.Exec(sqlStatement, detailId, quantityAssignedSale)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	insertTransactionalLog(s.enterprise, "purchase_order_detail", int(detailId), userId, "U")
	jsn, _ = json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "POST", string(jsn))

	// add quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row = db.QueryRow(sqlStatement, s.Order)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		if beginTrans {
			trans.Rollback()
		}
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(s.Product, warehouse, s.Quantity, s.enterprise, *trans)
	if !ok {
		if beginTrans {
			trans.Rollback()
		}
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	if beginTrans {
		///
		err = trans.Commit()
		return OkAndErrorCodeReturn{Ok: err == nil}, detailId
		///
	} else {
		return OkAndErrorCodeReturn{Ok: true}, detailId
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func associatePurchaseOrderWithPendingSalesOrders(purchaseDetailId int64, productId int32, quantity int32, enterpriseId int32, userId int32, trans sql.Tx) int32 {
	// associate pending sales order detail until there are no more quantity pending to be assigned, or there are no more pending sales order details
	sqlStatement := `SELECT id,quantity,"order" FROM sales_order_detail WHERE product=$1 AND status='A' ORDER BY (SELECT date_created FROM sales_order WHERE sales_order.id=sales_order_detail."order") ASC`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		log("DB", err.Error())
		return -1
	}
	defer rows.Close()

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
			_, err := trans.Exec(sqlStatement, salesDetailId, purchaseDetailId)
			if err == nil {
				quantityAssignedSale += salesQuantity
				setSalesOrderState(enterpriseId, orderId, userId, trans)
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

// ERROR CODES:
// 1. the detail is already invoiced
// 2. the detail has a delivery note generated
func (s *PurchaseOrderDetail) updatePurchaseOrderDetail(userId int32) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	detailInMemory := getPurchaseOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.enterprise != s.enterprise {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.QuantityInvoiced > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}
	if detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	}

	if detailInMemory.Quantity < s.Quantity { // increased quantity
		associatePurchaseOrderWithPendingSalesOrders(s.Id, s.Product, s.Quantity-detailInMemory.Quantity, s.enterprise, userId, *trans)
		s.QuantityAssignedSale = detailInMemory.QuantityAssignedSale
	} else if detailInMemory.Quantity > s.Quantity { // decreased quantity
		salesDetails := getSalesOrderDetailsFromPurchaseOrderDetail(s.Id, s.enterprise)
		sort.Slice(salesDetails, func(i, j int) bool { // sort by date_created ASC
			return salesDetails[i].DateCreated.Before(salesDetails[j].DateCreated)
		})

		sqlStatement := `UPDATE public.sales_order_detail SET status='A', purchase_order_detail=NULL WHERE id=$1`
		quantityAssignedSale := detailInMemory.QuantityAssignedSale

		for i := 0; i < len(salesDetails); i++ {
			_, err := trans.Exec(sqlStatement, salesDetails[i].Id)
			if err != nil {
				log("DB", err.Error())
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}
			quantityAssignedSale -= salesDetails[i].Quantity
			if quantityAssignedSale <= s.Quantity {
				break
			}
		}
		s.QuantityAssignedSale = quantityAssignedSale
	} else {
		s.QuantityAssignedSale = detailInMemory.QuantityAssignedSale
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	ok := addTotalProductsPurchaseOrder(s.Order, -detailInMemory.Price*float64(detailInMemory.Quantity), detailInMemory.VatPercent, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok = addTotalProductsPurchaseOrder(s.Order, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	sqlStatement := `UPDATE public.purchase_order_detail SET price = $2, quantity = $3, vat_percent = $4, total_amount = $5, quantity_assigned_sale = $6 WHERE id = $1`
	_, err = trans.Exec(sqlStatement, s.Id, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.QuantityAssignedSale)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	err = trans.Commit()
	if err != nil {
		log("DB", err.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a transaction.
//
// ERROR CODES:
// 1. the detail is already invoiced
// 2. the detail has a delivery note generated
func (s *PurchaseOrderDetail) deletePurchaseOrderDetail(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		var err error
		trans, err = db.Begin()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	detailInMemory := getPurchaseOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.enterprise != s.enterprise {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.QuantityInvoiced > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}
	if detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	}

	// roll back the state of the sale order details
	details := getSalesOrderDetailPurchaseOrderPending(s.Id)
	for i := 0; i < len(details); i++ {
		sqlStatement := `UPDATE sales_order_detail SET status='A',purchase_order_detail=NULL WHERE id=$1`
		_, err := trans.Exec(sqlStatement, details[i].Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
		ok := setSalesOrderState(detailInMemory.enterprise, details[i].Order, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(detailInMemory.enterprise, "purchase_order_detail", int(s.Id), userId, "D")
	json, _ := json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "DELETE", string(json))

	sqlStatement := `DELETE FROM public.purchase_order_detail WHERE id=$1`
	_, err := trans.Exec(sqlStatement, s.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsPurchaseOrder(detailInMemory.Order, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, detailInMemory.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	ok = removePurchaseOrderLinesNumber(detailInMemory.Order, detailInMemory.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	// substract quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailInMemory.Order)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(detailInMemory.Product, warehouse, -detailInMemory.Quantity, s.enterprise, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	err = trans.Commit()
	return OkAndErrorCodeReturn{Ok: err == nil}
	///
}

// Adds quantity to the field to prevent from other sale orders to use the quantity that is already reserved for order that are already waiting a purchase order.
// This function will substract if a negative quantity is given.
// THIS FUNCION DOES NOT OPEN A TRANSACTION
func addQuantityAssignedSalePurchaseOrder(detailId int64, quantity int32, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE purchase_order_detail SET quantity_assigned_sale=quantity_assigned_sale+$2 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, detailId, quantity)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	s := getPurchaseOrderDetailRowTransaction(detailId, trans)
	json, _ := json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "PUT", string(json))

	return true
}

// Adds an invoiced quantity to the purchase order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvoicedPurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32, trans sql.Tx) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET quantity_invoiced=quantity_invoiced+$2 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, detailId, quantity)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	s := getPurchaseOrderDetailRowTransaction(detailId, trans)
	json, _ := json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "PUT", string(json))

	detailAfter := getPurchaseOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity { // set as invoced
		ok := addPurchaseOrderInvoicedLines(detailBefore.Order, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok := removePurchaseOrderInvoicedLines(detailBefore.Order, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	}

	return err == nil && rows > 0
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNotePurchaseOrderDetail(detailId int64, quantity int32, enterpriseId int32, userId int32, trans sql.Tx) bool {
	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET quantity_delivery_note=quantity_delivery_note+$2 WHERE id=$1`
	res, err := trans.Exec(sqlStatement, detailId, quantity)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order_detail", int(detailId), userId, "U")
	s := getPurchaseOrderDetailRowTransaction(detailId, trans)
	json, _ := json.Marshal(s)
	go fireWebHook(s.enterprise, "purchase_order_detail", "PUT", string(json))

	detailAfter := getPurchaseOrderDetailRowTransaction(detailId, trans)
	if detailAfter.Id <= 0 {
		return false
	}

	if detailBefore.QuantityDeliveryNote != detailBefore.Quantity && detailAfter.QuantityDeliveryNote == detailAfter.Quantity { // set as delivery note generated
		ok := addPurchaseOrderDeliveryNoteLines(detailBefore.Order, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok := removePurchaseOrderDeliveryNoteLines(detailBefore.Order, enterpriseId, userId, trans)
		if !ok {
			return false
		}
	}

	if quantity > 0 { // the purchase order has been added to a delivery note, advance the status from the pending sales order details
		return setSalesOrderDetailStateAllPendingPurchaseOrder(detailId, enterpriseId, userId, trans) && setComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId, enterpriseId, userId, trans)
	} else { // the delivery note details has been removed, roll back the sales order detail status
		return undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId, enterpriseId, userId, trans) && undoComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId, enterpriseId, userId, trans)
	}
}

func cancelPurchaseOrderDetail(detailId int64, enterpriseId int32, userId int32) bool {
	detail := getPurchaseOrderDetailRow(detailId)
	if detail.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	if !detail.Cancelled {
		if detail.Quantity <= 0 || detail.QuantityInvoiced < 0 || detail.QuantityDeliveryNote > 0 {
			return false
		}

		sqlStatement := `UPDATE public.purchase_order_detail SET quantity_invoiced=quantity, quantity_delivery_note=quantity, cancelled=true WHERE id=$1 AND enterprise=$2`
		_, err := trans.Exec(sqlStatement, detailId, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		if err != nil {
			insertTransactionalLog(detail.enterprise, "purchase_order_detail", int(detailId), userId, "U")
			s := getSalesOrderDetailRow(detailId)
			json, _ := json.Marshal(s)
			go fireWebHook(s.enterprise, "purchase_order_detail", "PUT", string(json))
		}

		///
		err = trans.Commit()
		if err != nil {
			return false
		}
		///

		return err == nil
	} else {
		if detail.Quantity <= 0 || detail.QuantityInvoiced == 0 || detail.QuantityDeliveryNote == 0 {
			return false
		}

		sqlStatement := `UPDATE public.purchase_order_detail SET quantity_invoiced=0, quantity_delivery_note=0, cancelled=false, quantity_assigned_sale=0 WHERE id=$1 AND enterprise=$2`
		_, err := trans.Exec(sqlStatement, detailId, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		salesDetails := getSalesOrderDetailsFromPurchaseOrderDetail(detail.Id, detail.enterprise)

		sqlStatement = `UPDATE public.sales_order_detail SET status='A', purchase_order_detail=NULL WHERE id=$1`
		for i := 0; i < len(salesDetails); i++ {
			_, err := trans.Exec(sqlStatement, salesDetails[i].Id)
			if err != nil {
				log("DB", err.Error())
				trans.Rollback()
				return false
			}
		}

		if err != nil {
			insertTransactionalLog(detail.enterprise, "purchase_order_detail", int(detailId), userId, "U")
			s := getSalesOrderDetailRow(detailId)
			json, _ := json.Marshal(s)
			go fireWebHook(s.enterprise, "purchase_order_detail", "PUT", string(json))
		}

		///
		err = trans.Commit()
		if err != nil {
			return false
		}
		///

		return err == nil
	}
}

// All the purchase order detail has been added to a delivery note. Advance the status from all the pending sales details to "Sent to preparation".
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setSalesOrderDetailStateAllPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `SELECT "order" FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := db.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	defer rows.Close()

	sqlStatement = `UPDATE sales_order_detail SET status='E' WHERE purchase_order_detail=$1 AND status='B'`
	_, err = trans.Exec(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	for rows.Next() {
		var orderId int64
		rows.Scan(&orderId)
		setSalesOrderState(enterpriseId, orderId, userId, trans)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailId), userId, "U")
		s := getPurchaseOrderDetailRowTransaction(detailId, trans)
		json, _ := json.Marshal(s)
		go fireWebHook(s.enterprise, "sales_order_detail", "PUT", string(json))
	}

	return err == nil
}

// The purchase order detail was added to a delivery note and it advanced the status from the sales details, but the delivery note was deleted.
// Roll back the status change.
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func undoSalesOrderDetailStatueFromPendingPurchaseOrder(detailId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `SELECT "order" FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := trans.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	var orderIds []int64 = make([]int64, 0)

	for rows.Next() {
		var orderId int64
		rows.Scan(&orderId)
		orderIds = append(orderIds, orderId)
	}
	rows.Close()

	sqlStatement = `UPDATE sales_order_detail SET status='B' WHERE purchase_order_detail=$1 AND status='E'`
	_, err = trans.Exec(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	for i := 0; i < len(orderIds); i++ {
		setSalesOrderState(enterpriseId, orderIds[i], userId, trans)
		insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailId), userId, "U")
		s := getPurchaseOrderDetailRowTransaction(detailId, trans)
		json, _ := json.Marshal(s)
		go fireWebHook(s.enterprise, "sales_order_detail", "PUT", string(json))
	}

	return err == nil
}

// Gets all the sub orders from complex manufacturing orders that are waiting for this pending purchase order,
// and sets them to manufactured (and the parent order as manufacturable).
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `SELECT DISTINCT complex_manufacturing_order FROM public.complex_manufacturing_order_manufacturing_order WHERE purchase_order_detail = $1`
	rows, err := trans.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	defer rows.Close()

	var complexManufacturingOrders []int64 = make([]int64, 0)

	for rows.Next() {
		var complexManufacturingOrder int64
		rows.Scan(&complexManufacturingOrder)
	}

	for i := 0; i < len(complexManufacturingOrders); i++ {
		ok := setComplexManufacturingOrderManufacturingOrderManufactured(complexManufacturingOrders[i], true, enterpriseId, userId, &trans)
		if !ok {
			return false
		}
	}

	return true
}

// The purchase order detail was added to a delivery note and it advanced the status from the complex manufacturing orders, but the delivery note was deleted.
// Roll back the status change.
//
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func undoComplexManufacturingOrdersPendingPurchaseOrderManufactured(detailId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `SELECT DISTINCT complex_manufacturing_order FROM public.complex_manufacturing_order_manufacturing_order WHERE purchase_order_detail = $1`
	rows, err := db.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	defer rows.Close()

	var complexManufacturingOrders []int64 = make([]int64, 0)

	for rows.Next() {
		var complexManufacturingOrder int64
		rows.Scan(&complexManufacturingOrder)
	}

	for i := 0; i < len(complexManufacturingOrders); i++ {
		ok := setComplexManufacturingOrderManufacturingOrderManufactured(complexManufacturingOrders[i], false, enterpriseId, userId, &trans)
		if !ok {
			return false
		}
	}

	return true
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
	defer rows.Close()

	for rows.Next() {
		p := PurchaseSalesOrderDetail{}
		rows.Scan(&p.Id, &p.Order, &p.OrderName, &p.DateCreated, &p.CustomerName, &p.Quantity, &p.TotalAmount)
		purchaseSalesOrderDetail = append(purchaseSalesOrderDetail, p)
	}

	return purchaseSalesOrderDetail
}

type PurchaseComplexManufacturingOrder struct {
	Id           int64     `json:"id"`
	Type         int32     `json:"type"`
	Manufactured bool      `json:"manufactured"`
	DateCreated  time.Time `json:"dateCreated"`
	TypeName     string    `json:"typeName"`
}

func getComplexManufacturingOrdersFromPurchaseOrderDetail(detailId int64, enterpriseId int32) []PurchaseComplexManufacturingOrder {
	purchaseComplexManufacturingOrder := make([]PurchaseComplexManufacturingOrder, 0)
	sqlStatement := `SELECT DISTINCT complex_manufacturing_order.id,complex_manufacturing_order.type,complex_manufacturing_order.manufactured,complex_manufacturing_order.date_created,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=complex_manufacturing_order.type) FROM public.complex_manufacturing_order INNER JOIN complex_manufacturing_order_manufacturing_order ON complex_manufacturing_order_manufacturing_order.complex_manufacturing_order=complex_manufacturing_order.id WHERE complex_manufacturing_order_manufacturing_order.purchase_order_detail = $1 AND complex_manufacturing_order.enterprise = $2 ORDER BY complex_manufacturing_order.date_created ASC`
	rows, err := db.Query(sqlStatement, detailId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return purchaseComplexManufacturingOrder
	}
	defer rows.Close()

	for rows.Next() {
		o := PurchaseComplexManufacturingOrder{}
		rows.Scan(&o.Id, &o.Type, &o.Manufactured, &o.DateCreated, &o.TypeName)
		purchaseComplexManufacturingOrder = append(purchaseComplexManufacturingOrder, o)
	}
	return purchaseComplexManufacturingOrder
}

func filterPurchaseOrderDetails(input []PurchaseOrderDetail, test func(PurchaseOrderDetail) bool) (output []PurchaseOrderDetail) {
	for _, s := range input {
		if test(s) {
			output = append(output, s)
		}
	}
	return
}
