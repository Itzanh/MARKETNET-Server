package main

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ManufacturingOrder struct {
	Id                   int64      `json:"id"`
	OrderDetail          *int32     `json:"orderDetail"`
	Product              int32      `json:"product"`
	Type                 int16      `json:"type"`
	Uuid                 string     `json:"uuid"`
	DateCreated          time.Time  `json:"dateCreated"`
	DateLastUpdate       time.Time  `json:"dateLastUpdate"`
	Manufactured         bool       `json:"manufactured"`
	DateManufactured     *time.Time `json:"dateManufactured"`
	UserManufactured     *int16     `json:"userManufactured"`
	UserCreated          int16      `json:"userCreated"`
	TagPrinted           bool       `json:"tagPrinted"`
	DateTagPrinted       *time.Time `json:"dateTagPrinted"`
	Order                *int32     `json:"order"`
	UserTagPrinted       *int16     `json:"userTagPrinted"`
	TypeName             string     `json:"typeName"`
	ProductName          string     `json:"productName"`
	OrderName            string     `json:"orderName"`
	UserCreatedName      string     `json:"userCreatedName"`
	UserManufacturedName *string    `json:"userManufacturedName"`
	UserTagPrintedName   *string    `json:"userTagPrintedName"`
}

func getManufacturingOrder(orderTypeId int16) []ManufacturingOrder {
	if orderTypeId == 0 {
		return getAllManufacturingOrders()
	} else {
		return getManufacturingOrdersByType(orderTypeId)
	}
}

func getAllManufacturingOrders() []ManufacturingOrder {
	var orders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=manufacturing_order.type),(SELECT name FROM product WHERE product.id=manufacturing_order.product),(SELECT order_name FROM sales_order WHERE sales_order.id=manufacturing_order.order),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_tag_printed) FROM public.manufacturing_order ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.TypeName, &o.ProductName, &o.OrderName, &o.UserCreatedName, &o.UserManufacturedName, &o.UserTagPrintedName)
		orders = append(orders, o)
	}

	return orders
}

func getManufacturingOrdersByType(orderTypeId int16) []ManufacturingOrder {
	var orders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=manufacturing_order.type),(SELECT name FROM product WHERE product.id=manufacturing_order.product),(SELECT order_name FROM sales_order WHERE sales_order.id=manufacturing_order.order),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_tag_printed) FROM public.manufacturing_order WHERE type = $1 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, orderTypeId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.TypeName, &o.ProductName, &o.OrderName, &o.UserCreatedName, &o.UserManufacturedName, &o.UserTagPrintedName)
		orders = append(orders, o)
	}

	return orders
}

