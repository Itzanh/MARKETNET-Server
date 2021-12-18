package main

// The permission dictionary is a way to give or revoke very specific privileges to the groups of users in the ERP.
// For example, a permission in the dictionary could be: MANUALLY_CREATE_SALE_INVOICE, and if it's not given to the group, you can't press the "Add" button in the sale invoices menu.
// This permissions, and it's relation with the groups is stored in the database, so it's completely dynamic.
// A function will tell if a given user (int32) has a given permission (string) and will return a boolean, and it can be used however we want.

type PermissionDictionary struct {
	enterprise  int32
	Key         string `json:"key"`
	Description string `json:"description"`
}

func getPermissionDictionary(enterpriseId int32) []PermissionDictionary {
	var dictionary []PermissionDictionary = make([]PermissionDictionary, 0)
	sqlStatement := `SELECT * FROM public.permission_dictionary WHERE enterprise = $1 ORDER BY key ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return dictionary
	}

	for rows.Next() {
		d := PermissionDictionary{}
		rows.Scan(&d.enterprise, &d.Key, &d.Description)
		dictionary = append(dictionary, d)
	}

	return dictionary
}

type PermissionDictionaryGroup struct {
	Group         int32  `json:"group"`
	PermissionKey string `json:"permissionKey"`
	Description   string `json:"description"`
	enterprise    int32
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
	sqlStatement := `SELECT *,(SELECT description FROM permission_dictionary WHERE permission_dictionary.enterprise = permission_dictionary_group.enterprise AND permission_dictionary.key = permission_dictionary_group.permission_key) FROM public.permission_dictionary_group WHERE "group" = $1 AND enterprise = $2 ORDER BY permission_key ASC`
	rows, err := db.Query(sqlStatement, groupId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return dictionary
	}

	for rows.Next() {
		d := PermissionDictionaryGroup{}
		rows.Scan(&d.Group, &d.PermissionKey, &d.enterprise, &d.Description)
		dictionary = append(dictionary, d)
	}

	return dictionary
}

func getPermissionDictionaryGroupOut(enterpriseId int32, groupId int32) []PermissionDictionary {
	dictionary := getPermissionDictionary(enterpriseId)
	permissions := getPermissionDictionaryGroupIn(enterpriseId, groupId)

	for i := len(dictionary) - 1; i >= 0; i-- {
		d := dictionary[i]
		for j := 0; j < len(permissions); j++ {
			if permissions[j].PermissionKey == d.Key {
				dictionary = append(dictionary[:i], dictionary[i+1:]...)
				permissions = append(permissions[:j], permissions[j+1:]...) // little optimization :D
				break
			}
		}
	}

	return dictionary
}

func (p *PermissionDictionaryGroup) isValid() bool {
	return !(p.Group <= 0 || len(p.PermissionKey) == 0 || len(p.PermissionKey) > 150 || p.enterprise <= 0)
}

func (p *PermissionDictionaryGroup) insertPermissionDictionaryGroup() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.permission_dictionary_group("group", permission_key, enterprise) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, p.Group, p.PermissionKey, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

func (p *PermissionDictionaryGroup) deletePermissionDictionaryGroup() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `DELETE FROM public.permission_dictionary_group WHERE "group" = $1 AND permission_key = $2 AND enterprise = $3`
	_, err := db.Exec(sqlStatement, p.Group, p.PermissionKey, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	return err == nil
}

func getPermissionDictionaryUserGroupInForWebClient(userId int32) []string {
	var dictionary []string = make([]string, 0)
	sqlStatement := `SELECT permission_dictionary_group.permission_key FROM public.permission_dictionary_group INNER JOIN "group" ON "group".id=permission_dictionary_group."group" INNER JOIN user_group ON user_group."group"="group".id WHERE user_group."user" = $1`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		log("DB", err.Error())
		return dictionary
	}

	for rows.Next() {
		d := ""
		rows.Scan(&d)
		dictionary = append(dictionary, d)
	}

	return dictionary
}

func getUserPermission(permission string, enterpriseId int32, userId int32) bool {
	sqlStatement := `SELECT COUNT(*) FROM public.permission_dictionary_group WHERE "group" = $1 AND permission_key = $2 AND enterprise = $3`
	groups := getUserGroupsIn(userId)

	for i := 0; i < len(groups); i++ {
		row := db.QueryRow(sqlStatement, groups[i].Id, permission, enterpriseId)
		if row.Err() != nil {
			log("DB", row.Err().Error())
			return false
		}

		var count int16
		row.Scan(&count)
		if count > 0 {
			return true
		}
	}
	return false
}
