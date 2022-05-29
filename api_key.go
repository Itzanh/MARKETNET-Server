package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiKey struct {
	Id                int32             `json:"id"`
	Name              string            `json:"name" gorm:"type:character varying(140);not null:true"`
	DateCreated       time.Time         `json:"dateCreated" gorm:"type:timestamp(3) without time zone;not null:true"`
	UserCreatedId     int32             `json:"userCreatedId" gorm:"column:user_created;not null:true"`
	UserCreated       User              `json:"userCreated" gorm:"foreignkey:UserCreatedId,EnterpriseId;references:Id,EnterpriseId"`
	Off               bool              `json:"off" gorm:"not null:true"`
	UserId            int32             `json:"userId" gorm:"column:user;not null:true"`
	User              User              `json:"user" gorm:"foreignkey:UserId,EnterpriseId;references:Id,EnterpriseId"`
	Token             *string           `json:"token" gorm:"type:uuid;index:api_key_token,unique:true,where:token is not null"`
	EnterpriseId      int32             `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise        Settings          `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Auth              string            `json:"auth" gorm:"type:character(1);not null:true"` // P = Parameter, H = Header, B = Basic Auth, R = Bearer
	BasicAuthUser     *string           `json:"basicAuthUser" gorm:"type:character varying(20);index:api_key_basic_auth,unique:true,where:auth = 'B'"`
	BasicAuthPassword *string           `json:"basicAuthPassword" gorm:"type:character varying(20);index:api_key_basic_auth,unique:true,where:auth = 'B'"`
	Permissions       ApiKeyPermissions `json:"permissions" gorm:"-"`
	PermissionsJson   string            `json:"-" gorm:"column:permissions;type:json;not null:true"`
}

func (ak *ApiKey) TableName() string {
	return "api_key"
}

func getApiKeys(enterpriseId int32) []ApiKey {
	keys := make([]ApiKey, 0)
	// get all the api keys for the enterprise sorted by id ascending
	result := dbOrm.Where("api_key.enterprise = ?", enterpriseId).Preload("UserCreated").Preload("User").Order("api_key.id ASC").Find(&keys)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	for i := 0; i < len(keys); i++ {
		json.Unmarshal([]byte(keys[i].PermissionsJson), &keys[i].Permissions)
	}
	return keys
}

func (a *ApiKey) isValid() bool {
	return !(len(a.Name) == 0 || len(a.Name) > 140 || a.UserId <= 0 || (a.Auth != "P" && a.Auth != "H" && a.Auth != "B" && a.Auth != "R"))
}

