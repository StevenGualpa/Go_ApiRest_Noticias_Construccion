package models

import (
	"gorm.io/gorm"
	"time"
)

type UserCode struct {
	gorm.Model
	UserID     uint      `gorm:"not null" json:"userId"`
	Code       string    `gorm:"not null" json:"code"`
	Status     string    `gorm:"not null;default:'active'" json:"status"` // active, used, expired
	ExpiryTime time.Time `gorm:"not null" json:"expiryTime"`
	User       Usuario   `gorm:"foreignKey:UserID" json:"-"`
}

func (UserCode) TableName() string {
	return "usuariocode"
}
