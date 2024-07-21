package user_controller

import (
	"chatbox/pkg/jwt"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (u *UserController) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not logged in!",
			})
		}

		sub, err := jwt.ValidateToken(cookie.Value, u.Database.AccessTokenPublicKey)
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": "Failed to validate token: " + err.Error(),
			})
		}

		result, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(sub))
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": "Failed to get user data: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
			"user":   result,
		})
	}
}
