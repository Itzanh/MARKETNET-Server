package main

type SalesOrderDetail struct {
	Id                       int32   `json:"id"`
	Order                    int32   `json:"order"`
	Product                  int32   `json:"product"`
	Price                    float32 `json:"price"`
	Quantity                 int32   `json:"quantity"`
	VatPercent               float32 `json:"vatPercent"`
	TotalAmount              float32 `json:"totalAmount"`
	QuantityInvoiced         int32   `json:"quantityInvoiced"`
	QuantityDeliveryNote     int32   `json:"quantityDeliveryNote"`
	Status                   string  `json:"status"`
	QuantityPendingPackaging int32   `json:"quantityPendingPackaging"`
}

func getSalesOrderDetail(orderId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE "order" = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging)
		details = append(details, d)
	}

	return details
}

func getSalesOrderDetailRow(detailId int32) SalesOrderDetail {
	sqlStatement := `SELECT * FROM sales_order_detail WHERE id = $1 ORDER BY id ASC`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		return SalesOrderDetail{}
	}

	d := SalesOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging)

	return d
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

	sqlStatement := `INSERT INTO public.sales_order_detail("order", product, price, quantity, vat_percent, total_amount, status, quantity_pending_packaging) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	res, err := db.Exec(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Status, s.Quantity)
	if err != nil {
		trans.Rollback()
		return false
	}
	ok := addTotalProductsSalesOrder(s.Order, s.Price*float32(s.Quantity), s.VatPercent)
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
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return false
	}
	sqlStatement := `DELETE FROM public.sales_order_detail WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		trans.Rollback()
		return false
	}
	ok := addTotalProductsSalesOrder(detailInMemory.Order, -(detailInMemory.Price * float32(detailInMemory.Quantity)), detailInMemory.VatPercent)
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

type SalesOrderDetailDefaults struct {
	Price      float32 `json:"price"`
	VatPercent float32 `json:"vatPercent"`
}

// Adds an invoiced quantity to the sale order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvociedSalesOrderDetail(detailId int32, quantity int32) bool {

	detailBefore := getSalesOrderDetailRow(detailId)
	if detailBefore.Id <= 0 {
		return false
	}
	salesOrder := getSalesOrderRow(detailBefore.Order)

	sqlStatement := `UPDATE sales_order_detail SET quantity_invoiced = quantity_invoiced + $2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, detailId, quantity)
	rows, _ := res.RowsAffected()
	if err != nil && rows == 0 {
		return false
	}

	detailAfter := getSalesOrderDetailRow(detailId)
	if detailAfter.Id <= 0 {
		return false
	}

	var ok bool
	if detailBefore.QuantityInvoiced != detailBefore.Quantity && detailAfter.QuantityInvoiced == detailAfter.Quantity {
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, detailBefore.Quantity)
		// set the order detail state applying the workflow logic
		if ok {
			sqlStatement = `UPDATE sales_order_detail SET status = $2 WHERE id = $1`
			db.Exec(sqlStatement, detailId, detailBefore.computeStatus())
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity {
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, -detailBefore.Quantity)
		// reset order detail state to "Waiting for Payment"
		if ok {
			sqlStatement = `UPDATE sales_order_detail SET status = '_' WHERE id = $1`
			db.Exec(sqlStatement, detailId)
		}
	}

	if !ok {
		return false
	}

	return err == nil
}

func (s *SalesOrderDetail) computeStatus() string {
	product := getProductRow(s.Product)
	if product.Id <= 0 {
		return ""
	}
	order := getSalesOrderRow(s.Order)
	stock := getStockRow(s.Product, order.Warehouse)
	if stock.Quantity > 0 { // the product is in stock, send to preparation
		return "E"
	} else { // the product is not in stock, purchase or manufacture
		if product.Manufacturing {
			return "C"
		} else {
			// search for pending purchases
			// TODO: ADD PURCHASE ORDER FUNCTIONALITY
			return "A"
		}
	}
}
