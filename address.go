package main

type Address struct {
	Id                int32  `json:"id"`
	Customer          int32  `json:"customer"`
	Address           string `json:"address"`
	Address2          string `json:"address2"`
	City              int32  `json:"city"`
	Province          string `json:"province"`
	Country           int16  `json:"country"`
	PrivateOrBusiness string `json:"privateOrBusiness"`
	Notes             string `json:"notes"`
}

func getAddresses() []Address {
	var addresses []Address = make([]Address, 0)
	sqlStatement := `SELECT * FROM address ORDER BY id ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return addresses
	}
	for rows.Next() {
		a := Address{}
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.City, &a.Province, &a.Country, &a.PrivateOrBusiness, &a.Notes)
		addresses = append(addresses, a)
	}

	return addresses
}

func (a *Address) isValid() bool {
	return !(a.Customer <= 0 || len(a.Address) == 0 || len(a.Address) > 200 || len(a.Address2) > 200 || a.City <= 0 || a.Country <= 0 || (a.PrivateOrBusiness != "P" && a.PrivateOrBusiness != "B") || len(a.Notes) > 3000)
}

func (a *Address) insertAddress() bool {
	if !a.isValid() {
		return false
	}

	sqlStatement := `INSERT INTO address(customer, address, address_2, city, province, country, private_business, notes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	res, err := db.Exec(sqlStatement, a.Customer, a.Address, a.Address2, a.City, a.Province, a.Country, a.PrivateOrBusiness, a.Notes)
	if err != nil {
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (a *Address) updateAddress() bool {
	if a.Id <= 0 || !a.isValid() {
		return false
	}

	sqlStatement := `UPDATE address SET customer=$2, address=$3, address_2=$4, city=$5, province=$6, country=$7, private_business=$8, notes=$9 WHERE id = $1`
	res, err := db.Exec(sqlStatement, a.Id, a.Customer, a.Address, a.Address2, a.City, a.Province, a.Country, a.PrivateOrBusiness, a.Notes)
	if err != nil {
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
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

type AddressLocate struct {
	Id      int32  `json:"id"`
	Address string `json:"address"`
}

func locateAddress(customerId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	sqlStatement := `SELECT id, address FROM address WHERE customer = $1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, customerId)
	if err != nil {
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
		return ""
	}

	var address string
	row.Scan(&address)
	return address
}
