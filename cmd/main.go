package main

import (
	"chatbox/api/router"
	"chatbox/infrastructor"
	"github.com/labstack/echo/v4"
	"time"
)

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
