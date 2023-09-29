FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram-bot-url ./cmd/main

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/telegram-bot-url .

EXPOSE 8080

ENV DOCKER_ENV=1

CMD ["./telegram-bot-url"]