package repository

import (
	"fmt"
	"sanco_microservices/structs"
)

func (repo *GeneralRepository[T]) PaginateData(start int, length int, search string, orderBy string, orderDir string) (structs.PaginatedResponse[T], error) {
	var data []T
	var totalRecords int64
	var filteredRecords int64

	// Get total records count
	if err := repo.db.Model(&data).Count(&totalRecords).Error; err != nil {
		return structs.PaginatedResponse[T]{
			Message: "Failed to retrieve total records" + err.Error(),
			Code:    500,
		}, err
	}

	// Apply search filter if provided
	query := repo.db.Model(&data)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Get filtered records count
	if err := query.Count(&filteredRecords).Error; err != nil {
		return structs.PaginatedResponse[T]{
			Message: "Failed to retrieve filtered records" + err.Error(),
			Code:    500,
		}, err
	}

	// Apply ordering
	if orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))
	}

	// Apply pagination (start and length)
	if err := query.Offset(start).Limit(length).Find(&data).Error; err != nil {
		return structs.PaginatedResponse[T]{
			Message: "Failed to retrieve paginated data" + err.Error(),
			Code:    500,
		}, err
	}

	// Return the paginated response with success message
	return structs.PaginatedResponse[T]{
		TotalRecords:    totalRecords,
		FilteredRecords: filteredRecords,
		Data:            data,
		Message:         "Data retrieved successfully",
		Code:            200, // No error, leave empty
	}, nil
}
