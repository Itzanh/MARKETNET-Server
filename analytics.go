/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"time"
)

type MonthlySalesAmountQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *MonthlySalesAmountQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

// Monthly sales (amount)
type MonthlySalesAmount struct {
	Year   int16   `json:"year"`
	Month  int8    `json:"month"`
	Day    int8    `json:"day"`
	Amount float64 `json:"amount"`
}

// Monthly sales (amount)
func (q *MonthlySalesAmountQuery) monthlySalesAmount(enterpriseId int32) []MonthlySalesAmount {
	acounts := make([]MonthlySalesAmount, 0)
	var rows *sql.Rows
	var err error
	if q.isDefault() {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created),SUM(total_amount) FROM sales_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created),SUM(total_amount) FROM sales_order WHERE enterprise=$1 AND date_created >= $2 AND date_created <= $3 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),EXTRACT(DAY FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return acounts
	}
	defer rows.Close()

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
func (q *MonthlySalesAmountQuery) monthlySalesQuantity(enterpriseId int32) []MonthlySalesQuantity {
	quantity := make([]MonthlySalesQuantity, 0)
	var rows *sql.Rows
	var err error
	if q.isDefault() {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),COUNT(*) FROM sales_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),COUNT(*) FROM sales_order WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}
	defer rows.Close()

	for rows.Next() {
		q := MonthlySalesQuantity{}
		rows.Scan(&q.Year, &q.Month, &q.Quantity)
		quantity = append(quantity, q)
	}
	return quantity
}

// Sales of a product by months (quantity)
type SalesOfAProductQuantity struct {
	Year     int16           `json:"year"`
	Month    int8            `json:"month"`
	Quantity map[int32]int64 `json:"quantity"`
}

// Sales of a product by months (quantity)
func salesOfAProductQuantity(productIds []int32, enterpriseId int32) []SalesOfAProductQuantity {
	quantityResult := make([]SalesOfAProductQuantity, 0)
	var productId int32
	var year int16
	var month int8
	var quantity int64
	var found bool

	for i := 0; i < len(productIds); i++ {
		productId = productIds[i]
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),COUNT(*) FROM sales_order_detail WHERE product=$1 AND enterprise=$2 GROUP BY (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order)`
		rows, err := db.Query(sqlStatement, productId, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			return quantityResult
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&year, &month, &quantity)

			found = false
			for j := 0; j < len(quantityResult); j++ {
				if quantityResult[j].Year == year && quantityResult[j].Month == month {
					quantityResult[j].Quantity[productId] = quantity
					found = true
					break
				}
			}
			if !found {
				q := SalesOfAProductQuantity{
					Year:  year,
					Month: month,
					Quantity: map[int32]int64{
						productId: quantity,
					},
				}
				quantityResult = append(quantityResult, q)
			}
		}

	}
	return quantityResult
}

// Sales of a product by month (amount)
type SalesOfAProductAmount struct {
	Year   int16             `json:"year"`
	Month  int8              `json:"month"`
	Amount map[int32]float64 `json:"amount"`
}

// Sales of a product by month (amount)
func salesOfAProductAmount(productIds []int32, enterpriseId int32) []SalesOfAProductAmount {
	quantity := make([]SalesOfAProductAmount, 0)
	var productId int32
	var year int16
	var month int8
	var amount float64
	var found bool

	for i := 0; i < len(productIds); i++ {
		productId = productIds[i]
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),SUM(total_amount) FROM sales_order_detail WHERE product=$1 AND enterprise=$2 GROUP BY (SELECT EXTRACT(YEAR FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order),(SELECT EXTRACT(MONTH FROM date_created) FROM sales_order WHERE sales_order.id=sales_order_detail.order)`
		rows, err := db.Query(sqlStatement, productId, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			return quantity
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&year, &month, &amount)

			found = false
			for j := 0; j < len(quantity); j++ {
				if quantity[j].Year == year && quantity[j].Month == month {
					quantity[j].Amount[productId] = amount
					found = true
					break
				}
			}
			if !found {
				q := SalesOfAProductAmount{
					Year:  year,
					Month: month,
					Amount: map[int32]float64{
						productId: amount,
					},
				}
				quantity = append(quantity, q)
			}
		}
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
	defer rows.Close()

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
	defer rows.Close()

	for rows.Next() {
		d := DaysOfServicePurchaseOrders{}
		rows.Scan(&d.Year, &d.Month, &d.DaysAverage)
		days = append(days, d)
	}

	return days
}

type PurchaseOrdersByMonthQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *PurchaseOrdersByMonthQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

// Purchase orders by months (amount)
type PurchaseOrdersByMonthAmount struct {
	Year   int16   `json:"year"`
	Month  int8    `json:"month"`
	Amount float64 `json:"amount"`
}

