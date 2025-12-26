package service

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gen2brain/webp"
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

/**
 * WebPQuality defines the encoding quality for WebP conversion.
 * Range: 0-100, where 80 offers a good balance between size and quality.
 */
const WebPQuality = 80

/**
 * convertToWebP reads an image from the provided reader, decodes it,
 * and encodes it to WebP format with the configured quality setting.
 * Supports JPEG, PNG, and GIF input formats.
 *
 * @param reader - The source image data reader
 * @param originalName - Original filename to derive the WebP filename
 * @returns - WebP encoded bytes, new filename, and any error
 */
func convertToWebP(reader io.Reader, originalName string) (*bytes.Reader, string, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	opts := webp.Options{
		Quality: WebPQuality,
	}
	if err := webp.Encode(&buf, img, opts); err != nil {
		return nil, "", fmt.Errorf("failed to encode to WebP: %w", err)
	}

	baseName := strings.TrimSuffix(originalName, "."+getExtension(originalName))
	webpName := baseName + ".webp"

	return bytes.NewReader(buf.Bytes()), webpName, nil
}

/**
 * getExtension extracts the file extension from a filename.
 *
 * @param filename - The filename to extract extension from
 * @returns - The file extension without the leading dot
 */
func getExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

/**
 * InitUploadService initializes the upload service with a worker pool.
 *
 * @param maxWorkers - Maximum number of concurrent upload workers
 */
func InitUploadService(maxWorkers int) {
	uploadService = &UploadService{
		workerPool: make(chan struct{}, maxWorkers),
		jobQueue:   make(chan UploadJob, 100),
		results:    make(map[uint][]UploadResult),
	}

	for i := 0; i < maxWorkers; i++ {
		go uploadService.worker()
	}
}

/**
 * worker processes upload jobs from the job queue.
 */
func (us *UploadService) worker() {
	for job := range us.jobQueue {
		us.workerPool <- struct{}{}

		go func(j UploadJob) {
			defer func() { <-us.workerPool }()
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

/**
 * uploadSingleFile converts an image to WebP format and uploads it to ImgBB.
 * The conversion reduces file size significantly while maintaining quality.
 *
 * @param file - The multipart file header from the HTTP request
 * @param index - The index of the file in the upload batch
 * @returns - UploadResult containing the URL or error
 */
func (us *UploadService) uploadSingleFile(file *multipart.FileHeader, index int) UploadResult {
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("IMGBB_API_KEY environment variable not found"),
		}
	}

	src, err := file.Open()
	if err != nil {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("failed to open file: %w", err),
		}
	}
	defer src.Close()

	webpReader, webpFilename, err := convertToWebP(src, file.Filename)
	if err != nil {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("failed to convert to WebP: %w", err),
		}
	}

	resp, err := sharedRestyClient.R().
		SetFileReader("image", webpFilename, webpReader).
		SetFormData(map[string]string{"key": apiKey}).
		Post("https://api.imgbb.com/1/upload")

	if err != nil {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("upload failed: %w", err),
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("upload failed with status %d", resp.StatusCode()),
		}
	}

	url := gjson.Get(resp.String(), "data.url").String()
	if url == "" {
		return UploadResult{
			Index: index,
			Error: fmt.Errorf("URL not found in response"),
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

/**
 * UploadCategoryImageAsync uploads a single category image asynchronously,
 * converting it to WebP format before upload.
 *
 * @param categoryID - The ID of the category (currently unused but kept for interface consistency)
 * @param file - The multipart file header from the HTTP request
 * @returns - The uploaded image URL and any error
 */
func UploadCategoryImageAsync(categoryID uint, file *multipart.FileHeader) (string, error) {
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("IMGBB_API_KEY environment variable not found")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	webpReader, webpFilename, err := convertToWebP(src, file.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to convert to WebP: %w", err)
	}

	resultChan := make(chan UploadResult, 1)

	go func() {
		defer close(resultChan)

		resp, err := sharedRestyClient.R().
			SetFileReader("image", webpFilename, webpReader).
			SetFormData(map[string]string{"key": apiKey}).
			Post("https://api.imgbb.com/1/upload")

		if err != nil {
			resultChan <- UploadResult{Error: fmt.Errorf("upload failed: %w", err)}
			return
		}

		if resp.StatusCode() != http.StatusOK {
			resultChan <- UploadResult{Error: fmt.Errorf("upload failed with status %d", resp.StatusCode())}
			return
		}

		url := gjson.Get(resp.String(), "data.url").String()
		if url == "" {
			resultChan <- UploadResult{Error: fmt.Errorf("URL not found in response")}
			return
		}

		resultChan <- UploadResult{URL: url}
	}()

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
