package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TODO add transactional log
// TODO add recursivity

// TODO interactuar con pendiente de fabricar
// TODO mostrar en pedido de venta
// TODO mostrar en producto

// TODO modificar en lineas de pedido de compra

type ComplexManufacturingOrder struct {
	Id                         int64      `json:"id"`
	Type                       int32      `json:"type"`
	Manufactured               bool       `json:"manufactured"`
	DateManufactured           *time.Time `json:"dateManufactured"`
	UserManufactured           *int32     `json:"userManufactured"`
	QuantityPendingManufacture int32      `json:"quantityPendingManufacture"`
	QuantityManufactured       int32      `json:"quantityManufactured"`
	Warehouse                  string     `json:"warehouse"`
	DateCreated                time.Time  `json:"dateCreated"`
	Uuid                       string     `json:"uuid"`
	UserCreated                int32      `json:"userCreated"`
	TagPrinted                 bool       `json:"tagPrinted"`
	DateTagPrinted             *time.Time `json:"dateTagPrinted"`
	UserTagPrinted             *int32     `json:"userTagPrinted"`
	TypeName                   string     `json:"typeName"`
	UserCreatedName            string     `json:"userCreatedName"`
	UserManufacturedName       *string    `json:"userManufacturedName"`
	UserTagPrintedName         *string    `json:"userTagPrintedName"`
	enterprise                 int32
}

type ComplexManufacturingPaginationQuery struct {
	PaginationQuery
	OrderTypeId int32 `json:"orderTypeId"`
}

type ComplexManufacturingOrders struct {
	Rows                       int64                       `json:"rows"`
	ComplexManufacturingOrders []ComplexManufacturingOrder `json:"complexManufacturingOrder"`
}

func (q *ComplexManufacturingPaginationQuery) getComplexManufacturingOrder(enterpriseId int32) ComplexManufacturingOrders {
	if q.OrderTypeId == 0 {
		return (q.PaginationQuery).getAllComplexManufacturingOrders(enterpriseId)
	} else {
		return q.getComplexManufacturingOrdersByType(enterpriseId)
	}
}

func (q *PaginationQuery) getAllComplexManufacturingOrders(enterpriseId int32) ComplexManufacturingOrders {
	mo := ComplexManufacturingOrders{}
	mo.ComplexManufacturingOrders = make([]ComplexManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=complex_manufacturing_order.type),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_tag_printed) FROM public.complex_manufacturing_order WHERE enterprise=$3 ORDER BY date_created DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return mo
	}
	for rows.Next() {
		o := ComplexManufacturingOrder{}
		rows.Scan(&o.Id, &o.Type, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.enterprise, &o.QuantityPendingManufacture, &o.QuantityManufactured, &o.Warehouse, &o.DateCreated, &o.Uuid, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.UserTagPrinted, &o.TypeName, &o.UserCreatedName, &o.UserManufactured, &o.UserTagPrintedName)
		mo.ComplexManufacturingOrders = append(mo.ComplexManufacturingOrders, o)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.complex_manufacturing_order WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return mo
	}
	row.Scan(&mo.Rows)

	return mo
}

func (q *ComplexManufacturingPaginationQuery) getComplexManufacturingOrdersByType(enterpriseId int32) ComplexManufacturingOrders {
	mo := ComplexManufacturingOrders{}
	mo.ComplexManufacturingOrders = make([]ComplexManufacturingOrder, 0)
	sqlStatement := `SELECT *,(SELECT name FROM manufacturing_order_type WHERE manufacturing_order_type.id=complex_manufacturing_order.type),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_created),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_manufactured),(SELECT username FROM "user" WHERE "user".id=complex_manufacturing_order.user_tag_printed) FROM public.complex_manufacturing_order WHERE type=$1 AND enterprise=$4 ORDER BY date_created DESC OFFSET $2 LIMIT $3`
	rows, err := db.Query(sqlStatement, q.OrderTypeId, q.Offset, q.Limit, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return mo
	}

	for rows.Next() {
		o := ComplexManufacturingOrder{}
		rows.Scan(&o.Id, &o.Type, &o.Manufactured, &o.DateManufactured, &o.UserManufactured, &o.enterprise, &o.QuantityPendingManufacture, &o.QuantityManufactured, &o.Warehouse, &o.DateCreated, &o.Uuid, &o.UserCreated, &o.TagPrinted, &o.DateTagPrinted, &o.UserTagPrinted, &o.TypeName, &o.UserCreatedName, &o.UserManufactured, &o.UserTagPrintedName)
		mo.ComplexManufacturingOrders = append(mo.ComplexManufacturingOrders, o)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.complex_manufacturing_order WHERE type=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, q.OrderTypeId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return mo
	}
	row.Scan(&mo.Rows)

	return mo
}

