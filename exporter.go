package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const MAX_INT32 = 2147483647

type Table struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name         string `json:"name"`
	FieldType    string `json:"fieldType"` // N = Number, S = String, B = Boolean, D = DateTime
	IsForeignKey bool   `json:"isForeignKey"`
}

func getTableAndFieldInfo() []Table {
	tables := make([]Table, 0)

	sqlStatement := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND  schemaname != 'information_schema' ORDER BY tablename ASC`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return tables
	}

	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, Table{Name: name})
	}

	for i := 0; i < len(tables); i++ {
		// query the columns
		sqlStatement := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position ASC`
		rows, err := db.Query(sqlStatement, tables[i].Name)
		if err != nil {
			return tables
		}
		tables[i].Fields = make([]Field, 0)

		// create the field
		for rows.Next() {
			field := Field{}
			var column_name string
			var data_type string
			rows.Scan(&column_name, &data_type)

			var fieldType string
			switch data_type {
			case "smallint":
				fallthrough
			case "integer":
				fallthrough
			case "real":
				fallthrough
			case "bigint":
				fieldType = "N" // Number
			case "character":
				fallthrough
			case "text":
				fallthrough
			case "uuid":
				fallthrough
			case "character varying":
				fieldType = "S" // String
			case "boolean":
				fieldType = "B" // Boolean
			case "timestamp without time zone":
				fieldType = "D" // DateTime
			}

			field.Name = column_name
			field.FieldType = fieldType
			tables[i].Fields = append(tables[i].Fields, field)
		}

		// query the foreign keys
		sqlStatement = `SELECT kcu.column_name FROM information_schema.table_constraints AS tc JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1`
		rows, err = db.Query(sqlStatement, tables[i].Name)
		if err != nil {
			return tables
		}

		// set the fields as foreign keys
		for rows.Next() {
			var column_name string
			rows.Scan(&column_name)

			for j := 0; j < len(tables[i].Fields); j++ {
				if tables[i].Fields[j].Name == column_name {
					tables[i].Fields[j].IsForeignKey = true
					break
				}
			}
		}
	}

	return tables
}

type ExportInfo struct {
	Table     string            `json:"table"`
	Separator string            `json:"separator"`
	NewLine   string            `json:"newLine"`
	Fields    []ExportInfoField `json:"fields"`
}

type ExportInfoField struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
}

func (i *ExportInfo) isValid() bool {
	return !(len(i.Separator) == 0 || (i.NewLine != "CRLF" && i.NewLine != "CR" && i.NewLine != "LF") || len(i.Fields) == 0)
}

func (f *ExportInfoField) isValid() bool {
	return !(len(f.Name) == 0 || (f.Relation != "" && f.Relation != "I" && f.Relation != "N"))
}

