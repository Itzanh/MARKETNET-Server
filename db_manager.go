package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type DataBase struct {
	Tables []DBTable `json:"tables"`
}

type DBTable struct {
	Name        string         `json:"name"`
	Fields      []DBField      `json:"fields"`
	Indexes     []DBIndex      `json:"indexes"`
	Constraints []DBConstraint `json:"constraints"`
}

type DBField struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Length  int     `json:"length"`
	NotNULL bool    `json:"notNULL"`
	Default *string `json:"default"`
}

type DBIndex struct {
	Name     string `json:"name"`
	IndexDef string `json:"indexDef"`
}

type DBConstraint struct {
	Name          string `json:"name"`
	ConstraintDef string `json:"constraintDef"`
}

func (f *DBField) isEquals(other DBField) bool {
	if f.Name != other.Name {
		return false
	}
	return f.isEqualsType(other) && f.isEqualsNotNull(other) && f.isEqualsDefault(other)
}

func (f *DBField) isEqualsType(other DBField) bool {
	return f.Type == other.Type && f.Length == other.Length
}

func (f *DBField) isEqualsNotNull(other DBField) bool {
	return f.NotNULL == other.NotNULL
}

func (f *DBField) isEqualsDefault(other DBField) bool {
	if (f.Default == nil) != (other.Default == nil) {
		return false
	}
	if (f.Default != nil) && (other.Default != nil) {
		if *f.Default != *other.Default {
			return false
		}
	}
	return true
}

func (f *DBField) toString() string {
	if f.Length != 0 {
		if f.Type == "timestamp with time zone" {
			return `timestamp(` + strconv.Itoa(f.Length) + `) with time zone`
		} else if f.Type == "timestamp without time zone" {
			return `timestamp(` + strconv.Itoa(f.Length) + `) without time zone`
		} else {
			return f.Type + `(` + strconv.Itoa(f.Length) + `)`
		}
	} else {
		return f.Type
	}
}

// DEVELOPMENT ONLY
func generateSchemaJson() {
	var dbSchema DataBase = DataBase{}

	sqlStatement := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND  schemaname != 'information_schema'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log("Initialization", err.Error())
		return
	}

	var tableName string
	for rows.Next() {
		rows.Scan(&tableName)

		t := DBTable{
			Name: tableName,
		}

		// GET COLUMNS

		sqlStatement = `SELECT column_name,column_default,is_nullable,data_type,character_maximum_length,datetime_precision FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1 ORDER BY ordinal_position ASC`
		rowsFields, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			log("Initialization", err.Error())
			return
		}

		var isNullable string
		var charcterMaximumLength *int
		var datetimePrecision *int
		for rowsFields.Next() {
			f := DBField{}
			rowsFields.Scan(&f.Name, &f.Default, &isNullable, &f.Type, &charcterMaximumLength, &datetimePrecision)
			f.NotNULL = isNullable == "NO"
			if charcterMaximumLength != nil && *charcterMaximumLength != 0 {
				f.Length = *charcterMaximumLength
			} else if datetimePrecision != nil && *datetimePrecision != 0 {
				f.Length = *datetimePrecision
			}
			t.Fields = append(t.Fields, f)
		}

		// GET INDEXES

		sqlStatement := `SELECT indexname,indexdef FROM pg_indexes WHERE tablename = $1`
		rowsIndexes, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			log("Initialization", err.Error())
			return
		}

		for rowsIndexes.Next() {
			i := DBIndex{}
			rowsIndexes.Scan(&i.Name, &i.IndexDef)
			t.Indexes = append(t.Indexes, i)
		}

		// GET CONSTRAINTS

		sqlStatement = `SELECT DISTINCT pgc.conname AS constraint_name,pg_get_constraintdef(pgc.oid) FROM pg_constraint pgc JOIN pg_namespace nsp ON nsp.oid = pgc.connamespace JOIN pg_class  cls ON pgc.conrelid = cls.oid LEFT JOIN information_schema.constraint_column_usage ccu ON pgc.conname = ccu.constraint_name AND nsp.nspname = ccu.constraint_schema WHERE relname = $1 ORDER BY pgc.conname`
		rowsConstraints, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			log("Initialization", err.Error())
			return
		}

		for rowsConstraints.Next() {
			c := DBConstraint{}
			rowsConstraints.Scan(&c.Name, &c.ConstraintDef)
			t.Constraints = append(t.Constraints, c)
		}

		dbSchema.Tables = append(dbSchema.Tables, t)
	}

	data, err := json.Marshal(dbSchema)
	if err != nil {
		log("Initialization", err.Error())
		return
	}
	err = ioutil.WriteFile("schema.json", data, 0700)
	if err != nil {
		log("Initialization", err.Error())
		return
	}
}

