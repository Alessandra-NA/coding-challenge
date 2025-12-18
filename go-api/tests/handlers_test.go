package tests

import (
	"bytes"
	"encoding/json"
	"interseguro-challenge/go-api/handlers"
	"interseguro-challenge/go-api/middleware"
	"interseguro-challenge/go-api/models"
	"interseguro-challenge/go-api/utils"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// setupTestApp crea una aplicación Fiber para testing
func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware())
	return app
}

func TestLoginHandler(t *testing.T) {
	// Guardar valores originales de variables de entorno
	originalUsername := os.Getenv("AUTH_USERNAME")
	originalPassword := os.Getenv("AUTH_PASSWORD")
	defer func() {
		if originalUsername != "" {
			os.Setenv("AUTH_USERNAME", originalUsername)
		} else {
			os.Unsetenv("AUTH_USERNAME")
		}
		if originalPassword != "" {
			os.Setenv("AUTH_PASSWORD", originalPassword)
		} else {
			os.Unsetenv("AUTH_PASSWORD")
		}
	}()

	// Configurar credenciales de prueba
	os.Setenv("AUTH_USERNAME", "testuser")
	os.Setenv("AUTH_PASSWORD", "testpass")

	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON body",
			body:           "invalid json",
			expectedStatus: 400,
			expectedError:  "cuerpo de solicitud inválido",
		},
		{
			name: "missing username",
			body: models.LoginRequest{
				Password: "testpass",
			},
			expectedStatus: 400,
			expectedError:  "usuario y contraseña son requeridos",
		},
		{
			name: "missing password",
			body: models.LoginRequest{
				Username: "testuser",
			},
			expectedStatus: 400,
			expectedError:  "usuario y contraseña son requeridos",
		},
		{
			name: "missing both username and password",
			body: models.LoginRequest{
				Username: "",
				Password: "",
			},
			expectedStatus: 400,
			expectedError:  "usuario y contraseña son requeridos",
		},
		{
			name: "invalid credentials - wrong username",
			body: models.LoginRequest{
				Username: "wronguser",
				Password: "testpass",
			},
			expectedStatus: 401,
			expectedError:  "credenciales inválidas",
		},
		{
			name: "invalid credentials - wrong password",
			body: models.LoginRequest{
				Username: "testuser",
				Password: "wrongpass",
			},
			expectedStatus: 401,
			expectedError:  "credenciales inválidas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			app.Post("/auth/login", handlers.LoginHandler)

			var reqBody []byte
			var err error
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("Error marshaling request body: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Error decoding response: %v", err)
			}

			if errorMsg, ok := result["error"].(string); ok {
				if !strings.Contains(errorMsg, tt.expectedError) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectedError, errorMsg)
				}
			} else {
				t.Errorf("Expected error field in response, got %v", result)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	app := setupTestApp()
	app.Get("/health", handlers.HealthHandler)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result models.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if result.Status != "saludable" {
		t.Errorf("Expected status 'saludable', got %q", result.Status)
	}

	if result.Service != "go-api" {
		t.Errorf("Expected service 'go-api', got %q", result.Service)
	}
}

func TestRotateHandler(t *testing.T) {
	// Generar un token válido para las pruebas
	token, _, err := utils.GenerateJWT("testuser")
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	tests := []struct {
		name           string
		body           interface{}
		authToken      string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON body",
			body:           "invalid json",
			authToken:      token,
			expectedStatus: 400,
			expectedError:  "cuerpo de solicitud inválido",
		},
		{
			name: "empty matrix",
			body: models.MatrixRequest{
				Matrix: [][]int{},
			},
			authToken:      token,
			expectedStatus: 400,
			expectedError:  "La matriz no puede estar vacía",
		},
		{
			name: "non-rectangular matrix",
			body: models.MatrixRequest{
				Matrix: [][]int{{1, 2, 3}, {4, 5}},
			},
			authToken:      token,
			expectedStatus: 400,
			expectedError:  "La matriz debe ser rectangular",
		},
		{
			name: "matrix with empty row",
			body: models.MatrixRequest{
				Matrix: [][]int{{1, 2}, {}},
			},
			authToken:      token,
			expectedStatus: 400,
			expectedError:  "La matriz debe ser rectangular",
		},
		{
			name: "matrix exceeds size limit",
			body: models.MatrixRequest{
				Matrix: func() [][]int {
					matrix := make([][]int, 1001)
					for i := range matrix {
						matrix[i] = []int{1, 2}
					}
					return matrix
				}(),
			},
			authToken:      token,
			expectedStatus: 400,
			expectedError:  "La matriz tiene demasiadas filas",
		},
		{
			name:           "missing authorization header",
			body:           models.MatrixRequest{Matrix: [][]int{{1, 2}, {3, 4}}},
			authToken:      "",
			expectedStatus: 401,
			expectedError:  "falta el header de autorización",
		},
		{
			name:           "invalid authorization format",
			body:           models.MatrixRequest{Matrix: [][]int{{1, 2}, {3, 4}}},
			authToken:      "InvalidFormat ",
			expectedStatus: 401,
			expectedError:  "formato de header de autorización inválido",
		},
		{
			name:           "missing token after Bearer",
			body:           models.MatrixRequest{Matrix: [][]int{{1, 2}, {3, 4}}},
			authToken:      "Bearer",
			expectedStatus: 401,
			expectedError:  "formato de header de autorización inválido",
		},
		{
			name:           "invalid token",
			body:           models.MatrixRequest{Matrix: [][]int{{1, 2}, {3, 4}}},
			authToken:      "Bearer invalid.token.here",
			expectedStatus: 401,
			expectedError:  "token inválido o expirado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			protected := app.Group("")
			protected.Use(middleware.AuthMiddleware())
			protected.Post("/rotate", handlers.RotateHandler)

			var reqBody []byte
			var err error
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("Error marshaling request body: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/rotate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.authToken != "" {
				// Usar el token tal cual si ya tiene un formato específico (Bearer, InvalidFormat, etc.)
				// Solo agregar "Bearer " si es un token JWT simple sin prefijo
				authHeader := tt.authToken
				if strings.HasPrefix(authHeader, "Bearer ") || 
				   strings.HasPrefix(authHeader, "Bearer") ||
				   strings.HasPrefix(authHeader, "InvalidFormat") {
					// Ya tiene un formato específico, usar tal cual
				} else {
					// Es un token JWT simple, agregar "Bearer "
					authHeader = "Bearer " + authHeader
				}
				req.Header.Set("Authorization", authHeader)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("Error decoding response: %v", err)
			}

			if errorMsg, ok := result["error"].(string); ok {
				if !strings.Contains(errorMsg, tt.expectedError) {
					t.Errorf("Expected error message to contain %q, got %q", tt.expectedError, errorMsg)
				}
			} else {
				t.Errorf("Expected error field in response, got %v", result)
			}
		})
	}
}
