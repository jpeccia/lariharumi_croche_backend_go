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

	err := config.DB.Preload("Category").
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func SearchProductsByName(searchTerm string, limit, offset int) ([]model.Product, error) {
	var products []model.Product

	query := config.DB.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(searchTerm)+"%")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

func GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	var products []model.Product
	if err := config.DB.Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func UpdateProduct(product *model.Product) error {
	if err := config.DB.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func ParseImageUrls(imageUrls string) []string {
	return strings.Split(imageUrls, ",")
}

func JoinImageUrls(imagePaths []string) string {
	return strings.Join(imagePaths, ",")
}

func DeleteProduct(productID uint) error {
	var product model.Product
	if err := config.DB.Delete(&product, productID).Error; err != nil {
		return err
	}
	return nil
}