func getManufacturingOrderRow(manufacturingOrderId int64) ManufacturingOrder {
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE id = $1`
	row := db.QueryRow(sqlStatement, manufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ManufacturingOrder{}
	}

	o := ManufacturingOrder{}
	row.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted)

	return o
}

func (o *ManufacturingOrder) isValid() bool {
	return !((o.OrderDetail != nil && *o.OrderDetail <= 0) || o.Product <= 0 || (o.Order != nil && *o.Order <= 0))
}

func (o *ManufacturingOrder) insertManufacturingOrder() bool {
	if !o.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	o.Uuid = uuid.New().String()
	if o.Type <= 0 {
		product := getProductRow(o.Product)
		if product.Id <= 0 || !product.Manufacturing || product.ManufacturingOrderType == nil {
			return false
		}
		o.Type = *product.ManufacturingOrderType
	}
	sqlStatement := `INSERT INTO public.manufacturing_order(order_detail, product, type, uuid, user_created, "order") VALUES ($1, $2, $3, $4, $5, $6)`
	res, err := db.Exec(sqlStatement, o.OrderDetail, o.Product, o.Type, o.Uuid, o.UserCreated, o.Order)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	if o.OrderDetail != nil && *o.OrderDetail > 0 {
		sqlStatement = `UPDATE sales_order_detail SET status = 'D' WHERE id = $1`
		res, err = db.Exec(sqlStatement, o.OrderDetail)
		if err != nil {
			log("DB", err.Error())
			return false
		}
		ok := setSalesOrderState(*o.Order)
		if !ok {
			return false
		}
	}

	s := getSettingsRecord()
	ok := addQuantityPendingManufacture(o.Product, s.DefaultWarehouse, 1)
	if !ok {
		return false
	}

	///
	transErr = trans.Commit()
	if transErr != nil {
		return false
	}
	///

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (o *ManufacturingOrder) deleteManufacturingOrder() bool {
	if o.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	inMemoryManufacturingOrder := getManufacturingOrderRow(o.Id)
	if inMemoryManufacturingOrder.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.manufacturing_order WHERE id=$1`
	res, err := db.Exec(sqlStatement, o.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows <= 0 {
		trans.Rollback()
		return false
	}

	if inMemoryManufacturingOrder.OrderDetail != nil && *inMemoryManufacturingOrder.OrderDetail > 0 {
		sqlStatement = `UPDATE sales_order_detail SET status = 'C' WHERE id = $1`
		_, err = db.Exec(sqlStatement, inMemoryManufacturingOrder.OrderDetail)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
		ok := setSalesOrderState(*o.Order)
		if !ok {
			return false
		}
	}

	s := getSettingsRecord()
	ok := addQuantityPendingManufacture(inMemoryManufacturingOrder.Product, s.DefaultWarehouse, -1)
	if !ok {
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func toggleManufactuedManufacturingOrder(orderid int64, userId int16) bool {
	if orderid <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `UPDATE public.manufacturing_order SET manufactured = NOT manufactured, date_manufactured = CASE WHEN NOT manufactured THEN current_timestamp(3) ELSE NULL END, user_manufactured = CASE WHEN NOT manufactured THEN ` + strconv.Itoa(int(userId)) + ` ELSE NULL END WHERE id=$1`
	res, err := db.Exec(sqlStatement, orderid)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		trans.Rollback()
		return false
	}

	inMemoryManufacturingOrder := getManufacturingOrderRow(orderid)
	if inMemoryManufacturingOrder.Id <= 0 {
		return false
	}
	if inMemoryManufacturingOrder.OrderDetail != nil && *inMemoryManufacturingOrder.OrderDetail > 0 {
		var status string
		if inMemoryManufacturingOrder.Manufactured {
			status = "E"
		} else {
			status = "D"
		}
		sqlStatement = `UPDATE sales_order_detail SET status = $2 WHERE id = $1`
		_, err = db.Exec(sqlStatement, inMemoryManufacturingOrder.OrderDetail, status)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		sqlStatement = `SELECT "order" FROM sales_order_detail WHERE id = $1`
		row := db.QueryRow(sqlStatement, inMemoryManufacturingOrder.OrderDetail)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			trans.Rollback()
			return false
		}

		var orderId int32
		row.Scan(&orderId)
		if orderId <= 0 {
			trans.Rollback()
			return false
		}

		ok := setSalesOrderState(orderId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func manufacturingOrderAllSaleOrder(saleOrderId int32, userId int16) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	orderDetails := getSalesOrderDetail(saleOrderId)

	if saleOrder.Id <= 0 || len(orderDetails) == 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	for i := 0; i < len(orderDetails); i++ {
		if orderDetails[i].Status == "C" {
			orderDetail := orderDetails[i]
			o := ManufacturingOrder{}
			o.Product = orderDetail.Product
			o.OrderDetail = &orderDetail.Id
			o.Order = &saleOrder.Id
			o.UserCreated = userId
			ok := o.insertManufacturingOrder()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func (orderInfo *OrderDetailGenerate) manufacturingOrderPartiallySaleOrder(userId int16) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(orderInfo.OrderId)
	if saleOrder.Id <= 0 || len(orderInfo.Selection) == 0 {
		return false
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(orderInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(orderInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.Order != orderInfo.OrderId || orderInfo.Selection[i].Quantity == 0 || orderInfo.Selection[i].Quantity > orderDetail.Quantity {
			return false
		}
		if orderDetail.Status == "C" {
			saleOrderDetails = append(saleOrderDetails, orderDetail)
		}
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	for i := 0; i < len(saleOrderDetails); i++ {
		orderDetail := saleOrderDetails[i]
		o := ManufacturingOrder{}
		o.Product = orderDetail.Product
		o.OrderDetail = &orderDetail.Id
		o.Order = &saleOrder.Id
		o.UserCreated = userId
		ok := o.insertManufacturingOrder()
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func manufacturingOrderTagPrinted(orderId int64, userId int16) bool {
	if orderId <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order SET tag_printed = true, date_tag_printed = current_timestamp(3), user_tag_printed = $2 WHERE id=$1`
	_, err := db.Exec(sqlStatement, orderId, userId)
	return err == nil
}
