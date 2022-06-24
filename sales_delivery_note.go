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

type SalesDeliveryNote struct {
	Id                 int64         `json:"id" gorm:"index:sales_delivery_note_id_enterprise,unique:true,priority:1"`
	CustomerId         int32         `json:"customerId" gorm:"column:customer;not null"`
	Customer           Customer      `json:"customer" gorm:"foreignKey:CustomerId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated        time.Time     `json:"dateCreated" gorm:"column:date_created;not null;type:timestamp(3) with time zone;index:sales_delivery_note_date_created,sort:desc"`
	PaymentMethodId    int32         `json:"paymentMethodId" gorm:"column:payment_method;not null"`
	PaymentMethod      PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId    string        `json:"billingSeriesId" gorm:"column:billing_series;type:character(3);not null;index:sales_delivery_note_delivery_note_number,unique:true,priority:2"`
	BillingSeries      BillingSerie  `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	ShippingAddressId  int32         `json:"shippingAddressId" gorm:"column:shipping_address;not null"`
	ShippingAddress    Address       `json:"shippingAddress" gorm:"foreignKey:ShippingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	TotalProducts      float64       `json:"totalProducts" gorm:"column:total_products;not null;type:numeric(14,6)"`
	DiscountPercent    float64       `json:"discountPercent" gorm:"column:discount_percent;not null;type:numeric(14,6)"`
	FixDiscount        float64       `json:"fixDiscount" gorm:"column:fix_discount;not null;type:numeric(14,6)"`
	ShippingPrice      float64       `json:"shippingPrice" gorm:"column:shipping_price;not null;type:numeric(14,6)"`
	ShippingDiscount   float64       `json:"shippingDiscount" gorm:"column:shipping_discount;not null;type:numeric(14,6)"`
	TotalWithDiscount  float64       `json:"totalWithDiscount" gorm:"column:total_with_discount;not null;type:numeric(14,6)"`
	VatAmount          float64       `json:"vatAmount" gorm:"column:vat_amount;not null;type:numeric(14,6)"`
	TotalAmount        float64       `json:"totalAmount" gorm:"column:total_amount;not null;type:numeric(14,6)"`
	LinesNumber        int16         `json:"linesNumber" gorm:"column:lines_number;not null"`
	DeliveryNoteName   string        `json:"deliveryNoteName" gorm:"column:delivery_note_name;not null;type:character(15)"`
	DeliveryNoteNumber int32         `json:"deliveryNoteNumber" gorm:"column:delivery_note_number;not null;index:sales_delivery_note_delivery_note_number,unique:true,sort:desc,priority:3"`
	CurrencyId         int32         `json:"currencyId" gorm:"column:currency;not null"`
	Currency           Currency      `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange     float64       `json:"currencyChange" gorm:"column:currency_change;not null;type:numeric(14,6)"`
	EnterpriseId       int32         `json:"-" gorm:"column:enterprise;not null:true;index:sales_delivery_note_id_enterprise,unique:true,priority:2;;index:sales_delivery_note_delivery_note_number,unique:true,priority:1"`
	Enterprise         Settings      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (n *SalesDeliveryNote) TableName() string {
	return "sales_delivery_note"
}

type SalesDeliveryNotes struct {
	Rows   int64                   `json:"rows"`
	Notes  []SalesDeliveryNote     `json:"notes"`
	Footer SalesDeliveryNoteFooter `json:"footer"`
}

type SalesDeliveryNoteFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func (q *PaginationQuery) getSalesDeliveryNotes() SalesDeliveryNotes {
	sd := SalesDeliveryNotes{}
	if !q.isValid() {
		return sd
	}

	sd.Notes = make([]SalesDeliveryNote, 0)
	result := dbOrm.Model(&SalesDeliveryNote{}).Where("enterprise = ?", q.enterprise).Offset(int(q.Offset)).Limit(int(q.Limit)).Order("date_created DESC").Preload(clause.Associations).Find(&sd.Notes)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return sd
	}
	sd.Footer = SalesDeliveryNoteFooter{}
	result = dbOrm.Model(&SalesDeliveryNote{}).Count(&sd.Rows).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Where("enterprise = ?", q.enterprise).Scan(&sd.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return sd
	}

	return sd
}

