package main

import (
	"log"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
)

func main() {
	config.ConnectDB()

	log.Println("Servidor rodando na porta 8080....")
}
