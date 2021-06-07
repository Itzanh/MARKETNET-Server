package main

import (
	"strconv"
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
	Supplier               *int32    `json:"supplier"`
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
		rows.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier)
		products = append(products, p)
	}

	return products
}

func searchProduct(search string) []Product {
	var products []Product = make([]Product, 0)
	sqlStatement := `SELECT * FROM product WHERE name ILIKE $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%")
	if err != nil {
		return products
	}
	for rows.Next() {
		p := Product{}
		rows.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier)
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
	row.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier)

	return p
}

func (p *Product) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 150 || len(p.Reference) > 40 || (len(p.BarCode) != 0 && len(p.BarCode) != 13) || p.VatPercent < 0)
}

func (p *Product) insertProduct() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.product(name, reference, barcode, control_stock, weight, family, width, height, depth, off, stock, vat_percent, dsc, color, price, manufacturing, manufacturing_order_type, supplier) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	res, err := db.Exec(sqlStatement, p.Name, p.Reference, &p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType, p.Supplier)
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

	sqlStatement := `UPDATE public.product SET name=$2, reference=$3, barcode=$4, control_stock=$5, weight=$6, family=$7, width=$8, height=$9, depth=$10, off=$11, stock=$12, vat_percent=$13, dsc=$14, color=$15, price=$16, manufacturing=$17, manufacturing_order_type=$18, supplier=$19 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.Reference, p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType, p.Supplier)
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

// Get the sales order details with pending status, with the product specified.
func getProductSalesOrderDetailsPending(productId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE product=$1 AND quantity_delivery_note!=quantity`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail)
		details = append(details, d)
	}

	return details
}

// Get the purchase order details with pending status, with the product specified.
func getProductPurchaseOrderDetailsPending(productId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE product=$1 AND quantity_delivery_note!=quantity`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale)
		details = append(details, d)
	}

	return details
}

// Get the sales order details with the product specified.
func getProductSalesOrderDetails(productId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT * FROM sales_order_detail WHERE "order"=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail)
		details = append(details, d)
	}

	return details
}

// Get the purchase order details with the product specified.
func getProductPurchaseOrderDetails(productId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT * FROM purchase_order_detail WHERE product=$1`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale)
		details = append(details, d)
	}

	return details
}

// Get the warehouse movements with the product specified.
func getProductWarehouseMovement(productId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT * FROM warehouse_movement WHERE product=$1`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

func (p *Product) generateBarcode() bool {
	sqlStatement := `SELECT SUBSTRING(barcode,0,13) FROM product WHERE SUBSTRING(barcode,0,5) = $1 ORDER BY barcode DESC LIMIT 1`
	row := db.QueryRow(sqlStatement, getSettingsRecord().BarcodePrefix)
	if row.Err() != nil {
		return false
	}

	var barcode string
	row.Scan(&barcode)
	if len(barcode) == 0 {
		return false
	}

	code, err := strconv.Atoi(barcode)
	if err != nil {
		return false
	}
	code++
	barcode = strconv.Itoa(code)

	checkCode := 0
	for i := 1; i < len(barcode); i += 2 {
		j, _ := strconv.Atoi(barcode[i : i+1])
		checkCode += j
	}
	checkCode *= 3
	for i := 0; i < len(barcode); i += 2 {
		j, _ := strconv.Atoi(barcode[i : i+1])
		checkCode += j
	}

	p.BarCode = barcode + strconv.Itoa(10-(checkCode%10))
	return true
}
