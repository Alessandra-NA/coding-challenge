package middleware

import (
	"interseguro-challenge/go-api/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware crea un middleware que valida tokens JWT
// Espera el token en el header Authorization como "Bearer <token>"
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extraer token del header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "falta el header de autorización",
			})
		}

		// Verificar si el header comienza con "Bearer "
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "formato de header de autorización inválido",
			})
		}

		// Extraer token
		token := authHeader[len(bearerPrefix):]
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "falta el token",
			})
		}

		// Validar token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token inválido o expirado",
			})
		}

		// Almacenar claims en contexto para uso en handlers
		c.Locals("username", claims.Username)
		c.Locals("claims", claims)

		return c.Next()
	}
}

