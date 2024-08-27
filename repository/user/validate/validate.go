package validate

import (
	"chatbox/domain"
	"errors"
)

func IsInvalidUser(user *domain.User) error {
	if user.FullName == "" || user.Email == "" {
		return errors.New("the user's information cannot be empty")
	}
	return nil
}

func IsNilUsername(user *domain.User) error {
	if user.FullName == "" {
		return errors.New("the user's information cannot be empty")
	}
	return nil
}

func IsNilEmail(email string) error {
	if email == "" {
		return errors.New("the user's information cannot be empty")
	}
	return nil
}

func IsNilPasswordHash(user *domain.User) error {
	if user.Password == "" {
		return errors.New("the user's information cannot be empty")
	}
	return nil
}

func IsNilImage(avatarUrl string) error {
	if avatarUrl == "" {
		return errors.New("the user's information cannot be empty")
	}
	return nil
}
