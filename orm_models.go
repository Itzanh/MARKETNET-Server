package main

import "fmt"

func addORMModels() bool {
	// GORM does not do this automatically
	dbOrm.Exec("UPDATE pg_opclass SET opcdefault = true WHERE opcname='gin_trgm_ops'")

	dbOrm.Exec("ALTER TABLE public.sales_order_detail_digital_product_data DROP CONSTRAINT IF EXISTS fk_sales_order_detail_digital_product_data_detail")

	dbOrm.Exec("UPDATE connection_filter SET ip_address = '::1' WHERE ip_address = '['")
	dbOrm.Exec("ALTER TABLE connection_filter ALTER COLUMN ip_address type inet USING CAST(ip_address AS inet)")

	dbOrm.Exec("UPDATE connection_log SET ip_address = '::1' WHERE ip_address = '['")
	dbOrm.Exec("ALTER TABLE connection_log ALTER COLUMN ip_address type inet USING CAST(ip_address AS inet)")

	dbOrm.Exec("UPDATE login_tokens SET ip_address = '::1' WHERE ip_address = '['")
	dbOrm.Exec("ALTER TABLE login_tokens ALTER COLUMN ip_address type inet USING CAST(ip_address AS inet)")

	err := dbOrm.AutoMigrate(&Country{}, &Language{}, &Currency{}, &Warehouse{}, &Settings{}, &BillingSerie{}, &PaymentMethod{}, &Account{}, &Journal{},
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
		&SYOrderLineItem{}, &SYOrder{}, &SYProduct{}, &SYVariant{}, &WCCustomer{}, &WCOrderDetail{}, &WCOrder{}, &WCProductVariation{}, &WCProduct{}) // 109
	if err != nil {
		fmt.Println(err)
		log("AutoMigrate", err.Error())
		return false
	}

	// add account name field
	settingsRecords := getSettingsRecords()
	for i := 0; i < len(settingsRecords); i++ {
		accounts := getAccounts(settingsRecords[i].Id)
		for j := 0; j < len(accounts); j++ {
			account := accounts[j]
			account.updateAccount()
		}
	}

	// add enterprise in shipping status history
	shippingStatusHistoryRecords := make([]ShippingStatusHistory, 0)
	dbOrm.Model(&ShippingStatusHistory{}).Preload("Shipping").Find(&shippingStatusHistoryRecords)
	for i := 0; i < len(shippingStatusHistoryRecords); i++ {
		shippingStatusHistory := shippingStatusHistoryRecords[i]
		shippingStatusHistory.EnterpriseId = shippingStatusHistory.Shipping.EnterpriseId
		dbOrm.Save(&shippingStatusHistory)
	}

	// add enterprise in sales order detail digital product data
	salesOrderDetailDigitalProductDataRecords := make([]SalesOrderDetailDigitalProductData, 0)
	dbOrm.Model(&SalesOrderDetailDigitalProductData{}).Preload("Detail").Find(&salesOrderDetailDigitalProductDataRecords)
	for i := 0; i < len(salesOrderDetailDigitalProductDataRecords); i++ {
		salesOrderDetailDigitalProductData := salesOrderDetailDigitalProductDataRecords[i]
		salesOrderDetail := getSalesOrderDetailRow(salesOrderDetailDigitalProductData.DetailId)
		salesOrderDetailDigitalProductData.EnterpriseId = salesOrderDetail.EnterpriseId
		dbOrm.Save(&salesOrderDetailDigitalProductData)
	}

	return true
}
