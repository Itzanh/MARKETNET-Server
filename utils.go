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
