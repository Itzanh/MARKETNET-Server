package main

import "time"

type DocumentContainer struct {
	Id                  int32     `json:"id"`
	Name                string    `json:"name"`
	DateCreated         time.Time `json:"dateCreated"`
	Path                string    `json:"path"`
	MaxFileSize         int32     `json:"maxFileSize"`
	DisallowedMimeTypes string    `json:"disallowedMimeTypes"`
	AllowedMimeTypes    string    `json:"allowedMimeTypes"`
	enterprise          int32
}

func getDocumentContainer(enterpriseId int32) []DocumentContainer {
	var containters []DocumentContainer = make([]DocumentContainer, 0)
	sqlStatement := `SELECT * FROM document_container WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return containters
	}
	for rows.Next() {
		d := DocumentContainer{}
		rows.Scan(&d.Id, &d.Name, &d.DateCreated, &d.Path, &d.MaxFileSize, &d.DisallowedMimeTypes, &d.AllowedMimeTypes, &d.enterprise)
		containters = append(containters, d)
	}

	return containters
}

func getDocumentContainerRow(containerId int32) DocumentContainer {
	sqlStatement := `SELECT * FROM document_container WHERE id=$1`
	row := db.QueryRow(sqlStatement, containerId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return DocumentContainer{}
	}

	d := DocumentContainer{}
	row.Scan(&d.Id, &d.Name, &d.DateCreated, &d.Path, &d.MaxFileSize, &d.DisallowedMimeTypes, &d.AllowedMimeTypes, &d.enterprise)

	return d
}

func (c *DocumentContainer) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || len(c.Path) == 0 || len(c.Path) > 250 || c.MaxFileSize <= 0 || len(c.DisallowedMimeTypes) > 250 || len(c.AllowedMimeTypes) > 250)
}

func (d *DocumentContainer) insertDocumentContainer() bool {
	if !d.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.document_container(name, path, max_file_size, disallowed_mime_types, allowed_mime_types, enterprise) VALUES ($1, $2, $3, $4, $5, $6)`
	res, err := db.Exec(sqlStatement, d.Name, d.Path, d.MaxFileSize, d.DisallowedMimeTypes, d.AllowedMimeTypes, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (d *DocumentContainer) updateDocumentContainer() bool {
	if d.Id <= 0 || !d.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.document_container SET name=$2, path=$3, max_file_size=$4, disallowed_mime_types=$5, allowed_mime_types=$6, enterprise=$7 WHERE id=$1`
	res, err := db.Exec(sqlStatement, d.Id, d.Name, d.Path, d.MaxFileSize, d.DisallowedMimeTypes, d.AllowedMimeTypes, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (d *DocumentContainer) deleteDocumentContainer() bool {
	if d.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.document_container WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, d.Id, d.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type DocumentContainerLocate struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

func locateDocumentContainer(enterpriseId int32) []DocumentContainerLocate {
	var containters []DocumentContainerLocate = make([]DocumentContainerLocate, 0)
	sqlStatement := `SELECT id,name FROM document_container WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return containters
	}
	for rows.Next() {
		d := DocumentContainerLocate{}
		rows.Scan(&d.Id, &d.Name)
		containters = append(containters, d)
	}

	return containters
}
