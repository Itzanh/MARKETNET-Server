package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// created for every enterprise
type WebHookSettings struct {
	Id                                 int32  `json:"id"`
	Url                                string `json:"url"`
	AuthCode                           string `json:"authCode"`
	AuthMethod                         string `json:"authMethod"` // H = Header, P = Parameter
	SaleOrders                         bool   `json:"saleOrders"`
	SaleOrderDetails                   bool   `json:"saleOrderDetails"`
	SaleOrderDetailsDigitalProductData bool   `json:"saleOrderDetailsDigitalProductData"`
	SaleInvoices                       bool   `json:"saleInvoices"`
	SaleInvoiceDetails                 bool   `json:"saleInvoiceDetails"`
	SaleDeliveryNotes                  bool   `json:"saleDeliveryNotes"`
	PurchaseOrders                     bool   `json:"purchaseOrders"`
	PurchaseOrderDetails               bool   `json:"purchaseOrderDetails"`
	PurchaseInvoices                   bool   `json:"purchaseInvoices"`
	PurchaseInvoiceDetails             bool   `json:"purchaseInvoiceDetails"`
	PurchaseDeliveryNotes              bool   `json:"purchaseDeliveryNotes"`
	Customers                          bool   `json:"customers"`
	Suppliers                          bool   `json:"suppliers"`
	Products                           bool   `json:"products"`
	enterprise                         int32
}

