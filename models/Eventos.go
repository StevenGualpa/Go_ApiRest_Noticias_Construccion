package models

import (
	"gorm.io/gorm"
	"time"
)

type Evento struct {
	gorm.Model
	IDUsuario uint      `json:"ID_usuario" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null;size:255"`
	Time      time.Time `json:"time" gorm:"not null"`
	Date      time.Time `json:"date" gorm:"not null"`
	Usuario   Usuario   `json:"usuario" gorm:"foreignKey:IDUsuario"`
}

func (Evento) TableName() string {
	return "eventos"
}
