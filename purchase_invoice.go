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

type PurchaseInvoice struct {
	Id                   int64               `json:"id" gorm:"index:purchase_invoice_id_enterprise,unique:true,priority:1"`
	SupplierId           int32               `json:"supplierId" gorm:"column:supplier;not null:true"`
	Supplier             Supplier            `json:"supplier" gorm:"foreignkey:SupplierId,EnterpriseId;references:Id,EnterpriseId"`
	DateCreated          time.Time           `json:"dateCreated" gorm:"column:date_created;not null:true;type:timestamp(3) with time zone;index:purchase_invoice_date_created,sort:desc"`
	PaymentMethodId      int32               `json:"paymentMethodId" gorm:"column:payment_method;not null:true"`
	PaymentMethod        PaymentMethod       `json:"paymentMethod" gorm:"foreignkey:PaymentMethodId,EnterpriseId;references:Id,EnterpriseId"`
	BillingSeriesId      string              `json:"billingSeriesId" gorm:"column:billing_series;not null:true;index:purcahse_invoice_invoice_number,unique:true,priority:2"`
	BillingSeries        BillingSerie        `json:"billingSeries" gorm:"foreignkey:BillingSeriesId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyId           int32               `json:"currencyId" gorm:"column:currency;not null:true"`
	Currency             Currency            `json:"currency" gorm:"foreignkey:CurrencyId,EnterpriseId;references:Id,EnterpriseId"`
	CurrencyChange       float64             `json:"currencyChange" gorm:"column:currency_change;not null:true;type:numeric(14,6)"`
	BillingAddressId     int32               `json:"billingAddressId" gorm:"column:billing_address;not null:true"`
	BillingAddress       Address             `json:"billingAddress" gorm:"foreignkey:BillingAddressId,EnterpriseId;references:Id,EnterpriseId"`
	TotalProducts        float64             `json:"totalProducts" gorm:"column:total_products;not null:true;type:numeric(14,6)"`
	DiscountPercent      float64             `json:"discountPercent" gorm:"column:discount_percent;not null:true;type:numeric(14,6)"`
	FixDiscount          float64             `json:"fixDiscount" gorm:"column:fix_discount;not null:true;type:numeric(14,6)"`
	ShippingPrice        float64             `json:"shippingPrice" gorm:"column:shipping_price;not null:true;type:numeric(14,6)"`
	ShippingDiscount     float64             `json:"shippingDiscount" gorm:"column:shipping_discount;not null:true;type:numeric(14,6)"`
	TotalWithDiscount    float64             `json:"totalWithDiscount" gorm:"column:total_with_discount;not null:true;type:numeric(14,6)"`
	VatAmount            float64             `json:"vatAmount" gorm:"column:vat_amount;not null:true;type:numeric(14,6)"`
	TotalAmount          float64             `json:"totalAmount" gorm:"column:total_amount;not null:true;type:numeric(14,6)"`
	LinesNumber          int16               `json:"linesNumber" gorm:"column:lines_number;not null:true"`
	InvoiceNumber        int32               `json:"invoiceNumber" gorm:"column:invoice_number;not null:true;index:purcahse_invoice_invoice_number,unique:true,priority:3,sort:desc"`
	InvoiceName          string              `json:"invoiceName" gorm:"column:invoice_name;not null:true;type:character(15)"`
	AccountingMovementId *int64              `json:"accountingMovementId" gorm:"column:accounting_movement"`
	AccountingMovement   *AccountingMovement `json:"accountingMovement" gorm:"foreignkey:AccountingMovementId,EnterpriseId;references:Id,EnterpriseId"`
	EnterpriseId         int32               `json:"-" gorm:"column:enterprise;not null:true;index:purchase_invoice_id_enterprise,unique:true,priority:2;index:purcahse_invoice_invoice_number,unique:true,priority:1"`
	Enterprise           Settings            `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	Amending             bool                `json:"amending" gorm:"column:amending;not null:true"`
	AmendedInvoiceId     *int64              `json:"amendedInvoiceId" gorm:"column:amended_invoice"`
	AmendedInvoice       *PurchaseInvoice    `json:"amendedInvoice" gorm:"foreignkey:AmendedInvoiceId;references:Id"`
	IncomeTax            bool                `json:"incomeTax" gorm:"column:income_tax;not null:true"`
	IncomeTaxBase        float64             `json:"incomeTaxBase" gorm:"column:income_tax_base;not null:true;type:real"`
	IncomeTaxPercentage  float64             `json:"incomeTaxPercentage" gorm:"column:income_tax_percentage;not null:true;type:real"`
	IncomeTaxValue       float64             `json:"incomeTaxValue" gorm:"column:income_tax_value;not null:true;type:real"`
	Rent                 bool                `json:"rent" gorm:"column:rent;not null:true"`
	RentBase             float64             `json:"rentBase" gorm:"column:rent_base;not null:true;type:real"`
	RentPercentage       float64             `json:"rentPercentage" gorm:"column:rent_percentage;not null:true;type:real"`
	RentValue            float64             `json:"rentValue" gorm:"column:rent_value;not null:true;type:real"`
}

