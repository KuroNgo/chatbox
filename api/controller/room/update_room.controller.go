package room_controller

import (
	"chatbox/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
