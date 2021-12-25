package main

import (
	"time"
)

type ShippingTag struct {
	Id          int64     `json:"id"`
	Shipping    int64     `json:"shipping"`
	DateCreated time.Time `json:"dateCreated"`
	Label       []byte    `json:"label"`
	enterprise  int32
}

func getShippingTags(shippingId int64, enterpriseId int32) []ShippingTag {
	tags := make([]ShippingTag, 0)
	sqlStatement := `SELECT * FROM public.shipping_tag WHERE shipping=$1 AND enterprise=$2 ORDER BY date_created ASC`
	rows, err := db.Query(sqlStatement, shippingId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return tags
	}
	defer rows.Close()

	for rows.Next() {
		t := ShippingTag{}
		rows.Scan(&t.Id, &t.Shipping, &t.DateCreated, &t.Label, &t.enterprise)
		tags = append(tags, t)
	}

	return tags
}

func (t *ShippingTag) insertShippingTag() bool {
	sqlStatement := `INSERT INTO public.shipping_tag(shipping, label, enterprise) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, t.Shipping, t.Label, t.enterprise)
	return err == nil
}

func deleteAllShippingTags(enterpriseId int32) {
	sqlStatement := `DELETE FROM shipping_tag WHERE enterprise=$1`
	_, err := db.Exec(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
	}
}
