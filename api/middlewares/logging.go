package middlewares

import (
	log_controller "chatbox/api/controller/log"
	"chatbox/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func LoggerMiddleware(logger *zerolog.Logger, activity *log_controller.ActivityController) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			path := c.Request().URL.Path

			err := next(c) // Call the next handler

			param := struct {
				TimeStamp    time.Time
				Path         string
				ClientIP     string
				Method       string
				StatusCode   int
				Latency      time.Duration
				ErrorMessage string
			}{
				TimeStamp:  time.Now(),
				Path:       path,
				ClientIP:   c.RealIP(),
				Method:     c.Request().Method,
				StatusCode: c.Response().Status,
			}

			param.Latency = time.Since(start).Truncate(time.Millisecond)

			expireDuration := 30 * 24 * time.Hour
			currentTime := time.Now()
			expireValue := currentTime.Add(expireDuration)

			if param.StatusCode >= 500 || err != nil || (param.Method == "DELETE" && param.StatusCode == 200) {
				currentUser := c.Get("currentUser")
				user, _ := activity.UserUseCase.GetByID(c.Request().Context(), fmt.Sprintf("%s", currentUser))

				logger.Error().
					Str("client_id", param.ClientIP).
					Str("method", param.Method).
					Int("status_code", param.StatusCode).
					Int("body_size", int(c.Response().Size)).
					Str("path", param.Path).
					Str("latency", param.Latency.String()).
					Msg(param.ErrorMessage)

				newLog := domain.Logging{
					ID:           primitive.NewObjectID(),
					UserID:       user.ID,
					Method:       param.Method,
					StatusCode:   param.StatusCode,
					BodySize:     int(c.Response().Size),
					Path:         path,
					Latency:      param.Latency.String(),
					Error:        param.ErrorMessage,
					ActivityTime: param.TimeStamp,
					ExpireAt:     expireValue,
				}

				createErr := activity.ActivityUseCase.CreateOne(c.Request().Context(), newLog)
				if createErr != nil {
					logger.Error().Err(createErr).Msg("Failed to create activity log")
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"status": "error",
						"error":  "Failed to create activity log",
					})
				}
			}

			return err
		}
	}

}
