package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

/**
 * Shared HTTP client for image uploads.
 * Reusing a single client avoids connection overhead and memory leaks.
 */
var sharedRestyClient = resty.New().
	SetTimeout(30 * time.Second).
	SetRetryCount(2).
	SetRetryWaitTime(1 * time.Second)

// UploadResult representa o resultado de um upload
type UploadResult struct {
	URL   string `json:"url"`
	Error error  `json:"error,omitempty"`
	Index int    `json:"index"`
}

// UploadJob representa um trabalho de upload
type UploadJob struct {
	File      *multipart.FileHeader
	Index     int
	ProductID uint
}

// UploadService gerencia uploads assíncronos
type UploadService struct {
	workerPool chan struct{}
	jobQueue   chan UploadJob
	results    map[uint][]UploadResult
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

var uploadService *UploadService

// InitUploadService inicializa o serviço de upload
func InitUploadService(maxWorkers int) {
	uploadService = &UploadService{
		workerPool: make(chan struct{}, maxWorkers),
		jobQueue:   make(chan UploadJob, 100), // Buffer de 100 jobs
		results:    make(map[uint][]UploadResult),
	}

	// Inicia os workers
	for i := 0; i < maxWorkers; i++ {
		go uploadService.worker()
	}
}

// worker processa jobs de upload
func (us *UploadService) worker() {
	for job := range us.jobQueue {
		us.workerPool <- struct{}{} // Adquire um worker

		go func(j UploadJob) {
			defer func() { <-us.workerPool }() // Libera o worker
			defer us.wg.Done()

			result := us.uploadSingleFile(j.File, j.Index)

			us.mu.Lock()
			if us.results[j.ProductID] == nil {
				us.results[j.ProductID] = make([]UploadResult, 0)
			}
			us.results[j.ProductID] = append(us.results[j.ProductID], result)
			us.mu.Unlock()
		}(job)
	}
}

// uploadSingleFile faz upload de um arquivo para o ImgBB
func (us *UploadService) uploadSingleFile(file *multipart.FileHeader, index int) UploadResult {
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("chave da API ImgBB não encontrada"),
		}
	}

	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("erro ao abrir arquivo: %w", err),
		}
	}
	defer src.Close()

	// Faz upload usando resty
	resp, err := sharedRestyClient.R().
		SetFileReader("image", file.Filename, src).
		SetFormData(map[string]string{"key": apiKey}).
		Post("https://api.imgbb.com/1/upload")

	if err != nil {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("erro no upload: %w", err),
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("upload falhou com status %d", resp.StatusCode()),
		}
	}

	// Extrai URL da resposta
	url := gjson.Get(resp.String(), "data.url").String()
	if url == "" {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("URL não encontrada na resposta"),
		}
	}

	return UploadResult{
		URL:   url,
		Index: index,
	}
}

// UploadProductImagesAsync faz upload assíncrono de múltiplas imagens
func UploadProductImagesAsync(productID uint, files []*multipart.FileHeader) ([]UploadResult, error) {
	if uploadService == nil {
		InitUploadService(5) // 5 workers por padrão
	}

	// Limpa resultados anteriores para este produto
	uploadService.mu.Lock()
	uploadService.results[productID] = make([]UploadResult, 0)
	uploadService.mu.Unlock()

	// Adiciona jobs à fila
	uploadService.wg.Add(len(files))
	for i, file := range files {
		job := UploadJob{
			File:      file,
			Index:     i,
			ProductID: productID,
		}
		uploadService.jobQueue <- job
	}

	// Aguarda todos os uploads terminarem
	uploadService.wg.Wait()

	// Retorna os resultados
	uploadService.mu.RLock()
	results := uploadService.results[productID]
	uploadService.mu.RUnlock()

	return results, nil
}

// UploadCategoryImageAsync faz upload assíncrono de uma imagem de categoria
func UploadCategoryImageAsync(categoryID uint, file *multipart.FileHeader) (string, error) {
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("chave da API ImgBB não encontrada")
	}

	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer src.Close()

	// Cria um canal para receber o resultado
	resultChan := make(chan UploadResult, 1)

	// Faz upload em uma goroutine
	go func() {
		defer close(resultChan)

		resp, err := sharedRestyClient.R().
			SetFileReader("image", file.Filename, src).
			SetFormData(map[string]string{"key": apiKey}).
			Post("https://api.imgbb.com/1/upload")

		if err != nil {
			resultChan <- UploadResult{Error: fmt.Errorf("erro no upload: %w", err)}
			return
		}

		if resp.StatusCode() != http.StatusOK {
			resultChan <- UploadResult{Error: fmt.Errorf("upload falhou com status %d", resp.StatusCode())}
			return
		}

		// Extrai URL da resposta
		url := gjson.Get(resp.String(), "data.url").String()
		if url == "" {
			resultChan <- UploadResult{Error: fmt.Errorf("URL não encontrada na resposta")}
			return
		}

		resultChan <- UploadResult{URL: url}
	}()

	// Aguarda o resultado
	result := <-resultChan
	if result.Error != nil {
		return "", result.Error
	}

	return result.URL, nil
}

// GetUploadProgress retorna o progresso dos uploads para um produto
func GetUploadProgress(productID uint) []UploadResult {
	if uploadService == nil {
		return nil
	}

	uploadService.mu.RLock()
	defer uploadService.mu.RUnlock()

	return uploadService.results[productID]
}
