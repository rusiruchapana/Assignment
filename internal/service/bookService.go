package service

import (
	"book-api/internal/model"
	"book-api/internal/repository"
	"errors"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

type BookService interface {
	CreateBook(book model.Book) (*model.Book, error)
	GetAllBooks() ([]model.Book, error)
	GetBookByID(id string) (*model.Book, error)
	UpdateBook(id string, book model.Book) (*model.Book, error)
	DeleteBook(id string) error
	SearchBooks(query string) ([]model.Book, error)
	GetBooksPaginated(limit, offset int) ([]model.Book, error)
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

func (s *bookService) DeleteBook(id string) error {
	// Check if book exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *bookService) SearchBooks(query string) ([]model.Book, error) {
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if query == "" {
		return books, nil
	}

	// Convert query to lowercase once
	lowerQuery := strings.ToLower(query)

	// Create a channel to collect results
	resultChan := make(chan model.Book)
	var wg sync.WaitGroup

	// Split books into chunks for concurrent processing
	chunkSize := 100 // Adjust based on your expected dataset size
	chunks := chunkBooks(books, chunkSize)

	// Process each chunk in a goroutine
	for _, chunk := range chunks {
		wg.Add(1)
		go func(books []model.Book) {
			defer wg.Done()
			for _, book := range books {
				if strings.Contains(strings.ToLower(book.Title), lowerQuery) ||
					strings.Contains(strings.ToLower(book.Description), lowerQuery) {
					resultChan <- book
				}
			}
		}(chunk)
	}

	// Close the channel when all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from the channel
	var results []model.Book
	for book := range resultChan {
		results = append(results, book)
	}

	return results, nil
}

// Helper function to split books into chunks
func chunkBooks(books []model.Book, chunkSize int) [][]model.Book {
	var chunks [][]model.Book
	for i := 0; i < len(books); i += chunkSize {
		end := i + chunkSize
		if end > len(books) {
			end = len(books)
		}
		chunks = append(chunks, books[i:end])
	}
	return chunks
}

func (s *bookService) GetBooksPaginated(limit, offset int) ([]model.Book, error) {
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if offset >= len(books) {
		return []model.Book{}, nil
	}

	end := offset + limit
	if end > len(books) {
		end = len(books)
	}

	return books[offset:end], nil
}