// Purchase orders by months (amount)
func (q *PurchaseOrdersByMonthQuery) purchaseOrdersByMonthAmount(enterpriseId int32) []PurchaseOrdersByMonthAmount {
	amounts := make([]PurchaseOrdersByMonthAmount, 0)
	var rows *sql.Rows
	var err error
	if q.isDefault() {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),SUM(total_amount) FROM purchase_order WHERE enterprise=$1 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created),SUM(total_amount) FROM purchase_order WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created),EXTRACT(MONTH FROM date_created) ASC`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return amounts
	}
	defer rows.Close()

	for rows.Next() {
		a := PurchaseOrdersByMonthAmount{}
		rows.Scan(&a.Year, &a.Month, &a.Amount)
		amounts = append(amounts, a)
	}

	return amounts
}

// Payment methods of the sale orders
type PaymentMethodsSaleOrdersQuantityQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *PaymentMethodsSaleOrdersQuantityQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

type PaymentMethodsSaleOrders struct {
	Quantity []PaymentMethodsSaleOrdersQuantity `json:"quantity"`
	Amount   []PaymentMethodsSaleOrdersAmount   `json:"amount"`
}

type PaymentMethodsSaleOrdersQuantity struct {
	Quantity          int64  `json:"quantity"`
	PaymentMethod     int64  `json:"paymentMethod"`
	PaymentMethodName string `json:"paymentMethodName"`
}

type PaymentMethodsSaleOrdersAmount struct {
	Year   int16             `json:"year"`
	Month  int8              `json:"month"`
	Amount map[int32]float64 `json:"amount"`
}

// Payment methods of the sale orders
func (q *PaymentMethodsSaleOrdersQuantityQuery) paymentMethodsSaleOrdersAmount(enterpriseId int32) PaymentMethodsSaleOrders {
	var statistics PaymentMethodsSaleOrders = PaymentMethodsSaleOrders{}

	// GET QUANTITY

	quantity := make([]PaymentMethodsSaleOrdersQuantity, 0)
	var rows *sql.Rows
	var err error
	if q.isDefault() {
		sqlStatement := `SELECT COUNT(*),payment_method,(SELECT name FROM payment_method WHERE sales_order.payment_method=payment_method.id) FROM sales_order WHERE enterprise=$1 GROUP BY payment_method`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT COUNT(*),payment_method,(SELECT name FROM payment_method WHERE sales_order.payment_method=payment_method.id) FROM sales_order WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY payment_method`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	for rows.Next() {
		q := PaymentMethodsSaleOrdersQuantity{}
		rows.Scan(&q.Quantity, &q.PaymentMethod, &q.PaymentMethodName)
		quantity = append(quantity, q)
	}
	statistics.Quantity = quantity

	// GET AMOUNT

	amountResult := make([]PaymentMethodsSaleOrdersAmount, 0)
	var year int16
	var month int8
	var paymentMethodId int32
	var amount float64
	var found bool

	if q.isDefault() {
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),payment_method,SUM(total_amount) FROM sales_order WHERE enterprise=$1 GROUP BY (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),payment_method`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),payment_method,SUM(total_amount) FROM sales_order WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),payment_method`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&year, &month, &paymentMethodId, &amount)

		found = false
		for j := 0; j < len(amountResult); j++ {
			if amountResult[j].Year == year && amountResult[j].Month == month {
				amountResult[j].Amount[paymentMethodId] = amount
				found = true
				break
			}
		}
		if !found {
			q := PaymentMethodsSaleOrdersAmount{
				Year:  year,
				Month: month,
				Amount: map[int32]float64{
					paymentMethodId: amount,
				},
			}
			amountResult = append(amountResult, q)
		}
	}

	statistics.Amount = amountResult

	return statistics
}

// Sales by countries (amount)
type CountriesSaleOrdersQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *CountriesSaleOrdersQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

// Sales by countries (amount)

type CountriesSaleOrders struct {
	Amount  []CountriesSaleOrdersAmount        `json:"amount"`
	History []CountriesSaleOrdersAmountHistory `json:"history"`
}

type CountriesSaleOrdersAmount struct {
	Amount      float64 `json:"amount"`
	Country     int32   `json:"country"`
	CountryName string  `json:"countryName"`
}

type CountriesSaleOrdersAmountHistory struct {
	Year   int16             `json:"year"`
	Month  int8              `json:"month"`
	Amount map[int32]float64 `json:"amount"`
}

