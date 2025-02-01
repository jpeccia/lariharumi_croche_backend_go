package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{})
	if err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}
}

func ConnectDB() {
	LoadEnv()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados", err)
	}

	fmt.Println("Banco de dados conectado!")
	MigrateDB(DB)
}
