FROM golang:1.19-alpine AS build

RUN apk add alpine-sdk bash

WORKDIR /app

COPY . .

COPY vendor/ ./vendor/

RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags dynamic --mod=vendor -o main ./cmd/message

FROM alpine:3.14

COPY --from=build /app/main /usr/bin/message-service

LABEL name=message-service

ENTRYPOINT ["message-service"]
