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

func (s *PrizeService) InsertPrizes(ctx context.Context, prizes []domain.Prize) error {
	return s.repo.InsertPrizes(ctx, prizes)
}