func getSalesDeliveryNoteRow(deliveryNoteId int64) SalesDeliveryNote {
	n := SalesDeliveryNote{}
	result := dbOrm.Model(&SalesDeliveryNote{}).Where("id = ?", deliveryNoteId).Preload(clause.Associations).First(&n)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return n
	}
	return n
}

type SalesDeliveryNoteSearch struct {
	PaginatedSearch
	DateStart     *time.Time `json:"dateStart"`
	DateEnd       *time.Time `json:"dateEnd"`
	NotPosted     bool       `json:"notPosted"`
	BillingSeries *string    `json:"billingSeries"`
}

func (s *SalesDeliveryNoteSearch) searchSalesDelvieryNotes() SalesDeliveryNotes {
	sd := SalesDeliveryNotes{}
	if !s.isValid() {
		return sd
	}

	cursor := dbOrm.Model(&SalesDeliveryNote{}).Where("sales_delivery_note.enterprise = ?", s.enterprise)
	sd.Notes = make([]SalesDeliveryNote, 0)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("sales_delivery_note.delivery_note_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_delivery_note.customer").Where("sales_delivery_note.delivery_note_name LIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_delivery_note.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_delivery_note.date_created <= ?", s.DateEnd)
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("sales_delivery_note.billing_series = ?", *s.BillingSeries)
		}
	}
	result := cursor.Offset(int(s.Offset)).Limit(int(s.Limit)).Order("sales_delivery_note.date_created DESC").Preload(clause.Associations).Count(&sd.Rows).Find(&sd.Notes)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return sd
	}

	cursor = dbOrm.Model(&SalesDeliveryNote{}).Where("sales_delivery_note.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("sales_delivery_note.delivery_note_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_delivery_note.customer").Where("sales_delivery_note.delivery_note_name LIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_delivery_note.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_delivery_note.date_created <= ?", s.DateEnd)
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("sales_delivery_note.billing_series = ?", *s.BillingSeries)
		}
	}
	result = cursor.Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&sd.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return sd
}

func (n *SalesDeliveryNote) isValid() bool {
	return !(n.CustomerId <= 0 || n.PaymentMethodId <= 0 || len(n.BillingSeriesId) == 0 || len(n.BillingSeriesId) > 3 || n.ShippingAddressId <= 0)
}

func (s *SalesDeliveryNote) BeforeCreate(tx *gorm.DB) (err error) {
	var salesDeliveryNote SalesDeliveryNote
	tx.Model(&SalesDeliveryNote{}).Last(&salesDeliveryNote)
	s.Id = salesDeliveryNote.Id + 1
	return nil
}

