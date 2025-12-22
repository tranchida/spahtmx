package mongodb

import (
	"context"
	"spahtmx/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoRepository struct {
	DB *mongo.Database
}

type UserMongo struct {
	ID       bson.ObjectID `bson:"_id"`
	Username string        `bson:"username"`
	Email    string        `bson:"email"`
	Status   bool          `bson:"status"`
}

func (u *UserMongo) ToDomain() domain.User {

	return domain.User{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Status:   u.Status,
	}

}

func FromDomain(user domain.User) (*UserMongo, error) {

	uid, err := bson.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, err
	}
	return &UserMongo{
		ID:       uid,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

func (m MongoRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
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
		domainUsers = append(domainUsers, user.ToDomain())
	}
	return domainUsers, nil
}

func (m MongoRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	objid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}
	user := m.DB.Collection("users").FindOne(ctx, bson.D{{Key: "_id", Value: objid}}, nil)
	var u UserMongo
	err = user.Decode(&u)
	if err != nil {
		return domain.User{}, err
	}
	return u.ToDomain(), nil
}

func (m MongoRepository) CreateUser(ctx context.Context, user domain.User) error {
	um, err := FromDomain(user)
	if err != nil {
		return err
	}
	_, err = m.DB.Collection("users").InsertOne(ctx, um)
	return err
}

func (m MongoRepository) UpdateUser(ctx context.Context, user domain.User) error {
	u, err := FromDomain(user)
	if err != nil {
		return err
	}
	_, err = m.DB.Collection("users").ReplaceOne(ctx, bson.D{{Key: "_id", Value: u.ID}}, u)
	return err
}

func (m MongoRepository) GetUserCount(ctx context.Context) string {
	return "210"
}

func (m MongoRepository) GetPageView(ctx context.Context) string {
	return "12345"
}
