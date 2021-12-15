package main

// When adding fields to this table/struct, add them also to:
// The Permissions struct, and the getUserPermissions function at the bottom of this file.
// The getUserGroupsIn function at the beginning of the user_group.go file.
// The initialGroup function in the initial_data.go file, the admin always has all the permissions set.
type Group struct {
	Id            int32  `json:"id"`
	Name          string `json:"name"`
	Sales         bool   `json:"sales"`
	Purchases     bool   `json:"purchases"`
	Masters       bool   `json:"masters"`
	Warehouse     bool   `json:"warehouse"`
	Manufacturing bool   `json:"manufacturing"`
	Preparation   bool   `json:"preparation"`
	Admin         bool   `json:"admin"`
	PrestaShop    bool   `json:"prestashop"`
	Accounting    bool   `json:"accounting"`
	PointOfSale   bool   `json:"pointOfSale"`
	enterprise    int32
}

func getGroup(enterpriseId int32) []Group {
	var groups []Group = make([]Group, 0)
	sqlStatement := `SELECT * FROM "group" WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return groups
	}
	for rows.Next() {
		g := Group{}
		rows.Scan(&g.Id, &g.Name, &g.Sales, &g.Purchases, &g.Masters, &g.Warehouse, &g.Manufacturing, &g.Preparation, &g.Admin, &g.PrestaShop, &g.Accounting, &g.enterprise, &g.PointOfSale)
		groups = append(groups, g)
	}

	return groups
}

func (g *Group) isValid() bool {
	return !(len(g.Name) > 50)
}

func (g *Group) insertGroup() bool {
	if !g.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public."group"(name, sales, purchases, masters, warehouse, manufacturing, preparation, admin, prestashop, accounting, enterprise, point_of_sale) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	row := db.QueryRow(sqlStatement, g.Name, g.Sales, g.Purchases, g.Masters, g.Warehouse, g.Manufacturing, g.Preparation, g.Admin, g.PrestaShop, g.Accounting, g.enterprise, g.PointOfSale)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	row.Scan(&g.Id)

	return g.Id > 0
}

func (g *Group) updateGroup() bool {
	if g.Id <= 0 || !g.isValid() {
		return false
	}

	sqlStatement := `UPDATE public."group" SET name=$2, sales=$3, purchases=$4, masters=$5, warehouse=$6, manufacturing=$7, preparation=$8, admin=$9, prestashop=$10, accounting=$11, point_of_sale=$13 WHERE id=$1 AND enterprise=$12`
	res, err := db.Exec(sqlStatement, g.Id, g.Name, g.Sales, g.Purchases, g.Masters, g.Warehouse, g.Manufacturing, g.Preparation, g.Admin, &g.PrestaShop, g.Accounting, g.enterprise, g.PointOfSale)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (g *Group) deleteGroup() bool {
	if g.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public."group" WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, g.Id, g.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type Permissions struct {
	Sales         bool `json:"sales"`
	Purchases     bool `json:"purchases"`
	Masters       bool `json:"masters"`
	Warehouse     bool `json:"warehouse"`
	Manufacturing bool `json:"manufacturing"`
	Preparation   bool `json:"preparation"`
	Admin         bool `json:"admin"`
	PrestaShop    bool `json:"prestashop"`
	Accounting    bool `json:"accounting"`
	PointOfSale   bool `json:"pointOfSale"`
}

func getUserPermissions(userId int32, enterpriseId int32) Permissions {
	ug := getUserGroups(userId, enterpriseId)
	p := Permissions{}

	for i := 0; i < len(ug.GroupsIn); i++ {
		if ug.GroupsIn[i].Sales {
			p.Sales = true
		}
		if ug.GroupsIn[i].Purchases {
			p.Purchases = true
		}
		if ug.GroupsIn[i].Masters {
			p.Masters = true
		}
		if ug.GroupsIn[i].Warehouse {
			p.Warehouse = true
		}
		if ug.GroupsIn[i].Manufacturing {
			p.Manufacturing = true
		}
		if ug.GroupsIn[i].Preparation {
			p.Preparation = true
		}
		if ug.GroupsIn[i].Admin {
			p.Admin = true
		}
		if ug.GroupsIn[i].PrestaShop {
			p.PrestaShop = true
		}
		if ug.GroupsIn[i].Accounting {
			p.Accounting = true
		}
		if ug.GroupsIn[i].PointOfSale {
			p.PointOfSale = true
		}
	}

	return p
}
