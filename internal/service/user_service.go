package service

import (
	"errors"
	"log"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(name, email, password string) (*model.User, error) {
	// Tenta buscar um usuário com o email fornecido.
	existingUser, err := repository.GetUserByEmail(email)
	if err != nil {
		log.Printf("Erro ao buscar usuário: %v", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email já está em uso")
	}

	// Gera o hash da senha para armazenamento seguro.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		return nil, err
	}

	// Cria o objeto usuário.
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     model.UserRole, // Valor padrão para novos usuários.
	}

	// Salva o usuário no banco de dados.
	err = repository.CreateUser(user)
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		return nil, err
	}

	return user, nil
}
