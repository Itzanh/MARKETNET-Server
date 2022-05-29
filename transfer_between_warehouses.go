package main

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransferBetweenWarehouses struct {
	Id                     int64      `json:"id" gorm:"index:transfer_between_warehouses_id_enterprise,unique:true,priority:1"`
	WarehouseOriginId      string     `json:"warehouseOriginId" gorm:"column:warehouse_origin;type:character(2);not null"`
	WarehouseOrigin        Warehouse  `json:"warehouseOrigin" gorm:"foreignKey:WarehouseOriginId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseDestinationId string     `json:"warehouseDestinationId" gorm:"column:warehouse_destination;type:character(2);not null"`
	WarehouseDestination   Warehouse  `json:"warehouseDestination" gorm:"foreignKey:WarehouseDestinationId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId           int32      `json:"enterprise" gorm:"column:enterprise;type:integer;not null;index:transfer_between_warehouses_id_enterprise,unique:true,priority:2;index:transfer_between_warehouses_enterprise_finished_date_created,priority:1"`
	Enterprise             Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	DateCreated            time.Time  `json:"dateCreated" gorm:"column:date_created;type:timestamp(3) with time zone;not null;index:transfer_between_warehouses_enterprise_finished_date_created,priority:3"`
	DateFinished           *time.Time `json:"dateFinished" gorm:"column:date_finished;type:timestamp(3) with time zone"`
	Finished               bool       `json:"finished" gorm:"column:finished;type:boolean;not null;index:transfer_between_warehouses_enterprise_finished_date_created,priority:2"`
	LinesTransfered        int32      `json:"linesTransfered" gorm:"column:lines_transfered;type:integer;not null"`
	LinesTotal             int32      `json:"linesTotal" gorm:"column:lines_total;type:integer;not null"`
	Name                   string     `json:"name" gorm:"column:name;type:character varying(100);not null;index:transfer_between_warehouses_name,type:gin"`
}

func (t *TransferBetweenWarehouses) TableName() string {
	return "transfer_between_warehouses"
}

type TransferBetweenWarehousesQuery struct {
	Search     string     `json:"search"`
	DateStart  *time.Time `json:"dateStart"`
	DateEnd    *time.Time `json:"dateEnd"`
	Finished   bool       `json:"finished"`
	enterprise int32
}

func (q *TransferBetweenWarehousesQuery) searchTransferBetweenWarehouses() []TransferBetweenWarehouses {
	var transfers []TransferBetweenWarehouses = make([]TransferBetweenWarehouses, 0)

	cursor := dbOrm.Model(&TransferBetweenWarehouses{}).Where("transfer_between_warehouses.enterprise = ?", q.enterprise)
	if q.DateStart != nil {
		cursor.Where("transfer_between_warehouses.date_created >= ?", q.DateStart)
	}
	if q.DateEnd != nil {
		cursor.Where("transfer_between_warehouses.date_created <= ?", q.DateEnd)
	}
	cursor.Where("transfer_between_warehouses.finished = ?", q.Finished)
	if len(q.Search) > 0 {
		cursor.Where("transfer_between_warehouses.name LIKE ?", "%"+q.Search+"%")
	}

	result := cursor.Order("transfer_between_warehouses.date_created DESC").Preload(clause.Associations).Find(&transfers)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return transfers
}

