package room_controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// FetchManyRoom godoc
// @Summary Lấy danh sách phòng
// @Description API này trả về danh sách tất cả các phòng mà người dùng hiện tại có quyền truy cập.
// @Tags Room
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/fetch [get]
func (r *RoomController) FetchManyRoom() echo.HandlerFunc {
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

		data, err := r.RoomUseCase.FetchManyRoom(ctx, user.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "fail",
				"message": "server internal error",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
			"data":   data,
		})
	}
}

// FetchOneRoom godoc
// @Summary Lấy thông tin phòng theo ID
// @Description API này trả về thông tin chi tiết của một phòng dựa trên ID phòng.
// @Tags Room
// @Accept json
// @Produce json
// @Param _id query string true "ID của phòng cần lấy thông tin" example("605c72ef1f1b2c001f9b22a2")
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/1/fetch [get]
func (r *RoomController) FetchOneRoom() echo.HandlerFunc {
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
		data, err := r.RoomUseCase.FetchOneRoom(ctx, user.ID, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "fail",
				"message": "server internal error",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
			"data":   data,
		})
	}
}

// FetchOneByName godoc
// @Summary Lấy thông tin phòng theo tên
// @Description API này trả về thông tin chi tiết của một phòng dựa trên tên phòng.
// @Tags Room
// @Accept json
// @Produce json
// @Param name query string true "Tên của phòng cần lấy thông tin" example("Room1")
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/room/fetch/name [get]
func (r *RoomController) FetchOneByName() echo.HandlerFunc {
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

		name := c.QueryParam("name")
		data, err := r.RoomUseCase.FetchOneByName(ctx, user.ID, name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "fail",
				"message": "server internal error",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
			"data":   data,
		})
	}
}
