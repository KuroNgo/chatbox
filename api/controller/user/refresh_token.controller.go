package user_controller

import (
	"chatbox/pkg/jwt"
	"chatbox/pkg/review"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RefreshToken godoc
// @Summary Làm mới token truy cập
// @Description API này làm mới token truy cập dựa trên token làm mới được gửi qua cookie.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Security ApiKeyAuth
// @Router /api/user/refresh [get]
func (u *UserController) RefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		message := "could not refresh access token"

		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": message,
			})
		}

		sub, err := jwt.ValidateToken(cookie.Value, u.Database.RefreshTokenPublicKey)
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}

		user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(sub))
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": "the user belonging to this token no logger exists",
			})

		}

		access_token, err := jwt.CreateToken(u.Database.AccessTokenExpiresIn, user.ID, u.Database.AccessTokenPrivateKey)
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": err.Error(),
			})

		}

		gin_fake.SetCookie(c, "access_token", access_token, u.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
		gin_fake.SetCookie(c, "logged_in", "true", u.Database.AccessTokenMaxAge*60, "/", "localhost", false, false)

		return c.JSON(http.StatusOK, echo.Map{
			"status":       "success",
			"access_token": access_token,
		})
	}
}
