package main

type ReportTemplate struct {
	enterprise int32
	Key        string `json:"key"`
	Html       string `json:"html"`
}

func getReportTemplates(enterpriseId int32) []ReportTemplate {
	templates := make([]ReportTemplate, 0)
	sqlStatement := `SELECT * FROM public.report_template WHERE enterprise=$1`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return templates
	}

	for rows.Next() {
		t := ReportTemplate{}
		rows.Scan(&t.enterprise, &t.Key, &t.Html)
		templates = append(templates, t)
	}
	return templates
}

func getReportTemplate(enterpriseId int32, key string) ReportTemplate {
	t := ReportTemplate{}
	sqlStatement := `SELECT * FROM public.report_template WHERE enterprise=$1 AND key=$2`
	row := db.QueryRow(sqlStatement, enterpriseId, key)
	if row.Err() != nil {
		return t
	}

	row.Scan(&t.enterprise, &t.Key, &t.Html)
	return t
}

// Must NOT be callable from the web client!
func (r ReportTemplate) insertReportTemplate() {
	sqlStatement := `INSERT INTO public.report_template(enterprise, key, html) VALUES ($1, $2, $3)`
	db.Exec(sqlStatement, r.enterprise, r.Key, r.Html)
}

func (r *ReportTemplate) updateReportTemplate() bool {
	if len(r.Key) == 0 || len(r.Html) == 0 || len(r.Html) > 5000000 {
		return false
	}

	sqlStatement := `UPDATE public.report_template SET html=$3 WHERE enterprise=$1 AND key=$2`
	db.Exec(sqlStatement, r.enterprise, r.Key, r.Html)
	return false
}
