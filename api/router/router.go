package router

import (
	"chatbox/api/middlewares"
	"chatbox/bootstrap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func SetUp(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, echo *echo.Echo) {
	publicRouter := echo.Group("/api")
	swaggerRouter := echo.Group("")

	value := Activity(env, timeout, db)

	// Middleware
	publicRouter.Use(
		middlewares.CORSPublic(),
		middleware.Recover(),
		middlewares.LoggerMiddleware(&log.Logger, value),
	)

	// This is a CORS method for check IP validation
	publicRouter.OPTIONS("/*path", middlewares.OptionsMessage())

	// All Public APIs
	UserRouter(env, timeout, db, publicRouter)
	RoomRouter(env, timeout, db, publicRouter)
	ActivityRoute(env, timeout, db, publicRouter)
	MessageRouter(env, timeout, db, publicRouter)
	SwaggerRouter(swaggerRouter)
}
