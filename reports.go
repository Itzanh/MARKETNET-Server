package main

import (
	"encoding/base64"
	"fmt"
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
	idLangString, ok := r.URL.Query()["lang"]
	var idLang int
	if ok {
		idLang, err = strconv.Atoi(idLangString[0])
		if err != nil {
			return
		}
	}

	ok, enterpriseId := consumeToken(token[0])
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch report[0] {
	case "SALES_ORDER":
		w.Write(reportSalesOrder(id, forcePrint, enterpriseId, int32(idLang)))
	case "SALES_INVOICE":
		w.Write(reportSalesInvoice(id, forcePrint, enterpriseId))
	case "SALES_INVOICE_TICKET":
		w.Write(reportSalesInvoiceTicket(id, forcePrint, enterpriseId))
	case "SALES_DELIVERY_NOTE":
		w.Write(reportSalesDeliveryNote(id, forcePrint, enterpriseId))
	case "PURCHASE_ORDER":
		w.Write(reportPurchaseOrder(id, forcePrint, enterpriseId))
	case "BOX_CONTENT":
		w.Write(reportBoxContent(id, forcePrint, enterpriseId))
	case "PALLET_CONTENT":
		w.Write(reportPalletContent(id, forcePrint, enterpriseId))
	case "CARRIER_PALLET":
		w.Write(reportCarrierPallet(id, forcePrint, enterpriseId))
	}

}

func getEnterpriseLogoBase64(enterpriseId int32) string {
	logo, mimeType := getEnterpriseLogo(enterpriseId)
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(logo)
}

func reportSalesOrder(id int, forcePrint bool, enterpriseId int32, idLang int32) []byte {
	s := getSalesOrderRow(int64(id))

	paymentMethod := getNamePaymentMethod(s.PaymentMethod, enterpriseId)
	customer := getNameCustomer(s.Customer, enterpriseId)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getSalesOrderDetail(s.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "SALES_ORDER")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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

		product := getNameProduct(details[i].Product, enterpriseId)

		detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_unit_price$$", fmt.Sprintf("%.2f", details[i].Price), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_vat$$", fmt.Sprintf("%.2f", details[i].VatPercent), 1)
		detailHtml = strings.Replace(detailHtml, "$$detail_total$$", fmt.Sprintf("%.2f", details[i].TotalAmount), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	if idLang != 0 {
		html = translateReport(html, idLang, enterpriseId)
	}

	return []byte(html)
}

