package repositories

import (
	"errors"
	"noticias_uteq/models"
)

func RegisterCarrera(carrera *models.Carrera) error {
	// Verificar que no exista una carrera con el mismo nombre
	var existingCarrera models.Carrera
	DB.First(&existingCarrera, "nombre = ?", carrera.Nombre)
	if existingCarrera.ID != 0 {
		return errors.New("Ya existe una carrera con el mismo nombre")
	}

	// Registrar la nueva carrera
	return DB.Save(&carrera).Error
}

func UpdateCarrera(id string, carrera *models.Carrera) error {
	// Crear un objeto Carrera vacío
	var carreraToUpdate models.Carrera

	// Buscar la carrera en la base de datos
	result := DB.First(&carreraToUpdate, "id = ?", id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Actualizar los campos necesarios en el objeto carreraToUpdate con los datos del objeto carrera
	carreraToUpdate.Nombre = carrera.Nombre
	carreraToUpdate.ImageURL = carrera.ImageURL
	carreraToUpdate.UrlSitio = carrera.UrlSitio
	carreraToUpdate.Descripcion = carrera.Descripcion
	carreraToUpdate.FacultadID = carrera.FacultadID

	// Guardar el objeto carrera actualizado
	result = DB.Save(&carreraToUpdate)

	// Comprobar si se produjo algún error durante la actualización
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteCarrera(id string) error {
	// Crear un objeto Carrera vacío
	carrera := models.Carrera{}

	// Buscar la carrera en la base de datos utilizando su ID
	result := DB.First(&carrera, id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Eliminar la carrera
	result = DB.Delete(&carrera)

	// Comprobar si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllCarreras() ([]models.Carrera, int, error) {
	var carreras []models.Carrera
	result := DB.Find(&carreras)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return carreras, len(carreras), nil
}

func GetCarrerasByFacultadID(facultadID string) ([]models.Carrera, int, error) {
	var carreras []models.Carrera
	result := DB.Where("facultad_id = ?", facultadID).Find(&carreras)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return carreras, len(carreras), nil
}
