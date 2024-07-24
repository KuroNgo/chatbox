package middlewares

import (
	"chatbox/infrastructor"
	"chatbox/pkg/jwt"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func DeserializeUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var accessToken string
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return err
			}

			authorizationHeader := c.Request().Header.Get("Authorization")
			fields := strings.Fields(authorizationHeader)

			if len(fields) != 0 && fields[0] == "Bearer" {
				accessToken = fields[1]
			} else if err == nil {
				accessToken = cookie.Value
			}

			if accessToken == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"status":  "fail",
					"message": "You are not logged in",
				})
			}

			app := infrastructor.App()
			env := app.Env

			sub, err := jwt.ValidateToken(accessToken, env.AccessTokenPublicKey)
			if err != nil {
				fmt.Println("The err is: ", err)
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"status":  "fail",
					"message": err.Error(),
				})
			}

			c.Set("currentUser", sub)
			return next(c)
		}
	}
}
