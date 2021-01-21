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
```