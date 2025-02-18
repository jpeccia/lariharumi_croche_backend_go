package handler

import (
	"fmt"
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

// GetCategoryImage retorna a imagem de uma categoria (público)
func GetCategoryImage(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	// Conecte-se ao Supabase Storage
	client := getSupabaseClient()
	bucketName := "category-images"

	// Defina o caminho do arquivo para a categoria
	imagePath := fmt.Sprintf("categories/category_%d_", categoryID)

	// Verifique se o arquivo existe no Supabase Storage
	object, err := client.GetObject(bucketName, imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar a imagem no Supabase: " + err.Error()})
		return
	}

	// Construa a URL pública da imagem
	imageURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, object.Key)

	// Retorna a URL pública da imagem
	c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
}

// UploadCategoryImage realiza o upload de uma imagem para uma categoria
func UploadCategoryImage(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	// Obtém o arquivo da imagem
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagem não encontrada: " + err.Error()})
		return
	}

	// Define o nome do arquivo e os caminhos para salvar e retornar
	newFileName := fmt.Sprintf("category_%d_%s", categoryID, file.Filename)
	uploadPath := fmt.Sprintf("categories/%s", newFileName)

	// Abre o arquivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo"})
		return
	}
	defer src.Close()

	// Obtém o cliente Supabase
	client := getSupabaseClient()
	bucketName := "category-images"

	// Faz o upload para o Supabase Storage
	_, err = client.UploadFile(bucketName, uploadPath, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar para o Supabase: " + err.Error()})
		return
	}

	// Construa a URL pública da imagem
	imageURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, uploadPath)

	// Retorna a URL da imagem enviada
	c.JSON(http.StatusOK, gin.H{"message": "Imagem da categoria enviada!", "imageUrl": imageURL})
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
