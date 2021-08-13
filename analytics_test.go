package main

import "testing"

func TestMonthlySalesAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := monthlySalesAmount(nil)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan monthly sales amount")
		return
	}
}

func TestMonthlySalesQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := monthlySalesQuantity(nil)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan monthly sales quantity")
		return
	}
}

func TestSalesOfAProductQuantity(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := salesOfAProductQuantity(1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan sales of a product")
		return
	}
}

func TestSalesOfAProductAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := salesOfAProductAmount(1)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan sales of a product")
		return
	}
}

func TestDaysOfServiceSaleOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := daysOfServiceSaleOrders(nil)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan days of service")
		return
	}
}

func TestDaysOfServicePurchaseOrders(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := daysOfServicePurchaseOrders(nil)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan days of service")
		return
	}
}

func TestPurchaseOrdersByMonthAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := purchaseOrdersByMonthAmount(nil)
	if len(m) == 0 || m[0].Year <= 0 {
		t.Error("Can't scan purchase orders by month")
		return
	}
}

func TestPaymentMethodsSaleOrdersAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := paymentMethodsSaleOrdersAmount(nil)
	if len(m) == 0 || m[0].PaymentMethod <= 0 {
		t.Error("Can't scan amount by payment method")
		return
	}
}

func TestCountriesSaleOrdersAmount(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	q := CountriesSaleOrdersQuery{}
	m := q.countriesSaleOrdersAmount()
	if len(m) == 0 || m[0].Country <= 0 {
		t.Error("Can't scan amount of sales by country")
		return
	}
}

func TestManufacturingOrderCreatedManufacturedDaily(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := manufacturingOrderCreatedManufacturedDaily()
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

	m := dailyShippingQuantity()
	if len(m) == 0 || m[0].Quantity <= 0 {
		t.Error("Can't scan quantity shipping")
		return
	}
}

func TestShippingByCarriers(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := shippingByCarriers()
	if len(m) == 0 || m[0].Carrier <= 0 {
		t.Error("Can't scan shipping by carrier")
		return
	}
}
