package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/handler"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Grupo de autenticação
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
		auth.POST("/login", handler.LoginUser)
	}

	// Grupo público, sem autenticação
	products := r.Group("/products")
	{
		products.GET("/category/:id", handler.GetProductsByCategory) // Rota para obter produtos por categoria
		products.GET("/", handler.GetProducts)                       // GET products
		products.GET("/:id/images", handler.GetProductImages)        // GET product images
	}

	categories := r.Group("/categories")
	{
		categories.GET("/", handler.GetCategories)             // GET categories
		categories.GET("/:id/image", handler.GetCategoryImage) // GET category image
	}

	// Rota protegida para admin (somente para admin)
	admin := r.Group("").Use(middleware.AuthMiddleware("ADMIN"))
	{
		admin.POST("/products", handler.CreateProduct)          // POST create product
		admin.DELETE("/products/:id", handler.DeleteProduct)    // DELETE product
		admin.POST("/categories", handler.CreateCategory)       // POST create category
		admin.DELETE("/categories/:id", handler.DeleteCategory) // DELETE category
		// Rotas para upload de imagem
		admin.POST("/products/:id/upload-image", handler.UploadProductImage)
		admin.POST("/categories/:id/upload-image", handler.UploadCategoryImage)
	}

	return r
}
