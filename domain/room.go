package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionRoom = "room"
)

type Room struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Name   string             `bson:"name" json:"name"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

type Input struct {
	Name string `bson:"name" json:"name"`
}

type IRoomRepository interface {
	CreateRoom(ctx context.Context, room Room) error
	FetchManyRoom(ctx context.Context, userID primitive.ObjectID) ([]Room, error)
	FetchOneRoom(ctx context.Context, userID primitive.ObjectID, id string) (Room, error)
	FetchOneByName(ctx context.Context, userID primitive.ObjectID, name string) (Room, error)
	GetByName(ctx context.Context, userID primitive.ObjectID, name string) (*Room, error)
	UpdateRoom(ctx context.Context, userID primitive.ObjectID, room *Room) error
	DeleteRoom(ctx context.Context, userID primitive.ObjectID, id string) error
}

type IRoomUseCase interface {
	CreateRoom(ctx context.Context, room Room) error
	FetchManyRoom(ctx context.Context, userID primitive.ObjectID) ([]Room, error)
	FetchOneRoom(ctx context.Context, userID primitive.ObjectID, id string) (Room, error)
	FetchOneByName(ctx context.Context, userID primitive.ObjectID, name string) (Room, error)
	GetByName(ctx context.Context, userID primitive.ObjectID, name string) (*Room, error)
	UpdateRoom(ctx context.Context, userID primitive.ObjectID, room *Room) error
	DeleteRoom(ctx context.Context, userID primitive.ObjectID, id string) error
}
