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

/**
 * AddProductImage appends an image URL to the product's comma-separated list
 * and invalidates the corresponding image cache entry.
 *
 * @param productID - The target product ID
 * @param imagePath - The new image URL to append
 */
func AddProductImage(productID uint, imagePath string) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produto não encontrado")
	}

	if product.ImageUrls != "" {
		product.ImageUrls = product.ImageUrls + "," + imagePath
	} else {
		product.ImageUrls = imagePath
	}

	if err := repository.UpdateProduct(product); err != nil {
		return err
	}

	cacheService := &CacheService{}
	cacheService.InvalidateProductImageCache(productID)

	return nil
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

	cacheService := &CacheService{}

	if cached, found := cacheService.GetCachedPaginatedProducts(page, limit); found {
		return cached, nil
	}

	offset := (page - 1) * limit

	products, total, err := repository.GetPaginatedProductsWithCount(limit, offset)
	if err != nil {
		return nil, err
	}

	metadata := model.CalculatePagination(page, limit, total)

	response := &model.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}

	cacheService.SetCachedPaginatedProducts(response, page, limit)

	return response, nil
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

	cacheService := &CacheService{}

	if cached, found := cacheService.GetCachedSearchWithMetadata(searchTerm, page, limit); found {
		return cached, nil
	}

	offset := (page - 1) * limit

	products, total, err := repository.SearchProductsByNameWithCount(searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}

	metadata := model.CalculatePagination(page, limit, total)

	response := &model.PaginatedResponse{
		Data:     products,
		Metadata: metadata,
	}

	cacheService.SetCachedSearchWithMetadata(response, searchTerm, page, limit)

	return response, nil
}

/**
 * DeleteProductImage removes an image from the product by its index position
 * and invalidates the corresponding image cache entry.
 *
 * @param productID - The target product ID
 * @param index - Zero-based index of the image to remove
 */
func DeleteProductImage(productID uint, index int) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produto não encontrado")
	}

	imagePaths := repository.ParseImageUrls(product.ImageUrls)
	if index < 0 || index >= len(imagePaths) {
		return errors.New("índice de imagem inválido")
	}

	imagePaths = append(imagePaths[:index], imagePaths[index+1:]...)
	product.ImageUrls = repository.JoinImageUrls(imagePaths)

	if err := repository.UpdateProduct(product); err != nil {
		return err
	}

	cacheService := &CacheService{}
	cacheService.InvalidateProductImageCache(productID)

	return nil
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

/**
 * GetProductImages retrieves parsed image URLs for a product using cache-aside.
 * On cache hit, returns immediately without DB access.
 * On cache miss, fetches from DB, populates cache, and returns.
 *
 * @param productID - The product ID to fetch images for
 */
func GetProductImages(productID uint) ([]string, error) {
	cacheService := &CacheService{}

	if images, found := cacheService.GetCachedProductImages(productID); found {
		return images, nil
	}

	product, err := repository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produto não encontrado")
	}

	imagePaths := repository.ParseImageUrls(product.ImageUrls)

	cacheService.SetCachedProductImages(imagePaths, productID)

	return imagePaths, nil
}

/**
 * GetBatchProductImages retrieves images for multiple products in a single operation.
 * Uses cache-aside per product: cached products are served from memory,
 * cache misses are fetched in a single DB query using WHERE IN.
 * Time complexity: O(n) where n = number of requested product IDs.
 *
 * @param productIDs - Slice of product IDs to fetch images for
 */
func GetBatchProductImages(productIDs []uint) (map[uint][]string, error) {
	cacheService := &CacheService{}
	result := make(map[uint][]string, len(productIDs))
	var missIDs []uint

	for _, id := range productIDs {
		if images, found := cacheService.GetCachedProductImages(id); found {
			result[id] = images
		} else {
			missIDs = append(missIDs, id)
		}
	}

	if len(missIDs) == 0 {
		return result, nil
	}

	products, err := repository.GetProductsByIDs(missIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products batch: %w", err)
	}

	for i := range products {
		images := repository.ParseImageUrls(products[i].ImageUrls)
		result[products[i].ID] = images
		cacheService.SetCachedProductImages(images, products[i].ID)
	}

	return result, nil
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
