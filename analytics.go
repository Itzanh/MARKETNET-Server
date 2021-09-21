package main

import (
	"database/sql"
	"time"
)

// Monthly sales (amount)
type MonthlySalesAmount struct {
	Year   int16   `json:"year"`
	Month  int8    `json:"month"`
	Day    int8    `json:"day"`
	Amount float64 `json:"amount"`
}

// Monthly sales (amount)
func monthlySalesAmount(year *int16, enterpriseId int32) []MonthlySalesAmount {
	acounts := make([]MonthlySalesAmount, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created),SUM(total_amount) FROM sales_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created),SUM(total_amount) FROM sales_order WHERE enterprise=$2 AND EXTRACT(YEAR FROM date_created)=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return acounts
	}

	for rows.Next() {
		m := MonthlySalesAmount{}
		rows.Scan(&m.Year, &m.Month, &m.Day, &m.Amount)
		acounts = append(acounts, m)
	}
	return acounts
}

// Monthly sales (quantity)
type MonthlySalesQuantity struct {
	Year     int16 `json:"year"`
	Month    int8  `json:"month"`
	Quantity int32 `json:"quantity"`
}

// Monthly sales (quantity)
func monthlySalesQuantity(year *int16, enterpriseId int32) []MonthlySalesQuantity {
	quantity := make([]MonthlySalesQuantity, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),COUNT(*) FROM sales_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),COUNT(*) FROM sales_order WHERE enterprise=$2 AND EXTRACT(YEAR FROM date_created)=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		q := MonthlySalesQuantity{}
		rows.Scan(&q.Year, &q.Month, &q.Quantity)
		quantity = append(quantity, q)
	}
	return quantity
}

// Sales of a product by months (quantity)
type SalesOfAProductQuantity struct {
	Year     int16 `json:"year"`
	Month    int8  `json:"month"`
	Quantity int64 `json:"quantity"`
}

// Sales of a product by months (quantity)
func salesOfAProductQuantity(productId int32, enterpriseId int32) []SalesOfAProductQuantity {
	quantity := make([]SalesOfAProductQuantity, 0)
	sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),COUNT(*) FROM sales_order_detail WHERE product=$1 AND enterprise=$2 GROUP BY (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order)`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		q := SalesOfAProductQuantity{}
		rows.Scan(&q.Year, &q.Month, &q.Quantity)
		quantity = append(quantity, q)
	}

	return quantity
}

// Sales of a product by month (amount)
type SalesOfAProductAmount struct {
	Year   int16   `json:"year"`
	Month  int8    `json:"month"`
	Amount float64 `json:"amount"`
}

// Sales of a product by month (amount)
func salesOfAProductAmount(productId int32, enterpriseId int32) []SalesOfAProductAmount {
	quantity := make([]SalesOfAProductAmount, 0)
	sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),SUM(total_amount) FROM sales_order_detail WHERE product=$1 AND enterprise=$2 GROUP BY (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order)`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		q := SalesOfAProductAmount{}
		rows.Scan(&q.Year, &q.Month, &q.Amount)
		quantity = append(quantity, q)
	}

	return quantity
}

// Days of service of the sale orders
type DaysOfServiceSaleOrders struct {
	Year        int16 `json:"year"`
	Month       int8  `json:"month"`
	DaysAverage uint8 `json:"daysAverage"`
}

// Days of service of the sale orders
func daysOfServiceSaleOrders(year *int16, enterpriseId int32) []DaysOfServiceSaleOrders {
	days := make([]DaysOfServiceSaleOrders, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT EXTRACT(YEAR FROM sales_order.date_created),EXTRACT(MONTH FROM sales_order.date_created),AVG(EXTRACT(DAY FROM (shipping.date_sent - sales_order.date_created))) FROM sales_order INNER JOIN shipping ON shipping.order=sales_order.id WHERE shipping.date_sent IS NOT NULL AND sales_order.enterprise=$1 GROUP BY EXTRACT(YEAR FROM sales_order.date_created),EXTRACT(MONTH FROM sales_order.date_created)`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM sales_order.date_created),EXTRACT(MONTH FROM sales_order.date_created),AVG(EXTRACT(DAY FROM (shipping.date_sent - sales_order.date_created))) FROM sales_order INNER JOIN shipping ON shipping.order=sales_order.id WHERE shipping.date_sent IS NOT NULL AND EXTRACT(YEAR FROM sales_order.date_created)=$1 AND sales_order.enterprise=$2 GROUP BY EXTRACT(YEAR FROM sales_order.date_created),EXTRACT(MONTH FROM sales_order.date_created)`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return days
	}

	for rows.Next() {
		d := DaysOfServiceSaleOrders{}
		rows.Scan(&d.Year, &d.Month, &d.DaysAverage)
		days = append(days, d)
	}

	return days
}

