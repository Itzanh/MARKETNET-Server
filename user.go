package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id                        int32     `json:"id" gorm:"index:user_id_enterprise,unique:true,priority:1"`
	Username                  string    `json:"username" gorm:"column:username;type:character varying(40);not null:true;index:user_username,unique:true,priority:2"`
	FullName                  string    `json:"fullName" gorm:"column:full_name;type:character varying(150);not null:true"`
	Email                     string    `json:"email" gorm:"column:email;type:character varying(100);not null:true"`
	DateCreated               time.Time `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	DateLastPwd               time.Time `json:"dateLastPwd" gorm:"type:timestamp(3) with time zone;not null:true"`
	PwdNextLogin              bool      `json:"pwdNextLogin" gorm:"column:pwd_next_login;not null:true"`
	Off                       bool      `json:"off" gorm:"column:off;not null:true"`
	Pwd                       []byte    `json:"-" gorm:"column:pwd;type:bytea;not null:true"`
	Salt                      string    `json:"-" gorm:"column:salt;type:character(30);not null:true"`
	Iterations                int32     `json:"iterations" gorm:"column:iterations;not null:true"`
	Description               string    `json:"description" gorm:"column:dsc;not null:true;type:text"`
	DateLastLogin             time.Time `json:"dateLastLogin" gorm:"type:timestamp(3) with time zone;not null:true"`
	FailedLoginAttemps        int16     `json:"-" gorm:"not null:true"`
	Language                  string    `json:"language" gorm:"column:lang;type:character(2);not null:true"`
	EnterpriseId              int32     `json:"-" gorm:"column:config;not null:true;index:user_id_enterprise,unique:true,priority:2;index:user_username,unique:true,priority:1"`
	Enterprise                Settings  `json:"-" gorm:"foreignKey:config;references:Id"`
	UsesGoogleAuthenticator   bool      `json:"usesGoogleAuthenticator" gorm:"not null:true"`
	GoogleAuthenticatorSecret *string   `json:"-" gorm:"type:character(8)"`
}

func (u *User) TableName() string {
	return "user"
}

func getUser(enterpriseId int32) []User {
	var users []User = make([]User, 0)
	dbOrm.Model(&User{}).Where("config = ?", enterpriseId).Order("id ASC").Find(&users)
	return users
}

func getUserByUsername(enterpriseId int32, username string) User {
	var user User
	dbOrm.Model(&User{}).Where("config = ? AND username = ?", enterpriseId, username).First(&user)
	return user
}

func getUserRow(userId int32) User {
	var user User
	dbOrm.Model(&User{}).Where("id = ?", userId).First(&user)
	return user
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	var user User
	tx.Model(&User{}).Last(&user)
	u.Id = user.Id + 1
	return nil
}

func (u *UserInsert) insertUser(enterpriseId int32) bool {
	if !u.isValid() {
		return false
	}

	salt := generateSalt()
	passwd := hashPassword([]byte(salt+u.Password), settings.Server.HashIterations)

	var user User = User{}
	user.Username = u.Username
	user.FullName = u.FullName
	user.Email = ""
	user.DateCreated = time.Now()
	user.DateLastPwd = time.Now()
	user.Pwd = passwd
	user.Salt = salt
	user.Iterations = settings.Server.HashIterations
	user.DateLastLogin = time.Now()
	user.Language = u.Language
	user.EnterpriseId = enterpriseId

	result := dbOrm.Create(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (u *User) isValid() bool {
	return !(len(u.Username) == 0 || len(u.Username) > 40 || len(u.FullName) > 150 || len(u.Email) > 100 || len(u.Language) != 2 || len(u.Description) > 3000)
}

func (u *User) updateUser() bool {
	if u.Id <= 0 || !u.isValid() {
		return false
	}

	var user User
	result := dbOrm.Model(&User{}).Where("id = ? AND config = ?", u.Id, u.EnterpriseId).First(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	user.Username = u.Username
	user.FullName = u.FullName
	user.Email = u.Email
	user.Description = u.Description
	user.Language = u.Language

	result = dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (u *User) deleteUser() bool {
	if u.Id <= 0 {
		return false
	}

	user := getUserRow(u.Id)
	if user.Id <= 0 || user.EnterpriseId != u.EnterpriseId {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	result := trans.Where(`"user" = ?`, u.Id).Delete(&UserGroup{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	result = trans.Where("id = ? AND config = ?", u.Id, u.EnterpriseId).Delete(&User{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	trans.Commit()
	///
	return true
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

	var user User
	result := dbOrm.Model(&User{}).Where("id = ? AND config = ?", u.Id, enterpriseId).First(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	user.DateLastPwd = time.Now()
	user.Pwd = passwd
	user.Salt = salt
	user.Iterations = settings.Server.HashIterations
	user.PwdNextLogin = u.PwdNextLogin
	user.FailedLoginAttemps = 0

	result = dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type UserAutoPassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Function used by the user to change its own password.
func (u *UserAutoPassword) userAutoPassword(enterpriseId int32, userId int32) bool {
	if u.CurrentPassword == u.NewPassword {
		return false
	}

	passwordEvaluation := evaluatePasswordSecureCloud(enterpriseId, u.NewPassword)
	if (!passwordEvaluation.PasswordComplexity) || passwordEvaluation.PasswordInBlacklist || passwordEvaluation.PasswordHashInBlacklist {
		return false
	}

	user := getUserRow(userId)
	if user.Id <= 0 || user.EnterpriseId != enterpriseId {
		return false
	}
	oldPasswd := hashPassword([]byte(user.Salt+u.CurrentPassword), user.Iterations)

	if !comparePasswords(oldPasswd, user.Pwd) {
		return false
	}

	salt := generateSalt()
	passwd := hashPassword([]byte(salt+u.NewPassword), settings.Server.HashIterations)

	user.DateLastPwd = time.Now()
	user.Pwd = passwd
	user.Salt = salt
	user.Iterations = settings.Server.HashIterations
	user.PwdNextLogin = false
	user.FailedLoginAttemps = 0

	result := dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func checkForPasswordComplexity(enterpriseId int32, password string) bool {
	enterprise := getSettingsRecordById(enterpriseId)

	if len(password) < int(enterprise.PasswordMinimumLength) {
		return false
	}

	const ALPHA = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const NUMBERS = "0123456789"
	const UPPERCASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const LOWERCASE = "abcdefghijklmnopqrstuvwxyz"
	const ALL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	if enterprise.PasswordMinumumComplexity == "A" {
		if !characterSetFound(password, ALPHA) {
			return false
		}
	} else if enterprise.PasswordMinumumComplexity == "B" {
		if (!characterSetFound(password, ALPHA)) || (!characterSetFound(password, NUMBERS)) {
			return false
		}
	} else if enterprise.PasswordMinumumComplexity == "C" {
		if (!characterSetFound(password, UPPERCASE)) || (!characterSetFound(password, LOWERCASE)) || (!characterSetFound(password, NUMBERS)) {
			return false
		}
	} else if enterprise.PasswordMinumumComplexity == "D" {
		if (!characterSetFound(password, UPPERCASE)) || (!characterSetFound(password, LOWERCASE)) || (!characterSetFound(password, NUMBERS)) || (!characterSetNotFound(password, ALL)) {
			return false
		}
	} else {
		return false
	}
	return true
}

func characterSetFound(input string, characterSet string) bool {
	var characterSetFound = false
	for i := 0; i < len(characterSet); i++ {
		if strings.Contains(input, string(characterSet[i])) {
			characterSetFound = true
			break
		}
	}

	return characterSetFound
}

func characterSetNotFound(input string, characterSet string) bool {
	var characterSetFound = true
	for i := 0; i < len(characterSet); i++ {
		if !strings.Contains(input, string(characterSet[i])) {
			characterSetFound = false
			break
		}
	}

	return characterSetFound
}

func (u *User) offUser() bool {
	if u.Id <= 0 {
		return false
	}

	var user User
	result := dbOrm.Model(&User{}).Where("id = ? AND config = ?", u.Id, u.EnterpriseId).First(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	user.Off = !user.Off

	result = dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type UserLogin struct {
	Enterprise string `json:"enterprise"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
}

