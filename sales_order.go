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

type SaleOrder struct {
	Id                  int64         `json:"id" gorm:"index:sales_order_id_enterprise,unique:true,priority:1"`
	Reference           string        `json:"reference" gorm:"type:character varying(15);not null:true;index:sales_order_reference,type:gin"`
	CustomerId          int32         `json:"customerId" gorm:"type:integer;not null:true;column:customer"`
	Customer            Customer      `json:"customer" gorm:"foreignKey:CustomerId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated         time.Time     `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true;index:sales_order_date_created,sort:desc"`
	DatePaymentAccepted *time.Time    `json:"datePaymentAccepted" gorm:"type:timestamp(3) with time zone"`
	PaymentMethodId     int32         `json:"paymentMethodId" gorm:"column:payment_method;type:integer;not null:true"`
	PaymentMethod       PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId     string        `json:"billingSeriesId" gorm:"type:character(3);not null:true;column:billing_series;index:sales_order_order_number,unique:true,priority:2"`
	BillingSeries       BillingSerie  `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyId          int32         `json:"currencyId" gorm:"column:currency;type:integer;not null:true"`
	Currency            Currency      `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange      float64       `json:"currencyChange" gorm:"type:numeric(14,6);not null:true"`
	BillingAddressId    int32         `json:"billingAddressId" gorm:"type:integer;not null:true;column:billing_address"`
	BillingAddress      Address       `json:"billingAddress" gorm:"foreignKey:BillingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	ShippingAddressId   int32         `json:"shippingAddressId" gorm:"type:integer;not null:true;column:shipping_address"`
	ShippingAddress     Address       `json:"shippingAddress" gorm:"foreignKey:ShippingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	LinesNumber         int16         `json:"linesNumber" gorm:"not null:true"`
	InvoicedLines       int16         `json:"invoicedLines" gorm:"not null:true"`
	DeliveryNoteLines   int16         `json:"deliveryNoteLines" gorm:"not null:true"`
	TotalProducts       float64       `json:"totalProducts" gorm:"not null:true;type:numeric(14,6)"`
	DiscountPercent     float64       `json:"discountPercent" gorm:"not null:true;type:numeric(14,6)"`
	FixDiscount         float64       `json:"fixDiscount" gorm:"not null:true;type:numeric(14,6)"`
	ShippingPrice       float64       `json:"shippingPrice" gorm:"not null:true;type:numeric(14,6)"`
	ShippingDiscount    float64       `json:"shippingDiscount" gorm:"not null:true;type:numeric(14,6)"`
	TotalWithDiscount   float64       `json:"totalWithDiscount" gorm:"not null:true;type:numeric(14,6)"`
	VatAmount           float64       `json:"vatAmount" gorm:"not null:true;type:numeric(14,6)"`
	TotalAmount         float64       `json:"totalAmount" gorm:"not null:true;type:numeric(14,6)"`
	Description         string        `json:"description" gorm:"type:text;not null:true;column:dsc"`
	Notes               string        `json:"notes" gorm:"type:character varying(250);not null:true;column:notes"`
	Cancelled           bool          `json:"cancelled" gorm:"not null:true"`
	Status              string        `json:"status" gorm:"type:character(1);not null:true;column:status"` // _ = Waiting for payment, A = Waiting for purchase order, B = Purchase order pending, C = Waiting for manufacturing orders, D = Manufacturing orders pending, E = Sent to preparation, F = Awaiting for shipping, G = Shipped, H = Receiced by the customer, Z = Cancelled
	OrderNumber         int32         `json:"orderNumber" gorm:"not null:true;column:order_number;index:sales_order_order_number,unique:true,priority:3,sort:desc"`
	OrderName           string        `json:"orderName" gorm:"type:character(15);not null:true"`
	CarrierId           *int32        `json:"carrierId" gorm:"column:carrier"`
	Carrier             *Carrier      `json:"carrier" gorm:"foreignKey:CarrierId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopId        int32         `json:"-" gorm:"column:ps_id;not null:true;index:sales_order_ps_id,unique:true,priority:2,where:ps_id <> 0"`
	WooCommerceId       int32         `json:"-" gorm:"column:wc_id;not null:true;index:sales_order_wc_id,unique:true,priority:2,where:wc_id <> 0"`
	ShopifyId           int64         `json:"-" gorm:"column:sy_id;not null:true;index:sales_order_sy_id,unique:true,priority:2,where:sy_id <> 0"`
	ShopifyDraftId      int64         `json:"-" gorm:"column:sy_draft_id;not null:true;index:sales_order_sy_draft_id,unique:true,priority:2,where:sy_draft_id <> 0"`
	EnterpriseId        int32         `json:"-" gorm:"column:enterprise;not null:true;index:sales_order_id_enterprise,unique:true,priority:2;index:sales_order_order_number,unique:true,priority:1;index:sales_order_ps_id,unique:true,priority:1,where:ps_id <> 0;index:sales_order_sy_draft_id,unique:true,priority:1,where:sy_draft_id <> 0;index:sales_order_sy_id,unique:true,priority:1,where:sy_id <> 0;index:sales_order_wc_id,unique:true,priority:1,where:wc_id <> 0"`
	Enterprise          Settings      `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (so *SaleOrder) TableName() string {
	return "sales_order"
}

type SaleOrders struct {
	Rows   int64            `json:"rows"`
	Orders []SaleOrder      `json:"orders"`
	Footer SalesOrderFooter `json:"footer"`
}

type SalesOrderFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func (q *PaginationQuery) getSalesOrder(enterpriseId int32) SaleOrders {
	so := SaleOrders{}
	if !q.isValid() {
		return so
	}

	so.Orders = make([]SaleOrder, 0)
	// get all the sale order in the database for the current enterprise with the pagination sort by date created descending and limit the number of rows using dbOrm
	dbOrm.Where("enterprise = ?", enterpriseId).Order("date_created DESC").Limit(int(q.Limit)).Offset(int(q.Offset)).Preload(clause.Associations).Find(&so.Orders)
	// get the total number of sale order in the database for the current enterprise
	dbOrm.Model(&SaleOrder{}).Where("enterprise = ?", enterpriseId).Count(&so.Rows)
	// get the total amount of sale order in the database for the current enterprise
	dbOrm.Model(&SaleOrder{}).Where("enterprise = ?", enterpriseId).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&so.Footer)

	return so
}

