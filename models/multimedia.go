package models

import "gorm.io/gorm"

type Multimedia struct {
	gorm.Model
	Titulo        string  `json:"titulo"`
	Descripcion   string  `json:"descripcion"`
	URLImagen     string  `json:"url_imageb"`
	URLVideo      string  `json:"url_video"`
	TipoContenido string  `json:"tipo_contenido"`
	UsuarioID     uint    `json:"usuario_id"`
	Usuario       Usuario `gorm:"foreignKey:UsuarioID" json:"usuario"`
}

func (Multimedia) TableName() string {
	return "multimedia"
}
