package main

import (
	"database/sql"
)

type PurchaseInvoiceDetail struct {
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

func getPurchaseInvoiceDetail(invoiceId int64, enterpriseId int32) []PurchaseInvoiceDetail {
	var details []PurchaseInvoiceDetail = make([]PurchaseInvoiceDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=purchase_invoice_details.product) FROM purchase_invoice_details WHERE invoice=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, invoiceId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	defer rows.Close()

	for rows.Next() {
		d := PurchaseInvoiceDetail{}
		rows.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail, &d.enterprise, &d.Description, &d.ProductName)
		details = append(details, d)
	}

	return details
}

func getPurchaseInvoiceDetailRow(detailId int64) PurchaseInvoiceDetail {
	sqlStatement := `SELECT * FROM purchase_invoice_details WHERE id=$1 ORDER BY id ASC`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PurchaseInvoiceDetail{}
	}
	d := PurchaseInvoiceDetail{}
	row.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail, &d.enterprise, &d.Description)

	return d
}

func (d *PurchaseInvoiceDetail) isValid() bool {
	return !(d.Invoice <= 0 || (d.Product == nil && len(d.Description) == 0) || len(d.Description) > 150 || (d.Product != nil && *d.Product <= 0) || d.Quantity <= 0 || d.VatPercent < 0)
}

// ERROR CODES:
// 1. the product is deactivated
// 2. there is aleady a detail with this product
// 3. can't add details to a posted invoice
func (s *PurchaseInvoiceDetail) insertPurchaseInvoiceDetail(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
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
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	sqlStatement := `SELECT COUNT(purchase_invoice_details) FROM public.purchase_invoice_details WHERE invoice = $1 AND product = $2`
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
	invoice := getPurchaseInvoiceRowTransaction(s.Invoice, *trans)
	if invoice.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if invoice.AccountingMovement != nil {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
	}

	sqlStatement = `INSERT INTO public.purchase_invoice_details(invoice, product, price, quantity, vat_percent, total_amount, order_detail, enterprise, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row = trans.QueryRow(sqlStatement, s.Invoice, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.OrderDetail, s.enterprise, s.Description)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	var invoiceDetailId int64
	row.Scan(&invoiceDetailId)
	s.Id = invoiceDetailId

	insertTransactionalLog(s.enterprise, "purchase_invoice_details", int(invoiceDetailId), userId, "I")

	ok := addTotalProductsPurchaseInvoice(s.Invoice, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if s.OrderDetail != nil && *s.OrderDetail != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*s.OrderDetail, s.Quantity, s.enterprise, userId, *trans)
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

	return OkAndErrorCodeReturn{Ok: invoiceDetailId > 0}
}

// ERROR CODES
// 1. can't delete posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (d *PurchaseInvoiceDetail) deletePurchaseInvoiceDetail(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
	if d.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	detailInMemory := getPurchaseInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	i := getPurchaseInvoiceRow(detailInMemory.Invoice)
	if i.AccountingMovement != nil { // can't delete posted invoices
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(d.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		i := getPurchaseInvoiceRow(d.Invoice)
		invoiceNumber := getNextPurchaseInvoiceNumber(i.BillingSeries, i.enterprise)
		if invoiceNumber <= 0 || i.InvoiceNumber != (invoiceNumber-1) {
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
		}
	}

	insertTransactionalLog(detailInMemory.enterprise, "purchase_invoice_details", int(d.Id), userId, "D")

	sqlStatement := `DELETE FROM public.purchase_invoice_details WHERE id=$1 AND enterprise=$2`
	res, err := trans.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsPurchaseInvoice(detailInMemory.Invoice, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, d.enterprise, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.OrderDetail != nil && *detailInMemory.OrderDetail != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*detailInMemory.OrderDetail, -detailInMemory.Quantity, d.enterprise, userId, *trans)
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

	return OkAndErrorCodeReturn{Ok: rows > 0}
}
