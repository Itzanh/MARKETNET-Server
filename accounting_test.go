package main

import (
	"testing"
)

// ===== JOURNAL

func TestGetJournals(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	j := getJournals(1)

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
		Id:           123,
		Name:         "Test",
		Type:         "G",
		EnterpriseId: 1,
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

	a := getAccounts(1)

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
	a := s.searchAccounts(1)

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
		Id:           123,
		Name:         "Test",
		Type:         "G",
		EnterpriseId: 1,
	}
	j.insertJournal()

	a := Account{
		JournalId:    123,
		Name:         "Test",
		EnterpriseId: 1,
	}
	ok := a.insertAccount()
	if !ok {
		t.Error("Insert error, can't insert account")
		return
	}

	getAccountId := getAccountIdByAccountNumber(j.Id, a.AccountNumber, 1)
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

	a := locateAccountForCustomer(1)

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

	a := locateAccountForSupplier(1)

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

	a := locateAccountForBanks(1)

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

	m := getAccountingMovement(1)

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

	query := AccountingMovementSearch{}
	m := query.searchAccountingMovements(1)

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

	movements := getAccountingMovement(1)
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
		Type:           "N",
		BillingSerieId: "INT",
		EnterpriseId:   1,
	}

	ok := am.insertAccountingMovement(0, nil)
	if !ok {
		t.Error("Insert error, can't insert accounting movement")
		return
	}

	ok = am.deleteAccountingMovement(0, nil)
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
		return
	}
}

// ===== COLLECTION OPERATION

func TestGetPendingColletionOperations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	collectionOperation := getPendingColletionOperations(1)
	if len(collectionOperation) > 0 && collectionOperation[0].Id == 0 {
		t.Errorf("Can't scan collection operations")
	}
}

func TestSearchCollectionOperations(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	search := CollectionOperationPaymentTransactionSearch{
		Mode:      0,
		StartDate: nil,
		EndDate:   nil,
		Search:    "",
	}

	collectionOperation := searchCollectionOperations(search, 1)
	if len(collectionOperation) > 0 && collectionOperation[0].Id == 0 {
		t.Errorf("Can't scan collection operations")
	}
}

func TestGetColletionOperationRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	var id int32 = 1

	collectionOperations := getPendingColletionOperations(1)
	if len(collectionOperations) > 0 && collectionOperations[0].Id > 0 {
		id = collectionOperations[0].Id
	}

	collectionOperation := getColletionOperationRow(id)
	if collectionOperation.Id == 0 {
		t.Errorf("Can't scan collection operations")
	}
}

// ===== PAYMENT OPERATION

func TestGetPendingPaymentTransaction(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	paymentTransaction := getPendingPaymentTransaction(1)
	if len(paymentTransaction) > 0 && paymentTransaction[0].Id == 0 {
		t.Errorf("Can't scan payment transactions")
	}
}

func TestSearchPaymentTransactions(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	search := CollectionOperationPaymentTransactionSearch{
		Mode:      0,
		StartDate: nil,
		EndDate:   nil,
		Search:    "",
	}

	paymentTransaction := searchPaymentTransactions(search, 1)
	if len(paymentTransaction) > 0 && paymentTransaction[0].Id == 0 {
		t.Errorf("Can't scan payment transactions")
	}
}

func TestGetPaymentTransactionRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	var id int32 = 1

	paymentTransactions := getPendingPaymentTransaction(1)
	if len(paymentTransactions) > 0 && paymentTransactions[0].Id > 0 {
		id = paymentTransactions[0].Id
	}

	paymentTransaction := getPaymentTransactionRow(id)
	if paymentTransaction.Id == 0 {
		t.Errorf("Can't scan collection operations")
	}
}

// ===== POST SALE INVOICES

