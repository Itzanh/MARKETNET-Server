package main

import (
	"net/http"
)

type CustomFields struct {
	Id            int64    `json:"id"`
	Product       *int32   `json:"product"`
	Customer      *int32   `json:"customer"`
	Supplier      *int32   `json:"supplier"`
	Name          string   `json:"name"`
	FieldType     int16    `json:"fieldType"` // 1 = Short text, 2 = Long text, 3 = Number, 4 = Boolean, 5 = Image, 6 = File
	ValueString   *string  `json:"valueString"`
	ValueNumber   *float64 `json:"valueNumber"`
	ValueBoolean  *bool    `json:"valueBoolean"`
	ValueBinary   []byte   `json:"valueBinary"`
	FileName      *string  `json:"fileName"`
	FileSize      *int32   `json:"fileSize"`
	ImageMimeType *string  `json:"imageMimeType"`
	enterprise    int32
}

func (f *CustomFields) queryIsValid() bool {
	return !(f.Product == nil && f.Customer == nil && f.Supplier == nil)
}

func (f *CustomFields) getCustomFields() []CustomFields {
	var fields []CustomFields = make([]CustomFields, 0)
	if !f.queryIsValid() {
		return fields
	}

	var interfaces []interface{} = make([]interface{}, 0)

	sqlStatement := `SELECT * FROM public.custom_fields`
	if f.Product != nil {
		sqlStatement += ` WHERE product=$1`
		interfaces = append(interfaces, f.Product)
	} else if f.Customer != nil {
		sqlStatement += ` WHERE customer=$1`
		interfaces = append(interfaces, f.Customer)
	} else if f.Supplier != nil {
		sqlStatement += ` WHERE supplier=$1`
		interfaces = append(interfaces, f.Supplier)
	}
	sqlStatement += `AND enterprise=$2`
	interfaces = append(interfaces, f.enterprise)
	sqlStatement += ` ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, interfaces...)
	if err != nil {
		log("DB", err.Error())
		return fields
	}

	for rows.Next() {
		f := CustomFields{}
		rows.Scan(&f.Id, &f.enterprise, &f.Product, &f.Customer, &f.Supplier, &f.Name, &f.FieldType, &f.ValueString, &f.ValueNumber, &f.ValueBoolean, &f.ValueBinary, &f.FileName, &f.FileSize, &f.ImageMimeType)
		fields = append(fields, f)
	}

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

	return !((f.Product == nil && f.Customer == nil && f.Supplier == nil) || len(f.Name) == 0 || len(f.Name) > 255 || f.FieldType < 1 || f.FieldType > 6 || (f.FileName != nil && len(*f.FileName) > 255) || ((f.FieldType == 5 || f.FieldType == 6) && (f.FileName == nil || len(*f.FileName) == 0)) || (f.ValueString != nil && len(*f.ValueString) > 80000) || ((f.FieldType == 1 || f.FieldType == 2) && (f.ValueString == nil || len(*f.ValueString) == 0)))
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

	sqlStatement := `INSERT INTO public.custom_fields(enterprise, product, customer, supplier, name, field_type, value_string, value_number, value_boolean, value_binary, file_name, file_size, image_mime_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := db.Exec(sqlStatement, f.enterprise, f.Product, f.Customer, f.Supplier, f.Name, f.FieldType, f.ValueString, f.ValueNumber, f.ValueBoolean, f.ValueBinary, f.FileName, f.FileSize, f.ImageMimeType)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (f *CustomFields) updateCustomFields() bool {
	if !f.isValid() || f.Id <= 0 {
		return false
	}
	f.cleanUpCustomFields()

	sqlStatement := `UPDATE public.custom_fields SET product=$2, customer=$3, supplier=$4, name=$5, value_string=$6, value_number=$7, value_boolean=$8, value_binary=$9, file_name=$10, file_size=$11, image_mime_type=$13 WHERE id=$1 AND enterprise=$12`
	_, err := db.Exec(sqlStatement, f.Id, f.Product, f.Customer, f.Supplier, f.Name, f.ValueString, f.ValueNumber, f.ValueBoolean, f.ValueBinary, f.FileName, f.FileSize, f.enterprise, f.ImageMimeType)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (f *CustomFields) deleteCustomFields() bool {
	if f.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.custom_fields WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, f.Id, f.enterprise)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}
