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
}

type EnterpriseActivationRequest struct {
	EnterpriseKey string
	LicenseCode   string
	LicenseChance string
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
		ok = createNewEnterprise(newEnterprise.EnterpriseName, newEnterprise.EnterpriseDescription, newEnterprise.EnterpriseKey, newEnterprise.LicenseCode, newEnterprise.LicenseChance, newEnterprise.UserPassword)

	case "PUT":
		activateEnterprise := EnterpriseActivationRequest{}
		json.Unmarshal(body, &activateEnterprise)
		ok = activateEnterprise.reActivateEnterprise()

	case "DELETE":
		dactivateEnterprise := EnterpriseDeactivationRequest{}
		json.Unmarshal(body, &dactivateEnterprise)
		ok = deactivateEnterprise(dactivateEnterprise.EnterpriseKey)
	}

	json, _ := json.Marshal(ok)
	w.Write(json)
}

func (activateEnterprise *EnterpriseActivationRequest) reActivateEnterprise() bool {
	e := getSettingsRecordByEnterprise(activateEnterprise.EnterpriseKey)
	if e.Id <= 0 {
		return false
	}

	activation := ServerSettingsActivation{
		LicenseCode: activateEnterprise.LicenseCode,
		Chance:      &activateEnterprise.LicenseChance,
	}
	settings.Server.Activation[activateEnterprise.EnterpriseKey] = activation
	settings.setBackendSettings()
	return activation.activateEnterprise(e.Id)
}

func deactivateEnterprise(enterpriseKey string) bool {
	sqlStatement := `UPDATE config SET max_connections=0 WHERE enterprise_key=$1`
	res, err := db.Exec(sqlStatement, enterpriseKey)
	if err != nil {
		return false
	}
	rows, _ := res.RowsAffected()
	return rows > 0
}
