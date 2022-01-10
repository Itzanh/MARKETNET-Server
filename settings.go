package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
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
	Port                           uint16                              `json:"port"`
	HashIterations                 int32                               `json:"hashIterations"`
	TokenExpirationHours           int16                               `json:"tokenExpirationHours"`
	MaxLoginAttemps                int16                               `json:"maxLoginAttemps"`
	CronClearLogs                  string                              `json:"cronClearLogs"`
	MaxRequestsPerMinuteEnterprise int32                               `json:"maxRequestsPerMinuteEnterprise"`
	SaaSAccessToken                string                              `json:"SaaSAccessToken"`
	TLS                            ServerSettingsTLS                   `json:"tls"`
	Activation                     map[string]ServerSettingsActivation `json:"activation"`
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
	return settings, true
}

func (s *BackendSettings) setBackendSettings() bool {
	data, _ := json.MarshalIndent(s, "", "    ")
	err := ioutil.WriteFile("config.json", data, 0700)
	return err == nil
}

// Advanced settings stored in the database. Configurable by final users.
type Settings struct {
	Id                              int32      `json:"id"`
	DefaultVatPercent               float64    `json:"defaultVatPercent"`
	DefaultWarehouse                string     `json:"defaultWarehouse"`
	DefaultWarehouseName            string     `json:"defaultWarehouseName"`
	DateFormat                      string     `json:"dateFormat"`
	EnterpriseName                  string     `json:"enterpriseName"`
	EnterpriseDescription           string     `json:"enterpriseDescription"`
	Ecommerce                       string     `json:"ecommerce"` // "_" = None, "P" = PrestaShop, "M" = Magento, "W" = WooCommerce, "S" = Shopify
	PrestaShopUrl                   string     `json:"prestaShopUrl"`
	PrestaShopApiKey                string     `json:"prestaShopApiKey"`
	PrestaShopLanguageId            int32      `json:"prestaShopLanguageId"`
	PrestaShopExportSerie           *string    `json:"prestaShopExportSerie"`
	PrestaShopIntracommunitySerie   *string    `json:"prestaShopIntracommunitySerie"`
	PrestaShopInteriorSerie         *string    `json:"prestaShopInteriorSerie"`
	Email                           string     `json:"email"`    // "_" = None, "S" = SendGrid, "T" = SMTP
	Currency                        string     `json:"currency"` // "_" = None, "E" = European Central Bank
	CurrencyECBurl                  string     `json:"currencyECBurl"`
	BarcodePrefix                   string     `json:"barcodePrefix"`
	CronCurrency                    string     `json:"cronCurrency"`
	CronPrestaShop                  string     `json:"cronPrestaShop"`
	SendGridKey                     string     `json:"sendGridKey"`
	EmailFrom                       string     `json:"emailFrom"`
	NameFrom                        string     `json:"nameFrom"`
	PalletWeight                    float64    `json:"palletWeight"`
	PalletWidth                     float64    `json:"palletWidth"`
	PalletHeight                    float64    `json:"palletHeight"`
	PalletDepth                     float64    `json:"palletDepth"`
	MaxConnections                  int32      `json:"maxConnections"`
	PrestashopStatusPaymentAccepted int32      `json:"prestashopStatusPaymentAccepted"`
	PrestashopStatusShipped         int32      `json:"prestashopStatusShipped"`
	MinimumStockSalesPeriods        int16      `json:"minimumStockSalesPeriods"`
	MinimumStockSalesDays           int16      `json:"minimumStockSalesDays"`
	CustomerJournal                 *int32     `json:"customerJournal"`
	SalesJournal                    *int32     `json:"salesJournal"`
	SalesAccount                    *int32     `json:"salesAccount"`
	SupplierJournal                 *int32     `json:"supplierJournal"`
	PurchaseJournal                 *int32     `json:"purchaseJournal"`
	PurchaseAccount                 *int32     `json:"purchaseAccount"`
	EnableApiKey                    bool       `json:"enableApiKey"`
	CronClearLabels                 string     `json:"cronClearLabels"`
	LimitAccountingDate             *time.Time `json:"limitAccountingDate"`
	WooCommerceUrl                  string     `json:"woocommerceUrl"`
	WooCommerceConsumerKey          string     `json:"woocommerceConsumerKey"`
	WooCommerceConsumerSecret       string     `json:"woocommerceConsumerSecret"`
	WooCommerceExportSerie          *string    `json:"wooCommerceExportSerie"`
	WooCommerceIntracommunitySerie  *string    `json:"wooCommerceIntracommunitySerie"`
	WooCommerceInteriorSerie        *string    `json:"wooCommerceInteriorSerie"`
	WooCommerceDefaultPaymentMethod *int32     `json:"wooCommerceDefaultPaymentMethod"`
	ConnectionLog                   bool       `json:"connectionLog"`
	FilterConnections               bool       `json:"filterConnections"`
	ShopifyUrl                      string     `json:"shopifyUrl"`
	ShopifyToken                    string     `json:"shopifyToken"`
	ShopifyExportSerie              *string    `json:"shopifyExportSerie"`
	ShopifyIntracommunitySerie      *string    `json:"shopifyIntracommunitySerie"`
	ShopifyInteriorSerie            *string    `json:"shopifyInteriorSerie"`
	ShopifyDefaultPaymentMethod     *int32     `json:"shopifyDefaultPaymentMethod"`
	ShopifyShopLocationId           int64      `json:"shopifyShopLocationId"`
	EnterpriseKey                   string     `json:"enterpriseKey"`
	PasswordMinimumLength           int16      `json:"passwordMinimumLength"`
	PasswordMinumumComplexity       string     `json:"passwordMinumumComplexity"` // "A": Alphabetical, "B": Alphabetical + numbers, "C": Uppercase + lowercase + numbers, "D": Uppercase + lowercase + numbers + symbols
	InvoiceDeletePolicy             int16      `json:"invoiceDeletePolicy"`       // 0 = Allow invoice deletion, 1 = Only allow the deletion of the latest invoice in the billing serie, 2 = Never allow invoice deletion
	TransactionLog                  bool       `json:"transactionLog"`
	UndoManufacturingOrderSeconds   int16      `json:"undoManufacturingOrderSeconds"`
	CronSendcloudTracking           string     `json:"cronSendcloudTracking"`
	SMTPIdentity                    string     `json:"SMTPIdentity"`
	SMTPUsername                    string     `json:"SMTPUsername"`
	SMTPPassword                    string     `json:"SMTPPassword"`
	SMTPHostname                    string     `json:"SMTPHostname"`
	SMTPSTARTTLS                    bool       `json:"SMTPSTARTTLS"`
}

