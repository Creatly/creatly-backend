package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
)

type (
	Config struct {
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		FileStorage FileStorageConfig
		Email       EmailConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	FileStorageConfig struct {
		URL    string `mapstructure:"url"`
		Bucket string `mapstructure:"bucket"`
	}

	EmailConfig struct {
		ListID       string
		ClientID     string
		ClientSecret string
	}

	HTTPConfig struct {
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

// Init populates Config struct with values from config file
//located at filepath and environment variables
func Init(path string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := parseEnv(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("cache.ttl", &cfg.CacheTTL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("fileStorage", &cfg.FileStorage); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Mongo.URI = viper.GetString("uri")
	cfg.Mongo.User = viper.GetString("user")
	cfg.Mongo.Password = viper.GetString("pass")

	cfg.Auth.PasswordSalt = viper.GetString("salt")
	cfg.Auth.JWT.SigningKey = viper.GetString("signing_key")

	cfg.Email.ClientSecret = viper.GetString("secret")
	cfg.Email.ClientID = viper.GetString("id")
	cfg.Email.ListID = viper.GetString("listid")
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name

	return viper.ReadInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
}

func parseEnv() error {
	if err := parseMongoEnvVariables(); err != nil {
		return err
	}

	if err := parseJWTEnvVariables(); err != nil {
		return err
	}

	if err := parseSendPulseEnvVariables(); err != nil {
		return err
	}

	return parsePasswordEnvVariables()
}

func parseMongoEnvVariables() error {
	viper.SetEnvPrefix("mongo")
	if err := viper.BindEnv("uri"); err != nil {
		return err
	}

	if err := viper.BindEnv("user"); err != nil {
		return err
	}

	return viper.BindEnv("pass")
}

func parseSendPulseEnvVariables() error {
	viper.SetEnvPrefix("sendpulse")
	if err := viper.BindEnv("listid"); err != nil {
		return err
	}

	if err := viper.BindEnv("id"); err != nil {
		return err
	}

	return viper.BindEnv("secret")
}

func parsePasswordEnvVariables() error {
	viper.SetEnvPrefix("password")
	return viper.BindEnv("salt")
}

func parseJWTEnvVariables() error {
	viper.SetEnvPrefix("jwt")
	return viper.BindEnv("signing_key")
}
