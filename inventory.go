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
	"gorm.io/gorm/clause"
)

type Inventory struct {
	Id           int32      `json:"id" gorm:"index:inventory_id_enterprise,unique:true,priority:1"`
	EnterpriseId int32      `json:"-" gorm:"column:enterprise;not null:true;index:inventory_id_enterprise,unique:true,priority:2"`
	Enterprise   Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Name         string     `json:"name" gorm:"column:name;not null:true;type:character varying(50)"`
	DateCreated  time.Time  `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	Finished     bool       `json:"finished" gorm:"column:finished;not null:true"`
	DateFinished *time.Time `json:"dateFinished" gorm:"column:date_finished;type:timestamp(3) with time zone"`
	WarehouseId  string     `json:"warehouseId" gorm:"column:warehouse;not null:true;type:character(2)"`
	Warehouse    Warehouse  `json:"warehouse" gorm:"foreignKey:WarehouseId,EnterpriseId;references:Id,EnterpriseId"`
}

func (i *Inventory) TableName() string {
	return "inventory"
}

func getInventories(enterpriseId int32) []Inventory {
	var inventory []Inventory = make([]Inventory, 0)
	result := dbOrm.Model(&Inventory{}).Where("inventory.enterprise = ?", enterpriseId).Order("inventory.id DESC").Preload(clause.Associations).Find(&inventory)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return inventory
	}
	return inventory
}

func getInventoryRow(inventoryId int32) Inventory {
	i := Inventory{}
	result := dbOrm.Model(&Inventory{}).Where("inventory.id = ?", inventoryId).Preload(clause.Associations).First(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return i
	}
	return i
}

func (i *Inventory) isValid() bool {
	return !(len(i.Name) == 0 || len(i.Name) > 50 || len(i.WarehouseId) == 0 || len(i.WarehouseId) > 2)
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) (err error) {
	var inventory Inventory
	tx.Model(&Inventory{}).Last(&inventory)
	i.Id = inventory.Id + 1
	return nil
}

func (i *Inventory) insertInventory(enterpriseId int32) bool {
	if !i.isValid() {
		return false
	}

	i.EnterpriseId = enterpriseId
	i.DateCreated = time.Now()
	i.Finished = false
	i.DateFinished = nil

	result := dbOrm.Create(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (i *Inventory) deleteInventory(enterpriseId int32) OkAndErrorCodeReturn {
	if i.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	inventory := getInventoryRow(i.Id)
	if inventory.Id <= 0 || inventory.EnterpriseId != enterpriseId || inventory.Finished {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	details := getInventoryProducts(inventory.Id, enterpriseId)
	if len(details) > 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	result := dbOrm.Delete(&Inventory{}, "id = ? AND enterprise = ?", inventory.Id, enterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}
	return OkAndErrorCodeReturn{Ok: true}
}

func (i *Inventory) finishInventory(userId int32, enterpriseId int32) bool {
	inMemoyInventory := getInventoryRow(i.Id)
	if inMemoyInventory.Id <= 0 || inMemoyInventory.Finished || inMemoyInventory.EnterpriseId != enterpriseId {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	var ok bool
	products := getInventoryProducts(inMemoyInventory.Id, enterpriseId)
	for i := 0; i < len(products); i++ {
		p := products[i]

		wm := WarehouseMovement{
			EnterpriseId: enterpriseId,
			WarehouseId:  inMemoyInventory.WarehouseId,
			ProductId:    p.ProductId,
			Quantity:     p.Quantity,
			Type:         "R",
		}
		ok = wm.insertWarehouseMovement(userId, trans)
		if !ok {
			trans.Rollback()
			return false
		}

		result := trans.Model(&InventoryProducts{}).Where("inventory = ? AND product = ?", p.InventoryId, p.ProductId).Update("warehouse_movement", wm.Id)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	result := trans.Model(&Inventory{}).Where("id = ?", inMemoyInventory.Id).Updates(map[string]interface{}{
		"finished":      true,
		"date_finished": time.Now(),
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///
	return true
}

type InventoryProducts struct {
	InventoryId         int32              `json:"inventoryId" gorm:"primaryKey;column:inventory;not null:true"`
	Inventory           Inventory          `json:"inventory" gorm:"foreignKey:InventoryId,EnterpriseId;references:Id,EnterpriseId"`
	ProductId           int32              `json:"productId" gorm:"primaryKey;column:product;not null:true"`
	Product             Product            `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId        int32              `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise          Settings           `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Quantity            int32              `json:"quantity" gorm:"column:quantity;not null:true"`
	WarehouseMovementId *int64             `json:"warehouseMovementId" gorm:"column:warehouse_movement"`
	WarehouseMovement   *WarehouseMovement `json:"warehouseMovement" gorm:"foreignKey:WarehouseMovementId,EnterpriseId;references:Id,EnterpriseId"`
}

