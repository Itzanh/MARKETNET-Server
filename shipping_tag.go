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

type ShippingTag struct {
	Id           int64     `json:"id"`
	ShippingId   int64     `json:"shippingId" gorm:"column:shipping;not null:true"`
	Shipping     Shipping  `json:"shipping" gorm:"foreignKey:ShippingId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated  time.Time `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	Label        []byte    `json:"label" gorm:"column:label;not null:true"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (t *ShippingTag) TableName() string {
	return "shipping_tag"
}

func getShippingTags(shippingId int64, enterpriseId int32) []ShippingTag {
	tags := make([]ShippingTag, 0)
	result := dbOrm.Model(&ShippingTag{}).Where("shipping = ? AND enterprise = ?", shippingId, enterpriseId).Order("date_created DESC").Find(&tags)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return tags
}

func (t *ShippingTag) BeforeCreate(tx *gorm.DB) (err error) {
	var shippingTag ShippingTag
	tx.Model(&ShippingTag{}).Last(&shippingTag)
	t.Id = shippingTag.Id + 1
	return nil
}

func (t *ShippingTag) insertShippingTag() bool {
	t.DateCreated = time.Now()
	result := dbOrm.Create(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func deleteAllShippingTags(enterpriseId int32) {
	result := dbOrm.Where("enterprise = ?", enterpriseId).Delete(&ShippingTag{})
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}
