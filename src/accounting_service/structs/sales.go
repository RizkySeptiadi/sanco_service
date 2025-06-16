package structs

import (
	"time"

	"gorm.io/gorm"
)

type Date struct {
	time.Time
}

type SalesInvoice struct {
	gorm.Model
	ID                     int64  `json:"id"`
	SalesInvoiceRefID      int64  `json:"sales_invoice_ref_id" binding:"required"`
	SourceID               string `json:"source_id" binding:"required"`
	SalesOrderID           int64  `json:"sales_order_id" binding:"required"`
	Number                 string `json:"number" gorm:"unique" binding:"required"`
	Date                   string `json:"date" binding:"required" gorm:"type:date"`
	DebitChartOfAccountID  int64  `json:"debit_chart_of_account_id" binding:"required"`
	CreditChartOfAccountID int64  `json:"credit_chart_of_account_id" binding:"required"`
	CustomerID             int64  `json:"customer_id" binding:"required"`
	CustomerName           string `json:"customer_name" binding:"required"`

	CompanyID    *int64  `json:"company_id"`    // Nullable
	CompanyName  *string `json:"company_name"`  // Nullable
	SalesmanID   *int64  `json:"salesman_id"`   // Nullable
	SalesmanName *string `json:"salesman_name"` // Nullable

	TotalQuantity      int     `json:"total_quantity" binding:"required"`
	Subtotal           int     `json:"subtotal" binding:"required"`
	Discount           int     `json:"discount"`
	DiscountPercentage float64 `json:"discount_percentage"`
	TaxID              uint64  `json:"tax_id" binding:"required"`
	Tax                int     `json:"tax" binding:"required"`
	TaxPercentage      float64 `json:"tax_percentage"`
	Total              int     `json:"total" binding:"required"`
	TotalRemaining     float64 `json:"total_remaining"`
	TotalPaid          float64 `json:"total_paid"`

	Status string  `json:"status" gorm:"type:enum('draft','post')"`
	Note   *string `json:"note"` // Nullable

}
type SalesInvoiceItem struct {
	ID                int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SalesInvoiceID    int64    `json:"sales_invoice_id" binding:"required"`
	SalesInvoiceRefID int64    `json:"sales_invoice_ref_id" binding:"required"`
	SourceID          string   `json:"source_id" binding:"required"`
	ItemID            int64    `json:"item_id" binding:"required"`
	ItemName          string   `json:"item_name" binding:"required"`
	ItemPrice         float64  `json:"item_price" binding:"required"`
	AssemblyID        *int64   `json:"assembly_id"` // Nullable
	PN                *string  `json:"pn"`          // Nullable
	PNAlias           *string  `json:"pn_alias"`    // Nullable
	PNAME             *string  `json:"pname"`       // Nullable
	PNAMEAlias        *string  `json:"pname_alias"` // Nullable
	Unit              *string  `json:"unit"`        // Nullable
	ItemPriceAlias    *float64 `json:"item_price_alias"`

	Quantity           int     `json:"quantity" binding:"required"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Discount           float64 `json:"discount"`
	TaxPercentage      float64 `json:"tax_percentage"`
	Tax                float64 `json:"tax"`
	Total              float64 `json:"total"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestData struct {
	Parent        SalesInvoice                `json:"parent"`
	Details       map[string]SalesInvoiceItem `json:"detail"`
	DeletedDetail map[string]int              `json:"deleted_detail"`
}
