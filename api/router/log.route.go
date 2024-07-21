package router

import (
	log_controller "chatbox/api/controller/log"
	"chatbox/bootstrap"
	"chatbox/domain"
	"chatbox/repository"
	"chatbox/usecase"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *echo.Group) {
	ac := repository.NewActivityRepository(db, domain.CollectionLog, domain.CollectionUser)
	users := repository.NewUserRepository(db, domain.CollectionUser)

	activity := &log_controller.ActivityController{
		ActivityUseCase: usecase.NewActivityUseCase(ac, timeout),
		UserUseCase:     usecase.NewUserUseCase(users, timeout),
		Database:        env,
	}

	router := group.Group("/activity")
	router.GET("/fetch", activity.FetchManyActivity())
}
