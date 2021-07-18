package main

type Journal struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func getJournals() []Journal {
	journals := make([]Journal, 0)
	sqlStatement := `SELECT * FROM public.journal ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("DB", err.Error())
		return journals
	}

	for rows.Next() {
		j := Journal{}
		rows.Scan(&j.Id, &j.Name, &j.Type)
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

	sqlStatement := `INSERT INTO public.journal(id, name, type) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, j.Id, j.Name, j.Type)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (j *Journal) updateJournal() bool {
	if !j.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.journal SET name=$2, type=$3 WHERE id=$1`
	_, err := db.Exec(sqlStatement, j.Id, j.Name, j.Type)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (j *Journal) deleteJournal() bool {
	if j.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.journal WHERE id=$1`
	_, err := db.Exec(sqlStatement, j.Id)

	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
