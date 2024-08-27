package room_controller

import (
	"chatbox/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// CreateRoom godoc
// @Summary Tạo phòng mới
// @Description API này tạo một phòng mới với thông tin phòng từ yêu cầu.
// @Tags Room
// @Accept json
// @Produce json
// @Param room body domain.Input true "Thông tin phòng cần tạo"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/create [post]
func (r *RoomController) CreateRoom() echo.HandlerFunc {
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

		room := domain.Room{
			ID:     primitive.NewObjectID(),
			Name:   roomInput.Name,
			UserID: user.ID,
		}

		err = r.RoomUseCase.CreateRoom(ctx, room)
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
