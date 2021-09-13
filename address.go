package main

type Address struct {
	Id                int32   `json:"id"`
	Customer          *int32  `json:"customer"`
	Supplier          *int32  `json:"supplier"`
	Address           string  `json:"address"`
	Address2          string  `json:"address2"`
	City              string  `json:"city"`
	State             *int32  `json:"state"`
	Country           int16   `json:"country"`
	PrivateOrBusiness string  `json:"privateOrBusiness"` // P = Private, B = Business, _ = Not specified
	Notes             string  `json:"notes"`
	ZipCode           string  `json:"zipCode"`
	ContactName       string  `json:"contactName"`
	CountryName       string  `json:"countryName"`
	StateName         *string `json:"stateName"`
	PrestaShopId      int32
	ShopifyId         int64
}

type Addresses struct {
	Rows      int32     `json:"rows"`
	Addresses []Address `json:"addresses"`
}

func (q *PaginationQuery) getAddresses() Addresses {
	ad := Addresses{}
	if !q.isValid() {
		return ad
	}

	ad.Addresses = make([]Address, 0)
	sqlStatement := `SELECT *,CASE WHEN address.customer IS NOT NULL THEN (SELECT name FROM customer WHERE customer.id=address.customer) ELSE (SELECT name FROM suppliers WHERE suppliers.id=address.supplier) END,(SELECT name FROM country WHERE country.id=address.country),(SELECT name FROM state WHERE state.id=address.state) FROM address ORDER BY id ASC OFFSET $1 LIMIT $2`
	rows, err := db.Query(sqlStatement, q.Offset, q.Limit)
	if err != nil {
		log("DB", err.Error())
		return ad
	}
	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode, &a.ShopifyId, &a.ContactName, &a.CountryName, &a.StateName)
		ad.Addresses = append(ad.Addresses, a)
	}

	sqlStatement = `SELECT COUNT(*) FROM public.address`
	row := db.QueryRow(sqlStatement)
	row.Scan(&ad.Rows)

	return ad
}

func getAddressRow(addressId int32) Address {
	sqlStatement := `SELECT * FROM address WHERE id=$1`
	row := db.QueryRow(sqlStatement, addressId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return Address{}
	}

	a := Address{}
	row.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode, &a.ShopifyId)

	return a
}

func (s *PaginatedSearch) searchAddresses() Addresses {
	ad := Addresses{}
	if !s.isValid() {
		return ad
	}

	ad.Addresses = make([]Address, 0)
	sqlStatement := `SELECT address.*,CASE WHEN address.customer IS NOT NULL THEN (SELECT name FROM customer WHERE customer.id=address.customer) ELSE (SELECT name FROM suppliers WHERE suppliers.id=address.supplier) END,(SELECT name FROM country WHERE country.id=address.country),(SELECT name FROM state WHERE state.id=address.state) FROM address FULL JOIN customer ON customer.id=address.customer FULL JOIN state ON state.id=address.state FULL JOIN suppliers ON suppliers.id=address.supplier WHERE (address ILIKE $1 OR customer.name ILIKE $1 OR state.name ILIKE $1 OR suppliers.name ILIKE $1) AND (address.id > 0) ORDER BY id ASC OFFSET $2 LIMIT $3`
	rows, err := db.Query(sqlStatement, "%"+s.Search+"%", s.Offset, s.Limit)
	if err != nil {
		log("DB", err.Error())
		return ad
	}
	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode, &a.ShopifyId, &a.ContactName, &a.CountryName, &a.StateName)
		ad.Addresses = append(ad.Addresses, a)
	}

	sqlStatement = `SELECT COUNT(*) FROM address FULL JOIN customer ON customer.id=address.customer FULL JOIN state ON state.id=address.state FULL JOIN suppliers ON suppliers.id=address.supplier WHERE (address ILIKE $1 OR customer.name ILIKE $1 OR state.name ILIKE $1 OR suppliers.name ILIKE $1) AND (address.id > 0)`
	row := db.QueryRow(sqlStatement, "%"+s.Search+"%")
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ad
	}
	row.Scan(&ad.Rows)

	return ad
}

func (a *Address) isValid() bool {
	return !((a.Customer == nil && a.Supplier == nil) || (a.Customer != nil && *a.Customer <= 0) || (a.Supplier != nil && *a.Supplier <= 0) || len(a.Address) == 0 || len(a.Address) > 200 || len(a.Address2) > 200 || len(a.City) == 0 || len(a.City) > 100 || a.Country <= 0 || (a.PrivateOrBusiness != "P" && a.PrivateOrBusiness != "B" && a.PrivateOrBusiness != "_") || len(a.Notes) > 1000 || len(a.ZipCode) > 12)
}

func (a *Address) insertAddress() bool {
	if !a.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO address(customer, address, address_2, city, state, country, private_business, notes, supplier, ps_id, zip_code, sy_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	row := db.QueryRow(sqlStatement, a.Customer, a.Address, a.Address2, a.City, a.State, a.Country, a.PrivateOrBusiness, a.Notes, a.Supplier, a.PrestaShopId, a.ZipCode, a.ShopifyId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return false
	}

	var addressId int32
	row.Scan(&addressId)

	if a.Customer != nil {
		c := getCustomerRow(*a.Customer)
		var ok bool = false
		if c.MainAddress == nil {
			c.MainAddress = &addressId
			ok = true
		}
		if ok {
			c.updateCustomer()
		}
	}
	if a.Supplier != nil {
		s := getSupplierRow(*a.Supplier)
		var ok bool = false
		if s.MainAddress == nil {
			s.MainAddress = &addressId
			ok = true
		}
		if ok {
			s.updateSupplier()
		}
	}

	a.Id = addressId
	return true
}

func (a *Address) updateAddress() bool {
	if a.Id <= 0 || !a.isValid() {
		return false
	}

	sqlStatement := `UPDATE address SET customer=$2, address=$3, address_2=$4, city=$5, state=$6, country=$7, private_business=$8, notes=$9, supplier=$10, zip_code=$11 WHERE id = $1`
	res, err := db.Exec(sqlStatement, a.Id, a.Customer, a.Address, a.Address2, a.City, a.State, a.Country, a.PrivateOrBusiness, a.Notes, a.Supplier, a.ZipCode)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (a *Address) deleteAddress() bool {
	if a.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM address WHERE id = $1`
	res, err := db.Exec(sqlStatement, a.Id)
	if err != nil {
		log("DB", err.Error())
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type AddressLocate struct {
	Id      int32  `json:"id"`
	Address string `json:"address"`
}

func locateAddressByCustomer(customerId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	sqlStatement := `SELECT id, address FROM address WHERE customer = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, customerId)
	if err != nil {
		log("DB", err.Error())
		return addresses
	}
	for rows.Next() {
		a := AddressLocate{}
		rows.Scan(&a.Id, &a.Address)
		addresses = append(addresses, a)
	}

	return addresses
}

func locateAddressBySupplier(supplierId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	sqlStatement := `SELECT id, address FROM address WHERE supplier=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, supplierId)
	if err != nil {
		log("DB", err.Error())
		return addresses
	}
	for rows.Next() {
		a := AddressLocate{}
		rows.Scan(&a.Id, &a.Address)
		addresses = append(addresses, a)
	}

	return addresses
}

func getAddressName(addressId int32) string {
	sqlStatement := `SELECT address FROM address WHERE id = $1`
	row := db.QueryRow(sqlStatement, addressId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return ""
	}

	var address string
	row.Scan(&address)
	return address
}
