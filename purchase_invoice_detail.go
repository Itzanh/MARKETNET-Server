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

func (s *PurchaseInvoiceDetail) insertPurchaseInvoiceDetail(beginTransaction bool) bool {
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

	sqlStatement := `INSERT INTO public.purchase_invoice_details(invoice, product, price, quantity, vat_percent, total_amount, order_detail, enterprise, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	res, err := db.Exec(sqlStatement, s.Invoice, s.Product, s.Price, s.Quantity, s.VatPercent, s.TotalAmount, s.OrderDetail, s.enterprise, s.Description)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	ok := addTotalProductsPurchaseInvoice(s.Invoice, s.Price*float64(s.Quantity), s.VatPercent)
	if !ok {
		if beginTransaction {
			trans.Rollback()
		}
		return false
	}
	if s.OrderDetail != nil && *s.OrderDetail != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*s.OrderDetail, s.Quantity)
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

	return rows > 0
}

func (d *PurchaseInvoiceDetail) deletePurchaseInvoiceDetail() bool {
	if d.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	detailInMemory := getPurchaseInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return false
	}
	sqlStatement := `DELETE FROM public.purchase_invoice_details WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	ok := addTotalProductsPurchaseInvoice(detailInMemory.Invoice, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent)
	if !ok {
		trans.Rollback()
		return false
	}
	if detailInMemory.OrderDetail != nil && *detailInMemory.OrderDetail != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*detailInMemory.OrderDetail, -detailInMemory.Quantity)
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