// Days of service of the purchase orders
type DaysOfServicePurchaseOrders struct {
	Year        int16 `json:"year"`
	Month       int8  `json:"month"`
	DaysAverage uint8 `json:"daysAverage"`
}

// Days of service of the purchase orders
func daysOfServicePurchaseOrders(year *int16, enterpriseId int32) []DaysOfServicePurchaseOrders {
	days := make([]DaysOfServicePurchaseOrders, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT EXTRACT(YEAR FROM purchase_order.date_created),EXTRACT(MONTH FROM purchase_order.date_created),AVG(EXTRACT(DAY FROM (purchase_delivery_note.date_created - purchase_order.date_created))) FROM purchase_order INNER JOIN purchase_order_detail ON purchase_order_detail.order=purchase_order.id INNER JOIN warehouse_movement ON warehouse_movement.purchase_order_detail=purchase_order_detail.id INNER JOIN purchase_delivery_note ON purchase_delivery_note.id=warehouse_movement.purchase_delivery_note WHERE purchase_order.lines_number>0 AND purchase_order.lines_number=purchase_order.delivery_note_lines AND purchase_order.enterprise=$1 GROUP BY EXTRACT(YEAR FROM purchase_order.date_created),EXTRACT(MONTH FROM purchase_order.date_created)`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM purchase_order.date_created),EXTRACT(MONTH FROM purchase_order.date_created),AVG(EXTRACT(DAY FROM (purchase_delivery_note.date_created - purchase_order.date_created))) FROM purchase_order INNER JOIN purchase_order_detail ON purchase_order_detail.order=purchase_order.id INNER JOIN warehouse_movement ON warehouse_movement.purchase_order_detail=purchase_order_detail.id INNER JOIN purchase_delivery_note ON purchase_delivery_note.id=warehouse_movement.purchase_delivery_note WHERE purchase_order.lines_number>0 AND purchase_order.lines_number=purchase_order.delivery_note_lines AND EXTRACT(YEAR FROM purchase_order.date_created)=$1 AND purchase_order.enterprise=$2 GROUP BY EXTRACT(YEAR FROM purchase_order.date_created),EXTRACT(MONTH FROM purchase_order.date_created))`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return days
	}

	for rows.Next() {
		d := DaysOfServicePurchaseOrders{}
		rows.Scan(&d.Year, &d.Month, &d.DaysAverage)
		days = append(days, d)
	}

	return days
}

// Purchase orders by months (amount)
type PurchaseOrdersByMonthAmount struct {
	Year   int16   `json:"year"`
	Month  int8    `json:"month"`
	Amount float64 `json:"amount"`
}

// Purchase orders by months (amount)
func purchaseOrdersByMonthAmount(year *int16, enterpriseId int32) []PurchaseOrdersByMonthAmount {
	amounts := make([]PurchaseOrdersByMonthAmount, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),SUM(total_amount) FROM purchase_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created)`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),SUM(total_amount) FROM purchase_order WHERE enterprise=$2 AND EXTRACT(YEAR FROM date_created)=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created)`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return amounts
	}

	for rows.Next() {
		a := PurchaseOrdersByMonthAmount{}
		rows.Scan(&a.Year, &a.Month, &a.Amount)
		amounts = append(amounts, a)
	}

	return amounts
}

// Payment methods of the sale orders
type PaymentMethodsSaleOrdersQuantity struct {
	Quantity          int64  `json:"quantity"`
	PaymentMethod     int16  `json:"paymentMethod"`
	PaymentMethodName string `json:"paymentMethodName"`
}

