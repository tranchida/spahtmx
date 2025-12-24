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

type UserMongoRepository struct {
	DB *mongo.Database
}

type UserMongo struct {
	ID       bson.ObjectID `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Email    string        `bson:"email"`
	Status   bool          `bson:"status"`
}

func ToUserDomain(u UserMongo) domain.User {

	return domain.User{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Status:   u.Status,
	}

}

func FromUserDomain(user domain.User) (*UserMongo, error) {

	uid, err := bson.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, err
	}
	return &UserMongo{
		ID:       uid,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

func (m UserMongoRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	cursor, err := m.DB.Collection("users").Find(ctx, bson.D{}, options.Find().SetLimit(10))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var u []UserMongo
	err = cursor.All(ctx, &u)
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, user := range u {
		domainUsers = append(domainUsers, ToUserDomain(user))
	}
	return domainUsers, nil
}

func (m UserMongoRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	objid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, domain.ErrInvalidInput
	}
	user := m.DB.Collection("users").FindOne(ctx, bson.D{{Key: "_id", Value: objid}}, nil)
	var u UserMongo
	err = user.Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		slog.Error("Database error in GetUser", "error", err, "id", id)
		return domain.User{}, domain.ErrInternal
	}
	return ToUserDomain(u), nil
}

func (m UserMongoRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	res := m.DB.Collection("users").FindOne(ctx, bson.D{{Key: "username", Value: username}}, nil)
	var u UserMongo
	err := res.Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		slog.Error("Database error in GetByUsername", "error", err, "username", username)
		return domain.User{}, domain.ErrInternal
	}
	return ToUserDomain(u), nil
}

func (m UserMongoRepository) CreateUser(ctx context.Context, user domain.User) error {

	um, err := FromUserDomain(user)
	if err != nil {
		return err
	}
	_, err = m.DB.Collection("users").InsertOne(ctx, um)
	return err
}

func (m UserMongoRepository) UpdateUser(ctx context.Context, user domain.User) error {

	um, err := FromUserDomain(user)
	if err != nil {
		return err
	}
	_, err = m.DB.Collection("users").ReplaceOne(ctx, bson.D{{Key: "_id", Value: um.ID}}, um)
	return err
}

func (m UserMongoRepository) GetUserCount(ctx context.Context) string {
	return "210"
}

func (m UserMongoRepository) GetPageView(ctx context.Context) string {
	return "12345"
}
