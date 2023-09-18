package repositories

import (
	"errors"
	"noticias_uteq/models"
)

func Multimedia() ([]models.Multimedia, int, error) {
	var result []models.Multimedia

	err := DB.Preload("Usuario").Find(&result).Error
	if err != nil {
		return result, 0, err
	}

	return result, len(result), nil
}

func RegisterMultimedia(multimedia *models.Multimedia) error {
	var existingMultimedia models.Multimedia
	DB.First(&existingMultimedia, "titulo = ?", multimedia.Titulo)
	if existingMultimedia.ID != 0 {
		return errors.New("Ya existe un contenido multimedia con el mismo titulo")
	}
	return DB.Save(&multimedia).Error
}

func UpdateMultimedia(id string, multimedia *models.Multimedia) error {
	var multimediaToUpdate models.Multimedia

	result := DB.First(&multimediaToUpdate, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	multimediaToUpdate.Titulo = multimedia.Titulo
	multimediaToUpdate.Descripcion = multimedia.Descripcion
	multimediaToUpdate.URLImagen = multimedia.URLImagen
	multimediaToUpdate.URLVideo = multimedia.URLVideo
	multimediaToUpdate.TipoContenido = multimedia.TipoContenido
	multimediaToUpdate.UsuarioID = multimedia.UsuarioID

	result = DB.Save(&multimediaToUpdate)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteMultimedia(id string) error {
	multimedia := models.Multimedia{}
	result := DB.First(&multimedia, id)
	if result.Error != nil {
		return result.Error
	}
	result = DB.Delete(&multimedia)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllMultimedia() ([]models.Multimedia, int, error) {
	var multimedia []models.Multimedia
	result := DB.Find(&multimedia)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return multimedia, len(multimedia), nil
}

func GetMultimediaByContentType(contentType string) ([]models.Multimedia, int, error) {
	var multimedia []models.Multimedia
	result := DB.Where("tipo_contenido = ?", contentType).Find(&multimedia)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return multimedia, len(multimedia), nil
}
