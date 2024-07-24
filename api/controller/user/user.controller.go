package user_controller

import (
	"chatbox/bootstrap"
	"chatbox/domain"
)

type UserController struct {
	UserUseCase domain.IUserUseCase
	Database    *bootstrap.Database
}
