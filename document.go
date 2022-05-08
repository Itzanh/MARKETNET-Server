package main

// TODO: better error control
// enviar el tamaño del archivo al servidor y verificar que cabrá
// si el servidor devuelve un error durante la subida, mostrar ese error en el fontend

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentAccessToken struct {
	Uuid        string    `json:"token"`
	DateCreated time.Time `json:"dateCreated"`
	Enterprise  int32
}

var documentAccessTokens []DocumentAccessToken = make([]DocumentAccessToken, 0)

type Document struct {
	Id                     int32                 `json:"id"`
	Name                   string                `json:"name" gorm:"type:character varying(250);not null:true"`
	Uuid                   string                `json:"uuid" gorm:"type:uuid;not null:true;index:document_uuid,unique:true"`
	DateCreated            time.Time             `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	DateUpdated            time.Time             `json:"dateUpdated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Size                   int32                 `json:"size" gorm:"not null:true"`
	Container              int32                 `json:"container" gorm:"not null:true"`
	Description            string                `json:"description" gorm:"column:dsc;type:text;not null:true"`
	SalesOrderId           *int64                `json:"salesOrder" gorm:"column:sales_order"`
	SalesOrder             *SaleOrder            `json:"-" gorm:"foreignKey:SalesOrderId,EnterpriseId;references:Id,EnterpriseId"`
	SalesInvoiceId         *int64                `json:"salesInvoice" gorm:"column:sales_invoice"`
	SalesInvoice           *SalesInvoice         `json:"-" gorm:"foreignKey:SalesInvoiceId,EnterpriseId;references:Id,EnterpriseId"`
	SalesDeliveryNoteId    *int64                `json:"salesDeliveryNote" gorm:"column:sales_delivery_note"`
	SalesDeliveryNote      *SalesDeliveryNote    `json:"-" gorm:"foreignKey:SalesDeliveryNoteId,EnterpriseId;references:Id,EnterpriseId"`
	ShippingId             *int64                `json:"shipping" gorm:"column:shipping"`
	Shipping               *Shipping             `json:"-" gorm:"foreignKey:ShippingId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseOrderId        *int64                `json:"purchaseOrder" gorm:"column:purchase_order"`
	PurchaseOrder          *PurchaseOrder        `json:"-" gorm:"foreignKey:PurchaseOrderId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseInvoiceId      *int64                `json:"purchaseInvoice" gorm:"column:purchase_invoice"`
	PurchaseInvoice        *PurchaseInvoice      `json:"-" gorm:"foreignKey:PurchaseInvoiceId,EnterpriseId;references:Id,EnterpriseId"`
	PurchaseDeliveryNoteId *int64                `json:"purchaseDeliveryNote" gorm:"column:purchase_delivery_note"`
	PurchaseDeliveryNote   *PurchaseDeliveryNote `json:"-" gorm:"foreignKey:PurchaseDeliveryNoteId,EnterpriseId;references:Id,EnterpriseId"`
	MimeType               string                `json:"mimeType" gorm:"type:character varying(100);not null:true"`
	EnterpriseId           int32                 `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise             Settings              `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (d *Document) TableName() string {
	return "document"
}

func getDocuments(enterpriseId int32) []Document {
	var document []Document = make([]Document, 0)
	dbOrm.Model(&Document{}).Where("enterprise = ?", enterpriseId).Order("id DESC").Find(&document)
	return document
}

func (d *Document) getDocumentsRelations(enterpriseId int32) []Document {
	var documents []Document = make([]Document, 0)

	docDB := getDocumentRowById(d.Id)
	if docDB.EnterpriseId != enterpriseId {
		return documents
	}

	var query string
	var interfaces []interface{} = make([]interface{}, 0)
	if d.SalesOrderId != nil {
		query = `sales_order = ?`
		interfaces = append(interfaces, d.SalesOrderId)
	} else if d.SalesInvoiceId != nil {
		query = `sales_invoice = ?`
		interfaces = append(interfaces, d.SalesInvoiceId)
	} else if d.SalesDeliveryNoteId != nil {
		query = `sales_delivery_note = ?`
		interfaces = append(interfaces, d.SalesDeliveryNoteId)
	} else if d.ShippingId != nil {
		query = `shipping = ?`
		interfaces = append(interfaces, d.ShippingId)
	} else if d.PurchaseOrderId != nil {
		query = `purchase_order = ?`
		interfaces = append(interfaces, d.PurchaseOrderId)
	} else if d.PurchaseInvoiceId != nil {
		query = `purchase_invoice = ?`
		interfaces = append(interfaces, d.PurchaseInvoiceId)
	} else if d.PurchaseDeliveryNoteId != nil {
		query = `purchase_delivery_note = ?`
		interfaces = append(interfaces, d.PurchaseDeliveryNoteId)
	} else {
		return documents
	}
	dbOrm.Model(&Document{}).Where(query, interfaces...).Order("id DESC").Find(&documents)
	return documents
}

func getDocumentRow(uuid string) Document {
	d := Document{}
	dbOrm.Model(&Document{}).Where("uuid = ?", uuid).First(&d)
	return d
}

func getDocumentRowById(id int32) Document {
	d := Document{}
	dbOrm.Model(&Document{}).Where("id = ?", id).First(&d)
	return d
}

func (d *Document) isValid() bool {
	return !(len(d.Name) == 0 || len(d.Name) > 250 || d.Size <= 0 || d.Container <= 0 || len(d.Description) > 3000)
}

func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	var document Document
	tx.Model(&Document{}).Last(&document)
	d.Id = document.Id + 1
	return nil
}

