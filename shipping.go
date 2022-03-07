package main

import (
	"database/sql"
	"strconv"
	"time"
)

type Shipping struct {
	Id                int64      `json:"id"`
	Order             int64      `json:"order"`
	DeliveryNote      int64      `json:"deliveryNote"`
	DeliveryAddress   int32      `json:"deliveryAddress"`
	DateCreated       time.Time  `json:"dateCreated"`
	DateSent          *time.Time `json:"dateSent"`
	Sent              bool       `json:"sent"`
	Collected         bool       `json:"collected"`
	National          bool       `json:"national"`
	ShippingNumber    string     `json:"shippingNumber"`
	TrackingNumber    string     `json:"trackingNumber"`
	Carrier           int32      `json:"carrier"`
	Weight            float64    `json:"weight"`
	PackagesNumber    int16      `json:"packagesNumber"`
	CustomerName      string     `json:"customerName"`
	SaleOrderName     string     `json:"saleOrderName"`
	CarrierName       string     `json:"carrierName"`
	Incoterm          *int32     `json:"incoterm"`
	CarrierNotes      string     `json:"carrierNotes"`
	Description       string     `json:"description"`
	CarrierWebService string     `json:"carrierWebService"`
	Delivered         bool       `json:"delivered"`
	enterprise        int32
}

func getShippings(enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE enterprise=$1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return shippings
	}
	defer rows.Close()

	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.enterprise, &s.Delivered, &s.CustomerName, &s.SaleOrderName, &s.CarrierName, &s.CarrierWebService)
		shippings = append(shippings, s)
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
	var sqlStatement string

	var interfaces []interface{} = make([]interface{}, 0)
	orderNumber, err := strconv.Atoi(s.Search)

	if err == nil {
		sqlStatement = `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM shipping INNER JOIN sales_order ON sales_order.id=shipping."order" WHERE sales_order.order_number=$1 AND shipping.enterorise=$2`
		interfaces = append(interfaces, orderNumber)
	} else {
		sqlStatement = `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM shipping INNER JOIN sales_order ON shipping."order"=sales_order.id INNER JOIN customer ON customer.id=sales_order.customer WHERE (customer.name ILIKE $1) AND (shipping.enterprise=$2)`
		interfaces = append(interfaces, "%"+s.Search+"%")
	}
	interfaces = append(interfaces, enterpriseId)
	if s.DateStart != nil {
		sqlStatement += ` AND date_created>=$` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.DateStart)
	}
	if s.DateEnd != nil {
		sqlStatement += ` AND date_created<=$` + strconv.Itoa(len(interfaces)+1)
		interfaces = append(interfaces, s.DateEnd)
	}
	if s.Status == "S" {
		sqlStatement += ` AND sent=true`
	} else if s.Status == "N" {
		sqlStatement += ` AND sent=false`
	}
	sqlStatement += ` ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return shippings
	}
	defer rows.Close()

	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.enterprise, &s.Delivered, &s.CustomerName, &s.SaleOrderName, &s.CarrierName, &s.CarrierWebService)
		shippings = append(shippings, s)
	}

	return shippings
}

func getShippingRow(shippingId int64) Shipping {
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE id=$1`
	row := db.QueryRow(sqlStatement, shippingId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Shipping{}
	}

	s := Shipping{}
	row.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.enterprise, &s.Delivered, &s.CustomerName, &s.SaleOrderName, &s.CarrierName, &s.CarrierWebService)

	return s
}

func getShippingsPendingCollected(enterpriseId int32) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier),(SELECT webservice FROM carrier WHERE id=shipping.carrier) FROM public.shipping WHERE sent=true AND collected=false AND enterprise=$1 ORDER BY id DESC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return shippings
	}
	defer rows.Close()

	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.enterprise, &s.Delivered, &s.CustomerName, &s.SaleOrderName, &s.CarrierName, &s.CarrierWebService)
		shippings = append(shippings, s)
	}

	return shippings
}

func (s *Shipping) isValid() bool {
	return !(s.Order <= 0 || s.DeliveryAddress <= 0 || s.Carrier <= 0 || len(s.Description) > 3000)
}

