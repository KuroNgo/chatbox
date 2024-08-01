package repository

import (
	"chatbox/domain"
	"chatbox/pkg/cache"
	"chatbox/pkg/helper"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

var (
	messagesCache = cache.NewTTL[string, []domain.Message]()
	messageCache  = cache.NewTTL[string, domain.Message]()
)

const timeTL = 5 * time.Minute

func (m messageRepository) FetchMany(ctx context.Context, userID primitive.ObjectID, roomID string) ([]domain.Message, error) {
	errCh := make(chan error, 1)
	messagesCh := make(chan []domain.Message, 1)
	go func() {
		data, found := messagesCache.Get(userID.Hex() + roomID)
		if found {
			messagesCh <- data
			return
		}
	}()

	messagesData := <-messagesCh
	if !helper.IsZeroValue(messagesData) {
		return messagesData, nil
	}

	collectionMessage := m.database.Collection(m.collectionMessage)

	idRoom, _ := primitive.ObjectIDFromHex(roomID)
	filter := bson.M{"user_id": userID, "room_id": idRoom}
	cursor, err := collectionMessage.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
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

		wg.Add(1)
		go func(message domain.Message) {
			defer wg.Done()

			mu.Lock()
			messages = append(messages, message)
			mu.Unlock()
		}(message)
	}
	wg.Wait()

	messagesCache.Set(userID.Hex()+roomID, messages, timeTL)

	select {
	case err = <-errCh:
		return nil, err
	default:
		return messages, nil
	}
}

func (m messageRepository) FetchByOne(ctx context.Context, userID primitive.ObjectID, id string) (domain.Message, error) {
	messageCh := make(chan domain.Message)
	go func() {
		data, found := messageCache.Get(userID.Hex() + id)
		if found {
			messageCh <- data
			return
		}
	}()

	messageData := <-messageCh
	if !helper.IsZeroValue(messageData) {
		return messageData, nil
	}

	collectionMessage := m.database.Collection(m.collectionMessage)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"user_id": userID, "_id": ID}
	var message domain.Message
	err := collectionMessage.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		return domain.Message{}, nil
	}

	messageCache.Set(userID.Hex()+id, message, timeTL)
	return message, nil
}

func (m messageRepository) CreateOne(ctx context.Context, message domain.Message) error {
	collectionMessage := m.database.Collection(m.collectionMessage)
	_, err := collectionMessage.InsertOne(ctx, message)
	messagesCache.Clear()
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
	messageCache.Clear()
	return err
}

func (m messageRepository) DeleteOne(ctx context.Context, id string) error {
	collectionMessage := m.database.Collection(m.collectionMessage)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ID}
	_, err := collectionMessage.DeleteOne(ctx, filter)
	wg.Add(2)
	go func() {
		defer wg.Done()
		messageCache.Clear()
	}()
	go func() {
		defer wg.Done()
		messagesCache.Clear()
	}()
	wg.Wait()
	return err
}