func (pi *PurchaseInvoice) TableName() string {
	return "purchase_invoice"
}

type PurchaseInvoices struct {
	Rows     int64                 `json:"rows"`
	Invoices []PurchaseInvoice     `json:"invoices"`
	Footer   PurchaseInvoiceFooter `json:"footer"`
}

type PurchaseInvoiceFooter struct {
	TotalProducts float64 `json:"totalProducts"`
	TotalAmount   float64 `json:"totalAmount"`
}

func getPurchaseInvoices(enterpriseId int32) PurchaseInvoices {
	in := PurchaseInvoices{}
	in.Invoices = make([]PurchaseInvoice, 0)
	result := dbOrm.Model(&PurchaseInvoice{}).Where("purchase_invoice.enterprise = ?", enterpriseId).Order("purchase_invoice.date_created DESC").Preload(clause.Associations).Find(&in.Invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return in
	}
	result = dbOrm.Model(&PurchaseInvoice{}).Where("purchase_invoice.enterprise = ?", enterpriseId).Count(&in.Rows).Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&in.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return in
	}
	return in
}

type PurchaseInvoiceSearch struct {
	PaginatedSearch
	DateStart         *time.Time `json:"dateStart"`
	DateEnd           *time.Time `json:"dateEnd"`
	NotPosted         bool       `json:"notPosted"`
	PostedStatus      string     `json:"postedStatus"`      // "" = All, "P" = Posted, "N" = Not Posted
	SimplifiedInvoice string     `json:"simplifiedInvoice"` // "" = All, "S" = Simplified, "F" = Full
	Amending          string     `json:"amending"`          // "" = All, "A" = Amending, "R" = Regular
	BillingSeries     *string    `json:"billingSeries"`
}

