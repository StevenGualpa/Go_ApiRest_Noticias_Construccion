package repositories

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"noticias_uteq/models"
)

func EnvioNotificacion(title string, body string) error {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return errors.New("error getting Messaging client")
	}

	deviceTokens, err := GetAllDeviceTokenss()
	if err != nil {
		return err
	}

	// Usamos un mapa para almacenar los tokens únicos
	uniqueTokens := make(map[string]bool)
	for _, deviceToken := range deviceTokens {
		uniqueTokens[deviceToken.Token] = true
	}

	// Enviamos la notificación solo a los tokens únicos
	for token := range uniqueTokens {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: title, // Título personalizado
				Body:  body,  // Cuerpo personalizado
			},
			Token: token,
		}

		response, err := client.Send(ctx, message)
		if err != nil {
			fmt.Printf("error sending message: %v\n", err)
			continue
		}

		fmt.Printf("Successfully sent message: %v\n", response)
	}

	return nil
}

func ObtenerTokens() ([]models.DeviceToken, error) {
	var deviceTokens []models.DeviceToken
	result := DB.Find(&deviceTokens)
	if result.Error != nil {
		return nil, errors.New("no se pudo obtener los tokens de los dispositivos")
	}
	return deviceTokens, nil
}
