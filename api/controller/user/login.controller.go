package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/jwt"
	"chatbox/pkg/review"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (u *UserController) LoginUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		//  Lấy thông tin từ request
		var adminInput domain.SignIn
		if err := c.Bind(&adminInput); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)

		}

		var userInput domain.SignIn
		userInput.Email = adminInput.Email
		userInput.Password = adminInput.Password

		// Tìm kiếm user trong database
		user, err := u.UserUseCase.Login(ctx, userInput)
		if err == nil {
			// Generate token
			accessToken, err := jwt.CreateToken(u.Database.AccessTokenExpiresIn, user.ID, u.Database.AccessTokenPrivateKey)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "fail",
					"message": err.Error()},
				)

			}

			refreshToken, err := jwt.CreateToken(u.Database.RefreshTokenExpiresIn, user.ID, u.Database.RefreshTokenPrivateKey)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "fail",
					"message": err.Error()},
				)

			}

			review.SetCookie(c, "access_token", accessToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
			review.SetCookie(c, "refresh_token", refreshToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
			review.SetCookie(c, "logged_in", "true", u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

			err = c.JSON(http.StatusOK, echo.Map{
				"status":       "success",
				"message":      "Login successful with user role",
				"access_token": accessToken,
			})
		}

		// Trả về thông báo login không thành công
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
}
