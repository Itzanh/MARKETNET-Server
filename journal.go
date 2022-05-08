package main

type Journal struct {
	Id           int32    `json:"id" gorm:"primaryKey"`
	Name         string   `json:"name" gorm:"type:character varying(150);not null:true"`
	Type         string   `json:"type" gorm:"type:character(1);not null:true"` // S = Sale, P = Purchase, B = Bank, C = Cash, G = General
	EnterpriseId int32    `json:"-" gorm:"primaryKey;column:enterprise;not null:true"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (j *Journal) TableName() string {
	return "journal"
}

func getJournals(enterpriseId int32) []Journal {
	journals := make([]Journal, 0)
	dbOrm.Model(&Journal{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&journals)
	return journals
}

func (j *Journal) isValid() bool {
	return !(j.Id <= 0 || len(j.Name) == 0 || len(j.Name) == 150 || (j.Type != "S" && j.Type != "P" && j.Type != "B" && j.Type != "C" && j.Type != "G"))
}

func (j *Journal) insertJournal() bool {
	if !j.isValid() {
		return false
	}

	result := dbOrm.Create(j)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (j *Journal) updateJournal() bool {
	if !j.isValid() {
		return false
	}

	var journal Journal
	result := dbOrm.Where("id = ? AND enterprise = ?", j.Id, j.EnterpriseId).First(&journal)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	journal.Name = j.Name
	journal.Type = j.Type

	result = dbOrm.Save(&journal)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (j *Journal) deleteJournal() bool {
	if j.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", j.Id, j.EnterpriseId).Delete(&Journal{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
