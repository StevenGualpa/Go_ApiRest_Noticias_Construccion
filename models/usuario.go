package models

import "gorm.io/gorm"

type Usuario struct {
	gorm.Model
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Email      string `gorm:"unique" json:"email"` // Agregada la etiqueta 'gorm:"unique"'
	Password   string `json:"password,omitempty"`
	Rol        string `json:"rol,omitempty"`
	Verificado bool   `json:"verificado"`
}

func (Usuario) TableName() string {
	return "usuarios"
}
