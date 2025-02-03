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

// Claims com Role adicional
type CustomClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// GenerateToken gera um token JWT para o usuário
func GenerateToken(user *model.User) (string, error) {
	// Verifica a role antes de gerar o token
	fmt.Println("Gerando token para usuário com role:", user.Role)

	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "laribackend",
		},
		Role: string(user.Role), // Converte explicitamente para string
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken valida o token e retorna o ID do usuário
func ParseToken(tokenString string) (uint, string, error) {
	// Parseia o token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, "", err
	}

	// Extrai os claims personalizados (com role)
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", fmt.Errorf("não foi possível converter os claims")
	}

	// Converte o Subject (string) para uint
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", err
	}

	// Retorna o ID do usuário e a role
	return uint(id), claims.Role, nil
}
