package repository

import (
	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"gorm.io/gorm"
)

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Retorna nil para o usuário e nil para o erro, indicando que o usuário não foi encontrado.
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
