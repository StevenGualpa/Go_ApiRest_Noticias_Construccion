package repositories

import (
	"noticias_uteq/models"
)

func RegisterFacultad(facultad *models.Facultad) error {
	return DB.Save(&facultad).Error
}

/*Metodo para Actualizar una Facultad*/
func UpdateFacultad(id string, facultad *models.Facultad) error {
	// Crea un objeto Facultad vacío
	var facultadToUpdate models.Facultad

	// Busca a la facultad en la base de datos
	result := DB.First(&facultadToUpdate, "id = ?", id)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Actualiza los campos necesarios en el objeto facultadToUpdate con los datos del objeto facultad
	facultadToUpdate.Nombre = facultad.Nombre
	facultadToUpdate.UrlVideo = facultad.UrlVideo
	facultadToUpdate.Mision = facultad.Mision
	facultadToUpdate.Vision = facultad.Vision
	facultadToUpdate.Correo = facultad.Correo
	facultadToUpdate.Longitud = facultad.Longitud
	facultadToUpdate.Latitud = facultad.Latitud

	// Guarda el objeto facultad actualizado
	result = DB.Save(&facultadToUpdate)

	// Comprueba si se produjo algún error durante la actualización
	if result.Error != nil {
		return result.Error
	}

	return nil
}

/*Metodo para Eliminar un Usuario*/
func DeleteFacultad(id string) error {
	// Crea un objeto de usuario vacío
	facultad := models.Facultad{}

	// Busca al usuario en la base de datos utilizando su ID
	result := DB.First(&facultad, id)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Elimina al usuario
	result = DB.Delete(&facultad)

	// Comprueba si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllFacultades() ([]models.Facultad, int, error) {
	var facultades []models.Facultad
	result := DB.Find(&facultades)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return facultades, len(facultades), nil
}
