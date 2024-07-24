package message_controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