func (e *ExportInfo) export() string {
	// validate the data sent from the web
	if !e.isValid() {
		return ""
	}
	for i := 0; i < len(e.Fields); i++ {
		if !e.Fields[i].isValid() {
			return ""
		}
	}

	newLineCharacter := ""
	if e.NewLine == "CRLF" {
		newLineCharacter = "\r\n"
	} else if e.NewLine == "CR" {
		newLineCharacter = "\r"
	} else if e.NewLine == "LF" {
		newLineCharacter = "\n"
	}

	// We can't send to PostgreSQL the columns and tables in parameters, for example it won't accept 'SELECT * FROM $1' with 'product' as first parameter.
	// Validate that the table name and the fields exists in the database to prevent an SQL injection when concatenating the front-end strings.
	tableInfo := getTableAndFieldInfo()
	var ok bool = false

	for i := 0; i < len(tableInfo); i++ {
		if tableInfo[i].Name == e.Table {
			// table found :) check that all the fields exist
			for j := 0; j < len(e.Fields); j++ {
				ok = false
				for k := 0; k < len(tableInfo[i].Fields); k++ {
					if e.Fields[j].Name == tableInfo[i].Fields[k].Name {
						ok = true
						break
					}
				}
				if !ok {
					break
				}
			}
			break
		}
	}
	if !ok {
		return ""
	}

	// build a select statement that will retrieve the selected columns of all the table
	sqlStatement := `SELECT `
	for i := 0; i < len(e.Fields); i++ {
		if i > 0 {
			sqlStatement += `,`
		}
		fieldName := e.Fields[i].Name
		if e.Fields[i].Name == "order" {
			e.Fields[i].Name = "\"order\""
		}

		if e.Fields[i].Relation == "" || e.Fields[i].Relation == "I" {
			sqlStatement += e.Fields[i].Name
		} else if e.Fields[i].Relation == "N" {
			foreignTableName := findForeignTableName(e.Table, fieldName)

			columnName := "name"
			if foreignTableName == "address" {
				columnName = "address"
			} else if foreignTableName == "manufacturing_order" {
				columnName = "uuid"
			} else if foreignTableName == "packaging" {
				columnName = "id"
			} else if foreignTableName == "purchase_delivery_note" || foreignTableName == "sales_delivery_note" {
				columnName = "delivery_note_name"
			} else if foreignTableName == "purchase_invoice" || foreignTableName == "sales_invoice" {
				columnName = "invoice_name"
			} else if foreignTableName == "purchase_order" || foreignTableName == "sales_order" {
				columnName = "order_name"
			} else if foreignTableName == "purchase_order_detail" || foreignTableName == "warehouse_movement" {
				columnName = "id"
			}

			sqlStatement += `(SELECT ` + columnName + ` FROM ` + foreignTableName + ` WHERE ` + foreignTableName + `.id=` + e.Table + `.` + e.Fields[i].Name + `)`
		}
	}
	sqlStatement += ` FROM ` + e.Table

	// execute query, prepare arrays to store the result with all the columns as strings
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return ""
	}

	fields := make([]*string, len(e.Fields))
	data := make([]interface{}, len(e.Fields))

	for i := 0; i < len(e.Fields); i++ {
		data[i] = &fields[i]
	}

	// save line by line to a file
	fileId := uuid.New().String()
	f, err := os.OpenFile("./exported/"+fileId+".csv", os.O_CREATE, 0700)
	if err != nil {
		return ""
	}

	// build the file with the results
	for rows.Next() {
		rows.Scan(data...)

		line := ""
		for i := 0; i < len(e.Fields); i++ {
			if i > 0 {
				line += e.Separator
			}
			if fields[i] != nil {
				line += *fields[i]
			}
		}
		f.WriteString(line + newLineCharacter)
	}
	f.Close()
	return fileId
}

func findForeignTableName(tableName string, columnName string) string {
	sqlStatement := `SELECT ccu.table_name AS foreign_table_name FROM information_schema.table_constraints AS tc JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name=$1 AND kcu.column_name=$2`
	row := db.QueryRow(sqlStatement, tableName, columnName)
	if row.Err() != nil {
		return ""
	}

	var foreignTableName string
	row.Scan(&foreignTableName)
	return foreignTableName
}

