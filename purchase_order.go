package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseOrder struct {
	Id                int64         `json:"id" gorm:"index:purchase_order_id_enterprise,unique:true,priority:1"`
	SupplierReference string        `json:"supplierReference " gorm:"column:supplier_reference;type:character varying(40);not null:true;index:purchase_order_supplier_reference,type:gin"`
	SupplierId        int32         `json:"supplierId" gorm:"column:supplier;type:integer;not null:true"`
	Supplier          Supplier      `json:"supplier" gorm:"foreignKey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated       time.Time     `json:"dateCreated" gorm:"column:date_created;type:timestamp(3) with time zone;not null:true"`
	DatePaid          *time.Time    `json:"datePaid" gorm:"column:date_paid;type:timestamp(3) with time zone"`
	PaymentMethodId   int32         `json:"paymentMethodId" gorm:"column:payment_method;type:integer;not null:true"`
	PaymentMethod     PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId   string        `json:"billingSeriesId" gorm:"column:billing_series;type:character(3);not null:true;index:purchase_order_order_number,unique:true,priority:2"`
	BillingSeries     BillingSerie  `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyId        int32         `json:"currencyId" gorm:"column:currency;type:integer;not null:true"`
	Currency          Currency      `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange    float64       `json:"currencyChange" gorm:"column:currency_change;type:numeric(14,6);not null:true"`
	BillingAddressId  int32         `json:"billingAddressId" gorm:"column:billing_address;type:integer;not null:true"`
	BillingAddress    Address       `json:"billingAddress" gorm:"foreignKey:BillingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	ShippingAddressId int32         `json:"shippingAddressId" gorm:"column:shipping_address;type:integer;not null:true"`
	ShippingAddress   Address       `json:"shippingAddress" gorm:"foreignKey:ShippingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	LinesNumber       int16         `json:"linesNumber" gorm:"column:lines_number;type:smallint;not null:true"`
	InvoicedLines     int16         `json:"invoicedLines" gorm:"column:invoiced_lines;type:smallint;not null:true"`
	DeliveryNoteLines int16         `json:"deliveryNoteLines" gorm:"column:delivery_note_lines;type:smallint;not null:true"`
	TotalProducts     float64       `json:"totalProducts" gorm:"column:total_products;type:numeric(14,6);not null:true"`
	DiscountPercent   float64       `json:"discountPercent" gorm:"column:discount_percent;type:numeric(14,6);not null:true"`
	FixDiscount       float64       `json:"fixDiscount" gorm:"column:fix_discount;type:numeric(14,6);not null:true"`
	ShippingPrice     float64       `json:"shippingPrice" gorm:"column:shipping_price;type:numeric(14,6);not null:true"`
	ShippingDiscount  float64       `json:"shippingDiscount" gorm:"column:shipping_discount;type:numeric(14,6);not null:true"`
	TotalWithDiscount float64       `json:"totalWithDiscount" gorm:"column:total_with_discount;type:numeric(14,6);not null:true"`
	VatAmount         float64       `json:"vatAmount" gorm:"column:total_vat;type:numeric(14,6);not null:true"`
	TotalAmount       float64       `json:"totalAmount" gorm:"column:total_amount;type:numeric(14,6);not null:true"`
	Description       string        `json:"description" gorm:"column:dsc;type:text;not null:true"`
	Notes             string        `json:"notes" gorm:"column:notes;type:character varying(250);not null:true"`
	Cancelled         bool          `json:"cancelled" gorm:"column:cancelled;type:boolean;not null:true"`
	OrderNumber       int32         `json:"orderNumber" gorm:"column:order_number;type:integer;not null:true;index:purchase_order_order_number,unique:true,priority:3"`
	OrderName         string        `json:"orderName" gorm:"column:order_name;type:character(15);not null:true"`
	EnterpriseId      int32         `json:"-" gorm:"column:enterprise;not null:true;index:purchase_order_id_enterprise,unique:true,priority:2;index:purchase_order_order_number,unique:true,priority:1"`
	Enterprise        Settings      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (po *PurchaseOrder) TableName() string {
	return "purchase_order"
}

type PurchaseOrders struct {
	Rows   int64               `json:"rows"`
	Orders []PurchaseOrder     `json:"orders"`
	Footer PurchaseOrderFooter `json:"footer"`
}

type PurchaseOrderFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseOrder(enterpriseId int32) PurchaseOrders {
	o := PurchaseOrders{}
	o.Orders = make([]PurchaseOrder, 0)
	result := dbOrm.Model(&PurchaseOrder{}).Where("purchase_order.enterprise = ?", enterpriseId).Order("purchase_order.date_created DESC").Preload(clause.Associations).Find(&o.Orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}
	o.Footer = PurchaseOrderFooter{}
	result = dbOrm.Model(&PurchaseOrder{}).Where("purchase_order.enterprise = ?", enterpriseId).Count(&o.Rows).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&o.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}
	return o
}

func (s *OrderSearch) searchPurchaseOrder() PurchaseOrders {
	o := PurchaseOrders{}
	o.Orders = make([]PurchaseOrder, 0)

	cursor := dbOrm.Model(&PurchaseOrder{}).Where("purchase_order.enterprise = ?", s.enterprise)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("purchase_order.order_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_order.supplier").Where("suppliers.name ILIKE @search OR purchase_order.supplier_reference ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_order.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_order.date_created <= ?", s.DateEnd)
		}
	}
	result := cursor.Preload(clause.Associations).Order("purchase_order.date_created DESC").Count(&o.Rows).Find(&o.Orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}

	cursor = dbOrm.Model(&PurchaseOrder{}).Where("purchase_order.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("purchase_order.order_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_order.supplier").Where("suppliers.name ILIKE @search OR purchase_order.supplier_reference ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_order.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_order.date_created <= ?", s.DateEnd)
		}
	}
	result = cursor.Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&o.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return o
	}

	return o
}

func getPurchaseOrderRow(orderId int64) PurchaseOrder {
	p := PurchaseOrder{}
	result := dbOrm.Model(&PurchaseOrder{}).Where("purchase_order.id = ?", orderId).Preload(clause.Associations).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return p
	}
	return p
}

func (p *PurchaseOrder) isValid() bool {
	return !(len(p.SupplierReference) > 40 || p.SupplierId <= 0 || p.PaymentMethodId <= 0 || len(p.BillingSeriesId) == 0 || p.CurrencyId <= 0 || p.BillingAddressId <= 0 || p.ShippingAddressId <= 0 || len(p.Notes) > 250 || len(p.Description) > 3000)
}

func (o *PurchaseOrder) BeforeCreate(tx *gorm.DB) (err error) {
	var purchaseOrder PurchaseOrder
	tx.Model(&PurchaseOrder{}).Last(&purchaseOrder)
	o.Id = purchaseOrder.Id + 1
	return nil
}

func (p *PurchaseOrder) insertPurchaseOrder(userId int32, trans *gorm.DB) (bool, int64) {
	if !p.isValid() {
		return false, 0
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false, 0
		}
		///
	}

	p.OrderNumber = getNextPurchaseOrderNumber(p.BillingSeriesId, p.EnterpriseId)
	if p.OrderNumber <= 0 {
		return false, 0
	}
	p.CurrencyChange = getCurrencyExchange(p.CurrencyId)
	now := time.Now()
	p.OrderName = p.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", p.OrderNumber)

	p.DateCreated = time.Now()
	p.DatePaid = nil
	p.LinesNumber = 0
	p.InvoicedLines = 0
	p.DeliveryNoteLines = 0
	p.TotalProducts = 0
	p.TotalWithDiscount = p.ShippingPrice - p.ShippingDiscount - p.FixDiscount
	p.VatAmount = 0
	p.Cancelled = false

	result := trans.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(p.EnterpriseId, "purchase_order", int(p.Id), userId, "I")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "purchase_order", "POST", string(json))

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false, 0
		}
		///
	}

	return true, p.Id
}

func (p *PurchaseOrder) updatePurchaseOrder(userId int32) bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	inMemoryOrder := getPurchaseOrderRow(p.Id)
	if inMemoryOrder.Id <= 0 || inMemoryOrder.EnterpriseId != p.EnterpriseId {
		trans.Rollback()
		return false
	}

	if inMemoryOrder.InvoicedLines == 0 { // if the payment is pending, we allow to change more fields
		if p.CurrencyId != inMemoryOrder.CurrencyId {
			p.CurrencyChange = getCurrencyExchange(p.CurrencyId)
		}

		inMemoryOrder.SupplierId = p.SupplierId
		inMemoryOrder.PaymentMethodId = p.PaymentMethodId
		inMemoryOrder.CurrencyId = p.CurrencyId
		inMemoryOrder.CurrencyChange = p.CurrencyChange
		inMemoryOrder.BillingAddressId = p.BillingAddressId
		inMemoryOrder.ShippingAddressId = p.ShippingAddressId
		inMemoryOrder.DiscountPercent = p.DiscountPercent
		inMemoryOrder.FixDiscount = p.FixDiscount
		inMemoryOrder.ShippingPrice = p.ShippingPrice
		inMemoryOrder.ShippingDiscount = p.ShippingDiscount
		inMemoryOrder.Description = p.Description
		inMemoryOrder.Notes = p.Notes
		inMemoryOrder.SupplierReference = p.SupplierReference

		result := trans.Save(&inMemoryOrder)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		if p.DiscountPercent != inMemoryOrder.DiscountPercent || p.FixDiscount != inMemoryOrder.FixDiscount || p.ShippingPrice != inMemoryOrder.ShippingPrice || p.ShippingDiscount != inMemoryOrder.ShippingDiscount {
			ok := calcTotalsPurchaseOrder(p.Id, p.EnterpriseId, userId, *trans)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	} else {
		inMemoryOrder.SupplierId = p.SupplierId
		inMemoryOrder.BillingAddressId = p.BillingAddressId
		inMemoryOrder.ShippingAddressId = p.ShippingAddressId
		inMemoryOrder.Description = p.Description
		inMemoryOrder.Notes = p.Notes
		inMemoryOrder.SupplierReference = p.SupplierReference

		result := trans.Save(&inMemoryOrder)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	result := trans.Commit()
	if result.Error != nil {
		return false
	}
	///

	insertTransactionalLog(p.EnterpriseId, "purchase_order", int(p.Id), userId, "U")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// ERROR CODES:
// 1. Alerady invoiced
// 2. Delivery note generated
// 3. Error deleting detail <product>: <error>
func (p *PurchaseOrder) deletePurchaseOrder(userId int32) OkAndErrorCodeReturn {
	if p.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	inMemoryOrder := getPurchaseOrderRow(p.Id)
	if inMemoryOrder.EnterpriseId != p.EnterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}

	d := getPurchaseOrderDetail(p.Id, p.EnterpriseId)

	// prevent the order to be deleted if there is an invoice or a delivery note
	for i := 0; i < len(d); i++ {
		if d[i].QuantityInvoiced > 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
		if d[i].QuantityDeliveryNote > 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
		}
	}

	// delete details
	for i := 0; i < len(d); i++ {
		d[i].EnterpriseId = p.EnterpriseId
		ok := d[i].deletePurchaseOrderDetail(userId, trans)
		if !ok.Ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{strconv.Itoa(int(ok.ErrorCode)), d[i].Product.Name}}
		}
	}

	insertTransactionalLog(inMemoryOrder.EnterpriseId, "purchase_order", int(p.Id), userId, "D")
	json, _ := json.Marshal(p)
	go fireWebHook(p.EnterpriseId, "purchase_order", "DELETE", string(json))

	result := trans.Delete(&PurchaseOrder{}, "id = ?", p.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseOrder(orderId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.TotalProducts += totalAmount
	purchaseOrder.VatAmount += (totalAmount / 100) * vatPercent

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseOrder(orderId, enterpriseId, userId, trans)
}

// If the payment accepted date is null, sets it to the current date and time.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setDatePaymentAcceptedPurchaseOrder(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if purchaseOrder.DatePaid == nil {
		now := time.Now()
		purchaseOrder.DatePaid = &now

		result = trans.Save(&purchaseOrder)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
		json, _ := json.Marshal(purchaseOrder)
		go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))
	}

	return true
}

// Applies the logic to calculate the totals of the purchase order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseOrder(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.TotalWithDiscount = (purchaseOrder.TotalProducts - purchaseOrder.TotalProducts*(purchaseOrder.DiscountPercent/100)) - purchaseOrder.FixDiscount + purchaseOrder.ShippingPrice - purchaseOrder.ShippingDiscount
	purchaseOrder.TotalAmount = purchaseOrder.TotalWithDiscount + purchaseOrder.VatAmount

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

type PurchaseOrderDefaults struct {
	Warehouse     string `json:"warehouse"`
	WarehouseName string `json:"warehouseName"`
}

func getPurchaseOrderDefaults(enterpriseId int32) PurchaseOrderDefaults {
	s := getSettingsRecordById(enterpriseId)
	warehouseName := getNameWarehouse(s.DefaultWarehouseId, s.Id)

	return PurchaseOrderDefaults{Warehouse: s.DefaultWarehouseId, WarehouseName: warehouseName}
}

type PurchaseOrderRelations struct {
	Invoices      []PurchaseInvoice      `json:"invoices"`
	DeliveryNotes []PurchaseDeliveryNote `json:"deliveryNotes"`
}

func getPurchaseOrderRelations(orderId int64, enterpriseId int32) PurchaseOrderRelations {
	return PurchaseOrderRelations{
		Invoices:      getPurchaseOrderInvoices(orderId, enterpriseId),
		DeliveryNotes: getPurchaseOrderDeliveryNotes(orderId, enterpriseId),
	}
}

func getPurchaseOrderInvoices(orderId int64, enterpriseId int32) []PurchaseInvoice {
	purchaseOrder := getPurchaseOrderRow(orderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.EnterpriseId != enterpriseId {
		return []PurchaseInvoice{}
	}
	var invoices []PurchaseInvoice = make([]PurchaseInvoice, 0)
	purchaseOrderDetails := getPurchaseOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(purchaseOrderDetails); i++ {
		var invoiceDetails []PurchaseInvoiceDetail
		dbOrm.Model(&PurchaseInvoiceDetail{}).Where("order_detail = ?", purchaseOrderDetails[i].Id).Find(&invoiceDetails)
		for j := 0; j < len(invoiceDetails); j++ {
			// only append invoice to invoices if it doesn't already exist in the array searching by id
			var ok bool = true
			for k := 0; k < len(invoices); k++ {
				if invoices[k].Id == invoiceDetails[j].InvoiceId {
					ok = false
					break
				}
			}
			if ok {
				invoice := getPurchaseInvoiceRow(invoiceDetails[j].InvoiceId)
				invoices = append(invoices, invoice)
			}
		}
	}
	return invoices
}

func getPurchaseOrderDeliveryNotes(orderId int64, enterpriseId int32) []PurchaseDeliveryNote {
	purchaseOrder := getPurchaseOrderRow(orderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.EnterpriseId != enterpriseId {
		return []PurchaseDeliveryNote{}
	}
	var notes []PurchaseDeliveryNote = make([]PurchaseDeliveryNote, 0)
	purchaseOrderDetails := getPurchaseOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(purchaseOrderDetails); i++ {
		var noteDetails []WarehouseMovement
		dbOrm.Model(&WarehouseMovement{}).Where("purchase_order_detail = ?", purchaseOrderDetails[i].Id).Find(&noteDetails)
		for j := 0; j < len(noteDetails); j++ {
			// only append note to notes if it doesn't already exist in the array searching by id
			var ok bool = true
			for k := 0; k < len(notes); k++ {
				if notes[k].Id == *noteDetails[j].PurchaseDeliveryNoteId {
					ok = false
					break
				}
			}
			if ok {
				note := getPurchaseDeliveryNoteRow(*noteDetails[j].PurchaseDeliveryNoteId)
				notes = append(notes, note)
			}
		}
	}
	return notes
}

// Add an amount to the lines_number field in the purchase order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderLinesNumber(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.LinesNumber += 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// Takes out an amount to the lines_number field in the purchase order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderLinesNumber(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.LinesNumber -= 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// Add an amount to the invoiced_lines field in the purchase order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderInvoicedLines(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.InvoicedLines += 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// Takes out an amount to the invoiced_lines field in the purchase order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderInvoicedLines(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.InvoicedLines -= 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// Add an amount to the delivery_note_lines field in the purchase order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addPurchaseOrderDeliveryNoteLines(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.DeliveryNoteLines += 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}

// Takes out an amount to the delivery_note_lines field in the purchase order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removePurchaseOrderDeliveryNoteLines(orderId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseOrder PurchaseOrder
	result := trans.Model(&PurchaseOrder{}).Where("id = ?", orderId).First(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseOrder.DeliveryNoteLines -= 1

	result = trans.Save(&purchaseOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_order", int(orderId), userId, "U")
	json, _ := json.Marshal(purchaseOrder)
	go fireWebHook(enterpriseId, "purchase_order", "PUT", string(json))

	return true
}
