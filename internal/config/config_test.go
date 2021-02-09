package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	type env struct {
		mongoURI           string
		mongoUser          string
		mongoPass          string
		passwordSalt       string
		jwtSigningKey      string
		sendpulseListId    string
		sendpulseId        string
		sendpulseSecret    string
		host               string
		fondyMerchantId    string
		fondyMerchantPass  string
		paymentCallbackURL string
		paymentResponseURL string
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
		os.Setenv("SENDPULSE_LISTID", env.sendpulseListId)
		os.Setenv("SENDPULSE_ID", env.sendpulseId)
		os.Setenv("SENDPULSE_SECRET", env.sendpulseSecret)
		os.Setenv("HTTP_HOST", env.host)
		os.Setenv("FONDY_MERCHANT_ID", env.fondyMerchantId)
		os.Setenv("FONDY_MERCHANT_PASS", env.fondyMerchantPass)
		os.Setenv("PAYMENT_CALLBACK_URL", env.paymentCallbackURL)
		os.Setenv("PAYMENT_RESPONSE_URL", env.paymentResponseURL)
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
				path: "fixtures/test",
				env: env{
					mongoURI:           "mongodb://localhost:27017",
					mongoUser:          "admin",
					mongoPass:          "qwerty",
					passwordSalt:       "salt",
					jwtSigningKey:      "key",
					sendpulseSecret:    "secret",
					sendpulseId:        "id",
					sendpulseListId:    "listId",
					host:               "localhost",
					fondyMerchantId:    "123",
					fondyMerchantPass:  "fondy",
					paymentResponseURL: "https://zhashkevych.com/",
					paymentCallbackURL: "https://zhashkevych.com/callback",
				}},
			want: &Config{
				CacheTTL: time.Second * 3600,
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
				},
				Mongo: MongoConfig{
					Name:     "testDatabase",
					URI:      "mongodb://localhost:27017",
					User:     "admin",
					Password: "qwerty",
				},
				FileStorage: FileStorageConfig{
					URL:    "test.filestorage.com",
					Bucket: "test",
				},
				Email: EmailConfig{
					ListID:       "listId",
					ClientID:     "id",
					ClientSecret: "secret",
				},
				Payment: PaymentConfig{
					Fondy: FondyConfig{
						MerchantId:       "123",
						MerchantPassword: "fondy",
					},
					CallbackURL: "https://zhashkevych.com/callback",
					ResponseURL: "https://zhashkevych.com/",
				},
				Limiter: LimiterConfig{
					RPS:   10,
					Burst: 2,
					TTL:   time.Minute * 10,
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
