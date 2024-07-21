package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionMessage = "message"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	RoomID    primitive.ObjectID `bson:"room_id" json:"room_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ToUserID  primitive.ObjectID `bson:"to_user_id" json:"to_user_id"`
	Text      string             `bson:"text" json:"text"`
	TimeStamp time.Time          `bson:"time_stamp" json:"time_stamp"`
	Color     string             `bson:"color" json:"color"`
}

type IMessageRepository interface {
	FetchMany(ctx context.Context, userID primitive.ObjectID, roomID string) ([]Message, error)
	FetchByOne(ctx context.Context, userID primitive.ObjectID, id string) (Message, error)
	CreateOne(ctx context.Context, message Message) error
	UpdateOne(ctx context.Context, userID primitive.ObjectID, message Message) error
	DeleteOne(ctx context.Context, id string) error
}

type IMessageUseCase interface {
	FetchMany(ctx context.Context, userID primitive.ObjectID, roomID string) ([]Message, error)
	FetchByOne(ctx context.Context, userID primitive.ObjectID, id string) (Message, error)
	CreateOne(ctx context.Context, message Message) error
	UpdateOne(ctx context.Context, userID primitive.ObjectID, message Message) error
	DeleteOne(ctx context.Context, id string) error
}
