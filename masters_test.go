package main

import (
	"strings"
	"testing"
	"time"
)

// ===== CUSTOMERS

/* GET */

func TestGetCustomers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}
	c := q.getCustomers()

	for i := 0; i < len(c.Customers); i++ {
		if c.Customers[i].Id <= 0 {
			t.Error("Scan error, customers with ID 0.")
			return
		}
	}
}

func TestSearchCustomers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}}
	c := q.searchCustomers()

	for i := 0; i < len(c.Customers); i++ {
		if c.Customers[i].Id <= 0 {
			t.Error("Scan error, customers with ID 0.")
			return
		}
	}
}

func TestGetCustomerRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getCustomerRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, customer row with ID 0.")
		return
	}

}

func TestFindCustomerByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := findCustomerByName("", 1)
	if len(c) == 0 {
		t.Error("Can't find customers by name")
	}

	for i := 0; i < len(c); i++ {
		if c[i].Id <= 0 {
			t.Error("Scan error, fund customer by name with ID 0.")
			return
		}
	}
}

func TestGetNameCustomer(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	customerName := getNameCustomer(1, 1)
	if customerName == "" {
		t.Error("Can't get the name of the customers")
	}
}

/* INSERT - UPDATE - DELETE */

func TestCustomerInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	country := int32(55)
	language := int32(8)
	paymentMethod := int32(3)
	billingSeries := "EXP"

	c := Customer{
		Name:          "Jake Kaiser",
		Tradename:     "Jake Kaiser",
		FiscalName:    "Jake Kaiser",
		Phone:         "679681745",
		Email:         "jake.kaiser@gmail.com",
		Country:       &country,
		Language:      &language,
		PaymentMethod: &paymentMethod,
		BillingSeries: &billingSeries,
		enterprise:    1,
	}

	// insert
	res := c.insertCustomer(0)
	ok, customerId := res.Id > 0, res.Id
	if !ok {
		t.Error("Insert error, customer not inserted")
		return
	}
	c.Id = int32(customerId)

	c.TaxId = "ABCDEF1234"
	ok = c.updateCustomer(0)
	if !ok {
		t.Error("Update error, customer not updated")
		return
	}

	// check update
	c = getCustomerRow(c.Id)
	if c.TaxId != "ABCDEF1234" {
		t.Error("Update error, customer not successfully updated")
		return
	}

	// check defaults
	defaults := getCustomerDefaults(c.Id, 1)
	if (defaults.PaymentMethod == nil || *defaults.PaymentMethod != paymentMethod) || (defaults.BillingSeriesName == nil || len(*defaults.BillingSeriesName) == 0) || (defaults.BillingSeries == nil || *defaults.BillingSeries != billingSeries) {
		t.Error("Customer defaults are not correct")
		return
	}

	// delete
	ok = c.deleteCustomer(0)
	if !ok {
		t.Error("Delete error, customer not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestGetCustomerAddresses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}
	c := q.getCustomers()

	for i := 0; i < len(c.Customers); i++ {
		addresses := getCustomerAddresses(c.Customers[i].Id, 1)
		if len(addresses) > 0 {
			if addresses[0].Id <= 0 {
				t.Error("Customer addresses not scanned successfully.")
				return
			} else {
				return
			}
		}
	}
}

func TestGetCustomerSaleOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}
	c := q.getCustomers()

	for i := 0; i < len(c.Customers); i++ {
		orders := getCustomerSaleOrders(c.Customers[i].Id, 1)
		if len(orders) > 0 {
			if orders[0].Id <= 0 {
				t.Error("Customer sale orders not scanned successfully.")
				return
			} else {
				return
			}
		}
	}
}

func TestSetCustomerAccount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	country := int32(1)
	language := int32(8)
	paymentMethod := int32(3)
	billingSeries := "EXP"

	c := Customer{
		Name:          "Jake Kaiser",
		Tradename:     "Jake Kaiser",
		FiscalName:    "Jake Kaiser",
		Phone:         "679681745",
		Email:         "jake.kaiser@gmail.com",
		Country:       &country,
		Language:      &language,
		PaymentMethod: &paymentMethod,
		BillingSeries: &billingSeries,
		enterprise:    1,
	}

	// insert
	res := c.insertCustomer(0)
	ok, customerId := res.Id > 0, res.Id
	if !ok {
		t.Error("Insert error, customer not inserted")
		return
	}

	// update
	c.Id = int32(customerId)

	c.setCustomerAccount()

	c = getCustomerRow(c.Id)
	if c.Account == nil {
		t.Error("Customer account not set")
		return
	}

	ok = c.deleteCustomer(0)
	if !ok {
		t.Error("Customer not deleted")
	}
}

func TestLocateCustomers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := CustomerLocateQuery{Mode: 0, Value: ""}
	customers := q.locateCustomers(1)

	if len(customers) == 0 {
		t.Error("Customers can't be empty")
		return
	}

	for i := 0; i < len(customers); i++ {
		if customers[i].Id <= 0 {
			t.Error("Locate customers not scanned")
			return
		}
	}

	q = CustomerLocateQuery{Mode: 1, Value: ""}
	customers = q.locateCustomers(1)

	if len(customers) == 0 {
		t.Error("Customers can't be empty")
		return
	}

	for i := 0; i < len(customers); i++ {
		if customers[i].Id <= 0 {
			t.Error("Locate customers not scanned")
			return
		}
	}

	q = CustomerLocateQuery{Mode: 0, Value: "1"}
	customers = q.locateCustomers(1)

	if len(customers) == 0 {
		t.Error("Customers can't be empty")
		return
	}

	for i := 0; i < len(customers); i++ {
		if customers[i].Id <= 0 {
			t.Error("Locate customers not scanned")
			return
		}
	}
}

// ===== SUPPLIERS

/* GET */

