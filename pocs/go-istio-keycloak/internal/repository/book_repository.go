package repository

import (
	"database/sql"

	"book-service/internal/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(b *models.Book) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO books (title, author, isbn, year) VALUES (?, ?, ?, ?)`,
		b.Title, b.Author, b.ISBN, b.Year,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *BookRepository) GetByID(id int64) (*models.Book, error) {
	b := &models.Book{}
	err := r.db.QueryRow(
		`SELECT id, title, author, isbn, year, created_at FROM books WHERE id = ?`, id,
	).Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BookRepository) List(limit, offset int) ([]models.Book, error) {
	rows, err := r.db.Query(
		`SELECT id, title, author, isbn, year, created_at FROM books ORDER BY id LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (r *BookRepository) Count() (int64, error) {
	var total int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&total)
	return total, err
}

func (r *BookRepository) Delete(id int64) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
