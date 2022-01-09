package main

import (
	"net/smtp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendEmail(destinationAddress string, destinationAddressName string, subject string, innerText string, enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	if s.Email != "_" {
		el := EmailLog{EmailFrom: s.EmailFrom, NameFrom: s.NameFrom, DestinationEmail: destinationAddress, DestinationName: destinationAddressName, Subject: subject, Content: innerText, enterprise: enterpriseId}
		el.insertEmailLog()
	}

	if s.Email == "_" {
		return false
	} else if s.Email == "S" {
		sendEmailSendgrid(s.SendGridKey, s.EmailFrom, s.NameFrom, destinationAddress, destinationAddressName, subject, innerText)
	} else if s.Email == "T" {
		sendEmailSMTP(s.SMTPIdentity, s.SMTPUsername, s.SMTPPassword, s.SMTPHostname, destinationAddress, subject, innerText)
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

func sendEmailSMTP(identiy string, username string, password string, smtpServer string, destinationAddress string, subject string, innerText string) bool {
	auth := smtp.PlainAuth(identiy, username, password, smtpServer[:strings.Index(smtpServer, ":")])

	to := []string{destinationAddress}
	msg := []byte("To: " + destinationAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		innerText + "\r\n")

	err := smtp.SendMail(smtpServer, auth, username, to, msg)

	if err != nil {
		log("SMTP", err.Error())
	}

	return err == nil
}

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
