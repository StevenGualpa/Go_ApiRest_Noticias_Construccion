package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
	"strconv"
)

type NoticiasHandler struct {
}

type NoticiasResponse struct {
	Noticias []models.Noticia `json:"noticias"`
	Lenght   int              `json:"lenght"`
}

func (n *NoticiasHandler) GetNoticias(c *fiber.Ctx) error {
	noticias, length, err := repositories.Noticias()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&NoticiasResponse{
		Noticias: noticias,
		Lenght:   length,
	})
}

// este es el meotodo
func (n *NoticiasHandler) GetNoticiasByUsuario(c *fiber.Ctx) error {
	usuarioID, err := c.ParamsInt("usuarioID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al obtener el ID del usuario",
		})
	}

	noticias, err := repositories.GetNoticiasByUsuarioID(uint(usuarioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las noticias",
		})
	}

	return c.JSON(&NoticiasResponse{
		Noticias: noticias,
		Lenght:   len(noticias),
	})
}

func (n *NoticiasHandler) RegisterNoticias(c *fiber.Ctx) error {

	var noticia models.Noticia

	err := c.BodyParser(&noticia)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	// registramos la noticia
	err = repositories.RegisterNoticia(&noticia)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar la noticia",
		})
	}

	return c.Status(fiber.StatusOK).JSON(noticia)
}

func (n *NoticiasHandler) DeleteNoticias(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al obtener el id de la noticia",
		})
	}

	err = repositories.DeleteNoticia(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al eliminar el registro",
			"log":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (n *NoticiasHandler) Create(c *fiber.Ctx) error {
	var noticia models.Noticia

	// Deserializar el cuerpo de la petición HTTP en el struct Noticia
	if err := c.BodyParser(&noticia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se puede parsear el cuerpo de la solicitud",
		})
	}

	// Agregar la noticia a la base de datos
	if err := repositories.CreateNoticia(&noticia); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Problema al crear la noticia",
		})
	}

	// Enviar notificaciones a los usuarios
	if err := repositories.SendNotification(&noticia); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Problema al enviar las notificaciones",
		})
	}

	// Respuesta exitosa
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"noticia": noticia,
	})
}

func (n *NoticiasHandler) CreateNoticia(c *fiber.Ctx) error {
	// Crear el struct Noticia a partir del cuerpo de la solicitud
	var noticia models.Noticia
	if err := c.BodyParser(&noticia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el cuerpo de la solicitud",
		})
	}

	// Crear la noticia en la base de datos
	if err := repositories.CreateNoticia(&noticia); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear la noticia",
		})
	}

	// Enviar la notificación a todos los dispositivos registrados
	title := "Nueva Noticia"
	body := "Se ha publicado una nueva noticia: " + noticia.Titulo
	if err := repositories.EnvioNotificacion(title, body); err != nil {
		fmt.Println("Error al enviar la notificación:", err)
		// Puedes manejar este error como prefieras, dependiendo de si quieres que la falta de notificación sea crítica o no
	}

	// Devolver la noticia creada
	return c.JSON(noticia)
}

func (h NoticiasHandler) Delete(c *fiber.Ctx) error {
	token := c.Params("ID")

	err := repositories.DeleteDeviceToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (n *NoticiasHandler) CreateNoticiaFiltro(c *fiber.Ctx) error {
	var noticia models.Noticia
	if err := c.BodyParser(&noticia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el cuerpo de la solicitud",
		})
	}

	if err := repositories.CreateNoticiaFiltro(&noticia); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear la noticia",
		})
	}

	if err := repositories.SendNotificationFiltro(&noticia); err != nil {
		fmt.Println("Error al enviar la notificación:", err)
		// Puedes manejar este error como prefieras, dependiendo de si quieres que la falta de notificación sea crítica o no
	}

	return c.JSON(noticia)
}

type NotificacionesResponse struct {
	Notificaciones []models.Notificacion `json:"notificaciones"`
	Lenght         int                   `json:"lenght"`
}

func (n *NoticiasHandler) GetNotificacionesByUsuario(c *fiber.Ctx) error {
	usuarioID, err := c.ParamsInt("usuarioID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al obtener el ID del usuario",
		})
	}

	notificaciones, err := repositories.GetNotificacionesByUsuario(uint(usuarioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las notificaciones",
		})
	}

	return c.JSON(&NotificacionesResponse{
		Notificaciones: notificaciones,
		Lenght:         len(notificaciones),
	})
}

// Para actualizar la notificacion a disable
func (n *NoticiasHandler) UpdateNotificacionStatus(c *fiber.Ctx) error {
	notificacionID, err := c.ParamsInt("notificacionID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al obtener el ID de la notificación",
		})
	}

	err = repositories.UpdateNotificacionStatusToDisable(uint(notificacionID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el estado de la notificación",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Estado de notificación actualizado con éxito",
	})
}

// Parte de convocatoria:
func (n *NoticiasHandler) CreateConvocatoriaFiltro(c *fiber.Ctx) error {
	var convocatoria models.Convocatoria
	if err := c.BodyParser(&convocatoria); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el cuerpo de la solicitud",
		})
	}

	if err := repositories.CreateConvocatoria(&convocatoria); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo crear la convocatoria",
		})
	}

	if err := repositories.SendNotificationFiltroconvocatgoria(&convocatoria); err != nil {
		fmt.Println("Error al enviar la notificación:", err)
		// Puedes manejar este error como prefieras, dependiendo de si quieres que la falta de notificación sea crítica o no
	}

	return c.JSON(convocatoria)
}

func (n *NoticiasHandler) GetAllConvocatorias(ctx *fiber.Ctx) error {
	convocatorias, err := repositories.GetAllConvocatorias()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo obtener las convocatorias",
		})
	}

	return ctx.JSON(fiber.Map{
		"convocatorias": convocatorias,
	})
}

func (n *NoticiasHandler) GetConvocatoriaByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	convocatoria, err := repositories.GetConvocatoriaByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo obtener la convocatoria",
		})
	}

	return ctx.JSON(convocatoria)
}

func (n *NoticiasHandler) UpdateConvocatoria(ctx *fiber.Ctx) error {
	var convocatoriaUpdates models.Convocatoria
	if err := ctx.BodyParser(&convocatoriaUpdates); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear el cuerpo de la solicitud",
		})
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	if err := repositories.UpdateConvocatoria(uint(id), convocatoriaUpdates); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo actualizar la convocatoria",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Convocatoria actualizada exitosamente",
	})
}

func (n *NoticiasHandler) GetConvocatoriasByUsuario(c *fiber.Ctx) error {
	usuarioIDStr := c.Params("id")
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Error al obtener el ID del usuario",
			"detail": err.Error(),
		})
	}

	convocatorias, err := repositories.GetConvocatoriasByUsuarioID(uint(usuarioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Error al obtener las convocatorias",
			"detail": err.Error(),
		})
	}

	return c.JSON(&ConvocatoriasResponse{
		Convocatorias: convocatorias,
		Lenght:        len(convocatorias),
	})
}

type ConvocatoriasResponse struct {
	Convocatorias []models.Convocatoria `json:"convocatorias"`
	Lenght        int                   `json:"length"`
}
