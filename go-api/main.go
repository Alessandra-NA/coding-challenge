package main

import (
	"log"
	"os"

	"interseguro-challenge/go-api/handlers"
	"interseguro-challenge/go-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde archivo .env (si existe)
	_ = godotenv.Load()

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New()) // Recuperar de pánicos
	app.Use(middleware.LoggerMiddleware()) // Logging estructurado con IDs de solicitud

	// Configuración CORS para frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // En producción, especificar orígenes exactos
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		ExposeHeaders:    "X-Request-ID",
		AllowCredentials: true,
	}))

	// Rutas públicas
	app.Post("/auth/login", handlers.LoginHandler)
	app.Get("/health", handlers.HealthHandler)

	// Rutas protegidas (requieren JWT)
	protected := app.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Post("/rotate", handlers.RotateHandler)

	// Obtener puerto desde variable de entorno
	port := os.Getenv("PORT")
	if port == "" {
			port = os.Getenv("API1_PORT")
	}
	if port == "" {
			port = "8080"
	}

	// Iniciar servidor
	log.Printf("Iniciando servidor Go API en el puerto %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

