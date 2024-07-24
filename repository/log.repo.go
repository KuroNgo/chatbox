package repository

import (
	"chatbox/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type loggingRepository struct {
	database           *mongo.Database
	collectionActivity string
	collectionAdmin    string
}

func NewActivityRepository(db *mongo.Database, collectionActivity string, collectionAdmin string) domain.ILoggingRepository {
	return &loggingRepository{
		database:           db,
		collectionActivity: collectionActivity,
		collectionAdmin:    collectionAdmin,
	}
}

func (l loggingRepository) CreateOne(ctx context.Context, log domain.Logging) error {
	collectionActivity := l.database.Collection(l.collectionActivity)

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	expireTime := time.Date(tomorrow.Year(), tomorrow.Month()+1, tomorrow.Day(), 0, 0, 0, 0, time.UTC)

	log.ExpireAt = expireTime

	_, err := collectionActivity.InsertOne(ctx, &log)
	if err != nil {
		return err
	}

	// Tạo TTL Index
	index := mongo.IndexModel{
		Keys:    bson.M{"expire_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = collectionActivity.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
	}()
	wg.Wait()

	return nil
}

func (l loggingRepository) FetchMany(ctx context.Context, page string) ([]domain.Logging, error) {
	errCh := make(chan error, 1)

	collection := l.database.Collection(l.collectionActivity)
	collectionAdmin := l.database.Collection(l.collectionAdmin)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	//count, err := collection.CountDocuments(ctx, bson.D{})
	//if err != nil {
	//	return nil, err
	//}

	//totalPages := (count + int64(perPage) - 1) / int64(perPage)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
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

	var activities []domain.Logging
	activities = make([]domain.Logging, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var activity domain.Logging
		if err = cursor.Decode(&activity); err != nil {
			return nil, err
		}
		activity.ActivityTime = activity.ActivityTime.Add(7 * time.Hour)

		wg.Add(1)
		go func(activity domain.Logging) {
			defer wg.Done()
			var user domain.User
			filterUser := bson.M{"_id": activity.UserID}
			_ = collectionAdmin.FindOne(ctx, filterUser).Decode(&user)
			activity.UserID = user.ID

			// Thêm activity vào slice activities
			activities = append(activities, activity)
		}(activity)
	}
	wg.Wait()

	select {
	case err = <-errCh:
		return nil, err
	default:
		return activities, nil
	}
}
