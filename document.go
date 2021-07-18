package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DocumentAccessToken struct {
	Uuid        string    `json:"token"`
	DateCreated time.Time `json:"dateCreated"`
}

var documentAccessTokens []DocumentAccessToken = make([]DocumentAccessToken, 0)

type Document struct {
	Id                   int32     `json:"id"`
	Name                 string    `json:"name"`
	Uuid                 string    `json:"uuid"`
	DateCreated          time.Time `json:"dateCreated"`
	DateUpdated          time.Time `json:"dateUpdated"`
	Size                 int32     `json:"size"`
	Container            int16     `json:"container"`
	Description          string    `json:"description"`
	SalesOrder           *int32    `json:"salesOrder"`
	SalesInvoice         *int32    `json:"salesInvoice"`
	SalesDeliveryNote    *int32    `json:"salesDeliveryNote"`
	Shipping             *int32    `json:"shipping"`
	PurchaseOrder        *int32    `json:"purchaseOrder"`
	PurchaseInvoice      *int32    `json:"purchaseInvoice"`
	PurchaseDeliveryNote *int32    `json:"purchaseDeliveryNote"`
	MimeType             string    `json:"mimeType"`
}

func getDocuments() []Document {
	var document []Document = make([]Document, 0)
	sqlStatement := `SELECT * FROM document ORDER BY id DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return document
	}
	for rows.Next() {
		d := Document{}
		rows.Scan(&d.Id, &d.Name, &d.Uuid, &d.DateCreated, &d.DateUpdated, &d.Size, &d.Container, &d.Description, &d.SalesOrder, &d.SalesInvoice, &d.SalesDeliveryNote, &d.Shipping, &d.PurchaseOrder, &d.PurchaseInvoice, &d.PurchaseDeliveryNote, &d.MimeType)
		document = append(document, d)
	}

	return document
}

func (d *Document) getDocumentsRelations() []Document {
	var document []Document = make([]Document, 0)
	var rows *sql.Rows
	var err error
	if d.SalesOrder != nil {
		sqlStatement := `SELECT * FROM document WHERE sales_order=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.SalesOrder)
	} else if d.SalesInvoice != nil {
		sqlStatement := `SELECT * FROM document WHERE sales_invoice=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.SalesInvoice)
	} else if d.SalesDeliveryNote != nil {
		sqlStatement := `SELECT * FROM document WHERE sales_delivery_note=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.SalesDeliveryNote)
	} else if d.Shipping != nil {
		sqlStatement := `SELECT * FROM document WHERE shipping=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.Shipping)
	} else if d.PurchaseOrder != nil {
		sqlStatement := `SELECT * FROM document WHERE purchase_order=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.PurchaseOrder)
	} else if d.PurchaseInvoice != nil {
		sqlStatement := `SELECT * FROM document WHERE purchase_invoice=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.PurchaseInvoice)
	} else if d.PurchaseDeliveryNote != nil {
		sqlStatement := `SELECT * FROM document WHERE purchase_delivery_note=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, d.PurchaseDeliveryNote)
	} else {
		return document
	}
	if err != nil {
		log("DB", err.Error())
		return document
	}
	for rows.Next() {
		d := Document{}
		rows.Scan(&d.Id, &d.Name, &d.Uuid, &d.DateCreated, &d.DateUpdated, &d.Size, &d.Container, &d.Description, &d.SalesOrder, &d.SalesInvoice, &d.SalesDeliveryNote, &d.Shipping, &d.PurchaseOrder, &d.PurchaseInvoice, &d.PurchaseDeliveryNote, &d.MimeType)
		document = append(document, d)
	}

	return document
}

func getDocumentRow(uuid string) Document {
	sqlStatement := `SELECT * FROM document WHERE uuid=$1`
	row := db.QueryRow(sqlStatement, uuid)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Document{}
	}

	d := Document{}
	row.Scan(&d.Id, &d.Name, &d.Uuid, &d.DateCreated, &d.DateUpdated, &d.Size, &d.Container, &d.Description, &d.SalesOrder, &d.SalesInvoice, &d.SalesDeliveryNote, &d.Shipping, &d.PurchaseOrder, &d.PurchaseInvoice, &d.PurchaseDeliveryNote, &d.MimeType)

	return d
}

func getDocumentRowById(id int32) Document {
	sqlStatement := `SELECT * FROM document WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Document{}
	}

	d := Document{}
	row.Scan(&d.Id, &d.Name, &d.Uuid, &d.DateCreated, &d.DateUpdated, &d.Size, &d.Container, &d.Description, &d.SalesOrder, &d.SalesInvoice, &d.SalesDeliveryNote, &d.Shipping, &d.PurchaseOrder, &d.PurchaseInvoice, &d.PurchaseDeliveryNote, &d.MimeType)

	return d
}

func (d *Document) isValid() bool {
	return !(len(d.Name) == 0 || len(d.Name) > 250 || d.Size <= 0 || d.Container <= 0)
}

func (d *Document) insertDocument() bool {
	if !d.isValid() {
		return false
	}

	d.Uuid = uuid.New().String()
	sqlStatement := `INSERT INTO public.document(name, uuid, container, dsc, sales_order, sales_invoice, sales_delivery_note, shipping, purchase_order, purchase_invoice, purchase_delivery_note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	res, err := db.Exec(sqlStatement, d.Name, d.Uuid, d.Container, d.Description, d.SalesOrder, d.SalesInvoice, d.SalesDeliveryNote, d.Shipping, d.PurchaseOrder, d.PurchaseInvoice, d.PurchaseDeliveryNote)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (d *Document) deleteDocument() bool {
	if d.Id <= 0 {
		return false
	}

	inMemoryDocument := getDocumentRowById(d.Id)
	if inMemoryDocument.Id <= 0 {
		return false
	}
	container := getDocumentContainerRow(inMemoryDocument.Container)
	if container.Id <= 0 {
		return false
	}
	os.Remove(path.Join(container.Path, inMemoryDocument.Uuid))

	sqlStatement := `DELETE FROM public.document WHERE id=$1`
	res, err := db.Exec(sqlStatement, d.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
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
	ok := consumeToken(token)
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
		log("DB", err.Error())
		return nil, http.StatusInternalServerError
	}
	return content, http.StatusOK
}

func uploadDocument(token string, uuid string, document []byte) int {
	ok := consumeToken(token)
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
	err := ioutil.WriteFile(path.Join(container.Path, doc.Uuid), document, 0700)
	if err != nil {
		log("DB", err.Error())
		return http.StatusInternalServerError
	}

	mimeType := http.DetectContentType(document)

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

	sqlStatement := `UPDATE public.document SET date_updated=CURRENT_TIMESTAMP(3), size=$2, mime_type=$3 WHERE id=$1`
	db.Exec(sqlStatement, doc.Id, len(document), mimeType)
	return http.StatusOK
}

func grantDocumentAccessToken() DocumentAccessToken {
	t := DocumentAccessToken{}
	t.Uuid = uuid.New().String()
	t.DateCreated = time.Now()
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

func consumeToken(token string) bool {
	for i := 0; i < len(documentAccessTokens); i++ {
		if time.Until(documentAccessTokens[i].DateCreated).Seconds() <= 60 { // the token has not expired yet
			documentAccessTokens = append(documentAccessTokens[:i], documentAccessTokens[i+1:]...) // delete the token
			return true
		}
	}
	return false // the token was not found or is expired, let the cleaning function delete it
}
