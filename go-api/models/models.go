package models

import "github.com/golang-jwt/jwt/v5"

// MatrixRequest representa el cuerpo de solicitud para rotación de matriz
type MatrixRequest struct {
	Matrix [][]int `json:"matrix" validate:"required,min=1"`
}

// Statistics representa las estadísticas calculadas de API 2
type Statistics struct {
	MaxValue           int     `json:"max_value"`
	MinValue           int     `json:"min_value"`
	Average            float64 `json:"average"`
	TotalSum           int     `json:"total_sum"`
	OriginalIsDiagonal bool    `json:"original_is_diagonal"`
	RotatedIsDiagonal  bool    `json:"rotated_is_diagonal"`
	CalculationTimeMs  float64 `json:"calculation_time_ms,omitempty"`
}

// RotateResponse representa la respuesta final al cliente
type RotateResponse struct {
	OriginalMatrix   [][]int     `json:"original_matrix"`
	RotatedMatrix    [][]int     `json:"rotated_matrix"`
	Statistics       *Statistics `json:"statistics"`
	ProcessingTimeMs float64     `json:"processing_time_ms"`
}

// API2Request representa la solicitud enviada a API 2
type API2Request struct {
	OriginalMatrix [][]int `json:"original_matrix"`
	RotatedMatrix  [][]int `json:"rotated_matrix"`
	RequestID      string  `json:"request_id"`
}

// API2Response representa la respuesta de API 2
type API2Response struct {
	RotatedMatrix [][]int    `json:"rotated_matrix"`
	Statistics    Statistics `json:"statistics"`
	RequestID     string     `json:"request_id"`
}

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa la respuesta de login con token JWT
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// HealthResponse representa la respuesta del health check
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// JWTClaims representa los claims del token JWT
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

