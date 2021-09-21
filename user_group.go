package main

type UserGroup struct {
	User  int32 `json:"user"`
	Group int32 `json:"group"`
}

type UserGroups struct {
	GroupsIn  []Group `json:"groupsIn"`
	GroupsOut []Group `json:"groupsOut"`
}

func getUserGroups(userId int32, enterpriseId int32) UserGroups {
	groupsIn := getUserGroupsIn(userId)
	return UserGroups{
		GroupsIn:  groupsIn,
		GroupsOut: getUserGroupsOut(userId, groupsIn, enterpriseId),
	}
}

func getUserGroupsIn(userId int32) []Group {
	var groups []Group = make([]Group, 0)
	sqlStatement := `SELECT "group".* FROM "user" INNER JOIN user_group ON "user".id=user_group.user INNER JOIN "group" ON "group".id=user_group.group WHERE "user".id=$1 ORDER BY "group".id ASC`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		log("DB", err.Error())
		return groups
	}
	for rows.Next() {
		g := Group{}
		rows.Scan(&g.Id, &g.Name, &g.Sales, &g.Purchases, &g.Masters, &g.Warehouse, &g.Manufacturing, &g.Preparation, &g.Admin, &g.PrestaShop, &g.Accounting, &g.enterprise)
		groups = append(groups, g)
	}

	return groups
}

func getUserGroupsOut(userId int32, groupsIn []Group, enterpriseId int32) []Group {
	groups := getGroup(enterpriseId)

	for i := 0; i < len(groupsIn); i++ {
		for j := len(groups) - 1; j >= 0; j-- {
			if groupsIn[i].Id == groups[j].Id {
				groups = append(groups[0:j], groups[j+1:]...)
				break
			}
		}
	}

	return groups
}

func (u *UserGroup) isValid() bool {
	return !(u.User <= 0 || u.Group <= 0)
}

func (u *UserGroup) insertUserGroup() bool {
	if !u.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.user_group("user", "group") VALUES ($1, $2)`
	res, err := db.Exec(sqlStatement, u.User, u.Group)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (u *UserGroup) deleteUserGroup() bool {
	if !u.isValid() {
		return false
	}

	sqlStatement := `DELETE FROM public.user_group WHERE "user"=$1 AND "group"=$2`
	res, err := db.Exec(sqlStatement, u.User, u.Group)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
