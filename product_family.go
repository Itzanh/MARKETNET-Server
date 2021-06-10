package main

import (
	"strings"
)

type ProductFamily struct {
	Id        int16  `json:"id"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

func getProductFamilies() []ProductFamily {
	var families []ProductFamily = make([]ProductFamily, 0)
	sqlStatement := `SELECT * FROM public.product_family ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return families
	}
	for rows.Next() {
		f := ProductFamily{}
		rows.Scan(&f.Id, &f.Name, &f.Reference)
		families = append(families, f)
	}

	return families
}

func (f *ProductFamily) isValid() bool {
	return !(len(f.Name) == 0 || len(f.Name) > 100 || len(f.Reference) == 0 || len(f.Reference) > 40)
}

func (f *ProductFamily) insertProductFamily() bool {
	if !f.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.product_family(name, reference) VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, f.Name, f.Reference)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (f *ProductFamily) updateProductFamily() bool {
	if f.Id <= 0 || !f.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.product_family SET name=$2, reference=$3 WHERE id=$1`
	res, err := db.Exec(sqlStatement, f.Id, f.Name, f.Reference)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (f *ProductFamily) deleteProductFamily() bool {
	if f.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.product_family WHERE id=$1`
	res, err := db.Exec(sqlStatement, f.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findProductFamilyByName(productFamilyName string) []NameInt16 {
	var productFamily []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.product_family WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(productFamilyName))
	if err != nil {
		return productFamily
	}
	for rows.Next() {
		p := NameInt16{}
		rows.Scan(&p.Id, &p.Name)
		productFamily = append(productFamily, p)
	}

	return productFamily
}

func getNameProductFamily(id int16) string {
	sqlStatement := `SELECT name FROM public.product_family WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
