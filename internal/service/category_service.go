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

	// Invalida cache relacionado a categorias
	cacheService := &CacheService{}
	cacheService.InvalidateCategoryCache()

	return nil
}

// DeleteCategory deleta uma categoria pelo seu ID
func DeleteCategory(categoryID uint) error {
	// Opcional: verificar se a categoria possui produtos associados antes de deletar.
	err := repository.DeleteCategory(categoryID)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	// Invalida cache relacionado a categorias
	cacheService := &CacheService{}
	cacheService.InvalidateCategoryCache()

	return nil
}

// GetCategories retorna todas as categorias
func GetCategories() ([]model.Category, error) {
	cacheService := &CacheService{}

	// Tenta buscar no cache primeiro
	if categories, found := cacheService.GetCachedCategories(); found {
		return categories, nil
	}

	// Se não encontrou no cache, busca no banco
	categories, err := repository.GetCategories()
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	cacheService.SetCachedCategories(categories)

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

func AddCategoryImage(categoryID uint, imagePath string) error {
	category, err := repository.GetCategoryByID(categoryID)
	if err != nil || category == nil {
		return err
	}

	// Atualiza a imagem da categoria
	category.Image = imagePath
	if err := repository.UpdateCategory(category); err != nil {
		return err
	}

	return nil
}

func RemoveCategoryImage(categoryID uint) error {
	category, err := repository.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// Atualiza o banco de dados para remover a imagem da categoria
	category.Image = ""
	if err := repository.UpdateCategory(category); err != nil {
		return err
	}

	return nil
}

// UpdateCategory atualiza os dados de uma categoria no banco de dados
func UpdateCategory(categoryID uint, updatedCategory *model.Category) error {
	// Busca a categoria existente no banco de dados
	category, err := repository.GetCategoryByID(categoryID)
	if err != nil {
		return err // Retorna erro se não encontrar o produto
	}
	// Atualiza os campos da categoria
	category.Name = updatedCategory.Name
	category.Description = updatedCategory.Description
	category.Image = updatedCategory.Image

	// Salva as alterações no banco de dados
	if err := repository.UpdateCategory(category); err != nil {
		return err // Retorna erro se falhar ao atualizar no banco
	}

	// Invalida cache relacionado a categorias
	cacheService := &CacheService{}
	cacheService.InvalidateCategoryCache()

	return nil
}
