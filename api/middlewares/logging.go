package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Tiến hành xử lý yêu cầu
			err := next(c)

			// Ghi log thông tin
			latency := time.Since(start)
			method := c.Request().Method
			path := c.Request().URL.Path
			status := c.Response().Status
			fmt.Printf("%s %s %d %v\n", method, path, status, latency)

			return err
		}
	}

}
