package domain

import "context"

type UserRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	UpdateUserStatus(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type PrizeRepository interface {
	GetPrizes(ctx context.Context) ([]Prize, error)
	GetPrize(ctx context.Context, id string) (Prize, error)
	GetPrizesByYear(ctx context.Context, year string) ([]Prize, error)
	GetPrizesByCategory(ctx context.Context, category string) ([]Prize, error)
	GetPrizesByCategoryAndYear(ctx context.Context, category string, year string) ([]Prize, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetYears(ctx context.Context) ([]string, error)
}
