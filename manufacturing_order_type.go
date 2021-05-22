package main

type ManufacturingOrderType struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

func getManufacturingOrderType() []ManufacturingOrderType {
	var types []ManufacturingOrderType = make([]ManufacturingOrderType, 0)
	sqlStatement := `SELECT * FROM public.manufacturing_order_type ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return types
	}
	for rows.Next() {
		t := ManufacturingOrderType{}
		rows.Scan(&t.Id, &t.Name)
		types = append(types, t)
	}

	return types
}

func (t *ManufacturingOrderType) isValid() bool {
	return !(len(t.Name) == 0 || len(t.Name) > 100)
}

func (t *ManufacturingOrderType) insertManufacturingOrderType() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.manufacturing_order_type(name) VALUES ($1)`
	res, err := db.Exec(sqlStatement, t.Name)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (t *ManufacturingOrderType) updateManufacturingOrderType() bool {
	if t.Id <= 0 || !t.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order_type SET name=$2 WHERE id=$1`
	res, err := db.Exec(sqlStatement, t.Id, t.Name)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (t *ManufacturingOrderType) deleteManufacturingOrderType() bool {
	if t.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.manufacturing_order_type WHERE id=$1`
	res, err := db.Exec(sqlStatement, t.Id)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
