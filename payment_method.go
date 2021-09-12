package main

import "strings"

type PaymentMethod struct {
	Id                    int16  `json:"id"`
	Name                  string `json:"name"`
	PaidInAdvance         bool   `json:"paidInAdvance"`
	PrestashopModuleName  string `json:"prestashopModuleName"`
	DaysExpiration        int16  `json:"daysExpiration"`
	Bank                  *int32 `json:"bank"`
	WooCommerceModuleName string `json:"wooCommerceModuleName"`
	ShopifyModuleName     string `json:"shopifyModuleName"`
}

func getPaymentMethods() []PaymentMethod {
	var paymentMethod []PaymentMethod = make([]PaymentMethod, 0)
	sqlStatement := `SELECT * FROM public.payment_method ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return paymentMethod
	}
	for rows.Next() {
		p := PaymentMethod{}
		rows.Scan(&p.Id, &p.Name, &p.PaidInAdvance, &p.PrestashopModuleName, &p.DaysExpiration, &p.Bank, &p.WooCommerceModuleName, &p.ShopifyModuleName)
		paymentMethod = append(paymentMethod, p)
	}

	return paymentMethod
}

func getPaymentMethodRow(paymentMethodId int16) PaymentMethod {
	sqlStatement := `SELECT * FROM public.payment_method WHERE id=$1`
	row := db.QueryRow(sqlStatement, paymentMethodId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return PaymentMethod{}
	}

	p := PaymentMethod{}
	row.Scan(&p.Id, &p.Name, &p.PaidInAdvance, &p.PrestashopModuleName, &p.DaysExpiration, &p.Bank, &p.WooCommerceModuleName, &p.ShopifyModuleName)

	return p
}

func (p *PaymentMethod) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 100 || p.DaysExpiration < 0)
}

func (p *PaymentMethod) insertPaymentMethod() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.payment_method(name, paid_in_advance, prestashop_module_name, days_expiration, bank, woocommerce_module_name, shopify_module_name) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, p.Name, p.PaidInAdvance, p.PrestashopModuleName, p.DaysExpiration, p.Bank, p.WooCommerceModuleName, p.ShopifyModuleName)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *PaymentMethod) updatePaymentMethod() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.payment_method SET name=$2, paid_in_advance=$3, prestashop_module_name=$4, days_expiration=$5, bank=$6, woocommerce_module_name=$7, shopify_module_name=$8 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.PaidInAdvance, p.PrestashopModuleName, p.DaysExpiration, p.Bank, p.WooCommerceModuleName, p.ShopifyModuleName)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *PaymentMethod) deletePaymentMethod() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.payment_method WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func findPaymentMethodByName(paymentMethodName string) []NameInt16 {
	var paymentMethods []NameInt16 = make([]NameInt16, 0)
	sqlStatement := `SELECT id,name FROM public.payment_method WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(paymentMethodName))
	if err != nil {
		log("DB", err.Error())
		return paymentMethods
	}
	for rows.Next() {
		p := NameInt16{}
		rows.Scan(&p.Id, &p.Name)
		paymentMethods = append(paymentMethods, p)
	}

	return paymentMethods
}

func getNamePaymentMethod(id int16) string {
	sqlStatement := `SELECT name FROM public.payment_method WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
