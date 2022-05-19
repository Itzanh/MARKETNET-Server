package main

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesInvoiceDetail struct {
	Id            int64             `json:"id" gorm:"index:sales_invoice_detail_id_enterprise,unique:true,priority:1"`
	InvoiceId     int64             `json:"invoiceId" gorm:"column:invoice;not null:true;index:sales_invoice_detail_invoice_product,unique:true,priority:1"`
	Invoice       SalesInvoice      `json:"-" gorm:"foreignKey:InvoiceId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId     *int32            `json:"productId" gorm:"column:product;index:sales_invoice_detail_invoice_product,unique:true,priority:2"`
	Product       *Product          `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Price         float64           `json:"price" gorm:"column:price;not null:true;type:numeric(14,6)"`
	Quantity      int32             `json:"quantity" gorm:"column:quantity;not null:true"`
	VatPercent    float64           `json:"vatPercent" gorm:"column:vat_percent;not null:true;type:numeric(14,6)"`
	TotalAmount   float64           `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	OrderDetailId *int64            `json:"orderDetailId" gorm:"column:order_detail"`
	OrderDetail   *SalesOrderDetail `json:"orderDetail" gorm:"foreignKey:OrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId  int32             `json:"-" gorm:"column:enterprise;not null:true;index:sales_invoice_detail_id_enterprise,unique:true,priority:2"`
	Enterprise    Settings          `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Description   string            `json:"description" gorm:"column:description;not null:true;type:character varying(150)"`
}

func (d *SalesInvoiceDetail) TableName() string {
	return "sales_invoice_detail"
}

func getSalesInvoiceDetail(invoiceId int64, enterpriseId int32) []SalesInvoiceDetail {
	var details []SalesInvoiceDetail = make([]SalesInvoiceDetail, 0)
	// get all sale invoice detail for the sale invoice and enterprise using dbOrm
	result := dbOrm.Model(&SalesInvoiceDetail{}).Where("invoice = ? AND enterprise = ?", invoiceId, enterpriseId).Order("id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return details
}

func getSalesInvoiceDetailRow(detailId int64) SalesInvoiceDetail {
	var detail SalesInvoiceDetail = SalesInvoiceDetail{}
	// get the sale invoice detail using dbOrm
	result := dbOrm.Model(&SalesInvoiceDetail{}).Where("id = ?", detailId).Preload(clause.Associations).First(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return detail
}

func (d *SalesInvoiceDetail) isValid() bool {
	return !(d.InvoiceId <= 0 || (d.ProductId == nil && len(d.Description) == 0) || len(d.Description) > 150 || d.Quantity <= 0 || d.VatPercent < 0)
}

func (s *SalesInvoiceDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var salesInvoiceDetail SalesInvoiceDetail
	tx.Model(&SalesInvoiceDetail{}).Last(&salesInvoiceDetail)
	s.Id = salesInvoiceDetail.Id + 1
	return nil
}

// ERROR CODES:
// 1. the product is deactivated
// 2. there is aleady a detail with this product
// 3. can't add details to a posted invoice
func (s *SalesInvoiceDetail) insertSalesInvoiceDetail(trans *gorm.DB, userId int32) OkAndErrorCodeReturn {
	if !s.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	if s.ProductId != nil && *s.ProductId <= 0 {
		s.ProductId = nil
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
	if beginTransaction {
		///
		trn := dbOrm.Begin()
		if trn.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		trans = trn
		///
	}

	// the product and sale order are unique, there can't exist another detail for the same product in the same order
	var countProductInSaleOrder int64
	result := dbOrm.Model(&SalesInvoiceDetail{}).Where("invoice = ? AND product = ?", s.InvoiceId, s.ProductId).Count(&countProductInSaleOrder)
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
	invoice := getSalesInvoiceRowTransaction(s.InvoiceId, *trans)
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

	ok := addTotalProductsSalesInvoice(s.InvoiceId, s.Price*float64(s.Quantity), s.VatPercent, s.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if s.OrderDetailId != nil && *s.OrderDetailId != 0 {
		ok := addQuantityInvociedSalesOrderDetail(*s.OrderDetailId, s.Quantity, userId, *trans)
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

	insertTransactionalLog(s.EnterpriseId, "sales_invoice_detail", int(s.Id), userId, "I")
	json, _ := json.Marshal(s)
	go fireWebHook(s.EnterpriseId, "sales_invoice_detail", "POST", string(json))

	return OkAndErrorCodeReturn{Ok: true}
}

// ERROR CODES:
// 1. can't delete posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (d *SalesInvoiceDetail) deleteSalesInvoiceDetail(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if d.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	var beginTransaction bool = (trans == nil)
	if beginTransaction {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	detailInMemory := getSalesInvoiceDetailRow(d.Id)
	if detailInMemory.Id <= 0 {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	i := getSalesInvoiceRow(detailInMemory.InvoiceId)
	if i.AccountingMovementId != nil { // can't delete posted invoices
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(d.EnterpriseId)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoiceNumber := getNextSaleInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
		if invoiceNumber <= 0 || i.InvoiceNumber != (invoiceNumber-1) {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
		}
	}

	result := trans.Model(&SalesOrderDiscount{}).Where("sales_invoice_detail = ? AND enterprise = ?", d.Id, d.EnterpriseId).Update("sales_invoice_detail", nil)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	insertTransactionalLog(d.EnterpriseId, "sales_invoice_detail", int(d.Id), userId, "D")
	json, _ := json.Marshal(d)
	go fireWebHook(d.EnterpriseId, "sales_invoice_detail", "DELETE", string(json))

	result = trans.Delete(&SalesInvoiceDetail{}, "id = ? AND enterprise = ?", d.Id, d.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	ok := addTotalProductsSalesInvoice(detailInMemory.InvoiceId, -(detailInMemory.Price * float64(detailInMemory.Quantity)), detailInMemory.VatPercent, d.EnterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	if detailInMemory.OrderDetailId != nil && *detailInMemory.OrderDetailId != 0 {
		detail := getSalesOrderDetailRow(*detailInMemory.OrderDetailId)
		if detail.Id <= 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
		// if the detail had a purchase order pending, rollback the quantity assigned
		if detail.PurchaseOrderDetailId != nil {
			ok = addQuantityAssignedSalePurchaseOrder(*detail.PurchaseOrderDetailId, -detail.Quantity, detailInMemory.EnterpriseId, userId, *trans)
			if !ok {
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}
			trans.Model(&SalesOrderDetail{}).Where("id = ?", detail.Id).Update("purchase_order_detail", nil)
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}
		}
		// revert back the status
		ok := addQuantityInvociedSalesOrderDetail(*detailInMemory.OrderDetailId, -detailInMemory.Quantity, userId, *trans)
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
