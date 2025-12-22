package domain

import "context"

type UserRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
}

type PrizeRepository interface {
	GetPrizes(ctx context.Context) ([]Prize, error)
	GetPrize(ctx context.Context, id string) (Prize, error)
	InsertPrizes(ctx context.Context, prizes []Prize) error
}
