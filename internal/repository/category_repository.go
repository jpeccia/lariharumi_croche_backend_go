package repository

import (
	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"gorm.io/gorm"
)

// CreateCategory cria uma nova categoria no banco de dados
func CreateCategory(category *model.Category) error {
	if err := config.DB.Create(category).Error; err != nil {
		return err
	}
	return nil
}

// GetCategories retorna todas as categorias
func GetCategories() ([]model.Category, error) {
	var categories []model.Category
	if err := config.DB.Preload("Products").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retorna uma categoria pelo seu ID
func GetCategoryByID(categoryID uint) (*model.Category, error) {
	var category model.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deleta uma categoria pelo seu ID (soft delete)
func DeleteCategory(categoryID uint) error {
	if err := config.DB.Delete(&model.Category{}, categoryID).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteCategory deleta permanentemente uma categoria (apenas para admin)
func HardDeleteCategory(categoryID uint) error {
	if err := config.DB.Unscoped().Delete(&model.Category{}, categoryID).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCategory atualiza uma categoria no banco de dados
func UpdateCategory(category *model.Category) error {
	if err := config.DB.Save(category).Error; err != nil {
		return err
	}
	return nil
}
