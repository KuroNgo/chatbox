package usecase

import (
	"chatbox/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type messageUseCase struct {
	messageRepository domain.IMessageRepository
	contextTimeout    time.Duration
}

func NewMessageUseCase(messageRepository domain.IMessageRepository, contextTimeout time.Duration) domain.IMessageUseCase {
	return &messageUseCase{
		messageRepository: messageRepository,
		contextTimeout:    contextTimeout,
	}
}

func (m messageUseCase) FetchMany(ctx context.Context, userID primitive.ObjectID, roomID string) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	log, err := m.messageRepository.FetchMany(ctx, userID, roomID)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (m messageUseCase) FetchByOne(ctx context.Context, userID primitive.ObjectID, id string) (domain.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	log, err := m.messageRepository.FetchByOne(ctx, userID, id)
	if err != nil {
		return domain.Message{}, err
	}

	return log, nil
}

func (m messageUseCase) CreateOne(ctx context.Context, message domain.Message) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.messageRepository.CreateOne(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (m messageUseCase) UpdateOne(ctx context.Context, userID primitive.ObjectID, message domain.Message) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.messageRepository.UpdateOne(ctx, userID, message)
	if err != nil {
		return err
	}

	return nil
}

func (m messageUseCase) DeleteOne(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	err := m.messageRepository.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
