package main

import (
	"context"
	"time"
)

type Kind uint32

type Empty struct {
}

type User struct {
	ID         uint64
	username   string
	createdAt *time.Time
}

type ExampleService interface {
	Ping(ctx context.Context) (bool, error)
	GetUser(ctx context.Context, userID uint64) (*User, error)
}