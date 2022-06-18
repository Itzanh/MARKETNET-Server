package main

import (
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Shipping struct {
	Id                int64             `json:"id" gorm:"index:shipping_id_enterprise,unique:true,priority:1"`
	OrderId           int64             `json:"orderId" gorm:"column:order;not null:true"`
	Order             SaleOrder         `json:"order" gorm:"foreignKey:OrderId,EnterpriseId;references:Id,EnterpriseId"`
	DeliveryNoteId    int64             `json:"deliveryNoteId" gorm:"column:delivery_note;not null:true"`
	DeliveryNote      SalesDeliveryNote `json:"deliveryNote" gorm:"foreignKey:DeliveryNoteId,EnterpriseId;references:Id,EnterpriseId"`
	DeliveryAddressId int32             `json:"deliveryAddressId" gorm:"column:delivery_address;not null:true"`
	DeliveryAddress   Address           `json:"deliveryAddress" gorm:"foreignKey:DeliveryAddressId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated       time.Time         `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone"`
	DateSent          *time.Time        `json:"dateSent" gorm:"column:date_sent;type:timestamp(3) with time zone"`
	Sent              bool              `json:"sent" gorm:"column:sent;not null:true;index:shipping_sent_collected,priority:1;index:shipping_sent_collected_delivered,priority:1"`
	Collected         bool              `json:"collected" gorm:"column:collected;not null:true;index:shipping_sent_collected,priority:2;index:shipping_sent_collected_delivered,priority:2"`
	National          bool              `json:"national" gorm:"column:national;not null:true"`
	ShippingNumber    string            `json:"shippingNumber" gorm:"column:shipping_number;not null:true;type:character varying(50)"`
	TrackingNumber    string            `json:"trackingNumber" gorm:"column:tracking_number;not null:true;type:character varying(50)"`
	CarrierId         int32             `json:"carrierId" gorm:"column:carrier;not null:true"`
	Carrier           Carrier           `json:"carrier" gorm:"foreignKey:CarrierId,EnterpriseId;references:Id,EnterpriseId"`
	Weight            float64           `json:"weight" gorm:"column:weight;not null:true;type:numeric(14,6)"`
	PackagesNumber    int16             `json:"packagesNumber" gorm:"column:packages_number;not null:true"`
	IncotermId        *int32            `json:"incotermId" gorm:"column:incoterm"`
	Incoterm          *Incoterm         `json:"incoterm" gorm:"foreignKey:IncotermId,EnterpriseId;references:Id,EnterpriseId"`
	CarrierNotes      string            `json:"carrierNotes" gorm:"column:carrier_notes;not null:true;type:character varying(250)"`
	Description       string            `json:"description" gorm:"column:description;not null:true;type:text"`
	EnterpriseId      int32             `json:"-" gorm:"column:enterprise;not null:true;index:shipping_id_enterprise,unique:true,priority:2"`
	Enterprise        Settings          `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Delivered         bool              `json:"delivered" gorm:"column:delivered;not null:true;index:shipping_sent_collected_delivered,priority:3"`
}

func (s *Shipping) TableName() string {
	return "shipping"
}

func getShippings(enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	result := dbOrm.Model(&Shipping{}).Where("enterprise = ?", enterpriseId).Order("id DESC").Preload(clause.Associations).Preload("Order.Customer").Find(&shippings)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return shippings
}

type SearchShippings struct {
	Search    string     `json:"search"`
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
	Status    string     `json:"status"` // "" = All, "S" = Shipped, "N" = Not shipped
}

func (s *SearchShippings) searchShippings(enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)

	orderNumber, err := strconv.Atoi(s.Search)
	cursor := dbOrm.Model(&Shipping{}).Where("shipping.enterprise = ?", enterpriseId).Joins(`INNER JOIN sales_order ON sales_order.id=shipping."order"`)
	if err == nil {
		cursor = cursor.Where("sales_order.order_number = ?", orderNumber)
	} else {
		cursor = cursor.Where("customer.name ILIKE ?", "%"+s.Search+"%").Joins("INNER JOIN customer ON customer.id=sales_order.customer")
	}
	if s.DateStart != nil {
		cursor = cursor.Where("shipping.date_created >= ?", s.DateStart)
	}
	if s.DateEnd != nil {
		cursor = cursor.Where("shipping.date_created <= ?", s.DateEnd)
	}
	if s.Status == "S" {
		cursor = cursor.Where("shipping.sent=true")
	} else if s.Status == "N" {
		cursor = cursor.Where("shipping.sent=false")
	}
	result := cursor.Order("shipping.id DESC").Preload(clause.Associations).Preload("Order.Customer").Find(&shippings)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return shippings
}

func getShippingRow(shippingId int64) Shipping {
	s := Shipping{}
	result := dbOrm.Model(&Shipping{}).Where("id = ?", shippingId).Preload(clause.Associations).First(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return s
}

func getShippingsPendingCollected(enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	result := dbOrm.Model(&Shipping{}).Where("sent = true AND collected = false AND enterprise = ?", enterpriseId).Order("id DESC").Preload(clause.Associations).Preload("Order.Customer").Find(&shippings)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return shippings
}

func (s *Shipping) isValid() bool {
	return !(s.OrderId <= 0 || s.DeliveryAddressId <= 0 || s.CarrierId <= 0 || len(s.CarrierNotes) > 250 || len(s.Description) > 3000)
}

func (s *Shipping) BeforeCreate(tx *gorm.DB) (err error) {
	var shipping Shipping
	tx.Model(&Shipping{}).Last(&shipping)
	s.Id = shipping.Id + 1
	return nil
}

func (s *Shipping) insertShipping(userId int32, trans *gorm.DB) (bool, int64) {
	if !s.isValid() {
		return false, 0
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false, 0
		}
		///
	}

	s.DateCreated = time.Now()
	s.DateSent = nil
	s.Sent = false
	s.Collected = false
	s.ShippingNumber = ""
	s.TrackingNumber = ""
	s.Weight = 0
	s.PackagesNumber = 0
	s.Delivered = false

	result := trans.Create(&s)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(s.EnterpriseId, "shipping", int(s.Id), userId, "I")

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false, 0
		}
		///
	}

	return true, s.Id
}

func (s *Shipping) updateShipping(userId int32) bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	inDatabaseShipping := getShippingRow(s.Id)
	if inDatabaseShipping.Id <= 0 || inDatabaseShipping.EnterpriseId != s.EnterpriseId {
		return false
	}

	if inDatabaseShipping.Carrier.Webservice != "_" && (inDatabaseShipping.ShippingNumber != s.ShippingNumber || inDatabaseShipping.TrackingNumber != s.TrackingNumber) {
		return false
	}

	inDatabaseShipping.OrderId = s.OrderId
	inDatabaseShipping.DeliveryNoteId = s.DeliveryNoteId
	inDatabaseShipping.DeliveryAddressId = s.DeliveryAddressId
	inDatabaseShipping.CarrierId = s.CarrierId
	inDatabaseShipping.IncotermId = s.IncotermId
	inDatabaseShipping.CarrierNotes = s.CarrierNotes
	inDatabaseShipping.Description = s.Description
	inDatabaseShipping.ShippingNumber = s.ShippingNumber
	inDatabaseShipping.TrackingNumber = s.TrackingNumber

	result := dbOrm.Model(&Shipping{}).Where("id = ?", s.Id).Updates(map[string]interface{}{
		"order":            inDatabaseShipping.OrderId,
		"delivery_note":    inDatabaseShipping.DeliveryNoteId,
		"delivery_address": inDatabaseShipping.DeliveryAddressId,
		"carrier":          inDatabaseShipping.CarrierId,
		"incoterm":         inDatabaseShipping.IncotermId,
		"carrier_notes":    inDatabaseShipping.CarrierNotes,
		"description":      inDatabaseShipping.Description,
		"shipping_number":  inDatabaseShipping.ShippingNumber,
		"tracking_number":  inDatabaseShipping.TrackingNumber,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(s.EnterpriseId, "shipping", int(s.Id), userId, "U")

	if inDatabaseShipping.ShippingNumber != s.ShippingNumber || inDatabaseShipping.TrackingNumber != s.TrackingNumber {
		go ecommerceControllerUpdateTrackingNumber(inDatabaseShipping.OrderId, s.TrackingNumber, s.EnterpriseId)
	}

	return true
}

func (s *Shipping) deleteShipping(userId int32) bool {
	if s.Id <= 0 {
		return false
	}

	inDatabaseShipping := getShippingRow(s.Id)
	if inDatabaseShipping.Id <= 0 || inDatabaseShipping.EnterpriseId != s.EnterpriseId {
		return false
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	packaging := getPackagingByShipping(s.Id, s.EnterpriseId)
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detailsPackaged[j].OrderDetailId).Update("status", "E")
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(s.EnterpriseId, "sales_order_detail", int(detailsPackaged[j].OrderDetailId), userId, "U")

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetailId)
			ok := setSalesOrderState(s.EnterpriseId, saleOrderDetail.OrderId, userId, *trans)
			if !ok {
				trans.Rollback()
				return false
			}
		}

		result := trans.Model(&Packaging{}).Where("id = ?", packaging[i].Id).Update("shipping", nil)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(s.EnterpriseId, "packaging", int(packaging[i].Id), userId, "U")
	}

	insertTransactionalLog(s.EnterpriseId, "shipping", int(s.Id), userId, "D")

	result := trans.Delete(&Shipping{}, "id = ?", s.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	///
	result = trans.Commit()
	return result.Error == nil
	///
}

type OkAndErrorGenerateShipping struct {
	OkAndErrorCodeReturn
	Shipping Shipping `json:"shipping"`
}

// ERROR CODES:
// 1. No carrier selected in the order
// 2. A detail has not been completely packaged
// 3. Can't generate delivery note
func generateShippingFromSaleOrder(orderId int64, enterpriseId int32, userId int32) OkAndErrorGenerateShipping {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.EnterpriseId != enterpriseId {
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
	}
	packaging := getPackaging(orderId, enterpriseId)
	if saleOrder.Id <= 0 || len(packaging) == 0 {
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
	}
	if saleOrder.CarrierId == nil || *saleOrder.CarrierId <= 0 {
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}}
	}

	details := getSalesOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].QuantityPendingPackaging > 0 {
			return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false, ErrorCode: 2, ExtraData: []string{details[i].Product.Name}}}
		}
	}

	s := Shipping{}
	s.OrderId = saleOrder.Id
	s.DeliveryAddressId = saleOrder.ShippingAddressId
	s.CarrierId = *saleOrder.CarrierId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
	}
	///

	saleDeliveryNotes := getSalesOrderDeliveryNotes(orderId, enterpriseId)
	if len(saleDeliveryNotes) > 0 {
		s.DeliveryNoteId = saleDeliveryNotes[0].Id
	} else {
		ok, noteId := deliveryNoteAllSaleOrder(orderId, enterpriseId, userId, trans)
		if !ok.Ok || noteId <= 0 {
			trans.Rollback()
			return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}}
		}
		s.DeliveryNoteId = noteId
	}

	s.EnterpriseId = enterpriseId
	ok, shippingId := s.insertShipping(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
	}
	for i := 0; i < len(packaging); i++ {
		ok := associatePackagingToShipping(packaging[i].Id, shippingId, enterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
		}
	}
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			result := trans.Model(&SalesOrderDetail{}).Where("id = ?", detailsPackaged[j].OrderDetailId).Update("status", "F")
			if result.Error != nil {
				log("DB", result.Error.Error())
				trans.Rollback()
				return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
			}

			insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailsPackaged[j].OrderDetailId), userId, "U")

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetailId)
			ok := setSalesOrderState(enterpriseId, saleOrderDetail.OrderId, userId, *trans)
			if !ok {
				trans.Rollback()
				return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
			}
		}
	}

	///
	result := trans.Commit()
	if result.Error != nil {
		return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: false}}
	}

	return OkAndErrorGenerateShipping{OkAndErrorCodeReturn: OkAndErrorCodeReturn{Ok: true}, Shipping: getShippingRow(shippingId)}
	///
}

// THIS FUNCION DOES NOT OPEN A TRANSACTION
func associatePackagingToShipping(packagingId int64, shippingId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	result := trans.Model(&Packaging{}).Where("id = ?", packagingId).Update("shipping", shippingId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	var weight float64
	result = trans.Model(&Packaging{}).Where("id = ?", packagingId).Pluck("weight", &weight)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	var shipping Shipping
	result = trans.Model(&Shipping{}).Where("id = ?", shippingId).First(&shipping)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	shipping.Weight += weight
	shipping.PackagesNumber += 1

	result = trans.Updates(&shipping)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "shipping", int(shippingId), userId, "U")

	return true
}

type ToggleShippingSent struct {
	Ok             bool    `json:"ok"`
	ErrorMessage   *string `json:"errorMessage"`
	ShippingNumber *string `json:"shippingNumber"`
	TrackingNumber *string `json:"trackingNumber"`
}

func toggleShippingSent(shippingId int64, enterpriseId int32, userId int32) ToggleShippingSent {
	s := getShippingRow(shippingId)
	if s.EnterpriseId != enterpriseId {
		return ToggleShippingSent{Ok: false}
	}
	// it is not allowed to manually set as "sent" if the carrier has a webservice set.
	// it is not allowed to modify the "sent" field if the shipping was colletected by the carrier.
	if s.Id <= 0 || s.Collected {
		return ToggleShippingSent{Ok: false}
	}
	if s.Carrier.Webservice != "_" {
		ok, errorMessage := s.sendShipping(enterpriseId)
		if ok {
			s := getShippingRow(shippingId)
			go ecommerceControllerUpdateTrackingNumber(s.OrderId, s.TrackingNumber, enterpriseId)
			return ToggleShippingSent{Ok: ok, TrackingNumber: &s.TrackingNumber, ShippingNumber: &s.ShippingNumber}
		} else {
			return ToggleShippingSent{Ok: ok, ErrorMessage: errorMessage}
		}
	}

	s.Sent = !s.Sent
	if s.Sent {
		now := time.Now()
		s.DateSent = &now
	} else {
		s.DateSent = nil
	}

	result := dbOrm.Model(&Shipping{}).Where("id = ?", shippingId).Updates(map[string]interface{}{
		"sent":      s.Sent,
		"date_sent": s.DateSent,
	})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ToggleShippingSent{Ok: false}
	}

	insertTransactionalLog(enterpriseId, "shipping", int(shippingId), userId, "U")

	return ToggleShippingSent{Ok: true}
}

func (s *Shipping) sendShipping(enterpriseId int32) (bool, *string) {
	switch s.Carrier.Webservice {
	case "S":
		return s.sendShippingSendCloud(enterpriseId)
	default:
		return false, nil
	}
}

func (s *Shipping) sendShippingSendCloud(enterpriseId int32) (bool, *string) {
	ok, parcel := s.generateSendCloudParcel(enterpriseId)
	if !ok {
		return false, nil
	}
	return parcel.send(s)
}

func setShippingCollected(shippings []int64, enterpriseId int32, userId int32) bool {
	if len(shippings) == 0 {
		return false
	}
	for i := 0; i < len(shippings); i++ {
		if shippings[i] <= 0 {
			return false
		}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return false
	}
	///

	for i := 0; i < len(shippings); i++ {
		s := getShippingRow(shippings[i])
		if s.Id <= 0 || !s.Sent || s.EnterpriseId != enterpriseId {
			trans.Rollback()
			return false
		}

		result := trans.Model(&Shipping{}).Where("id = ?", shippings[i]).Update("collected", true)
		if result.Error != nil {
			log("DB", result.Error.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(enterpriseId, "shipping", int(shippings[i]), userId, "U")

		p := getPackagingByShipping(shippings[i], enterpriseId)
		for j := 0; j < len(p); j++ {
			for k := 0; k < len(p[j].DetailsPackaged); k++ {
				result = trans.Model(&SalesOrderDetail{}).Where("id = ?", p[j].DetailsPackaged[k].OrderDetailId).Where("status", "G")
				if result.Error != nil {
					log("DB", result.Error.Error())
					trans.Rollback()
					return false
				}
				setSalesOrderState(enterpriseId, p[j].DetailsPackaged[k].OrderDetail.OrderId, userId, *trans)

				insertTransactionalLog(enterpriseId, "sales_order_detail", int(p[j].DetailsPackaged[k].OrderDetailId), userId, "U")
			}
		}

	}

	///
	result := trans.Commit()
	return result.Error == nil
	///
}
