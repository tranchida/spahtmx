package domain

import "context"

type UserRepository interface {
	GetUsers() []User
	UpdateUserStatus(ctx context.Context, id string)
	GetUserCount() string
	GetPageView() string
}