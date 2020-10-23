package contract

import (
	"context"
)

type Author struct {
	ID       int64
	Name     string
	Metadata map[string]string
}

type BookID int64

type Book struct {
	ID      BookID
	Name    string
	Authors []Author
}

type Library interface {
	GetBooks(ctx context.Context) ([]Book, error)
	BorrowBook(ctx context.Context, BookID int64) error
	GetBookAuthor(ctx context.Context, BookID int64) (Author, error)
}
