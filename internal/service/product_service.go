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

	return nil
}

var productImagesDB = map[uint][]string{}

// AddProductImage adiciona um caminho de imagem ao produto
func AddProductImage(productID uint, imageURL string) error {
	if productID == 0 {
		return errors.New("ID do produto inválido")
	}

	productImagesDB[productID] = append(productImagesDB[productID], imageURL)
	return nil
}

func DeleteProductImage(productID uint, index int) error {
	images, exists := productImagesDB[productID]
	if !exists {
		return fmt.Errorf("Produto %d não encontrado", productID)
	}

	// Verifica se o índice é válido
	if index < 0 || index >= len(images) {
		return fmt.Errorf("Índice inválido")
	}

	// Remove a imagem com base no índice
	images = append(images[:index], images[index+1:]...)

	productImagesDB[productID] = images
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
	products, err := repository.GetProductsByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductImages(productID uint) ([]string, error) {
	images, exists := productImagesDB[productID]
	if !exists {
		return nil, fmt.Errorf("Nenhuma imagem encontrada para o produto %d", productID)
	}
	return images, nil
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

	return nil
}