func (p *InventoryProducts) TableName() string {
	return "inventory_products"
}

func getInventoryProducts(inventoryId int32, enterpriseId int32) []InventoryProducts {
	var inventoryProducts []InventoryProducts = make([]InventoryProducts, 0)
	result := dbOrm.Model(&InventoryProducts{}).Where("inventory_products.inventory = ? AND inventory_products.enterprise = ?", inventoryId, enterpriseId).Order("inventory_products.product ASC").Preload(clause.Associations).Find(&inventoryProducts)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return inventoryProducts
	}
	return inventoryProducts
}

func getInventoryProductsRow(inventoryId int32, productId int32, enterpriseId int32) InventoryProducts {
	ip := InventoryProducts{}
	result := dbOrm.Model(&InventoryProducts{}).Where("inventory_products.inventory = ? AND inventory_products.product = ? AND inventory_products.enterprise = ?", inventoryId, productId, enterpriseId).Preload(clause.Associations).First(&ip)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ip
	}
	return ip
}

func (i *InventoryProducts) isValid() bool {
	return !(i.InventoryId <= 0 || i.ProductId <= 0)
}

type InputInventoryProducts struct {
	Inventory         int32               `json:"inventory"`
	InventoryProducts []InventoryProducts `json:"inventoryProducts"`
	FamilyId          int32               `json:"familyId"`
}

