/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type EmailLog struct {
	Id               int64     `json:"id"`
	EmailFrom        string    `json:"emailFrom" gorm:"type:character varying(100);not null:true;index:email_log_trn,type:gin"`
	NameFrom         string    `json:"nameFrom" gorm:"type:character varying(100);not null:true;index:email_log_trn,type:gin"`
	DestinationEmail string    `json:"destinationEmail" gorm:"type:character varying(100);not null:true;index:email_log_trn,type:gin"`
	DestinationName  string    `json:"destinationName" gorm:"type:character varying(100);not null:true;index:email_log_trn,type:gin"`
	Subject          string    `json:"subject" gorm:"type:character varying(100);not null:true;index:email_log_trn,type:gin"`
	Content          string    `json:"content" gorm:"type:text;not null:true"`
	DateSent         time.Time `json:"dateSent" gorm:"type:timestamp(3) with time zone;not null:true;index:email_log_date_sent,sort:desc"`
	EnterpriseId     int32     `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise       Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (e *EmailLog) TableName() string {
	return "email_log"
}

type EmailLogSearch struct {
	SearchText    string     `json:"searchText"`
	DateSentStart *time.Time `json:"dateSentStart"`
	DateSentEnd   *time.Time `json:"dateSentEnd"`
}

func (search *EmailLogSearch) getEmailLogs(enterpriseId int32) []EmailLog {
	var emailLogs []EmailLog = make([]EmailLog, 0)

	// create a database cursor to query the email logs table using dbOrm
	cursor := dbOrm.Model(&EmailLog{}).Where("enterprise = ?", enterpriseId)

	if search.SearchText != "" {
		// search for the search text in the email logs table using named parameters
		cursor.Where("email_from LIKE @search OR destination_email LIKE @search OR subject LIKE @search", sql.Named("search", "%"+search.SearchText+"%"))
	}

	if search.DateSentStart != nil {
		cursor.Where("date_sent >= ?", search.DateSentStart)
	}

	if search.DateSentEnd != nil {
		cursor.Where("date_sent <= ?", search.DateSentEnd)
	}

	result := cursor.Order("date_sent DESC").Find(&emailLogs)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return emailLogs
	}
	return emailLogs
}

func (el *EmailLog) BeforeCreate(tx *gorm.DB) (err error) {
	var emailLog EmailLog
	tx.Model(&EmailLog{}).Last(&emailLog)
	el.Id = emailLog.Id + 1
	return nil
}

func (el *EmailLog) insertEmailLog() bool {
	el.DateSent = time.Now()
	result := dbOrm.Create(&el)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
