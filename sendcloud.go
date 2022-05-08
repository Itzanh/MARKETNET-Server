package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT = 75
const SENDCLOUD_EMAIL_ALLOWED_CHARACTER_SET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.!#$%&'*+-/=?^_`{|}~@"
const SENDCLOUD_COMMERCIAL_GOODS = int8(2)
const SENDCLOUD_MIN_WEIGHT_PARCEL_ITEMS = 0.00099

type Parcel struct {
	Name                    string         `json:"name"`
	CompanyName             string         `json:"company_name"`
	Address                 string         `json:"address"`
	Address_2               string         `json:"address_2"`
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
	Weight                  *float64       `json:"weight"`
	TotalOrderValue         *float64       `json:"total_order_value"`
	TotalOrderValueCurrency *string        `json:"total_order_value_currency"` // currency ISO-3
	Length                  *float64       `json:"length"`
	Width                   *float64       `json:"width"`
	Height                  *float64       `json:"height"`
	Shipment                ParcelShipment `json:"shipment"`
	RequestLabel            bool           `json:"request_label"`
}

type ParcelItem struct {
	Description   string  `json:"description"`
	Quantity      int32   `json:"quantity"`
	Weight        float64 `json:"weight"`
	Value         float64 `json:"value"`
	HSCode        string  `json:"hs_code"`
	OriginCountry *string `json:"origin_country"`
	SKU           *string `json:"sku"`
	ProductId     string  `json:"product_id"`
}

type ParcelShipment struct {
	Id   int32   `json:"id"`
	Name *string `json:"name"`
}

func (s *Shipping) generateSendCloudParcel(enterpriseId int32) (bool, *Parcel) {
	p := Parcel{}
	p.Quantity = 1
	p.RequestLabel = true

	// get the order
	o := getSalesOrderRow(s.OrderId)
	if o.Id <= 0 {
		return false, nil
	}

	carrier := getCarierRow(s.CarrierId)
	if carrier.Id <= 0 {
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. The order doesn't have a shipping. "+o.OrderName+"</p>", enterpriseId)
		}
		return false, nil
	}

	// customer name
	c := getCustomerRow(o.CustomerId)
	if c.Id <= 0 {
		return false, nil
	}
	p.Name = c.Name

	// company name
	settings := getSettingsRecordById(enterpriseId)
	p.CompanyName = settings.EnterpriseName

	// address
	a := getAddressRow(o.ShippingAddressId)
	p.Address = strings.TrimSpace(a.Address)
	if len(p.Address) > SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT {
		p.Address = p.Address[0:SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT]
	}
	if len(a.Address2) > 0 {
		address2 := strings.TrimSpace(a.Address2)
		if len(address2) > SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT {
			address2 = address2[0:SENDCLOUD_MAX_ADDRESS_CHARACTER_LIMIT]
		}
		p.Address_2 = address2
	} else {
		p.Address_2 = ""
	}
	p.City = a.City
	p.PostalCode = a.ZipCode
	p.Country = getCountryRow(a.CountryId, enterpriseId).Iso2 // country must have a ISO2 code!
	if a.StateId != nil {
		stateIsoCode := getStateRow(*a.StateId).IsoCode
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
	var weight float64 = 0
	p.ParcelItems = make([]ParcelItem, 0)
	details := getSalesOrderDetail(s.OrderId, enterpriseId)
	for i := 0; i < len(details); i++ {
		product := getProductRow(details[i].ProductId)

		pi := ParcelItem{}
		pi.Description = details[i].Product.Name
		pi.Quantity = details[i].Quantity
		pi.Weight = toFixed(math.Max(product.Weight*float64(details[i].Quantity), SENDCLOUD_MIN_WEIGHT_PARCEL_ITEMS), 3)
		weight += pi.Weight
		pi.Value = toFixed(details[i].TotalAmount, 2)
		if product.HSCodeId != nil {
			pi.HSCode = *product.HSCodeId
		}
		pi.OriginCountry = &product.OriginCountry
		pi.SKU = &product.BarCode
		pi.ProductId = strconv.Itoa(int(details[i].ProductId))

		p.ParcelItems = append(p.ParcelItems, pi)
	}

	// weight
	weight = toFixed(weight, 3)
	p.Weight = &weight

	// shipment
	p.Shipment = ParcelShipment{Id: carrier.SendcloudShippingMethod}
	// sender address
	p.SenderAddress = carrier.SendcloudSenderAddress

	// commercial invoice
	invoices := getSalesOrderInvoices(s.OrderId, enterpriseId)
	if len(invoices) > 0 {
		p.CustomsInvoiceNr = invoices[0].InvoiceName
		commercialGoods := SENDCLOUD_COMMERCIAL_GOODS
		p.CustomsShipmentType = &commercialGoods // Commercial Goods
	}

	return true, &p
}

