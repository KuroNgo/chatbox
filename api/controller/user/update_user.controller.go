package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/cloudinary"
	"chatbox/pkg/helper"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"net/http"
	"time"
)

func (u *UserController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		currentUser := c.Get("currentUser")
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "fail",
				"message": "You are not logged in!",
			})
		}
		user, err := u.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
		if err != nil || user == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status":  "Unauthorized",
				"message": "You are not authorized to perform this action!",
			})
		}

		fullName := c.Request().FormValue("full_name")
		phone := c.Request().FormValue("phone")

		file, err := c.FormFile("file")
		if err != nil {
			userResponse := domain.User{
				ID:        user.ID,
				FullName:  fullName,
				Phone:     phone,
				Role:      user.Role,
				UpdatedAt: time.Now(),
			}

			err = u.UserUseCase.Update(ctx, &userResponse)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"status":  "error",
					"message": err.Error(),
				})
			}

			return c.JSON(http.StatusOK, echo.Map{
				"status": "success",
			})
		}

		if !helper.IsImage(file.Filename) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid file format. Only images are allowed.",
			})
		}

		// Mở file để đọc dữ liệu
		f, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Error opening uploaded file",
			})
		}
		defer func(f multipart.File) {
			err = f.Close()
			if err != nil {
				return
			}
		}(f)

		imageURL, err := cloudinary.UploadImageToCloudinary(f, file.Filename, u.Database.CloudinaryUploadFolderUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": err.Error(),
			})
		}

		resultString, err := json.Marshal(user)
		userResponse := domain.User{
			FullName:  fullName,
			Phone:     phone,
			Role:      user.Role,
			AvatarURL: imageURL.ImageURL,
			AssetID:   imageURL.AssetID,
			UpdatedAt: time.Now(),
		}

		err = u.UserUseCase.Update(ctx, &userResponse)
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": string(resultString) + "the user belonging to this token no logger exists",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  "success",
			"message": "Updated user",
		})
	}

}
