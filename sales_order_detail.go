package main

import "time"

type SalesOrderDetail struct {
	Id                   int32   `json:"id"`
	Order                int32   `json:"order"`
	Product              int32   `json:"product"`
	Price                float32 `json:"price"`
	Quantity             int32   `json:"quantity"`
	VatPercent           float32 `json:"vatPercent"`
	TotalAmount          float32 `json:"totalAmount"`
	QuantityInvoiced     int32   `json:"quantityInvoiced"`
	QuantityDeliveryNote int32   `json:"quantityDeliveryNote"`
	// _ = Waiting for payment, A = Waiting for purchase order, B = Purchase order pending, C = Waiting for manufacturing orders, D = Manufacturing orders pending, E = Sent to preparation, F = Awaiting for shipping, G = Shipped, H = Receiced by the customer
	Status                   string `json:"status"`
	QuantityPendingPackaging int32  `json:"quantityPendingPackaging"`
	PurchaseOrderDetail      *int32 `json:"purchaseOrderDetail"`
	PrestaShopId             int32  `json:"prestaShopId"`
	ProductName              string `json:"productName"`
	Cancelled                bool   `json:"cancelled"`
	WooCommerceId            int32  `json:"wooCommerceId"`
	ShopifyId                int64
	ShopifyDraftId           int64
}

func getSalesOrderDetail(orderId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_order_detail.product) FROM sales_order_detail WHERE "order"=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.PrestaShopId, &d.Cancelled, &d.WooCommerceId, &d.ShopifyId, &d.ShopifyDraftId, &d.ProductName)
		details = append(details, d)
	}

	return details
}

func getSalesOrderDetailRow(detailId int32) SalesOrderDetail {
	sqlStatement := `SELECT * FROM sales_order_detail WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesOrderDetail{}
	}

	d := SalesOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.PrestaShopId, &d.Cancelled, &d.WooCommerceId, &d.ShopifyId, &d.ShopifyDraftId)

	return d
}

// Used for purchases
func getSalesOrderDetailWaitingForPurchaseOrder(productId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE product=$1 AND status='A'`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.PrestaShopId, &d.Cancelled, &d.WooCommerceId, &d.ShopifyId, &d.ShopifyDraftId)
		details = append(details, d)
	}

	return details
}

// Used for purchases
func getSalesOrderDetailPurchaseOrderPending(purchaseOrderDetail int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := db.Query(sqlStatement, purchaseOrderDetail)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.PrestaShopId, &d.Cancelled, &d.WooCommerceId, &d.ShopifyId, &d.ShopifyDraftId)
		details = append(details, d)
	}

	return details
}

func (s *SalesOrderDetail) isValid() bool {
	return !(s.Order <= 0 || s.Product <= 0 || s.Quantity <= 0 || s.VatPercent < 0)
}

