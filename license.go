package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	CHANCE_URL   = "https://license.marketneterp.io:12278/chance"
	ACTIVATE_URL = "https://license.marketneterp.io:12278/activate"
)

// The maximum limit of concurrent connections limited by the adquired license.
// key: int32: enterpriseId
// value: int16: maximum number of connections
var licenseMaxConnections map[int32]int16 = make(map[int32]int16)

// Attempt product activation by license code. If the activation fails, the server will shut down.
// This function will be called in a new thread when the server is started.
func activate() {
	// The license list can't be empty
	if len(settings.Server.Activation) == 0 {
		fmt.Println("There are no license codes in the config file.")
		fmt.Println("The application could not be activated by license and will shut down")
		os.Exit(2)
	}

	// There can't be duplicated license codes
	licenseCodes := make([]string, 0)
	for _, activation := range settings.Server.Activation {
		for i := 0; i < len(licenseCodes); i++ {
			if licenseCodes[i] == activation.LicenseCode {
				fmt.Println("There can't be duplicated license codes in the config file.")
				fmt.Println("The application could not be activated by license and will shut down")
				os.Exit(2)
			}
		}

		licenseCodes = append(licenseCodes, activation.LicenseCode)
	}

	// The enterprises that does not appear on the activation map must have 0 connections
	settingsRecords := getSettingsRecords()
	for enterpriseKey := range settings.Server.Activation {
		s := getSettingsRecordByEnterprise(enterpriseKey)
		if s.Id <= 0 {
			fmt.Println("Could not find enterprise with name " + enterpriseKey + ".")
			fmt.Println("The application could not be activated by license and will shut down")
			os.Exit(2)
		}

		// Remove the enterprises that are in the activation map from the array
		for i := 0; i < len(settingsRecords); i++ {
			if settingsRecords[i].Id == s.Id {
				settingsRecords = append(settingsRecords[:i], settingsRecords[i+1:]...)
			}
		}
	}
	// The remaining enterprises in the array are not in the activation map, set 0 as the maximum connections.
	for i := 0; i < len(settingsRecords); i++ {
		sqlStatement := `UPDATE config SET max_connections=0 WHERE id=$1`
		db.Exec(sqlStatement, settingsRecords[i].Id)
	}

	// Check the activation for each license of each enterprise
	for enterpriseKey, activation := range settings.Server.Activation {
		s := getSettingsRecordByEnterprise(enterpriseKey)
		if s.Id <= 0 {
			fmt.Println("Could not find enterprise with name " + enterpriseKey + ".")
			fmt.Println("The application could not be activated by license and will shut down")
			os.Exit(2)
		}

		// The license code must be a valid UUID
		_, err := uuid.Parse(activation.LicenseCode)
		if err != nil {
			fmt.Println("The license code in the config file is not a valid UUID.")
			fmt.Println("The application could not be activated by license and will shut down")
			os.Exit(2)
		}

		// If both the chance and the secret is null, the product can't be activated
		if activation.Chance == nil && (activation.Secret == nil || activation.InstallId == nil) {
			fmt.Println("There is no chance or secret in the config file")
			fmt.Println("The application could not be activated by license and will shut down")
			os.Exit(2)
		}

		// If there is a chance, activate the license of the product.
		if activation.Chance != nil {
			if !takeActivationChance(enterpriseKey, activation.LicenseCode, *activation.Chance) {
				fmt.Println("The application could not be activated by license and will shut down")
				os.Exit(2)
			}
		}

		// If there is an activation secret, check if it's correct
		if activation.Secret != nil && activation.InstallId != nil {
			if !checkActivation(enterpriseKey, s.Id, activation.LicenseCode, *activation.Secret, *activation.InstallId) {
				fmt.Println("This product is not activated")
				fmt.Println("The application could not be activated by license and will shut down")
				os.Exit(2)
			}
		} else {
			// Can't continue
			fmt.Println("There must be a secret and a installation ID in the config file to check the activation this product.")
			fmt.Println("The application could not be activated by license and will shut down")
			os.Exit(2)
		}
	}
}

