package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

/**
 * CacheControl returns a middleware that adds cache headers to responses.
 * Improves performance by allowing browsers and CDNs to cache responses.
 *
 * @param maxAge - Cache duration in seconds
 * @param isPublic - Whether the cache can be stored by shared caches (CDNs)
 */
func CacheControl(maxAge int, isPublic bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		visibility := "private"
		if isPublic {
			visibility = "public"
		}

		c.Header("Cache-Control", fmt.Sprintf("%s, max-age=%d", visibility, maxAge))
		c.Header("Vary", "Accept-Encoding")
		c.Next()
	}
}

/**
 * ETagMiddleware generates and validates ETags for response caching.
 * Returns 304 Not Modified if content hasn't changed.
 */
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() == 200 {
			timestamp := time.Now().Truncate(time.Minute).Unix()
			path := c.Request.URL.Path
			etag := generateETag(fmt.Sprintf("%s-%d", path, timestamp))
			c.Header("ETag", etag)

			if match := c.GetHeader("If-None-Match"); match == etag {
				c.AbortWithStatus(304)
				return
			}
		}
	}
}

/**
 * generateETag creates an MD5 hash for the ETag header.
 *
 * @param content - The content to hash
 * @returns - The ETag string with weak validator prefix
 */
func generateETag(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf(`W/"%s"`, hex.EncodeToString(hash[:]))
}
