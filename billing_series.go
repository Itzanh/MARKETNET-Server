/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"strings"
)

type BillingSerie struct {
	Id           string   `json:"id" gorm:"primaryKey;type:character(3)"`
	Name         string   `json:"name" gorm:"column:name;type:character varying(50);not null:true"`
	BillingType  string   `json:"billingType" gorm:"type:character(1);not null:true"`
	Year         int16    `json:"year" gorm:"not null:true"`
	EnterpriseId int32    `gorm:"primaryKey;column:enterprise"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (b *BillingSerie) TableName() string {
	return "billing_series"
}

func getBillingSeries(enterpriseId int32) []BillingSerie {
	var series []BillingSerie = make([]BillingSerie, 0)
	result := dbOrm.Model(&BillingSerie{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&series)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return series
	}

	return series
}

func (s *BillingSerie) isValid() bool {
	s.Id = strings.ToUpper(s.Id)
	return !(len(s.Id) != 3 || len(s.Name) == 0 || len(s.Name) > 50 || s.Year <= 0 || (s.BillingType != "S" && s.BillingType != "P"))
}

func (s *BillingSerie) insertBillingSerie() bool {
	if !s.isValid() {
		return false
	}

	result := dbOrm.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *BillingSerie) updateBillingSerie() bool {
	if s.Id == "" || !s.isValid() {
		return false
	}

	var serie BillingSerie
	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).First(&serie)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	serie.Name = s.Name
	serie.BillingType = s.BillingType
	serie.Year = s.Year

	result = dbOrm.Save(&serie)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *BillingSerie) deleteBillingSerie() bool {
	if s.Id == "" || len(s.Id) > 3 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, s.EnterpriseId).Delete(&BillingSerie{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func getNextSaleOrderNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&SaleOrder{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("order_number DESC").Limit(1).Select("order_number").Count(&rowsCount).Pluck("order_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func getNextSaleInvoiceNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&SalesInvoice{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("invoice_number DESC").Limit(1).Select("invoice_number").Count(&rowsCount).Pluck("invoice_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func getNextSaleDeliveryNoteNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&SalesDeliveryNote{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("delivery_note_number DESC").Limit(1).Select("delivery_note_number").Count(&rowsCount).Pluck("delivery_note_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func getNextPurchaseOrderNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&PurchaseOrder{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("order_number DESC").Limit(1).Select("order_number").Count(&rowsCount).Pluck("order_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func getNextPurchaseInvoiceNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&PurchaseInvoice{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("invoice_number DESC").Limit(1).Select("invoice_number").Count(&rowsCount).Pluck("invoice_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func getNextPurchaseDeliveryNoteNumber(billingSerieId string, enterpriseId int32) int32 {
	var orderNumber int32
	var rowsCount int64
	result := dbOrm.Model(&PurchaseDeliveryNote{}).Where("billing_series = ? AND enterprise = ?", billingSerieId, enterpriseId).Order("delivery_note_number DESC").Limit(1).Select("delivery_note_number").Count(&rowsCount).Pluck("delivery_note_number", &orderNumber)
	if rowsCount == 0 {
		return 1
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}
	return (orderNumber + 1)
}

func findBillingSerieByName(billingSerieName string, enterpriseId int32) []NameString {
	var billingSeries []NameString = make([]NameString, 0)
	result := dbOrm.Model(&Currency{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(billingSerieName), enterpriseId).Limit(10).Find(&billingSeries)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return billingSeries
	}

	return billingSeries
}

type LocateBillingSerie struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func locateBillingSeries(enterpriseId int32) []LocateBillingSerie {
	var series []LocateBillingSerie = make([]LocateBillingSerie, 0)
	result := dbOrm.Model(&BillingSerie{}).Where("enterprise = ?", enterpriseId).Find(&series)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return series
	}

	return series
}
