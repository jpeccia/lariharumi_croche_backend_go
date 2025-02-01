package repository

import (
	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := config.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}
