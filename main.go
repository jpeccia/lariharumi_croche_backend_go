package main

import (
	"log"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/router"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

func main() {
	config.ConnectDB()

	// Inicializa o serviço de upload assíncrono
	service.InitUploadService(5) // 5 workers para uploads

	r := router.SetupRouter()
	log.Println("Servidor rodando na porta 8080....")
	r.Run(":8080")
}
