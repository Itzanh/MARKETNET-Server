package main

type Stock struct {
	Product                    int32  `json:"product"`
	Warehouse                  string `json:"warehouse"`
	Quantity                   int32  `json:"quantity"`
	QuantityPendingReceived    int32  `json:"quantityPendingReceived"`
	QuantityPendingServed      int32  `json:"quantityPendingServed"`
	QuantityAvaialbe           int32  `json:"quantityAvaialbe"`
	QuantityPendingManufacture int32  `json:"quantityPendingManufacture"`
}

func getStock(productId int32) []Stock {
	var stock []Stock = make([]Stock, 0)
	sqlStatement := `SELECT * FROM stock WHERE product = $1 ORDER BY warehouse ASC`
	rows, err := db.Query(sqlStatement, productId)
	if err != nil {
		return stock
	}
	for rows.Next() {
		s := Stock{}
		rows.Scan(&s.Product, &s.Warehouse, &s.Quantity, &s.QuantityPendingReceived, &s.QuantityPendingServed, &s.QuantityAvaialbe, &s.QuantityPendingManufacture)
		stock = append(stock, s)
	}

	return stock
}

func getStockRow(productId int32, warehouseId string) Stock {
	sqlStatement := `SELECT * FROM stock WHERE product = $1 AND warehouse = $2`
	row := db.QueryRow(sqlStatement, productId, warehouseId)
	if row.Err() != nil {
		return Stock{}
	}

	s := Stock{}
	row.Scan(&s.Product, &s.Warehouse, &s.Quantity, &s.QuantityPendingReceived, &s.QuantityPendingServed, &s.QuantityAvaialbe, &s.QuantityPendingManufacture)

	return s
}

// Inserts a row with 0 stock in all columns
func createStockRow(productId int32, warehouseId string) bool {
	sqlStatement := `INSERT INTO stock (product,warehouse) VALUES ($1,$2)`
	res, err := db.Exec(sqlStatement, productId, warehouseId)
	rows, _ := res.RowsAffected()
	return rows > 0 && err == nil
}

// Adds an amount to the quantity pending of serving, and substract the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingServing(productId int32, warehouseId string, quantity int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_pending_served=quantity_pending_served+$3, quantity_available=quantity_available-$3 WHERE product=$1 AND warehouse=$2`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity)

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId) { // we create the row, and retry the operation
			return addQuantityPendingServing(productId, warehouseId, quantity)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	return err == nil
}

// Adds an amount to the quantity pending of receiving, and add to the amount from the quantity available.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityPendingReveiving(productId int32, warehouseId string, quantity int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity_pending_received=quantity_pending_received+$3, quantity_available=quantity_available+$3 WHERE product=$1 AND warehouse=$2`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity)

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId) { // we create the row, and retry the operation
			return addQuantityPendingReveiving(productId, warehouseId, quantity)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	return err == nil
}

// Add an amount to the stock column on the stock row for this product.
// This function will do this operation inversely if the parameter quantity is a negative number.
// Creates the stock row if it doesn't exists.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION
func addQuantityStock(productId int32, warehouseId string, quantity int32) bool {
	sqlStatement := `UPDATE public.stock SET quantity=quantity+$3 WHERE product=$1 AND warehouse=$2`
	res, err := db.Exec(sqlStatement, productId, warehouseId, quantity)

	rows, _ := res.RowsAffected()
	if rows == 0 && err == nil { // no error has ocurred, but the query hasn't affected any row. we assume that the stock row does not exist yet
		if createStockRow(productId, warehouseId) { // we create the row, and retry the operation
			return addQuantityStock(productId, warehouseId, quantity)
		} else {
			return false // the row could neither not be created or updated
		}
	}

	return err == nil
}
