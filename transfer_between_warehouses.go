package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type TransferBetweenWarehouses struct {
	Id                       int64      `json:"id"`
	WarehouseOrigin          string     `json:"warehouseOrigin"`
	WarehouseDestination     string     `json:"warehouseDestination"`
	DateCreated              time.Time  `json:"dateCreated"`
	DateFinished             *time.Time `json:"dateFinished"`
	Finished                 bool       `json:"finished"`
	LinesTransfered          int32      `json:"linesTransfered"`
	LinesTotal               int32      `json:"linesTotal"`
	Name                     string     `json:"name"`
	WarehouseOriginName      string     `json:"warehouseOriginName"`
	WarehouseDestinationName string     `json:"warehouseDestinationName"`
	enterprise               int32
}

type TransferBetweenWarehousesQuery struct {
	DateStart  *time.Time `json:"dateStart"`
	DateEnd    *time.Time `json:"dateEnd"`
	Finished   bool       `json:"finished"`
	enterprise int32
}

func (q *TransferBetweenWarehousesQuery) searchTransferBetweenWarehouses() []TransferBetweenWarehouses {
	var transfers []TransferBetweenWarehouses = make([]TransferBetweenWarehouses, 0)

	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id = transfer_between_warehouses.warehouse_origin AND warehouse.enterprise = transfer_between_warehouses.enterprise),(SELECT name FROM warehouse WHERE warehouse.id = transfer_between_warehouses.warehouse_destination AND warehouse.enterprise = transfer_between_warehouses.enterprise) FROM public.transfer_between_warehouses WHERE enterprise = $1 AND finished = $2`
	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, q.enterprise)
	interfaces = append(interfaces, q.Finished)

	if q.DateStart != nil {
		sqlStatement += ` AND date_created >= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, q.DateStart)
	}
	if q.DateEnd != nil {
		sqlStatement += ` AND date_created <= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, q.DateEnd)
	}

	sqlStatement += ` ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return transfers
	}
	defer rows.Close()

	for rows.Next() {
		t := TransferBetweenWarehouses{}
		rows.Scan(&t.Id, &t.WarehouseOrigin, &t.WarehouseDestination, &t.enterprise, &t.DateCreated, &t.DateFinished, &t.Finished, &t.LinesTransfered, &t.LinesTotal, &t.Name,
			&t.WarehouseOriginName, &t.WarehouseDestinationName)
		transfers = append(transfers, t)
	}

	return transfers
}

func getTransferBetweenWarehouses(transferBetweenWarehousesId int64) TransferBetweenWarehouses {
	sqlStatement := `SELECT * FROM public.transfer_between_warehouses WHERE id = $1 LIMIT 1`
	row := db.QueryRow(sqlStatement, transferBetweenWarehousesId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return TransferBetweenWarehouses{}
	}

	t := TransferBetweenWarehouses{}
	row.Scan(&t.Id, &t.WarehouseOrigin, &t.WarehouseDestination, &t.enterprise, &t.DateCreated, &t.DateFinished, &t.Finished, &t.LinesTransfered, &t.LinesTotal, &t.Name)
	return t
}

func (t *TransferBetweenWarehouses) isValid() bool {
	return !(len(t.WarehouseOrigin) == 0 || len(t.WarehouseDestination) == 0 || t.WarehouseOrigin == t.WarehouseDestination || len(t.Name) == 0 || len(t.Name) > 100)
}

