FROM golang:1.19-alpine AS build

RUN apk add alpine-sdk bash

WORKDIR /app

COPY . .

COPY vendor/ ./vendor/

RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags dynamic --mod=vendor -o main ./cmd/user

FROM alpine:3.14

COPY --from=build /app/main /usr/bin/user-service

LABEL name=user-service

ENTRYPOINT ["user-service"]
