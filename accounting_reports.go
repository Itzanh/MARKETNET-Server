package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

/* REPORT 111 */

type Form111Query struct {
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
}

type Form111Result struct {
	Elements          []Form111Element `json:"elements"`
	TotalWithDiscount float64          `json:"totalWithDiscount"`
	IncomeTaxBase     float64          `json:"incomeTaxBase"`
	IncomeTaxValue    float64          `json:"incomeTaxValue"`
	TotalAmount       float64          `json:"totalAmount"`
}

type Form111Element struct {
	PurchaseInvoice PurchaseInvoice `json:"purchaseInvoice"`
	Address         Address         `json:"address"`
	Supplier        Supplier        `json:"supplier"`
}

func (q *Form111Query) execReportForm111(enterpriseId int32) Form111Result {
	r := Form111Result{}
	r.Elements = make([]Form111Element, 0)
	sqlStatement := `SELECT id,supplier,billing_address FROM purchase_invoice WHERE income_tax AND date_created >= $1 AND date_created <= $2 AND enterprise = $3 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, q.DateStart, q.DateEnd, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return r
	}

	for rows.Next() {
		var id int64
		var supplier int32
		var address int32
		rows.Scan(&id, &supplier, &address)

		i := getPurchaseInvoiceRow(id)
		a := getAddressRow(address)

		r.Elements = append(r.Elements, Form111Element{
			PurchaseInvoice: i,
			Supplier:        getSupplierRow(supplier),
			Address:         a,
		})

		r.TotalWithDiscount += i.TotalWithDiscount
		r.IncomeTaxBase += i.IncomeTaxBase
		r.IncomeTaxValue += i.IncomeTaxValue
		r.TotalAmount += i.TotalAmount
	}

	return r
}

/* REPORT 115*/

type Form115Query struct {
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
}

type Form115Result struct {
	Elements          []Form115Element `json:"elements"`
	TotalWithDiscount float64          `json:"totalWithDiscount"`
	IncomeTaxBase     float64          `json:"incomeTaxBase"`
	IncomeTaxValue    float64          `json:"incomeTaxValue"`
	TotalAmount       float64          `json:"totalAmount"`
}

type Form115Element struct {
	PurchaseInvoice PurchaseInvoice `json:"purchaseInvoice"`
	Address         Address         `json:"address"`
	Supplier        Supplier        `json:"supplier"`
}

func (q *Form115Query) execReportForm115(enterpriseId int32) Form115Result {
	r := Form115Result{}
	r.Elements = make([]Form115Element, 0)
	sqlStatement := `SELECT id,supplier,billing_address FROM purchase_invoice WHERE rent AND date_created >= $1 AND date_created <= $2 AND enterprise=$3 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, q.DateStart, q.DateEnd, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return r
	}

	for rows.Next() {
		var id int64
		var supplier int32
		var address int32
		rows.Scan(&id, &supplier, &address)

		i := getPurchaseInvoiceRow(id)
		a := getAddressRow(address)

		r.Elements = append(r.Elements, Form115Element{
			PurchaseInvoice: i,
			Supplier:        getSupplierRow(supplier),
			Address:         a,
		})

		r.TotalWithDiscount += i.TotalWithDiscount
		r.IncomeTaxBase += i.IncomeTaxBase
		r.IncomeTaxValue += i.IncomeTaxValue
		r.TotalAmount += i.TotalAmount
	}

	return r
}

type InventoyValuationQuery struct {
	Date          time.Time `json:"date"`
	ProductFamily *int32    `json:"productFamily"`
}

type InventoyValuation struct {
	Product     int32   `json:"product"`
	ProductName string  `json:"productName"`
	Quantity    int32   `json:"quantity"`
	CostPrice   float64 `json:"costPrice"`
	Value       float64 `json:"value"`
}

