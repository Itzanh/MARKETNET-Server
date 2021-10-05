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

/*
func gAuth() {

	// maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// generate a random string - preferbly 6 or 8 characters

	randomStr := randStr(6)
	fmt.Println(randomStr)

	// For Google Authenticator purpose
	// for more details see
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	secret := base32.StdEncoding.EncodeToString([]byte(randomStr))
	fmt.Println(secret)

	// authentication link. Remember to replace SocketLoop with yours.
	// for more details see
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	authLink := "otpauth://totp/itzan@MARKETNET?secret=" + secret + "&issuer=MARKETNET"

	// Encode authLink to QR codes
	// qr.H = 65% redundant level
	// see https://godoc.org/code.google.com/p/rsc/qr#Level

	// otpauth://totp/itzan@MARKETNET?secret=JRWDERKFGE======&issuer=MARKETNET
	fmt.Println(authLink)
	/*code, err := qr.Encode(authLink, qr.L)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	imgByte := code.PNG()

	// convert byte to image for saving to file
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	err = imaging.Save(img, "./QRImgGA.png")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
*/
/*// everything ok
	fmt.Println("QR code generated and saved to QRimgGA.png. Please scan the QRImgGA.png with Google Authenticator App.")
	fmt.Println("NOTE : You need to remove the old entry for SocketLoop in Google Authenticator App each time.")

	tokenReader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter token to verify : ")

	// prompt user for input
	token, err := tokenReader.ReadString('\n')

	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}

	fmt.Println("Token : ", token)

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
		os.Exit(1)
	}

	fmt.Printf("Token string [%s] validation is : %v \n", trimmedToken, ok)

	fmt.Println("IMPORTANT : Once the user token is validated. Store the secret string into")
	fmt.Println("database and memory. Use the secret string associated with this")
	fmt.Println("user whenever 2FA is required.")
	fmt.Println("If the user decides to disable 2FA, remove the secret string from")
	fmt.Println("database and memory. Generate a new secret string when user")
	fmt.Println("re-enable 2FA.")

}
*/
