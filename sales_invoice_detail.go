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

func (s *SalesInvoiceDetail) insertSalesInvoiceDetail(beginTransaction bool, userId int32) bool {
	if !s.isValid() {
		return false
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

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

	sqlStatement := `INSERT INTO public.sales_invoice_detail(invoice, product, price, quantity, vat_percent, total_amount, order_detail, enterprise, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Invoice, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.OrderDetail, s.enterprise, s.Description)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var saleInvoiceDetailId int64
	row.Scan(&saleInvoiceDetailId)
	s.Id = saleInvoiceDetailId

	ok := addTotalProductsSalesInvoice(s.Invoice, s.Price*float64(s.Quantity), s.VatPercent, s.enterprise, userId)
	if !ok {
		if beginTransaction {
			trans.Rollback()
		}
		return false
	}
	if s.OrderDetail != nil && *s.OrderDetail != 0 {
		ok := addQuantityInvociedSalesOrderDetail(*s.OrderDetail, s.Quantity, userId)
		if !ok {
			if beginTransaction {
				trans.Rollback()
			}
			return false
		}
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false
		}
		///
	}

	if saleInvoiceDetailId > 0 {
		insertTransactionalLog(s.enterprise, "sales_invoice_detail", int(saleInvoiceDetailId), userId, "I")
	}

	return saleInvoiceDetailId > 0
}

func (d *SalesInvoiceDetail) deleteSalesInvoiceDetail(userId int32) bool {
	if d.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(d.enterprise)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return false
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		i := getSalesInvoiceRow(d.Invoice)
		invoiceNumber := getNextSaleInvoiceNumber(i.BillingSeries, i.enterprise)
		if invoiceNumber <= 0 || i.InvoiceNumber != (invoiceNumber-1) {
			return false
		}
	}

	detailInMemory := getSalesInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return false
	}

	insertTransactionalLog(d.enterprise, "sales_invoice_detail", int(d.Id), userId, "D")

	sqlStatement := `DELETE FROM public.sales_invoice_detail WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, d.Id, d.enterprise)
	rows, _ := res.RowsAffected()
	if err != nil {
		log("DB", err.Error())
		return false
	}

	// can't continue
	if rows == 0 {
		///
		err = trans.Commit()
		if err != nil {
			return false
		}
		///

		return rows > 0
	}

	ok := addTotalProductsSalesInvoice(detailInMemory.Invoice, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, d.enterprise, userId)
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
			ok = addQuantityAssignedSalePurchaseOrder(*detail.PurchaseOrderDetail, -detail.Quantity, detailInMemory.enterprise, userId)
			if !ok {
				trans.Rollback()
				return false
			}
		}
		// revert back the status
		ok := addQuantityInvociedSalesOrderDetail(*detailInMemory.OrderDetail, -detailInMemory.Quantity, userId)
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

	return rows > 0
}
