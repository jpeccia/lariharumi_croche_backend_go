package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

// CreateCategory cria uma nova categoria (exige token de admin)
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := model.Category{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
	}

	if err := service.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar categoria: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// DeleteCategory deleta uma categoria (exige token de admin)
func DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	if err := service.DeleteCategory(uint(categoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar categoria: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoria deletada com sucesso!"})
}

// GetCategories retorna todas as categorias (público)
func GetCategories(c *gin.Context) {
	categories, err := service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter categorias: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategoryImage(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	imageURL, err := service.GetCategoryImage(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter imagem da categoria: " + err.Error()})
		return
	}

	if imageURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Imagem não encontrada para essa categoria"})
		return
	}

	// Retorna a URL pública da imagem no ImgBB
	c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
}


func UploadCategoryImage(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagem não encontrada: " + err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir a imagem: " + err.Error()})
		return
	}
	defer src.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar formulário para upload: " + err.Error()})
		return
	}

	if _, err := io.Copy(part, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao copiar imagem: " + err.Error()})
		return
	}
	writer.Close()

	// Obtém a chave da API do ImgBB do .env
	imgBBKey := os.Getenv("IMGBB_API_KEY")
	if imgBBKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chave da API ImgBB não encontrada"})
		return
	}

	imgBBURL := fmt.Sprintf("https://api.imgbb.com/1/upload?key=%s", imgBBKey)
	req, err := http.NewRequest("POST", imgBBURL, &buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar requisição para ImgBB: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar imagem para ImgBB: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Decodifica a resposta do ImgBB
	var imgBBResp ImgBBResponse
	if err := json.NewDecoder(resp.Body).Decode(&imgBBResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar resposta do ImgBB: " + err.Error()})
		return
	}

	// Verifica se o upload foi bem-sucedido
	if !imgBBResp.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload para ImgBB falhou"})
		return
	}

	// Obtém a URL da imagem no ImgBB
	imageURL := imgBBResp.Data.URL

	// Atualiza a categoria no banco de dados com a URL da imagem
	if err := service.AddCategoryImage(uint(categoryID), imageURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria com imagem: " + err.Error()})
		return
	}

	// Retorna a URL da imagem para o frontend
	c.JSON(http.StatusOK, gin.H{
		"message":  "Imagem enviada com sucesso!",
		"imageUrl": imageURL,
	})
}

// DeleteCategoryImage deleta a imagem de uma categoria (exige token de admin)
func DeleteCategoryImage(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	// Obtém a imagem da categoria do banco de dados
	imagePath, err := service.GetCategoryImage(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar imagem da categoria: " + err.Error()})
		return
	}

	// Verifica se a imagem existe
	if imagePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Imagem da categoria não encontrada"})
		return
	}

	// Remove a imagem do sistema de arquivos
	localPath := "./uploads" + imagePath
	if err := os.Remove(localPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover a imagem do sistema de arquivos: " + err.Error()})
		return
	}

	// Remove o caminho da imagem do banco de dados
	if err := service.RemoveCategoryImage(uint(categoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover imagem do banco de dados: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagem da categoria deletada com sucesso!"})
}

// UpdateCategory atualiza os dados de uma categoria existente (exige token de admin)
func UpdateCategory(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	// Define a estrutura de dados para atualizar a categoria
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	// Valida os dados recebidos
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cria a categoria com os dados recebidos
	updatedCategory := model.Category{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
	}

	// Chama o serviço para atualizar a categoria no banco de dados
	if err := service.UpdateCategory(uint(categoryID), &updatedCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria: " + err.Error()})
		return
	}

	// Retorna a categoria atualizada
	c.JSON(http.StatusOK, updatedCategory)
}
