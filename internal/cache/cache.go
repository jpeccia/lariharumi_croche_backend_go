package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

/**
 * ErrCacheMiss is returned when a key is not found or has expired in the cache.
 */
var ErrCacheMiss = errors.New("cache: key not found or expired")

/**
 * cacheEntry stores a serialized value along with its expiration timestamp.
 */
type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

/**
 * isExpired returns true if the cache entry has passed its TTL.
 */
func (e *cacheEntry) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

var (
	store   map[string]cacheEntry
	storeMu sync.RWMutex
)

/**
 * InitCache initializes the in-memory cache store and starts a background
 * goroutine that evicts expired entries every 5 minutes.
 */
func InitCache() {
	store = make(map[string]cacheEntry)
	go evictionLoop()
	log.Println("In-memory cache initialized")
}

/**
 * evictionLoop runs every 5 minutes, scanning and removing expired entries.
 * Uses a write lock only when deletions are needed to minimize contention.
 */
func evictionLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		storeMu.Lock()
		for key, entry := range store {
			if entry.isExpired() {
				delete(store, key)
			}
		}
		storeMu.Unlock()
	}
}

/**
 * Set stores a value in the in-memory cache with the given TTL.
 * The value is JSON-serialized before storage.
 *
 * @param key - Cache key identifier
 * @param value - Any JSON-serializable value
 * @param ttl - Time-to-live duration for this entry
 */
func Set(key string, value interface{}, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache set marshal error: %w", err)
	}

	storeMu.Lock()
	store[key] = cacheEntry{
		value:     jsonValue,
		expiresAt: time.Now().Add(ttl),
	}
	storeMu.Unlock()

	return nil
}

/**
 * Get retrieves a value from the in-memory cache and deserializes it into dest.
 * Returns ErrCacheMiss if the key does not exist or has expired.
 *
 * @param key - Cache key identifier
 * @param dest - Pointer to the destination for JSON deserialization
 */
func Get(key string, dest interface{}) error {
	storeMu.RLock()
	entry, exists := store[key]
	storeMu.RUnlock()

	if !exists || entry.isExpired() {
		return ErrCacheMiss
	}

	return json.Unmarshal(entry.value, dest)
}

/**
 * Delete removes a single key from the cache.
 *
 * @param key - Cache key identifier to remove
 */
func Delete(key string) error {
	storeMu.Lock()
	delete(store, key)
	storeMu.Unlock()
	return nil
}

/**
 * DeletePattern removes all keys that start with the given prefix.
 * The prefix is derived by trimming any trailing wildcard (*) characters.
 *
 * @param pattern - Key prefix pattern (e.g. "products*")
 */
func DeletePattern(pattern string) error {
	prefix := strings.TrimRight(pattern, "*")

	storeMu.Lock()
	for key := range store {
		if strings.HasPrefix(key, prefix) {
			delete(store, key)
		}
	}
	storeMu.Unlock()

	return nil
}

/**
 * GenerateKey builds a cache key by joining a prefix with variadic parameters.
 *
 * @param prefix - Key namespace prefix
 * @param params - Additional segments appended with ":" separator
 */
func GenerateKey(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}
