package models

import "gorm.io/gorm"

type DeviceToken struct {
	gorm.Model
	Token     string `json:"token"`
	UsuarioID uint   `json:"usuario_id"`
}

func (DeviceToken) TableName() string {
	return "devicestoken"
}
