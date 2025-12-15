package app

import (
	"context"
	"spahtmx/internal/domain"
)

type UserService struct{
	repo domain.UserRepository
}

func NewUserService(r domain.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) GetUsers() []domain.User {
	return s.repo.GetUsers()
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id string) {
	s.repo.UpdateUserStatus(ctx, id)
}

func (s *UserService) GetUserCount() string {
	return s.repo.GetUserCount()
}

func (s *UserService) GetPageView() string {
	return s.repo.GetPageView()
}