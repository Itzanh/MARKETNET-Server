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

// returns:
// ok
// code
// 0 = parameter error / ok
// 1 = the input product has the same manufacturing order type as the component
// 2 = the output product doesn't have the same manufacturing order type as the component
// 3 = the product already exist in one of the components
func (c *ManufacturingOrderTypeComponents) isValid() (bool, uint8) {
	if c.Product <= 0 {
		return false, 0
	}
	// the manufacturing order type has to be the same as this one for the output, and different on the input to make sure that there are no recursivity errors
	product := getProductRow(c.Product)
	if product.Id <= 0 {
		return false, 0
	}
	if c.Type == "I" {
		if product.ManufacturingOrderType != nil && *product.ManufacturingOrderType == c.ManufacturingOrderType {
			return false, 1
		}
	} else if c.Type == "O" {
		if product.ManufacturingOrderType == nil || *product.ManufacturingOrderType != c.ManufacturingOrderType {
			return false, 2
		}
	} else {
		return false, 0
	}

	if c.Id > 0 { // update
		// check that the product has not been associated yet
		components := getManufacturingOrderTypeComponents(c.ManufacturingOrderType, c.enterprise)
		for i := 0; i < len(components); i++ {
			if components[i].Id != c.Id && components[i].Product == c.Product {
				return false, 3
			}
		}
	} else { // insert
		// check that the product has not been associated yet
		components := getManufacturingOrderTypeComponents(c.ManufacturingOrderType, c.enterprise)
		for i := 0; i < len(components); i++ {
			if components[i].Product == c.Product {
				return false, 3
			}
		}
	}

	return !(c.ManufacturingOrderType <= 0 || (c.Type != "I" && c.Type != "O") || c.Quantity <= 0), 0
}

func (c *ManufacturingOrderTypeComponents) insertManufacturingOrderTypeComponents() (bool, uint8) {
	ok, errorCode := c.isValid()
	if !ok {
		return false, errorCode
	}

	sqlStatement := `INSERT INTO public.manufacturing_order_type_components(manufacturing_order_type, type, product, quantity, enterprise) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, c.ManufacturingOrderType, c.Type, c.Product, c.Quantity, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false, 0
	}
	return true, 0
}

func (c *ManufacturingOrderTypeComponents) updateManufacturingOrderTypeComponents() (bool, uint8) {
	ok, errorCode := c.isValid()
	if c.Id <= 0 || !ok {
		return false, errorCode
	}

	sqlStatement := `UPDATE public.manufacturing_order_type_components SET manufacturing_order_type=$2, type=$3, product=$4, quantity=$5, enterprise=$6 WHERE id=$1`
	_, err := db.Exec(sqlStatement, c.Id, c.ManufacturingOrderType, c.Type, c.Product, c.Quantity, c.enterprise)
	if err != nil {
		log("DB", err.Error())
		return false, 0
	}
	return true, 0
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