func getSettingsRecordById(id int32) Settings {
	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id=config.default_warehouse AND warehouse.enterprise=config.id) FROM config WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return Settings{}
	}

	var s Settings
	row.Scan(&s.Id, &s.DefaultVatPercent, &s.DefaultWarehouse, &s.DateFormat, &s.EnterpriseName, &s.EnterpriseDescription, &s.Ecommerce, &s.Email, &s.Currency, &s.CurrencyECBurl, &s.BarcodePrefix, &s.PrestaShopUrl, &s.PrestaShopApiKey, &s.PrestaShopLanguageId, &s.PrestaShopExportSerie, &s.PrestaShopIntracommunitySerie, &s.PrestaShopInteriorSerie, &s.CronCurrency, &s.CronPrestaShop, &s.SendGridKey, &s.EmailFrom, &s.NameFrom, &s.PalletWeight, &s.PalletWidth, &s.PalletHeight, &s.PalletDepth, &s.MaxConnections, &s.PrestashopStatusPaymentAccepted, &s.PrestashopStatusShipped, &s.MinimumStockSalesPeriods, &s.MinimumStockSalesDays, &s.CustomerJournal, &s.SalesJournal, &s.SalesAccount, &s.SupplierJournal, &s.PurchaseJournal, &s.PurchaseAccount, &s.EnableApiKey, &s.CronClearLabels, &s.LimitAccountingDate, &s.WooCommerceUrl, &s.WooCommerceConsumerKey, &s.WooCommerceConsumerSecret, &s.WooCommerceExportSerie, &s.WooCommerceIntracommunitySerie, &s.WooCommerceInteriorSerie, &s.WooCommerceDefaultPaymentMethod, &s.ConnectionLog, &s.FilterConnections, &s.ShopifyUrl, &s.ShopifyToken, &s.ShopifyExportSerie, &s.ShopifyIntracommunitySerie, &s.ShopifyInteriorSerie, &s.ShopifyDefaultPaymentMethod, &s.ShopifyShopLocationId, &s.EnterpriseKey, &s.PasswordMinimumLength, &s.PasswordMinumumComplexity, &s.InvoiceDeletePolicy, &s.TransactionLog, &s.UndoManufacturingOrderSeconds, &s.CronSendcloudTracking, &s.SMTPIdentity, &s.SMTPUsername, &s.SMTPPassword, &s.SMTPHostname, &s.SMTPSTARTTLS, &s.DefaultWarehouseName)
	return s
}