type UserLoginResult struct {
	Ok                  bool         `json:"ok"`
	Token               string       `json:"token"`
	Permissions         *Permissions `json:"permissions"`
	Language            string       `json:"language"`
	GoogleAuthenticator bool         `json:"googleAuthenticator"`
	Reason              uint8        `json:"reason"` // 0 = Incorrect login, 1 = Connection filtered, 2 = Maximum number of connections reached
	ExtraData           []string     `json:"extraData"`
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
		t := LoginToken{UserId: user.Id, IpAddress: ipAddress}
		t.insertLoginToken()
		perm := getUserPermissions(user.Id, enterprise.Id)
		return UserLoginResult{Ok: true, Token: t.Name, Permissions: &perm, Language: user.Language}, user.Id, enterprise.Id
	} else { // the two arrays are different
		user.setUserFailedLoginAttemps(true)
		return UserLoginResult{Ok: false}, 0, 0
	}
}

func (r *UserLoginResult) checkUserConnection(userId int32, remoteAddr string, enterpriseId int32) bool {
	okFilter, filterName := userConnection(userId, remoteAddr, enterpriseId)
	if !okFilter {
		r.reset()
		r.Reason = 1
		r.ExtraData = []string{filterName}
		return false
	}
	s := getSettingsRecordById(enterpriseId)
	if len(getConnections(enterpriseId)) >= int(s.MaxConnections) {
		r.reset()
		r.Reason = 2
		r.ExtraData = []string{strconv.Itoa(int(s.MaxConnections))}
		return false
	}
	return true
}

