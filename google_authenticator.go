/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

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
	if user.Id <= 0 || user.EnterpriseId != enterpriseId {
		return GoogleAuthenticatorRegister{Ok: false, AuthLink: ""}
	}
	enterprise := getSettingsRecordById(user.EnterpriseId)
	if enterprise.Id <= 0 {
		return GoogleAuthenticatorRegister{Ok: false, AuthLink: ""}
	}
	// generate a random string - preferbly 6 or 8 characters
	randomStr := randStr(8)

	user.UsesGoogleAuthenticator = true
	user.GoogleAuthenticatorSecret = &randomStr

	result := dbOrm.Save(&user)
	if result.Error != nil {
		log("DB", result.Error.Error())
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
	if user.Id <= 0 || user.EnterpriseId != enterpriseId || user.GoogleAuthenticatorSecret == nil {
		return false
	}

	// For Google Authenticator purpose for more details see https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	googleAuthenticatorSecret := *user.GoogleAuthenticatorSecret
	secret := base32.StdEncoding.EncodeToString([]byte(googleAuthenticatorSecret))

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
	user := getUserRow(userId)

	user.UsesGoogleAuthenticator = false
	user.GoogleAuthenticatorSecret = nil

	result := dbOrm.Updates(user)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return result.RowsAffected > 0
}
