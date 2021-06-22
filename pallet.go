package main

type Pallet struct {
	Id         int32   `json:"id"`
	SalesOrder int32   `json:"salesOrder"`
	Weight     float32 `json:"weight"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Depth      float32 `json:"depth"`
	Name       string  `json:"name"`
}

type Pallets struct {
	HasPallets bool     `json:"hasPallets"`
	Pallets    []Pallet `json:"pallets"`
}

func getSalesOrderPallets(orderId int32) Pallets {
	sqlStatement := `SELECT pallets FROM sales_order INNER JOIN carrier ON carrier.id=sales_order.carrier WHERE sales_order.id=$1`
	row := db.QueryRow(sqlStatement, orderId)
	if row.Err() != nil {
		return Pallets{}
	}

	var hasPallets bool
	row.Scan(&hasPallets)
	if !hasPallets {
		return Pallets{HasPallets: false}
	}

	var pallets []Pallet = make([]Pallet, 0)
	sqlStatement = `SELECT * FROM public.pallets WHERE sales_order = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, orderId)
	if err != nil {
		return Pallets{}
	}
	for rows.Next() {
		p := Pallet{}
		rows.Scan(&p.Id, &p.SalesOrder, &p.Weight, &p.Width, &p.Height, &p.Depth, &p.Name)
		pallets = append(pallets, p)
	}

	return Pallets{HasPallets: true, Pallets: pallets}
}

func (p *Pallet) isValid() bool {
	return !(p.Weight <= 0 || p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 || len(p.Name) == 0 || len(p.Name) > 40)
}

func (p *Pallet) insertPallet() bool {
	if p.SalesOrder <= 0 || len(p.Name) == 0 || len(p.Name) > 40 {
		return false
	}

	s := getSettingsRecord()
	p.Weight = s.PalletWeight
	p.Width = s.PalletWidth
	p.Height = s.PalletHeight
	p.Depth = s.PalletDepth

	sqlStatement := `INSERT INTO public.pallets(sales_order, weight, width, height, depth, name) VALUES ($1, $2, $3, $4, $5, $6)`
	res, err := db.Exec(sqlStatement, p.SalesOrder, p.Weight, p.Width, p.Height, p.Depth, p.Name)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Pallet) updatePallet() bool {
	if p.Id <= 0 || !p.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.pallets SET weight=$2, width=$3, height=$4, depth=$5, name=$6 WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id, p.Weight, p.Width, p.Height, p.Depth, p.Name)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (p *Pallet) deletePallet() bool {
	if p.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.pallets WHERE id=$1`
	res, err := db.Exec(sqlStatement, p.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
