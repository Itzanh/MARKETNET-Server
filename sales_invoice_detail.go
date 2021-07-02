package main

import (
	"database/sql"
)

type SalesInvoiceDetail struct {
	Id          int32   `json:"id"`
	Invoice     int32   `json:"invoice"`
	Product     int32   `json:"product"`
	Price       float32 `json:"price"`
	Quantity    int32   `json:"quantity"`
	VatPercent  float32 `json:"vatPercent"`
	TotalAmount float32 `json:"totalAmount"`
	OrderDetail *int32  `json:"orderDetail"`
	ProductName string  `json:"productName"`
}

func getSalesInvoiceDetail(invoiceId int32) []SalesInvoiceDetail {
	var details []SalesInvoiceDetail = make([]SalesInvoiceDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_invoice_detail.product) FROM sales_invoice_detail WHERE invoice = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, invoiceId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := SalesInvoiceDetail{}
		rows.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail, &d.ProductName)
		details = append(details, d)
	}

	return details
}

func getSalesInvoiceDetailRow(detailId int32) SalesInvoiceDetail {
	sqlStatement := `SELECT * FROM sales_invoice_detail WHERE id = $1 ORDER BY id ASC`
	row := db.QueryRow(sqlStatement, detailId)
	if row.Err() != nil {
		return SalesInvoiceDetail{}
	}
	d := SalesInvoiceDetail{}
	row.Scan(&d.Id, &d.Invoice, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.OrderDetail)

	return d
}

func (d *SalesInvoiceDetail) isValid() bool {
	return !(d.Invoice <= 0 || d.Product <= 0 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (s *SalesInvoiceDetail) insertSalesInvoiceDetail(beginTransaction bool) bool {
	if !s.isValid() {
		return false
	}

	s.TotalAmount = (s.Price * float32(s.Quantity)) * (1 + (s.VatPercent / 100))

	var trans *sql.Tx
	if beginTransaction {
		///
		trn, err := db.Begin()
		if err != nil {
			return false
		}
		trans = trn
		///
	}

	sqlStatement := `INSERT INTO public.sales_invoice_detail(invoice, product, price, quantity, vat_percent, total_amount, order_detail) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, s.Invoice, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.OrderDetail)
	if err != nil {
		return false
	}
	ok := addTotalProductsSalesInvoice(s.Invoice, s.Price*float32(s.Quantity), s.VatPercent)
	if !ok {
		if beginTransaction {
			trans.Rollback()
		}
		return false
	}
	if s.OrderDetail != nil && *s.OrderDetail != 0 {
		ok := addQuantityInvociedSalesOrderDetail(*s.OrderDetail, s.Quantity)
		if !ok {
			if beginTransaction {
				trans.Rollback()
			}
			return false
		}
	}

	if beginTransaction {
		///
		err = trans.Commit()
		if err != nil {
			return false
		}
		///
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (d *SalesInvoiceDetail) deleteSalesInvoiceDetail() bool {
	if d.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	detailInMemory := getSalesInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return false
	}
	sqlStatement := `DELETE FROM public.sales_invoice_detail WHERE id=$1`
	res, err := db.Exec(sqlStatement, d.Id)
	if err != nil {
		return false
	}
	ok := addTotalProductsSalesInvoice(detailInMemory.Invoice, -(detailInMemory.Price * float32(detailInMemory.Quantity)), detailInMemory.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	if detailInMemory.OrderDetail != nil && *detailInMemory.OrderDetail != 0 {
		detail := getSalesOrderDetailRow(*detailInMemory.OrderDetail)
		if detail.Id <= 0 {
			trans.Rollback()
			return false
		}
		// if the detail had a purchase order pending, rollback the quantity assigned
		if detail.Status == "B" {
			ok = addQuantityAssignedSalePurchaseOrder(*detail.PurchaseOrderDetail, -detail.Quantity)
			if !ok {
				trans.Rollback()
				return false
			}
		}
		// revert back the status
		ok := addQuantityInvociedSalesOrderDetail(*detailInMemory.OrderDetail, -detailInMemory.Quantity)
		if !ok {
			trans.Rollback()
			return false
		}
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
