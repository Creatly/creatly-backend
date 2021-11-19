package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	type env struct {
		mongoURI              string
		mongoUser             string
		mongoPass             string
		passwordSalt          string
		jwtSigningKey         string
		host                  string
		fondyCallbackURL      string
		frontendUrl           string
		smtpPassword          string
		appEnv                string
		storageEndpoint       string
		storageBucket         string
		storageAccessKey      string
		storageSecretKey      string
		cloudflareApiKey      string
		cloudflareEmail       string
		cloudflareZoneEmail   string
		cloudflareCnameTarget string
	}

	type args struct {
		path string
		env  env
	}

	setEnv := func(env env) {
		os.Setenv("MONGO_URI", env.mongoURI)
		os.Setenv("MONGO_USER", env.mongoUser)
		os.Setenv("MONGO_PASS", env.mongoPass)
		os.Setenv("PASSWORD_SALT", env.passwordSalt)
		os.Setenv("JWT_SIGNING_KEY", env.jwtSigningKey)
		os.Setenv("HTTP_HOST", env.host)
		os.Setenv("FONDY_CALLBACK_URL", env.fondyCallbackURL)
		os.Setenv("FRONTEND_URL", env.frontendUrl)
		os.Setenv("SMTP_PASSWORD", env.smtpPassword)
		os.Setenv("APP_ENV", env.appEnv)
		os.Setenv("STORAGE_ENDPOINT", env.storageEndpoint)
		os.Setenv("STORAGE_BUCKET", env.storageBucket)
		os.Setenv("STORAGE_ACCESS_KEY", env.storageAccessKey)
		os.Setenv("STORAGE_SECRET_KEY", env.storageSecretKey)
		os.Setenv("CLOUDFLARE_API_KEY", env.cloudflareApiKey)
		os.Setenv("CLOUDFLARE_EMAIL", env.cloudflareEmail)
		os.Setenv("CLOUDFLARE_ZONE_EMAIL", env.cloudflareZoneEmail)
		os.Setenv("CLOUDFLARE_CNAME_TARGET", env.cloudflareCnameTarget)
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "test config",
			args: args{
				path: "fixtures",
				env: env{
					mongoURI:              "mongodb://localhost:27017",
					mongoUser:             "admin",
					mongoPass:             "qwerty",
					passwordSalt:          "salt",
					jwtSigningKey:         "key",
					host:                  "localhost",
					fondyCallbackURL:      "https://zhashkevych.com/callback",
					frontendUrl:           "http://localhost:1337",
					smtpPassword:          "qwerty123",
					appEnv:                "local",
					storageEndpoint:       "test.filestorage.com",
					storageBucket:         "test",
					storageAccessKey:      "qwerty123",
					storageSecretKey:      "qwerty123",
					cloudflareApiKey:      "api_key",
					cloudflareEmail:       "email",
					cloudflareZoneEmail:   "zone_email",
					cloudflareCnameTarget: "cname_target",
				},
			},
			want: &Config{
				Environment: "local",
				CacheTTL:    time.Second * 3600,
				HTTP: HTTPConfig{
					Host:               "localhost",
					MaxHeaderMegabytes: 1,
					Port:               "80",
					ReadTimeout:        time.Second * 10,
					WriteTimeout:       time.Second * 10,
				},
				Auth: AuthConfig{
					PasswordSalt: "salt",
					JWT: JWTConfig{
						RefreshTokenTTL: time.Minute * 30,
						AccessTokenTTL:  time.Minute * 15,
						SigningKey:      "key",
					},
					VerificationCodeLength: 10,
				},
				Mongo: MongoConfig{
					Name:     "testDatabase",
					URI:      "mongodb://localhost:27017",
					User:     "admin",
					Password: "qwerty",
				},
				FileStorage: FileStorageConfig{
					Endpoint:  "test.filestorage.com",
					Bucket:    "test",
					AccessKey: "qwerty123",
					SecretKey: "qwerty123",
				},
				Email: EmailConfig{
					Templates: EmailTemplates{
						Verification:       "./templates/verification_email.html",
						PurchaseSuccessful: "./templates/purchase_successful.html",
					},
					Subjects: EmailSubjects{
						Verification:       "Спасибо за регистрацию, %s!",
						PurchaseSuccessful: "Покупка прошла успешно!",
					},
				},
				Payment: PaymentConfig{
					FondyCallbackURL: "https://zhashkevych.com/callback",
				},
				Limiter: LimiterConfig{
					RPS:   10,
					Burst: 2,
					TTL:   time.Minute * 10,
				},
				SMTP: SMTPConfig{
					Host: "mail.privateemail.com",
					Port: 587,
					From: "maksim@zhashkevych.com",
					Pass: "qwerty123",
				},
				Cloudflare: CloudflareConfig{
					ApiKey:      "api_key",
					Email:       "email",
					CnameTarget: "cname_target",
					ZoneEmail:   "zone_email",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.args.env)

			got, err := Init(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() got = %v, want %v", got, tt.want)
			}
		})
	}
}
