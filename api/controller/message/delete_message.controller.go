package message_controller

import (
	constant "chatbox/pkg/const"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (m *MessageController) DeleteOne() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": constant.Unauthorized,
			})
		}

		ctx := c.Request().Context()
		user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": constant.Unauthenticated,
			})
		}

		id := c.QueryParam("_id")

		err = m.MessageUseCase.DeleteOne(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": constant.ServerInternalError + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}