func TestGetSuppliers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSuppliers(1)

	for i := 0; i < len(s); i++ {
		if s[i].Id <= 0 {
			t.Error("Scan error, suppliers with ID 0.")
			return
		}
	}
}

func TestSearchSuppliers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	s := searchSuppliers("", 1)

	for i := 0; i < len(s); i++ {
		if s[i].Id <= 0 {
			t.Error("Scan error, suppliers with ID 0.")
			return
		}
	}
}

func TestGetSuppliersRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getSupplierRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, supplier row with ID 0.")
		return
	}

}

func TestFindSupplierByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := findSupplierByName("", 1)
	if len(c) == 0 {
		t.Error("Can't find suppliers by name")
	}

	for i := 0; i < len(c); i++ {
		if c[i].Id <= 0 {
			t.Error("Scan error, fund supplier by name with ID 0.")
			return
		}
	}
}

func TestGetNameSupplier(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	supplierName := getNameSupplier(1, 1)
	if supplierName == "" {
		t.Error("Can't get the name of the suppliers")
	}
}

/* INSERT - UPDATE - DELETE */

func TestSupplierInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	country := int32(55)
	language := int32(8)
	paymentMethod := int32(3)
	billingSeries := "EXP"

	c := Supplier{
		Name:          "Jake Kaiser",
		Tradename:     "Jake Kaiser",
		FiscalName:    "Jake Kaiser",
		Phone:         "679681745",
		Email:         "jake.kaiser@gmail.com",
		Country:       &country,
		Language:      &language,
		PaymentMethod: &paymentMethod,
		BillingSeries: &billingSeries,
		enterprise:    1,
	}

	// insert
	ok := c.insertSupplier(0).Id > 0
	if !ok {
		t.Error("Insert error, supplier not inserted")
		return
	}

	// update
	suppliers := getSuppliers(1)
	c = suppliers[len(suppliers)-1]

	c.TaxId = "ABCDEF1234"
	ok = c.updateSupplier(0)
	if !ok {
		t.Error("Update error, supplier not updated")
		return
	}

	// check update
	c = getSupplierRow(c.Id)
	if c.TaxId != "ABCDEF1234" {
		t.Error("Update error, supplier not successfully updated")
		return
	}

	// check defaults
	defaults := getSupplierDefaults(c.Id, 1)
	if (defaults.PaymentMethod == nil || *defaults.PaymentMethod != paymentMethod) || (defaults.BillingSeriesName == nil || len(*defaults.BillingSeriesName) == 0) || (defaults.BillingSeries == nil || *defaults.BillingSeries != billingSeries) || (defaults.BillingSeriesName == nil || len(*defaults.BillingSeriesName) == 0) {
		t.Error("Supplier defaults are not correct")
		return
	}

	// delete
	ok = c.deleteSupplier(0)
	if !ok {
		t.Error("Delete error, supplier not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestGetSupplierAddresses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSuppliers(1)

	for i := 0; i < len(s); i++ {
		addresses := getSupplierAddresses(s[i].Id, 1)
		if len(addresses) > 0 {
			if addresses[0].Id <= 0 {
				t.Error("Supplier addresses not scanned successfully")
				return
			} else {
				return
			}
		}
	}
}

func TestGetSupplierSaleOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSuppliers(1)

	for i := 0; i < len(s); i++ {
		orders := getSupplierPurchaseOrders(s[i].Id, 1)
		if len(orders) > 0 {
			if orders[0].Id <= 0 {
				t.Error("Supplier purchase orders not scanned successfully")
				return
			} else {
				return
			}
		}
	}
}

func TestSetSupplierAccount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	country := int32(1)
	language := int32(8)
	paymentMethod := int32(3)
	billingSeries := "EXP"

	s := Supplier{
		Name:          "Jake Kaiser",
		Tradename:     "Jake Kaiser",
		FiscalName:    "Jake Kaiser",
		Phone:         "679681745",
		Email:         "jake.kaiser@gmail.com",
		Country:       &country,
		Language:      &language,
		PaymentMethod: &paymentMethod,
		BillingSeries: &billingSeries,
		enterprise:    1,
	}

	// insert
	ok := s.insertSupplier(0).Id > 0
	if !ok {
		t.Error("Insert error, supplier not inserted")
		return
	}

	// update
	suppliers := getSuppliers(1)
	s = suppliers[len(suppliers)-1]

	s.setSupplierAccount()

	s = getSupplierRow(s.Id)
	if s.Account == nil {
		t.Error("Supplier account not set")
		return
	}

	s.deleteSupplier(0)
}

func TestLocateSuppliers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := SupplierLocateQuery{Mode: 0, Value: ""}
	suppliers := q.locateSuppliers(1)

	if len(suppliers) == 0 {
		t.Error("Suppliers can't be empty")
		return
	}

	for i := 0; i < len(suppliers); i++ {
		if suppliers[i].Id <= 0 {
			t.Error("Locate suppliers not scanned")
			return
		}
	}

	q = SupplierLocateQuery{Mode: 1, Value: ""}
	suppliers = q.locateSuppliers(1)

	if len(suppliers) == 0 {
		t.Error("Suppliers can't be empty")
		return
	}

	for i := 0; i < len(suppliers); i++ {
		if suppliers[i].Id <= 0 {
			t.Error("Locate suppliers not scanned")
			return
		}
	}

	q = SupplierLocateQuery{Mode: 0, Value: "1"}
	suppliers = q.locateSuppliers(1)

	if len(suppliers) == 0 {
		t.Error("Suppliers can't be empty")
		return
	}

	for i := 0; i < len(suppliers); i++ {
		if suppliers[i].Id <= 0 {
			t.Error("Locate suppliers not scanned")
			return
		}
	}
}

// ===== PRODUCTS

/* GET */

func TestGetProduct(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	p := getProduct(1)

	for i := 0; i < len(p); i++ {
		if p[i].Id <= 0 {
			t.Error("Scan error, products with ID 0.")
			return
		}
	}
}

