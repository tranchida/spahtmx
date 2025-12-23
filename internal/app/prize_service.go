package app

import (
	"context"
	"spahtmx/internal/domain"
)

type PrizeService struct {
	repo domain.PrizeRepository
}

func NewPrizeService(r domain.PrizeRepository) *PrizeService {
	return &PrizeService{
		repo: r,
	}
}

func (s *PrizeService) GetPrizes(ctx context.Context) ([]domain.Prize, error) {
	return s.repo.GetPrizes(ctx)
}

func (s *PrizeService) GetPrize(ctx context.Context, id string) (domain.Prize, error) {
	return s.repo.GetPrize(ctx, id)
}

func (s *PrizeService) GetPrizesByYear(ctx context.Context, year string) ([]domain.Prize, error) {
	return s.repo.GetPrizesByYear(ctx, year)
}

func (s *PrizeService) GetPrizesByCategory(ctx context.Context, category string) ([]domain.Prize, error) {
	return s.repo.GetPrizesByCategory(ctx, category)
}

func (s *PrizeService) GetPrizesByCategoryAndYear(ctx context.Context, category string, year string) ([]domain.Prize, error) {
	return s.repo.GetPrizesByCategoryAndYear(ctx, category, year)
}

func (s *PrizeService) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}

func (s *PrizeService) GetYears(ctx context.Context) ([]string, error) {
	return s.repo.GetYears(ctx)
}
