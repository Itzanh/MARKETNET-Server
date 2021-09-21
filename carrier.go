package main

import "strings"

type Carrier struct {
	Id                      int32   `json:"id"`
	Name                    string  `json:"name"`
	MaxWeight               float32 `json:"maxWeight"`
	MaxWidth                float32 `json:"maxWidth"`
	MaxHeight               float32 `json:"maxHeight"`
	MaxDepth                float32 `json:"maxDepth"`
	MaxPackages             int16   `json:"maxPackages"`
	Phone                   string  `json:"phone"`
	Email                   string  `json:"email"`
	Web                     string  `json:"web"`
	Off                     bool    `json:"off"`
	PrestaShopId            int32   `json:"prestaShopId"`
	Pallets                 bool    `json:"pallets"`
	Webservice              string  `json:"webservice"`
	SendcloudUrl            string  `json:"sendcloudUrl"`
	SendcloudKey            string  `json:"sendcloudKey"`
	SendcloudSecret         string  `json:"sendcloudSecret"`
	SendcloudShippingMethod int32   `json:"sendcloudShippingMethod"`
	SendcloudSenderAddress  int64   `json:"sendcloudSenderAddress"`
	enterprise              int32
}

func getCariers(enterpriseId int32) []Carrier {
	var carriers []Carrier = make([]Carrier, 0)
	sqlStatement := `SELECT * FROM public.carrier WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return carriers
	}
	for rows.Next() {
		c := Carrier{}
		rows.Scan(&c.Id, &c.Name, &c.MaxWeight, &c.MaxWidth, &c.MaxHeight, &c.MaxDepth, &c.MaxPackages, &c.Phone, &c.Email, &c.Web, &c.Off, &c.PrestaShopId, &c.Pallets, &c.Webservice, &c.SendcloudUrl, &c.SendcloudKey, &c.SendcloudSecret, &c.SendcloudShippingMethod, &c.SendcloudSenderAddress, &c.enterprise)
		carriers = append(carriers, c)
	}

	return carriers
}

func getCarierRow(id int32) Carrier {
	sqlStatement := `SELECT * FROM public.carrier WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Carrier{}
	}

	c := Carrier{}
	row.Scan(&c.Id, &c.Name, &c.MaxWeight, &c.MaxWidth, &c.MaxHeight, &c.MaxDepth, &c.MaxPackages, &c.Phone, &c.Email, &c.Web, &c.Off, &c.PrestaShopId, &c.Pallets, &c.Webservice, &c.SendcloudUrl, &c.SendcloudKey, &c.SendcloudSecret, &c.SendcloudShippingMethod, &c.SendcloudSenderAddress, &c.enterprise)

	return c
}

func (c *Carrier) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || c.MaxWeight < 0 || c.MaxWidth < 0 || c.MaxHeight < 0 || c.MaxDepth < 0 || c.MaxPackages < 0 || len(c.Phone) > 15 || len(c.Email) > 100 || len(c.Web) > 100 || len(c.Webservice) != 1 || (c.Webservice != "_" && c.Webservice != "S") || len(c.SendcloudUrl) > 75 || (len(c.SendcloudKey) != 0 && len(c.SendcloudKey) != 32) || (len(c.SendcloudSecret) != 0 && len(c.SendcloudSecret) != 32) || c.SendcloudShippingMethod < 0 || c.SendcloudSenderAddress < 0)
}

func (c *Carrier) insertCarrier() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.carrier(name, max_weight, max_width, max_height, max_depth, max_packages, phone, email, web, off, ps_id, pallets, webservice, sendcloud_url, sendcloud_key, sendcloud_secret, sendcloud_shipping_method, sendcloud_sender_address, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	res, err := db.Exec(sqlStatement, c.Name, c.MaxWeight, c.MaxWidth, c.MaxHeight, c.MaxDepth, c.MaxPackages, c.Phone, c.Email, c.Web, c.Off, c.PrestaShopId, c.Pallets, c.Webservice, c.SendcloudUrl, c.SendcloudKey, c.SendcloudSecret, c.SendcloudShippingMethod, c.SendcloudSenderAddress, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Carrier) updateCarrier() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.carrier SET name=$2, max_weight=$3, max_width=$4, max_height=$5, max_depth=$6, max_packages=$7, phone=$8, email=$9, web=$10, off=$11, pallets=$12, webservice=$13, sendcloud_url=$14, sendcloud_key=$15, sendcloud_secret=$16, sendcloud_shipping_method=$17, sendcloud_sender_address=$18 WHERE id=$1 AND enterprise=$19`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.MaxWeight, c.MaxWidth, c.MaxHeight, c.MaxDepth, c.MaxPackages, c.Phone, c.Email, c.Web, c.Off, c.Pallets, c.Webservice, c.SendcloudUrl, c.SendcloudKey, c.SendcloudSecret, c.SendcloudShippingMethod, c.SendcloudSenderAddress, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Carrier) deleteCarrier() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.carrier WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCarrierByName(languageName string, enterpriseId int32) []NameInt16 {
	var carriers []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.carrier WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return carriers
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		carriers = append(carriers, c)
	}

	return carriers
}

func getNameCarrier(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.carrier WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
