package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB(uri string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	defer cancel()

	fmt.Println("Connect to MongoDB database successfully!")
	DB = client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("awayfromus").Collection(collectionName)
	return collection
}
