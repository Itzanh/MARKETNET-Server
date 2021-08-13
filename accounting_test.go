package main

import (
	"math"
	"testing"
)

// ===== JOURNAL

func TestGetJournals(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	j := getJournals()

	for i := 0; i < len(j); i++ {
		if j[i].Id <= 0 {
			t.Error("Scan error, journals with ID 0.")
			return
		}
	}
}

func TestJournalInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	j := Journal{
		Id:   123,
		Name: "Test",
		Type: "G",
	}

	ok := j.insertJournal()
	if !ok {
		t.Error("Insert error, can't insert journal")
		return
	}

	j.Name = "Test test"
	ok = j.updateJournal()
	if !ok {
		t.Error("Update error, can't update journal")
		return
	}

	ok = j.deleteJournal()
	if !ok {
		t.Error("Delete error, can't delete journal")
		return
	}
}

// ===== ACCOUNT

func TestGetAccounts(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	a := getAccounts()

	for i := 0; i < len(a); i++ {
		if a[i].Id <= 0 {
			t.Error("Scan error, accounts with ID 0.")
			return
		}
	}
}

func TestSearchAccounts(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := AccountSearch{Journal: 430, Search: ""}
	a := s.searchAccounts()

	for i := 0; i < len(a); i++ {
		if a[i].Id <= 0 {
			t.Error("Scan error, accounts with ID 0.")
			return
		}
	}
}

func TestGetAccountRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	a := getAccountRow(1)

	if a.Id <= 0 {
		t.Error("Scan error, account with ID 0.")
		return
	}
}

func TestAccountInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	j := Journal{
		Id:   123,
		Name: "Test",
		Type: "G",
	}
	j.insertJournal()

	a := Account{
		Journal: 123,
		Name:    "Test",
	}
	ok := a.insertAccount()
	if !ok {
		t.Error("Insert error, can't insert account")
		return
	}

	getAccountId := getAccountIdByAccountNumber(j.Id, a.AccountNumber)
	if getAccountId != a.Id {
		t.Error("Can't get account id by account number")
		return
	}

	a.Name = "Test test"
	ok = a.updateAccount()
	if !ok {
		t.Error("Update error, can't update account")
		return
	}

	ok = a.deleteAccount()
	if !ok {
		t.Error("Delete error, can't delete account")
		return
	}

	j.deleteJournal()
}

func TestLocateAccountForCustomer(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	a := locateAccountForCustomer()

	for i := 0; i < len(a); i++ {
		if a[i].Id <= 0 {
			t.Error("Scan error, accounts with ID 0.")
			return
		}
	}
}

func TestLocateAccountForSupplier(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	a := locateAccountForSupplier()

	for i := 0; i < len(a); i++ {
		if a[i].Id <= 0 {
			t.Error("Scan error, accounts with ID 0.")
			return
		}
	}
}

func TestLocateAccountForBanks(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	a := locateAccountForBanks()

	for i := 0; i < len(a); i++ {
		if a[i].Id <= 0 {
			t.Error("Scan error, accounts with ID 0.")
			return
		}
	}
}

// ===== ACCOUNTING MOVEMENT

func TestGetAccountingMovement(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := getAccountingMovement()

	for i := 0; i < len(m); i++ {
		if m[i].Id <= 0 {
			t.Error("Scan error, movements with ID 0.")
			return
		}
	}
}

func TestSearchAccountingMovements(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	m := searchAccountingMovements("")

	for i := 0; i < len(m); i++ {
		if m[i].Id <= 0 {
			t.Error("Scan error, movements with ID 0.")
			return
		}
	}
}

func TestGetAccountingMovementRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	movements := getAccountingMovement()
	m := getAccountingMovementRow(movements[0].Id)

	if m.Id <= 0 {
		t.Error("Scan error, movements with ID 0.")
		return
	}
}

func TestAccountingMovementInsertDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	am := AccountingMovement{
		Type:         "N",
		BillingSerie: "INT",
	}

	ok := am.insertAccountingMovement()
	if !ok {
		t.Error("Insert error, can't insert accounting movement")
		return
	}

	ok = am.deleteAccountingMovement()
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
		return
	}
}

// ===== COLLECTION OPERATION

// ===== PAYMENT OPERATION

// ===== POST SALE INVOICES

func TestSalesPostInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert invoice and details
	i := SalesInvoice{
		Customer:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 1,
	}
	_, invoiceId := i.insertSalesInvoice()
	d := SalesInvoiceDetail{
		Invoice:    invoiceId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
	}
	d.insertSalesInvoiceDetail(true)
	d = SalesInvoiceDetail{
		Invoice:    invoiceId,
		Product:    1,
		Price:      4.99,
		Quantity:   1,
		VatPercent: 21,
	}
	d.insertSalesInvoiceDetail(true)
	i = getSalesInvoiceRow(invoiceId)

	// post sale invoice
	result := salesPostInvoices([]int32{invoiceId})
	for i := 0; i < len(result); i++ {
		if !result[i].Ok {
			t.Error("Can't post sale invoice")
			return
		}
	}

	// check accounting movement
	movements := getAccountingMovement()
	am := movements[0]

	saleInvoices := getAccountingMovementSaleInvoices(am.Id)
	if len(saleInvoices) == 0 || saleInvoices[0].Id != invoiceId {
		t.Error("Can't scan sale invocies")
		return
	}

	movementDetails := getAccountingMovementDetail(am.Id)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the customer (debit)
	invoiceTotalAmount := i.TotalAmount
	var isCustomerAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 430 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the customer was not created successfully")
				return
			}
			isCustomerAccountPresent = true
			break
		}
	}
	if !isCustomerAccountPresent {
		t.Error("Can't find detail for the customer")
		return
	}
	// check detail for the sale (credit)
	invoiceTotalProducts := i.TotalProducts
	var isSalesAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 700 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != float32(math.Round(float64(invoiceTotalProducts)*100)/100) {
				t.Error("The detail for the sales was not created successfully")
				return
			}
			isSalesAccountPresent = true
			break
		}
	}
	if !isSalesAccountPresent {
		t.Error("Can't find detail for the sales")
		return
	}
	// check detail for the vat (credit)
	invoiceTotalVat := i.TotalAmount - i.TotalProducts
	var isVatAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 477 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != float32(math.Round(float64(invoiceTotalVat)*100)/100) {
				t.Error("The detail for the vat was not created successfully")
				return
			}
			isVatAccountPresent = true
			break
		}
	}
	if !isVatAccountPresent {
		t.Error("Can't find detail for the vat")
		return
	}

	// check charges
	collectionOperations := getColletionOperations(am.Id)
	if len(collectionOperations) == 0 || collectionOperations[0].Id <= 0 {
		t.Error("Charges not created, or can't scan charges")
		return
	}
	if collectionOperations[0].Paid > 0 || collectionOperations[0].Status != "P" {
		t.Error("Paid in advance")
		return
	}

	// add charges
	charge := Charges{
		CollectionOperation: collectionOperations[0].Id,
		Amount:              collectionOperations[0].Pending,
		Concept:             "TESTING...",
	}
	ok := charge.insertCharges()
	if !ok {
		t.Error("Can't insert charges")
		return
	}

	// test charges scan
	charges := getCharges(collectionOperations[0].Id)
	if len(charges) == 0 || charges[0].Id <= 0 {
		t.Error("Can't scan charges")
		return
	}
	charge = getChargesRow(charges[0].Id)
	if charge.Id <= 0 {
		t.Error("Can't scan charge row")
		return
	}

	// check that the collection operaton has been updated
	collectionOperations = getColletionOperations(am.Id)
	if collectionOperations[0].Paid != collectionOperations[0].Total || collectionOperations[0].Status != "C" {
		t.Error("Collection operation not updated")
		return
	}

	// a new accounting movement must have been generated
	newAccountingMovement := getAccountingMovementRow(am.Id + 1)
	if newAccountingMovement.Id <= 0 {
		t.Error("New accounting movement not generated")
		return
	}

	newMovementDetails := getAccountingMovementDetail(newAccountingMovement.Id)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the customer (credit)
	isCustomerAccountPresent = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].Journal == 430 {
			if newMovementDetails[i].Debit != 0 || newMovementDetails[i].Credit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the customer was not created successfully")
				return
			}
			isCustomerAccountPresent = true
			break
		}
	}
	if !isCustomerAccountPresent {
		t.Error("Can't find detail for the customer")
		return
	}
	// check detail for the bank (debit)
	var isBankAccountPresent bool = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].Journal == 570 {
			if newMovementDetails[i].Credit != 0 || newMovementDetails[i].Debit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the bank was not created successfully")
				return
			}
			isBankAccountPresent = true
			break
		}
	}
	if !isBankAccountPresent {
		t.Error("Can't find detail for the bank")
		return
	}

	// DELETE
	// delete the charge (that will delete the new accounting movement)
	ok = charge.deleteCharges()
	if !ok {
		t.Error("Delete error, can't delete charge")
		return
	}

	// delete accounting movement
	ok = am.deleteAccountingMovement()
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
		return
	}

	// delete invoioce and details
	ok = i.deleteSalesInvoice()
	if !ok {
		t.Error("Delete error, can't delete sale invoice")
		return
	}
}

// ===== POST PURCHASE INVOICES

func TestPurchasePostInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert invoice and details
	i := PurchaseInvoice{
		Supplier:       1,
		PaymentMethod:  1,
		BillingSeries:  "INT",
		Currency:       1,
		BillingAddress: 1,
	}
	_, invoiceId := i.insertPurchaseInvoice()
	d := PurchaseInvoiceDetail{
		Invoice:    invoiceId,
		Product:    4,
		Price:      9.99,
		Quantity:   2,
		VatPercent: 21,
	}
	d.insertPurchaseInvoiceDetail(true)
	d = PurchaseInvoiceDetail{
		Invoice:    invoiceId,
		Product:    1,
		Price:      4.99,
		Quantity:   1,
		VatPercent: 21,
	}
	d.insertPurchaseInvoiceDetail(true)
	i = getPurchaseInvoiceRow(invoiceId)

	// post purchase invoice
	result := purchasePostInvoices([]int32{invoiceId})
	for i := 0; i < len(result); i++ {
		if !result[i].Ok {
			t.Error("Can't post purchase invoice")
			return
		}
	}

	// check accounting movement
	movements := getAccountingMovement()
	am := movements[0]

	purchaseInvoices := getAccountingMovementPurchaseInvoices(am.Id)
	if len(purchaseInvoices) == 0 || purchaseInvoices[0].Id != invoiceId {
		t.Error("Can't scan purchase invocies")
		return
	}

	movementDetails := getAccountingMovementDetail(am.Id)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the supplier (credit)
	invoiceTotalAmount := i.TotalAmount
	var isSupplierAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 400 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the supplier was not created successfully")
				return
			}
			isSupplierAccountPresent = true
			break
		}
	}
	if !isSupplierAccountPresent {
		t.Error("Can't find detail for the supplier")
		return
	}
	// check detail for the purchase (debit)
	invoiceTotalProducts := i.TotalProducts
	var isPurchaseAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 600 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != float32(math.Round(float64(invoiceTotalProducts)*100)/100) {
				t.Error("The detail for the purchase was not created successfully")
				return
			}
			isPurchaseAccountPresent = true
			break
		}
	}
	if !isPurchaseAccountPresent {
		t.Error("Can't find detail for the purchase")
		return
	}
	// check detail for the vat (debit)
	invoiceTotalVat := i.TotalAmount - i.TotalProducts
	var isVatAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].Journal == 472 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != float32(math.Round(float64(invoiceTotalVat)*100)/100) {
				t.Error("The detail for the vat was not created successfully")
				return
			}
			isVatAccountPresent = true
			break
		}
	}
	if !isVatAccountPresent {
		t.Error("Can't find detail for the vat")
		return
	}

	// check payments
	paymentTransactions := getPaymentTransactions(am.Id)
	if len(paymentTransactions) == 0 || paymentTransactions[0].Id <= 0 {
		t.Error("Charges not created, or can't scan payments")
		return
	}
	if paymentTransactions[0].Paid > 0 || paymentTransactions[0].Status != "P" {
		t.Error("Paid in advance")
		return
	}

	// add payments
	payment := Payment{
		PaymentTransaction: paymentTransactions[0].Id,
		Amount:             paymentTransactions[0].Pending,
		Concept:            "TESTING...",
	}
	ok := payment.insertPayment()
	if !ok {
		t.Error("Can't insert payments")
		return
	}

	// test payments scan
	payments := getPayments(paymentTransactions[0].Id)
	if len(payments) == 0 || payments[0].Id <= 0 {
		t.Error("Can't scan payments")
		return
	}
	payment = getPaymentsRow(payments[0].Id)
	if payment.Id <= 0 {
		t.Error("Can't scan payment row")
		return
	}

	// check that the payment transaction has been updated
	paymentTransactions = getPaymentTransactions(am.Id)
	if paymentTransactions[0].Paid != paymentTransactions[0].Total || paymentTransactions[0].Status != "C" {
		t.Error("Payment transaction not updated")
		return
	}

	// a new accounting movement must have been generated
	newAccountingMovement := getAccountingMovementRow(am.Id + 1)
	if newAccountingMovement.Id <= 0 {
		t.Error("New accounting movement not generated")
		return
	}

	newMovementDetails := getAccountingMovementDetail(newAccountingMovement.Id)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the supplier (debit)
	isSupplierAccountPresent = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].Journal == 400 {
			if newMovementDetails[i].Credit != 0 || newMovementDetails[i].Debit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the supplier was not created successfully")
				return
			}
			isSupplierAccountPresent = true
			break
		}
	}
	if !isSupplierAccountPresent {
		t.Error("Can't find detail for the supplier")
		return
	}
	// check detail for the bank (credit)
	var isBankAccountPresent bool = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].Journal == 570 {
			if newMovementDetails[i].Debit != 0 || newMovementDetails[i].Credit != float32(math.Round(float64(invoiceTotalAmount)*100)/100) {
				t.Error("The detail for the bank was not created successfully")
				return
			}
			isBankAccountPresent = true
			break
		}
	}
	if !isBankAccountPresent {
		t.Error("Can't find detail for the bank")
		return
	}

	// DELETE
	// delete the payment (that will delete the new accounting movement)
	ok = payment.deletePayment()
	if !ok {
		t.Error("Delete error, can't delete payment")
		return
	}

	// delete accounting movement
	ok = am.deleteAccountingMovement()
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
	}

	// delete invoioce and details
	ok = i.deletePurchaseInvoice()
	if !ok {
		t.Error("Delete error, can't delete purchase invoice")
	}
}
