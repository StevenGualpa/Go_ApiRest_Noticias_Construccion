package repositories

import (
	"noticias_uteq/models"
	"time"
)

// Insertar una nueva entrada de Analitics
func InsertAnalitics(analitics *models.Analitics) error {
	analitics.Fecha = time.Now() // Establecer la fecha actual
	return DB.Save(&analitics).Error
}

// Obtener todos los registros de Analitics
func GetAllAnalitics() ([]models.Analitics, int, error) {
	var analitics []models.Analitics
	result := DB.Find(&analitics)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return analitics, len(analitics), nil
}

// Obtener todos los registros de Analitics por sección
func GetAnaliticsBySeccion(seccion string) ([]models.Analitics, int, error) {
	var analitics []models.Analitics
	result := DB.Where("seccion = ?", seccion).Find(&analitics)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return analitics, len(analitics), nil
}

// Obtener todos los registros de Analitics por fecha
func GetAnaliticsByFecha(fecha time.Time) ([]models.Analitics, int, error) {
	var analitics []models.Analitics
	result := DB.Where("fecha = ?", fecha).Find(&analitics)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return analitics, len(analitics), nil
}

// Obtener todos los registros de Analitics por sección y fecha
func GetAnaliticsBySeccionAndFecha(seccion string, fecha time.Time) ([]models.Analitics, int, error) {
	var analitics []models.Analitics
	result := DB.Where("seccion = ? AND fecha = ?", seccion, fecha).Find(&analitics)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return analitics, len(analitics), nil
}

// Este tipo representa la estructura del resultado de la consulta
type NombreCount struct {
	Nombre string `json:"nombre"`
	Count  int    `json:"count"`
}

// Obtener todos los registros de Analitics por sección y contar las repeticiones de cada nombre
func GetNombreCountBySeccion(seccion string) ([]NombreCount, error) {
	var results []NombreCount
	result := DB.Table("analitics").Select("nombre, COUNT(nombre) as count").Where("seccion = ?", seccion).Group("nombre").Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}