type SalesOrderSearch struct {
	PaginatedSearch
	DateStart          *time.Time `json:"dateStart"`
	DateEnd            *time.Time `json:"dateEnd"`
	Status             string     `json:"status"`
	InvoicedStatus     string     `json:"invoicedStatus"`     // "" = All, A = Invoiced, B = Not invoiced, C = Partially invoiced
	DeliveryNoteStatus string     `json:"deliveryNoteStatus"` // "" = All, A = Delivered, B = Not delivered, C = Partially delivered
}

func (s *SalesOrderSearch) searchSalesOrder() SaleOrders {
	so := SaleOrders{}
	if !s.isValid() {
		return so
	}

	so.Orders = make([]SaleOrder, 0)
	// get all the sale order in the database for the current enterprise with the pagination sort by date created descending and limit the number of rows using dbOrm
	cursor := dbOrm.Where("sales_order.enterprise = ?", s.enterprise)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("sales_order.order_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_order.customer").Where("sales_order.order_name LIKE @search OR sales_order.reference ILIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_order.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_order.date_created <= ?", s.DateEnd)
		}
		if s.Status != "" {
			cursor = cursor.Where("sales_order.sales_order.status = ?", s.Status)
		}
		if s.InvoicedStatus != "" {
			if s.InvoicedStatus == "A" {
				cursor = cursor.Where("sales_order.lines_number = sales_order.invoiced_lines")
			} else if s.InvoicedStatus == "B" {
				cursor = cursor.Where("sales_order.lines_number = 0")
			} else if s.InvoicedStatus == "C" {
				cursor = cursor.Where("sales_order.lines_number < sales_order.invoiced_lines")
			}
		}
		if s.DeliveryNoteStatus != "" {
			if s.DeliveryNoteStatus == "A" {
				cursor = cursor.Where("sales_order.lines_number = sales_order.delivery_note_lines")
			} else if s.DeliveryNoteStatus == "B" {
				cursor = cursor.Where("sales_order.lines_number = 0")
			} else if s.DeliveryNoteStatus == "C" {
				cursor = cursor.Where("sales_order.lines_number < sales_order.delivery_note_lines")
			}
		}
	}
	result := cursor.Order("sales_order.date_created DESC").Limit(int(s.Limit)).Offset(int(s.Offset)).Preload(clause.Associations).Find(&so.Orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return so
	}

	// get the total number of sale order in the database for the current enterprise
	// get the total amount of sale order in the database for the current enterprise
	cursor = dbOrm.Model(&SaleOrder{}).Where("sales_order.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("sales_order.order_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_order.customer").Where("sales_order.order_name LIKE @search OR sales_order.reference ILIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_order.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_order.date_created <= ?", s.DateEnd)
		}
		if s.Status != "" {
			cursor = cursor.Where("sales_order.sales_order.status = ?", s.Status)
		}
		if s.InvoicedStatus != "" {
			if s.InvoicedStatus == "A" {
				cursor = cursor.Where("sales_order.lines_number = sales_order.invoiced_lines")
			} else if s.InvoicedStatus == "B" {
				cursor = cursor.Where("sales_order.lines_number = 0")
			} else if s.InvoicedStatus == "C" {
				cursor = cursor.Where("sales_order.lines_number < sales_order.invoiced_lines")
			}
		}
		if s.DeliveryNoteStatus != "" {
			if s.DeliveryNoteStatus == "A" {
				cursor = cursor.Where("sales_order.lines_number = sales_order.delivery_note_lines")
			} else if s.DeliveryNoteStatus == "B" {
				cursor = cursor.Where("sales_order.lines_number = 0")
			} else if s.DeliveryNoteStatus == "C" {
				cursor = cursor.Where("sales_order.lines_number < sales_order.delivery_note_lines")
			}
		}
	}
	result = cursor.Count(&so.Rows).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&so.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return so
	}

	return so
}

