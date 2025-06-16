package structs

import "gorm.io/gorm"

type Companies struct {
	gorm.Model
	ID       int64
	RefID    int64  `json:"ref_id" binding:"required"`
	SourceID string `json:"source_id" binding:"required"`
	Name     string `json:"company_name"`
	Note     string `json:"note"`
	IsActive int32  `json:"is_active"`
}

type Suppliers struct {
	gorm.Model
	ID          int64
	RefID       int64  `json:"ref_id" binding:"required"`
	Name        string `json:"supplier_name" binding:"required"`
	SourceID    string `json:"source_id" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Note        string `json:"note"`
	IsActive    int32  `json:"is_active"`
}

type Customers struct {
	gorm.Model
	ID                      int64
	RefID                   int64  `json:"ref_id" binding:"required"`
	Name                    string `json:"customer_name" binding:"required"`
	SourceID                string `json:"source_id" binding:"required"`
	PhoneNumber             string `json:"phone_number"`
	Email                   string `json:"email"`
	Address                 string `json:"address"`
	TaxIdentificationNumber string `json:"tax_identification_number"`
	Note                    string `json:"note"`
	IsActive                int32  `json:"is_active"`
}