func reportSalesInvoice(id int, forcePrint bool, enterpriseId int32) []byte {
	i := getSalesInvoiceRow(int64(id))

	paymentMethod := getNamePaymentMethod(i.PaymentMethod, enterpriseId)
	customer := getNameCustomer(i.Customer, enterpriseId)
	address := getAddressRow(i.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getSalesInvoiceDetail(i.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "SALES_INVOICE")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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

		product := getNameProduct(*details[i].Product, enterpriseId)

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

func reportSalesInvoiceTicket(id int, forcePrint bool, enterpriseId int32) []byte {
	i := getSalesInvoiceRow(int64(id))

	paymentMethod := getNamePaymentMethod(i.PaymentMethod, enterpriseId)
	customer := getNameCustomer(i.Customer, enterpriseId)
	address := getAddressRow(i.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getSalesInvoiceDetail(i.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "SALES_INVOICE_TICKET")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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

		product := getNameProduct(*details[i].Product, enterpriseId)

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

func reportSalesDeliveryNote(id int, forcePrint bool, enterpriseId int32) []byte {
	n := getSalesDeliveryNoteRow(int64(id))

	paymentMethod := getNamePaymentMethod(n.PaymentMethod, enterpriseId)
	customer := getNameCustomer(n.Customer, enterpriseId)
	address := getAddressRow(n.ShippingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getWarehouseMovementBySalesDeliveryNote(n.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "SALES_DELIVERY_NOTE")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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
	html = strings.Replace(html, "$$note_vat_amount$$", fmt.Sprintf("%.2f", n.VatAmount), 1)
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

		product := getNameProduct(details[i].Product, enterpriseId)

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

func reportPurchaseOrder(id int, forcePrint bool, enterpriseId int32) []byte {
	s := getPurchaseOrderRow(int64(id))

	paymentMethod := getNamePaymentMethod(s.PaymentMethod, enterpriseId)
	supplier := getNameSupplier(s.Supplier, enterpriseId)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getPurchaseOrderDetail(s.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "PURCHASE_ORDER")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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

		product := getNameProduct(details[i].Product, enterpriseId)

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

func reportBoxContent(id int, forcePrint bool, enterpriseId int32) []byte {
	p := getPackagingRow(int64(id))
	_package := getPackagesRow(p.Package)

	template := getReportTemplate(enterpriseId, "BOX_CONTENT")

	html := template.Html

	html = strings.Replace(html, "$$box_name$$", _package.Name+" ("+fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth))+")", 1)
	html = strings.Replace(html, "$$box_weight$$", fmt.Sprintf("%.2f", p.Weight), 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	details := getSalesOrderDetailPackaged(p.Id, enterpriseId)

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(details); i++ {
		detailHtml := detailHtmlTemplate

		p := getProductRow(getSalesOrderDetailRow(details[i].OrderDetail).Product)

		detailHtml = strings.Replace(detailHtml, "$$product_name$$", p.Name, 1)
		detailHtml = strings.Replace(detailHtml, "$$product_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}

func reportPalletContent(id int, forcePrint bool, enterpriseId int32) []byte {
	p := getPalletsRow(int32(id))

	if p.Id <= 0 {
		return nil
	}

	template := getReportTemplate(enterpriseId, "PALLET_CONTENT")

	html := template.Html
	packaging := getPackaging(p.SalesOrder, enterpriseId)
	boxHtmlTemplate := html[strings.Index(html, "&&box&&")+len("&&box&&") : strings.Index(html, "&&--box--&&")]
	boxHtml := ""
	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]

	for i := 0; i < len(packaging); i++ {
		if *packaging[i].Pallet != p.Id {
			continue
		}
		detailHtml := boxHtmlTemplate

		_package := getPackagesRow(packaging[i].Package)

		detailHtml = strings.Replace(detailHtml, "$$box_name$$", _package.Name+" ("+fmt.Sprintf("%dx%dx%d", int(_package.Width), int(_package.Height), int(_package.Depth))+")", 1)
		detailHtml = strings.Replace(detailHtml, "$$box_weight$$", fmt.Sprintf("%.2f", packaging[i].Weight), 1)

		details := getSalesOrderDetailPackaged(packaging[i].Id, enterpriseId)

		boxDetailsHtml := ""
		for i := 0; i < len(details); i++ {
			dPackagedHtml := detailHtmlTemplate

			p := getProductRow(getSalesOrderDetailRow(details[i].OrderDetail).Product)

			dPackagedHtml = strings.Replace(dPackagedHtml, "$$product_name$$", p.Name, 1)
			dPackagedHtml = strings.Replace(dPackagedHtml, "$$product_quantity$$", strconv.Itoa(int(details[i].Quantity)), 1)

			boxDetailsHtml += dPackagedHtml
		}
		detailHtml = detailHtml[:strings.Index(detailHtml, "&&detail&&")] + boxDetailsHtml + detailHtml[strings.Index(detailHtml, "&&--detail--&&")+len("&&--detail--&&"):]

		boxHtml += detailHtml
	}

	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	html = html[:strings.Index(html, "&&box&&")] + boxHtml + html[strings.Index(html, "&&--box--&&")+len("&&--box--&&"):]

	return []byte(html)
}

func reportCarrierPallet(id int, forcePrint bool, enterpriseId int32) []byte {
	s := getSalesOrderRow(int64(id))

	customer := getNameCustomer(s.Customer, enterpriseId)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	pallets := getSalesOrderPallets(s.Id, enterpriseId).Pallets

	template := getReportTemplate(enterpriseId, "CARRIER_PALLET")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
	html = strings.Replace(html, "$$order_customer_name$$", customer, 1)
	html = strings.Replace(html, "$$address_address$$", address.Address, 1)
	html = strings.Replace(html, "$$address_address2$$", address.Address2, 1)
	html = strings.Replace(html, "$$address_city$$", address.City, 1)
	html = strings.Replace(html, "$$address_postcode$$", address.ZipCode, 1)
	html = strings.Replace(html, "$$address_state$$", stateName, 1)
	html = strings.Replace(html, "$$address_country$$", countryName, 1)
	if forcePrint {
		html = strings.Replace(html, "$$script$$", "window.print()", 1)
	} else {
		html = strings.Replace(html, "$$script$$", "", 1)
	}

	detailHtmlTemplate := html[strings.Index(html, "&&detail&&")+len("&&detail&&") : strings.Index(html, "&&--detail--&&")]
	detailsHtml := ""

	for i := 0; i < len(pallets); i++ {
		detailHtml := detailHtmlTemplate

		detailHtml = strings.Replace(detailHtml, "$$pallet_name$$", pallets[i].Name, 1)
		detailHtml = strings.Replace(detailHtml, "$$pallet_dimensions$$",
			fmt.Sprintf("%.2f x %.2f x %.2f", pallets[i].Width, pallets[i].Height, pallets[i].Depth), 1)
		detailHtml = strings.Replace(detailHtml, "$$pallet_weight$$", fmt.Sprintf("%.2f", pallets[i].Weight), 1)

		detailsHtml += detailHtml
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}

func reportSalesOrderDigitalProductDetails(id int, forcePrint bool, enterpriseId int32) []byte {
	s := getSalesOrderRow(int64(id))

	paymentMethod := getNamePaymentMethod(s.PaymentMethod, enterpriseId)
	customer := getNameCustomer(s.Customer, enterpriseId)
	address := getAddressRow(s.BillingAddress)
	stateName := ""
	if address.State != nil {
		stateName = getNameState(*address.State, enterpriseId)
	}
	countryName := getNameCountry(address.Country, enterpriseId)
	details := getSalesOrderDetail(s.Id, enterpriseId)

	template := getReportTemplate(enterpriseId, "SALES_ORDER_DIGITAL_PRODUCT_DATA")

	html := template.Html

	html = strings.Replace(html, "$$img_base64$$", getEnterpriseLogoBase64(enterpriseId), 1)
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
		digitalProductDetails := getSalesOrderDetailDigitalProductData(details[i].Id, enterpriseId)
		if len(digitalProductDetails) == 0 {
			continue
		}

		product := getNameProduct(details[i].Product, enterpriseId)

		for j := 0; j < len(digitalProductDetails); j++ {
			detailHtml := detailHtmlTemplate

			detailHtml = strings.Replace(detailHtml, "$$detail_product$$", product, 1)
			detailHtml = strings.Replace(detailHtml, "$$detail_key$$", digitalProductDetails[j].Key, 1)
			detailHtml = strings.Replace(detailHtml, "$$detail_value$$", digitalProductDetails[j].Key, 1)

			detailsHtml += detailHtml
		}
	}

	html = html[:strings.Index(html, "&&detail&&")] + detailsHtml + html[strings.Index(html, "&&--detail--&&")+len("&&--detail--&&"):]

	return []byte(html)
}