func (p *Parcel) send(s *Shipping) (bool, *string) {
	// get the carrier
	c := getCarierRow(s.CarrierId)
	if c.Id <= 0 {
		return false, nil
	}

	// make the request
	parcelObject := make(map[string]*Parcel)
	parcelObject["parcel"] = p
	jsonRequest, _ := json.Marshal(parcelObject)
	req, err := http.NewRequest("POST", c.SendcloudUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		s := getSettingsRecordById(s.EnterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", s.Id)
		}
		return false, nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SendcloudKey, c.SendcloudSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s := getSettingsRecordById(s.EnterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", s.Id)
		}
		return false, nil
	}
	// get the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s := getSettingsRecordById(s.EnterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", s.Id)
		}
		return false, nil
	}
	var response ParcelResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		s := getSettingsRecordById(s.EnterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p><p>"+string(body)+"</p>", s.Id)
		}
		return false, nil
	}
	if response.Parcel == nil {
		if response.Error == nil {
			return false, nil
		}
		parcelError := *response.Error
		log("SendCloud", string(jsonRequest)+parcelError.Message)
		s := getSettingsRecordById(s.EnterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order. There was an error connecting with SendCloud.</p><p>"+string(jsonRequest)+"</p><p>"+parcelError.Message+"</p>", s.Id)
		}
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
	return parcelResponse.saveLabel(c, s.Id, s.EnterpriseId), nil
}

func (p *ParcelResponse) saveLabel(c Carrier, shippingId int64, enterpriseId int32) bool {
	req, err := http.NewRequest("GET", p.Label.LabelPrinter, nil)
	if err != nil {
		log("SendCloud", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order and saving the label. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.SendcloudKey, c.SendcloudSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log("SendCloud", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order and saving the label. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", enterpriseId)
		}
		return false
	}
	// get the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("SendCloud", err.Error())
		s := getSettingsRecordById(enterpriseId)
		if len(s.EmailSendErrorSendCloud) > 0 {
			sendEmail(s.EmailSendErrorSendCloud, s.EmailSendErrorSendCloud, "SendCloud shipping error",
				"<p>Error when trying to ship the order and saving the label. There was an error connecting with SendCloud.</p><p>"+err.Error()+"</p>", enterpriseId)
		}
		return false
	}

	t := ShippingTag{}
	t.ShippingId = shippingId
	t.Label = body
	t.EnterpriseId = enterpriseId
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

func getShippingTrackingSendCloud(enterpriseId int32) {
	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return
	}
	///

	// webservice SendCloud, shipped, collected, but not delivered by the carrier yet
	sqlStatement := `SELECT id,shipping_number,carrier FROM shipping WHERE enterprise=$1 AND (SELECT webservice FROM carrier WHERE carrier.id=shipping.carrier) = 'S' AND sent = true AND collected = true AND delivered = false`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return
	}
	defer rows.Close()

	var shippingId int64
	var shippingNumber string
	var carrier int32
	for rows.Next() {
		rows.Scan(&shippingId, &shippingNumber, &carrier)
		c := getCarierRow(carrier)

		req, err := http.NewRequest("GET", c.SendcloudUrl+"/"+shippingNumber, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(c.SendcloudKey, c.SendcloudSecret)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		// get the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		parcel := ParcelGetContainer{}
		err = json.Unmarshal(body, &parcel)
		if err != nil {
			continue
		}

		if parcel.Parcel == nil || parcel.Parcel.Status == nil {
			continue
		}

		var delivered bool = (parcel.Parcel.Status.Id == 11 && parcel.Parcel.Status.Message == "Delivered")

		sqlStatement := `SELECT message FROM public.shipping_status_history WHERE shipping = $1 ORDER BY date_created DESC LIMIT 1`
		row := db.QueryRow(sqlStatement, shippingId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			continue
		}

		var lastMessageInDb string
		row.Scan(&lastMessageInDb)

		if parcel.Parcel.Status.Message != lastMessageInDb {
			sqlStatement := `INSERT INTO public.shipping_status_history(shipping, status_id, message, delivered) VALUES ($1, $2, $3, $4)`
			_, err := db.Exec(sqlStatement, shippingId, parcel.Parcel.Status.Id, parcel.Parcel.Status.Message, delivered)
			if err != nil {
				log("DB", err.Error())
				continue
			}
		}

		if delivered {
			sqlStatement := `UPDATE public.shipping SET delivered=true WHERE id=$1`
			_, err := db.Exec(sqlStatement, shippingId)
			if err != nil {
				log("DB", err.Error())
				continue
			}

			insertTransactionalLog(enterpriseId, "shipping", int(shippingId), 0, "U")
		}
	}

	///
	trans.Commit()
	///
}

type ParcelGetContainer struct {
	Parcel *ParcelGetParcel `json:"parcel"`
}

type ParcelGetParcel struct {
	Status *ParcelGetParcelStatus `json:"status"`
}

type ParcelGetParcelStatus struct {
	Id      int16  `json:"id"`
	Message string `json:"message"`
}
