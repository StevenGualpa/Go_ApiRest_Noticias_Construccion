package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"noticias_uteq/models"
	"noticias_uteq/repositories"
)

type UsuarioFacultadHandler struct {
}

type UsuarioFacultadResponse struct {
	UsuarioFacultades []models.UsuarioFacultad `json:"usuarioFacultades"`
	Lenght            int                      `json:"lenght"`
}
type UsuarioFacultadRpResponse struct {
	UsuarioFacultadesrp []models.UsuarioFacultadNombre `json:"usuarioFacultadesrp"`
	Lenght              int                            `json:"lenght"`
}

func (h *UsuarioFacultadHandler) CreateUsuarioFacultad(c *fiber.Ctx) error {
	var usuarioFacultad models.UsuarioFacultad
	err := c.BodyParser(&usuarioFacultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	err = repositories.CreateUsuarioFacultad(&usuarioFacultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar la preferencia de facultad",
		})
	}

	return c.JSON(fiber.Map{
		"msg":             "Preferencia de facultad registrada correctamente",
		"usuarioFacultad": usuarioFacultad,
	})
}

func (h *UsuarioFacultadHandler) GetFacultadPreferencesByUsuario(c *fiber.Ctx) error {
	usuarioID := c.Params("usuarioID")
	usuarioFacultades, length, err := repositories.GetAllFacultadPreferences(usuarioID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las preferencias de facultad",
		})
	}

	return c.JSON(&UsuarioFacultadResponse{
		UsuarioFacultades: usuarioFacultades,
		Lenght:            length,
	})
}

func (h *UsuarioFacultadHandler) UpsertUsuarioFacultad(c *fiber.Ctx) error {
	var usuarioFacultades []*models.UsuarioFacultad
	err := c.BodyParser(&usuarioFacultades)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	if len(usuarioFacultades) > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Solo puedes registrar un máximo de 3 preferencias",
		})
	}

	err = repositories.UpsertUsuarioFacultad(usuarioFacultades)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al actualizar las preferencias de facultad",
		})
	}

	return c.JSON(fiber.Map{
		"msg":             "Preferencias de facultad actualizadas correctamente",
		"usuarioFacultad": usuarioFacultades,
	})
}

// No se usa por ahora
func (h *UsuarioFacultadHandler) DeleteUsuarioFacultad(c *fiber.Ctx) error {
	id := c.Params("id")
	err := repositories.DeleteUsuarioFacultad(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar la preferencia de facultad",
		})
	}

	return c.JSON(fiber.Map{
		"success": "Preferencia de facultad eliminada correctamente",
	})
}

// Para relaciones
func (h *UsuarioFacultadHandler) CreateUsuarioFacultadRp(c *fiber.Ctx) error {
	var usuarioFacultad models.UsuarioFacultadNombre
	err := c.BodyParser(&usuarioFacultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	// TODO: Asegúrate de que el campo FacultadNombre se llena correctamente
	// y valida el nombre de la facultad aquí.

	err = repositories.CreateUsuarioFacultadRp(&usuarioFacultad)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al registrar la preferencia de facultad",
		})
	}

	return c.JSON(fiber.Map{
		"msg":             "Preferencia de facultad registrada correctamente",
		"usuarioFacultad": usuarioFacultad,
	})
}

func (h *UsuarioFacultadHandler) GetFacultadPreferencesByUsuarioRp(c *fiber.Ctx) error {
	usuarioID := c.Params("usuarioID")
	usuarioFacultades, length, err := repositories.GetAllFacultadPreferencesRp(usuarioID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las preferencias de facultad",
		})
	}

	return c.JSON(&UsuarioFacultadRpResponse{
		UsuarioFacultadesrp: usuarioFacultades,
		Lenght:              length,
	})
}

func (h *UsuarioFacultadHandler) UpsertUsuarioFacultadRp(c *fiber.Ctx) error {
	var usuarioFacultades []*models.UsuarioFacultadNombre
	err := c.BodyParser(&usuarioFacultades)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al convertir los datos",
		})
	}

	if len(usuarioFacultades) > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Solo puedes registrar un máximo de 3 preferencias",
		})
	}

	// TODO: Asegúrate de que los campos FacultadNombre se llenan correctamente
	// y valida los nombres de las facultades aquí.

	err = repositories.UpsertUsuarioFacultadRp(usuarioFacultades)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al actualizar las preferencias de facultad",
		})
	}

	return c.JSON(fiber.Map{
		"msg":             "Preferencias de facultad actualizadas correctamente",
		"usuarioFacultad": usuarioFacultades,
	})
}
