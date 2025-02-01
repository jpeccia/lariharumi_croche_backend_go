package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

func RegisterUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
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