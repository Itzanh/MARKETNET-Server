package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// A transactional log for tables like the invoces in necessary in order to comply with the laws in e-commerce and electronic invoicing.
// Nobody should be able to modify or alter these registers without it being properly registered in a log complete log.

type TransactionalLog struct {
	Id          int64 `json:"id"`
	enterprise  int32
	Table       string    `json:"table"`
	Register    string    `json:"register"`
	DateCreated time.Time `json:"dateCreated"`
	RegisterId  int64     `json:"registerId"`
	User        *int32    `json:"user"`
	Mode        string    `json:"mode"`
	UserName    *string   `json:"userName"`
}

type TransactionalLogQuery struct {
	enterpriseId int32
	TableName    string `json:"tableName"`
	RegisterId   int    `json:"registerId"`
}

func (query *TransactionalLogQuery) getRegisterTransactionalLogs() []TransactionalLog {
	var logs []TransactionalLog = make([]TransactionalLog, 0)
	sqlStatement := `SELECT *,(SELECT username FROM "user" WHERE "user".id=transactional_log."user") FROM public.transactional_log WHERE enterprise = $1 AND "table" = $2 AND register_id = $3 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, query.enterpriseId, query.TableName, query.RegisterId)
	if err != nil {
		log("DB", err.Error())
		return logs
	}

	for rows.Next() {
		s := TransactionalLog{}
		rows.Scan(&s.Id, &s.enterprise, &s.Table, &s.Register, &s.DateCreated, &s.RegisterId, &s.User, &s.Mode, &s.UserName)
		logs = append(logs, s)
	}

	return logs
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
	sqlStatement = `INSERT INTO public.transactional_log(enterprise, "table", register, register_id, "user", mode) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(sqlStatement, enterpriseId, tableName, JsonData, registerId, user, mode)
	if err != nil {
		log("DB", err.Error())
		fmt.Println(err.Error())
		fmt.Println(*user)
		return
	}

}
