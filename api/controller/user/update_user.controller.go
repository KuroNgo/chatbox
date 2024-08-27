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

// UpdateUser godoc
// @Summary Cập nhật thông tin người dùng
// @Description API này cho phép người dùng đã đăng nhập cập nhật thông tin của mình, bao gồm tên đầy đủ, số điện thoại và ảnh đại diện.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param full_name formData string true "Tên đầy đủ của người dùng"
// @Param phone formData string true "Số điện thoại của người dùng"
// @Param file formData file false "Ảnh đại diện của người dùng"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/user/update [put]
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
