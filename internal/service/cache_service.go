package service

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/cache"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

const (
	// TTL padrão para cache
	DefaultTTL = 15 * time.Minute
	// TTL para produtos (mais longo pois mudam menos)
	ProductTTL = 30 * time.Minute
	// TTL para categorias (mais longo pois mudam menos)
	CategoryTTL = 1 * time.Hour
)

// CacheService gerencia o cache de produtos e categorias
type CacheService struct{}

// GetCachedProducts busca produtos no cache
func (cs *CacheService) GetCachedProducts(limit, offset int) ([]model.Product, bool) {
	key := cache.GenerateKey("products", limit, offset)
	var products []model.Product

	err := cache.Get(key, &products)
	if err == redis.Nil {
		return nil, false // Cache miss
	}
	if err != nil {
		return nil, false // Erro no cache, não usar
	}

	return products, true // Cache hit
}

// SetCachedProducts armazena produtos no cache
func (cs *CacheService) SetCachedProducts(products []model.Product, limit, offset int) {
	key := cache.GenerateKey("products", limit, offset)
	cache.Set(key, products, ProductTTL)
}

// GetCachedProductsByCategory busca produtos por categoria no cache
func (cs *CacheService) GetCachedProductsByCategory(categoryID uint) ([]model.Product, bool) {
	key := cache.GenerateKey("products_category", categoryID)
	var products []model.Product

	err := cache.Get(key, &products)
	if err == redis.Nil {
		return nil, false // Cache miss
	}
	if err != nil {
		return nil, false // Erro no cache, não usar
	}

	return products, true // Cache hit
}

// SetCachedProductsByCategory armazena produtos por categoria no cache
func (cs *CacheService) SetCachedProductsByCategory(products []model.Product, categoryID uint) {
	key := cache.GenerateKey("products_category", categoryID)
	cache.Set(key, products, ProductTTL)
}

// GetCachedSearchProducts busca produtos pesquisados no cache
func (cs *CacheService) GetCachedSearchProducts(searchTerm string, limit, offset int) ([]model.Product, bool) {
	key := cache.GenerateKey("products_search", searchTerm, limit, offset)
	var products []model.Product

	err := cache.Get(key, &products)
	if err == redis.Nil {
		return nil, false // Cache miss
	}
	if err != nil {
		return nil, false // Erro no cache, não usar
	}

	return products, true // Cache hit
}

// SetCachedSearchProducts armazena produtos pesquisados no cache
func (cs *CacheService) SetCachedSearchProducts(products []model.Product, searchTerm string, limit, offset int) {
	key := cache.GenerateKey("products_search", searchTerm, limit, offset)
	cache.Set(key, products, ProductTTL)
}

// GetCachedCategories busca categorias no cache
func (cs *CacheService) GetCachedCategories() ([]model.Category, bool) {
	key := "categories"
	var categories []model.Category

	err := cache.Get(key, &categories)
	if err == redis.Nil {
		return nil, false // Cache miss
	}
	if err != nil {
		return nil, false // Erro no cache, não usar
	}

	return categories, true // Cache hit
}

// SetCachedCategories armazena categorias no cache
func (cs *CacheService) SetCachedCategories(categories []model.Category) {
	key := "categories"
	cache.Set(key, categories, CategoryTTL)
}

// InvalidateProductCache invalida cache relacionado a produtos
func (cs *CacheService) InvalidateProductCache() {
	// Remove todos os caches de produtos
	cache.DeletePattern("products*")
	cache.DeletePattern("products_category*")
	cache.DeletePattern("products_search*")
}

// InvalidateCategoryCache invalida cache relacionado a categorias
func (cs *CacheService) InvalidateCategoryCache() {
	// Remove cache de categorias
	cache.Delete("categories")
	// Também remove produtos por categoria pois podem ter mudado
	cache.DeletePattern("products_category*")
}

// InvalidateAllCache invalida todo o cache
func (cs *CacheService) InvalidateAllCache() {
	cache.DeletePattern("products*")
	cache.DeletePattern("categories*")
}
