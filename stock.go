package main

type Stock struct {
	Product                    int32  `json:"product"`
	Warehouse                  string `json:"warehouse"`
	Quantity                   int32  `json:"quantity"`
	QuantityPendingReceived    int32  `json:"quantityPendingReceived"`
	QuantityPendingServed      int32  `json:"quantityPendingServed"`
	QuantityAvaialbe           int32  `json:"quantityAvaialbe"`
	QuantityPendingManufacture int32  `json:"quantityPendingManufacture"`
	WarehouseName              string `json:"warehouseName"`
	enterprise                 int32
}

func getStock(productId int32, enterpriseId int32) []Stock {
	var stock []Stock = make([]Stock, 0)
	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id=stock.warehouse AND warehouse.enterprise=stock.enterprise) FROM stock WHERE product=$1 AND enterprise=$2 ORDER BY warehouse ASC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return stock
	}
	defer rows.Close()

	for rows.Next() {
		s := Stock{}
		rows.Scan(&s.Product, &s.Warehouse, &s.Quantity, &s.QuantityPendingReceived, &s.QuantityPendingServed, &s.QuantityAvaialbe, &s.QuantityPendingManufacture, &s.enterprise, &s.WarehouseName)
		stock = append(stock, s)
	}

	return stock
}

func getStockRow(productId int32, warehouseId string, enterpriseId int32) Stock {
	sqlStatement := `SELECT * FROM stock WHERE product=$1 AND warehouse=$2 AND enterprise=$3`
	row := db.QueryRow(sqlStatement, productId, warehouseId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Stock{}
	}

	s := Stock{}
	row.Scan(&s.Product, &s.Warehouse, &s.Quantity, &s.QuantityPendingReceived, &s.QuantityPendingServed, &s.QuantityAvaialbe, &s.QuantityPendingManufacture, &s.enterprise)

	return s
}

// Inserts a row with 0 stock in all columns
func createStockRow(productId int32, warehouseId string, enterpriseId int32) bool {
	sqlStatement := `INSERT INTO stock (product,warehouse,enterprise) VALUES ($1,$2,$3)`
	res, err := db.Exec(sqlStatement, productId, warehouseId, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0 && err == nil
}

// Adds an amount to the quantity pending of serving, and substract the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingServing(productId int32, warehouseId string, quantity int32, enterpriseId int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_pending_served=quantity_pending_served+$3, quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId) { // we create the row, and retry the operation
			return addQuantityPendingServing(productId, warehouseId, quantity, enterpriseId)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	ok := setQuantityAvailable(productId, warehouseId, enterpriseId)
	if !ok {
		return false
	}

	return err == nil
}

// Adds an amount to the quantity pending of receiving, and add to the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingReveiving(productId int32, warehouseId string, quantity int32, enterpriseId int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_pending_received=quantity_pending_received+$3, quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId) { // we create the row, and retry the operation
			return addQuantityPendingReveiving(productId, warehouseId, quantity, enterpriseId)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	ok := setQuantityAvailable(productId, warehouseId, enterpriseId)
	if !ok {
		return false
	}

	return err == nil
}

// Adds an amount to the quantity pending of manufacturing, and add to the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingManufacture(productId int32, warehouseId string, quantity int32, enterpriseId int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_pending_manufacture=quantity_pending_manufacture+$3, quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId) { // we create the row, and retry the operation
			return addQuantityPendingReveiving(productId, warehouseId, quantity, enterpriseId)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	ok := setQuantityAvailable(productId, warehouseId, enterpriseId)
	if !ok {
		return false
	}

	return err == nil
}

// Add an amount to the stock column on the stock row for this product.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityStock(productId int32, warehouseId string, quantity int32, enterpriseId int32) bool {
	productRow := getProductRow(productId)
	if productRow.enterprise != enterpriseId {
		return false
	}
	if !productRow.ControlStock {
		return true
	}

	sqlStatement := `UPDATE public.stock SET quantity=quantity+$3, quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId) { // we create the row, and retry the operation
			return addQuantityStock(productId, warehouseId, quantity, enterpriseId)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	if err == nil {
		return setQuantityAvailable(productId, warehouseId, enterpriseId) && setProductStockAllWarehouses(productId)
	} else {
		return false
	}
}

// Sets an amount to the stock column on the stock row for this product.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setQuantityStock(productId int32, warehouseId string, quantity int32, enterpriseId int32) bool {
	productRow := getProductRow(productId)
	if productRow.enterprise != enterpriseId {
		return false
	}
	if !productRow.ControlStock {
		return true
	}

	sqlStatement := `UPDATE public.stock SET quantity=$3, quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId, enterpriseId) { // we create the row, and retry the operation
			return setQuantityStock(productId, warehouseId, quantity, enterpriseId)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	if err == nil {
		return setQuantityAvailable(productId, warehouseId, enterpriseId) && setProductStockAllWarehouses(productId)
	} else {
		return false
	}
}

// Sets the "Quantity available" field on the stock in the products.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setQuantityAvailable(productId int32, warehouseId string, enterpriseId int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_available=quantity+quantity_pending_received-quantity_pending_served+quantity_pending_manufacture WHERE product=$1 AND warehouse=$2 AND enterprise=$3`
	res, err := db.Exec(sqlStatement, productId, warehouseId, enterpriseId)

	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0 && err == nil
}

// Sets the "stock" field on the product row, sum of the stocks in all the warehouses.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func setProductStockAllWarehouses(productId int32) bool {
	sqlStatement := `UPDATE product SET stock = (SELECT SUM(quantity) FROM stock WHERE product=$1) WHERE id=$1`
	_, err := db.Exec(sqlStatement, productId)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
