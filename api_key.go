package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id                int32             `json:"id"`
	Name              string            `json:"name"`
	DateCreated       time.Time         `json:"dateCreated"`
	UserCreated       int32             `json:"userCreated"`
	Off               bool              `json:"off"`
	User              int32             `json:"user"`
	Token             *string           `json:"token"`
	Auth              string            `json:"auth"` // P = Parameter, H = Header, B = Basic Auth, R = Bearer
	BasicAuthUser     *string           `json:"basicAuthUser"`
	BasicAuthPassword *string           `json:"basicAuthPassword"`
	UserCreatedName   string            `json:"userCreatedName"`
	UserName          string            `json:"userName"`
	Permissions       ApiKeyPermissions `json:"permissions"`
	permissionsJson   string
	enterprise        int32
}

func getApiKeys(enterpriseId int32) []ApiKey {
	keys := make([]ApiKey, 0)
	sqlStatement := `SELECT *,(SELECT username FROM "user" WHERE "user".id=api_key.user_created),(SELECT username FROM "user" WHERE "user".id=api_key."user") FROM public.api_key WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return keys
	}
	defer rows.Close()

	for rows.Next() {
		a := ApiKey{}
		rows.Scan(&a.Id, &a.Name, &a.DateCreated, &a.UserCreated, &a.Off, &a.User, &a.Token, &a.enterprise, &a.Auth, &a.BasicAuthUser, &a.BasicAuthPassword, &a.permissionsJson, &a.UserCreatedName, &a.UserName)
		json.Unmarshal([]byte(a.permissionsJson), &a.Permissions)
		keys = append(keys, a)
	}

	return keys
}

func (a *ApiKey) isValid() bool {
	return !(len(a.Name) == 0 || len(a.Name) > 64 || a.User <= 0 || (a.Auth != "P" && a.Auth != "H" && a.Auth != "B" && a.Auth != "R"))
}

func (a *ApiKey) insertApiKey() bool {
	if !a.isValid() {
		return false
	}

	if a.Auth == "P" || a.Auth == "H" || a.Auth == "R" {
		uuid := uuid.New().String()
		a.Token = &uuid
	} else if a.Auth == "B" {
		basicAuthUser := generateRandomString(20)
		a.BasicAuthUser = &basicAuthUser
		time.Sleep(1000)
		basicAuthPassword := generateRandomString(20)
		a.BasicAuthPassword = &basicAuthPassword
	}

	permissionsJson, _ := json.Marshal(a.Permissions)
	a.permissionsJson = string(permissionsJson)
	sqlStatement := `INSERT INTO public.api_key(name, user_created, "user", token, enterprise, auth, basic_auth_user, basic_auth_password, permissions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.Exec(sqlStatement, a.Name, a.UserCreated, a.User, a.Token, a.enterprise, a.Auth, a.BasicAuthUser, a.BasicAuthPassword, a.permissionsJson)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (a *ApiKey) updateApiKey() bool {
	if a.Id <= 0 || a.enterprise <= 0 || len(a.Name) == 0 || len(a.Name) > 64 {
		return false
	}

	permissionsJson, _ := json.Marshal(a.Permissions)
	a.permissionsJson = string(permissionsJson)
	sqlStatement := `UPDATE public.api_key SET name=$3, permissions=$4 WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise, a.Name, a.permissionsJson)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func generateRandomString(length int) string {
	const CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890,.-;*_!@#$%&"
	str := ""

	s := rand.NewSource(time.Now().UnixNano() + (time.Now().Unix() % 1000))
	r := rand.New(s)
	for i := 0; i < length; i++ {
		str += string(CHARSET[r.Intn(len(CHARSET))])
	}

	return str
}

func (a *ApiKey) deleteApiKey() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.api_key WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

func (a *ApiKey) offApiKey() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.api_key SET off=NOT off WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

// checks if the api key exists.
// Checks for different authentication methods available on the API Keys
// OK, user id, enterprise id
func checkApiKey(r *http.Request) (bool, int32, int32, *ApiKeyPermissions) {
	// Header
	token := r.Header.Get("X-Marketnet-Access-Token")
	if len(token) == 36 {
		ok, userId, enterpriseId, permission := checkApiKeyByTokenAuthType(token, "H")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId, permission
		}
	}

	// Basic auth
	token = r.Header.Get("Authorization")
	basicAuth := strings.Split(token, " ")
	if len(basicAuth) == 2 && basicAuth[0] == "Basic" {
		basicAuth, _ := base64.StdEncoding.DecodeString(basicAuth[1])
		usernamePassword := strings.Split(string(basicAuth), ":")
		if len(usernamePassword) == 2 && len(usernamePassword[0]) == 20 && len(usernamePassword[1]) == 20 {
			ok, userId, enterpriseId, permission := checkApiKeyByBasicAuthType(usernamePassword[0], usernamePassword[1])
			if ok && userId > 0 {
				return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId, permission
			}
		}
	} else if len(basicAuth) == 2 && basicAuth[0] == "Bearer" {
		ok, userId, enterpriseId, permission := checkApiKeyByTokenAuthType(basicAuth[1], "R")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId, permission
		}
	}

	// Parameter
	tokenParam, ok := r.URL.Query()["token"]
	if ok && len(tokenParam) > 0 && len(tokenParam[0]) == 36 {
		ok, userId, enterpriseId, permission := checkApiKeyByTokenAuthType(tokenParam[0], "P")
		if ok && userId > 0 {
			return checkMaxRequestsPerEnterprise(enterpriseId), userId, enterpriseId, permission
		}
	}

	// Givig up...
	return false, 0, 0, nil
}

func checkMaxRequestsPerEnterprise(enterpriseId int32) bool {
	// Is the enterprise activated? If the enterprise does not exist or there are no license to connect client, give up
	s := getSettingsRecordById(enterpriseId)
	if s.Id <= 0 || s.MaxConnections <= 0 {
		return false
	}

	requestsMade, ok := requestsPerMinuteEnterprise[enterpriseId]
	if !ok {
		requestsPerMinuteEnterprise[enterpriseId] = 1
		return true
	}
	if requestsMade >= settings.Server.MaxRequestsPerMinuteEnterprise {
		return false
	} else {
		requestsPerMinuteEnterprise[enterpriseId] = requestsPerMinuteEnterprise[enterpriseId] + 1
		return true
	}
}

func resetMaxRequestsPerEnterprise() {
	requestsPerMinuteEnterprise = make(map[int32]int32)
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// P = Parameter, H = Header, B = Basic Auth
// OK, user id, enterprise id
func checkApiKeyByTokenAuthType(token string, auth string) (bool, int32, int32, *ApiKeyPermissions) {
	sqlStatement := `SELECT "user",enterprise,permissions FROM public.api_key WHERE off=false AND token=$1 AND auth=$2`
	row := db.QueryRow(sqlStatement, token, auth)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0, 0, nil
	}

	var userId int32
	var enterpriseId int32
	var permissions string
	row.Scan(&userId, &enterpriseId, &permissions)

	p := ApiKeyPermissions{}
	if len(permissions) > 0 {
		json.Unmarshal([]byte(permissions), &p)
	}

	return true, userId, enterpriseId, &p
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// P = Parameter, H = Header, B = Basic Auth
// OK, user id, enterprise id
func checkApiKeyByBasicAuthType(basicAuthUser string, basicAuthPassword string) (bool, int32, int32, *ApiKeyPermissions) {
	sqlStatement := `SELECT "user",enterprise FROM public.api_key WHERE off=false AND basic_auth_user=$1 AND basic_auth_password=$2 AND auth='B'`
	row := db.QueryRow(sqlStatement, basicAuthUser, basicAuthPassword)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false, 0, 0, nil
	}

	var userId int32
	var enterpriseId int32
	var permissions string
	row.Scan(&userId, &enterpriseId, &permissions)

	p := ApiKeyPermissions{}
	if len(permissions) > 0 {
		json.Unmarshal([]byte(permissions), &p)
	}

	return true, userId, enterpriseId, &p
}

type ApiKeyPermissions struct {
	SaleOrders                           ApiKeyPermission `json:"saleOrders"`
	SaleOrderDetails                     ApiKeyPermission `json:"saleOrderDetails"`
	SaleOrderDetailsDigitalProductData   ApiKeyPermission `json:"saleOrderDetailsDigitalProductData"`
	SaleInvoices                         ApiKeyPermission `json:"saleInvoices"`
	SaleInvoiceDetails                   ApiKeyPermission `json:"saleInvoiceDetails"`
	SaleDeliveryNotes                    ApiKeyPermission `json:"saleDeliveryNotes"`
	PurchaseOrders                       ApiKeyPermission `json:"purchaseOrders"`
	PurchaseOrderDetails                 ApiKeyPermission `json:"purchaseOrderDetails"`
	PurchaseInvoices                     ApiKeyPermission `json:"purchaseInvoices"`
	PurchaseInvoiceDetails               ApiKeyPermission `json:"purchaseInvoiceDetails"`
	PurchaseDeliveryNotes                ApiKeyPermission `json:"purchaseDeliveryNotes"`
	Customers                            ApiKeyPermission `json:"customers"`
	Suppliers                            ApiKeyPermission `json:"suppliers"`
	Products                             ApiKeyPermission `json:"products"`
	Countries                            ApiKeyPermission `json:"countries"`
	States                               ApiKeyPermission `json:"states"`
	Colors                               ApiKeyPermission `json:"colors"`
	ProductFamilies                      ApiKeyPermission `json:"productFamilies"`
	Addresses                            ApiKeyPermission `json:"addresses"`
	Carriers                             ApiKeyPermission `json:"carriers"`
	BillingSeries                        ApiKeyPermission `json:"billingSeries"`
	Currencies                           ApiKeyPermission `json:"currencies"`
	PaymentMethods                       ApiKeyPermission `json:"paymentMethods"`
	Languages                            ApiKeyPermission `json:"languages"`
	Packages                             ApiKeyPermission `json:"packages"`
	Incoterms                            ApiKeyPermission `json:"incoterms"`
	Warehouses                           ApiKeyPermission `json:"warehouses"`
	WarehouseMovements                   ApiKeyPermission `json:"warehouseMovements"`
	ManufacturingOrders                  ApiKeyPermission `json:"manufacturingOrders"`
	ManufacturingOrderTypes              ApiKeyPermission `json:"manufacturingOrderTypes"`
	ComplexManufacturingOrders           ApiKeyPermission `json:"complexManufacturingOrders"`
	ComplexManufacturingOrdersComponents ApiKeyPermission `json:"complexManufacturingOrdersComponents"`
	ManufacturingOrderTypeComponents     ApiKeyPermission `json:"manufacturingOrderTypeComponents"`
	Shippings                            ApiKeyPermission `json:"shippings"`
	ShippingStatusHistory                ApiKeyPermission `json:"shippingStatusHistory"`
	Stock                                ApiKeyPermission `json:"stock"`
	Journal                              ApiKeyPermission `json:"journal"`
	Account                              ApiKeyPermission `json:"account"`
	AccountingMovement                   ApiKeyPermission `json:"accountingMovement"`
	AccountingMovementDetail             ApiKeyPermission `json:"accountingMovementDetail"`
	CollectionOperation                  ApiKeyPermission `json:"collectionOperation"`
	Charges                              ApiKeyPermission `json:"charges"`
	PaymentTransaction                   ApiKeyPermission `json:"paymentTransaction"`
	Payment                              ApiKeyPermission `json:"payment"`
	PostSaleInvoice                      ApiKeyPermission `json:"postSaleInvoice"`
	PostPurchaseInvoice                  ApiKeyPermission `json:"postPurchaseInvoice"`
}

type ApiKeyPermission struct {
	Get    bool `json:"get"`
	Post   bool `json:"post"`
	Put    bool `json:"put"`
	Delete bool `json:"delete"`
}
