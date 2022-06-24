/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

type ReportTemplate struct {
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Key          string   `json:"key" gorm:"primaryKey;column:key;not null:true;type:character varying(50)"`
	Html         string   `json:"html" gorm:"column:html;not null:true;type:text"`
}

func (r *ReportTemplate) TableName() string {
	return "report_template"
}

func getReportTemplates(enterpriseId int32) []ReportTemplate {
	templates := make([]ReportTemplate, 0)
	// get all the report templates for the enterprise using dbOrm
	result := dbOrm.Model(&ReportTemplate{}).Where("enterprise = ?", enterpriseId).Order("key ASC").Find(&templates)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return templates
}

func getReportTemplate(enterpriseId int32, key string) ReportTemplate {
	t := ReportTemplate{}
	// get the report template for the enterprise and the given key using dbOrm
	result := dbOrm.Model(&ReportTemplate{}).Where("enterprise = ? AND key = ?", enterpriseId, key).First(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return t
}

// Must NOT be callable from the web client!
func (r ReportTemplate) insertReportTemplate() {
	// insert the report template using dbOrm
	result := dbOrm.Create(&r)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
}

func (r *ReportTemplate) updateReportTemplate() bool {
	if len(r.Key) == 0 || len(r.Html) == 0 || len(r.Html) > 5000000 {
		return false
	}

	// get a single report template from the database for the given enterprise id and key using dbOrm
	var t ReportTemplate
	result := dbOrm.Model(&ReportTemplate{}).Where("enterprise = ? AND key = ?", r.EnterpriseId, r.Key).First(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	t.Html = r.Html

	// update the report template using dbOrm
	result = dbOrm.Save(&t)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}
	return true
}
