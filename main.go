package main

import (
	"log"
	"os"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/router"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

func main() {
	config.ConnectDB()

	// Inicializa o serviço de upload assíncrono
	service.InitUploadService(5) // 5 workers para uploads

	// Inicializa o cronjob para manter a aplicação ativa no Render
	service.InitCronService()

	r := router.SetupRouter()

	// Usa a porta do Render se disponível, senão 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s....", port)
	log.Println("🔄 Cronjob de health check ativo - aplicação será mantida viva no Render")
	r.Run(":" + port)
}
