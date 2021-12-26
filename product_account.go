package main

type ProductAccount struct {
	Id            int32  `json:"id"`
	Product       int32  `json:"product"`
	Account       int32  `json:"account"`
	Journal       int32  `json:"journal"`
	AccountNumber int32  `json:"accountNumber"`
	Type          string `json:"type"` // S = Sale, P = Purchase
	AccountName   string `json:"accountName"`
	enterprise    int32
}

func getProductAccounts(productId int32, enterpriseId int32) []ProductAccount {
	var accounts []ProductAccount = make([]ProductAccount, 0)
	sqlStatement := `SELECT *,(SELECT name FROM account WHERE account.id=product_account.account) FROM public.product_account WHERE product = $1 AND enterprise = $2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, productId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return accounts
	}
	defer rows.Close()

	for rows.Next() {
		a := ProductAccount{}
		rows.Scan(&a.Id, &a.Product, &a.Account, &a.Journal, &a.AccountNumber, &a.enterprise, &a.Type, &a.AccountName)
		accounts = append(accounts, a)
	}
	return accounts
}

// returns the specified account for the product, or nil if not specified by the user
// accType: S = Sales, p = Purchase
func getProductAccount(productId int32, accType string) *ProductAccount {
	if productId <= 0 || (accType != "S" && accType != "P") {
		return nil
	}

	sqlStatement := `SELECT COUNT(*) FROM public.product_account WHERE product = $1 AND type = $2`
	row := db.QueryRow(sqlStatement, productId, accType)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return nil
	}

	var rowsNumber int16
	row.Scan(&rowsNumber)

	if rowsNumber == 0 {
		return nil
	}

	sqlStatement = `SELECT * FROM public.product_account WHERE product = $1 AND type = $2`
	row = db.QueryRow(sqlStatement, productId, accType)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return nil
	}

	a := ProductAccount{}
	row.Scan(&a.Id, &a.Product, &a.Account, &a.Journal, &a.AccountNumber, &a.enterprise, &a.Type)
	return &a
}

func (a *ProductAccount) isValid() bool {
	return !(a.Product <= 0 || a.Account <= 0 || a.enterprise <= 0 || (a.Type != "S" && a.Type != "P"))
}

func (a *ProductAccount) insertProductAccount() bool {
	if !a.isValid() {
		return false
	}

	account := getAccountRow(a.Account)
	if account.Id <= 0 {
		return false
	}
	a.Journal = account.Journal
	a.AccountNumber = account.AccountNumber

	sqlStatement := `INSERT INTO public.product_account(product, account, jorunal, account_number, enterprise, type) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(sqlStatement, a.Product, a.Account, a.Journal, a.AccountNumber, a.enterprise, a.Type)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (a *ProductAccount) updateProductAccount() bool {
	if !a.isValid() || a.Id <= 0 {
		return false
	}

	account := getAccountRow(a.Account)
	if account.Id <= 0 {
		return false
	}
	a.Journal = account.Journal
	a.AccountNumber = account.AccountNumber

	sqlStatement := `UPDATE public.product_account SET account=$3, jorunal=$4, account_number=$5 WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise, a.Account, a.Journal, a.AccountNumber)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}

func (a *ProductAccount) deleteProductAccount() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.product_account WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, a.Id, a.enterprise)
	if err != nil {
		log("DB", err.Error())
	}
	return err == nil
}
