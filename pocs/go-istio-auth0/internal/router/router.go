package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"book-service/internal/config"
	"book-service/internal/handler"
	"book-service/internal/middleware"
	"book-service/internal/repository"
)

func New(cfg config.Config, repo *repository.BookRepository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "book-service",
			"endpoints": []gin.H{
				{"method": "GET", "path": "/health", "auth": false, "description": "liveness/readiness"},
				{"method": "GET", "path": "/", "auth": false, "description": "list endpoints"},
				{"method": "POST", "path": "/api/v1/books", "auth": true, "description": "create a book"},
				{"method": "GET", "path": "/api/v1/books?page=1&page_size=10", "auth": true, "description": "list books, paginated"},
				{"method": "GET", "path": "/api/v1/books/:id", "auth": true, "description": "get one book"},
				{"method": "DELETE", "path": "/api/v1/books/:id", "auth": true, "description": "delete a book"},
			},
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h := handler.NewBookHandler(repo)

	api := r.Group("/api/v1")
	api.Use(middleware.Auth(cfg.AuthHeader))
	{
		api.POST("/books", h.Create)
		api.GET("/books", h.List)
		api.GET("/books/:id", h.Get)
		api.DELETE("/books/:id", h.Delete)
	}

	return r
}