func exportToJSON(tableName string, enterpriseId int32) string {
	tableInfo := getTableAndFieldInfo()
	var ok bool = false

	for i := 0; i < len(tableInfo); i++ {
		if tableInfo[i].Name == tableName {
			ok = true
			break
		}
	}

	if !ok {
		return ""
	}

	fileId := uuid.New().String()
	f, err := os.OpenFile("./exported/"+fileId+".json", os.O_CREATE, 0700)
	if err != nil {
		return ""
	}

	var data []byte
	switch tableName {
	case "billing_series":
		data, _ = json.Marshal(getBillingSeries(enterpriseId))
	case "carrier":
		data, _ = json.Marshal(getCariers(enterpriseId))
	case "color":
		data, _ = json.Marshal(getColor(enterpriseId))
	case "config":
		data, _ = json.Marshal(getSettingsRecordById(enterpriseId))
	case "country":
		data, _ = json.Marshal(getCountries(enterpriseId))
	case "currency":
		data, _ = json.Marshal(getCurrencies(enterpriseId))
	case "document":
		data, _ = json.Marshal(getDocuments(enterpriseId))
	case "document_container":
		data, _ = json.Marshal(getDocumentContainer(enterpriseId))
	case "group":
		data, _ = json.Marshal(getGroup(enterpriseId))
	case "incoterm":
		data, _ = json.Marshal(getIncoterm(enterpriseId))
	case "language":
		data, _ = json.Marshal(getLanguages(enterpriseId))
	case "manufacturing_order":
		q := ManufacturingPaginationQuery{PaginationQuery: PaginationQuery{Offset: 0, Limit: MAX_INT32}, OrderTypeId: 0}
		data, _ = json.Marshal(q.getManufacturingOrder(enterpriseId))
	case "manufacturing_order_type":
		data, _ = json.Marshal(getManufacturingOrderType(enterpriseId))
	case "packages":
		data, _ = json.Marshal(getPackages(enterpriseId))
	case "payment_method":
		data, _ = json.Marshal(getPaymentMethods(enterpriseId))
	case "product":
		data, _ = json.Marshal(getProduct(enterpriseId))
	case "product_family":
		data, _ = json.Marshal(getProductFamilies(enterpriseId))
	case "purchase_delivery_note":
		data, _ = json.Marshal(getPurchaseDeliveryNotes(enterpriseId))
	case "purchase_invoice":
		data, _ = json.Marshal(getPurchaseInvoices(enterpriseId))
	case "purchase_order":
		data, _ = json.Marshal(getPurchaseOrder(enterpriseId))
	case "shipping":
		data, _ = json.Marshal(getShippings(enterpriseId))
	case "state":
		data, _ = json.Marshal(getStates(enterpriseId))
	case "suppliers":
		data, _ = json.Marshal(getSuppliers(enterpriseId))
	case "user":
		data, _ = json.Marshal(getUser(enterpriseId))
	case "warehouse":
		data, _ = json.Marshal(getWarehouses(enterpriseId))
	default:
		return ""
	}
	f.Write(data)
	f.Close()
	return fileId
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	uuidCsv := r.URL.Query()["uuid_csv"]
	uuidJson := r.URL.Query()["uuid_json"]
	if (len(uuidCsv) > 0 && len(uuidCsv[0]) != 0 && len(uuidCsv[0]) != 36) || (len(uuidJson) > 0 && len(uuidJson[0]) != 0 && len(uuidJson[0]) != 36) || (len(uuidCsv) == 0 && len(uuidJson) == 0) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(uuidCsv) > 0 {
		content, err := ioutil.ReadFile("./exported/" + uuidCsv[0] + ".csv")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+uuidCsv[0]+".csv")
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
		os.Remove("./exported/" + uuidCsv[0] + ".csv")
	} else {
		content, err := ioutil.ReadFile("./exported/" + uuidJson[0] + ".json")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+uuidJson[0]+".json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		os.Remove("./exported/" + uuidJson[0] + ".json")
	}
}

type ImportInfo struct {
	JsonData  string `json:"jsonData"`
	TableName string `json:"tableName"`
}

