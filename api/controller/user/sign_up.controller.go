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

// SignUp godoc
// @Summary Đăng ký người dùng mới
// @Description API này cho phép người dùng mới đăng ký tài khoản, bao gồm thông tin cá nhân, mật khẩu và ảnh đại diện.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Địa chỉ email của người dùng"
// @Param full_name formData string true "Tên đầy đủ của người dùng"
// @Param password formData string true "Mật khẩu của người dùng"
// @Param avatar_url formData string false "URL của ảnh đại diện người dùng"
// @Param phone formData string false "Số điện thoại của người dùng"
// @Param file formData file false "Ảnh đại diện của người dùng"
// @Success 200 {object} map[string]interface{} "Thành công"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/user/signup [post]
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
				if newUser.Verified == false {
					user, errE := u.UserUseCase.GetByEmail(ctx, newUser.Email)
					if errE != nil {
						return c.JSON(http.StatusBadRequest, echo.Map{
							"status":  "error",
							"message": err.Error()},
						)
					}

					if errDel := u.UserUseCase.Delete(ctx, user.ID.Hex()); errDel != nil {
						return c.JSON(http.StatusBadRequest, echo.Map{
							"status":  "error",
							"message": err.Error()},
						)
					}
				}

				errCreate := u.UserUseCase.Create(ctx, newUser)
				if errCreate != nil {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"status":  "error",
						"message": err.Error()},
					)
				}
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
				Subject:   "Your account verification code " + code,
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
