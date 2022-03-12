package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ManufacturingOrder struct {
	Id                   int64      `json:"id"`
	OrderDetail          *int64     `json:"orderDetail"`
	Product              int32      `json:"product"`
	Type                 int32      `json:"type"`
	Uuid                 string     `json:"uuid"`
	DateCreated          time.Time  `json:"dateCreated"`
	DateLastUpdate       time.Time  `json:"dateLastUpdate"`
	Manufactured         bool       `json:"manufactured"`
	DateManufactured     *time.Time `json:"dateManufactured"`
	UserManufactured     *int32     `json:"userManufactured"`
	UserCreated          int32      `json:"userCreated"`
	TagPrinted           bool       `json:"tagPrinted"`
	DateTagPrinted       *time.Time `json:"dateTagPrinted"`
	Order                *int64     `json:"order"`
	UserTagPrinted       *int32     `json:"userTagPrinted"`
	TypeName             string     `json:"typeName"`
	ProductName          string     `json:"productName"`
	OrderName            string     `json:"orderName"`
	UserCreatedName      string     `json:"userCreatedName"`
	UserManufacturedName *string    `json:"userManufacturedName"`
	UserTagPrintedName   *string    `json:"userTagPrintedName"`
	Warehouse            string     `json:"warehouse"`
	WarehouseMovement    *int64     `json:"warehouseMovement"`
	QuantityManufactured int32      `json:"quantityManufactured"`
	complex              bool
	enterprise           int32
}

type ManufacturingPaginationQuery struct {
	PaginationQuery
	OrderTypeId int32      `json:"orderTypeId"`
	DateStart   *time.Time `json:"dateStart"`
	DateEnd     *time.Time `json:"dateEnd"`
	Status      string     `json:"status"` // "" = All, "M" = Manufactured, "N" = Not manufactured
}

func (q *ManufacturingPaginationQuery) isDefault() bool {
	return q.OrderTypeId == 0 && q.DateStart == nil && q.DateEnd == nil && q.Status == ""
}

type ManufacturingOrders struct {
	Rows                int64                `json:"rows"`
	ManufacturingOrders []ManufacturingOrder `json:"manufacturingOrders"`
}

func (q *ManufacturingPaginationQuery) getManufacturingOrder(enterpriseId int32) ManufacturingOrders {
	if q.isDefault() {
		return (q.PaginationQuery).getAllManufacturingOrders(enterpriseId)
	} else {
		return q.getManufacturingOrdersByType(enterpriseId)
	}
}

func (q *PaginationQuery) getAllManufacturingOrders(enterpriseId int32) ManufacturingOrders {
	mo := ManufacturingOrders{}
	mo.ManufacturingOrders = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=manufacturing_order.type),(SELECT name FROM product WHERE product.id=manufacturing_order.product),(SELECT order_name FROM sales_order WHERE sales_order.id=manufacturing_order.order),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_tag_printed) FROM public.manufacturing_order WHERE enterprise=$3 ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return mo
	}
	defer rows.Close()

	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.complex, &o.TypeName, &o.ProductName, &o.OrderName, &o.UserCreatedName, &o.UserManufacturedName, &o.UserTagPrintedName)
		mo.ManufacturingOrders = append(mo.ManufacturingOrders, o)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.manufacturing_order WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return mo
	}
	row.Scan(&mo.Rows)

	return mo
}

