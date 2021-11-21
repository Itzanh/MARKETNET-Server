package main

import "time"

type ShippingStatusHistory struct {
	Id          int64     `json:"id"`
	Shipping    int64     `json:"shipping"`
	StatusId    int16     `json:"statusId"`
	Message     string    `json:"message"`
	Delivered   bool      `json:"delivered"`
	DateCreated time.Time `json:"dateCreated"`
}

func getShippingStatusHistory(enterpriseId int32, shippingId int64) []ShippingStatusHistory {
	var shippingStatusHistory []ShippingStatusHistory = make([]ShippingStatusHistory, 0)
	sqlStatement := `SELECT * FROM public.shipping_status_history WHERE (SELECT enterprise FROM shipping WHERE shipping.id=shipping_status_history.shipping)=$1 AND shipping=$2 ORDER BY date_created DESC`
	rows, err := db.Query(sqlStatement, enterpriseId, shippingId)
	if err != nil {
		log("DB", err.Error())
		return shippingStatusHistory
	}

	for rows.Next() {
		s := ShippingStatusHistory{}
		rows.Scan(&s.Id, &s.Shipping, &s.StatusId, &s.Message, &s.Delivered, &s.DateCreated)
		shippingStatusHistory = append(shippingStatusHistory, s)
	}

	return shippingStatusHistory
}
