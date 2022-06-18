package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesInvoice struct {
	Id                   int64               `json:"id" gorm:"index:sales_invoice_id_enterprise,unique:true,priority:1"`
	CustomerId           int32               `json:"customerId" gorm:"column:customer;not null:true"`
	Customer             Customer            `json:"customer" gorm:"foreignKey:CustomerId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated          time.Time           `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone;index:sales_invoice_date_created,sort:desc"`
	PaymentMethodId      int32               `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod        PaymentMethod       `json:"paymentMethod" gorm:"foreignKey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId      string              `json:"billingSeriesId" gorm:"column:billing_series;not null:true;type:character(3);index:sales_invoice_invoice_number,unique:true,priority:2"`
	BillingSeries        BillingSerie        `json:"billingSeries" gorm:"foreignKey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyId           int32               `json:"currencyId" gorm:"column:currency;not null:true"`
	Currency             Currency            `json:"currency" gorm:"foreignKey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange       float64             `json:"currencyChange" gorm:"type:numeric(14,6);not null:true"`
	BillingAddressId     int32               `json:"billingAddressId" gorm:"column:billing_address;not null:true"`
	BillingAddress       Address             `json:"billingAddress" gorm:"foreignKey:BillingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	TotalProducts        float64             `json:"totalProducts" gorm:"type:numeric(14,6);not null:true"`
	DiscountPercent      float64             `json:"discountPercent" gorm:"type:numeric(14,6);not null:true"`
	FixDiscount          float64             `json:"fixDiscount" gorm:"type:numeric(14,6);not null:true"`
	ShippingPrice        float64             `json:"shippingPrice" gorm:"type:numeric(14,6);not null:true"`
	ShippingDiscount     float64             `json:"shippingDiscount" gorm:"type:numeric(14,6);not null:true"`
	TotalWithDiscount    float64             `json:"totalWithDiscount" gorm:"type:numeric(14,6);not null:true"`
	VatAmount            float64             `json:"vatAmount" gorm:"type:numeric(14,6);not null:true"`
	TotalAmount          float64             `json:"totalAmount" gorm:"type:numeric(14,6);not null:true"`
	LinesNumber          int16               `json:"linesNumber" gorm:"column:lines_number;not null:true"`
	InvoiceNumber        int32               `json:"invoiceNumber" gorm:"column:invoice_number;not null:true;index:sales_invoice_invoice_number,unique:true,priority:3,sort:desc"`
	InvoiceName          string              `json:"invoiceName" gorm:"column:invoice_name;not null:true;type:character(15)"`
	AccountingMovementId *int64              `json:"accountingMovementId" gorm:"column:accounting_movement"`
	AccountingMovement   *AccountingMovement `json:"accountingMovement" gorm:"foreignKey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId         int32               `json:"-" gorm:"column:enterprise;not null:true;index:sales_invoice_id_enterprise,unique:true,priority:2;index:sales_invoice_invoice_number,unique:true,priority:1"`
	Enterprise           Settings            `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	SimplifiedInvoice    bool                `json:"simplifiedInvoice" gorm:"not null:true"`
	Amending             bool                `json:"amending" gorm:"not null:true"`
	AmendedInvoiceId     *int64              `json:"amendedInvoiceId" gorm:"column:amended_invoice"`
	AmendedInvoice       *SalesInvoice       `json:"amendedInvoice" gorm:"foreignKey:AmendedInvoiceId;references:Id"`
}

func (s SalesInvoice) TableName() string {
	return "sales_invoice"
}

type SaleInvoices struct {
	Rows     int64              `json:"rows"`
	Invoices []SalesInvoice     `json:"invoices"`
	Footer   SaleInvoicesFooter `json:"footer"`
}

type SaleInvoicesFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func (q *PaginationQuery) getSalesInvoices() SaleInvoices {
	si := SaleInvoices{}
	if !q.isValid() {
		return si
	}

	si.Invoices = make([]SalesInvoice, 0)
	si.Footer = SaleInvoicesFooter{}

	dbOrm.Model(&SalesInvoice{}).Where("enterprise = ?", q.enterprise).Order("date_created DESC").Offset(int(q.Offset)).Limit(int(q.Limit)).Preload(clause.Associations).Find(&si.Invoices)
	dbOrm.Model(&SalesInvoice{}).Where("enterprise = ?", q.enterprise).Count(&si.Rows).Select("SUM(total_amount) AS total_amount, SUM(total_products) AS total_products").Row().Scan(&si.Footer.TotalAmount, &si.Footer.TotalProducts)

	return si
}

type OrderSearch struct {
	PaginatedSearch
	DateStart *time.Time `json:"dateStart"`
	DateEnd   *time.Time `json:"dateEnd"`
	NotPosted bool       `json:"notPosted"`
}

type SalesInvoiceSearch struct {
	PaginatedSearch
	DateStart         *time.Time `json:"dateStart"`
	DateEnd           *time.Time `json:"dateEnd"`
	NotPosted         bool       `json:"notPosted"`
	PostedStatus      string     `json:"postedStatus"`      // "" = All, "P" = Posted, "N" = Not Posted
	SimplifiedInvoice string     `json:"simplifiedInvoice"` // "" = All, "S" = Simplified, "F" = Full
	Amending          string     `json:"amending"`          // "" = All, "A" = Amending, "R" = Regular
	BillingSeries     *string    `json:"billingSeries"`
}

func (s *SalesInvoiceSearch) searchSalesInvoices() SaleInvoices {
	si := SaleInvoices{}
	if !s.isValid() {
		return si
	}

	cursor := dbOrm.Model(&SalesInvoice{}).Where("sales_invoice.enterprise = ?", s.enterprise)

	si.Invoices = make([]SalesInvoice, 0)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("sales_invoice.invoice_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_invoice.customer").Where("sales_invoice.invoice_name LIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_invoice.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_invoice.date_created <= ?", s.DateEnd)
		}
		if s.NotPosted {
			cursor = cursor.Where("sales_invoice.accounting_movement IS NULL")
		}
		if s.PostedStatus != "" {
			if s.PostedStatus == "P" {
				cursor = cursor.Where("sales_invoice.accounting_movement IS NOT NULL")
			} else if s.PostedStatus == "N" {
				cursor = cursor.Where("sales_invoice.accounting_movement IS NULL")
			}
		}
		if s.SimplifiedInvoice != "" {
			if s.SimplifiedInvoice == "S" {
				cursor = cursor.Where("sales_invoice.simplified_invoice = ?", true)
			} else if s.SimplifiedInvoice == "F" {
				cursor = cursor.Where("sales_invoice.simplified_invoice = ?", false)
			}
		}
		if s.Amending != "" {
			if s.Amending == "A" {
				cursor = cursor.Where("sales_invoice.amending = ?", true)
			} else if s.Amending == "R" {
				cursor = cursor.Where("sales_invoice.amending = ?", false)
			}
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("sales_invoice.billing_series = ?", *s.BillingSeries)
		}
	}
	result := cursor.Order("sales_invoice.date_created DESC").Offset(int(s.Offset)).Limit(int(s.Limit)).Preload(clause.Associations).Find(&si.Invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return si
	}

	cursor = dbOrm.Model(&SalesInvoice{}).Where("sales_invoice.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("sales_invoice.invoice_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN customer ON customer.id=sales_invoice.customer").Where("sales_invoice.invoice_name LIKE @search OR customer.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("sales_invoice.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("sales_invoice.date_created <= ?", s.DateEnd)
		}
		if s.NotPosted {
			cursor = cursor.Where("sales_invoice.accounting_movement IS NULL")
		}
		if s.PostedStatus != "" {
			if s.PostedStatus == "P" {
				cursor = cursor.Where("sales_invoice.accounting_movement IS NOT NULL")
			} else if s.PostedStatus == "N" {
				cursor = cursor.Where("sales_invoice.accounting_movement IS NULL")
			}
		}
		if s.SimplifiedInvoice != "" {
			if s.SimplifiedInvoice == "S" {
				cursor = cursor.Where("sales_invoice.simplified_invoice = ?", true)
			} else if s.SimplifiedInvoice == "N" {
				cursor = cursor.Where("sales_invoice.simplified_invoice = ?", false)
			}
		}
		if s.Amending != "" {
			if s.Amending == "A" {
				cursor = cursor.Where("sales_invoice.amending = ?", true)
			} else if s.Amending == "R" {
				cursor = cursor.Where("sales_invoice.amending = ?", false)
			}
		}
		if s.BillingSeries != nil {
			cursor = cursor.Where("sales_invoice.billing_series = ?", *s.BillingSeries)
		}
	}
	result = cursor.Count(&si.Rows).Select("SUM(total_products) AS total_products,SUM(total_amount) AS total_amount").Scan(&si.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return si
	}

	return si
}

func getSalesInvoiceRow(invoiceId int64) SalesInvoice {
	si := SalesInvoice{}
	result := dbOrm.Model(&SalesInvoice{}).Where("id = ?", invoiceId).Preload(clause.Associations).First(&si)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return si
	}
	return si
}

func getSalesInvoiceRowTransaction(invoiceId int64, trans gorm.DB) SalesInvoice {
	si := SalesInvoice{}
	trans.Model(&SalesInvoice{}).Where("id = ?", invoiceId).First(&si)
	return si
}

func (i *SalesInvoice) isValid() bool {
	return !(i.CustomerId <= 0 || i.PaymentMethodId <= 0 || len(i.BillingSeriesId) == 0 || i.CurrencyId <= 0 || i.BillingAddressId <= 0)
}

func (c *SalesInvoice) BeforeCreate(tx *gorm.DB) (err error) {
	var salesInvoice SalesInvoice
	tx.Model(&SalesInvoice{}).Last(&salesInvoice)
	c.Id = salesInvoice.Id + 1
	return nil
}

func (i *SalesInvoice) insertSalesInvoice(userId int32, trans *gorm.DB) (bool, int64) {
	if !i.isValid() {
		return false, 0
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return false, 0
		}
		///
	}

	// get invoice name
	i.InvoiceNumber = getNextSaleInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
	if i.InvoiceNumber <= 0 {
		trans.Rollback()
		return false, 0
	}
	now := time.Now()
	i.InvoiceName = i.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	// get currency exchange
	i.CurrencyChange = getCurrencyExchange(i.CurrencyId)

	// simplified invoice
	address := getAddressRow(i.BillingAddressId)
	if address.Id <= 0 {
		trans.Rollback()
		return false, 0
	}
	country := getCountryRow(address.CountryId, i.EnterpriseId)
	if country.Id <= 0 {
		trans.Rollback()
		return false, 0
	}
	if country.Zone == "E" { // Export
		i.SimplifiedInvoice = false
	} else {
		customer := getCustomerRow(i.CustomerId)
		if country.Zone == "N" { // National
			i.SimplifiedInvoice = len(customer.TaxId) == 0
		} else { // European Union
			i.SimplifiedInvoice = len(customer.TaxId) == 0 && len(customer.VatNumber) == 0
		}
	}

	i.DateCreated = time.Now()
	i.TotalProducts = 0
	i.VatAmount = 0
	i.AccountingMovementId = nil
	i.Amending = false
	i.AmendedInvoiceId = nil

	result := trans.Create(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(i.EnterpriseId, "sales_invoice", int(i.Id), userId, "I")
	json, _ := json.Marshal(i)
	go fireWebHook(i.EnterpriseId, "sales_invoice", "POST", string(json))

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return false, 0
		}
		///
	}

	return true, i.Id
}

// 1. can't delete details in posted invoices
// 2. the invoice deletion is completely disallowed by policy
// 3. it is only allowed to delete the latest invoice of the billing series
func (i *SalesInvoice) deleteSalesInvoice(userId int32) OkAndErrorCodeReturn {
	if i.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice := getSalesInvoiceRow(i.Id)
	if invoice.AccountingMovementId != nil {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(i.EnterpriseId)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoiceNumber := getNextSaleInvoiceNumber(invoice.BillingSeriesId, invoice.EnterpriseId)
		if invoiceNumber <= 0 || invoice.InvoiceNumber != (invoiceNumber-1) {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
		}
	}

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	d := getSalesInvoiceDetail(i.Id, i.EnterpriseId)

	for i := 0; i < len(d); i++ {
		ok := d[i].deleteSalesInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(i.EnterpriseId, "sales_invoice", int(i.Id), userId, "D")
	json, _ := json.Marshal(i)
	go fireWebHook(i.EnterpriseId, "sales_invoice", "DELETE", string(json))

	result := trans.Delete(&SalesInvoice{}, "id = ? AND enterprise = ?", i.Id, i.EnterpriseId)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	///
	result = trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	return OkAndErrorCodeReturn{Ok: true}
}

func toggleSimplifiedInvoiceSalesInvoice(invoiceId int64, enterpriseId int32, userId int32) bool {
	invoice := getSalesInvoiceRow(invoiceId)
	if invoice.Id <= 0 {
		return false
	}
	invoice.SimplifiedInvoice = !invoice.SimplifiedInvoice
	result := dbOrm.Model(&SalesInvoice{}).Where("id = ?", invoiceId).Update("simplified_invoice", invoice.SimplifiedInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_invoice", int(invoiceId), userId, "U")
	json, _ := json.Marshal(invoice)
	go fireWebHook(enterpriseId, "sales_invoice", "PUT", string(json))

	return true
}

type MakeAmendingInvoice struct {
	InvoiceId   int64   `json:"invoiceId"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
}

