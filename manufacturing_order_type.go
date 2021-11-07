package main

type ManufacturingOrderType struct {
	Id                   int32  `json:"id"`
	Name                 string `json:"name"`
	QuantityManufactured int32  `json:"quantityManufactured"`
	enterprise           int32
}

func getManufacturingOrderType(enterpriseId int32) []ManufacturingOrderType {
	var types []ManufacturingOrderType = make([]ManufacturingOrderType, 0)
	sqlStatement := `SELECT * FROM public.manufacturing_order_type WHERE enterprise=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return types
	}
	for rows.Next() {
		t := ManufacturingOrderType{}
		rows.Scan(&t.Id, &t.Name, &t.enterprise, &t.QuantityManufactured)
		types = append(types, t)
	}

	return types
}

func getManufacturingOrderTypeRow(typeId int32) ManufacturingOrderType {
	sqlStatement := `SELECT * FROM public.manufacturing_order_type WHERE id=$1 ORDER BY id ASC`
	row := db.QueryRow(sqlStatement, typeId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ManufacturingOrderType{}
	}

	t := ManufacturingOrderType{}
	row.Scan(&t.Id, &t.Name, &t.enterprise, &t.QuantityManufactured)

	return t
}

func (t *ManufacturingOrderType) isValid() bool {
	return !(len(t.Name) == 0 || len(t.Name) > 100) || t.QuantityManufactured < 1
}

func (t *ManufacturingOrderType) insertManufacturingOrderType() bool {
	if !t.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.manufacturing_order_type(name, enterprise, quantity_manufactured) VALUES ($1, $2, $3) RETURNING id`
	row := db.QueryRow(sqlStatement, t.Name, t.enterprise, t.QuantityManufactured)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	row.Scan(&t.Id)

	return t.Id > 0
}

func (t *ManufacturingOrderType) updateManufacturingOrderType() bool {
	if t.Id <= 0 || !t.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order_type SET name=$2, quantity_manufactured=$4 WHERE id=$1 AND enterprise=$3`
	res, err := db.Exec(sqlStatement, t.Id, t.Name, t.enterprise, t.QuantityManufactured)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (t *ManufacturingOrderType) deleteManufacturingOrderType() bool {
	if t.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.manufacturing_order_type WHERE id=$1 AND enterprise=$2`
	res, err := db.Exec(sqlStatement, t.Id, t.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}
