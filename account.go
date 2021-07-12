package main

import "fmt"

type Account struct {
	Id            int32   `json:"id"`
	Journal       int16   `json:"journal"`
	Name          string  `json:"name"`
	Credit        float32 `json:"credit"`
	Debit         float32 `json:"debit"`
	Balance       float32 `json:"balance"`
	AccountNumber int32   `json:"accountNumber"`
}

func getAccounts() []Account {
	accounts := make([]Account, 0)
	sqlStatement := `SELECT * FROM public.account ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return accounts
	}

	for rows.Next() {
		a := Account{}
		rows.Scan(&a.Id, &a.Journal, &a.Name, &a.Credit, &a.Debit, &a.Balance, &a.AccountNumber)
		accounts = append(accounts, a)
	}

	return accounts
}

func getAccountRow(accountId int32) Account {
	sqlStatement := `SELECT * FROM public.account WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, accountId)
	if row.Err() != nil {
		return Account{}
	}

	a := Account{}
	row.Scan(&a.Id, &a.Journal, &a.Name, &a.Credit, &a.Debit, &a.Balance, &a.AccountNumber)

	return a
}

func (a *Account) isValid() bool {
	return !(a.Journal <= 0 || len(a.Name) == 0 || len(a.Name) > 150)
}

func (a *Account) insertAccount() bool {
	if !a.isValid() {
		return false
	}

	if a.AccountNumber <= 0 {
		a.AccountNumber = a.getNextAccountNumber()
	}

	sqlStatement := `INSERT INTO public.account(journal, name, account_number) VALUES ($1, $2, $3) RETURNING id`
	row := db.QueryRow(sqlStatement, a.Journal, a.Name, a.AccountNumber)
	if row.Err() != nil {
		return false
	}

	var accountId int32
	row.Scan(&accountId)
	a.Id = accountId

	return accountId > 0
}

func (a *Account) getNextAccountNumber() int32 {
	sqlStatement := `SELECT CASE WHEN MAX(account_number) IS NULL THEN 1 ELSE MAX(account_number) + 1 END FROM account WHERE journal=$1`
	row := db.QueryRow(sqlStatement, a.Journal)
	if row.Err() != nil {
		return 0
	}

	var accountNumber int32
	row.Scan(&accountNumber)
	return accountNumber
}

func (a *Account) updateAccount() bool {
	if a.Id <= 0 || !a.isValid() || a.AccountNumber <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.account SET journal=$2, name=$3, account_number=$4 WHERE id=$1`
	_, err := db.Exec(sqlStatement, a.Id, a.Journal, a.Name, a.AccountNumber)

	return err == nil
}

func (a *Account) deleteAccount() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.account WHERE id=$1`
	_, err := db.Exec(sqlStatement, a.Id)

	return err == nil
}

func getAccountIdByAccountNumber(journal int16, accountNumber int32) int32 {
	if journal <= 0 || accountNumber <= 0 {
		return 0
	}

	sqlStatement := `SELECT id FROM account WHERE account_number=$2 AND journal=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, journal, accountNumber)
	if row.Err() != nil {
		return 0
	}

	var accountId int32
	row.Scan(&accountId)
	return accountId
}

type AccountLocate struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func locateAccountForCustomer() []AccountLocate {
	s := getSettingsRecord()
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT id,journal,account_number,name FROM public.account WHERE journal=$1 ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement, s.CustomerJournal)
	if err != nil {
		return accounts
	}

	var journal int16
	var accountNumber int32
	var name string

	for rows.Next() {
		a := AccountLocate{}
		rows.Scan(&a.Id, &journal, &accountNumber, &name)
		a.Name = fmt.Sprintf("%d.%06d - %s", journal, accountNumber, name)
		accounts = append(accounts, a)
	}

	return accounts
}

func locateAccountForSupplier() []AccountLocate {
	s := getSettingsRecord()
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT id,journal,account_number,name FROM public.account WHERE journal=$1 ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement, s.SupplierJournal)
	if err != nil {
		return accounts
	}

	var journal int16
	var accountNumber int32
	var name string

	for rows.Next() {
		a := AccountLocate{}
		rows.Scan(&a.Id, &journal, &accountNumber, &name)
		a.Name = fmt.Sprintf("%d.%06d - %s", journal, accountNumber, name)
		accounts = append(accounts, a)
	}

	return accounts
}

func locateAccountForBanks() []AccountLocate {
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT account.id,account.journal,account.account_number,account.name FROM public.account INNER JOIN journal ON journal.id=account.journal WHERE journal.type='B' ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return accounts
	}

	var journal int16
	var accountNumber int32
	var name string

	for rows.Next() {
		a := AccountLocate{}
		rows.Scan(&a.Id, &journal, &accountNumber, &name)
		a.Name = fmt.Sprintf("%d.%06d - %s", journal, accountNumber, name)
		accounts = append(accounts, a)
	}

	return accounts
}

// Will add or take out credit and debit (if given a negative amount)
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func (a *Account) addCreditAndDebit(credit float32, debit float32) bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `UPDATE public.account SET debit=debit+$2,credit=credit+$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, a.Id, debit, credit)

	return err == nil
}
