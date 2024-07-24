package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
	"time"
)

const (
	maxRequests     = 5
	perMinutePeriod = 15 * time.Second
)

var (
	ipRequestsCounts = make(map[string]int) // can use some distrubuted db
	mutex            = &sync.Mutex{}
)

func RateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.Request().Header.Get("X-Real-IP") // Hoặc sử dụng c.Request().RemoteAddr để lấy IP trực tiếp
			if ip == "" {
				ip = c.Request().RemoteAddr
			}

			mutex.Lock()
			defer mutex.Unlock()
			count := ipRequestsCounts[ip]
			if count >= maxRequests {
				return c.JSON(http.StatusTooManyRequests, gin.H{
					"message": "Gửi quá nhiều request, vui lòng đợi 15 giây để thực hiện request tiếp theo!",
					"status":  "fail",
				})

			}

			ipRequestsCounts[ip] = count + 1
			time.AfterFunc(perMinutePeriod, func() {
				mutex.Lock()
				defer mutex.Unlock()

				ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
			})
			return next(c)
		}
	}
}
