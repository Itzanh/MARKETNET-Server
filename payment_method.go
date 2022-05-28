package main

import (
	"strings"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	Id                    int32    `json:"id" gorm:"primaryKey;index:payment_method_id_enterprise,unique:true,priority:1"`
	Name                  string   `json:"name" gorm:"type:character varying(100);not null:true"`
	PaidInAdvance         bool     `json:"paidInAdvance" gorm:"not null:true"`
	PrestashopModuleName  string   `json:"prestashopModuleName" gorm:"type:character varying(100);not null:true"`
	DaysExpiration        int16    `json:"daysExpiration" gorm:"not null:true"`
	Bank                  *int32   `json:"bank"`
	WooCommerceModuleName string   `json:"wooCommerceModuleName" gorm:"column:woocommerce_module_name;type:character varying(100);not null:true"`
	ShopifyModuleName     string   `json:"shopifyModuleName" gorm:"type:character varying(100);not null:true"`
	EnterpriseId          int32    `json:"-" gorm:"column:enterprise;not null:true;index:payment_method_id_enterprise,unique:true,priority:2"`
	Enterprise            Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *PaymentMethod) TableName() string {
	return "payment_method"
}

func getPaymentMethods(enterpriseId int32) []PaymentMethod {
	var paymentMethod []PaymentMethod = make([]PaymentMethod, 0)
	dbOrm.Model(&PaymentMethod{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&paymentMethod)
	return paymentMethod
}

func getPaymentMethodRow(paymentMethodId int32) PaymentMethod {
	p := PaymentMethod{}
	dbOrm.Model(&PaymentMethod{}).Where("id = ?", paymentMethodId).First(&p)
	return p
}

func (p *PaymentMethod) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 100 || p.DaysExpiration < 0 || len(p.PrestashopModuleName) > 100 || len(p.WooCommerceModuleName) > 100 || len(p.ShopifyModuleName) > 100)
}

func (p *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	var paymentMethod PaymentMethod
	tx.Model(&PaymentMethod{}).Last(&paymentMethod)
	p.Id = paymentMethod.Id + 1
	return nil
}

func (p *PaymentMethod) insertPaymentMethod() bool {
	if !p.isValid() {
		return false
	}

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *PaymentMethod) updatePaymentMethod() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	paymentMethod := PaymentMethod{}
	result := dbOrm.Model(&PaymentMethod{}).Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).First(&paymentMethod)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	paymentMethod.Name = p.Name
	paymentMethod.PaidInAdvance = p.PaidInAdvance
	paymentMethod.PrestashopModuleName = p.PrestashopModuleName
	paymentMethod.DaysExpiration = p.DaysExpiration
	paymentMethod.Bank = p.Bank
	paymentMethod.WooCommerceModuleName = p.WooCommerceModuleName
	paymentMethod.ShopifyModuleName = p.ShopifyModuleName

	result = dbOrm.Save(&paymentMethod)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *PaymentMethod) deletePaymentMethod() bool {
	if p.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).Delete(&PaymentMethod{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findPaymentMethodByName(paymentMethodName string, enterpriseId int32) []NameInt32 {
	var paymentMethod []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Country{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(paymentMethodName), enterpriseId).Limit(10).Find(&paymentMethod)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return paymentMethod
	}

	return paymentMethod
}

func getNamePaymentMethod(id int32, enterpriseId int32) string {
	var paymentMethod PaymentMethod
	result := dbOrm.Where("id = ? AND enterprise = ?", id, enterpriseId).First(&paymentMethod)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	return paymentMethod.Name
}

func locatePaymentMethods(enterpriseId int32) []NameInt32 {
	var paymentMethod []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&PaymentMethod{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&paymentMethod)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return paymentMethod
	}

	return paymentMethod
}
