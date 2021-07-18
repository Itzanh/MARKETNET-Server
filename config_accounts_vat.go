package main

type ConfigAccountsVat struct {
	VatPercent            float32 `json:"vatPercent"`
	AccountSale           int32   `json:"accountSale"`
	AccountPurchase       int32   `json:"accountPurchase"`
	AccountSaleNumber     int32   `json:"accountSaleNumber"`
	JournalSale           int16   `json:"journalSale"`
	AccountPurchaseNumber int32   `json:"accountPurchaseNumber"`
	JournalPurchase       int16   `json:"journalPurchase"`
}

func getConfigAccountsVat() []ConfigAccountsVat {
	configAccountsVat := make([]ConfigAccountsVat, 0)
	sqlStatement := `SELECT *,(SELECT journal FROM account WHERE account.id=config_accounts_vat.account_sale),(SELECT account_number FROM account WHERE account.id=config_accounts_vat.account_sale),(SELECT journal FROM account WHERE account.id=config_accounts_vat.account_purchase),(SELECT account_number FROM account WHERE account.id=config_accounts_vat.account_purchase) FROM public.config_accounts_vat ORDER BY vat_percent ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return configAccountsVat
	}

	for rows.Next() {
		c := ConfigAccountsVat{}
		rows.Scan(&c.VatPercent, &c.AccountSale, &c.AccountPurchase, &c.JournalSale, &c.AccountSaleNumber, &c.JournalPurchase, &c.AccountPurchaseNumber)
		configAccountsVat = append(configAccountsVat, c)
	}

	return configAccountsVat
}

func getConfigAccountsVatSaleRow(vatPercent float32) (int16, int32) {
	sqlStatement := `SELECT (SELECT journal FROM account WHERE account.id=config_accounts_vat.account_sale),(SELECT account_number FROM account WHERE account.id=config_accounts_vat.account_sale) FROM config_accounts_vat WHERE vat_percent=$1`
	row := db.QueryRow(sqlStatement, vatPercent)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0, 0
	}

	var journal int16
	var accountNumber int32
	row.Scan(&journal, &accountNumber)
	return journal, accountNumber
}

func getConfigAccountsVatPurchaseRow(vatPercent float32) (int16, int32) {
	sqlStatement := `SELECT (SELECT journal FROM account WHERE account.id=config_accounts_vat.account_purchase),(SELECT account_number FROM account WHERE account.id=config_accounts_vat.account_purchase) FROM config_accounts_vat WHERE vat_percent=$1`
	row := db.QueryRow(sqlStatement, vatPercent)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0, 0
	}

	var journal int16
	var accountNumber int32
	row.Scan(&journal, &accountNumber)
	return journal, accountNumber
}

func (c *ConfigAccountsVat) isValid() bool {
	return !(c.VatPercent <= 0 || c.AccountSaleNumber <= 0 || c.JournalSale <= 0 || c.AccountPurchaseNumber <= 0 || c.JournalPurchase <= 0)
}

func (c *ConfigAccountsVat) insertConfigAccountsVat() bool {
	if !c.isValid() {
		return false
	}

	c.AccountSale = getAccountIdByAccountNumber(c.JournalSale, c.AccountSaleNumber)
	c.AccountPurchase = getAccountIdByAccountNumber(c.JournalPurchase, c.AccountPurchaseNumber)
	if c.AccountSale <= 0 || c.AccountPurchase <= 0 {
		return false
	}

	sqlStatement := `INSERT INTO public.config_accounts_vat(vat_percent, account_sale, account_purchase) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, c.VatPercent, c.AccountSale, c.AccountPurchase)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (c *ConfigAccountsVat) deleteConfigAccountsVat() bool {
	sqlStatement := `DELETE FROM public.config_accounts_vat WHERE vat_percent=$1`
	_, err := db.Exec(sqlStatement, c.VatPercent)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
