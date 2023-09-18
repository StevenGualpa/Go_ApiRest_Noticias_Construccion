package models

import (
	"gorm.io/gorm"
	"time"
)

type Convocatoria struct {
	gorm.Model
	Titulo       string    `json:"titulo"`
	Asunto       string    `json:"asunto"`
	FechaInicio  time.Time `json:"fecha_inicio"`
	FechaFin     time.Time `json:"fecha_fin"`
	Lugar        string    `json:"lugar"`
	Emisor       string    `json:"emisor"`
	Portada      string    `json:"portada"`
	Destinatario string    `json:"destinatario"`
}

func (Convocatoria) TableName() string {
	return "convocatorias"
}
