package middleware

import (
	"github.com/gin-gonic/gin"
)

/**
 * SecurityHeaders adds essential HTTP security headers to all responses.
 * Protects against common web vulnerabilities:
 * - X-Content-Type-Options: Prevents MIME type sniffing
 * - X-Frame-Options: Prevents clickjacking attacks
 * - X-XSS-Protection: Legacy XSS protection for older browsers
 * - Referrer-Policy: Controls referrer information leakage
 * - Permissions-Policy: Restricts browser feature access
 */
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}
