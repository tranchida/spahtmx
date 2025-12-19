package app

import (
	"context"
	"spahtmx/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(r domain.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) GetUsers(ctx context.Context) []domain.User {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id string) {
	s.repo.UpdateUserStatus(ctx, id)
}

func (s *UserService) GetUserCount(ctx context.Context) string {
	return s.repo.GetUserCount(ctx)
}

func (s *UserService) GetPageView(ctx context.Context) string {
	return s.repo.GetPageView(ctx)
}
