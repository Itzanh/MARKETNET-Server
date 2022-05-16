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
		return false
	}

	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN ecommerce")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN email")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_url")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_api_key")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_language_id")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_export_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_intracommunity_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_interior_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN sendgrid_key")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN email_from")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN name_from")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_status_payment_accepted")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN prestashop_status_shipped")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_url")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_consumer_key")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_consumer_secret")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_export_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_intracommunity_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_interior_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN woocommerce_default_payment_method")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_url")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_token")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_export_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_intracommunity_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_interior_serie")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_default_payment_method")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN shopify_shop_location_id")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_identity")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_username")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_password")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_hostname")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_starttls")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN smtp_reply_to")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN email_send_error_ecommerce")
	dbOrm.Exec("ALTER TABLE public.config DROP COLUMN email_send_error_sendcloud")

	return true
}
