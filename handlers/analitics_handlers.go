package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
	"time"
)

type AnaliticsHandler struct {
}

func (h *AnaliticsHandler) InsertAnalitics(c *fiber.Ctx) error {
	var analitics models.Analitics
	if err := c.BodyParser(&analitics); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err := repositories.InsertAnalitics(&analitics)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al insertar los datos de analítica",
		})
	}

	return c.JSON(analitics)
}

func (h *AnaliticsHandler) GetAllAnalitics(c *fiber.Ctx) error {
	analitics, length, err := repositories.GetAllAnalitics()
	if err != nil {
		log.Println(err)
	}

	return c.JSON(fiber.Map{
		"analitics": analitics,
		"length":    length,
	})
}

func (h *AnaliticsHandler) GetAnaliticsBySeccion(c *fiber.Ctx) error {
	seccion := c.Params("seccion")
	analitics, length, err := repositories.GetAnaliticsBySeccion(seccion)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(fiber.Map{
		"analitics": analitics,
		"length":    length,
	})
}

func (h *AnaliticsHandler) GetAnaliticsByFecha(c *fiber.Ctx) error {
	dateStr := c.Params("fecha")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fecha no válido",
		})
	}

	analitics, length, err := repositories.GetAnaliticsByFecha(date)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(fiber.Map{
		"analitics": analitics,
		"length":    length,
	})
}

func (h *AnaliticsHandler) GetAnaliticsBySeccionAndFecha(c *fiber.Ctx) error {
	seccion := c.Params("seccion")
	dateStr := c.Params("fecha")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fecha no válido",
		})
	}

	analitics, length, err := repositories.GetAnaliticsBySeccionAndFecha(seccion, date)
	if err != nil {
		log.Println(err)
	}

	return c.JSON(fiber.Map{
		"analitics": analitics,
		"length":    length,
	})
}

func (h *AnaliticsHandler) GetNombreCountBySeccion(c *fiber.Ctx) error {
	seccion := c.Params("seccion")
	results, err := repositories.GetNombreCountBySeccion(seccion)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener los datos",
		})
	}

	return c.JSON(fiber.Map{
		"results": results,
	})
}
