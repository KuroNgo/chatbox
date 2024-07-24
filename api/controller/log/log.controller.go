package log_controller

import (
	"chatbox/bootstrap"
	"chatbox/domain"
)

type ActivityController struct {
	ActivityUseCase domain.IActivityUseCase
	UserUseCase     domain.IUserUseCase
	Database        *bootstrap.Database
}