func getComplexManufacturingOrderRow(complexManufacturingOrderId int64) ComplexManufacturingOrder {
	c := ComplexManufacturingOrder{}

	sqlStatement := `SELECT * FROM public.complex_manufacturing_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, complexManufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return c
	}

	row.Scan(&c.Id, &c.Type, &c.Manufactured, &c.DateManufactured, &c.UserManufactured, &c.enterprise, &c.QuantityPendingManufacture, &c.QuantityManufactured, &c.Warehouse, &c.DateCreated, &c.Uuid, &c.UserCreated, &c.TagPrinted, &c.DateTagPrinted, &c.UserTagPrinted)

	return c
}

// Specify a negative number to substract
// DOES NOT OPEN A TRANSACTION
func addQuantityPendingManufactureComplexManufacturingOrder(complexManufacturingOrderId int64, quantity int32) bool {
	sqlStatement := `UPDATE public.complex_manufacturing_order SET quantity_pending_manufacture=$2 WHERE id=$1`
	_, err := db.Exec(sqlStatement, complexManufacturingOrderId, quantity)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

// Specify a negative number to substract
// DOES NOT OPEN A TRANSACTION
func addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId int64, quantity int32) bool {
	sqlStatement := `UPDATE public.complex_manufacturing_order SET quantity_manufactured=$2 WHERE id=$1`
	_, err := db.Exec(sqlStatement, complexManufacturingOrderId, quantity)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func complexManufacturingOrderAllSaleOrder(saleOrderId int64, userId int32, enterpriseId int32) bool {
	saleOrder := getSalesOrderRow(saleOrderId)
	details := getSalesOrderDetail(saleOrderId, enterpriseId)

	return complexManufacturingOrerGeneration(saleOrderId, userId, enterpriseId, saleOrder, details)
}

func (orderInfo *OrderDetailGenerate) complexManufacturingOrderPartiallySaleOrder(userId int32, enterpriseId int32) bool {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(orderInfo.OrderId)
	if saleOrder.Id <= 0 || saleOrder.enterprise != enterpriseId || len(orderInfo.Selection) == 0 {
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

	return complexManufacturingOrerGeneration(orderInfo.OrderId, userId, enterpriseId, saleOrder, saleOrderDetails)
}

func complexManufacturingOrerGeneration(saleOrderId int64, userId int32, enterpriseId int32, saleOrder SaleOrder, details []SalesOrderDetail) bool {
	for i := 0; i < len(details); i++ {
		if details[i].Status != "C" {
			continue
		}
		orderDetail := details[i]

		product := getProductRow(orderDetail.Product)
		if product.Id <= 0 || !product.Manufacturing || product.ManufacturingOrderType == nil || *product.ManufacturingOrderType == 0 {
			continue
		}
		manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderType)
		if manufacturingOrderType.Id <= 0 || manufacturingOrderType.QuantityManufactured <= 0 || !manufacturingOrderType.Complex {
			continue
		}

		components := getManufacturingOrderTypeComponents(manufacturingOrderType.Id, enterpriseId)
		var component *ManufacturingOrderTypeComponents = nil
		for i := 0; i < len(components); i++ {
			if components[i].Type == "O" && components[i].Product == product.Id {
				component = &components[i]
				break
			}
		}
		if component == nil {
			return false
		}

		for j := 0; j < int(orderDetail.Quantity); j += int(component.Quantity) {
			cmo := ComplexManufacturingOrder{
				Type:       manufacturingOrderType.Id,
				enterprise: enterpriseId,
				Warehouse:  saleOrder.Warehouse,
			}
			ok := cmo.insertComplexManufacturingOrder(1, true)
			if !ok {
				return false
			}

			id := getPendingComplexManufacturingOrderOutputsWithoutSaleOrderDetail(product.Id)
			if id == nil || *id <= 0 {
				continue
			}

			sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET sale_order_detail=$2 WHERE id=$1`
			_, err := db.Exec(sqlStatement, id, orderDetail.Id)
			if err != nil {
				log("DB", err.Error())
				return false
			}

			sqlStatement = `UPDATE sales_order_detail SET status = 'D' WHERE id = $1`
			_, err = db.Exec(sqlStatement, orderDetail.Id)
			if err != nil {
				log("DB", err.Error())
				return false
			}

			insertTransactionalLog(enterpriseId, "sales_order_detail", int(orderDetail.Id), userId, "U")

			ok = setSalesOrderState(enterpriseId, saleOrderId, userId)
			if !ok {
				return false
			}
		} // for j := 0; j < int(orderDetail.Quantity); j += int(component.Quantity) {
	} // for i := 0; i < len(details); i++

	return true
}

