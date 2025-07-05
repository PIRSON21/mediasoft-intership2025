FROM golang:1.24.4-alpine AS builder

WORKDIR /var/app

RUN apk update --no-cache

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir build && go build -o build/app -ldflags="-s -w" -trimpath ./cmd/intership


FROM alpine:3.22

WORKDIR /var/app

RUN apk update && apk get curl

COPY --from=builder ./build/app .

ENTRYPOINT ["./app"]


