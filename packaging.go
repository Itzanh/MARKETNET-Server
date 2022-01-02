package main

import (
	"database/sql"
	"fmt"
)

type Packaging struct {
	Id              int64                      `json:"id"`
	Package         int32                      `json:"package"`
	PackageName     string                     `json:"packageName"`
	SalesOrder      int64                      `json:"salesOrder"`
	Weight          float64                    `json:"weight"`
	Shipping        *int64                     `json:"shipping"`
	DetailsPackaged []SalesOrderDetailPackaged `json:"detailsPackaged"`
	Pallet          *int32                     `json:"pallet"`
	enterprise      int32
}

func getPackaging(salesOrderId int64, enterpriseId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	sqlStatement := `SELECT * FROM public.packaging WHERE sales_order=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, salesOrderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return packaging
	}
	defer rows.Close()

	for rows.Next() {
		p := Packaging{}
		rows.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet, &p.enterprise)
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id, enterpriseId)

		_package := getPackagesRow(p.Package)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging = append(packaging, p)
	}

	return packaging
}

func getPackagingByShipping(shippingId int64, enterpriseId int32) []Packaging {
	var packaging []Packaging = make([]Packaging, 0)
	sqlStatement := `SELECT * FROM public.packaging WHERE shipping=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, shippingId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return packaging
	}
	for rows.Next() {
		p := Packaging{}
		rows.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet, &p.enterprise)
		p.DetailsPackaged = getSalesOrderDetailPackaged(p.Id, enterpriseId)

		_package := getPackagesRow(p.Package)
		p.PackageName = _package.Name + " (" + fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth)) + ")"

		packaging = append(packaging, p)
	}

	return packaging
}

func getPackagingRow(packagingId int64) Packaging {
	sqlStatement := `SELECT * FROM public.packaging WHERE id=$1`
	row := db.QueryRow(sqlStatement, packagingId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Packaging{}
	}

	p := Packaging{}
	row.Scan(&p.Id, &p.Package, &p.SalesOrder, &p.Weight, &p.Shipping, &p.Pallet, &p.enterprise)

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
	sqlStatement := `INSERT INTO public.packaging("package", sales_order, weight, pallet, enterprise) VALUES ($1, $2, $3, $4, $5)`
	res, err := trans.Exec(sqlStatement, p.Package, p.SalesOrder, p.Weight, p.Pallet, p.enterprise)
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
	addQuantityStock(_package.Product, s.Warehouse, -1, p.enterprise, *trans)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func (p *Packaging) deletePackaging(enterpriseId int32, userId int32) bool {
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
	if inMemoryPackaging.Id <= 0 || inMemoryPackaging.enterprise != enterpriseId {
		trans.Rollback()
		return false
	}

	detailsPackaged := getSalesOrderDetailPackaged(p.Id, enterpriseId)
	for i := 0; i < len(detailsPackaged); i++ {
		detailsPackaged[i].enterprise = enterpriseId
		ok := detailsPackaged[i].deleteSalesOrderDetailPackaged(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM packaging WHERE id=$1 AND enterprise=$2`
	_, err := trans.Exec(sqlStatement, p.Id, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	_package := getPackagesRow(inMemoryPackaging.Package)
	s := getSalesOrderRow(inMemoryPackaging.SalesOrder)
	addQuantityStock(_package.Product, s.Warehouse, 1, p.enterprise, *trans)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addWeightPackaging(packagingId int64, weight float64, trans sql.Tx) bool {
	sqlStatement := `UPDATE packaging SET weight = weight + $2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, packagingId, weight)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}
	rows, _ := res.RowsAffected()

	return rows > 0 && err == nil
}
