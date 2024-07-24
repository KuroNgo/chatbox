package repository

import (
	"chatbox/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepository struct {
	database          *mongo.Database
	collectionMessage string
}

func NewMessageRepository(database *mongo.Database, collectionMessage string) domain.IMessageRepository {
	return &messageRepository{
		database:          database,
		collectionMessage: collectionMessage,
	}
}

func (m messageRepository) FetchMany(ctx context.Context, userID primitive.ObjectID, roomID string) ([]domain.Message, error) {
	collectionMessage := m.database.Collection(m.collectionMessage)

	idRoom, _ := primitive.ObjectIDFromHex(roomID)
	filter := bson.M{"user_id": userID, "room_id": idRoom}
	cursor, err := collectionMessage.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var messages []domain.Message
	messages = make([]domain.Message, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var message domain.Message
		err = cursor.Decode(&message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (m messageRepository) FetchByOne(ctx context.Context, userID primitive.ObjectID, id string) (domain.Message, error) {
	collectionMessage := m.database.Collection(m.collectionMessage)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"user_id": userID, "_id": ID}
	var message domain.Message
	err := collectionMessage.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		return domain.Message{}, nil
	}

	return message, nil
}

func (m messageRepository) CreateOne(ctx context.Context, message domain.Message) error {
	collectionMessage := m.database.Collection(m.collectionMessage)
	_, err := collectionMessage.InsertOne(ctx, message)
	return err
}

func (m messageRepository) UpdateOne(ctx context.Context, userID primitive.ObjectID, message domain.Message) error {
	collectionMessage := m.database.Collection(m.collectionMessage)

	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{
		"text":       message.Text,
		"time_stamp": message.TimeStamp,
		"color":      message.Color,
	}}

	_, err := collectionMessage.UpdateOne(ctx, filter, update)
	return err
}

func (m messageRepository) DeleteOne(ctx context.Context, id string) error {
	collectionMessage := m.database.Collection(m.collectionMessage)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	_, err := collectionMessage.DeleteOne(ctx, filter)

	return err
}
