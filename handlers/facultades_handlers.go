package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
)

type FacultadHandler struct {
}
type FacultadResponse struct {
	Facultades []models.Facultad `json:"facultad"`
	Lenght     int               `json:"lenght"`
}

func (h *FacultadHandler) CreateFacultad(c *fiber.Ctx) error {
	var facultad models.Facultad
	err := c.BodyParser(&facultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.RegisterFacultad(&facultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar la facultad",
		})
	}

	return c.JSON(fiber.Map{
		"msg":      "Facultad registrada correctamente",
		"facultad": facultad,
	})
}

func (h *FacultadHandler) UpdateFacultad(c *fiber.Ctx) error {
	// Obtén el ID de la facultad de los parámetros de la ruta
	id := c.Params("id")

	// Crea un objeto de Facultad para almacenar los datos de la solicitud
	var facultad models.Facultad
	if err := c.BodyParser(&facultad); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se puede analizar el JSON",
		})
	}

	// Llama a la función UpdateFacultad del repositorio
	err := repositories.UpdateFacultad(id, &facultad)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar la facultad",
		})
	}

	// Si todo va bien, devuelve la facultad actualizada
	return c.JSON(facultad)
}

func (h *FacultadHandler) DeleteFacultad(c *fiber.Ctx) error {
	// Obtén el ID de la facultad del parámetro de la ruta
	id := c.Params("id")

	// Llama a la función DeleteFacultad del repositorio
	err := repositories.DeleteFacultad(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar la facultad",
		})
	}

	// Si todo va bien, devolver un mensaje de éxito
	return c.JSON(fiber.Map{
		"success": "Facultad eliminada correctamente",
	})
}

func (u *FacultadHandler) GetAllFacultades(c *fiber.Ctx) error {
	facultades, length, err := repositories.GetAllFacultades()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&FacultadResponse{
		Facultades: facultades,
		Lenght:     length,
	})
}
