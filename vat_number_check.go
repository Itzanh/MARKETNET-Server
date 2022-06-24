/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

const EU_VAT_NUMBER_WEBSERVICE_URL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

type CheckVatNumber struct {
	CountryIsoCode2 string `json:"countryIsoCode2"`
	VATNumber       string `json:"VATNumber"`
}

func (c *CheckVatNumber) isValid() bool {
	return !(len(c.CountryIsoCode2) != 2 || len(c.VATNumber) < 3 || len(c.VATNumber) > 50)
}

func checkVatNumber(countryIsoCode2 string, vatNumber string) OkAndErrorCodeReturn {
	xml := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
	<soapenv:Header/>
	<soapenv:Body>
	   <urn:checkVat>
		  <urn:countryCode>%1</urn:countryCode>
		  <urn:vatNumber>%2</urn:vatNumber>
	   </urn:checkVat>
	</soapenv:Body>
</soapenv:Envelope>`
	xml = strings.Replace(xml, "%1", countryIsoCode2, 1)
	xml = strings.Replace(xml, "%2", vatNumber, 1)

	resp, err := http.Post(EU_VAT_NUMBER_WEBSERVICE_URL, "application/xml", bytes.NewBuffer([]byte(xml)))
	if err != nil {
		log("VIES", err.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log("VIES", err.Error())
		return OkAndErrorCodeReturn{Ok: false}
	}

	if strings.Contains(string(result), "<valid>true</valid>") {
		return OkAndErrorCodeReturn{Ok: true, ErrorCode: 1}
	} else if strings.Contains(string(result), "<valid>false</valid>") {
		return OkAndErrorCodeReturn{Ok: true, ErrorCode: 2}
	} else {
		return OkAndErrorCodeReturn{Ok: false}
	}
}
