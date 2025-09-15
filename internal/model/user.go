package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Não expor senha no JSON
	Role     Role   `json:"role" gorm:"default:USER"`
}
