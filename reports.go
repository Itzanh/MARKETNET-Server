package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func generateReport(w http.ResponseWriter, r *http.Request) {
	report, ok := r.URL.Query()["report"]
	if !ok {
		return
	}
	force_print, ok := r.URL.Query()["force_print"]
	forcePrint := false
	if ok {
		forcePrint = force_print[0] == "1"
	}

	if report[0] == "LOGO" {
		data := getLogo()
		w.Write(data)
		mimeType := http.DetectContentType(data)
		w.Header().Add("Content-Type", mimeType)
		return
	}

	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) != 36 {
		return
	}
	idString, ok := r.URL.Query()["id"]
	if !ok {
		return
	}
	id, err := strconv.Atoi(idString[0])
	if err != nil {
		return
	}

	ok = consumeToken(token[0])
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch report[0] {
	case "SALES_ORDER":
		w.Write(reportSalesOrder(id, forcePrint))
	case "SALES_INVOICE":
		w.Write(reportSalesInvoice(id, forcePrint))
	case "SALES_DELIVERY_NOTE":
		w.Write(reportSalesDeliveryNote(id, forcePrint))
	case "PURCHASE_ORDER":
		w.Write(reportPurchaseOrder(id, forcePrint))
	}

}

func getLogo() []byte {
	content, err := ioutil.ReadFile("./reports/logo.png")
	if err != nil {
		return nil
	}
	return content
}

