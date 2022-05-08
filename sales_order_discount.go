package main

import "gorm.io/gorm"

type SalesOrderDiscount struct {
	Id                   int32               `json:"id"`
	OrderId              int64               `json:"orderId" gorm:"column:order;not null:true"`
	Order                SaleOrder           `json:"-" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	Name                 string              `json:"name" gorm:"column:name;not null:true;type:character varying(100)"`
	ValueTaxIncluded     float64             `json:"valueTaxIncluded" gorm:"column:value_tax_included;not null:true;type:numeric(14,6)"`
	ValueTaxExcluded     float64             `json:"valueTaxExcluded" gorm:"column:value_tax_excluded;not null:true;type:numeric(14,6)"`
	EnterpriseId         int32               `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise           Settings            `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	SalesInvoiceDetailId *int32              `json:"salesInvoiceDetailId" gorm:"column:sales_invoice_detail;type:integer"`
	SalesInvoiceDetail   *SalesInvoiceDetail `json:"-" gorm:"foreignKey:SalesInvoiceDetailId,EnterpriseId;references:Id,EnterpriseId"`
}

func (s *SalesOrderDiscount) TableName() string {
	return "sales_order_discount"
}

func getSalesOrderDiscounts(orderId int64, enterpriseId int32) []SalesOrderDiscount {
	var discounts []SalesOrderDiscount = make([]SalesOrderDiscount, 0)
	result := dbOrm.Model(&SalesOrderDiscount{}).Where("\"order\" = ? AND enterprise = ?", orderId, enterpriseId).Order("id ASC").Find(&discounts)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return discounts
	}
	return discounts
}

func getSalesOrderDiscountsRow(discountId int32) SalesOrderDiscount {
	d := SalesOrderDiscount{}
	result := dbOrm.Model(&SalesOrderDiscount{}).Where("id = ?", discountId).First(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return d
	}
	return d
}

func (d *SalesOrderDiscount) isValid() bool {
	return !(d.OrderId <= 0 || len(d.Name) == 0 || len(d.Name) > 100 || d.ValueTaxIncluded <= 0 || d.ValueTaxExcluded <= 0)
}

func (d *SalesOrderDiscount) BeforeCreate(tx *gorm.DB) (err error) {
	var salesOrderDiscount SalesOrderDiscount
	tx.Model(&SalesOrderDiscount{}).Last(&salesOrderDiscount)
	d.Id = salesOrderDiscount.Id + 1
	return nil
}

func (d *SalesOrderDiscount) insertSalesOrderDiscount(userId int32) bool {
	if !d.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	d.SalesInvoiceDetailId = nil

	result := trans.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	ok := addDiscountsSalesOrder(d.EnterpriseId, d.OrderId, userId, d.ValueTaxExcluded, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func (d *SalesOrderDiscount) deleteSalesOrderDiscount(userId int32) bool {
	if d.Id <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	inMemoryDiscount := getSalesOrderDiscountsRow(d.Id)
	if inMemoryDiscount.Id <= 0 || inMemoryDiscount.EnterpriseId != d.EnterpriseId {
		trans.Rollback()
		return false
	}

	result := trans.Delete(&SalesOrderDiscount{}, "id = ? AND enterprise = ?", d.Id, d.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	ok := addDiscountsSalesOrder(d.EnterpriseId, inMemoryDiscount.OrderId, userId, -inMemoryDiscount.ValueTaxExcluded, *trans)
	if !ok {
		trans.Rollback()
		return false
	}

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

func invoiceSalesOrderDiscounts(orderId int64, invoiceId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	rows, err := dbOrm.Model(&SalesOrderDiscount{}).Where("\"order\" = ? AND enterprise = ? AND sales_invoice_detail IS NULL", orderId, enterpriseId).Order("id ASC").Select("id,name,value_tax_excluded").Rows()
	if err != nil {
		log("DB", err.Error())
		return false
	}

	for rows.Next() {
		var id int32
		var name string
		var valueTaxExcluded float64
		rows.Scan(&id, &name, &valueTaxExcluded)

		invoiceDetal := SalesInvoiceDetail{}
		invoiceDetal.InvoiceId = invoiceId
		invoiceDetal.Description = name
		invoiceDetal.Price = -valueTaxExcluded
		invoiceDetal.Quantity = 1
		invoiceDetal.TotalAmount = -valueTaxExcluded
		invoiceDetal.VatPercent = 0
		invoiceDetal.EnterpriseId = enterpriseId
		ok := invoiceDetal.insertSalesInvoiceDetail(&trans, userId)
		if !ok.Ok {
			trans.Rollback()
			return false
		}

		result := trans.Where("id = ?", id).Update("sales_invoice_detail", invoiceDetal.Id)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	return true
}
