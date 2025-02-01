package service

import (
	"errors"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(name, email, password string) (*model.User, error) {
	existingUser, _ := repository.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.UserRole, 
	}

	err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
