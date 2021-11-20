package main

type SalesOrderDiscount struct {
	Id               int32   `json:"id"`
	Order            int64   `json:"order"`
	Name             string  `json:"name"`
	ValueTaxIncluded float64 `json:"valueTaxIncluded"`
	ValueTaxExcluded float64 `json:"valueTaxExcluded"`
	enterprise       int32
}

func getSalesOrderDiscounts(orderId int32, enterpriseId int32) []SalesOrderDiscount {
	var discounts []SalesOrderDiscount = make([]SalesOrderDiscount, 0)
	sqlStatement := `SELECT * FROM public.sales_order_discount WHERE "order"=$1 AND enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return discounts
	}
	for rows.Next() {
		d := SalesOrderDiscount{}
		rows.Scan(&d.Id, &d.Order, &d.Name, &d.ValueTaxIncluded, &d.ValueTaxExcluded, &d.enterprise)
		discounts = append(discounts, d)
	}

	return discounts
}

func getSalesOrderDiscountsRow(discountId int32) SalesOrderDiscount {
	sqlStatement := `SELECT * FROM public.sales_order_discount WHERE id=$1`
	row := db.QueryRow(sqlStatement, discountId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return SalesOrderDiscount{}
	}

	d := SalesOrderDiscount{}
	row.Scan(&d.Id, &d.Order, &d.Name, &d.ValueTaxIncluded, &d.ValueTaxExcluded, &d.enterprise)

	return d
}

func (d *SalesOrderDiscount) isValid() bool {
	return !(d.Order <= 0 || len(d.Name) == 0 || len(d.Name) > 100 || d.ValueTaxIncluded <= 0 || d.ValueTaxExcluded <= 0)
}

func (d *SalesOrderDiscount) insertSalesOrderDiscount(userId int32) bool {
	if !d.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.sales_order_discount("order", name, value_tax_included, value_tax_excluded, enterprise) VALUES ($1, $2, $3, $4, $5)`
	res, err := db.Exec(sqlStatement, d.Order, d.Name, d.ValueTaxIncluded, d.ValueTaxExcluded, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return false
	}

	ok := addDiscountsSalesOrder(d.enterprise, d.Order, userId, d.ValueTaxExcluded)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	if err != nil {
		return false
	}
	///

	return rows > 0
}

func (d *SalesOrderDiscount) deleteSalesOrderDiscount(userId int32) bool {
	if d.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	inMemoryDiscount := getSalesOrderDiscountsRow(d.Id)
	if inMemoryDiscount.Id <= 0 || inMemoryDiscount.enterprise != d.enterprise {
		trans.Rollback()
		return false
	}

	sqlStatement := `DELETE FROM public.sales_order_discount WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	ok := addDiscountsSalesOrder(d.enterprise, inMemoryDiscount.Order, userId, -inMemoryDiscount.ValueTaxExcluded)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	if err != nil {
		return false
	}
	///

	rows, _ := res.RowsAffected()
	return rows > 0
}
