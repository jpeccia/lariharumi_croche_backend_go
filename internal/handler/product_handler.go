package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
)

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

func DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	if err := service.DeleteProduct(uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso!"})
}

type ImgBBResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Success bool `json:"success"`
}

func UploadProductImages(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar o formulário: " + err.Error()})
		return
	}

	files := form.File["images[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhuma imagem enviada"})
		return
	}

	// Faz upload assíncrono
	results, err := service.UploadProductImagesAsync(uint(productID), files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar upload: " + err.Error()})
		return
	}

	// Processa os resultados
	var uploadedUrls []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("Imagem %d: %v", result.Index, result.Error))
		} else {
			uploadedUrls = append(uploadedUrls, result.URL)
			// Salva a URL no banco de dados
			if err := service.AddProductImage(uint(productID), result.URL); err != nil {
				log.Printf("Erro ao salvar imagem %s no banco: %v", result.URL, err)
			}
		}
	}

	response := gin.H{
		"message": "Upload processado",
		"success": len(uploadedUrls),
		"failed":  len(errors),
		"urls":    uploadedUrls,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

func UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	existingProduct, err := repository.GetProductByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

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
		PriceRange:  req.Price,
		CategoryID:  req.CategoryID,
	}

	if req.Image != "" {
		product.ImageUrls = req.Image
	} else {
		product.ImageUrls = existingProduct.ImageUrls
	}

	if err := service.UpdateProduct(uint(productID), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

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

func SearchProducts(c *gin.Context) {
	searchTerm := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// Usa a nova função com metadados de paginação
	paginatedResponse, err := service.SearchProductsWithMetadata(searchTerm, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produtos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

func GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// Usa a nova função com metadados de paginação
	paginatedResponse, err := service.GetPaginatedProductsWithMetadata(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter produtos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

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

func GetProductImages(c *gin.Context) {
	productIDStr := c.Param("id")
	log.Printf("Recebendo requisição para imagens do produto ID: %s", productIDStr)

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		log.Println("ID do produto inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	images, err := service.GetProductImages(uint(productID))
	if err != nil {
		log.Printf("Erro ao obter imagens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter imagens do produto: " + err.Error()})
		return
	}

	log.Printf("Enviando imagens para o frontend: %v", images)
	c.JSON(http.StatusOK, images)
}

// GetUploadProgress retorna o progresso dos uploads para um produto
func GetUploadProgress(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	progress := service.GetUploadProgress(uint(productID))
	c.JSON(http.StatusOK, gin.H{
		"productId": productID,
		"progress":  progress,
	})
}
