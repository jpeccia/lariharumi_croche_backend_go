package config

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura o middleware de CORS
func CORSMiddleware() gin.HandlerFunc {
	frontendURL := os.Getenv("FRONTEND_URL") // Certifique-se de que FRONTEND_URL está definido no seu .env
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // Use o valor padrão se não estiver no .env
	}

	config := cors.Config{
		AllowOrigins:     []string{frontendURL}, // Permite requisições da URL do frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true, // Permite cookies ou credenciais nas requisições
	}

	return cors.New(config)
}
