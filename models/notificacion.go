package models

import "gorm.io/gorm"

type Notificacion struct {
	gorm.Model
	Titulo    string `json:"titulo"`
	Contenido string `json:"contenido"`
	UsuarioID uint   `json:"usuario_id"`
	Status    string `gorm:"default:'active'" json:"status"` // active, disable
}

func (Notificacion) TableName() string {
	return "notificaciones"
}