func (n *SalesDeliveryNote) insertSalesDeliveryNotes(userId int32, trans *gorm.DB) (bool, int64) {
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

	n.DeliveryNoteNumber = getNextSaleDeliveryNoteNumber(n.BillingSeriesId, n.EnterpriseId)
	if n.DeliveryNoteNumber <= 0 {
		return false, 0
	}
	n.CurrencyChange = getCurrencyExchange(n.CurrencyId)
	now := time.Now()
	n.DeliveryNoteName = n.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", n.DeliveryNoteNumber)

	n.DateCreated = time.Now()
	n.TotalProducts = 0
	n.DiscountPercent = 0
	n.FixDiscount = 0
	n.ShippingPrice = 0
	n.ShippingDiscount = 0
	n.TotalWithDiscount = 0
	n.VatAmount = 0
	n.TotalAmount = 0
	n.LinesNumber = 0

	result := trans.Create(&n)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(n.EnterpriseId, "sales_delivery_note", int(n.Id), userId, "I")
	json, _ := json.Marshal(n)
	go fireWebHook(n.EnterpriseId, "sales_delivery_note", "POST", string(json))

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

// ERROR CODES:
// 1. A shipping is associated to this delivery note
func (n *SalesDeliveryNote) deleteSalesDeliveryNotes(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if n.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	shipping := getSalesDeliveryNoteShippings(n.Id)
	if len(shipping) > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	d := getWarehouseMovementBySalesDeliveryNote(n.Id, n.EnterpriseId)
	for i := 0; i < len(d); i++ {
		ok := d[i].deleteWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(n.EnterpriseId, "sales_delivery_note", int(n.Id), userId, "D")
	json, _ := json.Marshal(n)
	go fireWebHook(n.EnterpriseId, "sales_delivery_note", "DELETE", string(json))

	result := trans.Delete(&SalesDeliveryNote{}, "id = ? AND enterprise = ?", n.Id, n.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	return OkAndErrorCodeReturn{Ok: true}
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. There are no details to generate the delivery note
func deliveryNoteAllSaleOrder(saleOrderId int64, enterpriseId int32, userId int32, trans *gorm.DB) (OkAndErrorCodeReturn, int64) {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.EnterpriseId != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if saleOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	if saleOrder.DeliveryNoteLines >= saleOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}, 0
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.EnterpriseId)
	filterSalesOrderDetails(orderDetails, func(sod SalesOrderDetail) bool { return sod.QuantityDeliveryNote < sod.Quantity })
	if len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}, 0
	}

	// create a delivery note for that order
	n := SalesDeliveryNote{}
	n.CustomerId = saleOrder.CustomerId
	n.ShippingAddressId = saleOrder.ShippingAddressId
	n.CurrencyId = saleOrder.CurrencyId
	n.PaymentMethodId = saleOrder.PaymentMethodId
	n.BillingSeriesId = saleOrder.BillingSeriesId

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
		///
	}

	n.EnterpriseId = enterpriseId
	ok, deliveryNoteId := n.insertSalesDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}, 0
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "O"
		movement.WarehouseId = orderDetail.WarehouseId
		movement.ProductId = orderDetail.ProductId
		movement.Quantity = -(orderDetail.Quantity - orderDetail.QuantityDeliveryNote)
		movement.SalesDeliveryNoteId = &deliveryNoteId
		movement.SalesOrderDetailId = &orderDetail.Id
		movement.SalesOrderId = &saleOrder.Id
		movement.Price = orderDetail.Price
		movement.VatPercent = orderDetail.VatPercent
		movement.EnterpriseId = enterpriseId
		ok = movement.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
	}

	if beginTransaction {
		///
		result := trans.Commit()
		if result.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}, 0
		}
		///
	}
	return OkAndErrorCodeReturn{Ok: true}, deliveryNoteId
}

