package usecase

import (
	"chatbox/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type roomUseCase struct {
	roomRepository domain.IRoomRepository
	contextTimeout time.Duration
}

func (r roomUseCase) GetByName(ctx context.Context, userID primitive.ObjectID, name string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	data, err := r.roomRepository.GetByName(ctx, userID, name)

	if err != nil {
		return &domain.Room{}, err
	}

	return data, nil
}

func NewRoomUseCase(roomRepository domain.IRoomRepository, timeout time.Duration) domain.IRoomUseCase {
	return &roomUseCase{
		roomRepository: roomRepository,
		contextTimeout: timeout,
	}
}

func (r roomUseCase) CreateRoom(ctx context.Context, room domain.Room) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	err := r.roomRepository.CreateRoom(ctx, room)

	if err != nil {
		return err
	}

	return nil
}

func (r roomUseCase) FetchManyRoom(ctx context.Context, userID primitive.ObjectID) ([]domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	data, err := r.roomRepository.FetchManyRoom(ctx, userID)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r roomUseCase) FetchOneRoom(ctx context.Context, userID primitive.ObjectID, id string) (domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	data, err := r.roomRepository.FetchOneRoom(ctx, userID, id)

	if err != nil {
		return domain.Room{}, err
	}

	return data, nil
}

func (r roomUseCase) FetchOneByName(ctx context.Context, userID primitive.ObjectID, name string) (domain.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	data, err := r.roomRepository.FetchOneByName(ctx, userID, name)

	if err != nil {
		return domain.Room{}, err
	}

	return data, nil
}

func (r roomUseCase) UpdateRoom(ctx context.Context, userID primitive.ObjectID, room *domain.Room) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	err := r.roomRepository.UpdateRoom(ctx, userID, room)

	if err != nil {
		return err
	}

	return nil
}

func (r roomUseCase) DeleteRoom(ctx context.Context, userID primitive.ObjectID, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()
	err := r.roomRepository.DeleteRoom(ctx, userID, id)

	if err != nil {
		return err
	}

	return nil
}
