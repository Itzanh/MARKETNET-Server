package main

import "strings"

type ReportTemplateTranslation struct {
	enterprise   int32
	Key          string `json:"key"`
	Language     int32  `json:"language"`
	Translation  string `json:"translation"`
	LanguageName string `json:"languageName"`
}

func getReportTemplateTranslations(enterpriseId int32) []ReportTemplateTranslation {
	var translations []ReportTemplateTranslation = make([]ReportTemplateTranslation, 0)
	sqlStatement := `SELECT *,(SELECT name FROM language WHERE language.id=report_template_translation.language) FROM public.report_template_translation WHERE enterprise = $1 ORDER BY key ASC, language ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return translations
	}

	for rows.Next() {
		t := ReportTemplateTranslation{}
		rows.Scan(&t.enterprise, &t.Key, &t.Language, &t.Translation, &t.LanguageName)
		translations = append(translations, t)
	}
	return translations
}

func (t *ReportTemplateTranslation) isValid() bool {
	return !(t.enterprise <= 0 || len(t.Key) == 0 || len(t.Key) > 50 || t.Language <= 0 || len(t.Translation) == 0 || len(t.Translation) > 255)
}

func (t *ReportTemplateTranslation) insertReportTemplateTranslation() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.report_template_translation(enterprise, key, language, translation) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, t.enterprise, t.Key, t.Language, t.Translation)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (t *ReportTemplateTranslation) updateReportTemplateTranslation() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.report_template_translation SET translation=$4 WHERE enterprise=$1 AND key=$2 AND language=$3`
	_, err := db.Exec(sqlStatement, t.enterprise, t.Key, t.Language, t.Translation)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (t *ReportTemplateTranslation) deleteReportTemplateTranslation() bool {
	if t.enterprise <= 0 || len(t.Key) == 0 || t.Language <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.report_template_translation WHERE enterprise=$1 AND key=$2 AND language=$3`
	_, err := db.Exec(sqlStatement, t.enterprise, t.Key, t.Language)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func translateReport(reportContent string, languageId int32, enterpriseId int32) string {
	sqlStatement := `SELECT key, translation FROM public.report_template_translation WHERE enterprise = $1 AND language = $2`
	rows, err := db.Query(sqlStatement, enterpriseId, languageId)
	if err != nil {
		log("DB", err.Error())
		return reportContent
	}

	var key string
	var translation string
	for rows.Next() {
		rows.Scan(&key, &translation)

		reportContent = strings.ReplaceAll(reportContent, "{{"+key+"}}", translation)
	}

	return reportContent
}
