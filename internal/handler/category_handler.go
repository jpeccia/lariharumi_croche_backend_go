package handler

import (
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{"image": image})
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

	// Define um caminho para salvar a imagem (ajuste conforme sua necessidade)
	filename := filepath.Base(file.Filename)
	uploadPath := "./uploads/categories/" + filename

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar a imagem: " + err.Error()})
		return
	}

	// Atualiza a categoria com a nova imagem. Aqui, você pode chamar uma função no service.
	if err := service.AddCategoryImage(uint(categoryID), uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar categoria com imagem: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagem enviada com sucesso!", "path": uploadPath})
}
