package test

import (
	"context"
	"regexp"	
	"github.com/vcilabs/hubs/data/presenter"
)

type BookID int64

type Empty struct {
}

type Author struct {
	ID       int64
	Name     string
	Metadata map[string]interface{}
}

type Book struct {
	ID  BookID
	Name    string
	Authors []Author
}

type Library interface {
	GetBooks(ctx context.Context) ([]*Book, string, error)
	BorrowBook(ctx context.Context, BookID int64) error
	GetBookAuthor(ctx context.Context, BookID int64) (Author, map[string]interface{}, regexp.Regexp, error)
	GetFile(ctx context.Context, id int64) (presenter.File, error)
}