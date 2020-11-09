package typescript

import (
	"context"
)

type Kind uint32

type Empty struct {
}

type User struct {
	id         uint64
	username   string
	role       *Kind
	meta       map[string]interface{}
	internalID uint64
}

type Page struct {
	num uint32
}

type ExampleService interface {
	Ping(ctx context.Context) (bool, error)
	GetUser(ctx context.Context, userID uint64) (*User, error)
	FindUsers(ctx context.Context, q string) (*Page, []*User, error)
}
