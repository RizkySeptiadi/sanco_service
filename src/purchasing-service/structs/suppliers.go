package structs

import (
	"time"

	"gorm.io/gorm"
)

type Date struct {
	time.Time
}

type Sanco_Purchase_Invoices struct {
	gorm.Model
	ID           int64     `json:"ID"`
	Supplier_id  int64     `gorm:"unique" binding:"required" json:"Supplier_id"`
	Date         time.Time `json:"Date" gorm:"type:date"` // Date only
	NumberManual string    `json:"NumberManual"`

	Purchase_invoice_number string  `json:"Purchase_invoice_number"`
	Quantity                float32 `binding:"required" json:"Quantity"`
	Total                   float32 `binding:"required" json:"Total"`
	Discount                float32 `binding:"required" json:"Discount"`
	Tax_percentage          float32 `binding:"required" json:"Tax_percentage"`
	Tax                     float32 `json:"Tax"`
	Grand_total             float32 `json:"Grand_total"`
	Total_payment           float32 `json:"Total_payment"`
	Post                    int8    `json:"Post"`
	Account_payable_cart    int8    `json:"Account_payable_cart"`
}

type Sanco_Purchase_Invoice_details struct {
	gorm.Model
	ID                  int64   `json:"ID"`
	Supplier_id         int64   `gorm:"unique" binding:"required" json:"Supplier_id"`
	Purchase_invoice_id int64   `gorm:"not null" binding:"required" json:"Purchase_invoice_id"`
	Pn                  string  `json:"Pn"`
	Pname               string  `json:"Pname"`
	Quantity            float32 `binding:"required" json:"Quantity"`
	Price               float32 `binding:"required" json:"Price"`
	Discount            float32 `binding:"required" json:"Discount"`
	Subtotal            float32 `binding:"required" json:"Subtotal"`
}
type Sanco_purchase_invoice_numberings struct {
	gorm.Model
	ID                  int64
	Purchase_invoice_id int64  `gorm:"unique"  binding:"required"`
	Number              string `binding:"required"`
}
type PurchasingIndex struct {
	gorm.Model
	SupplierID            int64   // ID from Sanco_Suppliers
	SupplierCode          string  // Code from Sanco_Suppliers
	Date                  string  // Code from Sanco_Suppliers
	NumberManual          string  // Code from Sanco_Suppliers
	SupplierName          string  // Name from Sanco_Suppliers
	PurchaseInvoiceID     int64   // ID from Sanco_Purchase_Invoices
	PurchaseInvoiceNumber string  // Purchase invoice number from Sanco_Purchase_Invoices
	Quantity              int     // Quantity from Sanco_Purchase_Invoices
	Total                 float32 // Total from Sanco_Purchase_Invoices
	Discount              float32 // Discount from Sanco_Purchase_Invoices
	TaxPercentage         float32 // Tax percentage from Sanco_Purchase_Invoices
	Tax                   float32 // Tax from Sanco_Purchase_Invoices
	GrandTotal            float32 // Grand total from Sanco_Purchase_Invoices
	TotalPayment          float32 // Total payment from Sanco_Purchase_Invoices
	TotalRemaining        float32 // Total remaining from Sanco_Purchase_Invoices
	Post                  int8    // Post from Sanco_Purchase_Invoices
	AccountPayableCart    int8    // Account payable cart from Sanco_Purchase_Invoices
}

type CustomDate struct {
	time.Time
}

// Custom UnmarshalJSON to handle the date format "2006-01-02"
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	// Remove quotes from JSON string
	str := string(b)
	str = str[1 : len(str)-1]

	// Parse the date
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}

	*cd = CustomDate{date}
	return nil
}

type RequestData struct {
	Parent        Sanco_Purchase_Invoices                   `json:"parent"`
	Details       map[string]Sanco_Purchase_Invoice_details `json:"detail"`
	DeletedDetail map[string]int                            `json:"deleted_detail"`
}
