package main

import (
	"time"
)

type ShippingTag struct {
	Id          int64     `json:"id"`
	Shipping    int32     `json:"shipping"`
	DateCreated time.Time `json:"dateCreated"`
	Label       []byte    `json:"label"`
}

func getShippingTags(shippingId int32) []ShippingTag {
	tags := make([]ShippingTag, 0)
	sqlStatement := `SELECT * FROM public.shipping_tag WHERE shipping=$1 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, shippingId)
	if err != nil {
		log("DB", err.Error())
		return tags
	}
	for rows.Next() {
		t := ShippingTag{}
		rows.Scan(&t.Id, &t.Shipping, &t.DateCreated, &t.Label)
		tags = append(tags, t)
	}

	return tags
}

func (t *ShippingTag) insertShippingTag() bool {
	sqlStatement := `INSERT INTO public.shipping_tag(shipping, label) VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, t.Shipping, t.Label)
	return err == nil
}

func deleteAllShippingTags() {
	sqlStatement := `DELETE FROM shipping_tag`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		log("DB", err.Error())
	}
}
