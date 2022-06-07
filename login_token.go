package main

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const LOGIN_TOKEN_LENGTH = 128

type LoginToken struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name" gorm:"type:character(128);not null:true;index:login_tokens_name_ip_address,unique:true"`
	DateLastUsed time.Time `json:"dateLastUsed" gorm:"type:timestamp(3) with time zone;not null:true"`
	UserId       int32     `json:"user" gorm:"column:user;type:integer;not null:true"`
	User         User      `json:"-" gorm:"foreignKey:UserId;references:Id"`
	IpAddress    string    `json:"ipAddress" gorm:"type:inet;not null:true;index:login_tokens_name_ip_address,unique:true"`
}

func (l *LoginToken) TableName() string {
	return "login_tokens"
}

func (lt *LoginToken) BeforeCreate(tx *gorm.DB) (err error) {
	var loginToken LoginToken
	tx.Model(&LoginToken{}).Last(&loginToken)
	lt.Id = loginToken.Id + 1
	return nil
}

func (t *LoginToken) insertLoginToken() bool {
	t.generateRandomToken()

	t.DateLastUsed = time.Now()

	// insert the login token into the database using dbOrm
	result := dbOrm.Create(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (t *LoginToken) generateRandomToken() {
	const CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	t.Name = ""

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < LOGIN_TOKEN_LENGTH; i++ {
		t.Name += string(CHARSET[r.Intn(len(CHARSET))])
	}
}

// The token has been used, and is not rolled to a different security token
func (t *LoginToken) rollToken() string {
	t.generateRandomToken()
	result := dbOrm.Model(&LoginToken{}).Where("id = ?", t.Id).Update("name", t.Name)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return t.Name
}

// Ok, user permissions, user id, enterprise id, rolled token
func (t *LoginToken) checkLoginToken() (bool, *Permissions, int32, int32, string) {
	if len(t.Name) != LOGIN_TOKEN_LENGTH {
		return false, nil, 0, 0, ""
	}

	// get a single login token from the database there the name and ip address are the same as the ones passed using dbOrm
	tok := LoginToken{}
	result := dbOrm.Where("name = ? AND ip_address = ?", t.Name, t.IpAddress).Preload("User").First(&tok)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, nil, 0, 0, ""
	}

	if tok.Id <= 0 {
		return false, nil, 0, 0, ""
	}

	if time.Until(tok.DateLastUsed).Hours() > float64(settings.Server.TokenExpirationHours) { // the token has expired, delete it and return an error
		result = dbOrm.Delete(&tok)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false, nil, 0, 0, ""
		}
		return false, nil, 0, 0, ""
	} else { // the token is still valid, renew the token and return OK
		tok.DateLastUsed = time.Now()
		result := dbOrm.Model(&LoginToken{}).Where("id = ?", tok.Id).Update("date_last_used", tok.DateLastUsed)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false, nil, 0, 0, ""
		}
		perm := getUserPermissions(tok.UserId, tok.User.EnterpriseId)
		return true, &perm, tok.UserId, tok.User.EnterpriseId, tok.rollToken()
	}
}

func deleteLoginTokensFromUser(userId int32, enterpriseId int32) bool {
	// get user row
	user := getUserRow(userId)
	// check the user's enterprise
	if user.EnterpriseId != enterpriseId {
		return false
	}

	result := dbOrm.Model(&LoginToken{}).Where(`"user" = ?`, userId).Delete(&LoginToken{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func deleteAllLoginTokens(enterpriseId int32) bool {
	result := dbOrm.Model(&LoginToken{}).Where(`"user" IN (SELECT id FROM "user" WHERE config = ?)`, enterpriseId).Delete(&LoginToken{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
