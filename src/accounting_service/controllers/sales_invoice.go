package controllers

import (
	"accounting_service/database"
	"accounting_service/repository"
	"accounting_service/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SalesInvoiceController[T any] struct {
	repo *repository.SalesInvoiceRepository[T]
	db   *gorm.DB
}

func NewSalesInvoiceController[T any](repo *repository.SalesInvoiceRepository[T]) *SalesInvoiceController[T] {
	return &SalesInvoiceController[T]{db: database.DbConnection,
		repo: repo}
}

func (ctrl *SalesInvoiceController[T]) Create(c *gin.Context) {
	var requestData structs.RequestData

	// Bind the incoming JSON to the RequestData struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": err.Error()})
		return
	}

	// Begin transaction
	tx := ctrl.db.Begin()

	if err := ctrl.repo.CreatePurchase(tx, &requestData.Parent); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed To store parent data", "error": err.Error()})

		return
	}

	// Now that the parent has been inserted, the ID is populated
	parentID := requestData.Parent.ID

	// Extract details into a slice for easier iteration
	var details []structs.SalesInvoiceItem
	for _, detail := range requestData.Details {
		detail.SalesInvoiceID = parentID // Assign the newly created Purchase Invoice ID
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
