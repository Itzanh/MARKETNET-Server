package main

// In the first versions, there was only a PrestaShop integration, but there will be more.
// This code chooses from PrestaShop or other e-commerce integrations set up in the settings.
// The app will call this generic piece of code, and it will call the corresponding e-commerce integration. In future e-commerce integrations, add the calls here to add more code to the app.

func ecommerceControllerImportFromEcommerce() {
	s := getSettingsRecord()

	switch s.Ecommerce {
	case "P":
		importFromPrestaShop()
	case "W":
		importFromWooCommerce()
	}
}

func ecommerceControllerUpdateTrackingNumber(salesOrderId int32, trackingNumber string) bool {
	s := getSettingsRecord()

	switch s.Ecommerce {
	case "P":
		return updateTrackingNumberPrestaShopOrder(salesOrderId, trackingNumber)
	case "W":
		return updateTrackingNumberWooCommerceOrder(salesOrderId, trackingNumber)
	}

	return false
}

func ecommerceControllerupdateStatusPaymentAccepted(salesOrderId int32) bool {
	s := getSettingsRecord()

	switch s.Ecommerce {
	case "P":
		return updateStatusPaymentAcceptedPrestaShop(salesOrderId)
	case "W":
		return updateStatusPaymentAcceptedWooCommerce(salesOrderId)
	}

	return false
}
