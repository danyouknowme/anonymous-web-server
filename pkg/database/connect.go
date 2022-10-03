package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	utils.LoadConfig()

	client, err := mongo.NewClient(options.Client().ApplyURI(utils.AppConfig.MongoUri))
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
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("awayfromus").Collection(collectionName)
	return collection
}
