package factory

import (
	"database/sql"
	"library-api-book/internal/controllers"
	"library-api-book/internal/repositories"
	"library-api-book/internal/services"
)

type Provider struct {
	BookProvider controllers.BookController
	BookService  services.BookService
}

func InitFactory(db *sql.DB) *Provider {

	bookRepo := repositories.NewBookRepository()
	borrowRepo := repositories.NewBorrowRepository()

	bookService := services.NewBookService(db, bookRepo, borrowRepo)
	bookController := controllers.NewBookController(bookService)

	return &Provider{
		BookProvider: bookController,
		BookService:  bookService,
	}
}