func (c *ComplexManufacturingOrder) isValid() bool {
	return !(c.Type <= 0 || c.enterprise == 0)
}

func (c *ComplexManufacturingOrder) insertComplexManufacturingOrder(userId int32, beginTransaction bool) bool {
	if !c.isValid() {
		return false
	}

	var trans *sql.Tx
	var err error
	if beginTransaction {
		///
		trans, err = db.Begin()
		if err != nil {
			return false
		}
		///
	}

	// generate uuid
	c.Uuid = uuid.New().String()

	// set the warehouse
	if len(c.Warehouse) == 0 {
		s := getSettingsRecordById(c.enterprise)
		c.Warehouse = s.DefaultWarehouse
	}

	sqlStatement := `INSERT INTO public.complex_manufacturing_order(type, enterprise, warehouse, uuid, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := db.QueryRow(sqlStatement, c.Type, c.enterprise, c.Warehouse, c.Uuid, userId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		if beginTransaction {
			trans.Rollback()
		}
		return false
	}

	var complexManufacturingOrderId int64
	row.Scan(&complexManufacturingOrderId)
	complexManufacturingOrder := getComplexManufacturingOrderRow(complexManufacturingOrderId)

	components := getManufacturingOrderTypeComponents(c.Type, c.enterprise)

	var subOrders []ComplexManufacturingOrderManufacturingOrder = make([]ComplexManufacturingOrderManufacturingOrder, 0)
	for i := 0; i < len(components); i++ {
		if components[i].Type != "I" { // Only Input
			continue
		}

		manufacturingOrderTypeComponent := components[i]
		if manufacturingOrderTypeComponent.Id <= 0 || manufacturingOrderTypeComponent.enterprise != c.enterprise || manufacturingOrderTypeComponent.Quantity <= 0 {
			if beginTransaction {
				trans.Rollback()
			}
			return false
		}

		stock := getStockRow(manufacturingOrderTypeComponent.Product, complexManufacturingOrder.Warehouse, c.enterprise)
		if stock.QuantityAvaialbe >= manufacturingOrderTypeComponent.Quantity {
			// there is stock for the manufacturing, we make a manufacturing order to reserve the stock
			wm := WarehouseMovement{
				Warehouse:  complexManufacturingOrder.Warehouse,
				Product:    manufacturingOrderTypeComponent.Product,
				Quantity:   manufacturingOrderTypeComponent.Quantity,
				Type:       "O",
				enterprise: c.enterprise,
			}
			ok := wm.insertWarehouseMovement(userId)
			if !ok {
				if beginTransaction {
					trans.Rollback()
				}
				return false
			}

			c := ComplexManufacturingOrderManufacturingOrder{
				Type:                            "I",
				ComplexManufacturingOrder:       complexManufacturingOrder.Id,
				enterprise:                      complexManufacturingOrder.enterprise,
				WarehouseMovement:               &wm.Id,
				Product:                         manufacturingOrderTypeComponent.Product,
				ManufacturingOrderTypeComponent: manufacturingOrderTypeComponent.Id,
				Manufactured:                    true,
			}
			subOrders = append(subOrders, c)
		} else { // if stock.QuantityAvaialbe >= manufacturingOrderTypeComponent.Quantity {
			// the product is from a supplier or from manufacturing?
			product := getProductRow(manufacturingOrderTypeComponent.Product)
			if product.Manufacturing {
				// we search existing orders to make stock (without order and without complex order)
				manufacturingOrders := getManufacturingOrdersForStockPending(c.enterprise, manufacturingOrderTypeComponent.Product)
				var quantityManufacturedForStock int32 = 0
				for i := 0; i < len(manufacturingOrders); i++ {
					quantityManufacturedForStock += manufacturingOrders[0].QuantityManufactured
				}

				// associate with the existing orders
				if quantityManufacturedForStock >= manufacturingOrderTypeComponent.Quantity {
					var quantityAdded int32 = 0
					// the orders come sorted by date_created ASC, so the ones that are older are first (the ones we expect to manufacture before)
					for i := 0; i < len(manufacturingOrders); i++ {
						c := ComplexManufacturingOrderManufacturingOrder{
							Type:                            "I",
							ComplexManufacturingOrder:       complexManufacturingOrder.Id,
							enterprise:                      complexManufacturingOrder.enterprise,
							ManufacturingOrder:              &manufacturingOrders[i].Id,
							Product:                         manufacturingOrderTypeComponent.Product,
							ManufacturingOrderTypeComponent: manufacturingOrderTypeComponent.Id,
							Manufactured:                    false,
						}
						subOrders = append(subOrders, c)
						// set the manufacturing order as complex, so it doesn't count as stock
						sqlStatement := `UPDATE public.manufacturing_order SET complex=true WHERE id=$1`
						_, err := db.Exec(sqlStatement, manufacturingOrders[i].Id)
						if err != nil {
							log("DB", err.Error())
							if beginTransaction {
								trans.Rollback()
							}
							return false
						}
						// stop the loop as soon as we get enought quantity
						quantityAdded += manufacturingOrders[i].QuantityManufactured
						if quantityAdded >= manufacturingOrderTypeComponent.Quantity {
							break
						}
					}
				} else { // there are no stock orders, create a new one
					manufacturingOrderType := getManufacturingOrderTypeRow(*product.ManufacturingOrderType)
					for i := 0; i < int(manufacturingOrderTypeComponent.Quantity); i += int(manufacturingOrderType.QuantityManufactured) {
						mo := ManufacturingOrder{
							Product:    manufacturingOrderTypeComponent.Product,
							Type:       manufacturingOrderTypeComponent.ManufacturingOrderType,
							enterprise: complexManufacturingOrder.enterprise,
							Warehouse:  complexManufacturingOrder.Warehouse,
							complex:    true,
						}
						mo.insertManufacturingOrder(userId)
						c := ComplexManufacturingOrderManufacturingOrder{
							Type:                            "I",
							ComplexManufacturingOrder:       complexManufacturingOrder.Id,
							enterprise:                      complexManufacturingOrder.enterprise,
							ManufacturingOrder:              &mo.Id,
							Product:                         manufacturingOrderTypeComponent.Product,
							ManufacturingOrderTypeComponent: manufacturingOrderTypeComponent.Id,
							Manufactured:                    false,
						}
						subOrders = append(subOrders, c)
					}
				}
			} else { // if product.Manufacturing
				var purchaseDetailId int64 = 0
				// search for a pending purchase order detail
				sqlStatement := `SELECT id FROM purchase_order_detail WHERE product=$1 AND quantity_delivery_note = 0 AND quantity - quantity_assigned_sale >= $2 ORDER BY (SELECT date_created FROM purchase_order WHERE purchase_order.id=purchase_order_detail."order") ASC LIMIT 1`
				row := db.QueryRow(sqlStatement, manufacturingOrderTypeComponent.Product, manufacturingOrderTypeComponent.Quantity)
				if row.Err() != nil {
					log("DB", row.Err().Error())
					if beginTransaction {
						trans.Rollback()
					}
					return false
				}

				row.Scan(&purchaseDetailId)

				c := ComplexManufacturingOrderManufacturingOrder{
					Type:                            "I",
					ComplexManufacturingOrder:       complexManufacturingOrder.Id,
					enterprise:                      complexManufacturingOrder.enterprise,
					PurchaseOrderDetail:             &purchaseDetailId,
					Product:                         manufacturingOrderTypeComponent.Product,
					ManufacturingOrderTypeComponent: manufacturingOrderTypeComponent.Id,
					Manufactured:                    false,
				}
				subOrders = append(subOrders, c)

				// there are no pending purchase order details, return error
				if purchaseDetailId == 0 {
					if beginTransaction {
						trans.Rollback()
					}
					return false
				} else {
					// add quantity assigned to sale orders
					// TODO view the quantity assigled to the complex manufacturing orders from the frontend
					ok := addQuantityAssignedSalePurchaseOrder(purchaseDetailId, manufacturingOrderTypeComponent.Quantity, complexManufacturingOrder.enterprise, userId)
					if !ok {
						if beginTransaction {
							trans.Rollback()
						}
						return false
					}
				}
			}
		} // if stock.QuantityAvaialbe >= manufacturingOrderTypeComponent.Quantity {

	} // for i := 0; i < len(components); i++ {

	for i := 0; i < len(components); i++ {
		if components[i].Type != "O" { // Only Output
			continue
		}

		manufacturingOrderTypeComponent := components[i]
		if manufacturingOrderTypeComponent.Id <= 0 || manufacturingOrderTypeComponent.enterprise != c.enterprise || manufacturingOrderTypeComponent.Quantity <= 0 {
			if beginTransaction {
				trans.Rollback()
			}
			return false
		}

		c := ComplexManufacturingOrderManufacturingOrder{
			Type:                            "O",
			ComplexManufacturingOrder:       complexManufacturingOrder.Id,
			enterprise:                      complexManufacturingOrder.enterprise,
			Product:                         manufacturingOrderTypeComponent.Product,
			ManufacturingOrderTypeComponent: manufacturingOrderTypeComponent.Id,
			Manufactured:                    false,
		}
		subOrders = append(subOrders, c)
	} // for i := 0; i < len(components); i++ {

	for i := 0; i < len(subOrders); i++ {
		ok := subOrders[i].insertComplexManufacturingOrderManufacturingOrder(userId)
		if !ok {
			if beginTransaction {
				trans.Rollback()
			}
			return false
		}
	} // for i := 0; i < len(subOrders); i++ {

	return true
}

func getPendingComplexManufacturingOrderOutputsWithoutSaleOrderDetail(productId int32) *int64 {
	sqlStatement := `SELECT id FROM complex_manufacturing_order_manufacturing_order WHERE (product = $1) AND (NOT manufactured) AND (type = 'O') AND (sale_order_detail IS NULL) ORDER BY id ASC LIMIT 1`
	row := db.QueryRow(sqlStatement, productId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return nil
	}

	var id int64
	row.Scan(&id)
	return &id
}

func (c *ComplexManufacturingOrder) deleteComplexManufacturingOrder(userId int32) bool {

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	orderInMemory := getComplexManufacturingOrderRow(c.Id)
	if orderInMemory.Id <= 0 || orderInMemory.enterprise != c.enterprise {
		return false
	}

	components := getComplexManufacturingOrderManufacturingOrder(c.Id, c.enterprise)

	for i := 0; i < len(components); i++ {
		ok := components[i].deleteComplexManufacturingOrderManufacturingOrder(userId)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.complex_manufacturing_order WHERE id=$1`
	db.Exec(sqlStatement, c.Id)
	_, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		trans.Rollback()
		log("DB", err.Error())
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func toggleManufactuedComplexManufacturingOrder(orderid int64, userId int32, enterpriseId int32) bool {
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

	inMemoryComplexManufacturingOrder := getComplexManufacturingOrderRow(orderid)
	if inMemoryComplexManufacturingOrder.enterprise != enterpriseId {
		return false
	}

	// validation
	if inMemoryComplexManufacturingOrder.Manufactured && inMemoryComplexManufacturingOrder.DateManufactured != nil && int64(time.Since(*inMemoryComplexManufacturingOrder.DateManufactured).Seconds()) > int64(settings.UndoManufacturingOrderSeconds) {
		return false
	}
	if !inMemoryComplexManufacturingOrder.Manufactured && inMemoryComplexManufacturingOrder.QuantityManufactured != inMemoryComplexManufacturingOrder.QuantityPendingManufacture {
		return false
	}

	sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET manufactured = $2 WHERE complex_manufacturing_order=$1 AND type = 'O'`
	_, err := db.Exec(sqlStatement, orderid, !inMemoryComplexManufacturingOrder.Manufactured)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	cmomo := getComplexManufacturingOrderManufacturingOrder(orderid, enterpriseId)
	if !inMemoryComplexManufacturingOrder.Manufactured {
		for i := 0; i < len(cmomo); i++ {
			if cmomo[i].Type == "I" {
				continue
			}

			com := getManufacturingOrderTypeComponentRow(cmomo[i].ManufacturingOrderTypeComponent)

			wm := WarehouseMovement{
				Product:    cmomo[i].Product,
				Warehouse:  inMemoryComplexManufacturingOrder.Warehouse,
				Quantity:   com.Quantity,
				Type:       "O",
				enterprise: enterpriseId,
			}
			wm.insertWarehouseMovement(userId)

			sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET warehouse_movement=$2 WHERE id=$1`
			_, err := db.Exec(sqlStatement, cmomo[i].Id, wm.Id)
			if err != nil {
				log("DB", err.Error())
				return false
			}

			if cmomo[i].SaleOrderDetail != nil {
				sqlStatement := `SELECT COUNT(*) FROM public.complex_manufacturing_order_manufacturing_order WHERE sale_order_detail=$1 AND NOT manufactured`
				row := db.QueryRow(sqlStatement, cmomo[i].SaleOrderDetail)
				if row.Err() != nil {
					log("DB", row.Err().Error())
					return false
				}

				var ordersPending int32
				row.Scan(&ordersPending)

				if ordersPending == 0 {
					sqlStatement = `UPDATE sales_order_detail SET status = 'E' WHERE id = $1`
					_, err = db.Exec(sqlStatement, cmomo[i].SaleOrderDetail)
					if err != nil {
						log("DB", err.Error())
						return false
					}
				}
			}
		} // for i := 0; i < len(cmomo); i++ {

		sqlStatement := `UPDATE public.complex_manufacturing_order SET manufactured=true, date_manufactured=CURRENT_TIMESTAMP(3), user_manufactured=$2 WHERE id=$1`
		_, err := db.Exec(sqlStatement, inMemoryComplexManufacturingOrder.Id, userId)
		if err != nil {
			log("DB", err.Error())
		}
	} else { // if !inMemoryComplexManufacturingOrder.Manufactured {
		for i := 0; i < len(cmomo); i++ {
			if cmomo[i].Type == "I" {
				continue
			}

			sqlStatement := `UPDATE public.complex_manufacturing_order_manufacturing_order SET warehouse_movement=NULL WHERE id=$1`
			_, err := db.Exec(sqlStatement, cmomo[i].Id)
			if err != nil {
				log("DB", err.Error())
				return false
			}

			if cmomo[i].WarehouseMovement != nil {
				wm := getWarehouseMovementRow(*cmomo[i].WarehouseMovement)
				ok := wm.deleteWarehouseMovement(userId)
				if !ok {
					return false
				}
			}

			if cmomo[i].SaleOrderDetail != nil {
				sqlStatement = `UPDATE sales_order_detail SET status = 'D' WHERE id = $1`
				_, err = db.Exec(sqlStatement, cmomo[i].SaleOrderDetail)
				if err != nil {
					log("DB", err.Error())
					return false
				}
			}
		}

		sqlStatement := `UPDATE public.complex_manufacturing_order SET manufactured=false, date_manufactured=NULL, user_manufactured=NULL WHERE id=$1`
		_, err := db.Exec(sqlStatement, inMemoryComplexManufacturingOrder.Id)
		if err != nil {
			log("DB", err.Error())
		}
	} // } else { // if !inMemoryComplexManufacturingOrder.Manufactured {

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

