package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"math/rand"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
	"strconv"
	"time"
)

var lettersAndDigits = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type UsuarioHandler struct {
}
type UsersResponse struct {
	Usuarios []models.Usuario `json:"usuarios"`
	Lenght   int              `json:"lenght"`
}

type UsersCodeResponse struct {
	UsuariosCode []models.UserCode `json:"usuarioscode"`
	Lenght       int               `json:"lenght"`
}

func (u *UsuarioHandler) Login(c *fiber.Ctx) error {
	var request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DeviceToken string `json:"deviceToken"`
	}

	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	usuario := models.Usuario{
		Email:    request.Email,
		Password: request.Password,
	}

	err = repositories.Login(&usuario, request.DeviceToken)
	if err != nil {
		if err.Error() == "Cuenta no verificada" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Cuenta por verificar, valide su cuenta",
				"usuario": usuario,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"jwt":     "bearer asdfkjslkdfj",
		"usuario": usuario,
	})
}

func (u *UsuarioHandler) Register(c *fiber.Ctx) error {

	var usuario models.Usuario
	err := c.BodyParser(&usuario)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.RegisterUsuario(&usuario)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar al usuario",
		})
	}
	code := generateRandomCode(10)
	sendEmail(usuario.Email, code)

	return c.JSON(fiber.Map{
		"msg":     "Usuario registrado correctamente",
		"usuario": usuario,
		"correo":  usuario.Email,
	})
}

func (u *UsuarioHandler) UpdateUser(c *fiber.Ctx) error {
	// Obtener el ID del usuario del parámetro de la ruta
	id := c.Params("id")

	// Crear un nuevo objeto Usuario para almacenar los datos de la solicitud
	var data models.Usuario
	// Parsear el cuerpo de la solicitud en el objeto Usuario
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Llamar a la función UpdateUser del repositorio
	err := repositories.UpdateUsuario(id, &data)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el usuario",
		})
	}

	// Si todo va bien, devolver el usuario actualizado
	return c.JSON(fiber.Map{
		"updatedUser": data,
	})
}

func (u *UsuarioHandler) DeleteUser(c *fiber.Ctx) error {
	// Obtener el ID del usuario del parámetro de la ruta
	id := c.Params("id")
	// Llamar a la función DeleteUser del repositorio
	err := repositories.DeleteUsuario(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar el usuario",
		})
	}

	// Si todo va bien, devolver un mensaje de éxito
	return c.JSON(fiber.Map{
		"message": "Usuario eliminado con éxito",
	})
}

func (u *UsuarioHandler) GetUsers(c *fiber.Ctx) error {
	usuarios, length, err := repositories.GetAllUsers()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&UsersResponse{
		Usuarios: usuarios,
		Lenght:   length,
	})
}

func (u *UsuarioHandler) GetUsersByRole(c *fiber.Ctx) error {
	// Obtener el rol del parámetro de la ruta
	role := c.Params("role")

	// Llamar a la función GetUserByRole del repositorio
	usuarios, length, err := repositories.GetUsuariosByRole(role)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener usuarios",
		})
	}

	// Si todo va bien, devolver los usuarios
	return c.JSON(&UsersResponse{
		Usuarios: usuarios,
		Lenght:   length,
	})
}

func (u *UsuarioHandler) RegisterUser(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	usuario := models.Usuario{
		Email: data["email"],
		// Aquí puedes agregar todos los campos que se necesiten para el registro
	}

	err := repositories.RegisterUsuario(&usuario)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Hubo un error al registrar el usuario")
	}

	// Genera un código aleatorio
	code := generateRandomCode(10)

	// Enviar el correo electrónico
	sendEmail(usuario.Email, code)

	return c.JSON(fiber.Map{"status": "success", "message": "Usuario registrado correctamente :v", "data": usuario})
}