func (r *UserLoginResult) reset() {
	r.Ok = false
	r.Token = ""
	r.Permissions = nil
	r.Language = ""
	r.GoogleAuthenticator = false
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
	user := getUserRow(u.Id)

	if addOrReset {
		user.FailedLoginAttemps = user.FailedLoginAttemps + 1
	} else {
		user.FailedLoginAttemps = 0
	}

	result := dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func setUserDateLastLogin(userId int32) {
	user := getUserRow(userId)

	user.DateLastLogin = time.Now()

	result := dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return
	}
}

type SecureCloudResult struct {
	PasswordComplexity      bool `json:"passwordComplexity"`
	PasswordInBlacklist     bool `json:"passwordInBlacklist"`
	PasswordHashInBlacklist bool `json:"passwordHashInBlacklist"`
}

func evaluatePasswordSecureCloud(enterpriseId int32, password string) SecureCloudResult {
	result := SecureCloudResult{}
	result.PasswordComplexity = checkForPasswordComplexity(enterpriseId, password)
	result.PasswordInBlacklist = searchPasswordInBlackList(password)
	result.PasswordHashInBlacklist = searchPasswordHashInBlackList(password)
	return result
}

func searchPasswordInBlackList(password string) bool {
	var rowsFound int64
	result := dbOrm.Model(&PwdBlacklist{}).Where("pwd = ?", password).Count(&rowsFound)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return rowsFound > 0
}

func searchPasswordHashInBlackList(password string) bool {
	hasher := sha1.New()
	var pwd []byte = []byte(password)

	hasher.Write(pwd)
	pwd = hasher.Sum(nil)

	var rowsFound int64
	result := dbOrm.Model(&PwdSHA1Blacklist{}).Where("hash = ?", pwd).Count(&rowsFound)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return rowsFound > 0
}

func addPasswordsToBlacklist() {
	pwdBlackList, ok := getParameterValue("pwd-blacklist-file")
	if ok {
		if len(pwdBlackList) == 0 {
			fmt.Println("You must specify a .txt file path")
			return
		}

		insertPwdBlacklist(pwdBlackList)
		return
	}

	pwdHashBlackList, ok := getParameterValue("pwd-hash-blacklist-file")
	if ok {
		if len(pwdHashBlackList) == 0 {
			fmt.Println("You must specify a .txt file path")
			return
		}

		insertPwdBlacklistHash(pwdHashBlackList)
		return
	}

	singlePwd, ok := getParameterValue("single-pwd")
	if ok {
		if len(singlePwd) == 0 {
			fmt.Println("You must specify a password")
			return
		}

		insertSinglePwdBlacklistHash(singlePwd)
		return
	}

	fmt.Println("Option not recognised")
}

type PwdBlacklist struct {
	Pwd string `json:"pwd" gorm:"primaryKey;column:pwd;type:character varying(255)"`
}

func (p *PwdBlacklist) TableName() string {
	return "pwd_blacklist"
}

// inserts an entire password dictionary in txt format in the blacklist
// format is a list of passwords separated by a new line
// example: "C:\\Users\\Itzan\\Desktop\\rockyou.txt"
func insertPwdBlacklist(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var password string
	var result *gorm.DB
	for scanner.Scan() {
		password = scanner.Text()

		result = dbOrm.Create(&PwdBlacklist{Pwd: password})
		if result.Error != nil {
			log("DB", result.Error.Error())
			return
		}
	}
}

type PwdSHA1Blacklist struct {
	Hash []byte `json:"hash" gorm:"primaryKey;column:hash;type:bytea"`
}

func (p *PwdSHA1Blacklist) TableName() string {
	return "pwd_sha1_blacklist"
}

// inserts a list of SHA-1 hashes in the hashed passwords blacklist
// format is a list of SHA-1 hashes separated by a new line
// example: "C:\\Users\\Itzan\\Desktop\\pwned-passwords-sha1-ordered-by-count-v7.txt"
func insertPwdBlacklistHash(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var password string
	var hash []byte
	var result *gorm.DB
	for scanner.Scan() {
		password = scanner.Text()[:40]
		hash, _ = hex.DecodeString(password)
		result = dbOrm.Create(&PwdSHA1Blacklist{Hash: hash})
		if result.Error != nil {
			log("DB", result.Error.Error())
			return
		}
	}
}

func insertSinglePwdBlacklistHash(pass string) {
	hasher := sha1.New()

	hasher.Write([]byte(pass))
	hash := hasher.Sum(nil)

	result := dbOrm.Create(&PwdSHA1Blacklist{Hash: hash})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return
	}
}
