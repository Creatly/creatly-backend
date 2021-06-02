package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute
	defaultVerificationCodeLength = 8

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		FileStorage FileStorageConfig
		Email       EmailConfig
		Payment     PaymentConfig
		Limiter     LimiterConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
		FrontendURL string
		SMTP        SMTPConfig
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	FileStorageConfig struct {
		Endpoint  string
		Bucket    string
		AccessKey string
		SecretKey string
	}

	EmailConfig struct {
		SendPulse SendPulseConfig
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	SendPulseConfig struct {
		ListID       string
		ClientID     string
		ClientSecret string
	}

	EmailTemplates struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	EmailSubjects struct {
		Verification       string `mapstructure:"verification_email"`
		PurchaseSuccessful string `mapstructure:"purchase_successful"`
	}

	PaymentConfig struct {
		Fondy       FondyConfig
		CallbackURL string
		ResponseURL string
	}

	FondyConfig struct {
		MerchantId       string
		MerchantPassword string
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}

	SMTPConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		From string `mapstructure:"from"`
		Pass string
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseEnv(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, viper.GetString("env")); err != nil {
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

	if err := viper.UnmarshalKey("auth.verificationCodeLength", &cfg.Auth.VerificationCodeLength); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("fileStorage", &cfg.FileStorage); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("email.templates", &cfg.Email.Templates); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("email.subjects", &cfg.Email.Subjects); err != nil {
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

	cfg.Email.SendPulse.ClientSecret = viper.GetString("secret")
	cfg.Email.SendPulse.ClientID = viper.GetString("id")
	cfg.Email.SendPulse.ListID = viper.GetString("listid")

	cfg.HTTP.Host = viper.GetString("host")

	cfg.Payment.Fondy.MerchantId = viper.GetString("merchant_id")
	cfg.Payment.Fondy.MerchantPassword = viper.GetString("merchant_pass")
	cfg.Payment.CallbackURL = viper.GetString("callback_url")
	cfg.Payment.ResponseURL = viper.GetString("response_url")

	cfg.FrontendURL = viper.GetString("url")

	cfg.SMTP.Pass = viper.GetString("password")

	cfg.Environment = viper.GetString("env")

	cfg.FileStorage.Endpoint = viper.GetString("endpoint")
	cfg.FileStorage.AccessKey = viper.GetString("access_key")
	cfg.FileStorage.SecretKey = viper.GetString("secret_key")
	cfg.FileStorage.Bucket = viper.GetString("bucket")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}

func parseEnv() error {
	if err := parseMongoEnvVariables(); err != nil {
		return err
	}

	if err := parseJWTFromEnv(); err != nil {
		return err
	}

	if err := parseSendPulseEnvVariables(); err != nil {
		return err
	}

	if err := parseHostFromEnv(); err != nil {
		return err
	}

	if err := parseFondyEnvVariables(); err != nil {
		return err
	}

	if err := parsePaymentEnvVariables(); err != nil {
		return err
	}

	if err := parseFrontendHostFromEnv(); err != nil {
		return err
	}

	if err := parseSMTPPassFromEnv(); err != nil {
		return err
	}

	if err := parseAppEnvFromEnv(); err != nil {
		return err
	}

	if err := parseStorageEnvVariables(); err != nil {
		return err
	}

	return parsePasswordFromEnv()
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

func parseFondyEnvVariables() error {
	viper.SetEnvPrefix("fondy")

	if err := viper.BindEnv("merchant_id"); err != nil {
		return err
	}

	if err := viper.BindEnv("merchant_pass"); err != nil {
		return err
	}

	if err := viper.BindEnv("callback_url"); err != nil {
		return err
	}

	return viper.BindEnv("response_url")
}

func parsePaymentEnvVariables() error {
	viper.SetEnvPrefix("payment")

	if err := viper.BindEnv("callback_url"); err != nil {
		return err
	}

	return viper.BindEnv("response_url")
}

func parsePasswordFromEnv() error {
	viper.SetEnvPrefix("password")

	return viper.BindEnv("salt")
}

func parseJWTFromEnv() error {
	viper.SetEnvPrefix("jwt")

	return viper.BindEnv("signing_key")
}

func parseHostFromEnv() error {
	viper.SetEnvPrefix("http")

	return viper.BindEnv("host")
}

func parseFrontendHostFromEnv() error {
	viper.SetEnvPrefix("frontend")

	return viper.BindEnv("url")
}

func parseSMTPPassFromEnv() error {
	viper.SetEnvPrefix("smtp")

	return viper.BindEnv("password")
}

func parseAppEnvFromEnv() error {
	viper.SetEnvPrefix("app")

	return viper.BindEnv("env")
}

func parseStorageEnvVariables() error {
	viper.SetEnvPrefix("storage")

	if err := viper.BindEnv("bucket"); err != nil {
		return err
	}

	if err := viper.BindEnv("endpoint"); err != nil {
		return err
	}

	if err := viper.BindEnv("access_key"); err != nil {
		return err
	}

	return viper.BindEnv("secret_key")
}
