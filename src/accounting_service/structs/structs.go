package structs

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
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

type Role struct {
	gorm.Model

	ID        int64
	Name      string  `gorm:"unique"`
	Users     []Users `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
