package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"interseguro-challenge/go-api/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

// getJWTSecret obtiene el secreto JWT desde variable de entorno
// Por defecto usa un secreto de desarrollo si no está configurado (debe cambiarse en producción)
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "development-secret-key-change-in-production"
	}
	return secret
}

// GetJWTExpiration obtiene el tiempo de expiración JWT en segundos desde variable de entorno
// Por defecto 3600 segundos (1 hora)
func GetJWTExpiration() time.Duration {
	expirationStr := os.Getenv("JWT_EXPIRATION")
	if expirationStr == "" {
		return 3600 * time.Second
	}

	// Parsear como segundos
	var expiration int
	if _, err := fmt.Sscanf(expirationStr, "%d", &expiration); err == nil {
		return time.Duration(expiration) * time.Second
	}

	return 3600 * time.Second
}

// GenerateJWT genera un token JWT para el nombre de usuario dado
// Retorna el string del token y el tiempo de expiración en segundos
func GenerateJWT(username string) (string, int, error) {
	expiration := GetJWTExpiration()
	expirationSeconds := int(expiration.Seconds())

	claims := models.JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-api",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationSeconds, nil
}

// ValidateJWT valida un token JWT y retorna los claims
// Retorna los claims si es válido, error en caso contrario
func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

