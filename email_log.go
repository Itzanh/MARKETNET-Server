package main

import (
	"strconv"
	"time"
)

type EmailLog struct {
	Id               int64     `json:"id"`
	EmailFrom        string    `json:"emailFrom"`
	NameFrom         string    `json:"nameFrom"`
	DestinationEmail string    `json:"destinationEmail"`
	DestinationName  string    `json:"destinationName"`
	Subject          string    `json:"subject"`
	Content          string    `json:"content"`
	DateSent         time.Time `json:"dateSent"`
	enterprise       int32
}

type EmailLogSearch struct {
	SearchText    string     `json:"searchText"`
	DateSentStart *time.Time `json:"dateSentStart"`
	DateSentEnd   *time.Time `json:"dateSentEnd"`
}

func (search *EmailLogSearch) getEmailLogs(enterpriseId int32) []EmailLog {
	var emailLogs []EmailLog = make([]EmailLog, 0)
	var interfaces []interface{} = make([]interface{}, 0)
	interfaces = append(interfaces, enterpriseId)

	sqlStatement := `SELECT * FROM public.email_log WHERE enterprise=$1`

	if search.SearchText != "" {
		sqlStatement += ` AND (email_from ILIKE $2 OR name_from ILIKE $2 OR destination_email ILIKE $2 OR destination_name ILIKE $2 OR subject ILIKE $2)`
		interfaces = append(interfaces, "%"+search.SearchText+"%")
	}

	if search.DateSentStart != nil {
		sqlStatement += ` AND date_sent >= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.DateSentStart)
	}

	if search.DateSentEnd != nil {
		sqlStatement += ` AND date_sent <= $` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, search.DateSentEnd)
	}

	sqlStatement += ` ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return emailLogs
	}

	for rows.Next() {
		p := EmailLog{}
		rows.Scan(&p.Id, &p.EmailFrom, &p.NameFrom, &p.DestinationEmail, &p.DestinationName, &p.Subject, &p.Content, &p.DateSent, &p.enterprise)
		emailLogs = append(emailLogs, p)
	}

	return emailLogs
}

func (el *EmailLog) insertEmailLog() bool {
	sqlStatement := `INSERT INTO public.email_log(email_from, name_from, destination_email, destination_name, subject, content, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, el.EmailFrom, el.NameFrom, el.DestinationEmail, el.DestinationName, el.Subject, el.Content, el.enterprise)
	return err == nil
}
