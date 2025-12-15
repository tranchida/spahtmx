package app

import (
	"context"
	"spahtmx/internal/domain"
)

type userService struct{
	repo domain.UserRepository
}

func NewUserService(r domain.UserRepository) domain.UserRepository {
	return &userService{
		repo: r,
	}
}

func (s *userService) GetUsers() []domain.User {
	return s.repo.GetUsers()
}

func (s *userService) UpdateUserStatus(ctx context.Context, id string) {
	s.repo.UpdateUserStatus(ctx, id)
}

func (s *userService) GetUserCount() string {
	return s.repo.GetUserCount()
}

func (s *userService) GetPageView() string {
	return s.repo.GetPageView()
}