package main

// The permission dictionary is a way to give or revoke very specific privileges to the groups of users in the ERP.
// For example, a permission in the dictionary could be: MANUALLY_CREATE_SALE_INVOICE, and if it's not given to the group, you can't press the "Add" button in the sale invoices menu.
// This permissions, and it's relation with the groups is stored in the database, so it's completely dynamic.
// A function will tell if a given user (int32) has a given permission (string) and will return a boolean, and it can be used however we want.

type PermissionDictionary struct {
	EnterpriseId int32    `json:"-" gorm:"primaryKey;not null:true;column:enterprise"`
	Key          string   `json:"key" gorm:"primaryKey;not null:true;type:character varying(150)"`
	Description  string   `json:"description" gorm:"not null:true;type:character varying(250)"`
	Enterprise   Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (pd *PermissionDictionary) TableName() string {
	return "permission_dictionary"
}

func getPermissionDictionary(enterpriseId int32) []PermissionDictionary {
	var dictionary []PermissionDictionary = make([]PermissionDictionary, 0)
	dbOrm.Model(&PermissionDictionary{}).Where("enterprise = ?", enterpriseId).Order("key ASC").Find(&dictionary)
	return dictionary
}

type PermissionDictionaryGroup struct {
	GroupId         int32                `json:"group" gorm:"primaryKey;column:group;not null:true"`
	Group           Group                `json:"-" gorm:"foreignKey:GroupId,EnterpriseId;references:Id,EnterpriseId"`
	PermissionKeyId string               `json:"permissionKey" gorm:"primaryKey;column:permission_key;not null:true;type:character varying(150)"`
	PermissionKey   PermissionDictionary `json:"-" gorm:"foreignKey:PermissionKeyId,EnterpriseId;references:Key,EnterpriseId"`
	EnterpriseId    int32                `json:"-" gorm:"primaryKey;not null:true;column:enterprise"`
	Enterprise      Settings             `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
}

func (pdg *PermissionDictionaryGroup) TableName() string {
	return "permission_dictionary_group"
}

type PermissionDictionaryGroupInOut struct {
	In  []PermissionDictionaryGroup `json:"in"`
	Out []PermissionDictionary      `json:"out"`
}

func getGroupPermissionDictionary(enterpriseId int32, groupId int32) PermissionDictionaryGroupInOut {
	return PermissionDictionaryGroupInOut{
		In:  getPermissionDictionaryGroupIn(enterpriseId, groupId),
		Out: getPermissionDictionaryGroupOut(enterpriseId, groupId),
	}
}

func getPermissionDictionaryGroupIn(enterpriseId int32, groupId int32) []PermissionDictionaryGroup {
	var dictionary []PermissionDictionaryGroup = make([]PermissionDictionaryGroup, 0)
	result := dbOrm.Model(&PermissionDictionaryGroup{}).Where("permission_dictionary_group.\"group\" = ? AND permission_dictionary_group.enterprise = ?", groupId, enterpriseId).Joins("Group").Joins("PermissionKey").Joins("Enterprise").Order("permission_key ASC").Find(&dictionary)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return dictionary
	}
	return dictionary
}

func getPermissionDictionaryGroupOut(enterpriseId int32, groupId int32) []PermissionDictionary {
	dictionary := getPermissionDictionary(enterpriseId)
	permissions := getPermissionDictionaryGroupIn(enterpriseId, groupId)

	for i := len(dictionary) - 1; i >= 0; i-- {
		d := dictionary[i]
		for j := 0; j < len(permissions); j++ {
			if permissions[j].PermissionKeyId == d.Key {
				dictionary = append(dictionary[:i], dictionary[i+1:]...)
				permissions = append(permissions[:j], permissions[j+1:]...) // little optimization :D
				break
			}
		}
	}

	return dictionary
}

func (p *PermissionDictionaryGroup) isValid() bool {
	return !(p.GroupId <= 0 || len(p.PermissionKeyId) == 0 || len(p.PermissionKeyId) > 150 || p.EnterpriseId <= 0)
}

func (p *PermissionDictionaryGroup) insertPermissionDictionaryGroup() bool {
	if !p.isValid() {
		return false
	}

	result := dbOrm.Create(&p)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (p *PermissionDictionaryGroup) deletePermissionDictionaryGroup() bool {
	if !p.isValid() {
		return false
	}

	result := dbOrm.Where(`"group" = ? AND permission_key = ? AND enterprise = ?`, p.GroupId, p.PermissionKeyId, p.EnterpriseId).Delete(&PermissionDictionaryGroup{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func getPermissionDictionaryUserGroupInForWebClient(userId int32) []string {
	var dictionary []string = make([]string, 0)
	groups := getUserGroupsIn(userId)

	for i := 0; i < len(groups); i++ {
		var permissionsInGroup []PermissionDictionaryGroup = make([]PermissionDictionaryGroup, 0)
		dbOrm.Model(&PermissionDictionaryGroup{}).Where("\"group\" = ?", groups[i].Id).Find(&permissionsInGroup)

		for j := 0; j < len(permissionsInGroup); j++ {
			dictionary = append(dictionary, permissionsInGroup[j].PermissionKeyId)
		}
	}

	return dictionary
}

func getUserPermission(permission string, enterpriseId int32, userId int32) bool {
	groups := getUserGroupsIn(userId)

	for i := 0; i < len(groups); i++ {
		var count int64
		dbOrm.Model(&PermissionDictionaryGroup{}).Where(`"group" = ? AND permission_key = ? AND enterprise = ?`, groups[i].Id, permission, enterpriseId).Count(&count)
		if count > 0 {
			return true
		}
	}
	return false
}
