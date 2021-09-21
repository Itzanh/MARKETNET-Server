package main

import (
	"crypto/sha512"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	Id                 int32     `json:"id"`
	Username           string    `json:"username"`
	FullName           string    `json:"fullName"`
	Email              string    `json:"email"`
	DateCreated        time.Time `json:"dateCreated"`
	DateLastPwd        time.Time `json:"dateLastPwd"`
	PwdNextLogin       bool      `json:"pwdNextLogin"`
	Off                bool      `json:"off"`
	Pwd                []byte
	Salt               string
	Iterations         int32     `json:"iterations"`
	Description        string    `json:"description"`
	DateLastLogin      time.Time `json:"dateLastLogin"`
	FailedLoginAttemps int16
	Language           string `json:"language"`
	enterprise         int32
}

func getUser(enterpriseId int32) []User {
	var users []User = make([]User, 0)
	sqlStatement := `SELECT * FROM "user" WHERE config=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return users
	}
	for rows.Next() {
		u := User{}
		rows.Scan(&u.Id, &u.Username, &u.FullName, &u.Email, &u.DateCreated, &u.DateLastPwd, &u.PwdNextLogin, &u.Off, &u.Pwd, &u.Salt, &u.Iterations, &u.Description, &u.DateLastLogin, &u.FailedLoginAttemps, &u.Language, &u.enterprise)
		users = append(users, u)
	}

	return users
}

func getUserByUsername(enterpriseId int32, username string) User {
	sqlStatement := `SELECT * FROM "user" WHERE config=$1 AND username=$2`
	row := db.QueryRow(sqlStatement, enterpriseId, username)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return User{}
	}

	u := User{}
	row.Scan(&u.Id, &u.Username, &u.FullName, &u.Email, &u.DateCreated, &u.DateLastPwd, &u.PwdNextLogin, &u.Off, &u.Pwd, &u.Salt, &u.Iterations, &u.Description, &u.DateLastLogin, &u.FailedLoginAttemps, &u.Language, &u.enterprise)

	return u
}

func getUserRow(userId int32) User {
	sqlStatement := `SELECT * FROM "user" WHERE id=$1`
	row := db.QueryRow(sqlStatement, userId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return User{}
	}

	u := User{}
	row.Scan(&u.Id, &u.Username, &u.FullName, &u.Email, &u.DateCreated, &u.DateLastPwd, &u.PwdNextLogin, &u.Off, &u.Pwd, &u.Salt, &u.Iterations, &u.Description, &u.DateLastLogin, &u.FailedLoginAttemps, &u.Language, &u.enterprise)

	return u
}

type UserInsert struct {
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	Language string `json:"language"`
}

func (u *UserInsert) isValid() bool {
	return !(len(u.Username) == 0 || len(u.Username) > 40 || len(u.Password) < 8 || len(u.Language) != 2)
}

func generateSalt() string {
	const CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890,.-;:_!@#$%&"
	salt := ""

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 30; i++ {
		salt += string(CHARSET[r.Intn(len(CHARSET))])
	}

	return salt
}

func hashPassword(password []byte, iterations int32) []byte {
	hasher := sha512.New()
	var pwd []byte = password

	var i int32
	for i = 0; i < iterations; i++ {
		hasher.Write(pwd)
		pwd = hasher.Sum(nil)
	}

	return pwd
}

func (u *UserInsert) insertUser(enterpriseId int32) bool {
	if !u.isValid() {
		return false
	}

	salt := generateSalt()
	passwd := hashPassword([]byte(salt+u.Password), settings.Server.HashIterations)

	sqlStatement := `INSERT INTO public."user"(username, full_name, pwd, salt, iterations, lang, config) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, u.Username, u.FullName, passwd, salt, settings.Server.HashIterations, u.Language, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (u *User) isValid() bool {
	return !(len(u.Username) == 0 || len(u.Username) > 40 || len(u.FullName) > 150 || len(u.Email) > 100 || len(u.Language) != 2)
}

func (u *User) updateUser() bool {
	if u.Id <= 0 || !u.isValid() {
		return false
	}

	sqlStatement := `UPDATE public."user" SET username=$2, full_name=$3, email=$4, dsc=$5, lang=$6 WHERE id=$1 AND config=$7`
	res, err := db.Exec(sqlStatement, u.Id, u.Username, u.FullName, u.Email, u.Description, u.Language, u.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (u *User) deleteUser() bool {
	if u.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public."user" WHERE id=$1 AND config=$2`
	res, err := db.Exec(sqlStatement, u.Id, u.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type UserPassword struct {
	Id           int32  `json:"id"`
	Password     string `json:"password"`
	PwdNextLogin bool   `json:"pwdNextLogin"`
}

func (u *UserPassword) userPassword(enterpriseId int32) bool {
	if u.Id <= 0 || len(u.Password) < 8 {
		return false
	}

	salt := generateSalt()
	passwd := hashPassword([]byte(salt+u.Password), settings.Server.HashIterations)

	sqlStatement := `UPDATE public."user" SET date_last_pwd=CURRENT_TIMESTAMP(3), pwd=$2, salt=$3, iterations=$4, pwd_next_login=$5 WHERE id=$1 AND config=$6`
	res, err := db.Exec(sqlStatement, u.Id, passwd, salt, settings.Server.HashIterations, u.PwdNextLogin, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (u *User) offUser() bool {
	if u.Id <= 0 {
		return false
	}
	sqlStatement := `UPDATE public."user" SET off = NOT off WHERE id=$1 AND config=$2`
	res, err := db.Exec(sqlStatement, u.Id, u.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type UserLogin struct {
	Enterprise string `json:"enterprise"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
}

type UserLoginResult struct {
	Ok          bool         `json:"ok"`
	Token       string       `json:"token"`
	Permissions *Permissions `json:"permissions"`
	Language    string       `json:"language"`
}

// Result, user id, enterprise id
func (u *UserLogin) login(ipAddress string) (UserLoginResult, int32, int32) {
	if len(u.Username) == 0 || len(u.Username) > 50 || len(u.Password) < 8 {
		return UserLoginResult{Ok: false}, 0, 0
	}

	enterprise := getSettingsRecordByEnterprise(strings.ToUpper(u.Enterprise))
	if enterprise.Id <= 0 {
		return UserLoginResult{Ok: false}, 0, 0
	}

	user := getUserByUsername(enterprise.Id, u.Username)
	if user.Id <= 0 || user.Off || user.FailedLoginAttemps >= settings.Server.MaxLoginAttemps {
		return UserLoginResult{Ok: false}, 0, 0
	}

	passwd := hashPassword([]byte(user.Salt+u.Password), user.Iterations)

	if comparePasswords(passwd, user.Pwd) {
		user.setUserFailedLoginAttemps(false)
		t := LoginToken{User: user.Id, IpAddress: ipAddress}
		t.insertLoginToken()
		perm := getUserPermissions(user.Id, enterprise.Id)
		return UserLoginResult{Ok: true, Token: t.Name, Permissions: &perm, Language: user.Language}, user.Id, enterprise.Id
	} else { // the two arrays are different
		user.setUserFailedLoginAttemps(true)
		return UserLoginResult{Ok: false}, 0, 0
	}
}

func comparePasswords(passwordInput []byte, passwordOutput []byte) bool {
	if len(passwordInput) != len(passwordOutput) {
		return false
	}

	for i := 0; i < len(passwordInput); i++ {
		if passwordInput[i] != passwordOutput[i] {
			return false
		}
	}

	return true
}

// Adds or resets the amounts of failed login attemps for one user.
// addOrReset = true: Add one failed attemps
// addOrReset = false: Resets the failed attemps
func (u *User) setUserFailedLoginAttemps(addOrReset bool) bool {
	sqlStatement := ""
	if addOrReset {
		sqlStatement = `UPDATE "user" SET failed_login_attemps=failed_login_attemps+1 WHERE id=$1`
	} else {
		sqlStatement = `UPDATE "user" SET failed_login_attemps=0 WHERE id=$1`
	}

	res, err := db.Exec(sqlStatement, u.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
