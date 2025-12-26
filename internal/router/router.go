package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/handler"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/middleware"
)

/**
 * SetupRouter configures the Gin router with all routes and middleware.
 * Includes gzip compression for optimized response payloads.
 */
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(config.CORSMiddleware())

	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(204)
	})

	r.Static("/uploads", "./uploads")

	// Health check endpoints (públicos)
	r.GET("/health", handler.HealthCheck)
	r.GET("/ping", handler.Ping)

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
		auth.POST("/login", handler.LoginUser)
	}

	categories := r.Group("/categories")
	categories.Use(middleware.CacheControl(60, true))
	{
		categories.GET("", handler.GetCategories)
		categories.GET("/:id/image", handler.GetCategoryImage)
	}

	products := r.Group("/products")
	products.Use(middleware.CacheControl(60, true))
	{
		products.GET("/category/:id", handler.GetProductsByCategory)
		products.GET("", handler.GetProducts)
		products.GET("/search", handler.SearchProducts)
		products.GET("/:id/images", handler.GetProductImages)
	}

	r.GET("/promotion", handler.GetPromotion)

	admin := r.Group("").Use(middleware.AuthMiddleware("ADMIN"))
	{
		admin.POST("/products", handler.CreateProduct)
		admin.PATCH("/products/:id", handler.UpdateProduct)
		admin.PUT("/categories/:id", handler.UpdateCategory)
		admin.DELETE("/products/:id", handler.DeleteProduct)
		admin.POST("/categories", handler.CreateCategory)
		admin.DELETE("/categories/:id", handler.DeleteCategory)
		admin.POST("/products/:id/upload-images", handler.UploadProductImages)
		admin.POST("/categories/:id/upload-image", handler.UploadCategoryImage)
		admin.DELETE("/products/:id/images/:index", handler.DeleteProductImage)
		admin.DELETE("/categories/:id/image", handler.DeleteCategoryImage)
		admin.GET("/products/:id/upload-progress", handler.GetUploadProgress)

		admin.PUT("/promotion", handler.UpdatePromotion)
	}

	return r
}
