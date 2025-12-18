package model

import "gorm.io/gorm"

/**
 * Product represents a crochet product in the catalog.
 * Includes indexes for optimized name search and category filtering.
 */
type Product struct {
	gorm.Model
	Name        string   `json:"name" gorm:"index:idx_product_name"`
	Description string   `json:"description"`
	ImageUrls   string   `json:"imageUrls"`
	PriceRange  string   `json:"priceRange"`
	CategoryID  uint     `json:"categoryId" gorm:"index:idx_product_category"`
	Category    Category `json:"category"`
}
