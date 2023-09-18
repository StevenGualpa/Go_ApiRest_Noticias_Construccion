package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
)

type CarreraHandler struct {
}

type CarreraResponse struct {
	Carreras []models.Carrera `json:"carreras"`
	Lenght   int              `json:"lenght"`
}

func (h *CarreraHandler) CreateCarrera(c *fiber.Ctx) error {
	var carrera models.Carrera
	err := c.BodyParser(&carrera)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.RegisterCarrera(&carrera)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(), // Devuelve el mensaje de error específico del repositorio
		})
	}

	return c.JSON(fiber.Map{
		"msg":     "Carrera registrada correctamente",
		"carrera": carrera,
	})
}

func (h *CarreraHandler) UpdateCarrera(c *fiber.Ctx) error {
	// Obtén el ID de la carrera de los parámetros de la ruta
	id := c.Params("id")

	// Crea un objeto de Carrera para almacenar los datos de la solicitud
	var carrera models.Carrera
	if err := c.BodyParser(&carrera); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se puede analizar el JSON",
		})
	}

	// Llama a la función UpdateCarrera del repositorio
	err := repositories.UpdateCarrera(id, &carrera)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar la carrera",
		})
	}

	// Si todo va bien, devuelve la carrera actualizada
	return c.JSON(carrera)
}

func (h *CarreraHandler) DeleteCarrera(c *fiber.Ctx) error {
	// Obtén el ID de la carrera del parámetro de la ruta
	id := c.Params("id")

	// Llama a la función DeleteCarrera del repositorio
	err := repositories.DeleteCarrera(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar la carrera",
		})
	}

	// Si todo va bien, devolver un mensaje de éxito
	return c.JSON(fiber.Map{
		"success": "Carrera eliminada correctamente",
	})
}

func (h *CarreraHandler) GetAllCarreras(c *fiber.Ctx) error {
	carreras, length, err := repositories.GetAllCarreras()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(&CarreraResponse{
		Carreras: carreras,
		Lenght:   length,
	})
}

func (h *CarreraHandler) GetCarrerasByFacultadID(c *fiber.Ctx) error {
	facultadID := c.Params("facultadID")

	carreras, length, err := repositories.GetCarrerasByFacultadID(facultadID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las carreras",
		})
	}

	return c.JSON(&CarreraResponse{
		Carreras: carreras,
		Lenght:   length,
	})
}