func getWebHookSettings(enterpriseId int32) []WebHookSettings {
	var settings []WebHookSettings = make([]WebHookSettings, 0)
	sqlStatement := `SELECT * FROM public.webhook_settings WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return settings
	}

	for rows.Next() {
		s := WebHookSettings{}
		rows.Scan(&s.Id, &s.enterprise, &s.Url, &s.AuthCode, &s.AuthMethod, &s.SaleOrders, &s.SaleOrderDetails, &s.SaleOrderDetailsDigitalProductData, &s.SaleInvoices, &s.SaleInvoiceDetails, &s.SaleDeliveryNotes, &s.PurchaseOrders, &s.PurchaseOrderDetails, &s.PurchaseInvoices, &s.PurchaseInvoiceDetails, &s.PurchaseDeliveryNotes, &s.Customers, &s.Suppliers, &s.Products)
		settings = append(settings, s)
	}

	return settings
}

func (s *WebHookSettings) isValid() bool {
	return !(len(s.Url) == 0 || len(s.Url) > 255 || (s.AuthMethod != "H" && s.AuthMethod != "P"))
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
	sqlStatement := `INSERT INTO public.webhook_settings(enterprise, url, auth_code, auth_method, sale_orders, sale_order_details, sale_order_details_digital_product_data, sale_invoices, sale_invoice_details, sale_delivery_notes, purchase_orders, purchase_order_details, purchase_invoices, purchase_invoice_details, purchase_delivery_notes, customers, suppliers, products) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := db.Exec(sqlStatement, enterpriseId, s.Url, s.AuthCode, s.AuthMethod, s.SaleOrders, s.SaleOrderDetails, s.SaleOrderDetailsDigitalProductData, s.SaleInvoices, s.SaleInvoiceDetails, s.SaleDeliveryNotes, s.PurchaseOrders, s.PurchaseOrderDetails, s.PurchaseInvoices, s.PurchaseInvoiceDetails, s.PurchaseDeliveryNotes, s.Customers, s.Suppliers, s.Products)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (s *WebHookSettings) updateWebHookSettings(enterpriseId int32) bool {
	if !s.isValid() || s.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.webhook_settings SET url=$2, auth_method=$3, sale_orders=$4, sale_order_details=$5, sale_order_details_digital_product_data=$6, sale_invoices=$7, sale_invoice_details=$8, sale_delivery_notes=$9, purchase_orders=$10, purchase_order_details=$11, purchase_invoices=$12, purchase_invoice_details=$13, purchase_delivery_notes=$14, customers=$15, suppliers=$16, products=$17 WHERE id=$1 AND enterprise=$18`
	_, err := db.Exec(sqlStatement, s.Id, s.Url, s.AuthMethod, s.SaleOrders, s.SaleOrderDetails, s.SaleOrderDetailsDigitalProductData, s.SaleInvoices, s.SaleInvoiceDetails, s.SaleDeliveryNotes, s.PurchaseOrders, s.PurchaseOrderDetails, s.PurchaseInvoices, s.PurchaseInvoiceDetails, s.PurchaseDeliveryNotes, s.Customers, s.Suppliers, s.Products, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (s *WebHookSettings) deleteWebHookSettings(enterpriseId int32) bool {
	if s.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.webhook_logs WHERE webhook = $1 AND enterprise = $2`
	_, err := db.Exec(sqlStatement, s.Id, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `DELETE FROM public.webhook_queue WHERE webhook = $1 AND enterprise = $2`
	_, err = db.Exec(sqlStatement, s.Id, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	sqlStatement = `DELETE FROM public.webhook_settings WHERE id = $1 AND enterprise = $2`
	_, err = db.Exec(sqlStatement, s.Id, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (s *WebHookSettings) renewAuthToken(enterpriseId int32) string {
	if s.Id <= 0 {
		return ""
	}

	s.AuthCode = uuid.New().String()
	sqlStatement := `UPDATE public.webhook_settings SET auth_code=$3 WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, s.Id, enterpriseId, s.AuthCode)
	if err != nil {
		log("DB", err.Error())
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
			enterprise: enterpriseId,
			WebHook:    w.Id,
			Url:        w.Url + "/" + resource,
			AuthCode:   w.AuthCode,
			AuthMethod: w.AuthMethod,
			Send:       data,
			Method:     method,
		}
		r.sendWebHookRequest(true)
	}
}

// logs a call made to the customer's web service
type WebHookLog struct {
	Id               int64     `json:"id"`
	WebHook          int32     `json:"webhook"`
	Url              string    `json:"url"`
	AuthCode         string    `json:"authCode"`
	AuthMethod       string    `json:"authMethod"` // H = Header, P = Parameter
	DateCreated      time.Time `json:"dateCreated"`
	Sent             string    `json:"sent"`
	Received         string    `json:"received"`
	ReceivedHttpCode int16     `json:"receivedHttpCode"`
	Method           string    `json:"method"`
	enterprise       int32
}

func getWebHookLogs(enterpriseId int32, webHookId int32) []WebHookLog {
	var logs []WebHookLog = make([]WebHookLog, 0)
	sqlStatement := `SELECT * FROM public.webhook_logs WHERE enterprise = $1 AND webhook = $2 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId, webHookId)
	if err != nil {
		log("DB", err.Error())
		return logs
	}

	for rows.Next() {
		l := WebHookLog{}
		rows.Scan(&l.Id, &l.WebHook, &l.enterprise, &l.Url, &l.AuthCode, &l.AuthMethod, &l.DateCreated, &l.Sent, &l.Received, &l.ReceivedHttpCode, &l.Method)
		logs = append(logs, l)
	}

	return logs
}

// INTERNAL USE ONLY
func (l *WebHookLog) insertWebHookLog() bool {
	sqlStatement := `INSERT INTO public.webhook_logs(webhook, enterprise, url, auth_code, auth_method, sent, received, received_http_code, method) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.Exec(sqlStatement, l.WebHook, l.enterprise, l.Url, l.AuthCode, l.AuthMethod, l.Sent, l.Received, l.ReceivedHttpCode, l.Method)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

// Represents a request to be made.
// It can be sent inmediately, or it can be stored in the webhook_queue table waiting to be sent.
type WebHookRequest struct {
	Id          string    `json:"id"`
	WebHook     int32     `json:"webhook"`
	Url         string    `json:"url"`
	AuthCode    string    `json:"authCode"`
	AuthMethod  string    `json:"authMethod"` // H = Header, P = Parameter
	DateCreated time.Time `json:"dateCreated"`
	Send        string    `json:"send"`
	Method      string    `json:"method"`
	enterprise  int32
}

func getWebHookRequestQueue(enterpriseId int32, webHookId int32) []WebHookRequest {
	var requests []WebHookRequest = make([]WebHookRequest, 0)
	sqlStatement := `SELECT * FROM webhook_queue WHERE webhook = $1 AND enterprise = $2 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, webHookId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return requests
	}

	for rows.Next() {
		r := WebHookRequest{}
		rows.Scan(&r.Id, &r.WebHook, &r.enterprise, &r.Url, &r.AuthCode, &r.AuthMethod, &r.DateCreated, &r.Send, &r.Method)
		requests = append(requests, r)
	}

	return requests
}

// INTERNAL USE ONLY
func getWebHookRequestQueueInternal(webHookId int32) []WebHookRequest {
	var requests []WebHookRequest = make([]WebHookRequest, 0)
	sqlStatement := `SELECT * FROM webhook_queue WHERE webhook = $1 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, webHookId)
	if err != nil {
		log("DB", err.Error())
		return requests
	}

	for rows.Next() {
		r := WebHookRequest{}
		rows.Scan(&r.Id, &r.WebHook, &r.enterprise, &r.Url, &r.AuthCode, &r.AuthMethod, &r.DateCreated, &r.Send, &r.Method)
		requests = append(requests, r)
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
		enterprise:       r.enterprise,
		WebHook:          r.WebHook,
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
	queue := getWebHookRequestQueueInternal(r.WebHook)
	if len(queue) > int(settings.Server.MaxQueueSizePerWebHook) {
		return false
	}

	sqlStatement := `INSERT INTO public.webhook_queue(webhook, enterprise, url, auth_code, auth_method, send, method) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(sqlStatement, r.WebHook, r.enterprise, r.Url, r.AuthCode, r.AuthMethod, r.Send, r.Method)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (r *WebHookRequest) dequeueWebHookRequest() bool {
	sqlStatement := `DELETE FROM public.webhook_queue WHERE id = $1`
	_, err := db.Exec(sqlStatement, r.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

// CRON TASK
func attemptToSendQueuedWebHooks() {
	sqlStatement := `SELECT DISTINCT webhook FROM public.webhook_queue`
	rows, err := db.Query(sqlStatement)
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
