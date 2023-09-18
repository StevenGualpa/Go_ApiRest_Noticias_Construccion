package models

import "gorm.io/gorm"

type Sesion struct {
	gorm.Model
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
	Codigo        string `json:"codigo"`
}

func (Sesion) TableName() string {
	return "sesiones"
}
