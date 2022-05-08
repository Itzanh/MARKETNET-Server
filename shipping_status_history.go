package main

import "time"

type ShippingStatusHistory struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	ShippingId  int64     `json:"shippingId" gorm:"primaryKey;column:shipping;not null:true"`
	Shipping    Shipping  `json:"shipping" gorm:"foreignKey:ShippingId;references:Id"`
	StatusId    int16     `json:"statusId" gorm:"column:status_id;not null:true"`
	Message     string    `json:"message" gorm:"column:message;not null:true;type:text"`
	Delivered   bool      `json:"delivered" gorm:"column:delivered;not null:true"`
	DateCreated time.Time `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
}

func (ShippingStatusHistory) TableName() string {
	return "shipping_status_history"
}

func getShippingStatusHistory(enterpriseId int32, shippingId int64) []ShippingStatusHistory {
	var shipping Shipping
	result := dbOrm.Model(&Shipping{}).Where("id = ?", shippingId).Find(&shipping)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	if shipping.EnterpriseId != enterpriseId {
		return nil
	}

	var shippingStatusHistory []ShippingStatusHistory = make([]ShippingStatusHistory, 0)
	result = dbOrm.Model(&ShippingStatusHistory{}).Where("shipping = ?", shippingId).Order("date_created DESC").Find(&shippingStatusHistory)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return shippingStatusHistory
}