func (s *PurchaseInvoiceSearch) searchPurchaseInvoice() PurchaseInvoices {
	in := PurchaseInvoices{}
	in.Invoices = make([]PurchaseInvoice, 0)

	cursor := dbOrm.Model(&PurchaseInvoice{}).Where("purchase_invoice.enterprise = ?", s.enterprise)
	orderNumber, err := strconv.Atoi(s.Search)
	if err == nil {
		cursor = cursor.Where("purchase_invoice.invoice_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_invoice.supplier").Where("purchase_invoice.invoice_name LIKE @search OR suppliers.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_invoice.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_invoice.date_created <= ?", s.DateEnd)
		}
		if s.NotPosted {
			cursor = cursor.Where("purchase_invoice.accounting_movement IS NULL")
		}
		if s.PostedStatus != "" {
			if s.PostedStatus == "P" {
				cursor = cursor.Where("purchase_invoice.accounting_movement IS NOT NULL")
			} else if s.PostedStatus == "N" {
				cursor = cursor.Where("purchase_invoice.accounting_movement IS NULL")
			}
		}
		if s.SimplifiedInvoice != "" {
			if s.SimplifiedInvoice == "S" {
				cursor = cursor.Where("purchase_invoice.simplified_invoice = ?", true)
			} else if s.SimplifiedInvoice == "F" {
				cursor = cursor.Where("purchase_invoice.simplified_invoice = ?", false)
			}
		}
		if s.Amending != "" {
			if s.Amending == "A" {
				cursor = cursor.Where("purchase_invoice.amending = ?", true)
			} else if s.Amending == "R" {
				cursor = cursor.Where("purchase_invoice.amending = ?", false)
			}
		}
	}
	result := cursor.Order("purchase_invoice.date_created DESC").Limit(int(s.Limit)).Offset(int(s.Offset)).Count(&in.Rows).Preload(clause.Associations).Find(&in.Invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return in
	}

	cursor = dbOrm.Model(&PurchaseInvoice{}).Where("purchase_invoice.enterprise = ?", s.enterprise)
	if err == nil {
		cursor = cursor.Where("purchase_invoice.invoice_number = ?", orderNumber)
	} else {
		cursor = cursor.Joins("INNER JOIN suppliers ON suppliers.id=purchase_invoice.supplier").Where("purchase_invoice.invoice_name LIKE @search OR suppliers.name ILIKE @search", sql.Named("search", "%"+s.Search+"%"))
		if s.DateStart != nil {
			cursor = cursor.Where("purchase_invoice.date_created >= ?", s.DateStart)
		}
		if s.DateEnd != nil {
			cursor = cursor.Where("purchase_invoice.date_created <= ?", s.DateEnd)
		}
		if s.NotPosted {
			cursor = cursor.Where("purchase_invoice.accounting_movement IS NULL")
		}
		if s.PostedStatus != "" {
			if s.PostedStatus == "P" {
				cursor = cursor.Where("purchase_invoice.accounting_movement IS NOT NULL")
			} else if s.PostedStatus == "N" {
				cursor = cursor.Where("purchase_invoice.accounting_movement IS NULL")
			}
		}
		if s.SimplifiedInvoice != "" {
			if s.SimplifiedInvoice == "S" {
				cursor = cursor.Where("purchase_invoice.simplified_invoice = ?", true)
			} else if s.SimplifiedInvoice == "F" {
				cursor = cursor.Where("purchase_invoice.simplified_invoice = ?", false)
			}
		}
		if s.Amending != "" {
			if s.Amending == "A" {
				cursor = cursor.Where("purchase_invoice.amending = ?", true)
			} else if s.Amending == "R" {
				cursor = cursor.Where("purchase_invoice.amending = ?", false)
			}
		}
	}
	result = cursor.Select("SUM(total_products) as total_products, SUM(total_amount) as total_amount").Scan(&in.Footer)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}

	return in
}

func getPurchaseInvoiceRow(invoiceId int64) PurchaseInvoice {
	i := PurchaseInvoice{}
	result := dbOrm.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).Preload(clause.Associations).First(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return i
}

func getPurchaseInvoiceRowTransaction(invoiceId int64, trans gorm.DB) PurchaseInvoice {
	i := PurchaseInvoice{}
	result := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).Preload(clause.Associations).First(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
	}
	return i
}

func (i *PurchaseInvoice) isValid() bool {
	return !(i.SupplierId <= 0 || i.PaymentMethodId <= 0 || len(i.BillingSeriesId) == 0 || i.CurrencyId <= 0 || i.BillingAddressId <= 0 || i.IncomeTaxBase < 0 || i.IncomeTaxPercentage < 0 || i.RentBase < 0 || i.RentPercentage < 0)
}

