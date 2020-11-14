package test

//go:generate webrpc-gen -schema=contract.go -target=ts -client -server -out ./typescript.gen.ts
//go:generate webrpc-gen -schema=contract.go -target=go -client -server -pkg=main -out ./golang.gen.go
//go:generate webrpc-gen -schema=contract.go -target=js -client -server -extra=noexports -out ./javascript.gen.js

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
	ID      BookID
	Name    string
	Authors []Author
}

type Library interface {
	GetBooks(ctx context.Context) ([]*Book, string, error)
	BorrowBook(ctx context.Context, BookID int64) error
	GetBookAuthor(ctx context.Context, BookID int64) (Author, map[string]interface{}, regexp.Regexp, error)
	GetFile(ctx context.Context, id int64) (presenter.File, error)
}