func getSalesOrderPreparation(enterpriseId int32) []SaleOrder {
	return getSalesOrderStatus("E", enterpriseId)
}

func getSalesOrderAwaitingShipping(enterpriseId int32) []SaleOrder {
	return getSalesOrderStatus("F", enterpriseId)
}

func getSalesOrderStatus(status string, enterpriseId int32) []SaleOrder {
	var sales []SaleOrder = make([]SaleOrder, 0)
	result := dbOrm.Where("status = ? AND enterprise = ?", status, enterpriseId).Order("date_created DESC").Preload(clause.Associations).Find(&sales)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return sales
	}

	return sales
}

func getSalesOrderRow(id int64) SaleOrder {
	var so SaleOrder = SaleOrder{}
	result := dbOrm.Where("id = ?", id).Preload(clause.Associations).First(&so)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return so
	}

	return so
}

func getSalesOrderRowTransaction(id int64, trans gorm.DB) SaleOrder {
	// get a single sale order row from the database using the transaction object
	s := SaleOrder{}
	trans.Model(&SaleOrder{}).Where("id = ?", id).First(&s)
	return s
}

func (s *SaleOrder) isValid() bool {
	return !(len(s.Reference) > 15 || s.CustomerId <= 0 || s.PaymentMethodId <= 0 || len(s.BillingSeriesId) == 0 || s.CurrencyId <= 0 || s.BillingAddressId <= 0 || s.ShippingAddressId <= 0 || len(s.Notes) > 250 || len(s.Description) > 3000)
}

