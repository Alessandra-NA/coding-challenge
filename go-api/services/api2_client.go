package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"interseguro-challenge/go-api/models"
)

// API2Client maneja la comunicación HTTP con API 2
type API2Client struct {
	client    *http.Client
	baseURL   string
	maxRetries int
}

// NewAPI2Client crea un nuevo cliente API 2 con cliente HTTP configurado
// Usa connection pooling y timeouts para estar listo para producción
func NewAPI2Client() *API2Client {
	baseURL := os.Getenv("API2_URL")
	if baseURL == "" {
		// En producción, API2_URL debe estar configurado
		// Fallback solo para desarrollo local fuera de Docker
		baseURL = "http://localhost:3001"
		fmt.Println("ADVERTENCIA: API2_URL no está configurado, usando fallback localhost:3001")
	}

	maxRetries := 3
	if maxRetriesEnv := os.Getenv("API2_MAX_RETRIES"); maxRetriesEnv != "" {
		fmt.Sscanf(maxRetriesEnv, "%d", &maxRetries)
	}

	// Configurar cliente HTTP con connection pooling y timeouts
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second, // Timeout de 5 segundos para solicitudes
	}

	return &API2Client{
		client:     client,
		baseURL:    baseURL,
		maxRetries: maxRetries,
	}
}

// ProcessMatrix envía matrices original y rotada a API 2 y recibe estadísticas
// Implementa lógica de reintento con exponential backoff
func (c *API2Client) ProcessMatrix(originalMatrix [][]int, rotatedMatrix [][]int, requestID string) (*models.API2Response, error) {
	request := models.API2Request{
		OriginalMatrix: originalMatrix,
		RotatedMatrix:  rotatedMatrix,
		RequestID:      requestID,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error al serializar la solicitud: %w", err)
	}

	var lastErr error
	backoff := 100 * time.Millisecond // Backoff inicial

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: esperar antes de reintentar
			time.Sleep(backoff)
			backoff *= 2 // Duplicar el backoff para el siguiente reintento
		}

		resp, err := c.client.Post(
			c.baseURL+"/process",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			lastErr = fmt.Errorf("error en la solicitud: %w", err)
			continue // Reintentar
		}

		defer resp.Body.Close()

		// Verificar código de estado
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API 2 retornó estado %d: %s", resp.StatusCode, string(body))
			
			// No reintentar en errores 4xx (errores del cliente)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				break
			}
			continue // Reintentar en errores 5xx
		}

		// Parsear respuesta
		var api2Response models.API2Response
		if err := json.NewDecoder(resp.Body).Decode(&api2Response); err != nil {
			lastErr = fmt.Errorf("error al decodificar la respuesta: %w", err)
			continue // Reintentar en errores de decodificación (pueden ser transitorios)
		}

		return &api2Response, nil
	}

	return nil, fmt.Errorf("falló después de %d intentos: %w", c.maxRetries+1, lastErr)
}