func reportSalesOrder(id int, forcePrint bool) []byte {
	s := getSalesOrderRow(int32(id))

	paymentMethod := getNamePaymentMethod(s.PaymentMethod)
	customer := getNameCustomer(s.Customer)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State)
	}
	countryName := getNameCountry(address.Country)
	details := getSalesOrderDetail(s.Id)

	content, err := ioutil.ReadFile("./reports/sales_order.html")
	if err != nil {
		return nil
	}

	html := string(content)

	html = strings.Replace(html, "$$order_number$$", s.OrderName, 1)
	html = strings.Replace(html, "$$order_date$$", s.DateCreated.Format("2006-01-02 15:04:05"), 1)
	html = strings.Replace(html, "$$order_reference$$", s.Reference, 1)
	html = strings.Replace(html, "$$order_payment_method_name$$", paymentMethod, 1)
	html = strings.Replace(html, "$$order_customer_name$$", customer, 1)
	html = strings.Replace(html, "$$address_address$$", address.Address, 1)
	html = strings.Replace(html, "$$address_address2$$", address.Address2, 1)
	html = strings.Replace(html, "$$address_city$$", address.City, 1)
	html = strings.Replace(html, "$$address_postcode$$", address.ZipCode, 1)
	html = strings.Replace(html, "$$address_state$$", stateName, 1)
	html = strings.Replace(html, "$$address_country$$", countryName, 1)
	html = strings.Replace(html, "$$order_notes$$", s.Notes, 1)
	html = strings.Replace(html, "$$order_total_products$$", fmt.Sprintf("%.2f", s.TotalProducts), 1)
	html = strings.Replace(html, "$$order_vat_amount$$", fmt.Sprintf("%.2f", s.VatAmount), 1)
	html = strings.Replace(html, "$$order_discount_percent$$", fmt.Sprintf("%.2f", s.DiscountPercent), 1)
	html = strings.Replace(html, "$$order_fix_discount$$", fmt.Sprintf("%.2f", s.FixDiscount), 1)
	html = strings.Replace(html, "$$order_shipping_price$$", fmt.Sprintf("%.2f", s.ShippingPrice), 1)
	html = strings.Replace(html, "$$order_shipping_discount$$", fmt.Sprintf("%.2f", s.ShippingDiscount), 1)
	html = strings.Replace(html, "$$order_total_with_discount$$", fmt.Sprintf("%.2f", s.TotalWithDiscount), 1)
	html = strings.Replace(html, "$$order_total_amount$$", fmt.Sprintf("%.2f", s.TotalAmount), 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(details); i++ {
		detailHtml := detailHtmlTemplate

		product := getNameProduct(details[i].Product)

		detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_unit_price$$", fmt.Sprintf("%.2f", details[i].Price), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_vat$$", fmt.Sprintf("%.2f", details[i].VatPercent), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_total$$", fmt.Sprintf("%.2f", details[i].TotalAmount), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}

func reportSalesInvoice(id int, forcePrint bool) []byte {
	i := getSalesInvoiceRow(int32(id))

	paymentMethod := getNamePaymentMethod(i.PaymentMethod)
	customer := getNameCustomer(i.Customer)
	address := getAddressRow(i.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State)
	}
	countryName := getNameCountry(address.Country)
	details := getSalesOrderDetail(i.Id)

	content, err := ioutil.ReadFile("./reports/sales_invoice.html")
	if err != nil {
		return nil
	}

	html := string(content)

	html = strings.Replace(html, "$$invoice_number$$", i.InvoiceName, 1)
	html = strings.Replace(html, "$$invoice_date$$", i.DateCreated.Format("2006-01-02 15:04:05"), 1)
	html = strings.Replace(html, "$$invoice_payment_method_name$$", paymentMethod, 1)
	html = strings.Replace(html, "$$invoice_customer_name$$", customer, 1)
	html = strings.Replace(html, "$$address_address$$", address.Address, 1)
	html = strings.Replace(html, "$$address_address2$$", address.Address2, 1)
	html = strings.Replace(html, "$$address_city$$", address.City, 1)
	html = strings.Replace(html, "$$address_postcode$$", address.ZipCode, 1)
	html = strings.Replace(html, "$$address_state$$", stateName, 1)
	html = strings.Replace(html, "$$address_country$$", countryName, 1)
	html = strings.Replace(html, "$$invoice_total_products$$", fmt.Sprintf("%.2f", i.TotalProducts), 1)
	html = strings.Replace(html, "$$invoice_vat_amount$$", fmt.Sprintf("%.2f", i.VatAmount), 1)
	html = strings.Replace(html, "$$invoice_discount_percent$$", fmt.Sprintf("%.2f", i.DiscountPercent), 1)
	html = strings.Replace(html, "$$invoice_fix_discount$$", fmt.Sprintf("%.2f", i.FixDiscount), 1)
	html = strings.Replace(html, "$$invoice_shipping_price$$", fmt.Sprintf("%.2f", i.ShippingPrice), 1)
	html = strings.Replace(html, "$$invoice_shipping_discount$$", fmt.Sprintf("%.2f", i.ShippingDiscount), 1)
	html = strings.Replace(html, "$$invoice_total_with_discount$$", fmt.Sprintf("%.2f", i.TotalWithDiscount), 1)
	html = strings.Replace(html, "$$invoice_total_amount$$", fmt.Sprintf("%.2f", i.TotalAmount), 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(details); i++ {
		detailHtml := detailHtmlTemplate

		product := getNameProduct(details[i].Product)

		detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_unit_price$$", fmt.Sprintf("%.2f", details[i].Price), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_vat$$", fmt.Sprintf("%.2f", details[i].VatPercent), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_total$$", fmt.Sprintf("%.2f", details[i].TotalAmount), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}

func reportSalesDeliveryNote(id int, forcePrint bool) []byte {
	n := getSalesDeliveryNoteRow(int32(id))

	paymentMethod := getNamePaymentMethod(n.PaymentMethod)
	customer := getNameCustomer(n.Customer)
	address := getAddressRow(n.ShippingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State)
	}
	countryName := getNameCountry(address.Country)
	details := getSalesOrderDetail(n.Id)

	content, err := ioutil.ReadFile("./reports/sales_delivery_note.html")
	if err != nil {
		return nil
	}

	html := string(content)

	html = strings.Replace(html, "$$note_number$$", n.DeliveryNoteName, 1)
	html = strings.Replace(html, "$$note_date$$", n.DateCreated.Format("2006-01-02 15:04:05"), 1)
	html = strings.Replace(html, "$$note_payment_method_name$$", paymentMethod, 1)
	html = strings.Replace(html, "$$note_customer_name$$", customer, 1)
	html = strings.Replace(html, "$$address_address$$", address.Address, 1)
	html = strings.Replace(html, "$$address_address2$$", address.Address2, 1)
	html = strings.Replace(html, "$$address_city$$", address.City, 1)
	html = strings.Replace(html, "$$address_postcode$$", address.ZipCode, 1)
	html = strings.Replace(html, "$$address_state$$", stateName, 1)
	html = strings.Replace(html, "$$address_country$$", countryName, 1)
	html = strings.Replace(html, "$$note_total_products$$", fmt.Sprintf("%.2f", n.TotalProducts), 1)
	html = strings.Replace(html, "$$note_vat_amount$$", fmt.Sprintf("%.2f", n.TotalVat), 1)
	html = strings.Replace(html, "$$note_discount_percent$$", fmt.Sprintf("%.2f", n.DiscountPercent), 1)
	html = strings.Replace(html, "$$note_fix_discount$$", fmt.Sprintf("%.2f", n.FixDiscount), 1)
	html = strings.Replace(html, "$$note_shipping_price$$", fmt.Sprintf("%.2f", n.ShippingPrice), 1)
	html = strings.Replace(html, "$$note_shipping_discount$$", fmt.Sprintf("%.2f", n.ShippingDiscount), 1)
	html = strings.Replace(html, "$$note_total_with_discount$$", fmt.Sprintf("%.2f", n.TotalWithDiscount), 1)
	html = strings.Replace(html, "$$note_total_amount$$", fmt.Sprintf("%.2f", n.TotalAmount), 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(details); i++ {
		detailHtml := detailHtmlTemplate

		product := getNameProduct(details[i].Product)

		detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_unit_price$$", fmt.Sprintf("%.2f", details[i].Price), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_vat$$", fmt.Sprintf("%.2f", details[i].VatPercent), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_total$$", fmt.Sprintf("%.2f", details[i].TotalAmount), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}

func reportPurchaseOrder(id int, forcePrint bool) []byte {
	s := getPurchaseOrderRow(int32(id))

	paymentMethod := getNamePaymentMethod(s.PaymentMethod)
	supplier := getNameSupplier(s.Supplier)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State)
	}
	countryName := getNameCountry(address.Country)
	details := getSalesOrderDetail(s.Id)

	content, err := ioutil.ReadFile("./reports/purchase_order.html")
	if err != nil {
		return nil
	}

	html := string(content)

	html = strings.Replace(html, "$$order_number$$", s.OrderName, 1)
	html = strings.Replace(html, "$$order_date$$", s.DateCreated.Format("2006-01-02 15:04:05"), 1)
	html = strings.Replace(html, "$$order_reference$$", s.SupplierReference, 1)
	html = strings.Replace(html, "$$order_payment_method_name$$", paymentMethod, 1)
	html = strings.Replace(html, "$$order_customer_name$$", supplier, 1)
	html = strings.Replace(html, "$$address_address$$", address.Address, 1)
	html = strings.Replace(html, "$$address_address2$$", address.Address2, 1)
	html = strings.Replace(html, "$$address_city$$", address.City, 1)
	html = strings.Replace(html, "$$address_postcode$$", address.ZipCode, 1)
	html = strings.Replace(html, "$$address_state$$", stateName, 1)
	html = strings.Replace(html, "$$address_country$$", countryName, 1)
	html = strings.Replace(html, "$$order_notes$$", s.Notes, 1)
	html = strings.Replace(html, "$$order_total_products$$", fmt.Sprintf("%.2f", s.TotalProducts), 1)
	html = strings.Replace(html, "$$order_vat_amount$$", fmt.Sprintf("%.2f", s.VatAmount), 1)
	html = strings.Replace(html, "$$order_discount_percent$$", fmt.Sprintf("%.2f", s.DiscountPercent), 1)
	html = strings.Replace(html, "$$order_fix_discount$$", fmt.Sprintf("%.2f", s.FixDiscount), 1)
	html = strings.Replace(html, "$$order_shipping_price$$", fmt.Sprintf("%.2f", s.ShippingPrice), 1)
	html = strings.Replace(html, "$$order_shipping_discount$$", fmt.Sprintf("%.2f", s.ShippingDiscount), 1)
	html = strings.Replace(html, "$$order_total_with_discount$$", fmt.Sprintf("%.2f", s.TotalWithDiscount), 1)
	html = strings.Replace(html, "$$order_total_amount$$", fmt.Sprintf("%.2f", s.TotalAmount), 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(details); i++ {
		detailHtml := detailHtmlTemplate

		product := getNameProduct(details[i].Product)

		detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_unit_price$$", fmt.Sprintf("%.2f", details[i].Price), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_vat$$", fmt.Sprintf("%.2f", details[i].VatPercent), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_total$$", fmt.Sprintf("%.2f", details[i].TotalAmount), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}
