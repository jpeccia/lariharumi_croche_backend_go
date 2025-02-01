package service

import (
	"errors"
	"fmt"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
)

// CreateCategory cria uma nova categoria
func CreateCategory(category *model.Category) error {
	if category.Name == "" {
		return errors.New("nome da categoria é obrigatório")
	}

	err := repository.CreateCategory(category)
	if err != nil {
		return fmt.Errorf("erro ao criar categoria: %w", err)
	}

	return nil
}

// DeleteCategory deleta uma categoria pelo seu ID
func DeleteCategory(categoryID uint) error {
	// Opcional: verificar se a categoria possui produtos associados antes de deletar.
	err := repository.DeleteCategory(categoryID)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}
	return nil
}

// GetCategories retorna todas as categorias
func GetCategories() ([]model.Category, error) {
	categories, err := repository.GetCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryImage retorna a imagem da categoria
func GetCategoryImage(categoryID uint) (string, error) {
	category, err := repository.GetCategoryByID(categoryID)
	if err != nil {
		return "", err
	}
	if category == nil {
		return "", errors.New("categoria não encontrada")
	}

	return category.Image, nil
}
