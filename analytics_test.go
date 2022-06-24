/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "testing"

func TestMonthlySalesAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := MonthlySalesAmountQuery{}
	m := q.monthlySalesAmount(1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan monthly sales amount")
		return
	}
}

func TestMonthlySalesQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := MonthlySalesAmountQuery{}
	m := q.monthlySalesQuantity(1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan monthly sales quantity")
		return
	}
}

func TestSalesOfAProductQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := salesOfAProductQuantity([]int32{1}, 1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan sales of a product")
		return
	}
}

func TestSalesOfAProductAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := salesOfAProductAmount([]int32{1}, 1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan sales of a product")
		return
	}
}

func TestDaysOfServiceSaleOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := daysOfServiceSaleOrders(nil, 1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan days of service")
		return
	}
}

func TestDaysOfServicePurchaseOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := daysOfServicePurchaseOrders(nil, 1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan days of service")
		return
	}
}

func TestPurchaseOrdersByMonthAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PurchaseOrdersByMonthQuery{}
	m := q.purchaseOrdersByMonthAmount(1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan purchase orders by month")
		return
	}
}

func TestPaymentMethodsSaleOrdersAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := PaymentMethodsSaleOrdersQuantityQuery{}
	m := q.paymentMethodsSaleOrdersAmount(1)
	if len(m.Quantity) == 0 || m.Quantity[0].PaymentMethod <= 0 {
		t.Error("Can't scan amount by payment method")
		return
	}
}

func TestCountriesSaleOrdersAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := CountriesSaleOrdersQuery{}
	m := q.countriesSaleOrdersAmount(1)
	if len(m.Amount) == 0 || m.Amount[0].Country <= 0 {
		t.Error("Can't scan amount of sales by country")
		return
	}
}

func TestManufacturingOrderCreatedManufacturedDaily(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := ManufacturingOrderCreatedManufacturedDailyQuery{}
	m := q.manufacturingOrderCreatedManufacturedDaily(1)
	if len(m.Created) == 0 || m.Created[0].Quantity <= 0 {
		t.Error("Can't scan manufacturing orders created")
		return
	}
	if len(m.Manufactured) == 0 || m.Manufactured[0].Quantity <= 0 {
		t.Error("Can't scan manufacturing orders manufactures")
		return
	}
}

func TestDailyShippingQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := dailyShippingQuantity(1)
	if len(m) == 0 || m[0].Quantity <= 0 {
		t.Error("Can't scan quantity shipping")
		return
	}
}

func TestShippingByCarriers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := ShippingByCarriersQuery{}
	m := q.shippingByCarriers(1)
	if len(m.Quantity) == 0 || m.Quantity[0].Carrier <= 0 {
		t.Error("Can't scan shipping by carrier")
		return
	}
}
