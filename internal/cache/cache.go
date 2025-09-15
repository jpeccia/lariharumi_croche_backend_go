package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis inicializa a conexão com o Redis
func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: os.Getenv("REDIS_PASSWORD"), // sem senha por padrão
		DB:       0,                           // usar o banco padrão
	})

	// Testa a conexão
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Aviso: Redis não disponível, continuando sem cache: %v", err)
		RedisClient = nil
		return
	}

	log.Println("Redis conectado com sucesso!")
}

// Set armazena um valor no cache com TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	if RedisClient == nil {
		return nil // Redis não disponível, não é erro
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, jsonValue, ttl).Err()
}

// Get recupera um valor do cache
func Get(key string, dest interface{}) error {
	if RedisClient == nil {
		return redis.Nil // Redis não disponível
	}

	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete remove uma chave do cache
func Delete(key string) error {
	if RedisClient == nil {
		return nil // Redis não disponível, não é erro
	}

	return RedisClient.Del(ctx, key).Err()
}

// DeletePattern remove todas as chaves que correspondem ao padrão
func DeletePattern(pattern string) error {
	if RedisClient == nil {
		return nil // Redis não disponível, não é erro
	}

	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// GenerateKey gera uma chave de cache baseada em parâmetros
func GenerateKey(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}