func (t *TransferBetweenWarehouses) insertTransferBetweenWarehouses() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.transfer_between_warehouses(warehouse_origin, warehouse_destination, enterprise, name) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, t.WarehouseOrigin, t.WarehouseDestination, t.enterprise, t.Name)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (t *TransferBetweenWarehouses) deleteTransferBetweenWarehouses() bool {
	if t.Id <= 0 {
		return false
	}

	details := getTransferBetweenWarehousesDetail(t.Id, t.enterprise)
	for i := 0; i < len(details); i++ {
		if details[i].QuantityTransfered > 0 {
			return false
		}
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	for i := 0; i < len(details); i++ {
		if !details[i].deleteTransferBetweenWarehousesDetail(trans) {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.transfer_between_warehouses WHERE id = $1 AND enterprise = $2`
	_, err := trans.Exec(sqlStatement, t.Id, t.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	return true
	///
}

type TransferBetweenWarehousesDetail struct {
	Id                        int64  `json:"id"`
	TransferBetweenWarehouses int64  `json:"transferBetweenWarehouses"`
	Product                   int32  `json:"product"`
	Quantity                  int32  `json:"quantity"`
	QuantityTransfered        int32  `json:"quantityTransfered"`
	Finished                  bool   `json:"finished"`
	ProductReference          string `json:"productReference"`
	WarehouseMovementOut      *int64 `json:"warehouseMovementOut"`
	WarehouseMovementIn       *int64 `json:"warehouseMovementIn"`
	ProductName               string `json:"productName"`
	enterprise                int32
}

func getTransferBetweenWarehousesDetail(transferBetweenWarehousesId int64, enterpriseId int32) []TransferBetweenWarehousesDetail {
	var details []TransferBetweenWarehousesDetail = make([]TransferBetweenWarehousesDetail, 0)

	sqlStatement := `SELECT *,(SELECT reference FROM product WHERE product.id = transfer_between_warehouses_detail.product),(SELECT name FROM product WHERE product.id = transfer_between_warehouses_detail.product) FROM public.transfer_between_warehouses_detail WHERE transfer_between_warehouses = $1 AND enterprise = $2 ORDER BY product ASC, id ASC`
	rows, err := db.Query(sqlStatement, transferBetweenWarehousesId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return details
	}

	for rows.Next() {
		d := TransferBetweenWarehousesDetail{}
		rows.Scan(&d.Id, &d.TransferBetweenWarehouses, &d.enterprise, &d.Product, &d.Quantity, &d.QuantityTransfered, &d.Finished, &d.WarehouseMovementOut, &d.WarehouseMovementIn, &d.ProductReference, &d.ProductName)
		details = append(details, d)
	}

	return details
}

// For internal use only
func getTransferBetweenWarehousesDetailRow(transferBetweenWarehousesDetailId int64) TransferBetweenWarehousesDetail {
	sqlStatement := `SELECT * FROM public.transfer_between_warehouses_detail WHERE id = $1 LIMIT 1`
	row := db.QueryRow(sqlStatement, transferBetweenWarehousesDetailId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return TransferBetweenWarehousesDetail{}
	}

	d := TransferBetweenWarehousesDetail{}
	row.Scan(&d.Id, &d.TransferBetweenWarehouses, &d.enterprise, &d.Product, &d.Quantity, &d.QuantityTransfered, &d.Finished, &d.WarehouseMovementOut, &d.WarehouseMovementIn)
	return d
}

func (d *TransferBetweenWarehousesDetail) isValid() bool {
	return !(d.TransferBetweenWarehouses <= 0 || d.Product <= 0 || d.Quantity <= 0)
}

func (d *TransferBetweenWarehousesDetail) insertTransferBetweenWarehousesDetail() bool {
	if !d.isValid() {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	transfer := getTransferBetweenWarehouses(d.TransferBetweenWarehouses)
	if transfer.Id <= 0 || transfer.enterprise != d.enterprise || transfer.Finished {
		return false
	}

	sqlStatement := `INSERT INTO public.transfer_between_warehouses_detail(transfer_between_warehouses, enterprise, product, quantity) VALUES ($1, $2, $3, $4)`
	_, err := trans.Exec(sqlStatement, d.TransferBetweenWarehouses, d.enterprise, d.Product, d.Quantity)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE public.transfer_between_warehouses SET lines_total = lines_total + 1 WHERE id = $1`
	_, err = trans.Exec(sqlStatement, d.TransferBetweenWarehouses)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	return true
	///
}

func (d *TransferBetweenWarehousesDetail) deleteTransferBetweenWarehousesDetail(trans *sql.Tx) bool {
	if d.Id <= 0 {
		return false
	}

	detailInMemory := getTransferBetweenWarehousesDetailRow(d.Id)
	if detailInMemory.Id <= 0 || detailInMemory.enterprise != d.enterprise || detailInMemory.QuantityTransfered > 0 {
		return false
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false
		}
		///
	}

	sqlStatement := `DELETE FROM public.transfer_between_warehouses_detail WHERE id = $1 AND enterprise = $2`
	_, err := trans.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `UPDATE public.transfer_between_warehouses SET lines_total = lines_total + 1 WHERE id = $1`
	_, err = trans.Exec(sqlStatement, d.TransferBetweenWarehouses)
	if err != nil {
		log("DB", err.Error())
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

func (q *TransferBetweenWarehousesDetailBarCodeQuery) transferBetweenWarehousesDetailBarCode(enterpriseId int32, userId int32) bool {
	if !q.isValid() {
		return false
	}

	if len(q.BarCode) != 13 {
		q.BarCode = fmt.Sprintf("%013s", q.BarCode)
	}

	sqlStatement := `SELECT id FROM public.transfer_between_warehouses_detail WHERE enterprise = $1 AND transfer_between_warehouses = $2 AND quantity_transfered < quantity AND product = (SELECT id FROM product WHERE product.enterprise = $1 AND product.barCode = $3 LIMIT 1) ORDER BY id ASC LIMIT 1`
	row := db.QueryRow(sqlStatement, enterpriseId, q.TransferBetweenWarehousesId, q.BarCode)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var transferBetweenWarehousesDetailId int64
	row.Scan(&transferBetweenWarehousesDetailId)
	if transferBetweenWarehousesDetailId <= 0 {
		return false
	}

	detail := getTransferBetweenWarehousesDetailRow(transferBetweenWarehousesDetailId)
	if detail.Id <= 0 || detail.enterprise != enterpriseId {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement = `UPDATE public.transfer_between_warehouses_detail SET quantity_transfered=quantity_transfered+1, finished=(quantity_transfered+1)=quantity WHERE id=$1`
	_, err := trans.Exec(sqlStatement, transferBetweenWarehousesDetailId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	if detail.Quantity == detail.QuantityTransfered+1 {
		sqlStatement := `UPDATE public.transfer_between_warehouses SET date_finished=CASE WHEN (lines_transfered+1=lines_total) THEN CURRENT_TIMESTAMP(3) ELSE NULL END, finished=lines_transfered+1=lines_total, lines_transfered=lines_transfered+1 WHERE id=$1`
		_, err = trans.Exec(sqlStatement, detail.TransferBetweenWarehouses)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		transfer := getTransferBetweenWarehouses(detail.TransferBetweenWarehouses)
		wmOut := WarehouseMovement{
			Warehouse:  transfer.WarehouseOrigin,
			Product:    detail.Product,
			Quantity:   detail.Quantity,
			Type:       "O",
			enterprise: detail.enterprise,
		}
		if !wmOut.insertWarehouseMovement(userId, trans) {
			trans.Rollback()
			return false
		}

		wmIn := WarehouseMovement{
			Warehouse:  transfer.WarehouseDestination,
			Product:    detail.Product,
			Quantity:   detail.Quantity,
			Type:       "I",
			enterprise: detail.enterprise,
		}
		if !wmIn.insertWarehouseMovement(userId, trans) {
			trans.Rollback()
			return false
		}

		sqlStatement = `UPDATE public.transfer_between_warehouses_detail SET warehouse_movement_out=$2, warehouse_movement_in=$3 WHERE id=$1`
		_, err = trans.Exec(sqlStatement, detail.Id, wmOut.Id, wmIn.Id)
		if err != nil {
			log("DB", err.Error())
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
	if detail.Id <= 0 || detail.enterprise != enterpriseId || detail.QuantityTransfered+q.Quantity > detail.Quantity {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `UPDATE public.transfer_between_warehouses_detail SET quantity_transfered=quantity_transfered+1, finished=(quantity_transfered+1)=quantity WHERE id=$1`
	_, err := trans.Exec(sqlStatement, detail.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	if detail.Quantity == detail.QuantityTransfered+1 {
		sqlStatement := `UPDATE public.transfer_between_warehouses SET date_finished=CASE WHEN (lines_transfered+1=lines_total) THEN CURRENT_TIMESTAMP(3) ELSE NULL END, finished=lines_transfered+1=lines_total, lines_transfered=lines_transfered+1 WHERE id=$1`
		_, err = trans.Exec(sqlStatement, detail.TransferBetweenWarehouses)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		transfer := getTransferBetweenWarehouses(detail.TransferBetweenWarehouses)
		wmOut := WarehouseMovement{
			Warehouse:  transfer.WarehouseOrigin,
			Product:    detail.Product,
			Quantity:   detail.Quantity,
			Type:       "O",
			enterprise: detail.enterprise,
		}
		if !wmOut.insertWarehouseMovement(userId, trans) {
			trans.Rollback()
			return false
		}

		wmIn := WarehouseMovement{
			Warehouse:  transfer.WarehouseDestination,
			Product:    detail.Product,
			Quantity:   detail.Quantity,
			Type:       "I",
			enterprise: detail.enterprise,
		}
		if !wmIn.insertWarehouseMovement(userId, trans) {
			trans.Rollback()
			return false
		}

		sqlStatement = `UPDATE public.transfer_between_warehouses_detail SET warehouse_movement_out=$2, warehouse_movement_in=$3 WHERE id=$1`
		_, err = trans.Exec(sqlStatement, detail.Id, wmOut.Id, wmIn.Id)
		if err != nil {
			log("DB", err.Error())
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

	transfer := getTransferBetweenWarehouses(transferBetweenWarehousesId)
	if transfer.Id <= 0 || transfer.enterprise != enterpriseId {
		return movements
	}

	details := getTransferBetweenWarehousesDetail(transferBetweenWarehousesId, enterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].WarehouseMovementOut != nil {
			m := getWarehouseMovementRow(*details[i].WarehouseMovementOut)
			m.WarehouseName = getNameWarehouse(m.Warehouse, enterpriseId)
			m.ProductName = getNameProduct(m.Product, enterpriseId)
			movements = append(movements, m)
		}
		if details[i].WarehouseMovementIn != nil {
			m := getWarehouseMovementRow(*details[i].WarehouseMovementIn)
			m.WarehouseName = getNameWarehouse(m.Warehouse, enterpriseId)
			m.ProductName = getNameProduct(m.Product, enterpriseId)
			movements = append(movements, m)
		}
	}

	return movements
}
