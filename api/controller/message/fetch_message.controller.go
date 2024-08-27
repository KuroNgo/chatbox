package message_controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// FetchMany godoc
// @Summary Lấy nhiều tin nhắn
// @Description API này trả về nhiều tin nhắn từ một phòng dựa trên ID phòng và ID người dùng.
// @Tags Message
// @Accept json
// @Produce json
// @Param room_id query string true "ID của phòng cần lấy tin nhắn"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/message/fetch [get]
func (m *MessageController) FetchMany() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not login!",
			})
		}

		ctx := c.Request().Context()
		user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not authorize in perform this action",
			})
		}

		roomIDStr := c.QueryParam("room_id")
		if roomIDStr == "" {
			return c.JSON(http.StatusBadRequest, "room_id is required")
		}

		data, err := m.MessageUseCase.FetchMany(ctx, user.ID, roomIDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Internal Server Error:" + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": data,
		})
	}
}
