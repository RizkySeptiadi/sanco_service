package repository

import (
	"accounting_service/database"
	"accounting_service/structs"

	"gorm.io/gorm"
)

type SalesInvoiceRepository[T any] struct {
	db     *gorm.DB
	module string
}

// Initialize GeneralRepository
func NewGeneralSalesInvoiceRepository[T any](module string) *SalesInvoiceRepository[T] {
	return &SalesInvoiceRepository[T]{
		db:     database.DbConnection,
		module: module,
	}
}

func (repo *SalesInvoiceRepository[T]) GetDataDetailByID(Purchase_invoice_id int64) (structs.Response[T], error) {
	var data T
	// Retrieve the record by ID without checking for soft deletes
	if err := repo.db.Where("id = ?", Purchase_invoice_id).First(&data).Error; err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}
func (repo *SalesInvoiceRepository[T]) CreatePurchase(tx *gorm.DB, parent *structs.SalesInvoice) error {
	return tx.Create(parent).Error
}

// CreateDetail inserts the detail record into the database
func (repo *SalesInvoiceRepository[T]) CreateDetail(tx *gorm.DB, detail *structs.SalesInvoiceItem) error {
	return tx.Create(detail).Error
}

// func (repo *SalesInvoiceRepository[T]) UpdatePurchase(tx *gorm.DB, purchase *structs.SalesInvoice) error {
// 	// Use the Updates method to specify the fields to update
// 	// Create a temporary struct to hold the fields you want to update
// 	updateData := structs.SalesInvoice{
// 		Supplier_id:             purchase.Supplier_id,
// 		Purchase_invoice_number: purchase.Purchase_invoice_number,
// 		Quantity:                purchase.Quantity,
// 		Total:                   purchase.Total,
// 		Discount:                purchase.Discount,
// 		Tax_percentage:          purchase.Tax_percentage,
// 		Tax:                     purchase.Tax,
// 		Grand_total:             purchase.Grand_total,
// 		Total_payment:           purchase.Total_payment,
// 		Post:                    purchase.Post,
// 		Account_payable_cart:    purchase.Account_payable_cart,
// 		// Exclude created_at and updated_at
// 	}

// 	// Update the purchase record while excluding created_at and updated_at
// 	return tx.Model(&structs.SalesInvoice{}).Where("id = ?", purchase.ID).Updates(updateData).Error
// }

// func (repo *SalesInvoiceRepository[T]) UpdateDetail(tx *gorm.DB, detail *structs.Sanco_Purchase_Invoice_details) error {
// 	// Create a temporary struct to hold the fields you want to update
// 	updateData := structs.Sanco_Purchase_Invoice_details{
// 		Supplier_id:         detail.Supplier_id,
// 		Purchase_invoice_id: detail.Purchase_invoice_id,
// 		Pn:                  detail.Pn,
// 		Pname:               detail.Pname,
// 		Quantity:            detail.Quantity,
// 		Price:               detail.Price,
// 		Discount:            detail.Discount,
// 		Subtotal:            detail.Subtotal,
// 		// Exclude created_at and updated_at
// 	}

//		// Update the detail record while excluding created_at and updated_at
//		return tx.Model(&structs.Sanco_Purchase_Invoice_details{}).Where("id = ?", detail.ID).Updates(updateData).Error
//	}
// func (repo *SalesInvoiceRepository[T]) DeleteDetail(tx *gorm.DB, detailID int) error {
// 	return tx.Where("id = ?", detailID).Delete(&structs.Sanco_Purchase_Invoice_details{}).Error
// }

// GetDetailByID retrieves a detail record by its ID
func (repo *SalesInvoiceRepository[T]) GetByID(id int64) (structs.Response[T], error) {
	var data T
	err := repo.db.First(&data, id).Error
	if err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}

func (repo *SalesInvoiceRepository[T]) CreateData(data *T) (structs.Response[T], error) {
	if err := repo.db.Create(data).Error; err != nil {
		return structs.NewResponse(500, "Failed to create "+repo.module, *new(T)), err
	}

	return structs.NewResponse(201, repo.module+" created successfully", *data), nil
}
