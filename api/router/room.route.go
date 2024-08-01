package router

import (
	room_controller "chatbox/api/controller/room"
	"chatbox/api/middlewares"
	"chatbox/bootstrap"
	"chatbox/domain"
	"chatbox/repository"
	"chatbox/usecase"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func RoomRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *echo.Group) {
	r := repository.NewRoomRepository(db, domain.CollectionRoom)
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	room := &room_controller.RoomController{
		RoomUseCase: usecase.NewRoomUseCase(r, timeout),
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/v1/room")
	router.Use(middlewares.DeserializeUser())
	router.GET("/fetch", room.FetchManyRoom())
	router.GET("/1/fetch", room.FetchOneRoom())
	router.GET("/fetch/name", room.FetchOneByName())
	router.POST("/create", room.CreateRoom())
	router.PATCH("/update", room.UpdateRoom())
	router.DELETE("/delete", room.DeleteRoom())
}
