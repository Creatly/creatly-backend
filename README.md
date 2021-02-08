# Backend app for Online Course Constructor Platform

## Build & Run (Locally)
### Prerequisites
- go 1.15
- docker
- golangci-lint (<i>optional</i>, used to run code checks)

Create .env file in root directory and add following values:
```dotenv
MONGO_URI=mongodb://mongodb:27017
MONGO_USER=admin
MONGO_PASS=qwerty
PASSWORD_SALT=<random string>
JWT_SIGNING_KEY=<random string>
SENDPULSE_LISTID=<list id>
SENDPULSE_ID=<client id>
SENDPULSE_SECRET=<client secret>
HTTP_HOST=localhost
FONDY_MERCHANT_ID=1396424
FONDY_MERCHANT_PASS=test
PAYMENT_CALLBACK_URL=<host>/api/v1/callback/fondy
PAYMENT_REDIRECT_URL=https://example.com/
```

Use `make run` to build&run project, `make lint` to check code with linter.

# TODO: Describe project architecture