func (s *Shipping) insertShipping(userId int32, trans *sql.Tx) (bool, int64) {
	if !s.isValid() {
		return false, 0
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		var transErr error
		trans, transErr = db.Begin()
		if transErr != nil {
			return false, 0
		}
		///
	}

	sqlStatement := `INSERT INTO public.shipping("order", delivery_note, delivery_address, "national", carrier, incoterm, carrier_notes, description, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	row := trans.QueryRow(sqlStatement, s.Order, s.DeliveryNote, s.DeliveryAddress, s.National, s.Carrier, s.Incoterm, s.CarrierNotes, s.Description, s.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false, 0
	}

	var shippingId int64
	row.Scan(&shippingId)
	s.Id = shippingId

	if shippingId > 0 {
		insertTransactionalLog(s.enterprise, "shipping", int(shippingId), userId, "I")
	}

	if beginTransaction {
		///
		err := trans.Commit()
		if err != nil {
			return false, 0
		}
		///
	}

	return shippingId > 0, shippingId
}

func (s *Shipping) updateShipping(userId int32) bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	inDatabaseShipping := getShippingRow(s.Id)
	if inDatabaseShipping.Id <= 0 || inDatabaseShipping.enterprise != s.enterprise {
		return false
	}

	if inDatabaseShipping.CarrierWebService != "_" && (inDatabaseShipping.ShippingNumber != s.ShippingNumber || inDatabaseShipping.TrackingNumber != s.TrackingNumber) {
		return false
	}

	sqlStatement := `UPDATE public.shipping SET "order"=$2, delivery_note=$3, delivery_address=$4, carrier=$5, incoterm=$6, carrier_notes=$7, description=$8, shipping_number=$9, tracking_number=$10 WHERE id=$1 AND enterprise=$11`
	res, err := db.Exec(sqlStatement, s.Id, s.Order, s.DeliveryNote, s.DeliveryAddress, s.Carrier, s.Incoterm, s.CarrierNotes, s.Description, s.ShippingNumber, s.TrackingNumber, s.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	insertTransactionalLog(s.enterprise, "shipping", int(s.Id), userId, "U")

	if inDatabaseShipping.ShippingNumber != s.ShippingNumber || inDatabaseShipping.TrackingNumber != s.TrackingNumber {
		go ecommerceControllerUpdateTrackingNumber(inDatabaseShipping.Order, s.TrackingNumber, s.enterprise)
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Shipping) deleteShipping(userId int32) bool {
	if s.Id <= 0 {
		return false
	}

	inDatabaseShipping := getShippingRow(s.Id)
	if inDatabaseShipping.Id <= 0 || inDatabaseShipping.enterprise != s.enterprise {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	packaging := getPackagingByShipping(s.Id, s.enterprise)
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			sqlStatement := `UPDATE sales_order_detail SET status='E' WHERE id=$1`
			_, err := trans.Exec(sqlStatement, detailsPackaged[j].OrderDetail)
			if err != nil {
				log("DB", err.Error())
				trans.Rollback()
				return false
			}

			insertTransactionalLog(s.enterprise, "sales_order_detail", int(detailsPackaged[j].OrderDetail), userId, "U")

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetail)
			ok := setSalesOrderState(s.enterprise, saleOrderDetail.Order, userId, *trans)
			if !ok {
				trans.Rollback()
				return false
			}
		}

		sqlStatement := `UPDATE packaging SET shipping=NULL WHERE id=$1`
		_, err := trans.Exec(sqlStatement, packaging[i].Id)
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}

		insertTransactionalLog(s.enterprise, "packaging", int(packaging[i].Id), userId, "U")
	}

	insertTransactionalLog(s.enterprise, "shipping", int(s.Id), userId, "D")

	sqlStatement := `DELETE FROM public.shipping WHERE id=$1`
	_, err := trans.Exec(sqlStatement, s.Id)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

// ERROR CODES:
// 1. No carrier selected in the order
// 2. A detail has not been completely packaged
// 3. Can't generate delivery note
func generateShippingFromSaleOrder(orderId int64, enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	saleOrder := getSalesOrderRow(orderId)
	if saleOrder.enterprise != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}
	packaging := getPackaging(orderId, enterpriseId)
	if saleOrder.Id <= 0 || len(packaging) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if saleOrder.Carrier == nil || *saleOrder.Carrier <= 0 {
		return OkAndErrorCodeReturn{Ok: false, ErorCode: 1}
	}

	details := getSalesOrderDetail(orderId, enterpriseId)
	for i := 0; i < len(details); i++ {
		if details[i].QuantityPendingPackaging > 0 {
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 2, ExtraData: []string{details[i].ProductName}}
		}
	}

	s := Shipping{}
	s.Order = saleOrder.Id
	s.DeliveryAddress = saleOrder.ShippingAddress
	s.Carrier = *saleOrder.Carrier

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	saleDeliveryNotes := getSalesOrderDeliveryNotes(orderId, enterpriseId)
	if len(saleDeliveryNotes) > 0 {
		s.DeliveryNote = saleDeliveryNotes[0].Id
	} else {
		ok, noteId := deliveryNoteAllSaleOrder(orderId, enterpriseId, userId, trans)
		if !ok.Ok || noteId <= 0 {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false, ErorCode: 3}
		}
		s.DeliveryNote = noteId
	}

	s.enterprise = enterpriseId
	ok, shippingId := s.insertShipping(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(packaging); i++ {
		ok := associatePackagingToShipping(packaging[i].Id, shippingId, enterpriseId, userId, *trans)
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			sqlStatement := `UPDATE sales_order_detail SET status='F' WHERE id=$1`
			_, err := trans.Exec(sqlStatement, detailsPackaged[j].OrderDetail)
			if err != nil {
				log("DB", err.Error())
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}

			insertTransactionalLog(enterpriseId, "sales_order_detail", int(detailsPackaged[j].OrderDetail), userId, "U")

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetail)
			ok := setSalesOrderState(enterpriseId, saleOrderDetail.Order, userId, *trans)
			if !ok {
				trans.Rollback()
				return OkAndErrorCodeReturn{Ok: false}
			}
		}
	}

	///
	transErr = trans.Commit()
	return OkAndErrorCodeReturn{Ok: transErr == nil}
	///
}

// THIS FUNCION DOES NOT OPEN A TRANSACTION
func associatePackagingToShipping(packagingId int64, shippingId int64, enterpriseId int32, userId int32, trans sql.Tx) bool {
	sqlStatement := `UPDATE public.packaging SET shipping=$2 WHERE id=$1`
	_, err := trans.Exec(sqlStatement, packagingId, shippingId)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	sqlStatement = `SELECT weight FROM public.packaging WHERE id=$1`
	row := db.QueryRow(sqlStatement, packagingId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		trans.Rollback()
		return false
	}

	var weight float64
	row.Scan(&weight)

	sqlStatement = `UPDATE public.shipping SET weight=$2, packages_number=packages_number+1 WHERE id=$1`
	_, err = db.Exec(sqlStatement, shippingId, weight)
	if err != nil {
		log("DB", err.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "shipping", int(shippingId), userId, "U")

	return true
}

type ToggleShippingSent struct {
	Ok           bool    `json:"ok"`
	ErrorMessage *string `json:"errorMessage"`
}

func toggleShippingSent(shippingId int64, enterpriseId int32, userId int32) ToggleShippingSent {
	s := getShippingRow(shippingId)
	if s.enterprise != enterpriseId {
		return ToggleShippingSent{Ok: false}
	}
	// it is not allowed to manually set as "sent" if the carrier has a webservice set.
	// it is not allowed to modify the "sent" field if the shipping was colletected by the carrier.
	if s.Id <= 0 || s.Collected {
		return ToggleShippingSent{Ok: false}
	}
	if s.CarrierWebService != "_" {
		ok, errorMessage := s.sendShipping(enterpriseId)
		if ok {
			go ecommerceControllerUpdateTrackingNumber(s.Order, s.TrackingNumber, enterpriseId)
		}
		return ToggleShippingSent{Ok: ok, ErrorMessage: errorMessage}
	}

	sqlStatement := `UPDATE shipping SET sent = NOT sent, date_sent = CASE sent WHEN false THEN CURRENT_TIMESTAMP(3) ELSE NULL END WHERE id = $1`
	_, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		log("DB", err.Error())
		return ToggleShippingSent{Ok: false}
	}

	insertTransactionalLog(enterpriseId, "shipping", int(shippingId), userId, "U")

	return ToggleShippingSent{Ok: true}
}

func (s *Shipping) sendShipping(enterpriseId int32) (bool, *string) {
	switch s.CarrierWebService {
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
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	sqlStatement := `UPDATE shipping SET collected=true WHERE id=$1`
	for i := 0; i < len(shippings); i++ {
		s := getShippingRow(shippings[i])
		if s.Id <= 0 || !s.Sent || s.enterprise != enterpriseId {
			trans.Rollback()
			return false
		}

		_, err := trans.Exec(sqlStatement, shippings[i])
		if err != nil {
			log("DB", err.Error())
			trans.Rollback()
			return false
		}
		insertTransactionalLog(enterpriseId, "shipping", int(shippings[i]), userId, "U")

		p := getPackagingByShipping(shippings[i], enterpriseId)
		for j := 0; j < len(p); j++ {
			for k := 0; k < len(p[j].DetailsPackaged); k++ {
				sqlStatement := `UPDATE sales_order_detail SET status='G' WHERE id=$1`
				_, err := trans.Exec(sqlStatement, p[j].DetailsPackaged[k].OrderDetail)
				if err != nil {
					log("DB", err.Error())
					trans.Rollback()
					return false
				}

				insertTransactionalLog(enterpriseId, "sales_order_detail", int(p[j].DetailsPackaged[k].OrderDetail), userId, "U")
			}
		}

	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}
