package service

import (
	"errors"
	"fmt"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
)

// CreateProduct cria um novo produto
func CreateProduct(product *model.Product) error {
	// Você pode adicionar validações, como verificar se a categoria existe.
	if product.CategoryID == 0 {
		return errors.New("categoria inválida")
	}

	// Chama o repositório para criar o produto.
	err := repository.CreateProduct(product)
	if err != nil {
		return fmt.Errorf("erro ao criar produto: %w", err)
	}

	// Invalida cache relacionado a produtos
	cacheService := &CacheService{}
	cacheService.InvalidateProductCache()

	return nil
}

// AddProductImage adiciona um caminho de imagem ao produto
func AddProductImage(productID uint, imagePath string) error {
	// Buscar o produto pelo ID (opcional, para verificar se existe)
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produto não encontrado")
	}

	// Aqui supomos que o campo ImageUrls é uma string que armazena os caminhos separados por vírgula.
	// Em uma implementação real, pode ser um array ou uma tabela associada.
	if product.ImageUrls != "" {
		product.ImageUrls = product.ImageUrls + "," + imagePath
	} else {
		product.ImageUrls = imagePath
	}

	// Atualiza o produto com o novo caminho de imagem.
	return repository.UpdateProduct(product)
}

func GetPaginatedProducts(limit int, offset int) ([]model.Product, error) {
	cacheService := &CacheService{}

	// Tenta buscar no cache primeiro
	if products, found := cacheService.GetCachedProducts(limit, offset); found {
		return products, nil
	}

	// Se não encontrou no cache, busca no banco
	products, err := repository.GetPaginatedProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	cacheService.SetCachedProducts(products, limit, offset)

	return products, nil
}

// GetPaginatedProductsWithMetadata retorna produtos paginados com metadados
func GetPaginatedProductsWithMetadata(page, limit int) (*model.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Busca produtos e contagem total
	products, total, err := repository.GetPaginatedProductsWithCount(limit, offset)
	if err != nil {
		return nil, err
	}

	// Calcula metadados de paginação
	metadata := model.CalculatePagination(page, limit, total)

	return &model.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}, nil
}

func SearchProducts(searchTerm string, limit, offset int) ([]model.Product, error) {
	cacheService := &CacheService{}

	// Tenta buscar no cache primeiro
	if products, found := cacheService.GetCachedSearchProducts(searchTerm, limit, offset); found {
		return products, nil
	}

	// Se não encontrou no cache, busca no banco
	products, err := repository.SearchProductsByName(searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	cacheService.SetCachedSearchProducts(products, searchTerm, limit, offset)

	return products, nil
}

// SearchProductsWithMetadata retorna produtos pesquisados com metadados de paginação
func SearchProductsWithMetadata(searchTerm string, page, limit int) (*model.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Busca produtos e contagem total
	products, total, err := repository.SearchProductsByNameWithCount(searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}

	// Calcula metadados de paginação
	metadata := model.CalculatePagination(page, limit, total)

	return &model.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}, nil
}

// DeleteProductImage remove uma imagem do produto dado seu índice (posição na lista)
func DeleteProductImage(productID uint, index int) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produto não encontrado")
	}

	// Supondo que as imagens estejam armazenadas como string separada por vírgula.
	imagePaths := repository.ParseImageUrls(product.ImageUrls)
	if index < 0 || index >= len(imagePaths) {
		return errors.New("índice de imagem inválido")
	}

	// Remove a imagem do slice
	imagePaths = append(imagePaths[:index], imagePaths[index+1:]...)
	// Atualiza o campo ImageUrls
	product.ImageUrls = repository.JoinImageUrls(imagePaths)

	return repository.UpdateProduct(product)
}

// GetProducts retorna todos os produtos
func GetProducts() ([]model.Product, error) {
	products, err := repository.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsByCategory retorna os produtos filtrados por categoria
func GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	cacheService := &CacheService{}

	// Tenta buscar no cache primeiro
	if products, found := cacheService.GetCachedProductsByCategory(categoryID); found {
		return products, nil
	}

	// Se não encontrou no cache, busca no banco
	products, err := repository.GetProductsByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	cacheService.SetCachedProductsByCategory(products, categoryID)

	return products, nil
}

// GetProductImages retorna as imagens de um produto
// Se o campo ImageUrls for uma string separada por vírgula, podemos utilizar essa lógica para transformá-la em slice.
func GetProductImages(productID uint) ([]string, error) {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produto não encontrado")
	}
	imagePaths := repository.ParseImageUrls(product.ImageUrls)
	return imagePaths, nil
}

// DeleteProduct deleta um produto do banco de dados
func DeleteProduct(productID uint) error {
	// Tenta encontrar o produto
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}

	// Verifica se o produto existe
	if product == nil {
		return errors.New("produto não encontrado")
	}

	// Deleta o produto
	if err := repository.DeleteProduct(productID); err != nil {
		return err
	}

	// Invalida cache relacionado a produtos
	cacheService := &CacheService{}
	cacheService.InvalidateProductCache()

	return nil
}

// UpdateProduct no serviço
func UpdateProduct(productID uint, updatedProduct *model.Product) error {
	// Primeiro, tenta buscar o produto pelo ID
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err // Retorna erro se não encontrar o produto
	}

	// Atualiza os campos do produto
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.ImageUrls = updatedProduct.ImageUrls
	product.PriceRange = updatedProduct.PriceRange
	product.CategoryID = updatedProduct.CategoryID

	// Atualiza o produto no banco de dados
	if err := repository.UpdateProduct(product); err != nil {
		return err // Retorna erro se falhar ao atualizar no banco
	}

	// Invalida cache relacionado a produtos
	cacheService := &CacheService{}
	cacheService.InvalidateProductCache()

	return nil
}
