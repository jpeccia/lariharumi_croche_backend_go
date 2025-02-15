package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	Products []Product `json:"products" gorm:"foreignKey:CategoryID"`
}