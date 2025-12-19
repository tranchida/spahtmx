package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID  bson.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	City string `bson:"city"`
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

	myusers := []User{
		{ID: bson.NewObjectID(), Name: "John Doe", City: "New York"},
		{ID: bson.NewObjectID(), Name: "Jane Smith", City: "Los Angeles"},
		{ID: bson.NewObjectID(), Name: "Mike Johnson", City: "Chicago"},
		{ID: bson.NewObjectID(), Name: "Emily Davis", City: "Houston"},
		{ID: bson.NewObjectID(), Name: "David Wilson", City: "Phoenix"},
		{ID: bson.NewObjectID(), Name: "Sarah Brown", City: "Philadelphia"},
		{ID: bson.NewObjectID(), Name: "Chris Lee", City: "San Antonio"},
		{ID: bson.NewObjectID(), Name: "Anna Garcia", City: "San Diego"},
		{ID: bson.NewObjectID(), Name: "James Martinez", City: "Dallas"},
		{ID: bson.NewObjectID(), Name: "Laura Rodriguez", City: "San Jose"},
		{ID: bson.NewObjectID(), Name: "Robert Hernandez", City: "Austin"},
		{ID: bson.NewObjectID(), Name: "Linda Lopez", City: "Jacksonville"},
		{ID: bson.NewObjectID(), Name: "Michael Gonzalez", City: "Fort Worth"},
		{ID: bson.NewObjectID(), Name: "Barbara Wilson", City: "Columbus"},
		{ID: bson.NewObjectID(), Name: "William Anderson", City: "Charlotte"},
		{ID: bson.NewObjectID(), Name: "Elizabeth Thomas", City: "San Francisco"},
		{ID: bson.NewObjectID(), Name: "David Taylor", City: "Indianapolis"},
		{ID: bson.NewObjectID(), Name: "Jennifer Moore", City: "Seattle"},
		{ID: bson.NewObjectID(), Name: "Richard Jackson", City: "Denver"},
		{ID: bson.NewObjectID(), Name: "Susan White", City: "Washington"},
		{ID: bson.NewObjectID(), Name: "Joseph Harris", City: "Boston"},
		{ID: bson.NewObjectID(), Name: "Margaret Martin", City: "El Paso"},
		{ID: bson.NewObjectID(), Name: "Thomas Thompson", City: "Nashville"},
		{ID: bson.NewObjectID(), Name: "Dorothy Garcia", City: "Detroit"},
		{ID: bson.NewObjectID(), Name: "Charles Martinez", City: "Oklahoma City"},
		{ID: bson.NewObjectID(), Name: "Karen Robinson", City: "Portland"},
	}

	colusers := client.Database("test").Collection("users")

	_, err = colusers.DeleteMany(ctx, bson.D{})

	cusers, err := colusers.InsertMany(ctx, myusers)
	if err != nil {
		log.Fatal("Failed to insert users", err)
	}
	log.Printf("Inserted user IDs: %v\n", cusers.InsertedIDs)

}
