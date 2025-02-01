package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/util"
)

// AuthMiddleware é um middleware que verifica o token JWT nas requisições
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token da requisição
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// Remove o prefixo "Bearer " se existir
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Valida o token
		userID, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Adiciona o ID do usuário ao contexto
		c.Set("userID", userID)

		// Continuar o processamento
		c.Next()
	}
}
