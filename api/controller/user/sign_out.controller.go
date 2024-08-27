package user_controller

import (
	gin_fake "chatbox/pkg/review"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LogoutUser godoc
// @Summary Đăng xuất người dùng
// @Description API này xóa các cookie liên quan đến phiên làm việc của người dùng và trả về phản hồi thành công.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/user/logout [get]
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

		gin_fake.SetCookie(c, "access_token", "", -1, "/", "localhost", false, true)
		gin_fake.SetCookie(c, "refresh_token", "", -1, "/", "localhost", false, true)
		gin_fake.SetCookie(c, "logged_in", "", -1, "/", "localhost", false, false)

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})

	}
}