type ComplexManufacturingOrderManufacturingOrder struct {
	Id                              int64   `json:"id"`
	ManufacturingOrder              *int64  `json:"manufacturingOrder"`
	Type                            string  `json:"type"` // I = Input, O = Output
	ComplexManufacturingOrder       int64   `json:"complexManufacturingOrder"`
	WarehouseMovement               *int64  `json:"warehouseMovement"`
	Manufactured                    bool    `json:"manufactured"`
	Product                         int32   `json:"product"`
	ManufacturingOrderTypeComponent int32   `json:"manufacturingOrderTypeComponent"`
	PurchaseOrderDetail             *int64  `json:"purchaseOrderDetail"`
	SaleOrderDetail                 *int64  `json:"saleOrderDetail"`
	ProductName                     *string `json:"productName"`
	PurchaseOrderName               *string `json:"purchaseOrderName"`
	SaleOrderName                   *string `json:"saleOrderName"`
	enterprise                      int32
}

func getComplexManufacturingOrderManufacturingOrder(complexManufacturingOrderId int64, enterpriseId int32) []ComplexManufacturingOrderManufacturingOrder {
	var orders []ComplexManufacturingOrderManufacturingOrder = make([]ComplexManufacturingOrderManufacturingOrder, 0)

	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=complex_manufacturing_order_manufacturing_order.product),(SELECT order_name FROM purchase_order WHERE purchase_order.id=(SELECT "order" FROM purchase_order_detail WHERE purchase_order_detail.id=complex_manufacturing_order_manufacturing_order.purchase_order_detail)),(SELECT order_name FROM sales_order WHERE sales_order.id=(SELECT "order" FROM sales_order_detail WHERE sales_order_detail.id=complex_manufacturing_order_manufacturing_order.sale_order_detail)) FROM public.complex_manufacturing_order_manufacturing_order WHERE complex_manufacturing_order=$1 AND enterprise=$2`
	rows, err := db.Query(sqlStatement, complexManufacturingOrderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return orders
	}

	for rows.Next() {
		c := ComplexManufacturingOrderManufacturingOrder{}
		rows.Scan(&c.Id, &c.ManufacturingOrder, &c.Type, &c.ComplexManufacturingOrder, &c.enterprise, &c.WarehouseMovement, &c.Manufactured, &c.Product, &c.ManufacturingOrderTypeComponent, &c.PurchaseOrderDetail, &c.SaleOrderDetail, &c.ProductName, &c.PurchaseOrderName, &c.SaleOrderName)
		orders = append(orders, c)
	}

	return orders
}

func getComplexManufacturingOrderManufacturingOrderRow(complexManufacturingOrderManufacturingOrderId int64) ComplexManufacturingOrderManufacturingOrder {
	c := ComplexManufacturingOrderManufacturingOrder{}

	sqlStatement := `SELECT * FROM public.complex_manufacturing_order_manufacturing_order WHERE id=$1`
	row := db.QueryRow(sqlStatement, complexManufacturingOrderManufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return c
	}

	row.Scan(&c.Id, &c.ManufacturingOrder, &c.Type, &c.ComplexManufacturingOrder, &c.enterprise, &c.WarehouseMovement, &c.Manufactured, &c.Product, &c.ManufacturingOrderTypeComponent, &c.PurchaseOrderDetail, &c.SaleOrderDetail)
	return c
}

func (c *ComplexManufacturingOrderManufacturingOrder) isValid() bool {
	return !(c.Product <= 0 || (c.Type != "I" && c.Type != "O") || c.ComplexManufacturingOrder <= 0 || c.ManufacturingOrderTypeComponent <= 0)
}

// DOES NOT OPEN A TRANSACTION
func (c *ComplexManufacturingOrderManufacturingOrder) insertComplexManufacturingOrderManufacturingOrder(userId int32) bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.complex_manufacturing_order_manufacturing_order(manufacturing_order, type, complex_manufacturing_order, enterprise, warehouse_movement, product, manufacturing_order_type_component, purchase_order_detail, sale_order_detail, manufactured) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := db.Exec(sqlStatement, c.ManufacturingOrder, c.Type, c.ComplexManufacturingOrder, c.enterprise, c.WarehouseMovement, c.Product, c.ManufacturingOrderTypeComponent, c.PurchaseOrderDetail, c.SaleOrderDetail, c.Manufactured)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	ok := addQuantityPendingManufactureComplexManufacturingOrder(c.ComplexManufacturingOrder, 1)
	if ok && c.WarehouseMovement != nil {
		return addQuantityManufacturedComplexManufacturingOrder(c.ComplexManufacturingOrder, 1)
	}
	return ok
}

// DOES NOT OPEN A TRANSACTION
func (c *ComplexManufacturingOrderManufacturingOrder) deleteComplexManufacturingOrderManufacturingOrder(userId int32) bool {
	if c.Id <= 0 {
		return false
	}

	comInMemory := getComplexManufacturingOrderManufacturingOrderRow(c.Id)
	if comInMemory.Id <= 0 || comInMemory.enterprise != c.enterprise {
		return false
	}

	sqlStatement := `DELETE FROM public.complex_manufacturing_order_manufacturing_order WHERE id=$1`
	_, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	ok := addQuantityPendingManufactureComplexManufacturingOrder(c.ComplexManufacturingOrder, -1)
	if !ok {
		return false
	}

	if comInMemory.ManufacturingOrder != nil {
		mo := getManufacturingOrderRow(*comInMemory.ManufacturingOrder)
		ok := mo.deleteManufacturingOrder(userId)
		if !ok {
			return false
		}
	}

	if comInMemory.WarehouseMovement != nil {
		wm := getWarehouseMovementRow(*comInMemory.WarehouseMovement)
		ok := wm.deleteWarehouseMovement(userId)
		if !ok {
			return false
		}
	}

	if comInMemory.PurchaseOrderDetail != nil {
		component := getManufacturingOrderTypeComponentRow(comInMemory.ManufacturingOrderTypeComponent)
		ok := addQuantityAssignedSalePurchaseOrder(*comInMemory.PurchaseOrderDetail, component.Quantity, comInMemory.enterprise, userId)
		if !ok {
			return false
		}
	}

	if comInMemory.SaleOrderDetail != nil {
		sqlStatement = `UPDATE sales_order_detail SET status = 'C' WHERE id = $1`
		_, err := db.Exec(sqlStatement, comInMemory.SaleOrderDetail)
		if err != nil {
			log("DB", err.Error())
			return false
		}

		insertTransactionalLog(c.enterprise, "sales_order_detail", int(*comInMemory.SaleOrderDetail), userId, "U")

		ok := setSalesOrderState(c.enterprise, *comInMemory.SaleOrderDetail, userId)
		if !ok {
			return false
		}
	}

	return true
}

func setComplexManufacturingOrderManufacturingOrderManufactured(manufacturingOrderId int64, manufactured bool) bool {
	sqlStatement := `SELECT complex_manufacturing_order, manufactured FROM public.complex_manufacturing_order_manufacturing_order WHERE manufacturing_order=$1 AND type = 'O'`
	row := db.QueryRow(sqlStatement, manufacturingOrderId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var complexManufacturingOrderId int64
	var orderManufactured bool
	row.Scan(&complexManufacturingOrderId, &orderManufactured)

	if complexManufacturingOrderId <= 0 || manufactured == orderManufactured {
		return false
	}

	sqlStatement = `UPDATE public.complex_manufacturing_order_manufacturing_order SET manufactured = $2 WHERE manufacturing_order=$1 AND type = 'O'`
	_, err := db.Exec(sqlStatement, manufacturingOrderId, manufactured)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	if !orderManufactured == manufactured {
		return addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId, 1)
	} else if orderManufactured == !manufactured {
		return addQuantityManufacturedComplexManufacturingOrder(complexManufacturingOrderId, -1)
	}

	return false
}
