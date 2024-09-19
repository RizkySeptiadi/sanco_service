package controllers

import (
	"net/http"
	"sanco_microservices/repository"
	"sanco_microservices/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuppliersController struct {
	Repo *repository.SuppliersRepository
}

func NewSuppliersController() *SuppliersController {
	return &SuppliersController{
		Repo: repository.NewSuppliersRepository(),
	}
}

func (controller *SuppliersController) Create(c *gin.Context) {
	var Suppliers structs.Sanco_Suppliers
	if err := c.ShouldBindJSON(&Suppliers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Repo.CreateSuppliers(&Suppliers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room type"})
		return
	}

	c.JSON(http.StatusOK, Suppliers)
}

func (controller *SuppliersController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room type ID"})
		return
	}

	Suppliers, err := controller.Repo.GetSuppliersByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room type not found"})
		return
	}

	c.JSON(http.StatusOK, Suppliers)
}

func (controller *SuppliersController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room type ID"})
		return
	}

	var updatedSuppliers structs.Sanco_Suppliers
	if err := c.ShouldBindJSON(&updatedSuppliers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Repo.UpdateSuppliers(id, &updatedSuppliers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room type updated successfully"})
}

func (controller *SuppliersController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room type ID"})
		return
	}

	if err := controller.Repo.DeleteSuppliers(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room type deleted successfully"})
}

func (controller *SuppliersController) Get(c *gin.Context) {
	Supplierss, err := controller.Repo.GetAllSupplierss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room types"})
		return
	}

	c.JSON(http.StatusOK, Supplierss)
}
