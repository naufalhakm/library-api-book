package services

import (
	"context"
	"database/sql"
	"library-api-book/internal/commons/response"
	"library-api-book/internal/models"
	"library-api-book/internal/params"
	"library-api-book/internal/repositories"
	"time"
)

type BookService interface {
	CreateBook(ctx context.Context, req *params.BookRequest) *response.CustomError
	GetDetailBook(ctx context.Context, id uint64) (*params.BookResponse, *response.CustomError)
	UpdateBook(ctx context.Context, id uint64, req *params.BookRequest) *response.CustomError
	DeleteBook(ctx context.Context, id uint64) *response.CustomError
	GetAllBooks(ctx context.Context, pagination *models.Pagination, search string) ([]*params.BookResponse, *response.CustomError)
	GetRecommendationBook(ctx context.Context, id uint64) ([]*params.BookResponse, *response.CustomError)
	DecreaseStock(ctx context.Context, bookID uint64) *response.CustomError
	IncreaseStock(ctx context.Context, bookID uint64) *response.CustomError
}

type BookServiceImpl struct {
	DB             *sql.DB
	BookRepository repositories.BookRepository
}

func NewBookService(db *sql.DB, bookRepository repositories.BookRepository) BookService {
	return &BookServiceImpl{
		DB:             db,
		BookRepository: bookRepository,
	}
}

func (service *BookServiceImpl) CreateBook(ctx context.Context, req *params.BookRequest) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed Connection to database errors: " + err.Error())
	}
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var book = models.Book{
		AuthorID:  req.AuthorID,
		Title:     req.Title,
		Stock:     req.Stock,
		PublishAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.BookRepository.CreateBook(ctx, tx, &book)

	if err != nil {
		return response.GeneralError(err.Error())
	}

	return nil
}

func (service *BookServiceImpl) GetDetailBook(ctx context.Context, id uint64) (*params.BookResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, response.GeneralError("Failed Connection to database errors: " + err.Error())
	}
	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, id)
	if err != nil {
		return nil, response.NotFoundError("Book not found")
	}

	bookResponse := &params.BookResponse{
		ID:        book.ID,
		AuthorID:  book.AuthorID,
		Title:     book.Title,
		Stock:     book.Stock,
		PublishAt: book.PublishAt,
		UpdatedAt: book.UpdatedAt,
	}

	return bookResponse, nil
}

func (service *BookServiceImpl) UpdateBook(ctx context.Context, id uint64, req *params.BookRequest) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	book := models.Book{
		ID:        id,
		AuthorID:  req.AuthorID,
		Title:     req.Title,
		Stock:     req.Stock,
		UpdatedAt: time.Now(),
	}

	err = service.BookRepository.UpdateBook(ctx, tx, &book)
	if err != nil {
		tx.Rollback()
		return response.GeneralError("Failed to update book: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) DeleteBook(ctx context.Context, id uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = service.BookRepository.DeleteBook(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		return response.GeneralError("Failed to delete book: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) GetAllBooks(ctx context.Context, pagination *models.Pagination, search string) ([]*params.BookResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	pagination.Offset = (pagination.Page - 1) * pagination.PageSize

	books, err := service.BookRepository.GetAllBooks(ctx, tx, pagination, search)
	if err != nil {
		return nil, response.GeneralError("Failed to fetch books: " + err.Error())
	}

	bookResponses := make([]*params.BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = &params.BookResponse{
			ID:        book.ID,
			AuthorID:  book.AuthorID,
			Title:     book.Title,
			Stock:     book.Stock,
			PublishAt: book.PublishAt,
			UpdatedAt: book.UpdatedAt,
		}
	}

	pagination.PageCount = (pagination.TotalCount + pagination.PageSize - 1) / pagination.PageSize

	return bookResponses, nil
}

func (service *BookServiceImpl) GetRecommendationBook(ctx context.Context, id uint64) ([]*params.BookResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	books, err := service.BookRepository.GetRecommendationBooks(ctx, tx, id)
	if err != nil {
		return nil, response.GeneralError("Failed to fetch recommended books: " + err.Error())
	}

	bookResponses := make([]*params.BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = &params.BookResponse{
			ID:        book.ID,
			AuthorID:  book.AuthorID,
			Title:     book.Title,
			Stock:     book.Stock,
			PublishAt: book.PublishAt,
			UpdatedAt: book.UpdatedAt,
		}
	}

	return bookResponses, nil
}

func (service *BookServiceImpl) DecreaseStock(ctx context.Context, bookID uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, bookID)
	if err != nil {
		return response.GeneralError("Failed to find book: " + err.Error())
	}
	if book.Stock <= 0 {
		return response.BadRequestError("Book is out of stock")
	}

	book.Stock--
	err = service.BookRepository.UpdateBook(ctx, tx, book)
	if err != nil {
		return response.GeneralError("Failed to update book stock: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) IncreaseStock(ctx context.Context, bookID uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, bookID)
	if err != nil {
		return response.GeneralError("Failed to find book: " + err.Error())
	}
	if book.Stock <= 0 {
		return response.BadRequestError("Book is out of stock")
	}

	book.Stock++
	err = service.BookRepository.UpdateBook(ctx, tx, book)
	if err != nil {
		return response.GeneralError("Failed to update book stock: " + err.Error())
	}

	return nil
}
