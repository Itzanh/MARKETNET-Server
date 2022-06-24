/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendEmail(destinationAddress string, destinationAddressName string, subject string, innerText string, enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	if s.SettingsEmail.Email != "_" {
		el := EmailLog{EmailFrom: s.SettingsEmail.EmailFrom, NameFrom: s.SettingsEmail.NameFrom, DestinationEmail: destinationAddress, DestinationName: destinationAddressName, Subject: subject, Content: innerText, EnterpriseId: enterpriseId}
		el.insertEmailLog()
	}

	if s.SettingsEmail.Email == "_" {
		return false
	} else if s.SettingsEmail.Email == "S" {
		sendEmailSendgrid(s.SettingsEmail.SendGridKey, s.SettingsEmail.EmailFrom, s.SettingsEmail.NameFrom, destinationAddress, destinationAddressName, subject, innerText)
	} else if s.SettingsEmail.Email == "T" {
		if s.SettingsEmail.SMTPSTARTTLS {
			sendEmailSMTPwithSTARTTLS(s.SettingsEmail.SMTPIdentity, s.SettingsEmail.SMTPUsername, s.SettingsEmail.SMTPPassword, s.SettingsEmail.SMTPHostname, destinationAddress, subject, innerText, s.SettingsEmail.SMTPReplyTo)
		} else {
			sendEmailSMTPPlainAuth(s.SettingsEmail.SMTPIdentity, s.SettingsEmail.SMTPUsername, s.SettingsEmail.SMTPPassword, s.SettingsEmail.SMTPHostname, destinationAddress, subject, innerText, s.SettingsEmail.SMTPReplyTo)
		}
	}
	return false
}

func sendEmailSendgrid(key string, fromAddress string, fromAddressName string, destinationAddress string, destinationAddressName string, subject string, innerText string) bool {
	from := mail.NewEmail(fromAddressName, fromAddress)
	to := mail.NewEmail(destinationAddressName, destinationAddress)
	message := mail.NewSingleEmail(from, subject, to, strip.StripTags(innerText), innerText)
	client := sendgrid.NewSendClient(key)
	_, err := client.Send(message)

	if err != nil {
		log("SENDGRID", err.Error())
	}

	return err == nil
}

func sendEmailSMTPPlainAuth(identiy, username, password, smtpServer, destinationAddress, subject, innerText, replyTo string) bool {
	auth := smtp.PlainAuth(identiy, username, password, smtpServer[:strings.Index(smtpServer, ":")])

	if len(replyTo) > 0 {
		replyTo = "Reply-To: " + replyTo + "\r\n"
	}

	to := []string{destinationAddress}
	msg := []byte("From: " + username + "\r\n" +
		"To: " + destinationAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		replyTo +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		innerText + "\r\n")

	err := smtp.SendMail(smtpServer, auth, username, to, msg)

	if err != nil {
		log("SMTP", err.Error())
	}

	return err == nil
}

/* SMTP EMAIL WITH START TLS */

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}

func sendEmailSMTPwithSTARTTLS(identiy, username, password, smtpServer, destinationAddress, subject, innerText, replyTo string) bool {
	conn, err := net.Dial("tcp", smtpServer)
	if err != nil {
		log("SMTP", err.Error())
	}

	if len(replyTo) > 0 {
		replyTo = "Reply-To: " + replyTo + "\r\n"
	}

	c, err := smtp.NewClient(conn, smtpServer[:strings.Index(smtpServer, ":")])
	if err != nil {
		log("SMTP", err.Error())
	}

	tlsconfig := &tls.Config{
		ServerName: smtpServer[:strings.Index(smtpServer, ":")],
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		log("SMTP", err.Error())
	}

	auth := LoginAuth(username, password)

	if err = c.Auth(auth); err != nil {
		log("SMTP", err.Error())
	}

	to := []string{destinationAddress}
	msg := []byte("From: " + username + "\r\n" +
		"To: " + destinationAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		replyTo +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		innerText + "\r\n")

	err = smtp.SendMail(smtpServer, auth, username, to, msg)
	if err != nil {
		log("SMTP", err.Error())
		return false
	}

	return err == nil
}

/*/ SMTP EMAIL WITH START TLS /*/

type EmailInfo struct {
	DestinationAddress     string `json:"destinationAddress"`
	DestinationAddressName string `json:"destinationAddressName"`
	Subject                string `json:"subject"`
	ReportId               string `json:"reportId"`
	ReportDataId           int32  `json:"reportDataId"`
	Language               int32  `json:"language"` // can be 0
}

func (e *EmailInfo) isValid() bool {
	return !(len(e.DestinationAddress) == 0 || len(e.DestinationAddressName) == 0 || len(e.Subject) == 0 || len(e.ReportId) == 0 || e.ReportDataId <= 0 || e.Language < 0)
}

func (e *EmailInfo) sendEmail(enterpriseId int32) bool {
	if !e.isValid() {
		return false
	}

	var report []byte
	switch e.ReportId {
	case "SALES_ORDER":
		report = reportSalesOrder(int(e.ReportDataId), false, enterpriseId, e.Language)
	case "SALES_INVOICE":
		report = reportSalesInvoice(int(e.ReportDataId), false, enterpriseId)
	case "SALES_DELIVERY_NOTE":
		report = reportSalesDeliveryNote(int(e.ReportDataId), false, enterpriseId)
	case "PURCHASE_ORDER":
		report = reportPurchaseOrder(int(e.ReportDataId), false, enterpriseId)
	case "SALES_ORDER_DIGITAL_PRODUCT_DATA":
		report = reportSalesOrderDigitalProductDetails(int(e.ReportDataId), false, enterpriseId)
	default:
		return false
	}

	return sendEmail(e.DestinationAddress, e.DestinationAddressName, e.Subject, string(report), enterpriseId)
}