func makeAmendingSaleInvoice(invoiceId int64, enterpriseId int32, quantity float64, description string, userId int32) bool {
	i := getSalesInvoiceRow(invoiceId)
	if i.Id <= 0 || i.EnterpriseId != enterpriseId {
		return false
	}

	// we can't make an amending invoice the same day as the original invoice
	now := time.Now()
	if i.DateCreated.Year() == now.Year() && i.DateCreated.YearDay() == now.YearDay() {
		return false
	}
	// we can't make an amending invoice with a greater amount that the original invoice
	if quantity <= 0 || quantity > i.TotalAmount {
		return false
	}

	settings := getSettingsRecordById(enterpriseId)

	// get invoice name
	invoiceNumber := getNextSaleInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
	if i.InvoiceNumber <= 0 {
		return false
	}
	invoiceName := i.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	var detailAmount float64
	var vatPercent float64
	// VAT excluded invoice, when amending the invoice, we are not returning any tax money
	if i.VatAmount == 0 {
		detailAmount = quantity
		vatPercent = 0
	} else { // Invoice with VAT, we return the quantity without the tax to the customer, and then we add the tax percent, so the total of the invoice is the amount we want to return (in taxes and to the customer)
		detailAmount = quantity / (1 + (settings.DefaultVatPercent / 100))
		vatPercent = settings.DefaultVatPercent
	}
	var vatAmount float64 = (quantity / 100) * vatPercent

	var amendingInvoice SalesInvoice = SalesInvoice{}

	amendingInvoice.CustomerId = i.CustomerId
	amendingInvoice.DateCreated = time.Now()
	amendingInvoice.PaymentMethodId = i.PaymentMethodId
	amendingInvoice.BillingSeriesId = i.BillingSeriesId
	amendingInvoice.CurrencyId = i.CurrencyId
	amendingInvoice.CurrencyChange = i.CurrencyChange
	amendingInvoice.BillingAddressId = i.BillingAddressId
	amendingInvoice.SimplifiedInvoice = i.SimplifiedInvoice
	amendingInvoice.Amending = true
	amendingInvoice.AmendedInvoiceId = &i.Id
	amendingInvoice.InvoiceNumber = invoiceNumber
	amendingInvoice.InvoiceName = invoiceName
	amendingInvoice.TotalProducts = -detailAmount
	amendingInvoice.TotalWithDiscount = -detailAmount
	amendingInvoice.VatAmount = -vatAmount
	amendingInvoice.TotalAmount = -quantity
	amendingInvoice.DiscountPercent = 0
	amendingInvoice.FixDiscount = 0
	amendingInvoice.ShippingPrice = 0
	amendingInvoice.ShippingDiscount = 0
	amendingInvoice.EnterpriseId = enterpriseId

	result := dbOrm.Create(&amendingInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	detail := SalesInvoiceDetail{}

	detail.InvoiceId = amendingInvoice.Id
	detail.Description = description
	detail.Price = -detailAmount
	detail.Quantity = 1
	detail.VatPercent = vatPercent
	detail.TotalAmount = -quantity
	detail.EnterpriseId = enterpriseId

	result = dbOrm.Create(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_invoice", int(invoiceId), userId, "I")
	jsn, _ := json.Marshal(amendingInvoice)
	go fireWebHook(amendingInvoice.EnterpriseId, "sales_invoice", "PUT", string(jsn))

	insertTransactionalLog(enterpriseId, "sales_invoice_detail", int(detail.Id), userId, "I")
	jsn, _ = json.Marshal(detail)
	go fireWebHook(detail.EnterpriseId, "sales_invoice_detail", "PUT", string(jsn))

	return true
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsSalesInvoice(invoiceId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var saleInvoice SalesInvoice
	result := trans.Model(&SalesInvoice{}).Where("id = ?", invoiceId).First(&saleInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleInvoice.TotalProducts += totalAmount
	saleInvoice.VatAmount += (totalAmount / 100) * vatPercent

	result = trans.Save(&saleInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsSaleInvoice(enterpriseId, invoiceId, userId, trans)
}

// Applies the logic to calculate the totals of the sales invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsSaleInvoice(enterpriseId int32, invoiceId int64, userId int32, trans gorm.DB) bool {
	var saleInvoice SalesInvoice
	result := trans.Model(&SalesInvoice{}).Where("id = ?", invoiceId).First(&saleInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	saleInvoice.TotalWithDiscount = (saleInvoice.TotalProducts - saleInvoice.TotalProducts*(saleInvoice.DiscountPercent/100)) - saleInvoice.FixDiscount + saleInvoice.ShippingPrice - saleInvoice.ShippingDiscount
	saleInvoice.TotalAmount = saleInvoice.TotalWithDiscount + saleInvoice.VatAmount

	result = trans.Save(&saleInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "sales_invoice", int(invoiceId), userId, "U")
	json, _ := json.Marshal(saleInvoice)
	go fireWebHook(enterpriseId, "sales_invoice", "PUT", string(json))

	return true
}

// ERROR CODES:
// 1. The order is already invoiced
// 2. There are no details to invoice
func invoiceAllSaleOrder(saleOrderId int64, enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(saleOrderId)
	if saleOrder.EnterpriseId != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if saleOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if saleOrder.InvoicedLines >= saleOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}
	orderDetails := getSalesOrderDetail(saleOrderId, saleOrder.EnterpriseId)
	filterSalesOrderDetails(orderDetails, func(sod SalesOrderDetail) bool { return sod.QuantityInvoiced < sod.Quantity })
	if len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	// create an invoice for that order
	invoice := SalesInvoice{}
	invoice.CustomerId = saleOrder.CustomerId
	invoice.BillingAddressId = saleOrder.BillingAddressId
	invoice.BillingSeriesId = saleOrder.BillingSeriesId
	invoice.CurrencyId = saleOrder.CurrencyId
	invoice.PaymentMethodId = saleOrder.PaymentMethodId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedSalesOrder(enterpriseId, saleOrder.Id, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.EnterpriseId = saleOrder.EnterpriseId
	ok, invoiceId := invoice.insertSalesInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		invoiceDetal := SalesInvoiceDetail{}
		invoiceDetal.InvoiceId = invoiceId
		invoiceDetal.OrderDetailId = &orderDetail.Id
		invoiceDetal.Price = orderDetail.Price
		invoiceDetal.ProductId = &orderDetail.ProductId
		invoiceDetal.Quantity = (orderDetail.Quantity - orderDetail.QuantityInvoiced)
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		invoiceDetal.EnterpriseId = invoice.EnterpriseId
		ok := invoiceDetal.insertSalesInvoiceDetail(trans, userId)
		if !ok.Ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	invoiceSalesOrderDiscounts(saleOrderId, invoiceId, enterpriseId, userId, *trans)

	go ecommerceControllerupdateStatusPaymentAccepted(saleOrderId, invoice.EnterpriseId)

	///
	result := trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	return OkAndErrorCodeReturn{Ok: true}
	///
}

type OrderDetailGenerate struct {
	OrderId   int64                          `json:"orderId"`
	Selection []OrderDetailGenerateSelection `json:"selection"`
}

type OrderDetailGenerateSelection struct {
	Id       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

// ERROR CODES:
// 1. The order is aleady invoiced
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail is already invoiced
// 4. The selected quantity is greater than the quantity pending of invoicing in the detail
func (invoiceInfo *OrderDetailGenerate) invoicePartiallySaleOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the sale order and it's details
	saleOrder := getSalesOrderRow(invoiceInfo.OrderId)
	if saleOrder.Id <= 0 || saleOrder.EnterpriseId != enterpriseId || len(invoiceInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if saleOrder.InvoicedLines >= saleOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	var saleOrderDetails []SalesOrderDetail = make([]SalesOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getSalesOrderDetailRow(invoiceInfo.Selection[i].Id)
		if orderDetail.Id <= 0 || orderDetail.OrderId != invoiceInfo.OrderId || invoiceInfo.Selection[i].Quantity == 0 {
			return OkAndErrorCodeReturn{Ok: false}
		}
		if invoiceInfo.Selection[i].Quantity > orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2, ExtraData: []string{product.Name}}
		}
		if orderDetail.QuantityInvoiced >= orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3, ExtraData: []string{product.Name}}
		}
		if (invoiceInfo.Selection[i].Quantity + orderDetail.QuantityInvoiced) > orderDetail.Quantity {
			product := getProductRow(orderDetail.ProductId)
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 4, ExtraData: []string{product.Name}}
		}
		saleOrderDetails = append(saleOrderDetails, orderDetail)
	}

	// create an invoice for that order
	invoice := SalesInvoice{}
	invoice.CustomerId = saleOrder.CustomerId
	invoice.BillingAddressId = saleOrder.BillingAddressId
	invoice.BillingSeriesId = saleOrder.BillingSeriesId
	invoice.CurrencyId = saleOrder.CurrencyId
	invoice.PaymentMethodId = saleOrder.PaymentMethodId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedSalesOrder(enterpriseId, saleOrder.Id, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.EnterpriseId = saleOrder.EnterpriseId
	ok, invoiceId := invoice.insertSalesInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(saleOrderDetails); i++ {
		orderDetail := saleOrderDetails[i]
		invoiceDetal := SalesInvoiceDetail{}
		invoiceDetal.InvoiceId = invoiceId
		invoiceDetal.OrderDetailId = &orderDetail.Id
		invoiceDetal.Price = orderDetail.Price
		invoiceDetal.ProductId = &orderDetail.ProductId
		invoiceDetal.Quantity = invoiceInfo.Selection[i].Quantity
		invoiceDetal.TotalAmount = orderDetail.TotalAmount
		invoiceDetal.VatPercent = orderDetail.VatPercent
		invoiceDetal.EnterpriseId = invoice.EnterpriseId
		ok = invoiceDetal.insertSalesInvoiceDetail(trans, userId).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	invoiceSalesOrderDiscounts(saleOrder.Id, invoiceId, enterpriseId, userId, *trans)

	go ecommerceControllerupdateStatusPaymentAccepted(invoiceInfo.OrderId, invoice.EnterpriseId)

	///
	result := trans.Commit()
	if result.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	return OkAndErrorCodeReturn{Ok: true}
	///
}

type SalesInvoiceRelations struct {
	Orders        []SaleOrder         `json:"orders"`
	DeliveryNotes []SalesDeliveryNote `json:"notes"`
	Invoices      []SalesInvoice      `json:"invoices"`
}

func getSalesInvoiceRelations(invoiceId int64, enterpriseId int32) SalesInvoiceRelations {
	return SalesInvoiceRelations{
		Orders:        getSalesInvoiceOrders(invoiceId, enterpriseId),
		DeliveryNotes: getSalesInvoiceDeliveryNotes(invoiceId, enterpriseId),
		Invoices:      getSalesInvoiceAmendingAmendedInvoices(invoiceId, enterpriseId),
	}
}

func getSalesInvoiceOrders(invoiceId int64, enterpriseId int32) []SaleOrder {
	saleInvoice := getSalesInvoiceRow(invoiceId)
	if saleInvoice.Id <= 0 || saleInvoice.EnterpriseId != enterpriseId {
		return []SaleOrder{}
	}
	var orders []SaleOrder = make([]SaleOrder, 0)
	invoiceDetails := getSalesInvoiceDetail(invoiceId, enterpriseId)
	for i := 0; i < len(invoiceDetails); i++ {
		var orderDetail SalesOrderDetail
		if invoiceDetails[i].OrderDetailId != nil {
			orderDetail = getSalesOrderDetailRow(*invoiceDetails[i].OrderDetailId)
		} else {
			continue
		}
		var ok bool = true
		for j := 0; j < len(orders); j++ {
			if orders[j].Id == orderDetail.OrderId {
				ok = false
				break
			}
		}
		if ok {
			order := getSalesOrderRow(orderDetail.OrderId)
			orders = append(orders, order)
		}
	}
	return orders
}

func getSalesInvoiceDeliveryNotes(invoiceId int64, enterpriseId int32) []SalesDeliveryNote {
	var notes []SalesDeliveryNote = make([]SalesDeliveryNote, 0)

	orders := getSalesInvoiceOrders(invoiceId, enterpriseId)
	for i := 0; i < len(orders); i++ {
		relations := getSalesOrderDeliveryNotes(orders[i].Id, enterpriseId)
		notes = append(notes, relations...)
	}

	return notes
}

func getSalesInvoiceAmendingAmendedInvoices(invoiceId int64, enterpriseId int32) []SalesInvoice {
	invoices := make([]SalesInvoice, 0)

	i := getSalesInvoiceRow(invoiceId)
	if i.EnterpriseId != enterpriseId {
		return invoices
	}

	result := dbOrm.Model(&SalesInvoice{}).Where("amended_invoice = ? AND enterprise = ?", i.Id, enterpriseId).Find(&invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return invoices
	}

	if i.Amending && i.AmendedInvoiceId != nil {
		invoices = append(invoices, getSalesInvoiceRow(*i.AmendedInvoiceId))
	}

	return invoices
}
