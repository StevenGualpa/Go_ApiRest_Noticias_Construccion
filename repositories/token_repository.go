package repositories

import (
	"gorm.io/gorm"
)

type DeviceTokenRepository struct {
	DB *gorm.DB
}