func (s *SalesOrderDetail) insertSalesOrderDetail() bool {
	if !s.isValid() {
		return false
	}

	s.TotalAmount = (s.Price * float32(s.Quantity)) * (1 + (s.VatPercent / 100))
	s.Status = "_"

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.sales_order_detail("order", product, price, quantity, vat_percent, total_amount, status, quantity_pending_packaging, ps_id, wc_id, sy_draft_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	res, err := db.Exec(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Status, s.Quantity, s.PrestaShopId, s.WooCommerceId, s.ShopifyDraftId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	ok := addTotalProductsSalesOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = setSalesOrderState(s.Order)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = addSalesOrderLinesNumber(s.Order)
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

func (s *SalesOrderDetail) updateSalesOrderDetail() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	inMemoryDetail := getSalesOrderDetailRow(s.Id)
	if inMemoryDetail.Id <= 0 {
		return false
	}
	if inMemoryDetail.QuantityInvoiced > 0 {
		return false
	}

	s.TotalAmount = (s.Price * float32(s.Quantity)) * (1 + (s.VatPercent / 100))
	sqlStatement := `UPDATE sales_order_detail SET product=$2,price=$3,quantity=$4,vat_percent=$5,total_amount=$6,sy_id=$7 WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.ShopifyId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	// take out the old value
	ok := addTotalProductsSalesOrder(inMemoryDetail.Order, -(inMemoryDetail.Price * float32(inMemoryDetail.Quantity)), inMemoryDetail.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	// add the new value
	ok = addTotalProductsSalesOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
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
func (s *SalesOrderDetail) deleteSalesOrderDetail() bool {
	if s.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	detailInMemory := getSalesOrderDetailRow(s.Id)
	if detailInMemory.Id <= 0 || detailInMemory.QuantityInvoiced > 0 || detailInMemory.QuantityDeliveryNote > 0 {
		trans.Rollback()
		return false
	}
	sqlStatement := `DELETE FROM public.sales_order_detail WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	ok := addTotalProductsSalesOrder(detailInMemory.Order, -(detailInMemory.Price * float32(detailInMemory.Quantity)), detailInMemory.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = setSalesOrderState(detailInMemory.Order)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = removeSalesOrderLinesNumber(detailInMemory.Order)
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

// Adds an invoiced quantity to the sale order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvociedSalesOrderDetail(detailId int32, quantity int32) bool {
	detailBefore := getSalesOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}
	salesOrder := getSalesOrderRow(detailBefore.Order)

	sqlStatement := `UPDATE sales_order_detail SET quantity_invoiced=quantity_invoiced+$2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()
	if err != nil && rows == 0 {
		log("DB", err.Error())
		return false
	}

	detailAfter := getSalesOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	var ok bool
	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity { // set as invoced
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, detailBefore.Quantity)
		// set the order detail state applying the workflow logic
		if ok {
			status, purchaseOrderDetail := detailBefore.computeStatus()
			sqlStatement = `UPDATE sales_order_detail SET status=$2,purchase_order_detail=$3 WHERE id=$1`
			_, err := db.Exec(sqlStatement, detailId, status, purchaseOrderDetail)
			if err != nil {
				log("DB", err.Error())
				return false
			}
		}
		if !ok {
			return false
		}
		ok = addSalesOrderInvoicedLines(detailBefore.Order)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, -detailBefore.Quantity)
		// reset order detail state to "Waiting for Payment"
		if ok {
			sqlStatement = `UPDATE sales_order_detail SET status='_',purchase_order_detail=NULL WHERE id=$1`
			_, err := db.Exec(sqlStatement, detailId)
			if err != nil {
				log("DB", err.Error())
				return false
			}
		}
		if !ok {
			return false
		}
		ok = removeSalesOrderInvoicedLines(detailBefore.Order)
		if !ok {
			return false
		}
	}

	ok = setSalesOrderState(salesOrder.Id)
	if !ok {
		return false
	}

	return err == nil
}

func (s *SalesOrderDetail) computeStatus() (string, *int32) {
	product := getProductRow(s.Product)
	if product.Id <= 0 {
		return "", nil
	}
	order := getSalesOrderRow(s.Order)
	stock := getStockRow(s.Product, order.Warehouse)
	if stock.Quantity > 0 { // the product is in stock, send to preparation
		return "E", nil
	} else { // the product is not in stock, purchase or manufacture
		if product.Manufacturing {
			return "C", nil
		} else {
			// search for pending purchases
			sqlStatement := `SELECT id FROM purchase_order_detail WHERE product=$1 AND quantity_delivery_note = 0 AND quantity - quantity_assigned_sale >= $2 ORDER BY (SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order") ASC LIMIT 1`
			row := db.QueryRow(sqlStatement, s.Product, s.Quantity)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				return "A", nil
			}
			var purchaseDetailId int32
			row.Scan(&purchaseDetailId)
			if purchaseDetailId <= 0 {
				return "A", nil
			}

			// add quantity assigned to sale orders
			ok := addQuantityAssignedSalePurchaseOrder(purchaseDetailId, s.Quantity)
			if !ok {
				return "A", nil
			}

			// set the purchase order detail
			return "B", &purchaseDetailId
		}
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityPendingPackagingSaleOrderDetail(detailId int32, quantity int32) bool {
	sqlStatement := `UPDATE sales_order_detail SET quantity_pending_packaging = quantity_pending_packaging + $2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
	}

	ok := rows > 0 && err == nil
	if !ok {
		return false
	}

	detail := getSalesOrderDetailRow(detailId)
	var status string
	if detail.QuantityPendingPackaging <= 0 {
		status = "F"
	} else {
		status = "E"
	}
	sqlStatement = `UPDATE sales_order_detail SET status=$2 WHERE id=$1`
	res, err = db.Exec(sqlStatement, detailId, status)
	rows, _ = res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
		return false
	}

	ok = rows > 0 && err == nil
	if !ok {
		return false
	}

	return setSalesOrderState(detail.Order)
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNoteSalesOrderDetail(detailId int32, quantity int32) bool {

	detailBefore := getSalesOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE sales_order_detail SET quantity_delivery_note = quantity_delivery_note + $2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()

	if err != nil {
		log("DB", err.Error())
		return false
	}

	detailAfter := getSalesOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	var ok bool
	if detailBefore.QuantityDeliveryNote != detailBefore.Quantity && detailAfter.QuantityDeliveryNote == detailAfter.Quantity { // set as delivery note generated
		ok = addSalesOrderDeliveryNoteLines(detailBefore.Order)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok = removeSalesOrderDeliveryNoteLines(detailBefore.Order)
		if !ok {
			return false
		}
	}

	if err != nil && rows == 0 {
		return false
	}
	return true
}

func cancelSalesOrderDetail(detailId int32) bool {
	detail := getSalesOrderDetailRow(detailId)
	if detail.Id <= 0 {
		return false
	}

	if !detail.Cancelled {
		if detail.Quantity <= 0 || detail.QuantityInvoiced < 0 || detail.QuantityDeliveryNote > 0 {
			return false
		}

		sqlStatement := `UPDATE public.sales_order_detail SET quantity_invoiced=quantity, quantity_delivery_note=quantity, status='Z', cancelled=true WHERE id=$1`
		_, err := db.Exec(sqlStatement, detailId)

		if err != nil {
			log("DB", err.Error())
			return false
		}

		ok := setSalesOrderState(detail.Order)
		if !ok {
			return false
		}

		return err == nil
	} else {
		if detail.Quantity <= 0 || detail.QuantityInvoiced == 0 || detail.QuantityDeliveryNote == 0 {
			return false
		}

		sqlStatement := `UPDATE public.sales_order_detail SET quantity_invoiced=0, quantity_delivery_note=0, cancelled=false WHERE id=$1`
		_, err := db.Exec(sqlStatement, detailId)

		if err != nil {
			log("DB", err.Error())
		}

		status, purchaseOrderDetail := detail.computeStatus()
		sqlStatement = `UPDATE sales_order_detail SET status=$2,purchase_order_detail=$3 WHERE id=$1`
		_, err = db.Exec(sqlStatement, detailId, status, purchaseOrderDetail)
		if err != nil {
			log("DB", err.Error())
		}

		ok := setSalesOrderState(detail.Order)
		if !ok {
			return false
		}

		return err == nil
	}
}

type SalePurchasesOrderDetail struct {
	Id           int32     `json:"id"`
	Order        int32     `json:"order"`
	OrderName    string    `json:"orderName"`
	DateCreated  time.Time `json:"dateCreated"`
	SupplierName string    `json:"supplierName"`
	Quantity     int32     `json:"quantity"`
	TotalAmount  float32   `json:"totalAmount"`
}

func getPurchasesOrderDetailsFromSaleOrderDetail(detailId int32) []SalePurchasesOrderDetail {
	salePurchasesOrderDetail := make([]SalePurchasesOrderDetail, 0)
	sqlStatement := `SELECT purchase_order_detail.id,"order",(SELECT order_name FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order"),(SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order"),(SELECT name FROM suppliers WHERE suppliers.id=(SELECT supplier FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order")),quantity,total_amount FROM purchase_order_detail WHERE purchase_order_detail.id=(SELECT purchase_order_detail FROM sales_order_detail WHERE id=$1)`
	rows, err := db.Query(sqlStatement, detailId)
	if err != nil {
		log("DB", err.Error())
		return salePurchasesOrderDetail
	}

	for rows.Next() {
		p := SalePurchasesOrderDetail{}
		rows.Scan(&p.Id, &p.Order, &p.OrderName, &p.DateCreated, &p.SupplierName, &p.Quantity, &p.TotalAmount)
		salePurchasesOrderDetail = append(salePurchasesOrderDetail, p)
	}

	return salePurchasesOrderDetail
}
