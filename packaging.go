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
	Pallet          *int32                     `json:"pallet"`
}

func getPackaging(salesOrderId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	sqlStatement := `SELECT * FROM public.packaging WHERE sales_order=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, salesOrderId)
	if err != nil {
		log("DB", err.Error())
		return packaging
	}
	for rows.Next() {
		p := Packaging{}
		rows.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet)
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id)

		_package := getPackagesRow(p.Package)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging = append(packaging, p)
	}

	return packaging
}

func getPackagingByShipping(shippingId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	sqlStatement := `SELECT * FROM public.packaging WHERE shipping=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, shippingId)
	if err != nil {
		log("DB", err.Error())
		return packaging
	}
	for rows.Next() {
		p := Packaging{}
		rows.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet)
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id)

		_package := getPackagesRow(p.Package)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging = append(packaging, p)
	}

	return packaging
}

func getPackagingRow(packagingId int32) Packaging {
	sqlStatement := `SELECT * FROM public.packaging WHERE id=$1`
	row := db.QueryRow(sqlStatement, packagingId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Packaging{}
	}

	p := Packaging{}
	row.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet)

	return p
}

func (p *Packaging) isValid() bool {
	return !(p.Package <= 0 || p.SalesOrder <= 0)
}

func (p *Packaging) insertPackaging() bool {
	if !p.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	_package := getPackagesRow(p.Package)
	if _package.Id <= 0 {
		trans.Rollback()
		return false
	}
	p.Weight = _package.Weight
	sqlStatement := `INSERT INTO public.packaging("package", sales_order, weight, pallet) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, p.Package, p.SalesOrder, p.Weight, p.Pallet)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		trans.Rollback()
		return false
	}

	s := getSalesOrderRow(p.SalesOrder)
	addQuantityStock(_package.Product, s.Warehouse, -1)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
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

	inMemoryPackaging := getPackagingRow(p.Id)
	if inMemoryPackaging.Id <= 0 {
		trans.Rollback()
		return false
	}

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
		log("DB", err.Error())
		return false
	}

	_package := getPackagesRow(inMemoryPackaging.Package)
	s := getSalesOrderRow(inMemoryPackaging.SalesOrder)
	addQuantityStock(_package.Product, s.Warehouse, 1)

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

	if err != nil {
		log("DB", err.Error())
	}

	return rows > 0 && err == nil
}
