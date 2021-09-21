package main

type Incoterm struct {
	Id         int32  `json:"id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	enterprise int32
}

func getIncoterm(enterpriseId int32) []Incoterm {
	var incoterms []Incoterm = make([]Incoterm, 0)
	sqlStatement := `SELECT * FROM public.incoterm WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return incoterms
	}
	for rows.Next() {
		i := Incoterm{}
		rows.Scan(&i.Id, &i.Key, &i.Name, &i.enterprise)
		incoterms = append(incoterms, i)
	}

	return incoterms
}

func (i *Incoterm) isValid() bool {
	return !(len(i.Key) == 0 || len(i.Key) > 3 || len(i.Name) == 0 || len(i.Name) > 50)
}

func (i *Incoterm) insertIncoterm() bool {
	if !i.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.incoterm(key, name, enterprise) VALUES ($1, $2, $3)`
	res, err := db.Exec(sqlStatement, i.Key, i.Name, i.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (i *Incoterm) updateIncoterm() bool {
	if i.Id <= 0 || !i.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.incoterm SET key=$2, name=$3 WHERE id=$1 AND enterprise=$4`
	res, err := db.Exec(sqlStatement, i.Id, i.Key, i.Name, i.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (i *Incoterm) deleteIncoterm() bool {
	if i.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.incoterm WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, i.Id, i.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
