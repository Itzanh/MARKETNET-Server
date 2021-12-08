package main

import (
	"time"
)

type SalesOrderDetail struct {
	Id                       int64   `json:"id"`
	Order                    int64   `json:"order"`
	Product                  int32   `json:"product"`
	Price                    float64 `json:"price"`
	Quantity                 int32   `json:"quantity"`
	VatPercent               float64 `json:"vatPercent"`
	TotalAmount              float64 `json:"totalAmount"`
	QuantityInvoiced         int32   `json:"quantityInvoiced"`
	QuantityDeliveryNote     int32   `json:"quantityDeliveryNote"`
	Status                   string  `json:"status"` // _ = Waiting for payment, A = Waiting for purchase order, B = Purchase order pending, C = Waiting for manufacturing orders, D = Manufacturing orders pending, E = Sent to preparation, F = Awaiting for shipping, G = Shipped, H = Receiced by the customer
	QuantityPendingPackaging int32   `json:"quantityPendingPackaging"`
	PurchaseOrderDetail      *int64  `json:"purchaseOrderDetail"`
	ProductName              string  `json:"productName"`
	Cancelled                bool    `json:"cancelled"`
	DigitalProduct           bool    `json:"digitalProduct"`
	prestaShopId             int32
	wooCommerceId            int32
	shopifyId                int64
	shopifyDraftId           int64
	enterprise               int32
}

func getSalesOrderDetail(orderId int64, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_order_detail.product),(SELECT digital_product FROM product WHERE product.id=sales_order_detail.product) FROM sales_order_detail WHERE "order"=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise, &d.ProductName, &d.DigitalProduct)
		details = append(details, d)
	}

	return details
}

func getSalesOrderDetailRow(detailId int64) SalesOrderDetail {
	sqlStatement := `SELECT * FROM sales_order_detail WHERE id=$1`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesOrderDetail{}
	}

	d := SalesOrderDetail{}
	row.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise)

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
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise)
		details = append(details, d)
	}

	return details
}

// Used for purchases
func getSalesOrderDetailPurchaseOrderPending(purchaseOrderDetail int64) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE purchase_order_detail=$1 AND status='B'`
	rows, err := db.Query(sqlStatement, purchaseOrderDetail)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise)
		details = append(details, d)
	}

	return details
}

func (s *SalesOrderDetail) isValid() bool {
	return !(s.Order <= 0 || s.Product <= 0 || s.Quantity <= 0 || s.VatPercent < 0)
}

func (s *SalesOrderDetail) insertSalesOrderDetail(userId int32) bool {
	if !s.isValid() {
		return false
	}

	p := getProductRow(s.Product)
	if p.Id <= 0 || p.Off {
		return false
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))
	s.Status = "_"

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.sales_order_detail("order", product, price, quantity, vat_percent, total_amount, status, quantity_pending_packaging, ps_id, wc_id, sy_draft_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Order, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.Status, s.Quantity, s.prestaShopId, s.wooCommerceId, s.shopifyDraftId, s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false
	}

	var detail int64
	row.Scan(&detail)
	s.Id = detail

	ok := addTotalProductsSalesOrder(s.enterprise, s.Order, userId, s.Price*float64(s.Quantity), s.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = setSalesOrderState(s.enterprise, s.Order, userId)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = addSalesOrderLinesNumber(s.enterprise, s.Order, userId)
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

	if detail > 0 {
		insertTransactionalLog(s.enterprise, "sales_order_detail", int(detail), userId, "I")
	}

	return detail > 0
}

func (s *SalesOrderDetail) updateSalesOrderDetail(userId int32) bool {
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

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))
	sqlStatement := `UPDATE sales_order_detail SET product=$2,price=$3,quantity=$4,vat_percent=$5,total_amount=$6,sy_id=$7 WHERE id=$1 AND enterprise=$8`
	res, err := db.Exec(sqlStatement, s.Id, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.shopifyId, s.enterprise)
	rows, _ := res.RowsAffected()
	if err != nil {
		log("DB", err.Error())
		return false
	}

	if rows == 0 {
		///
		err = trans.Commit()
		if err != nil {
			return false
		}
		///

		return rows > 0
	}

	// take out the old value
	ok := addTotalProductsSalesOrder(s.enterprise, inMemoryDetail.Order, userId, -(inMemoryDetail.Price * float64(inMemoryDetail.Quantity)), inMemoryDetail.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	// add the new value
	ok = addTotalProductsSalesOrder(s.enterprise, s.Order, userId, s.Price*float64(s.Quantity), s.VatPercent)
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

	if rows > 0 {
		insertTransactionalLog(s.enterprise, "sales_order_detail", int(s.Id), userId, "U")
	}

	return rows > 0
}

// Deletes an order detail, substracting the stock and the amount from the order total. All the operations are done under a transaction.
func (s *SalesOrderDetail) deleteSalesOrderDetail(userId int32) bool {
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

	insertTransactionalLog(s.enterprise, "sales_order_detail", int(s.Id), userId, "D")

	sqlStatement := `DELETE FROM public.sales_order_detail WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, s.Id, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		///
		err = trans.Rollback()
		if err != nil {
			return false
		}
		///

		return false
	}

	ok := addTotalProductsSalesOrder(s.enterprise, detailInMemory.Order, userId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = setSalesOrderState(detailInMemory.enterprise, detailInMemory.Order, userId)
	if !ok {
		trans.Rollback()
		return false
	}
	ok = removeSalesOrderLinesNumber(detailInMemory.enterprise, detailInMemory.Order, userId)
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

	return rows > 0
}