func (f *ImportInfo) importJson(enterpriseId int32) bool {
	if len(f.JsonData) == 0 || len(f.TableName) == 0 {
		return false
	}

	tableInfo := getTableAndFieldInfo()
	var ok bool = false

	for i := 0; i < len(tableInfo); i++ {
		if tableInfo[i].Name == f.TableName {
			ok = true
			break
		}
	}

	if !ok {
		return false
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false
	}
	///

	var jsonData []byte = []byte(f.JsonData)
	switch f.TableName {
	case "address":
		var address []Address
		json.Unmarshal(jsonData, &address)
		for i := 0; i < len(address); i++ {
			address[i].enterprise = enterpriseId
			ok = address[i].insertAddress()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "billing_series":
		var serie []BillingSerie
		json.Unmarshal(jsonData, &serie)
		for i := 0; i < len(serie); i++ {
			serie[i].enterprise = enterpriseId
			ok = serie[i].insertBillingSerie()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "carrier":
		var carrier []Carrier
		json.Unmarshal(jsonData, &carrier)
		for i := 0; i < len(carrier); i++ {
			carrier[i].enterprise = enterpriseId
			ok = carrier[i].insertCarrier()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "color":
		var color []Color
		json.Unmarshal(jsonData, &color)
		for i := 0; i < len(color); i++ {
			color[i].enterprise = enterpriseId
			ok = color[i].insertColor()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "config":
		var settings Settings
		json.Unmarshal(jsonData, &settings)
		settings.Id = enterpriseId
		ok = settings.updateSettingsRecord()
		if !ok {
			trans.Rollback()
			return false
		}
	case "country":
		var country []Country
		json.Unmarshal(jsonData, &country)
		for i := 0; i < len(country); i++ {
			country[i].enterprise = enterpriseId
			ok = country[i].insertCountry()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "currency":
		var currency []Currency
		json.Unmarshal(jsonData, &currency)
		for i := 0; i < len(currency); i++ {
			currency[i].enterprise = enterpriseId
			ok = currency[i].insertCurrency()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "customer":
		var customer []Customer
		json.Unmarshal(jsonData, &customer)
		for i := 0; i < len(customer); i++ {
			customer[i].enterprise = enterpriseId
			ok, _ = customer[i].insertCustomer()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "document_container":
		var documentContainer []DocumentContainer
		json.Unmarshal(jsonData, &documentContainer)
		for i := 0; i < len(documentContainer); i++ {
			documentContainer[i].enterprise = enterpriseId
			ok = documentContainer[i].insertDocumentContainer()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "group":
		var group []Group
		json.Unmarshal(jsonData, &group)
		for i := 0; i < len(group); i++ {
			group[i].enterprise = enterpriseId
			ok = group[i].insertGroup()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "incoterm":
		var incoterm []Incoterm
		json.Unmarshal(jsonData, &incoterm)
		for i := 0; i < len(incoterm); i++ {
			incoterm[i].enterprise = enterpriseId
			ok = incoterm[i].insertIncoterm()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "language":
		var language []Language
		json.Unmarshal(jsonData, &language)
		for i := 0; i < len(language); i++ {
			language[i].enterprise = enterpriseId
			ok = language[i].insertLanguage()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "manufacturing_order":
		var manufacturingOrder []ManufacturingOrder
		json.Unmarshal(jsonData, &manufacturingOrder)
		for i := 0; i < len(manufacturingOrder); i++ {
			manufacturingOrder[i].UserCreated = 1
			manufacturingOrder[i].enterprise = enterpriseId
			ok = manufacturingOrder[i].insertManufacturingOrder()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "manufacturing_order_type":
		var manufacturingOrderType []ManufacturingOrderType
		json.Unmarshal(jsonData, &manufacturingOrderType)
		for i := 0; i < len(manufacturingOrderType); i++ {
			manufacturingOrderType[i].enterprise = enterpriseId
			ok = manufacturingOrderType[i].insertManufacturingOrderType()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "packages":
		var packages []Packages
		json.Unmarshal(jsonData, &packages)
		for i := 0; i < len(packages); i++ {
			packages[i].enterprise = enterpriseId
			ok = packages[i].insertPackage()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "payment_method":
		var paymentMethod []PaymentMethod
		json.Unmarshal(jsonData, &paymentMethod)
		for i := 0; i < len(paymentMethod); i++ {
			paymentMethod[i].enterprise = enterpriseId
			ok = paymentMethod[i].insertPaymentMethod()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "product":
		var product []Product
		json.Unmarshal(jsonData, &product)
		for i := 0; i < len(product); i++ {
			product[i].enterprise = enterpriseId
			ok = product[i].insertProduct()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "product_family":
		var productFamily []ProductFamily
		json.Unmarshal(jsonData, &productFamily)
		for i := 0; i < len(productFamily); i++ {
			productFamily[i].enterprise = enterpriseId
			ok = productFamily[i].insertProductFamily()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "purchase_delivery_note":
		var purchaseDeliveryNote []PurchaseDeliveryNote
		json.Unmarshal(jsonData, &purchaseDeliveryNote)
		for i := 0; i < len(purchaseDeliveryNote); i++ {
			purchaseDeliveryNote[i].enterprise = enterpriseId
			ok, _ = purchaseDeliveryNote[i].insertPurchaseDeliveryNotes()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "purchase_invoice":
		var purchaseInvoice []PurchaseInvoice
		json.Unmarshal(jsonData, &purchaseInvoice)
		for i := 0; i < len(purchaseInvoice); i++ {
			purchaseInvoice[i].enterprise = enterpriseId
			ok, _ = purchaseInvoice[i].insertPurchaseInvoice()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "purchase_order":
		var purchaseOrder []PurchaseOrder
		json.Unmarshal(jsonData, &purchaseOrder)
		for i := 0; i < len(purchaseOrder); i++ {
			purchaseOrder[i].enterprise = enterpriseId
			ok, _ = purchaseOrder[i].insertPurchaseOrder()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "sales_delivery_note":
		var salesDeliveryNote []SalesDeliveryNote
		json.Unmarshal(jsonData, &salesDeliveryNote)
		for i := 0; i < len(salesDeliveryNote); i++ {
			salesDeliveryNote[i].enterprise = enterpriseId
			ok, _ = salesDeliveryNote[i].insertSalesDeliveryNotes()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "sales_invoice":
		var saleInvoice []SalesInvoice
		json.Unmarshal(jsonData, &saleInvoice)
		for i := 0; i < len(saleInvoice); i++ {
			saleInvoice[i].enterprise = enterpriseId
			ok, _ = saleInvoice[i].insertSalesInvoice()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "sales_order":
		var saleOrder []SaleOrder
		json.Unmarshal(jsonData, &saleOrder)
		for i := 0; i < len(saleOrder); i++ {
			saleOrder[i].enterprise = enterpriseId
			ok, _ = saleOrder[i].insertSalesOrder()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "shipping":
		var shipping []Shipping
		json.Unmarshal(jsonData, &shipping)
		for i := 0; i < len(shipping); i++ {
			shipping[i].enterprise = enterpriseId
			ok, _ = shipping[i].insertShipping()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "state":
		var state []State
		json.Unmarshal(jsonData, &state)
		for i := 0; i < len(state); i++ {
			state[i].enterprise = enterpriseId
			ok = state[i].insertState()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "suppliers":
		var supplier []Supplier
		json.Unmarshal(jsonData, &supplier)
		for i := 0; i < len(supplier); i++ {
			supplier[i].enterprise = enterpriseId
			ok = supplier[i].insertSupplier()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "warehouse":
		var warehouse []Warehouse
		json.Unmarshal(jsonData, &warehouse)
		for i := 0; i < len(warehouse); i++ {
			warehouse[i].enterprise = enterpriseId
			ok = warehouse[i].insertWarehouse()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	case "warehouse_movement":
		var warehouseMovement []WarehouseMovement
		json.Unmarshal(jsonData, &warehouseMovement)
		for i := 0; i < len(warehouseMovement); i++ {
			warehouseMovement[i].enterprise = enterpriseId
			ok = warehouseMovement[i].insertWarehouseMovement()
			if !ok {
				trans.Rollback()
				return false
			}
		}
	default:
		trans.Rollback()
		return false
	}

	///
	err = trans.Commit()
	return err == nil
	///
}
