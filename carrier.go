package main

import (
	"strings"

	"gorm.io/gorm"
)

type Carrier struct {
	Id                      int32    `json:"id" gorm:"index:carrier_id_enterprise,unique:true,priority:1"`
	Name                    string   `json:"name" gorm:"type:character varying(50);not null:true"`
	MaxWeight               float64  `json:"maxWeight" gorm:"type:numeric(14,6);not null:true"`
	MaxWidth                float64  `json:"maxWidth" gorm:"type:numeric(14,6);not null:true"`
	MaxHeight               float64  `json:"maxHeight" gorm:"type:numeric(14,6);not null:true"`
	MaxDepth                float64  `json:"maxDepth" gorm:"type:numeric(14,6);not null:true"`
	MaxPackages             int16    `json:"maxPackages" gorm:"not null:true"`
	Phone                   string   `json:"phone" gorm:"type:character varying(15);not null:true"`
	Email                   string   `json:"email" gorm:"type:character varying(100);not null:true"`
	Web                     string   `json:"web" gorm:"type:character varying(100);not null:true"`
	Off                     bool     `json:"off" gorm:"not null:true"`
	PrestaShopId            int32    `json:"prestaShopId" gorm:"column:ps_id;not null:true;index:carrier_ps_id,unique:true,where:ps_id <> 0"`
	Pallets                 bool     `json:"pallets" gorm:"not null:true"`
	Webservice              string   `json:"webservice" gorm:"type:character(1);not null:true"`
	SendcloudUrl            string   `json:"sendcloudUrl" gorm:"type:character varying(75);not null:true"`
	SendcloudKey            string   `json:"sendcloudKey" gorm:"type:character varying(32);not null:true"`
	SendcloudSecret         string   `json:"sendcloudSecret" gorm:"type:character varying(32);not null:true"`
	SendcloudShippingMethod int32    `json:"sendcloudShippingMethod" gorm:"not null:true"`
	SendcloudSenderAddress  int64    `json:"sendcloudSenderAddress" gorm:"not null:true"`
	EnterpriseId            int32    `json:"-" gorm:"column:enterprise;not null:true;index:carrier_id_enterprise,unique:true,priority:2"`
	Enterprise              Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Carrier) TableName() string {
	return "carrier"
}

func getCariers(enterpriseId int32) []Carrier {
	var carriers []Carrier = make([]Carrier, 0)
	dbOrm.Model(&Carrier{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&carriers)
	return carriers
}

func getCarierRow(id int32) Carrier {
	c := Carrier{}
	dbOrm.Model(&Carrier{}).Where("id = ?", id).First(&c)
	return c
}

func (c *Carrier) BeforeCreate(tx *gorm.DB) (err error) {
	var carrier Carrier
	tx.Model(&Carrier{}).Last(&carrier)
	c.Id = carrier.Id + 1
	return nil
}

func (c *Carrier) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || c.MaxWeight < 0 || c.MaxWidth < 0 || c.MaxHeight < 0 || c.MaxDepth < 0 || c.MaxPackages < 0 || len(c.Phone) > 15 || len(c.Email) > 100 || len(c.Web) > 100 || len(c.Webservice) != 1 || (c.Webservice != "_" && c.Webservice != "S") || len(c.SendcloudUrl) > 75 || (len(c.SendcloudKey) != 0 && len(c.SendcloudKey) != 32) || (len(c.SendcloudSecret) != 0 && len(c.SendcloudSecret) != 32) || c.SendcloudShippingMethod < 0 || c.SendcloudSenderAddress < 0)
}

func (c *Carrier) insertCarrier() bool {
	if !c.isValid() {
		return false
	}

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Carrier) updateCarrier() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	var carrier Carrier
	result := dbOrm.Model(&Carrier{}).Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&carrier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	carrier.Name = c.Name
	carrier.MaxWeight = c.MaxWeight
	carrier.MaxWidth = c.MaxWidth
	carrier.MaxHeight = c.MaxHeight
	carrier.MaxDepth = c.MaxDepth
	carrier.MaxPackages = c.MaxPackages
	carrier.Phone = c.Phone
	carrier.Email = c.Email
	carrier.Web = c.Web
	carrier.Off = c.Off
	carrier.PrestaShopId = c.PrestaShopId
	carrier.Pallets = c.Pallets
	carrier.Webservice = c.Webservice
	carrier.SendcloudUrl = c.SendcloudUrl
	carrier.SendcloudSecret = c.SendcloudSecret
	carrier.SendcloudShippingMethod = c.SendcloudShippingMethod
	carrier.SendcloudSenderAddress = c.SendcloudSenderAddress

	result = dbOrm.Save(&carrier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Carrier) deleteCarrier() bool {
	if c.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Carrier{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func findCarrierByName(carrierName string, enterpriseId int32) []NameInt32 {
	var carriers []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Carrier{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(carrierName), enterpriseId).Limit(10).Find(&carriers)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return carriers
	}

	return carriers
}

func getNameCarrier(id int32, enterpriseId int32) string {
	var carrier Carrier
	result := dbOrm.Where("id = ? AND enterprise = ?", id, enterpriseId).First(&carrier)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	return carrier.Name
}

func locateCarriers(enterpriseId int32) []NameInt32 {
	var carriers []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Carrier{}).Where("enterprise = ?", enterpriseId).Find(&carriers)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return carriers
	}

	return carriers
}
