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
	u := s.repo.GetUser(ctx, id)
	u.Status = !u.Status
	s.repo.UpdateUser(ctx, u)
}

func (s *UserService) GetUserCount(ctx context.Context) string {
	return "1234"
}

func (s *UserService) GetPageView(ctx context.Context) string {
	return "1212121"
}
