package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CONFIG_VER = "1.0"

// Basic, static, server settings such as the DB password or the port.
type BackendSettings struct {
	Version string           `json:"version"`
	Db      DatabaseSettings `json:"db"`
	Server  ServerSettings   `json:"server"`
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
	Port                           uint16                              `json:"port"`
	HashIterations                 int32                               `json:"hashIterations"`
	TokenExpirationHours           int16                               `json:"tokenExpirationHours"`
	MaxLoginAttemps                int16                               `json:"maxLoginAttemps"`
	CronClearLogs                  string                              `json:"cronClearLogs"`
	MaxRequestsPerMinuteEnterprise int32                               `json:"maxRequestsPerMinuteEnterprise"`
	SaaSAccessToken                string                              `json:"SaaSAccessToken"`
	MaxWebHooksPerEnterprise       uint16                              `json:"maxWebHooksPerEnterprise"`
	MaxQueueSizePerWebHook         int32                               `json:"maxQueueSizePerWebHook"`
	WebSecurity                    ServerSettingsWebSecurity           `json:"webSecurity"`
	TLS                            ServerSettingsTLS                   `json:"tls"`
	Activation                     map[string]ServerSettingsActivation `json:"activation"`
}

type ServerSettingsWebSecurity struct {
	ReadTimeoutSeconds        uint8 `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds       uint8 `json:"writeTimeoutSeconds"`
	MaxLimitApiQueries        int64 `json:"maxLimitApiQueries"`
	MaxHeaderBytes            int   `json:"maxHeaderBytes"`
	MaxRequestBodyLength      int64 `json:"maxRequestBodyLength"`
	MaxLengthWebSocketMessage int64 `json:"maxLengthWebSocketMessage"`
}

// SSL settings for the web server.
type ServerSettingsTLS struct {
	UseTLS  bool   `json:"useTLS"`
	CrtPath string `json:"crtPath"`
	KeyPath string `json:"keyPath"`
}

// License activation.
type ServerSettingsActivation struct {
	LicenseCode string  `json:"licenseCode"`
	Chance      *string `json:"chance"`
	Secret      *string `json:"secret"`
	InstallId   *string `json:"installId"`
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

	if settings.isConfigUpgradeRequired() {
		if !settings.upgradeConfig() {
			return BackendSettings{}, false
		}
	}

	return settings, true
}

func (s *BackendSettings) setBackendSettings() bool {
	data, _ := json.MarshalIndent(s, "", "    ")
	err := ioutil.WriteFile("config.json", data, 0700)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func (s *BackendSettings) setBackupBackendSettings() bool {
	data, _ := json.MarshalIndent(s, "", "    ")
	err := ioutil.WriteFile("config.backup.json", data, 0700)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func (s *BackendSettings) isConfigUpgradeRequired() bool {
	return s.Version != CONFIG_VER
}

func (s *BackendSettings) upgradeConfig() bool {
	fmt.Println("Upgrading config.json file. You might want to take a look at the MARKETNET wiki and check out the new config.json properties.")
	if !s.setBackupBackendSettings() {
		fmt.Println("Could not back up JSON file. Error creating config.backup.json.")
		return false
	}
	s.Server.WebSecurity = ServerSettingsWebSecurity{ // Default values
		ReadTimeoutSeconds:        60,
		WriteTimeoutSeconds:       60,
		MaxLimitApiQueries:        1000,
		MaxHeaderBytes:            50000,
		MaxRequestBodyLength:      24000000,
		MaxLengthWebSocketMessage: 24000000,
	}
	s.Version = CONFIG_VER
	if !s.setBackendSettings() {
		fmt.Println("Upgrade failed. Can't write config.json. Please, rename 'config.backup.json' to 'config.json' to restore the aplication's previous status.")
		return false
	} else {
		err := os.Remove("config.backup.json")
		if err != nil {
			fmt.Println("WARNING: Can't delete the temporary backup file: 'config.backup.json'.", err)
		}
		fmt.Println("Successfully upgraded config.json.")
		return true
	}
}

// Advanced settings stored in the database. Configurable by final users.
type Settings struct {
	Id                            int32              `json:"id" gorm:"primaryKey"`
	DefaultVatPercent             float64            `json:"defaultVatPercent" gorm:"type:numeric(14,6);not null:true"`
	DefaultWarehouseId            string             `json:"defaultWarehouseId" gorm:"column:default_warehouse;type:character(2)"`
	DefaultWarehouse              *Warehouse         `json:"defaultWarehouse" gorm:"foreignKey:DefaultWarehouseId,Id;references:Id,EnterpriseId"`
	DateFormat                    string             `json:"dateFormat" gorm:"type:character varying(25);not null:true"`
	EnterpriseName                string             `json:"enterpriseName" gorm:"type:character varying(50);not null:true"`
	EnterpriseDescription         string             `json:"enterpriseDescription" gorm:"type:character varying(250);not null:true"`
	Currency                      string             `json:"currency" gorm:"type:character(1);not null:true"` // "_" = None, "E" = European Central Bank
	CurrencyECBurl                string             `json:"currencyECBurl" gorm:"column:currency_ecb_url;type:character varying(100);not null:true"`
	BarcodePrefix                 string             `json:"barcodePrefix" gorm:"type:character varying(4);not null:true"`
	SettingsEcommerce             *SettingsEcommerce `json:"settingsEcommerce" gorm:"foreignKey:Id;references:EnterpriseId"`
	CronCurrency                  string             `json:"cronCurrency" gorm:"type:character varying(25);not null:true"`
	CronPrestaShop                string             `json:"cronPrestaShop" gorm:"column:cron_prestashop;type:character varying(25);not null:true"`
	PalletWeight                  float64            `json:"palletWeight" gorm:"type:numeric(14,6);not null:true"`
	PalletWidth                   float64            `json:"palletWidth" gorm:"type:numeric(14,6);not null:true"`
	PalletHeight                  float64            `json:"palletHeight" gorm:"type:numeric(14,6);not null:true"`
	PalletDepth                   float64            `json:"palletDepth" gorm:"type:numeric(14,6);not null:true"`
	MaxConnections                int32              `json:"maxConnections" gorm:"not null:true"`
	MinimumStockSalesPeriods      int16              `json:"minimumStockSalesPeriods" gorm:"not null:true"`
	MinimumStockSalesDays         int16              `json:"minimumStockSalesDays" gorm:"not null:true"`
	CustomerJournalId             *int32             `json:"customerJournalId" gorm:"column:customer_journal"`
	CustomerJournal               *Journal           `json:"customerJournal" gorm:"foreignKey:CustomerJournalId,Id;references:Id,EnterpriseId"`
	SalesJournalId                *int32             `json:"salesJournalId" gorm:"column:sales_journal"`
	SalesJournal                  *Journal           `json:"salesJournal" gorm:"foreignKey:SalesJournalId,Id;references:Id,EnterpriseId"`
	SalesAccountId                *int32             `json:"salesAccountId" gorm:"column:sales_account"`
	SalesAccount                  *Account           `json:"salesAccount" gorm:"foreignKey:SalesAccountId,Id;references:Id,EnterpriseId"`
	SupplierJournalId             *int32             `json:"supplierJournalId" gorm:"column:supplier_journal"`
	SupplierJournal               *Journal           `json:"supplierJournal" gorm:"foreignKey:SupplierJournalId,Id;references:Id,EnterpriseId"`
	PurchaseJournalId             *int32             `json:"purchaseJournalId" gorm:"column:purchase_journal"`
	PurchaseJournal               *Journal           `json:"purchaseJournal" gorm:"foreignKey:PurchaseJournalId,Id;references:Id,EnterpriseId"`
	PurchaseAccountId             *int32             `json:"purchaseAccountId" gorm:"column:purchase_account"`
	PurchaseAccount               *Account           `json:"purchaseAccount" gorm:"foreignKey:PurchaseAccountId,Id;references:Id,EnterpriseId"`
	EnableApiKey                  bool               `json:"enableApiKey" gorm:"not null:true"`
	CronClearLabels               string             `json:"cronClearLabels" gorm:"type:character varying(25);not null:true"`
	LimitAccountingDate           *time.Time         `json:"limitAccountingDate" gorm:"type:timestamp(0) with time zone"`
	ConnectionLog                 bool               `json:"connectionLog" gorm:"not null:true"`
	FilterConnections             bool               `json:"filterConnections" gorm:"not null:true"`
	EnterpriseKey                 string             `json:"enterpriseKey" gorm:"type:character varying(25);not null:true;index:config_enterprise_key,unique:true,priority:1"`
	PasswordMinimumLength         int16              `json:"passwordMinimumLength" gorm:"not null:true"`
	PasswordMinumumComplexity     string             `json:"passwordMinumumComplexity" gorm:"type:character(1);not null:true"` // "A": Alphabetical, "B": Alphabetical + numbers, "C": Uppercase + lowercase + numbers, "D": Uppercase + lowercase + numbers + symbols
	InvoiceDeletePolicy           int16              `json:"invoiceDeletePolicy" gorm:"not null:true"`                         // 0 = Allow invoice deletion, 1 = Only allow the deletion of the latest invoice in the billing serie, 2 = Never allow invoice deletion
	TransactionLog                bool               `json:"transactionLog" gorm:"not null:true"`
	UndoManufacturingOrderSeconds int16              `json:"undoManufacturingOrderSeconds" gorm:"not null:true"`
	CronSendCloudTracking         string             `json:"cronSendCloudTracking" gorm:"column:cron_sendcloud_tracking;type:character varying(25);not null:true"`
	SettingsEmail                 *SettingsEmail     `json:"settingsEmail" gorm:"foreignKey:Id;references:EnterpriseId"`
	SettingsCleanUp               *SettingsCleanUp   `json:"settingsCleanUp" gorm:"foreignKey:Id;references:EnterpriseId"`
}

func (s *Settings) TableName() string {
	return "config"
}

func getSettingsRecordById(id int32) Settings {
	var s Settings
	result := dbOrm.Where("config.id = ?", id).Preload(clause.Associations).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

func getSettingsRecords() []Settings {
	settings := make([]Settings, 0)
	result := dbOrm.Model(&Settings{}).Preload(clause.Associations).Find(&settings)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return settings
}

func getSettingsRecordByEnterprise(enterpriseKey string) Settings {
	var s Settings
	result := dbOrm.Where("enterprise_key = ?", enterpriseKey).Preload(clause.Associations).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

func (s *Settings) isValid() bool {
	return !(s.DefaultVatPercent < 0 || len(s.DefaultWarehouseId) != 2 || len(s.DateFormat) == 0 || len(s.DateFormat) > 25 || len(s.EnterpriseName) == 0 || len(s.EnterpriseName) > 50 || len(s.EnterpriseDescription) > 250 || (s.Currency != "_" && s.Currency != "E") || len(s.CurrencyECBurl) > 100 || (s.Currency == "E" && len(s.CurrencyECBurl) == 0) || len(s.BarcodePrefix) > 4 || len(s.CronCurrency) > 25 || len(s.CronPrestaShop) > 25 || s.PalletWeight < 0 || s.PalletWidth < 0 || s.PalletHeight < 0 || s.PalletDepth < 0 || s.MaxConnections < 0 || s.MinimumStockSalesPeriods < 0 || s.MinimumStockSalesDays < 0 || s.PasswordMinimumLength < 6 || (s.PasswordMinumumComplexity != "A" && s.PasswordMinumumComplexity != "B" && s.PasswordMinumumComplexity != "C" && s.PasswordMinumumComplexity != "D") || s.InvoiceDeletePolicy < 0 || s.InvoiceDeletePolicy > 2 || s.UndoManufacturingOrderSeconds < 0 || len(s.CronSendCloudTracking) > 25)
}

func (s *Settings) updateSettingsRecord() bool {
	if !s.isValid() {
		return false
	}

	var salesAccount *int32
	if s.SalesJournalId != nil && *s.SalesJournalId > 0 {
		acc := getAccountIdByAccountNumber(*s.SalesJournalId, 1, s.Id)
		if acc > 0 {
			salesAccount = &acc
		}
	}
	var purchaseAccount *int32
	if s.PurchaseJournalId != nil && *s.PurchaseJournalId > 0 {
		acc := getAccountIdByAccountNumber(*s.PurchaseJournalId, 1, s.Id)
		if acc > 0 {
			purchaseAccount = &acc
		}
	}

	// licensing
	// not in the license map
	_, ok := licenseMaxConnections[s.Id]
	if !ok {
		s.MaxConnections = 0
	}
	// don't let to set more connections than the allowed in the license
	if s.MaxConnections <= 0 {
		s.MaxConnections = int32(licenseMaxConnections[s.Id])
	} else {
		s.MaxConnections = int32(math.Min(float64(s.MaxConnections), float64(licenseMaxConnections[s.Id])))
	}

	// limit accounting date
	if s.LimitAccountingDate != nil && (*s.LimitAccountingDate).After(time.Now()) {
		return false
	}

	// connection log
	if !s.ConnectionLog {
		s.FilterConnections = false
	}

	// check crons
	if s.Currency != "_" {
		_, err := cron.ParseStandard(s.CronCurrency)
		if err != nil {
			return false
		}
	}
	if s.SettingsEcommerce.Ecommerce != "_" {
		_, err := cron.ParseStandard(s.CronPrestaShop)
		if err != nil {
			return false
		}
	}
	_, err := cron.ParseStandard(s.CronClearLabels)
	if err != nil {
		return false
	}

	// ¿has the cron changed?
	settingsInMemory := getSettingsRecordById(s.Id)
	if settingsInMemory.CronClearLabels != s.CronClearLabels || settingsInMemory.Currency != s.Currency || settingsInMemory.CronCurrency != s.CronCurrency || settingsInMemory.SettingsEcommerce.Ecommerce != s.SettingsEcommerce.Ecommerce || settingsInMemory.CronPrestaShop != s.CronPrestaShop {
		refreshRunningCrons(settingsInMemory, *s)
	}

	var settingsInDisk Settings
	dbOrm.Model(&Settings{}).Where("id = ?", s.Id).First(&settingsInDisk)

	settingsInDisk.DefaultVatPercent = s.DefaultVatPercent
	settingsInDisk.DefaultWarehouseId = s.DefaultWarehouseId
	settingsInDisk.DateFormat = s.DateFormat
	settingsInDisk.EnterpriseName = s.EnterpriseName
	settingsInDisk.EnterpriseDescription = s.EnterpriseDescription
	settingsInDisk.Currency = s.Currency
	if s.Currency == "E" {
		settingsInDisk.CurrencyECBurl = s.CurrencyECBurl
	}
	settingsInDisk.BarcodePrefix = s.BarcodePrefix
	settingsInDisk.CronCurrency = s.CronCurrency
	settingsInDisk.CronPrestaShop = s.CronPrestaShop
	settingsInDisk.PalletWeight = s.PalletWeight
	settingsInDisk.PalletHeight = s.PalletHeight
	settingsInDisk.PalletDepth = s.PalletDepth
	settingsInDisk.MaxConnections = s.MaxConnections
	settingsInDisk.MinimumStockSalesPeriods = s.MinimumStockSalesPeriods
	settingsInDisk.MinimumStockSalesDays = s.MinimumStockSalesDays
	settingsInDisk.CustomerJournalId = s.CustomerJournalId
	settingsInDisk.SalesJournalId = s.SalesJournalId
	settingsInDisk.SalesAccountId = salesAccount
	settingsInDisk.SupplierJournalId = s.SupplierJournalId
	settingsInDisk.PurchaseJournalId = s.PurchaseJournalId
	settingsInDisk.PurchaseAccountId = purchaseAccount
	settingsInDisk.EnableApiKey = s.EnableApiKey
	settingsInDisk.CronClearLabels = s.CronClearLabels
	settingsInDisk.LimitAccountingDate = s.LimitAccountingDate
	settingsInDisk.ConnectionLog = s.ConnectionLog
	settingsInDisk.FilterConnections = s.FilterConnections
	settingsInDisk.PasswordMinimumLength = s.PasswordMinimumLength
	settingsInDisk.PasswordMinumumComplexity = s.PasswordMinumumComplexity
	settingsInDisk.InvoiceDeletePolicy = s.InvoiceDeletePolicy
	settingsInDisk.TransactionLog = s.TransactionLog
	settingsInDisk.UndoManufacturingOrderSeconds = s.UndoManufacturingOrderSeconds
	settingsInDisk.CronSendCloudTracking = s.CronSendCloudTracking

	trans := dbOrm.Begin()

	result := trans.Save(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	s.SettingsEcommerce.EnterpriseId = s.Id
	if !s.SettingsEcommerce.updateSettingsEcommerce(trans) {
		trans.Rollback()
		return false
	}

	s.SettingsEmail.EnterpriseId = s.Id
	if !s.SettingsEmail.updateSettingsEmail(trans) {
		trans.Rollback()
		return false
	}

	s.SettingsCleanUp.EnterpriseId = s.Id
	if !s.SettingsCleanUp.updateSettingsCleanUp(trans) {
		trans.Rollback()
		return false
	}

	trans.Commit()
	return true
}

type SettingsEcommerce struct {
	Id                                int32          `json:"id" gorm:"primaryKey"`
	Ecommerce                         string         `json:"ecommerce" gorm:"type:character(1);not null:true"` // "_" = None, "P" = PrestaShop, "M" = Magento, "W" = WooCommerce, "S" = Shopify
	PrestaShopUrl                     string         `json:"prestaShopUrl" gorm:"column:prestashop_url;type:character varying(100);not null:true"`
	PrestaShopApiKey                  string         `json:"prestaShopApiKey" gorm:"column:prestashop_api_key;type:character varying(32);not null:true"`
	PrestaShopLanguageId              int32          `json:"prestaShopLanguageId" gorm:"column:prestashop_language_id;not null:true"`
	PrestaShopExportSerieId           *string        `json:"prestaShopExportSerieId" gorm:"column:prestashop_export_serie;type:character(3)"`
	PrestaShopExportSerie             *BillingSerie  `json:"-" gorm:"foreignKey:PrestaShopExportSerieId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopIntracommunitySerieId   *string        `json:"prestaShopIntracommunitySerieId" gorm:"column:prestashop_intracommunity_serie;type:character(3)"`
	PrestaShopIntracommunitySerie     *BillingSerie  `json:"-" gorm:"foreignKey:PrestaShopIntracommunitySerieId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopInteriorSerieId         *string        `json:"prestaShopInteriorSerieId" gorm:"column:prestashop_interior_serie;type:character(3)"`
	PrestaShopInteriorSerie           *BillingSerie  `json:"-" gorm:"foreignKey:PrestaShopInteriorSerieId,EnterpriseId;references:Id,EnterpriseId"`
	PrestaShopStatusPaymentAccepted   int32          `json:"prestashopStatusPaymentAccepted" gorm:"column:prestashop_status_payment_accepted;not null:true"`
	PrestaShopStatusShipped           int32          `json:"prestashopStatusShipped" gorm:"column:prestashop_status_shipped;not null:true"`
	WooCommerceUrl                    string         `json:"woocommerceUrl" gorm:"column:woocommerce_url;type:character varying(100);not null:true"`
	WooCommerceConsumerKey            string         `json:"woocommerceConsumerKey" gorm:"column:woocommerce_consumer_key;type:character varying(50);not null:true"`
	WooCommerceConsumerSecret         string         `json:"woocommerceConsumerSecret" gorm:"column:woocommerce_consumer_secret;type:character varying(50);not null:true"`
	WooCommerceExportSerieId          *string        `json:"wooCommerceExportSerieId" gorm:"column:woocommerce_export_serie;type:character(3)"`
	WooCommerceExportSerie            *BillingSerie  `json:"-" gorm:"foreignKey:WooCommerceExportSerieId,EnterpriseId;references:Id,EnterpriseId"`
	WooCommerceIntracommunitySerieId  *string        `json:"wooCommerceIntracommunitySerieId" gorm:"column:woocommerce_intracommunity_serie;type:character(3)"`
	WooCommerceIntracommunitySerie    *BillingSerie  `json:"-" gorm:"foreignKey:WooCommerceIntracommunitySerieId,EnterpriseId;references:Id,EnterpriseId"`
	WooCommerceInteriorSerieId        *string        `json:"wooCommerceInteriorSerieId" gorm:"column:woocommerce_interior_serie;type:character(3)"`
	WooCommerceInteriorSerie          *BillingSerie  `json:"-" gorm:"foreignKey:WooCommerceInteriorSerieId,EnterpriseId;references:Id,EnterpriseId"`
	WooCommerceDefaultPaymentMethodId *int32         `json:"wooCommerceDefaultPaymentMethodId" gorm:"column:woocommerce_default_payment_method"`
	WooCommerceDefaultPaymentMethod   *PaymentMethod `json:"-" gorm:"foreignKey:WooCommerceDefaultPaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	ShopifyUrl                        string         `json:"shopifyUrl" gorm:"type:character varying(100);not null:true"`
	ShopifyToken                      string         `json:"shopifyToken" gorm:"type:character varying(50);not null:true"`
	ShopifyExportSerieId              *string        `json:"shopifyExportSerieId" gorm:"type:character(3);column:shopify_export_serie"`
	ShopifyExportSerie                *BillingSerie  `json:"-" gorm:"foreignKey:ShopifyExportSerieId,EnterpriseId;references:Id,EnterpriseId"`
	ShopifyIntracommunitySerieId      *string        `json:"shopifyIntracommunitySerieId" gorm:"type:character(3);column:shopify_intracommunity_serie"`
	ShopifyIntracommunitySerie        *BillingSerie  `json:"-" gorm:"foreignKey:ShopifyIntracommunitySerieId,EnterpriseId;references:Id,EnterpriseId"`
	ShopifyInteriorSerieId            *string        `json:"shopifyInteriorSerieId" gorm:"type:character(3);column:shopify_interior_serie"`
	ShopifyInteriorSerie              *BillingSerie  `json:"-" gorm:"foreignKey:ShopifyInteriorSerieId,EnterpriseId;references:Id,EnterpriseId"`
	ShopifyDefaultPaymentMethodId     *int32         `json:"shopifyDefaultPaymentMethodId" gorm:"column:shopify_default_payment_method"`
	ShopifyDefaultPaymentMethod       *PaymentMethod `json:"-" gorm:"foreignKey:ShopifyDefaultPaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	ShopifyShopLocationId             int64          `json:"shopifyShopLocationId" gorm:"not null:true"`
	EnterpriseId                      int32          `json:"-" gorm:"column:enterprise;not null:true;index:config_ecommerce_enterprise,unique:true"`
	Enterprise                        Settings       `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *SettingsEcommerce) TableName() string {
	return "config_ecommerce"
}

