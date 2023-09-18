package models

import "gorm.io/gorm"

type Facultad struct {
	gorm.Model
	Nombre      string `gorm:"not null;unique" json:"nombre"`
	UrlVideo    string `gorm:"not null" json:"urlVideo"`
	UrlFacebook string `gorm:"not null" json:"urlFacebook"`
	UrlSitio    string `gorm:"not null" json:"UrlSitio"`
	Mision      string `gorm:"not null" json:"mision"`
	Correo      string `gorm:"not null" json:"correo"`
	Latitud     string `gorm:"not null" json:"latitud"`
	Longitud    string `gorm:"not null" json:"longitud"`
	Vision      string `gorm:"not null" json:"vision"`
}

func (Facultad) TableName() string {
	return "facultad"
}
