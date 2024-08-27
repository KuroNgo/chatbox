package message_controller

import (
	"chatbox/domain"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

var upgrade = websocket.Upgrader{}

// Setup godoc
// @Summary Thiết lập kết nối WebSocket
// @Description API này thiết lập kết nối WebSocket cho việc gửi và nhận tin nhắn trong một phòng cụ thể.
// @Tags Message
// @Accept json
// @Produce json
// @Param room_id query string true "ID của phòng" example("605c72ef1f1b2c001f9b22a2")
// @Param to_user_id query string true "ID của người nhận" example("605c72ef1f1b2c001f9b22a3")
// @Success 101 {object} map[string]interface{} "Kết nối WebSocket thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/message/ws [get]
func (m *MessageController) Setup() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		//currentUser := c.Get("currentUser")
		//if currentUser == nil {
		//	return c.JSON(http.StatusUnauthorized, echo.Map{
		//		"status":  "fail",
		//		"message": "You are not login!",
		//	})
		//}
		//
		//user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		//if err != nil || user == nil {
		//	return c.JSON(http.StatusUnauthorized, echo.Map{
		//		"status":  "fail",
		//		"message": "You are not authorize in perform this action",
		//	})
		//}

		roomID := c.QueryParam("room_id")
		idRoom, _ := primitive.ObjectIDFromHex(roomID)
		toUserID := c.QueryParam("to_user_id")
		idToUser, _ := primitive.ObjectIDFromHex(toUserID)

		upgrade = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Cho phép tất cả các nguồn gốc hoặc điều chỉnh theo yêu cầu bảo mật của bạn
			},
		}

		ws, err := upgrade.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				return
			}
		}(ws)

		for {

			mt, message, err := ws.ReadMessage()
			if err != nil {
				break
			}

			data := domain.Message{
				ID:     primitive.NewObjectID(),
				RoomID: idRoom,
				//UserID:    user.ID,
				ToUserID:  idToUser,
				Text:      string(message),
				TimeStamp: time.Now(),
				Color:     "black",
			}

			if err = m.MessageUseCase.CreateOne(ctx, data); err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"status":  "fail",
					"message": "Internal Server Error: " + err.Error(),
				})
			}

			err = ws.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
		return nil
	}
}
