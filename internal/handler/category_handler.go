package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	image, err := service.GetCategoryImage(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter imagem da categoria: " + err.Error()})
		return
	}

	// Aqui você gera a URL completa para a imagem
	imageURL := fmt.Sprintf("http://localhost:8080%s", image)
	c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
}

// UploadCategoryImage realiza o upload de uma imagem para uma categoria
func UploadCategoryImage(c *gin.Context) {
	// Obtém o id da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	// Obtém o arquivo do formulário multipart
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagem não encontrada: " + err.Error()})
		return
	}

	// Define o nome do arquivo e os caminhos para salvar e retornar
	filename := filepath.Base(file.Filename)
	localPath := "./uploads/categories/" + filename   // Caminho para salvar no sistema de arquivos
	relativePath := "/uploads/categories/" + filename // Caminho relativo para a URL

	// Salva o arquivo no diretório local
	if err := c.SaveUploadedFile(file, localPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar a imagem: " + err.Error()})
		return
	}

	// Atualiza a categoria com a nova imagem (salva o caminho relativo no banco de dados, por exemplo)
	if err := service.AddCategoryImage(uint(categoryID), relativePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria com imagem: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Imagem enviada com sucesso!",
		"imageUrl": "http://localhost:8080/uploads/categories/" + filename, // Remover o caminho redundante
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