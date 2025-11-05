package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

// GET /promotion (público)
func GetPromotion(c *gin.Context) {
	p, err := service.GetActivePromotion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter promoção: " + err.Error()})
		return
	}
	if p == nil {
		c.Status(http.StatusNoContent) // 204 quando não há promoção ativa/definida
		return
	}
	c.JSON(http.StatusOK, p)
}

// PUT /promotion (admin)
func UpdatePromotion(c *gin.Context) {
	var req model.Promotion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	saved, err := service.SavePromotion(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400 para validações
		return
	}

	c.JSON(http.StatusOK, saved)
}
