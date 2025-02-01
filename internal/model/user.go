package model

import "gorm.io/gorm"

type UserRole string

const (
	AdminRole UserRole = "ADMIN"
	UserRole  UserRole = "USER"
)

type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Role     UserRole `json:"role"` // Role do usuário: ADMIN ou USER
}