func (s *SettingsEcommerce) isValid() bool {
	return !(s.Ecommerce != "_" && s.Ecommerce != "P" && s.Ecommerce != "M" && s.Ecommerce != "W" && s.Ecommerce != "S") || len(s.PrestaShopUrl) > 100 || len(s.PrestaShopApiKey) > 32 || len(s.WooCommerceUrl) > 100 || len(s.WooCommerceConsumerKey) > 50 || len(s.WooCommerceConsumerSecret) > 50 || len(s.ShopifyUrl) > 100 || len(s.ShopifyToken) > 50 || s.PrestaShopLanguageId < 0 || s.PrestaShopStatusPaymentAccepted < 0 || s.PrestaShopStatusShipped < 0 || (s.Ecommerce == "P" && (len(s.PrestaShopUrl) == 0 || len(s.PrestaShopApiKey) == 0 || s.PrestaShopLanguageId == 0 || s.PrestaShopExportSerieId == nil || s.PrestaShopIntracommunitySerieId == nil || s.PrestaShopInteriorSerieId == nil || s.PrestaShopStatusPaymentAccepted == 0 || s.PrestaShopStatusShipped == 0)) || (s.Ecommerce == "W" && (len(s.WooCommerceUrl) == 0 || len(s.WooCommerceConsumerKey) == 0 || len(s.WooCommerceConsumerSecret) == 0 || s.WooCommerceDefaultPaymentMethodId == nil || s.WooCommerceExportSerieId == nil || s.WooCommerceInteriorSerieId == nil || s.WooCommerceIntracommunitySerieId == nil)) || (s.Ecommerce == "S" && (len(s.ShopifyUrl) == 0 || len(s.ShopifyToken) == 0 || s.ShopifyDefaultPaymentMethodId == nil || s.ShopifyExportSerieId == nil || s.ShopifyInteriorSerieId == nil || s.ShopifyIntracommunitySerieId == nil))
}

