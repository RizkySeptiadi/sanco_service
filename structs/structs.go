package structs

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID                int64
	Username          string `gorm:"unique"`
	Name              string
	Email             string `gorm:"unique"`
	Email_verified_at time.Time
	Password          string
	Remember_token    string
	Created_at        time.Time
	Updated_at        time.Time
}
type Sanco_Suppliers struct {
	gorm.Model
	ID                  int64
	Code                string `gorm:"unique"  binding:"required"`
	Name                string `binding:"required"`
	Bank_account_number string
	Bank_account_name   string
	Contact             string
	Address             string
	Phone_num           string
	Email               string
	Status              int32
}

type Role struct {
	gorm.Model

	ID        int64
	Name      string `gorm:"unique"`
	Users     []User `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
