package main

import "fmt"

type Packaging struct {
	Id              int32                      `json:"id"`
	Package         int16                      `json:"package"`
	PackageName     string                     `json:"packageName"`
	SalesOrder      int32                      `json:"salesOrder"`
	Weight          float32                    `json:"weight"`
	Shipping        *int32                     `json:"shipping"`
	DetailsPackaged []SalesOrderDetailPackaged `json:"detailsPackaged"`
}

func getPackaging(salesOrderId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	sqlStatement := `SELECT * FROM public.packaging WHERE sales_order=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, salesOrderId)
	if err != nil {
		return packaging
	}
	for rows.Next() {
		p := Packaging{}
		rows.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping)
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id)

		_package := getPackagesRow(p.Package)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging = append(packaging, p)
	}

	return packaging
}

func (p *Packaging) isValid() bool {
	return !(p.Package <= 0 || p.SalesOrder <= 0)
}

func (p *Packaging) insertPackaging() bool {
	if !p.isValid() {
		return false
	}

	_package := getPackagesRow(p.Package)
	if _package.Id <= 0 {
		return false
	}
	p.Weight = _package.Weight
	sqlStatement := `INSERT INTO public.packaging("package", sales_order, weight) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, p.Package, p.SalesOrder, p.Weight)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Packaging) deletePackaging() bool {
	if p.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	detailsPackaged := getSalesOrderDetailPackaged(p.Id)
	for i := 0; i < len(detailsPackaged); i++ {
		ok := detailsPackaged[i].deleteSalesOrderDetailPackaged(false)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM packaging WHERE id=$1`
	_, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addWeightPackaging(packagingId int32, weight float32) bool {
	sqlStatement := `UPDATE packaging SET weight = weight + $2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, packagingId, weight)
	rows, _ := res.RowsAffected()

	return rows > 0 && err == nil
}
