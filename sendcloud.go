package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT = 75
const SENDCLOUD_EMAIL_ALLOWED_CHARACTER_SET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.!#$%&'*+-/=?^_`{|}~"
const SENDCLOUD_COMMERCIAL_GOODS = int8(2)

type Parcel struct {
	Name                    string         `json:"name"`
	CompanyName             string         `json:"company_name"`
	Address                 string         `json:"address"`
	Address_2               *string        `json:"address_2"`
	City                    string         `json:"city"`
	PostalCode              string         `json:"postal_code"`
	Country                 string         `json:"country"` // country ISO-2 code
	CountryState            *string        `json:"country_state"`
	Telephone               *int64         `json:"telephone"`
	Email                   string         `json:"email"`
	SenderAddress           int64          `json:"sender_address"`
	CustomsInvoiceNr        string         `json:"customs_invoice_nr"`
	CustomsShipmentType     *int8          `json:"customs_shipment_type"`
	ExternalReference       *string        `json:"external_reference"`
	Quantity                int8           `json:"quantity"`
	OrderNumber             *string        `json:"order_number"`
	ParcelItems             []ParcelItem   `json:"parcel_items"`
	Weight                  *float32       `json:"weight"`
	TotalOrderValue         *float32       `json:"total_order_value"`
	TotalOrderValueCurrency *string        `json:"total_order_value_currency"` // currency ISO-3
	Length                  *float32       `json:"length"`
	Width                   *float32       `json:"width"`
	Height                  *float32       `json:"height"`
	Shipment                ParcelShipment `json:"shipment"`
	RequestLabel            bool           `json:"request_label"`
}

type ParcelItem struct {
	Description   string  `json:"description"`
	Quantity      int8    `json:"quantity"`
	Weight        float32 `json:"weight"`
	Value         float32 `json:"value"`
	HSCode        string  `json:"hs_code"`
	OriginCountry *string `json:"origin_country"`
}

type ParcelShipment struct {
	Id   int32   `json:"id"`
	Name *string `json:"name"`
}

func (s *Shipping) generateSendCloudParcel() (bool, *Parcel) {
	p := Parcel{}
	p.Quantity = 1
	p.RequestLabel = true

	carrier := getCarierRow(s.Carrier)
	if carrier.Id <= 0 {
		return false, nil
	}

	// get the order
	o := getSalesOrderRow(s.Order)
	if o.Id <= 0 {
		return false, nil
	}

	// customer name
	c := getCustomerRow(o.Customer)
	if c.Id <= 0 {
		return false, nil
	}
	p.Name = c.Name

	// company name
	settings := getSettingsRecord()
	p.CompanyName = settings.EnterpriseName

	// address
	a := getAddressRow(o.ShippingAddress)
	p.Address = strings.TrimSpace(a.Address)
	if len(p.Address) > SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT {
		p.Address = p.Address[0:SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT]
	}
	if len(a.Address2) > 0 {
		address2 := strings.TrimSpace(a.Address2)
		if len(address2) > SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT {
			address2 = address2[0:SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT]
		}
		p.Address_2 = &address2
	}
	p.City = a.City
	p.PostalCode = a.ZipCode
	p.Country = getCountryRow(a.Country).Iso2 // country must have a ISO2 code!
	if a.State != nil {
		stateIsoCode := getStateRow(*a.State).IsoCode
		p.CountryState = &stateIsoCode
	}

	// telephone+email
	telephone, err := strconv.Atoi(c.Phone)
	if err == nil {
		tlf64 := int64(telephone)
		p.Telephone = &tlf64
	}

	p.Email = c.Email
	for i := 0; i < len(p.Email); i++ {
		if !strings.Contains(SENDCLOUD_EMAIL_ALLOWED_CHARACTER_SET, string(p.Email[i])) {
			p.Email = strings.ReplaceAll(p.Email, string(p.Email[i]), "")
		}
	}

	// external reference / order number
	p.ExternalReference = &o.OrderName
	p.OrderNumber = &o.OrderName

	// parcel items
	p.ParcelItems = make([]ParcelItem, 0)
	packaging := getPackagingByShipping(s.Id)
	for i := 0; i < len(packaging); i++ {
		pi := ParcelItem{}
		pi.Description = packaging[i].PackageName
		pi.Weight = packaging[i].Weight

		details := getSalesOrderDetailPackaged(packaging[i].Id)
		for j := 0; j < len(details); j++ {
			pi.Quantity += int8(details[j].Quantity)
			pi.Value += getSalesOrderDetailRow(details[j].OrderDetail).Price * float32(details[j].Quantity)
		}

		p.ParcelItems = append(p.ParcelItems, pi)
	}

	// weight
	p.Weight = &s.Weight

	// shipment
	p.Shipment = ParcelShipment{Id: carrier.SendcloudShippingMethod}
	// sender address
	p.SenderAddress = carrier.SendcloudSenderAddress

	// commercial invoice
	invoices := getSalesOrderInvoices(s.Order)
	if len(invoices) > 0 {
		p.CustomsInvoiceNr = invoices[0].InvoiceName
		commercialGoods := SENDCLOUD_COMMERCIAL_GOODS
		p.CustomsShipmentType = &commercialGoods // Commercial Goods
	}

	return true, &p
}

