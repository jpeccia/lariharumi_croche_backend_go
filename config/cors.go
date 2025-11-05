package config

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	frontendURL := os.Getenv("FRONTEND_URL")

	log.Println("CORS AllowOrigins:", frontendURL)

	config := cors.Config{
		AllowOrigins:     filterEmpty([]string{frontendURL, "https://larifazcroche.vercel.app"}),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "XMLHttpRequest"},
		AllowCredentials: true,
	}

	return cors.New(config)
}

func filterEmpty(in []string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if v != "" {
			out = append(out, v)
		}
	}
	log.Println("CORS AllowOrigins:", out)
	return out
}
