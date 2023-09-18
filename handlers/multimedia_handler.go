package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
)

type MultimediaHandler struct {
}

type MultimediaResponse struct {
	Multimedias []models.Multimedia `json:"multimedias"`
	Lenght      int                 `json:"lenght"`
}

func (h *MultimediaHandler) CreateMultimedia1(c *fiber.Ctx) error {
	var multimedia models.Multimedia
	err := c.BodyParser(&multimedia)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.RegisterMultimedia(&multimedia)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar el multimedia",
		})
	}

	// Enviar la notificación a todos los dispositivos registrados
	title := "Nuevo Contenido Multimedia"
	body := "Se ha publicado un nuevo contenido multimedia: " + multimedia.Titulo
	if err := repositories.EnvioNotificacion(title, body); err != nil {
		fmt.Println("Error al enviar la notificación:", err)
		// Puedes manejar este error como prefieras, dependiendo de si quieres que la falta de notificación sea crítica o no
	}

	return c.JSON(fiber.Map{
		"msg":        "Multimedia registrado correctamente",
		"multimedia": multimedia,
	})
}

func (h *MultimediaHandler) UpdateMultimedia(c *fiber.Ctx) error {
	id := c.Params("id")
	var multimedia models.Multimedia
	if err := c.BodyParser(&multimedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se puede analizar el JSON",
		})
	}

	err := repositories.UpdateMultimedia(id, &multimedia)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el multimedia",
		})
	}

	return c.JSON(multimedia)
}

func (h *MultimediaHandler) DeleteMultimedia(c *fiber.Ctx) error {
	id := c.Params("id")
	err := repositories.DeleteMultimedia(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar el multimedia",
		})
	}

	return c.JSON(fiber.Map{
		"success": "Multimedia eliminado correctamente",
	})
}

func (h *MultimediaHandler) GetAllMultimedia(c *fiber.Ctx) error {
	multimedias, length, err := repositories.GetAllMultimedia()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&MultimediaResponse{
		Multimedias: multimedias,
		Lenght:      length,
	})
}

func (h *MultimediaHandler) GetMultimediaByContentType(c *fiber.Ctx) error {
	contentType := c.Params("contentType")

	multimedias, length, err := repositories.GetMultimediaByContentType(contentType)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las multimedia",
		})
	}

	return c.JSON(&MultimediaResponse{
		Multimedias: multimedias,
		Lenght:      length,
	})
}

func (h *MultimediaHandler) CreateMultimedia(c *fiber.Ctx) error {
	var multimedia models.Multimedia
	err := c.BodyParser(&multimedia)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.RegisterMultimedia(&multimedia)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar el multimedia",
		})
	}

	// Enviar la notificación a todos los dispositivos registrados
	title := "Nuevo Contenido Multimedia"
	body := "Se ha publicado un nuevo contenido multimedia: " + multimedia.Titulo
	if err := repositories.EnvioNotificacion(title, body); err != nil {
		fmt.Println("Error al enviar la notificación:", err)
		// Puedes manejar este error como prefieras, dependiendo de si quieres que la falta de notificación sea crítica o no
	}

	// Guardar la notificación en la base de datos
	// Asumiendo que tienes una función SaveNotification que toma un título, cuerpo y usuarioID
	err = repositories.SaveNotification(title, body, multimedia.UsuarioID)
	if err != nil {
		fmt.Println("Error al guardar la notificación:", err)
	}

	return c.JSON(fiber.Map{
		"msg":        "Multimedia registrado correctamente",
		"multimedia": multimedia,
	})
}
