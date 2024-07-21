package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// CORSPublic sets CORS headers for specific origins
func CORSPublic() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin == "http://localhost:3000" || origin == "http://localhost:5173" {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length,Accept-Encoding,"+
					"X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With,Access-Control-Allow-Origin,access-control-allow-methods")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")

				if c.Request().Method == http.MethodOptions {
					return c.NoContent(http.StatusNoContent)
				}
			}
			return next(c)
		}
	}
}

func CORSPrivate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := "http://localhost:3000"
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With,Access-Control-Allow-Origin,access-control-allow-methods")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}
}

func OptionsMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "http://localhost:5173" {
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With,Access-Control-Allow-Origin,access-control-allow-methods")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTION")

			if c.Request().Method == http.MethodOptions {
				err := c.NoContent(http.StatusNoContent)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
