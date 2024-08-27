package main

import (
	"chatbox/api/router"
	"chatbox/infrastructor"
	"github.com/labstack/echo/v4"
	"time"
)

// @title Chatbox with Echo
// @version 1.0
// @description This is a Chathox server for Echo.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	app := infrastructor.App()
	env := app.Env
	db := app.MongoDB.Database(env.DBName)
	defer app.CloseDBConnection()
	timeout := time.Duration(env.ContextTimeout) * time.Second

	_echo := echo.New()
	router.SetUp(env, timeout, db, _echo)
	err := _echo.Start(env.ServerAddress)
	if err != nil {
		return
	}
}
