package main

import (
	"strconv"
	"strings"
	"time"
)

type Product struct {
	Id                      int32     `json:"id"`
	Name                    string    `json:"name"`
	Reference               string    `json:"reference"`
	BarCode                 string    `json:"barCode"`
	ControlStock            bool      `json:"controlStock"`
	Weight                  float64   `json:"weight"`
	Family                  *int32    `json:"family"`
	Width                   float64   `json:"width"`
	Height                  float64   `json:"height"`
	Depth                   float64   `json:"depth"`
	Off                     bool      `json:"off"`
	Stock                   int32     `json:"stock"`
	VatPercent              float64   `json:"vatPercent"`
	DateCreated             time.Time `json:"dateCreated"`
	Description             string    `json:"description"`
	Color                   *int32    `json:"color"`
	Price                   float64   `json:"price"`
	Manufacturing           bool      `json:"manufacturing"`
	ManufacturingOrderType  *int32    `json:"manufacturingOrderType"`
	Supplier                *int32    `json:"supplier"`
	FamilyName              *string   `json:"familyName"`
	MinimumStock            int32     `json:"minimumStock"`
	TrackMinimumStock       bool      `json:"trackMinimumStock"`
	prestaShopId            int32
	prestaShopCombinationId int32
	wooCommerceId           int32
	wooCommerceVariationId  int32
	shopifyId               int64
	shopifyVariantId        int64
	enterprise              int32
}

func getProduct(enterpriseId int32) []Product {
	var products []Product = make([]Product, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product_family WHERE product_family.id=product.family) FROM public.product WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := Product{}
		rows.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier, &p.prestaShopId, &p.prestaShopCombinationId, &p.MinimumStock, &p.TrackMinimumStock, &p.wooCommerceId, &p.wooCommerceVariationId, &p.shopifyId, &p.shopifyVariantId, &p.enterprise, &p.FamilyName)
		products = append(products, p)
	}

	return products
}

type ProductSearch struct {
	Search            string `json:"search"`
	TrackMinimumStock bool   `json:"trackMinimumStock"`
}

func (search *ProductSearch) searchProduct(enterpriseId int32) []Product {
	var products []Product = make([]Product, 0)
	sqlStatement := ""
	if search.TrackMinimumStock {
		sqlStatement = `SELECT *,(SELECT name FROM product_family WHERE product_family.id=product.family) FROM product WHERE (name ILIKE $1 AND track_minimum_stock=true) AND (enterprise=$2) ORDER BY id ASC`
	} else {
		sqlStatement = `SELECT *,(SELECT name FROM product_family WHERE product_family.id=product.family) FROM product WHERE (name ILIKE $1) AND (enterprise=$2) ORDER BY id ASC`
	}
	rows, err := db.Query(sqlStatement, "%"+search.Search+"%", enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := Product{}
		rows.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier, &p.prestaShopId, &p.prestaShopCombinationId, &p.MinimumStock, &p.TrackMinimumStock, &p.wooCommerceId, &p.wooCommerceVariationId, &p.shopifyId, &p.shopifyVariantId, &p.enterprise, &p.FamilyName)
		products = append(products, p)
	}

	return products
}

func getProductRow(productId int32) Product {
	sqlStatement := `SELECT * FROM public.product WHERE id=$1`
	row := db.QueryRow(sqlStatement, productId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Product{}
	}

	p := Product{}
	row.Scan(&p.Id, &p.Name, &p.Reference, &p.BarCode, &p.ControlStock, &p.Weight, &p.Family, &p.Width, &p.Height, &p.Depth, &p.Off, &p.Stock, &p.VatPercent, &p.DateCreated, &p.Description, &p.Color, &p.Price, &p.Manufacturing, &p.ManufacturingOrderType, &p.Supplier, &p.prestaShopId, &p.prestaShopCombinationId, &p.MinimumStock, &p.TrackMinimumStock, &p.wooCommerceId, &p.wooCommerceVariationId, &p.shopifyId, &p.shopifyVariantId, &p.enterprise)

	return p
}

