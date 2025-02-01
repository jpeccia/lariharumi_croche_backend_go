package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/auth")
	{
		api.POST("/register", handler.RegisterUser)
	}

	return r
}
