package configs

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	PrefixUrl string

	AppName string

	LogSavePath string
	LogSaveName string
	LogFileExt  string

	MaxLogFiles int

	ImageStaticPath string
	ImageSavePath   string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MySql struct {
	ConnString string
	Name       string

	MySqlBase MySqlBase
}

type MySqlBase struct {
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type MongoDB struct {
	ConnString string
	Name       string
}

type Jwt struct {
	Secret         string
	ExpirationDays int
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Service struct {
	UserInfo     string
	UserAuth     string
	Notification string
}

type Config struct {
	App     App
	Server  Server
	MySql   MySql
	MongoDB MongoDB
	Jwt     Jwt
	Redis   Redis
	Service Service
}

var C Config

// Setup initialize the configuration instance
func Setup() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found")
		} else {
			log.Fatalf("Config file was found but another error was produced")
		}
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	godotenv.Load()

	applyEnvVariables()
}

func applyEnvVariables() {
	C.App.PrefixUrl = viper.GetString("PREFIX_URL")

	C.Server.RunMode = viper.GetString("RUN_MODE")

	C.MySql.ConnString = viper.GetString("DB_CONNSTRING")
	C.MySql.Name = viper.GetString("DB_NAME")

	C.MongoDB.ConnString = viper.GetString("MONGO_DB_COONSTRING")
	C.MongoDB.Name = viper.GetString("MONGO_DB_NAME")

	C.Redis.Addr = viper.GetString("REDIS_ADDR")
	C.Redis.Password = viper.GetString("REDIS_PASSWORD")
	C.Redis.DB = viper.GetInt("REDIS_DB")

	C.Jwt.Secret = viper.GetString("JWT_SECRET")
	C.Jwt.ExpirationDays = viper.GetInt("JWT_EXPIRATION_DAYS")

	C.Service.UserInfo = viper.GetString("USER_INFO_SERVICE_ADDR")
	C.Service.UserAuth = viper.GetString("USER_AUTH_SERVICE_ADDR")
	C.Service.Notification = viper.GetString("NOTIFICATION_SERVICE_ADDR")
}