func (q *ManufacturingPaginationQuery) getManufacturingOrdersByType(enterpriseId int32) ManufacturingOrders {
	mo := ManufacturingOrders{}
	mo.ManufacturingOrders = make([]ManufacturingOrder, 0)

	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, q.Offset)
	interfaces = append(interfaces, q.Limit)
	interfaces = append(interfaces, enterpriseId)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=manufacturing_order.type),(SELECT name FROM product WHERE product.id=manufacturing_order.product),(SELECT order_name FROM sales_order WHERE sales_order.id=manufacturing_order.order),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=manufacturing_order.user_tag_printed) FROM public.manufacturing_order WHERE enterprise=$3`
	if q.OrderTypeId != 0 {
		sqlStatement += ` AND type=$` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, q.OrderTypeId)
	}
	if q.DateStart != nil {
		sqlStatement += ` AND date_created>=$` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, q.DateStart)
	}
	if q.DateEnd != nil {
		sqlStatement += ` AND date_created<=$` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, q.DateEnd)
	}
	if q.Status == "M" {
		sqlStatement += ` AND manufactured=true`
	} else if q.Status == "N" {
		sqlStatement += ` AND manufactured=false`
	}
	sqlStatement += ` ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		fmt.Println(err)
		return mo
	}
	defer rows.Close()

	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.complex, &o.TypeName, &o.ProductName, &o.OrderName, &o.UserCreatedName, &o.UserManufacturedName, &o.UserTagPrintedName)
		mo.ManufacturingOrders = append(mo.ManufacturingOrders, o)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.manufacturing_order WHERE type=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, q.OrderTypeId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return mo
	}
	row.Scan(&mo.Rows)

	return mo
}

func getManufacturingOrderRow(manufacturingOrderId int64) ManufacturingOrder {
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE id = $1`
	row := db.QueryRow(sqlStatement, manufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ManufacturingOrder{}
	}

	o := ManufacturingOrder{}
	row.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.complex)

	return o
}

func getManufacturingOrderRowTransaction(manufacturingOrderId int64, trans sql.Tx) ManufacturingOrder {
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE id = $1`
	row := trans.QueryRow(sqlStatement, manufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ManufacturingOrder{}
	}

	o := ManufacturingOrder{}
	row.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.complex)

	return o
}

