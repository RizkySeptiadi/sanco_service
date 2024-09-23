package controllers

import (
	"net/http"
	"sanco_microservices/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GeneralController[T any] struct {
	repo *repository.GeneralRepository[T]
}

func NewGeneralController[T any](repo *repository.GeneralRepository[T]) *GeneralController[T] {
	return &GeneralController[T]{repo: repo}
}

func (ctrl *GeneralController[T]) Create(c *gin.Context) {
	var data T

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.CreateData(&data)
	if err != nil {
		// Include both the message and error details in the response
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	// On success, return the successful response
	c.JSON(response.Code, response)
}

func (ctrl *GeneralController[T]) GetTables(c *gin.Context) {
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

	// Call the repository to fetch paginated data
	response, err := ctrl.repo.PaginateData(start, length, search, orderBy, orderDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated response
	c.JSON(http.StatusOK, response)
}

func (ctrl *GeneralController[T]) Get(c *gin.Context) {
	response, err := ctrl.repo.GetAllData() // Assuming GetAllData method is implemented in the repository
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	c.JSON(response.Code, response)
}

func (ctrl *GeneralController[T]) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.GetDataByID(id)
	if err != nil {
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	c.JSON(response.Code, response)
}

func (ctrl *GeneralController[T]) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid ID", "error": err.Error()})
		return
	}

	var updatedData T
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid Request", "error": err.Error()})
		return
	}

	response, err := ctrl.repo.UpdateData(id, &updatedData)
	if err != nil {
		// Include both the message and error details in the response
		c.JSON(response.Code, gin.H{"code": response.Code, "message": response.Message, "error": err.Error()})
		return
	}

	// On success, just return the response
	c.JSON(response.Code, response)
}

func (ctrl *GeneralController[T]) Delete(c *gin.Context) {
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