func (s *SaleOrder) BeforeCreate(tx *gorm.DB) (err error) {
	var saleOrder SaleOrder
	tx.Model(&SaleOrder{}).Last(&saleOrder)
	s.Id = saleOrder.Id + 1
	return nil
}

func (s *SaleOrder) insertSalesOrder(userId int32) (bool, int64) {
	if !s.isValid() {
		return false, 0
	}

	s.OrderNumber = getNextSaleOrderNumber(s.BillingSeriesId, s.EnterpriseId)
	if s.OrderNumber <= 0 {
		return false, 0
	}
	s.CurrencyChange = getCurrencyExchange(s.CurrencyId)
	now := time.Now()
	s.OrderName = s.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", s.OrderNumber)

	s.DateCreated = time.Now()
	s.DatePaymentAccepted = nil
	s.TotalWithDiscount = s.ShippingPrice - s.ShippingDiscount - s.FixDiscount
	s.VatAmount = 0
	s.TotalAmount = s.TotalWithDiscount + s.VatAmount
	s.Cancelled = false
	s.Status = "_"
	s.LinesNumber = 0
	s.InvoicedLines = 0
	s.DeliveryNoteLines = 0
	s.TotalProducts = 0

	result := dbOrm.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0
	}

	insertTransactionalLog(s.EnterpriseId, "sales_order", int(s.Id), userId, "I")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order", "POST", string(json))

	return true, s.Id
}

