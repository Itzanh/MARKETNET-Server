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
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging)
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
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging)

	return d
}

func (d *PurchaseOrderDetail) isValid() bool {
	return !(d.Order <= 0 || d.Product <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (s *PurchaseOrderDetail) insertPurchaseOrderDetail() bool {
	fmt.Println("insertPurchaseOrderDetail")
	if !s.isValid() {
		fmt.Println("INVALID")
		return false
	}

	s.TotalAmount = (s.Price * float32(s.Quantity)) * (1 + (s.VatPercent / 100))

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.purchase_order_detail("order", product, price, quantity, vat_percent, total_amount, quantity_pending_packaging) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Quantity)
	if err != nil {
		fmt.Println(err)
		trans.Rollback()
		return false
	}
	ok := addTotalProductsPurchaseOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
	if !ok {
		fmt.Println("ERROR 1")
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
	sqlStatement := `DELETE FROM public.purchase_order_detail WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		trans.Rollback()
		return false
	}
	ok := addTotalProductsPurchaseOrder(detailInMemory.Order, -(detailInMemory.Price * float32(detailInMemory.Quantity)), detailInMemory.VatPercent)
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