// Adds an invoiced quantity to the sale order detail. This function will subsctract from the quantity if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityInvociedSalesOrderDetail(detailId int64, quantity int32, userId int32) bool {
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
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, detailBefore.Quantity, detailBefore.enterprise)
		// set the order detail state applying the workflow logic
		if ok {
			status, purchaseOrderDetail := detailBefore.computeStatus(userId)
			sqlStatement := `UPDATE sales_order_detail SET status=$2,purchase_order_detail=$3 WHERE id=$1`
			_, err := db.Exec(sqlStatement, detailId, status, purchaseOrderDetail)
			if err != nil {
				log("DB", err.Error())
				return false
			}

		}
		if !ok {
			return false
		}
		ok = addSalesOrderInvoicedLines(salesOrder.enterprise, detailBefore.Order, userId)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityInvoiced == detailBefore.Quantity && detailAfter.QuantityInvoiced != detailAfter.Quantity { // undo invoiced
		ok = addQuantityPendingServing(detailBefore.Product, salesOrder.Warehouse, -detailBefore.Quantity, detailBefore.enterprise)
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
		ok = removeSalesOrderInvoicedLines(salesOrder.enterprise, detailBefore.Order, userId)
		if !ok {
			return false
		}

		// reset relations
		if detailBefore.PurchaseOrderDetail != nil {
			ok := addQuantityAssignedSalePurchaseOrder(*detailBefore.PurchaseOrderDetail, detailBefore.Quantity, detailBefore.enterprise, userId)
			if !ok {
				return false
			}
		}
		orders := getSalesOrderManufacturingOrders(salesOrder.Id, salesOrder.enterprise)
		for i := 0; i < len(orders); i++ {
			if orders[i].OrderDetail != nil || *orders[i].OrderDetail != detailBefore.Id {
				continue
			}

			ok := orders[i].deleteManufacturingOrder(userId)
			if !ok {
				return false
			}
		}
		sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET sale_order_detail = NULL WHERE sale_order_detail = $1`
		_, err := db.Exec(sqlStatement, detailBefore.Id)
		if err != nil {
			log("DB", err.Error())
			return false
		}
		// -- reset relations
	}

	ok = setSalesOrderState(salesOrder.enterprise, salesOrder.Id, userId)
	if !ok {
		return false
	}

	if err == nil {
		insertTransactionalLog(detailBefore.enterprise, "sales_order_detail", int(detailId), userId, "U")
	}

	return err == nil
}

// returns: status, purchase order detail id
func (s *SalesOrderDetail) computeStatus(userId int32) (string, *int64) {
	product := getProductRow(s.Product)
	if product.Id <= 0 {
		return "", nil
	}

	order := getSalesOrderRow(s.Order)
	stock := getStockRow(s.Product, order.Warehouse, s.enterprise)
	if !product.ControlStock {
		return "E", nil
	} else if stock.Quantity > 0 { // the product is in stock, send to preparation
		return "E", nil
	} else { // the product is not in stock, purchase or manufacture
		if product.Manufacturing {
			// search for pending manufacturing order for stock
			manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderType)
			if manufacturingOrderType.Complex {
				sqlStatement := `SELECT id, manufacturing_order_type_component FROM public.complex_manufacturing_order_manufacturing_order WHERE product = $1 AND type = 'O' AND manufactured = false AND sale_order_detail IS NULL ORDER BY id ASC`
				rows, err := db.Query(sqlStatement, product.Id)
				if err != nil {
					log("DB", err.Error())
					// fallback
					return "C", nil
				}

				var orders []int64 = make([]int64, 0)
				var quantities []int32 = make([]int32, 0)
				var totalQuantityManufactured int32 = 0

				for rows.Next() {
					var complexManufacturingOrderForStockId int64
					var manufacturingOrderTypeComponentId int32
					rows.Scan(&complexManufacturingOrderForStockId, &manufacturingOrderTypeComponentId)
					orders = append(orders, complexManufacturingOrderForStockId)

					com := getManufacturingOrderTypeComponentRow(manufacturingOrderTypeComponentId)
					quantities = append(quantities, com.Quantity)
					totalQuantityManufactured += com.Quantity
				}

				if totalQuantityManufactured >= s.Quantity {
					var quantityAssigned int32 = 0
					for i := 0; i < len(orders); i++ {
						sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET sale_order_detail=$2 WHERE id=$1`
						_, err := db.Exec(sqlStatement, orders[i], s.Id)
						if err != nil {
							log("DB", err.Error())
							// fallback
							return "C", nil
						}

						quantityAssigned += quantities[i]
						if quantityAssigned >= s.Quantity {
							break
						}
					}

					return "D", nil
				} else {
					return "C", nil
				}
			} else {
				sqlStatement := `SELECT id, quantity_manufactured FROM manufacturing_order WHERE manufactured = false AND product = $1 AND complex = false ORDER BY date_created ASC`
				rows, err := db.Query(sqlStatement, product.Id)
				if err != nil {
					log("DB", err.Error())
					// fallback
					return "C", nil
				}

				var totalQuantityManufactured int32 = 0
				var orders []int64 = make([]int64, 0)
				var quantities []int32 = make([]int32, 0)

				for rows.Next() {
					var manufacturingOrderForStockId int64
					var quantityManufactured int32
					rows.Scan(&manufacturingOrderForStockId, &quantityManufactured)
					totalQuantityManufactured += quantityManufactured
					orders = append(orders, manufacturingOrderForStockId)
					quantities = append(quantities, quantityManufactured)
				}

				if totalQuantityManufactured < s.Quantity {
					return "C", nil
				} else {
					var quantityAssigned int32 = 0
					for i := 0; i < len(orders); i++ {
						sqlStatement := `UPDATE public.manufacturing_order SET sale_order_detail=$2 WHERE id=$1`
						_, err := db.Exec(sqlStatement, orders[i], s.Id)
						if err != nil {
							log("DB", err.Error())
							// fallback
							return "C", nil
						}

						quantityAssigned += quantities[i]
						if quantityAssigned >= s.Quantity {
							break
						}
					}
					return "D", nil
				}
			}
		} else {
			// search for pending purchases
			sqlStatement := `SELECT id FROM purchase_order_detail WHERE product=$1 AND quantity_delivery_note = 0 AND quantity - quantity_assigned_sale >= $2 ORDER BY (SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order") ASC LIMIT 1`
			row := db.QueryRow(sqlStatement, s.Product, s.Quantity)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				return "A", nil
			}
			var purchaseDetailId int64
			row.Scan(&purchaseDetailId)
			if purchaseDetailId <= 0 {
				return "A", nil
			}

			// add quantity assigned to sale orders
			ok := addQuantityAssignedSalePurchaseOrder(purchaseDetailId, s.Quantity, order.enterprise, userId)
			if !ok {
				return "A", nil
			}

			// set the purchase order detail
			return "B", &purchaseDetailId
		}
	}
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityPendingPackagingSaleOrderDetail(detailId int64, quantity int32, userId int32) bool {
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

	if setSalesOrderState(detail.enterprise, detail.Order, userId) {
		insertTransactionalLog(detail.enterprise, "sales_order_detail", int(detailId), userId, "U")
		return true
	}
	return false
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addQuantityDeliveryNoteSalesOrderDetail(detailId int64, quantity int32, userId int32) bool {

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
		ok = addSalesOrderDeliveryNoteLines(detailBefore.enterprise, detailBefore.Order, userId)
		if !ok {
			return false
		}
	} else if detailBefore.QuantityDeliveryNote == detailBefore.Quantity && detailAfter.QuantityDeliveryNote != detailAfter.Quantity { // undo delivery note generated
		ok = removeSalesOrderDeliveryNoteLines(detailBefore.enterprise, detailBefore.Order, userId)
		if !ok {
			return false
		}
	}

	if err != nil && rows == 0 {
		return false
	} else {
		insertTransactionalLog(detailBefore.enterprise, "sales_order_detail", int(detailId), userId, "U")
		return true
	}
}

