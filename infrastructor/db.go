package infrastructor

import (
	"chatbox/bootstrap"
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

	return client
}

func CloseMongoDBConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connect to Mongo closed")
}
