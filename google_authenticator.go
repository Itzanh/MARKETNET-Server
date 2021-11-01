package main

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/dgryski/dgoogauth"
)

func randStr(strSize int) string {
	const dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

type GoogleAuthenticatorRegister struct {
	Ok       bool   `json:"ok"`
	AuthLink string `json:"authLink"`
}

func registerUserInGoogleAuthenticator(userId int32, enterpriseId int32) GoogleAuthenticatorRegister {
	user := getUserRow(userId)
	if user.Id <= 0 || user.enterprise != enterpriseId {
		return GoogleAuthenticatorRegister{Ok: false, AuthLink: ""}
	}
	enterprise := getSettingsRecordById(user.enterprise)
	if enterprise.Id <= 0 {
		return GoogleAuthenticatorRegister{Ok: false, AuthLink: ""}
	}
	// generate a random string - preferbly 6 or 8 characters
	randomStr := randStr(8)

	sqlStatement := `UPDATE public."user" SET uses_google_authenticator=true, google_authenticator_secret=$3 WHERE id=$1 AND config=$2`
	_, err := db.Exec(sqlStatement, user.Id, enterpriseId, randomStr)
	if err != nil {
		return GoogleAuthenticatorRegister{Ok: false, AuthLink: ""}
	}

	// For Google Authenticator purpose for more details see https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	secret := base32.StdEncoding.EncodeToString([]byte(randomStr))

	// Authentication link. For more details see https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	authLink := "otpauth://totp/" + user.Username + "@" + enterprise.EnterpriseKey + "?secret=" + secret + "&issuer=MARKETNET"
	return GoogleAuthenticatorRegister{Ok: true, AuthLink: authLink}
}

func authenticateUserInGoogleAuthenticator(userId int32, enterpriseId int32, token string) bool {
	user := getUserRow(userId)
	if user.Id <= 0 || user.enterprise != enterpriseId {
		return false
	}

	// For Google Authenticator purpose for more details see https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	secret := base32.StdEncoding.EncodeToString([]byte(user.googleAuthenticatorSecret))

	// setup the one-time-password configuration.
	otpConfig := &dgoogauth.OTPConfig{
		Secret:      strings.TrimSpace(secret),
		WindowSize:  3,
		HotpCounter: 0,
	}

	// get rid of the extra \n from the token string
	// otherwise the validation will fail
	trimmedToken := strings.TrimSpace(token)

	// Validate token
	ok, err := otpConfig.Authenticate(trimmedToken)

	if err != nil {
		fmt.Println(err)
	}

	return ok
}

func removeUserFromGoogleAuthenticator(userId int32, enterpriseId int32) bool {
	sqlStatement := `UPDATE public."user" SET uses_google_authenticator=false, google_authenticator_secret=NULL WHERE id=$1 AND config=$2`
	_, err := db.Exec(sqlStatement, userId, enterpriseId)
	return err == nil
}
