package repositories

import "noticias_uteq/models"

func Analtics() ([]models.Analitics, int, error) {
	var analitics []models.Analitics

	err := DB.Preload("Usuario").Preload("Tags").Find(&analitics).Error
	if err != nil {
		return analitics, 0, err
	}

	return analitics, len(analitics), nil
}
