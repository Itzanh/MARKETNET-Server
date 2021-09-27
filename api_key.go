package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id                int32     `json:"id"`
	Name              string    `json:"name"`
	DateCreated       time.Time `json:"dateCreated"`
	UserCreated       int32     `json:"userCreated"`
	Off               bool      `json:"off"`
	User              int32     `json:"user"`
	Token             *string   `json:"token"`
	Auth              string    `json:"auth"` // P = Parameter, H = Header, B = Basic Auth, R = Bearer
	BasicAuthUser     *string   `json:"basicAuthUser"`
	BasicAuthPassword *string   `json:"basicAuthPassword"`
	UserCreatedName   string    `json:"userCreatedName"`
	UserName          string    `json:"userName"`
	enterprise        int32
}

func getApiKeys(enterpriseId int32) []ApiKey {
	keys := make([]ApiKey, 0)
	sqlStatement := `SELECT *,(SELECT username FROM "user" WHERE "user".id=api_key.user_created),(SELECT username FROM "user" WHERE "user".id=api_key."user") FROM public.api_key WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return keys
	}
	for rows.Next() {
		a := ApiKey{}
		rows.Scan(&a.Id, &a.Name, &a.DateCreated, &a.UserCreated, &a.Off, &a.User, &a.Token, &a.enterprise, &a.Auth, &a.BasicAuthUser, &a.BasicAuthPassword, &a.UserCreatedName, &a.UserName)
		keys = append(keys, a)
	}

	return keys
}

func (a *ApiKey) isValid() bool {
	return !(len(a.Name) == 0 || len(a.Name) > 64 || a.User <= 0 || (a.Auth != "P" && a.Auth != "H" && a.Auth != "B" && a.Auth != "R"))
}

func (a *ApiKey) insertApiKey() bool {
	if !a.isValid() {
		return false
	}

	if a.Auth == "P" || a.Auth == "H" || a.Auth == "R" {
		uuid := uuid.New().String()
		a.Token = &uuid
	} else if a.Auth == "B" {
		basicAuthUser := generateRandomString(20)
		a.BasicAuthUser = &basicAuthUser
		basicAuthPassword := generateRandomString(20)
		a.BasicAuthPassword = &basicAuthPassword
	}
	sqlStatement := `INSERT INTO public.api_key(name, user_created, "user", token, enterprise, auth, basic_auth_user, basic_auth_password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(sqlStatement, a.Name, a.UserCreated, a.User, a.Token, a.enterprise, a.Auth, a.BasicAuthUser, a.BasicAuthPassword)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func generateRandomString(length int) string {
	const CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890,.-;:_!@#$%&"
	str := ""

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < length; i++ {
		str += string(CHARSET[r.Intn(len(CHARSET))])
	}

	return str
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
// Checks for different authentication methods available on the API Keys
// OK, user id, enterprise id
func checkApiKey(r *http.Request) (bool, int32, int32) {
	// Header
	token := r.Header.Get("X-Marketnet-Access-Token")
	if len(token) == 36 {
		ok, userId, enterpriseId := checkApiKeyByTokenAuthType(token, "H")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId
		}
	}

	// Basic auth
	token = r.Header.Get("Authorization")
	basicAuth := strings.Split(token, " ")
	if len(basicAuth) == 2 && basicAuth[0] == "Basic" {
		basicAuth, _ := base64.StdEncoding.DecodeString(basicAuth[1])
		usernamePassword := strings.Split(string(basicAuth), ":")
		if len(usernamePassword) == 2 && len(usernamePassword[0]) == 20 && len(usernamePassword[1]) == 20 {
			ok, userId, enterpriseId := checkApiKeyByBasicAuthType(usernamePassword[0], usernamePassword[1])
			if ok && userId > 0 {
				return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId
			}
		}
	} else if len(basicAuth) == 2 && basicAuth[0] == "Bearer" {
		ok, userId, enterpriseId := checkApiKeyByTokenAuthType(basicAuth[1], "R")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId
		}
	}

	// Parameter
	tokenParam, ok := r.URL.Query()["token"]
	if ok && len(tokenParam) > 0 && len(tokenParam[0]) == 36 {
		ok, userId, enterpriseId := checkApiKeyByTokenAuthType(tokenParam[0], "P")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId
		}
	}

	// Givig up...
	return false, 0, 0
}

func checkMaxRequestsPerEnterprise(enterpriseId int32) bool {
	requestsMade, ok := requestsPerMinuteEnterprise[enterpriseId]
	if !ok {
		requestsPerMinuteEnterprise[enterpriseId] = 1
		return true
	}
	if requestsMade >= settings.Server.MaxRequestsPerMinuteEnterprise {
		return false
	} else {
		requestsPerMinuteEnterprise[enterpriseId] = requestsPerMinuteEnterprise[enterpriseId] + 1
		return true
	}
}

func resetMaxRequestsPerEnterprise() {
	requestsPerMinuteEnterprise = make(map[int32]int32)
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// P = Parameter, H = Header, B = Basic Auth
// OK, user id, enterprise id
func checkApiKeyByTokenAuthType(token string, auth string) (bool, int32, int32) {
	sqlStatement := `SELECT "user",enterprise FROM public.api_key WHERE off=false AND token=$1 AND auth=$2`
	row := db.QueryRow(sqlStatement, token, auth)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0, 0
	}

	var userId int32
	var enterpriseId int32
	row.Scan(&userId, &enterpriseId)

	return true, userId, enterpriseId
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// P = Parameter, H = Header, B = Basic Auth
// OK, user id, enterprise id
func checkApiKeyByBasicAuthType(basicAuthUser string, basicAuthPassword string) (bool, int32, int32) {
	sqlStatement := `SELECT "user",enterprise FROM public.api_key WHERE off=false AND basic_auth_user=$1 AND basic_auth_password=$2 AND auth='B'`
	row := db.QueryRow(sqlStatement, basicAuthUser, basicAuthPassword)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0, 0
	}

	var userId int32
	var enterpriseId int32
	row.Scan(&userId, &enterpriseId)

	return true, userId, enterpriseId
}
