package repository

import (
	"chatbox/domain"
	"chatbox/pkg/cache"
	"chatbox/pkg/helper"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type userRepository struct {
	database       *mongo.Database
	collectionUser string
}

func NewUserRepository(db *mongo.Database, collectionUser string) domain.IUserRepository {
	return &userRepository{
		database:       db,
		collectionUser: collectionUser,
	}
}

var (
	userCache  = cache.NewTTL[string, *domain.User]()
	usersCache = cache.NewTTL[string, domain.UserResponse]()

	wg sync.WaitGroup
	mu sync.Mutex
)

const (
	cacheTTL = 5 * time.Minute
)

func (u *userRepository) UpdateVerifyForChangePassword(ctx context.Context, user *domain.User) (*mongo.UpdateResult, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":   user.Verified,
		"updated_at": user.UpdatedAt,
	}}}

	data, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"password":          user.Password,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	filterUnique := bson.M{"email": user.Email}
	count, err := collectionUser.CountDocuments(ctx, filterUnique)
	if count > 0 {
		return errors.New("the email must be unique")
	}
	_, err = collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"phone":      user.Phone,
		"updated_at": user.UpdatedAt,
	}}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()
	return nil
}

func (u *userRepository) CheckVerify(ctx context.Context, verificationCode string) bool {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collectionUser.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		return false
	}

	return true
}

func (u *userRepository) GetByVerificationCode(ctx context.Context, verificationCode string) (*domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}

	var user domain.User
	err := collectionUser.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateImage(c context.Context, userID string, imageURL string) error {
	collectionUser := u.database.Collection(u.collectionUser)
	doc, err := helper.ToDoc(imageURL)
	objID, err := primitive.ObjectIDFromHex(userID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collectionUser.UpdateOne(c, filter, update)

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) UpdateVerify(ctx context.Context, user *domain.User) (*mongo.UpdateResult, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":          user.Verified,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	data, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) Create(c context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"email": user.Email}
	count, err := collectionUser.CountDocuments(c, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the email do not unique")
	}
	_, err = collectionUser.InsertOne(c, &user)

	wg.Add(1)
	go func() {
		defer wg.Done()
		usersCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) FetchMany(c context.Context) (domain.UserResponse, error) {
	usersCh := make(chan domain.UserResponse)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := usersCache.Get("users")
		if found {
			usersCh <- data
			return
		}
	}()

	go func() {
		defer close(usersCh)
		wg.Wait()
	}()

	userData := <-usersCh
	if !helper.IsZeroValue(userData) {
		return userData, nil
	}

	collectionUser := u.database.Collection(u.collectionUser)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collectionUser.Find(c, bson.D{}, opts)

	if err != nil {
		return domain.UserResponse{}, err
	}

	var users []domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return domain.UserResponse{}, err
	}

	response := domain.UserResponse{
		User: users,
	}

	usersCache.Set("users", response, 5*time.Minute)
	return response, err
}

func (u *userRepository) DeleteOne(c context.Context, userID string) error {
	collectionUser := u.database.Collection(u.collectionUser)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	_, err = collectionUser.DeleteOne(c, filter)

	wg.Add(2)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()

	go func() {
		defer wg.Done()
		usersCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) GetByEmail(c context.Context, email string) (*domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)
	var user domain.User
	err := collectionUser.FindOne(c, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (u *userRepository) Login(c context.Context, request domain.SignIn) (*domain.User, error) {
	userCh := make(chan *domain.User)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := userCache.Get(request.Email + request.Password)
		if found {
			userCh <- data
			return
		}
	}()

	go func() {
		defer close(userCh)
		wg.Wait()
	}()

	userData := <-userCh
	if !helper.IsZeroValue(userData) {
		return userData, nil
	}

	user, err := u.GetByEmail(c, request.Email)

	// Kiểm tra xem mật khẩu đã nhập có đúng với mật khẩu đã hash trong cơ sở dữ liệu không
	if err = helper.VerifyPassword(user.Password, request.Password); err != nil {
		return &domain.User{}, errors.New("email or password not found! ")
	}

	userCache.Set(request.Email+request.Password, user, cacheTTL)
	return user, nil
}

func (u *userRepository) GetByID(c context.Context, id string) (*domain.User, error) {
	userCh := make(chan *domain.User)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := userCache.Get(id)
		if found {
			userCh <- data
			return
		}
	}()

	go func() {
		defer close(userCh)
		wg.Wait()
	}()

	userData := <-userCh
	if !helper.IsZeroValue(userData) {
		return userData, nil
	}

	collectionUser := u.database.Collection(u.collectionUser)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = collectionUser.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return nil, err
	}

	userCache.Set(id, &user, cacheTTL)
	return &user, nil
}

func (u *userRepository) UpsertOne(c context.Context, email string, user *domain.User) (*domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"phone":      user.Phone,
		"provider":   user.Provider,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"role":       user.Role,
	}}}
	res := collectionUser.FindOneAndUpdate(c, filter, update, opts)

	var updatedUser *domain.User
	if err := res.Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()

	return updatedUser, nil
}

func (u *userRepository) UniqueVerificationCode(ctx context.Context, verificationCode string) bool {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collectionUser.CountDocuments(ctx, filter)
	if err != nil || count > 0 {
		return false
	}
	return true
}
