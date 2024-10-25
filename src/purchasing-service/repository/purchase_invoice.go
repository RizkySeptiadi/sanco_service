package repository

import (
	"fmt"
	"purchasing_service/database"
	"purchasing_service/structs"

	"gorm.io/gorm"
)

type PurchaseInvoiceRepository[T any] struct {
	db     *gorm.DB
	module string
}

// Initialize GeneralRepository
func NewGeneralPurchaseInvoiceRepository[T any](module string) *PurchaseInvoiceRepository[T] {
	return &PurchaseInvoiceRepository[T]{
		db:     database.DbConnection,
		module: module,
	}
}

// PaginateDataSupp retrieves paginated data from the combined structs
func (repo *PurchaseInvoiceRepository[T]) PaginateDataSupp(start int, length int, search string, orderBy string, orderDir string, status string) (structs.PaginatedResponse[structs.PurchasingIndex], error) {
	var data []structs.PurchasingIndex
	var totalRecords int64
	var filteredRecords int64

	// Base query
	query := repo.db.Model(&structs.Sanco_Purchase_Invoices{}).
		Joins("JOIN sanco_purchase_invoice_numberings ON sanco_purchase_invoice_numberings.purchase_invoice_id = sanco_purchase_invoices.id").
		Joins("JOIN sanco_suppliers ON sanco_suppliers.id = sanco_purchase_invoices.supplier_id").
		Where("sanco_purchase_invoices.deleted_at IS NULL") // Explicitly check for deleted_at

	// Get total records count
	if err := query.Count(&totalRecords).Error; err != nil {
		return structs.PaginatedResponse[structs.PurchasingIndex]{
			Message: "Failed to retrieve total records: " + err.Error(),
			Code:    500,
		}, err
	}

	// Apply search filter if provided
	query = query.Select("sanco_purchase_invoices.*, sanco_purchase_invoice_numberings.number AS purchase_invoice_number, sanco_suppliers.name AS supplier_name, sanco_suppliers.id AS supplier_id")

	if search != "" {
		query = query.Where("sanco_purchase_invoice_numberings.number LIKE ?", "%"+search+"%").
			Or("sanco_suppliers.name LIKE ?", "%"+search+"%").
			Or("sanco_purchase_invoices.purchase_invoice_number LIKE ?", "%"+search+"%")
	}

	// Get filtered records count
	if err := query.Count(&filteredRecords).Error; err != nil {
		return structs.PaginatedResponse[structs.PurchasingIndex]{
			Message: "Failed to retrieve filtered records: " + err.Error(),
			Code:    500,
		}, err
	}

	// Apply ordering
	if orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))
	}

	// Apply pagination (start and length)
	if err := query.Offset(start).Limit(length).Find(&data).Error; err != nil {
		return structs.PaginatedResponse[structs.PurchasingIndex]{
			Message: "Failed to retrieve paginated data: " + err.Error(),
			Code:    500,
		}, err
	}

	// Return the paginated response
	return structs.PaginatedResponse[structs.PurchasingIndex]{
		TotalRecords:    totalRecords,
		FilteredRecords: filteredRecords,
		Data:            data,
		Message:         "Data retrieved successfully",
		Code:            200,
	}, nil
}

func (repo *PurchaseInvoiceRepository[T]) GetDataDetailByID(Purchase_invoice_id int64) (structs.Response[T], error) {
	var data T
	// Retrieve the record by ID without checking for soft deletes
	if err := repo.db.Where("id = ?", Purchase_invoice_id).First(&data).Error; err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}
func (repo *PurchaseInvoiceRepository[T]) CreatePurchase(tx *gorm.DB, parent *structs.Sanco_Purchase_Invoices) error {
	return tx.Create(parent).Error
}

// CreateDetail inserts the detail record into the database
func (repo *PurchaseInvoiceRepository[T]) CreateDetail(tx *gorm.DB, detail *structs.Sanco_Purchase_Invoice_details) error {
	return tx.Create(detail).Error
}
func (repo *PurchaseInvoiceRepository[T]) UpdatePurchase(tx *gorm.DB, purchase *structs.Sanco_Purchase_Invoices) error {
	// Use the Updates method to specify the fields to update
	// Create a temporary struct to hold the fields you want to update
	updateData := structs.Sanco_Purchase_Invoices{
		Supplier_id:             purchase.Supplier_id,
		Purchase_invoice_number: purchase.Purchase_invoice_number,
		Quantity:                purchase.Quantity,
		Total:                   purchase.Total,
		Discount:                purchase.Discount,
		Tax_percentage:          purchase.Tax_percentage,
		Tax:                     purchase.Tax,
		Grand_total:             purchase.Grand_total,
		Total_payment:           purchase.Total_payment,
		Post:                    purchase.Post,
		Account_payable_cart:    purchase.Account_payable_cart,
		// Exclude created_at and updated_at
	}

	// Update the purchase record while excluding created_at and updated_at
	return tx.Model(&structs.Sanco_Purchase_Invoices{}).Where("id = ?", purchase.ID).Updates(updateData).Error
}

func (repo *PurchaseInvoiceRepository[T]) UpdateDetail(tx *gorm.DB, detail *structs.Sanco_Purchase_Invoice_details) error {
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
func (repo *PurchaseInvoiceRepository[T]) DeleteDetail(tx *gorm.DB, detailID int) error {
	return tx.Where("id = ?", detailID).Delete(&structs.Sanco_Purchase_Invoice_details{}).Error
}

// GetDetailByID retrieves a detail record by its ID
func (repo *PurchaseInvoiceRepository[T]) GetByID(id int64) (structs.Response[T], error) {
	var data T
	err := repo.db.First(&data, id).Error
	if err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}
func (repo *PurchaseInvoiceRepository[T]) GetDetailByID(tx *gorm.DB, id int64) (*structs.Sanco_Purchase_Invoice_details, error) {
	var detail structs.Sanco_Purchase_Invoice_details
	err := tx.First(&detail, id).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
func (repo *PurchaseInvoiceRepository[T]) GetAllDataDetailByID(Purchase_invoice_id int64) (structs.Response[[]T], error) {
	var data []T
	// Retrieve the record by ID without checking for soft deletes
	if err := repo.db.Where("purchase_invoice_id = ?", Purchase_invoice_id).Find(&data).Error; err != nil {
		return structs.NewResponse(500, "Failed to retrieve "+repo.module, []T{}), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}
func (repo *PurchaseInvoiceRepository[T]) CreateData(data *T) (structs.Response[T], error) {
	if err := repo.db.Create(data).Error; err != nil {
		return structs.NewResponse(500, "Failed to create "+repo.module, *new(T)), err
	}

	return structs.NewResponse(201, repo.module+" created successfully", *data), nil
}
func (repo *PurchaseInvoiceRepository[T]) GenerateNewNumber() (string, error) {
	var newNumber string

	// Build the custom SQL query to generate the new number
	query := `
		SELECT CONCAT(
			"SCO",
			DATE_FORMAT(CURRENT_DATE, "%y%m"), "PI",
			LPAD(IFNULL(MAX(CAST(SUBSTRING(number, 10, 4) AS UNSIGNED)), 0) + 1, 4, "0")
		) AS new_number
		FROM sanco_purchase_invoice_numberings
		WHERE SUBSTRING(number, 4, 4) = DATE_FORMAT(CURRENT_DATE, "%y%m")
	`

	// Execute the raw SQL query and retrieve the result into newNumber
	err := repo.db.Raw(query).Scan(&newNumber).Error
	if err != nil {
		return "", err
	}

	return newNumber, nil
}
