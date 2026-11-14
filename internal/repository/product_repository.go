package repository

import (
	"strings"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"gorm.io/gorm"
)

// CreateProduct cria um novo produto no banco de dados
func CreateProduct(product *model.Product) error {
	if err := config.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// GetProductByID retorna um produto pelo seu ID
func GetProductByID(productID uint) (*model.Product, error) {
	var product model.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetProducts retorna todos os produtos
func GetProducts() ([]model.Product, error) {
	var products []model.Product
	if err := config.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func GetPaginatedProducts(limit int, offset int) ([]model.Product, error) {
	var products []model.Product

	err := config.DB.Joins("Category").
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetPaginatedProductsWithCount retorna produtos paginados com contagem total
func GetPaginatedProductsWithCount(limit int, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Conta o total de produtos
	if err := config.DB.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Busca os produtos com preload
	err := config.DB.Joins("Category").
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func SearchProductsByName(searchTerm string, limit, offset int) ([]model.Product, error) {
	var products []model.Product

	query := config.DB.Joins("Category").Where("products.name ILIKE ?", "%"+searchTerm+"%")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

// SearchProductsByNameWithCount retorna produtos pesquisados com contagem total
func SearchProductsByNameWithCount(searchTerm string, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	// Conta o total de produtos que correspondem à pesquisa
	if err := config.DB.Model(&model.Product{}).Where("products.name ILIKE ?", "%"+searchTerm+"%").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Busca os produtos com preload
	query := config.DB.Joins("Category").Where("products.name ILIKE ?", "%"+searchTerm+"%")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, total, err
}

func GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	var products []model.Product
	if err := config.DB.Joins("Category").Where("products.category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsByCategoryWithCount retorna produtos de uma categoria paginados,
// com a contagem total, ordenados por nome (usa o índice composto cat+nome).
func GetProductsByCategoryWithCount(categoryID uint, limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := config.DB.Model(&model.Product{}).
		Where("products.category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := config.DB.Joins("Category").
		Where("products.category_id = ?", categoryID).
		Order("products.name ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func UpdateProduct(product *model.Product) error {
	if err := config.DB.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func ParseImageUrls(imageUrls string) []string {
	if imageUrls == "" {
		return []string{}
	}
	return strings.Split(imageUrls, ",")
}

func JoinImageUrls(imagePaths []string) string {
	return strings.Join(imagePaths, ",")
}

// DeleteProduct deleta um produto pelo seu ID (soft delete)
func DeleteProduct(productID uint) error {
	var product model.Product
	if err := config.DB.Delete(&product, productID).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteProduct deleta permanentemente um produto (apenas para admin)
func HardDeleteProduct(productID uint) error {
	var product model.Product
	if err := config.DB.Unscoped().Delete(&product, productID).Error; err != nil {
		return err
	}
	return nil
}

/**
 * GetProductsByIDs fetches multiple products in a single query using WHERE IN.
 * This avoids N+1 queries when loading images for a batch of products.
 *
 * @param productIDs - Slice of product IDs to fetch
 */
func GetProductsByIDs(productIDs []uint) ([]model.Product, error) {
	var products []model.Product
	if err := config.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
