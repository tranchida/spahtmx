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

func FromDomain(user domain.User) *UserMongo {
	return nil
}

func (m MongoRepository) GetUsers(ctx context.Context) []domain.User {
	//TODO implement me
	users, err := m.DB.Collection("users").Find(ctx, bson.D{}, options.Find().SetLimit(10))
	if err != nil {
		panic(err)

	}
	var u []UserMongo
	err = users.All(context.Background(), &u)
	var domainUsers []domain.User
	for _, user := range u {
		domainUsers = append(domainUsers, user.ToDomain())
	}
	return domainUsers
}

func (m MongoRepository) UpdateUserStatus(ctx context.Context, id string) {

	objid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	user := m.DB.Collection("users").FindOne(ctx, bson.D{{"_id", objid}}, nil)
	var u UserMongo
	err = user.Decode(&u)
	if err != nil {
		panic(err)
	}
	u.Status = !u.Status
	_, err = m.DB.Collection("users").ReplaceOne(ctx, bson.D{{"_id", u.ID}}, u)
	if err != nil {
		panic(err)
	}
}

func (m MongoRepository) GetUserCount(ctx context.Context) string {
	return "210"
}

func (m MongoRepository) GetPageView(ctx context.Context) string {
	return "12345"
}
