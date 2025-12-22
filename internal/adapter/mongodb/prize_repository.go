package mongodb

import (
	"context"
	"errors"
	"spahtmx/internal/domain"

	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PrizeMongoRepository struct {
	DB *mongo.Database
}

type PrizeMongo struct {
	ID                bson.ObjectID   `bson:"_id"`
	Year              string          `bson:"year"`
	Category          string          `bson:"category"`
	OverallMotivation string          `bson:"overallMotivation,omitempty"`
	Laureates         []LaureateMongo `bson:"laureates"`
}

type LaureateMongo struct {
	Firstname  string `bson:"firstname"`
	Surname    string `bson:"surname"`
	Motivation string `bson:"motivation"`
	Share      string `bson:"share"`
}

func ToPrizeDomain(p PrizeMongo) domain.Prize {

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

func FromPrizeDomain(prize domain.Prize) (*PrizeMongo, error) {

	var uid bson.ObjectID
	var err error

	if prize.ID == "" {
		uid = bson.NewObjectID()
	} else {
		uid, err = bson.ObjectIDFromHex(prize.ID)
		if err != nil {
			return nil, err
		}
	}

	return &PrizeMongo{
		ID:                uid,
		Year:              prize.Year,
		Category:          prize.Category,
		OverallMotivation: prize.OverallMotivation,
		Laureates: func() []LaureateMongo {
			var laureates []LaureateMongo
			for _, l := range prize.Laureates {
				laureates = append(laureates, LaureateMongo{
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

func (m PrizeMongoRepository) GetPrizes(ctx context.Context) ([]domain.Prize, error) {
	cursor, err := m.DB.Collection("prize").Find(ctx, bson.D{}, options.Find().SetLimit(10))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var u []PrizeMongo
	err = cursor.All(ctx, &u)
	if err != nil {
		return nil, err
	}

	var domainPrizes []domain.Prize
	for _, prize := range u {
		domainPrizes = append(domainPrizes, ToPrizeDomain(prize))
	}
	return domainPrizes, nil
}

func (m PrizeMongoRepository) GetPrize(ctx context.Context, id string) (domain.Prize, error) {
	objid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.Prize{}, domain.ErrInvalidInput
	}
	user := m.DB.Collection("prize").FindOne(ctx, bson.D{{Key: "_id", Value: objid}}, nil)
	var p PrizeMongo
	err = user.Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Prize{}, domain.PrizeNotFound
		}
		slog.Error("Database error in GetPrize", "error", err, "id", id)
		return domain.Prize{}, domain.ErrInternal
	}
	return ToPrizeDomain(p), nil
}

func (m PrizeMongoRepository) InsertPrizes(ctx context.Context, prizes []domain.Prize) error {
	//TODO implement me
	panic("implement me")
}
