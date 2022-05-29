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
		&TransferBetweenWarehousesMinimumStock{}) // 112
	if err != nil {
		fmt.Println("AutoMigrate", err)
		log("AutoMigrate", err.Error())
		return false
	}

	sqlStatement := `SELECT id,warehouse FROM public.sales_order ORDER BY id ASC`
	rows, _ := dbOrm.Raw(sqlStatement).Rows()
	for rows.Next() {
		var id int
		var warehouseId string
		rows.Scan(&id, &warehouseId)

		sqlStatement = `UPDATE public.sales_order_detail SET warehouse = $2 WHERE "order" = $1 AND warehouse IS NULL`
		db.Exec(sqlStatement, id, warehouseId)
	}

	sqlStatement = `SELECT id,warehouse FROM public.purchase_order ORDER BY id ASC`
	rows, _ = dbOrm.Raw(sqlStatement).Rows()
	for rows.Next() {
		var id int
		var warehouseId string
		rows.Scan(&id, &warehouseId)

		sqlStatement = `UPDATE public.purchase_order_detail SET warehouse = $2 WHERE "order" = $1 AND warehouse IS NULL`
		db.Exec(sqlStatement, id, warehouseId)
	}

	return true
}
