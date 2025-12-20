package domain

import "context"

type UserRepository interface {
	GetUsers(ctx context.Context) []User
	GetUser(ctx context.Context, id string) User
	CreateUser(ctx context.Context, user User)
	UpdateUser(ctx context.Context, user User)
}
