package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	Id                int32  `json:"id"`
	Customer          *int32 `json:"customer"`
	Supplier          *int32 `json:"supplier"`
	Address           string `json:"address"`
	Address2          string `json:"address2"`
	City              string `json:"city"`
	State             *int32 `json:"state"`
	Country           int16  `json:"country"`
	PrivateOrBusiness string `json:"privateOrBusiness"`
	Notes             string `json:"notes"`
	PrestaShopId      int32  `json:"prestaShopId"`
	ZipCode           string `json:"zipCode"`
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
		rows.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode)
		addresses = append(addresses, a)
	}

	return addresses
}

func getAddressRow(addressId int32) Address {
	sqlStatement := `SELECT * FROM address WHERE id=$1`
	row := db.QueryRow(sqlStatement, addressId)
	if row.Err() != nil {
		return Address{}
	}

	a := Address{}
	row.Scan(&a.Id, &a.Customer, &a.Address, &a.Address2, &a.State, &a.City, &a.Country, &a.PrivateOrBusiness, &a.Notes, &a.Supplier, &a.PrestaShopId, &a.ZipCode)

	return a
}

func (a *Address) isValid() bool {
	return !((a.Customer == nil && a.Supplier == nil) || (a.Customer != nil && *a.Customer <= 0) || (a.Supplier != nil && *a.Supplier <= 0) || len(a.Address) == 0 || len(a.Address) > 200 || len(a.Address2) > 200 || len(a.City) == 0 || len(a.City) > 100 || a.Country <= 0 || (a.PrivateOrBusiness != "P" && a.PrivateOrBusiness != "B" && a.PrivateOrBusiness != "_") || len(a.Notes) > 1000 || len(a.ZipCode) > 12)
}

func (a *Address) insertAddress() bool {
	if !a.isValid() {
		fmt.Println("INVALID")
		data, _ := json.Marshal(a)
		fmt.Println(string(data))
		return false
	}

	sqlStatement := `INSERT INTO address(customer, address, address_2, city, state, country, private_business, notes, supplier, ps_id, zip_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	res, err := db.Exec(sqlStatement, a.Customer, a.Address, a.Address2, a.City, a.State, a.Country, a.PrivateOrBusiness, a.Notes, a.Supplier, a.PrestaShopId, a.ZipCode)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, _ := res.RowsAffected()
	return rows > 0
}

func (a *Address) updateAddress() bool {
	if a.Id <= 0 || !a.isValid() {
		return false
	}

	sqlStatement := `UPDATE address SET customer=$2, address=$3, address_2=$4, city=$5, state=$6, country=$7, private_business=$8, notes=$9, supplier=$10, zip_code=$11 WHERE id = $1`
	res, err := db.Exec(sqlStatement, a.Id, a.Customer, a.Address, a.Address2, a.City, a.State, a.Country, a.PrivateOrBusiness, a.Notes, a.Supplier, a.ZipCode)
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

func locateAddressByCustomer(customerId int32) []AddressLocate {
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

func locateAddressBySupplier(supplierId int32) []AddressLocate {
	var addresses []AddressLocate = make([]AddressLocate, 0)
	sqlStatement := `SELECT id, address FROM address WHERE supplier=$1 ORDER BY id ASC`
	rows, err := db.Query(sqlStatement, supplierId)
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
