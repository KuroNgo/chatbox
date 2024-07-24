package log_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *ActivityController) FetchManyActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "You are not logged in!",
			})

		}

		user, err := a.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			return c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Unauthorized",
				"message": "You are not authorized to perform this action!",
			})
		}

		page := c.QueryParam("page")
		if page == "" {
			page = "golang" // Giá trị mặc định
		}

		activity, err := a.ActivityUseCase.FetchMany(ctx, page)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"activity_log": activity,
		})
	}
}
