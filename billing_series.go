package main

import (
	"strings"
)

type BillingSerie struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	BillingType string `json:"billingType"`
	Year        int16  `json:"year"`
	enterprise  int32
}

func getBillingSeries(enterpriseId int32) []BillingSerie {
	var series []BillingSerie = make([]BillingSerie, 0)
	sqlStatement := `SELECT * FROM public.billing_series WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return series
	}
	for rows.Next() {
		s := BillingSerie{}
		rows.Scan(&s.Id, &s.Name, &s.BillingType, &s.Year, &s.enterprise)
		series = append(series, s)
	}

	return series
}

func (s *BillingSerie) isValid() bool {
	return !(len(s.Id) == 0 || len(s.Id) > 3 || len(s.Name) == 0 || len(s.Name) > 50 || s.Year <= 0 || (s.BillingType != "S" && s.BillingType != "P"))
}

func (s *BillingSerie) insertBillingSerie() bool {
	if !s.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.billing_series(id, name, billing_type, year, enterprise) VALUES ($1, $2, $3, $4, $5)`
	res, err := db.Exec(sqlStatement, s.Id, s.Name, s.BillingType, s.Year, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *BillingSerie) updateBillingSerie() bool {
	if s.Id == "" || !s.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.billing_series SET name=$2, billing_type=$3, year=$4 WHERE id=$1 AND enterprise=$5`
	res, err := db.Exec(sqlStatement, s.Id, s.Name, s.BillingType, s.Year, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *BillingSerie) deleteBillingSerie() bool {
	if s.Id == "" {
		return false
	}

	sqlStatement := `DELETE FROM billing_series WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, s.Id, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func getNextSaleOrderNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_order.order_number) END AS id FROM sales_order WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func getNextSaleInvoiceNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_invoice.invoice_number) END AS id FROM sales_invoice WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func getNextSaleDeliveryNoteNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(sales_delivery_note.delivery_note_number) END AS id FROM sales_delivery_note WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func getNextPurchaseOrderNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_order.order_number) END AS id FROM purchase_order WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func getNextPurchaseInvoiceNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_invoice.invoice_number) END AS id FROM purchase_invoice WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func getNextPurchaseDeliveryNoteNumber(billingSerieId string, enterpriseId int32) int32 {
	sqlStatement := `SELECT CASE COUNT(*) WHEN 0 THEN 0 ELSE MAX(purchase_delivery_note.delivery_note_number) END AS id FROM purchase_delivery_note WHERE "billing_series"=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, billingSerieId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var orderNumber int32
	row.Scan(&orderNumber)
	return (orderNumber + 1)
}

func findBillingSerieByName(billingSerieName string, enterpriseId int32) []NameString {
	var billingSeries []NameString = make([]NameString, 0)
	sqlStatement := `SELECT id,name FROM public.billing_series WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(billingSerieName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return billingSeries
	}
	for rows.Next() {
		b := NameString{}
		rows.Scan(&b.Id, &b.Name)
		billingSeries = append(billingSeries, b)
	}

	return billingSeries
}

func getNameBillingSerie(id string, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.billing_series WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type LocateBillingSerie struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func locateBillingSeries(enterpriseId int32) []LocateBillingSerie {
	var series []LocateBillingSerie = make([]LocateBillingSerie, 0)
	sqlStatement := `SELECT id,name FROM public.billing_series WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return series
	}
	for rows.Next() {
		s := LocateBillingSerie{}
		rows.Scan(&s.Id, &s.Name)
		series = append(series, s)
	}

	return series
}
