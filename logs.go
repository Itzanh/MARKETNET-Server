package main

import (
	"time"
)

func log(title string, info string) bool {
	sqlStatement := `INSERT INTO public.logs(title, info) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, title, info)
	return err == nil
}

// clear logs older than one month
// runned by cron, run every month
func clearLogs() {
	now := time.Now()
	now = now.AddDate(0, -1, 0)
	sqlStatement := `DELETE FROM public.logs WHERE date_created <= $1`
	db.Exec(sqlStatement, now)
}
