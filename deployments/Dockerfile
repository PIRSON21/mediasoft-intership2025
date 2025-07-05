FROM golang:1.24.4-alpine AS builder
WORKDIR /var/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN mkdir -p build && go build -o build/app -ldflags="-s -w" -trimpath ./cmd/intership


FROM alpine:3.22
WORKDIR /var/app

RUN apk update && apk add curl

COPY --from=builder /var/app/build/app .

ENTRYPOINT ["./app"]


