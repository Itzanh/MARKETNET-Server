package main

import (
	"database/sql"
	"strconv"
	"time"
)

type Shipping struct {
	Id              int32      `json:"id"`
	Order           int32      `json:"order"`
	DeliveryNote    int32      `json:"deliveryNote"`
	DeliveryAddress int32      `json:"deliveryAddress"`
	DateCreated     time.Time  `json:"dateCreated"`
	DateSent        *time.Time `json:"dateSent"`
	Sent            bool       `json:"sent"`
	Collected       bool       `json:"collected"`
	National        bool       `json:"national"`
	ShippingNumber  string     `json:"shippingNumber"`
	TrackingNumber  string     `json:"trackingNumber"`
	Carrier         int16      `json:"carrier"`
	Weight          float32    `json:"weight"`
	PackagesNumber  int16      `json:"packagesNumber"`
	CustomerName    string     `json:"customerName"`
	SaleOrderName   string     `json:"saleOrderName"`
	CarrierName     string     `json:"carrierName"`
	Incoterm        *int16     `json:"incoterm"`
	CarrierNotes    string     `json:"carrierNotes"`
	Description     string     `json:"description"`
}

func getShippings() []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier) FROM public.shipping ORDER BY id DESC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return shippings
	}
	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.CustomerName, &s.SaleOrderName, &s.CarrierName)
		shippings = append(shippings, s)
	}

	return shippings
}

func searchShippings(search string) []Shipping {
	var shippings []Shipping = make([]Shipping, 0)
	var rows *sql.Rows
	orderNumber, err := strconv.Atoi(search)
	if err == nil {
		sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier) FROM shipping INNER JOIN sales_order ON sales_order.id=shipping."order" WHERE sales_order.order_number=$1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, orderNumber)
	} else {
		sqlStatement := `SELECT shipping.*,(SELECT name FROM customer WHERE id=(SELECT customer FROM sales_order WHERE id=shipping."order")),(SELECT order_name FROM sales_order WHERE id=shipping."order"),(SELECT name FROM carrier WHERE id=shipping.carrier) FROM shipping INNER JOIN sales_order ON shipping."order"=sales_order.id INNER JOIN customer ON customer.id=sales_order.customer WHERE customer.name ILIKE $1 ORDER BY id DESC`
		rows, err = db.Query(sqlStatement, "%"+search+"%")
	}
	if err != nil {
		return shippings
	}
	for rows.Next() {
		s := Shipping{}
		rows.Scan(&s.Id, &s.Order, &s.DeliveryNote, &s.DeliveryAddress, &s.DateCreated, &s.DateSent, &s.Sent, &s.Collected, &s.National, &s.ShippingNumber, &s.TrackingNumber, &s.Carrier, &s.Weight, &s.PackagesNumber, &s.Incoterm, &s.CarrierNotes, &s.Description, &s.CustomerName, &s.SaleOrderName, &s.CarrierName)
		shippings = append(shippings, s)
	}

	return shippings
}

func (s *Shipping) isValid() bool {
	return !(s.Order <= 0 || s.DeliveryAddress <= 0 || s.Carrier <= 0)
}

