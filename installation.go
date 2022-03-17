package main

import (
	"fmt"
)

// This file contains functions that will help in the deploy of this software and later maintenance.

// The repository where is software is located always contains a file named db.sql, that contains the database schema without data.
// Also contains a file named update.sql, that contains the SQL code necessary to update the schema from the prior version.
// This function returns true if the database already exists, or if it has been installed or updated successfully.
// Returns false if the database could not be created or updated.
func installDB() bool {
	/*// Does the database have tables? Or is it empty?
	sqlStatement := `SELECT COUNT(*) FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND  schemaname != 'information_schema'`
	row := db.QueryRow(sqlStatement)
	if row.Err() != nil {
		return false
	}

	var tables int32
	row.Scan(&tables)

	if tables == 0 {
		content, err := ioutil.ReadFile("db.sql")
		if err != nil {
			fmt.Println("Count not read file bd.sql", err)
			return false
		}

		_, err = db.Exec(string(content))
		if err != nil {
			fmt.Println("Could not copy database schema", err)
			return false
		}

		// truncate the file on successfull update
		updateFile, _ := os.OpenFile("update.sql", os.O_RDWR, 0666)
		updateFile.Truncate(0)
	} else {
		content, err := ioutil.ReadFile("update.sql")
		if err != nil {
			fmt.Println("Could not read file update.sql", err)
			return false
		}

		// there is no pending updates
		if len(content) == 0 {
			return true
		}

		_, err = db.Exec(string(content))
		if err != nil {
			fmt.Println("Could not update database schema", err)
			return false
		}

		// truncate the file on successfull update
		updateFile, _ := os.OpenFile("update.sql", os.O_RDWR, 0666)
		updateFile.Truncate(0)
	}
	return true*/

	fmt.Println("Upgrading the database schema...")
	ok := upradeDataBaseSchema()
	if ok {
		fmt.Println("Upgrade OK. Database ready!")
	} else {
		fmt.Println("There was an error upgrading the database schema. The SQL error returned by PostgreSQL is shown above this message. The server will exit now with error code 1.")
	}
	return ok
}
