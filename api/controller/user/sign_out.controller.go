package user_controller

import (
	"chatbox/pkg/review"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (u *UserController) LogoutUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not login!",
			})
		}

		user, err := u.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			err = c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "Unauthorized",
				"message": "You are not authorized to perform this action!",
			})

		}

		review.SetCookie(c, "access_token", "", -1, "/", "localhost", false, true)
		review.SetCookie(c, "refresh_token", "", -1, "/", "localhost", false, true)
		review.SetCookie(c, "logged_in", "", -1, "/", "localhost", false, false)

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})

	}
}