func upradeDataBaseSchema() bool {
	data, err := ioutil.ReadFile("schema.json")
	if err != nil {
		fmt.Println(err)
		return false
	}
	var dbSchema DataBase = DataBase{}
	err = json.Unmarshal(data, &dbSchema)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var t DBTable
	var tableCount int
	for i := 0; i < len(dbSchema.Tables); i++ { // TABLES
		t = dbSchema.Tables[i]

		sqlStatement := `SELECT COUNT(*) FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND  schemaname != 'information_schema' AND tablename = $1`
		row := db.QueryRow(sqlStatement, t.Name)
		if row.Err() != nil {
			fmt.Println(row.Err())
			return false
		}
		row.Scan(&tableCount)

		// CREATE TABLE

		if tableCount == 0 {
			sqlStatement = `CREATE TABLE public.` + t.Name + `()`
			_, err := db.Exec(sqlStatement)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}

		// CHECK COLUMNS

		sqlStatement = `SELECT column_name,column_default,is_nullable,data_type,character_maximum_length,datetime_precision FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1 ORDER BY ordinal_position ASC`
		rowsFields, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			fmt.Println(err)
			return false
		}

		var fieldsNow []DBField = make([]DBField, 0)
		var isNullable string
		var charcterMaximumLength *int
		var datetimePrecision *int
		for rowsFields.Next() {
			f := DBField{}
			rowsFields.Scan(&f.Name, &f.Default, &isNullable, &f.Type, &charcterMaximumLength, &datetimePrecision)
			f.NotNULL = isNullable == "NO"
			if charcterMaximumLength != nil && *charcterMaximumLength != 0 {
				f.Length = *charcterMaximumLength
			} else if datetimePrecision != nil && *datetimePrecision != 0 {
				f.Length = *datetimePrecision
			}
			fieldsNow = append(fieldsNow, f)
		}

		// check if there are fields in schema.json that are not present in the current database schema
		// and update the definition if it is different in schema.json
		var found bool
		var f DBField
		var fOld DBField
		for j := 0; j < len(t.Fields); j++ { // FIELDS
			f = t.Fields[j]
			found = false
			for k := 0; k < len(fieldsNow); k++ {
				if f.Name == fieldsNow[k].Name {
					found = true
					fOld = fieldsNow[k]
					break
				}
			}
			if !found {
				// CREATE FIELD
				sqlStatement = `ALTER TABLE public.` + t.Name + ` ADD COLUMN ` + f.Name + ` `
				sqlStatement += f.toString()
				if f.NotNULL {
					sqlStatement += ` NOT NULL`
				}
				if f.Default != nil {
					sqlStatement += ` DEFAULT ` + *f.Default
				}
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
			} else if !f.isEquals(fOld) {
				// UPDATE FIELD
				if !f.isEqualsType(fOld) {
					sqlStatement = `ALTER TABLE public.` + t.Name + ` ALTER COLUMN ` + f.Name + ` TYPE ` + f.toString()
					_, err := db.Exec(sqlStatement)
					if err != nil {
						fmt.Println(err)
						return false
					}
				} // if !f.isEqualsType(fOld) {
				if !f.isEqualsNotNull(fOld) {
					if f.NotNULL && !fOld.NotNULL {
						sqlStatement = `ALTER TABLE public.` + t.Name + ` ALTER COLUMN ` + f.Name + ` SET NOT NULL`
						_, err := db.Exec(sqlStatement)
						if err != nil {
							fmt.Println(err)
							return false
						}
					} else {
						sqlStatement = `ALTER TABLE public.` + t.Name + ` ALTER COLUMN ` + f.Name + ` DROP NOT NULL`
						_, err := db.Exec(sqlStatement)
						if err != nil {
							fmt.Println(err)
							return false
						}
					}
				} // if !f.isEqualsNotNull(fOld) {
				if !f.isEqualsDefault(fOld) {
					sqlStatement = `ALTER TABLE public.` + t.Name + ` ALTER COLUMN ` + f.Name + ` SET DEFAULT ` + *f.Default
					_, err := db.Exec(sqlStatement)
					if err != nil {
						fmt.Println(err)
						return false
					}
				} // if !f.isEqualsDefault(fOld) {
			}
		} // for j := 0; j < len(t.Fields); j++ { // FIELDS

		// delete fields that are present in the current database schema but are not present in schema.json
		for j := 0; j < len(fieldsNow); j++ { // FIELDS
			found = false
			f = fieldsNow[j]
			for k := 0; k < len(t.Fields); k++ {
				if f.Name == t.Fields[k].Name {
					found = true
					break
				}
			}
			if !found {
				sqlStatement = `ALTER TABLE public.` + t.Name + ` DROP COLUMN ` + f.Name
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
			}
		} // for j := 0; j < len(fieldsNow); j++ { // FIELDS

		// CHECK INDEXES

		sqlStatement = `SELECT indexname,indexdef FROM pg_indexes WHERE tablename = $1`
		rowsIndexes, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			fmt.Println(err)
			return false
		}

		var indexesNow []DBIndex = make([]DBIndex, 0)
		for rowsIndexes.Next() {
			i := DBIndex{}
			rowsIndexes.Scan(&i.Name, &i.IndexDef)
			indexesNow = append(indexesNow, i)
		}

		var indexDef string
		var idx DBIndex
		for i := 0; i < len(t.Indexes); i++ { // INDEXES
			idx = t.Indexes[i]
			found = false
			indexDef = ""
			for k := 0; k < len(indexesNow); k++ {
				if idx.Name == indexesNow[k].Name {
					found = true
					indexDef = indexesNow[k].IndexDef
					break
				}
			}
			if !found {
				_, err := db.Exec(idx.IndexDef)
				if err != nil {
					fmt.Println(err)
					return false
				}
			} else if idx.IndexDef != indexDef {
				sqlStatement = `DROP INDEX public.` + idx.Name
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
				_, err = db.Exec(idx.IndexDef)
				if err != nil {
					fmt.Println(err)
					return false
				}
			}
		} // for i := 0; i < len(t.Indexes); i++ { // INDEXES
		for i := 0; i < len(indexesNow); i++ { // INDEXES
			idx = indexesNow[i]
			found = false
			for k := 0; k < len(t.Indexes); k++ {
				if idx.Name == t.Indexes[k].Name {
					found = true
					break
				}
			}
			if !found {
				sqlStatement = `DROP INDEX public.` + idx.Name
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
			}
		} // for i := 0; i < len(indexesNow); i++ { // INDEXES

		// CHECK CONSTRAINTS

		sqlStatement = `SELECT DISTINCT pgc.conname AS constraint_name,pg_get_constraintdef(pgc.oid) FROM pg_constraint pgc JOIN pg_namespace nsp ON nsp.oid = pgc.connamespace JOIN pg_class  cls ON pgc.conrelid = cls.oid LEFT JOIN information_schema.constraint_column_usage ccu ON pgc.conname = ccu.constraint_name AND nsp.nspname = ccu.constraint_schema WHERE relname = $1 ORDER BY pgc.conname`
		rowsConstraints, err := db.Query(sqlStatement, t.Name)
		if err != nil {
			fmt.Println(err)
			return false
		}

		var constraintsNow []DBConstraint = make([]DBConstraint, 0)
		for rowsConstraints.Next() {
			c := DBConstraint{}
			rowsConstraints.Scan(&c.Name, &c.ConstraintDef)
			constraintsNow = append(constraintsNow, c)
		}

		var constraintDef string
		var c DBConstraint
		for i := 0; i < len(t.Constraints); i++ { // CONSTRAINTS
			c = t.Constraints[i]
			found = false
			constraintDef = ""
			for k := 0; k < len(constraintsNow); k++ {
				if c.Name == constraintsNow[k].Name {
					found = true
					constraintDef = constraintsNow[k].ConstraintDef
					break
				}
			}
			if !found {
				_, err := db.Exec(`ALTER TABLE public.` + t.Name + ` ADD CONSTRAINT ` + c.Name + ` ` + c.ConstraintDef)
				if err != nil {
					fmt.Println(err)
					return false
				}
			} else if c.ConstraintDef != constraintDef {
				sqlStatement = `ALTER TABLE public.` + t.Name + ` DROP CONSTRAINT ` + c.Name
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
				_, err = db.Exec(c.ConstraintDef)
				if err != nil {
					fmt.Println(err)
					return false
				}
			}
		} // for i := 0; i < len(t.Indexes); i++ { // CONSTRAINTS
		for i := 0; i < len(constraintsNow); i++ { // CONSTRAINTS
			c = constraintsNow[i]
			found = false
			for k := 0; k < len(t.Constraints); k++ {
				if c.Name == t.Constraints[k].Name {
					found = true
					break
				}
			}
			if !found {
				sqlStatement = `ALTER TABLE public.` + t.Name + ` DROP CONSTRAINT ` + c.Name
				_, err := db.Exec(sqlStatement)
				if err != nil {
					fmt.Println(err)
					return false
				}
			}
		} // for i := 0; i < len(indexesNow); i++ { // CONSTRAINTS

	} // for i := 0; i < len(dbSchema.Tables); i++ { // TABLES

	// DELETE TABLES THAT ARE NOT IN SCHEMA.JSON

	sqlStatement := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND  schemaname != 'information_schema'`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var tableName string
	var found bool
	for rows.Next() {
		rows.Scan(&tableName)

		found = false
		for i := 0; i < len(dbSchema.Tables); i++ { // TABLES
			if dbSchema.Tables[i].Name == tableName {
				found = true
				break
			}
		}
		if !found {
			sqlStatement = `DROP TABLE public.` + tableName
			_, err := db.Exec(sqlStatement)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}
	}

	return true
}
