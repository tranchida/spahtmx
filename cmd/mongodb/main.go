package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
	City string        `bson:"city"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		options.Client().ApplyURI("mongodb://localhost:27017"),
		options.Client().SetAuth(options.Credential{
			Username: "root",
			Password: "example",
		}),
	)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect MongoDB", err)
		}
	}()

	cusers := client.Database("test").Collection("users")

	cursor, err := cusers.Find(ctx, bson.D{},
		options.Find().SetLimit(20))
	if err != nil {
		log.Fatal(err)
	}
	var u []User
	err = cursor.All(ctx, &u)
	if err != nil {
		log.Fatal(err)
	}

	for count, user := range u {
		jsonData, _ := json.Marshal(user)
		log.Printf("count %d id : %s\n", count, user.ID.Hex())
		log.Printf("User from MongoDB: %s\n", jsonData)
	}

}