func TestSalesPostInvoices(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	var product int32 = 4
	var product2 int32 = 1

	// insert invoice and details
	i := SalesInvoice{
		CustomerId:       1,
		PaymentMethodId:  1,
		BillingSeriesId:  "INT",
		CurrencyId:       1,
		BillingAddressId: 1,
		EnterpriseId:     1,
	}
	_, invoiceId := i.insertSalesInvoice(0, nil)
	d := SalesInvoiceDetail{
		InvoiceId:    invoiceId,
		ProductId:    &product,
		Price:        9.99,
		Quantity:     2,
		VatPercent:   21,
		EnterpriseId: 1,
	}
	d.insertSalesInvoiceDetail(nil, 0)
	d = SalesInvoiceDetail{
		InvoiceId:    invoiceId,
		ProductId:    &product2,
		Price:        4.99,
		Quantity:     1,
		VatPercent:   21,
		EnterpriseId: 1,
	}
	d.insertSalesInvoiceDetail(nil, 0)
	i = getSalesInvoiceRow(invoiceId)

	// post sale invoice
	result := salesPostInvoices([]int64{invoiceId}, 1, 0)
	for i := 0; i < len(result); i++ {
		if !result[i].Ok {
			t.Error("Can't post sale invoice")
			return
		}
	}

	// check accounting movement
	movements := getAccountingMovement(1)
	am := movements[0]

	saleInvoices := getAccountingMovementSaleInvoices(am.Id)
	if len(saleInvoices) == 0 || saleInvoices[0].Id != invoiceId {
		t.Error("Can't scan sale invocies")
		return
	}

	movementDetails := getAccountingMovementDetail(am.Id, 1)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the customer (debit)
	invoiceTotalAmount := i.TotalAmount
	var isCustomerAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].JournalId == 430 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != invoiceTotalAmount {
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
		if movementDetails[i].JournalId == 700 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != invoiceTotalProducts {
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
	invoiceTotalVat := i.VatAmount
	var isVatAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].JournalId == 477 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != invoiceTotalVat {
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
	collectionOperations := getColletionOperations(am.Id, 1)
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
		CollectionOperationId: collectionOperations[0].Id,
		Amount:                collectionOperations[0].Pending,
		Concept:               "TESTING...",
		EnterpriseId:          1,
	}
	ok := charge.insertCharges(0)
	if !ok {
		t.Error("Can't insert charges")
		return
	}

	// test charges scan
	charges := getCharges(collectionOperations[0].Id, 1)
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
	collectionOperations = getColletionOperations(am.Id, 1)
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

	newMovementDetails := getAccountingMovementDetail(newAccountingMovement.Id, 1)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the customer (credit)
	isCustomerAccountPresent = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].JournalId == 430 {
			if newMovementDetails[i].Debit != 0 || newMovementDetails[i].Credit != invoiceTotalAmount {
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
		if newMovementDetails[i].JournalId == 570 {
			if newMovementDetails[i].Credit != 0 || newMovementDetails[i].Debit != invoiceTotalAmount {
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
	charge.EnterpriseId = 1
	ok = charge.deleteCharges(0)
	if !ok {
		t.Error("Delete error, can't delete charge")
		return
	}

	// delete accounting movement
	ok = am.deleteAccountingMovement(0, nil)
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
		return
	}

	// delete invoioce and details
	ok = i.deleteSalesInvoice(0).Ok
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

	var product int32 = 4
	var product2 int32 = 1

	// insert invoice and details
	i := PurchaseInvoice{
		SupplierId:       1,
		PaymentMethodId:  1,
		BillingSeriesId:  "INT",
		CurrencyId:       1,
		BillingAddressId: 1,
		EnterpriseId:     1,
	}
	_, invoiceId := i.insertPurchaseInvoice(0, nil)
	d := PurchaseInvoiceDetail{
		InvoiceId:    invoiceId,
		ProductId:    &product,
		Price:        9.99,
		Quantity:     2,
		VatPercent:   21,
		EnterpriseId: 1,
	}
	d.insertPurchaseInvoiceDetail(0, nil)
	d = PurchaseInvoiceDetail{
		InvoiceId:    invoiceId,
		ProductId:    &product2,
		Price:        4.99,
		Quantity:     1,
		VatPercent:   21,
		EnterpriseId: 1,
	}
	d.insertPurchaseInvoiceDetail(0, nil)
	i = getPurchaseInvoiceRow(invoiceId)

	// post purchase invoice
	result := purchasePostInvoices([]int64{invoiceId}, 1, 0)
	for i := 0; i < len(result); i++ {
		if !result[i].Ok {
			t.Error("Can't post purchase invoice")
			return
		}
	}

	// check accounting movement
	movements := getAccountingMovement(1)
	am := movements[0]

	purchaseInvoices := getAccountingMovementPurchaseInvoices(am.Id)
	if len(purchaseInvoices) == 0 || purchaseInvoices[0].Id != invoiceId {
		t.Error("Can't scan purchase invocies")
		return
	}

	movementDetails := getAccountingMovementDetail(am.Id, 1)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the supplier (credit)
	invoiceTotalAmount := i.TotalAmount
	var isSupplierAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].JournalId == 400 {
			if movementDetails[i].Debit != 0 || movementDetails[i].Credit != invoiceTotalAmount {
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
		if movementDetails[i].JournalId == 600 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != invoiceTotalProducts {
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
	invoiceTotalVat := i.VatAmount
	var isVatAccountPresent bool = false
	for i := 0; i < len(movementDetails); i++ {
		if movementDetails[i].JournalId == 472 {
			if movementDetails[i].Credit != 0 || movementDetails[i].Debit != invoiceTotalVat {
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
	paymentTransactions := getPaymentTransactions(am.Id, 1)
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
		PaymentTransactionId: paymentTransactions[0].Id,
		Amount:               paymentTransactions[0].Pending,
		Concept:              "TESTING...",
		EnterpriseId:         1,
	}
	ok := payment.insertPayment(0)
	if !ok {
		t.Error("Can't insert payments")
		return
	}

	// test payments scan
	payments := getPayments(paymentTransactions[0].Id, 1)
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
	paymentTransactions = getPaymentTransactions(am.Id, 1)
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

	newMovementDetails := getAccountingMovementDetail(newAccountingMovement.Id, 1)
	if len(movementDetails) == 0 || movementDetails[0].Id <= 0 {
		t.Error("Can't scan movement details")
		return
	}

	// check details
	// check detail for the supplier (debit)
	isSupplierAccountPresent = false
	for i := 0; i < len(newMovementDetails); i++ {
		if newMovementDetails[i].JournalId == 400 {
			if newMovementDetails[i].Credit != 0 || newMovementDetails[i].Debit != invoiceTotalAmount {
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
		if newMovementDetails[i].JournalId == 570 {
			if newMovementDetails[i].Debit != 0 || newMovementDetails[i].Credit != invoiceTotalAmount {
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
	ok = payment.deletePayment(0)
	if !ok {
		t.Error("Delete error, can't delete payment")
		return
	}

	// delete accounting movement
	ok = am.deleteAccountingMovement(0, nil)
	if !ok {
		t.Error("Delete error, can't delete accounting movement")
	}

	// delete invoioce and details
	ok = i.deletePurchaseInvoice(0, nil).Ok
	if !ok {
		t.Error("Delete error, can't delete purchase invoice")
	}
}
