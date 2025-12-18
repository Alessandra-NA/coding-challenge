package handlers

import (
	"interseguro-challenge/go-api/middleware"
	"interseguro-challenge/go-api/models"
	"interseguro-challenge/go-api/services"
	"interseguro-challenge/go-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RotateHandler maneja las solicitudes POST /rotate
// Rota la matriz 90° en sentido horario y obtiene estadísticas de API 2
func RotateHandler(c *fiber.Ctx) error {
	startTime := time.Now()

	// Obtener ID de solicitud del middleware
	requestID, _ := c.Locals(middleware.RequestIDKey).(string)

	// Parsear cuerpo de solicitud
	var req models.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "cuerpo de solicitud inválido",
			"request_id": requestID,
		})
	}

	// Validar matriz
	if err := utils.ValidateMatrix(req.Matrix); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": requestID,
		})
	}

	// Validar tamaño de matriz (prevenir problemas de memoria)
	maxSize := 1000
	if err := utils.ValidateMatrixSize(req.Matrix, maxSize, maxSize); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": requestID,
		})
	}

	// Rotar matriz
	rotatedMatrix := services.RotateMatrix90Clockwise(req.Matrix)

	// Enviar ambas matrices a API 2 para estadísticas
	api2Client := services.NewAPI2Client()
	api2Response, err := api2Client.ProcessMatrix(req.Matrix, rotatedMatrix, requestID)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":      "error al procesar matriz con API 2: " + err.Error(),
			"request_id": requestID,
		})
	}

	// Calcular tiempo de procesamiento
	processingTime := time.Since(startTime)

	// Construir respuesta
	response := models.RotateResponse{
		OriginalMatrix:   req.Matrix,
		RotatedMatrix:    rotatedMatrix,
		Statistics:       &api2Response.Statistics,
		ProcessingTimeMs: float64(processingTime.Nanoseconds()) / 1e6, // Convertir a milisegundos
	}

	return c.JSON(response)
}

// HealthHandler maneja las solicitudes GET /health
func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(models.HealthResponse{
		Status:  "saludable",
		Service: "go-api",
	})
}

