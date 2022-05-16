package main

import "fmt"

func addORMModels() bool {
	// GORM does not do this automatically
	dbOrm.Exec("UPDATE pg_opclass SET opcdefault = true WHERE opcname='gin_trgm_ops'")

	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_width")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_height")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_size")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_margin_top")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_margin_bottom")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_margin_left")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN IF EXISTS product_barcode_label_margin_right")

	err := dbOrm.AutoMigrate(&Country{}, &Language{}, &Currency{}, &Warehouse{}, &Settings{}, &SettingsEcommerce{}, &BillingSerie{}, &PaymentMethod{}, &Account{}, &Journal{},
		&User{}, &Group{}, &UserGroup{}, &PermissionDictionary{}, &PermissionDictionaryGroup{}, &DocumentContainer{}, &Document{}, &ProductFamily{}, &ProductAccount{},
		&Carrier{}, &State{}, &Color{}, &Packages{}, &Incoterm{}, &Product{}, &ConfigAccountsVat{}, &ManufacturingOrderType{}, &ManufacturingOrderTypeComponents{},
		&Supplier{}, &Address{}, &Customer{}, &AccountingMovement{}, &AccountingMovementDetail{}, &CollectionOperation{}, &PaymentTransaction{}, &Charges{}, &Payment{},
		&ApiKey{}, &ConnectionLog{}, &ConnectionFilter{}, &ConnectionFilterUser{}, &EmailLog{}, &CustomFields{}, &HSCode{}, &LoginToken{}, &POSTerminal{}, &WebHookSettings{},
		&WebHookLog{}, &WebHookRequest{}, &ReportTemplate{}, &ReportTemplateTranslation{}, &TransactionalLog{}, &SaleOrder{}, &SalesOrderDetail{}, &SalesInvoice{}, &SalesInvoiceDetail{},
		&SalesDeliveryNote{}, &PurchaseOrder{}, &PurchaseOrderDetail{}, &PurchaseInvoice{}, &PurchaseInvoiceDetail{}, &PurchaseDeliveryNote{}, &WarehouseMovement{}, &Stock{}, &Shipping{},
		&TransferBetweenWarehouses{}, &TransferBetweenWarehousesDetail{}, &Inventory{}, &InventoryProducts{}, &ManufacturingOrder{}, &ComplexManufacturingOrder{},
		&ComplexManufacturingOrderManufacturingOrder{}, &EnterpriseLogo{}, &Pallet{}, &Packaging{}, &SalesOrderDetailPackaged{}, &SalesOrderDetailDigitalProductData{}, &SalesOrderDiscount{},
		&ShippingStatusHistory{}, &ShippingTag{}, &ProductImage{}, &PwdBlacklist{}, &PwdSHA1Blacklist{}, &PSAddress{}, &PSCarrier{}, &PSCountry{}, &PSCurrency{}, &PSCustomer{}, &PSLanguage{},
		&PSOrder{}, &PSOrderDetail{}, &PSProduct{}, &PSProductCombination{}, &PSProductOptionValue{}, &PSState{}, &PSZone{}, &SYAddress{}, &SYCustomer{}, &SYDraftOrderLineItem{}, &SYDraftOrder{},
		&SYOrderLineItem{}, &SYOrder{}, &SYProduct{}, &SYVariant{}, &WCCustomer{}, &WCOrderDetail{}, &WCOrder{}, &WCProductVariation{}, &WCProduct{}, &LabelPrinterProfile{}) // 111
	if err != nil {
		fmt.Println("AutoMigrate", err)
		log("AutoMigrate", err.Error())
		//return false
	}

	/*settingsRecords := getSettingsRecords()
	for i := 0; i < len(settingsRecords); i++ {
		settings := settingsRecords[i]
		/*ecommerceConfig := SettingsEcommerce{
			EnterpriseId:                      settings.Id,
			Ecommerce:                         settings.Ecommerce,
			PrestaShopUrl:                     settings.PrestaShopUrl,
			PrestaShopApiKey:                  settings.PrestaShopApiKey,
			PrestaShopLanguageId:              settings.PrestaShopLanguageId,
			PrestaShopExportSerieId:           settings.PrestaShopExportSerieId,
			PrestaShopIntracommunitySerieId:   settings.PrestaShopIntracommunitySerieId,
			PrestaShopInteriorSerieId:         settings.PrestaShopInteriorSerieId,
			PrestaShopStatusPaymentAccepted:   settings.PrestaShopStatusPaymentAccepted,
			PrestaShopStatusShipped:           settings.PrestaShopStatusShipped,
			WooCommerceUrl:                    settings.WooCommerceUrl,
			WooCommerceConsumerKey:            settings.WooCommerceConsumerKey,
			WooCommerceConsumerSecret:         settings.WooCommerceConsumerSecret,
			WooCommerceExportSerieId:          settings.WooCommerceExportSerieId,
			WooCommerceIntracommunitySerieId:  settings.WooCommerceIntracommunitySerieId,
			WooCommerceInteriorSerieId:        settings.WooCommerceInteriorSerieId,
			WooCommerceDefaultPaymentMethodId: settings.WooCommerceDefaultPaymentMethodId,
			ShopifyUrl:                        settings.ShopifyUrl,
			ShopifyToken:                      settings.ShopifyToken,
			ShopifyExportSerieId:              settings.ShopifyExportSerieId,
			ShopifyIntracommunitySerieId:      settings.ShopifyIntracommunitySerieId,
			ShopifyInteriorSerieId:            settings.ShopifyInteriorSerieId,
			ShopifyDefaultPaymentMethodId:     settings.ShopifyDefaultPaymentMethodId,
			ShopifyShopLocationId:             settings.ShopifyShopLocationId,
		}
		dbOrm.Save(&ecommerceConfig)*/
	/*emailConfig := SettingsEmail{
		EnterpriseId:            settings.Id,
		Email:                   settings.Email,
		SendGridKey:             settings.SendGridKey,
		EmailFrom:               settings.EmailFrom,
		NameFrom:                settings.NameFrom,
		SMTPIdentity:            settings.SMTPIdentity,
		SMTPUsername:            settings.SMTPUsername,
		SMTPPassword:            settings.SMTPPassword,
		SMTPHostname:            settings.SMTPHostname,
		SMTPSTARTTLS:            settings.SMTPSTARTTLS,
		SMTPReplyTo:             settings.SMTPReplyTo,
		EmailSendErrorEcommerce: settings.EmailSendErrorEcommerce,
		EmailSendErrorSendCloud: settings.EmailSendErrorSendCloud,
	}
	dbOrm.Save(&emailConfig)*/
	//}

	return true
}
