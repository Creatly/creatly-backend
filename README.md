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
PASSWORD_SALT=<any salt here>
JWT_SIGNING_KEY=<any key here>
SENDPULSE_LISTID=1154579
SENDPULSE_ID=ba3789b39d2f2768353b278dd3bad9fe
SENDPULSE_SECRET=3a955b125b089d71f26c041507310c1c
```

Use `make run` to build&run project, `make lint` to check code with linter.

# TODO: Describe project architecture