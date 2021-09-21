package main

type Packages struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	Weight     float32 `json:"weight"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Depth      float32 `json:"depth"`
	Product    int32   `json:"product"`
	enterprise int32
}

func getPackages(enterpriseId int32) []Packages {
	var products []Packages = make([]Packages, 0)
	sqlStatement := `SELECT * FROM public.packages WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return products
	}
	for rows.Next() {
		p := Packages{}
		rows.Scan(&p.Id, &p.Name, &p.Weight, &p.Width, &p.Height, &p.Depth, &p.Product, &p.enterprise)
		products = append(products, p)
	}

	return products
}

func getPackagesRow(packageId int32) Packages {
	sqlStatement := `SELECT * FROM public.packages WHERE id=$1`
	row := db.QueryRow(sqlStatement, packageId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Packages{}
	}

	p := Packages{}
	row.Scan(&p.Id, &p.Name, &p.Weight, &p.Width, &p.Height, &p.Depth, &p.Product, &p.enterprise)

	return p
}

func (p *Packages) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 50 || p.Weight < 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 || p.Product <= 0)
}

func (p *Packages) insertPackage() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.packages(name, weight, width, height, depth, product, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, p.Name, p.Weight, p.Width, p.Height, p.Depth, p.Product, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Packages) updatePackage() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.packages SET name=$2, weight=$3, width=$4, height=$5, depth=$6, product=$7, enterprise=$8 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.Weight, p.Width, p.Height, p.Depth, p.Product, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Packages) deletePackage() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.packages WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, p.Id, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
