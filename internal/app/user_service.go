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

func (s *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}
	u.Status = !u.Status
	return s.repo.UpdateUser(ctx, u)
}

func (s *UserService) GetUserCount(ctx context.Context) string {
	return "1234"
}

func (s *UserService) GetPageView(ctx context.Context) string {
	return "1212121"
}
