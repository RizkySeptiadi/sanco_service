package controllers

import (
	"net/http"
	"purchasing_service/database"
	"purchasing_service/repository"
	"purchasing_service/structs"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseInvoiceController[T any] struct {
	repo *repository.PurchaseInvoiceRepository[T]
	db   *gorm.DB
}

func NewPurchaseInvoiceController[T any](repo *repository.PurchaseInvoiceRepository[T]) *PurchaseInvoiceController[T] {
	return &PurchaseInvoiceController[T]{db: database.DbConnection,
		repo: repo}
}
func (ctrl *PurchaseInvoiceController[T]) GetTables(c *gin.Context) {
	// Get query parameters for pagination, sorting, and searching
	start, err := strconv.Atoi(c.DefaultQuery("start", "0")) // Starting index
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": "Invalid start parameter"})
		return
	}

	length, err := strconv.Atoi(c.DefaultQuery("length", "10")) // Number of records to fetch
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": "Invalid length parameter"})
		return
	}

	search := c.DefaultQuery("search", "")        // Search term
	orderBy := c.DefaultQuery("orderBy", "id")    // Column to order by
	orderDir := c.DefaultQuery("orderDir", "asc") // Order direction (asc or desc)
	status := c.DefaultQuery("status", "")

	// Call the repository to fetch paginated data
	response, err := ctrl.repo.PaginateDataSupp(start, length, search, orderBy, orderDir, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated response
	c.JSON(http.StatusOK, response)
}
func (ctrl *PurchaseInvoiceController[T]) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.GetByID(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	c.JSON(response.Code, response)
}
func (ctrl *PurchaseInvoiceController[T]) GetDataDetailByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.GetDataDetailByID(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	c.JSON(response.Code, response)
}
func (ctrl *PurchaseInvoiceController[T]) GetAllDataDetailByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.GetAllDataDetailByID(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	c.JSON(response.Code, response)
}

func (ctrl *PurchaseInvoiceController[T]) Create(c *gin.Context) {
	var requestData structs.RequestData

	// Bind the incoming JSON to the RequestData struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": err.Error()})
		return
	}

	// Begin transaction
	tx := ctrl.db.Begin()

	newNumber, err := ctrl.repo.GenerateNewNumber()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to generate purchase invoice number", "error": err.Error()})
		return
	}
	requestData.Parent.Purchase_invoice_number = newNumber
	if err := ctrl.repo.CreatePurchase(tx, &requestData.Parent); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed To store parent data", "error": err.Error()})

		return
	}
	invoiceNumbering := structs.Sanco_purchase_invoice_numberings{
		Purchase_invoice_id: requestData.Parent.ID,
		Number:              newNumber,
	}
	if err := tx.Create(&invoiceNumbering).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to save invoice numbering", "error": err.Error()})
		return
	}

	// Now that the parent has been inserted, the ID is populated
	parentID := requestData.Parent.ID

	// Extract details into a slice for easier iteration
	var details []structs.Sanco_Purchase_Invoice_details
	for _, detail := range requestData.Details {
		detail.Purchase_invoice_id = parentID // Assign the newly created Purchase Invoice ID
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
func (ctrl *PurchaseInvoiceController[T]) Update(c *gin.Context) {
	var requestData structs.RequestData

	// Get the ID from the URL parameters
	idParam := c.Param("id") // Adjust according to your route setup (e.g., /purchase-invoices/:id)
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Missing ID in URL"})
		return
	}

	// Bind the incoming JSON to the RequestData struct
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": err.Error()})
		return
	}

	// Begin transaction
	tx := ctrl.db.Begin()

	// Set the ID for the parent purchase invoice from the URL
	id, err := strconv.ParseInt(idParam, 10, 64) // Base 10 and 64-bit integer
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid id", "error": err.Error()})

		return
	}

	requestData.Parent.ID = id
	// Update the parent purchase invoice
	if err := ctrl.repo.UpdatePurchase(tx, &requestData.Parent); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to update parent data", "error": err.Error()})
		return
	}

	// Update the detail records
	for _, detail := range requestData.Details {
		// Check if detail already exists based on Purchase_invoice_id and other criteria if necessary
		existingDetail, err := ctrl.repo.GetDetailByID(tx, detail.ID) // Assuming detail.ID holds the ID
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to retrieve detail data", "error": err.Error()})
			return
		}

		// If the detail exists, update it
		if existingDetail != nil {
			detail.ID = existingDetail.ID // Ensure to keep the existing ID for the update
			if err := ctrl.repo.UpdateDetail(tx, &detail); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to update detail data", "error": err.Error()})
				return
			}
		} else {
			// If the detail does not exist, you may choose to insert it or handle as needed
			detail.Purchase_invoice_id = requestData.Parent.ID // Ensure the correct parent ID
			if err := ctrl.repo.CreateDetail(tx, &detail); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to store detail data", "error": err.Error()})
				return
			}
		}
	}
	for _, id := range requestData.DeletedDetail {
		if err := ctrl.repo.DeleteDetail(tx, id); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to delete detail", "error": err.Error()})
			return
		}
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to commit transaction, please try again!", "error": err.Error()})
		return
	}

	// On success, return the successful response
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Purchase invoice updated successfully", "data": requestData})
}

func (ctrl *PurchaseInvoiceController[T]) UpdateState(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}
	response, err := ctrl.repo.UpdateState(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}
	c.JSON(response.Code, response)
}
func (ctrl *PurchaseInvoiceController[T]) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}
	response, err := ctrl.repo.DeleteData(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}
	c.JSON(response.Code, response)
}
