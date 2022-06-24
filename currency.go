/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Currency struct {
	Id           int32     `json:"id" gorm:"primaryKey;index:currency_id_enterprise,unique:true,priority:1"`
	Name         string    `json:"name" gorm:"column:name;type:character varying(150);not null:true"`
	Sign         string    `json:"sign" gorm:"column:sign;type:character(3);not null:true"`
	IsoCode      string    `json:"isoCode" gorm:"column:iso_code;type:character(3);not null:true;index:currency_iso_code,unique:true,priority:2"`
	IsoNum       int16     `json:"isoNum" gorm:"not null:true;index:currency_num,unique:true,priority:2"`
	Change       float64   `json:"change" gorm:"column:exchange;type:numeric(14,6);not null:true"`
	ExchangeDate time.Time `json:"exchangeDate" gorm:"column:exchange_date;type:date;not null:true"`
	EnterpriseId int32     `json:"-" gorm:"column:enterprise;not null:true;index:currency_id_enterprise,unique:true,priority:2;index:currency_iso_code,unique:true,priority:1;index:currency_num,unique:true,priority:1"`
	Enterprise   Settings  `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (c *Currency) TableName() string {
	return "currency"
}

func getCurrencies(enterpriseId int32) []Currency {
	var currencies []Currency = make([]Currency, 0)
	result := dbOrm.Model(&Currency{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&currencies)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return currencies
	}

	return currencies
}

func (c *Currency) isValid() bool {
	return !(len(c.Name) == 0 || len(c.Name) > 150 || len([]rune(c.Sign)) > 3 || len(c.IsoCode) > 3 || c.IsoNum < 0 || c.Change <= 0)
}

func (c *Currency) BeforeCreate(tx *gorm.DB) (err error) {
	var currency Currency
	tx.Model(&Currency{}).Last(&currency)
	c.Id = currency.Id + 1
	return nil
}

func (c *Currency) insertCurrency() bool {
	if !c.isValid() {
		return false
	}

	// DEFAULTS
	c.ExchangeDate = time.Now()

	result := dbOrm.Create(&c)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Currency) updateCurrency() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	var currency Currency
	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).First(&currency)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	currency.Name = c.Name
	currency.Sign = c.Sign
	currency.IsoCode = c.IsoCode
	currency.IsoNum = c.IsoNum
	currency.Change = c.Change

	result = dbOrm.Save(&currency)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (c *Currency) deleteCurrency() bool {
	if c.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", c.Id, c.EnterpriseId).Delete(&Currency{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func getCurrencyExchange(currencyId int32) float64 {
	var currency Currency
	result := dbOrm.Where("id = ?", currencyId).First(&currency)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return 0
	}

	return currency.Change
}

func findCurrencyByName(currencyName string, enterpriseId int32) []NameInt32 {
	var currencies []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Currency{}).Where("(UPPER(name) LIKE ? || '%') AND enterprise = ?", strings.ToUpper(currencyName), enterpriseId).Limit(10).Find(&currencies)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return currencies
	}

	return currencies
}

func locateCurrency(enterpriseId int32) []NameInt32 {
	var currencies []NameInt32 = make([]NameInt32, 0)
	result := dbOrm.Model(&Currency{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Limit(10).Find(&currencies)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return currencies
	}

	return currencies
}

func updateCurrencyExchange(enterpriseId int32) {
	if getSettingsRecordById(enterpriseId).Currency != "E" {
		return
	}
	currencies := getCurrencies(enterpriseId)

	for i := 0; i < len(currencies); i++ {
		if len(currencies[i].IsoCode) == 0 || currencies[i].IsoCode == "EUR" {
			continue
		}

		now := time.Now()
		now = now.AddDate(0, 0, -1)
		currentDate := now.Format("2006-01-02")
		resp, err := http.Get(getSettingsRecordById(enterpriseId).CurrencyECBurl + "D." + currencies[i].IsoCode + ".EUR.SP00.A?startPeriod=" + currentDate + "&endPeriod=" + currentDate)
		if err != nil {
			log("ECB Exchange", err.Error())
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log("ECB Exchange", err.Error())
			return
		}

		if len(body) == 0 { // there is no data for this date
			continue
		}
		xml := string(body)

		// check that there is a field with the date
		dateIndexOf := strings.Index(xml, "<generic:ObsDimension value=\"")
		if dateIndexOf <= 0 {
			continue
		}
		dateIndexOf += len("<generic:ObsDimension value=\"")

		// chech that there is a field with the exchange rate
		valueIndexOf := strings.Index(xml, "<generic:ObsValue value=\"")
		if valueIndexOf <= 0 {
			continue
		}
		valueIndexOf += len("<generic:ObsValue value=\"")

		// extract the field and parse the data
		dateString := xml[dateIndexOf : dateIndexOf+10]
		valueString := xml[valueIndexOf : valueIndexOf+strings.Index(xml[valueIndexOf:], "\"")]

		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			log("ECB Exchange", err.Error())
			continue
		}
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			log("ECB Exchange", err.Error())
			continue
		}

		currency := currencies[i]
		currency.Change = value
		currency.ExchangeDate = date
		currency.updateCurrency()
	}
}
