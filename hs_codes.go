package main

type HSCode struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type HSCodeQuery struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (q *HSCodeQuery) getHSCodes() []HSCode {
	var codes []HSCode = make([]HSCode, 0)
	sqlStatement := `SELECT * FROM public.hs_codes WHERE id LIKE $1 AND name ILIKE $2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, "%"+q.Id, "%"+q.Name+"%")
	if err != nil {
		log("DB", err.Error())
		return codes
	}

	for rows.Next() {
		c := HSCode{}
		rows.Scan(&c.Id, &c.Name)
		codes = append(codes, c)
	}

	return codes
}
