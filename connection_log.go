/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type ConnectionLog struct {
	Id               int64      `json:"id"`
	DateConnected    time.Time  `json:"dateConnected" gorm:"type:timestamp(3) with time zone;not null:true;index:connection_log_date_connected,sort:desc;index:connection_log_user_date_connected,sort:desc,priority:2,where:date_disconnected IS NULL"`
	DateDisconnected *time.Time `json:"dateDisconnected" gorm:"type:timestamp(3) with time zone"`
	UserId           int32      `json:"userId" gorm:"column:user;not null:true;index:connection_log_user_date_connected,priority:1,where:date_disconnected IS NULL"`
	User             User       `json:"user" gorm:"foreignkey:UserId,EnterpriseId;references:Id,EnterpriseId"`
	Ok               bool       `json:"ok" gorm:"not null:true"`
	IpAddress        string     `json:"ipAddress" gorm:"not null:true;type:inet"`
	EnterpriseId     int32      `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise       Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (cl *ConnectionLog) TableName() string {
	return "connection_log"
}

type ConnectionLogQuery struct {
	enterprise int32
	Offset     int64      `json:"offset"`
	Limit      int64      `json:"limit"`
	DateStart  *time.Time `json:"dateStart"`
	DateEnd    *time.Time `json:"dateEnd"`
	Ok         *bool      `json:"ok"`
	UserId     *int32     `json:"userId"`
}

type ConnectionLogs struct {
	Logs []ConnectionLog `json:"logs"`
	Rows int64           `json:"rows"`
}

func (q *ConnectionLogQuery) getConnectionLogs() ConnectionLogs {
	logs := make([]ConnectionLog, 0)
	// get all connection logs from the database for the current enterprise and pagination using dbOrm
	cursor := dbOrm.Where("connection_log.enterprise = ?", q.enterprise)
	if q.DateStart != nil {
		cursor = cursor.Where("date_connected >= ?", *q.DateStart)
	}
	if q.DateEnd != nil {
		cursor = cursor.Where("date_connected <= ?", *q.DateEnd)
	}
	if q.Ok != nil {
		cursor = cursor.Where("ok = ?", *q.Ok)
	}
	if q.UserId != nil {
		cursor = cursor.Where(`"user" = ?`, *q.UserId)
	}
	result := cursor.Offset(int(q.Offset)).Limit(int(q.Limit)).Preload("User").Order("connection_log.date_connected DESC").Find(&logs)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	// get all the connection logs count from the database using dbOrm
	var count int64
	cursor = dbOrm.Model(&ConnectionLog{}).Where("connection_log.enterprise = ?", q.enterprise)
	if q.DateStart != nil {
		cursor = cursor.Where("date_connected >= ?", *q.DateStart)
	}
	if q.DateEnd != nil {
		cursor = cursor.Where("date_connected <= ?", *q.DateEnd)
	}
	if q.Ok != nil {
		cursor = cursor.Where("ok = ?", *q.Ok)
	}
	if q.UserId != nil {
		cursor = cursor.Where(`"user" = ?`, *q.UserId)
	}
	result = cursor.Count(&count)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return ConnectionLogs{Logs: logs, Rows: count}
}

func (l *ConnectionLog) BeforeCreate(tx *gorm.DB) (err error) {
	var connectionLog ConnectionLog
	tx.Model(&ConnectionLog{}).Last(&connectionLog)
	l.Id = connectionLog.Id + 1
	return nil
}

func (l *ConnectionLog) insertConnectionLog() {
	l.DateConnected = time.Now()
	result := dbOrm.Create(&l)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}

// Called during the client login.
// 1. Logs the user connection
// 2. Filters the user connection
func userConnection(userId int32, ipAddress string, enterpriseId int32) (bool, string) {
	s := getSettingsRecordById(enterpriseId)
	if !s.ConnectionLog {
		return true, ""
	}

	// Remote the port from the address
	if strings.Contains(ipAddress, ":") {
		ipAddress = ipAddress[:strings.Index(ipAddress, ":")]
	}
	l := ConnectionLog{UserId: userId, IpAddress: ipAddress, EnterpriseId: enterpriseId}

	// the default user ("marketnet") is not filtered
	if userId == 1 {
		l.Ok = true
		l.insertConnectionLog()
		return true, ""
	}

	if s.FilterConnections {
		filters := getConnectionFiltersByUser(userId)
		for i := 0; i < len(filters); i++ {
			if filters[i].Type == "I" {
				if *filters[i].IpAddress != ipAddress {
					l.Ok = false
					l.insertConnectionLog()
					return false, filters[i].Name
				}
			} else if filters[i].Type == "S" {
				now := time.Now()
				// 0000-01-01 19:03:54 +0000 UTC
				time := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
				if time.Before(*filters[i].TimeStart) || time.After(*filters[i].TimeEnd) {
					l.Ok = false
					l.insertConnectionLog()
					return false, filters[i].Name
				}
			}
		}
	}

	l.Ok = true
	l.insertConnectionLog()
	return true, ""
}

func userDisconnected(user int32) {
	// get a single connection log for the user where the data disconnected is null sorted by date connected descending using dbOrm
	var connectionLog ConnectionLog
	result := dbOrm.Where(`"user" = ? AND date_disconnected IS NULL`, user).Order("date_connected DESC").First(&connectionLog)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return
	}

	now := time.Now()
	connectionLog.DateDisconnected = &now

	// update the connection log with the date disconnected using dbOrm
	result = dbOrm.Save(&connectionLog)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}

func cleanUpConnectionLogs(enterpriseId int32) {
	settings := getSettingsRecordById(enterpriseId)

	result := dbOrm.Model(&ConnectionLog{}).Where("enterprise = ? AND date_connected < ?", enterpriseId, time.Now().Add(-time.Duration(settings.SettingsCleanUp.ConnectionLogDays)*time.Hour*24)).Delete(&ConnectionLog{})
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}

type ConnectionFilter struct {
	Id           int32      `json:"id"`
	Name         string     `json:"name" gorm:"not null:true;type:character varying(100)"`
	Type         string     `json:"type" gorm:"type:character(1);not null:true"` // I = IP, S = Schedule
	IpAddress    *string    `json:"ipAddress" gorm:"type:inet"`
	TimeStart    *time.Time `json:"timeStart" gorm:"column:time_start;type:timestamp(0) with time zone"`
	TimeEnd      *time.Time `json:"timeEnd" gorm:"column:time_end;type:timestamp(0) with time zone"`
	EnterpriseId int32      `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise   Settings   `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (cf *ConnectionFilter) TableName() string {
	return "connection_filter"
}

func getConnectionFilters(enterpriseId int32) []ConnectionFilter {
	filters := make([]ConnectionFilter, 0)
	// get all connection filters for the current enterprise using dbOrm sorted by id ascending
	result := dbOrm.Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&filters)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return filters
}

func getConnectionFilterRow(id int32) ConnectionFilter {
	f := ConnectionFilter{}
	// get a single connection filter row by id using dbOrm
	dbOrm.Where("id = ?", id).First(&f)
	return f
}

func getConnectionFiltersByUser(userId int32) []ConnectionFilter {
	// get the connetion filter users by user id using dbOrm
	var connectionFilterUsers []ConnectionFilterUser
	dbOrm.Where("user_id = ?", userId).Find(&connectionFilterUsers)
	// get the connection filter from the connection filter users using dbOrm
	var filters []ConnectionFilter
	for i := 0; i < len(connectionFilterUsers); i++ {
		filters = append(filters, getConnectionFilterRow(connectionFilterUsers[i].ConnectionFilterId))
	}
	return filters
}

func (f *ConnectionFilter) BeforeCreate(tx *gorm.DB) (err error) {
	var connectionFilter ConnectionFilter
	tx.Model(&ConnectionFilter{}).Last(&connectionFilter)
	f.Id = connectionFilter.Id + 1
	return nil
}

func (f *ConnectionFilter) isValid() bool {
	return !(len(f.Name) == 0 || len(f.Name) > 100 || (f.Type != "I" && f.Type != "S") || (f.Type == "I" && (f.IpAddress == nil || f.TimeStart != nil || f.TimeEnd != nil)) || (f.Type == "S" && (f.IpAddress != nil || f.TimeStart == nil || f.TimeEnd == nil)))
}

func (f *ConnectionFilter) cleanUp() {
	if f.Type == "S" {
		timeStart := time.Date(0, 1, 1, f.TimeStart.Hour(), f.TimeStart.Minute(), f.TimeStart.Second(), 0, time.UTC)
		f.TimeStart = &timeStart

		timeEnd := time.Date(0, 1, 1, f.TimeEnd.Hour(), f.TimeEnd.Minute(), f.TimeEnd.Second(), 0, time.UTC)
		f.TimeEnd = &timeEnd

		f.IpAddress = nil
	} else if f.Type == "I" {
		f.TimeStart = nil
		f.TimeEnd = nil
	}
}

func (f *ConnectionFilter) insertConnectionFilter() bool {
	if !f.isValid() {
		return false
	}
	f.cleanUp()

	result := dbOrm.Create(&f)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (f *ConnectionFilter) updateConnectionFilter() bool {
	filter := getConnectionFilterRow(f.Id)
	if filter.Id < 0 {
		return false
	}
	if f.Type != filter.Type {
		return false
	}
	if !f.isValid() {
		return false
	}
	f.cleanUp()

	// get a single connection filter row by id using dbOrm
	var connectionFilter ConnectionFilter
	result := dbOrm.Where("id = ?", f.Id).First(&connectionFilter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// copy all the attributes from the f object to the connection filter object
	connectionFilter.Name = f.Name
	connectionFilter.Type = f.Type
	connectionFilter.TimeStart = f.TimeStart
	connectionFilter.TimeEnd = f.TimeEnd

	// update the connection filter using dbOrm
	result = dbOrm.Save(&connectionFilter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (f *ConnectionFilter) deleteConnectionFilter() bool {
	// delete a single connetion filter by id and enterprise using dbOrm
	result := dbOrm.Where("id = ? AND enterprise = ?", f.Id, f.EnterpriseId).Delete(ConnectionFilter{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

type ConnectionFilterUser struct {
	ConnectionFilterId int32            `json:"connectionFilterId" gorm:"primaryKey;column:connection_filter;not null:true"`
	ConnectionFilter   ConnectionFilter `json:"connectionFilter" gorm:"foreignKey:ConnectionFilterId;references:Id"`
	UserId             int32            `json:"userId" gorm:"primaryKey;column:user;not null:true"`
	User               User             `json:"user" gorm:"foreignKey:UserId;references:Id"`
}

func (cfu *ConnectionFilterUser) TableName() string {
	return "connection_filter_user"
}

func getConnectionFilterUser(filterId int32, enterpriseId int32) []ConnectionFilterUser {
	// get a single connection filter row by id using dbOrm
	var connectionFilter ConnectionFilter
	result := dbOrm.Where("id = ? AND enterprise = ?", filterId, enterpriseId).First(&connectionFilter)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	if connectionFilter.EnterpriseId != enterpriseId {
		return nil
	}

	// get all connection filter users for the current filter using dbOrm sorted by id ascending
	var connectionFilterUsers []ConnectionFilterUser
	result = dbOrm.Where("connection_filter_user.connection_filter = ?", filterId).Order(`connection_filter_user.connection_filter,connection_filter_user."user" ASC`).Preload("User").Find(&connectionFilterUsers)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return connectionFilterUsers
}

func getConnectionFilterUserByUser(userId int32, enterpriseId int32) []ConnectionFilterUser {
	// get a single connection filter row by id using dbOrm
	user := getUserRow(userId)
	if user.EnterpriseId != enterpriseId {
		return nil
	}

	// get all connection filter users for the current filter using dbOrm sorted by id ascending
	var connectionFilterUsers []ConnectionFilterUser
	result := dbOrm.Where(`connection_filter_user."user" = ?`, userId).Order(`connection_filter_user.connection_filter,connection_filter_user."user" ASC`).Preload("ConnectionFilter").Find(&connectionFilterUsers)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return connectionFilterUsers
}

func (f *ConnectionFilterUser) insertConnectionFilterUser(enterpriseId int32) bool {
	filterInMemory := getConnectionFilterRow(f.ConnectionFilterId)
	if filterInMemory.EnterpriseId != enterpriseId {
		return false
	}
	userInMemory := getUserRow(f.UserId)
	if userInMemory.EnterpriseId != enterpriseId {
		return false
	}

	result := dbOrm.Create(&f)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (f *ConnectionFilterUser) deleteConnectionFilterUser(enterpriseId int32) bool {
	filterInMemory := getConnectionFilterRow(f.ConnectionFilterId)
	if filterInMemory.EnterpriseId != enterpriseId {
		return false
	}
	userInMemory := getUserRow(f.UserId)
	if userInMemory.EnterpriseId != enterpriseId {
		return false
	}

	result := dbOrm.Where(`connection_filter_user.connection_filter = ? AND connection_filter_user."user" = ?`, f.ConnectionFilterId, f.UserId).Delete(&ConnectionFilterUser{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
