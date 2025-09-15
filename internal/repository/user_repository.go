package repository

import (
	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"gorm.io/gorm"
)

// CreateUser cria um novo usuário no banco de dados
func CreateUser(user *model.User) error {
	if err := config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID retorna um usuário pelo seu ID
func GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retorna um usuário pelo seu email
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUsers retorna todos os usuários
func GetUsers() ([]model.User, error) {
	var users []model.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser atualiza as informações de um usuário
func UpdateUser(user *model.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser deleta um usuário pelo seu ID (soft delete)
func DeleteUser(userID uint) error {
	if err := config.DB.Delete(&model.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

// HardDeleteUser deleta permanentemente um usuário (apenas para admin)
func HardDeleteUser(userID uint) error {
	if err := config.DB.Unscoped().Delete(&model.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