func (s *SettingsEcommerce) updateSettingsEcommerce(trans *gorm.DB) bool {
	if !s.isValid() {
		return false
	}

	var settingsInDisk SettingsEcommerce
	result := dbOrm.Model(&SettingsEcommerce{}).Where("enterprise = ?", s.EnterpriseId).First(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	settingsInDisk.Ecommerce = s.Ecommerce
	if s.Ecommerce == "P" {
		settingsInDisk.PrestaShopUrl = s.PrestaShopUrl
		settingsInDisk.PrestaShopApiKey = s.PrestaShopApiKey
		settingsInDisk.PrestaShopLanguageId = s.PrestaShopLanguageId
		settingsInDisk.PrestaShopExportSerieId = s.PrestaShopExportSerieId
		settingsInDisk.PrestaShopIntracommunitySerieId = s.PrestaShopIntracommunitySerieId
		settingsInDisk.PrestaShopInteriorSerieId = s.PrestaShopInteriorSerieId
		settingsInDisk.PrestaShopStatusPaymentAccepted = s.PrestaShopStatusPaymentAccepted
		settingsInDisk.PrestaShopStatusShipped = s.PrestaShopStatusShipped
	} else if s.Ecommerce == "W" {
		settingsInDisk.WooCommerceUrl = s.WooCommerceUrl
		settingsInDisk.WooCommerceConsumerKey = s.WooCommerceConsumerKey
		settingsInDisk.WooCommerceConsumerSecret = s.WooCommerceConsumerSecret
		settingsInDisk.WooCommerceExportSerieId = s.WooCommerceExportSerieId
		settingsInDisk.WooCommerceIntracommunitySerieId = s.WooCommerceIntracommunitySerieId
		settingsInDisk.WooCommerceInteriorSerieId = s.WooCommerceInteriorSerieId
		settingsInDisk.WooCommerceDefaultPaymentMethodId = s.WooCommerceDefaultPaymentMethodId
	} else if s.Ecommerce == "S" {
		settingsInDisk.ShopifyUrl = s.ShopifyUrl
		settingsInDisk.ShopifyToken = s.ShopifyToken
		settingsInDisk.ShopifyExportSerieId = s.ShopifyExportSerieId
		settingsInDisk.ShopifyIntracommunitySerieId = s.ShopifyIntracommunitySerieId
		settingsInDisk.ShopifyInteriorSerieId = s.ShopifyInteriorSerieId
		settingsInDisk.ShopifyDefaultPaymentMethodId = s.ShopifyDefaultPaymentMethodId
		settingsInDisk.ShopifyShopLocationId = s.ShopifyShopLocationId
	}

	result = trans.Save(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type SettingsEmail struct {
	Id                      int32    `json:"id" gorm:"primaryKey"`
	Email                   string   `json:"email" gorm:"type:character(1);not null:true"` // "_" = None, "S" = SendGrid, "T" = SMTP
	SendGridKey             string   `json:"sendGridKey" gorm:"column:sendgrid_key;type:character varying(75);not null:true"`
	EmailFrom               string   `json:"emailFrom" gorm:"type:character varying(50);not null:true"`
	NameFrom                string   `json:"nameFrom" gorm:"type:character varying(50);not null:true"`
	SMTPIdentity            string   `json:"SMTPIdentity" gorm:"type:character varying(50);not null:true"`
	SMTPUsername            string   `json:"SMTPUsername" gorm:"type:character varying(50);not null:true"`
	SMTPPassword            string   `json:"SMTPPassword" gorm:"type:character varying(50);not null:true"`
	SMTPHostname            string   `json:"SMTPHostname" gorm:"type:character varying(50);not null:true"`
	SMTPSTARTTLS            bool     `json:"SMTPSTARTTLS" gorm:"column:smtp_starttls;not null:true"`
	SMTPReplyTo             string   `json:"SMTPReplyTo" gorm:"type:character varying(50);not null:true"`
	EmailSendErrorEcommerce string   `json:"emailSendErrorEcommerce" gorm:"type:character varying(150);not null:true"`
	EmailSendErrorSendCloud string   `json:"emailSendErrorSendCloud" gorm:"column:email_send_error_sendcloud;type:character varying(150);not null:true"`
	EnterpriseId            int32    `json:"-" gorm:"column:enterprise;not null:true;index:config_email_enterprise,unique:true"`
	Enterprise              Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *SettingsEmail) TableName() string {
	return "config_email"
}

func (s *SettingsEmail) isValid() bool {
	return !((s.Email != "_" && s.Email != "S" && s.Email != "T") || len(s.SendGridKey) > 75 || len(s.EmailFrom) > 50 || len(s.NameFrom) > 50 || len(s.SMTPIdentity) > 50 || len(s.SMTPUsername) > 50 || len(s.SMTPPassword) > 50 || len(s.SMTPHostname) > 50 || len(s.SMTPReplyTo) > 50 || len(s.EmailSendErrorEcommerce) > 150 || len(s.EmailSendErrorSendCloud) > 150 || (s.Email == "S" && (len(s.SendGridKey) == 0 || len(s.EmailFrom) == 0 || len(s.NameFrom) == 0 || !emailIsValid(s.EmailFrom))) || (s.Email == "T" && (len(s.SMTPUsername) == 0 || len(s.SMTPPassword) == 0 || len(s.SMTPHostname) == 0 || !emailIsValid(s.SMTPUsername) || !hostnameWithPortValid(s.SMTPHostname))))
}

func (s *SettingsEmail) updateSettingsEmail(trans *gorm.DB) bool {
	if !s.isValid() {
		return false
	}

	var settingsInDisk SettingsEmail
	result := dbOrm.Model(&SettingsEmail{}).Where("enterprise = ?", s.EnterpriseId).First(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	settingsInDisk.Email = s.Email
	if s.Email == "S" {
		settingsInDisk.SendGridKey = s.SendGridKey
		settingsInDisk.EmailFrom = s.EmailFrom
		settingsInDisk.NameFrom = s.NameFrom
	} else if s.Email == "T" {
		settingsInDisk.SMTPIdentity = s.SMTPIdentity
		settingsInDisk.SMTPUsername = s.SMTPUsername
		settingsInDisk.SMTPPassword = s.SMTPPassword
		settingsInDisk.SMTPHostname = s.SMTPHostname
		settingsInDisk.SMTPSTARTTLS = s.SMTPSTARTTLS
		settingsInDisk.SMTPReplyTo = s.SMTPReplyTo
	}

	result = trans.Save(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type SettingsCleanUp struct {
	CronCleanTransactionalLog string   `json:"cronCleanTransactionalLog" gorm:"type:character varying(25);not null:true"`
	TransactionalLogDays      int16    `json:"transactionalLogDays" gorm:"type:smallint;not null:true"`
	CronCleanConnectionLog    string   `json:"cronCleanConnectionLog" gorm:"type:character varying(25);not null:true"`
	ConnectionLogDays         int16    `json:"connectionLogDays" gorm:"type:smallint;not null:true"`
	CronCleanLoginToken       string   `json:"cronCleanLoginToken" gorm:"type:character varying(25);not null:true"`
	EnterpriseId              int32    `json:"-" gorm:"column:enterprise;not null:true;index:config_cleanup_enterprise,unique:true"`
	Enterprise                Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (s *SettingsCleanUp) TableName() string {
	return "config_cleanup"
}

func (s *SettingsCleanUp) isValid() bool {
	return !(len(s.CronCleanTransactionalLog) > 25 || len(s.CronCleanConnectionLog) > 25 || len(s.CronCleanLoginToken) > 25 || s.TransactionalLogDays < 0 || s.ConnectionLogDays < 0)
}

func (s *SettingsCleanUp) updateSettingsCleanUp(trans *gorm.DB) bool {
	if !s.isValid() {
		return false
	}

	var settingsInDisk SettingsCleanUp
	result := dbOrm.Model(&SettingsCleanUp{}).Where("enterprise = ?", s.EnterpriseId).First(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	settingsInDisk.CronCleanTransactionalLog = s.CronCleanTransactionalLog
	settingsInDisk.TransactionalLogDays = s.TransactionalLogDays
	settingsInDisk.CronCleanConnectionLog = s.CronCleanConnectionLog
	settingsInDisk.ConnectionLogDays = s.ConnectionLogDays
	settingsInDisk.CronCleanLoginToken = s.CronCleanLoginToken

	result = trans.Save(&settingsInDisk)
	if result.Error != nil {
		fmt.Println(result.Error)
		log("DB", result.Error.Error())
		return false
	}

	return true
}

// Don't allow every client to get the secret data, like API keys.
// This object holds the config that every client has to know, and the "Settings" object contains admin information.
type ClientSettings struct {
	DefaultVatPercent             float64              `json:"defaultVatPercent"`
	DefaultWarehouse              string               `json:"defaultWarehouse"`
	DefaultWarehouseName          string               `json:"defaultWarehouseName"`
	DateFormat                    string               `json:"dateFormat"`
	Ecommerce                     string               `json:"ecommerce"`           // "_" = None, "P" = PrestaShop, "M" = Magento
	InvoiceDeletePolicy           int16                `json:"invoiceDeletePolicy"` // 0 = Allow invoice deletion, 1 = Only allow the deletion of the latest invoice in the billing serie, 2 = Never allow invoice deletion
	LabelPrinterProfileEAN13      *LabelPrinterProfile `json:"labelPrinterProfileEAN13"`
	LabelPrinterProfileCode128    *LabelPrinterProfile `json:"labelPrinterProfileCode128"`
	LabelPrinterProfileDataMatrix *LabelPrinterProfile `json:"labelPrinterProfileDataMatrix"`
}

func (s Settings) censorSettings() ClientSettings {
	c := ClientSettings{}
	c.DefaultVatPercent = s.DefaultVatPercent
	c.DefaultWarehouse = s.DefaultWarehouseId
	c.DefaultWarehouseName = s.DefaultWarehouse.Name
	c.DateFormat = s.DateFormat
	c.Ecommerce = s.SettingsEcommerce.Ecommerce
	c.InvoiceDeletePolicy = s.InvoiceDeletePolicy

	c.LabelPrinterProfileEAN13 = getLabelPrinterProfileByEnterpriseTypeAndActive(s.Id, "E")
	c.LabelPrinterProfileCode128 = getLabelPrinterProfileByEnterpriseTypeAndActive(s.Id, "C")
	c.LabelPrinterProfileDataMatrix = getLabelPrinterProfileByEnterpriseTypeAndActive(s.Id, "D")

	return c
}

type EnterpriseCronInfo struct {
	CronClearLabels       cron.EntryID
	CronCurrency          *cron.EntryID
	CronPrestaShop        *cron.EntryID
	CronSendcloudTracking *cron.EntryID
}

func refreshRunningCrons(oldSettings Settings, newSettings Settings) {
	runningCronsMutex.Lock()
	enterpriseCronInfo := runningCrons[oldSettings.Id]

	if oldSettings.CronClearLabels != newSettings.CronClearLabels {
		c.Remove(enterpriseCronInfo.CronClearLabels)
		cronId, err := c.AddFunc(newSettings.CronClearLabels, func() {
			deleteAllShippingTags(oldSettings.Id)
		})
		if err != nil {
			enterpriseCronInfo.CronClearLabels = cronId
		}
	}
	if oldSettings.Currency != newSettings.Currency || oldSettings.CronCurrency != newSettings.CronCurrency {
		if enterpriseCronInfo.CronCurrency != nil {
			c.Remove(*enterpriseCronInfo.CronCurrency)
		}
		if newSettings.Currency != "_" {
			cronId, err := c.AddFunc(newSettings.CronCurrency, func() {
				updateCurrencyExchange(oldSettings.Id)
			})
			if err != nil {
				enterpriseCronInfo.CronCurrency = &cronId
			}
		}
	}
	if oldSettings.SettingsEcommerce.Ecommerce != newSettings.SettingsEcommerce.Ecommerce || oldSettings.CronPrestaShop != newSettings.CronPrestaShop {
		if enterpriseCronInfo.CronPrestaShop != nil {
			c.Remove(*enterpriseCronInfo.CronPrestaShop)
		}
		if newSettings.SettingsEcommerce.Ecommerce != "_" {
			e := ECommerce{Enterprise: oldSettings.Id}
			cronId, err := c.AddFunc(newSettings.CronPrestaShop, e.ecommerceControllerImportFromEcommerce)
			if err != nil {
				enterpriseCronInfo.CronPrestaShop = &cronId
			}
		}
	}
	if oldSettings.CronSendCloudTracking != newSettings.CronSendCloudTracking {
		if enterpriseCronInfo.CronSendcloudTracking != nil {
			c.Remove(*enterpriseCronInfo.CronSendcloudTracking)
		}
		if newSettings.CronSendCloudTracking != "" {
			cronId, err := c.AddFunc(newSettings.CronSendCloudTracking, func() {
				getShippingTrackingSendCloud(oldSettings.Id)
			})
			if err != nil {
				enterpriseCronInfo.CronSendcloudTracking = &cronId
			}
		}
	}

	runningCrons[oldSettings.Id] = enterpriseCronInfo
	runningCronsMutex.Unlock()
}

func (s *Settings) BeforeCreate(tx *gorm.DB) (err error) {
	var settings Settings
	tx.Model(&Settings{}).Last(&settings)
	s.Id = settings.Id + 1
	return nil
}

func addEnterpriseFromParameters() bool {
	enterpriseKey, ok := getParameterValue("enterprise_key")
	if !ok {
		fmt.Println("Error, parameter enterprise_key not present.")
		return false
	}
	enterpriseName, ok := getParameterValue("enterprise_name")
	if !ok {
		fmt.Println("Error, parameter enterprise_name not present.")
		return false
	}
	enterpriseDesc, ok := getParameterValue("enterprise_desc")
	if !ok {
		fmt.Println("Error, parameter enterprise_desc not present.")
		return false
	}
	userPassword, ok := getParameterValue("user_password")
	if !ok {
		fmt.Println("Error, parameter user_password not present.")
		return false
	}
	licenseCode, ok := getParameterValue("license_code")
	if !ok {
		fmt.Println("Error, parameter license_code not present.")
		return false
	}
	licenseChance, ok := getParameterValue("license_chance")
	if !ok {
		fmt.Println("Error, parameter license_chance not present.")
		return false
	}

	return createNewEnterprise(enterpriseName, enterpriseDesc, enterpriseKey, licenseCode, licenseChance, userPassword, 0)
}

func createNewEnterprise(enterpriseName string, enterpriseDesc string, enterpriseKey string, licenseCode string, licenseChance string, userPassword string, documentSpace float64) bool {
	if len(enterpriseKey) == 0 || len(enterpriseName) == 0 || len(userPassword) < 8 || len(licenseCode) == 0 || len(licenseChance) == 0 {
		fmt.Println("Error: Invalid data in parameters. Check all the parameters in the documentation.")
		return false
	}

	ok, enterpriseId := initialConfigCreateEnterprise(enterpriseName, enterpriseDesc, strings.ToUpper(enterpriseKey))
	if !ok || enterpriseId <= 0 {
		fmt.Println("Error: Could not create the enterprise in the database.")
		return false
	}

	initialData(enterpriseId)

	sqlStatement := `UPDATE config SET default_warehouse=$1 WHERE id=$2`
	_, err := db.Exec(sqlStatement, "W1", enterpriseId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	config := getSettingsRecordById(enterpriseId)
	ecommerceExportSerie := "EXP"
	ecommerceIntracommunitySerie := "IEU"
	ecommerceInteriorSerie := "INT"
	config.SettingsEcommerce = new(SettingsEcommerce)
	config.SettingsEcommerce.PrestaShopExportSerieId = &ecommerceExportSerie
	config.SettingsEcommerce.PrestaShopIntracommunitySerieId = &ecommerceIntracommunitySerie
	config.SettingsEcommerce.PrestaShopInteriorSerieId = &ecommerceInteriorSerie
	config.SettingsEcommerce.PrestaShopStatusPaymentAccepted = 2
	config.SettingsEcommerce.PrestaShopStatusShipped = 4
	config.SettingsEcommerce.WooCommerceExportSerieId = &ecommerceExportSerie
	config.SettingsEcommerce.WooCommerceIntracommunitySerieId = &ecommerceIntracommunitySerie
	config.SettingsEcommerce.WooCommerceInteriorSerieId = &ecommerceInteriorSerie
	config.SettingsEcommerce.ShopifyExportSerieId = &ecommerceExportSerie
	config.SettingsEcommerce.ShopifyIntracommunitySerieId = &ecommerceIntracommunitySerie
	config.SettingsEcommerce.ShopifyInteriorSerieId = &ecommerceInteriorSerie
	if !config.updateSettingsRecord() {
		fmt.Println("Error: Could not update the settings record.")
		return false
	}

	activation := ServerSettingsActivation{
		LicenseCode: licenseCode,
		Chance:      &licenseChance,
	}
	settings.Server.Activation[enterpriseKey] = activation
	settings.setBackendSettings()
	if !activation.activateEnterprise(enterpriseId) {
		fmt.Println("Error: Could not activate by license the new enterprise.")
		return false
	}

	insert := UserInsert{
		Username: "marketnet",
		FullName: "MARKETNET ADMINISTRATOR",
		Password: userPassword,
		Language: "en",
	}
	if !insert.insertUser(enterpriseId) {
		fmt.Println("Error: Cloud not create the new admin user.")
		return false
	} else {
		fmt.Println("Generated admin user")
	}

	group := Group{
		Name:          "Administrators",
		Sales:         true,
		Purchases:     true,
		Masters:       true,
		Warehouse:     true,
		Manufacturing: true,
		Preparation:   true,
		Admin:         true,
		PrestaShop:    true,
		Accounting:    true,
		EnterpriseId:  enterpriseId,
	}
	if !group.insertGroup() {
		fmt.Println("Error: Could not create the admin group.")
		return false
	} else {
		fmt.Println("Generated admin group")
	}

	users := getUser(enterpriseId)
	user := users[len(users)-1]

	ug := UserGroup{
		UserId:  user.Id,
		GroupId: group.Id,
	}
	if !ug.insertUserGroup() {
		fmt.Println("Error: Could not assign the admin user to the admin group.")
		return false
	} else {
		fmt.Println("Added the admin user to the admin group")
	}

	if documentSpace > 0 {
		documentContainerUUID := uuid.NewString()
		workingDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return false
		}

		err = os.Mkdir(path.Join(workingDir, documentContainerUUID), 0755)
		if err != nil {
			fmt.Println(err)
			return false
		}

		dc := DocumentContainer{
			Name:         "Default Document Container",
			Path:         path.Join(workingDir, documentContainerUUID),
			MaxStorage:   int64(documentSpace) * 1000000000, // Gb to bytes
			EnterpriseId: enterpriseId,
		}
		if !dc.insertDocumentContainer(true) {
			fmt.Println("Could not create document container")
			return false
		}
	}

	return true
}
