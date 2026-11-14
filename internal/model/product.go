package model

import "gorm.io/gorm"

/**
 * Product represents a crochet product in the catalog.
 * Indexes:
 * - idx_product_name: Optimizes name-based searches
 * - idx_product_category: Optimizes category filtering
 * - idx_product_category_name: Composite index for category + name sorting
 */
type Product struct {
	gorm.Model
	Name        string `json:"name" gorm:"index:idx_product_name;index:idx_product_category_name,priority:2"`
	Description string `json:"description"`
	ImageUrls   string `json:"imageUrls" gorm:"type:text"`
	// ThumbUrls guarda, em paralelo a ImageUrls (mesma ordem), uma versão
	// reduzida de cada imagem (ImgBB display/medium). O grid do catálogo usa
	// essas miniaturas para não baixar os originais full-res. Imagens antigas
	// ficam com ThumbUrls vazio; nesse caso o front usa ImageUrls como fallback.
	ThumbUrls   string   `json:"thumbUrls" gorm:"type:text"`
	PriceRange  string   `json:"priceRange"`
	CategoryID  uint     `json:"categoryId" gorm:"index:idx_product_category;index:idx_product_category_name,priority:1"`
	Category    Category `json:"category"`
}