func (i *PurchaseInvoice) BeforeCreate(tx *gorm.DB) (err error) {
	var purchaseInvoice PurchaseInvoice
	tx.Model(&PurchaseInvoice{}).Last(&purchaseInvoice)
	i.Id = purchaseInvoice.Id + 1
	return nil
}

func (i *PurchaseInvoice) insertPurchaseInvoice(userId int32, trans *gorm.DB) (bool, int64) {
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

	i.InvoiceNumber = getNextPurchaseInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
	if i.InvoiceNumber <= 0 {
		return false, 0
	}
	i.CurrencyChange = getCurrencyExchange(i.CurrencyId)
	now := time.Now()
	i.InvoiceName = i.BillingSeriesId + "/" + strconv.Itoa(now.Year()) + "/" + fmt.Sprintf("%06d", i.InvoiceNumber)

	i.DateCreated = time.Now()
	i.TotalProducts = 0
	i.VatAmount = 0
	i.LinesNumber = 0
	i.AccountingMovementId = nil
	i.Amending = false
	i.AmendedInvoiceId = nil
	i.IncomeTaxBase = 0
	i.IncomeTaxValue = 0
	i.RentBase = 0
	i.RentValue = 0

	result := trans.Create(&i)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false, 0
	}

	insertTransactionalLog(i.EnterpriseId, "purchase_invoice", int(i.Id), userId, "I")
	json, _ := json.Marshal(i)
	go fireWebHook(i.EnterpriseId, "purchase_invoice", "POST", string(json))

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
func (i *PurchaseInvoice) deletePurchaseInvoice(userId int32, trans *gorm.DB) OkAndErrorCodeReturn {
	if i.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice := getPurchaseInvoiceRow(i.Id)
	if invoice.EnterpriseId != i.EnterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if invoice.AccountingMovementId != nil {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	// INVOICE DELETION POLICY
	s := getSettingsRecordById(i.EnterpriseId)
	if s.InvoiceDeletePolicy == 2 { // Don't allow to delete
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	} else if s.InvoiceDeletePolicy == 1 { // Allow to delete only the latest invoice of the billing series
		invoice := getPurchaseInvoiceRow(i.Id)
		invoiceNumber := getNextPurchaseInvoiceNumber(invoice.BillingSeriesId, invoice.EnterpriseId)
		if invoiceNumber <= 0 || invoice.InvoiceNumber != (invoiceNumber-1) {
			return OkAndErrorCodeReturn{Ok: false, ErrorCode: 3}
		}
	}

	var beginTransaction bool = (trans == nil)
	if trans == nil {
		///
		trans = dbOrm.Begin()
		if trans.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	d := getPurchaseInvoiceDetail(i.Id, i.EnterpriseId)
	for i := 0; i < len(d); i++ {
		ok := d[i].deletePurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	insertTransactionalLog(invoice.EnterpriseId, "purchase_invoice", int(i.Id), userId, "D")
	json, _ := json.Marshal(i)
	go fireWebHook(i.EnterpriseId, "purchase_invoice", "DELETE", string(json))

	result := trans.Delete(&PurchaseInvoice{}, "id = ?", i.Id)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	if beginTransaction {
		///
		result = trans.Commit()
		if result.Error != nil {
			return OkAndErrorCodeReturn{Ok: false}
		}
		///
	}

	return OkAndErrorCodeReturn{Ok: true}
}

func makeAmendingPurchaseInvoice(invoiceId int64, enterpriseId int32, quantity float64, description string, userId int32) bool {
	i := getPurchaseInvoiceRow(invoiceId)
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
	invoiceNumber := getNextPurchaseInvoiceNumber(i.BillingSeriesId, i.EnterpriseId)
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

	var amendingInvoice PurchaseInvoice = PurchaseInvoice{}

	amendingInvoice.SupplierId = i.SupplierId
	amendingInvoice.DateCreated = time.Now()
	amendingInvoice.PaymentMethodId = i.PaymentMethodId
	amendingInvoice.BillingSeriesId = i.BillingSeriesId
	amendingInvoice.CurrencyId = i.CurrencyId
	amendingInvoice.CurrencyChange = i.CurrencyChange
	amendingInvoice.BillingAddressId = i.BillingAddressId
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
	amendingInvoice.IncomeTax = i.IncomeTax
	amendingInvoice.IncomeTaxPercentage = i.IncomeTaxPercentage
	amendingInvoice.Rent = i.Rent
	amendingInvoice.RentPercentage = i.RentPercentage
	amendingInvoice.EnterpriseId = enterpriseId

	result := dbOrm.Create(&amendingInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_invoice", int(amendingInvoice.Id), userId, "I")
	jsonData, _ := json.Marshal(amendingInvoice)
	go fireWebHook(i.EnterpriseId, "purchase_invoice", "POST", string(jsonData))

	detail := PurchaseInvoiceDetail{}

	detail.InvoiceId = amendingInvoice.Id
	detail.Description = description
	detail.Price = -detailAmount
	detail.Quantity = 1
	detail.VatPercent = vatPercent
	detail.TotalAmount = -quantity
	detail.IncomeTax = i.IncomeTax
	detail.Rent = i.Rent
	detail.EnterpriseId = enterpriseId

	result = dbOrm.Create(&detail)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_invoice_details", int(detail.Id), userId, "I")
	jsonData, _ = json.Marshal(detail)
	go fireWebHook(enterpriseId, "purchase_invoice_details", "POST", string(jsonData))

	return true
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addTotalProductsPurchaseInvoice(invoiceId int64, totalAmount float64, vatPercent float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseInvoice PurchaseInvoice
	result := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).First(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseInvoice.TotalProducts += totalAmount
	purchaseInvoice.VatAmount += (totalAmount / 100) * vatPercent

	result = trans.Save(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// Applies the logic to calculate the totals of the purchase invoice and the discounts.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func calcTotalsPurchaseInvoice(invoiceId int64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseInvoice PurchaseInvoice
	result := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).First(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseInvoice.TotalWithDiscount = (purchaseInvoice.TotalProducts - purchaseInvoice.TotalProducts*(purchaseInvoice.DiscountPercent/100)) - purchaseInvoice.FixDiscount + purchaseInvoice.ShippingPrice - purchaseInvoice.ShippingDiscount
	if purchaseInvoice.IncomeTax {
		purchaseInvoice.IncomeTaxValue = (purchaseInvoice.IncomeTaxBase / 100) * purchaseInvoice.IncomeTaxPercentage
	} else {
		purchaseInvoice.IncomeTaxValue = 0
	}
	if purchaseInvoice.Rent {
		purchaseInvoice.RentValue = (purchaseInvoice.RentBase / 100) * purchaseInvoice.RentPercentage
	} else {
		purchaseInvoice.RentValue = 0
	}
	purchaseInvoice.TotalAmount = purchaseInvoice.TotalWithDiscount + purchaseInvoice.VatAmount - purchaseInvoice.IncomeTaxValue - purchaseInvoice.RentValue

	result = trans.Save(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	insertTransactionalLog(enterpriseId, "purchase_invoice", int(invoiceId), userId, "U")
	json, _ := json.Marshal(purchaseInvoice)
	go fireWebHook(enterpriseId, "purchase_invoice", "PUT", string(json))

	return true
}

// Adds a income tax base amount to the income tax base. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addIncomeTaxBasePurchaseInvoice(invoiceId int64, totalAmount float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseInvoice PurchaseInvoice
	result := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).First(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseInvoice.IncomeTaxBase += totalAmount

	result = trans.Save(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// Adds a total amount to the invoice total. This function will subsctract from the total if the totalAmount is negative.
// THIS FUNCTION DOES NOT OPEN A TRANSACTION.
func addRentBaseProductsPurchaseInvoice(invoiceId int64, totalAmount float64, enterpriseId int32, userId int32, trans gorm.DB) bool {
	var purchaseInvoice PurchaseInvoice
	result := trans.Model(&PurchaseInvoice{}).Where("id = ?", invoiceId).First(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	purchaseInvoice.RentBase += totalAmount

	result = trans.Save(&purchaseInvoice)
	if result.Error != nil {
		log("DB", result.Error.Error())
		trans.Rollback()
		return false
	}

	return calcTotalsPurchaseInvoice(invoiceId, enterpriseId, userId, trans)
}

// ERROR CODES:
// 1. The order is already invoiced
// 2. There are no details to invoice
func invoiceAllPurchaseOrder(purchaseOrderId int64, enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the purchase order and it's details
	purchaseOrder := getPurchaseOrderRow(purchaseOrderId)
	if purchaseOrder.EnterpriseId != enterpriseId {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.Id <= 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.InvoicedLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}
	orderDetails := getPurchaseOrderDetail(purchaseOrderId, purchaseOrder.EnterpriseId)
	filterPurchaseOrderDetails(orderDetails, func(pod PurchaseOrderDetail) bool { return pod.QuantityInvoiced < pod.Quantity })
	if purchaseOrder.Id <= 0 || len(orderDetails) == 0 {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 2}
	}

	// create an invoice for that order
	invoice := PurchaseInvoice{}
	invoice.SupplierId = purchaseOrder.SupplierId
	invoice.BillingAddressId = purchaseOrder.BillingAddressId
	invoice.BillingSeriesId = purchaseOrder.BillingSeriesId
	invoice.CurrencyId = purchaseOrder.CurrencyId
	invoice.PaymentMethodId = purchaseOrder.PaymentMethodId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id, enterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.EnterpriseId = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(orderDetails); i++ {
		orderDetail := orderDetails[i]
		invoiceDetail := PurchaseInvoiceDetail{}
		invoiceDetail.InvoiceId = invoiceId
		invoiceDetail.OrderDetailId = &orderDetail.Id
		invoiceDetail.Price = orderDetail.Price
		invoiceDetail.ProductId = &orderDetail.ProductId
		invoiceDetail.Quantity = orderDetail.Quantity
		invoiceDetail.TotalAmount = orderDetail.TotalAmount
		invoiceDetail.VatPercent = orderDetail.VatPercent
		invoiceDetail.EnterpriseId = enterpriseId
		ok = invoiceDetail.insertPurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}
	///
}

// ERROR CODES:
// 1. The order is aleady invoiced
// 2. The selected quantity is greater than the quantity in the detail
// 3. The detail is already invoiced
// 4. The selected quantity is greater than the quantity pending of invoicing in the detail
func (invoiceInfo *OrderDetailGenerate) invoicePartiallyPurchaseOrder(enterpriseId int32, userId int32) OkAndErrorCodeReturn {
	// get the sale order and it's details
	purchaseOrder := getPurchaseOrderRow(invoiceInfo.OrderId)
	if purchaseOrder.Id <= 0 || purchaseOrder.EnterpriseId != enterpriseId || len(invoiceInfo.Selection) == 0 {
		return OkAndErrorCodeReturn{Ok: false}
	}
	if purchaseOrder.InvoicedLines >= purchaseOrder.LinesNumber {
		return OkAndErrorCodeReturn{Ok: false, ErrorCode: 1}
	}

	var purchaseOrderDetails []PurchaseOrderDetail = make([]PurchaseOrderDetail, 0)
	for i := 0; i < len(invoiceInfo.Selection); i++ {
		orderDetail := getPurchaseOrderDetailRow(invoiceInfo.Selection[i].Id)
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
		purchaseOrderDetails = append(purchaseOrderDetails, orderDetail)
	}

	// create an invoice for that order
	invoice := PurchaseInvoice{}
	invoice.SupplierId = purchaseOrder.SupplierId
	invoice.BillingAddressId = purchaseOrder.BillingAddressId
	invoice.BillingSeriesId = purchaseOrder.BillingSeriesId
	invoice.CurrencyId = purchaseOrder.CurrencyId
	invoice.PaymentMethodId = purchaseOrder.PaymentMethodId

	///
	trans := dbOrm.Begin()
	if trans.Error != nil {
		return OkAndErrorCodeReturn{Ok: false}
	}
	///

	ok := setDatePaymentAcceptedPurchaseOrder(purchaseOrder.Id, enterpriseId, userId, *trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}

	invoice.EnterpriseId = enterpriseId
	ok, invoiceId := invoice.insertPurchaseInvoice(userId, trans)
	if !ok {
		trans.Rollback()
		return OkAndErrorCodeReturn{Ok: false}
	}
	for i := 0; i < len(purchaseOrderDetails); i++ {
		orderDetail := purchaseOrderDetails[i]
		invoiceDetail := PurchaseInvoiceDetail{}
		invoiceDetail.InvoiceId = invoiceId
		invoiceDetail.OrderDetailId = &orderDetail.Id
		invoiceDetail.Price = orderDetail.Price
		invoiceDetail.ProductId = &orderDetail.ProductId
		invoiceDetail.Quantity = invoiceInfo.Selection[i].Quantity
		invoiceDetail.TotalAmount = orderDetail.TotalAmount
		invoiceDetail.VatPercent = orderDetail.VatPercent
		invoiceDetail.EnterpriseId = enterpriseId
		ok = invoiceDetail.insertPurchaseInvoiceDetail(userId, trans).Ok
		if !ok {
			trans.Rollback()
			return OkAndErrorCodeReturn{Ok: false}
		}
	}

	///
	trans.Commit()
	return OkAndErrorCodeReturn{Ok: true}
	///
}

type PurchaseInvoiceRelations struct {
	Orders   []PurchaseOrder   `json:"orders"`
	Invoices []PurchaseInvoice `json:"invoices"`
}

func getPurchaseInvoiceRelations(invoiceId int64, enterpriseId int32) PurchaseInvoiceRelations {
	return PurchaseInvoiceRelations{
		Orders:   getPurchaseInvoiceOrders(invoiceId, enterpriseId),
		Invoices: getPurchaseInvoiceAmendingAmendedInvoices(invoiceId, enterpriseId),
	}
}

func getPurchaseInvoiceOrders(invoiceId int64, enterpriseId int32) []PurchaseOrder {
	purchaseInvoice := getPurchaseInvoiceRow(invoiceId)
	if purchaseInvoice.Id <= 0 || purchaseInvoice.EnterpriseId != enterpriseId {
		return []PurchaseOrder{}
	}
	var orders []PurchaseOrder = make([]PurchaseOrder, 0)
	invoiceDetails := getPurchaseInvoiceDetail(invoiceId, enterpriseId)
	for i := 0; i < len(invoiceDetails); i++ {
		var orderDetail PurchaseOrderDetail
		if invoiceDetails[i].OrderDetailId != nil {
			orderDetail = getPurchaseOrderDetailRow(*invoiceDetails[i].OrderDetailId)
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
			order := getPurchaseOrderRow(orderDetail.OrderId)
			orders = append(orders, order)
		}
	}
	return orders
}

func getPurchaseInvoiceAmendingAmendedInvoices(invoiceId int64, enterpriseId int32) []PurchaseInvoice {
	invoices := make([]PurchaseInvoice, 0)

	i := getPurchaseInvoiceRow(invoiceId)
	if i.EnterpriseId != enterpriseId {
		return invoices
	}

	if i.Amending && i.AmendedInvoiceId != nil {
		invoices = append(invoices, getPurchaseInvoiceRow(*i.AmendedInvoiceId))
	}

	result := dbOrm.Model(&PurchaseInvoice{}).Where("amended_invoice = ?", i.Id).First(&invoices)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return invoices
	}

	return invoices
}
