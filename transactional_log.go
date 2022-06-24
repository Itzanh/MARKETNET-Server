/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// A transactional log for tables like the invoces in necessary in order to comply with the laws in e-commerce and electronic invoicing.
// Nobody should be able to modify or alter these registers without it being properly registered in a log complete log.

type TransactionalLog struct {
	Id           int64     `json:"id"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise;not null:true;index:transactional_log_enterprise_table_register_id,priority:1"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Table        string    `json:"table" gorm:"column:table;not null:true;type:character varying(150);index:transactional_log_enterprise_table_register_id,priority:2"`
	Register     string    `json:"register" gorm:"column:register;not null:true;type:jsonb"`
	DateCreated  time.Time `json:"dateCreated" gorm:"type:timestamp(3) without time zone;not null:true"`
	RegisterId   int64     `json:"registerId" gorm:"column:register_id;not null:true;index:transactional_log_enterprise_table_register_id,priority:3"`
	UserId       *int32    `json:"userId" gorm:"column:user"`
	User         *User     `json:"user" gorm:"foreignKey:UserId,EnterpriseId;references:Id,EnterpriseId"`
	Mode         string    `json:"mode" gorm:"column:mode;not null:true;type:character(1)"`
}

func (t *TransactionalLog) TableName() string {
	return "transactional_log"
}

type TransactionalLogQuery struct {
	enterpriseId int32
	TableName    string `json:"tableName"`
	RegisterId   int    `json:"registerId"`
}

func (query *TransactionalLogQuery) getRegisterTransactionalLogs() []TransactionalLog {
	var logs []TransactionalLog = make([]TransactionalLog, 0)
	// get all the transactional logs from the database for this enterprise id, table name and register id sorted by date created ascending using dbOrm
	result := dbOrm.Where("transactional_log.enterprise = ? AND transactional_log.table = ? AND transactional_log.register_id = ?", query.enterpriseId, query.TableName, query.RegisterId).Preload("User").Order("transactional_log.date_created ASC").Find(&logs)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return logs
}

func (tl *TransactionalLog) BeforeCreate(tx *gorm.DB) (err error) {
	var transactionalLog TransactionalLog
	tx.Model(&TransactionalLog{}).Last(&transactionalLog)
	tl.Id = transactionalLog.Id + 1
	return nil
}

// enterprise id
// table name eg: sales_order
// value of the "id" field in the table
// user ID modifier or "0" for automatic/unattended/cron processes
// mode: I = Insert, U = Update, D = Delete
func insertTransactionalLog(enterpriseId int32, tableName string, registerId int, userId int32, mode string) {
	s := getSettingsRecordById(enterpriseId)
	if !s.TransactionLog {
		return
	}

	if mode != "I" && mode != "U" && mode != "D" {
		return
	}

	// query the columns
	sqlStatement := `SELECT column_name FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position ASC`
	rows, err := db.Query(sqlStatement, tableName)
	if err != nil {
		log("DB", err.Error())
		return
	}
	defer rows.Close()

	columnNames := make([]string, 0)
	for rows.Next() {
		var columnName string
		rows.Scan(&columnName)
		columnNames = append(columnNames, columnName)
	}

	// query the row
	sqlStatement = `SELECT * FROM ` + tableName + ` WHERE id=$1`
	row := db.QueryRow(sqlStatement, registerId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return
	}

	fields := make([]*string, len(columnNames))
	data := make([]interface{}, len(columnNames))

	for i := 0; i < len(columnNames); i++ {
		data[i] = &fields[i]
	}

	err = row.Scan(data...)
	if err != nil {
		fmt.Println(err)
		return
	}

	register := make(map[string]interface{})

	for i := 0; i < len(columnNames); i++ {
		register[columnNames[i]] = data[i]
	}

	JsonData, _ := json.Marshal(register)

	var user *int32 = nil
	if userId != 0 {
		user = &userId
	}

	// insert into transactional log
	transactionalLog := TransactionalLog{
		EnterpriseId: enterpriseId,
		Table:        tableName,
		Register:     string(JsonData),
		DateCreated:  time.Now(),
		RegisterId:   int64(registerId),
		UserId:       user,
		Mode:         mode,
	}
	// insert into transactional log
	result := dbOrm.Create(&transactionalLog)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}

func cleanUpTransactionalLog(enterpriseId int32) {
	settings := getSettingsRecordById(enterpriseId)

	result := dbOrm.Model(&TransactionalLog{}).Where("enterprise = ? AND date_created < ?", enterpriseId, time.Now().Add(-time.Duration(settings.SettingsCleanUp.TransactionalLogDays)*time.Hour*24)).Delete(&TransactionalLog{})
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}
