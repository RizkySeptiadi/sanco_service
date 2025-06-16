package repository

import (
	"accounting_service/database"
	"accounting_service/structs"
	"errors"
	"fmt"
	"reflect"

	"time"

	"gorm.io/gorm"
)

type GeneralRepository[T any] struct {
	db     *gorm.DB
	module string
}

// Initialize GeneralRepository
func NewGeneralRepository[T any](module string) *GeneralRepository[T] {
	return &GeneralRepository[T]{
		db:     database.DbConnection,
		module: module,
	}
}

//DON'T TOUCH THESE UNLESS YOU KNOW WHAT ARE YOU DOING

// C R U D    O P E R A T I O N S

// Read All Datas
func (repo *GeneralRepository[T]) GetAllData() (structs.Response[[]T], error) {
	var data []T

	if err := repo.db.Find(&data).Error; err != nil {
		return structs.NewResponse(500, "Failed to retrieve "+repo.module, []T{}), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}

// Read one Data by ID
func (repo *GeneralRepository[T]) GetDataByID(id int64) (structs.Response[T], error) {
	var data T
	if err := repo.db.First(&data, id).Error; err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" retrieved successfully", data), nil
}

// Create Data
func (repo *GeneralRepository[T]) CreateData(data *T) (structs.Response[T], error) {
	if err := repo.db.Create(data).Error; err != nil {
		return structs.NewResponse(500, "Failed to create "+repo.module, *new(T)), err
	}

	return structs.NewResponse(201, repo.module+" created successfully", *data), nil
}

// Update Data
func (repo *GeneralRepository[T]) UpdateData(id int64, updatedData *T) (structs.Response[T], error) {
	var existingData T

	// Use reflection to get ref_id and source_id fields from the updatedData
	val := reflect.ValueOf(*updatedData)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	refIDField := val.FieldByName("RefID")
	sourceIDField := val.FieldByName("SourceID")

	if !refIDField.IsValid() || !sourceIDField.IsValid() {
		return structs.NewResponse[T](400, "Missing RefID or SourceID field", *new(T)), fmt.Errorf("required fields not found")
	}

	refID := refIDField.Interface()
	sourceID := sourceIDField.Interface()

	// Use WHERE ref_id + source_id for record lookup
	if err := repo.db.Where("ref_id = ? AND source_id = ?", refID, sourceID).First(&existingData).Error; err != nil {
		return structs.NewResponse[T](404, repo.module+" not found", *new(T)), nil
	}

	// Perform the update
	if err := repo.db.Model(&existingData).Updates(updatedData).Error; err != nil {
		return structs.NewResponse[T](500, "Failed to update "+repo.module, *new(T)), err
	}

	return structs.NewResponse[T](200, repo.module+" updated successfully", *updatedData), nil
}

// func (repo *GeneralRepository[T]) UpdateData(id int64, updatedData *T) (structs.Response[T], error) {
// 	var existingData T
// 	if err := repo.db.First(&existingData, id).Error; err != nil {
// 		// Return a response with error details when the record is not found
// 		return structs.NewResponse[T](404, repo.module+" not found", *new(T)), nil
// 	}

// 	if err := repo.db.Model(&existingData).Updates(updatedData).Error; err != nil {
// 		// Return a response with error details for failed update
// 		return structs.NewResponse[T](500, "Failed to update "+repo.module, *new(T)), err
// 	}

// 	// Return a success response with updated data
// 	return structs.NewResponse[T](200, repo.module+" updated successfully", *updatedData), nil
// }

// Soft Delete
func (repo *GeneralRepository[T]) DeleteData(id int64) (structs.Response[T], error) {
	var data T
	relatedTable := "sanco_purchase_invoices"
	// Find the record to delete
	if err := repo.db.First(&data, id).Error; err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}

	// Check if there are any related records in the relatedTable
	var relatedCount int64
	if err := repo.db.Table(relatedTable).Where("supplier_id = ?", id).Count(&relatedCount).Error; err != nil {
		return structs.NewResponse(500, "Failed to check related records in "+relatedTable, *new(T)), err
	}

	// If related records exist, prevent deletion
	if relatedCount > 0 {
		return structs.NewResponse(400, "Cannot delete data because related records exist in purchase invoices.", *new(T)), nil
	}

	// Perform the soft delete by updating the 'deleted_at' column
	if err := repo.db.Model(&data).Update("deleted_at", time.Now()).Error; err != nil {
		return structs.NewResponse(500, "Failed to soft delete "+repo.module, *new(T)), err
	}

	return structs.NewResponse(200, repo.module+" soft deleted successfully", data), nil
}

func (repo *GeneralRepository[T]) UpdateState(id int64) (structs.Response[T], error) {
	var data T
	if err := repo.db.First(&data, id).Error; err != nil {
		return structs.NewResponse(404, repo.module+" not found", *new(T)), err
	}
	var currentStatus int
	if err := repo.db.Model(&data).Select("status").Where("id = ?", id).Scan(&currentStatus).Error; err != nil {
		return structs.NewResponse(500, "Failed to retrieve current status", *new(T)), err
	}
	newStatus := 0
	if currentStatus != 1 {
		newStatus = 1
	}
	if err := repo.db.Model(&data).Update("status", newStatus).Error; err != nil {
		return structs.NewResponse(500, "Failed to update status", *new(T)), err
	}
	return structs.NewResponse(200, repo.module+" status updated successfully", data), nil
}
func (r *GeneralRepository[T]) UpsertBatch(dataList *[]T) (*structs.Response[T], error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return &structs.Response[T]{Code: 500, Message: "Failed to begin transaction"}, tx.Error
	}

	for _, data := range *dataList {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		companyIDField := val.FieldByName("RefID")
		sourceIDField := val.FieldByName("SourceID")

		if !companyIDField.IsValid() || !sourceIDField.IsValid() {
			tx.Rollback()
			return &structs.Response[T]{
				Code:    400,
				Message: "Missing RefID or SourceID field",
			}, fmt.Errorf("required fields not found")
		}

		refID := companyIDField.Interface()
		sourceID := sourceIDField.Interface()

		var existing T
		err := tx.Where("ref_id = ? AND source_id = ?", refID, sourceID).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Insert
				if err := tx.Create(&data).Error; err != nil {
					tx.Rollback()
					return &structs.Response[T]{Code: 500, Message: "Failed to insert"}, err
				}
			} else {
				tx.Rollback()
				return &structs.Response[T]{Code: 500, Message: "Failed to query"}, err
			}
		} else {
			// Update
			if err := tx.Model(&existing).Updates(data).Error; err != nil {
				tx.Rollback()
				return &structs.Response[T]{Code: 500, Message: "Failed to update"}, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &structs.Response[T]{Code: 500, Message: "Failed to commit transaction"}, err
	}

	return &structs.Response[T]{Code: 200, Message: "Batch upsert successful"}, nil
}
