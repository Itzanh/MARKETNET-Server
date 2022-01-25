package main

import (
	"database/sql"
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
		a.CountryName = getCountryRow(a.Country, enterpriseId).Name
		if a.State != nil {
			stateName := getNameState(*a.State, enterpriseId)
			a.StateName = &stateName
		}

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
		a.CountryName = getCountryRow(a.Country, enterpriseId).Name
		if a.State != nil {
			stateName := getNameState(*a.State, enterpriseId)
			a.StateName = &stateName
		}

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
