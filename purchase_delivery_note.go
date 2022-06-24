/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

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

type PurchaseDeliveryNote struct {
	Id                 int64         `json:"id" gorm:"index:purchase_delivery_note_id_enterprise,unique:true,priority:1"`
	SupplierId         int32         `json:"supplierId" gorm:"column:supplier;not null:true"`
	Supplier           Supplier      `json:"supplier" gorm:"foreignKey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated        time.Time     `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone;index:purchase_delivery_note_date_created,sort:desc"`
	PaymentMethodId    int32         `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod      PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId    string        `json:"billingSeriesId" gorm:"column:billing_series;type:character(3);not null:true;index:purchase_delivery_note_delivery_note_number,unique:true,priority:2"`
	BillingSeries      BillingSerie  `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	ShippingAddressId  int32         `json:"shippingAddressId" gorm:"column:shipping_address;not null:true"`
	ShippingAddress    Address       `json:"shippingAddress" gorm:"foreignKey:ShippingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	TotalProducts      float64       `json:"totalProducts" gorm:"column:total_products;not null:true;type:numeric(14,6)"`
	DiscountPercent    float64       `json:"discountPercent" gorm:"column:discount_percent;not null:true;type:numeric(14,6)"`
	FixDiscount        float64       `json:"fixDiscount" gorm:"column:fix_discount;not null:true;type:numeric(14,6)"`
	ShippingPrice      float64       `json:"shippingPrice" gorm:"column:shipping_price;not null:true;type:numeric(14,6)"`
	ShippingDiscount   float64       `json:"shippingDiscount" gorm:"column:shipping_discount;not null:true;type:numeric(14,6)"`
	TotalWithDiscount  float64       `json:"totalWithDiscount" gorm:"column:total_with_discount;not null:true;type:numeric(14,6)"`
	TotalVat           float64       `json:"totalVat" gorm:"column:total_vat;not null:true;type:numeric(14,6)"`
	TotalAmount        float64       `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	LinesNumber        int16         `json:"linesNumber" gorm:"column:lines_number;not null:true"`
	DeliveryNoteName   string        `json:"deliveryNoteName" gorm:"column:delivery_note_name;not null:true;type:character(15)"`
	DeliveryNoteNumber int32         `json:"deliveryNoteNumber" gorm:"column:delivery_note_number;not null:true;index:purchase_delivery_note_delivery_note_number,unique:true,priority:3"`
	CurrencyId         int32         `json:"currencyId" gorm:"column:currency;not null:true"`
	Currency           Currency      `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange     float64       `json:"currencyChange" gorm:"column:currency_change;not null:true;type:numeric(14,6)"`
	EnterpriseId       int32         `json:"-" gorm:"column:enterprise;not null:true;index:purchase_delivery_note_id_enterprise,unique:true,priority:2;index:purchase_delivery_note_delivery_note_number,unique:true,priority:1"`
	Enterprise         Settings      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *PurchaseDeliveryNote) TableName() string {
	return "purchase_delivery_note"
}

type PurchaseDeliveryNotes struct {
	Rows   int64                      `json:"rows"`
	Notes  []PurchaseDeliveryNote     `json:"notes"`
	Footer PurchaseDeliveryNoteFooter `json:"footer"`
}

type PurchaseDeliveryNoteFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseDeliveryNotes(enterpriseId int32) PurchaseDeliveryNotes {
	dn := PurchaseDeliveryNotes{}
	dn.Notes = make([]PurchaseDeliveryNote, 0)
	result := dbOrm.Model(&PurchaseDeliveryNote{}).Where("purchase_delivery_note.enterprise = ?", enterpriseId).Order("purchase_delivery_note.date_created DESC").Preload(clause.Associations).Find(&dn.Notes)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return dn
	}
	dn.Footer = PurchaseDeliveryNoteFooter{}
	result = dbOrm.Model(&PurchaseDeliveryNote{}).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Where("purchase_delivery_note.enterprise = ?", enterpriseId).Count(&dn.Rows).Scan(&dn.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return dn
	}
	return dn
}

type PurchaseDeliveryNoteSearch struct {
	PaginatedSearch
	DateStart     *time.Time `json:"dateStart"`
	DateEnd       *time.Time `json:"dateEnd"`
	NotPosted     bool       `json:"notPosted"`
	BillingSeries *string    `json:"billingSeries"`
}

func (s *PurchaseDeliveryNoteSearch) searchPurchaseDeliveryNote() PurchaseDeliveryNotes {
	dn := PurchaseDeliveryNotes{}
	dn.Notes = make([]PurchaseDeliveryNote, 0)
	cursor := dbOrm.Model(&PurchaseDeliveryNote{}).Where("purchase_delivery_note.enterprise = ?", s.enterprise)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("purchase_delivery_note.delivery_note_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_delivery_note.supplier").Where("purchase_delivery_note.delivery_note_name LIKE @search OR suppliers.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_delivery_note.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_delivery_note.date_created <= ?", s.DateEnd)
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("purchase_delivery_note.billing_series = ?", *s.BillingSeries)
		}
	}
	result := cursor.Order("purchase_delivery_note.date_created DESC").Preload(clause.Associations).Count(&dn.Rows).Find(&dn.Notes)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return dn
	}

	cursor = dbOrm.Model(&PurchaseDeliveryNote{}).Where("purchase_delivery_note.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("purchase_delivery_note.delivery_note_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_delivery_note.supplier").Where("purchase_delivery_note.delivery_note_name LIKE @search OR suppliers.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_delivery_note.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_delivery_note.date_created <= ?", s.DateEnd)
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("purchase_delivery_note.billing_series = ?", *s.BillingSeries)
		}
	}
	result = cursor.Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&dn.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return dn
}

func getPurchaseDeliveryNoteRow(deliveryNoteId int64) PurchaseDeliveryNote {
	p := PurchaseDeliveryNote{}
	result := dbOrm.Model(&PurchaseDeliveryNote{}).Where("purchase_delivery_note.id = ?", deliveryNoteId).Preload(clause.Associations).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return p
	}
	return p
}

func (n *PurchaseDeliveryNote) isValid() bool {
	return !(n.SupplierId <= 0 || n.PaymentMethodId <= 0 || len(n.BillingSeriesId) == 0 || len(n.BillingSeriesId) > 3 || n.ShippingAddressId <= 0)
}

func (c *PurchaseDeliveryNote) BeforeCreate(tx *gorm.DB) (err error) {
	var purchaseDeliveryNote PurchaseDeliveryNote
	tx.Model(&PurchaseDeliveryNote{}).Last(&purchaseDeliveryNote)
	c.Id = purchaseDeliveryNote.Id + 1
	return nil
}

func (n *PurchaseDeliveryNote) insertPurchaseDeliveryNotes(userId int32, trans *gorm.DB) (bool, int64) {
	if !n.isValid() {
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

	n.DeliveryNoteNumber = getNextPurchaseDeliveryNoteNumber(n.BillingSeriesId, n.EnterpriseId)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.CurrencyId)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)
	n.DateCreated = now
	n.TotalProducts = 0
	n.TotalWithDiscount = 0
	n.TotalVat = 0
	n.TotalAmount = 0
	n.LinesNumber = 0

	result := trans.Create(&n)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(n.EnterpriseId, "purchase_delivery_note", int(n.Id), userId, "I")
	json, _ := json.Marshal(n)
	go fireWebHook(n.EnterpriseId, "purchase_delivery_note", "POST", string(json))

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false, 0
		}
		///
	}

	return true, n.Id
}

func (n *PurchaseDeliveryNote) deletePurchaseDeliveryNotes(userId int32, trans *gorm.DB) bool {
	if n.Id <= 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	inMemoryNote := getPurchaseDeliveryNoteRow(n.Id)
	if inMemoryNote.EnterpriseId != n.EnterpriseId {
		return false
	}

	d := getWarehouseMovementByPurchaseDeliveryNote(n.Id, n.EnterpriseId)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	insertTransactionalLog(n.EnterpriseId, "purchase_delivery_note", int(n.Id), userId, "D")
	json, _ := json.Marshal(n)
	go fireWebHook(n.EnterpriseId, "purchase_delivery_note", "DELETE", string(json))

	result := trans.Delete(&PurchaseDeliveryNote{}, "id = ? AND enterprise = ?", n.Id, n.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false
		}
		///
	}

	return true
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. There are no details to generate the delivery note
func deliveryNoteAllPurchaseOrder(purchaseOrderId int64, enterpriseId int32, userId int32) (OkAndErrorCodeReturn, int64) {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	if purchaseOrder.EnterpriseId != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if purchaseOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if purchaseOrder.DeliveryNoteLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}, 0
	}
	orderDetails := getPurchaseOrderDetail(purchaseOrderId, purchaseOrder.EnterpriseId)
	filterPurchaseOrderDetails(orderDetails, func(pod PurchaseOrderDetail) bool { return pod.QuantityDeliveryNote < pod.Quantity })
	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}

	// create a delivery note for that order
	n := PurchaseDeliveryNote{}
	n.SupplierId = purchaseOrder.SupplierId
	n.ShippingAddressId = purchaseOrder.ShippingAddressId
	n.CurrencyId = purchaseOrder.CurrencyId
	n.PaymentMethodId = purchaseOrder.PaymentMethodId
	n.BillingSeriesId = purchaseOrder.BillingSeriesId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	///

	n.EnterpriseId = enterpriseId
	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "I"
		movement.WarehouseId = orderDetail.WarehouseId
		movement.ProductId = orderDetail.ProductId
		movement.Quantity = orderDetail.Quantity
		movement.PurchaseDeliveryNoteId = &deliveryNoteId
		movement.PurchaseOrderDetailId = &orderDetail.Id
		movement.PurchaseOrderId = &purchaseOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.EnterpriseId = enterpriseId
		ok = movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
	}

	///
	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}, deliveryNoteId
	///
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail has a delivery note generated
// 4. The selected quantity is greater than the quantity pending of delivery note generation in the detail
func (noteInfo *OrderDetailGenerate) deliveryNotePartiallyPurchaseOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(noteInfo.OrderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.EnterpriseId != enterpriseId || len(noteInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.DeliveryNoteLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(noteInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.OrderId != noteInfo.OrderId || noteInfo.Selection[i].Quantity == 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if noteInfo.Selection[i].Quantity > orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2, ExtraData: []string{product.Name}}
		}
		if orderDetail.QuantityDeliveryNote >= orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{product.Name}}
		}
		if (noteInfo.Selection[i].Quantity + orderDetail.QuantityDeliveryNote) > orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 4, ExtraData: []string{product.Name}}
		}
		purchaseOrderDetails = append(purchaseOrderDetails, orderDetail)
	}

	// create a delivery note for that order
	n := PurchaseDeliveryNote{}
	n.SupplierId = purchaseOrder.SupplierId
	n.ShippingAddressId = purchaseOrder.ShippingAddressId
	n.CurrencyId = purchaseOrder.CurrencyId
	n.PaymentMethodId = purchaseOrder.PaymentMethodId
	n.BillingSeriesId = purchaseOrder.BillingSeriesId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	n.EnterpriseId = enterpriseId
	ok, deliveryNoteId := n.insertPurchaseDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(purchaseOrderDetails); i++ {
		orderDetail := purchaseOrderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "I"
		movement.WarehouseId = orderDetail.WarehouseId
		movement.ProductId = orderDetail.ProductId
		movement.Quantity = noteInfo.Selection[i].Quantity
		movement.PurchaseDeliveryNoteId = &deliveryNoteId
		movement.PurchaseOrderDetailId = &orderDetail.Id
		movement.PurchaseOrderId = &purchaseOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.EnterpriseId = enterpriseId
		ok = movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}
	///
}

type PurchaseDeliveryNoteRelation struct {
	Orders []PurchaseOrder `json:"orders"`
}

func getPurchaseDeliveryNoteRelations(noteId int64, enterpriseId int32) PurchaseDeliveryNoteRelation {
	return PurchaseDeliveryNoteRelation{
		Orders: getPurchaseDeliveryNoteOrders(noteId, enterpriseId),
	}
}

func getPurchaseDeliveryNoteOrders(noteId int64, enterpriseId int32) []PurchaseOrder {
	deliveryNote := getPurchaseDeliveryNoteRow(noteId)
	if deliveryNote.Id <= 0 || deliveryNote.EnterpriseId != enterpriseId {
		return make([]PurchaseOrder, 0)
	}
	var details []WarehouseMovement = make([]WarehouseMovement, 0)
	dbOrm.Model(&WarehouseMovement{}).Where("purchase_delivery_note = ?", noteId).Distinct("purchase_order").Find(&details)
	var orders []PurchaseOrder = make([]PurchaseOrder, 0)
	for i := 0; i < len(details); i++ {
		if details[i].PurchaseOrderId != nil {
			var order PurchaseOrder = getPurchaseOrderRow(*details[i].PurchaseOrderId)
			orders = append(orders, order)
		}
	}
	return orders
}

// Adds a total amount to the delivery note total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseDeliveryNote(noteId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseDeliveryNote PurchaseDeliveryNote
	result := trans.Model(&PurchaseDeliveryNote{}).Where("id = ?", noteId).First(&purchaseDeliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseDeliveryNote.TotalProducts += totalAmount
	purchaseDeliveryNote.TotalVat += (totalAmount / 100) * vatPercent

	result = trans.Save(&purchaseDeliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseDeliveryNote(noteId, enterpriseId, userId, trans)
}

// Applies the logic to calculate the totals of the delivery note.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseDeliveryNote(noteId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var deliveryNote PurchaseDeliveryNote
	result := trans.Model(&PurchaseDeliveryNote{}).Where("id = ?", noteId).First(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	deliveryNote.TotalWithDiscount = deliveryNote.TotalProducts - (deliveryNote.TotalProducts * (deliveryNote.DiscountPercent / 100)) - deliveryNote.FixDiscount + deliveryNote.ShippingPrice - deliveryNote.ShippingDiscount
	deliveryNote.TotalAmount = deliveryNote.TotalWithDiscount + deliveryNote.TotalVat

	result = trans.Save(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_delivery_note", int(noteId), userId, "U")
	json, _ := json.Marshal(deliveryNote)
	go fireWebHook(enterpriseId, "purchase_delivery_note", "PUT", string(json))

	return true
}
