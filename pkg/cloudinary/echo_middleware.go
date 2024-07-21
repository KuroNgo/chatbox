package cloudinary

import (
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"net/http"
)

func FileUploadMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file, header, err := c.Request().FormFile("files")
			if err != nil {
				err := c.JSON(http.StatusBadRequest, echo.Map{
					"error": err.Error(),
				})
				if err != nil {
					return err
				}
			}
			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					return
				}
			}(file) // close file properly

			c.Set("filePath", header.Filename)
			c.Set("file", file)

			return next(c)
		}
	}
}
