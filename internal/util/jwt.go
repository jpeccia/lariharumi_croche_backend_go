package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

var jwtSecret []byte
var refreshSecret []byte

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

/**
 * init validates JWT secrets on startup.
 * Requires JWT_SECRET with minimum 32 characters for security.
 */
func init() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Println("WARNING: JWT_SECRET should have at least 32 characters for security")
	}
	if secret == "" {
		secret = "development-secret-key-change-in-production"
		log.Println("WARNING: Using default JWT_SECRET. Set JWT_SECRET environment variable in production!")
	}
	jwtSecret = []byte(secret)

	refreshSecretEnv := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecretEnv == "" {
		refreshSecretEnv = secret + "-refresh"
	}
	refreshSecret = []byte(refreshSecretEnv)
}

/**
 * CustomClaims extends JWT claims with user role.
 */
type CustomClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

/**
 * TokenPair contains both access and refresh tokens.
 */
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

/**
 * GenerateTokenPair creates both access and refresh tokens for a user.
 *
 * @param user - The user to generate tokens for
 * @returns - TokenPair with access and refresh tokens
 */
func GenerateTokenPair(user *model.User) (*TokenPair, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
	}, nil
}

/**
 * GenerateToken generates a JWT access token for the user (legacy support).
 *
 * @param user - The user to generate token for
 * @returns - JWT token string
 */
func GenerateToken(user *model.User) (string, error) {
	return generateAccessToken(user)
}

func generateAccessToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(AccessTokenDuration)

	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "laribackend",
		},
		Role: string(user.Role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(RefreshTokenDuration)

	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "laribackend-refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

/**
 * ParseToken validates the access token and returns user ID and role.
 *
 * @param tokenString - The JWT token to parse
 * @returns - User ID, role, and any error
 */
func ParseToken(tokenString string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", errors.New("invalid claims format")
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", err
	}

	return uint(id), claims.Role, nil
}

/**
 * ParseRefreshToken validates a refresh token and returns the user ID.
 *
 * @param tokenString - The refresh token to parse
 * @returns - User ID and any error
 */
func ParseRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("invalid claims format")
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
