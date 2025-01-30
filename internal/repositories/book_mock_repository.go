package repositories

import (
	"context"
	"database/sql"
	"library-api-book/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error {
	args := m.Called(ctx, tx, book)
	return args.Error(0)
}

func (m *MockBookRepository) FindBookByID(ctx context.Context, tx *sql.Tx, id uint64) (*models.Book, error) {
	args := m.Called(ctx, tx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Book), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBookRepository) UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error {
	args := m.Called(ctx, tx, book)
	return args.Error(0)
}

func (m *MockBookRepository) DeleteBook(ctx context.Context, tx *sql.Tx, id uint64) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *MockBookRepository) GetAllBooks(ctx context.Context, tx *sql.Tx, pagination *models.Pagination, searchQuery string) ([]*models.Book, error) {
	args := m.Called(ctx, tx, pagination, searchQuery)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Book), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBookRepository) GetRecommendationBooks(ctx context.Context, tx *sql.Tx, userID uint64) ([]*models.Book, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Book), args.Error(1)
	}
	return nil, args.Error(1)
}
