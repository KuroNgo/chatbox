package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/helper"
	mailk "chatbox/pkg/mail"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/thanhpk/randstr"
	"net/http"
	"time"
)

func (u *UserController) ForgetPasswordInUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var forgetInput domain.ForgetPassword
		if err := c.Bind(&forgetInput); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		// Check if email exists
		user, err := u.UserUseCase.GetByEmail(ctx, forgetInput.Email)
		if err != nil || user == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  "error",
				"message": "User not found or error occurred",
			})
		}

		if user.Provider != "fe-it" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": fmt.Sprintf("Sorry, the %s you provided is not supported.", user.Provider),
			})
		}

		// Generate unique verification code
		var code string
		for {
			code = randstr.String(6)
			if u.UserUseCase.UniqueVerificationCode(ctx, code) {
				break
			}
		}

		updUser := domain.User{
			ID:               user.ID,
			Verified:         true,
			VerificationCode: code,
			UpdatedAt:        time.Now(),
		}

		// Update User in Database
		_, err = u.UserUseCase.UpdateVerify(ctx, &updUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		emailData := mailk.EmailData{
			Code:      code,
			FirstName: user.FullName,
			Subject:   "Khôi phục mật khẩu",
		}

		err = mailk.SendEmail(&emailData, user.Email, "forget_password.html")
		if err != nil {
			return c.JSON(http.StatusBadGateway, echo.Map{
				"status":  "error",
				"message": "There was an error sending the email",
			})
		}

		// Success response
		return c.JSON(http.StatusOK, echo.Map{
			"status":  "success",
			"message": "We sent an email with a verification code to your email",
		})
	}
}

func (u *UserController) VerificationCodeForChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var verificationCode domain.VerificationCode
		if err := c.Bind(&verificationCode); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		user, err := u.UserUseCase.GetByVerificationCode(ctx, verificationCode.VerificationCode)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		res := u.UserUseCase.CheckVerify(ctx, verificationCode.VerificationCode)
		if !res {
			return c.JSON(http.StatusNotModified, echo.Map{
				"status":  "error",
				"message": "Verification code check failed",
			})
		}

		updUser := domain.User{
			ID:        user.ID,
			Verified:  true,
			UpdatedAt: time.Now(),
		}

		// Update User in Database
		if _, err := u.UserUseCase.UpdateVerifyForChangePassword(ctx, &updUser); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error(),
			})

		}

		cookie := &http.Cookie{
			Name:     "verification_code",
			Value:    verificationCode.VerificationCode,
			MaxAge:   u.Database.AccessTokenMaxAge * 60,
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: false,
		}

		c.SetCookie(cookie)

		// Trả về phản hồi thành công
		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}

}

func (u *UserController) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		cookie, err := c.Cookie("verification_code")
		if err != nil {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status":  "fail",
				"message": "Verification code is missing!",
			})
		}

		var changePasswordInput domain.ChangePassword
		if err := c.Bind(&changePasswordInput); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)

		}

		if changePasswordInput.Password != changePasswordInput.PasswordCompare {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": "The passwords provided do not match.",
			})
		}

		// Đối với change password, không clear giá trị verification Code ở phía client và backend
		user, err := u.UserUseCase.GetByVerificationCode(ctx, cookie.Value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": "verification code do not match"},
			)
		}

		changePasswordInput.Password, err = helper.HashPassword(changePasswordInput.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		updateUser := &domain.User{
			ID:       user.ID,
			Password: changePasswordInput.Password,
		}

		err = u.UserUseCase.UpdatePassword(ctx, updateUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status":  "error",
				"message": err.Error()},
			)
		}

		cookies := new(http.Cookie)
		cookie.Name = "verification_code"
		cookie.Value = ""
		cookie.MaxAge = -1
		cookie.Path = "/"
		cookie.Domain = "localhost"
		cookie.Secure = false
		cookie.HttpOnly = false

		c.SetCookie(cookies)

		// Trả về phản hồi thành công
		return c.JSON(http.StatusOK, echo.Map{
			"status": "success",
		})
	}
}