func (input *InputInventoryProducts) insertUpdateDeleteInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId || i.Finished {
		return false
	}

	for i := 0; i < len(input.InventoryProducts); i++ {
		ip := input.InventoryProducts[i]
		if !ip.isValid() {
			return false
		}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	// input data optimization
	arayExistentInventoryProducts := getInventoryProducts(input.Inventory, enterpriseId)
	var existentInventoryProducts map[int32]InventoryProducts = make(map[int32]InventoryProducts)
	for i := 0; i < len(arayExistentInventoryProducts); i++ {
		a := arayExistentInventoryProducts[i]
		existentInventoryProducts[a.ProductId] = a
	}

	// cross data
	var toInsert []InventoryProducts = make([]InventoryProducts, 0)
	var toUpdate []InventoryProducts = make([]InventoryProducts, 0)
	for i := 0; i < len(input.InventoryProducts); i++ {
		newIp := input.InventoryProducts[i]

		oldIp, ok := existentInventoryProducts[newIp.ProductId]
		delete(existentInventoryProducts, newIp.ProductId)
		if !ok {
			toInsert = append(toInsert, newIp)
		} else if oldIp.Quantity != newIp.Quantity {
			toUpdate = append(toUpdate, newIp)
		}
	}

	// insert data
	for i := 0; i < len(toInsert); i++ {
		pi := toInsert[i]
		pi.EnterpriseId = enterpriseId
		pi.WarehouseMovementId = nil
		result := trans.Create(&pi)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	// update data
	for i := 0; i < len(toUpdate); i++ {
		pi := toUpdate[i]
		result := trans.Model(&InventoryProducts{}).Where("inventory = ? AND product = ?", input.Inventory, pi.ProductId).Update("quantity", pi.Quantity)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	// delete the remaining data in the map
	for k := range existentInventoryProducts {
		result := trans.Delete(&InventoryProducts{}, "inventory = ? AND product = ?", input.Inventory, k)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) insertProductFamilyInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId || i.Finished {
		return false
	}

	var enterprise int32
	result := dbOrm.Model(&ProductFamily{}).Where("id = ?", input.FamilyId).Pluck("enterprise", &enterprise)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	rows, err := dbOrm.Model(&Product{}).Where("family = ? AND enterprise = ?", input.FamilyId, enterpriseId).Select("id").Order("id ASC").Rows()
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	var productId int32
	for rows.Next() {
		rows.Scan(&productId)

		var inventoryProduct InventoryProducts = InventoryProducts{
			InventoryId:         input.Inventory,
			ProductId:           productId,
			EnterpriseId:        enterpriseId,
			Quantity:            0,
			WarehouseMovementId: nil,
		}

		result := trans.Create(&inventoryProduct)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) insertAllProductsInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId || i.Finished {
		return false
	}

	products := getProduct(enterpriseId)

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	var productId int32
	for i := 0; i < len(products); i++ {
		if products[i].Off {
			continue
		}
		productId = products[i].Id

		var inventoryProduct InventoryProducts = InventoryProducts{
			InventoryId:         input.Inventory,
			ProductId:           productId,
			EnterpriseId:        enterpriseId,
			Quantity:            0,
			WarehouseMovementId: nil,
		}

		result := trans.Create(&inventoryProduct)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}

func (input *InputInventoryProducts) deleteAllProductsInventoryProducts(enterpriseId int32) bool {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId || i.Finished {
		return false
	}

	result := dbOrm.Model(&InventoryProducts{}).Where("inventory = ? AND enterprise = ?", input.Inventory, enterpriseId).Delete(InventoryProducts{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

type BarCodeInputInventoryProducts struct {
	Inventory int32  `json:"inventory"`
	BarCode   string `json:"barCode"`
}

type BarCodeInputInventoryProductsResult struct {
	Ok               bool   `json:"ok"`
	ProductReference string `json:"productReference"`
	ProductName      string `json:"productName"`
	Quantity         int32  `json:"quantity"`
}

func (input *BarCodeInputInventoryProducts) insertOrCountInventoryProductsByBarcode(enterpriseId int32) BarCodeInputInventoryProductsResult {
	i := getInventoryRow(input.Inventory)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId || i.Finished {
		return BarCodeInputInventoryProductsResult{}
	}

	product := getProductByBarcode(input.BarCode, enterpriseId)
	if product.Id <= 0 {
		return BarCodeInputInventoryProductsResult{}
	}

	var rowCount int64
	result := dbOrm.Model(&InventoryProducts{}).Where("inventory = ? AND product = ?", input.Inventory, product.Id).Count(&rowCount)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return BarCodeInputInventoryProductsResult{}
	}

	if rowCount == 0 {
		var inventoryProduct InventoryProducts = InventoryProducts{
			InventoryId:         input.Inventory,
			ProductId:           product.Id,
			EnterpriseId:        enterpriseId,
			Quantity:            1,
			WarehouseMovementId: nil,
		}

		result := dbOrm.Create(&inventoryProduct)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return BarCodeInputInventoryProductsResult{}
		}
	} else {
		var quantity int32
		result := dbOrm.Model(&InventoryProducts{}).Where("inventory = ? AND product = ?", input.Inventory, product.Id).Select("quantity").Pluck("quantity", &quantity)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return BarCodeInputInventoryProductsResult{}
		}

		quantity += 1

		result = dbOrm.Model(&InventoryProducts{}).Where("inventory = ? AND product = ?", input.Inventory, product.Id).Update("quantity", quantity)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return BarCodeInputInventoryProductsResult{}
		}
	}

	inventoryProduct := getInventoryProductsRow(input.Inventory, product.Id, enterpriseId)
	return BarCodeInputInventoryProductsResult{Ok: true, ProductReference: product.Reference, ProductName: product.Name, Quantity: inventoryProduct.Quantity}
}
