package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionLog = "log"
)

type Logging struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	RoomID       primitive.ObjectID `json:"room_id" bson:"room_id"`
	Method       string             `json:"method" bson:"method"`
	StatusCode   int                `json:"status_code" bson:"status_code"`
	BodySize     int                `json:"body_size" bson:"body_size"`
	Path         string             `json:"path" bson:"path"`
	Latency      string             `json:"latency" bson:"latency"`
	Error        string             `json:"error" bson:"error"`
	ActivityTime time.Time          `json:"activity_time" bson:"activity_time"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
}

type ILoggingRepository interface {
	CreateOne(ctx context.Context, log Logging) error
	FetchMany(ctx context.Context, page string) ([]Logging, error)
}

type IActivityUseCase interface {
	FetchMany(ctx context.Context, page string) ([]Logging, error)
	CreateOne(ctx context.Context, log Logging) error
}
