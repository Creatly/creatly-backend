# Creatly backend application

## Build & Run (Locally)
### Prerequisites
- go 1.15
- docker
- golangci-lint (<i>optional</i>, used to run code checks)

Create .env file in root directory and add following values:
```dotenv
APP_ENV=local

MONGO_URI=mongodb://mongodb:27017
MONGO_USER=admin
MONGO_PASS=qwerty

PASSWORD_SALT=<random string>
JWT_SIGNING_KEY=<random string>

SENDPULSE_LISTID=
SENDPULSE_ID=
SENDPULSE_SECRET=

HTTP_HOST=localhost

FONDY_MERCHANT_ID=1396424
FONDY_MERCHANT_PASS=test
PAYMENT_CALLBACK_URL=<host>/api/v1/callback/fondy
PAYMENT_REDIRECT_URL=https://example.com/

SMTP_PASSWORD=<password>

STORAGE_ENDPOINT=
STORAGE_BUCKET=
STORAGE_ACCESS_KEY=
STORAGE_SECRET_KEY=

CLOUDFLARE_API_KEY=
CLOUDFLARE_EMAIL=
CLOUDFLARE_ZONE_EMAIL=
CLOUDFLARE_CNAME_TARGET=
```

Use `make run` to build&run project, `make lint` to check code with linter.