func getSettingsRecords() []Settings {
	settings := make([]Settings, 0)
	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id=config.default_warehouse AND warehouse.enterprise=config.id) FROM config`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return settings
	}
	defer rows.Close()

	for rows.Next() {
		var s Settings
		rows.Scan(&s.Id, &s.DefaultVatPercent, &s.DefaultWarehouse, &s.DateFormat, &s.EnterpriseName, &s.EnterpriseDescription, &s.Ecommerce, &s.Email, &s.Currency, &s.CurrencyECBurl, &s.BarcodePrefix, &s.PrestaShopUrl, &s.PrestaShopApiKey, &s.PrestaShopLanguageId, &s.PrestaShopExportSerie, &s.PrestaShopIntracommunitySerie, &s.PrestaShopInteriorSerie, &s.CronCurrency, &s.CronPrestaShop, &s.SendGridKey, &s.EmailFrom, &s.NameFrom, &s.PalletWeight, &s.PalletWidth, &s.PalletHeight, &s.PalletDepth, &s.MaxConnections, &s.PrestashopStatusPaymentAccepted, &s.PrestashopStatusShipped, &s.MinimumStockSalesPeriods, &s.MinimumStockSalesDays, &s.CustomerJournal, &s.SalesJournal, &s.SalesAccount, &s.SupplierJournal, &s.PurchaseJournal, &s.PurchaseAccount, &s.EnableApiKey, &s.CronClearLabels, &s.LimitAccountingDate, &s.WooCommerceUrl, &s.WooCommerceConsumerKey, &s.WooCommerceConsumerSecret, &s.WooCommerceExportSerie, &s.WooCommerceIntracommunitySerie, &s.WooCommerceInteriorSerie, &s.WooCommerceDefaultPaymentMethod, &s.ConnectionLog, &s.FilterConnections, &s.ShopifyUrl, &s.ShopifyToken, &s.ShopifyExportSerie, &s.ShopifyIntracommunitySerie, &s.ShopifyInteriorSerie, &s.ShopifyDefaultPaymentMethod, &s.ShopifyShopLocationId, &s.EnterpriseKey, &s.PasswordMinimumLength, &s.PasswordMinumumComplexity, &s.InvoiceDeletePolicy, &s.TransactionLog, &s.UndoManufacturingOrderSeconds, &s.CronSendcloudTracking, &s.SMTPIdentity, &s.SMTPUsername, &s.SMTPPassword, &s.SMTPHostname, &s.SMTPSTARTTLS, &s.DefaultWarehouseName)
		settings = append(settings, s)
	}

	return settings
}

func getSettingsRecordByEnterprise(enterpriseKey string) Settings {
	sqlStatement := `SELECT *,(SELECT name FROM warehouse WHERE warehouse.id=config.default_warehouse AND warehouse.enterprise=config.id) FROM config WHERE enterprise_key=$1`
	row := db.QueryRow(sqlStatement, strings.ToUpper(enterpriseKey))
	if row.Err() != nil {
		return Settings{}
	}

	var s Settings
	row.Scan(&s.Id, &s.DefaultVatPercent, &s.DefaultWarehouse, &s.DateFormat, &s.EnterpriseName, &s.EnterpriseDescription, &s.Ecommerce, &s.Email, &s.Currency, &s.CurrencyECBurl, &s.BarcodePrefix, &s.PrestaShopUrl, &s.PrestaShopApiKey, &s.PrestaShopLanguageId, &s.PrestaShopExportSerie, &s.PrestaShopIntracommunitySerie, &s.PrestaShopInteriorSerie, &s.CronCurrency, &s.CronPrestaShop, &s.SendGridKey, &s.EmailFrom, &s.NameFrom, &s.PalletWeight, &s.PalletWidth, &s.PalletHeight, &s.PalletDepth, &s.MaxConnections, &s.PrestashopStatusPaymentAccepted, &s.PrestashopStatusShipped, &s.MinimumStockSalesPeriods, &s.MinimumStockSalesDays, &s.CustomerJournal, &s.SalesJournal, &s.SalesAccount, &s.SupplierJournal, &s.PurchaseJournal, &s.PurchaseAccount, &s.EnableApiKey, &s.CronClearLabels, &s.LimitAccountingDate, &s.WooCommerceUrl, &s.WooCommerceConsumerKey, &s.WooCommerceConsumerSecret, &s.WooCommerceExportSerie, &s.WooCommerceIntracommunitySerie, &s.WooCommerceInteriorSerie, &s.WooCommerceDefaultPaymentMethod, &s.ConnectionLog, &s.FilterConnections, &s.ShopifyUrl, &s.ShopifyToken, &s.ShopifyExportSerie, &s.ShopifyIntracommunitySerie, &s.ShopifyInteriorSerie, &s.ShopifyDefaultPaymentMethod, &s.ShopifyShopLocationId, &s.EnterpriseKey, &s.PasswordMinimumLength, &s.PasswordMinumumComplexity, &s.InvoiceDeletePolicy, &s.TransactionLog, &s.UndoManufacturingOrderSeconds, &s.CronSendcloudTracking, &s.SMTPIdentity, &s.SMTPUsername, &s.SMTPPassword, &s.SMTPHostname, &s.SMTPSTARTTLS, &s.DefaultWarehouseName)
	return s
}

func (s *Settings) isValid() bool {
	return !(s.DefaultVatPercent < 0 || len(s.DefaultWarehouse) != 2 || len(s.DateFormat) == 0 || len(s.DateFormat) > 25 || len(s.EnterpriseName) == 0 || len(s.EnterpriseName) > 50 || len(s.EnterpriseDescription) > 250 || (s.Ecommerce != "_" && s.Ecommerce != "P" && s.Ecommerce != "M" && s.Ecommerce != "W" && s.Ecommerce != "S") || (s.Email != "_" && s.Email != "S" && s.Email != "T") || (s.Currency != "_" && s.Currency != "E") || len(s.CurrencyECBurl) > 100 || len(s.BarcodePrefix) > 4 || len(s.PrestaShopUrl) > 100 || len(s.PrestaShopApiKey) > 32 || s.PrestaShopLanguageId < 0 || len(s.CronCurrency) > 25 || len(s.CronPrestaShop) > 25 || len(s.SendGridKey) > 75 || len(s.EmailFrom) > 50 || len(s.NameFrom) > 50 || s.PalletWeight < 0 || s.PalletWidth < 0 || s.PalletHeight < 0 || s.PalletDepth < 0 || s.MaxConnections < 0 || s.PrestashopStatusPaymentAccepted < 0 || s.PrestashopStatusShipped < 0 || s.MinimumStockSalesPeriods < 0 || s.MinimumStockSalesDays < 0 || (s.Ecommerce == "P" && (s.PrestaShopLanguageId == 0 || s.PrestaShopExportSerie == nil || s.PrestaShopIntracommunitySerie == nil || s.PrestaShopInteriorSerie == nil || s.PrestashopStatusPaymentAccepted == 0 || s.PrestashopStatusShipped == 0)) || (s.Ecommerce == "W" && (s.WooCommerceDefaultPaymentMethod == nil || s.WooCommerceExportSerie == nil || s.WooCommerceInteriorSerie == nil || s.WooCommerceIntracommunitySerie == nil)) || (s.Ecommerce == "S" && (s.ShopifyDefaultPaymentMethod == nil || s.ShopifyExportSerie == nil || s.ShopifyInteriorSerie == nil || s.ShopifyIntracommunitySerie == nil)) || s.PasswordMinimumLength < 6 || (s.PasswordMinumumComplexity != "A" && s.PasswordMinumumComplexity != "B" && s.PasswordMinumumComplexity != "C" && s.PasswordMinumumComplexity != "D") || s.InvoiceDeletePolicy < 0 || s.InvoiceDeletePolicy > 2 || s.UndoManufacturingOrderSeconds < 0 || len(s.CronSendcloudTracking) > 25 || (s.Email == "S" && (len(s.SendGridKey) == 0 || len(s.EmailFrom) == 0 || len(s.NameFrom) == 0 || !emailIsValid(s.EmailFrom))) || (s.Email == "T" && (len(s.SMTPUsername) == 0 || len(s.SMTPPassword) == 0 || len(s.SMTPHostname) == 0 || !emailIsValid(s.SMTPUsername) || !hostnameWithPortValid(s.SMTPHostname))))
}

func (s *Settings) updateSettingsRecord() bool {
	if !s.isValid() {
		return false
	}

	var salesAccount *int32
	if s.SalesJournal != nil && *s.SalesJournal > 0 {
		acc := getAccountIdByAccountNumber(*s.SalesJournal, 1, s.Id)
		if acc > 0 {
			salesAccount = &acc
		}
	}
	var purchaseAccount *int32
	if s.PurchaseJournal != nil && *s.PurchaseJournal > 0 {
		acc := getAccountIdByAccountNumber(*s.PurchaseJournal, 1, s.Id)
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
	if s.Ecommerce != "_" {
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
	if settingsInMemory.CronClearLabels != s.CronClearLabels || settingsInMemory.Currency != s.Currency || settingsInMemory.CronCurrency != s.CronCurrency || settingsInMemory.Ecommerce != s.Ecommerce || settingsInMemory.CronPrestaShop != s.CronPrestaShop {
		refreshRunningCrons(settingsInMemory, *s)
	}

	sqlStatement := `UPDATE public.config SET default_vat_percent=$1, default_warehouse=$2, date_format=$3, enterprise_name=$4, enterprise_description=$5, ecommerce=$6, email=$7, currency=$8, currency_ecb_url=$9, barcode_prefix=$10, prestashop_url=$11, prestashop_api_key=$12, prestashop_language_id=$13, prestashop_export_serie=$14, prestashop_intracommunity_serie=$15, prestashop_interior_serie=$16, cron_currency=$17, cron_prestashop=$18, sendgrid_key=$19, email_from=$20, name_from=$21, pallet_weight=$22, pallet_width=$23, pallet_height=$24, pallet_depth=$25, max_connections=$26, prestashop_status_payment_accepted=$27, prestashop_status_shipped=$28, minimum_stock_sales_periods=$29, minimum_stock_sales_days=$30, customer_journal=$31, sales_journal=$32, sales_account=$33, supplier_journal=$34, purchase_journal=$35, purchase_account=$36, enable_api_key=$37, cron_clear_labels=$38, limit_accounting_date=$39, woocommerce_url=$40, woocommerce_consumer_key=$41, woocommerce_consumer_secret=$42, woocommerce_export_serie=$43, woocommerce_intracommunity_serie=$44, woocommerce_interior_serie=$45, woocommerce_default_payment_method=$46, connection_log=$47, filter_connections=$48, shopify_url=$49, shopify_token=$50, shopify_export_serie=$51, shopify_intracommunity_serie=$52, shopify_interior_serie=$53, shopify_default_payment_method=$54, shopify_shop_location_id=$55, password_minimum_length=$57, password_minumum_complexity=$58, invoice_delete_policy=$59, transaction_log=$60, undo_manufacturing_order_seconds=$61, cron_sendcloud_tracking=$62, smtp_identity=$63, smtp_username=$64, smtp_password=$65, smtp_hostname=$66, smtp_starttls=$67 WHERE id=$56`
	res, err := db.Exec(sqlStatement, s.DefaultVatPercent, s.DefaultWarehouse, s.DateFormat, s.EnterpriseName, s.EnterpriseDescription, s.Ecommerce, s.Email, s.Currency, s.CurrencyECBurl, s.BarcodePrefix, s.PrestaShopUrl, s.PrestaShopApiKey, s.PrestaShopLanguageId, s.PrestaShopExportSerie, s.PrestaShopIntracommunitySerie, s.PrestaShopInteriorSerie, s.CronCurrency, s.CronPrestaShop, s.SendGridKey, s.EmailFrom, s.NameFrom, s.PalletWeight, s.PalletWidth, s.PalletHeight, s.PalletDepth, s.MaxConnections, s.PrestashopStatusPaymentAccepted, s.PrestashopStatusShipped, s.MinimumStockSalesPeriods, s.MinimumStockSalesDays, s.CustomerJournal, s.SalesJournal, salesAccount, s.SupplierJournal, s.PurchaseJournal, purchaseAccount, s.EnableApiKey, s.CronClearLabels, s.LimitAccountingDate, s.WooCommerceUrl, s.WooCommerceConsumerKey, s.WooCommerceConsumerSecret, s.WooCommerceExportSerie, s.WooCommerceIntracommunitySerie, s.WooCommerceInteriorSerie, s.WooCommerceDefaultPaymentMethod, s.ConnectionLog, s.FilterConnections, s.ShopifyUrl, s.ShopifyToken, s.ShopifyExportSerie, s.ShopifyIntracommunitySerie, s.ShopifyInteriorSerie, s.ShopifyDefaultPaymentMethod, s.ShopifyShopLocationId, s.Id, s.PasswordMinimumLength, s.PasswordMinumumComplexity, s.InvoiceDeletePolicy, s.TransactionLog, s.UndoManufacturingOrderSeconds, s.CronSendcloudTracking, s.SMTPIdentity, s.SMTPUsername, s.SMTPPassword, s.SMTPHostname, s.SMTPSTARTTLS)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

// Don't allow every client to get the secret data, like API keys.
// This object holds the config that every client has to know, and the "Settings" object contains admin information.
type ClientSettings struct {
	DefaultVatPercent    float64 `json:"defaultVatPercent"`
	DefaultWarehouse     string  `json:"defaultWarehouse"`
	DefaultWarehouseName string  `json:"defaultWarehouseName"`
	DateFormat           string  `json:"dateFormat"`
	Ecommerce            string  `json:"ecommerce"`           // "_" = None, "P" = PrestaShop, "M" = Magento
	InvoiceDeletePolicy  int16   `json:"invoiceDeletePolicy"` // 0 = Allow invoice deletion, 1 = Only allow the deletion of the latest invoice in the billing serie, 2 = Never allow invoice deletion
}

func (s Settings) censorSettings() ClientSettings {
	c := ClientSettings{}
	c.DefaultVatPercent = s.DefaultVatPercent
	c.DefaultWarehouse = s.DefaultWarehouse
	c.DefaultWarehouseName = s.DefaultWarehouseName
	c.DateFormat = s.DateFormat
	c.Ecommerce = s.Ecommerce
	c.InvoiceDeletePolicy = s.InvoiceDeletePolicy
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
	if oldSettings.Ecommerce != newSettings.Ecommerce || oldSettings.CronPrestaShop != newSettings.CronPrestaShop {
		if enterpriseCronInfo.CronPrestaShop != nil {
			c.Remove(*enterpriseCronInfo.CronPrestaShop)
		}
		if newSettings.Ecommerce != "_" {
			e := ECommerce{Enterprise: oldSettings.Id}
			cronId, err := c.AddFunc(newSettings.CronPrestaShop, e.ecommerceControllerImportFromEcommerce)
			if err != nil {
				enterpriseCronInfo.CronPrestaShop = &cronId
			}
		}
	}
	if oldSettings.CronSendcloudTracking != newSettings.CronSendcloudTracking {
		if enterpriseCronInfo.CronSendcloudTracking != nil {
			c.Remove(*enterpriseCronInfo.CronSendcloudTracking)
		}
		if newSettings.CronSendcloudTracking != "" {
			cronId, err := c.AddFunc(newSettings.CronSendcloudTracking, func() {
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

func addEnterpriseFromParameters() bool {
	enterpriseKey, ok := getParameterValue("enterprise_key")
	if !ok {
		return false
	}
	enterpriseName, ok := getParameterValue("enterprise_name")
	if !ok {
		return false
	}
	enterpriseDesc, ok := getParameterValue("enterprise_desc")
	if !ok {
		return false
	}
	userPassword, ok := getParameterValue("user_password")
	if !ok {
		return false
	}
	licenseCode, ok := getParameterValue("license_code")
	if !ok {
		return false
	}
	licenseChance, ok := getParameterValue("license_chance")
	if !ok {
		return false
	}

	return createNewEnterprise(enterpriseName, enterpriseDesc, enterpriseKey, licenseCode, licenseChance, userPassword)
}

func createNewEnterprise(enterpriseName string, enterpriseDesc string, enterpriseKey string, licenseCode string, licenseChance string, userPassword string) bool {
	if len(enterpriseKey) == 0 || len(enterpriseName) == 0 || len(userPassword) < 8 || len(licenseCode) == 0 || len(licenseChance) == 0 {
		return false
	}

	ok, enterpriseId := initialConfigCreateEnterprise(enterpriseName, enterpriseDesc, strings.ToUpper(enterpriseKey))
	if !ok || enterpriseId <= 0 {
		return false
	}

	initialData(enterpriseId)

	sqlStatement := `UPDATE config SET default_warehouse=$1 WHERE id=$2`
	db.Exec(sqlStatement, "W1", enterpriseId)

	config := getSettingsRecordById(enterpriseId)
	ecommerceExportSerie := "EXP"
	ecommerceIntracommunitySerie := "IEU"
	ecommerceInteriorSerie := "INT"
	config.PrestaShopExportSerie = &ecommerceExportSerie
	config.PrestaShopIntracommunitySerie = &ecommerceIntracommunitySerie
	config.PrestaShopInteriorSerie = &ecommerceInteriorSerie
	config.PrestashopStatusPaymentAccepted = 2
	config.PrestashopStatusShipped = 4
	config.WooCommerceExportSerie = &ecommerceExportSerie
	config.WooCommerceIntracommunitySerie = &ecommerceIntracommunitySerie
	config.WooCommerceInteriorSerie = &ecommerceInteriorSerie
	config.ShopifyExportSerie = &ecommerceExportSerie
	config.ShopifyIntracommunitySerie = &ecommerceIntracommunitySerie
	config.ShopifyInteriorSerie = &ecommerceInteriorSerie
	if !config.updateSettingsRecord() {
		return false
	}

	activation := ServerSettingsActivation{
		LicenseCode: licenseCode,
		Chance:      &licenseChance,
	}
	settings.Server.Activation[enterpriseKey] = activation
	settings.setBackendSettings()
	if !activation.activateEnterprise(enterpriseId) {
		return false
	}

	insert := UserInsert{
		Username: "marketnet",
		FullName: "MARKETNET ADMINISTRATOR",
		Password: userPassword,
		Language: "en",
	}
	if !insert.insertUser(enterpriseId) {
		return false
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
		enterprise:    enterpriseId,
	}
	if !group.insertGroup() {
		return false
	}

	users := getUser(enterpriseId)
	user := users[len(users)-1]

	ug := UserGroup{
		User:  user.Id,
		Group: group.Id,
	}
	return ug.insertUserGroup()
}
