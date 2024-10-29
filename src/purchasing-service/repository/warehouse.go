package repository

import (
	"purchasing_service/database"
	"purchasing_service/structs"

	"gorm.io/gorm"
)

type WarehouseRepository[T any] struct {
	db     *gorm.DB
	module string
}

// Initialize GeneralRepository
func NewGeneralWarehouseRepository[T any](module string) *WarehouseRepository[T] {
	return &WarehouseRepository[T]{
		db:     database.DbConnection,
		module: module,
	}
}

// func (repo *WarehouseRepository[T]) Create(tx *gorm.DB, parent *structs.Sanco_Purchase_Invoices) error {
// 	return tx.Create(parent).Error
// }

// CreateDetail inserts the detail record into the database
func (repo *WarehouseRepository[T]) CreateDetail(tx *gorm.DB, detail *structs.Sanco_Purchase_Invoice_details) error {
	return tx.Create(detail).Error
}
func (repo *WarehouseRepository[T]) Create(tx *gorm.DB, purchase *structs.Sanco_Purchase_Invoices) error {
	// Prepare data for insertion
	newData := structs.Incoming_Order{
		IDSupplier: purchase.Supplier_id,
		IDPo:       purchase.Purchase_invoice_number,
		PoDate:     purchase.Date,
		Tgljt:      purchase.Date,

		TotalQty:   purchase.Quantity,
		TotalPrice: purchase.Total,
		TotalDisc:  purchase.Discount,
		PpnPercent: purchase.Tax_percentage,
		TotalPpn:   purchase.Tax,
		GrandTotal: purchase.Grand_total,

		State: "0",
		// created_at and updated_at will be automatically handled by GORM
	}

	// Insert a new record into the Incoming_Order table
	return tx.Create(&newData).Error
}
func (repo *WarehouseRepository[T]) InsertDetail(tx *gorm.DB, invoice *structs.Sanco_Purchase_Invoices, detail *structs.Sanco_Purchase_Invoice_details) error {
	// Prepare the new IncomingOrderDetail data
	incomingDetail := structs.Incoming_order_detail{
		IDPo: invoice.Purchase_invoice_number,
		// SupplierName: fmt.Sprintf("%s (Sanko)", detail.SupplierName),
		PoDate:   invoice.Date,
		PN:       detail.Pn,
		PNas:     detail.Pn,
		PNAME:    detail.Pname,
		PNAMEas:  detail.Pname,
		Qty:      detail.Quantity,
		QtyReal:  0,
		QtyIn:    0,
		QtyOrder: detail.Quantity,
		Price:    0,
		PriceUSD: 0,
		// User:     fmt.Sprintf("%s (Sanko)", auth.GetUser().Name), // Assuming `auth.GetUser()` retrieves the authenticated user
	}

	// Insert the new detail record
	return tx.Create(&incomingDetail).Error
}
func (repo *WarehouseRepository[T]) UpdateDetail(tx *gorm.DB, detail *structs.Sanco_Purchase_Invoice_details) error {
	// Create a temporary struct to hold the fields you want to update
	updateData := structs.Sanco_Purchase_Invoice_details{
		Supplier_id:         detail.Supplier_id,
		Purchase_invoice_id: detail.Purchase_invoice_id,
		Pn:                  detail.Pn,
		Pname:               detail.Pname,
		Quantity:            detail.Quantity,
		Price:               detail.Price,
		Discount:            detail.Discount,
		Subtotal:            detail.Subtotal,
		// Exclude created_at and updated_at
	}

	// Update the detail record while excluding created_at and updated_at
	return tx.Model(&structs.Sanco_Purchase_Invoice_details{}).Where("id = ?", detail.ID).Updates(updateData).Error
}