func (p *Product) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 150 || len(p.Reference) > 40 || (len(p.BarCode) != 0 && len(p.BarCode) != 13) || p.VatPercent < 0 || p.Price < 0 || p.Weight < 0 || p.Width < 0 || p.Height < 0 || p.Depth < 0)
}

func (p *Product) insertProduct() bool {
	if !p.isValid() {
		return false
	}

	// Check that the format for EAN13 barcodes is correct
	if len(p.BarCode) == 13 {
		if !checkEan13(p.BarCode) {
			return false
		}
	}

	sqlStatement := `INSERT INTO public.product(name, reference, barcode, control_stock, weight, family, width, height, depth, off, stock, vat_percent, dsc, color, price, manufacturing, manufacturing_order_type, supplier, ps_id, ps_combination_id, minimum_stock, track_minimum_stock, wc_id, wc_variation_id, sy_id, sy_variant_id, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27) RETURNING id`
	row := db.QueryRow(sqlStatement, p.Name, p.Reference, &p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType, p.Supplier, p.prestaShopId, p.prestaShopCombinationId, p.MinimumStock, p.TrackMinimumStock, &p.wooCommerceId, &p.wooCommerceVariationId, p.shopifyId, p.shopifyVariantId, p.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var productId int32
	row.Scan(&productId)
	p.Id = productId

	return productId > 0
}

func (p *Product) updateProduct() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	// Check that the format for EAN13 barcodes is correct
	if len(p.BarCode) == 13 {
		if !checkEan13(p.BarCode) {
			return false
		}
	}

	sqlStatement := `UPDATE public.product SET name=$2, reference=$3, barcode=$4, control_stock=$5, weight=$6, family=$7, width=$8, height=$9, depth=$10, off=$11, stock=$12, vat_percent=$13, dsc=$14, color=$15, price=$16, manufacturing=$17, manufacturing_order_type=$18, supplier=$19, minimum_stock=$20, track_minimum_stock=$21 WHERE id=$1 AND enterprise=$22`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.Reference, p.BarCode, p.ControlStock, p.Weight, p.Family, p.Width, p.Height, p.Depth, p.Off, p.Stock, p.VatPercent, p.Description, p.Color, p.Price, p.Manufacturing, p.ManufacturingOrderType, p.Supplier, p.MinimumStock, p.TrackMinimumStock, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Product) deleteProduct() bool {
	if p.Id <= 0 {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	sqlStatement := `DELETE FROM stock WHERE product=$1 AND (SELECT enterprise FROM product WHERE product.id=stock.product)=$2`
	_, err = db.Exec(sqlStatement, p.Id, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `DELETE FROM public.product WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, p.Id, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
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

func findProductByName(productName string, enterpriseId int32) []NameInt32 {
	var products []NameInt32 = make([]NameInt32, 0)
	sqlStatement := `SELECT id,name FROM public.product WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 AND off=false ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(productName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := NameInt32{}
		rows.Scan(&p.Id, &p.Name)
		products = append(products, p)
	}

	return products
}

func getNameProduct(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.product WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type OrderDetailDefaults struct {
	Price      float64 `json:"price"`
	VatPercent float64 `json:"vatPercent"`
}

func getOrderDetailDefaults(productId int32, enterpriseId int32) OrderDetailDefaults {
	sqlStatement := `SELECT price, vat_percent FROM product WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, productId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return OrderDetailDefaults{}
	}
	s := OrderDetailDefaults{}
	row.Scan(&s.Price, &s.VatPercent)
	return s
}

// Get the sales order details with pending status, with the product specified.
func getProductSalesOrderDetailsPending(productId int32, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_order_detail.product) FROM sales_order_detail WHERE product=$1 AND quantity_delivery_note!=quantity AND enterprise=$2 ORDER BY sales_order_detail.id DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise, &d.ProductName)
		details = append(details, d)
	}

	return details
}

// Get the purchase order details with pending status, with the product specified.
func getProductPurchaseOrderDetailsPending(productId int32, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=purchase_order_detail.product) FROM purchase_order_detail WHERE product=$1 AND quantity_delivery_note!=quantity AND enterprise=$2 ORDER BY purchase_order_detail.id DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.ProductName)
		details = append(details, d)
	}

	return details
}

// Get the sales order details with the product specified.
func getProductSalesOrderDetails(productId int32, enterpriseId int32) []SalesOrderDetail {
	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=sales_order_detail.product) FROM sales_order_detail WHERE product=$1 AND enterprise=$2 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := SalesOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.Status, &d.QuantityPendingPackaging, &d.PurchaseOrderDetail, &d.prestaShopId, &d.Cancelled, &d.wooCommerceId, &d.shopifyId, &d.shopifyDraftId, &d.enterprise, &d.ProductName)
		details = append(details, d)
	}

	return details
}

// Get the purchase order details with the product specified.
func getProductPurchaseOrderDetails(productId int32, enterpriseId int32) []PurchaseOrderDetail {
	var details []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=purchase_order_detail.product) FROM purchase_order_detail WHERE product=$1 AND enterprise=$2 ORDER BY purchase_order_detail.id DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}
	for rows.Next() {
		d := PurchaseOrderDetail{}
		rows.Scan(&d.Id, &d.Order, &d.Product, &d.Price, &d.Quantity, &d.VatPercent, &d.TotalAmount, &d.QuantityInvoiced, &d.QuantityDeliveryNote, &d.QuantityPendingPackaging, &d.QuantityAssignedSale, &d.enterprise, &d.ProductName)
		details = append(details, d)
	}

	return details
}

// Get the warehouse movements with the product specified.
func getProductWarehouseMovement(productId int32, enterpriseId int32) []WarehouseMovement {
	var warehouseMovements []WarehouseMovement = make([]WarehouseMovement, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=warehouse_movement.product),(SELECT name FROM warehouse WHERE warehouse.id=warehouse_movement.warehouse AND warehouse.enterprise=warehouse_movement.enterprise) FROM warehouse_movement WHERE product=$1 AND enterprise=$2 ORDER BY warehouse_movement.id DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return warehouseMovements
	}
	for rows.Next() {
		m := WarehouseMovement{}
		rows.Scan(&m.Id, &m.Warehouse, &m.Product, &m.Quantity, &m.DateCreated, &m.Type, &m.SalesOrder, &m.SalesOrderDetail, &m.SalesInvoice, &m.SalesInvoiceDetail, &m.SalesDeliveryNote, &m.Description, &m.PurchaseOrder, &m.PurchaseOrderDetail, &m.PurchaseInvoice, &m.PurchaseInvoiceDetail, &m.PurchaseDeliveryNote, &m.DraggedStock, &m.Price, &m.VatPercent, &m.TotalAmount, &m.enterprise, &m.ProductName, &m.WarehouseName)
		warehouseMovements = append(warehouseMovements, m)
	}

	return warehouseMovements
}

// Get the manufacturing orders with the product specified.
func getProductManufacturingOrders(productId int32, enterpriseId int32) []ManufacturingOrder {
	manufacturingOrders := make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=manufacturing_order.type),(SELECT name FROM product WHERE product.id=manufacturing_order.product),(SELECT order_name FROM sales_order WHERE sales_order.id=manufacturing_order.order),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_tag_printed) FROM public.manufacturing_order WHERE product=$1 AND enterprise=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return manufacturingOrders
	}
	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.TypeName, &o.ProductName, &o.OrderName, &o.UserCreatedName, &o.UserManufacturedName, &o.UserTagPrintedName)
		manufacturingOrders = append(manufacturingOrders, o)
	}

	return manufacturingOrders
}

func (p *Product) generateBarcode(enterpriseId int32) bool {
	sqlStatement := `SELECT SUBSTRING(barcode,0,13) FROM product WHERE enterprise=$1 AND SUBSTRING(barcode,0,5)=$2 ORDER BY barcode DESC LIMIT 1`
	row := db.QueryRow(sqlStatement, enterpriseId, getSettingsRecordById(enterpriseId).BarcodePrefix)
	if row.Err() != nil {
		log("DB", row.Err().Error())
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

// Check if the bar code is a valid EAN13 product code. (Check the verification digit)
func checkEan13(barcode string) bool {
	// barcode must be a number
	_, err := strconv.Atoi(barcode)
	if err != nil {
		return false
	}

	// get the first 12 digits (remove the 13 character, which is the control digit), and reverse the string
	barcode12 := barcode[0:12]
	barcode12 = Reverse(barcode12)

	// add the numbers in the odd positions
	var controlNumber uint16
	for i := 0; i < len(barcode12); i += 2 {
		digit, _ := strconv.Atoi(string(barcode12[i]))
		controlNumber += uint16(digit)
	}

	// multiply by 3
	controlNumber *= 3

	// add the numbers in the pair positions
	for i := 1; i < len(barcode12); i += 2 {
		digit, _ := strconv.Atoi(string(barcode12[i]))
		controlNumber += uint16(digit)
	}

	// immediately higher ten
	var controlDigit uint16 = (10 - (controlNumber % 10)) % 10

	// check the control digits are the same
	inputControlDigit, _ := strconv.Atoi(string(barcode[12]))
	return controlDigit == uint16(inputControlDigit)
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type ProductImage struct {
	Id      int32  `json:"id"`
	Product int32  `json:"product"`
	URL     string `json:"url"`
}

func getProductImages(productId int32, enterpriseId int32) []ProductImage {
	var image []ProductImage = make([]ProductImage, 0)
	sqlStatement := `SELECT * FROM public.product_image WHERE product=$1 AND (SELECT enterprise FROM product WHERE product.id=product_image.product)=$2`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return image
	}
	for rows.Next() {
		d := ProductImage{}
		rows.Scan(&d.Id, &d.Product, &d.URL)
		image = append(image, d)
	}

	return image
}

func (i *ProductImage) isValid() bool {
	return !(len(i.URL) == 0 || len(i.URL) > 255)
}

func (i *ProductImage) insertProductImage(enterpriseId int32) bool {
	if !i.isValid() || i.Product <= 0 {
		return false
	}

	p := getProductRow(i.Product)
	if p.enterprise != enterpriseId {
		return false
	}

	sqlStatement := `INSERT INTO public.product_image(product, url) VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, i.Product, i.URL)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (i *ProductImage) updateProductImage(enterpriseId int32) bool {
	if i.Id <= 0 || !i.isValid() {
		return false
	}

	sqlStatement := `SELECT product FROM public.product_image WHERE id=$1`
	row := db.QueryRow(sqlStatement, i.Id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var productId int32
	row.Scan(&productId)
	p := getProductRow(productId)
	if p.enterprise != enterpriseId {
		return false
	}

	sqlStatement = `UPDATE public.product_image SET url=$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, i.Id, i.URL)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (i *ProductImage) deleteProductImage(enterpriseId int32) bool {
	if i.Id <= 0 {
		return false
	}

	sqlStatement := `SELECT product FROM public.product_image WHERE id=$1`
	row := db.QueryRow(sqlStatement, i.Id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var productId int32
	row.Scan(&productId)
	p := getProductRow(productId)
	if p.enterprise != enterpriseId {
		return false
	}

	sqlStatement = `DELETE FROM public.product_image WHERE id=$1`
	res, err := db.Exec(sqlStatement, i.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func calculateMinimumStock(enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)
	t := time.Now()
	if s.MinimumStockSalesPeriods <= 0 || s.MinimumStockSalesDays <= 0 {
		return false
	}
	t = t.AddDate(0, 0, -int(s.MinimumStockSalesDays))

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	sqlStatement := `SELECT id FROM product WHERE track_minimum_stock=true AND enterprise=$1 AND off=false`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	for rows.Next() {
		var productId int32
		rows.Scan(&productId)

		sqlStatement := `SELECT SUM(sales_order_detail.quantity) FROM sales_order_detail INNER JOIN sales_order ON sales_order.id=sales_order_detail.order WHERE sales_order_detail.product=$1 AND sales_order.date_created >= $2`
		row := db.QueryRow(sqlStatement, productId, t)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			trans.Rollback()
			return false
		}

		var quantitySold int32
		row.Scan(&quantitySold)

		sqlStatement = `UPDATE product SET minimum_stock=$2 WHERE id=$1`
		_, err := db.Exec(sqlStatement, productId, quantitySold/int32(s.MinimumStockSalesPeriods))
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
	}

	///
	err = trans.Commit()
	return err == nil
	///
}

func generateManufacturingOrPurchaseOrdersMinimumStock(userId int32, enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)
	var generadedPurchaseOrders map[int32]PurchaseOrder = make(map[int32]PurchaseOrder) // Key: supplier ID, Value: generated purchase order

	sqlStatement := `SELECT product.id,stock.quantity_available,product.minimum_stock,product.manufacturing,product.manufacturing_order_type,product.supplier FROM product INNER JOIN stock ON stock.product=product.id WHERE product.track_minimum_stock=true AND stock.quantity_available < (product.minimum_stock*2) AND product.enterprise=$1 AND off=false`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	// iterate over the list of product that don't have covered the minimum stock
	for rows.Next() {
		var productId int32
		var quantityAvailable int32
		var minimumStock int32
		var manufacturing bool
		var manufacturingOrderType *int32
		var supplier *int32
		rows.Scan(&productId, &quantityAvailable, &minimumStock, &manufacturing, &manufacturingOrderType, &supplier)

		if manufacturing { // if the product is from manufacture, generate the manufacturing orders
			// generate manufacturing order or purchase orders until the available quantity is equal to the minimum stock * 2
			for i := quantityAvailable; i < (minimumStock * 2); i++ {

				o := ManufacturingOrder{Product: productId, Type: *manufacturingOrderType}
				o.UserCreated = userId
				o.enterprise = enterpriseId
				ok := o.insertManufacturingOrder()
				if !ok {
					trans.Rollback()
					return false
				}
			}
		} else { // if the product is not from manufacture, generate the purchase order to the supplier
			o, ok := generadedPurchaseOrders[*supplier]
			if !ok { // there is no purchase order generated for this supplier, create it and add to the map
				d := getSupplierDefaults(*supplier, enterpriseId)
				if d.BillingSeries == nil || d.Currency == nil || d.MainBillingAddress == nil || d.MainShippingAddress == nil || d.PaymentMethod == nil {
					continue
				}
				p := PurchaseOrder{}
				p.Warehouse = s.DefaultWarehouse
				p.Supplier = *supplier
				p.BillingSeries = *d.BillingSeries
				p.Currency = *d.Currency
				p.BillingAddress = *d.MainBillingAddress
				p.ShippingAddress = *d.MainShippingAddress
				p.PaymentMethod = *d.PaymentMethod

				p.enterprise = enterpriseId
				ok, purchaseOrderId := p.insertPurchaseOrder()
				if !ok {
					trans.Rollback()
					return false
				}
				p.Id = purchaseOrderId
				generadedPurchaseOrders[*supplier] = p

				// generate the needs as a detail
				product := getProductRow(productId)
				det := PurchaseOrderDetail{Order: p.Id, Product: productId, Quantity: (minimumStock * 2) - quantityAvailable, Price: product.Price, VatPercent: product.VatPercent, enterprise: enterpriseId}
				ok, _ = det.insertPurchaseOrderDetail(false)
				if !ok {
					trans.Rollback()
					return false
				}
			} else { // it already exists a purchase order for this supplier, add the needs as details
				// generate the needs as a detail
				product := getProductRow(productId)
				det := PurchaseOrderDetail{Order: o.Id, Product: productId, Quantity: (minimumStock * 2) - quantityAvailable, Price: product.Price, VatPercent: product.VatPercent, enterprise: enterpriseId}
				ok, _ = det.insertPurchaseOrderDetail(false)
				if !ok {
					trans.Rollback()
					return false
				}
			}
		}
	}

	///
	err = trans.Commit()
	return err == nil
	///
}

type ProductLocate struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

type ProductLocateQuery struct {
	Mode  int32  `json:"mode"` // 0 = ID, 1 = Name, 2 = Reference
	Value string `json:"value"`
}

func (q *ProductLocateQuery) locateProduct(enterpriseId int32) []ProductLocate {
	var products []ProductLocate = make([]ProductLocate, 0)
	sqlStatement := ``
	parameters := make([]interface{}, 0)
	if q.Value == "" {
		sqlStatement = `SELECT id,name,reference FROM product WHERE enterprise=$1 AND off=false ORDER BY id ASC`
		parameters = append(parameters, enterpriseId)
	} else if q.Mode == 0 {
		id, err := strconv.Atoi(q.Value)
		if err != nil {
			sqlStatement = `SELECT id,name,reference FROM product WHERE enterprise=$1 AND off=false ORDER BY id ASC`
			parameters = append(parameters, enterpriseId)
		} else {
			sqlStatement = `SELECT id,name,reference FROM product WHERE id=$1 AND enterprise=$2 AND off=false`
			parameters = append(parameters, id)
			parameters = append(parameters, enterpriseId)
		}
	} else if q.Mode == 1 {
		sqlStatement = `SELECT id,name,reference FROM product WHERE name ILIKE $1 AND enterprise=$2 AND off=false ORDER BY id ASC`
		parameters = append(parameters, "%"+q.Value+"%")
		parameters = append(parameters, enterpriseId)
	} else if q.Mode == 2 {
		sqlStatement = `SELECT id,name,reference FROM product WHERE reference ILIKE $1 AND enterprise=$2 AND off=false ORDER BY id ASC`
		parameters = append(parameters, "%"+q.Value+"%")
		parameters = append(parameters, enterpriseId)
	}
	rows, err := db.Query(sqlStatement, parameters...)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := ProductLocate{}
		rows.Scan(&p.Id, &p.Name, &p.Reference)
		products = append(products, p)
	}

	return products
}

/* PRODUCT GENERATOR */

type ProductGenerator struct {
	Products                   []ProductGenerate `json:"products"`
	ManufacturingOrderTypeMode int16             `json:"manufacturingOrderTypeMode"` // 0 = No manufacturing, 1 = Create a manufacturing order type with the same name, 2 = Use the manufacturing order name
	ManufacturingOrderTypeName *string           `json:"manufacturingOrderTypeName"`
}

type ProductGenerate struct {
	Name            string  `json:"name"`
	Reference       string  `json:"reference"`
	GenerateBarCode bool    `json:"generateBarCode"`
	BarCode         string  `json:"barCode"`
	Weight          float64 `json:"weight"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	Depth           float64 `json:"depth"`
	Price           float64 `json:"price"`
	Manufacturing   bool    `json:"manufacturing"`
	InitialStock    int32   `json:"initialStock"`
}

func (g *ProductGenerator) productGenerator(enterpriseId int32) bool {
	var manufacturingOrderTypeId int32
	if g.ManufacturingOrderTypeMode == 2 {
		if g.ManufacturingOrderTypeName == nil {
			return false
		}

		mot := ManufacturingOrderType{
			Name:       *g.ManufacturingOrderTypeName,
			enterprise: enterpriseId,
		}
		mot.insertManufacturingOrderType()
		manufacturingOrderTypeId = mot.Id
	}

	for i := 0; i < len(g.Products); i++ {
		product := g.Products[i]

		p := Product{
			Name:          product.Name,
			Reference:     product.Reference,
			BarCode:       product.BarCode,
			Weight:        product.Weight,
			Width:         product.Width,
			Height:        product.Height,
			Depth:         product.Depth,
			Price:         product.Price,
			Manufacturing: product.Manufacturing,
			enterprise:    enterpriseId,
		}

		if product.Manufacturing && g.ManufacturingOrderTypeMode == 1 {
			mot := ManufacturingOrderType{
				Name:       product.Name,
				enterprise: enterpriseId,
			}
			mot.insertManufacturingOrderType()
			p.ManufacturingOrderType = &mot.Id
		} else if product.Manufacturing && g.ManufacturingOrderTypeMode == 2 {
			p.ManufacturingOrderType = &manufacturingOrderTypeId
		}

		ok := p.insertProduct()
		if !ok {
			return false
		}

		if product.GenerateBarCode {
			p := getProductRow(p.Id)
			p.generateBarcode(enterpriseId)
			p.updateProduct()
		}

		if product.InitialStock != 0 {
			s := getSettingsRecordById(enterpriseId)
			wm := WarehouseMovement{
				Warehouse:  s.DefaultWarehouse,
				Product:    p.Id,
				Quantity:   product.InitialStock,
				Type:       "R",
				enterprise: enterpriseId,
			}
			wm.insertWarehouseMovement()
		}
	}
	return true
}
