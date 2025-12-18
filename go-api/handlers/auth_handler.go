package handlers

import (
	"interseguro-challenge/go-api/models"
	"interseguro-challenge/go-api/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

// LoginHandler maneja las solicitudes POST /auth/login
// Valida credenciales y retorna token JWT
// Para propósitos de demo, usa variables de entorno para credenciales
func LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parsear cuerpo de solicitud
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuerpo de solicitud inválido",
		})
	}

	// Validar campos requeridos
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "usuario y contraseña son requeridos",
		})
	}

	// Validar credenciales (usando variables de entorno)
	// En producción, esto debería consultar una base de datos
	validUsername := os.Getenv("AUTH_USERNAME")
	if validUsername == "" {
		validUsername = "admin" // Por defecto para desarrollo
	}

	validPassword := os.Getenv("AUTH_PASSWORD")
	if validPassword == "" {
		validPassword = "admin123" // Por defecto para desarrollo
	}

	if req.Username != validUsername || req.Password != validPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "credenciales inválidas",
		})
	}

	// Generar token JWT
	token, expiresIn, err := utils.GenerateJWT(req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error al generar el token",
		})
	}

	// Retornar token
	response := models.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}

	return c.JSON(response)
}

