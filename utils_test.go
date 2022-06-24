/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

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
