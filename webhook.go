package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// created for every enterprise
type WebHookSettings struct {
	Id                                 int32    `json:"id" gorm:"index:webhook_settings_id_enterprise,unique:true,priority:1"`
	EnterpriseId                       int32    `json:"-" gorm:"column:enterprise;not null:true;index:webhook_settings_id_enterprise,unique:true,priority:2"`
	Enterprise                         Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Url                                string   `json:"url" gorm:"column:url;not null:true;type:character varying(255)"`
	AuthCode                           string   `json:"authCode" gorm:"type:uuid;not null:true"`
	AuthMethod                         string   `json:"authMethod" gorm:"type:character (1);not null:true"` // H = Header, P = Parameter
	SaleOrders                         bool     `json:"saleOrders" gorm:"type:boolean;not null:true"`
	SaleOrderDetails                   bool     `json:"saleOrderDetails" gorm:"type:boolean;not null:true"`
	SaleOrderDetailsDigitalProductData bool     `json:"saleOrderDetailsDigitalProductData" gorm:"type:boolean;not null:true"`
	SaleInvoices                       bool     `json:"saleInvoices" gorm:"type:boolean;not null:true"`
	SaleInvoiceDetails                 bool     `json:"saleInvoiceDetails" gorm:"type:boolean;not null:true"`
	SaleDeliveryNotes                  bool     `json:"saleDeliveryNotes" gorm:"type:boolean;not null:true"`
	PurchaseOrders                     bool     `json:"purchaseOrders" gorm:"type:boolean;not null:true"`
	PurchaseOrderDetails               bool     `json:"purchaseOrderDetails" gorm:"type:boolean;not null:true"`
	PurchaseInvoices                   bool     `json:"purchaseInvoices" gorm:"type:boolean;not null:true"`
	PurchaseInvoiceDetails             bool     `json:"purchaseInvoiceDetails" gorm:"type:boolean;not null:true"`
	PurchaseDeliveryNotes              bool     `json:"purchaseDeliveryNotes" gorm:"type:boolean;not null:true"`
	Customers                          bool     `json:"customers" gorm:"type:boolean;not null:true"`
	Suppliers                          bool     `json:"suppliers" gorm:"type:boolean;not null:true"`
	Products                           bool     `json:"products" gorm:"type:boolean;not null:true"`
}

func (s *WebHookSettings) TableName() string {
	return "webhook_settings"
}

func getWebHookSettings(enterpriseId int32) []WebHookSettings {
	var settings []WebHookSettings = make([]WebHookSettings, 0)
	// get all webhooks settings from the database for the given enterprise sorted by id ascending using dbOrm
	dbOrm.Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&settings)
	return settings
}

func (s *WebHookSettings) isValid() bool {
	return !(len(s.Url) == 0 || len(s.Url) > 255 || (s.AuthMethod != "H" && s.AuthMethod != "P"))
}

func (s *WebHookSettings) BeforeCreate(tx *gorm.DB) (err error) {
	var webHookSettings WebHookSettings
	tx.Model(&WebHookSettings{}).Last(&webHookSettings)
	s.Id = webHookSettings.Id + 1
	return nil
}

