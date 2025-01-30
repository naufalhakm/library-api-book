package services

import (
	"context"
	"library-api-book/internal/models"
	"library-api-book/internal/params"
	"library-api-book/internal/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBook_Success(t *testing.T) {
	mockBookRepo := new(repositories.MockBookRepository)

	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockDB.ExpectBegin()
	mockBookRepo.On("CreateBook", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDB.ExpectCommit()

	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockBookRepo,
	}

	req := params.BookRequest{
		AuthorID: 1,
		Title:    "Test Book",
		Stock:    100,
	}

	errCus := service.CreateBook(context.Background(), &req)

	assert.Nil(t, errCus)
	mockBookRepo.AssertExpectations(t)
}

func TestGetDetailBook_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("FindBookByID", mock.Anything, mock.Anything, uint64(1)).Return(&models.Book{
		ID:        1,
		AuthorID:  1,
		Title:     "Test Book",
		Stock:     10,
		PublishAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)
	mockDB.ExpectCommit()

	bookResponse, errResponse := service.GetDetailBook(context.Background(), 1)

	assert.Nil(t, errResponse)
	assert.Equal(t, uint64(1), bookResponse.ID)
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestUpdateBook_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("UpdateBook", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDB.ExpectCommit()

	req := &params.BookRequest{
		AuthorID: 1,
		Title:    "Updated Book",
		Stock:    5,
	}
	errResponse := service.UpdateBook(context.Background(), 1, req)

	assert.Nil(t, errResponse)
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestDeleteBook_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("DeleteBook", mock.Anything, mock.Anything, uint64(1)).Return(nil)
	mockDB.ExpectCommit()

	errResponse := service.DeleteBook(context.Background(), 1)

	assert.Nil(t, errResponse)
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestGetAllBooks_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("GetAllBooks", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*models.Book{
		{
			ID:        1,
			AuthorID:  1,
			Title:     "Test Book 1",
			Stock:     10,
			PublishAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)
	mockDB.ExpectCommit()

	pagination := &models.Pagination{
		Page:     1,
		PageSize: 10,
	}
	books, errResponse := service.GetAllBooks(context.Background(), pagination, "")

	assert.Nil(t, errResponse)
	assert.Equal(t, 1, len(books))
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestGetRecommendationBook_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("GetRecommendationBooks", mock.Anything, mock.Anything, uint64(1)).Return([]*models.Book{
		{
			ID:        1,
			AuthorID:  1,
			Title:     "Recommended Book",
			Stock:     10,
			PublishAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)
	mockDB.ExpectCommit()

	books, errResponse := service.GetRecommendationBook(context.Background(), 1)

	assert.Nil(t, errResponse)
	assert.Equal(t, 1, len(books))
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestDecreaseStock_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("FindBookByID", mock.Anything, mock.Anything, uint64(1)).Return(&models.Book{
		ID:    1,
		Stock: 10,
	}, nil)
	mockRepo.On("UpdateBook", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDB.ExpectCommit()

	errResponse := service.DecreaseStock(context.Background(), 1)

	assert.Nil(t, errResponse)
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}

func TestIncreaseStock_Success(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	mockRepo := new(repositories.MockBookRepository)
	service := &BookServiceImpl{
		DB:             db,
		BookRepository: mockRepo,
	}

	mockDB.ExpectBegin()
	mockRepo.On("FindBookByID", mock.Anything, mock.Anything, uint64(1)).Return(&models.Book{
		ID:    1,
		Stock: 10,
	}, nil)
	mockRepo.On("UpdateBook", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockDB.ExpectCommit()

	errResponse := service.IncreaseStock(context.Background(), 1)

	assert.Nil(t, errResponse)
	mockRepo.AssertExpectations(t)
	mockDB.ExpectationsWereMet()
}
