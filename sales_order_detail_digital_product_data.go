package main

type SalesOrderDetailDigitalProductData struct {
	Id     int32  `json:"id"`
	Detail int64  `json:"detail"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func getSalesOrderDetailDigitalProductData(salesOrderDetailId int64, enterpriseId int32) []SalesOrderDetailDigitalProductData {
	productData := make([]SalesOrderDetailDigitalProductData, 0)

	detailRow := getSalesOrderDetailRow(salesOrderDetailId)
	if detailRow.enterprise != enterpriseId {
		return productData
	}

	sqlStatement := `SELECT * FROM public.sales_order_detail_digital_product_data WHERE detail = $1`
	rows, err := db.Query(sqlStatement, salesOrderDetailId)
	if err != nil {
		log("DB", err.Error())
		return productData
	}
	defer rows.Close()

	for rows.Next() {
		pd := SalesOrderDetailDigitalProductData{}
		rows.Scan(&pd.Id, &pd.Detail, &pd.Key, &pd.Value)
		productData = append(productData, pd)
	}

	return productData
}

func (d *SalesOrderDetailDigitalProductData) isValid() bool {
	return !(d.Detail <= 0 || len(d.Key) == 0 || len(d.Value) == 0)
}

func (d *SalesOrderDetailDigitalProductData) insertSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if !d.isValid() {
		return false
	}

	detailRow := getSalesOrderDetailRow(d.Detail)
	if detailRow.enterprise != enterpriseId || detailRow.Status != "E" {
		return false
	}
	productRow := getProductRow(detailRow.Product)
	if !productRow.DigitalProduct {
		return false
	}

	sqlStatement := `INSERT INTO public.sales_order_detail_digital_product_data(detail, key, value) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, d.Detail, d.Key, d.Value)
	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (d *SalesOrderDetailDigitalProductData) updateSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if !d.isValid() || d.Id <= 0 {
		return false
	}

	detailRow := getSalesOrderDetailRow(d.Detail)
	if detailRow.enterprise != enterpriseId || detailRow.Status != "E" {
		return false
	}

	sqlStatement := `UPDATE public.sales_order_detail_digital_product_data SET detail=$2, key=$3, value=$4 WHERE id=$1`
	_, err := db.Exec(sqlStatement, d.Id, d.Detail, d.Key, d.Value)
	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

func (d *SalesOrderDetailDigitalProductData) deleteSalesOrderDetailDigitalProductData(enterpriseId int32) bool {
	if d.Id <= 0 {
		return false
	}

	sqlStatement := `DELETE FROM public.sales_order_detail_digital_product_data WHERE id=$1 AND (SELECT enterprise FROM sales_order_detail WHERE sales_order_detail.id=sales_order_detail_digital_product_data.detail)=$2 AND (SELECT status FROM sales_order_detail WHERE sales_order_detail.id=sales_order_detail_digital_product_data.detail)='E'`
	_, err := db.Exec(sqlStatement, d.Id, enterpriseId)
	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}

type SetDigitalSalesOrderDetailAsSent struct {
	Detail                 int64  `json:"detail"`
	SendEmail              bool   `json:"sendEmail"`
	DestinationAddress     string `json:"destinationAddress"`
	DestinationAddressName string `json:"destinationAddressName"`
	Subject                string `json:"subject"`
}

func (data *SetDigitalSalesOrderDetailAsSent) setDigitalSalesOrderDetailAsSent(enterpriseId int32) bool {
	detail := getSalesOrderDetailRow(data.Detail)
	if detail.enterprise != enterpriseId || detail.Status != "E" {
		return false
	}
	digitalProductData := getSalesOrderDetailDigitalProductData(detail.Id, enterpriseId)
	if len(digitalProductData) == 0 {
		return false
	}

	if data.SendEmail {
		ei := EmailInfo{
			DestinationAddress:     data.DestinationAddress,
			DestinationAddressName: data.DestinationAddressName,
			Subject:                data.Subject,
			ReportId:               "SALES_ORDER_DIGITAL_PRODUCT_DATA",
			ReportDataId:           int32(detail.Order),
		}
		ei.sendEmail(enterpriseId)
	}

	sqlStatement := `UPDATE sales_order_detail SET status='G' WHERE id=$1`
	_, err := db.Exec(sqlStatement, data.Detail)
	if err != nil {
		log("DB", err.Error())
	}

	return err == nil
}