func TestSearchProduct(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// search all
	q := ProductSearch{Search: "", TrackMinimumStock: false}
	p := q.searchProduct(1)

	for i := 0; i < len(p); i++ {
		if p[i].Id <= 0 {
			t.Error("Scan error, products with ID 0.")
			return
		}
	}

	// search track minimum stock
	q = ProductSearch{Search: "", TrackMinimumStock: true}
	p = q.searchProduct(1)

	for i := 0; i < len(p); i++ {
		if p[i].Id <= 0 {
			t.Error("Scan error, products with ID 0.")
			return
		}
	}
}

func TestGetProductRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := getProductRow(1)
	if o.Id <= 0 {
		t.Error("Scan error, product row with ID 0.")
		return
	}

}

func TestGetNameProduct(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	productName := getNameProduct(1, 1)
	if productName == "" {
		t.Error("Could not get the name of the product")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestProductInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(2)
	supplier := int32(1)

	p := Product{
		Name:                   "Glass Office Desk",
		Reference:              "OF-DSK",
		BarCode:                "1234067891236",
		ControlStock:           true,
		Weight:                 30,
		Family:                 &family,
		Width:                  160,
		Height:                 100,
		Depth:                  40,
		VatPercent:             21,
		Price:                  65,
		Manufacturing:          true,
		ManufacturingOrderType: &manufacturingOrderType,
		Supplier:               &supplier,
		TrackMinimumStock:      true,
		prestaShopId:           1,
		enterprise:             1,
	}

	ok := p.insertProduct(0).Ok
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	p.Name = "Wooden Office Desk"
	ok = p.updateProduct(0).Ok
	if !ok {
		t.Error("Update error, could not update product")
		return
	}

	p = getProductRow(p.Id)
	if p.Name != "Wooden Office Desk" {
		t.Error("Update error, product update not successful")
		return
	}

	ok = p.deleteProduct(0).Ok
	if !ok {
		t.Error("Delete error, could not delete product")
		return
	}
}

/* FUNCTIONALITY */

func TestGetOrderDetailDefaults(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	detauls := getOrderDetailDefaults(1, 1)
	if detauls.Price == 0 || detauls.VatPercent == 0 {
		t.Error("Order details defaults lot loaded")
		return
	}
}

func TestProductRelations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}
	products := getProduct(1)

	for i := 0; i < len(products); i++ {
		saleOrders := getProductSalesOrderDetailsPending(products[i].Id, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Sale orders with ID 0 on product")
				return
			} else {
				break
			}
		}
	}

	for i := 0; i < len(products); i++ {
		saleOrders := getProductPurchaseOrderDetailsPending(products[i].Id, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Purchase orders with ID 0 on product")
				return
			} else {
				break
			}
		}
	}

	for i := 0; i < len(products); i++ {
		saleOrders := getProductSalesOrderDetails(ProductSalesOrderDetailsQuery{ProductId: products[i].Id}, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Sale orders with ID 0 on product")
				return
			} else {
				break
			}
		}
	}

	for i := 0; i < len(products); i++ {
		saleOrders := getProductPurchaseOrderDetails(ProductPurchaseOrderDetailsQuery{ProductId: products[i].Id}, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Purchase orders with ID 0 on product")
				return
			} else {
				break
			}
		}
	}

	for i := 0; i < len(products); i++ {
		saleOrders := getProductWarehouseMovement(ProductPurchaseOrderDetailsQuery{ProductId: products[i].Id}, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Warehouse movement with ID 0 on product.")
				return
			} else {
				break
			}
		}
	}

	for i := 0; i < len(products); i++ {
		saleOrders := getProductManufacturingOrders(ProductManufacturingOrdersQuery{ProductId: products[i].Id}, 1)
		if len(saleOrders) > 0 {
			if saleOrders[0].Id <= 0 {
				t.Error("Manufacturing orders with ID 0 on product.")
				return
			} else {
				break
			}
		}
	}
}

func TestGenerateBarcode(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := int32(1)
	manufacturingOrderType := int32(2)
	supplier := int32(1)

	p := Product{
		Name:                   "Glass Office Desk",
		Reference:              "OF-DSK",
		BarCode:                "1234067891236",
		ControlStock:           true,
		Weight:                 30,
		Family:                 &family,
		Width:                  160,
		Height:                 100,
		Depth:                  40,
		VatPercent:             21,
		Price:                  65,
		Manufacturing:          true,
		ManufacturingOrderType: &manufacturingOrderType,
		Supplier:               &supplier,
		TrackMinimumStock:      true,
		enterprise:             1,
	}

	ok := p.insertProduct(0).Ok
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	products := getProduct(1)
	p = products[len(products)-1]

	ok = p.generateBarcode(1)
	if !ok {
		t.Error("Error generating product barcode")
		return
	}

	p = getProductRow(p.Id)
	if len(p.BarCode) == 0 {
		t.Error("EAN13 barcode not generated")
		return
	}

	ok = p.deleteProduct(0).Ok
	if !ok {
		t.Error("Delete error, could not delete product")
		return
	}
}

func TestGetProductImages(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	products := getProduct(1)

	for i := 0; i < len(products); i++ {
		images := getProductImages(products[i].Id, 1)
		if len(images) > 0 {
			if images[0].Id <= 0 {
				t.Error("Product images not scanned successfully")
				return
			} else {
				return
			}
		}
	}
}

func TestProductImageInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	pi := ProductImage{
		Product: 1,
		URL:     "http://example.domain/picture.png",
	}

	ok := pi.insertProductImage(1)
	if !ok {
		t.Error("Product image not inserted")
		return
	}

	images := getProductImages(1, 1)
	pi = images[len(images)-1]

	pi.URL = "https://example.domain/picture.png"
	pi.updateProductImage(1)
	if !ok {
		t.Error("Product image not updated")
		return
	}

	pi.deleteProductImage(1)
	if !ok {
		t.Error("Product image not deleted")
		return
	}
}

func TestCalculateMinimumStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	ok := calculateMinimumStock(1, 0)
	if !ok {
		t.Error("Calculate minimum stock not successful")
	}
}

func TestGenerateManufacturingOrPurchaseOrdersMinimumStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	ok := generateManufacturingOrPurchaseOrdersMinimumStock(1, 1)
	if !ok {
		t.Error("Could not generate manufacturing or purchase orders to cover minumum stock")
		return
	}
}

func TestLocateProduct(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := ProductLocateQuery{Mode: 0, Value: ""}
	produts := q.locateProduct(1)

	if len(produts) == 0 || produts[0].Id <= 0 {
		t.Error("Could not scan products")
		return
	}

	for i := 0; i < len(produts); i++ {
		p := getProductRow(produts[i].Id)
		if p.Off {
			t.Error("The locator is showing off products")
			return
		}
	}

	q = ProductLocateQuery{Mode: 0, Value: "1"}
	produts = q.locateProduct(1)

	if len(produts) == 0 || produts[0].Id <= 0 {
		t.Error("Could not scan products")
		return
	}

	q = ProductLocateQuery{Mode: 1, Value: ""}
	produts = q.locateProduct(1)

	if len(produts) == 0 || produts[0].Id <= 0 {
		t.Error("Could not scan products")
		return
	}

	q = ProductLocateQuery{Mode: 2, Value: ""}
	produts = q.locateProduct(1)

	if len(produts) == 0 || produts[0].Id <= 0 {
		t.Error("Could not scan products")
		return
	}
}

func TestProductOffSaleOrderDetailInsert(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder(1)

	p := Product{
		Name:              "Glass Office Desk",
		Reference:         "OF-DSK",
		BarCode:           "1234067891236",
		ControlStock:      true,
		Weight:            30,
		Width:             160,
		Height:            100,
		Depth:             40,
		VatPercent:        21,
		Price:             65,
		Manufacturing:     false,
		TrackMinimumStock: true,
		prestaShopId:      1,
		enterprise:        1,
		Off:               true,
	}

	ok := p.insertProduct(0).Ok
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    p.Id,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	// test insert
	ok = d.insertSalesOrderDetail(0).Ok
	if ok {
		t.Error("Insert error, sale order detail inserted with an off product")
		return
	}

	ok = p.deleteProduct(1).Ok
	if !ok {
		t.Error("Delete error, can't delete the temp product")
		return
	}

	o.Id = orderId
	o.enterprise = 1
	ok = o.deleteSalesOrder(1).Ok
	if !ok {
		t.Error("Delete error, can't delete the temp sale order")
		return
	}
}

func TestCheckEan13(t *testing.T) {
	if !checkEan13("1234567890418") {
		t.Error("The EAN13 should be OK")
	}

	if checkEan13("1234567890417") {
		t.Error("The EAN13 should not be OK")
	}

	if checkEan13("123456789041") {
		t.Error("The EAN13 should not be OK")
	}
}

func TestSaleOrderDetailNoStock(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	o := SaleOrder{
		Warehouse:       "W1",
		Customer:        1,
		PaymentMethod:   3,
		BillingSeries:   "EXP",
		Currency:        1,
		BillingAddress:  1,
		ShippingAddress: 1,
		Description:     "",
		Notes:           "",
		enterprise:      1,
	}

	_, orderId := o.insertSalesOrder(1)

	p := Product{
		Name:              "Glass Office Desk",
		Reference:         "OF-DSK",
		BarCode:           "1234067891236",
		ControlStock:      false,
		Weight:            30,
		Width:             160,
		Height:            100,
		Depth:             40,
		VatPercent:        21,
		Price:             65,
		Manufacturing:     true,
		TrackMinimumStock: true,
		prestaShopId:      1,
		enterprise:        1,
		DigitalProduct:    true,
	}

	ok := p.insertProduct(0).Ok
	if !ok {
		t.Error("Insert error, could not insert product")
		return
	}

	d := SalesOrderDetail{
		Order:      orderId,
		Product:    p.Id,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
		enterprise: 1,
	}

	// test insert
	ok = d.insertSalesOrderDetail(0).Ok
	if !ok {
		t.Error("Insert error, sale order detail not inserted")
		return
	}

	invoiceAllSaleOrder(orderId, 1, 0)

	// check if the order has the correct status
	details := getSalesOrderDetail(orderId, 1)
	if details[0].Status != "E" {
		t.Error("The status is not 'E' for no control stock")
		return
	}

	// CRUD digital product data
	dpd := SalesOrderDetailDigitalProductData{
		Detail: d.Id,
		Key:    "test",
		Value:  "1234-5678-90",
	}
	dpd.insertSalesOrderDetailDigitalProductData(1)

	dgDatas := getSalesOrderDetailDigitalProductData(d.Id, 1)
	if len(dgDatas) == 0 || dgDatas[0].Id <= 0 {
		t.Error("Can't scan digital product data")
		return
	}

	dgDatas[0].Key = "license_code"
	dgDatas[0].updateSalesOrderDetailDigitalProductData(1)

	dgDatas = getSalesOrderDetailDigitalProductData(d.Id, 1)
	if dgDatas[0].Key != "license_code" {
		t.Error("Digital product data not updated")
		return
	}

	// set as sent
	setAsSent := SetDigitalSalesOrderDetailAsSent{
		Detail:                 d.Id,
		SendEmail:              false,
		DestinationAddress:     "admin@marketneterp.io",
		DestinationAddressName: "admin@marketneterp.io",
		Subject:                "UNIT TEST",
	}
	setAsSent.setDigitalSalesOrderDetailAsSent(1)

	// check if the order has the correct status
	details = getSalesOrderDetail(orderId, 1)
	if details[0].Status != "G" {
		t.Error("The status is not 'G' for digital product sent")
		return
	}

	// attempt delete
	sqlStatement := `UPDATE sales_order_detail SET status='E' WHERE id=$1`
	_, err := db.Exec(sqlStatement, d.Id)
	if err != nil {
		log("DB", err.Error())
	}

	dgDatas[0].deleteSalesOrderDetailDigitalProductData(1)
	dgDatas = getSalesOrderDetailDigitalProductData(d.Id, 1)
	if len(dgDatas) != 0 {
		t.Error("Can't delete digital product data")
		return
	}

	inv := getSalesOrderRelations(orderId, 1).Invoices[0]
	inv.deleteSalesInvoice(0)

	details[0].enterprise = 1
	ok = details[0].deleteSalesOrderDetail(1, nil).Ok
	if !ok {
		t.Error("Delete error, sale order detail not deleted")
		return
	}

	// check that the sale order has been updated correctly
	inMemoryOrder := getSalesOrderRow(orderId)
	if inMemoryOrder.LinesNumber != 0 {
		t.Error("The sale order number of lines is not updated upon delete")
		return
	}

	o.Id = orderId
	o.enterprise = 1
	o.deleteSalesOrder(0)

	p.enterprise = 1
	p.deleteProduct(0)
}

