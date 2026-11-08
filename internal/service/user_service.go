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

// LoginUser verifica o email e senha do usuário e gera um par de tokens (acesso + refresh)
func LoginUser(email, password string) (*util.TokenPair, error) {
	// Busca o usuário no banco de dados
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	// GetUserByEmail retorna (nil, nil) quando não encontra — sem este check,
	// o acesso a user.Password abaixo geraria panic (500) em vez de 401.
	if user == nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica a senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Gera o par de tokens JWT
	tokens, err := util.GenerateTokenPair(user)
	if err != nil {
		return nil, errors.New("erro ao gerar token")
	}

	return tokens, nil
}

// RefreshTokens valida um refresh token e emite um novo par de tokens.
func RefreshTokens(refreshToken string) (*util.TokenPair, error) {
	userID, err := util.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token inválido")
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	return util.GenerateTokenPair(user)
}