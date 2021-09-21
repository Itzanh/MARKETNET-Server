package main

import (
	"strings"
)

type Language struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Iso2       string `json:"iso2"`
	Iso3       string `json:"iso3"`
	enterprise int32
}

func getLanguages(enterpriseId int32) []Language {
	var languages []Language = make([]Language, 0)
	sqlStatement := `SELECT * FROM public.language WHERE enterprise=$1 ORDER BY id ASC `
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return languages
	}
	for rows.Next() {
		l := Language{}
		rows.Scan(&l.Id, &l.Name, &l.Iso2, &l.Iso3, &l.enterprise)
		languages = append(languages, l)
	}

	return languages
}

func (l *Language) isValid() bool {
	return !(len(l.Name) == 0 || len(l.Name) > 50 || len(l.Iso2) != 2 || (len(l.Iso3) != 0 && len(l.Iso3) != 3))
}

func searchLanguages(search string, enterpriseId int32) []Language {
	var languages []Language = make([]Language, 0)
	sqlStatement := `SELECT * FROM language WHERE (name ILIKE $1 OR iso_2 = UPPER($2) OR iso_3 = UPPER($2)) AND enterprise=$3 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+search+"%", search, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return languages
	}
	for rows.Next() {
		l := Language{}
		rows.Scan(&l.Id, &l.Name, &l.Iso2, &l.Iso3, &l.enterprise)
		languages = append(languages, l)
	}

	return languages
}

func (l *Language) insertLanguage() bool {
	if !l.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.language(name, iso_2, iso_3, enterprise) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, l.Name, l.Iso2, l.Iso3, l.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (l *Language) updateLanguage() bool {
	if l.Id <= 0 || !l.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.language SET name=$2, iso_2=$3, iso_3=$4 WHERE id=$1 AND enterprise=$5`
	res, err := db.Exec(sqlStatement, l.Id, l.Name, l.Iso2, l.Iso3, l.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (l *Language) deleteLanguage() bool {
	if l.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM language WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, l.Id, l.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findLanguageByName(languageName string, enterpriseId int32) []NameInt16 {
	var languages []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.language WHERE (UPPER(name) LIKE $1 || '%') AND enterprise=$2 ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(languageName), enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return languages
	}
	for rows.Next() {
		l := NameInt16{}
		rows.Scan(&l.Id, &l.Name)
		languages = append(languages, l)
	}

	return languages
}

func getNameLanguage(id int32, enterpriseId int32) string {
	sqlStatement := `SELECT name FROM public.language WHERE id=$1 AND enterprise=$2`
	row := db.QueryRow(sqlStatement, id, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
