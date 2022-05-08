package main

import "gorm.io/gorm"

type SalesOrderDetailDigitalProductData struct {
	Id       int32            `json:"id"`
	DetailId int64            `json:"detailId" gorm:"column:detail;not null:true"`
	Detail   SalesOrderDetail `json:"detail" gorm:"foreignKey:DetailId;reference:Id"`
	Key      string           `json:"key" gorm:"column:key;not null:true;type:character varying(50)"`
	Value    string           `json:"value" gorm:"column:value;not null:true;type:character varying(250)"`
}

func (s *SalesOrderDetailDigitalProductData) TableName() string {
	return "sales_order_detail_digital_product_data"
}

func getSalesOrderDetailDigitalProductData(salesOrderDetailId int64, enterpriseId int32) []SalesOrderDetailDigitalProductData {
	productData := make([]SalesOrderDetailDigitalProductData, 0)

	detailRow := getSalesOrderDetailRow(salesOrderDetailId)
	if detailRow.EnterpriseId != enterpriseId {
		return productData
	}

	result := dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Where("detail = ?", salesOrderDetailId).Find(&productData)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return productData
}

func (d *SalesOrderDetailDigitalProductData) isValid() bool {
	return !(d.DetailId <= 0 || len(d.Key) == 0 || len(d.Value) == 0)
}

func (d *SalesOrderDetailDigitalProductData) BeforeCreate(tx *gorm.DB) (err error) {
	var salesOrderDetailDigitalProductData SalesOrderDetailDigitalProductData
	tx.Model(&SalesOrderDetailDigitalProductData{}).Last(&salesOrderDetailDigitalProductData)
	d.Id = salesOrderDetailDigitalProductData.Id + 1
	return nil
}

func (d *SalesOrderDetailDigitalProductData) insertSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if !d.isValid() {
		return false
	}

	detailRow := getSalesOrderDetailRow(d.DetailId)
	if detailRow.EnterpriseId != enterpriseId || detailRow.Status != "E" {
		return false
	}
	productRow := getProductRow(detailRow.ProductId)
	if !productRow.DigitalProduct {
		return false
	}

	result := dbOrm.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (d *SalesOrderDetailDigitalProductData) updateSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if !d.isValid() || d.Id <= 0 {
		return false
	}

	detailRow := getSalesOrderDetailRow(d.DetailId)
	if detailRow.EnterpriseId != enterpriseId || detailRow.Status != "E" {
		return false
	}

	result := dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Where("id = ?", d.Id).Updates(map[string]interface{}{
		"key":   d.Key,
		"value": d.Value,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (d *SalesOrderDetailDigitalProductData) deleteSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if d.Id <= 0 {
		return false
	}

	var digitalData SalesOrderDetailDigitalProductData
	result := dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Where("id = ?", d.Id).First(&digitalData)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	detailRow := getSalesOrderDetailRow(digitalData.DetailId)
	if detailRow.EnterpriseId != enterpriseId || detailRow.Status != "E" {
		return false
	}

	result = dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Where("id = ?", d.Id).Delete(&SalesOrderDetailDigitalProductData{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type SetDigitalSalesOrderDetailAsSent struct {
	Detail                 int64  `json:"detail"`
	SendEmail              bool   `json:"sendEmail"`
	DestinationAddress     string `json:"destinationAddress"`
	DestinationAddressName string `json:"destinationAddressName"`
	Subject                string `json:"subject"`
}

func (data *SetDigitalSalesOrderDetailAsSent) setDigitalSalesOrderDetailAsSent(enterpriseId int32, userId int32) bool {
	detail := getSalesOrderDetailRow(data.Detail)
	if detail.EnterpriseId != enterpriseId || detail.Status != "E" {
		return false
	}
	digitalProductData := getSalesOrderDetailDigitalProductData(detail.Id, enterpriseId)
	if len(digitalProductData) == 0 {
		return false
	}

	if data.SendEmail {
		ei := EmailInfo{
			DestinationAddress:     data.DestinationAddress,
			DestinationAddressName: data.DestinationAddressName,
			Subject:                data.Subject,
			ReportId:               "SALES_ORDER_DIGITAL_PRODUCT_DATA",
			ReportDataId:           int32(detail.OrderId),
		}
		ei.sendEmail(enterpriseId)
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	result := trans.Model(&SalesOrderDetail{}).Where("id = ?", data.Detail).Update("status", "G")
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	setSalesOrderState(enterpriseId, detail.OrderId, userId, *trans)

	///
	trans.Commit()
	///
	return true
}
