package service

import (
	"time"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/cache"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

const (
	DefaultTTL  = 15 * time.Minute
	ProductTTL  = 30 * time.Minute
	CategoryTTL = 1 * time.Hour
)

/**
 * CacheService manages in-memory caching for products and categories.
 */
type CacheService struct{}

/**
 * GetCachedProducts retrieves paginated products from cache.
 * Returns the products slice and a boolean indicating cache hit.
 */
func (cs *CacheService) GetCachedProducts(limit, offset int) ([]model.Product, bool) {
	key := cache.GenerateKey("products", limit, offset)
	var products []model.Product

	err := cache.Get(key, &products)
	if err != nil {
		return nil, false
	}

	return products, true
}

/**
 * SetCachedProducts stores paginated products in cache with ProductTTL.
 */
func (cs *CacheService) SetCachedProducts(products []model.Product, limit, offset int) {
	key := cache.GenerateKey("products", limit, offset)
	cache.Set(key, products, ProductTTL)
}

/**
 * GetCachedProductsByCategory retrieves category-filtered products from cache.
 */
func (cs *CacheService) GetCachedProductsByCategory(categoryID uint) ([]model.Product, bool) {
	key := cache.GenerateKey("products_category", categoryID)
	var products []model.Product

	err := cache.Get(key, &products)
	if err != nil {
		return nil, false
	}

	return products, true
}

/**
 * SetCachedProductsByCategory stores category-filtered products in cache.
 */
func (cs *CacheService) SetCachedProductsByCategory(products []model.Product, categoryID uint) {
	key := cache.GenerateKey("products_category", categoryID)
	cache.Set(key, products, ProductTTL)
}

/**
 * GetCachedSearchProducts retrieves search results from cache.
 */
func (cs *CacheService) GetCachedSearchProducts(searchTerm string, limit, offset int) ([]model.Product, bool) {
	key := cache.GenerateKey("products_search", searchTerm, limit, offset)
	var products []model.Product

	err := cache.Get(key, &products)
	if err != nil {
		return nil, false
	}

	return products, true
}

/**
 * SetCachedSearchProducts stores search results in cache.
 */
func (cs *CacheService) SetCachedSearchProducts(products []model.Product, searchTerm string, limit, offset int) {
	key := cache.GenerateKey("products_search", searchTerm, limit, offset)
	cache.Set(key, products, ProductTTL)
}

/**
 * GetCachedCategories retrieves all categories from cache.
 */
func (cs *CacheService) GetCachedCategories() ([]model.Category, bool) {
	key := "categories"
	var categories []model.Category

	err := cache.Get(key, &categories)
	if err != nil {
		return nil, false
	}

	return categories, true
}

/**
 * SetCachedCategories stores all categories in cache with CategoryTTL.
 */
func (cs *CacheService) SetCachedCategories(categories []model.Category) {
	key := "categories"
	cache.Set(key, categories, CategoryTTL)
}

/**
 * GetCachedPaginatedProducts retrieves a full paginated response from cache.
 */
func (cs *CacheService) GetCachedPaginatedProducts(page, limit int) (*model.PaginatedResponse, bool) {
	key := cache.GenerateKey("paginated_products", page, limit)
	var response model.PaginatedResponse

	err := cache.Get(key, &response)
	if err != nil {
		return nil, false
	}

	return &response, true
}

/**
 * SetCachedPaginatedProducts stores a full paginated response in cache.
 */
func (cs *CacheService) SetCachedPaginatedProducts(response *model.PaginatedResponse, page, limit int) {
	key := cache.GenerateKey("paginated_products", page, limit)
	cache.Set(key, response, ProductTTL)
}

/**
 * GetCachedSearchWithMetadata retrieves a paginated search response from cache.
 */
func (cs *CacheService) GetCachedSearchWithMetadata(searchTerm string, page, limit int) (*model.PaginatedResponse, bool) {
	key := cache.GenerateKey("search_products", searchTerm, page, limit)
	var response model.PaginatedResponse

	err := cache.Get(key, &response)
	if err != nil {
		return nil, false
	}

	return &response, true
}

/**
 * SetCachedSearchWithMetadata stores a paginated search response in cache.
 */
func (cs *CacheService) SetCachedSearchWithMetadata(response *model.PaginatedResponse, searchTerm string, page, limit int) {
	key := cache.GenerateKey("search_products", searchTerm, page, limit)
	cache.Set(key, response, ProductTTL)
}

/**
 * GetCachedProductImages retrieves parsed image URLs for a product from cache.
 * Returns the image slice and a boolean indicating cache hit.
 *
 * @param productID - The product identifier
 */
func (cs *CacheService) GetCachedProductImages(productID uint) ([]string, bool) {
	key := cache.GenerateKey("product_images", productID)
	var images []string

	err := cache.Get(key, &images)
	if err != nil {
		return nil, false
	}

	return images, true
}

/**
 * SetCachedProductImages stores parsed image URLs for a product in cache.
 *
 * @param images - The parsed image URL slice
 * @param productID - The product identifier
 */
func (cs *CacheService) SetCachedProductImages(images []string, productID uint) {
	key := cache.GenerateKey("product_images", productID)
	cache.Set(key, images, ProductTTL)
}

/**
 * InvalidateProductImageCache removes the image cache for a specific product.
 *
 * @param productID - The product identifier whose image cache should be cleared
 */
func (cs *CacheService) InvalidateProductImageCache(productID uint) {
	key := cache.GenerateKey("product_images", productID)
	cache.Delete(key)
}

/**
 * InvalidateProductCache removes all product-related cache entries.
 */
func (cs *CacheService) InvalidateProductCache() {
	cache.DeletePattern("products*")
	cache.DeletePattern("products_category*")
	cache.DeletePattern("products_search*")
	cache.DeletePattern("paginated_products*")
	cache.DeletePattern("search_products*")
	cache.DeletePattern("product_images*")
}

/**
 * InvalidateCategoryCache removes category cache and related product caches.
 */
func (cs *CacheService) InvalidateCategoryCache() {
	cache.Delete("categories")
	cache.DeletePattern("products_category*")
}

/**
 * InvalidateAllCache removes every cached entry.
 */
func (cs *CacheService) InvalidateAllCache() {
	cache.DeletePattern("products*")
	cache.DeletePattern("categories*")
	cache.DeletePattern("paginated_products*")
	cache.DeletePattern("search_products*")
	cache.DeletePattern("product_images*")
}
