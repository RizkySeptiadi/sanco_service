package structs

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseInvoice struct {
	gorm.Model
	ID                     uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseInvoiceRefID   int64     `json:"purchase_invoice_ref_id"`
	SourceID               string    `json:"source_id"`
	PurchaseOrderID        uint64    `json:"purchase_order_id" binding:"required"`
	Number                 string    `json:"number" binding:"required"`
	Date                   time.Time `json:"date" binding:"required"`
	DebitChartOfAccountID  uint64    `json:"debit_chart_of_account_id" binding:"required"`
	CreditChartOfAccountID uint64    `json:"credit_chart_of_account_id" binding:"required"`
	SupplierID             uint64    `json:"supplier_id" binding:"required"`
	SupplierName           string    `json:"supplier_name" binding:"required"`
	CompanyID              uint64    `json:"company_id"`
	CompanyName            string    `json:"company_name"`
	SalesmanID             uint64    `json:"salesman_id"`
	SalesmanName           string    `json:"salesman_name"`
	TotalQuantity          int       `json:"total_quantity" binding:"required"`
	Subtotal               int       `json:"subtotal" binding:"required"`
	Discount               int       `json:"discount"`
	DiscountPercentage     float64   `json:"discount_percentage"`
	TaxID                  uint64    `json:"tax_id" binding:"required"`
	Tax                    int       `json:"tax" binding:"required"`
	TaxPercentage          float64   `json:"tax_percentage"`
	Total                  int       `json:"total" binding:"required"`
	TotalRemaining         float64   `json:"total_remaining"`
	TotalPaid              float64   `json:"total_paid"`
	Status                 string    `json:"status" binding:"required"` // enum: draft/post
	Note                   string    `json:"note"`
}

type PurchaseInvoiceItem struct {
	gorm.Model
	ID                   uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseInvoiceRefID int64      `json:"purchase_invoice_ref_id"`
	SourceID             string     `json:"source_id"`
	PurchaseInvoiceID    uint64     `json:"purchase_invoice_id" binding:"required"`
	ItemID               uint64     `json:"item_id" binding:"required"`
	ItemName             string     `json:"item_name" binding:"required"`
	ItemPrice            float64    `json:"item_price" binding:"required"`
	PN                   string     `json:"pn"`
	PNAlias              string     `json:"pn_alias"`
	Pname                string     `json:"pname"`
	PnameAlias           string     `json:"pname_alias"`
	Unit                 string     `json:"unit"`
	ItemPriceAlias       float64    `json:"item_price_alias"`
	Quantity             int        `json:"quantity" binding:"required"`
	DiscountPercentage   float64    `json:"discount_percentage" binding:"required"`
	Discount             float64    `json:"discount" binding:"required"`
	TaxPercentage        float64    `json:"tax_percentage" binding:"required"`
	Tax                  float64    `json:"tax" binding:"required"`
	Total                float64    `json:"total" binding:"required"`
	CreatedAt            *time.Time `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	DeletedAt            *time.Time `gorm:"index" json:"deleted_at"`
}

// TableName overrides the table name used by GORM
func (PurchaseInvoice) TableName() string {
	return "purchase_invoices"
}

type RequestDataP struct {
	Parent        PurchaseInvoice                `json:"parent"`
	Details       map[string]PurchaseInvoiceItem `json:"detail"`
	DeletedDetail map[string]int                 `json:"deleted_detail"`
}
