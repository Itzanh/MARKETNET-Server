/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"gorm.io/gorm"
)

type Pallet struct {
	Id           int32     `json:"id" gorm:"index:pallet_id_enterprise,unique:true,priority:1"`
	SalesOrderId int64     `json:"salesOrderId" gorm:"column:sales_order;not null:true"`
	SalesOrder   SaleOrder `json:"salesOrder" gorm:"foreignKey:SalesOrderId,EnterpriseId;references:Id,EnterpriseId"`
	Weight       float64   `json:"weight" gorm:"column:weight;not null:true;type:numeric(14,6)"`
	Width        float64   `json:"width" gorm:"column:width;not null:true;type:numeric(14,6)"`
	Height       float64   `json:"height" gorm:"column:height;not null:true;type:numeric(14,6)"`
	Depth        float64   `json:"depth" gorm:"column:depth;not null:true;type:numeric(14,6)"`
	Name         string    `json:"name" gorm:"column:name;not null:true;type:character varying(40)"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise;not null:true;index:pallet_id_enterprise,unique:true,priority:2"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (p *Pallet) TableName() string {
	return "pallets"
}

type Pallets struct {
	HasPallets bool     `json:"hasPallets"`
	Pallets    []Pallet `json:"pallets"`
}

func getSalesOrderPallets(orderId int64, enterpriseId int32) Pallets {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.Id <= 0 {
		return Pallets{}
	}
	if saleOrder.Carrier == nil || !saleOrder.Carrier.Pallets {
		return Pallets{HasPallets: false}
	}

	var pallets []Pallet = make([]Pallet, 0)
	result := dbOrm.Model(&Pallet{}).Where("sales_order = ? AND enterprise = ?", orderId, enterpriseId).Find(&pallets)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return Pallets{}
	}
	return Pallets{HasPallets: true, Pallets: pallets}
}

func getPalletsRow(palletId int32) Pallet {
	p := Pallet{}
	result := dbOrm.Model(&p).Where("id = ?", palletId).First(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return Pallet{}
	}
	return p
}

func (p *Pallet) isValid() bool {
	return !(p.Weight <= 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 || len(p.Name) == 0 || len(p.Name) > 40)
}

func (p *Pallet) BeforeCreate(tx *gorm.DB) (err error) {
	var pallet Pallet
	tx.Model(&Pallet{}).Last(&pallet)
	p.Id = pallet.Id + 1
	return nil
}

func (p *Pallet) insertPallet() bool {
	if p.SalesOrderId <= 0 || len(p.Name) == 0 || len(p.Name) > 40 {
		return false
	}

	s := getSettingsRecordById(p.EnterpriseId)
	p.Weight = s.PalletWeight
	p.Width = s.PalletWidth
	p.Height = s.PalletHeight
	p.Depth = s.PalletDepth

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *Pallet) updatePallet() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	var pallet Pallet
	result := dbOrm.Model(&p).Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).First(&pallet)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	pallet.Name = p.Name
	pallet.Weight = p.Weight
	pallet.Width = p.Width
	pallet.Height = p.Height
	pallet.Depth = p.Depth

	result = dbOrm.Updates(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *Pallet) deletePallet() bool {
	if p.Id <= 0 {
		return false
	}

	result := dbOrm.Delete(&Pallet{}, "id = ? AND enterprise = ?", p.Id, p.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
