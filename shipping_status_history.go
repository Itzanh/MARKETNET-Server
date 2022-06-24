/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	"gorm.io/gorm"
)

type ShippingStatusHistory struct {
	Id           int64     `json:"id" gorm:"primaryKey"`
	ShippingId   int64     `json:"shippingId" gorm:"column:shipping;not null:true"`
	Shipping     Shipping  `json:"-" gorm:"foreignKey:ShippingId;references:Id"`
	StatusId     int16     `json:"statusId" gorm:"column:status_id;not null:true"`
	Message      string    `json:"message" gorm:"column:message;not null:true;type:text"`
	Delivered    bool      `json:"delivered" gorm:"column:delivered;not null:true"`
	DateCreated  time.Time `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (ShippingStatusHistory) TableName() string {
	return "shipping_status_history"
}

func (ssh *ShippingStatusHistory) BeforeCreate(tx *gorm.DB) (err error) {
	var shippingStatusHistory ShippingStatusHistory
	tx.Model(&ShippingStatusHistory{}).Last(&shippingStatusHistory)
	ssh.Id = shippingStatusHistory.Id + 1
	return nil
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