// ERROR CODES:
// 1. The order already has a delivery note generated
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail has a delivery note generated
// 4. The selected quantity is greater than the quantity pending of delivery note generation in the detail
func (noteInfo *OrderDetailGenerate) deliveryNotePartiallySaleOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(noteInfo.OrderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId || len(noteInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if saleOrder.DeliveryNoteLines >= saleOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(noteInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(noteInfo.Selection[i].Id)
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
		saleOrderDetails = append(saleOrderDetails, orderDetail)
	}

	// create a delivery note for that order
	n := SalesDeliveryNote{}
	n.CustomerId = saleOrder.CustomerId
	n.ShippingAddressId = saleOrder.ShippingAddressId
	n.CurrencyId = saleOrder.CurrencyId
	n.PaymentMethodId = saleOrder.PaymentMethodId
	n.BillingSeriesId = saleOrder.BillingSeriesId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	n.EnterpriseId = enterpriseId
	ok, deliveryNoteId := n.insertSalesDeliveryNotes(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(saleOrderDetails); i++ {
		orderDetail := saleOrderDetails[i]
		movement := WarehouseMovement{}
		movement.Type = "O"
		movement.WarehouseId = orderDetail.WarehouseId
		movement.ProductId = orderDetail.ProductId
		movement.Quantity = -noteInfo.Selection[i].Quantity
		movement.SalesDeliveryNoteId = &deliveryNoteId
		movement.SalesOrderDetailId = &orderDetail.Id
		movement.SalesOrderId = &saleOrder.Id
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

type SalesDeliveryNoteLocate struct {
	Id               int64     `json:"id" gorm:"column:date_created;not null;type:timestamp(3) with time zone"`
	CustomerId       int32     `json:"customerId" gorm:"column:customer;not null"`
	Customer         Customer  `json:"customer"`
	DateCreated      time.Time `json:"dateCreated"`
	DeliveryNoteName string    `json:"deliveryNoteName" gorm:"column:delivery_note_name;not null;type:character(15)"`
}

func locateSalesDeliveryNotesBySalesOrder(orderId int64, enterpriseId int32) []SalesDeliveryNoteLocate {
	var deliveryNotes []SalesDeliveryNoteLocate = make([]SalesDeliveryNoteLocate, 0)
	relations := getSalesOrderDeliveryNotes(orderId, enterpriseId)
	for i := 0; i < len(relations); i++ {
		deliveryNotes = append(deliveryNotes, SalesDeliveryNoteLocate{
			Id:               relations[i].Id,
			CustomerId:       relations[i].CustomerId,
			Customer:         getCustomerRow(relations[i].CustomerId),
			DateCreated:      relations[i].DateCreated,
			DeliveryNoteName: relations[i].DeliveryNoteName,
		})
	}
	return deliveryNotes
}

type SalesDeliveryNoteRelation struct {
	Orders    []SaleOrder `json:"orders"`
	Shippings []Shipping  `json:"shippings"`
}

func getSalesDeliveryNoteRelations(noteId int64, enterpriseId int32) SalesDeliveryNoteRelation {
	return SalesDeliveryNoteRelation{
		Orders:    getSalesDeliveryNoteOrders(noteId, enterpriseId),
		Shippings: getSalesDeliveryNoteShippings(noteId),
	}
}

func getSalesDeliveryNoteOrders(noteId int64, enterpriseId int32) []SaleOrder {
	deliveryNote := getSalesDeliveryNoteRow(noteId)
	if deliveryNote.Id <= 0 || deliveryNote.EnterpriseId != enterpriseId {
		return make([]SaleOrder, 0)
	}
	var details []WarehouseMovement = make([]WarehouseMovement, 0)
	dbOrm.Model(&WarehouseMovement{}).Where("sales_delivery_note = ?", noteId).Distinct("sales_order").Find(&details)
	var sales []SaleOrder = make([]SaleOrder, 0)
	for i := 0; i < len(details); i++ {
		if details[i].SalesOrderId != nil {
			var order SaleOrder = getSalesOrderRow(*details[i].SalesOrderId)
			sales = append(sales, order)
		}
	}
	return sales
}

func getSalesDeliveryNoteShippings(noteId int64) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	result := dbOrm.Model(&Shipping{}).Where("delivery_note = ?", noteId).Order("id ASC").Preload(clause.Associations).Find(&shippings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return shippings
	}
	return shippings
}

// Adds a total amount to the delivery note total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesDeliveryNote(noteId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var deliveryNote SalesDeliveryNote
	result := trans.Model(&SalesDeliveryNote{}).Where("id = ?", noteId).First(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	deliveryNote.TotalProducts += totalAmount
	deliveryNote.VatAmount += (totalAmount / 100) * vatPercent

	result = trans.Save(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsSaleDeliveryNote(noteId, enterpriseId, userId, trans)
}

// Applies the logic to calculate the totals of the sales delivery note.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleDeliveryNote(noteId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var deliveryNote SalesDeliveryNote
	result := trans.Model(&SalesDeliveryNote{}).Where("id = ?", noteId).First(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	deliveryNote.TotalWithDiscount = deliveryNote.TotalProducts - (deliveryNote.TotalProducts * (deliveryNote.DiscountPercent / 100)) - deliveryNote.FixDiscount + deliveryNote.ShippingPrice - deliveryNote.ShippingDiscount
	deliveryNote.TotalAmount = deliveryNote.TotalWithDiscount + deliveryNote.VatAmount

	result = trans.Save(&deliveryNote)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_delivery_note", int(noteId), userId, "U")
	json, _ := json.Marshal(deliveryNote)
	go fireWebHook(enterpriseId, "sales_delivery_note", "PUT", string(json))

	return true
}
