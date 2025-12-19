package domain

import "context"

type UserRepository interface {
	GetUsers(ctx context.Context) []User
	UpdateUserStatus(ctx context.Context, id string)
	GetUserCount(ctx context.Context) string
	GetPageView(ctx context.Context) string
}
