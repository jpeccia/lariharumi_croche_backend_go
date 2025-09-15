package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/handler"


)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(config.CORSMiddleware())

	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(204)
	})

	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
		auth.POST("/login", handler.LoginUser)
	}

	categories := r.Group("/categories")
	{
		categories.GET("", handler.GetCategories)
		categories.GET("/:id/image", handler.GetCategoryImage)
	}

	products := r.Group("/products")
	{
		products.GET("/category/:id", handler.GetProductsByCategory)
		products.GET("", handler.GetProducts)
		products.GET("/:id/images", handler.GetProductImages)
	}

	admin := r.Group("").Use(middleware.AuthMiddleware("ADMIN"))
	{
		admin.POST("/products", handler.CreateProduct)
		admin.GET("/products/search", handler.SearchProducts)
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
	}

	return r
}
