/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "strings"

type ReportTemplateTranslation struct {
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Key          string   `json:"key" gorm:"primaryKey;column:key;not null:true;type:character varying(50)"`
	LanguageId   int32    `json:"languageId" gorm:"primaryKey;column:language;not null:true"`
	Language     Language `json:"language" gorm:"foreignKey:LanguageId,EnterpriseId;references:Id,EnterpriseId"`
	Translation  string   `json:"translation" gorm:"column:translation;not null:true;type:character varying(255)"`
}

func (t *ReportTemplateTranslation) TableName() string {
	return "report_template_translation"
}

func getReportTemplateTranslations(enterpriseId int32) []ReportTemplateTranslation {
	var translations []ReportTemplateTranslation = make([]ReportTemplateTranslation, 0)
	// get the report template translations from the database for the given enterprise id sorted by key and language ascending using dbOrm
	result := dbOrm.Where("enterprise = ?", enterpriseId).Order("key ASC, language ASC").Preload("Language").Find(&translations)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return translations
}

func (t *ReportTemplateTranslation) isValid() bool {
	return !(t.EnterpriseId <= 0 || len(t.Key) == 0 || len(t.Key) > 50 || t.LanguageId <= 0 || len(t.Translation) == 0 || len(t.Translation) > 255)
}

func (t *ReportTemplateTranslation) insertReportTemplateTranslation() bool {
	if !t.isValid() {
		return false
	}

	// insert the report template translation into the database using dbOrm
	result := dbOrm.Create(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}

func (t *ReportTemplateTranslation) updateReportTemplateTranslation() bool {
	if !t.isValid() {
		return false
	}

	// update the report template translation in the database using dbOrm
	result := dbOrm.Model(&ReportTemplateTranslation{}).Where("enterprise = ? AND key = ? AND language = ?", t.EnterpriseId, t.Key, t.LanguageId).UpdateColumn("translation", t.Translation)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (t *ReportTemplateTranslation) deleteReportTemplateTranslation() bool {
	if t.EnterpriseId <= 0 || len(t.Key) == 0 || t.LanguageId <= 0 {
		return false
	}

	// delete the report template translation from the database using dbOrm
	result := dbOrm.Where("enterprise = ? AND key = ? AND language = ?", t.EnterpriseId, t.Key, t.LanguageId).Delete(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func translateReport(reportContent string, languageId int32, enterpriseId int32) string {
	// get a single report template translation from the database for the given enterprise id and language id using dbOrm
	// get all the translations for the given enterprise id and language id using dbOrm
	// translate the report content using the translations
	// return the translated report content
	var translation ReportTemplateTranslation
	var translationCount int64
	result := dbOrm.Where("enterprise = ? AND language = ?", enterpriseId, languageId).Count(&translationCount).First(&translation)
	if translationCount == 0 {
		result = dbOrm.Where("enterprise = ? AND language = ?", enterpriseId, 1).First(&translation)
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	var translations []ReportTemplateTranslation
	result = dbOrm.Where("enterprise = ? AND language = ?", enterpriseId, languageId).Count(&translationCount).Find(&translations)
	if translationCount == 0 {
		result = dbOrm.Where("enterprise = ? AND language = ?", enterpriseId, 1).Find(&translations)
	}
	if result.Error != nil {
		log("DB", result.Error.Error())
		return ""
	}

	for _, translation := range translations {
		reportContent = strings.Replace(reportContent, "{"+translation.Key+"}", translation.Translation, -1)
	}

	return reportContent
}
