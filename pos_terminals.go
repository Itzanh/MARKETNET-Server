package main

import (
	"github.com/google/uuid"
)

type POSTerminal struct {
	Id                        int64   `json:"id"`
	Uuid                      string  `json:"uuid"`
	Name                      string  `json:"name"`
	OrdersCustomer            *int32  `json:"ordersCustomer"`
	OrdersInvoiceAddress      *int32  `json:"ordersInvoiceAddress"`
	OrdersDeliveryAddress     *int32  `json:"ordersDeliveryAddress"`
	OrdersPaymentMethod       *int32  `json:"ordersPaymentMethod"`
	OrdersBillingSeries       *string `json:"ordersBillingSeries"`
	OrdersWarehouse           *string `json:"ordersWarehouse"`
	OrdersCurrency            *int32  `json:"ordersCurrency"`
	enterprise                int32
	OrdersCustomerName        *string `json:"ordersCustomerName"`
	OrdersInvoiceAddressName  *string `json:"ordersInvoiceAddressName"`
	OrdersDeliveryAddressName *string `json:"ordersDeliveryAddressName"`
}

func getPOSTerminals(enterpriseId int32) []POSTerminal {
	var posTerminals []POSTerminal = make([]POSTerminal, 0)
	sqlStatement := `SELECT *,(SELECT name FROM customer WHERE customer.id=pos_terminals.orders_customer),(SELECT address FROM address WHERE address.id=pos_terminals.orders_invoice_address),(SELECT address FROM address WHERE address.id=pos_terminals.orders_delivery_address) FROM public.pos_terminals WHERE enterprise = $1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return posTerminals
	}
	defer rows.Close()

	for rows.Next() {
		var t POSTerminal
		rows.Scan(&t.Id, &t.Uuid, &t.Name, &t.OrdersCustomer, &t.OrdersInvoiceAddress, &t.OrdersDeliveryAddress, &t.OrdersPaymentMethod, &t.OrdersBillingSeries, &t.OrdersWarehouse, &t.OrdersCurrency, &t.enterprise, &t.OrdersCustomerName, &t.OrdersInvoiceAddressName, &t.OrdersDeliveryAddressName)
		posTerminals = append(posTerminals, t)
	}
	return posTerminals
}

func getPOSTerminalByUUID(uuid string, enterpriseId int32) POSTerminal {
	sqlStatement := `SELECT * FROM public.pos_terminals WHERE uuid = $1 AND enterprise = $2`
	row := db.QueryRow(sqlStatement, uuid, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return POSTerminal{}
	}

	var t POSTerminal
	row.Scan(&t.Id, &t.Uuid, &t.Name, &t.OrdersCustomer, &t.OrdersInvoiceAddress, &t.OrdersDeliveryAddress, &t.OrdersPaymentMethod, &t.OrdersBillingSeries, &t.OrdersWarehouse, &t.OrdersCurrency, &t.enterprise)

	return t
}

func (t *POSTerminal) isValid() bool {
	return !(len(t.Name) == 0 || len(t.Name) > 150 || t.enterprise <= 0)
}

func (t *POSTerminal) isReady() bool {
	return t.Id > 0 && t.isValid() && !(t.OrdersCustomer == nil || t.OrdersInvoiceAddress == nil || t.OrdersDeliveryAddress == nil || t.OrdersPaymentMethod == nil || t.OrdersBillingSeries == nil || t.OrdersWarehouse == nil || t.OrdersCurrency == nil)
}

type TerminalRegisterResult struct {
	Ok      bool   `json:"ok"`
	Uuid    string `json:"uuid"`
	IsReady bool   `json:"isReady"`
}

func posTerminalRequest(terminal string, enterpriseId int32) TerminalRegisterResult {
	if len(terminal) == 0 {
		t := POSTerminal{
			enterprise: enterpriseId,
		}
		ok := t.insertPOSTerminal()
		if !ok {
			return TerminalRegisterResult{Ok: false}
		}
		return TerminalRegisterResult{Ok: true, Uuid: t.Uuid, IsReady: false}
	} else {
		t := getPOSTerminalByUUID(terminal, enterpriseId)
		if t.Id <= 0 {
			t := POSTerminal{
				enterprise: enterpriseId,
			}
			ok := t.insertPOSTerminal()
			if !ok {
				return TerminalRegisterResult{Ok: false}
			}
			return TerminalRegisterResult{Ok: true, Uuid: t.Uuid, IsReady: false}
		} else {
			return TerminalRegisterResult{Ok: true, Uuid: t.Uuid, IsReady: t.isReady()}
		}
	}
}

func (t *POSTerminal) insertPOSTerminal() bool {
	t.Uuid = uuid.New().String()
	t.Name = t.Uuid

	if !t.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.pos_terminals(uuid, name, enterprise) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, t.Uuid, t.Uuid, t.enterprise)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (t *POSTerminal) updatePOSTerminal() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.pos_terminals SET name=$3, orders_customer=$4, orders_invoice_address=$5, orders_delivery_address=$6, orders_payment_method=$7, orders_billing_series=$8, orders_warehouse=$9, orders_currency=$10 WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, t.Id, t.enterprise, t.Name, t.OrdersCustomer, t.OrdersInvoiceAddress, t.OrdersDeliveryAddress, t.OrdersPaymentMethod, t.OrdersBillingSeries, t.OrdersWarehouse, t.OrdersCurrency)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func posInsertNewSaleOrder(terminal string, enterpriseId int32, userId int32) SaleOrder {
	posTerminal := getPOSTerminalByUUID(terminal, enterpriseId)
	if !posTerminal.isReady() {
		return SaleOrder{}
	}

	o := SaleOrder{
		Warehouse:       *posTerminal.OrdersWarehouse,
		Customer:        *posTerminal.OrdersCustomer,
		PaymentMethod:   *posTerminal.OrdersPaymentMethod,
		BillingSeries:   *posTerminal.OrdersBillingSeries,
		Currency:        *posTerminal.OrdersCurrency,
		BillingAddress:  *posTerminal.OrdersInvoiceAddress,
		ShippingAddress: *posTerminal.OrdersDeliveryAddress,
		enterprise:      enterpriseId,
	}
	_, orderId := o.insertSalesOrder(userId)
	o.Id = orderId
	return o
}

type InsertNewSaleOrderDetail struct {
	Terminal string `json:"terminal"`
	Order    int64  `json:"order"`
	BarCode  string `json:"barCode"`
	Quantity int32  `json:"quantity"`
}

func (i *InsertNewSaleOrderDetail) posInsertNewSaleOrderDetail(enterpriseId int32, userId int32) bool {
	posTerminal := getPOSTerminalByUUID(i.Terminal, enterpriseId)
	if !posTerminal.isReady() {
		return false
	}

	product := getProductByBarcode(i.BarCode, enterpriseId)
	if product.Id <= 0 {
		return false
	}

	// increment by one the quantity if we scan the same code again
	details := getSalesOrderDetail(i.Order, enterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].Product == product.Id {
			details[i].Quantity++
			return details[i].updateSalesOrderDetail(userId).Ok
		}
	}

	d := SalesOrderDetail{
		Order:      i.Order,
		Product:    product.Id,
		Quantity:   i.Quantity,
		Price:      product.Price,
		VatPercent: product.VatPercent,
		enterprise: enterpriseId,
	}
	return d.insertSalesOrderDetail(userId).Ok
}

type POSServeSaleOrder struct {
	Ok        bool  `json:"ok"`
	InvoiceId int64 `json:"invoiceId"`
}

func posServeSaleOrder(orderId int64, enterpriseId int32, userId int32) POSServeSaleOrder {
	ok := invoiceAllSaleOrder(orderId, enterpriseId, userId).Ok
	if !ok {
		return POSServeSaleOrder{Ok: false}
	}

	okAndErr, _ := deliveryNoteAllSaleOrder(orderId, enterpriseId, userId, nil)
	if !okAndErr.Ok {
		return POSServeSaleOrder{Ok: false}
	}

	inv := getSalesOrderInvoices(orderId, enterpriseId)
	if len(inv) == 0 {
		return POSServeSaleOrder{Ok: false}
	}
	return POSServeSaleOrder{Ok: true, InvoiceId: inv[0].Id}
}
