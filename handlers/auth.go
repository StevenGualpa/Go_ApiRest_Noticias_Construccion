package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "https://noticias-uteq-4c62c24e7cc5.herokuapp.com/usuarios/google/callback",
	ClientID:     "371719052873-09o0nk4p9gf6j6229lb1v4ttg1ree5fc.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GoogleLoginHandler inicia el flujo OAuth2 con Google
func GoogleLoginHandler(c *fiber.Ctx) error {
	// Define una estructura para tu cuerpo de solicitud
	type RequestBody struct {
		Codigo string `json:"codigo"`
	}

	// Instancia tu estructura
	var body RequestBody

	// Parsea el cuerpo de la solicitud entrante y almacena el valor en tu estructura
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al parsear el cuerpo de la solicitud",
		})
	}

	// Obtiene el código de tu estructura
	codigo := body.Codigo

	// Si no hay un código, devuelve un error (puedes manejar esto como prefieras)
	if codigo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Código requerido",
		})
	}

	// Inicia el flujo OAuth con el código como el state
	url := googleOauthConfig.AuthCodeURL(codigo, oauth2.AccessTypeOffline)
	return c.JSON(fiber.Map{
		"url": url,
	})
}

// GoogleCallbackHandler maneja el callback de Google y obtiene el token
func GoogleCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	estado := c.Query("state") // Este será tu código enviado desde el cliente

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("No se pudo obtener el token de Google: " + err.Error())
	}

	// Obtén la información del perfil del usuario
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("No se pudo obtener la información del usuario: " + err.Error())
	}
	defer resp.Body.Close()

	// Decodifica la respuesta JSON
	var userinfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("No se pudo decodificar la información del usuario: " + err.Error())
	}

	// Aquí guardamos la información del usuario y el código en la tabla de sesión
	sesion := models.Sesion{
		Email:         userinfo["email"].(string),
		EmailVerified: userinfo["email_verified"].(bool),
		FamilyName:    userinfo["family_name"].(string),
		GivenName:     userinfo["given_name"].(string),
		Locale:        userinfo["locale"].(string),
		Name:          userinfo["name"].(string),
		Picture:       userinfo["picture"].(string),
		Sub:           userinfo["sub"].(string),
		Codigo:        estado, // Aquí usamos el código que obtuvimos desde el cliente
	}
	err = repositories.CreateSesion(&sesion)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creando sesión: " + err.Error())
	}

	return c.SendString("Sesión creada con éxito")
}

type AuthHandler struct {
}

type SesionsResponse struct {
	Sesiones []models.Sesion `json:"sesiones"`
	Lenght   int             `json:"lenght"`
}

func (u *AuthHandler) GetUsersAuth(c *fiber.Ctx) error {
	sesiones, length, err := repositories.GetAllSesiones()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&SesionsResponse{
		Sesiones: sesiones,
		Lenght:   length,
	})
}

func (u *AuthHandler) GoogleSignInHandler(c *fiber.Ctx) error {
	// Define una estructura para tu cuerpo de solicitud
	type RequestBody struct {
		Codigo      string `json:"codigo"`
		DeviceToken string `json:"deviceToken"`
	}

	// Instancia tu estructura
	var body RequestBody

	// Parsea el cuerpo de la solicitud entrante y almacena el valor en tu estructura
	if err := c.BodyParser(&body); err != nil {
		fmt.Println("Error al parsear el cuerpo de la solicitud:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error al parsear el cuerpo de la solicitud")
	}

	// Obtiene el código y el token del dispositivo de tu estructura
	codigo := body.Codigo
	deviceToken := body.DeviceToken

	// Si no hay un código o un token de dispositivo, devuelve un error
	if codigo == "" || deviceToken == "" {
		fmt.Println("Código y token del dispositivo requeridos")
		return c.Status(fiber.StatusBadRequest).SendString("Código y token del dispositivo requeridos")
	}

	// Intenta iniciar sesión usando el método del repositorio
	usuario, err := repositories.Logingoogle(codigo, deviceToken)
	if err != nil {
		// Si hay un error, devuelve un mensaje de error apropiado
		fmt.Println("Error durante el inicio de sesión:", err)
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	// Si todo sale bien, devuelve un mensaje de éxito junto con los datos del usuario
	return c.JSON(fiber.Map{
		"message": "Inicio de sesión exitoso",
		"usuario": usuario,
	})
}
