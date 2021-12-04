package main

import (
	"net/mail"
	"os"
	"strconv"
	"strings"
)

func emailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
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