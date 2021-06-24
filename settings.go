package main

import (
	"encoding/json"
	"io/ioutil"
)

// Basic, static, server settings such as the DB password or the port.
type BackendSettings struct {
	Db     DatabaseSettings `json:"db"`
	Server ServerSettings   `json:"server"`
}

// Credentials for connecting to PostgreSQL.
type DatabaseSettings struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// Basic info for the app.
type ServerSettings struct {
	Port                 uint16 `json:"port"`
	HashIterations       int32  `json:"hashIterations"`
	TokenExpirationHours int16  `json:"tokenExpirationHours"`
}

func getBackendSettings() (BackendSettings, bool) {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return BackendSettings{}, false
	}

	var settings BackendSettings
	err = json.Unmarshal(content, &settings)
	if err != nil {
		return BackendSettings{}, false
	}
	return settings, true
}

// Advanced settings stored in the database. Configurable by final users.
type Settings struct {
	DefaultVatPercent             float32 `json:"defaultVatPercent"`
	DefaultWarehouse              string  `json:"defaultWarehouse"`
	DefaultWarehouseName          string  `json:"defaultWarehouseName"`
	DateFormat                    string  `json:"dateFormat"`
	EnterpriseName                string  `json:"enterpriseName"`
	EnterpriseDescription         string  `json:"enterpriseDescription"`
	Ecommerce                     string  `json:"ecommerce"`
	PrestaShopUrl                 string  `json:"prestaShopUrl"`
	PrestaShopApiKey              string  `json:"prestaShopApiKey"`
	PrestaShopLanguageId          int32   `json:"prestaShopLanguageId"`
	PrestaShopExportSerie         string  `json:"prestaShopExportSerie"`
	PrestaShopIntracommunitySerie string  `json:"prestaShopIntracommunitySerie"`
	PrestaShopInteriorSerie       string  `json:"prestaShopInteriorSerie"`
	Email                         string  `json:"email"`
	Currency                      string  `json:"currency"`
	CurrencyECBurl                string  `json:"currencyECBurl"`
	BarcodePrefix                 string  `json:"barcodePrefix"`
	CronCurrency                  string  `json:"cronCurrency"`
	CronPrestaShop                string  `json:"cronPrestaShop"`
	SendGridKey                   string  `json:"sendGridKey"`
	EmailFrom                     string  `json:"emailFrom"`
	NameFrom                      string  `json:"nameFrom"`
	PalletWeight                  float32 `json:"palletWeight"`
	PalletWidth                   float32 `json:"palletWidth"`
	PalletHeight                  float32 `json:"palletHeight"`
	PalletDepth                   float32 `json:"palletDepth"`
	MaxConnections                int32   `json:"maxConnections"`
}

func getSettingsRecord() Settings {
	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id=config.default_warehouse) FROM config WHERE id=1`
	row := db.QueryRow(sqlStatement)
	if row.Err() != nil {
		return Settings{}
	}

	var s Settings
	var id int32
	row.Scan(&id, &s.DefaultVatPercent, &s.DefaultWarehouse, &s.DateFormat, &s.EnterpriseName, &s.EnterpriseDescription, &s.Ecommerce, &s.Email, &s.Currency, &s.CurrencyECBurl, &s.BarcodePrefix, &s.PrestaShopUrl, &s.PrestaShopApiKey, &s.PrestaShopLanguageId, &s.PrestaShopExportSerie, &s.PrestaShopIntracommunitySerie, &s.PrestaShopInteriorSerie, &s.CronCurrency, &s.CronPrestaShop, &s.SendGridKey, &s.EmailFrom, &s.NameFrom, &s.PalletWeight, &s.PalletWidth, &s.PalletHeight, &s.PalletDepth, &s.MaxConnections, &s.DefaultWarehouseName)
	return s
}

func (s *Settings) isValid() bool {
	return !(s.DefaultVatPercent < 0 || len(s.DefaultWarehouse) != 2 || len(s.DateFormat) == 0 || len(s.DateFormat) > 25 || len(s.EnterpriseName) == 0 || len(s.EnterpriseName) > 50 || len(s.EnterpriseDescription) > 250 || (s.Ecommerce != "_" && s.Ecommerce != "P" && s.Ecommerce != "M") || (s.Email != "_" && s.Email != "S" && s.Email != "T") || (s.Currency != "_" && s.Currency != "E") || len(s.CurrencyECBurl) > 100 || len(s.BarcodePrefix) > 4 || len(s.PrestaShopUrl) > 100 || len(s.PrestaShopApiKey) > 32 || s.PrestaShopLanguageId < 0 || len(s.CronCurrency) > 25 || len(s.CronPrestaShop) > 25 || len(s.SendGridKey) > 75 || len(s.EmailFrom) > 50 || len(s.NameFrom) > 50 || s.PalletWeight < 0 || s.PalletWidth < 0 || s.PalletHeight < 0 || s.PalletDepth < 0 || s.MaxConnections < 0)
}

func (s *Settings) updateSettingsRecord() bool {
	if !s.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.config SET default_vat_percent=$1, default_warehouse=$2, date_format=$3, enterprise_name=$4, enterprise_description=$5, ecommerce=$6, email=$7, currency=$8, currency_ecb_url=$9, barcode_prefix=$10, prestashop_url=$11, prestashop_api_key=$12, prestashop_language_id=$13, prestashop_export_serie=$14, prestashop_intracommunity_serie=$15, prestashop_interior_serie=$16, cron_currency=$17, cron_prestashop=$18, sendgrid_key=$19, email_from=$20, name_from=$21, pallet_weight=$22, pallet_width=$23, pallet_height=$24, pallet_depth=$25, max_connections=$26 WHERE id=1`
	res, err := db.Exec(sqlStatement, s.DefaultVatPercent, s.DefaultWarehouse, s.DateFormat, s.EnterpriseName, s.EnterpriseDescription, s.Ecommerce, s.Email, s.Currency, s.CurrencyECBurl, s.BarcodePrefix, s.PrestaShopUrl, s.PrestaShopApiKey, s.PrestaShopLanguageId, s.PrestaShopExportSerie, s.PrestaShopIntracommunitySerie, s.PrestaShopInteriorSerie, s.CronCurrency, s.CronPrestaShop, s.SendGridKey, s.EmailFrom, s.NameFrom, s.PalletWeight, s.PalletWidth, s.PalletHeight, s.PalletDepth, s.MaxConnections)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
