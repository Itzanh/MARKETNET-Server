package main

type Packages struct {
	Id     int16   `json:"id"`
	Name   string  `json:"name"`
	Weight float32 `json:"weight"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Depth  float32 `json:"depth"`
}

func getPackages() []Packages {
	var products []Packages = make([]Packages, 0)
	sqlStatement := `SELECT * FROM public.packages ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return products
	}
	for rows.Next() {
		p := Packages{}
		rows.Scan(&p.Id, &p.Name, &p.Weight, &p.Width, &p.Height, &p.Depth)
		products = append(products, p)
	}

	return products
}

func getPackagesRow(packageId int16) Packages {
	sqlStatement := `SELECT * FROM public.packages WHERE id=$1`
	row := db.QueryRow(sqlStatement, packageId)
	if row.Err() != nil {
		return Packages{}
	}

	p := Packages{}
	row.Scan(&p.Id, &p.Name, &p.Weight, &p.Width, &p.Height, &p.Depth)

	return p
}

func (p *Packages) isValid() bool {
	return !(len(p.Name) == 0 || len(p.Name) > 50 || p.Weight < 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0)
}

func (p *Packages) insertPackage() bool {
	if !p.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.packages(name, weight, width, height, depth) VALUES ($1, $2, $3, $4, $5)`
	res, err := db.Exec(sqlStatement, p.Name, p.Weight, p.Width, p.Height, p.Depth)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Packages) updatePackage() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.packages SET name=$2, weight=$3, width=$4, height=$5, depth=$6 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Name, p.Weight, p.Width, p.Height, p.Depth)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Packages) deletePackage() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.packages WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
