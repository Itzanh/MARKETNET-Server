package main

import "fmt"

type PurchaseOrderDetail struct {
	Id                       int32   `json:"id"`
	Order                    int32   `json:"order"`
	Product                  int32   `json:"product"`
	Price                    float32 `json:"price"`
	Quantity                 int32   `json:"quantity"`
	VatPercent               float32 `json:"vatPercent"`
	TotalAmount              float32 `json:"totalAmount"`
	QuantityInvoiced         int32   `json:"quantityInvoiced"`
	QuantityDeliveryNote     int32   `json:"quantityDeliveryNote"`
	QuantityPendingPackaging int32   `json:"quantityPendingPackaging"`
	QuantityAssignedSale     int32   `json:"quantityAssignedSale"`
}

func getPurchaseOrderDetail(orderId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE "order"=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale)
		details = append(details, d)
	}

	return details
}

func getPurchaseOrderDetailRow(detailId int32) PurchaseOrderDetail {
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		return PurchaseOrderDetail{}
	}

	d := PurchaseOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale)

	return d
}

func (d *PurchaseOrderDetail) isValid() bool {
	return !(d.Order <= 0 || d.Product <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (s *PurchaseOrderDetail) insertPurchaseOrderDetail() (bool, int32) {
	if !s.isValid() {
		return false, 0
	}

	s.TotalAmount = (s.Price * float32(s.Quantity)) * (1 + (s.VatPercent / 100))

	///
	trans, err := db.Begin()
	if err != nil {
		return false, 0
	}
	///

	sqlStatement := `INSERT INTO public.purchase_order_detail("order", product, price, quantity, vat_percent, total_amount, quantity_pending_packaging) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Quantity)
	if row.Err() != nil {
		trans.Rollback()
		return false, 0
	}
	var detailId int32
	row.Scan(&detailId)

	ok := addTotalProductsPurchaseOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
	if !ok {
		trans.Rollback()
		return false, 0
	}

	quantityAssignedSale := associatePurchaseOrderWithPendingSalesOrders(detailId, s.Product, s.Quantity)
	sqlStatement = `UPDATE purchase_order_detail SET quantity_assigned_sale=$2 WHERE id=$1`
	_, err = db.Exec(sqlStatement, detailId, quantityAssignedSale)
	if err != nil {
		trans.Rollback()
		return false, 0
	}

	// add quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row = db.QueryRow(sqlStatement, s.Order)
	if row.Err() != nil {
		fmt.Println(row.Err())
		trans.Rollback()
		return false, 0
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(s.Product, warehouse, s.Quantity)
	if !ok {
		trans.Rollback()
		return false, 0
	}

	///
	err = trans.Commit()
	return err == nil, detailId
	///
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func associatePurchaseOrderWithPendingSalesOrders(purchaseDetailId int32, productId int32, quantity int32) int32 {
	// associate pending sales order detail until there are no more quantity pending to be assigned, or there are no more pending sales order details
	sqlStatement := `SELECT id,quantity,"order" FROM sales_order_detail WHERE product=$1 AND status='A' ORDER BY (SELECT date_created FROM sales_order WHERE sales_order.id=sales_order_detail."order") ASC`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return 0
	}

	var quantityAssignedSale int32
	for quantityAssignedSale < quantity {
		if rows.Next() {
			var salesDetailId int32
			var salesQuantity int32
			var orderId int32
			rows.Scan(&salesDetailId, &salesQuantity, &orderId)

			if quantityAssignedSale+salesQuantity > quantity { // no more rows to proecss
				return quantityAssignedSale
			}

			sqlStatement := `UPDATE sales_order_detail SET status='B',purchase_order_detail=$2 WHERE id=$1`
			_, err := db.Exec(sqlStatement, salesDetailId, purchaseDetailId)
			if err == nil {
				quantityAssignedSale += salesQuantity
				setSalesOrderState(orderId)
			}
		} else { // no more rows to proecss
			return quantityAssignedSale
		}
	}
	return quantityAssignedSale
}

func (s *PurchaseOrderDetail) updatePurchaseOrderDetail() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	order := getPurchaseOrderRow(s.Order)
	if order.Id <= 0 || order.InvoicedLines != 0 {
		return false
	}
	inMemoryDetail := getPurchaseOrderDetailRow(s.Id)
	if inMemoryDetail.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET product=$2,price=$3,quantity=$4,vat_percent=$5 WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id, s.Product, s.Price, s.Quantity, s.VatPercent)
	if err != nil {
		return false
	}

	// take out the old value
	ok := addTotalProductsPurchaseOrder(inMemoryDetail.Order, -(inMemoryDetail.Price * float32(inMemoryDetail.Quantity)), inMemoryDetail.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	// add the new value
	ok = addTotalProductsPurchaseOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}

	// update quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, s.Order)
	if row.Err() != nil {
		trans.Rollback()
		return false
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(s.Product, warehouse, -inMemoryDetail.Quantity)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = addQuantityPendingReveiving(s.Product, warehouse, s.Quantity)
	if !ok {
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

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a transaction.
func (s *PurchaseOrderDetail) deletePurchaseOrderDetail() bool {
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
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return false
	}

	// roll back the state of the sale order details
	details := getSalesOrderDetailPurchaseOrderPending(s.Id)
	for i := 0; i < len(details); i++ {
		sqlStatement := `UPDATE sales_order_detail SET status='A',purchase_order_detail=NULL WHERE id=$1`
		_, err = db.Exec(sqlStatement, details[i].Id)
		if err != nil {
			trans.Rollback()
			return false
		}
		ok := setSalesOrderState(details[i].Order)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.purchase_order_detail WHERE id=$1`
	_, err = db.Exec(sqlStatement, s.Id)
	if err != nil {
		trans.Rollback()
		return false
	}
	ok := addTotalProductsPurchaseOrder(detailInMemory.Order, -(detailInMemory.Price * float32(detailInMemory.Quantity)), detailInMemory.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}

	// substract quantity pending receiving
	sqlStatement = `SELECT warehouse FROM purchase_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailInMemory.Order)
	if row.Err() != nil {
		trans.Rollback()
		return false
	}
	var warehouse string
	row.Scan(&warehouse)
	ok = addQuantityPendingReveiving(detailInMemory.Product, warehouse, -detailInMemory.Quantity)
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
func addQuantityAssignedSalePurchaseOrder(detailId int32, quantity int32) bool {
	sqlStatement := `UPDATE purchase_order_detail SET quantity_assigned_sale=quantity_assigned_sale+$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()
	if err != nil || rows == 0 {
		return false
	}
	return true
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNotePurchaseOrderDetail(detailId int32, quantity int32) bool {

	detailBefore := getPurchaseOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE purchase_order_detail SET quantity_delivery_note=quantity_delivery_note+$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	if err != nil {
		fmt.Println(err)
	}
	rows, _ := res.RowsAffected()
	if err != nil && rows == 0 {
		return false
	}
	return true
}