func (s *SaleOrder) updateSalesOrder(userId int32) bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	// get a sale order from the database by id and enterprise id using dbOrm
	var inMemoryOrder SaleOrder
	result := dbOrm.Model(&SaleOrder{}).Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).First(&inMemoryOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if inMemoryOrder.Status == "_" { // if the payment is pending, we allow to change more fields
		if s.CurrencyId != inMemoryOrder.CurrencyId {
			s.CurrencyChange = getCurrencyExchange(s.CurrencyId)
		} else {
			s.CurrencyChange = inMemoryOrder.CurrencyChange
		}

		inMemoryOrder.CustomerId = s.CustomerId
		inMemoryOrder.PaymentMethodId = s.PaymentMethodId
		inMemoryOrder.CurrencyId = s.CurrencyId
		inMemoryOrder.CurrencyChange = s.CurrencyChange
		inMemoryOrder.BillingAddressId = s.BillingAddressId
		inMemoryOrder.ShippingAddressId = s.ShippingAddressId
		inMemoryOrder.DiscountPercent = s.CurrencyChange
		inMemoryOrder.FixDiscount = s.FixDiscount
		inMemoryOrder.ShippingPrice = s.ShippingPrice
		inMemoryOrder.ShippingDiscount = s.ShippingDiscount
		inMemoryOrder.Description = s.Description
		inMemoryOrder.Notes = s.Notes
		inMemoryOrder.Reference = s.Reference
		inMemoryOrder.CarrierId = s.CarrierId
		inMemoryOrder.ShopifyId = s.ShopifyId

	} else {
		inMemoryOrder.CustomerId = s.CustomerId
		inMemoryOrder.BillingAddressId = s.BillingAddressId
		inMemoryOrder.ShippingAddressId = s.ShippingAddressId
		inMemoryOrder.Description = s.Description
		inMemoryOrder.Notes = s.Notes
		inMemoryOrder.Reference = s.Reference
		inMemoryOrder.CarrierId = s.CarrierId
	}

	result = trans.Save(&inMemoryOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if inMemoryOrder.Status == "_" && s.DiscountPercent != inMemoryOrder.DiscountPercent || s.FixDiscount != inMemoryOrder.FixDiscount || s.ShippingPrice != inMemoryOrder.ShippingPrice || s.ShippingDiscount != inMemoryOrder.ShippingDiscount {
		ok := calcTotalsSaleOrder(s.EnterpriseId, s.Id, userId, *trans)
		if !ok {
			trans.Rollback()
			return false
		}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return false
	}
	///

	insertTransactionalLog(s.EnterpriseId, "sales_order", int(s.Id), userId, "U")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_order", "PUT", string(json))

	return true
}

// ERROR CODES
// 1. Alerady invoiced
// 2. Delivery note generated
// 3. Error deleting detail <product>: <error>
func (s *SaleOrder) deleteSalesOrder(userId int32) OkAndErrorCodeReturn {
	if s.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	saleOrderInMemory := getSalesOrderRow(s.Id)
	if saleOrderInMemory.Id <= 0 || saleOrderInMemory.EnterpriseId != s.EnterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	d := getSalesOrderDetail(s.Id, s.EnterpriseId)

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
		d[i].EnterpriseId = s.EnterpriseId
		ok := d[i].deleteSalesOrderDetail(userId, trans)
		if !ok.Ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{strconv.Itoa(int(ok.ErrorCode)), d[i].Product.Name}}
		}
	}

	// delete sales order detail packaged using dbOrm

	for i := 0; i < len(d); i++ {
		result := trans.Delete(&SalesOrderDetailPackaged{}, "order_detail = ? AND enterprise = ?", d[i].Id, s.EnterpriseId)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	// delete packaging using dbOrm
	result := trans.Delete(&Packaging{}, "sales_order = ? AND enterprise = ?", s.Id, s.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	// delete pallets using dbOrm
	result = trans.Delete(&Pallet{}, "sales_order = ? AND enterprise = ?", s.Id, s.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	// delete discounts
	discounts := getSalesOrderDiscounts(s.Id, s.EnterpriseId)
	for i := 0; i < len(discounts); i++ {
		ok := discounts[i].deleteSalesOrderDiscount(userId)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(s.EnterpriseId, "sales_order", int(s.Id), userId, "D")
	inMemoryOrder := getSalesOrderRow(s.Id)
	json, _ := json.Marshal(inMemoryOrder)
	go fireWebHook(s.EnterpriseId, "sales_order", "DELETE", string(json))

	// delete sale order
	result = trans.Model(&s).Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).Delete(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

// Adds a total amount to the order total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesOrder(enterpriseId int32, orderId int64, userId int32, totalAmount float64, vatPercent float64, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.TotalProducts += totalAmount
	saleOrder.VatAmount += (totalAmount / 100) * vatPercent

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsSaleOrder(enterpriseId, orderId, userId, trans)
}

// Adds the discounts to the fix discount of the order. This function will substract if the amount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addDiscountsSalesOrder(enterpriseId int32, orderId int64, userId int32, amountTaxExcluded float64, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.FixDiscount += amountTaxExcluded

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsSaleOrder(enterpriseId, orderId, userId, trans)
}

// If the payment accepted date is null, sets it to the current date and time.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setDatePaymentAcceptedSalesOrder(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	if saleOrder.DatePaymentAccepted != nil {
		now := time.Now()
		saleOrder.DatePaymentAccepted = &now

		// save sale order to the database using dbOrm
		result = trans.Save(&saleOrder)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
		json, _ := json.Marshal(saleOrder)
		go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))
	}

	return true
}

// Applies the logic to calculate the totals of the sales order and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleOrder(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.TotalWithDiscount = (saleOrder.TotalProducts - saleOrder.TotalProducts*(saleOrder.DiscountPercent/100)) - saleOrder.FixDiscount + saleOrder.ShippingPrice - saleOrder.ShippingDiscount
	saleOrder.TotalAmount = saleOrder.TotalWithDiscount + saleOrder.VatAmount

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

type SalesOrderRelations struct {
	Invoices                   []SalesInvoice              `json:"invoices"`
	ManufacturingOrders        []ManufacturingOrder        `json:"manufacturingOrders"`
	ComplexManufacturingOrders []ComplexManufacturingOrder `json:"complexManufacturingOrders"`
	DeliveryNotes              []SalesDeliveryNote         `json:"deliveryNotes"`
	Shippings                  []Shipping                  `json:"shippings"`
}

