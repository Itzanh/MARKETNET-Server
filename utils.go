/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func emailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func phoneIsValid(phone string) bool {
	const VALID_CHARACTERS = "0123456789()-+. "
	if len(phone) == 0 {
		return false
	}
	for i := 0; i < len(phone); i++ {
		if !strings.Contains(VALID_CHARACTERS, string(phone[i])) {
			return false
		}
	}
	return true
}

// Abs returns the absolute value of x.
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// Abs returns the absolute value of x.
func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func isParameterPresent(parameter string) bool {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == parameter {
			return true
		}
	}
	return false
}

func getParameterValue(parameter string) (string, bool) {
	for i := 1; i < len(os.Args); i++ {
		parameterValue := strings.Split(os.Args[i], "=")
		if len(parameterValue) == 2 && parameterValue[0] == parameter {
			return parameterValue[1], true
		}
	}
	return "", false
}

func hostnameWithPortValid(hostname string) bool {
	colonCount := strings.Count(hostname, ":")
	if colonCount != 1 {
		return false
	}

	colonIndexOf := strings.Index(hostname, ":")
	if colonIndexOf < 1 {
		return false
	}

	host := hostname[:colonIndexOf]
	if len(host) == 0 {
		return false
	}

	port := hostname[colonIndexOf+1:]
	if len(port) == 0 {
		return false
	}

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return false
	}
	return true
}

// Check if the bar code is a valid EAN13 product code. (Check the verification digit)
func checkEan13(barcode string) bool {
	if len(barcode) != 13 {
		return false
	}
	// barcode must be a number
	_, err := strconv.Atoi(barcode)
	if err != nil {
		return false
	}

	// get the first 12 digits (remove the 13 character, which is the control digit), and reverse the string
	barcode12 := barcode[0:12]
	barcode12 = Reverse(barcode12)

	// add the numbers in the odd positions
	var controlNumber uint16
	for i := 0; i < len(barcode12); i += 2 {
		digit, _ := strconv.Atoi(string(barcode12[i]))
		controlNumber += uint16(digit)
	}

	// multiply by 3
	controlNumber *= 3

	// add the numbers in the pair positions
	for i := 1; i < len(barcode12); i += 2 {
		digit, _ := strconv.Atoi(string(barcode12[i]))
		controlNumber += uint16(digit)
	}

	// immediately higher ten
	var controlDigit uint16 = (10 - (controlNumber % 10)) % 10

	// check the control digits are the same
	inputControlDigit, _ := strconv.Atoi(string(barcode[12]))
	return controlDigit == uint16(inputControlDigit)
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func checkUUID(Uuid string) bool {
	_, err := uuid.Parse(Uuid)
	return err == nil
}

type OkAndErrorCodeReturn struct {
	Ok        bool     `json:"ok"`
	ErrorCode uint8    `json:"errorCode"`
	ExtraData []string `json:"extraData"`
}

// https://www.socketloop.com/tutorials/golang-chunk-split-or-divide-a-string-into-smaller-chunk-example
// Does the same as the chunk_split PHP funtion
// Formats data for the RFC 2045 semantics
func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}
	var result string = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		// charSlice
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}

	return result
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func stringArrayToString(stringArray []string) string {
	jsonData, _ := json.Marshal(stringArray)
	return string(jsonData)
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func checkBase64(base64String string) bool {
	_, err := base64.StdEncoding.DecodeString(base64String)
	return err == nil
}

func base64ToUuid(base64String string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	uuid, err := uuid.FromBytes(decoded)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func checkHex(hexString string) bool {
	_, err := hex.DecodeString(hexString)
	return err == nil
}

func checkUrl(urlString string) bool {
	_, err := url.Parse(urlString)
	return err == nil
}
