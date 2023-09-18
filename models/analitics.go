package models

import (
	"gorm.io/gorm"
	"time"
)

type Analitics struct {
	gorm.Model
	Seccion string    `json:"seccion"`
	Nombre  string    `json:"nombre"`
	Fecha   time.Time `json:"fecha"`
}

func (Analitics) TableName() string {
	return "analitics"
}
