package room_controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// DeleteRoom godoc
// @Summary Xóa phòng
// @Description API này xóa một phòng dựa trên ID của phòng và ID của người dùng.
// @Tags Room
// @Accept json
// @Produce json
// @Param _id query string true "ID của phòng cần xóa" example("605c72ef1f1b2c001f9b22a2")
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/delete [delete]
func (r *RoomController) DeleteRoom() echo.HandlerFunc {
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

		id := c.QueryParam("_id")
		err = r.RoomUseCase.DeleteRoom(ctx, user.ID, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "fail",
				"message": "server internal error: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}
