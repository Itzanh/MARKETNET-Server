package main

import (
	"strings"
	"time"
)

type ConnectionLog struct {
	Id               int64      `json:"id"`
	DateConnected    time.Time  `json:"dateConnected"`
	DateDisconnected *time.Time `json:"dateDisconnected"`
	User             int16      `json:"user"`
	Ok               bool       `json:"ok"`
	IpAddress        string     `json:"ipAddress"`
	UserName         string     `json:"userName"`
}

type ConnectionLogs struct {
	Logs []ConnectionLog `json:"logs"`
	Rows int64           `json:"rows"`
}

func (q *PaginationQuery) getConnectionLogs() ConnectionLogs {
	logs := make([]ConnectionLog, 0)
	sqlStatement := `SELECT *,(SELECT username FROM "user" WHERE "user".id=connection_log."user") FROM public.connection_log ORDER BY date_connected DESC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit)
	if err != nil {
		log("DB", err.Error())
		return ConnectionLogs{}
	}

	for rows.Next() {
		l := ConnectionLog{}
		rows.Scan(&l.Id, &l.DateConnected, &l.DateDisconnected, &l.User, &l.Ok, &l.IpAddress, &l.UserName)
		logs = append(logs, l)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.connection_log`
	row := db.QueryRow(sqlStatement)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ConnectionLogs{}
	}

	var rowCount int64
	row.Scan(&rowCount)

	return ConnectionLogs{Logs: logs, Rows: rowCount}
}

func (l *ConnectionLog) insertConnectionLog() {
	sqlStatement := `INSERT INTO public.connection_log("user", ok, ip_address) VALUES ($1, $2, $3)`
	db.Exec(sqlStatement, l.User, l.Ok, l.IpAddress)
}

// Called during the client login.
// 1. Logs the user connection
// 2. Filters the user connection
func userConnection(userId int16, ipAddress string) bool {
	s := getSettingsRecord()
	if !s.ConnectionLog {
		return true
	}

	// Remote the port from the address
	if strings.Contains(ipAddress, ":") {
		ipAddress = ipAddress[:strings.Index(ipAddress, ":")]
	}
	l := ConnectionLog{User: userId, IpAddress: ipAddress}

	// the default user ("marketnet") is not filtered
	if userId == 1 {
		l.Ok = true
		l.insertConnectionLog()
		return true
	}

	if s.FilterConnections {
		filters := getConnectionFiltersByUser(userId)
		for i := 0; i < len(filters); i++ {
			if filters[i].Type == "I" {
				if *filters[i].IpAddress != ipAddress {
					l.Ok = false
					l.insertConnectionLog()
					return false
				}
			} else if filters[i].Type == "S" {
				now := time.Now()
				// 0000-01-01 19:03:54 +0000 UTC
				time := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
				if time.Before(*filters[i].TimeStart) || time.After(*filters[i].TimeEnd) {
					l.Ok = false
					l.insertConnectionLog()
					return false
				}
			}
		}
	}

	l.Ok = true
	l.insertConnectionLog()
	return true
}

func userDisconnected(user int16) {
	sqlStatement := `UPDATE connection_log SET date_disconnected=CURRENT_TIMESTAMP(3) WHERE id=(SELECT id FROM public.connection_log WHERE "user"=$1 AND date_disconnected IS NULL ORDER BY date_connected DESC LIMIT 1)`
	db.Exec(sqlStatement, user)
}

type ConnectionFilter struct {
	Id        int16      `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // I = IP, S = Schedule
	IpAddress *string    `json:"ipAddress"`
	TimeStart *time.Time `json:"timeStart"`
	TimeEnd   *time.Time `json:"timeEnd"`
}

func getConnectionFilters() []ConnectionFilter {
	filters := make([]ConnectionFilter, 0)
	sqlStatement := `SELECT * FROM public.connection_filter ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return filters
	}

	for rows.Next() {
		f := ConnectionFilter{}
		rows.Scan(&f.Id, &f.Name, &f.Type, &f.IpAddress, &f.TimeStart, &f.TimeEnd)
		filters = append(filters, f)
	}
	return filters
}

