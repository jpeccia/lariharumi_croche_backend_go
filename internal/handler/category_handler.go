package handler

import (
	"net/http"
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
