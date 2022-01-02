package main

import (
	"database/sql"
)

type SalesOrderDetailPackaged struct {
	OrderDetail int64  `json:"orderDetail"`
	ProductName string `json:"productName"`
	Packaging   int64  `json:"packaging"`
	Quantity    int32  `json:"quantity"`
	enterprise  int32
}

func getSalesOrderDetailPackaged(packagingId int64, enterpriseId int32) []SalesOrderDetailPackaged {
	var packaged []SalesOrderDetailPackaged = make([]SalesOrderDetailPackaged, 0)
	sqlStatement := `SELECT * FROM public.sales_order_detail_packaged WHERE packaging=$1 AND enterprise=$2 ORDER BY order_detail ASC`
	rows, err := db.Query(sqlStatement, packagingId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return packaged
	}
	defer rows.Close()

	for rows.Next() {
		p := SalesOrderDetailPackaged{}
		rows.Scan(&p.OrderDetail, &p.Packaging, &p.Quantity, &p.enterprise)
		detail := getSalesOrderDetailRow(p.OrderDetail)
		p.ProductName = getNameProduct(detail.Product, enterpriseId)
		packaged = append(packaged, p)
	}

	return packaged
}

func getSalesOrderDetailPackagedRow(orderDetailId int64, packagingId int64) SalesOrderDetailPackaged {
	sqlStatement := `SELECT * FROM public.sales_order_detail_packaged WHERE packaging=$1 AND order_detail=$2`
	row := db.QueryRow(sqlStatement, packagingId, orderDetailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesOrderDetailPackaged{}
	}

	p := SalesOrderDetailPackaged{}
	row.Scan(&p.OrderDetail, &p.Packaging, &p.Quantity, &p.enterprise)

	return p
}

func (p *SalesOrderDetailPackaged) isValid() bool {
	return !(p.OrderDetail <= 0 || p.Packaging <= 0 || p.Quantity <= 0)
}

func (p *SalesOrderDetailPackaged) insertSalesOrderDetailPackaged(userId int32) bool {
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

	sqlStatement := `INSERT INTO public.sales_order_detail_packaged(order_detail, packaging, quantity, enterprise) VALUES ($1, $2, $3, $4)`
	res, err := trans.Exec(sqlStatement, p.OrderDetail, p.Packaging, p.Quantity, p.enterprise)
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

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetail, -p.Quantity, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	product := getProductRow(detail.Product)
	ok = addWeightPackaging(p.Packaging, product.Weight, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func (p *SalesOrderDetailPackaged) deleteSalesOrderDetailPackaged(userId int32, trans *sql.Tx) bool {
	if p.OrderDetail <= 0 || p.Packaging <= 0 {
		return false
	}

	inMemoryPackage := getSalesOrderDetailPackagedRow(p.OrderDetail, p.Packaging)
	if inMemoryPackage.OrderDetail <= 0 || inMemoryPackage.enterprise != p.enterprise || inMemoryPackage.Packaging <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false
		}
		///
	}

	sqlStatement := `DELETE FROM sales_order_detail_packaged WHERE order_detail=$1 AND packaging=$2 AND enterprise=$3`
	res, err := trans.Exec(sqlStatement, p.OrderDetail, p.Packaging, p.enterprise)
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

	ok := addQuantityPendingPackagingSaleOrderDetail(p.OrderDetail, inMemoryPackage.Quantity, userId, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	detail := getSalesOrderDetailRow(p.OrderDetail)
	product := getProductRow(detail.Product)
	ok = addWeightPackaging(p.Packaging, -product.Weight, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false
		}
		///
	}
	return true
}

type SalesOrderDetailPackagedEAN13 struct {
	SalesOrder int64  `json:"salesOrder"`
	EAN13      string `json:"ean13"`
	Packaging  int64  `json:"packaging"`
	Quantity   int32  `json:"quantity"`
}

func (d *SalesOrderDetailPackagedEAN13) isValid() bool {
	return !(d.SalesOrder <= 0 || len(d.EAN13) != 13 || d.Packaging <= 0 || d.Quantity <= 0)
}

func (d *SalesOrderDetailPackagedEAN13) insertSalesOrderDetailPackagedEAN13(enterpriseId int32, userId int32) bool {
	if !d.isValid() {
		return false
	}

	sqlStatement := `SELECT sales_order_detail.id FROM sales_order_detail INNER JOIN product ON product.id=sales_order_detail.product WHERE sales_order_detail."order"=$1 AND sales_order_detail.quantity_pending_packaging>0 AND product.barcode=$2 AND product.enterprise=$3`
	row := db.QueryRow(sqlStatement, d.SalesOrder, d.EAN13, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var salesOrderDetailId int64
	row.Scan(&salesOrderDetailId)
	if salesOrderDetailId <= 0 {
		return false
	}

	p := SalesOrderDetailPackaged{}
	p.OrderDetail = salesOrderDetailId
	p.Packaging = d.Packaging
	p.Quantity = d.Quantity
	p.enterprise = enterpriseId

	return p.insertSalesOrderDetailPackaged(userId)
}
