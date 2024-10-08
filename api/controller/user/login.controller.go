package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/jwt"
	"chatbox/pkg/review"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoginUser godoc
// @Summary Đăng nhập người dùng
// @Description API này cho phép người dùng đăng nhập vào hệ thống bằng địa chỉ email và mật khẩu của họ. Nếu đăng nhập thành công, hệ thống sẽ tạo và trả về access token và refresh token.
// @Tags User
// @Accept json
// @Produce json
// @Param body body domain.SignIn true "Thông tin đăng nhập của người dùng"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/user/login [post]
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

			gin_fake.SetCookie(c, "access_token", accessToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
			gin_fake.SetCookie(c, "refresh_token", refreshToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
			gin_fake.SetCookie(c, "logged_in", "true", u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

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
