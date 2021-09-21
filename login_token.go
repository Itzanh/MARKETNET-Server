package main

import (
	"math/rand"
	"time"
)

const LOGIN_TOKEN_LENGTH = 128

type LoginToken struct {
	Id           int32     `json:"id"`
	Name         string    `json:"name"`
	DateLastUsed time.Time `json:"dateLastUsed"`
	User         int32     `json:"user"`
	IpAddress    string    `json:"ipAddress"`
}

func (t *LoginToken) insertLoginToken() bool {
	t.generateRandomToken()

	sqlStatement := `INSERT INTO public.login_tokens(name, "user", ip_address) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, t.Name, t.User, t.IpAddress)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
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

// Ok, user permissions, user id, enterprise id
func (t *LoginToken) checkLoginToken() (bool, *Permissions, int32, int32) {
	if len(t.Name) != LOGIN_TOKEN_LENGTH {
		return false, nil, 0, 0
	}

	sqlStatement := `SELECT *,(SELECT config FROM "user" WHERE "user".id=login_tokens."user") FROM login_tokens WHERE name=$1 AND ip_address=$2 LIMIT 1`
	row := db.QueryRow(sqlStatement, t.Name, t.IpAddress)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, nil, 0, 0
	}

	tok := LoginToken{}
	var enterpriseId int32
	row.Scan(&tok.Id, &tok.Name, &tok.DateLastUsed, &tok.User, &tok.IpAddress, &enterpriseId)
	if tok.Id <= 0 {
		return false, nil, 0, 0
	}

	if time.Until(tok.DateLastUsed).Hours() > float64(settings.Server.TokenExpirationHours) { // the token has expired, delete it and return an error
		sqlStatement := `DELETE FROM login_tokens WHERE name=$1 AND ip_address=$2`
		db.Exec(sqlStatement, t.Name, t.IpAddress)
		return false, nil, 0, 0
	} else { // the token is still valid, renew the token and return OK
		sqlStatement := `UPDATE login_tokens SET date_last_used=CURRENT_TIMESTAMP(3) WHERE name=$1 AND ip_address=$2`
		db.Exec(sqlStatement, t.Name, t.IpAddress)
		perm := getUserPermissions(tok.User, enterpriseId)
		return true, &perm, tok.User, enterpriseId
	}
}
