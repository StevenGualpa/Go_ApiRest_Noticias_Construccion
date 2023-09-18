package repositories

import (
	"errors"
	"noticias_uteq/models"
)

func CreateUsuarioFacultad(usuarioFacultad *models.UsuarioFacultad) error {
	// Comprobar si el usuario ya tiene 3 facultades favoritas
	var count int64
	DB.Model(&models.UsuarioFacultad{}).Where("usuario_id = ?", usuarioFacultad.UsuarioID).Count(&count)
	if count >= 3 {
		return errors.New("El usuario ya tiene 3 facultades favoritas")
	}

	// Comprobar si el usuario ya tiene esta facultad como favorita
	var existingUsuarioFacultad models.UsuarioFacultad
	DB.First(&existingUsuarioFacultad, "usuario_id = ? AND facultad_id = ?", usuarioFacultad.UsuarioID, usuarioFacultad.FacultadID)
	if existingUsuarioFacultad.ID != 0 {
		return errors.New("El usuario ya tiene esta facultad como favorita")
	}

	// Crear la preferencia de facultad para el usuario
	return DB.Save(&usuarioFacultad).Error
}

// No se usa por ahora
func DeleteUsuarioFacultad(id string) error {
	// Crear un objeto UsuarioFacultad vacío
	usuarioFacultad := models.UsuarioFacultad{}

	// Buscar la preferencia de facultad del usuario en la base de datos utilizando su ID
	result := DB.First(&usuarioFacultad, id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Eliminar la preferencia de facultad del usuario
	result = DB.Delete(&usuarioFacultad)

	// Comprobar si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllFacultadPreferences(usuarioID string) ([]models.UsuarioFacultad, int, error) {
	var usuarioFacultades []models.UsuarioFacultad
	result := DB.Preload("Usuario").Preload("Facultad").Where("usuario_id = ?", usuarioID).Find(&usuarioFacultades)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return usuarioFacultades, len(usuarioFacultades), nil
}

func UpsertUsuarioFacultad(usuarioFacultades []*models.UsuarioFacultad) error {
	// Si el usuario no tiene preferencias registradas, registra hasta un máximo de tres
	if len(usuarioFacultades) == 0 {
		return errors.New("No hay preferencias para registrar")
	}

	// Eliminar las preferencias existentes del usuario
	err := DeleteUsuarioFacultadByUsuarioID(usuarioFacultades[0].UsuarioID)
	if err != nil {
		return err
	}

	// Registrar las nuevas preferencias
	for _, usuarioFacultad := range usuarioFacultades {
		err := DB.Save(&usuarioFacultad).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Eliminar todas las preferencias de facultad del usuario basadas en su UsuarioID
func DeleteUsuarioFacultadByUsuarioID(usuarioID uint) error {
	// Eliminar todas las preferencias de facultad del usuario basadas en su UsuarioID
	result := DB.Where("usuario_id = ?", usuarioID).Delete(&models.UsuarioFacultad{})

	// Comprobar si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Metodos que usan endpoint de Relaciones Publicas Repository
func CreateUsuarioFacultadRp(usuarioFacultad *models.UsuarioFacultadNombre) error {
	// Comprobar si el usuario ya tiene 3 facultades favoritas
	var count int64
	DB.Model(&models.UsuarioFacultad{}).Where("usuario_id = ?", usuarioFacultad.UsuarioID).Count(&count)
	if count >= 3 {
		return errors.New("El usuario ya tiene 3 facultades favoritas")
	}

	// Comprobar si el usuario ya tiene esta facultad como favorita
	var existingUsuarioFacultad models.UsuarioFacultad
	DB.First(&existingUsuarioFacultad, "usuario_id = ? AND facultad_nombre = ?", usuarioFacultad.UsuarioID, usuarioFacultad.FacultadNombre)
	if existingUsuarioFacultad.ID != 0 {
		return errors.New("El usuario ya tiene esta facultad como favorita")
	}

	// Crear la preferencia de facultad para el usuario
	return DB.Save(&usuarioFacultad).Error
}

func UpsertUsuarioFacultadRp(usuarioFacultades []*models.UsuarioFacultadNombre) error {
	// Si el usuario no tiene preferencias registradas, registra hasta un máximo de tres
	if len(usuarioFacultades) == 0 {
		return errors.New("No hay preferencias para registrar")
	}

	// Eliminar las preferencias existentes del usuario
	err := DeleteUsuarioFacultadByUsuarioIDRp(usuarioFacultades[0].UsuarioID)
	if err != nil {
		return err
	}

	// Registrar las nuevas preferencias
	for _, usuarioFacultad := range usuarioFacultades {
		err := DB.Save(&usuarioFacultad).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteUsuarioFacultadByUsuarioIDRp(usuarioID uint) error {
	// Eliminar todas las preferencias de facultad del usuario basadas en su UsuarioID
	result := DB.Where("usuario_id = ?", usuarioID).Delete(&models.UsuarioFacultadNombre{})

	// Comprobar si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllFacultadPreferencesRp(usuarioID string) ([]models.UsuarioFacultadNombre, int, error) {
	var usuarioFacultades []models.UsuarioFacultadNombre
	result := DB.Preload("Usuario").Preload("Facultad").Where("usuario_id = ?", usuarioID).Find(&usuarioFacultades)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return usuarioFacultades, len(usuarioFacultades), nil
}
