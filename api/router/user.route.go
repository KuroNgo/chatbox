package router

import (
	user_controller "chatbox/api/controller/user"
	"chatbox/api/middlewares"
	"chatbox/bootstrap"
	"chatbox/domain"
	"chatbox/repository"
	"chatbox/usecase"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UserRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *echo.Group) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	user := &user_controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/v1/user")
	router.POST("/signup", user.SignUp())
	router.PATCH("/update", user.UpdateUser(), middlewares.DeserializeUser())
	router.PATCH("/verify", user.VerificationCode())
	router.PATCH("/verify/password", user.VerificationCodeForChangePassword())
	router.PATCH("/password/forget", user.ChangePassword())
	router.POST("/forget", user.ForgetPasswordInUser())
	router.GET("/info", user.GetMe())
	router.GET("/google/callback", user.GoogleLoginWithUser())
	router.POST("/login", user.LoginUser(), middlewares.RateLimiter())
	router.GET("/refresh", user.RefreshToken())
	router.GET("/logout", user.LogoutUser(), middlewares.DeserializeUser())
}