package main

import "fmt"

func addORMModels() bool {
	// GORM does not do this automatically
	dbOrm.Exec("UPDATE pg_opclass SET opcdefault = true WHERE opcname='gin_trgm_ops'")

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
		&SYOrderLineItem{}, &SYOrder{}, &SYProduct{}, &SYVariant{}, &WCCustomer{}, &WCOrderDetail{}, &WCOrder{}, &WCProductVariation{}, &WCProduct{}, &LabelPrinterProfile{},
		&TransferBetweenWarehousesMinimumStock{}, &ProductIncludedProduct{}, &ProductIncludedProductSalesOrderDetail{}, &SettingsCleanUp{}) // 115
	if err != nil {
		fmt.Println("AutoMigrate", err)
		log("AutoMigrate", err.Error())
		return false
	}

	// create default rows for Settings Clean Up
	settingsRecords := getSettingsRecords()
	for _, settingsRecord := range settingsRecords {
		settingsCleanUp := SettingsCleanUp{
			CronCleanTransactionalLog: "@monthly",
			TransactionalLogDays:      90,
			CronCleanConnectionLog:    "@weekly",
			ConnectionLogDays:         30,
			CronCleanLoginToken:       "@weekly",
			EnterpriseId:              settingsRecord.Id,
		}
		dbOrm.Create(&settingsCleanUp)
	}

	return true
}
