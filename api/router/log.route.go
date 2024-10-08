package router

import (
	log_controller "chatbox/api/controller/log"
	"chatbox/api/middlewares"
	"chatbox/bootstrap"
	"chatbox/domain"
	repository2 "chatbox/repository/log/repository"
	"chatbox/repository/user/repository"
	"chatbox/usecase"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *echo.Group) {
	ac := repository2.NewActivityRepository(db, domain.CollectionLog, domain.CollectionUser)
	users := user_repository.NewUserRepository(db, domain.CollectionUser)

	activity := &log_controller.ActivityController{
		ActivityUseCase: usecase.NewActivityUseCase(ac, timeout),
		UserUseCase:     usecase.NewUserUseCase(users, timeout),
		Database:        env,
	}

	router := group.Group("/activity")
	router.Use(middlewares.DeserializeUser())
	router.GET("/fetch", activity.FetchManyActivity())
}

func Activity(env *bootstrap.Database, timeout time.Duration, db *mongo.Database) *log_controller.ActivityController {
	ac := repository2.NewActivityRepository(db, domain.CollectionLog, domain.CollectionUser)
	users := user_repository.NewUserRepository(db, domain.CollectionUser)

	activity := &log_controller.ActivityController{
		ActivityUseCase: usecase.NewActivityUseCase(ac, timeout),
		UserUseCase:     usecase.NewUserUseCase(users, timeout),
		Database:        env,
	}

	return activity
}
