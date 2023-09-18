package repositories

import "noticias_uteq/models"

func Revistas() ([]models.Revista, error) {
	var revistas []models.Revista
	err := DB.Find(&revistas).Error
	return revistas, err
}

func RegistrarRevista(revista *models.Revista) error {
	return nil
}

func DeleteRevista(id int) error {
	return nil
}
