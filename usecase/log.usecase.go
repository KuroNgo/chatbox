package usecase

import (
	"chatbox/domain"
	"context"
	"time"
)

type loggingUseCase struct {
	loggingRepository domain.ILoggingRepository
	contextTimeout    time.Duration
}

func NewActivityUseCase(activityRepository domain.ILoggingRepository, timeout time.Duration) domain.IActivityUseCase {
	return &loggingUseCase{
		loggingRepository: activityRepository,
		contextTimeout:    timeout,
	}
}

func (l loggingUseCase) FetchMany(ctx context.Context, page string) ([]domain.Logging, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	log, err := l.loggingRepository.FetchMany(ctx, page)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (l loggingUseCase) CreateOne(ctx context.Context, log domain.Logging) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	err := l.loggingRepository.CreateOne(ctx, log)
	if err != nil {
		return err
	}

	return nil
}