// Attempt product activation by license code. If the activation fails, the server will shut down.
// This function will be called in a new thread when the server is started.
func (activation *ServerSettingsActivation) activateEnterprise(enterpriseId int32) bool {
	// There can't be duplicated license codes
	licenseCodes := make([]string, 0)
	for _, activation := range settings.Server.Activation {
		for i := 0; i < len(licenseCodes); i++ {
			if licenseCodes[i] == activation.LicenseCode {
				fmt.Println("There can't be duplicated license codes in the config file.")
				fmt.Println("The application could not be activated by license and will shut down")
				return false
			}
		}

		licenseCodes = append(licenseCodes, activation.LicenseCode)
	}

	// Check the activation for this license of this enterprise
	s := getSettingsRecordById(enterpriseId)
	if s.Id <= 0 {
		fmt.Println("Could not find enterprise with name " + s.EnterpriseKey + ".")
		fmt.Println("The application could not be activated by license and will shut down")
		return false
	}

	// The license code must be a valid UUID
	_, err := uuid.Parse(activation.LicenseCode)
	if err != nil {
		fmt.Println("The license code in the config file is not a valid UUID.")
		fmt.Println("The application could not be activated by license and will shut down")
		return false
	}

	// If both the chance and the secret is null, the product can't be activated
	if activation.Chance == nil && (activation.Secret == nil || activation.InstallId == nil) {
		fmt.Println("There is no chance or secret in the config file")
		fmt.Println("The application could not be activated by license and will shut down")
		return false
	}

	// If there is a chance, activate the license of the product.
	if activation.Chance != nil {
		if !takeActivationChance(s.EnterpriseKey, activation.LicenseCode, *activation.Chance) {
			fmt.Println("The application could not be activated by license and will shut down")
			return false
		}
		act := settings.Server.Activation[s.EnterpriseKey]
		activation = &act
	}

	// If there is an activation secret, check if it's correct
	if activation.Secret != nil && activation.InstallId != nil {
		if !checkActivation(s.EnterpriseKey, s.Id, activation.LicenseCode, *activation.Secret, *activation.InstallId) {
			fmt.Println("This product is not activated")
			fmt.Println("The application could not be activated by license and will shut down")
			return false
		}
	} else {
		// Can't continue
		fmt.Println("There must be a secret and a installation ID in the config file to check the activation this product.")
		fmt.Println("The application could not be activated by license and will shut down")
		return false
	}
	return true
}

type TakeChanceActivation struct {
	LicenseCode string `json:"licenseCode"`
	Chance      string `json:"chance"`
	InstallId   string `json:"installId"`
}

// Change the secret by a secret code
func takeActivationChance(enterpriseKey string, licenseCode string, chance string) bool {
	// Generate a random installation ID for this server. This code gets generated by the client for every activation.
	takeActivationChance := TakeChanceActivation{
		LicenseCode: licenseCode,
		Chance:      chance,
		InstallId:   uuid.New().String(),
	}
	data, _ := json.Marshal(takeActivationChance)
	resp, err := http.Post(CHANCE_URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Get the license secret. This code gets generated by the server for every activation.
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	secret := string(response)
	if len(secret) == 30 {
		activation := settings.Server.Activation[enterpriseKey]
		activation.Chance = nil
		activation.Secret = &secret
		activation.InstallId = &takeActivationChance.InstallId
		settings.Server.Activation[enterpriseKey] = activation
		ok := settings.setBackendSettings()

		if !ok {
			fmt.Println("Could not save the activation secret to the config file.")
		}

		return ok
	} else {
		fmt.Println("The server has refused the activation of the license. Check the chance code again.")
		return false
	}
}

type ServerActivation struct {
	LicenseCode string `json:"licenseCode"`
	Secret      string `json:"secret"`
	InstallId   string `json:"installId"`
}

type ServerActivationResult struct {
	Ok             bool  `json:"ok"`
	MaxConnections int16 `json:"maxConnections"`
}

func checkActivation(enterpriseKey string, enterpriseId int32, licenseCode string, secret string, installId string) bool {
	a := ServerActivation{
		LicenseCode: licenseCode,
		Secret:      secret,
		InstallId:   installId,
	}
	data, _ := json.Marshal(a)
	resp, err := http.Post(ACTIVATE_URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Get the check
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var result ServerActivationResult
	json.Unmarshal(response, &result)

	if result.Ok { // Successfully activated
		setLicenseMaxConnectionsLimit(result.MaxConnections, enterpriseId)
		return true
	} else {
		// The license code is incorrect or (probably) the license code has been used in another marketnet installation so the secret or install id has changed.
		// Remove the activation data from the settings and prevent the server from starting.
		activation := settings.Server.Activation[enterpriseKey]
		activation.Chance = nil
		activation.Secret = nil
		activation.InstallId = nil
		settings.setBackendSettings()
		return false
	}
}

// Changes done here must algo be done in the updateSettingsRecord function in settings.go.
func setLicenseMaxConnectionsLimit(maxConnections int16, enterpriseId int32) {
	licenseMaxConnections[enterpriseId] = maxConnections
	s := getSettingsRecordById(enterpriseId)
	if s.MaxConnections == 0 {
		s.MaxConnections = int32(maxConnections)
	} else {
		s.MaxConnections = int32(math.Min(float64(s.MaxConnections), float64(maxConnections)))
	}
	sqlStatement := `UPDATE config SET max_connections=$2 WHERE id=$1`
	db.Exec(sqlStatement, enterpriseId, s.MaxConnections)
}
