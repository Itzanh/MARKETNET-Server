package main

import "strings"

type Color struct {
	Id       int16  `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

func getColor() []Color {
	var color []Color = make([]Color, 0)
	sqlStatement := `SELECT * FROM public.color ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return color
	}
	for rows.Next() {
		c := Color{}
		rows.Scan(&c.Id, &c.Name, &c.HexColor)
		color = append(color, c)
	}

	return color
}

func (c *Color) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || len(c.HexColor) > 6)
}

func (c *Color) insertColor() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.color(name, hex_color) VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, c.Name, c.HexColor)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Color) updateColor() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.color SET name=$2, hex_color=$3 WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id, c.Name, c.HexColor)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (c *Color) deleteColor() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.color WHERE id=$1`
	res, err := db.Exec(sqlStatement, c.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findColorByName(colorName string) []NameInt16 {
	var colors []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.color WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(colorName))
	if err != nil {
		return colors
	}
	for rows.Next() {
		c := NameInt16{}
		rows.Scan(&c.Id, &c.Name)
		colors = append(colors, c)
	}

	return colors
}

func getNameColor(id int16) string {
	sqlStatement := `SELECT name FROM public.color WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