func (s *Shipping) insertShipping() (bool, int32) {
	if !s.isValid() {
		return false, 0
	}

	sqlStatement := `INSERT INTO public.shipping("order", delivery_note, delivery_address, "national", carrier, incoterm, carrier_notes, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	row := db.QueryRow(sqlStatement, s.Order, s.DeliveryNote, s.DeliveryAddress, s.National, s.Carrier, s.Incoterm, s.CarrierNotes, s.Description)
	if row.Err() != nil {
		return false, 0
	}

	var shippingId int32
	row.Scan(&shippingId)
	return shippingId > 0, shippingId
}

func (s *Shipping) updateShipping() bool {
	if s.Id <= 0 || !s.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.shipping SET "order"=$2, delivery_note=$3, delivery_address=$4, carrier=$5, incoterm=$6, carrier_notes=$7, description=$8 WHERE id=$1`
	res, err := db.Exec(sqlStatement, s.Id, s.Order, s.DeliveryNote, s.DeliveryAddress, s.Carrier, s.Incoterm, s.CarrierNotes, s.Description)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Shipping) deleteShipping() bool {
	if s.Id <= 0 {
		return false
	}

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	packaging := getPackagingByShipping(s.Id)
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			sqlStatement := `UPDATE sales_order_detail SET status='E' WHERE id=$1`
			_, err := db.Exec(sqlStatement, detailsPackaged[j].OrderDetail)
			if err != nil {
				trans.Rollback()
				return false
			}

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetail)
			ok := setSalesOrderState(saleOrderDetail.Order)
			if !ok {
				trans.Rollback()
				return false
			}
		}

		sqlStatement := `UPDATE packaging SET shipping=NULL WHERE id=$1`
		_, err := db.Exec(sqlStatement, packaging[i].Id)
		if err != nil {
			trans.Rollback()
			return false
		}
	}

	sqlStatement := `DELETE FROM public.shipping WHERE id=$1`
	_, err := db.Exec(sqlStatement, s.Id)
	if err != nil {
		return false
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

func generateShippingFromSaleOrder(orderId int32) bool {
	saleOrder := getSalesOrderRow(orderId)
	packaging := getPackaging(orderId)

	if saleOrder.Id <= 0 || len(packaging) == 0 || saleOrder.Carrier == nil || *saleOrder.Carrier <= 0 {
		return false
	}

	s := Shipping{}
	s.Order = saleOrder.Id
	s.DeliveryAddress = saleOrder.ShippingAddress
	s.Carrier = *saleOrder.Carrier

	///
	trans, transErr := db.Begin()
	if transErr != nil {
		return false
	}
	///

	saleDeliveryNotes := getSalesOrderDeliveryNotes(orderId)
	if len(saleDeliveryNotes) > 0 {
		s.DeliveryNote = saleDeliveryNotes[0].Id
	} else {
		ok, noteId := deliveryNoteAllSaleOrder(orderId)
		if !ok || noteId <= 0 {
			trans.Rollback()
			return false
		}
		s.DeliveryNote = noteId
	}

	ok, shippingId := s.insertShipping()
	if !ok {
		trans.Rollback()
		return false
	}
	for i := 0; i < len(packaging); i++ {
		ok := associatePackagingToShipping(packaging[i].Id, shippingId)
		if !ok {
			trans.Rollback()
			return false
		}
	}
	for i := 0; i < len(packaging); i++ {
		detailsPackaged := packaging[i].DetailsPackaged
		for j := 0; j < len(detailsPackaged); j++ {
			sqlStatement := `UPDATE sales_order_detail SET status='F' WHERE id=$1`
			_, err := db.Exec(sqlStatement, detailsPackaged[j].OrderDetail)
			if err != nil {
				trans.Rollback()
				return false
			}

			saleOrderDetail := getSalesOrderDetailRow(detailsPackaged[j].OrderDetail)
			ok := setSalesOrderState(saleOrderDetail.Order)
			if !ok {
				trans.Rollback()
				return false
			}
		}
	}

	///
	transErr = trans.Commit()
	return transErr == nil
	///
}

// THIS FUNCION DOES NOT OPEN A TRANSACTION
func associatePackagingToShipping(packagingId int32, shippingId int32) bool {
	sqlStatement := `UPDATE public.packaging SET shipping=$2 WHERE id=$1`
	_, err := db.Exec(sqlStatement, packagingId, shippingId)
	if err != nil {
		return false
	}

	sqlStatement = `SELECT weight FROM public.packaging WHERE id=$1`
	row := db.QueryRow(sqlStatement, packagingId)
	if row.Err() != nil {
		return false
	}

	var weight float32
	row.Scan(&weight)

	sqlStatement = `UPDATE public.shipping SET weight=$2, packages_number=packages_number+1 WHERE id=$1`
	_, err = db.Exec(sqlStatement, shippingId, weight)
	return err == nil
}

func toggleShippingSent(shippingId int32) bool {
	sqlStatement := `UPDATE shipping SET sent = NOT sent, date_sent = CASE sent WHEN false THEN CURRENT_TIMESTAMP(3) ELSE NULL END WHERE id = $1`
	_, err := db.Exec(sqlStatement, shippingId)
	return err == nil
}