func getTransferBetweenWarehousesRow(transferBetweenWarehousesId int64) TransferBetweenWarehouses {
	t := TransferBetweenWarehouses{}
	result := dbOrm.Where("id = ?", transferBetweenWarehousesId).First(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return t
}

func (t *TransferBetweenWarehouses) isValid() bool {
	return !(len(t.WarehouseOriginId) != 2 || len(t.WarehouseDestinationId) != 2 || t.WarehouseOriginId == t.WarehouseDestinationId || len(t.Name) == 0 || len(t.Name) > 100)
}

func (t *TransferBetweenWarehouses) BeforeCreate(tx *gorm.DB) (err error) {
	var transferBetweenWarehouses TransferBetweenWarehouses
	tx.Model(&TransferBetweenWarehouses{}).Last(&transferBetweenWarehouses)
	t.Id = transferBetweenWarehouses.Id + 1
	return nil
}

func (t *TransferBetweenWarehouses) insertTransferBetweenWarehouses() bool {
	if !t.isValid() {
		return false
	}

	t.DateCreated = time.Now()
	t.DateFinished = nil
	t.Finished = false
	t.LinesTransfered = 0
	t.LinesTotal = 0

	result := dbOrm.Create(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (t *TransferBetweenWarehouses) deleteTransferBetweenWarehouses() bool {
	if t.Id <= 0 {
		return false
	}

	details := getTransferBetweenWarehousesDetails(t.Id, t.EnterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].QuantityTransferred > 0 {
			return false
		}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	for i := 0; i < len(details); i++ {
		if !details[i].deleteTransferBetweenWarehousesDetail(trans) {
			trans.Rollback()
			return false
		}
	}

	result := trans.Delete(&TransferBetweenWarehouses{}, "id = ? AND enterprise = ?", t.Id, t.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	///
	trans.Commit()
	return true
	///
}

type TransferBetweenWarehousesDetail struct {
	Id                          int64                     `json:"id"`
	TransferBetweenWarehousesId int64                     `json:"transferBetweenWarehousesId" gorm:"column:transfer_between_warehouses;type:bigint;not null;index:transfer_between_warehouses_detail_barcode,priority:2,where:quantity_transferred < quantity"`
	TransferBetweenWarehouses   TransferBetweenWarehouses `json:"transferBetweenWarehouses" gorm:"foreignKey:TransferBetweenWarehousesId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId                int32                     `json:"-" gorm:"column:enterprise;type:integer;not null;index:transfer_between_warehouses_detail_barcode,priority:1,where:quantity_transferred < quantity"`
	Enterprise                  Settings                  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	ProductId                   int32                     `json:"productId" gorm:"column:product;type:integer;not null;index:transfer_between_warehouses_detail_barcode,priority:3,where:quantity_transferred < quantity"`
	Product                     Product                   `json:"product" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	Quantity                    int32                     `json:"quantity" gorm:"column:quantity;type:integer;not null"`
	QuantityTransferred         int32                     `json:"quantityTransferred" gorm:"column:quantity_transferred;type:integer;not null"`
	Finished                    bool                      `json:"finished" gorm:"column:finished;type:boolean;not null"`
	WarehouseMovementOutId      *int64                    `json:"warehouseMovementOutId" gorm:"column:warehouse_movement_out;type:bigint"`
	WarehouseMovementOut        *WarehouseMovement        `json:"warehouseMovementOut" gorm:"foreignKey:WarehouseMovementOutId,EnterpriseId;references:Id,EnterpriseId"`
	WarehouseMovementInId       *int64                    `json:"warehouseMovementInId" gorm:"column:warehouse_movement_in;type:bigint"`
	WarehouseMovementIn         *WarehouseMovement        `json:"warehouseMovementIn" gorm:"foreignKey:WarehouseMovementInId,EnterpriseId;references:Id,EnterpriseId"`
	SalesOrderDetailId          *int64                    `json:"salesOrderDetailId" gorm:"column:sales_order_detail;type:bigint"`
	SalesOrderDetail            *SalesOrderDetail         `json:"salesOrderDetail" gorm:"foreignKey:SalesOrderDetailId,EnterpriseId;references:Id,EnterpriseId"`
}

func (t *TransferBetweenWarehousesDetail) TableName() string {
	return "transfer_between_warehouses_detail"
}

func getTransferBetweenWarehousesDetails(transferBetweenWarehousesId int64, enterpriseId int32) []TransferBetweenWarehousesDetail {
	var details []TransferBetweenWarehousesDetail = make([]TransferBetweenWarehousesDetail, 0)
	result := dbOrm.Model(&TransferBetweenWarehousesDetail{}).Where("transfer_between_warehouses = ? AND enterprise = ?", transferBetweenWarehousesId, enterpriseId).Order("product ASC, id ASC").Preload(clause.Associations).Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return details
}

// For internal use only
func getTransferBetweenWarehousesDetailRow(transferBetweenWarehousesDetailId int64) TransferBetweenWarehousesDetail {
	d := TransferBetweenWarehousesDetail{}
	result := dbOrm.Where("id = ?", transferBetweenWarehousesDetailId).First(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return d
}

func (d *TransferBetweenWarehousesDetail) isValid() bool {
	return !(d.TransferBetweenWarehousesId <= 0 || d.ProductId <= 0 || d.Quantity <= 0)
}

func (d *TransferBetweenWarehousesDetail) BeforeCreate(tx *gorm.DB) (err error) {
	var transferBetweenWarehousesDetail TransferBetweenWarehousesDetail
	tx.Model(&TransferBetweenWarehousesDetail{}).Last(&transferBetweenWarehousesDetail)
	d.Id = transferBetweenWarehousesDetail.Id + 1
	return nil
}

func (d *TransferBetweenWarehousesDetail) insertTransferBetweenWarehousesDetail() bool {
	if !d.isValid() {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	transfer := getTransferBetweenWarehousesRow(d.TransferBetweenWarehousesId)
	if transfer.Id <= 0 || transfer.EnterpriseId != d.EnterpriseId || transfer.Finished {
		return false
	}

	d.QuantityTransferred = 0
	d.Finished = false
	d.WarehouseMovementInId = nil
	d.WarehouseMovementOutId = nil

	result := trans.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	transfer.LinesTotal += 1

	result = trans.Updates(&transfer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	return true
	///
}

func (d *TransferBetweenWarehousesDetail) deleteTransferBetweenWarehousesDetail(trans *gorm.DB) bool {
	if d.Id <= 0 {
		return false
	}

	detailInMemory := getTransferBetweenWarehousesDetailRow(d.Id)
	if detailInMemory.Id <= 0 || detailInMemory.EnterpriseId != d.EnterpriseId || detailInMemory.QuantityTransferred > 0 {
		return false
	}

	transfer := getTransferBetweenWarehousesRow(detailInMemory.TransferBetweenWarehousesId)
	if transfer.Id <= 0 || transfer.EnterpriseId != d.EnterpriseId || transfer.Finished {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false
		}
		///
	}

	result := trans.Model(&TransferBetweenWarehousesDetail{}).Where("id = ? AND enterprise = ?", d.Id, d.EnterpriseId).Delete(&TransferBetweenWarehousesDetail{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	transfer.LinesTotal -= 1

	result = trans.Updates(&transfer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	if beginTransaction {
		trans.Commit()
	}
	return true
	///
}

type TransferBetweenWarehousesDetailBarCodeQuery struct {
	TransferBetweenWarehousesId int64  `json:"transferBetweenWarehousesId"`
	BarCode                     string `json:"barCode"`
}

func (q *TransferBetweenWarehousesDetailBarCodeQuery) isValid() bool {
	return !(q.TransferBetweenWarehousesId <= 0 || len(q.BarCode) == 0 || len(q.BarCode) > 13)
}

func (detail *TransferBetweenWarehousesDetail) finishDetail(trans *gorm.DB, userId int32) bool {
	transfer := getTransferBetweenWarehousesRow(detail.TransferBetweenWarehousesId)

	// add 1 line transfered, set as finished
	transfer.LinesTransfered += 1
	transfer.Finished = transfer.LinesTransfered == transfer.LinesTotal
	if transfer.Finished {
		now := time.Now()
		transfer.DateFinished = &now
	} else {
		transfer.DateFinished = nil
	}

	// make an output warehouse movement from the origin warehouse
	wmOut := WarehouseMovement{
		WarehouseId:  transfer.WarehouseOriginId,
		ProductId:    detail.ProductId,
		Quantity:     detail.Quantity,
		Type:         "O",
		EnterpriseId: detail.EnterpriseId,
	}
	if !wmOut.insertWarehouseMovement(userId, trans) {
		trans.Rollback()
		return false
	}

	// make an input warehouse movement to the destination warehouse
	wmIn := WarehouseMovement{
		WarehouseId:  transfer.WarehouseDestinationId,
		ProductId:    detail.ProductId,
		Quantity:     detail.Quantity,
		Type:         "I",
		EnterpriseId: detail.EnterpriseId,
	}
	if !wmIn.insertWarehouseMovement(userId, trans) {
		trans.Rollback()
		return false
	}

	// save the transfer detail
	detail.WarehouseMovementOutId = &wmOut.Id
	detail.WarehouseMovementInId = &wmIn.Id

	result := trans.Updates(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// save the transfer
	result = trans.Model(&TransferBetweenWarehouses{}).Where("id = ?", transfer.Id).Updates(&transfer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	// if a sale order detail is attached, move the sale order detail to the destinarion warehouse
	if detail.SalesOrderDetailId != nil {
		result = dbOrm.Model(&SalesOrderDetail{}).Where("id = ?", detail.SalesOrderDetailId).Updates(map[string]interface{}{
			"warehouse": transfer.WarehouseDestinationId,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	return true
}

func (q *TransferBetweenWarehousesDetailBarCodeQuery) transferBetweenWarehousesDetailBarCode(enterpriseId int32, userId int32) bool {
	if !q.isValid() {
		return false
	}

	if len(q.BarCode) != 13 {
		q.BarCode = fmt.Sprintf("%013s", q.BarCode)
	}

	var transferBetweenWarehousesDetailId int64
	result := dbOrm.Model(&TransferBetweenWarehousesDetail{}).Where("enterprise = @enterpriseId AND transfer_between_warehouses = @transferBetweenWarehousesId AND quantity_transferred < quantity AND product = (SELECT id FROM product WHERE product.enterprise = @enterpriseId AND product.barCode = @barCode LIMIT 1)", sql.Named("enterpriseId", enterpriseId), sql.Named("transferBetweenWarehousesId", q.TransferBetweenWarehousesId), sql.Named("barCode", q.BarCode)).Order("id ASC").Limit(1).Pluck("id", &transferBetweenWarehousesDetailId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	detail := getTransferBetweenWarehousesDetailRow(transferBetweenWarehousesDetailId)
	if detail.Id <= 0 || detail.EnterpriseId != enterpriseId {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	detail.QuantityTransferred += 1
	detail.Finished = detail.QuantityTransferred == detail.Quantity

	if detail.Finished {
		ok := detail.finishDetail(trans, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	} else {
		result := trans.Model(&TransferBetweenWarehousesDetail{}).Where("id = ?", detail.Id).Update("quantity_transferred", detail.QuantityTransferred)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	return true
	///
}

type TransferBetweenWarehousesDetailQuantityQuery struct {
	TransferBetweenWarehousesDetailId int64 `json:"transferBetweenWarehousesDetailId"`
	Quantity                          int32 `json:"quantity"`
}

func (q *TransferBetweenWarehousesDetailQuantityQuery) isValid() bool {
	return !(q.TransferBetweenWarehousesDetailId <= 0 || q.Quantity <= 0)
}

func (q *TransferBetweenWarehousesDetailQuantityQuery) transferBetweenWarehousesDetailQuantity(enterpriseId int32, userId int32) bool {
	if !q.isValid() {
		return false
	}

	detail := getTransferBetweenWarehousesDetailRow(q.TransferBetweenWarehousesDetailId)
	if detail.Id <= 0 || detail.EnterpriseId != enterpriseId || detail.QuantityTransferred+q.Quantity > detail.Quantity {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	detail.QuantityTransferred += q.Quantity
	detail.Finished = detail.QuantityTransferred == detail.Quantity

	if detail.Finished {
		ok := detail.finishDetail(trans, userId)
		if !ok {
			trans.Rollback()
			return false
		}
	} else {
		result := trans.Model(&TransferBetweenWarehousesDetail{}).Where("id = ?", q.TransferBetweenWarehousesDetailId).Update("quantity_transferred", detail.QuantityTransferred)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}
	}

	///
	trans.Commit()
	return true
	///
}

func getTransferBetweenWarehousesWarehouseMovements(transferBetweenWarehousesId int64, enterpriseId int32) []WarehouseMovement {
	var movements []WarehouseMovement = make([]WarehouseMovement, 0)
	if transferBetweenWarehousesId <= 0 {
		return movements
	}

	transfer := getTransferBetweenWarehousesRow(transferBetweenWarehousesId)
	if transfer.Id <= 0 || transfer.EnterpriseId != enterpriseId {
		return movements
	}

	details := getTransferBetweenWarehousesDetails(transferBetweenWarehousesId, enterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].WarehouseMovementOutId != nil {
			m := getWarehouseMovementRow(*details[i].WarehouseMovementOutId)
			movements = append(movements, m)
		}
		if details[i].WarehouseMovementInId != nil {
			m := getWarehouseMovementRow(*details[i].WarehouseMovementInId)
			movements = append(movements, m)
		}
	}

	return movements
}

type TransferBetweenWarehousesToSentToPreparationOrders struct {
	WarehouseOriginId      *string `json:"warehouseOriginId"` // nil = All other warehouses
	WarehouseDestinationId string  `json:"warehouseDestinationId"`
	Name                   string  `json:"name"`
}

func (t *TransferBetweenWarehousesToSentToPreparationOrders) isValid() bool {
	return !((t.WarehouseOriginId != nil && len(*t.WarehouseOriginId) != 2) || len(t.WarehouseDestinationId) == 0 || (t.WarehouseOriginId != nil && *t.WarehouseOriginId == t.WarehouseDestinationId) || len(t.Name) == 0 || len(t.Name) > 100)
}

func (t *TransferBetweenWarehousesToSentToPreparationOrders) doTransfer(enterpriseId int32) bool {
	if !t.isValid() {
		return false
	}

	var details []SalesOrderDetail = make([]SalesOrderDetail, 0)
	cursor := dbOrm.Model(&SalesOrderDetail{}).Where(`(SELECT status FROM sales_order WHERE sales_order.id=sales_order_detail."order") = 'E'`)
	if t.WarehouseOriginId == nil {
		cursor = cursor.Where("warehouse != ?", t.WarehouseDestinationId)
	} else {
		cursor = cursor.Where("warehouse = ?", *t.WarehouseOriginId)
	}
	result := cursor.Order("id ASC").Find(&details)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	///
	trans := dbOrm.Begin()
	///

	var transfers map[string]TransferBetweenWarehouses = make(map[string]TransferBetweenWarehouses)
	if t.WarehouseOriginId != nil {
		transferBetweenWarehouses := TransferBetweenWarehouses{
			WarehouseOriginId:      *t.WarehouseOriginId,
			WarehouseDestinationId: t.WarehouseDestinationId,
			Name:                   t.Name,
			EnterpriseId:           enterpriseId,
		}

		result = trans.Create(&transferBetweenWarehouses)
		if result.Error != nil {
			trans.Rollback()
			log("DB", result.Error.Error())
			return false
		}
		transfers[*t.WarehouseOriginId] = transferBetweenWarehouses
	}

	for i := 0; i < len(details); i++ {
		var detail = details[i]

		_, ok := transfers[detail.WarehouseId]
		if !ok {
			transferBetweenWarehouses := TransferBetweenWarehouses{
				WarehouseOriginId:      detail.WarehouseId,
				WarehouseDestinationId: t.WarehouseDestinationId,
				Name:                   t.Name,
				EnterpriseId:           enterpriseId,
			}

			result = trans.Create(&transferBetweenWarehouses)
			if result.Error != nil {
				trans.Rollback()
				log("DB", result.Error.Error())
				return false
			}
			transfers[detail.WarehouseId] = transferBetweenWarehouses
		}

		var transferDetail = TransferBetweenWarehousesDetail{
			TransferBetweenWarehousesId: transfers[detail.WarehouseId].Id,
			ProductId:                   detail.ProductId,
			Quantity:                    detail.Quantity,
			EnterpriseId:                enterpriseId,
			SalesOrderDetailId:          &detail.Id,
		}
		result = trans.Create(&transferDetail)
		if result.Error != nil {
			trans.Rollback()
			log("DB", result.Error.Error())
			return false
		}
	}

	///
	trans.Commit()
	///
	return true
}