func (s *WebHookSettings) insertWebHookSettings(enterpriseId int32) bool {
	if !s.isValid() {
		return false
	}

	webHooks := getWebHookSettings(enterpriseId)
	if settings.Server.MaxWebHooksPerEnterprise > 0 && len(webHooks) > int(settings.Server.MaxWebHooksPerEnterprise) {
		return false
	}

	s.AuthCode = uuid.New().String()
	s.EnterpriseId = enterpriseId

	// insert the webhook settings into the database using dbOrm
	result := dbOrm.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *WebHookSettings) updateWebHookSettings(enterpriseId int32) bool {
	if !s.isValid() || s.Id <= 0 {
		return false
	}

	// get a single webhook settings from the database by id and enterprise using dbOrm
	var webHookSettings WebHookSettings
	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, enterpriseId).First(&webHookSettings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// copy all the fields from the webhook settings to the webhook settings in the database
	webHookSettings.Url = s.Url
	webHookSettings.AuthMethod = s.AuthMethod
	webHookSettings.SaleOrders = s.SaleOrders
	webHookSettings.SaleOrderDetails = s.SaleOrderDetails
	webHookSettings.SaleOrderDetailsDigitalProductData = s.SaleOrderDetailsDigitalProductData
	webHookSettings.SaleInvoices = s.SaleInvoices
	webHookSettings.SaleInvoiceDetails = s.SaleInvoiceDetails
	webHookSettings.SaleDeliveryNotes = s.SaleDeliveryNotes
	webHookSettings.PurchaseOrders = s.PurchaseOrders
	webHookSettings.PurchaseOrderDetails = s.PurchaseOrderDetails
	webHookSettings.PurchaseInvoices = s.PurchaseInvoices
	webHookSettings.PurchaseInvoiceDetails = s.PurchaseInvoiceDetails
	webHookSettings.PurchaseDeliveryNotes = s.PurchaseDeliveryNotes
	webHookSettings.Customers = s.Customers
	webHookSettings.Suppliers = s.Suppliers
	webHookSettings.Products = s.Products

	// update the webhook settings in the database using dbOrm
	result = dbOrm.Save(&webHookSettings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *WebHookSettings) deleteWebHookSettings(enterpriseId int32) bool {
	if s.Id <= 0 {
		return false
	}

	// delete all the webhook logs for the webhook settings and enterprise id using dbOrm
	result := dbOrm.Where("webhook = ? AND enterprise = ?", s.Id, enterpriseId).Delete(&WebHookLog{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// delete all the webhook request for the webhook settings and enterprise id using dbOrm
	result = dbOrm.Where("webhook = ? AND enterprise = ?", s.Id, enterpriseId).Delete(&WebHookRequest{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// delete a single webhook settings from the database by id and enterprise using dbOrm
	result = dbOrm.Where("id = ? AND enterprise = ?", s.Id, enterpriseId).Delete(&WebHookSettings{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (s *WebHookSettings) renewAuthToken(enterpriseId int32) string {
	if s.Id <= 0 {
		return ""
	}

	// get a single webhook settings from the database by id and enterprise using dbOrm
	var webHookSettings WebHookSettings
	result := dbOrm.Where("id = ? AND enterprise = ?", s.Id, enterpriseId).First(&webHookSettings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	s.AuthCode = uuid.New().String()
	webHookSettings.AuthCode = s.AuthCode

	// update the webhook settings in the database using dbOrm
	result = dbOrm.Save(&webHookSettings)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}
	return s.AuthCode
}

// ALWAYS CALL IN A NEW THREAD - INTERNAL USE ONLY
// Send webhook calls, is required.
// enterpriseId,
// resource: for example "sales_order" or "customer",
// method: GET / POST / PUT / DELETE,
// data: JSON string to send
func fireWebHook(enterpriseId int32, resource string, method string, data string) {
	webhooks := getWebHookSettings(enterpriseId)
	for i := 0; i < len(webhooks); i++ {
		w := webhooks[i]
		r := WebHookRequest{
			EnterpriseId: enterpriseId,
			WebHookId:    w.Id,
			Url:          w.Url + "/" + resource,
			AuthCode:     w.AuthCode,
			AuthMethod:   w.AuthMethod,
			Send:         data,
			Method:       method,
		}
		r.sendWebHookRequest(true)
	}
}

// logs a call made to the customer's web service
type WebHookLog struct {
	Id               int64           `json:"id"`
	WebHookId        int32           `json:"webhookId" gorm:"column:webhook;not null:true"`
	WebHook          WebHookSettings `json:"-" gorm:"foreignKey:WebHookId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId     int32           `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise       Settings        `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Url              string          `json:"url" gorm:"column:url;not null:true;type:character varying(255)"`
	AuthCode         string          `json:"authCode" gorm:"type:uuid;not null:true"`
	AuthMethod       string          `json:"authMethod" gorm:"type:character (1);not null:true"` // H = Header, P = Parameter
	DateCreated      time.Time       `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	Sent             string          `json:"sent" gorm:"type:text;not null:true"`
	Received         string          `json:"received" gorm:"type:text;not null:true"`
	ReceivedHttpCode int16           `json:"receivedHttpCode" gorm:"column:received_http_code;not null:true"`
	Method           string          `json:"method" gorm:"type:character varying (10);not null:true"`
}

func (l *WebHookLog) TableName() string {
	return "webhook_logs"
}

func getWebHookLogs(enterpriseId int32, webHookId int32) []WebHookLog {
	var logs []WebHookLog = make([]WebHookLog, 0)
	// get all webhook logs from the database using dbOrm
	result := dbOrm.Where("enterprise = ? AND webhook = ?", enterpriseId, webHookId).Order("id DESC").Find(&logs)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return logs
	}

	return logs
}

func (l *WebHookLog) BeforeCreate(tx *gorm.DB) (err error) {
	var webHookLog WebHookLog
	tx.Model(&WebHookLog{}).Last(&webHookLog)
	l.Id = webHookLog.Id + 1
	return nil
}

// INTERNAL USE ONLY
func (l *WebHookLog) insertWebHookLog() bool {
	// insert a single webhook log into the database using dbOrm
	result := dbOrm.Create(&l)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

// Represents a request to be made.
// It can be sent inmediately, or it can be stored in the webhook_queue table waiting to be sent.
type WebHookRequest struct {
	Id           string          `json:"id" gorm:"primary_key;type:uuid;not null:true"`
	WebHookId    int32           `json:"webhookId" gorm:"column:webhook;not null:true"`
	WebHook      WebHookSettings `json:"-" gorm:"foreignKey:WebHookId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId int32           `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise   Settings        `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Url          string          `json:"url" gorm:"column:url;not null:true;type:character varying(255)"`
	AuthCode     string          `json:"authCode" gorm:"type:uuid;not null:true"`
	AuthMethod   string          `json:"authMethod" gorm:"type:character (1);not null:true"` // H = Header, P = Parameter
	DateCreated  time.Time       `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	Send         string          `json:"send" gorm:"type:text;not null:true"`
	Method       string          `json:"method" gorm:"type:character varying (10);not null:true"`
}

func (r *WebHookRequest) TableName() string {
	return "webhook_queue"
}

func getWebHookRequestQueue(enterpriseId int32, webHookId int32) []WebHookRequest {
	var requests []WebHookRequest = make([]WebHookRequest, 0)
	// get all webhook requests from the database using dbOrm
	result := dbOrm.Where("enterprise = ? AND webhook = ?", enterpriseId, webHookId).Order("date_created ASC").Find(&requests)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return requests
	}
	return requests
}

// INTERNAL USE ONLY
func getWebHookRequestQueueInternal(webHookId int32) []WebHookRequest {
	var requests []WebHookRequest = make([]WebHookRequest, 0)
	// get all webhook requests from the database using dbOrm
	result := dbOrm.Where("webhook = ?", webHookId).Order("date_created ASC").Find(&requests)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return requests
	}

	return requests
}

// Send data to the customer's WebHook
// storeWhenError: true when sending the request the first time, false when calling from the queue
func (r *WebHookRequest) sendWebHookRequest(storeWhenError bool) bool {
	if r.AuthMethod == "P" && strings.Contains(r.Url, "?accessToken=") {
		r.Url += "?accessToken=" + r.AuthCode
	}
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewBuffer([]byte(r.Send)))
	if err != nil {
		if storeWhenError {
			r.queueWebHookRequest()
		}
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	if r.AuthMethod == "H" {
		req.Header.Set("X-Marketnet-Access-Token", r.AuthCode)
	}

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if storeWhenError {
			r.queueWebHookRequest()
		}
		return false
	}
	// get the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if storeWhenError {
			r.queueWebHookRequest()
		}
		return false
	}

	l := WebHookLog{
		EnterpriseId:     r.EnterpriseId,
		WebHookId:        r.WebHookId,
		Url:              r.Url,
		AuthCode:         r.AuthCode,
		AuthMethod:       r.AuthMethod,
		Sent:             r.Send,
		Received:         string(body),
		ReceivedHttpCode: int16(resp.StatusCode),
		Method:           r.Method,
	}
	l.insertWebHookLog()

	return true
}

func (r *WebHookRequest) queueWebHookRequest() bool {
	queue := getWebHookRequestQueueInternal(r.WebHookId)
	if len(queue) > int(settings.Server.MaxQueueSizePerWebHook) {
		return false
	}

	r.Id = uuid.NewString()
	r.DateCreated = time.Now()

	result := dbOrm.Create(&r)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (r *WebHookRequest) dequeueWebHookRequest() bool {
	result := dbOrm.Delete(&r)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

// CRON TASK
func attemptToSendQueuedWebHooks() {
	// get all webhook requests from the database using dbOrm
	rows, err := dbOrm.Model(&WebHookRequest{}).Select("webhook").Distinct().Rows()
	if err != nil {
		log("DB", err.Error())
		return
	}

	var webHookId int32
	var ok bool
	for rows.Next() {
		rows.Scan(&webHookId)

		queuedRequests := getWebHookRequestQueueInternal(webHookId)

		for i := 0; i < len(queuedRequests); i++ {
			ok = queuedRequests[i].sendWebHookRequest(false)
			if ok {
				queuedRequests[i].dequeueWebHookRequest()
			} else {
				break
			}
		}
	}
}