func getSalesOrderRelations(orderId int64, enterpriseId int32) SalesOrderRelations {
	return SalesOrderRelations{
		Invoices:                   getSalesOrderInvoices(orderId, enterpriseId),
		ManufacturingOrders:        getSalesOrderManufacturingOrders(orderId, enterpriseId),
		DeliveryNotes:              getSalesOrderDeliveryNotes(orderId, enterpriseId),
		Shippings:                  getSalesOrderShippings(orderId, enterpriseId),
		ComplexManufacturingOrders: getSalesOrderComplexManufacturingOrders(orderId, enterpriseId),
	}
}

func getSalesOrderInvoices(orderId int64, enterpriseId int32) []SalesInvoice {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId {
		return []SalesInvoice{}
	}
	var invoices []SalesInvoice = make([]SalesInvoice, 0)
	salesOrderDetails := getSalesOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(salesOrderDetails); i++ {
		var invoiceDetails []SalesInvoiceDetail
		dbOrm.Model(&SalesInvoiceDetail{}).Where("order_detail = ?", salesOrderDetails[i].Id).Find(&invoiceDetails)
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
				invoice := getSalesInvoiceRow(invoiceDetails[j].InvoiceId)
				invoices = append(invoices, invoice)
			}
		}
	}
	return invoices
}

func getSalesOrderManufacturingOrders(orderId int64, enterpriseId int32) []ManufacturingOrder {
	// MANUFACTURING ORDER
	var manufacturingOrders []ManufacturingOrder = make([]ManufacturingOrder, 0)
	// get manufacturing orders for this order using dbOrm
	result := dbOrm.Model(&ManufacturingOrder{}).Where("\"order\" = ? AND enterprise = ?", orderId, enterpriseId).Order("date_created,id DESC").Preload(clause.Associations).Find(&manufacturingOrders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return manufacturingOrders
	}
	return manufacturingOrders
}

func getSalesOrderComplexManufacturingOrders(orderId int64, enterpriseId int32) []ComplexManufacturingOrder {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId {
		return []ComplexManufacturingOrder{}
	}
	var complexManufacturingOrders []ComplexManufacturingOrder = make([]ComplexManufacturingOrder, 0)
	salesOrderDetails := getSalesOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(salesOrderDetails); i++ {
		var complexManufacturingOrderDetails []ComplexManufacturingOrderManufacturingOrder
		dbOrm.Model(&ComplexManufacturingOrderManufacturingOrder{}).Where("sales_order_detail = ?", salesOrderDetails[i].Id).Find(&complexManufacturingOrderDetails)
		for j := 0; j < len(complexManufacturingOrderDetails); j++ {
			// only append if it doesn't already exist in the array searching by id
			var ok bool = true
			for k := 0; k < len(complexManufacturingOrders); k++ {
				if complexManufacturingOrders[k].Id == complexManufacturingOrderDetails[j].ComplexManufacturingOrderId {
					ok = false
					break
				}
			}
			if ok {
				complexManufacturingOrder := getComplexManufacturingOrderRow(complexManufacturingOrderDetails[j].ComplexManufacturingOrderId)
				complexManufacturingOrders = append(complexManufacturingOrders, complexManufacturingOrder)
			}
		}
	}
	return complexManufacturingOrders
}

