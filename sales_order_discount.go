package main

type SalesOrderDiscount struct {
	Id               int32   `json:"id"`
	Order            int32   `json:"order"`
	Name             string  `json:"name"`
	ValueTaxIncluded float32 `json:"valueTaxIncluded"`
	ValueTaxExcluded float32 `json:"valueTaxExcluded"`
}

func getSalesOrderDiscounts(orderId int32) []SalesOrderDiscount {
	var discounts []SalesOrderDiscount = make([]SalesOrderDiscount, 0)
	sqlStatement := `SELECT * FROM public.sales_order_discount WHERE "order" = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return discounts
	}
	for rows.Next() {
		d := SalesOrderDiscount{}
		rows.Scan(&d.Id, &d.Order, &d.Name, &d.ValueTaxIncluded, &d.ValueTaxExcluded)
		discounts = append(discounts, d)
	}

	return discounts
}

func getSalesOrderDiscountsRow(discountId int32) SalesOrderDiscount {
	sqlStatement := `SELECT * FROM public.sales_order_discount WHERE id = $1`
	row := db.QueryRow(sqlStatement, discountId)
	if row.Err() != nil {
		return SalesOrderDiscount{}
	}

	d := SalesOrderDiscount{}
	row.Scan(&d.Id, &d.Order, &d.Name, &d.ValueTaxIncluded, &d.ValueTaxExcluded)

	return d
}

func (d *SalesOrderDiscount) isValid() bool {
	return !(d.Order <= 0 || len(d.Name) == 0 || len(d.Name) > 100 || d.ValueTaxIncluded <= 0 || d.ValueTaxExcluded <= 0)
}

func (d *SalesOrderDiscount) insertSalesOrderDiscount() bool {
	if !d.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `INSERT INTO public.sales_order_discount("order", name, value_tax_included, value_tax_excluded) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, d.Order, d.Name, d.ValueTaxIncluded, d.ValueTaxExcluded)
	if err != nil {
		trans.Rollback()
		return false
	}

	ok := addDiscountsSalesOrder(d.Order, d.ValueTaxExcluded)
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

func (d *SalesOrderDiscount) deleteSalesOrderDiscount() bool {
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
	if inMemoryDiscount.Id <= 0 {
		trans.Rollback()
		return false
	}

	sqlStatement := `DELETE FROM public.sales_order_discount WHERE id = $1`
	res, err := db.Exec(sqlStatement, d.Id)
	if err != nil {
		trans.Rollback()
		return false
	}

	ok := addDiscountsSalesOrder(inMemoryDiscount.Order, -inMemoryDiscount.ValueTaxExcluded)
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