func (d *Document) insertDocument() OkAndErrorCodeReturn {
	if !d.isValid() {
		return OkAndErrorCodeReturn{Ok: false}
	}

	documentContainer := getDocumentContainerRow(d.Container)
	if documentContainer.MaxStorage > 0 && documentContainer.UsedStorage+int64(d.Size) > documentContainer.MaxStorage {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1, ExtraData: []string{
			strconv.Itoa(int(documentContainer.UsedStorage)),
			strconv.Itoa(int(documentContainer.MaxStorage)),
		}}
	}
	if d.Size > documentContainer.MaxFileSize {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2, ExtraData: []string{
			strconv.Itoa(int(documentContainer.MaxFileSize)),
		}}
	}
	if int64(d.Size) > settings.Server.WebSecurity.MaxRequestBodyLength {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{
			strconv.Itoa(int(settings.Server.WebSecurity.MaxRequestBodyLength)),
		}}
	}

	d.Uuid = uuid.New().String()
	d.DateCreated = time.Now()
	d.DateUpdated = time.Now()
	d.Size = 0
	d.MimeType = ""

	result := dbOrm.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	return OkAndErrorCodeReturn{Ok: true, ExtraData: []string{
		d.Uuid,
	}}
}

func (d *Document) deleteDocument() bool {
	if d.Id <= 0 {
		return false
	}

	inMemoryDocument := getDocumentRowById(d.Id)
	if inMemoryDocument.Id <= 0 || d.EnterpriseId != inMemoryDocument.EnterpriseId {
		return false
	}
	container := getDocumentContainerRow(inMemoryDocument.Container)
	if container.Id <= 0 {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	result := trans.Where("id = ? AND enterprise = ?", d.Id, d.EnterpriseId).Delete(&Document{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	if !container.updateUsedStorage(-inMemoryDocument.Size, trans) {
		trans.Rollback()
		return false
	}

	err := os.Remove(path.Join(container.Path, inMemoryDocument.Uuid))
	if err != nil {
		log("FS", err.Error())
		trans.Rollback()
		return false
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
	///
}

func handleDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type")
	switch r.Method {
	case "GET":
		uuid, ok := r.URL.Query()["uuid"]
		if !ok || len(uuid[0]) != 36 {
			return
		}
		token, ok := r.URL.Query()["token"]
		if !ok || len(token[0]) != 36 {
			return
		}
		content, statusCode := downloadDocument(token[0], uuid[0])
		w.WriteHeader(statusCode)
		w.Write(content)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		uuid, ok := r.URL.Query()["uuid"]
		if !ok || len(uuid[0]) != 36 {
			return
		}
		token, ok := r.URL.Query()["token"]
		if !ok || len(token[0]) != 36 {
			return
		}
		statusCode := uploadDocument(token[0], uuid[0], body)
		w.WriteHeader(statusCode)
	}
}

func downloadDocument(token string, uuid string) ([]byte, int) {
	ok, _ := consumeToken(token)
	if !ok {
		return nil, http.StatusUnauthorized
	}
	doc := getDocumentRow(uuid)
	if doc.Id <= 0 {
		return nil, http.StatusNotFound
	}
	container := getDocumentContainerRow(doc.Container)
	if container.Id <= 0 {
		return nil, http.StatusNotFound
	}
	content, err := ioutil.ReadFile(path.Join(container.Path, doc.Uuid))
	if err != nil {
		log("FS", err.Error())
		return nil, http.StatusInternalServerError
	}
	return content, http.StatusOK
}

func uploadDocument(token string, uuid string, document []byte) int {
	ok, _ := consumeToken(token)
	if !ok {
		return http.StatusUnauthorized
	}
	doc := getDocumentRow(uuid)
	if doc.Id <= 0 {
		return http.StatusNotFound
	}
	container := getDocumentContainerRow(doc.Container)
	if container.Id <= 0 {
		return http.StatusNotFound
	}

	if container.MaxStorage > 0 && container.UsedStorage+int64(len(document)) > container.MaxStorage {
		return http.StatusRequestEntityTooLarge
	}

	mimeType := http.DetectContentType(document)

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return http.StatusInternalServerError
	}
	///

	if !container.updateUsedStorage(-doc.Size, trans) {
		trans.Rollback()
		return http.StatusInternalServerError
	}

	doc.DateUpdated = time.Now()
	doc.Size = int32(len(document))
	doc.MimeType = mimeType

	result := trans.Save(&doc)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return http.StatusInternalServerError
	}
	if !container.updateUsedStorage(int32(len(document)), trans) {
		trans.Rollback()
		return http.StatusInternalServerError
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return http.StatusInternalServerError
	}
	///

	err := ioutil.WriteFile(path.Join(container.Path, doc.Uuid), document, 0700)
	if err != nil {
		log("FS", err.Error())
		return http.StatusInternalServerError
	}

	if len(container.AllowedMimeTypes) > 0 {
		allowedMimeTypes := strings.Split(container.AllowedMimeTypes, ",")
		var ok bool = false
		for i := 0; i < len(allowedMimeTypes); i++ {
			if allowedMimeTypes[i] == mimeType {
				ok = true
				break
			}
		}
		if !ok {
			return http.StatusNotAcceptable
		}
	} else {
		disallowedMimeTypes := strings.Split(container.DisallowedMimeTypes, ",")
		for i := 0; i < len(disallowedMimeTypes); i++ {
			if disallowedMimeTypes[i] == mimeType {
				return http.StatusNotAcceptable
			}
		}
	}

	return http.StatusOK
}

func grantDocumentAccessToken(enterpriseId int32) DocumentAccessToken {
	t := DocumentAccessToken{}
	t.Uuid = uuid.New().String()
	t.DateCreated = time.Now()
	t.Enterprise = enterpriseId
	documentAccessTokens = append(documentAccessTokens, t)
	return t
}

func cleanDocumentTokens() {
	for {
		time.Sleep(60000)
		for i := len(documentAccessTokens) - 1; i >= 0; i-- {
			if time.Until(documentAccessTokens[i].DateCreated).Seconds() > 60 {
				documentAccessTokens = append(documentAccessTokens[:i], documentAccessTokens[i+1:]...)
			}
		}
	}
}

func consumeToken(token string) (bool, int32) {
	for i := 0; i < len(documentAccessTokens); i++ {
		if time.Until(documentAccessTokens[i].DateCreated).Seconds() <= 60 { // the token has not expired yet
			enterpriseId := documentAccessTokens[i].Enterprise
			documentAccessTokens = append(documentAccessTokens[:i], documentAccessTokens[i+1:]...) // delete the token
			return true, enterpriseId
		}
	}
	return false, 0 // the token was not found or is expired, let the cleaning function delete it
}
