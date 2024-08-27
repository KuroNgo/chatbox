package data_seeder

import (
	"chatbox/domain"
	"chatbox/pkg/helper"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var user1 = domain.User{
	ID:        primitive.NewObjectID(),
	FullName:  "admin",
	Email:     "admin@admin.com",
	Password:  "12345",
	Phone:     "0329245971",
	Verified:  true,
	Provider:  "app",
	Role:      "employee",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var user2 = domain.User{
	ID:        primitive.NewObjectID(),
	FullName:  "admin2",
	Email:     "admin2@admin.com",
	Password:  "12345",
	Phone:     "0329245971",
	Verified:  true,
	Provider:  "app",
	Role:      "employee",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func SeedUser(ctx context.Context, client *mongo.Client) error {
	collectionUser := client.Database("Chatbox").Collection("user")

	count, err := collectionUser.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	user1.Password, err = helper.HashPassword(user1.Password)
	if err != nil {
		return err
	}

	user2.Password, err = helper.HashPassword(user2.Password)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = collectionUser.InsertOne(ctx, user1)
		if err != nil {
			return err
		}

		_, err = collectionUser.InsertOne(ctx, user2)
		if err != nil {
			return err
		}
	}

	return nil
}
