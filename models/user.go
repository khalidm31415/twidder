package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint
	Username    string    `gorm:"type:varchar(100);not null;unique"`
	DisplayName string    `gorm:"type:varchar(100);not null"`
	Email       string    `gorm:"type:varchar(254);not null;unique"`
	Password    string    `gorm:"type:text;not null" json:"-"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	DeletedAt   gorm.DeletedAt
}