func (u *UsuarioHandler) RegisterUserCode(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	usuario := models.Usuario{
		Email:      data["email"],
		Verificado: false, // Establecer Verificado como falso por defecto
		// Aquí puedes agregar todos los campos que se necesiten para el registro
		Nombre:   data["nombre"],
		Apellido: data["apellido"],
		Password: data["password"],
		Rol:      data["rol"],
	}

	// Verificar si el correo electrónico ya está registrado
	exists, err := repositories.EmailExists(usuario.Email)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Hubo un error al verificar el correo electrónico")
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El correo electrónico ya está registrado",
		})
	}

	// Registrar al usuario
	err = repositories.RegisterUsuario(&usuario)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Hubo un error al registrar el usuario")
	}

	// Genera un código aleatorio
	code := generateRandomCode(10)

	// Crea un nuevo registro de UserCode con el código generado y el ID del usuario recién creado
	userCode := models.UserCode{
		Code:       code,
		Status:     "active",
		ExpiryTime: time.Now().Add(5 * time.Minute),
		User:       usuario,
	}

	err = repositories.CreateUserCode(&userCode)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Hubo un error al crear el código de usuario")
	}

	// Enviar el correo electrónico
	err = sendEmail(usuario.Email, code)
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Hubo un error al enviar el correo electrónico")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Usuario registrado correctamente", "data": usuario})
}

/*Generar Codigo*/
func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = lettersAndDigits[rand.Intn(len(lettersAndDigits))]
	}
	return string(b)
}

/*Enviar Correro*/
func sendEmail(email string, code string) error {
	from := mail.NewEmail("StevenGualpa", "stevengualpa@gmail.com")
	subject := "Confirmación de registro"
	to := mail.NewEmail("User", email)
	plainTextContent := "Por favor confirme su correo electrónico y Vamos a ver Barbie"
	htmlContent := fmt.Sprintf("<p>Gracias por registrarte. Por favor confirma tu dirección de correo electrónico usando el siguiente código: <strong>%s</strong></p>", code)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return nil
	}
}

/*Obtener todos los codigos*/
func (u *UsuarioHandler) GetAllUserCodes(c *fiber.Ctx) error {
	usuarioscodes, length, err := repositories.GetAllUserCodes()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&UsersCodeResponse{
		UsuariosCode: usuarioscodes,
		Lenght:       length,
	})
}

/*Si desea crear un nuevo codigo, verificar que no este caducado o que no exita*/
func (u *UsuarioHandler) GenerateNewCode(c *fiber.Ctx) error {
	// Obtener el ID del usuario del parámetro de la ruta
	id := c.Params("id")
	// Convertir el ID del usuario a uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de usuario inválido",
		})
	}

	// Obtener el usuario para asegurarse de que exista
	usuario, err := repositories.GetUsuarioByID(uint(userID))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Llamar a la función GetUserCodeByUserID del repositorio
	userCode, err := repositories.GetUserCodeByUserID(uint(userID))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener el código del usuario",
		})
	}

	// Comprobar si el código del usuario ha expirado
	if userCode.ExpiryTime.Before(time.Now()) {
		// Si el código ha expirado, eliminar el código anterior
		err = repositories.DeleteUserCode(userCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al eliminar el código del usuario",
			})
		}

		// Generar un nuevo código aleatorio
		newCode := generateRandomCode(10)

		// Crear un nuevo UserCode con el nuevo código
		newUserCode := models.UserCode{
			UserID:     uint(userID),
			Code:       newCode,
			Status:     "active",
			ExpiryTime: time.Now().Add(5 * time.Minute),
		}

		// Guardar el nuevo UserCode en la base de datos
		err = repositories.CreateUserCode(&newUserCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al crear el nuevo código del usuario",
			})
		}

		// Enviar el correo electrónico
		err = sendEmail(usuario.Email, newCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al enviar el correo electrónico",
			})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Nuevo código generado y enviado por correo electrónico correctamente"})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ya existe un código activo para este usuario",
		})
	}
}

func (u *UsuarioHandler) GetAllDeviceTokens(c *fiber.Ctx) error {
	deviceTokens, err := repositories.GetAllDeviceTokens()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener tokens de dispositivo",
		})
	}

	return c.JSON(deviceTokens)
}

