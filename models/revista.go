package models

import "gorm.io/gorm"

type Revista struct {
	gorm.Model
	Titulo  string `json:"titulo"`
	Date    string `json:"date"`
	Portada string `json:"portada"`
	Url     string `json:"url"`
}

func (Revista) TableName() string {
	return "revistas"
}
