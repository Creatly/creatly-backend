FROM golang:1.15-alpine3.12 AS builder

RUN go version

COPY . /github.com/zhashkevych/creatly-backend/
WORKDIR /github.com/zhashkevych/creatly-backend/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/zhashkevych/creatly-backend/.bin/app .
COPY --from=0 /github.com/zhashkevych/creatly-backend/configs configs/
COPY --from=0 /github.com/zhashkevych/creatly-backend/templates templates/

CMD ["./app"]