package configs

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	conf *Configuration
	once sync.Once
)

func Config() *Configuration {
	once.Do(func() {
		conf = load()
	})
	return conf
}

type Configuration struct {
	HTTPPort    string
	LogLevel    string
	AppName     string
	Environment string

	ServerPort int
	ServerHost string

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	RedisAddr     string
	RedisPassword string

	CasbinConfigPath           string
	MiddlewareRolesPath        string
	AccessTokenDuration        time.Duration
	RefreshTokenDuration       time.Duration
	RefreshPasswdTokenDuration time.Duration

	MaxFileSizeInMBs   int64
	MediaServiceKey    string
	MediaServiceUrl    string
	MediaServiceSecret string

	BucketName  string
	Credentials string

	// context timeout in seconds

	JWTSecretKey string
}

func load() *Configuration {

	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading env file: ", err)
	}

	var config Configuration

	v := viper.New()
	v.AutomaticEnv()

	// v.SetDefault("CASBIN_CONFIG_PATH", "db/rbac_model.conf")
	// v.SetDefault("MIDDLEWARE_ROLES_PATH", "db/models.csv")
	// v.SetDefault("CREDENTIALS", "db/credentials.json")

	// v.SetDefault("MEDIA_SERVICE_URL", "https://media.lavina.tech/media")
	// v.SetDefault("MEDIA_SERVICE_SECRET", "some secret key")
	// v.SetDefault("MEDIA_SERVICE_KEY", "key")

	config.Environment = v.GetString("ENVIRONMENT")
	config.HTTPPort = v.GetString("HTTP_PORT")

	config.LogLevel = v.GetString("LOG_LEVEL")

	config.PostgresDatabase = v.GetString("POSTGRES_DB")
	config.PostgresUser = v.GetString("POSTGRES_USER")
	config.PostgresPassword = v.GetString("POSTGRES_PASSWORD")
	config.PostgresHost = v.GetString("POSTGRES_HOST")
	config.PostgresPort = v.GetInt("POSTGRES_PORT")
	config.RedisAddr = v.GetString("REDIS_ADDR")
	config.RedisPassword = v.GetString("REDIS_PASSWORD")

	// config.CasbinConfigPath = v.GetString("CASBIN_CONFIG_PATH")
	// config.MiddlewareRolesPath = v.GetString("MIDDLEWARE_ROLES_PATH")
	config.JWTSecretKey = v.GetString("JWT_SECRET_KEY")
	config.Credentials = v.GetString("CREDENTIALS")

	// config.MediaServiceSecret = v.GetString("MEDIA_SERVICE_SECRET")
	// config.MediaServiceKey = v.GetString("MEDIA_SERVICE_KEY")
	// config.MediaServiceUrl = v.GetString("MEDIA_SERVICE_URL")
	// config.MaxFileSizeInMBs = v.GetInt64("MAX_FILE_SIZE_MB")

	// config.CodeToIgnore = v.GetString("CODE_TO_IGNORE") //used for testing purpose
	// config.PhoneToIgnore = "+998900000000" //used for testing purpose

	//validate the configuration
	err = config.validate()
	if err != nil {
		log.Fatal("error validating config: ", err)
	}

	return &config
}

func (c *Configuration) validate() error {
	if c.HTTPPort == "" {
		return errors.New("http_port required")
	}
	if c.PostgresDatabase == "" {
		return errors.New("PostgresDatabase required")
	}
	if c.PostgresUser == "" {
		return errors.New("PostgresUser required")
	}
	if c.PostgresPassword == "" {
		return errors.New("PostgresPassword required")
	}
	// ....

	return nil
}
