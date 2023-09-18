package models

import "gorm.io/gorm"

type Carrera struct {
	gorm.Model
	Nombre      string `gorm:"not null" json:"nombre"`
	ImageURL    string `gorm:"not null" json:"imageURL"`
	UrlSitio    string `gorm:"not null" json:"UrlSitio"`
	Descripcion string `gorm:"not null" json:"descripcion"`
	FacultadID  uint   `gorm:"not null" json:"facultadID"`
}

func (Carrera) TableName() string {
	return "carrera"
}
