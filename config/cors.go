package config

import (
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Configura o CORS para a aplicação
func SetupCors(r *gin.Engine) {
	// Obtém a URL do front-end do arquivo .env
	frontEndURL := os.Getenv("FRONTEND_URL")

	// Configuração do CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontEndURL},    // URL do Front-end a partir do .env
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},  // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Cabeçalhos permitidos
		AllowCredentials: true,  // Permite enviar cookies e credenciais
	}))
}
