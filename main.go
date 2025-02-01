package main

import (
	"log"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/router"
)

func main() {
	config.ConnectDB()

	r := router.SetupRouter()

	log.Println("Servidor rodando na porta 8080....")
	r.Run(":8080")
}
