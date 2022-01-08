package main

import (
	"database/sql"
)

type SalesInvoiceDetail struct {
	Id          int64   `json:"id"`
	Invoice     int64   `json:"invoice"`
	Product     *int32  `json:"product"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
	VatPercent  float64 `json:"vatPercent"`
	TotalAmount float64 `json:"totalAmount"`
	OrderDetail *int64  `json:"orderDetail"`
	ProductName string  `json:"productName"`
	Description string  `json:"description"`
	enterprise  int32
}

func getSalesInvoiceDetail(invoiceId int64, enterpriseId int32) []SalesInvoiceDetail {
	var details []SalesInvoiceDetail = make([]SalesInvoiceDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_invoice_detail.product) FROM sales_invoice_detail WHERE invoice=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, invoiceId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	defer rows.Close()

	for rows.Next() {
		d := SalesInvoiceDetail{}
		rows.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail, &d.enterprise, &d.Description, &d.ProductName)
		details = append(details, d)
	}

	return details
}

func getSalesInvoiceDetailRow(detailId int64) SalesInvoiceDetail {
	sqlStatement := `SELECT * FROM sales_invoice_detail WHERE id = $1 ORDER BY id ASC`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesInvoiceDetail{}
	}
	d := SalesInvoiceDetail{}
	row.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail, &d.enterprise, &d.Description)

	return d
}

func (d *SalesInvoiceDetail) isValid() bool {
	return !(d.Invoice <= 0 || (d.Product == nil && len(d.Description) == 0) || len(d.Description) > 150 || (d.Product != nil && *d.Product <= 0) || d.Quantity <= 0 || d.VatPercent < 0)
}

// ERROR CODES:
// 1. the product is deactivated
// 2. there is aleady a detail with this product
// 3. can't add details to a posted invoice
func (s *SalesInvoiceDetail) insertSalesInvoiceDetail(trans *sql.Tx, userId int32) OkAndErrorCodeReturn {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	if s.Product != nil {
		p := getProductRow(*s.Product)
		if p.Id <= 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if p.Off {
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
		}
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		trn, err := db.Begin()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		trans = trn
		///
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	sqlStatement := `SELECT COUNT(sales_invoice_detail) FROM public.sales_invoice_detail WHERE invoice = $1 AND product = $2`
	row := db.QueryRow(sqlStatement, s.Invoice, s.Product)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	var countProductInSaleOrder int16
	row.Scan(&countProductInSaleOrder)
	if countProductInSaleOrder > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	}

	// can't add details to a posted invoice
	invoice := getSalesInvoiceRowTransaction(s.Invoice, *trans)
	if invoice.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if invoice.AccountingMovement != nil {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
	}

	sqlStatement = `INSERT INTO public.sales_invoice_detail(invoice, product, price, quantity, vat_percent, total_amount, order_detail, enterprise, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row = trans.QueryRow(sqlStatement, s.Invoice, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.OrderDetail, s.enterprise, s.Description)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	var saleInvoiceDetailId int64
	row.Scan(&saleInvoiceDetailId)
	s.Id = saleInvoiceDetailId

	ok := addTotalProductsSalesInvoice(s.Invoice, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if s.OrderDetail != nil && *s.OrderDetail != 0 {
		ok := addQuantityInvociedSalesOrderDetail(*s.OrderDetail, s.Quantity, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	if saleInvoiceDetailId > 0 {
		insertTransactionalLog(s.enterprise, "sales_invoice_detail", int(saleInvoiceDetailId), userId, "I")
	}

	return OkAndErrorCodeReturn{Ok: saleInvoiceDetailId > 0}
}

// ERROR CODES:
// 1. can't delete posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (d *SalesInvoiceDetail) deleteSalesInvoiceDetail(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
	if d.Id <= 0 {
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

	detailInMemory := getSalesInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	i := getSalesInvoiceRow(detailInMemory.Invoice)
	if i.AccountingMovement != nil { // can't delete posted invoices
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(d.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoiceNumber := getNextSaleInvoiceNumber(i.BillingSeries, i.enterprise)
		if invoiceNumber <= 0 || i.InvoiceNumber != (invoiceNumber-1) {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
		}
	}

	sqlStatement := `UPDATE public.sales_order_discount SET sales_invoice_detail = NULL WHERE sales_invoice_detail=$1 AND enterprise=$2`
	_, err := trans.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		trans.Rollback()
		log("DB", err.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(d.enterprise, "sales_invoice_detail", int(d.Id), userId, "D")

	sqlStatement = `DELETE FROM public.sales_invoice_detail WHERE id=$1 AND enterprise=$2`
	res, err := trans.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		trans.Rollback()
		log("DB", err.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	rows, _ := res.RowsAffected()

	// can't continue
	if rows == 0 {
		///
		err = trans.Rollback()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///

		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsSalesInvoice(detailInMemory.Invoice, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, d.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.OrderDetail != nil && *detailInMemory.OrderDetail != 0 {
		detail := getSalesOrderDetailRow(*detailInMemory.OrderDetail)
		if detail.Id <= 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
		// if the detail had a purchase order pending, rollback the quantity assigned
		if detail.Status == "B" {
			ok = addQuantityAssignedSalePurchaseOrder(*detail.PurchaseOrderDetail, -detail.Quantity, detailInMemory.enterprise, userId, *trans)
			if !ok {
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}
		}
		// revert back the status
		ok := addQuantityInvociedSalesOrderDetail(*detailInMemory.OrderDetail, -detailInMemory.Quantity, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	if beginTransaction {
		///
		err = trans.Commit()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	return OkAndErrorCodeReturn{Ok: rows > 0}
}
