package room_controller

import (
	"chatbox/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UpdateRoom godoc
// @Summary Cập nhật thông tin phòng
// @Description API này cập nhật thông tin của một phòng dựa trên dữ liệu đầu vào được cung cấp.
// @Tags Room
// @Accept json
// @Produce json
// @Param room body domain.Input true "Dữ liệu phòng cần cập nhật"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/update [put]
func (r *RoomController) UpdateRoom() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not login!",
			})
		}

		ctx := c.Request().Context()
		user, err := r.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not authorize in perform this action",
			})
		}

		var roomInput domain.Input
		if err = c.Bind(&roomInput); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "fail",
				"message": "Bad request: " + err.Error(),
			})
		}

		data, err := r.RoomUseCase.GetByName(ctx, user.ID, roomInput.Name)
		if err = c.Bind(&roomInput); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "fail",
				"message": "Bad request: " + err.Error(),
			})
		}

		err = r.RoomUseCase.UpdateRoom(ctx, user.ID, data)
		if err = c.Bind(&roomInput); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "fail",
				"message": "Bad request: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}
