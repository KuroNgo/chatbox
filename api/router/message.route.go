package router

import (
	message_controller "chatbox/api/controller/message"
	"chatbox/bootstrap"
	"chatbox/domain"
	repository2 "chatbox/repository/message/repository"
	"chatbox/repository/user/repository"
	"chatbox/usecase"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func MessageRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *echo.Group) {
	m := repository2.NewMessageRepository(db, domain.CollectionMessage)
	ur := user_repository.NewUserRepository(db, domain.CollectionUser)

	message := &message_controller.MessageController{
		MessageUseCase: usecase.NewMessageUseCase(m, timeout),
		UserUseCase:    usecase.NewUserUseCase(ur, timeout),
		Database:       env,
	}

	router := group.Group("/message")
	//router.Use(middlewares.DeserializeUser())
	router.GET("/ws", message.Setup())
	router.GET("/fetch", message.FetchMany())
	router.DELETE("/delete", message.DeleteOne())
}
