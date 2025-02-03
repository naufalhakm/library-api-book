package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"library-api-book/internal/commons/response"
	"library-api-book/internal/logger"
	"library-api-book/internal/models"
	"library-api-book/internal/params"
	"library-api-book/internal/repositories"
	"time"

	"github.com/redis/go-redis/v9"
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
	RedisClient    *redis.Client
	Logger         logger.Logger
}

func NewBookService(db *sql.DB, redisClient *redis.Client, bookRepository repositories.BookRepository, log logger.Logger) BookService {
	return &BookServiceImpl{
		DB:             db,
		BookRepository: bookRepository,
		RedisClient:    redisClient,
		Logger:         log,
	}
}

func (service *BookServiceImpl) CreateBook(ctx context.Context, req *params.BookRequest) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - CreateBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - CreateBook", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - CreateBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	book := models.Book{
		AuthorID:  req.AuthorID,
		Title:     req.Title,
		Stock:     req.Stock,
		PublishAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = service.BookRepository.CreateBook(ctx, tx, &book)
	if err != nil {
		service.Logger.Error("[BookService] Failed to create book - CreateBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to create book: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) GetDetailBook(ctx context.Context, id uint64) (*params.BookResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - GetDetailBook", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - GetDetailBook", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - GetDetailBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, id)
	if err != nil {
		service.Logger.Error("[BookService] Failed to find book by ID - GetDetailBook", map[string]interface{}{
			"book_id": id,
			"error":   err.Error(),
		})
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
		service.Logger.Error("[BookService] Failed to begin transaction - UpdateBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - UpdateBook", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - UpdateBook", map[string]interface{}{
				"error": err.Error(),
			})
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
		service.Logger.Error("[BookService] Failed to update book - UpdateBook", map[string]interface{}{
			"book_id": id,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to update book: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) DeleteBook(ctx context.Context, id uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - DeleteBook", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - DeleteBook", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - DeleteBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	err = service.BookRepository.DeleteBook(ctx, tx, id)
	if err != nil {
		service.Logger.Error("[BookService] Failed to delete book - DeleteBook", map[string]interface{}{
			"book_id": id,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to delete book: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) GetAllBooks(ctx context.Context, pagination *models.Pagination, search string) ([]*params.BookResponse, *response.CustomError) {
	cacheKey := fmt.Sprintf("books:%d:%d:%s", pagination.Page, pagination.PageSize, search)

	cachedData, err := service.RedisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var bookResponses []*params.BookResponse
		if err := json.Unmarshal(cachedData, &bookResponses); err == nil {
			service.Logger.Info("[BookService] Retrieved books from cache", map[string]interface{}{
				"cache_key": cacheKey,
			})
			return bookResponses, nil
		}
	}

	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - GetAllBooks", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - GetAllBooks", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - GetAllBooks", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	pagination.Offset = (pagination.Page - 1) * pagination.PageSize

	books, err := service.BookRepository.GetAllBooks(ctx, tx, pagination, search)
	if err != nil {
		service.Logger.Error("[BookService] Failed to fetch books - GetAllBooks", map[string]interface{}{
			"error": err.Error(),
		})
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

	serializedData, err := json.Marshal(bookResponses)
	if err != nil {
		service.Logger.Error("[BookService] Failed to serialize books for caching - GetAllBooks", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		if err := service.RedisClient.Set(ctx, cacheKey, serializedData, 5*time.Minute).Err(); err != nil {
			service.Logger.Error("[BookService] Failed to cache books - GetAllBooks", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			service.Logger.Info("[BookService] Cached books successfully", map[string]interface{}{
				"cache_key": cacheKey,
			})
		}
	}

	return bookResponses, nil
}

func (service *BookServiceImpl) GetRecommendationBook(ctx context.Context, id uint64) ([]*params.BookResponse, *response.CustomError) {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - GetRecommendationBook", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - GetRecommendationBook", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - GetRecommendationBook", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	books, err := service.BookRepository.GetRecommendationBooks(ctx, tx, id)
	if err != nil {
		service.Logger.Error("[BookService] Failed to fetch recommended books - GetRecommendationBook", map[string]interface{}{
			"error": err.Error(),
		})
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
		service.Logger.Error("[BookService] Failed to begin transaction - DecreaseStock", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - DecreaseStock", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - DecreaseStock", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, bookID)
	if err != nil {
		service.Logger.Error("[BookService] Failed to find book - DecreaseStock", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to find book: " + err.Error())
	}
	if book.Stock <= 0 {
		service.Logger.Warn("[BookService] Book is out of stock - DecreaseStock", map[string]interface{}{
			"book_id": bookID,
		})
		return response.BadRequestError("Book is out of stock")
	}

	book.Stock--
	err = service.BookRepository.UpdateBook(ctx, tx, book)
	if err != nil {
		service.Logger.Error("[BookService] Failed to update book stock - DecreaseStock", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to update book stock: " + err.Error())
	}

	return nil
}

func (service *BookServiceImpl) IncreaseStock(ctx context.Context, bookID uint64) *response.CustomError {
	tx, err := service.DB.Begin()
	if err != nil {
		service.Logger.Error("[BookService] Failed to begin transaction - IncreaseStock", map[string]interface{}{
			"error": err.Error(),
		})
		return response.GeneralError("Failed to connect to the database: " + err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to panic - IncreaseStock", map[string]interface{}{
				"error": p,
			})
		} else if err != nil {
			tx.Rollback()
			service.Logger.Error("[BookService] Transaction rolled back due to error - IncreaseStock", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			tx.Commit()
		}
	}()

	book, err := service.BookRepository.FindBookByID(ctx, tx, bookID)
	if err != nil {
		service.Logger.Error("[BookService] Failed to find book - IncreaseStock", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to find book: " + err.Error())
	}

	book.Stock++
	err = service.BookRepository.UpdateBook(ctx, tx, book)
	if err != nil {
		service.Logger.Error("[BookService] Failed to update book stock - IncreaseStock", map[string]interface{}{
			"book_id": bookID,
			"error":   err.Error(),
		})
		return response.GeneralError("Failed to update book stock: " + err.Error())
	}

	return nil
}
