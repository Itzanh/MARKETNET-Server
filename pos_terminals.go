package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type POSTerminal struct {
	Id                      int64          `json:"id"`
	Uuid                    string         `json:"uuid" gorm:"type:uuid;not null;true;index:pos_terminals_uuid,unique:true"`
	Name                    string         `json:"name" gorm:"type:varchar(150);not null;true"`
	OrdersCustomerId        *int32         `json:"ordersCustomerId" gorm:"column:orders_customer"`
	OrdersCustomer          *Customer      `json:"ordersCustomer" gorm:"foreignKey:OrdersCustomerId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersInvoiceAddressId  *int32         `json:"ordersInvoiceAddressId" gorm:"column:orders_invoice_address"`
	OrdersInvoiceAddress    *Address       `json:"ordersInvoiceAddress" gorm:"foreignKey:OrdersInvoiceAddressId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersDeliveryAddressId *int32         `json:"ordersDeliveryAddressId" gorm:"column:orders_delivery_address"`
	OrdersDeliveryAddress   *Address       `json:"ordersDeliveryAddress" gorm:"foreignKey:OrdersDeliveryAddressId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersPaymentMethodId   *int32         `json:"ordersPaymentMethodId" gorm:"column:orders_payment_method"`
	OrdersPaymentMethod     *PaymentMethod `json:"ordersPaymentMethod" gorm:"foreignKey:OrdersPaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersBillingSeriesId   *string        `json:"ordersBillingSeriesId" gorm:"column:orders_billing_series"`
	OrdersBillingSeries     *BillingSerie  `json:"ordersBillingSeries" gorm:"foreignKey:OrdersBillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersWarehouseId       *string        `json:"ordersWarehouseId" gorm:"column:orders_warehouse"`
	OrdersWarehouse         *Warehouse     `json:"ordersWarehouse" gorm:"foreignKey:OrdersWarehouseId,EnterpriseId;references:Id,EnterpriseId"`
	OrdersCurrencyId        *int32         `json:"ordersCurrencyId" gorm:"column:orders_currency"`
	OrdersCurrency          *Currency      `json:"ordersCurrency" gorm:"foreignKey:OrdersCurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId            int32          `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise              Settings       `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (t *POSTerminal) TableName() string {
	return "pos_terminals"
}

func getPOSTerminals(enterpriseId int32) []POSTerminal {
	var posTerminals []POSTerminal = make([]POSTerminal, 0)
	// get all the pos terminals for the enterprise id sorted by id asc using dbOrm
	dbOrm.Where("enterprise = ?", enterpriseId).Preload(clause.Associations).Order("id asc").Find(&posTerminals)
	return posTerminals
}

func getPOSTerminalByUUID(uuid string, enterpriseId int32) POSTerminal {
	var t POSTerminal
	result := dbOrm.Where("uuid = ? AND enterprise = ?", uuid, enterpriseId).First(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return t
}

func (t *POSTerminal) isValid() bool {
	return !(len(t.Name) == 0 || len(t.Name) > 150 || t.EnterpriseId <= 0)
}

func (t *POSTerminal) isReady() bool {
	return t.Id > 0 && t.isValid() && !(t.OrdersCustomerId == nil || t.OrdersInvoiceAddressId == nil || t.OrdersDeliveryAddressId == nil || t.OrdersPaymentMethodId == nil || t.OrdersBillingSeriesId == nil || t.OrdersWarehouseId == nil || t.OrdersCurrencyId == nil)
}

type TerminalRegisterResult struct {
	Ok      bool   `json:"ok"`
	Uuid    string `json:"uuid"`
	IsReady bool   `json:"isReady"`
}

func posTerminalRequest(terminal string, enterpriseId int32) TerminalRegisterResult {
	if len(terminal) == 0 {
		t := POSTerminal{
			EnterpriseId: enterpriseId,
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
				EnterpriseId: enterpriseId,
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

func (t *POSTerminal) BeforeCreate(tx *gorm.DB) (err error) {
	var posTerminal POSTerminal
	tx.Model(&POSTerminal{}).Last(&posTerminal)
	t.Id = posTerminal.Id + 1
	return nil
}

func (t *POSTerminal) insertPOSTerminal() bool {
	t.Uuid = uuid.New().String()
	t.Name = t.Uuid

	if !t.isValid() {
		return false
	}

	var terminal POSTerminal

	terminal.Uuid = t.Uuid
	terminal.Name = t.Name
	terminal.EnterpriseId = t.EnterpriseId

	result := dbOrm.Create(&terminal)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (t *POSTerminal) updatePOSTerminal() bool {
	if !t.isValid() {
		return false
	}

	// get a single pos terminal by id and enterprise id
	var terminal POSTerminal
	result := dbOrm.Where("id = ? AND enterprise = ?", t.Id, t.EnterpriseId).First(&terminal)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// update the pos terminal
	terminal.Name = t.Name
	terminal.OrdersCustomerId = t.OrdersCustomerId
	terminal.OrdersInvoiceAddressId = t.OrdersInvoiceAddressId
	terminal.OrdersDeliveryAddressId = t.OrdersDeliveryAddressId
	terminal.OrdersPaymentMethodId = t.OrdersPaymentMethodId
	terminal.OrdersBillingSeriesId = t.OrdersBillingSeriesId
	terminal.OrdersWarehouseId = t.OrdersWarehouseId
	terminal.OrdersCurrencyId = t.OrdersCurrencyId

	result = dbOrm.Save(&terminal)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func posInsertNewSaleOrder(terminal string, enterpriseId int32, userId int32) SaleOrder {
	posTerminal := getPOSTerminalByUUID(terminal, enterpriseId)
	if !posTerminal.isReady() {
		return SaleOrder{}
	}

	o := SaleOrder{
		WarehouseId:       *posTerminal.OrdersWarehouseId,
		CustomerId:        *posTerminal.OrdersCustomerId,
		PaymentMethodId:   *posTerminal.OrdersPaymentMethodId,
		BillingSeriesId:   *posTerminal.OrdersBillingSeriesId,
		CurrencyId:        *posTerminal.OrdersCurrencyId,
		BillingAddressId:  *posTerminal.OrdersInvoiceAddressId,
		ShippingAddressId: *posTerminal.OrdersDeliveryAddressId,
		EnterpriseId:      enterpriseId,
	}
	_, orderId := o.insertSalesOrder(userId)
	o.Id = orderId
	return o
}

func deletePOSTerminal(terminal string, enterpriseId int32) bool {
	// delete a single pos terminal by id and enterprise id
	result := dbOrm.Where("uuid = ? AND enterprise = ?", terminal, enterpriseId).Delete(&POSTerminal{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
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
		if details[i].ProductId == product.Id {
			details[i].Quantity++
			return details[i].updateSalesOrderDetail(userId).Ok
		}
	}

	d := SalesOrderDetail{
		OrderId:      i.Order,
		ProductId:    product.Id,
		Quantity:     i.Quantity,
		Price:        product.Price,
		VatPercent:   product.VatPercent,
		EnterpriseId: enterpriseId,
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
