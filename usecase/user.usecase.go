package usecase

import (
	"chatbox/domain"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type userUseCase struct {
	userRepository domain.IUserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(userRepository domain.IUserRepository, timeout time.Duration) domain.IUserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (u *userUseCase) UpdatePassword(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepository.UpdatePassword(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) GetByVerificationCode(ctx context.Context, verificationCode string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	user, err := u.userRepository.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (u *userUseCase) CheckVerify(ctx context.Context, verificationCode string) bool {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res := u.userRepository.CheckVerify(ctx, verificationCode)
	return res
}

func (u *userUseCase) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepository.Update(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) UpdateImage(c context.Context, userID string, imageURL string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	err := u.userRepository.UpdateImage(ctx, imageURL, userID)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	err := u.userRepository.Create(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Login(c context.Context, request domain.SignIn) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUseCase) Fetch(c context.Context) (domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.FetchMany(ctx)
	if err != nil {
		return domain.UserResponse{}, err
	}

	return user, err
}

func (u *userUseCase) GetByEmail(c context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUseCase) GetByID(c context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userUseCase) UpdateVerify(ctx context.Context, user *domain.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.userRepository.UpdateVerify(ctx, user)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *userUseCase) Delete(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepository.DeleteOne(ctx, userID)

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) UniqueVerificationCode(ctx context.Context, verificationCode string) bool {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	res := u.userRepository.UniqueVerificationCode(ctx, verificationCode)
	return res
}

func (u *userUseCase) UpdateVerifyForChangePassword(ctx context.Context, user *domain.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.userRepository.UpdateVerifyForChangePassword(ctx, user)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *userUseCase) UpsertUser(ctx context.Context, email string, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.userRepository.UpsertOne(ctx, email, user)

	if err != nil {
		return &domain.User{}, err
	}

	return data, nil
}

func (u *userUseCase) FetchMany(ctx context.Context) (domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	data, err := u.userRepository.FetchMany(ctx)

	if err != nil {
		return domain.UserResponse{}, err
	}

	return data, nil
}
