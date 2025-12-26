package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/cache"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

var DB *gorm.DB

/**
 * LoadEnv loads environment variables from .env file if available.
 */
func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Aviso: Não foi possível carregar o .env, mas o sistema continuará usando variáveis do ambiente.")
		}
	} else {
		log.Println("Executando sem .env, utilizando variáveis do ambiente.")
	}
}

/**
 * MigrateDB runs auto-migration for all registered models.
 * @param db The GORM database instance.
 */
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{}, &model.Promotion{})
	if err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}
}

/**
 * ConnectDB establishes a PostgreSQL connection with optimized pooling settings.
 * Configures: 10 idle connections, 100 max open connections, 1-hour connection lifetime.
 */
func ConnectDB() {
	LoadEnv()

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), sslMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Falha ao obter conexão sql.DB", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Banco de dados conectado com connection pooling!")
	MigrateDB(DB)

	cache.InitRedis()
}
