package user_controller

import (
	"chatbox/domain"
	"chatbox/pkg/jwt"
	oauth2_google "chatbox/pkg/oauth2"
	"chatbox/pkg/review"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

func (u *UserController) GoogleLoginWithUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var googleOauthConfig = &oauth2.Config{}
		googleOauthConfig = &oauth2.Config{
			ClientID:     u.Database.GoogleClientID,
			ClientSecret: u.Database.GoogleClientSecret,
			RedirectURL:  u.Database.GoogleOAuthRedirectUrl,
			Scopes:       []string{"profile", "email"}, // Adjust scopes as needed
			Endpoint:     google.Endpoint,
		}

		code := c.QueryParam("code")
		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			fmt.Println("Error exchanging code: " + err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		userInfo, err := oauth2_google.GetUserInfo(token.AccessToken)
		if err != nil {
			fmt.Println("Error getting user info: " + err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})

		}

		// Giả sử userInfo là một map[string]interface{}
		email := userInfo["email"].(string)
		fullName := userInfo["name"].(string)
		avatarURL := userInfo["picture"].(string)
		verifiedEmail := userInfo["verified_email"].(bool)
		resBody := &domain.User{
			ID:        primitive.NewObjectID(),
			Email:     email,
			FullName:  fullName,
			AvatarURL: avatarURL,
			Provider:  "google",
			Role:      "user",
			Verified:  verifiedEmail,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updateUser, err := u.UserUseCase.UpsertUser(ctx, resBody.Email, resBody)
		if err != nil {
			return c.JSON(http.StatusBadGateway, echo.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}

		signedToken, err := oauth2_google.SignJWT(userInfo)
		if err != nil {
			fmt.Println("Error signing token: " + err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		accessTokenCh := make(chan string)
		refreshTokenCh := make(chan string)

		go func() {
			defer close(accessTokenCh)
			// Generate token
			accessToken, err := jwt.CreateToken(u.Database.AccessTokenExpiresIn, updateUser.ID, u.Database.AccessTokenPrivateKey)
			if err != nil {
				err := c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "fail",
					"message": err.Error()},
				)
				if err != nil {
					return
				}
				return
			}
			accessTokenCh <- accessToken
		}()

		go func() {
			defer close(refreshTokenCh)
			refreshToken, err := jwt.CreateToken(u.Database.RefreshTokenExpiresIn, updateUser.ID, u.Database.RefreshTokenPrivateKey)
			if err != nil {
				err := c.JSON(http.StatusBadRequest, echo.Map{
					"status":  "fail",
					"message": err.Error()},
				)
				if err != nil {
					return
				}
				return
			}
			refreshTokenCh <- refreshToken
		}()

		accessToken := <-accessTokenCh
		refreshToken := <-refreshTokenCh

		review.SetCookie(c, "access_token", accessToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		review.SetCookie(c, "refresh_token", refreshToken, u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		review.SetCookie(c, "logged_in", "true", u.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

		return c.JSON(http.StatusOK, echo.Map{"token": signedToken})
	}
}
