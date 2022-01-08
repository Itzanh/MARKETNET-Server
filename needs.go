package main

import (
	"sort"
)

type Need struct {
	Product      int32   `json:"product"`
	ProductName  string  `json:"productName"`
	SupplierName *string `json:"supplierName"`
	Quantity     int32   `json:"quantity"`
}

func getNeeds(enterpriseId int32) []Need {
	var needs []Need = make([]Need, 0)
	sqlStatement := `SELECT product,(SELECT name FROM product WHERE product.id=sales_order_detail.product),(SELECT name FROM suppliers WHERE suppliers.id=(SELECT supplier FROM product WHERE product.id=sales_order_detail.product)),SUM(quantity) FROM sales_order_detail WHERE status='A' AND enterprise=$1 GROUP BY product`
	rows, err := db.Query(sqlStatement, enterpriseId)
	if err != nil {
		log("DB", err.Error())
		return needs
	}
	defer rows.Close()

	for rows.Next() {
		n := Need{}
		rows.Scan(&n.Product, &n.ProductName, &n.SupplierName, &n.Quantity)
		needs = append(needs, n)
	}

	return needs
}

func getNeedRow(productId int32) int32 {
	sqlStatement := `SELECT SUM(quantity) FROM sales_order_detail WHERE status='A' AND product=$1 GROUP BY product`
	row := db.QueryRow(sqlStatement, productId)
	if row.Err() != nil {
		log("DB", row.Err().Error())
		return 0
	}

	var quantity int32
	row.Scan(&quantity)
	return quantity
}

type PurchaseNeed struct {
	ProductId int32 `json:"product"`
	Quantity  int32 `json:"quantity"`
	product   Product
	supplier  Supplier
}

type PurchaseNeeds []PurchaseNeed

func (n PurchaseNeeds) Len() int {
	return len(n)
}
func (n PurchaseNeeds) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
func (n PurchaseNeeds) Less(i, j int) bool {
	return n[i].supplier.Id < n[j].supplier.Id
}

// returns:
// ok
// error code:
// 0 = internal error
// 1 = no needs selected
// 2 = the product selected is a manufacturing product
// 3 = the product does not have a supplier
// 4 = no quantity specified
// 5 = the supplier does not have a main billing address
// 6 = the supplier does not have a main shipping address
// 7 = the supplier does not have a payment method
// 8 = the supplier does not have a billing series
func generatePurchaseOrdersFromNeeds(needs []PurchaseNeed, enterpriseId int32, userId int32) (bool, uint8) {
	if len(needs) == 0 {
		return false, 1
	}

	///
	trans, err := db.Begin()
	if err != nil {
		return false, 0
	}
	///

	for i := 0; i < len(needs); i++ {
		product := getProductRow(needs[i].ProductId)
		if product.enterprise != enterpriseId {
			continue
		}
		quantityNeeded := getNeedRow(product.Id)
		if product.Manufacturing {
			trans.Rollback()
			return false, 2
		}
		if product.Supplier == nil || *product.Supplier <= 0 {
			trans.Rollback()
			return false, 3
		}
		if needs[i].Quantity <= 0 || quantityNeeded > needs[i].Quantity {
			trans.Rollback()
			return false, 4
		}

		supplier := getSupplierRow(*product.Supplier)
		needs[i].product = product
		needs[i].supplier = supplier
	}

	sort.Sort(PurchaseNeeds(needs))

	// multisplit the "needs" array into smaller "supplierNeeds" arrays, with the needs of the products of the same supplier
	// create one purchase order for each supplier, with every need as an order detail
	supplierNeeds := make([]PurchaseNeed, 0)
	for i := 0; i < len(needs); i++ {
		var ok bool = true
		if i == 0 || needs[i].supplier.Id == needs[i-1].supplier.Id {
			supplierNeeds = append(supplierNeeds, needs[i])
			ok = false
		} else if i != len(needs)-1 {
			i--
		}
		if ok || i == len(needs)-1 {
			if supplierNeeds[0].supplier.MainBillingAddress == nil {
				trans.Rollback()
				return false, 5
			}
			if supplierNeeds[0].supplier.MainShippingAddress == nil {
				trans.Rollback()
				return false, 6
			}
			if supplierNeeds[0].supplier.PaymentMethod == nil {
				trans.Rollback()
				return false, 7
			}
			if supplierNeeds[0].supplier.BillingSeries == nil {
				trans.Rollback()
				return false, 8
			}
			// "supplierNeeds" is now an array with the needs of the same supplier, create the purchase order and the detail from the supplier and the needs
			o := PurchaseOrder{}
			o.Supplier = supplierNeeds[0].supplier.Id
			o.BillingAddress = *supplierNeeds[0].supplier.MainBillingAddress
			o.ShippingAddress = *supplierNeeds[0].supplier.MainShippingAddress
			o.PaymentMethod = *supplierNeeds[0].supplier.PaymentMethod
			o.BillingSeries = *supplierNeeds[0].supplier.BillingSeries
			o.Warehouse = getPurchaseOrderDefaults(enterpriseId).Warehouse
			o.Currency = *getSupplierDefaults(supplierNeeds[0].supplier.Id, enterpriseId).Currency
			o.enterprise = enterpriseId
			ok, orderId := o.insertPurchaseOrder(userId, trans)
			if !ok || orderId <= 0 {
				trans.Rollback()
				return false, 0
			}

			for j := 0; j < len(supplierNeeds); j++ {
				d := PurchaseOrderDetail{}
				d.Order = orderId
				d.Product = supplierNeeds[j].product.Id
				d.Price = supplierNeeds[j].product.Price
				d.Quantity = supplierNeeds[j].Quantity
				d.VatPercent = supplierNeeds[j].product.VatPercent
				d.enterprise = enterpriseId
				ok, detailId := d.insertPurchaseOrderDetail(userId, trans)
				if !ok.Ok {
					trans.Rollback()
					return false, 0
				}

				// advance the status to "Purchase order pending" of the pending sales order details
				details := getSalesOrderDetailWaitingForPurchaseOrder(supplierNeeds[j].product.Id)
				for k := 0; k < len(details); k++ {
					sqlStatement := `UPDATE sales_order_detail SET status='B',purchase_order_detail=$2 WHERE id=$1`
					_, err := db.Exec(sqlStatement, details[k].Id, detailId)
					if err != nil {
						log("DB", err.Error())
						trans.Rollback()
						return false, 0
					}
					ok := setSalesOrderState(enterpriseId, details[k].Order, userId, *trans)
					if !ok {
						trans.Rollback()
						return false, 0
					}
				}
			}

			supplierNeeds = make([]PurchaseNeed, 0)
		}
	}

	///
	err = trans.Commit()
	return err == nil, 0
	///
}
