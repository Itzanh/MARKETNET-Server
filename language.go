package main

import (
	"strings"
)

type Language struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
	Iso3 string `json:"iso3"`
}

func getLanguages() []Language {
	var languages []Language = make([]Language, 0)
	sqlStatement := `SELECT * FROM public.language ORDER BY id ASC `
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return languages
	}
	for rows.Next() {
		l := Language{}
		rows.Scan(&l.Id, &l.Name, &l.Iso2, &l.Iso3)
		languages = append(languages, l)
	}

	return languages
}

func (l *Language) isValid() bool {
	return !(len(l.Name) == 0 || len(l.Name) > 50 || len(l.Iso2) != 2 || (len(l.Iso3) != 0 && len(l.Iso3) != 3))
}

func (l *Language) insertLanguage() bool {
	if !l.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.language(name, iso_2, iso_3) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, l.Name, l.Iso2, l.Iso3)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (l *Language) updateLanguage() bool {
	if l.Id <= 0 || !l.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.language SET name=$2, iso_2=$3, iso_3=$4 WHERE id=$1`
	res, err := db.Exec(sqlStatement, l.Id, l.Name, l.Iso2, l.Iso3)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (l *Language) deleteLanguage() bool {
	if l.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM language WHERE id = $1`
	res, err := db.Exec(sqlStatement, l.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findLanguageByName(languageName string) []NameInt16 {
	var languages []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.language WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName))
	if err != nil {
		return languages
	}
	for rows.Next() {
		l := NameInt16{}
		rows.Scan(&l.Id, &l.Name)
		languages = append(languages, l)
	}

	return languages
}

func getNameLanguage(id int16) string {
	sqlStatement := `SELECT name FROM public.language WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
