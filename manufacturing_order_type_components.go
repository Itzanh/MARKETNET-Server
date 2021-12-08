package main

type ManufacturingOrderTypeComponents struct {
	Id                     int32  `json:"id"`
	ManufacturingOrderType int32  `json:"manufacturingOrderType"`
	Type                   string `json:"type"` // I = Input, O = Output
	Product                int32  `json:"product"`
	Quantity               int32  `json:"quantity"`
	ProductName            string `json:"productName"`
	enterprise             int32
}

func getManufacturingOrderTypeComponents(manfuacturingOrderTypeId int32, enterpriserId int32) []ManufacturingOrderTypeComponents {
	var components []ManufacturingOrderTypeComponents = make([]ManufacturingOrderTypeComponents, 0)
	manufacturingOrderType := getManufacturingOrderTypeRow(manfuacturingOrderTypeId)
	if manufacturingOrderType.enterprise != enterpriserId {
		return components
	}

	sqlStatement := `SELECT *,(SELECT name FROM product WHERE product.id=manufacturing_order_type_components.product) FROM public.manufacturing_order_type_components WHERE manufacturing_order_type=$1 ORDER BY product ASC`
	rows, err := db.Query(sqlStatement, manfuacturingOrderTypeId)
	if err != nil {
		log("DB", err.Error())
		return components
	}

	for rows.Next() {
		var c ManufacturingOrderTypeComponents
		rows.Scan(&c.Id, &c.ManufacturingOrderType, &c.Type, &c.Product, &c.Quantity, &c.enterprise, &c.ProductName)
		components = append(components, c)
	}

	return components
}

func getManufacturingOrderTypeComponentRow(manfuacturingOrderTypeId int32) ManufacturingOrderTypeComponents {
	c := ManufacturingOrderTypeComponents{}

	sqlStatement := `SELECT * FROM manufacturing_order_type_components WHERE id=$1`
	row := db.QueryRow(sqlStatement, manfuacturingOrderTypeId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return c
	}

	row.Scan(&c.Id, &c.ManufacturingOrderType, &c.Type, &c.Product, &c.Quantity, &c.enterprise)
	return c
}

func (c *ManufacturingOrderTypeComponents) isValid() bool {
	return !(c.ManufacturingOrderType <= 0 || (c.Type != "I" && c.Type != "O") || c.Product <= 0 || c.Quantity <= 0)
}

func (c *ManufacturingOrderTypeComponents) insertManufacturingOrderTypeComponents() bool {
	if !c.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO public.manufacturing_order_type_components(manufacturing_order_type, type, product, quantity, enterprise) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, c.ManufacturingOrderType, c.Type, c.Product, c.Quantity, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (c *ManufacturingOrderTypeComponents) updateManufacturingOrderTypeComponents() bool {
	if c.Id <= 0 || !c.isValid() {
		return false
	}

	sqlStatement := `UPDATE public.manufacturing_order_type_components SET manufacturing_order_type=$2, type=$3, product=$4, quantity=$5, enterprise=$6 WHERE id=$1`
	_, err := db.Exec(sqlStatement, c.Id, c.ManufacturingOrderType, c.Type, c.Product, c.Quantity, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}

func (c *ManufacturingOrderTypeComponents) deleteManufacturingOrderTypeComponents() bool {
	if c.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.manufacturing_order_type_components WHERE id=$1 AND enterprise=$2`
	_, err := db.Exec(sqlStatement, c.Id, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false
	}
	return true
}
