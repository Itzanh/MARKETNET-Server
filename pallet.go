package main

type Pallet struct {
	Id         int32   `json:"id"`
	SalesOrder int64   `json:"salesOrder"`
	Weight     float64 `json:"weight"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Depth      float64 `json:"depth"`
	Name       string  `json:"name"`
	enterprise int32
}

type Pallets struct {
	HasPallets bool     `json:"hasPallets"`
	Pallets    []Pallet `json:"pallets"`
}

func getSalesOrderPallets(orderId int64, enterpriseId int32) Pallets {
	sqlStatement := `SELECT pallets FROM sales_order INNER JOIN carrier ON carrier.id=sales_order.carrier WHERE sales_order.id=$1 AND sales_order.enterprise=$2`
	row := db.QueryRow(sqlStatement, orderId, enterpriseId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Pallets{}
	}

	var hasPallets bool
	row.Scan(&hasPallets)
	if !hasPallets {
		return Pallets{HasPallets: false}
	}

	var pallets []Pallet = make([]Pallet, 0)
	sqlStatement = `SELECT * FROM public.pallets WHERE sales_order=$1 AND pallets.enterprise=$2 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return Pallets{}
	}
	for rows.Next() {
		p := Pallet{}
		rows.Scan(&p.Id, &p.SalesOrder, &p.Weight, &p.Width, &p.Height, &p.Depth, &p.Name, &p.enterprise)
		pallets = append(pallets, p)
	}

	return Pallets{HasPallets: true, Pallets: pallets}
}

func getPalletsRow(palletId int32) Pallet {
	sqlStatement := `SELECT * FROM public.pallets WHERE id=$1`
	row := db.QueryRow(sqlStatement, palletId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Pallet{}
	}

	p := Pallet{}
	row.Scan(&p.Id, &p.SalesOrder, &p.Weight, &p.Width, &p.Height, &p.Depth, &p.Name, &p.enterprise)

	return p
}

func (p *Pallet) isValid() bool {
	return !(p.Weight <= 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 || len(p.Name) == 0 || len(p.Name) > 40)
}

func (p *Pallet) insertPallet() bool {
	if p.SalesOrder <= 0 || len(p.Name) == 0 || len(p.Name) > 40 {
		return false
	}

	s := getSettingsRecordById(p.enterprise)
	p.Weight = s.PalletWeight
	p.Width = s.PalletWidth
	p.Height = s.PalletHeight
	p.Depth = s.PalletDepth

	sqlStatement := `INSERT INTO public.pallets(sales_order, weight, width, height, depth, name, enterprise) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := db.Exec(sqlStatement, p.SalesOrder, p.Weight, p.Width, p.Height, p.Depth, p.Name, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Pallet) updatePallet() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.pallets SET weight=$2, width=$3, height=$4, depth=$5, name=$6 WHERE id=$1 AND enterprise=$7`
	res, err := db.Exec(sqlStatement, p.Id, p.Weight, p.Width, p.Height, p.Depth, p.Name, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Pallet) deletePallet() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.pallets WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, p.Id, p.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
