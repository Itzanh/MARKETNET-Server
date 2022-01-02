package main

import (
	"io/ioutil"
	"math/rand"
	"strings"
)

const (
	ALFA   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUMBER = "0123456789"
)

// This file is for developmenmt only.

func generateDemoData(enterpriseId int32) {
	generateCustomers(enterpriseId)
	generateAddresses(enterpriseId)
	generateSaleOrders(enterpriseId)
	generateInvoiceAllSalesOrders()
}

func generateCustomers(enterpriseId int32) {
	INT := "INT"
	EXP := "EXP"
	IEU := "IEU"

	countries := getCountries(enterpriseId)
	languages := getLanguages(enterpriseId)
	paymentMethods := getPaymentMethods(enterpriseId)

	content, err := ioutil.ReadFile("customer_names.txt")
	if err != nil {
		return
	}

	names := strings.Split(string(content), "\n")
	for i := 0; i < len(names); i++ {
		name := names[i]
		c := Customer{}
		// name
		c.Name = name
		c.Tradename = name
		c.FiscalName = name

		// tax id
		if rand.Intn(100) <= 15 {
			c.TaxId = ""
			c.TaxId += string(NUMBER[rand.Intn(len(NUMBER))])
			for j := 0; j < 8; j++ {
				c.TaxId += string(ALFA[rand.Intn(len(ALFA))])
			}
		}

		// country
		country := countries[rand.Intn(len(countries))]
		c.Country = &country.Id

		// state
		states := getStatesByCountry(*c.Country, 1)
		if len(states) > 0 {
			c.State = &states[rand.Intn(len(states))].Id
		}

		// vat number
		if country.Zone == "U" && rand.Intn(100) <= 10 {
			c.VatNumber = ""
			c.VatNumber += string(NUMBER[rand.Intn(len(NUMBER))])
			for j := 0; j < 8; j++ {
				c.VatNumber += string(ALFA[rand.Intn(len(ALFA))])
			}
		}

		// phone
		for i := 0; i < 9; i++ {
			c.Phone += string(NUMBER[rand.Intn(len(NUMBER))])
		}

		// email
		c.Email = strings.ToLower(strings.ReplaceAll(name, " ", "."))
		if rand.Intn(100) <= 50 {
			c.Email += "@gmail.com"
		} else {
			c.Email += "@hotmail.com"
		}

		// billing serie
		if country.Zone == "N" || country.Zone == "U" && c.VatNumber == "" {
			c.BillingSeries = &INT
		} else if country.Zone == "E" {
			c.BillingSeries = &EXP
		} else {
			c.BillingSeries = &IEU
		}

		// language
		c.Language = &languages[rand.Intn(len(languages))].Id

		// paymet method
		c.PaymentMethod = &paymentMethods[rand.Intn(len(paymentMethods))].Id

		c.enterprise = 1
		c.insertCustomer(0)
	}
}

func generateAddresses(enterpriseId int32) {
	countries := getCountries(1)
	q := PaginationQuery{Offset: 0, Limit: 0, enterprise: 1}
	customers := q.getCustomers()

	for i := 0; i < len(customers.Customers); i++ {
		customer := customers.Customers[i]

		addressesGenerate := rand.Intn(4) + 1
		for l := 0; l < addressesGenerate; l++ {
			a := Address{}
			a.Customer = &customer.Id

			// address 2
			words := rand.Intn(4) + 1
			for j := 0; j < words; j++ {
				characters := rand.Intn(24) + 3
				for k := 0; k < characters; k++ {
					a.Address += string(ALFA[rand.Intn(len(ALFA))])
				}
				a.Address += " "
			}

			// address 2
			if rand.Intn(100) >= 50 {
				words := rand.Intn(4) + 1
				for j := 0; j < words; j++ {
					characters := rand.Intn(24) + 3
					for k := 0; k < characters; k++ {
						a.Address2 += string(ALFA[rand.Intn(len(ALFA))])
					}
					a.Address2 += " "
				}
			}

			// country
			country := countries[rand.Intn(len(countries))]
			a.Country = country.Id

			// state
			states := getStatesByCountry(a.Country, 1)
			if len(states) > 0 {
				a.State = &states[rand.Intn(len(states))].Id
			}

			// city
			characters := rand.Intn(24) + 3
			for k := 0; k < characters; k++ {
				a.City += string(ALFA[rand.Intn(len(ALFA))])
			}

			// zip code
			if rand.Intn(100) >= 50 {
				for k := 0; k < 5; k++ {
					a.ZipCode += string(NUMBER[rand.Intn(len(NUMBER))])
				}
			} else {
				if rand.Intn(100) >= 50 {
					for k := 0; k < 2; k++ {
						a.ZipCode += string(ALFA[rand.Intn(len(ALFA))])
					}
					for k := 0; k < 3; k++ {
						a.ZipCode += string(NUMBER[rand.Intn(len(NUMBER))])
					}
				} else {
					for k := 0; k < 5; k++ {
						if rand.Intn(100) >= 40 {
							a.ZipCode += string(NUMBER[rand.Intn(len(NUMBER))])
						} else {
							a.ZipCode += string(ALFA[rand.Intn(len(ALFA))])
						}
					}
				}
			}

			a.PrivateOrBusiness = "_"
			a.enterprise = enterpriseId
			a.insertAddress(0)
		}
	}

}

func generateSaleOrders(enterpriseId int32) {
	currencies := getCurrencies(enterpriseId)
	products := getProduct(1)
	q := PaginationQuery{Offset: 0, Limit: 0}
	customers := q.getCustomers()

	for i := 0; i < len(customers.Customers); i++ {
		customer := customers.Customers[i]
		addresses := getCustomerAddresses(customer.Id, 1)

		if customer.PaymentMethod == nil || customer.BillingSeries == nil {
			continue
		}

		ordersGenerate := rand.Intn(10) + 1
		for l := 0; l < ordersGenerate; l++ {
			o := SaleOrder{}
			o.Warehouse = "W1"
			o.Customer = customer.Id
			o.PaymentMethod = *customer.PaymentMethod
			o.BillingSeries = *customer.BillingSeries
			o.Currency = currencies[rand.Intn(len(currencies))].Id
			o.BillingAddress = addresses[rand.Intn(len(addresses))].Id
			o.ShippingAddress = addresses[rand.Intn(len(addresses))].Id
			o.enterprise = 1
			ok, id := o.insertSalesOrder(1)

			if !ok {
				continue
			}

			detailsGenerate := rand.Intn(15) + 1
			for l := 0; l < detailsGenerate; l++ {
				product := products[rand.Intn(len(products))]

				d := SalesOrderDetail{}
				d.Order = id
				d.Product = product.Id
				d.Price = product.Price
				d.Quantity = int32(rand.Intn(10) + 1)
				d.VatPercent = product.VatPercent
				d.enterprise = 1
				d.insertSalesOrderDetail(0)
			}
		}
	}
}

func generateInvoiceAllSalesOrders() {
	q := PaginationQuery{Offset: 0, Limit: 100000}
	o := q.getSalesOrder(1)

	for i := 0; i < len(o.Orders); i++ {
		invoiceAllSaleOrder(o.Orders[i].Id, 1, 0)
		deliveryNoteAllSaleOrder(o.Orders[i].Id, 1, 0, nil)
	}
}
