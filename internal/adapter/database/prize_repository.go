package database

import (
	"context"
	"spahtmx/internal/domain"
	"strconv"

	"github.com/uptrace/bun"
)

type PrizeBunRepository struct {
	DB *bun.DB
}

type PrizeBun struct {
	bun.BaseModel `bun:"table:prizes"`

	ID                int64         `bun:"id,pk,autoincrement"`
	Year              string        `bun:"year"`
	Category          string        `bun:"category"`
	OverallMotivation string        `bun:"overall_motivation"`
	Laureates         []LaureateBun `bun:"laureates,rel:has-many,join:id=prize_id"`
}

type LaureateBun struct {
	bun.BaseModel `bun:"table:laureates"`

	ID         int64  `bun:"id,pk,autoincrement"`
	Firstname  string `bun:"firstname"`
	Surname    string `bun:"surname"`
	Motivation string `bun:"motivation"`
	Share      string `bun:"share"`
	PrizeID    int64  `bun:"prize_id"`
}

func ToPrizeDomain(p PrizeBun) domain.Prize {

	return domain.Prize{
		Year:              p.Year,
		Category:          p.Category,
		OverallMotivation: p.OverallMotivation,
		Laureates: func() []domain.Laureate {
			var laureates []domain.Laureate
			for _, l := range p.Laureates {
				laureates = append(laureates, domain.Laureate{
					Firstname:  l.Firstname,
					Surname:    l.Surname,
					Motivation: l.Motivation,
					Share:      l.Share,
				})
			}
			return laureates
		}(),
	}

}

func FromPrizeDomain(prize domain.Prize) (*PrizeBun, error) {

	return &PrizeBun{
		Year:              prize.Year,
		Category:          prize.Category,
		OverallMotivation: prize.OverallMotivation,
		Laureates: func() []LaureateBun {
			var laureates []LaureateBun
			for _, l := range prize.Laureates {
				laureates = append(laureates, LaureateBun{
					Firstname:  l.Firstname,
					Surname:    l.Surname,
					Motivation: l.Motivation,
					Share:      l.Share,
				})
			}
			return laureates
		}(),
	}, nil
}

func (r *PrizeBunRepository) Save(ctx context.Context, prize domain.Prize) error {

	prizeBun, err := FromPrizeDomain(prize)
	if err != nil {
		return err
	}

	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(prizeBun).Exec(ctx)
		if err != nil {
			return err
		}

		if len(prizeBun.Laureates) == 0 {
			return nil
		}

		for i := range prizeBun.Laureates {
			prizeBun.Laureates[i].PrizeID = prizeBun.ID
		}

		_, err = tx.NewInsert().Model(&prizeBun.Laureates).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *PrizeBunRepository) FindAll(ctx context.Context) ([]domain.Prize, error) {

	var prizes []PrizeBun
	err := r.DB.NewSelect().Model(&prizes).Relation("Laureates").Scan(ctx)
	if err != nil {
		return nil, err
	}

	var domainPrizes []domain.Prize
	for _, p := range prizes {
		domainPrizes = append(domainPrizes, ToPrizeDomain(p))
	}

	return domainPrizes, nil
}

func (r *PrizeBunRepository) FindByID(ctx context.Context, id int64) (*domain.Prize, error) {

	var prize PrizeBun
	err := r.DB.NewSelect().Model(&prize).Relation("Laureates").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	domainPrize := ToPrizeDomain(prize)
	return &domainPrize, nil
}

func (r *PrizeBunRepository) DeleteByID(ctx context.Context, id int64) error {

	_, err := r.DB.NewDelete().Model((*PrizeBun)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrizeBunRepository) Update(ctx context.Context, prize domain.Prize) error {

	prizeBun, err := FromPrizeDomain(prize)
	if err != nil {
		return err
	}

	_, err = r.DB.NewUpdate().Model(prizeBun).Where("id = ?", prize.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrizeBunRepository) FindByYear(ctx context.Context, year string) ([]domain.Prize, error) {

	var prizes []PrizeBun
	err := r.DB.NewSelect().Model(&prizes).Relation("Laureates").Where("year = ?", year).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var domainPrizes []domain.Prize
	for _, p := range prizes {
		domainPrizes = append(domainPrizes, ToPrizeDomain(p))
	}

	return domainPrizes, nil
}

func (r *PrizeBunRepository) FindByCategory(ctx context.Context, category string) ([]domain.Prize, error) {

	var prizes []PrizeBun
	err := r.DB.NewSelect().Model(&prizes).Relation("Laureates").Where("category = ?", category).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var domainPrizes []domain.Prize
	for _, p := range prizes {
		domainPrizes = append(domainPrizes, ToPrizeDomain(p))
	}

	return domainPrizes, nil
}

func (r *PrizeBunRepository) FindByCategoryAndYear(ctx context.Context, category string, year string) ([]domain.Prize, error) {

	var prizes []PrizeBun
	err := r.DB.NewSelect().Model(&prizes).Relation("Laureates").Where("category = ? AND year = ?", category, year).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var domainPrizes []domain.Prize
	for _, p := range prizes {
		domainPrizes = append(domainPrizes, ToPrizeDomain(p))
	}

	return domainPrizes, nil
}

func (r *PrizeBunRepository) GetCategories(ctx context.Context) ([]string, error) {

	var categories []string
	err := r.DB.NewSelect().Model((*PrizeBun)(nil)).Column("category").Distinct().Scan(ctx, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PrizeBunRepository) GetYears(ctx context.Context) ([]string, error) {

	var years []string
	err := r.DB.NewSelect().Model((*PrizeBun)(nil)).Column("year").Distinct().Scan(ctx, &years)
	if err != nil {
		return nil, err
	}

	return years, nil
}

func (r *PrizeBunRepository) GetPrizes(ctx context.Context) ([]domain.Prize, error) {
	return r.FindAll(ctx)
}

func (r *PrizeBunRepository) GetPrize(ctx context.Context, id string) (domain.Prize, error) {
	prizeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return domain.Prize{}, err
	}

	prize, err := r.FindByID(ctx, prizeID)
	if err != nil {
		return domain.Prize{}, err
	}

	return *prize, nil
}

func (r *PrizeBunRepository) GetPrizesByYear(ctx context.Context, year string) ([]domain.Prize, error) {
	return r.FindByYear(ctx, year)
}

func (r *PrizeBunRepository) GetPrizesByCategory(ctx context.Context, category string) ([]domain.Prize, error) {
	return r.FindByCategory(ctx, category)
}

func (r *PrizeBunRepository) GetPrizesByCategoryAndYear(ctx context.Context, category string, year string) ([]domain.Prize, error) {
	return r.FindByCategoryAndYear(ctx, category, year)
}