func (q *InventoyValuationQuery) getInventoyValuation(enterpriseId int32) []InventoyValuation {
	var inventoyValuation []InventoyValuation = make([]InventoyValuation, 0)

	var rows *sql.Rows
	var err error
	if q.ProductFamily == nil {
		sqlStatement := `SELECT id, cost_price, name FROM product WHERE enterprise = $1 ORDER BY id ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT id, cost_price, name FROM product WHERE enterprise = $1 AND family = $2 ORDER BY id ASC`
		rows, err = db.Query(sqlStatement, enterpriseId, q.ProductFamily)
	}

	if err != nil {
		log("DB", err.Error())
		return inventoyValuation
	}

	sqlStatement := `SELECT dragged_stock FROM warehouse_movement WHERE product = $1 ORDER BY id DESC LIMIT 1`
	var productId int32
	var costPrice float64
	var productName string
	var draggedStock int32
	for rows.Next() {
		rows.Scan(&productId, &costPrice, &productName)

		// get inventory at date:
		row := db.QueryRow(sqlStatement, productId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			draggedStock = 0
		} else {
			row.Scan(&draggedStock)
		}

		i := InventoyValuation{
			Product:     productId,
			ProductName: productName,
			Quantity:    draggedStock,
			CostPrice:   costPrice,
			Value:       costPrice * float64(draggedStock),
		}
		inventoyValuation = append(inventoyValuation, i)
	}

	return inventoyValuation
}

type IntrastatReportQuery struct {
	DateStart         time.Time `json:"dateStart"`
	DateEnd           time.Time `json:"dateEnd"`
	CountryOriginCode string    `json:"countryOriginCode"`
	StateOriginCode   int       `json:"stateOriginCode"`
	enterpriseId      int32
}

type IntrastatReport struct {
	ReportSales    string `json:"reportSales"`
	ReportPurchase string `json:"reportPurchase"`
}

// Internal use only
type IntrastatReportDetail struct {
	HSCode        string
	IsoCode2      string
	Weight        float64
	Quantity      int64
	TotalInvoiced float64
}

