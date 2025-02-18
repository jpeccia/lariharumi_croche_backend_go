package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/service"
	storage "github.com/supabase-community/storage-go"
)

func getSupabaseClient() *storage.Client {
	supabaseURL := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	client := storage.NewClient(supabaseURL+"/storage/v1", apiKey, nil)
	return client
}

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
func UploadProductImages(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	// Processa o formulário multipart
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar o formulário"})
		return
	}

	// Verifica se há imagens enviadas
	files := form.File["images[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhuma imagem enviada"})
		return
	}

	// Configura o cliente do Supabase
	client := getSupabaseClient()
	bucketName := "product-images"

	var uploadedPaths []string
	for _, file := range files {
		// Verifica o tipo de arquivo (apenas imagens)
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não é uma imagem"})
			return
		}

		// Cria um nome único para o arquivo
		newFileName := fmt.Sprintf("product_%d_%s", productID, file.Filename)
		uploadPath := fmt.Sprintf("products/%s", newFileName)

		// Abre o arquivo para o upload
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo"})
			return
		}
		defer src.Close()

		// Faz o upload para o Supabase
		_, err = client.UploadFile(bucketName, uploadPath, src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar para o Supabase: " + err.Error()})
			return
		}

		// Gera a URL pública da imagem
		imageURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", os.Getenv("SUPABASE_URL"), bucketName, uploadPath)
		uploadedPaths = append(uploadedPaths, imageURL)
	}

	// Retorna as URLs das imagens enviadas
	c.JSON(http.StatusOK, gin.H{"message": "Imagens enviadas!", "paths": uploadedPaths})
}

// UpdateProduct atualiza as informações de um produto existente (exige token de admin)
func UpdateProduct(c *gin.Context) {
	// Obtém o ID do produto a ser atualizado
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	// Busca o produto existente no banco de dados
	existingProduct, err := repository.GetProductByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	// Recebe os dados atualizados para o produto
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

	// Atualiza os campos necessários, mantendo a imagem anterior se o campo estiver vazio
	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		PriceRange:  req.Price,
		CategoryID:  req.CategoryID,
	}

	// Se `image` for vazio, mantém a imagem anterior
	if req.Image != "" {
		product.ImageUrls = req.Image
	} else {
		product.ImageUrls = existingProduct.ImageUrls
	}

	// Chama o serviço de atualização
	if err := service.UpdateProduct(uint(productID), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produto: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProductImage(c *gin.Context) {
	productIDStr := c.Param("id")
	indexStr := c.Param("index")

	// Conversão do ID do produto para uint
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
		return
	}

	// Conversão do índice para o tipo int
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index inválido"})
		return
	}

	// Chama a função de exclusão com o índice
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

	// Busca imagens no banco de dados
	images, err := service.GetProductImages(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter imagens do produto"})
		return
	}

	c.JSON(http.StatusOK, images)
}
