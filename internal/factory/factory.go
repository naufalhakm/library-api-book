package factory

import (
	"database/sql"
	"library-api-book/internal/controllers"
	"library-api-book/internal/logger"
	"library-api-book/internal/repositories"
	"library-api-book/internal/services"
	"log"

	"github.com/redis/go-redis/v9"
)

type Provider struct {
	BookProvider controllers.BookController
	BookService  services.BookService
}

func InitFactory(db *sql.DB, redis *redis.Client) *Provider {
	newLog, err := logger.NewLogger("./var/log/book.log")
	if err != nil {
		log.Fatalf("[Logger] Failed to initialize book service logger: %v", err)
	}

	bookRepo := repositories.NewBookRepository()

	bookService := services.NewBookService(db, redis, bookRepo, newLog)
	bookController := controllers.NewBookController(bookService)

	return &Provider{
		BookProvider: bookController,
		BookService:  bookService,
	}
}
