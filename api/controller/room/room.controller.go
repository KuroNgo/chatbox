package room_controller

import (
	"chatbox/bootstrap"
	"chatbox/domain"
)

type RoomController struct {
	RoomUseCase domain.IRoomUseCase
	UserUseCase domain.IUserUseCase
	Database    *bootstrap.Database
}
