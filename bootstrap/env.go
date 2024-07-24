package bootstrap

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Database struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	// implement the Cloudinary
	CloudinaryCloudName        string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryAPIKey           string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryAPISecret        string `mapstructure:"CLOUDINARY_API_SECRET"`
	CloudinaryUploadFolderUser string `mapstructure:"CLOUDINARY_UPLOAD_FOLDER_AUDIO_VOCABULARY"`

	// implement the Google Oauth
	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
}

func NewEnv() *Database {
	env := Database{}
	viper.SetConfigFile("./config/app.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Can't find the app.env:", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalln("Environment can't be loaded: ", err)
	}

	return &env
}
