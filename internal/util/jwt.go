package util

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Segredo para assinar o token

// GenerateToken gera um token JWT para o usuário
func GenerateToken(user *model.User) (string, error) {
	// Define o tempo de expiração do token (1 hora)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Criação do token
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID), // Converte o ID para string
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    "laribackend", // Emissor do token
	}

	// Criação do token com a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken valida o token e retorna o ID do usuário
func ParseToken(tokenString string) (uint, error) {
	// Parseia o token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	// Extrai os claims (dados do usuário)
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, fmt.Errorf("não foi possível converter os claims")
	}

	// Converte o Subject (string) para uint
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
