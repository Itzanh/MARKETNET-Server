package main

import (
	"strings"
)

type Warehouse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func getWarehouses() []Warehouse {
	var warehouses []Warehouse = make([]Warehouse, 0)
	sqlStatement := `SELECT * FROM public.warehouse ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return warehouses
	}
	for rows.Next() {
		w := Warehouse{}
		rows.Scan(&w.Id, &w.Name)
		warehouses = append(warehouses, w)
	}

	return warehouses
}

func (w *Warehouse) isValid() bool {
	return !(len(w.Id) == 0 || len(w.Id) > 2 || len(w.Name) == 0 || len(w.Name) > 50)
}

func (w *Warehouse) insertWarehouse() bool {
	if !w.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.warehouse(id, name) VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, w.Id, w.Name)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (w *Warehouse) updateWarehouse() bool {
	if !w.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.warehouse SET name=$2 WHERE id = $1`
	res, err := db.Exec(sqlStatement, w.Id, w.Name)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (w *Warehouse) deleteWarehouse() bool {
	if w.Id == "" || len(w.Id) != 2 {
		return false
	}

	sqlStatement := `DELETE FROM warehouse WHERE id = $1`
	res, err := db.Exec(sqlStatement, w.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findWarehouseByName(languageName string) []NameString {
	var warehouses []NameString = make([]NameString, 0)
	sqlStatement := `SELECT id,name FROM public.warehouse WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		log("DB", err.Error())
		return warehouses
	}
	for rows.Next() {
		w := NameString{}
		rows.Scan(&w.Id, &w.Name)
		warehouses = append(warehouses, w)
	}

	return warehouses
}

func getNameWarehouse(id string) string {
	sqlStatement := `SELECT name FROM public.warehouse WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

// Regenerates the stock of the product for all the products in the database.
// This "stock" field is the sum of the stock in all the warehouses.
func regenerateProductStock() bool {
	sqlStatement := `UPDATE product SET stock = CASE WHEN (SELECT SUM(quantity) FROM stock WHERE stock.product=product.id) IS NULL THEN 0 ELSE (SELECT SUM(quantity) FROM stock WHERE stock.product=product.id) END`
	_, err := db.Exec(sqlStatement)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
