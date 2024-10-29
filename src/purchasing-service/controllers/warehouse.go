package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"purchasing_service/database"
	"purchasing_service/repository"
	"purchasing_service/structs"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseController[T any] struct {
	repo *repository.WarehouseRepository[T]
	db   *gorm.DB
}

func NewWarehouseController[T any](repo *repository.WarehouseRepository[T]) *WarehouseController[T] {
	return &WarehouseController[T]{db: database.DbConnection,
		repo: repo}
}

func (ctrl *WarehouseController[T]) Create(c *gin.Context) {
	var requestData structs.RequestData

	// Bind the incoming JSON to the RequestData struct

	// Begin transaction
	tx := ctrl.db.Begin()

	parentID := c.Param("id") // Get the ID parameter from the request
	parentDataURL := fmt.Sprintf("http://localhost:8080/api/purchase_invoice/get_parent/%s", parentID)

	resp, err := http.Get(parentDataURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to fetch parent data", "error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Decode the parent data into the parentResponse.Parent struct
	var parentResponse struct {
		Parent structs.Sanco_Purchase_Invoices // define the struct with relevant fields
	}
	if err := json.NewDecoder(resp.Body).Decode(&parentResponse.Parent); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to decode parent data", "error": err.Error()})
		return
	}

	// Store the parent data
	if err := ctrl.repo.Create(tx, &parentResponse.Parent); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to store parent data", "error": err.Error()})
		return
	}

	// Now that the parent has been inserted, the ID is populated
	convertedID, err := strconv.ParseInt(parentID, 10, 64)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Invalid ID", "error": err.Error()})
		return
	}

	// Extract details into a slice for easier iteration
	var details []structs.Sanco_Purchase_Invoice_details
	for _, detail := range requestData.Details {
		detail.Purchase_invoice_id = convertedID // Assign the newly created Purchase Invoice ID
		details = append(details, detail)
	}

	// Insert each detail record into the database
	for _, detail := range details {
		if err := ctrl.repo.CreateDetail(tx, &detail); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed To store detail data", "error": err.Error()})

			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed To Commit transaction, pleasy try again !", "error": err.Error()})

		return
	}

	// On success, return the successful response

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Purchase invoice created successfully", "data": requestData})

}
