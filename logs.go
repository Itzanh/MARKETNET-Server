package main

import (
	"time"

	"github.com/go-errors/errors"
)

func log(title string, info string) bool {
	errTrc := errors.Errorf(info)
	stackTrace := errTrc.ErrorStack()
	sqlStatement := `INSERT INTO public.logs(title, info, stacktrace) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, title, info, stackTrace)
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
