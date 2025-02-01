package service

import (
	"errors"
	"log"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/util"
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

// LoginUser verifica o email e senha do usuário, e gera um token JWT se forem válidos
func LoginUser(email, password string) (string, error) {
	// Busca o usuário no banco de dados
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("usuário não encontrado")
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("senha incorreta")
	}

	// Gera o token JWT
	token, err := util.GenerateToken(user)
	if err != nil {
		return "", errors.New("erro ao gerar token")
	}

	return token, nil
}