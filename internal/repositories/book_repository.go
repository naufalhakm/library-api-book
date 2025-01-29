package repositories

import (
	"context"
	"database/sql"
	"errors"
	"library-api-book/internal/models"
)

type BookRepository interface {
	CreateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error
	FindBookByID(ctx context.Context, tx *sql.Tx, id uint64) (*models.Book, error)
	UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error
	DeleteBook(ctx context.Context, tx *sql.Tx, id uint64) error
	GetAllBooks(ctx context.Context, tx *sql.Tx, pagination *models.Pagination, searchQuery string) ([]*models.Book, error)
	GetRecommendationBooks(ctx context.Context, tx *sql.Tx, userID uint64) ([]*models.Book, error)
}

type BookRepositoryImpl struct {
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{}
}

func (repository *BookRepositoryImpl) CreateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error {
	query := `INSERT INTO books (author_id, title, stock, publish_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	response, err := tx.ExecContext(ctx, query, book.AuthorID, book.Title, book.Stock, book.PublishAt, book.UpdatedAt)
	if err != nil || response == nil {
		return errors.New("Failed to create a book, transaction rolled back. Reason: " + err.Error())
	}

	return nil
}

func (repository *BookRepositoryImpl) FindBookByID(ctx context.Context, tx *sql.Tx, id uint64) (*models.Book, error) {
	query := "SELECT id, author_id, title, stock, publish_at, updated_at FROM books WHERE id = $1"
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var book = models.Book{}
	if rows.Next() {
		err := rows.Scan(&book.ID, &book.AuthorID, &book.Title, &book.Stock, &book.PublishAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &book, nil
	} else {
		return nil, errors.New("book is not found")
	}
}

func (repository *BookRepositoryImpl) UpdateBook(ctx context.Context, tx *sql.Tx, book *models.Book) error {
	query := `UPDATE books SET author_id = $1, title = $2, stock = $3, updated_at = $4 WHERE id = $5`

	_, err := tx.ExecContext(ctx, query,
		book.AuthorID,
		book.Title,
		book.Stock,
		book.UpdatedAt,
		book.ID,
	)
	if err != nil {
		return errors.New("Failed to update a book, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *BookRepositoryImpl) DeleteBook(ctx context.Context, tx *sql.Tx, id uint64) error {
	SQL := `DELETE FROM books WHERE id = $1`

	_, err := tx.ExecContext(ctx, SQL, id)
	if err != nil {
		return errors.New("Failed to update a book, transaction rolled back. Reason: " + err.Error())
	}
	return nil
}

func (repository *BookRepositoryImpl) GetAllBooks(ctx context.Context, tx *sql.Tx, pagination *models.Pagination, searchQuery string) ([]*models.Book, error) {
	query := `
		SELECT id, author_id, title, stock, publish_at, updated_at 
		FROM books 
	`

	var params []interface{}
	params = append(params, pagination.PageSize, pagination.Offset)

	if searchQuery != "" {
		query += ` WHERE title ILIKE $3`
		params = append(params, "%"+searchQuery+"%")
	}

	query += ` ORDER BY title ASC, publish_at DESC LIMIT $1 OFFSET $2`

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.AuthorID, &book.Title, &book.Stock, &book.PublishAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)

	}
	return books, nil
}
func (repository *BookRepositoryImpl) GetRecommendationBooks(ctx context.Context, tx *sql.Tx, userID uint64) ([]*models.Book, error) {
	query := `SELECT b.id, b.title, b.author_id, b.stock, b.publish_at, b.updated_at
		FROM books b
		JOIN book_categories bc ON b.id = bc.book_id
		WHERE bc.category_id IN (
			SELECT bc.category_id
			FROM user_activities ua
			JOIN book_categories bc ON ua.book_id = bc.book_id
			WHERE ua.user_id = $1
		)
		AND b.id NOT IN (
			SELECT book_id FROM borrows WHERE user_id = $1 AND returned_at IS NULL
		)
		ORDER BY RANDOM()
		LIMIT 10;`
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.AuthorID, &book.Title, &book.Stock, &book.PublishAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)

	}
	return books, nil
}
