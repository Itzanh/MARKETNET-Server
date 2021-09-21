package main

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"dateCreated"`
	UserCreated int32     `json:"userCreated"`
	Off         bool      `json:"off"`
	User        int32     `json:"user"`
	Token       string    `json:"token"`
	enterprise  int32
}

func getApiKeys(enterpriseId int32) []ApiKey {
	keys := make([]ApiKey, 0)
	sqlStatement := `SELECT * FROM public.api_key WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return keys
	}
	for rows.Next() {
		a := ApiKey{}
		rows.Scan(&a.Id, &a.Name, &a.DateCreated, &a.UserCreated, &a.Off, &a.User, &a.Token, &a.enterprise)
		keys = append(keys, a)
	}

	return keys
}

func (a *ApiKey) isValid() bool {
	return !(len(a.Name) == 0 || len(a.Name) > 64 || a.User <= 0)
}

func (a *ApiKey) insertApiKey() bool {
	if !a.isValid() {
		return false
	}

	a.Token = uuid.New().String()
	sqlStatement := `INSERT INTO public.api_key(name, user_created, "user", token, enterprise) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, a.Name, a.UserCreated, a.User, a.Token, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (a *ApiKey) deleteApiKey() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.api_key WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

func (a *ApiKey) offApiKey() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.api_key SET off=NOT off WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// OK, user id, enterprise id
func checkApiKey(token string) (bool, int32, int32) {

	sqlStatement := `SELECT "user",enterprise FROM public.api_key WHERE off=false AND token=$1`
	row := db.QueryRow(sqlStatement, token)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0, 0
	}

	var userId int32
	var enterpriseId int32
	row.Scan(&userId, &enterpriseId)

	return true, userId, 0
}
