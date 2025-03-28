package repository

import (
	"book-api/internal/model"
	"encoding/json"
	"os"
	"sync"
)

type BookRepository interface {
	Create(book model.Book) error
	GetAll() ([]model.Book, error)
}

type fileBookRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewFileBookRepository(filePath string) BookRepository {
	return &fileBookRepository{filePath: filePath}
}

func (r *fileBookRepository) readBooks() ([]model.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Book{}, nil
		}
		return nil, err
	}

	var books []model.Book
	if len(data) == 0 {
		return books, nil
	}

	err = json.Unmarshal(data, &books)
	return books, err
}

func (r *fileBookRepository) writeBooks(books []model.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *fileBookRepository) Create(book model.Book) error {
	books, err := r.readBooks()
	if err != nil {
		return err
	}

	books = append(books, book)
	return r.writeBooks(books)
}

func (r *fileBookRepository) GetAll() ([]model.Book, error) {
	return r.readBooks()
}
