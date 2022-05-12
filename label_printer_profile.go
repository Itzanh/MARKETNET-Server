package main

import (
	"gorm.io/gorm"
)

type LabelPrinterProfile struct {
	Id                              int32    `json:"id"`
	Type                            string   `json:"type" gorm:"type:character(1);not null:true;index:label_printer_profile_type_enterprise,unique:true,priority:1,where:active"` // E = EAN13, C = Code128, D = DataMatrix
	Active                          bool     `json:"active" gorm:"not null:true"`
	ProductBarCodeLabelWidth        int16    `json:"productBarCodeLabelWidth" gorm:"column:product_barcode_label_width;not null:true"`
	ProductBarCodeLabelHeight       int16    `json:"productBarCodeLabelHeight" gorm:"column:product_barcode_label_height;not null:true"`
	ProductBarCodeLabelSize         int16    `json:"productBarCodeLabelSize" gorm:"column:product_barcode_label_size;not null:true"`
	ProductBarCodeLabelMarginTop    int16    `json:"productBarCodeLabelMarginTop" gorm:"column:product_barcode_label_margin_top;not null:true"`
	ProductBarCodeLabelMarginBottom int16    `json:"productBarCodeLabelMarginBottom" gorm:"column:product_barcode_label_margin_bottom;not null:true"`
	ProductBarCodeLabelMarginLeft   int16    `json:"productBarCodeLabelMarginLeft" gorm:"column:product_barcode_label_margin_left;not null:true"`
	ProductBarCodeLabelMarginRight  int16    `json:"productBarCodeLabelMarginRight" gorm:"column:product_barcode_label_margin_right;not null:true"`
	EnterpriseId                    int32    `json:"-" gorm:"column:enterprise;not null:true;index:label_printer_profile_type_enterprise,unique:true,priority:2,where:active"`
	Enterprise                      Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (LabelPrinterProfile) TableName() string {
	return "label_printer_profile"
}

func getLabelPrinterProfiles(enterpriseId int32) []LabelPrinterProfile {
	var labelPrinterProfiles []LabelPrinterProfile
	result := dbOrm.Where("enterprise = ?", enterpriseId).Find(&labelPrinterProfiles)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return labelPrinterProfiles
}

func getLabelPrinterProfileByEnterpriseTypeAndActive(enterpriseId int32, profileType string) *LabelPrinterProfile {
	var labelPrinterProfile *LabelPrinterProfile
	var rowCount int64
	result := dbOrm.Model(&LabelPrinterProfile{}).Where("enterprise = ? AND type = ? AND active = ?", enterpriseId, profileType, true).Count(&rowCount).First(&labelPrinterProfile)
	if rowCount == 0 {
		return nil
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	return labelPrinterProfile
}

func (l *LabelPrinterProfile) isValid() bool {
	return !((l.Type != "E" && l.Type != "C" && l.Type != "D") || l.ProductBarCodeLabelWidth < 0 || l.ProductBarCodeLabelHeight < 0 || l.ProductBarCodeLabelSize < 0 || l.ProductBarCodeLabelMarginTop < 0 || l.ProductBarCodeLabelMarginBottom < 0 || l.ProductBarCodeLabelMarginLeft < 0 || l.ProductBarCodeLabelMarginRight < 0)
}

func (p *LabelPrinterProfile) BeforeCreate(tx *gorm.DB) (err error) {
	var labelPrinterProfile LabelPrinterProfile
	tx.Model(&LabelPrinterProfile{}).Last(&labelPrinterProfile)
	p.Id = labelPrinterProfile.Id + 1
	return nil
}

func (l *LabelPrinterProfile) insertLabelPrinterProfile() bool {
	if !l.isValid() {
		return false
	}

	result := dbOrm.Create(&l)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (p *LabelPrinterProfile) updateLabelPrinterProfile() bool {
	var labelPrinterProfile LabelPrinterProfile
	result := dbOrm.Model(&LabelPrinterProfile{}).Where("id = ? AND enterprise = ?", p.Id, p.EnterpriseId).First(&labelPrinterProfile)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	labelPrinterProfile.Active = p.Active
	labelPrinterProfile.ProductBarCodeLabelWidth = p.ProductBarCodeLabelWidth
	labelPrinterProfile.ProductBarCodeLabelHeight = p.ProductBarCodeLabelHeight
	labelPrinterProfile.ProductBarCodeLabelSize = p.ProductBarCodeLabelSize
	labelPrinterProfile.ProductBarCodeLabelMarginTop = p.ProductBarCodeLabelMarginTop
	labelPrinterProfile.ProductBarCodeLabelMarginBottom = p.ProductBarCodeLabelMarginBottom
	labelPrinterProfile.ProductBarCodeLabelMarginLeft = p.ProductBarCodeLabelMarginLeft
	labelPrinterProfile.ProductBarCodeLabelMarginRight = p.ProductBarCodeLabelMarginRight

	result = dbOrm.Save(&labelPrinterProfile)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (l *LabelPrinterProfile) deleteLabelPrinterProfile() bool {
	result := dbOrm.Model(&LabelPrinterProfile{}).Where("id = ? AND enterprise = ?", l.Id, l.EnterpriseId).Delete(&LabelPrinterProfile{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
