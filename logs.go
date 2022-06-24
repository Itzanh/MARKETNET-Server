/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type Log struct {
	Id          int64
	DateCreated time.Time `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	Title       string    `json:"title" gorm:"column:title;not null:true;type:character varying(255)"`
	Info        string    `json:"info" gorm:"column:info;not null:true;type:text"`
	Stacktrace  string    `json:"stacktrace" gorm:"column:stacktrace;not null:true;type:text"`
}

func (Log) TableName() string {
	return "logs"
}

func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	var log Log
	tx.Model(&Log{}).Last(&log)
	l.Id = log.Id + 1
	return nil
}

func log(title string, info string) {
	errTrc := errors.Errorf(info)
	stackTrace := errTrc.ErrorStack()
	log := Log{
		DateCreated: time.Now(),
		Title:       title,
		Info:        info,
		Stacktrace:  stackTrace,
	}
	dbOrm.Create(&log)
}

// clear logs older than one month
// runned by cron, run every month
func clearLogs() {
	now := time.Now()
	now = now.AddDate(0, -1, 0)
	sqlStatement := `DELETE FROM public.logs WHERE date_created <= $1`
	db.Exec(sqlStatement, now)
}
