package main

import "strings"

type Carrier struct {
	Id           int16   `json:"id"`
	Name         string  `json:"name"`
	MaxWeight    float32 `json:"maxWeight"`
	MaxWidth     float32 `json:"maxWidth"`
	MaxHeight    float32 `json:"maxHeight"`
	MaxDepth     float32 `json:"maxDepth"`
	MaxPackages  int16   `json:"maxPackages"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Web          string  `json:"web"`
	Off          bool    `json:"off"`
	PrestaShopId int32   `json:"prestaShopId"`
	Pallets      bool    `json:"pallets"`
	Webservice   string  `json:"webservice"`
}

func getCariers() []Carrier {
	var carriers []Carrier = make([]Carrier, 0)
	sqlStatement := `SELECT * FROM public.carrier ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return carriers
	}
	for rows.Next() {
		c := Carrier{}
		rows.Scan(&c.Id, &c.Name, &c.MaxWeight, &c.MaxWidth, &c.MaxHeight, &c.MaxDepth, &c.MaxPackages, &c.Phone, &c.Email, &c.Web, &c.Off, &c.PrestaShopId, &c.Pallets, &c.Webservice)
		carriers = append(carriers, c)
	}

	return carriers
}

func (c *Carrier) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || c.MaxWeight < 0 || c.MaxWidth < 0 || c.MaxHeight < 0 || c.MaxDepth < 0 || c.MaxPackages < 0 || len(c.Phone) > 15 || len(c.Email) > 100 || len(c.Web) > 100 || len(c.Webservice) != 1 || (c.Webservice != "_"))
}

func (c *Carrier) insertCarrier() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.carrier(name, max_weight, max_width, max_height, max_depth, max_packages, phone, email, web, off, ps_id, pallets, webservice) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	res, err := db.Exec(sqlStatement, c.Name, c.MaxWeight, c.MaxWidth, c.MaxHeight, c.MaxDepth, c.MaxPackages, c.Phone, c.Email, c.Web, c.Off, c.PrestaShopId, c.Pallets, c.Webservice)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Carrier) updateCarrier() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.carrier SET name=$2, max_weight=$3, max_width=$4, max_height=$5, max_depth=$6, max_packages=$7, phone=$8, email=$9, web=$10, off=$11, pallets=$12, webservice=$13 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.MaxWeight, c.MaxWidth, c.MaxHeight, c.MaxDepth, c.MaxPackages, c.Phone, c.Email, c.Web, c.Off, c.Pallets, c.Webservice)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Carrier) deleteCarrier() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.carrier WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findCarrierByName(languageName string) []NameInt16 {
	var carriers []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.carrier WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return carriers
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		carriers = append(carriers, c)
	}

	return carriers
}

func getNameCarrier(id int16) string {
	sqlStatement := `SELECT name FROM public.carrier WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
