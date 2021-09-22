package main

import (
	"database/sql"
	"fmt"
)

type Account struct {
	Id            int32   `json:"id"`
	Journal       int16   `json:"journal"`
	Name          string  `json:"name"`
	Credit        float32 `json:"credit"`
	Debit         float32 `json:"debit"`
	Balance       float32 `json:"balance"`
	AccountNumber int32   `json:"accountNumber"`
	enterprise    int32
}

func getAccounts(enterpriseId int32) []Account {
	accounts := make([]Account, 0)
	sqlStatement := `SELECT * FROM public.account WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return accounts
	}

	for rows.Next() {
		a := Account{}
		rows.Scan(&a.Id, &a.Journal, &a.Name, &a.Credit, &a.Debit, &a.Balance, &a.AccountNumber, &a.enterprise)
		accounts = append(accounts, a)
	}

	return accounts
}

type AccountSearch struct {
	Journal int32  `json:"journal"`
	Search  string `json:"search"`
}

func (s *AccountSearch) searchAccounts(enterpriseId int32) []Account {
	accounts := make([]Account, 0)
	var rows *sql.Rows
	var err error
	if s.Journal <= 0 {
		sqlStatement := `SELECT * FROM public.account WHERE (name ILIKE $1) AND (enterprise=$2) ORDER BY id ASC`
		rows, err = db.Query(sqlStatement, "%"+s.Search+"%", enterpriseId)
	} else {
		sqlStatement := `SELECT * FROM public.account WHERE (name ILIKE $1) AND (journal=$2) AND (enterprise=$3) ORDER BY id ASC`
		rows, err = db.Query(sqlStatement, "%"+s.Search+"%", s.Journal, enterpriseId)
	}
	if err != nil {
		log("DB", err.Error())
		return accounts
	}

	for rows.Next() {
		a := Account{}
		rows.Scan(&a.Id, &a.Journal, &a.Name, &a.Credit, &a.Debit, &a.Balance, &a.AccountNumber, &a.enterprise)
		accounts = append(accounts, a)
	}

	return accounts
}

func getAccountRow(accountId int32) Account {
	sqlStatement := `SELECT * FROM public.account WHERE id=$1 LIMIT 1`
	row := db.QueryRow(sqlStatement, accountId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Account{}
	}

	a := Account{}
	row.Scan(&a.Id, &a.Journal, &a.Name, &a.Credit, &a.Debit, &a.Balance, &a.AccountNumber, &a.enterprise)

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
		if a.AccountNumber <= 0 {
			return false
		}
	}

	sqlStatement := `INSERT INTO public.account(journal, name, account_number, enterprise) VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(sqlStatement, a.Journal, a.Name, a.AccountNumber, a.enterprise)
	if row.Err() != nil {
		log("DB", row.Err().Error())
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
		log("DB", row.Err().Error())
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

	sqlStatement := `UPDATE public.account SET journal=$2, name=$3, account_number=$4 WHERE id=$1 AND enterprise=$5`
	_, err := db.Exec(sqlStatement, a.Id, a.Journal, a.Name, a.AccountNumber, a.enterprise)

	return err == nil
}

func (a *Account) deleteAccount() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.account WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func getAccountIdByAccountNumber(journal int16, accountNumber int32, enterpriseId int32) int32 {
	if journal <= 0 || accountNumber <= 0 || enterpriseId <= 0 {
		return 0
	}

	sqlStatement := `SELECT id FROM account WHERE account_number=$2 AND journal=$1 AND enterprise=$3 LIMIT 1`
	row := db.QueryRow(sqlStatement, journal, accountNumber, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
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

func locateAccountForCustomer(enterpriseId int32) []AccountLocate {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT id,journal,account_number,name FROM public.account WHERE journal=$1 AND enterprise=$2 ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement, s.CustomerJournal, enterpriseId)
	if err != nil {
		log("DB", err.Error())
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

func locateAccountForSupplier(enterpriseId int32) []AccountLocate {
	s := getSettingsRecordById(enterpriseId)
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT id,journal,account_number,name FROM public.account WHERE journal=$1 AND enterprise=$2 ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement, s.SupplierJournal, enterpriseId)
	if err != nil {
		log("DB", err.Error())
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

func locateAccountForBanks(enterpriseId int32) []AccountLocate {
	accounts := make([]AccountLocate, 0)
	sqlStatement := `SELECT account.id,account.journal,account.account_number,account.name FROM public.account INNER JOIN journal ON journal.id=account.journal WHERE journal.type='B' AND account.enterprise=$1 ORDER BY account_number ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
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

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