func TestProductGenerator(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	pg := ProductGenerator{
		Products: []ProductGenerate{
			{
				Name:            "Test product generator",
				Reference:       "TST_PROD_GEN",
				GenerateBarCode: true,
				Weight:          10,
				Width:           10,
				Height:          10,
				Depth:           10,
				Price:           100,
				Manufacturing:   true,
				InitialStock:    12,
			},
		},
		ManufacturingOrderTypeMode: 1,
	}
	ok := pg.productGenerator(1, 0)
	if !ok {
		t.Error("The product generador returned an error")
		return
	}

	products := getProduct(1)
	p := products[len(products)-1]
	if p.Reference != "TST_PROD_GEN" {
		t.Error("Product not generated")
		return
	}

	if p.Name != "Test product generator" {
		t.Error("Incorrect name")
		return
	}

	if p.Reference != "TST_PROD_GEN" {
		t.Error("Incorrect reference")
		return
	}

	if len(strings.Trim(p.BarCode, "")) == 0 {
		t.Error("Barcode not generated")
		return
	}

	if p.Weight != 10 || p.Width != 10 || p.Height != 10 || p.Depth != 10 {
		t.Error("Weight and dimensions not set")
		return
	}

	if p.Price != 100 {
		t.Error("Price not set")
		return
	}

	if !p.Manufacturing || p.ManufacturingOrderType == nil || *p.ManufacturingOrderType <= 0 {
		t.Error("Manufacturing information not set")
		return
	}

	mot := getManufacturingOrderTypeRow(*p.ManufacturingOrderType)
	if mot.Id <= 0 {
		t.Error("Can't scan manufacturing order created")
		return
	}

	if mot.Name != "Test product generator" {
		t.Error("Manufacturing order name not set")
		return
	}

	if mot.QuantityManufactured <= 0 {
		t.Error("Quantity manufactured incorrect")
		return
	}

	p.deleteProduct(0)

	mot.deleteManufacturingOrderType()
}

// ===== COUNTRY

/* GET */

func TestGetCountries(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	countries := getCountries(1)
	if len(countries) == 0 || countries[0].Id <= 0 {
		t.Error("Can't scan countries")
		return
	}
}

func TestGetCountryRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	country := getCountryRow(1, 1)
	if country.Id <= 0 {
		t.Error("Can't scan country row")
		return
	}
}

