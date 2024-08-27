package router

import (
	"chatbox/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func SwaggerRouter(group *echo.Group) {
	router := group.Group("")

	docs.SwaggerInfo.BasePath = ""

	handler := echoSwagger.EchoWrapHandler(
		echoSwagger.URL("http://localhost:8080/swagger/doc.json"), // Đường dẫn đến tài liệu Swagger
	)

	echoSwagger.WrapHandler = handler
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	// Thực hiện tự động chuyển hướng khi người dùng truy cập vào root "/"
	router.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusFound, "/swagger/index.html")
	})
}
