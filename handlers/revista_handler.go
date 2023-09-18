package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/repositories"
)

type RevistaHandler struct {
}

// GetRevistas obtiene todos los registros de revista de la base de datos.
func (revista *RevistaHandler) GetRevistas(c *fiber.Ctx) error {

	resultado, _ := repositories.Revistas()
	return c.Status(fiber.StatusOK).JSON(resultado)
}

// PostRevista este handler crea un registo para una revista.
func (revista *RevistaHandler) PostRevista(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("La revista se ha creado correctamente")
}

// DeleteRevista elimina el registro de la base de datos.
func (revista *RevistaHandler) DeleteRevista(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al obtener el id de la revista",
		})
	}

	log.Println(id)
	return c.Status(fiber.StatusOK).SendString("La revista se ha eliminado correctamente")
}