func TestSearchCountries(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	countries := searchCountries("", 1)
	if len(countries) == 0 || countries[0].Id <= 0 {
		t.Error("Can't scan countries")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestCountryInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := Country{
		Name:        "Test",
		Iso2:        "XX",
		Iso3:        "XYZ",
		UNCode:      12345,
		Zone:        "E",
		PhonePrefix: 4321,
		enterprise:  1,
	}

	ok := c.insertCountry()
	if !ok {
		t.Error("Insert error, can't insert country")
		return
	}

	countries := getCountries(1)
	c = countries[len(countries)-1]

	c.Name = "Test test"
	ok = c.updateCountry()
	if !ok {
		t.Error("Update error, country not updated")
		return
	}

	c = getCountryRow(c.Id, 1)
	if c.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = c.deleteCountry()
	if !ok {
		t.Error("Delete error, country not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindCountryByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	countries := findCountryByName("", 1)
	if len(countries) == 0 || countries[0].Id <= 0 {
		t.Error("Can't scan countries")
		return
	}
}

func TestGetNameCountry(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	countryName := getNameCountry(1, 1)
	if countryName == "" {
		t.Error("Can't scan state name")
	}
}

// ===== STATE

/* GET */

func TestGetStates(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	states := getStates(1)
	if len(states) == 0 || states[0].Id <= 0 {
		t.Error("Can't scan states")
		return
	}
}

func TestGetStateRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	state := getStateRow(1)
	if state.Id <= 0 {
		t.Error("Can't scan state row")
		return
	}
}

func TestGetStatesByCountry(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	states := getStatesByCountry(1, 1)
	if len(states) == 0 || states[0].Id <= 0 {
		t.Error("Can't scan states")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestStateInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := State{
		Country:    1,
		Name:       "Test",
		IsoCode:    "XYZ",
		enterprise: 1,
	}

	ok := s.insertState()
	if !ok {
		t.Error("Insert error, can't insert state")
		return
	}

	states := getStates(1)
	s = states[len(states)-1]

	s.Name = "Test test"
	ok = s.updateState()
	if !ok {
		t.Error("Update error, state not updated")
		return
	}

	s = getStateRow(s.Id)
	if s.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = s.deleteState()
	if !ok {
		t.Error("Delete error, state not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindStateByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	countries := findCountryByName("", 1)
	if len(countries) == 0 || countries[0].Id <= 0 {
		t.Error("Can't scan countries")
		return
	}
}

func TestGetNameState(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	stateName := getNameState(1, 1)
	if stateName == "" {
		t.Error("Can't scan state name")
	}
}

// ===== COLOR

/* GET */

func TestGetColor(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	color := getColor(1)
	if len(color) == 0 || color[0].Id <= 0 {
		t.Error("Can't scan colors")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestColorInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := Color{
		Name:       "Test",
		HexColor:   "123456",
		enterprise: 1,
	}

	ok := c.insertColor()
	if !ok {
		t.Error("Insert error, can't insert color")
		return
	}

	colors := getColor(1)
	c = colors[len(colors)-1]

	c.Name = "Test test"
	ok = c.updateColor()
	if !ok {
		t.Error("Update error, color not updated")
		return
	}

	colors = getColor(1)
	c = colors[len(colors)-1]

	if c.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = c.deleteColor()
	if !ok {
		t.Error("Delete error, color not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindColorByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	colors := findColorByName("", 1)
	if len(colors) == 0 || colors[0].Id <= 0 {
		t.Error("Can't scan colors")
		return
	}
}

func TestGetNameColor(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	color := getNameColor(1, 1)
	if color == "" {
		t.Error("Can't scan color name")
	}
}

// ===== PRODUCT FAMILY

/* GET */

func TestGetProductFamilies(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	family := getProductFamilies(1)
	if len(family) == 0 || family[0].Id <= 0 {
		t.Error("Can't scan family")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestProductFamilyInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	pf := ProductFamily{
		Name:       "Test",
		Reference:  "XYZ",
		enterprise: 1,
	}

	ok := pf.insertProductFamily()
	if !ok {
		t.Error("Insert error, can't insert product family")
		return
	}

	families := getProductFamilies(1)
	pf = families[len(families)-1]

	pf.Name = "Test test"
	ok = pf.updateProductFamily()
	if !ok {
		t.Error("Update error, product family not updated")
		return
	}

	families = getProductFamilies(1)
	pf = families[len(families)-1]

	if pf.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = pf.deleteProductFamily()
	if !ok {
		t.Error("Delete error, product family not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindProductFamilyByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	pf := findProductFamilyByName("", 1)
	if len(pf) == 0 || pf[0].Id <= 0 {
		t.Error("Can't scan product families")
		return
	}
}

func TestGetNameProductFamily(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	color := getNameProductFamily(1, 1)
	if color == "" {
		t.Error("Can't scan product family name")
	}
}

// ===== ADDRESS

/* GET */

func TestGetAddresses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginationQuery{Offset: 0, Limit: 1, enterprise: 1}
	addresses := q.getAddresses()
	if len(addresses.Addresses) == 0 || addresses.Addresses[0].Id <= 0 {
		t.Error("Can't scan addresses")
		return
	}
}

func TestGetAddressRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	address := getAddressRow(1)
	if address.Id <= 0 {
		t.Error("Can't scan address row")
		return
	}
}

func TestSearchAddresses(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaginatedSearch{PaginationQuery: PaginationQuery{Offset: 0, Limit: 1, enterprise: 1}, Search: ""}
	addresses := q.searchAddresses()
	if len(addresses.Addresses) == 0 || addresses.Addresses[0].Id <= 0 {
		t.Error("Can't scan addresses")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestAddressInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	customer := int32(1)
	a := Address{
		Customer:          &customer,
		Supplier:          nil,
		Address:           "DVY NPPVHLE WFPZKKIBFAIYMMR RYFPAIBTBYENHAGGJPNNT",
		Address2:          "GUULBOTQGDPGHYTZKZNRT",
		City:              "NKTCH",
		Country:           1,
		PrivateOrBusiness: "_",
		ZipCode:           "AWS13",
		enterprise:        1,
	}

	ok := a.insertAddress(0).Id > 0
	if !ok {
		t.Error("Insert error, can't insert address")
		return
	}

	q := PaginationQuery{Offset: 0, Limit: MAX_INT32, enterprise: 1}
	addresses := q.getAddresses().Addresses
	a = addresses[len(addresses)-1]

	a.Address = "Test test"
	ok = a.updateAddress()
	if !ok {
		t.Error("Update error, address not updated")
		return
	}

	a = getAddressRow(a.Id)
	if a.Address != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = a.deleteAddress()
	if !ok {
		t.Error("Delete error, address not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestLocateAddressByCustomer(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	address := locateAddressByCustomer(1, 1)
	if len(address) == 0 || address[0].Id <= 0 {
		t.Error("Can't scan addresses")
		return
	}
}

func TestLocateAddressBySupplier(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	address := locateAddressBySupplier(1, 1)
	if len(address) == 0 || address[0].Id <= 0 {
		t.Error("Can't scan addresses")
		return
	}
}

func TestGetAddressName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	addressName := getAddressName(1, 1)
	if addressName == "" {
		t.Error("Can't scan address name")
	}
}

// ===== CARRIER

/* GET */

func TestGetCariers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carriers := getCariers(1)
	if len(carriers) == 0 || carriers[0].Id <= 0 {
		t.Error("Can't scan carriers")
		return
	}
}

func TestGetCarierRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carrier := getCarierRow(1)
	if carrier.Id <= 0 {
		t.Error("Can't scan carrier row")
		return
	}
}

func TestLocateCarriers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carriers := locateCarriers(1)
	if len(carriers) == 0 || carriers[0].Id <= 0 {
		t.Error("Can't scan carriers")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestCarrierInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := Carrier{
		Name:       "ACME Corp",
		MaxWeight:  35,
		MaxWidth:   150,
		MaxHeight:  150,
		MaxDepth:   150,
		Phone:      "987654321",
		Email:      "contact@acme.com",
		Web:        "acmecorp.com",
		Webservice: "_",
		enterprise: 1,
	}

	ok := c.insertCarrier()
	if !ok {
		t.Error("Insert error, can't insert carrier")
		return
	}

	carriers := getCariers(1)
	c = carriers[len(carriers)-1]

	c.Name = "Test test"
	ok = c.updateCarrier()
	if !ok {
		t.Error("Update error, carrier not updated")
		return
	}

	c = getCarierRow(c.Id)
	if c.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = c.deleteCarrier()
	if !ok {
		t.Error("Delete error, carrier not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindCarrierByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carriers := findCarrierByName("", 1)
	if len(carriers) == 0 || carriers[0].Id <= 0 {
		t.Error("Can't scan carriers")
		return
	}
}

func TestGetNameCarrier(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carrierName := getNameCarrier(1, 1)
	if carrierName == "" {
		t.Error("Can't scan carrier name")
	}
}

// ===== BILLING SERIE

/* GET */

func TestGetBillingSeries(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	billingSeries := getBillingSeries(1)
	if len(billingSeries) == 0 || len(billingSeries[0].Id) == 0 {
		t.Error("Can't scan billing series")
		return
	}
}

func TestLocateBillingSeries(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	billingSeries := locateBillingSeries(1)
	if len(billingSeries) == 0 || len(billingSeries[0].Id) == 0 {
		t.Error("Can't scan billing series")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestBillingSeriesInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	b := BillingSerie{
		Id:          "EXA",
		Name:        "Example series",
		BillingType: "S",
		Year:        2021,
		enterprise:  1,
	}

	ok := b.insertBillingSerie()
	if !ok {
		t.Error("Insert error, can't insert billing series")
		return
	}

	billingSeries := getBillingSeries(1)
	for i := 0; i < len(billingSeries); i++ {
		if billingSeries[i].Id == "EXA" {
			b = billingSeries[i]
			break
		}
	}

	b.Name = "Test test"
	ok = b.updateBillingSerie()
	if !ok {
		t.Error("Update error, billing series not updated")
		return
	}

	billingSeries = getBillingSeries(1)
	for i := 0; i < len(billingSeries); i++ {
		if billingSeries[i].Id == "EXA" {
			b = billingSeries[i]
			break
		}
	}

	if b.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = b.deleteBillingSerie()
	if !ok {
		t.Error("Delete error, billing serie not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindBillingSerieByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	billingSeries := findBillingSerieByName("", 1)
	if len(billingSeries) == 0 || len(billingSeries[0].Id) == 0 {
		t.Error("Can't scan billing series")
		return
	}
}

func TestGetNameBillingSerie(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	billingSeries := getNameBillingSerie("INT", 1)
	if billingSeries == "" {
		t.Error("Can't scan billing series name")
	}
}

func TestGetNextNumberBillingSeries(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	number := getNextSaleOrderNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next sale order number")
	}

	number = getNextSaleInvoiceNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next sale invoice number")
	}

	number = getNextSaleDeliveryNoteNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next sale delivery note number")
	}

	number = getNextPurchaseOrderNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next purchase order number")
	}

	number = getNextPurchaseInvoiceNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next purchase invoicer number")
	}

	number = getNextPurchaseDeliveryNoteNumber("INT", 1)
	if number <= 0 {
		t.Error("Could not get next purchase delivery note number")
	}
}

// ===== CURRENCY

/* GET */

func TestGetCurrencies(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	currencies := getCurrencies(1)
	if len(currencies) == 0 || currencies[0].Id <= 0 {
		t.Error("Can't scan currencies")
		return
	}
}

func TestLocateCurrency(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	currencies := locateCurrency(1)
	if len(currencies) == 0 || currencies[0].Id <= 0 {
		t.Error("Can't scan currencies")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestCurrenciesInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	c := Currency{
		Name:         "Bitcoin",
		Sign:         "₿",
		IsoCode:      "BTC",
		IsoNum:       12345,
		Change:       1,
		ExchangeDate: time.Now(),
		enterprise:   1,
	}

	ok := c.insertCurrency()
	if !ok {
		t.Error("Insert error, can't insert currencies")
		return
	}

	currencies := getCurrencies(1)
	c = currencies[len(currencies)-1]

	c.Name = "Test test"
	ok = c.updateCurrency()
	if !ok {
		t.Error("Update error, currency not updated")
		return
	}

	currencies = getCurrencies(1)
	c = currencies[len(currencies)-1]

	if c.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = c.deleteCurrency()
	if !ok {
		t.Error("Delete error, currency not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindCurrencyByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	currency := findCurrencyByName("", 1)
	if len(currency) == 0 || currency[0].Id <= 0 {
		t.Error("Can't scan currency")
		return
	}
}

func TestGetNameCurrency(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	currencyName := getNameCurrency(1, 1)
	if currencyName == "" {
		t.Error("Can't scan currency name")
	}
}

func TestGetCurrencyExchange(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	currencyExchange := getCurrencyExchange(1)
	if currencyExchange <= 0 {
		t.Error("Can't scan currency exchange")
	}
}

// ===== PAYMENT METHOD

/* GET */

func TestGetPaymentMethods(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	paymentMethods := getPaymentMethods(1)
	if len(paymentMethods) == 0 || paymentMethods[0].Id <= 0 {
		t.Error("Can't scan payment methods")
		return
	}
}

func TestGetPaymentMethodRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carrier := getPaymentMethodRow(1)
	if carrier.Id <= 0 {
		t.Error("Can't scan payment method row")
		return
	}
}

func TestLocatePaymentMethods(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	paymentMethods := locatePaymentMethods(1)
	if len(paymentMethods) == 0 || paymentMethods[0].Id <= 0 {
		t.Error("Can't scan payment methods")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestPaymentMethodInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	pm := PaymentMethod{
		Name:                 "Bitcoin",
		PaidInAdvance:        true,
		PrestashopModuleName: "btc",
		DaysExpiration:       0,
		enterprise:           1,
	}

	ok := pm.insertPaymentMethod()
	if !ok {
		t.Error("Insert error, can't insert payment method")
		return
	}

	paymentMethods := getPaymentMethods(1)
	pm = paymentMethods[len(paymentMethods)-1]

	pm.Name = "Test test"
	ok = pm.updatePaymentMethod()
	if !ok {
		t.Error("Update error, payment method not updated")
		return
	}

	paymentMethods = getPaymentMethods(1)
	pm = paymentMethods[len(paymentMethods)-1]

	if pm.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = pm.deletePaymentMethod()
	if !ok {
		t.Error("Delete error, payment method not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindPaymentMethodByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	paymentMethod := findPaymentMethodByName("", 1)
	if len(paymentMethod) == 0 || paymentMethod[0].Id <= 0 {
		t.Error("Can't scan payment methods")
		return
	}
}

func TestGetNamePaymentMethod(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	paymentMethodName := getNamePaymentMethod(1, 1)
	if paymentMethodName == "" {
		t.Error("Can't scan payment method name")
	}
}

// ===== LANGUAGE

/* GET */

func TestGetLanguages(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	languages := getLanguages(1)
	if len(languages) == 0 || languages[0].Id <= 0 {
		t.Error("Can't scan payment methods")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestLanguageInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	l := Language{
		Name:       "ACME Corp. Super Secrete Language",
		Iso2:       "AC",
		Iso3:       "ACM",
		enterprise: 1,
	}

	ok := l.insertLanguage()
	if !ok {
		t.Error("Insert error, can't insert language")
		return
	}

	language := getLanguages(1)
	l = language[len(language)-1]

	l.Name = "Test test"
	ok = l.updateLanguage()
	if !ok {
		t.Error("Update error, language not updated")
		return
	}

	language = getLanguages(1)
	l = language[len(language)-1]

	if l.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = l.deleteLanguage()
	if !ok {
		t.Error("Delete error, language not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestFindLanguageByName(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	language := findLanguageByName("", 1)
	if len(language) == 0 || language[0].Id <= 0 {
		t.Error("Can't scan language")
		return
	}
}

func TestGetNameLanguage(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	language := getNameLanguage(1, 1)
	if language == "" {
		t.Error("Can't scan language name")
	}
}

// ===== PACKAGES

/* GET */

func TestGetPackages(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	packages := getPackages(1)
	if len(packages) == 0 || packages[0].Id <= 0 {
		t.Error("Can't scan packages")
		return
	}
}

func TestGetPackagesRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	carrier := getPackagesRow(1)
	if carrier.Id <= 0 {
		t.Error("Can't scan packages row")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestPackagesInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	p := Packages{
		Name:       "Test box",
		Weight:     1,
		Width:      40,
		Height:     40,
		Depth:      40,
		Product:    1,
		enterprise: 1,
	}

	ok := p.insertPackage()
	if !ok {
		t.Error("Insert error, can't insert package")
		return
	}

	language := getPackages(1)
	p = language[len(language)-1]

	p.Name = "Test test"
	ok = p.updatePackage()
	if !ok {
		t.Error("Update error, package not updated")
		return
	}

	language = getPackages(1)
	p = language[len(language)-1]

	if p.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = p.deletePackage()
	if !ok {
		t.Error("Delete error, package not deleted")
		return
	}
}

// ===== INCOTERMS

/* GET */

func TestGetIncoterm(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	incoterms := getIncoterm(1)
	if len(incoterms) == 0 || incoterms[0].Id <= 0 {
		t.Error("Can't scan incoterms")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestIncotermsInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	i := Incoterm{
		Name:       "Test incoterm",
		Key:        "TST",
		enterprise: 1,
	}

	ok := i.insertIncoterm()
	if !ok {
		t.Error("Insert error, can't insert incoterm")
		return
	}

	incoterms := getIncoterm(1)
	i = incoterms[len(incoterms)-1]

	i.Name = "Test test"
	ok = i.updateIncoterm()
	if !ok {
		t.Error("Update error, incoterm not updated")
		return
	}

	incoterms = getIncoterm(1)
	i = incoterms[len(incoterms)-1]

	if i.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = i.deleteIncoterm()
	if !ok {
		t.Error("Delete error, incoterm not deleted")
		return
	}
}

// ===== DOCUMENT CONTAINER

/* GET */

func TestGetDocumentContainer(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	docContainers := getDocumentContainer(1)
	if len(docContainers) == 0 || docContainers[0].Id <= 0 {
		t.Error("Can't scan document containers")
		return
	}
}

func TestGetDocumentContainerRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	docContainer := getDocumentContainerRow(1)
	if docContainer.Id <= 0 {
		t.Error("Can't scan document container row")
		return
	}
}

/* INSERT - UPDATE - DELETE */

func TestDocumentContainerInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	p := DocumentContainer{
		Name:        "Test",
		Path:        "/marketnet/docs/",
		MaxFileSize: 1000,
		enterprise:  1,
	}

	ok := p.insertDocumentContainer()
	if !ok {
		t.Error("Insert error, can't insert document container")
		return
	}

	docContainers := getDocumentContainer(1)
	p = docContainers[len(docContainers)-1]

	p.Name = "Test test"
	ok = p.updateDocumentContainer()
	if !ok {
		t.Error("Update error, document container not updated")
		return
	}

	docContainers = getDocumentContainer(1)
	p = docContainers[len(docContainers)-1]

	if p.Name != "Test test" {
		t.Error("Update not successful")
		return
	}

	ok = p.deleteDocumentContainer()
	if !ok {
		t.Error("Delete error, document container not deleted")
		return
	}
}

/* FUNCTIONALITY */

func TestLocateDocumentContainer(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	docContainers := locateDocumentContainer(1)
	if len(docContainers) == 0 || docContainers[0].Id <= 0 {
		t.Error("Can't scan document containers")
		return
	}
}
