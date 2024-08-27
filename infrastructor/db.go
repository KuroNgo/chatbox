package infrastructor

import (
	"chatbox/bootstrap"
	"chatbox/repository/user/data_seeder"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewMongoDatabase(env *bootstrap.Database) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := fmt.Sprintf("mongodb://localhost:27017/")
	mongoConn := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	// migration
	err = Migrations(ctx, client)
	if err != nil {
		return nil
	}

	return client
}

func CloseMongoDBConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connect to Mongo closed")
}

func Migrations(ctx context.Context, client *mongo.Client) error {
	// migration
	err := data_seeder.SeedUser(ctx, client)
	if err != nil {
		return nil
	}

	return nil
}
