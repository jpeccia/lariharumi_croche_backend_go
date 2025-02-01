package handler

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

// CreateProduct cria um novo produto (exige token de admin)
func CreateProduct(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Price       string `json:"price"`
		CategoryID  uint   `json:"categoryId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		ImageUrls:   req.Image,
		PriceRange:  req.Price,
		CategoryID:  req.CategoryID,
	}

	if err := service.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// DeleteProduct deleta um produto pelo ID (exige token de admin)
func DeleteProduct(c *gin.Context) {
	// Obtém o id do produto da URL
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	// Chama o serviço para deletar o produto
	if err := service.DeleteProduct(uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso!"})
}


// UploadProductImage realiza o upload de uma imagem para um produto
func UploadProductImage(c *gin.Context) {
	// Obtém o id do produto da URL
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
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
	uploadPath := "./uploads/products/" + filename

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar a imagem: " + err.Error()})
		return
	}

	// Atualiza o produto com a nova imagem. Aqui, você pode chamar uma função no service.
	if err := service.AddProductImage(uint(productID), uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produto com imagem: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagem enviada com sucesso!", "path": uploadPath})
}

// DeleteProductImage remove uma imagem do produto. O endpoint é definido como /products/:id/images/:index
func DeleteProductImage(c *gin.Context) {
	productIDStr := c.Param("id")
	indexStr := c.Param("index")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index inválido"})
		return
	}

	if err := service.DeleteProductImage(uint(productID), index); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar imagem: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagem deletada com sucesso!"})
}

// GetProducts retorna todos os produtos (público)
func GetProducts(c *gin.Context) {
	products, err := service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter produtos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductsByCategory retorna os produtos filtrados por categoria (público)
func GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
		return
	}

	products, err := service.GetProductsByCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter produtos da categoria: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductImages retorna as imagens de um produto (público)
func GetProductImages(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	images, err := service.GetProductImages(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter imagens do produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

