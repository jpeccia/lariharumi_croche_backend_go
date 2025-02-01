package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImageUrls   string   `json:"imageUrls"`
	PriceRange  string   `json:"priceRange"`
	CategoryID  uint     `json:"categoryId"`
	Category    Category `json:"category"`
}
