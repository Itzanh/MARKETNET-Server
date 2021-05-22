package main

import "database/sql"

type SalesOrderDetailPackaged struct {
	OrderDetail int32  `json:"orderDetail"`
	ProductName string `json:"productName"`
	Packaging   int32  `json:"packaging"`
	Quantity    int32  `json:"quantity"`
}

func getSalesOrderDetailPackaged(packagingId int32) []SalesOrderDetailPackaged {
	var packaged []SalesOrderDetailPackaged = make([]SalesOrderDetailPackaged, 0)
	sqlStatement := `SELECT * FROM public.sales_order_detail_packaged WHERE packaging=$1 ORDER BY order_detail ASC`
	rows, err := db.Query(sqlStatement, packagingId)
	if err != nil {
		return packaged
	}
	for rows.Next() {
		p := SalesOrderDetailPackaged{}
		rows.Scan(&p.OrderDetail, &p.Packaging, &p.Quantity)
		detail := getSalesOrderDetailRow(p.OrderDetail)
		p.ProductName = getNameProduct(detail.Product)
		packaged = append(packaged, p)
	}

	return packaged
}

func getSalesOrderDetailPackagedRow(orderDetailId int32, packagingId int32) SalesOrderDetailPackaged {
	sqlStatement := `SELECT * FROM public.sales_order_detail_packaged WHERE packaging=$1 AND order_detail=$2`
	row := db.QueryRow(sqlStatement, packagingId, orderDetailId)
	if row.Err() != nil {
		return SalesOrderDetailPackaged{}
	}

	p := SalesOrderDetailPackaged{}
	row.Scan(&p.OrderDetail, &p.Packaging, &p.Quantity)

	return p
}

func (p *SalesOrderDetailPackaged) isValid() bool {
	return !(p.OrderDetail <= 0 || p.Packaging <= 0 || p.Quantity <= 0)
}

func (p *SalesOrderDetailPackaged) insertSalesOrderDetailPackaged() bool {
	if !p.isValid() {
		return false
	}

	detail := getSalesOrderDetailRow(p.OrderDetail)
	if detail.QuantityPendingPackaging <= 0 || p.Quantity > detail.QuantityPendingPackaging {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.sales_order_detail_packaged(order_detail, packaging, quantity) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, p.OrderDetail, p.Packaging, p.Quantity)
	if err != nil {
		trans.Rollback()
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		trans.Rollback()
		return false
	}

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetail, -p.Quantity)
	if !ok {
		trans.Rollback()
		return false
	}

	product := getProductRow(detail.Product)
	ok = addWeightPackaging(p.Packaging, product.Weight)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func (p *SalesOrderDetailPackaged) deleteSalesOrderDetailPackaged(openTransaction bool) bool {
	if p.OrderDetail <= 0 || p.Packaging <= 0 {
		return false
	}

	inMemoryPackage := getSalesOrderDetailPackagedRow(p.OrderDetail, p.Packaging)
	if inMemoryPackage.OrderDetail <= 0 || inMemoryPackage.Packaging <= 0 {
		return false
	}

	var trans *sql.Tx
	if openTransaction {
		///
		trn, transErr := db.Begin()
		if transErr != nil {
			return false
		}
		trans = trn
		///
	}

	sqlStatement := `DELETE FROM sales_order_detail_packaged WHERE order_detail=$1 AND packaging=$2`
	res, err := db.Exec(sqlStatement, p.OrderDetail, p.Packaging)
	if err != nil {
		if openTransaction {
			trans.Rollback()
		}
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		if openTransaction {
			trans.Rollback()
		}
		return false
	}

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetail, inMemoryPackage.Quantity)
	if !ok {
		if openTransaction {
			trans.Rollback()
		}
		return false
	}

	detail := getSalesOrderDetailRow(p.OrderDetail)
	product := getProductRow(detail.Product)
	ok = addWeightPackaging(p.Packaging, -product.Weight)
	if !ok {
		if openTransaction {
			trans.Rollback()
		}
		return false
	}

	if openTransaction {
		///
		transErr := trans.Commit()
		return transErr == nil
		///
	} else {
		return true
	}

}
