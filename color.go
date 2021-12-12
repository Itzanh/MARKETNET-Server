package main

import "strings"

type Color struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	HexColor   string `json:"hexColor"`
	enterprise int32
}

func getColor(enterpriseId int32) []Color {
	var color []Color = make([]Color, 0)
	sqlStatement := `SELECT * FROM public.color WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return color
	}
	for rows.Next() {
		c := Color{}
		rows.Scan(&c.Id, &c.Name, &c.HexColor, &c.enterprise)
		color = append(color, c)
	}

	return color
}

func (c *Color) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 100 || len(c.HexColor) > 6)
}

func (c *Color) insertColor() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.color(name, hex_color, enterprise) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, c.Name, c.HexColor, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Color) updateColor() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.color SET name=$2, hex_color=$3 WHERE id=$1 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.HexColor, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Color) deleteColor() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.color WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findColorByName(colorName string, enterpriseId int32) []NameInt16 {
	var colors []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.color WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(colorName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return colors
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		colors = append(colors, c)
	}

	return colors
}

func getNameColor(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.color WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}

type ColorLocate struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func locateColor(enterpriseId int32) []ColorLocate {
	var color []ColorLocate = make([]ColorLocate, 0)
	sqlStatement := `SELECT id,name FROM public.color WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return color
	}
	for rows.Next() {
		c := ColorLocate{}
		rows.Scan(&c.Id, &c.Name)
		color = append(color, c)
	}

	return color
}