func getConnectionFilterRow(id int16) ConnectionFilter {
	f := ConnectionFilter{}
	sqlStatement := `SELECT * FROM public.connection_filter WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return f
	}

	row.Scan(&f.Id, &f.Name, &f.Type, &f.IpAddress, &f.TimeStart, &f.TimeEnd)

	return f
}

func getConnectionFiltersByUser(userId int16) []ConnectionFilter {
	filters := make([]ConnectionFilter, 0)
	sqlStatement := `SELECT connection_filter.* FROM public.connection_filter INNER JOIN connection_filter_user ON connection_filter_user.connection_filter=connection_filter.id WHERE connection_filter_user."user"=$1`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		log("DB", err.Error())
		return filters
	}

	for rows.Next() {
		f := ConnectionFilter{}
		rows.Scan(&f.Id, &f.Name, &f.Type, &f.IpAddress, &f.TimeStart, &f.TimeEnd)
		filters = append(filters, f)
	}
	return filters
}

func (f *ConnectionFilter) insertConnectionFilter() bool {
	if len(f.Name) == 0 || len(f.Name) > 100 || (f.Type != "I" && f.Type != "S") || (f.Type == "I" && (f.IpAddress == nil || f.TimeStart != nil || f.TimeEnd != nil)) || (f.Type == "S" && (f.IpAddress != nil || f.TimeStart == nil || f.TimeEnd == nil)) {
		return false
	}

	sqlStatement := `INSERT INTO public.connection_filter(name, type, ip_address, time_start, time_end) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, f.Name, f.Type, f.IpAddress, f.TimeStart, f.TimeEnd)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (f *ConnectionFilter) updateConnectionFilter() bool {
	filter := getConnectionFilterRow(f.Id)
	if filter.Id < 0 {
		return false
	}
	if f.Type != filter.Type {
		return false
	}
	if len(f.Name) == 0 || len(f.Name) > 100 || (f.Type != "I" && f.Type != "S") || (f.Type == "I" && (f.IpAddress == nil || f.TimeStart != nil || f.TimeEnd != nil)) || (f.Type == "S" && (f.IpAddress != nil || f.TimeStart == nil || f.TimeEnd == nil)) {
		return false
	}

	sqlStatement := `UPDATE public.connection_filter SET name=$2, ip_address=$3, time_start=$4, time_end=$5 WHERE id=$1`
	_, err := db.Exec(sqlStatement, f.Id, f.Name, f.IpAddress, f.TimeStart, f.TimeEnd)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (f *ConnectionFilter) deleteConnectionFilter() bool {
	sqlStatement := `DELETE FROM public.connection_filter WHERE id=$1`
	_, err := db.Exec(sqlStatement, f.Id)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

type ConnectionFilterUser struct {
	ConnectionFilter int16  `json:"connectionFilter"`
	User             int16  `json:"user"`
	UserName         string `json:"userName"`
}

func getConnectionFilterUser(filterId int16) []ConnectionFilterUser {
	filters := make([]ConnectionFilterUser, 0)
	sqlStatement := `SELECT *,(SELECT username FROM "user" WHERE "user".id=connection_filter_user."user") FROM public.connection_filter_user WHERE connection_filter=$1 ORDER BY connection_filter ASC, "user" ASC`
	rows, err := db.Query(sqlStatement, filterId)
	if err != nil {
		log("DB", err.Error())
		return filters
	}

	for rows.Next() {
		f := ConnectionFilterUser{}
		rows.Scan(&f.ConnectionFilter, &f.User, &f.UserName)
		filters = append(filters, f)
	}
	return filters
}

func (f *ConnectionFilterUser) insertConnectionFilterUser() bool {
	sqlStatement := `INSERT INTO public.connection_filter_user(connection_filter, "user") VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, f.ConnectionFilter, f.User)
	return err == nil
}

func (f *ConnectionFilterUser) deleteConnectionFilterUser() bool {
	sqlStatement := `DELETE FROM public.connection_filter_user WHERE connection_filter=$1 AND "user"=$2`
	_, err := db.Exec(sqlStatement, f.ConnectionFilter, f.User)
	return err == nil
}
