/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type NewEnterpriseRequest struct {
	EnterpriseKey         string
	EnterpriseName        string
	EnterpriseDescription string
	UserPassword          string
	LicenseCode           string
	LicenseChance         string
	DocumentSpace         float64
}

type EnterpriseActivationRequest struct {
	EnterpriseKey string
	LicenseCode   string
	LicenseChance string
	DocumentSpace float64
}

type EnterpriseDeactivationRequest struct {
	EnterpriseKey string
}

func handleEnterprise(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")

	// Header
	token := r.Header.Get("X-Marketnet-Access-Token")
	if token != settings.Server.SaaSAccessToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var ok bool

	switch r.Method {
	case "POST":
		newEnterprise := NewEnterpriseRequest{}
		json.Unmarshal(body, &newEnterprise)
		ok = createNewEnterprise(newEnterprise.EnterpriseName, newEnterprise.EnterpriseDescription, newEnterprise.EnterpriseKey, newEnterprise.LicenseCode, newEnterprise.LicenseChance, newEnterprise.UserPassword, newEnterprise.DocumentSpace)
	}

	json, _ := json.Marshal(ok)
	w.Write(json)
}