func (p *Parcel) send(s *Shipping) (bool, *string) {
	// get the carrier
	c := getCarierRow(s.Carrier)
	if c.Id <= 0 {
		return false, nil
	}

	// make the request
	parcelObject := make(map[string]*Parcel)
	parcelObject["parcel"] = p
	jsonRequest, _ := json.Marshal(parcelObject)
	req, err := http.NewRequest("POST", c.SendcloudUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return false, nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SendcloudKey, c.SendcloudSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	// get the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil
	}
	var response ParcelResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, nil
	}
	if response.Parcel == nil {
		if response.Error == nil {
			return false, nil
		}
		parcelError := *response.Error
		return false, &parcelError.Message
	}
	parcelResponse := *response.Parcel

	if parcelResponse.Id <= 0 {
		return false, nil
	}

	// update the shipping
	s.TrackingNumber = parcelResponse.TrackingNumber
	s.ShippingNumber = strconv.Itoa(int(parcelResponse.Id))

	sqlStatement := `UPDATE shipping SET sent = NOT sent, date_sent = CASE sent WHEN false THEN CURRENT_TIMESTAMP(3) ELSE NULL END, tracking_number=$2, shipping_number=$3 WHERE id = $1`
	_, err = db.Exec(sqlStatement, s.Id, s.TrackingNumber, s.ShippingNumber)

	if err != nil {
		log("DB", err.Error())
		return false, nil
	}

	// save the label
	return parcelResponse.saveLabel(c, s.Id), nil
}

func (p *ParcelResponse) saveLabel(c Carrier, shippingId int32) bool {
	req, err := http.NewRequest("GET", p.Label.LabelPrinter, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SendcloudKey, c.SendcloudSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	// get the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	t := ShippingTag{}
	t.Shipping = shippingId
	t.Label = body
	return t.insertShippingTag()
}

type ParcelResponseBody struct {
	Parcel *ParcelResponse `json:"parcel"`
	Error  *ParcelError    `json:"error"`
}

// Does not contain all the attriutes from the response, only the ones that we want to extract.
type ParcelResponse struct {
	Id             int64               `json:"id"`
	TrackingNumber string              `json:"tracking_number"`
	Label          ParcelResponseLabel `json:"label"`
	ColliUuid      string              `json:"colli_uuid"`
	TrackingUrl    string              `json:"tracking_url"`
}

type ParcelResponseLabel struct {
	LabelPrinter string `json:"label_printer"`
}

type ParcelError struct {
	Message string `json:"message"`
}
