package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

// CronService gerencia tarefas agendadas
type CronService struct {
	baseURL string
	client  *http.Client
}

var cronService *CronService

// InitCronService inicializa o serviço de cron
func InitCronService() {
	baseURL := os.Getenv("BASEURL")
	if baseURL == "" {
		// Usa a porta do Render se disponível, senão localhost
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		baseURL = "http://localhost:" + port
	}

	cronService = &CronService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Inicia o cronjob em uma goroutine
	go cronService.startHealthCheckCron()
}

// startHealthCheckCron executa ping a cada 25 segundos para manter a aplicação ativa
func (cs *CronService) startHealthCheckCron() {
	ticker := time.NewTicker(25 * time.Second) // A cada 25 segundos (menos que 30s do Render)
	defer ticker.Stop()

	log.Println("🔄 Cronjob de health check iniciado - ping a cada 25 segundos")

	for {
		select {
		case <-ticker.C:
			cs.performHealthCheck()
		}
	}
}

// performHealthCheck faz uma requisição para manter a aplicação ativa
func (cs *CronService) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tenta diferentes endpoints para manter a aplicação ativa
	endpoints := []string{
		"/categories", // Endpoint público
		"/products",   // Endpoint público
		"/health",     // Endpoint de health check (se existir)
	}

	for _, endpoint := range endpoints {
		url := cs.baseURL + endpoint

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Printf("❌ Erro ao criar requisição para %s: %v", url, err)
			continue
		}

		resp, err := cs.client.Do(req)
		if err != nil {
			log.Printf("❌ Erro ao fazer ping para %s: %v", url, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("✅ Ping bem-sucedido para %s (status: %d)", endpoint, resp.StatusCode)
			return // Sucesso, não precisa tentar outros endpoints
		} else {
			log.Printf("⚠️ Ping com status inesperado para %s: %d", endpoint, resp.StatusCode)
		}
	}
}

// StopCronService para o serviço de cron (para testes ou shutdown graceful)
func StopCronService() {
	if cronService != nil {
		log.Println("🛑 Parando serviço de cron...")
		// Aqui você pode implementar lógica de parada se necessário
	}
}
