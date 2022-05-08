package main

import (
	"net/http"
	"strings"
)

const LOGO_MAX_SIZE = 1000000                                                     // 1 Mb
const LOGO_ALLOWED_MIME_TYPES = "image/jpeg;image/png;image/svg+xml;image/x-icon" // Allowed mime types separated by ";"
const LOGO_ALLOWED_MIME_TYPES_SEP = ";"

type EnterpriseLogo struct {
	EnterpriseId int32    `gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Logo         []byte   `gorm:"column:logo;not null:true"`
	MimeType     string   `gorm:"column:mime_type;not null:true;type:character varying(150)"`
}

func (l *EnterpriseLogo) TableName() string {
	return "enterprise_logo"
}

// returns: image, mime type
func getEnterpriseLogo(enterpriseId int32) ([]byte, string) {
	var logo EnterpriseLogo
	result := dbOrm.Model(&EnterpriseLogo{}).Where("enterprise = ?", enterpriseId).First(&logo)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil, ""
	}
	return logo.Logo, logo.MimeType
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
	var count int64
	result := dbOrm.Model(&EnterpriseLogo{}).Where("enterprise = ?", enterpriseId).Count(&count)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	if count > 0 { // update the existing row, the user already has an image selected and it's changing it
		result = dbOrm.Model(&EnterpriseLogo{}).Where("enterprise = ?", enterpriseId).Updates(map[string]interface{}{
			"logo":      logo,
			"mime_type": mimeType,
		})
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false
		}
	} else { // create a new row, the user did not have an image selected and it's adding one
		var logo EnterpriseLogo = EnterpriseLogo{
			EnterpriseId: enterpriseId,
			Logo:         logo,
			MimeType:     mimeType,
		}
		result = dbOrm.Create(&logo)
		if result.Error != nil {
			log("DB", result.Error.Error())
			return false
		}
	}
	return true
}

func deleteEnterpriseLogo(enterpriseId int32) bool {
	result := dbOrm.Model(&EnterpriseLogo{}).Where("enterprise = ?", enterpriseId).Delete(&EnterpriseLogo{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
