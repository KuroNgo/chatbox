package log_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
)

// FetchManyActivity godoc
// @Summary Lấy danh sách hoạt động
// @Description API này trả về danh sách hoạt động dựa trên trang yêu cầu.
// @Tags Activity
// @Accept json
// @Produce json
// @Param page query string true "Trang cần lấy dữ liệu" default(1)
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Lỗi hệ thống"
// @Security ApiKeyAuth
// @Router /api/activity/fetch [get]
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
			page = "1" // Giá trị mặc định
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
