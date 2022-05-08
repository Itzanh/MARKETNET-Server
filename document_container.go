package main

import (
	"time"

	"gorm.io/gorm"
)

type DocumentContainer struct {
	Id                  int32     `json:"id" gorm:"index:_document_container_id_enterprise,unique:true,priority:1"`
	Name                string    `json:"name" gorm:"type:character varying(50);not null:true"`
	DateCreated         time.Time `json:"dateCreated" gorm:"type:timestamp(3) with time zone;not null:true"`
	Path                string    `json:"path" gorm:"type:character varying(250);not null:true"`
	MaxFileSize         int32     `json:"maxFileSize" gorm:"not null:true"`
	DisallowedMimeTypes string    `json:"disallowedMimeTypes" gorm:"type:character varying(250);not null:true"`
	AllowedMimeTypes    string    `json:"allowedMimeTypes" gorm:"type:character varying(250);not null:true"`
	EnterpriseId        int32     `json:"-" gorm:"column:enterprise;not null:true;index:_document_container_id_enterprise,unique:true,priority:2"`
	Enterprise          Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	UsedStorage         int64     `json:"usedStorage" gorm:"not null:true"`
	MaxStorage          int64     `json:"maxStorage" gorm:"not null:true"`
}

func (dc *DocumentContainer) TableName() string {
	return "document_container"
}

func getDocumentContainer(enterpriseId int32) []DocumentContainer {
	var containters []DocumentContainer = make([]DocumentContainer, 0)
	dbOrm.Model(&DocumentContainer{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&containters)
	return containters
}

func getDocumentContainerRow(containerId int32) DocumentContainer {
	d := DocumentContainer{}
	dbOrm.Model(&DocumentContainer{}).Where("id = ?", containerId).First(&d)
	return d
}

func (c *DocumentContainer) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 50 || len(c.Path) == 0 || len(c.Path) > 250 || c.MaxFileSize <= 0 || len(c.DisallowedMimeTypes) > 250 || len(c.AllowedMimeTypes) > 250)
}

func (d *DocumentContainer) BeforeCreate(tx *gorm.DB) (err error) {
	var documentContainer DocumentContainer
	tx.Model(&DocumentContainer{}).Last(&documentContainer)
	d.Id = documentContainer.Id + 1
	return nil
}

func (d *DocumentContainer) insertDocumentContainer() bool {
	if !d.isValid() {
		return false
	}

	if isParameterPresent("--saas") {
		return false
	}

	d.DateCreated = time.Now()
	d.UsedStorage = 0
	d.MaxStorage = 0

	result := dbOrm.Create(&d)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (d *DocumentContainer) updateDocumentContainer() bool {
	if d.Id <= 0 || !d.isValid() {
		return false
	}

	if isParameterPresent("--saas") {
		d.MaxStorage = getDocumentContainerRow(d.Id).MaxStorage
	}

	documentContainer := DocumentContainer{}
	result := dbOrm.Model(&DocumentContainer{}).Where("id = ? AND enterprise = ?", d.Id, d.EnterpriseId).First(&documentContainer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	documentContainer.Name = d.Name
	documentContainer.Path = d.Path
	documentContainer.MaxFileSize = d.MaxFileSize
	documentContainer.DisallowedMimeTypes = d.DisallowedMimeTypes
	documentContainer.AllowedMimeTypes = d.AllowedMimeTypes
	documentContainer.MaxStorage = d.MaxStorage

	result = dbOrm.Save(&documentContainer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (d *DocumentContainer) deleteDocumentContainer() bool {
	if d.Id <= 0 {
		return false
	}

	if isParameterPresent("--saas") {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", d.Id, d.EnterpriseId).Delete(&DocumentContainer{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (d *DocumentContainer) updateUsedStorage(sizeInBytes int32, trans *gorm.DB) bool {
	if d.Id <= 0 {
		return false
	}

	documentContainer := DocumentContainer{}
	result := trans.Model(&DocumentContainer{}).Where("id = ?", d.Id).First(&documentContainer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	documentContainer.UsedStorage += int64(sizeInBytes)

	result = trans.Save(&documentContainer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}
	return true
}

type DocumentContainerLocate struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

func locateDocumentContainer(enterpriseId int32) []NameInt32 {
	var containters []NameInt32 = make([]NameInt32, 0)
	dbOrm.Model(&DocumentContainer{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&containters)
	return containters
}