// FR;31;FOB;11;3;;85182190;CN;1;115;162;15,37;15,37
// DE;28;CIF;11;1;0811;85182190;US;1;2459;1982;4589,46;4589,46
// IT;12;FOB;11;3;;02012030;;1;800;;987,00;890,45
//
// ES;11;DDP;11;3;;84261100;AT;1;0.000;8;191.20;191.20
// ES;11;DDP;11;3;;84261100;BG;1;120.000;10;629.20;629.20
// ES;11;DDP;11;3;;84261100;RO;1;0.000;9;215.10;215.10
func (q *IntrastatReportQuery) intrastatReport() IntrastatReport {
	var report IntrastatReport = IntrastatReport{}

	// SALES
	var salesDetails []IntrastatReportDetail = make([]IntrastatReportDetail, 0)

	sqlStatement := `SELECT id,(SELECT iso_2 FROM country WHERE country.id=(SELECT country FROM address WHERE address.id=sales_order.shipping_address)) FROM sales_order WHERE enterprise = $1 AND date_created >= $2 AND date_created <= $3 AND lines_number = delivery_note_lines AND (SELECT zone FROM country WHERE country.id=(SELECT country FROM address WHERE address.id=sales_order.shipping_address)) = 'U' ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, q.enterpriseId, q.DateStart, q.DateEnd)
	if err != nil {
		log("DB", err.Error())
		return report
	}

	var ok bool = false
	// order
	var saleOrderId int64
	var destinationCountryIso2 string
	// details
	var weight float64
	var hsCode *string
	var quantity int32
	var totalAmount float64
	for rows.Next() {
		rows.Scan(&saleOrderId, &destinationCountryIso2)

		// get the details
		sqlStatement = `SELECT (SELECT weight FROM product WHERE product.id=sales_order_detail.product),(SELECT hs_code FROM product WHERE product.id=sales_order_detail.product),quantity,total_amount FROM sales_order_detail WHERE "order" = $1`
		rowsDetails, err := db.Query(sqlStatement, saleOrderId)
		if err != nil {
			log("DB", err.Error())
			return report
		}

		for rowsDetails.Next() {
			rowsDetails.Scan(&weight, &hsCode, &quantity, &totalAmount)

			if hsCode == nil {
				hsCode = new(string)
			}

			ok = false
			for i := 0; i < len(salesDetails); i++ {
				var salesDetail IntrastatReportDetail = salesDetails[i]
				if salesDetail.HSCode == *hsCode && salesDetail.IsoCode2 == destinationCountryIso2 {
					salesDetail.Weight = weight * float64(quantity)
					salesDetail.Quantity += int64(quantity)
					salesDetail.TotalInvoiced += totalAmount
					ok = true
					break
				}
			}
			if !ok {
				salesDetails = append(salesDetails, IntrastatReportDetail{
					HSCode:        *hsCode,
					IsoCode2:      destinationCountryIso2,
					Weight:        weight * float64(quantity),
					Quantity:      int64(quantity),
					TotalInvoiced: totalAmount,
				})
			}
		}
	}

	// PURCHASES
	var purchasesDetails []IntrastatReportDetail = make([]IntrastatReportDetail, 0)

	sqlStatement = `SELECT id,(SELECT iso_2 FROM country WHERE country.id=(SELECT country FROM address WHERE address.id=purchase_order.shipping_address)) FROM purchase_order WHERE enterprise = $1 AND date_created >= $2 AND date_created <= $3 AND lines_number = delivery_note_lines AND (SELECT zone FROM country WHERE country.id=(SELECT country FROM address WHERE address.id=purchase_order.shipping_address)) = 'U' ORDER BY id ASC`
	rows, err = db.Query(sqlStatement, q.enterpriseId, q.DateStart, q.DateEnd)
	if err != nil {
		log("DB", err.Error())
		return report
	}

	// order
	var purchaseOrderId int64
	for rows.Next() {
		rows.Scan(&purchaseOrderId, &destinationCountryIso2)

		// get the details
		sqlStatement = `SELECT (SELECT weight FROM product WHERE product.id=purchase_order_detail.product),(SELECT hs_code FROM product WHERE product.id=purchase_order_detail.product),quantity,total_amount FROM purchase_order_detail WHERE "order" = $1`
		rowsDetails, err := db.Query(sqlStatement, purchaseOrderId)
		if err != nil {
			log("DB", err.Error())
			return report
		}

		for rowsDetails.Next() {
			rowsDetails.Scan(&weight, &hsCode, &quantity, &totalAmount)

			if hsCode == nil {
				hsCode = new(string)
			}

			ok = false
			for i := 0; i < len(purchasesDetails); i++ {
				var purchasesDetail IntrastatReportDetail = purchasesDetails[i]
				if purchasesDetail.HSCode == *hsCode && purchasesDetail.IsoCode2 == destinationCountryIso2 {
					purchasesDetail.Weight = weight * float64(quantity)
					purchasesDetail.Quantity += int64(quantity)
					purchasesDetail.TotalInvoiced += totalAmount
					ok = true
					break
				}
			}
			if !ok {
				purchasesDetails = append(purchasesDetails, IntrastatReportDetail{
					HSCode:        *hsCode,
					IsoCode2:      destinationCountryIso2,
					Weight:        weight * float64(quantity),
					Quantity:      int64(quantity),
					TotalInvoiced: totalAmount,
				})
			}
		}
	}

	// GENERATE REPORT: SALES
	for i := 0; i < len(salesDetails); i++ {
		var salesDetail IntrastatReportDetail = salesDetails[i]
		report.ReportSales += q.CountryOriginCode + ";" + strconv.Itoa(q.StateOriginCode) + ";DDP;11;3;;" + salesDetail.HSCode + ";" + salesDetail.IsoCode2 + ";1;" + fmt.Sprintf("%.3f", salesDetail.Weight) + ";" + strconv.Itoa(int(salesDetail.Quantity)) + ";" + fmt.Sprintf("%.2f", salesDetail.TotalInvoiced) + ";" + fmt.Sprintf("%.2f", salesDetail.TotalInvoiced) + "\n"
	}

	// GENERATE REPORT: PURCHASES
	for i := 0; i < len(purchasesDetails); i++ {
		var purchasesDetail IntrastatReportDetail = purchasesDetails[i]
		report.ReportPurchase += q.CountryOriginCode + ";" + strconv.Itoa(q.StateOriginCode) + ";DDP;11;3;;" + purchasesDetail.HSCode + ";" + purchasesDetail.IsoCode2 + ";1;" + fmt.Sprintf("%.3f", purchasesDetail.Weight) + ";" + strconv.Itoa(int(purchasesDetail.Quantity)) + ";" + fmt.Sprintf("%.2f", purchasesDetail.TotalInvoiced) + ";" + fmt.Sprintf("%.2f", purchasesDetail.TotalInvoiced) + "\n"
	}

	return report
}
