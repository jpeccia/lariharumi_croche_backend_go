package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/util"
)

// AuthMiddleware é um middleware que verifica o token JWT e a role do usuário
func AuthMiddleware(roleRequired string) gin.HandlerFunc {
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

		// Valida o token e extrai o userID
		userID, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Busca o usuário no banco de dados com base no userID
		user, err := repository.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
			c.Abort()
			return
		}

		// Verifica se o usuário tem a role correta
		if roleRequired != "" && user.Role != model.AdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado: permissões insuficientes"})
			c.Abort()
			return
		}

		// Adiciona o usuário ao contexto para outras operações
		c.Set("userID", userID)

		// Continua a execução
		c.Next()
	}
}