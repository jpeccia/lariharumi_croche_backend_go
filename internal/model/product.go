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
	Name        string   `json:"name" gorm:"index:idx_product_name;index:idx_product_category_name,priority:2"`
	Description string   `json:"description"`
	ImageUrls   string   `json:"imageUrls"`
	PriceRange  string   `json:"priceRange"`
	CategoryID  uint     `json:"categoryId" gorm:"index:idx_product_category;index:idx_product_category_name,priority:1"`
	Category    Category `json:"category"`
}
