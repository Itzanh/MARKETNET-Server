package main

import (
	"strings"
)

type ProductFamily struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Reference  string `json:"reference"`
	enterprise int32
}

func getProductFamilies(enterpriseId int32) []ProductFamily {
	var families []ProductFamily = make([]ProductFamily, 0)
	sqlStatement := `SELECT * FROM public.product_family WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return families
	}
	for rows.Next() {
		f := ProductFamily{}
		rows.Scan(&f.Id, &f.Name, &f.Reference, &f.enterprise)
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

	sqlStatement := `INSERT INTO public.product_family(name, reference, enterprise) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, f.Name, f.Reference, f.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (f *ProductFamily) updateProductFamily() bool {
	if f.Id <= 0 || !f.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.product_family SET name=$2, reference=$3 WHERE id=$1 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, f.Id, f.Name, f.Reference, f.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (f *ProductFamily) deleteProductFamily() bool {
	if f.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.product_family WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, f.Id, f.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findProductFamilyByName(productFamilyName string, enterpriseId int32) []NameInt16 {
	var productFamily []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.product_family WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(productFamilyName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return productFamily
	}
	for rows.Next() {
		p := NameInt16{}
		rows.Scan(&p.Id, &p.Name)
		productFamily = append(productFamily, p)
	}

	return productFamily
}

func getNameProductFamily(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.product_family WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
