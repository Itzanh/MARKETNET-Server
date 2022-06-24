/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

// In the first versions, there was only a PrestaShop integration, but there will be more.
// This code chooses from PrestaShop or other e-commerce integrations set up in the settings.
// The app will call this generic piece of code, and it will call the corresponding e-commerce integration. In future e-commerce integrations, add the calls here to add more code to the app.

type ECommerce struct {
	Enterprise int32
}

func (e *ECommerce) ecommerceControllerImportFromEcommerce() {
	s := getSettingsRecordById(e.Enterprise)

	switch s.SettingsEcommerce.Ecommerce {
	case "P":
		importFromPrestaShop(e.Enterprise)
	case "W":
		importFromWooCommerce(e.Enterprise)
	case "S":
		importFromShopify(e.Enterprise)
	}
}

func ecommerceControllerUpdateTrackingNumber(salesOrderId int64, trackingNumber string, enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	switch s.SettingsEcommerce.Ecommerce {
	case "P":
		return updateTrackingNumberPrestaShopOrder(salesOrderId, trackingNumber, enterpriseId)
	case "W":
		return updateTrackingNumberWooCommerceOrder(salesOrderId, trackingNumber, enterpriseId)
	case "S":
		return updateTrackingNumberShopifyOrder(salesOrderId, trackingNumber, enterpriseId)
	}

	return false
}

func ecommerceControllerupdateStatusPaymentAccepted(salesOrderId int64, enterpriseId int32) bool {
	s := getSettingsRecordById(enterpriseId)

	switch s.SettingsEcommerce.Ecommerce {
	case "P":
		return updateStatusPaymentAcceptedPrestaShop(salesOrderId, enterpriseId)
	case "W":
		return updateStatusPaymentAcceptedWooCommerce(salesOrderId, enterpriseId)
	case "S":
		return updateStatusPaymentAcceptedShopify(salesOrderId, enterpriseId)
	}

	return false
}
