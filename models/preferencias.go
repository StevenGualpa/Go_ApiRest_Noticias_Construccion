package models

import "gorm.io/gorm"

type UsuarioFacultad struct {
	gorm.Model
	UsuarioID  uint
	Usuario    Usuario
	FacultadID uint
	Facultad   Facultad
}

func (UsuarioFacultad) TableName() string {
	return "usuarioFacultad"
}

type UsuarioFacultadNombre struct {
	gorm.Model
	UsuarioID      uint
	Usuario        Usuario
	FacultadNombre string
	Facultad       Facultad `gorm:"foreignKey:FacultadNombre;references:Nombre"`
}

func (UsuarioFacultadNombre) TableName() string {
	return "usuarioFacultadNombre"
}