// Sales by countries (amount)
func (q *CountriesSaleOrdersQuery) countriesSaleOrdersAmount(enterpriseId int32) CountriesSaleOrders {
	var statistics CountriesSaleOrders = CountriesSaleOrders{}

	// global amount
	amountResult := make([]CountriesSaleOrdersAmount, 0)
	var rows *sql.Rows
	var err error
	sqlStatement := ``
	if q.isDefault() {
		sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address INNER JOIN country ON country.id=address.country WHERE sales_order.enterprise=$1 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement = `SELECT SUM(sales_order.total_amount),country.id,country.name FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address INNER JOIN country ON country.id=address.country WHERE sales_order.enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY country.id ORDER BY SUM(sales_order.total_amount) DESC LIMIT 10`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	for rows.Next() {
		q := CountriesSaleOrdersAmount{}
		rows.Scan(&q.Amount, &q.Country, &q.CountryName)
		amountResult = append(amountResult, q)
	}
	statistics.Amount = amountResult

	// history

	history := make([]CountriesSaleOrdersAmountHistory, 0)
	var year int16
	var month int8
	var countryId int32
	var amount float64
	var found bool

	if q.isDefault() {
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),country,SUM(total_amount) FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address WHERE sales_order.enterprise=$1 GROUP BY (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),country`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement := `SELECT (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),country,SUM(total_amount) FROM sales_order INNER JOIN address ON address.id=sales_order.shipping_address WHERE sales_order.enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY (SELECT EXTRACT(YEAR FROM date_created)),(SELECT EXTRACT(MONTH FROM date_created)),country`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&year, &month, &countryId, &amount)

		found = false
		for j := 0; j < len(history); j++ {
			if history[j].Year == year && history[j].Month == month {
				history[j].Amount[countryId] = amount
				found = true
				break
			}
		}
		if !found {
			q := CountriesSaleOrdersAmountHistory{
				Year:  year,
				Month: month,
				Amount: map[int32]float64{
					countryId: amount,
				},
			}
			history = append(history, q)
		}
	}

	statistics.History = history

	return statistics
}

type ManufacturingOrderCreatedManufacturedDailyQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *ManufacturingOrderCreatedManufacturedDailyQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

// Manufacturing orders created/manufactured daily
type ManufacturingOrderCreatedManufacturedDaily struct {
	Date     time.Time `json:"date"`
	Quantity int32     `json:"quantity"`
}

// Manufacturing orders created/manufactured daily
type ManufacturingOrderCreatedManufactured struct {
	Created      []ManufacturingOrderCreatedManufacturedDaily `json:"created"`
	Manufactured []ManufacturingOrderCreatedManufacturedDaily `json:"manufactured"`
}

// Manufacturing orders created/manufactured daily
func (q *ManufacturingOrderCreatedManufacturedDailyQuery) manufacturingOrderCreatedManufacturedDaily(enterpriseId int32) ManufacturingOrderCreatedManufactured {
	quantity := ManufacturingOrderCreatedManufactured{}
	quantity.Created = make([]ManufacturingOrderCreatedManufacturedDaily, 0)
	quantity.Manufactured = make([]ManufacturingOrderCreatedManufacturedDaily, 0)
	var sqlStatement string
	var rows *sql.Rows
	var err error

	if q.isDefault() {
		sqlStatement = `SELECT date_created::date,COUNT(*) FROM manufacturing_order WHERE enterprise=$1 GROUP BY date_created::date ORDER BY date_created::date`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement = `SELECT date_created::date,COUNT(*) FROM manufacturing_order WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY date_created::date ORDER BY date_created::date`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}
	defer rows.Close()

	for rows.Next() {
		d := ManufacturingOrderCreatedManufacturedDaily{}
		rows.Scan(&d.Date, &d.Quantity)
		quantity.Created = append(quantity.Created, d)
	}

	if q.isDefault() {
		sqlStatement = `SELECT date_manufactured::date,COUNT(*) FROM manufacturing_order WHERE date_manufactured IS NOT NULL AND enterprise=$1 GROUP BY date_manufactured::date ORDER BY date_manufactured::date`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement = `SELECT date_manufactured::date,COUNT(*) FROM manufacturing_order WHERE date_manufactured IS NOT NULL AND enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY date_manufactured::date ORDER BY date_manufactured::date`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return quantity
	}
	defer rows.Close()

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
	Quantity int32     `json:"quantity"`
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
	defer rows.Close()

	for rows.Next() {
		d := DailyShippingQuantity{}
		rows.Scan(&d.Date, &d.Quantity)
		quantity = append(quantity, d)
	}

	return quantity
}

type ShippingByCarriersQuery struct {
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
}

func (q *ShippingByCarriersQuery) isDefault() bool {
	return (q.DateStart == nil && q.DateEnd == nil)
}

// Shippings by carrier
type ShippingByCarriersResult struct {
	Quantity []ShippingByCarriers        `json:"quantity"`
	History  []ShippingByCarriersHistory `json:"history"`
}

type ShippingByCarriers struct {
	Quantity    int32  `json:"quantity"`
	Carrier     int32  `json:"carrier"`
	CarrierName string `json:"carrierName"`
}

type ShippingByCarriersHistory struct {
	Date     time.Time       `json:"date"`
	Quantity map[int32]int32 `json:"quantity"`
}

// Shippings by carrier
func (q *ShippingByCarriersQuery) shippingByCarriers(enterpriseId int32) ShippingByCarriersResult {
	var statistics ShippingByCarriersResult = ShippingByCarriersResult{}
	quantity := make([]ShippingByCarriers, 0)
	var sqlStatement string
	var rows *sql.Rows
	var err error

	if q.isDefault() {
		sqlStatement = `SELECT COUNT(*),carrier,(SELECT name FROM carrier WHERE carrier.id=shipping.carrier) FROM shipping WHERE enterprise=$1 GROUP BY carrier`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement = `SELECT COUNT(*),carrier,(SELECT name FROM carrier WHERE carrier.id=shipping.carrier) FROM shipping WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY carrier`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	for rows.Next() {
		s := ShippingByCarriers{}
		rows.Scan(&s.Quantity, &s.Carrier, &s.CarrierName)
		quantity = append(quantity, s)
	}
	statistics.Quantity = quantity

	// HISTORY

	history := make([]ShippingByCarriersHistory, 0)
	if q.isDefault() {
		sqlStatement = `SELECT date_sent::date,COUNT(*),carrier FROM shipping WHERE enterprise=$1 GROUP BY date_sent::date,carrier ORDER BY date_sent::date ASC`
		rows, err = db.Query(sqlStatement, enterpriseId)
	} else {
		sqlStatement = `SELECT date_sent::date,COUNT(*),carrier FROM shipping WHERE enterprise=$1 AND date_created>=$2 AND date_created<=$3 GROUP BY date_sent::date,carrier ORDER BY date_sent::date ASC`
		rows, err = db.Query(sqlStatement, enterpriseId, q.DateStart, q.DateEnd)
	}
	if err != nil {
		log("DB", err.Error())
		return statistics
	}
	defer rows.Close()

	var dateSend time.Time
	var count int32
	var carrierId int32
	var found bool
	for rows.Next() {
		rows.Scan(&dateSend, &count, &carrierId)

		found = false
		for j := 0; j < len(history); j++ {
			if history[j].Date == dateSend {
				history[j].Quantity[carrierId] = count
				found = true
				break
			}
		}
		if !found {
			s := ShippingByCarriersHistory{
				Date: dateSend,
				Quantity: map[int32]int32{
					carrierId: count,
				},
			}
			history = append(history, s)
		}
	}
	statistics.History = history

	return statistics
}

type BenefitsStatisticsQuery struct {
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
	Sales     bool      `json:"sales"`
	Purchases bool      `json:"purchases"`
}

type BenefitsStatistics struct {
	Sales     []BenefitsStatisticsValue `json:"sales"`
	Purchases []BenefitsStatisticsValue `json:"purchases"`
}

type BenefitsStatisticsValue struct {
	Year  int16   `json:"year"`
	Month uint8   `json:"month"`
	Value float64 `json:"value"`
}

func (q *BenefitsStatisticsQuery) benefitsStatistics(enterpriseId int32) BenefitsStatistics {
	var benefits BenefitsStatistics = BenefitsStatistics{}
	benefits.Sales = make([]BenefitsStatisticsValue, 0)
	benefits.Purchases = make([]BenefitsStatisticsValue, 0)

	if q.Purchases {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created), SUM(total_amount) FROM purchase_invoice WHERE date_created >= $1 AND date_created <= $2 AND enterprise = $3 GROUP BY EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created) ASC`
		rows, err := db.Query(sqlStatement, q.DateStart, q.DateEnd, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			return benefits
		}

		for rows.Next() {
			v := BenefitsStatisticsValue{}
			rows.Scan(&v.Year, &v.Month, &v.Value)
			benefits.Purchases = append(benefits.Purchases, v)
		}
	}

	if q.Sales {
		sqlStatement := `SELECT EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created), SUM(total_amount) FROM sales_invoice WHERE date_created >= $1 AND date_created <= $2 AND enterprise = $3 GROUP BY EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created) ORDER BY EXTRACT(YEAR FROM date_created), EXTRACT(MONTH FROM date_created) ASC`
		rows, err := db.Query(sqlStatement, q.DateStart, q.DateEnd, enterpriseId)
		if err != nil {
			log("DB", err.Error())
			return benefits
		}

		for rows.Next() {
			v := BenefitsStatisticsValue{}
			rows.Scan(&v.Year, &v.Month, &v.Value)
			benefits.Sales = append(benefits.Sales, v)
		}
	}

	return benefits
}
