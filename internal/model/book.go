package model

import (
	"encoding/json"
	"errors"
	"time"
)

type Book struct {
	BookID          string    `json:"bookId"`
	AuthorID        string    `json:"authorId"`
	PublisherID     string    `json:"publisherId"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publicationDate"`
	ISBN            string    `json:"isbn"`
	Pages           int       `json:"pages"`
	Genre           string    `json:"genre"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Quantity        int       `json:"quantity"`
}

// Custom UnmarshalJSON to handle date-only format
func (b *Book) UnmarshalJSON(data []byte) error {
	type Alias Book // Create alias to avoid infinite recursion
	aux := &struct {
		PublicationDate string `json:"publicationDate"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.PublicationDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.PublicationDate)
		if err != nil {
			return errors.New("invalid publicationDate format, expected YYYY-MM-DD")
		}
		b.PublicationDate = parsedDate
	}

	return nil
}

// Custom MarshalJSON to format date correctly when serializing
func (b Book) MarshalJSON() ([]byte, error) {
	type Alias Book
	return json.Marshal(&struct {
		PublicationDate string `json:"publicationDate"`
		*Alias
	}{
		PublicationDate: b.PublicationDate.Format("2006-01-02"),
		Alias:           (*Alias)(&b),
	})
}
