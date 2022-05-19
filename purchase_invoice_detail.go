package main

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseInvoiceDetail struct {
	Id            int64                `json:"id" gorm:"index:purchase_invoice_details_id_enterprise,unique:true,priority:1"`
	InvoiceId     int64                `json:"invoiceId" gorm:"column:invoice;not null:true"`
	Invoice       PurchaseInvoice      `json:"-" gorm:"foreignKey:InvoiceId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId     *int32               `json:"productId" gorm:"column:product"`
	Product       *Product             `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Price         float64              `json:"price" gorm:"column:price;not null:true;type:numeric(14,6)"`
	Quantity      int32                `json:"quantity" gorm:"column:quantity;not null:true"`
	VatPercent    float64              `json:"vatPercent" gorm:"column:vat_percent;not null:true;type:numeric(14,6)"`
	TotalAmount   float64              `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	OrderDetailId *int64               `json:"orderDetailId" gorm:"column:order_detail"`
	OrderDetail   *PurchaseOrderDetail `json:"orderDetail" gorm:"foreignKey:OrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId  int32                `json:"-" gorm:"column:enterprise;not null:true;index:purchase_invoice_details_id_enterprise,unique:true,priority:2"`
	Enterprise    Settings             `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Description   string               `json:"description" gorm:"column:description;not null:true;type:character varying(150)"`
	IncomeTax     bool                 `json:"incomeTax" gorm:"column:income_tax;not null:true"`
	Rent          bool                 `json:"rent" gorm:"column:rent;not null:true"`
}

func (*PurchaseInvoiceDetail) TableName() string {
	return "purchase_invoice_details"
}

func getPurchaseInvoiceDetail(invoiceId int64, enterpriseId int32) []PurchaseInvoiceDetail {
	var details []PurchaseInvoiceDetail = make([]PurchaseInvoiceDetail, 0)
	result := dbOrm.Model(&PurchaseInvoiceDetail{}).Where("invoice = ? AND enterprise = ?", invoiceId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return details
}

func getPurchaseInvoiceDetailRow(detailId int64) PurchaseInvoiceDetail {
	d := PurchaseInvoiceDetail{}
	result := dbOrm.Model(&PurchaseInvoiceDetail{}).Where("id = ?", detailId).Preload(clause.Associations).First(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return d
}

func (d *PurchaseInvoiceDetail) isValid() bool {
	return !(d.InvoiceId <= 0 || (d.ProductId == nil && len(d.Description) == 0) || len(d.Description) > 150 || (d.ProductId != nil && *d.ProductId <= 0) || d.Quantity <= 0 || d.VatPercent < 0)
}

func (c *PurchaseInvoiceDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var purchaseInvoiceDetail PurchaseInvoiceDetail
	tx.Model(&PurchaseInvoiceDetail{}).Last(&purchaseInvoiceDetail)
	c.Id = purchaseInvoiceDetail.Id + 1
	return nil
}

// ERROR CODES:
// 1. the product is deactivated
// 2. there is aleady a detail with this product
// 3. can't add details to a posted invoice
func (s *PurchaseInvoiceDetail) insertPurchaseInvoiceDetail(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	if s.ProductId != nil {
		p := getProductRow(*s.ProductId)
		if p.Id <= 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if p.Off {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
		}
	}

	s.TotalAmount = (s.Price * float64(s.Quantity)) * (1 + (s.VatPercent / 100))

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	var countProductInSaleOrder int64
	result := dbOrm.Model(&PurchaseInvoiceDetail{}).Where("invoice = ? AND product = ?", s.InvoiceId, s.ProductId).Count(&countProductInSaleOrder)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if countProductInSaleOrder > 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	// can't add details to a posted invoice
	invoice := getPurchaseInvoiceRowTransaction(s.InvoiceId, *trans)
	if invoice.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if invoice.AccountingMovementId != nil {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
	}

	result = trans.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(s.EnterpriseId, "purchase_invoice_details", int(s.Id), userId, "I")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "purchase_invoice_details", "POST", string(json))

	ok := addTotalProductsPurchaseInvoice(s.InvoiceId, s.Price*float64(s.Quantity), s.VatPercent, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if s.OrderDetailId != nil && *s.OrderDetailId != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*s.OrderDetailId, s.Quantity, s.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}
	if s.IncomeTax {
		ok = addIncomeTaxBasePurchaseInvoice(s.InvoiceId, s.Price*float64(s.Quantity), s.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}
	if s.Rent {
		ok = addRentBaseProductsPurchaseInvoice(s.InvoiceId, s.Price*float64(s.Quantity), s.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
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

// ERROR CODES
// 1. can't delete posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (d *PurchaseInvoiceDetail) deletePurchaseInvoiceDetail(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if d.Id <= 0 {
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

	detailInMemory := getPurchaseInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	i := getPurchaseInvoiceRow(detailInMemory.InvoiceId)
	if i.AccountingMovementId != nil { // can't delete posted invoices
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(d.EnterpriseId)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		i := getPurchaseInvoiceRow(d.InvoiceId)
		invoiceNumber := getNextPurchaseInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
		if invoiceNumber <= 0 || i.InvoiceNumber != (invoiceNumber-1) {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
		}
	}

	insertTransactionalLog(detailInMemory.EnterpriseId, "purchase_invoice_details", int(d.Id), userId, "D")
	json, _ := json.Marshal(d)
	go fireWebHook(d.EnterpriseId, "purchase_invoice_details", "DELETE", string(json))

	result := trans.Delete(&PurchaseInvoiceDetail{}, "id = ? AND enterprise = ?", d.Id, d.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsPurchaseInvoice(detailInMemory.InvoiceId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, d.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.OrderDetailId != nil && *detailInMemory.OrderDetailId != 0 {
		ok := addQuantityInvoicedPurchaseOrderDetail(*detailInMemory.OrderDetailId, -detailInMemory.Quantity, d.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}
	if detailInMemory.IncomeTax {
		ok = addIncomeTaxBasePurchaseInvoice(detailInMemory.InvoiceId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), d.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}
	if detailInMemory.Rent {
		ok = addRentBaseProductsPurchaseInvoice(detailInMemory.InvoiceId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), d.EnterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
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