func (a *ApiKey) BeforeCreate(tx *gorm.DB) (err error) {
	var apiKey ApiKey
	tx.Model(&ApiKey{}).Last(&apiKey)
	a.Id = apiKey.Id + 1
	return nil
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
	a.PermissionsJson = string(permissionsJson)

	a.DateCreated = time.Now()

	result := dbOrm.Create(&a)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (a *ApiKey) updateApiKey() bool {
	if a.Id <= 0 || a.EnterpriseId <= 0 || len(a.Name) == 0 || len(a.Name) > 64 {
		return false
	}

	permissionsJson, _ := json.Marshal(a.Permissions)
	a.PermissionsJson = string(permissionsJson)

	// get a single api key from the database using dbOrm where the id is the same as the id of the api key passed in and the enterprise id is the same as the enterprise id of the api key passed in
	var dbApiKey ApiKey
	dbOrm.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).First(&dbApiKey)

	dbApiKey.Name = a.Name
	dbApiKey.PermissionsJson = a.PermissionsJson

	// update the api key in the database
	result := dbOrm.Save(&dbApiKey)
	if result.Error != nil {
		log("DB", result.Error.Error())
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

	// delete a single api key from the database using dbOrm where the id is the same as passed in and the enterprise id is the same as the enterprise id of the api key passed in
	result := dbOrm.Delete(ApiKey{}, "id = ? AND enterprise = ?", a.Id, a.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (a *ApiKey) offApiKey() bool {
	if a.Id <= 0 {
		return false
	}

	// get a single api key from the database using dbOrm where the id is the same as the id of the api key passed in and the enterprise id is the same as the enterprise id of the api key passed in
	var dbApiKey ApiKey
	result := dbOrm.Where("id = ? AND enterprise = ?", a.Id, a.EnterpriseId).First(&dbApiKey)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	dbApiKey.Off = !dbApiKey.Off

	// update the api key in the database
	result = dbOrm.Save(&dbApiKey)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
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
	if settings.Server.MaxRequestsPerMinuteEnterprise > 0 && requestsMade >= settings.Server.MaxRequestsPerMinuteEnterprise {
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
	// get the api keys from the database where the token and auth type are the same as the ones passed in
	var dbApiKey ApiKey
	var rowCount int64
	result := dbOrm.Where("token = ? AND auth = ? AND off = false", token, auth).Limit(1).First(&dbApiKey).Count(&rowCount)
	if rowCount == 0 {
		return false, 0, 0, nil
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0, 0, nil
	}

	p := ApiKeyPermissions{}
	if len(dbApiKey.PermissionsJson) > 0 {
		json.Unmarshal([]byte(dbApiKey.PermissionsJson), &p)
	}

	return true, dbApiKey.UserId, dbApiKey.EnterpriseId, &p
}

// checks if the api key exists.
// returns is there exists and active key with this uuid, and if exists, returns also the userId
// P = Parameter, H = Header, B = Basic Auth
// OK, user id, enterprise id
func checkApiKeyByBasicAuthType(basicAuthUser string, basicAuthPassword string) (bool, int32, int32, *ApiKeyPermissions) {
	// get a single api key from the database where the basic_auth_user and basic_auth_password are the same as the ones passed in and the key is not off
	var dbApiKey ApiKey
	var rowCount int64
	result := dbOrm.Where("basic_auth_user = ? AND basic_auth_password = ? AND off = false", basicAuthUser, basicAuthPassword).Limit(1).First(&dbApiKey).Count(&rowCount)
	if rowCount == 0 {
		return false, 0, 0, nil
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false, 0, 0, nil
	}

	p := ApiKeyPermissions{}
	if len(dbApiKey.PermissionsJson) > 0 {
		json.Unmarshal([]byte(dbApiKey.PermissionsJson), &p)
	}

	return true, dbApiKey.UserId, dbApiKey.EnterpriseId, &p
}

type ApiKeyPermissions struct {
	SaleOrders                            ApiKeyPermission `json:"saleOrders"`
	SaleOrderDetails                      ApiKeyPermission `json:"saleOrderDetails"`
	SaleOrderDetailsDigitalProductData    ApiKeyPermission `json:"saleOrderDetailsDigitalProductData"`
	SaleInvoices                          ApiKeyPermission `json:"saleInvoices"`
	SaleInvoiceDetails                    ApiKeyPermission `json:"saleInvoiceDetails"`
	SaleDeliveryNotes                     ApiKeyPermission `json:"saleDeliveryNotes"`
	PurchaseOrders                        ApiKeyPermission `json:"purchaseOrders"`
	PurchaseOrderDetails                  ApiKeyPermission `json:"purchaseOrderDetails"`
	PurchaseInvoices                      ApiKeyPermission `json:"purchaseInvoices"`
	PurchaseInvoiceDetails                ApiKeyPermission `json:"purchaseInvoiceDetails"`
	PurchaseDeliveryNotes                 ApiKeyPermission `json:"purchaseDeliveryNotes"`
	Customers                             ApiKeyPermission `json:"customers"`
	Suppliers                             ApiKeyPermission `json:"suppliers"`
	Products                              ApiKeyPermission `json:"products"`
	Countries                             ApiKeyPermission `json:"countries"`
	States                                ApiKeyPermission `json:"states"`
	Colors                                ApiKeyPermission `json:"colors"`
	ProductFamilies                       ApiKeyPermission `json:"productFamilies"`
	Addresses                             ApiKeyPermission `json:"addresses"`
	Carriers                              ApiKeyPermission `json:"carriers"`
	BillingSeries                         ApiKeyPermission `json:"billingSeries"`
	Currencies                            ApiKeyPermission `json:"currencies"`
	PaymentMethods                        ApiKeyPermission `json:"paymentMethods"`
	Languages                             ApiKeyPermission `json:"languages"`
	Packages                              ApiKeyPermission `json:"packages"`
	Incoterms                             ApiKeyPermission `json:"incoterms"`
	Warehouses                            ApiKeyPermission `json:"warehouses"`
	WarehouseMovements                    ApiKeyPermission `json:"warehouseMovements"`
	TransferBetweenWarehouses             ApiKeyPermission `json:"transferBetweenWarehouses"`
	TransferBetweenWarehousesDetail       ApiKeyPermission `json:"transferBetweenWarehousesDetail"`
	TransferBetweenWarehousesMinimumStock ApiKeyPermission `json:"transferBetweenWarehousesMinimumStock"`
	ManufacturingOrders                   ApiKeyPermission `json:"manufacturingOrders"`
	ManufacturingOrderTypes               ApiKeyPermission `json:"manufacturingOrderTypes"`
	ComplexManufacturingOrders            ApiKeyPermission `json:"complexManufacturingOrders"`
	ComplexManufacturingOrdersComponents  ApiKeyPermission `json:"complexManufacturingOrdersComponents"`
	ManufacturingOrderTypeComponents      ApiKeyPermission `json:"manufacturingOrderTypeComponents"`
	Shippings                             ApiKeyPermission `json:"shippings"`
	ShippingStatusHistory                 ApiKeyPermission `json:"shippingStatusHistory"`
	Stock                                 ApiKeyPermission `json:"stock"`
	Journal                               ApiKeyPermission `json:"journal"`
	Account                               ApiKeyPermission `json:"account"`
	AccountingMovement                    ApiKeyPermission `json:"accountingMovement"`
	AccountingMovementDetail              ApiKeyPermission `json:"accountingMovementDetail"`
	CollectionOperation                   ApiKeyPermission `json:"collectionOperation"`
	Charges                               ApiKeyPermission `json:"charges"`
	PaymentTransaction                    ApiKeyPermission `json:"paymentTransaction"`
	Payment                               ApiKeyPermission `json:"payment"`
	PostSaleInvoice                       ApiKeyPermission `json:"postSaleInvoice"`
	PostPurchaseInvoice                   ApiKeyPermission `json:"postPurchaseInvoice"`
}

type ApiKeyPermission struct {
	Get    bool `json:"get"`
	Post   bool `json:"post"`
	Put    bool `json:"put"`
	Delete bool `json:"delete"`
}
