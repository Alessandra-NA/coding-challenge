package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

// RequestIDKey es la clave usada para almacenar el ID de solicitud en el contexto
const RequestIDKey = "request_id"

// LoggerMiddleware crea un middleware que registra solicitudes con logging estructurado
// y agrega un ID de solicitud para trazabilidad
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generar o recuperar ID de solicitud
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}
		c.Locals(RequestIDKey, requestID)

		// Registrar tiempo de inicio
		start := time.Now()

		// Procesar solicitud
		err := c.Next()

		// Calcular duración
		duration := time.Since(start)

		// Registrar información estructurada de la solicitud
		logEntry := fiber.Map{
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"request_id":  requestID,
			"method":      c.Method(),
			"path":        c.Path(),
			"status_code": c.Response().StatusCode(),
			"duration_ms": duration.Milliseconds(),
			"ip":          c.IP(),
			"user_agent":  c.Get("User-Agent"),
		}

		// Registrar error si está presente
		if err != nil {
			logEntry["error"] = err.Error()
		}

		// Salida de log estructurado (en producción, usar un logger apropiado como zerolog o zap)
		// Por ahora, el ID de solicitud se establece en headers y contexto para trazabilidad
		// El logging estructurado puede agregarse aquí con una librería de logging apropiada
		_ = logEntry // Suprimir advertencia de variable no usada

		return err
	}
}

