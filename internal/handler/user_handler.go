package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

func RegisterUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required,min=2,max=100"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao registrar usuário"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// LoginUser recebe as credenciais e retorna um token JWT
func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := service.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Retorna o token para o usuário
	c.JSON(http.StatusOK, gin.H{"token": token})
}
