/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "database/sql"

type HSCode struct {
	Id   string `json:"id" gorm:"primaryKey;type:character varying(8);not null:true"`
	Name string `json:"name" gorm:"type:character varying(255);not null:true"`
}

func (h *HSCode) TableName() string {
	return "hs_codes"
}

type HSCodeQuery struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (q *HSCodeQuery) getHSCodes() []HSCode {
	var codes []HSCode = make([]HSCode, 0)
	// get the HS Codes from the database where id or name matches the query using named arguments order by id ascending using dbOrm
	dbOrm.Model(HSCode{}).Where("id LIKE @search OR name LIKE @search", sql.Named("search", "%"+q.Id+"%")).Order("id ASC").Find(&codes)
	return codes
}
