package models

import gorm "gorm.io/gorm"

type Noticia struct {
	gorm.Model
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Portada     string `json:"portada"`
	URL         string `json:"url"`
	Categoria   string `json:"categoria"`
}

func (Noticia) TableName() string {
	return "noticias"
}
