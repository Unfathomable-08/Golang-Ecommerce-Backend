package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client *mongo.Client
	DB *mongo.Database
)

func ConnectDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)		// for setting the time limit of the connection to the database 
	defer cancel()		// defer is used to clean up resources when the function returns

	// connect DB and store it in cleint
	client, err := mongo.Connect(
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")),
	)

	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	DB = client.Database(os.Getenv("DB_NAME"))

	fmt.Println("MongoDB Connected")

	return nil
}