// Payment methods of the sale orders
func paymentMethodsSaleOrdersAmount(year *int16, enterpriseId int32) []PaymentMethodsSaleOrdersQuantity {
	quantity := make([]PaymentMethodsSaleOrdersQuantity, 0)
	var rows *sql.Rows
	var err error
	if year == nil {
		sqlStatement := `SELECT COUNT(*),payment_method,(SELECT name FROM payment_method WHERE sales_order.payment_method=payment_method.id) FROM sales_order WHERE enterprise=$1 GROUP BY payment_method`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT COUNT(*),payment_method,(SELECT name FROM payment_method WHERE sales_order.payment_method=payment_method.id) FROM sales_order WHERE enterprise=$2 AND EXTRACT(YEAR FROM date_created)=$1 GROUP BY payment_method`
		rows, err = db.Query(sqlStatement, year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		q := PaymentMethodsSaleOrdersQuantity{}
		rows.Scan(&q.Quantity, &q.PaymentMethod, &q.PaymentMethodName)
		quantity = append(quantity, q)
	}

	return quantity
}

// Sales by countries (amount)
type CountriesSaleOrdersQuery struct {
	Year            *int16 `json:"year"`
	ShippingAddress bool   `json:"shippingAddress"`
}

// Sales by countries (amount)
type CountriesSaleOrdersAmount struct {
	Amount      float64 `json:"quantity"`
	Country     int16   `json:"country"`
	CountryName string  `json:"countryName"`
}

// Sales by countries (amount)
func (q *CountriesSaleOrdersQuery) countriesSaleOrdersAmount(enterpriseId int32) []CountriesSaleOrdersAmount {
	quantity := make([]CountriesSaleOrdersAmount, 0)
	var rows *sql.Rows
	var err error
	sqlStatement := ``
	if q.Year == nil {
		if q.ShippingAddress {
			sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address INNER JOIN country ON country.id=address.country WHERE sales_order.enterprise=$1 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		} else {
			sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.billing_address INNER JOIN country ON country.id=address.country WHERE sales_order.enterprise=$1 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		}
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		if q.ShippingAddress {
			sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address INNER JOIN country ON country.id=address.country WHERE EXTRACT(YEAR FROM date_created)=$1 AND enterprise=$2 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		} else {
			sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.billing_address INNER JOIN country ON country.id=address.country WHERE EXTRACT(YEAR FROM date_created)=$1 AND enterprise=$2 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		}
		rows, err = db.Query(sqlStatement, q.Year, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		q := CountriesSaleOrdersAmount{}
		rows.Scan(&q.Amount, &q.Country, &q.CountryName)
		quantity = append(quantity, q)
	}

	return quantity
}

// Manufacturing orders created/manufactured daily
type ManufacturingOrderCreatedManufacturedDaily struct {
	Date     time.Time `json:"date"`
	Quantity uint16    `json:"quantity"`
}

// Manufacturing orders created/manufactured daily
type ManufacturingOrderCreatedManufactured struct {
	Created      []ManufacturingOrderCreatedManufacturedDaily `json:"created"`
	Manufactured []ManufacturingOrderCreatedManufacturedDaily `json:"manufactured"`
}

// Manufacturing orders created/manufactured daily
func manufacturingOrderCreatedManufacturedDaily(enterpriseId int32) ManufacturingOrderCreatedManufactured {
	quantity := ManufacturingOrderCreatedManufactured{}
	quantity.Created = make([]ManufacturingOrderCreatedManufacturedDaily, 0)
	quantity.Manufactured = make([]ManufacturingOrderCreatedManufacturedDaily, 0)
	sqlStatement := `SELECT date_created::date,COUNT(*) FROM manufacturing_order WHERE enterprise=$1 GROUP BY date_created::date ORDER BY date_created::date`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		d := ManufacturingOrderCreatedManufacturedDaily{}
		rows.Scan(&d.Date, &d.Quantity)
		quantity.Created = append(quantity.Created, d)
	}

	sqlStatement = `SELECT date_manufactured::date,COUNT(*) FROM manufacturing_order WHERE date_manufactured IS NOT NULL AND enterprise=$1 GROUP BY date_manufactured::date ORDER BY date_manufactured::date`
	rows, err = db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		d := ManufacturingOrderCreatedManufacturedDaily{}
		rows.Scan(&d.Date, &d.Quantity)
		quantity.Manufactured = append(quantity.Created, d)
	}

	return quantity
}

// Daily shipping (quantity)
type DailyShippingQuantity struct {
	Date     time.Time `json:"date"`
	Quantity uint16    `json:"quantity"`
}

// Daily shipping (quantity)
func dailyShippingQuantity(enterpriseId int32) []DailyShippingQuantity {
	quantity := make([]DailyShippingQuantity, 0)
	sqlStatement := `SELECT date_created::date,COUNT(*) FROM shipping WHERE enterprise=$1 GROUP BY date_created::date ORDER BY date_created::date`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		d := DailyShippingQuantity{}
		rows.Scan(&d.Date, &d.Quantity)
		quantity = append(quantity, d)
	}

	return quantity
}

// Shippings by carrier
type ShippingByCarriers struct {
	Quantity    int32  `json:"quantity"`
	Carrier     int16  `json:"carrier"`
	CarrierName string `json:"carrierName"`
}

// Shippings by carrier
func shippingByCarriers(enterpriseId int32) []ShippingByCarriers {
	quantity := make([]ShippingByCarriers, 0)
	sqlStatement := `SELECT COUNT(*),carrier,(SELECT name FROM carrier WHERE carrier.id=shipping.carrier) FROM shipping WHERE enterprise=$1 GROUP BY carrier`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return quantity
	}

	for rows.Next() {
		s := ShippingByCarriers{}
		rows.Scan(&s.Quantity, &s.Carrier, &s.CarrierName)
		quantity = append(quantity, s)
	}

	return quantity
}
