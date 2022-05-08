package main

import (
	"net/http"

	"gorm.io/gorm"
)

type CustomFields struct {
	Id            int64     `json:"id"`
	EnterpriseId  int32     `json:"-" gorm:"column:enterprise;not null:true"`
	Enterprise    Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	ProductId     *int32    `json:"product" gorm:"column:product"`
	Product       *Product  `json:"-" gorm:"foreignKey:ProductId,EnterpriseId;references:Id,EnterpriseId"`
	CustomerId    *int32    `json:"customer" gorm:"column:customer"`
	Customer      *Customer `json:"-" gorm:"foreignKey:CustomerId,EnterpriseId;references:Id,EnterpriseId"`
	SupplierId    *int32    `json:"supplier" gorm:"column:supplier"`
	Supplier      *Supplier `json:"-" gorm:"foreignKey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	Name          string    `json:"name" gorm:"type:character varying(255);not null:true"`
	FieldType     int16     `json:"fieldType" gorm:"not null:true"` // 1 = Short text, 2 = Long text, 3 = Number, 4 = Boolean, 5 = Image, 6 = File
	ValueString   *string   `json:"valueString" gorm:"type:text"`
	ValueNumber   *float64  `json:"valueNumber" gorm:"type:numeric(14,6)"`
	ValueBoolean  *bool     `json:"valueBoolean"`
	ValueBinary   []byte    `json:"valueBinary" gorm:"type:bytea"`
	FileName      *string   `json:"fileName" gorm:"type:character varying(255)"`
	FileSize      *int32    `json:"fileSize" gorm:"type:integer"`
	ImageMimeType *string   `json:"imageMimeType" gorm:"type:character varying(255)"`
}

func (CustomFields) TableName() string {
	return "custom_fields"
}

func (f *CustomFields) queryIsValid() bool {
	return !(f.ProductId == nil && f.CustomerId == nil && f.SupplierId == nil)
}

func (f *CustomFields) getCustomFields() []CustomFields {
	var fields []CustomFields = make([]CustomFields, 0)
	if !f.queryIsValid() {
		return fields
	}

	// create a database cursor to get the custom fields using dbOrm
	cursor := dbOrm.Model(CustomFields{}).Where("enterprise = ?", f.EnterpriseId)

	if f.ProductId != nil {
		cursor = cursor.Where("product = ?", f.ProductId)
	} else if f.CustomerId != nil {
		cursor = cursor.Where("customer = ?", f.CustomerId)
	} else if f.SupplierId != nil {
		cursor = cursor.Where("supplier = ?", f.SupplierId)
	}
	cursor.Order("id ASC").Find(&fields)
	return fields
}

func (f *CustomFields) isValid() bool {
	if f.FieldType == 5 || f.FieldType == 6 {
		if len(f.ValueBinary) > 5000000 { // 5Mb
			return false
		}
	}

	if f.FieldType == 5 {
		mimeType := http.DetectContentType(f.ValueBinary)
		if mimeType != "image/jpeg" && mimeType != "image/jpg" && mimeType != "image/png" && mimeType != "image/gif" && mimeType != "image/bmp" && mimeType != "image/webp" {
			return false
		}
		f.ImageMimeType = &mimeType
	}

	return !((f.ProductId == nil && f.CustomerId == nil && f.SupplierId == nil) || len(f.Name) == 0 || len(f.Name) > 255 || f.FieldType < 1 || f.FieldType > 6 || (f.FileName != nil && len(*f.FileName) > 255) || ((f.FieldType == 5 || f.FieldType == 6) && (f.FileName == nil || len(*f.FileName) == 0)) || (f.ValueString != nil && len(*f.ValueString) > 80000) || ((f.FieldType == 1 || f.FieldType == 2) && (f.ValueString == nil || len(*f.ValueString) == 0)))
}

func (f *CustomFields) cleanUpCustomFields() {
	if f.FieldType != 1 && f.FieldType != 2 {
		f.ValueString = nil
	} else if f.FieldType != 3 {
		f.ValueNumber = nil
	} else if f.FieldType != 4 {
		f.ValueBoolean = nil
	} else if f.FieldType != 5 && f.FieldType != 6 {
		f.ValueBinary = nil
		f.FileName = nil
		f.FileSize = nil
		f.ImageMimeType = nil
	}

	if f.FieldType == 5 || f.FieldType == 6 {
		valueBinaryLen := int32(len(f.ValueBinary))
		f.FileSize = &valueBinaryLen
	}
}

func (f *CustomFields) BeforeCreate(tx *gorm.DB) (err error) {
	var customFields CustomFields
	tx.Model(&CustomFields{}).Last(&customFields)
	f.Id = customFields.Id + 1
	return nil
}

func (f *CustomFields) insertCustomFields() bool {
	if !f.isValid() {
		return false
	}
	f.cleanUpCustomFields()

	// No more than 3 binary content files are allowed
	customFields := f.getCustomFields()
	var countBinaryFiles uint8
	for i := 0; i < len(customFields); i++ {
		if customFields[i].FieldType == 5 || customFields[i].FieldType == 6 {
			countBinaryFiles++
		}
	}
	if countBinaryFiles >= 3 {
		return false
	}

	result := dbOrm.Create(&f)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (f *CustomFields) updateCustomFields() bool {
	if !f.isValid() || f.Id <= 0 {
		return false
	}
	f.cleanUpCustomFields()

	// get a single custom field from the database by id and enterprise using dbOrm
	var customField CustomFields
	result := dbOrm.Model(CustomFields{}).Where("enterprise = ? AND id = ?", f.EnterpriseId, f.Id).First(&customField)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	// copy all the fields from the database to the new custom field
	customField.Name = f.Name
	customField.FieldType = f.FieldType
	customField.ValueString = f.ValueString
	customField.ValueNumber = f.ValueNumber
	customField.ValueBoolean = f.ValueBoolean
	customField.ValueBinary = f.ValueBinary
	customField.FileName = f.FileName
	customField.FileSize = f.FileSize
	customField.ImageMimeType = f.ImageMimeType

	// update the custom field in the database using dbOrm
	result = dbOrm.Save(&customField)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (f *CustomFields) deleteCustomFields() bool {
	if f.Id <= 0 {
		return false
	}

	// delete a single custom field from the database by id and enterprise using dbOrm
	result := dbOrm.Where("id = ? AND enterprise = ?", f.Id, f.EnterpriseId).Delete(CustomFields{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
