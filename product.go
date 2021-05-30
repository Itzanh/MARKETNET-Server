package main

import (
	"strings"
	"time"
)

type Product struct {
	Id                     int32     `json:"id"`
	Name                   string    `json:"name"`
	Reference              string    `json:"reference"`
	BarCode                string    `json:"barCode"`
	ControlStock           bool      `json:"controlStock"`
	Weight                 float32   `json:"weight"`
	Family                 *int16    `json:"family"`
	Width                  float32   `json:"width"`
	Height                 float32   `json:"height"`
	Depth                  float32   `json:"depth"`
	Off                    bool      `json:"off"`
	Stock                  int32     `json:"stock"`
	VatPercent             float32   `json:"vatPercent"`
	DateCreated            time.Time `json:"dateCreated"`
	Description            string    `json:"description"`
	Color                  *int16    `json:"color"`
	Price                  float32   `json:"price"`
	Manufacturing          bool      `json:"manufacturing"`
	ManufacturingOrderType *int16    `json:"manufacturingOrderType"`
}

func getProduct() []Product {
	var products []Product = make([]Product, 0)
	sqlStatement := `SELECT * FROM public.product ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return products
	}
	for rows.Next() {
		p := Product{}
		rows.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType)
		products = append(products, p)
	}

	return products
}

func getProductRow(productId int32) Product {
	sqlStatement := `SELECT * FROM public.product WHERE id = $1`
	row := db.QueryRow(sqlStatement, productId)
	if row.Err() != nil {
		return Product{}
	}

	p := Product{}
	row.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType)

	return p
}

func (p *Product) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 150 || len(p.Reference) == 0 || len(p.Reference) > 40 || len(p.BarCode) != 13 || p.VatPercent < 0)
}

func (p *Product) insertProduct() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.product(name, reference, barcode, control_stock, weight, family, width, height, depth, off, stock, vat_percent, dsc, color, price, manufacturing, manufacturing_order_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, &17)`
	res, err := db.Exec(sqlStatement, p.Name, p.Reference, &p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Product) updateProduct() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.product SET name=$2, reference=$3, barcode=$4, control_stock=$5, weight=$6, family=$7, width=$8, height=$9, depth=$10, off=$11, stock=$12, vat_percent=$13, dsc=$14, color=$15, price=$16, manufacturing=$17, manufacturing_order_type=$18 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.Reference, p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Product) deleteProduct() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.product WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type ProductName struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func findProductByName(languageName string) []ProductName {
	var products []ProductName = make([]ProductName, 0)
	sqlStatement := `SELECT id,name FROM public.product WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return products
	}
	for rows.Next() {
		p := ProductName{}
		rows.Scan(&p.Id, &p.Name)
		products = append(products, p)
	}

	return products
}

func getNameProduct(id int32) string {
	sqlStatement := `SELECT name FROM public.product WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type OrderDetailDefaults struct {
	Price      float32 `json:"price"`
	VatPercent float32 `json:"vatPercent"`
}

func getOrderDetailDefaults(roductId int32) OrderDetailDefaults {
	sqlStatement := `SELECT price, vat_percent FROM product WHERE id = $1`
	row := db.QueryRow(sqlStatement, roductId)
	if row.Err() != nil {
		return OrderDetailDefaults{}
	}
	s := OrderDetailDefaults{}
	row.Scan(&s.Price, &s.VatPercent)
	return s
}
