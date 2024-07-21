package message_controller

import (
	"chatbox/bootstrap"
	"chatbox/domain"
)

type MessageController struct {
	MessageUseCase domain.IMessageUseCase
	UserUseCase    domain.IUserUseCase
	Database       *bootstrap.Database
}
