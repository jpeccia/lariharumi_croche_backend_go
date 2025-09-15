package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/config"
)

// HealthResponse representa a resposta do health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Cache     string    `json:"cache"`
}

// HealthCheck retorna o status de saúde da aplicação
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Database:  "connected",
		Cache:     "available",
	}

	// Verifica conexão com o banco de dados
	if config.DB == nil {
		response.Status = "unhealthy"
		response.Database = "disconnected"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Testa uma query simples no banco
	sqlDB, err := config.DB.DB()
	if err != nil {
		response.Status = "unhealthy"
		response.Database = "error"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Status = "unhealthy"
		response.Database = "ping_failed"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Verifica cache (Redis)
	if !isCacheAvailable() {
		response.Cache = "unavailable"
		// Não marca como unhealthy se apenas o cache estiver indisponível
	}

	c.JSON(http.StatusOK, response)
}

// isCacheAvailable verifica se o cache está disponível
func isCacheAvailable() bool {
	// Aqui você pode implementar uma verificação específica do Redis
	// Por enquanto, retorna true assumindo que o cache está funcionando
	return true
}

// Ping é um endpoint simples para manter a aplicação ativa
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now(),
		"status":    "alive",
	})
}