func getSalesOrderDeliveryNotes(orderId int64, enterpriseId int32) []SalesDeliveryNote {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId {
		return []SalesDeliveryNote{}
	}
	var notes []SalesDeliveryNote = make([]SalesDeliveryNote, 0)
	salesOrderDetails := getSalesOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(salesOrderDetails); i++ {
		var noteDetails []WarehouseMovement
		dbOrm.Model(&WarehouseMovement{}).Where("sales_order_detail = ?", salesOrderDetails[i].Id).Find(&noteDetails)
		for j := 0; j < len(noteDetails); j++ {
			// only append note to notes if it doesn't already exist in the array searching by id
			var ok bool = true
			for k := 0; k < len(notes); k++ {
				if notes[k].Id == *noteDetails[j].SalesDeliveryNoteId {
					ok = false
					break
				}
			}
			if ok {
				note := getSalesDeliveryNoteRow(*noteDetails[j].SalesDeliveryNoteId)
				notes = append(notes, note)
			}
		}
	}
	return notes
}

func getSalesOrderShippings(orderId int64, enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	// get the shippings for this order and enterprise id using dbOrm
	result := dbOrm.Model(&Shipping{}).Where("\"order\" = ? AND enterprise = ?", orderId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&shippings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return shippings
	}

	return shippings
}

// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func setSalesOrderState(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	var status string
	result := trans.Model(&SalesOrderDetail{}).Where("\"order\"", orderId).Order("status ASC").Limit(1).Pluck("status", &status)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Model(&SaleOrder{}).Where("id = ?", orderId).Update("status", status)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")

	return true
}

type SaleOrderLocateReturn struct {
	Orders []SaleOrderLocate `json:"orders"`
	Rows   int64             `json:"rows"`
}

type SaleOrderLocate struct {
	Id          int64     `json:"id"`
	CustomerId  int32     `json:"customerId" gorm:"type:integer;not null:true;column:customer"`
	Customer    Customer  `json:"customer"`
	OrderName   string    `json:"orderName"`
	DateCreated time.Time `json:"dateCreated"`
}

type SaleOrderLocateQuery struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

func (query *SaleOrderLocateQuery) locateSaleOrder(enterpriseId int32) SaleOrderLocateReturn {
	res := SaleOrderLocateReturn{}
	res.Orders = make([]SaleOrderLocate, 0)
	// get all sale orders for the current enterprise sorted by date created desc using dbOrm
	result := dbOrm.Model(&SaleOrder{}).Where("enterprise = ?", enterpriseId).Order("date_created DESC").Offset(int(query.Offset)).Limit(int(query.Limit)).Find(&res.Orders)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return res
	}

	for i := 0; i < len(res.Orders); i++ {
		// get the customer for this sale order using dbOrm
		res.Orders[i].Customer = getCustomerRow(res.Orders[i].CustomerId)
	}

	result = dbOrm.Model(&SaleOrder{}).Where("enterprise = ?", enterpriseId).Count(&res.Rows)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return res
}

// Add an amount to the lines_number field in the sale order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderLinesNumber(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.LinesNumber += 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

// Takes out an amount to the lines_number field in the sale order. This number represents the total of lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderLinesNumber(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.LinesNumber -= 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

// Add an amount to the invoiced_lines field in the sale order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderInvoicedLines(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.InvoicedLines += 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

// Takes out an amount to the invoiced_lines field in the sale order. This number represents the total of invoiced lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderInvoicedLines(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.InvoicedLines -= 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

// Add an amount to the delivery_note_lines field in the sale order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addSalesOrderDeliveryNoteLines(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.DeliveryNoteLines += 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}

// Takes out an amount to the delivery_note_lines field in the sale order. This number represents the total of delivery note lines.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func removeSalesOrderDeliveryNoteLines(enterpriseId int32, orderId int64, userId int32, trans gorm.DB) bool {
	// get a single sales order row from the database using dbOrm
	var saleOrder SaleOrder
	result := trans.Model(&SaleOrder{}).Where("id = ?", orderId).First(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleOrder.DeliveryNoteLines -= 1

	// save sale order to the database using dbOrm
	result = trans.Save(&saleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_order", int(orderId), userId, "U")
	json, _ := json.Marshal(saleOrder)
	go fireWebHook(enterpriseId, "sales_order", "PUT", string(json))

	return true
}
