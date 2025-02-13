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
func GetProducts(limit, offset int) ([]model.Product, error) {
	var products []model.Product
	err := config.DB.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsByCategory retorna produtos filtrados por categoria
func GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	var products []model.Product
	if err := config.DB.Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct atualiza as informações de um produto
func UpdateProduct(product *model.Product) error {
	if err := config.DB.Save(product).Error; err != nil {
		return err
	}
	return nil
}

// ParseImageUrls converte o campo de imagens (string separada por vírgula) em um slice de strings
func ParseImageUrls(imageUrls string) []string {
	return strings.Split(imageUrls, ",")
}

// JoinImageUrls junta um slice de imagens em uma string separada por vírgula
func JoinImageUrls(imagePaths []string) string {
	return strings.Join(imagePaths, ",")
}

// DeleteProduct deleta um produto pelo ID
func DeleteProduct(productID uint) error {
	var product model.Product
	if err := config.DB.Delete(&product, productID).Error; err != nil {
		return err
	}
	return nil
}
