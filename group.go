package main

type Group struct {
	Id            int16  `json:"id"`
	Name          string `json:"name"`
	Sales         bool   `json:"sales"`
	Purchases     bool   `json:"purchases"`
	Masters       bool   `json:"masters"`
	Warehouse     bool   `json:"warehouse"`
	Manufacturing bool   `json:"manufacturing"`
	Preparation   bool   `json:"preparation"`
	Admin         bool   `json:"admin"`
}

func getGroup() []Group {
	var groups []Group = make([]Group, 0)
	sqlStatement := `SELECT * FROM "group" ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return groups
	}
	for rows.Next() {
		g := Group{}
		rows.Scan(&g.Id, &g.Name, &g.Sales, &g.Purchases, &g.Masters, &g.Warehouse, &g.Manufacturing, &g.Preparation, &g.Admin)
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

	sqlStatement := `INSERT INTO public."group"(name, sales, purchases, masters, warehouse, manufacturing, preparation, admin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	res, err := db.Exec(sqlStatement, g.Name, g.Sales, g.Purchases, g.Masters, g.Warehouse, g.Manufacturing, g.Preparation, g.Admin)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (g *Group) updateGroup() bool {
	if g.Id <= 0 || !g.isValid() {
		return false
	}

	sqlStatement := `UPDATE public."group" SET name=$2, sales=$3, purchases=$4, masters=$5, warehouse=$6, manufacturing=$7, preparation=$8, admin=$9 WHERE id=$1`
	res, err := db.Exec(sqlStatement, g.Id, g.Name, g.Sales, g.Purchases, g.Masters, g.Warehouse, g.Manufacturing, g.Preparation, g.Admin)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (g *Group) deleteGroup() bool {
	if g.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public."group" WHERE id=$1`
	res, err := db.Exec(sqlStatement, g.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
