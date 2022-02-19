package main

import (
	"time"
)

type Inventory struct {
	Id           int32      `json:"id"`
	Name         string     `json:"name"`
	DateCreated  time.Time  `json:"dateCreated"`
	Finished     bool       `json:"finished"`
	DateFinished *time.Time `json:"dateFinished"`
	Warehouse    string     `json:"warehouse"`
	enterprise   int32
}

func getInventories(enterpriseId int32) []Inventory {
	var inventory []Inventory = make([]Inventory, 0)
	sqlStatement := `SELECT * FROM public.inventory WHERE enterprise = $1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return inventory
	}

	for rows.Next() {
		i := Inventory{}
		rows.Scan(&i.Id, &i.enterprise, &i.Name, &i.DateCreated, &i.Finished, &i.DateFinished, &i.Warehouse)
		inventory = append(inventory, i)
	}

	return inventory
}

func getInventoryRow(inventoryId int32) Inventory {
	sqlStatement := `SELECT * FROM public.inventory WHERE id = $1 ORDER BY id DESC`
	row := db.QueryRow(sqlStatement, inventoryId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Inventory{}
	}

	i := Inventory{}
	row.Scan(&i.Id, &i.enterprise, &i.Name, &i.DateCreated, &i.Finished, &i.DateFinished, &i.Warehouse)
	return i
}

func (i *Inventory) isValid() bool {
	return !(len(i.Name) == 0 || len(i.Name) > 50 || len(i.Warehouse) == 0 || len(i.Warehouse) > 2)
}

func (i *Inventory) insertInventory(enterpriseId int32) bool {
	if !i.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.inventory(enterprise, name, warehouse) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, enterpriseId, i.Name, i.Warehouse)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (i *Inventory) deleteInventory(enterpriseId int32) bool {
	if i.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.inventory WHERE id = $1 AND enterprise = $2`
	_, err := db.Exec(sqlStatement, i.Id, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (i *Inventory) finishInventory(userId int32, enterpriseId int32) bool {
	inMemoyInventory := getInventoryRow(i.Id)
	if inMemoyInventory.Id <= 0 || inMemoyInventory.Finished || inMemoyInventory.enterprise != enterpriseId {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `UPDATE public.inventory_products SET warehouse_movement=$3 WHERE inventory=$1 AND product=$2`
	var ok bool
	products := getInventoryProducts(inMemoyInventory.Id, enterpriseId)
	for i := 0; i < len(products); i++ {
		p := products[i]

		wm := WarehouseMovement{
			enterprise: enterpriseId,
			Warehouse:  inMemoyInventory.Warehouse,
			Product:    p.Product,
			Quantity:   p.Quantity,
			Type:       "R",
		}
		ok = wm.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		_, err := trans.Exec(sqlStatement, p.Inventory, p.Product, wm.Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
	}

	sqlStatement = `UPDATE public.inventory SET finished=true, date_finished=CURRENT_TIMESTAMP(3) WHERE id=$1`
	_, err := trans.Exec(sqlStatement, inMemoyInventory.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///
	return true
}

type InventoryProducts struct {
	Inventory         int32  `json:"inventory"`
	Product           int32  `json:"product"`
	Quantity          int32  `json:"quantity"`
	WarehouseMovement *int64 `json:"warehouseMovement"`
	ProductName       string `json:"productName"`
	enterprise        int32
}

func getInventoryProducts(inventoryId int32, enterpriseId int32) []InventoryProducts {
	var inventoryProducts []InventoryProducts = make([]InventoryProducts, 0)
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=inventory_products.product) FROM public.inventory_products WHERE inventory = $1 AND enterprise = $2 ORDER BY product ASC`
	rows, err := db.Query(sqlStatement, inventoryId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return inventoryProducts
	}

	for rows.Next() {
		ip := InventoryProducts{}
		rows.Scan(&ip.Inventory, &ip.Product, &ip.enterprise, &ip.Quantity, &ip.WarehouseMovement, &ip.ProductName)
		inventoryProducts = append(inventoryProducts, ip)
	}

	return inventoryProducts
}

func getInventoryProductsRow(inventoryId int32, productId int32, enterpriseId int32) InventoryProducts {
	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=inventory_products.product) FROM public.inventory_products WHERE inventory = $1 AND product = $2 AND enterprise = $3 ORDER BY product ASC`
	row := db.QueryRow(sqlStatement, inventoryId, productId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return InventoryProducts{}
	}

	ip := InventoryProducts{}
	row.Scan(&ip.Inventory, &ip.Product, &ip.enterprise, &ip.Quantity, &ip.WarehouseMovement, &ip.ProductName)

	return ip
}

func (i *InventoryProducts) isValid() bool {
	return !(i.Inventory <= 0 || i.Product <= 0)
}

type InputInventoryProducts struct {
	Inventory         int32               `json:"inventory"`
	InventoryProducts []InventoryProducts `json:"inventoryProducts"`
	FamilyId          int32               `json:"familyId"`
}

func (input *InputInventoryProducts) insertUpdateDeleteInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.enterprise != enterpriseId || i.Finished {
		return false
	}

	for i := 0; i < len(input.InventoryProducts); i++ {
		ip := input.InventoryProducts[i]
		if !ip.isValid() {
			return false
		}
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	// input data optimization
	arayExistentInventoryProducts := getInventoryProducts(input.Inventory, enterpriseId)
	var existentInventoryProducts map[int32]InventoryProducts = make(map[int32]InventoryProducts)
	for i := 0; i < len(arayExistentInventoryProducts); i++ {
		a := arayExistentInventoryProducts[i]
		existentInventoryProducts[a.Product] = a
	}

	// cross data
	var toInsert []InventoryProducts = make([]InventoryProducts, 0)
	var toUpdate []InventoryProducts = make([]InventoryProducts, 0)
	for i := 0; i < len(input.InventoryProducts); i++ {
		newIp := input.InventoryProducts[i]

		oldIp, ok := existentInventoryProducts[newIp.Product]
		delete(existentInventoryProducts, newIp.Product)
		if !ok {
			toInsert = append(toInsert, newIp)
		} else if oldIp.Quantity != newIp.Quantity {
			toUpdate = append(toUpdate, newIp)
		}
	}

	// insert data
	for i := 0; i < len(toInsert); i++ {
		pi := toInsert[i]
		sqlStatement := `INSERT INTO public.inventory_products(inventory, product, enterprise, quantity) VALUES ($1, $2, $3, $4)`
		_, err := trans.Exec(sqlStatement, input.Inventory, pi.Product, enterpriseId, pi.Quantity)
		if err != nil {
			trans.Rollback()
			log("DB", err.Error())
			return false
		}
	}

	// update data
	for i := 0; i < len(toUpdate); i++ {
		pi := toUpdate[i]
		sqlStatement := `UPDATE public.inventory_products SET quantity=$3 WHERE inventory=$1 AND product=$2`
		_, err := trans.Exec(sqlStatement, input.Inventory, pi.Product, pi.Quantity)
		if err != nil {
			trans.Rollback()
			log("DB", err.Error())
			return false
		}
	}

	// delete the remaining data in the map
	for k := range existentInventoryProducts {
		sqlStatement := `DELETE FROM public.inventory_products WHERE inventory=$1 AND product=$2`
		_, err := trans.Exec(sqlStatement, input.Inventory, k)
		if err != nil {
			trans.Rollback()
			log("DB", err.Error())
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) insertProductFamilyInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.enterprise != enterpriseId || i.Finished {
		return false
	}

	sqlStatement := `SELECT enterprise FROM product_family WHERE id=$1`
	row := db.QueryRow(sqlStatement, input.FamilyId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var enterprise int32
	row.Scan(&enterprise)
	if enterprise != enterpriseId {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement = `SELECT id FROM product WHERE family = $1 AND enterprise = $2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, input.FamilyId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `INSERT INTO public.inventory_products(inventory, product, enterprise, quantity) VALUES ($1, $2, $3, $4)`
	var productId int32
	for rows.Next() {
		rows.Scan(&productId)

		_, err := trans.Exec(sqlStatement, input.Inventory, productId, enterpriseId, 0)
		if err != nil {
			trans.Rollback()
			log("DB", err.Error())
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) insertAllProductsInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.enterprise != enterpriseId || i.Finished {
		return false
	}

	products := getProduct(enterpriseId)

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.inventory_products(inventory, product, enterprise, quantity) VALUES ($1, $2, $3, $4)`
	var productId int32
	for i := 0; i < len(products); i++ {
		if products[i].Off {
			continue
		}
		productId = products[i].Id

		_, err := trans.Exec(sqlStatement, input.Inventory, productId, enterpriseId, 0)
		if err != nil {
			trans.Rollback()
			log("DB", err.Error())
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) deleteAllProductsInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.enterprise != enterpriseId || i.Finished {
		return false
	}

	sqlStatement := `DELETE FROM public.inventory_products WHERE inventory = $1 AND enterprise = $2`
	_, err := db.Exec(sqlStatement, input.Inventory, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

type BarCodeInputInventoryProducts struct {
	Inventory int32  `json:"inventory"`
	BarCode   string `json:"barCode"`
}

type BarCodeInputInventoryProductsResult struct {
	Ok               bool   `json:"ok"`
	ProductReference string `json:"productReference"`
	ProductName      string `json:"productName"`
	Quantity         int32  `json:"quantity"`
}

func (input *BarCodeInputInventoryProducts) insertOrCountInventoryProductsByBarcode(enterpriseId int32) BarCodeInputInventoryProductsResult {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.enterprise != enterpriseId || i.Finished {
		return BarCodeInputInventoryProductsResult{}
	}

	product := getProductByBarcode(input.BarCode, enterpriseId)
	if product.Id <= 0 {
		return BarCodeInputInventoryProductsResult{}
	}

	sqlStatement := `SELECT COUNT(*) FROM public.inventory_products WHERE inventory = $1 AND product =$2`
	row := db.QueryRow(sqlStatement, input.Inventory, product.Id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return BarCodeInputInventoryProductsResult{}
	}

	var rowCount int16
	row.Scan(&rowCount)
	if rowCount == 0 {
		sqlStatement = `INSERT INTO public.inventory_products(inventory, product, enterprise, quantity) VALUES ($1, $2, $3, $4)`
		_, err := db.Exec(sqlStatement, input.Inventory, product.Id, enterpriseId, 1)
		if err != nil {
			log("DB", err.Error())
			return BarCodeInputInventoryProductsResult{}
		}
	} else {
		sqlStatement = `UPDATE public.inventory_products SET quantity = quantity + 1 WHERE inventory = $1 AND product = $2`
		_, err := db.Exec(sqlStatement, input.Inventory, product.Id)
		if err != nil {
			log("DB", err.Error())
			return BarCodeInputInventoryProductsResult{}
		}
	}

	inventoryProduct := getInventoryProductsRow(input.Inventory, product.Id, enterpriseId)
	return BarCodeInputInventoryProductsResult{Ok: true, ProductReference: product.Reference, ProductName: product.Name, Quantity: inventoryProduct.Quantity}
}
