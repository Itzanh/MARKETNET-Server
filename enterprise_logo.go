package main

import (
	"net/http"
	"strings"
)

const LOGO_MAX_SIZE = 1000000                                                     // 1 Mb
const LOGO_ALLOWED_MIME_TYPES = "image/jpeg;image/png;image/svg+xml;image/x-icon" // Allowed mime types separated by ";"
const LOGO_ALLOWED_MIME_TYPES_SEP = ";"

// returns: image, mime type
func getEnterpriseLogo(enterpriseId int32) ([]byte, string) {
	sqlStatement := `SELECT logo, mime_type FROM public.enterprise_logo WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return nil, ""
	}

	var logo []byte
	var mimeType string
	row.Scan(&logo, &mimeType)
	return logo, mimeType
}

func setEnterpriseLogo(enterpriseId int32, logo []byte) bool {
	// check size
	if len(logo) > LOGO_MAX_SIZE || len(logo) == 0 {
		return false
	}

	// check the mine type
	mimeType := http.DetectContentType(logo)
	allowedMimeTypes := strings.Split(LOGO_ALLOWED_MIME_TYPES, LOGO_ALLOWED_MIME_TYPES_SEP)
	var mimeTypeFound bool = false
	for i := 0; i < len(allowedMimeTypes); i++ {
		if allowedMimeTypes[i] == mimeType {
			mimeTypeFound = true
			break
		}
	}
	if !mimeTypeFound {
		return false
	}

	// the logo row already exists? used for creating/updating
	sqlStatement := `SELECT COUNT(*) FROM public.enterprise_logo WHERE enterprise=$1`
	row := db.QueryRow(sqlStatement, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var count int32
	row.Scan(&count)

	if count > 0 { // update the existing row, the user already has an image selected and it's changing it
		sqlStatement := `UPDATE public.enterprise_logo SET logo=$2, mime_type=$3 WHERE enterprise=$1`
		_, err := db.Exec(sqlStatement, enterpriseId, logo, mimeType)
		if err != nil {
			log("DB", err.Error())
			return false
		}
	} else { // create a new row, the user did not have an image selected and it's adding one
		sqlStatement := `INSERT INTO public.enterprise_logo(enterprise, logo, mime_type) VALUES ($1, $2, $3)`
		_, err := db.Exec(sqlStatement, enterpriseId, logo, mimeType)
		if err != nil {
			log("DB", err.Error())
			return false
		}
	}
	return true
}

func deleteEnterpriseLogo(enterpriseId int32) bool {
	sqlStatement := `DELETE FROM public.enterprise_logo WHERE enterprise=$1`
	_, err := db.Exec(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}