/*Verificar Codigo*/
func (u *UsuarioHandler) VerifyUserCode(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de usuario inválido",
		})
	}

	code := c.Params("code")

	// Obtener el usuario para asegurarse de que exista
	usuario, err := repositories.GetUsuarioByID(uint(userID))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	_, _, err = repositories.VerifyAndExpireUserCode1(uint(userID), code)
	if err != nil {
		if err.Error() == "El código ha expirado" {
			// Generar un nuevo código aleatorio
			newCode := generateRandomCode(10)

			// Crear un nuevo UserCode con el nuevo código
			newUserCode := models.UserCode{
				UserID:     uint(userID),
				Code:       newCode,
				Status:     "active",
				ExpiryTime: time.Now().Add(5 * time.Minute),
			}

			// Guardar el nuevo UserCode en la base de datos
			if err := repositories.CreateUserCode(&newUserCode); err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error al crear el nuevo código del usuario",
				})
			}

			// Enviar el correo electrónico
			if err := sendEmail(usuario.Email, newCode); err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error al enviar el correo electrónico",
				})
			}

			return c.JSON(fiber.Map{
				"message": "Código caducado. Se ha generado un nuevo código.",
				"usuaio":  usuario,
			})
		} else if err.Error() == "record not found" {
			return c.JSON(fiber.Map{
				"message": "Verificación fallida",
				"error":   "El código de verificación no es correcto",
			})
		} else {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Éxito en la verificación
	return c.JSON(fiber.Map{
		"message": "Cuenta verificada",
		"usuario": usuario,
	})
}

// Resetear Codigo
func (u *UsuarioHandler) ResetPassword(c *fiber.Ctx) error {
	// Obtener el correo electrónico del usuario del parámetro de la ruta
	email := c.Params("email")

	// Obtener el usuario para asegurarse de que exista usando el correo electrónico
	usuario, err := repositories.GetUsuarioByEmail(email)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Llamar a la función GetUserCodeByUserID del repositorio
	userCode, err := repositories.GetUserCodeByUserID(usuario.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener el código del usuario",
		})
	}

	// Comprobar si el código del usuario ha expirado
	if userCode.ExpiryTime.Before(time.Now()) {
		// Si el código ha expirado, eliminar el código anterior
		err = repositories.DeleteUserCode(userCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al eliminar el código del usuario",
			})
		}

		// Generar un nuevo código aleatorio
		newCode := generateRandomCode(10)

		// Crear un nuevo UserCode con el nuevo código
		newUserCode := models.UserCode{
			UserID:     usuario.ID,
			Code:       newCode,
			Status:     "active",
			ExpiryTime: time.Now().Add(5 * time.Minute),
		}

		// Guardar el nuevo UserCode en la base de datos
		err = repositories.CreateUserCode(&newUserCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al crear el nuevo código del usuario",
			})
		}

		// Enviar el correo electrónico
		err = sendEmail(usuario.Email, newCode)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al enviar el correo electrónico",
			})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Nuevo código generado y enviado por correo electrónico correctamente"})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ya existe un código activo para este usuario",
		})
	}
}

func (u *UsuarioHandler) VerifyUserCodeCorreo(c *fiber.Ctx) error {
	// Definir una estructura para recoger correo electrónico y código del cuerpo
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	// Parsear el cuerpo de la petición
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al analizar el cuerpo de la petición",
		})
	}

	// Obtener usuario por correo electrónico
	usuario, err := repositories.GetUsuarioByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Verificar el código
	_, _, err = repositories.VerifyAndExpireUserCode1(usuario.ID, body.Code)
	if err != nil {
		// Si el código es incorrecto o ha expirado, generamos un nuevo código y lo enviamos
		newCode := generateRandomCode(10)
		err = sendEmail(body.Email, newCode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al enviar el correo electrónico",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código es incorrecto o ha expirado. Se ha enviado un nuevo código al correo electrónico.",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Código verificado con éxito",
	})
}

func (u *UsuarioHandler) ChangePassword(c *fiber.Ctx) error {
	// Definir una estructura para recoger correo electrónico y la nueva contraseña del cuerpo
	type Request struct {
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}

	// Parsear el cuerpo de la petición
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al analizar el cuerpo de la petición",
		})
	}

	// Obtener usuario por correo electrónico
	usuario, err := repositories.GetUsuarioByEmail(body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Usuario no encontrado",
		})
	}

	// Actualizar la contraseña del usuario
	usuario.Password = body.NewPassword
	err = repositories.UpdateUsuario(strconv.Itoa(int(usuario.ID)), usuario)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar la contraseña del usuario",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Contraseña actualizada con éxito",
	})
}
