package main

type Journal struct {
	Id         int16  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"` // S = Sale, P = Purchase, B = Bank, C = Cash, G = General
	enterprise int32
}

func getJournals(enterpriseId int32) []Journal {
	journals := make([]Journal, 0)
	sqlStatement := `SELECT * FROM public.journal WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return journals
	}

	for rows.Next() {
		j := Journal{}
		rows.Scan(&j.Id, &j.Name, &j.Type, &j.enterprise)
		journals = append(journals, j)
	}

	return journals
}

func (j *Journal) isValid() bool {
	return !(j.Id <= 0 || len(j.Name) == 0 || len(j.Name) == 150 || (j.Type != "S" && j.Type != "P" && j.Type != "B" && j.Type != "C" && j.Type != "G"))
}

func (j *Journal) insertJournal() bool {
	if !j.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.journal(id, name, type, enterprise) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, j.Id, j.Name, j.Type, j.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (j *Journal) updateJournal() bool {
	if !j.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.journal SET name=$2, type=$3 WHERE id=$1 AND enterprise=$4`
	_, err := db.Exec(sqlStatement, j.Id, j.Name, j.Type, j.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (j *Journal) deleteJournal() bool {
	if j.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.journal WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, j.Id, j.enterprise)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