func cancelSalesOrderDetail(detailId int64, enterpriseId int32, userId int32) bool {
	detail := getSalesOrderDetailRow(detailId)
	if detail.Id <= 0 {
		return false
	}

	if !detail.Cancelled {
		if detail.Quantity <= 0 || detail.QuantityInvoiced < 0 || detail.QuantityDeliveryNote > 0 {
			return false
		}

		sqlStatement := `UPDATE public.sales_order_detail SET quantity_invoiced=quantity, quantity_delivery_note=quantity, status='Z', cancelled=true WHERE id=$1 AND enterprise=$2`
		_, err := db.Exec(sqlStatement, detailId, enterpriseId)

		if err != nil {
			log("DB", err.Error())
			return false
		}

		ok := setSalesOrderState(enterpriseId, detail.Order, userId)
		if !ok {
			return false
		}

		if err != nil {
			insertTransactionalLog(detail.enterprise, "sales_order_detail", int(detailId), userId, "U")
		}

		return err == nil
	} else {
		if detail.Quantity <= 0 || detail.QuantityInvoiced == 0 || detail.QuantityDeliveryNote == 0 {
			return false
		}

		sqlStatement := `UPDATE public.sales_order_detail SET quantity_invoiced=0, quantity_delivery_note=0, cancelled=false WHERE id=$1 AND enterprise=$2`
		_, err := db.Exec(sqlStatement, detailId, enterpriseId)

		if err != nil {
			log("DB", err.Error())
		}

		status, purchaseOrderDetail := detail.computeStatus(userId)
		sqlStatement = `UPDATE sales_order_detail SET status=$2,purchase_order_detail=$3 WHERE id=$1 AND enterprise=$4`
		_, err = db.Exec(sqlStatement, detailId, status, purchaseOrderDetail, enterpriseId)
		if err != nil {
			log("DB", err.Error())
		}

		ok := setSalesOrderState(enterpriseId, detail.Order, userId)
		if !ok {
			return false
		}

		if err != nil {
			insertTransactionalLog(detail.enterprise, "sales_order_detail", int(detailId), userId, "U")
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
	TotalAmount  float64   `json:"totalAmount"`
}

func getPurchasesOrderDetailsFromSaleOrderDetail(detailId int32, enterpriseId int32) []SalePurchasesOrderDetail {
	salePurchasesOrderDetail := make([]SalePurchasesOrderDetail, 0)
	sqlStatement := `SELECT purchase_order_detail.id,"order",(SELECT order_name FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order"),(SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order"),(SELECT name FROM suppliers WHERE suppliers.id=(SELECT supplier FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order")),quantity,total_amount FROM purchase_order_detail WHERE purchase_order_detail.id=(SELECT purchase_order_detail FROM sales_order_detail WHERE id=$1 AND enterprise=$2)`
	rows, err := db.Query(sqlStatement, detailId, enterpriseId)
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
