package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Log struct {
	id          int64
	DateCreated time.Time `json:"dateCreated"`
	Title       string    `json:"title"`
	Info        string    `json:"info"`
	Stacktrace  string    `json:"stacktrace"`
}

type CrashReporter struct {
	LicenseCode string `json:"licenseCode"`
	Logs        []Log  `json:"logs"`
}

const CRASHREPORTER_URL = "https://license.marketneterp.io:12278/crash_reports"

func crashreporter() {
	var logs []Log = make([]Log, 0)
	sqlStatement := `SELECT * FROM public.logs ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		l := Log{}
		rows.Scan(&l.id, &l.DateCreated, &l.Title, &l.Info, &l.Stacktrace)
		logs = append(logs, l)
	}

	if len(logs) == 0 {
		return
	}

	crashreport := CrashReporter{
		Logs: logs,
	}
	data, _ := json.Marshal(crashreport)
	_, err = http.Post(CRASHREPORTER_URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	sqlStatement = `DELETE FROM public.logs WHERE id >= 1 AND id <= $1`
	db.Exec(sqlStatement, logs[len(logs)-1].id)
}
