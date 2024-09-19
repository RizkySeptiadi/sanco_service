package repository

import (
	"sanco_microservices/database"
	"sanco_microservices/structs"

	"gorm.io/gorm"
)

type SuppliersRepository struct {
	db *gorm.DB
}

func NewSuppliersRepository() *SuppliersRepository {
	return &SuppliersRepository{
		db: database.DbConnection,
	}
}

func (repo *SuppliersRepository) CreateSuppliers(roomType *structs.Sanco_Suppliers) error {
	return repo.db.Create(roomType).Error
}

func (repo *SuppliersRepository) GetSuppliersByID(id int64) (*structs.Sanco_Suppliers, error) {
	var roomType structs.Sanco_Suppliers
	if err := repo.db.First(&roomType, id).Error; err != nil {
		return nil, err
	}
	return &roomType, nil
}

func (repo *SuppliersRepository) UpdateSuppliers(id int64, updatedRoomType *structs.Sanco_Suppliers) error {
	var roomType structs.Sanco_Suppliers

	// Check if the room type exists
	if err := repo.db.First(&roomType, id).Error; err != nil {
		return err
	}

	// Update the existing room type with the new data
	return repo.db.Model(&roomType).Updates(updatedRoomType).Error
}

func (repo *SuppliersRepository) DeleteSuppliers(id int64) error {
	return repo.db.Delete(&structs.Sanco_Suppliers{}, id).Error
}

func (repo *SuppliersRepository) GetAllSupplierss() ([]structs.Sanco_Suppliers, error) {
	var roomTypes []structs.Sanco_Suppliers
	if err := repo.db.Find(&roomTypes).Error; err != nil {
		return nil, err
	}
	return roomTypes, nil
}
