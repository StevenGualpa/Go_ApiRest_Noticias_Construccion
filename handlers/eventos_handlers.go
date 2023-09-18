package handlers

import (
	"github.com/gofiber/fiber/v2"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
	"strconv"
	"time"
)

type EventoHandler struct{}

func (h *EventoHandler) CreateEvento(c *fiber.Ctx) error {
	var evento models.Evento
	if err := c.BodyParser(&evento); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al analizar el cuerpo de la solicitud",
		})
	}

	err := repositories.CreateEvento(&evento)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al crear el evento",
		})
	}

	return c.JSON(evento)
}

func (h *EventoHandler) GetAllEventos(c *fiber.Ctx) error {
	eventos, length, err := repositories.GetAllEventos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener todos los eventos",
		})
	}

	return c.JSON(fiber.Map{
		"eventos": eventos,
		"length":  length,
	})
}

func (h *EventoHandler) GetEventoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	evento, err := repositories.GetEventoByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener el evento por ID",
		})
	}

	return c.JSON(evento)
}

func (h *EventoHandler) UpdateEvento(c *fiber.Ctx) error {
	var updatedEvento models.Evento
	if err := c.BodyParser(&updatedEvento); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al analizar el cuerpo de la solicitud",
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	err = repositories.UpdateEvento(uint(id), &updatedEvento)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar el evento",
		})
	}

	return c.JSON(updatedEvento)
}

func (h *EventoHandler) DeleteEvento(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	err = repositories.DeleteEvento(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar el evento",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Evento eliminado exitosamente",
	})
}

func (h *EventoHandler) GetEventosByFecha(c *fiber.Ctx) error {
	dateStr := c.Params("fecha")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fecha no válido",
		})
	}

	eventos, length, err := repositories.GetEventosByFecha(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener eventos por fecha",
		})
	}

	return c.JSON(fiber.Map{
		"eventos": eventos,
		"length":  length,
	})
}

func (h *EventoHandler) GetEventosByIDUsuario(c *fiber.Ctx) error {
	idUsuario, err := strconv.Atoi(c.Params("idUsuario"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de usuario inválido",
		})
	}

	eventos, length, err := repositories.GetEventosByIDUsuario(uint(idUsuario))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener eventos por ID de usuario",
		})
	}

	return c.JSON(fiber.Map{
		"eventos": eventos,
		"length":  length,
	})
}
