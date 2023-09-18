package repositories

import (
	"noticias_uteq/models"
	"time"
)

// Crear un nuevo evento
func CreateEvento(evento *models.Evento) error {
	return DB.Create(&evento).Error
}

// Obtener todos los eventos
func GetAllEventos() ([]models.Evento, int, error) {
	var eventos []models.Evento
	result := DB.Find(&eventos)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return eventos, len(eventos), nil
}

// Obtener un evento por ID
func GetEventoByID(id uint) (models.Evento, error) {
	var evento models.Evento
	result := DB.First(&evento, id)
	if result.Error != nil {
		return models.Evento{}, result.Error
	}
	return evento, nil
}

// Actualizar un evento
func UpdateEvento(id uint, updatedEvento *models.Evento) error {
	evento, err := GetEventoByID(id)
	if err != nil {
		return err
	}

	evento.Name = updatedEvento.Name
	evento.Time = updatedEvento.Time
	evento.Date = updatedEvento.Date
	evento.IDUsuario = updatedEvento.IDUsuario

	return DB.Save(&evento).Error
}

// Eliminar un evento
func DeleteEvento(id uint) error {
	evento, err := GetEventoByID(id)
	if err != nil {
		return err
	}

	return DB.Delete(&evento).Error
}

// Obtener todos los eventos por fecha
func GetEventosByFecha(fecha time.Time) ([]models.Evento, int, error) {
	var eventos []models.Evento
	result := DB.Where("date = ?", fecha).Find(&eventos)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return eventos, len(eventos), nil
}

// Obtener todos los eventos por ID de usuario
func GetEventosByIDUsuario(idUsuario uint) ([]models.Evento, int, error) {
	var eventos []models.Evento
	result := DB.Where("ID_usuario = ?", idUsuario).Find(&eventos)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return eventos, len(eventos), nil
}