func getManufacturingOrdersForStockPending(enterpriseId int32, productId int32) []ManufacturingOrder {
	var orders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	sqlStatement := `SELECT * FROM public.manufacturing_order WHERE enterprise=$1 AND product=$2 AND NOT manufactured AND order_detail IS NULL AND NOT complex ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, enterpriseId, productId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}
	defer rows.Close()

	for rows.Next() {
		o := ManufacturingOrder{}
		rows.Scan(&o.Id, &o.OrderDetail, &o.Product, &o.Type, &o.Uuid, &o.DateCreated, &o.DateLastUpdate, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.Order, &o.UserTagPrinted, &o.enterprise, &o.Warehouse, &o.WarehouseMovement, &o.QuantityManufactured, &o.complex)
		orders = append(orders, o)
	}

	return orders
}

func (o *ManufacturingOrder) isValid() bool {
	return !((o.OrderDetail != nil && *o.OrderDetail <= 0) || o.Product <= 0 || (o.Order != nil && *o.Order <= 0))
}

// ERROR CODES:
// 1. There is no manufacturing order type in the product
// 2. The product is deactivated
func (o *ManufacturingOrder) insertManufacturingOrder(userId int32, trans *sql.Tx) OkAndErrorCodeReturn {
	if !o.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	// generate uuid
	o.Uuid = uuid.New().String()

	product := getProductRow(o.Product)
	if product.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	// don't allow deactivated products
	if product.Off {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}
	// get type if it's not specified
	if o.Type <= 0 {
		if !product.Manufacturing || product.ManufacturingOrderType == nil {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
		o.Type = *product.ManufacturingOrderType
	}

	// get quantity manufactured from the type if it's not specified
	mType := getManufacturingOrderTypeRow(o.Type)
	if mType.Id <= 0 || mType.enterprise != o.enterprise {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	o.QuantityManufactured = mType.QuantityManufactured

	// set the warehouse
	if len(o.Warehouse) == 0 {
		s := getSettingsRecordById(o.enterprise)
		o.Warehouse = s.DefaultWarehouse
	}

	sqlStatement := `INSERT INTO public.manufacturing_order(order_detail, product, type, uuid, user_created, "order", enterprise, warehouse, quantity_manufactured, complex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	row := trans.QueryRow(sqlStatement, o.OrderDetail, o.Product, o.Type, o.Uuid, o.UserCreated, o.Order, o.enterprise, o.Warehouse, o.QuantityManufactured, o.complex)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	var manufacturingOrderId int64
	row.Scan(&manufacturingOrderId)
	o.Id = manufacturingOrderId

	if o.OrderDetail != nil && *o.OrderDetail > 0 {
		sqlStatement = `UPDATE sales_order_detail SET status = 'D' WHERE id = $1`
		_, err := trans.Exec(sqlStatement, o.OrderDetail)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}

		insertTransactionalLog(o.enterprise, "sales_order_detail", int(*o.OrderDetail), userId, "U")

		ok := setSalesOrderState(o.enterprise, *o.Order, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	ok := addQuantityPendingManufacture(o.Product, o.Warehouse, 1, o.enterprise, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	if manufacturingOrderId > 0 {
		insertTransactionalLog(o.enterprise, "manufacturing_order", int(manufacturingOrderId), userId, "I")
	}

	return OkAndErrorCodeReturn{Ok: manufacturingOrderId > 0}
}

func (o *ManufacturingOrder) deleteManufacturingOrder(userId int32, trans *sql.Tx) bool {
	if o.Id <= 0 {
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

	inMemoryManufacturingOrder := getManufacturingOrderRow(o.Id)
	if inMemoryManufacturingOrder.Id <= 0 || inMemoryManufacturingOrder.enterprise != o.enterprise {
		return false
	}

	insertTransactionalLog(inMemoryManufacturingOrder.enterprise, "manufacturing_order", int(o.Id), userId, "D")

	sqlStatement := `DELETE FROM public.manufacturing_order WHERE id=$1`
	res, err := trans.Exec(sqlStatement, o.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	rows, _ := res.RowsAffected()
	if rows <= 0 {
		trans.Rollback()
		return false
	}

	if inMemoryManufacturingOrder.OrderDetail != nil && *inMemoryManufacturingOrder.OrderDetail > 0 && inMemoryManufacturingOrder.Order != nil && *inMemoryManufacturingOrder.Order > 0 {
		sqlStatement = `UPDATE sales_order_detail SET status = 'C' WHERE id = $1`
		_, err = trans.Exec(sqlStatement, inMemoryManufacturingOrder.OrderDetail)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(o.enterprise, "sales_order_detail", int(*inMemoryManufacturingOrder.OrderDetail), userId, "U")

		ok := setSalesOrderState(inMemoryManufacturingOrder.enterprise, *inMemoryManufacturingOrder.Order, userId, *trans)
		if !ok {
			return false
		}
	}

	ok := addQuantityPendingManufacture(inMemoryManufacturingOrder.Product, inMemoryManufacturingOrder.Warehouse, -1, inMemoryManufacturingOrder.enterprise, *trans)
	if !ok {
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

func toggleManufactuedManufacturingOrder(orderid int64, userId int32, enterpriseId int32) bool {
	if orderid <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	settings := getSettingsRecordById(enterpriseId)

	inMemoryManufacturingOrder := getManufacturingOrderRowTransaction(orderid, *trans)
	if inMemoryManufacturingOrder.enterprise != enterpriseId {
		return false
	}

	// validation
	if inMemoryManufacturingOrder.Manufactured && inMemoryManufacturingOrder.DateManufactured != nil && int64(time.Since(*inMemoryManufacturingOrder.DateManufactured).Seconds()) > int64(settings.UndoManufacturingOrderSeconds) {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order SET manufactured = NOT manufactured, date_manufactured = CASE WHEN NOT manufactured THEN current_timestamp(3) ELSE NULL END, user_manufactured = CASE WHEN NOT manufactured THEN ` + strconv.Itoa(int(userId)) + ` ELSE NULL END WHERE id=$1`
	res, err := trans.Exec(sqlStatement, orderid)
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

	insertTransactionalLog(inMemoryManufacturingOrder.enterprise, "manufacturing_order", int(orderid), userId, "U")

	inMemoryManufacturingOrder = getManufacturingOrderRowTransaction(orderid, *trans)
	if inMemoryManufacturingOrder.Id <= 0 {
		return false
	}
	if inMemoryManufacturingOrder.OrderDetail != nil && *inMemoryManufacturingOrder.OrderDetail > 0 {
		var status string
		if inMemoryManufacturingOrder.Manufactured {
			// are all the manufacturing orders manufactured?

			sqlStatement := `SELECT COUNT(id) FROM public.manufacturing_order WHERE order_detail = $1 AND manufactured`
			row := trans.QueryRow(sqlStatement, inMemoryManufacturingOrder.OrderDetail)
			if row.Err() != nil {
				log("DB", row.Err().Error())
				trans.Rollback()
				return false
			}

			var manufacturedOrders int32
			row.Scan(&manufacturedOrders)

			orderDetail := getSalesOrderDetailRow(*inMemoryManufacturingOrder.OrderDetail)

			if manufacturedOrders >= orderDetail.Quantity {
				status = "E"
			} else {
				status = "D"
			}
		} else {
			status = "D"
		}
		sqlStatement = `UPDATE sales_order_detail SET status = $2 WHERE id = $1`
		_, err = trans.Exec(sqlStatement, inMemoryManufacturingOrder.OrderDetail, status)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.enterprise, "sales_order_detail", int(*inMemoryManufacturingOrder.OrderDetail), userId, "U")

		sqlStatement = `SELECT "order" FROM sales_order_detail WHERE id = $1`
		row := trans.QueryRow(sqlStatement, inMemoryManufacturingOrder.OrderDetail)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			trans.Rollback()
			return false
		}

		var orderId int64
		row.Scan(&orderId)
		if orderId <= 0 {
			trans.Rollback()
			return false
		}

		ok := setSalesOrderState(enterpriseId, orderId, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// Create / delete warehouse movement
	if inMemoryManufacturingOrder.Manufactured {
		movement := WarehouseMovement{
			Warehouse:  inMemoryManufacturingOrder.Warehouse,
			Product:    inMemoryManufacturingOrder.Product,
			Quantity:   inMemoryManufacturingOrder.QuantityManufactured,
			Type:       "I", // Input
			enterprise: enterpriseId,
		}
		ok := movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		sqlStatement := `UPDATE public.manufacturing_order SET warehouse_movement=$2 WHERE id=$1`
		_, err := trans.Exec(sqlStatement, inMemoryManufacturingOrder.Id, movement.Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.enterprise, "manufacturing_order", int(inMemoryManufacturingOrder.Id), userId, "U")

		ok = addQuantityPendingManufacture(inMemoryManufacturingOrder.Product, inMemoryManufacturingOrder.Warehouse, -1, inMemoryManufacturingOrder.enterprise, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	} else {
		sqlStatement := `UPDATE public.manufacturing_order SET warehouse_movement=NULL WHERE id=$1`
		_, err := trans.Exec(sqlStatement, inMemoryManufacturingOrder.Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(inMemoryManufacturingOrder.enterprise, "manufacturing_order", int(inMemoryManufacturingOrder.Id), userId, "U")

		movement := getWarehouseMovementRow(*inMemoryManufacturingOrder.WarehouseMovement)
		ok := movement.deleteWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		ok = addQuantityPendingManufacture(inMemoryManufacturingOrder.Product, inMemoryManufacturingOrder.Warehouse, 1, inMemoryManufacturingOrder.enterprise, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	// manufacture / undo complex manufacturing orders
	setComplexManufacturingOrderManufacturingOrderManufactured(inMemoryManufacturingOrder.Id, inMemoryManufacturingOrder.Manufactured, inMemoryManufacturingOrder.enterprise, userId, trans)

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func manufacturingOrderAllSaleOrder(saleOrderId int64, userId int32, enterpriseId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.enterprise != enterpriseId {
		return false
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.enterprise)

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

			product := getProductRow(orderDetail.Product)
			if product.Id <= 0 || !product.Manufacturing || product.ManufacturingOrderType == nil || *product.ManufacturingOrderType == 0 {
				continue
			}
			manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderType)
			if manufacturingOrderType.Id <= 0 || manufacturingOrderType.QuantityManufactured <= 0 || manufacturingOrderType.Complex {
				continue
			}

			for j := 0; j < int(orderDetail.Quantity); j += int(manufacturingOrderType.QuantityManufactured) {
				o := ManufacturingOrder{}
				o.Product = orderDetail.Product
				o.OrderDetail = &orderDetail.Id
				o.Order = &saleOrder.Id
				o.UserCreated = userId
				o.enterprise = enterpriseId
				o.Warehouse = saleOrder.Warehouse
				ok := o.insertManufacturingOrder(userId, trans).Ok
				if !ok {
					trans.Rollback()
					return false
				}
			}
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type ManufacturingOrderGenerate struct {
	Selection []ManufacturingOrderGenerateSelection `json:"selection"`
}

type ManufacturingOrderGenerateSelection struct {
	OrderId  int64 `json:"orderId"`
	Id       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

func (orderInfo *ManufacturingOrderGenerate) manufacturingOrderPartiallySaleOrder(userId int32, enterpriseId int32) bool {
	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	for i := 0; i < len(orderInfo.Selection); i++ {
		orderInfoSelection := orderInfo.Selection[i]
		// get the sale order and it's details
		saleOrder := getSalesOrderRow(orderInfoSelection.OrderId)
		if saleOrder.Id <= 0 || saleOrder.enterprise != enterpriseId || len(orderInfo.Selection) == 0 {
			trans.Rollback()
			return false
		}

		// get the details
		orderDetail := getSalesOrderDetailRow(orderInfoSelection.Id)
		if orderDetail.Id <= 0 || orderDetail.Order != orderInfoSelection.OrderId || orderInfoSelection.Quantity == 0 || orderInfoSelection.Quantity > orderDetail.Quantity {
			trans.Rollback()
			return false
		}
		if orderDetail.Status == "C" {
			o := ManufacturingOrder{}
			o.Product = orderDetail.Product
			o.OrderDetail = &orderDetail.Id
			o.Order = &orderDetail.Order
			o.UserCreated = userId
			o.enterprise = enterpriseId
			o.Warehouse = saleOrder.Warehouse
			ok := o.insertManufacturingOrder(userId, trans).Ok
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

func manufacturingOrderTagPrinted(orderId int64, userId int32, enterpriseId int32) bool {
	if orderId <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order SET tag_printed = true, date_tag_printed = current_timestamp(3), user_tag_printed = $2 WHERE id=$1 AND enterprise=$3`
	_, err := db.Exec(sqlStatement, orderId, userId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
	}

	if err == nil {
		insertTransactionalLog(enterpriseId, "manufacturing_order", int(orderId), userId, "U")
	}

	return err == nil
}

type MultipleManufacturingOrders struct {
	Order   ManufacturingOrder `json:"order"`
	Quantiy int                `json:"quantity"`
}

func (o *MultipleManufacturingOrders) insertMultipleManufacturingOrders(userId int32) OkAndErrorCodeReturn {
	if !o.Order.isValid() || o.Quantiy <= 0 || o.Quantiy > 10000 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	trans, transErr := db.Begin()
	if transErr != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}

	for i := 0; i < o.Quantiy; i++ {
		ok := o.Order.insertManufacturingOrder(userId, trans)
		if !ok.Ok {
			trans.Rollback()
			return ok
		}
	}

	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}
}
