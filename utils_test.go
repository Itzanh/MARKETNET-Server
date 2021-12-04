package main

import "testing"

func TestEmailIsValid(t *testing.T) {
	if !emailIsValid("user@enterprise.com") {
		t.Error("The email should be OK")
	}

	if !emailIsValid("user.surname@enterprise.com") {
		t.Error("The email should be OK")
	}

	if emailIsValid("@enterprise.com") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("enterprise.com") {
		t.Error("The email should not be OK")
	}

	if !emailIsValid("user@enterprise") {
		t.Error("The email should be OK")
	}

	if emailIsValid("user@enterprise.") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("user@.com") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("user@.") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("user@") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("@") {
		t.Error("The email should not be OK")
	}

	if emailIsValid("") {
		t.Error("The email should not be OK")
	}
}

func TestIsParameterPresent(t *testing.T) {
	if isParameterPresent("totallyNotCallableParameter") {
		t.Error("The parameter should not be OK")
	}
}

func TestGetParameterValue(t *testing.T) {
	str, ok := getParameterValue("totallyNotCallableParameter")
	if str != "" || ok {
		t.Error("Parameter value should not be OK")
	}
}
