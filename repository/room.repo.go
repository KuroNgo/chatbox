package repository

import (
	"chatbox/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type roomRepository struct {
	database       *mongo.Database
	collectionRoom string
}

func NewRoomRepository(db *mongo.Database, collectionRoom string) domain.IRoomRepository {
	return &roomRepository{
		database:       db,
		collectionRoom: collectionRoom,
	}
}

func (r roomRepository) GetByName(ctx context.Context, userID primitive.ObjectID, name string) (*domain.Room, error) {
	collectionRoom := r.database.Collection(r.collectionRoom)

	filter := bson.M{"user_id": userID, "name": name}
	var room *domain.Room
	err := collectionRoom.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return &domain.Room{}, err
	}

	return room, nil
}

func (r roomRepository) CreateRoom(ctx context.Context, room domain.Room) error {
	collectionRoom := r.database.Collection(r.collectionRoom)
	_, err := collectionRoom.InsertOne(ctx, room)
	return err
}

func (r roomRepository) FetchManyRoom(ctx context.Context, userID primitive.ObjectID) ([]domain.Room, error) {
	collectionRoom := r.database.Collection(r.collectionRoom)

	filter := bson.M{"user_id": userID}
	cursor, err := collectionRoom.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var rooms []domain.Room
	for cursor.Next(ctx) {
		var room domain.Room
		err = cursor.Decode(&room)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r roomRepository) FetchOneRoom(ctx context.Context, userID primitive.ObjectID, id string) (domain.Room, error) {
	collectionRoom := r.database.Collection(r.collectionRoom)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"user_id": userID, "_id": ID}

	var room domain.Room
	err := collectionRoom.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return domain.Room{}, err
	}

	return room, nil
}

func (r roomRepository) FetchOneByName(ctx context.Context, userID primitive.ObjectID, name string) (domain.Room, error) {
	collectionRoom := r.database.Collection(r.collectionRoom)

	filter := bson.M{"user_id": userID, "name": name}
	var room domain.Room
	err := collectionRoom.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return domain.Room{}, err
	}

	return room, nil
}

func (r roomRepository) UpdateRoom(ctx context.Context, userID primitive.ObjectID, room *domain.Room) error {
	collectionRoom := r.database.Collection(r.collectionRoom)

	filter := bson.M{"user_id": userID, "room": room}
	upd := bson.M{"$set": bson.M{
		"name": room.Name,
	}}

	_, err := collectionRoom.UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}

	return nil
}

func (r roomRepository) DeleteRoom(ctx context.Context, userID primitive.ObjectID, id string) error {
	collectionRoom := r.database.Collection(r.collectionRoom)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"user_id": userID, "_id": ID}
	_, err := collectionRoom.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
