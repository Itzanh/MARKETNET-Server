package main

import "strings"

type PaymentMethod struct {
	Id            int16  `json:"id"`
	Name          string `json:"name"`
	PaidInAdvance bool   `json:"paidInAdvance"`
}

func getPaymentMethods() []PaymentMethod {
	var paymentMethod []PaymentMethod = make([]PaymentMethod, 0)
	sqlStatement := `SELECT * FROM public.payment_method ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return paymentMethod
	}
	for rows.Next() {
		p := PaymentMethod{}
		rows.Scan(&p.Id, &p.Name, &p.PaidInAdvance)
		paymentMethod = append(paymentMethod, p)
	}

	return paymentMethod
}

func (p *PaymentMethod) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 100)
}

func (p *PaymentMethod) insertPaymentMethod() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.payment_method(name, paid_in_advance) VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, p.Name, p.PaidInAdvance)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *PaymentMethod) updatePaymentMethod() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.payment_method SET name=$2, paid_in_advance=$3 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.PaidInAdvance)
	if err != nil {
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
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type PaymentMethodName struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

func findPaymentMethodByName(paymentMethodName string) []PaymentMethodName {
	var paymentMethods []PaymentMethodName = make([]PaymentMethodName, 0)
	sqlStatement := `SELECT id,name FROM public.payment_method WHERE UPPER(name) LIKE $1 || '%' ORDER BY id ASC LIMIT 10`
	rows, err := db.Query(sqlStatement, strings.ToUpper(paymentMethodName))
	if err != nil {
		return paymentMethods
	}
	for rows.Next() {
		p := PaymentMethodName{}
		rows.Scan(&p.Id, &p.Name)
		paymentMethods = append(paymentMethods, p)
	}

	return paymentMethods
}

func getNamePaymentMethod(id int16) string {
	sqlStatement := `SELECT name FROM public.payment_method WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return ""
	}
	name := ""
	row.Scan(&name)
	return name
}
