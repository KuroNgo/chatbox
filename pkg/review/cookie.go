package gin_fake

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Auth struct {
	Database struct {
		AccessTokenMaxAge int
	}
}

func SetCookie(c echo.Context, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.MaxAge = maxAge
	cookie.Path = path
	cookie.Domain = domain
	cookie.Secure = secure
	cookie.HttpOnly = httpOnly
	c.SetCookie(cookie)
}
