package service

import (
	"book-api/internal/model"
	"book-api/internal/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

type BookService interface {
	CreateBook(book model.Book) (*model.Book, error)
	GetAllBooks() ([]model.Book, error)
	GetBookByID(id string) (*model.Book, error)
	UpdateBook(id string, book model.Book) (*model.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(book model.Book) (*model.Book, error) {
	// Validate required fields
	if book.Title == "" || book.AuthorID == "" || book.PublisherID == "" || book.ISBN == "" {
		return nil, errors.New("missing required fields")
	}

	// Generate a new UUID if not provided
	if book.BookID == "" {
		book.BookID = uuid.New().String()
	}

	// Set default quantity if not provided
	if book.Quantity == 0 {
		book.Quantity = 1
	}

	// Set current time if publication date is zero
	if book.PublicationDate.IsZero() {
		book.PublicationDate = time.Now()
	}

	err := s.repo.Create(book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *bookService) GetAllBooks() ([]model.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetBookByID(id string) (*model.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) UpdateBook(id string, book model.Book) (*model.Book, error) {
	// Ensure the ID in the path matches the book ID
	if id != book.BookID {
		return nil, errors.New("ID in path does not match book ID")
	}

	// Check if book exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update the book
	err = s.repo.Update(id, book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}
