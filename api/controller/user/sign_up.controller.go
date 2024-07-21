package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/cloudinary"
	"chatbox/pkg/helper"
	mailk "chatbox/pkg/mail"
	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"time"
)

func (u *UserController) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		email := c.Request().FormValue("email")
		fullName := c.Request().FormValue("full_name")
		password := c.Request().FormValue("password")
		avatarUrl := c.Request().FormValue("avatar_url")
		phone := c.Request().FormValue("phone")

		if !helper.EmailValid(email) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": "Email Invalid !",
			})
		}

		// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
		if !helper.PasswordStrong(password) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "error",
				"message": "Password must have at least 8 characters, " +
					"including uppercase letters, lowercase letters and numbers!",
			})
		}

		// Băm mật khẩu
		hashedPassword, err := helper.HashPassword(password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		password = hashedPassword
		password = helper.Santize(password)
		email = helper.Santize(email)
		file, err := c.FormFile("file")
		if err != nil {
			newUser := &domain.User{
				ID:        primitive.NewObjectID(),
				FullName:  fullName,
				AvatarURL: avatarUrl,
				Email:     email,
				Password:  hashedPassword,
				Verified:  false,
				Provider:  "fe-it",
				Role:      "user",
				Phone:     phone,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// thực hiện đăng ký người dùng
			err = u.UserUseCase.Create(ctx, newUser)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "error",
					"message": err.Error()},
				)
			}

			var code string
			for {
				code = randstr.Dec(6)
				if u.UserUseCase.UniqueVerificationCode(ctx, code) {
					break
				}
			}

			updUser := domain.User{
				ID:               newUser.ID,
				VerificationCode: code,
				Verified:         false,
				UpdatedAt:        time.Now(),
			}

			// Update User in Database
			_, err = u.UserUseCase.UpdateVerify(ctx, &updUser)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "error",
					"message": err.Error()},
				)
			}

			emailData := mailk.EmailData{
				Code:      code,
				FirstName: newUser.FullName,
				Subject:   "Your account verification code",
			}

			err = mailk.SendEmail(&emailData, newUser.Email, "sign_up.html")
			if err != nil {
				return c.JSON(http.StatusBadGateway, echo.Map{
					"status":  "fail",
					"message": "There was an error sending email" + err.Error(),
				})
			}

			// Trả về phản hồi thành công
			return c.JSON(http.StatusOK, echo.Map{
				"status":  "success",
				"message": "We sent an email with a verification code to your email",
			})
		}

		// Kiểm tra xem file có phải là hình ảnh không
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
			err := f.Close()
			if err != nil {

			}
		}(f)

		imageURL, err := cloudinary.UploadImageToCloudinary(f, file.Filename, u.Database.CloudinaryUploadFolderUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": err.Error(),
			})
		}
		newUser := domain.User{
			ID:        primitive.NewObjectID(),
			FullName:  fullName,
			AvatarURL: imageURL.ImageURL,
			AssetID:   imageURL.AssetID,
			Email:     email,
			Password:  hashedPassword,
			Verified:  false,
			Provider:  "fe-it",
			Role:      "user",
			Phone:     phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// thực hiện đăng ký người dùng
		err = u.UserUseCase.Create(ctx, &newUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		// Trả về phản hồi thành công
		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}

func (u *UserController) VerificationCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var verificationCode domain.VerificationCode
		if err := c.Bind(&verificationCode); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": err.Error()},
			)
		}

		user, err := u.UserUseCase.GetByVerificationCode(ctx, verificationCode.VerificationCode)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		res := u.UserUseCase.CheckVerify(ctx, verificationCode.VerificationCode)
		if res != true {
			return c.JSON(http.StatusNotModified, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		updUser := domain.User{
			ID:        user.ID,
			Verified:  true,
			UpdatedAt: time.Now(),
		}

		// Update User in Database
		_, err = u.UserUseCase.UpdateVerify(ctx, &updUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		// Trả về phản hồi thành công
		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}
